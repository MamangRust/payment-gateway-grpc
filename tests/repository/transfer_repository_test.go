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

type TransferRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.TransferRepository
}

func (s *TransferRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewTransferRepository(queries)
}

func (s *TransferRepositoryTestSuite) TearDownSuite() {
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *TransferRepositoryTestSuite) TestCreateTransfer() {
	// Seed source card
	user1Req := &requests.CreateUserRequest{
		FirstName: "Sender",
		LastName:  "Test",
		Email:     "sender@example.com",
		Password:  "password123",
	}
	queries := db.New(s.dbPool)
	userRepo := repository.NewUserRepository(queries)
	user1, _ := userRepo.CreateUser(context.Background(), user1Req)

	cardRepo := repository.NewCardRepository(queries)
	card1, _ := cardRepo.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       int(user1.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "111",
		CardProvider: "visa",
	})

	// Seed destination card
	user2Req := &requests.CreateUserRequest{
		FirstName: "Receiver",
		LastName:  "Test",
		Email:     "receiver@example.com",
		Password:  "password123",
	}
	user2, _ := userRepo.CreateUser(context.Background(), user2Req)
	card2, _ := cardRepo.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       int(user2.UserID),
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "222",
		CardProvider: "mastercard",
	})

	req := &requests.CreateTransferRequest{
		TransferFrom:   card1.CardNumber,
		TransferTo:     card2.CardNumber,
		TransferAmount: 25000,
	}

	transfer, err := s.repo.CreateTransfer(context.Background(), req)
	s.NoError(err)
	s.NotNil(transfer)
	s.Equal(int32(req.TransferAmount), transfer.TransferAmount)
}

func TestTransferRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransferRepositoryTestSuite))
}
