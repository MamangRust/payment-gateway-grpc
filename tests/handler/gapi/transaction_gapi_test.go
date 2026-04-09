package gapi_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	transaction_cache "MamangRust/paymentgatewaygrpc/internal/cache/transaction"
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

type TransactionGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.TransactionServiceClient
	conn        *grpc.ClientConn
	repos       *repository.Repositories

	cardNumber string
	merchantID int32
	apiKey     string
}

func (s *TransactionGapiTestSuite) SetupSuite() {
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

	transactionService := service.NewTransactionService(service.TransactionServiceDeps{
		TransactionRepo: s.repos.Transaction,
		CardRepo:        s.repos.Card,
		MerchantRepo:    s.repos.Merchant,
		SaldoRepo:       s.repos.Saldo,
		Logger:          log,
		Observability:   obs,
		Cache:           cTransaction,
	})

	// Seed User, Card, Merchant, Saldo
	user, _ := s.repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Trans", LastName: "Gapi", Email: "trans.gapi@test.com", Password: "password123",
	})
	card, _ := s.repos.Card.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID: int(user.UserID), CardType: "debit", ExpireDate: time.Now().AddDate(1, 0, 0), CVV: "123", CardProvider: "visa",
	})
	s.cardNumber = card.CardNumber
	merchant, _ := s.repos.Merchant.CreateMerchant(context.Background(), &requests.CreateMerchantRequest{
		UserID: int(user.UserID), Name: "Gapi Merchant",
	})
	s.merchantID = merchant.MerchantID
	s.apiKey = merchant.ApiKey
	s.repos.Saldo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber: s.cardNumber, TotalBalance: 1000000,
	})

	// Start gRPC Server
	transactionHandler := gapi.NewTransactionHandleGrpc(transactionService)
	server := grpc.NewServer()
	pb.RegisterTransactionServiceServer(server, transactionHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)
	go func() { _ = server.Serve(lis) }()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewTransactionServiceClient(conn)
}

func (s *TransactionGapiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *TransactionGapiTestSuite) TestTransactionLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateTransactionRequest{
		ApiKey:          s.apiKey,
		CardNumber:      s.cardNumber,
		Amount:          100000,
		PaymentMethod:   "bri",
		MerchantId:      s.merchantID,
		TransactionTime: timestamppb.New(time.Now()),
	}
	res, err := s.client.CreateTransaction(ctx, createReq)
	s.NoError(err)
	s.Equal(int32(100000), res.Data.Amount)

	transactionID := res.Data.Id

	// 2. Find By ID
	found, err := s.client.FindByIdTransaction(ctx, &pb.FindByIdTransactionRequest{TransactionId: transactionID})
	s.NoError(err)
	s.Equal(transactionID, found.Data.Id)

	// 3. Update
	updateReq := &pb.UpdateTransactionRequest{
		ApiKey:          s.apiKey,
		TransactionId:   transactionID,
		CardNumber:      s.cardNumber,
		Amount:          150000,
		PaymentMethod:   "bri",
		MerchantId:      s.merchantID,
		TransactionTime: timestamppb.New(time.Now()),
	}
	updated, err := s.client.UpdateTransaction(ctx, updateReq)
	s.NoError(err)
	s.Equal(int32(150000), updated.Data.Amount)

	// 4. Delete (Trash)
	_, err = s.client.TrashedTransaction(ctx, &pb.FindByIdTransactionRequest{TransactionId: transactionID})
	s.NoError(err)

	// 5. Restore
	_, err = s.client.RestoreTransaction(ctx, &pb.FindByIdTransactionRequest{TransactionId: transactionID})
	s.NoError(err)

	// 6. Permanent Delete
	_, err = s.client.DeleteTransactionPermanent(ctx, &pb.FindByIdTransactionRequest{TransactionId: transactionID})
	s.NoError(err)
}

func TestTransactionGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionGapiTestSuite))
}
