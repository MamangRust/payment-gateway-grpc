package service_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	transfer_cache "MamangRust/paymentgatewaygrpc/internal/cache/transfer"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/internal/service"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
	"MamangRust/paymentgatewaygrpc/tests"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type TransferServiceTestSuite struct {
	suite.Suite
	ts              *tests.TestSuite
	dbPool          *pgxpool.Pool
	redisClient     *redis.Client
	transferService service.TransferService
	userRepo        repository.UserRepository
	cardRepo        repository.CardRepository
	saldoRepo       repository.SaldoRepository
}

func (s *TransferServiceTestSuite) SetupSuite() {
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
	s.userRepo = repos.User
	s.cardRepo = repos.Card
	s.saldoRepo = repos.Saldo

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	obs, _ := observability.NewObservability("test", log)
	cacheStore := cache.NewCacheStore(s.redisClient, log, nil)
	cacheTransfer := transfer_cache.NewTransferMencache(cacheStore)

	s.transferService = service.NewTransferService(service.TransferServiceDeps{
		TransferRepo:  repos.Transfer,
		UserRepo:      repos.User,
		CardRepo:      repos.Card,
		SaldoRepo:     repos.Saldo,
		Logger:        log,
		Observability: obs,
		Cache:         cacheTransfer,
	})
}

func (s *TransferServiceTestSuite) TearDownSuite() {
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *TransferServiceTestSuite) TestTransferLifecycle() {
	ctx := context.Background()

	// Seed Users and Cards
	user1, _ := s.userRepo.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Sender",
		LastName:  "User",
		Email:     "sender@example.com",
		Password:  "password123",
	})
	card1, _ := s.cardRepo.CreateCard(ctx, &requests.CreateCardRequest{
		UserID:       int(user1.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	s.saldoRepo.CreateSaldo(ctx, &requests.CreateSaldoRequest{
		CardNumber:   card1.CardNumber,
		TotalBalance: 1000000,
	})

	user2, _ := s.userRepo.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Receiver",
		LastName:  "User",
		Email:     "receiver@example.com",
		Password:  "password123",
	})
	card2, _ := s.cardRepo.CreateCard(ctx, &requests.CreateCardRequest{
		UserID:       int(user2.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	s.saldoRepo.CreateSaldo(ctx, &requests.CreateSaldoRequest{
		CardNumber:   card2.CardNumber,
		TotalBalance: 100000,
	})

	// 1. Create Transfer
	req := &requests.CreateTransferRequest{
		TransferFrom:   card1.CardNumber,
		TransferTo:     card2.CardNumber,
		TransferAmount: 50000,
	}
	transfer, err := s.transferService.CreateTransaction(ctx, req)
	s.NoError(err)
	s.NotNil(transfer)
	s.Equal(int32(req.TransferAmount), transfer.TransferAmount)

	// 2. Find By ID
	found, err := s.transferService.FindById(ctx, int(transfer.TransferID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(transfer.TransferID, found.TransferID)
}

func TestTransferServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransferServiceTestSuite))
}
