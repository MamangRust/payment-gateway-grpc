package service_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	topup_cache "MamangRust/paymentgatewaygrpc/internal/cache/topup"
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

type TopupServiceTestSuite struct {
	suite.Suite
	ts           *tests.TestSuite
	dbPool       *pgxpool.Pool
	redisClient  *redis.Client
	topupService service.TopupService
	userRepo     repository.UserRepository
	cardRepo     repository.CardRepository
	saldoRepo    repository.SaldoRepository
}

func (s *TopupServiceTestSuite) SetupSuite() {
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
	cacheTopup := topup_cache.NewTopupMencache(cacheStore)

	s.topupService = service.NewTopupService(service.TopupServiceDeps{
		TopupRepo:     repos.Topup,
		CardRepo:      repos.Card,
		SaldoRepo:     repos.Saldo,
		Logger:        log,
		Observability: obs,
		Cache:         cacheTopup,
	})
}

func (s *TopupServiceTestSuite) TearDownSuite() {
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *TopupServiceTestSuite) TestTopupLifecycle() {
	ctx := context.Background()

	// Seed User and Card
	user, _ := s.userRepo.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Topup",
		LastName:  "Owner",
		Email:     "topup.service@example.com",
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
		TotalBalance: 0,
	})

	// 1. Create Topup
	req := &requests.CreateTopupRequest{
		CardNumber:  card.CardNumber,
		TopupAmount: 100000,
		TopupMethod: "visa",
	}
	topup, err := s.topupService.CreateTopup(ctx, req)
	s.NoError(err)
	s.NotNil(topup)
	s.Equal(int32(req.TopupAmount), topup.TopupAmount)

	// 2. Find By ID
	found, err := s.topupService.FindById(ctx, int(topup.TopupID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(topup.TopupAmount, found.TopupAmount)
}

func TestTopupServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TopupServiceTestSuite))
}
