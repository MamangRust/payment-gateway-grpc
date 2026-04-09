package gapi_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	withdraw_cache "MamangRust/paymentgatewaygrpc/internal/cache/withdraw"
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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WithdrawGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.WithdrawServiceClient
	conn        *grpc.ClientConn
	repos       *repository.Repositories

	cardNumber string
}

func (s *WithdrawGapiTestSuite) SetupSuite() {
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
	cWithdraw := withdraw_cache.NewWithdrawMencache(cacheStore)

	withdrawService := service.NewWithdrawService(service.WithdrawServiceDeps{
		WithdrawRepo:  s.repos.Withdraw,
		UserRepo:      s.repos.User,
		SaldoRepo:     s.repos.Saldo,
		Logger:        log,
		Observability: obs,
		Cache:         cWithdraw,
	})

	// Seed User, Card, Saldo
	user, _ := s.repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Withdraw", LastName: "Gapi", Email: "withdraw.gapi@test.com", Password: "password123",
	})
	card, _ := s.repos.Card.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID: int(user.UserID), CardType: "debit", ExpireDate: time.Now().AddDate(1, 0, 0), CVV: "999", CardProvider: "visa",
	})
	s.cardNumber = card.CardNumber
	s.repos.Saldo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber: s.cardNumber, TotalBalance: 1000000,
	})

	// Start gRPC Server
	withdrawHandler := gapi.NewWithdrawHandleGrpc(withdrawService)
	server := grpc.NewServer()
	pb.RegisterWithdrawServiceServer(server, withdrawHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)
	go func() { _ = server.Serve(lis) }()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewWithdrawServiceClient(conn)
}

func (s *WithdrawGapiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *WithdrawGapiTestSuite) TestWithdrawLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateWithdrawRequest{
		CardNumber:     s.cardNumber,
		WithdrawAmount: 100000,
		WithdrawTime:   timestamppb.New(time.Now()),
	}
	res, err := s.client.CreateWithdraw(ctx, createReq)
	s.NoError(err)
	s.Equal(int32(100000), res.Data.WithdrawAmount)

	withdrawID := res.Data.WithdrawId

	// Verify balance
	saldo, _ := s.repos.Saldo.FindByCardNumber(ctx, s.cardNumber)
	s.Equal(int32(900000), saldo.TotalBalance)

	// 2. Find By ID
	found, err := s.client.FindByIdWithdraw(ctx, &pb.FindByIdWithdrawRequest{WithdrawId: withdrawID})
	s.NoError(err)
	s.Equal(withdrawID, found.Data.WithdrawId)

	// 3. Update
	updateReq := &pb.UpdateWithdrawRequest{
		WithdrawId:     withdrawID,
		CardNumber:     s.cardNumber,
		WithdrawAmount: 150000,
		WithdrawTime:   timestamppb.New(time.Now()),
	}
	updated, err := s.client.UpdateWithdraw(ctx, updateReq)
	s.NoError(err)
	s.Equal(int32(150000), updated.Data.WithdrawAmount)

	// Verify adjusted balance (900k - 50k = 850k)
	saldo, _ = s.repos.Saldo.FindByCardNumber(ctx, s.cardNumber)
	s.Equal(int32(850000), saldo.TotalBalance)

	// 4. Delete (Trash)
	_, err = s.client.TrashedWithdraw(ctx, &pb.FindByIdWithdrawRequest{WithdrawId: withdrawID})
	s.NoError(err)

	// 5. Restore
	_, err = s.client.RestoreWithdraw(ctx, &pb.FindByIdWithdrawRequest{WithdrawId: withdrawID})
	s.NoError(err)

	// 6. Permanent Delete
	_, err = s.client.DeleteWithdrawPermanent(ctx, &pb.FindByIdWithdrawRequest{WithdrawId: withdrawID})
	s.NoError(err)
}

func TestWithdrawGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(WithdrawGapiTestSuite))
}
