package repository_test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/tests"
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type MerchantRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.MerchantRepository
}

func (s *MerchantRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewMerchantRepository(queries)
}

func (s *MerchantRepositoryTestSuite) TearDownSuite() {
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *MerchantRepositoryTestSuite) TestCreateMerchant() {
	// Seed user first
	userReq := &requests.CreateUserRequest{
		FirstName: "Merchant",
		LastName:  "Owner",
		Email:     "merchantowner@example.com",
		Password:  "password123",
	}
	queries := db.New(s.dbPool)
	userRepo := repository.NewUserRepository(queries)
	user, err := userRepo.CreateUser(context.Background(), userReq)
	s.Require().NoError(err)

	req := &requests.CreateMerchantRequest{
		Name:   "Test Merchant",
		UserID: int(user.UserID),
	}

	merchant, err := s.repo.CreateMerchant(context.Background(), req)
	s.NoError(err)
	s.Require().NotNil(merchant)
	s.Equal(req.Name, merchant.Name)
}

func (s *MerchantRepositoryTestSuite) TestFindById() {
	userReq := &requests.CreateUserRequest{
		FirstName: "Merchant",
		LastName:  "Owner2",
		Email:     "merchantowner2@example.com",
		Password:  "password123",
	}
	queries := db.New(s.dbPool)
	userRepo := repository.NewUserRepository(queries)
	user, _ := userRepo.CreateUser(context.Background(), userReq)

	req := &requests.CreateMerchantRequest{
		Name:   "Find Me Merchant",
		UserID: int(user.UserID),
	}
	m, _ := s.repo.CreateMerchant(context.Background(), req)

	found, err := s.repo.FindById(context.Background(), int(m.MerchantID))
	s.NoError(err)
	s.Require().NotNil(found)
	s.Equal(m.Name, found.Name)
}

func TestMerchantRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantRepositoryTestSuite))
}
