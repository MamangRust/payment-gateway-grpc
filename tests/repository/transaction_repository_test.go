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

type TransactionRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.TransactionRepository
}

func (s *TransactionRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewTransactionRepository(queries)
}

func (s *TransactionRepositoryTestSuite) TearDownSuite() {
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *TransactionRepositoryTestSuite) TestCreateTransaction() {
	// Seed card and merchant first
	userReq := &requests.CreateUserRequest{
		FirstName: "Transaction",
		LastName:  "Owner",
		Email:     "txowner@example.com",
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

	merchantRepo := repository.NewMerchantRepository(queries)
	merchant, err := merchantRepo.CreateMerchant(context.Background(), &requests.CreateMerchantRequest{
		Name:   "Tx Merchant",
		UserID: int(user.UserID),
	})
	s.Require().NoError(err)

	merchantID := int(merchant.MerchantID)
	req := &requests.CreateTransactionRequest{
		CardNumber:      card.CardNumber,
		Amount:          100000,
		MerchantID:      &merchantID,
		PaymentMethod:   "visa",
		TransactionTime: time.Now(),
	}

	transaction, err := s.repo.CreateTransaction(context.Background(), req)
	s.NoError(err)
	s.Require().NotNil(transaction)
	s.Equal(int32(req.Amount), transaction.Amount)
}

func TestTransactionRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionRepositoryTestSuite))
}
