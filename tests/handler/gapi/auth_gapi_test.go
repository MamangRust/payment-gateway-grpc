package gapi_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	auth_cache "MamangRust/paymentgatewaygrpc/internal/cache/auth"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/handler/gapi"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/auth"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/hash"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
	"MamangRust/paymentgatewaygrpc/tests"
	"context"
	"net"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.AuthServiceClient
	conn        *grpc.ClientConn
}

func (s *AuthGapiTestSuite) SetupSuite() {
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
	tokenManager, _ := auth.NewManager("test-secret-key-that-is-long-enough-32")
	hasher := hash.NewHashingPassword()
	obs, _ := observability.NewObservability("test", log)
	cacheStore := cache.NewCacheStore(s.redisClient, log, nil)
	cacheAuth := auth_cache.NewMencache(cacheStore)

	authService := service.NewAuthService(service.AuthServiceDeps{
		AuthRepo:         repos.User,
		RefreshTokenRepo: repos.RefreshToken,
		RoleRepo:         repos.Role,
		UserRoleRepo:     repos.UserRole,
		Hash:             hasher,
		Token:            tokenManager,
		Logger:           log,
		Tracer:           obs,
		CacheIdentity:    cacheAuth.IdentityCache,
		CacheLogin:       cacheAuth.LoginCache,
	})

	// Seed ROLE_ADMIN
	_, _ = repos.Role.CreateRole(context.Background(), &requests.CreateRoleRequest{Name: "ROLE_ADMIN"})

	// Start gRPC Server
	authHandler := gapi.NewAuthHandleGrpc(authService)
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, authHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		if err := server.Serve(lis); err != nil {
			return
		}
	}()

	// Create Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewAuthServiceClient(conn)
}

func (s *AuthGapiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *AuthGapiTestSuite) TestRegisterUser() {
	req := &pb.RegisterRequest{
		Firstname:       "Gapi",
		Lastname:        "Test",
		Email:           "gapi@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	res, err := s.client.RegisterUser(context.Background(), req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.Equal(req.Email, res.Data.Email)
}

func TestAuthGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthGapiTestSuite))
}
