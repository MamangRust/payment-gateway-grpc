package service_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	withdraw_cache "MamangRust/paymentgatewaygrpc/internal/cache/withdraw"
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

type WithdrawServiceTestSuite struct {
	suite.Suite
	ts              *tests.TestSuite
	dbPool          *pgxpool.Pool
	redisClient     *redis.Client
	withdrawService service.WithdrawService
	userRepo        repository.UserRepository
	cardRepo        repository.CardRepository
	saldoRepo       repository.SaldoRepository
}

func (s *WithdrawServiceTestSuite) SetupSuite() {
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
	cacheWithdraw := withdraw_cache.NewWithdrawMencache(cacheStore)

	s.withdrawService = service.NewWithdrawService(service.WithdrawServiceDeps{
		WithdrawRepo:  repos.Withdraw,
		SaldoRepo:     repos.Saldo,
		UserRepo:      repos.User,
		Logger:        log,
		Observability: obs,
		Cache:         cacheWithdraw,
	})
}

func (s *WithdrawServiceTestSuite) TearDownSuite() {
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *WithdrawServiceTestSuite) TestWithdrawLifecycle() {
	ctx := context.Background()

	// Seed User and Card
	user, _ := s.userRepo.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Withdraw",
		LastName:  "User",
		Email:     "withdraw.service@example.com",
		Password:  "password123",
	})
	card, _ := s.cardRepo.CreateCard(ctx, &requests.CreateCardRequest{
		UserID:       int(user.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	s.saldoRepo.CreateSaldo(ctx, &requests.CreateSaldoRequest{
		CardNumber:   card.CardNumber,
		TotalBalance: 500000,
	})

	// 1. Create Withdraw
	req := &requests.CreateWithdrawRequest{
		CardNumber:     card.CardNumber,
		WithdrawAmount: 100000,
		WithdrawTime:   time.Now(),
	}
	withdraw, err := s.withdrawService.Create(ctx, req)
	s.NoError(err)
	s.NotNil(withdraw)
	s.Equal(int32(req.WithdrawAmount), withdraw.WithdrawAmount)

	// 2. Find By ID
	found, err := s.withdrawService.FindById(ctx, int(withdraw.WithdrawID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(withdraw.WithdrawID, found.WithdrawID)
}

func TestWithdrawServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(WithdrawServiceTestSuite))
}
