package service_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	merchant_cache "MamangRust/paymentgatewaygrpc/internal/cache/merchant"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/internal/service"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
	"MamangRust/paymentgatewaygrpc/tests"
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type MerchantServiceTestSuite struct {
	suite.Suite
	ts              *tests.TestSuite
	dbPool          *pgxpool.Pool
	redisClient     *redis.Client
	merchantService service.MerchantService
	userRepo        repository.UserRepository
}

func (s *MerchantServiceTestSuite) SetupSuite() {
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

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	obs, _ := observability.NewObservability("test", log)
	cacheStore := cache.NewCacheStore(s.redisClient, log, nil)
	cacheMerchant := merchant_cache.NewMerchantMencache(cacheStore)

	s.merchantService = service.NewMerchantService(service.MerchantServiceDeps{
		MerchantRepo:  repos.Merchant,
		UserRepo:      repos.User,
		Logger:        log,
		Observability: obs,
		Cache:         cacheMerchant,
	})
}

func (s *MerchantServiceTestSuite) TearDownSuite() {
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *MerchantServiceTestSuite) TestMerchantLifecycle() {
	ctx := context.Background()

	// Seed User
	user, _ := s.userRepo.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Merchant",
		LastName:  "Owner",
		Email:     "merchant.service@example.com",
		Password:  "password123",
	})

	// 1. Create Merchant
	req := &requests.CreateMerchantRequest{
		Name:   "Service Merchant Ltd",
		UserID: int(user.UserID),
	}
	merchant, err := s.merchantService.CreateMerchant(ctx, req)
	s.NoError(err)
	s.NotNil(merchant)
	s.Equal(req.Name, merchant.Name)

	// 2. Find By ID
	found, err := s.merchantService.FindById(ctx, int(merchant.MerchantID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(merchant.ApiKey, found.ApiKey)

	// 3. Find By API Key
	foundByApi, err := s.merchantService.FindByApiKey(ctx, merchant.ApiKey)
	s.NoError(err)
	s.NotNil(foundByApi)
	s.Equal(merchant.MerchantID, foundByApi.MerchantID)
}

func TestMerchantServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantServiceTestSuite))
}
