package service_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	transaction_cache "MamangRust/paymentgatewaygrpc/internal/cache/transaction"
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

type TransactionServiceTestSuite struct {
	suite.Suite
	ts                 *tests.TestSuite
	dbPool             *pgxpool.Pool
	redisClient        *redis.Client
	transactionService service.TransactionService
	userRepo           repository.UserRepository
	cardRepo           repository.CardRepository
	merchantRepo       repository.MerchantRepository
	saldoRepo          repository.SaldoRepository
}

func (s *TransactionServiceTestSuite) SetupSuite() {
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
	s.merchantRepo = repos.Merchant
	s.saldoRepo = repos.Saldo

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	cacheStore := cache.NewCacheStore(s.redisClient, log, nil)
	cacheTx := transaction_cache.NewTransactionMencache(cacheStore)

	obs, _ := observability.NewObservability("test", log)
	s.transactionService = service.NewTransactionService(service.TransactionServiceDeps{
		TransactionRepo: repos.Transaction,
		MerchantRepo:    repos.Merchant,
		CardRepo:        repos.Card,
		SaldoRepo:       repos.Saldo,
		Logger:          log,
		Observability:   obs,
		Cache:           cacheTx,
	})
}

func (s *TransactionServiceTestSuite) TearDownSuite() {
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *TransactionServiceTestSuite) TestTransactionLifecycle() {
	ctx := context.Background()

	// Seed User, Card, Merchant
	user, _ := s.userRepo.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Tx",
		LastName:  "Owner",
		Email:     "tx.service@example.com",
		Password:  "password123",
	})
	card, _ := s.cardRepo.CreateCard(ctx, &requests.CreateCardRequest{
		UserID:       int(user.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		Name:   "Service Merchant",
		UserID: int(user.UserID),
	})

	s.saldoRepo.CreateSaldo(ctx, &requests.CreateSaldoRequest{
		CardNumber:   card.CardNumber,
		TotalBalance: 1000000,
	})

	merchantCard, _ := s.cardRepo.CreateCard(ctx, &requests.CreateCardRequest{
		UserID:       int(user.UserID), // Using the same user for merchant card for simplicity in seed
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "456",
		CardProvider: "mastercard",
	})
	s.saldoRepo.CreateSaldo(ctx, &requests.CreateSaldoRequest{
		CardNumber:   merchantCard.CardNumber,
		TotalBalance: 0,
	})

	// 1. Create Transaction
	merchantID := int(merchant.MerchantID)
	req := &requests.CreateTransactionRequest{
		CardNumber:      card.CardNumber,
		Amount:          100000,
		PaymentMethod:   "visa",
		MerchantID:      &merchantID,
		TransactionTime: time.Now(),
	}
	tx, err := s.transactionService.Create(ctx, merchant.ApiKey, req)
	s.NoError(err)
	s.NotNil(tx)
	s.Equal(int32(req.Amount), tx.Amount)

	// 2. Find By ID
	found, err := s.transactionService.FindById(ctx, int(tx.TransactionID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(tx.TransactionID, found.TransactionID)
}

func TestTransactionServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionServiceTestSuite))
}
