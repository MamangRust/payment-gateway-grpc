package api_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	api_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/auth"
	auth_cache "MamangRust/paymentgatewaygrpc/internal/cache/auth"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/handler/api"
	"MamangRust/paymentgatewaygrpc/internal/handler/gapi"
	apimapper "MamangRust/paymentgatewaygrpc/internal/mapper"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/internal/service"
	auth_manager "MamangRust/paymentgatewaygrpc/pkg/auth"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	app_errors "MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/hash"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
	"MamangRust/paymentgatewaygrpc/tests"
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.AuthServiceClient
	conn        *grpc.ClientConn
	router      *echo.Echo
}

func (s *AuthHandlerTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	s.redisClient = redis.NewClient(&redis.Options{
		Addr: s.ts.RedisURL,
	})

	queries := db.New(pool)
	repos := repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	tokenM, _ := auth_manager.NewManager("test-secret-key-that-is-long-enough-32")
	hasher := hash.NewHashingPassword()
	obs, _ := observability.NewObservability("test", log)
	cacheStore := cache.NewCacheStore(s.redisClient, log, nil)
	cAuth := auth_cache.NewMencache(cacheStore)

	authService := service.NewAuthService(service.AuthServiceDeps{
		AuthRepo:         repos.User,
		RefreshTokenRepo: repos.RefreshToken,
		RoleRepo:         repos.Role,
		UserRoleRepo:     repos.UserRole,
		Hash:             hasher,
		Token:            tokenM,
		Logger:           log,
		Tracer:           obs,
		CacheIdentity:    cAuth.IdentityCache,
		CacheLogin:       cAuth.LoginCache,
	})

	// Seed ROLE_ADMIN
	_, _ = repos.Role.CreateRole(context.Background(), &requests.CreateRoleRequest{Name: "ROLE_ADMIN"})

	// Start gRPC Server
	authGapiHandler := gapi.NewAuthHandleGrpc(authService)
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, authGapiHandler)
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
	s.client = pb.NewAuthServiceClient(conn)

	// Setup Echo
	s.router = echo.New()
	authMapper := apimapper.NewAuthResponseMapper()
	apiErrorHandler := app_errors.NewApiHandler(obs, log)
	apiCacheAuth := api_cache.NewMencache(cacheStore)

	api.NewHandlerAuth(s.router, s.client, log, authMapper, apiErrorHandler, apiCacheAuth)
}

func (s *AuthHandlerTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *AuthHandlerTestSuite) TestRegister() {
	req := requests.CreateUserRequest{
		FirstName:       "Handler",
		LastName:        "Test",
		Email:           "handler@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	body, _ := json.Marshal(req)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, request)

	s.Equal(http.StatusCreated, rec.Code)
}

func TestAuthHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthHandlerTestSuite))
}
