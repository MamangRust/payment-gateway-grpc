package api_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	api_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/card"
	card_cache "MamangRust/paymentgatewaygrpc/internal/cache/card"
	user_cache "MamangRust/paymentgatewaygrpc/internal/cache/user"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/handler/api"
	"MamangRust/paymentgatewaygrpc/internal/handler/gapi"
	apimapper "MamangRust/paymentgatewaygrpc/internal/mapper"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/internal/service"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	app_errors "MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/hash"
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

type CardHandlerTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.CardServiceClient
	conn        *grpc.ClientConn
	router      *echo.Echo
	userID      int
}

func (s *CardHandlerTestSuite) SetupSuite() {
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
	hasher := hash.NewHashingPassword()
	_ = hasher
	obs, _ := observability.NewObservability("test", log)
	cacheStore := cache.NewCacheStore(s.redisClient, log, nil)

	cCard := card_cache.NewCardMencache(cacheStore)
	_ = user_cache.NewUserMencache(cacheStore)

	cardService := service.NewCardService(service.CardServiceDeps{
		CardRepo:      repos.Card,
		UserRepo:      repos.User,
		Logger:        log,
		Observability: obs,
		Cache:         cCard,
	})

	// Seed User
	user, err := repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Card",
		LastName:  "Owner",
		Email:     "card.owner@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// Start gRPC Server
	cardGapiHandler := gapi.NewCardHandleGrpc(cardService)
	server := grpc.NewServer()
	pb.RegisterCardServiceServer(server, cardGapiHandler)
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
	s.client = pb.NewCardServiceClient(conn)

	// Setup Echo
	s.router = echo.New()
	cardMapper := apimapper.NewCardResponseMapper()
	apiErrorHandler := app_errors.NewApiHandler(obs, log)
	apiCacheCard := api_cache.NewCardMencache(cacheStore)

	api.NewHandlerCard(s.client, s.router, log, apiErrorHandler, cardMapper, apiCacheCard)
}

func (s *CardHandlerTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *CardHandlerTestSuite) TestCardLifecycle() {
	// 1. Create Card
	req := requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	}
	body, _ := json.Marshal(req)

	request := httptest.NewRequest(http.MethodPost, "/api/card/create", bytes.NewBuffer(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	cardData := createRes["data"].(map[string]interface{})
	cardID := int(cardData["id"].(float64))

	// 2. Find By ID
	request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/card/%d", cardID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Update
	updateReq := requests.UpdateCardRequest{
		CardID:       cardID,
		UserID:       s.userID,
		CardType:     "credit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "123",
		CardProvider: "mastercard",
	}
	updateBody, _ := json.Marshal(updateReq)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/card/update/%d", cardID), bytes.NewBuffer(updateBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Delete Permanent
	request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/card/permanent/%d", cardID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func TestCardHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CardHandlerTestSuite))
}
