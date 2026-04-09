package gapi_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	transfer_cache "MamangRust/paymentgatewaygrpc/internal/cache/transfer"
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

type TransferGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.TransferServiceClient
	conn        *grpc.ClientConn
	repos       *repository.Repositories

	senderCard   string
	receiverCard string
}

func (s *TransferGapiTestSuite) SetupSuite() {
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
		UserRepo:      s.repos.User,
		CardRepo:      s.repos.Card,
		SaldoRepo:     s.repos.Saldo,
		Logger:        log,
		Observability: obs,
		Cache:         cTransfer,
	})

	// Seed Sender
	sUser, _ := s.repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Sender", LastName: "Gapi", Email: "sender.gapi@test.com", Password: "password123",
	})
	sCard, _ := s.repos.Card.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID: int(sUser.UserID), CardType: "debit", ExpireDate: time.Now().AddDate(1, 0, 0), CVV: "111", CardProvider: "visa",
	})
	s.senderCard = sCard.CardNumber
	s.repos.Saldo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber: s.senderCard, TotalBalance: 1000000,
	})

	// Seed Receiver
	rUser, _ := s.repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Receiver", LastName: "Gapi", Email: "receiver.gapi@test.com", Password: "password123",
	})
	rCard, _ := s.repos.Card.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID: int(rUser.UserID), CardType: "debit", ExpireDate: time.Now().AddDate(1, 0, 0), CVV: "222", CardProvider: "mastercard",
	})
	s.receiverCard = rCard.CardNumber
	s.repos.Saldo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber: s.receiverCard, TotalBalance: 0,
	})

	// Start gRPC Server
	transferHandler := gapi.NewTransferHandleGrpc(transferService)
	server := grpc.NewServer()
	pb.RegisterTransferServiceServer(server, transferHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)
	go func() { _ = server.Serve(lis) }()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewTransferServiceClient(conn)
}

func (s *TransferGapiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *TransferGapiTestSuite) TestTransferLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateTransferRequest{
		TransferFrom:   s.senderCard,
		TransferTo:     s.receiverCard,
		TransferAmount: 100000,
	}
	res, err := s.client.CreateTransfer(ctx, createReq)
	s.NoError(err)
	s.Equal(int32(100000), res.Data.TransferAmount)

	transferID := res.Data.Id

	// Verify balances
	sSaldo, _ := s.repos.Saldo.FindByCardNumber(ctx, s.senderCard)
	s.Equal(int32(900000), sSaldo.TotalBalance)
	rSaldo, _ := s.repos.Saldo.FindByCardNumber(ctx, s.receiverCard)
	s.Equal(int32(100000), rSaldo.TotalBalance)

	// 2. Find By ID
	found, err := s.client.FindByIdTransfer(ctx, &pb.FindByIdTransferRequest{TransferId: transferID})
	s.NoError(err)
	s.Equal(transferID, found.Data.Id)

	// 3. Update
	updateReq := &pb.UpdateTransferRequest{
		TransferId:     transferID,
		TransferFrom:   s.senderCard,
		TransferTo:     s.receiverCard,
		TransferAmount: 150000,
	}
	updated, err := s.client.UpdateTransfer(ctx, updateReq)
	s.NoError(err)
	s.Equal(int32(150000), updated.Data.TransferAmount)

	// Verify adjusted balances (Sender 900k - 50k = 850k, Receiver 100k + 50k = 150k)
	sSaldo, _ = s.repos.Saldo.FindByCardNumber(ctx, s.senderCard)
	s.Equal(int32(850000), sSaldo.TotalBalance)
	rSaldo, _ = s.repos.Saldo.FindByCardNumber(ctx, s.receiverCard)
	s.Equal(int32(150000), rSaldo.TotalBalance)

	// 4. Delete (Trash)
	_, err = s.client.TrashedTransfer(ctx, &pb.FindByIdTransferRequest{TransferId: transferID})
	s.NoError(err)

	// 5. Restore
	_, err = s.client.RestoreTransfer(ctx, &pb.FindByIdTransferRequest{TransferId: transferID})
	s.NoError(err)

	// 6. Permanent Delete
	_, err = s.client.DeleteTransferPermanent(ctx, &pb.FindByIdTransferRequest{TransferId: transferID})
	s.NoError(err)
}

func TestTransferGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransferGapiTestSuite))
}
