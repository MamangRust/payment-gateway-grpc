package api_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	api_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/saldo"
	saldo_cache "MamangRust/paymentgatewaygrpc/internal/cache/saldo"
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

type SaldoHandlerTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.SaldoServiceClient
	conn        *grpc.ClientConn
	router      *echo.Echo
	cardNumber  string
}

func (s *SaldoHandlerTestSuite) SetupSuite() {
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
	repos := repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	obs, _ := observability.NewObservability("test", log)
	cacheStore := cache.NewCacheStore(s.redisClient, log, nil)

	cSaldo := saldo_cache.NewSaldoMencache(cacheStore)

	saldoService := service.NewSaldoService(service.SaldoServiceDeps{
		SaldoRepo:     repos.Saldo,
		CardRepo:      repos.Card,
		Logger:        log,
		Observability: obs,
		Cache:         cSaldo,
	})

	// Seed User and Card
	user, err := repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Saldo",
		LastName:  "Owner",
		Email:     "saldo.owner@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	card, err := repos.Card.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       int(user.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.cardNumber = card.CardNumber

	// Start gRPC Server
	saldoGapiHandler := gapi.NewSaldoHandleGrpc(saldoService)
	server := grpc.NewServer()
	pb.RegisterSaldoServiceServer(server, saldoGapiHandler)
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
	s.client = pb.NewSaldoServiceClient(conn)

	// Setup Echo
	s.router = echo.New()
	saldoMapper := apimapper.NewSaldoResponseMapper()
	apiErrorHandler := app_errors.NewApiHandler(obs, log)
	apiCacheSaldo := api_cache.NewSaldoMencache(cacheStore)

	api.NewHandlerSaldo(s.client, s.router, log, saldoMapper, apiErrorHandler, apiCacheSaldo)
}

func (s *SaldoHandlerTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *SaldoHandlerTestSuite) TestSaldoLifecycle() {
	// 1. Create Saldo
	req := requests.CreateSaldoRequest{
		CardNumber:   s.cardNumber,
		TotalBalance: 100000,
	}
	body, _ := json.Marshal(req)

	request := httptest.NewRequest(http.MethodPost, "/api/saldos/create", bytes.NewBuffer(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, request)
	if rec.Code != http.StatusOK {
		s.T().Logf("Create Saldo response: %s", rec.Body.String())
	}
	s.Equal(http.StatusOK, rec.Code)

	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	saldoData := createRes["data"].(map[string]interface{})
	saldoID := int(saldoData["id"].(float64))

	// 2. Find By ID
	request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/saldos/%d", saldoID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Find By Card Number
	request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/saldos/card_number/%s", s.cardNumber), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Update
	updateReq := requests.UpdateSaldoRequest{
		CardNumber:   s.cardNumber,
		TotalBalance: 150000,
	}
	updateBody, _ := json.Marshal(updateReq)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/saldos/update/%d", saldoID), bytes.NewBuffer(updateBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Delete Permanent
	request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/saldos/permanent/%d", saldoID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	if rec.Code != http.StatusOK {
		s.T().Logf("Delete Saldo response: %s", rec.Body.String())
	}
	s.Equal(http.StatusOK, rec.Code)
}

func TestSaldoHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SaldoHandlerTestSuite))
}
