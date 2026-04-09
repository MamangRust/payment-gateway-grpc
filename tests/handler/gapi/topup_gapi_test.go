package gapi_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	topup_cache "MamangRust/paymentgatewaygrpc/internal/cache/topup"
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
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TopupGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.TopupServiceClient
	conn        *grpc.ClientConn
	repos       *repository.Repositories

	cardNumber string
}

func (s *TopupGapiTestSuite) SetupSuite() {
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
	cTopup := topup_cache.NewTopupMencache(cacheStore)

	topupService := service.NewTopupService(service.TopupServiceDeps{
		TopupRepo:     s.repos.Topup,
		CardRepo:      s.repos.Card,
		SaldoRepo:     s.repos.Saldo,
		Logger:        log,
		Observability: obs,
		Cache:         cTopup,
	})

	// Seed User, Card, Saldo
	user, _ := s.repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Topup", LastName: "Gapi", Email: "topup.gapi@test.com", Password: "password123",
	})
	card, _ := s.repos.Card.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID: int(user.UserID), CardType: "debit", ExpireDate: time.Now().AddDate(1, 0, 0), CVV: "444", CardProvider: "visa",
	})
	s.cardNumber = card.CardNumber
	s.repos.Saldo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber: s.cardNumber, TotalBalance: 0,
	})

	// Start gRPC Server
	topupHandler := gapi.NewTopupHandleGrpc(topupService)
	server := grpc.NewServer()
	pb.RegisterTopupServiceServer(server, topupHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)
	go func() { _ = server.Serve(lis) }()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewTopupServiceClient(conn)
}

func (s *TopupGapiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *TopupGapiTestSuite) TestTopupLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateTopupRequest{
		CardNumber:  s.cardNumber,
		TopupAmount: 100000,
		TopupMethod: "bri",
	}
	res, err := s.client.CreateTopup(ctx, createReq)
	s.NoError(err)
	s.Equal(int32(100000), res.Data.TopupAmount)

	topupID := res.Data.Id

	// Verify balance
	saldo, _ := s.repos.Saldo.FindByCardNumber(ctx, s.cardNumber)
	s.Equal(int32(100000), saldo.TotalBalance)

	// 2. Find By ID
	found, err := s.client.FindByIdTopup(ctx, &pb.FindByIdTopupRequest{TopupId: topupID})
	s.NoError(err)
	s.Equal(topupID, found.Data.Id)

	// 3. Update
	updateReq := &pb.UpdateTopupRequest{
		TopupId:     topupID,
		CardNumber:  s.cardNumber,
		TopupAmount: 150000,
		TopupMethod: "bri",
	}
	updated, err := s.client.UpdateTopup(ctx, updateReq)
	s.NoError(err)
	s.Equal(int32(150000), updated.Data.TopupAmount)

	// Verify adjusted balance (100k + 50k = 150k)
	saldo, _ = s.repos.Saldo.FindByCardNumber(ctx, s.cardNumber)
	s.Equal(int32(150000), saldo.TotalBalance)

	// 4. Delete (Trash)
	_, err = s.client.TrashedTopup(ctx, &pb.FindByIdTopupRequest{TopupId: topupID})
	s.NoError(err)

	// 5. Restore
	_, err = s.client.RestoreTopup(ctx, &pb.FindByIdTopupRequest{TopupId: topupID})
	s.NoError(err)

	// 6. Permanent Delete
	_, err = s.client.DeleteTopupPermanent(ctx, &pb.FindByIdTopupRequest{TopupId: topupID})
	s.NoError(err)
}

func TestTopupGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TopupGapiTestSuite))
}
