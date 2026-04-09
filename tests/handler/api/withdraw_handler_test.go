package api_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	api_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/withdraw"
	withdraw_cache "MamangRust/paymentgatewaygrpc/internal/cache/withdraw"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/handler/api"
	"MamangRust/paymentgatewaygrpc/internal/handler/gapi"
	apimapper "MamangRust/paymentgatewaygrpc/internal/mapper"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/internal/service"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	app_errors "MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
	"MamangRust/paymentgatewaygrpc/tests"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WithdrawHandlerTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.WithdrawServiceClient
	conn        *grpc.ClientConn
	router      *echo.Echo
	repos       *repository.Repositories

	customerCardNumber string
}

func (s *WithdrawHandlerTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repos = repository.NewRepositories(queries)

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	obs, _ := observability.NewObservability("test", log)
	cacheStore := cache.NewCacheStore(s.redisClient, log, nil)

	cWithdraw := withdraw_cache.NewWithdrawMencache(cacheStore)

	withdrawService := service.NewWithdrawService(service.WithdrawServiceDeps{
		UserRepo:      s.repos.User,
		SaldoRepo:     s.repos.Saldo,
		WithdrawRepo:  s.repos.Withdraw,
		Logger:        log,
		Observability: obs,
		Cache:         cWithdraw,
	})

	// Seed Customer
	customer, err := s.repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Withdraw",
		LastName:  "Customer",
		Email:     "withdraw@test.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	card, err := s.repos.Card.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       int(customer.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		CVV:          "999",
		CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.customerCardNumber = card.CardNumber

	_, err = s.repos.Saldo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   s.customerCardNumber,
		TotalBalance: 1000000,
	})
	s.Require().NoError(err)

	// Start gRPC Server
	withdrawGapiHandler := gapi.NewWithdrawHandleGrpc(withdrawService)

	server := grpc.NewServer()
	pb.RegisterWithdrawServiceServer(server, withdrawGapiHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// Create gRPC Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewWithdrawServiceClient(conn)

	// Setup Echo
	s.router = echo.New()
	withdrawMapper := apimapper.NewWithdrawResponseMapper()
	apiErrorHandler := app_errors.NewApiHandler(obs, log)
	apiCacheWithdraw := api_cache.NewWithdrawMencache(cacheStore)

	api.NewHandlerWithdraw(s.client, s.router, log, withdrawMapper, apiErrorHandler, apiCacheWithdraw)
}

func (s *WithdrawHandlerTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *WithdrawHandlerTestSuite) TestWithdrawLifecycle() {
	// 1. Create Withdraw
	createReq := map[string]interface{}{
		"card_number":     s.customerCardNumber,
		"withdraw_amount": 100000,
		"withdraw_time":   time.Now().Format(time.RFC3339),
	}
	body, _ := json.Marshal(createReq)

	request := httptest.NewRequest(http.MethodPost, "/api/withdraws/create", bytes.NewBuffer(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, request)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())

	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	withdrawData := createRes["data"].(map[string]interface{})
	withdrawID := int(withdrawData["id"].(float64))

	// Verify balance
	customerSaldo, _ := s.repos.Saldo.FindByCardNumber(context.Background(), s.customerCardNumber)
	s.Equal(int32(900000), customerSaldo.TotalBalance)

	// 2. Find By ID
	request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/withdraws/%d", withdrawID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Find All
	request = httptest.NewRequest(http.MethodGet, "/api/withdraws?page=1&page_size=10", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Update
	updateReq := map[string]interface{}{
		"card_number":     s.customerCardNumber,
		"withdraw_id":     withdrawID,
		"withdraw_amount": 150000, // Increase by 50000
		"withdraw_time":   time.Now().Format(time.RFC3339),
	}
	updateBody, _ := json.Marshal(updateReq)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/withdraws/update/%d", withdrawID), bytes.NewBuffer(updateBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())

	// Verify adjusted balance (900k - 50k = 850k)
	customerSaldo, _ = s.repos.Saldo.FindByCardNumber(context.Background(), s.customerCardNumber)
	s.Equal(int32(850000), customerSaldo.TotalBalance)

	// 5. Trashed
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/withdraws/trashed/%d", withdrawID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Restore
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/withdraws/restore/%d", withdrawID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Permanent Delete
	request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/withdraws/permanent/%d", withdrawID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func TestWithdrawHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(WithdrawHandlerTestSuite))
}
