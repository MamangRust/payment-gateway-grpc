package repository_test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/tests"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type SaldoRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.SaldoRepository
}

func (s *SaldoRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewSaldoRepository(queries)
}

func (s *SaldoRepositoryTestSuite) TearDownSuite() {
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *SaldoRepositoryTestSuite) TestCreateSaldo() {
	// Seed card first
	userReq := &requests.CreateUserRequest{
		FirstName: "Saldo",
		LastName:  "Owner",
		Email:     "saldoowner@example.com",
		Password:  "password123",
	}
	queries := db.New(s.dbPool)
	userRepo := repository.NewUserRepository(queries)
	user, err := userRepo.CreateUser(context.Background(), userReq)
	s.Require().NoError(err)

	cardRepo := repository.NewCardRepository(queries)
	card, err := cardRepo.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       int(user.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	s.Require().NoError(err)

	req := &requests.CreateSaldoRequest{
		CardNumber:   card.CardNumber,
		TotalBalance: 1000000,
	}

	saldo, err := s.repo.CreateSaldo(context.Background(), req)
	s.NoError(err)
	s.Require().NotNil(saldo)
	s.Equal(int32(req.TotalBalance), saldo.TotalBalance)
}

func TestSaldoRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SaldoRepositoryTestSuite))
}
