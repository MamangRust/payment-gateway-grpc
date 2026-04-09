package api_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	api_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/transfer"
	transfer_cache "MamangRust/paymentgatewaygrpc/internal/cache/transfer"
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

type TransferHandlerTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.TransferServiceClient
	conn        *grpc.ClientConn
	router      *echo.Echo
	repos       *repository.Repositories

	senderCardNumber   string
	receiverCardNumber string
}

func (s *TransferHandlerTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	queries := db.New(pool)
	s.repos = repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	obs, _ := observability.NewObservability("test", log)
	cacheStore := cache.NewCacheStore(s.redisClient, log, nil)

	cTransfer := transfer_cache.NewTransferMencache(cacheStore)

	transferService := service.NewTransferService(service.TransferServiceDeps{
		TransferRepo:  s.repos.Transfer,
		CardRepo:      s.repos.Card,
		SaldoRepo:     s.repos.Saldo,
		Logger:        log,
		Observability: obs,
		Cache:         cTransfer,
	})

	// Seed Sender
	sender, err := s.repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Sender",
		LastName:  "User",
		Email:     "sender@transfer.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	sCard, err := s.repos.Card.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       int(sender.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		CVV:          "111",
		CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.senderCardNumber = sCard.CardNumber

	_, err = s.repos.Saldo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   s.senderCardNumber,
		TotalBalance: 1000000,
	})
	s.Require().NoError(err)

	// Seed Receiver
	receiver, err := s.repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Receiver",
		LastName:  "User",
		Email:     "receiver@transfer.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	rCard, err := s.repos.Card.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       int(receiver.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		CVV:          "222",
		CardProvider: "mastercard",
	})
	s.Require().NoError(err)
	s.receiverCardNumber = rCard.CardNumber

	_, err = s.repos.Saldo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   s.receiverCardNumber,
		TotalBalance: 0,
	})
	s.Require().NoError(err)

	// Start gRPC Server
	transferGapiHandler := gapi.NewTransferHandleGrpc(transferService)

	server := grpc.NewServer()
	pb.RegisterTransferServiceServer(server, transferGapiHandler)
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
	s.client = pb.NewTransferServiceClient(conn)

	// Setup Echo
	s.router = echo.New()
	transferMapper := apimapper.NewTransferResponseMapper()
	apiErrorHandler := app_errors.NewApiHandler(obs, log)
	apiCacheTransfer := api_cache.NewTransferMencache(cacheStore)

	api.NewHandlerTransfer(s.client, s.router, log, transferMapper, apiErrorHandler, apiCacheTransfer)
}

func (s *TransferHandlerTestSuite) TearDownSuite() {
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

func (s *TransferHandlerTestSuite) TestTransferLifecycle() {
	// 1. Create Transfer
	createReq := map[string]interface{}{
		"transfer_from":   s.senderCardNumber,
		"transfer_to":     s.receiverCardNumber,
		"transfer_amount": 100000,
		"transfer_time":   time.Now().Format(time.RFC3339),
	}
	body, _ := json.Marshal(createReq)

	request := httptest.NewRequest(http.MethodPost, "/api/transfers/create", bytes.NewBuffer(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, request)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())

	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	transferData := createRes["data"].(map[string]interface{})
	transferID := int(transferData["id"].(float64))

	// Verify balances
	senderSaldo, _ := s.repos.Saldo.FindByCardNumber(context.Background(), s.senderCardNumber)
	s.Equal(int32(900000), senderSaldo.TotalBalance)

	receiverSaldo, _ := s.repos.Saldo.FindByCardNumber(context.Background(), s.receiverCardNumber)
	s.Equal(int32(100000), receiverSaldo.TotalBalance)

	// 2. Find By ID
	request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/transfers/%d", transferID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Find All
	request = httptest.NewRequest(http.MethodGet, "/api/transfers?page=1&page_size=10", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Update
	updateReq := map[string]interface{}{
		"transfer_from":   s.senderCardNumber,
		"transfer_to":     s.receiverCardNumber,
		"transfer_amount": 150000, // Increase by 50000
		"transfer_time":   time.Now().Format(time.RFC3339),
	}
	updateBody, _ := json.Marshal(updateReq)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transfers/update/%d", transferID), bytes.NewBuffer(updateBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())

	// Verify adjusted balances (Sender 900k - 50k = 850k, Receiver 100k + 50k = 150k)
	senderSaldo, _ = s.repos.Saldo.FindByCardNumber(context.Background(), s.senderCardNumber)
	s.Equal(int32(850000), senderSaldo.TotalBalance)

	receiverSaldo, _ = s.repos.Saldo.FindByCardNumber(context.Background(), s.receiverCardNumber)
	s.Equal(int32(150000), receiverSaldo.TotalBalance)

	// 5. Trashed
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transfers/trashed/%d", transferID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Restore
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transfers/restore/%d", transferID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Permanent Delete
	request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/transfers/permanent/%d", transferID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func TestTransferHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransferHandlerTestSuite))
}
