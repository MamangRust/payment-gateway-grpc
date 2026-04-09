package gapi_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	merchant_cache "MamangRust/paymentgatewaygrpc/internal/cache/merchant"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/handler/gapi"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/internal/service"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
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

type MerchantGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.MerchantServiceClient
	conn        *grpc.ClientConn
	repos       *repository.Repositories
	userID      int32
}

func (s *MerchantGapiTestSuite) SetupSuite() {
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
	cMerchant := merchant_cache.NewMerchantMencache(cacheStore)

	merchantService := service.NewMerchantService(service.MerchantServiceDeps{
		MerchantRepo:  s.repos.Merchant,
		UserRepo:      s.repos.User,
		Logger:        log,
		Observability: obs,
		Cache:         cMerchant,
	})

	// Seed User
	user, err := s.repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Merchant",
		LastName:  "Gapi",
		Email:     "merchant.gapi@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = user.UserID

	// Start gRPC Server
	merchantHandler := gapi.NewMerchantHandleGrpc(merchantService)
	server := grpc.NewServer()
	pb.RegisterMerchantServiceServer(server, merchantHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// Create Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewMerchantServiceClient(conn)
}

func (s *MerchantGapiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *MerchantGapiTestSuite) TestMerchantLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateMerchantRequest{
		UserId: s.userID,
		Name:   "Gapi Store",
	}
	res, err := s.client.CreateMerchant(ctx, createReq)
	s.NoError(err)
	s.Equal("Gapi Store", res.Data.Name)

	merchantID := res.Data.Id

	// 2. Find By ID
	found, err := s.client.FindByIdMerchant(ctx, &pb.FindByIdMerchantRequest{MerchantId: merchantID})
	s.NoError(err)
	s.Equal(merchantID, found.Data.Id)

	// 3. Update
	updateReq := &pb.UpdateMerchantRequest{
		MerchantId: merchantID,
		UserId:     s.userID,
		Name:       "Gapi Store Updated",
		Status:     "active",
	}
	updated, err := s.client.UpdateMerchant(ctx, updateReq)
	s.NoError(err)
	s.Equal("Gapi Store Updated", updated.Data.Name)

	// 4. Delete (Trash)
	_, err = s.client.TrashedMerchant(ctx, &pb.FindByIdMerchantRequest{MerchantId: merchantID})
	s.NoError(err)

	// 5. Restore
	_, err = s.client.RestoreMerchant(ctx, &pb.FindByIdMerchantRequest{MerchantId: merchantID})
	s.NoError(err)

	// 6. Permanent Delete
	_, err = s.client.DeleteMerchantPermanent(ctx, &pb.FindByIdMerchantRequest{MerchantId: merchantID})
	s.NoError(err)
}

func TestMerchantGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantGapiTestSuite))
}
