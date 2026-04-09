package api_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	api_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/transaction"
	merchant_cache "MamangRust/paymentgatewaygrpc/internal/cache/merchant"
	transaction_cache "MamangRust/paymentgatewaygrpc/internal/cache/transaction"
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

type TransactionHandlerTestSuite struct {
	suite.Suite
	ts             *tests.TestSuite
	dbPool         *pgxpool.Pool
	redisClient    *redis.Client
	grpcServer     *grpc.Server
	client         pb.TransactionServiceClient
	merchantClient pb.MerchantServiceClient
	conn           *grpc.ClientConn
	router         *echo.Echo
	repos          *repository.Repositories

	customerCardNumber string
	merchantApiKey     string
	merchantID         int
	merchantCardNumber string
}

func (s *TransactionHandlerTestSuite) SetupSuite() {
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

	cTransaction := transaction_cache.NewTransactionMencache(cacheStore)
	cMerchant := merchant_cache.NewMerchantMencache(cacheStore)

	transactionService := service.NewTransactionService(service.TransactionServiceDeps{
		TransactionRepo: s.repos.Transaction,
		MerchantRepo:    s.repos.Merchant,
		CardRepo:        s.repos.Card,
		SaldoRepo:       s.repos.Saldo,
		Logger:          log,
		Observability:   obs,
		Cache:           cTransaction,
	})

	merchantService := service.NewMerchantService(service.MerchantServiceDeps{
		UserRepo:      s.repos.User,
		MerchantRepo:  s.repos.Merchant,
		Cache:         cMerchant,
		Logger:        log,
		Observability: obs,
	})

	// Seed Customer
	customer, err := s.repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Transaction",
		LastName:  "Customer",
		Email:     "customer@transaction.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	card, err := s.repos.Card.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       int(customer.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.customerCardNumber = card.CardNumber

	_, err = s.repos.Saldo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   s.customerCardNumber,
		TotalBalance: 1000000,
	})
	s.Require().NoError(err)

	// Seed Merchant
	mOwner, err := s.repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Merchant",
		LastName:  "Owner",
		Email:     "merchant.owner@transaction.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	merchant, err := s.repos.Merchant.CreateMerchant(context.Background(), &requests.CreateMerchantRequest{
		UserID: int(mOwner.UserID),
		Name:   "Transaction Merchant",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	_, err = s.repos.Merchant.UpdateMerchantStatus(context.Background(), &requests.UpdateMerchantStatus{
		MerchantID: s.merchantID,
		Status:     "active",
	})
	s.Require().NoError(err)

	mFull, _ := s.repos.Merchant.FindById(context.Background(), s.merchantID)
	s.merchantApiKey = mFull.ApiKey

	mCard, err := s.repos.Card.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       int(mOwner.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		CVV:          "321",
		CardProvider: "mastercard",
	})
	s.Require().NoError(err)
	s.merchantCardNumber = mCard.CardNumber

	_, err = s.repos.Saldo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   s.merchantCardNumber,
		TotalBalance: 0,
	})
	s.Require().NoError(err)

	// Start gRPC Server
	transactionGapiHandler := gapi.NewTransactionHandleGrpc(transactionService)
	merchantGapiHandler := gapi.NewMerchantHandleGrpc(merchantService)

	server := grpc.NewServer()
	pb.RegisterTransactionServiceServer(server, transactionGapiHandler)
	pb.RegisterMerchantServiceServer(server, merchantGapiHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// Create gRPC Clients
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewTransactionServiceClient(conn)
	s.merchantClient = pb.NewMerchantServiceClient(conn)

	// Setup Echo
	s.router = echo.New()
	transactionMapper := apimapper.NewTransactionResponseMapper()
	apiErrorHandler := app_errors.NewApiHandler(obs, log)
	apiCacheTransaction := api_cache.NewTransactionMencache(cacheStore)

	api.NewHandlerTransaction(s.client, s.merchantClient, s.router, log, transactionMapper, apiErrorHandler, apiCacheTransaction)
}

func (s *TransactionHandlerTestSuite) TearDownSuite() {
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

func (s *TransactionHandlerTestSuite) TestTransactionLifecycle() {
	// 1. Create Transaction
	createReq := map[string]interface{}{
		"card_number":      s.customerCardNumber,
		"amount":           50000,
		"payment_method":   "visa",
		"merchant_id":      s.merchantID,
		"transaction_time": time.Now().Format(time.RFC3339),
	}
	body, _ := json.Marshal(createReq)

	request := httptest.NewRequest(http.MethodPost, "/api/transactions/create", bytes.NewBuffer(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set("X-API-Key", s.merchantApiKey)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, request)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())

	var createRes map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &createRes)
	s.Require().NoError(err)
	transactionData, ok := createRes["data"].(map[string]interface{})
	s.Require().True(ok, "response data should be a map")
	transactionID := int(transactionData["id"].(float64))

	// Verify balance update
	customerSaldo, _ := s.repos.Saldo.FindByCardNumber(context.Background(), s.customerCardNumber)
	s.Equal(int32(950000), customerSaldo.TotalBalance)

	merchantSaldo, _ := s.repos.Saldo.FindByCardNumber(context.Background(), s.merchantCardNumber)
	s.Equal(int32(50000), merchantSaldo.TotalBalance)

	// 2. Find By ID
	request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/transactions/%d", transactionID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Find All
	request = httptest.NewRequest(http.MethodGet, "/api/transactions?page=1&page_size=10", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code, rec.Body.String())

	// 4. Update
	updateReq := map[string]interface{}{
		"card_number":      s.customerCardNumber, // Fixed: Added missing card_number
		"amount":           60000,                // Increase amount
		"payment_method":   "visa",
		"merchant_id":      s.merchantID,
		"transaction_time": time.Now().Format(time.RFC3339),
	}
	updateBody, _ := json.Marshal(updateReq)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transactions/update/%d", transactionID), bytes.NewBuffer(updateBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set("X-API-Key", s.merchantApiKey)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())

	// Verify adjusted balance (950000 - 10000 = 940000)
	customerSaldo, _ = s.repos.Saldo.FindByCardNumber(context.Background(), s.customerCardNumber)
	s.Equal(int32(940000), customerSaldo.TotalBalance)

	// 5. Trashed (Soft Delete)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transactions/trashed/%d", transactionID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Restore
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transactions/%d/restore", transactionID), nil) // Fixed path if needed
	// Wait, I used "/api/transactions/restore/:id" in NewHandlerTransaction
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transactions/restore/%d", transactionID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Permanent Delete
	request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/transactions/permanent/%d", transactionID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func TestTransactionHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionHandlerTestSuite))
}
