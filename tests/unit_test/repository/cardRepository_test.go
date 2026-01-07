package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CardRepositorySuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	mockRepo *mocks.MockCardRepository
}

func (suite *CardRepositorySuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockCardRepository(suite.mockCtrl)
}

func (suite *CardRepositorySuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *CardRepositorySuite) TestFindAllCards_Success() {
	req := &requests.FindAllCards{Search: "debit", Page: 1, PageSize: 10}
	cards := []*record.CardRecord{{ID: 1, CardNumber: "1111-xxxx-xxxx-1111"}}
	total := 1

	suite.mockRepo.EXPECT().FindAllCards(req).Return(cards, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllCards(req)

	suite.NoError(err)
	suite.Equal(cards, result)
	suite.Equal(1, *totalRes)
}

func (suite *CardRepositorySuite) TestFindByActive_Success() {
	req := &requests.FindAllCards{Search: "active", Page: 1, PageSize: 10}
	activeCards := []*record.CardRecord{{ID: 1, CardNumber: "2222-xxxx-xxxx-2222"}}
	total := 1

	suite.mockRepo.EXPECT().FindByActive(req).Return(activeCards, &total, nil)

	result, totalRes, err := suite.mockRepo.FindByActive(req)

	suite.NoError(err)
	suite.Equal(activeCards, result)
	suite.Equal(1, *totalRes)
}

func (suite *CardRepositorySuite) TestFindByTrashed_Success() {
	req := &requests.FindAllCards{Page: 1, PageSize: 10}
	deletedAt := "2024-01-03T00:00:00Z"
	trashedCards := []*record.CardRecord{{ID: 2, CardNumber: "3333-xxxx-xxxx-3333", DeletedAt: &deletedAt}}
	total := 1

	suite.mockRepo.EXPECT().FindByTrashed(req).Return(trashedCards, &total, nil)

	result, totalRes, err := suite.mockRepo.FindByTrashed(req)

	suite.NoError(err)
	suite.Equal(trashedCards, result)
	suite.Equal(1, *totalRes)
}

func (suite *CardRepositorySuite) TestFindById_Success() {
	cardID := 1
	card := &record.CardRecord{ID: cardID, CardNumber: "1111-xxxx-xxxx-1111"}

	suite.mockRepo.EXPECT().FindById(cardID).Return(card, nil)

	result, err := suite.mockRepo.FindById(cardID)

	suite.NoError(err)
	suite.Equal(card, result)
}

func (suite *CardRepositorySuite) TestFindById_NotFound() {
	cardID := 99
	expectedErr := errors.New("card not found")

	suite.mockRepo.EXPECT().FindById(cardID).Return(nil, expectedErr)

	result, err := suite.mockRepo.FindById(cardID)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(expectedErr, err)
}

func (suite *CardRepositorySuite) TestFindCardByUserId_Success() {
	userID := 123
	card := &record.CardRecord{UserID: userID, CardNumber: "4444-xxxx-xxxx-4444"}

	suite.mockRepo.EXPECT().FindCardByUserId(userID).Return(card, nil)

	result, err := suite.mockRepo.FindCardByUserId(userID)

	suite.NoError(err)
	suite.Equal(card, result)
}

func (suite *CardRepositorySuite) TestCreateCard_Success() {
	expireDate := time.Now().AddDate(4, 0, 0)
	req := &requests.CreateCardRequest{
		UserID:       123,
		CardType:     "debit",
		ExpireDate:   expireDate,
		CVV:          "123",
		CardProvider: "VISA",
	}
	createdCard := &record.CardRecord{ID: 10, UserID: 123, CardNumber: "1234-5678-...", ExpireDate: expireDate.Format("2006-01-02")}

	suite.mockRepo.EXPECT().CreateCard(req).Return(createdCard, nil)

	result, err := suite.mockRepo.CreateCard(req)

	suite.NoError(err)
	suite.Equal(createdCard, result)
}

func (suite *CardRepositorySuite) TestUpdateCard_Success() {
	cardID := 1
	req := &requests.UpdateCardRequest{CardID: cardID, UserID: 456}
	updatedCard := &record.CardRecord{ID: cardID, UserID: 456}

	suite.mockRepo.EXPECT().UpdateCard(req).Return(updatedCard, nil)

	result, err := suite.mockRepo.UpdateCard(req)

	suite.NoError(err)
	suite.Equal(updatedCard, result)
}

func (suite *CardRepositorySuite) TestTrashedCard_Success() {
	cardID := 1
	deletedAt := "2024-01-03T00:00:00Z"
	trashedCard := &record.CardRecord{ID: cardID, DeletedAt: &deletedAt}

	suite.mockRepo.EXPECT().TrashedCard(cardID).Return(trashedCard, nil)

	result, err := suite.mockRepo.TrashedCard(cardID)

	suite.NoError(err)
	suite.NotNil(result.DeletedAt)
	suite.Equal(trashedCard, result)
}

func (suite *CardRepositorySuite) TestRestoreCard_Success() {
	cardID := 1
	restoredCard := &record.CardRecord{ID: cardID, DeletedAt: nil}

	suite.mockRepo.EXPECT().RestoreCard(cardID).Return(restoredCard, nil)

	result, err := suite.mockRepo.RestoreCard(cardID)

	suite.NoError(err)
	suite.Nil(result.DeletedAt)
	suite.Equal(restoredCard, result)
}

func (suite *CardRepositorySuite) TestDeleteCardPermanent_Success() {
	cardID := 1

	suite.mockRepo.EXPECT().DeleteCardPermanent(cardID).Return(true, nil)

	result, err := suite.mockRepo.DeleteCardPermanent(cardID)

	suite.NoError(err)
	suite.True(result)
}

func (suite *CardRepositorySuite) TestRestoreAllCard_Success() {
	suite.mockRepo.EXPECT().RestoreAllCard().Return(true, nil)

	result, err := suite.mockRepo.RestoreAllCard()

	suite.NoError(err)
	suite.True(result)
}

func (suite *CardRepositorySuite) TestDeleteAllCardPermanent_Success() {
	suite.mockRepo.EXPECT().DeleteAllCardPermanent().Return(true, nil)

	result, err := suite.mockRepo.DeleteAllCardPermanent()

	suite.NoError(err)
	suite.True(result)
}

func (suite *CardRepositorySuite) TestGetTotalBalances_Success() {
	totalBalance := int64(5000000)

	suite.mockRepo.EXPECT().GetTotalBalances().Return(&totalBalance, nil)

	result, err := suite.mockRepo.GetTotalBalances()

	suite.NoError(err)
	suite.Equal(int64(5000000), *result)
}

func (suite *CardRepositorySuite) TestGetTotalTopAmount_Success() {
	totalTopup := int64(1200000)

	suite.mockRepo.EXPECT().GetTotalTopAmount().Return(&totalTopup, nil)

	result, err := suite.mockRepo.GetTotalTopAmount()

	suite.NoError(err)
	suite.Equal(int64(1200000), *result)
}

func (suite *CardRepositorySuite) TestGetTotalBalanceByCardNumber_Success() {
	cardNumber := "1234-5678-9012-3456"
	balance := int64(120000)

	suite.mockRepo.EXPECT().GetTotalBalanceByCardNumber(cardNumber).Return(&balance, nil)

	result, err := suite.mockRepo.GetTotalBalanceByCardNumber(cardNumber)

	suite.NoError(err)
	suite.Equal(int64(120000), *result)
}

func (suite *CardRepositorySuite) TestGetTotalTransferAmountBySender_Success() {
	senderCardNumber := "1111-2222-3333-4444"
	totalTransfer := int64(800000)

	suite.mockRepo.EXPECT().GetTotalTransferAmountBySender(senderCardNumber).Return(&totalTransfer, nil)

	result, err := suite.mockRepo.GetTotalTransferAmountBySender(senderCardNumber)

	suite.NoError(err)
	suite.Equal(int64(800000), *result)
}


func (suite *CardRepositorySuite) TestGetMonthlyBalance_Success() {
	year := 2024
	monthlyBalances := []*record.CardMonthBalance{
		{Month: "2024-01", TotalBalance: 1000000},
		{Month: "2024-02", TotalBalance: 1500000},
	}

	suite.mockRepo.EXPECT().GetMonthlyBalance(year).Return(monthlyBalances, nil)

	result, err := suite.mockRepo.GetMonthlyBalance(year)

	suite.NoError(err)
	suite.Equal(monthlyBalances, result)
}

func (suite *CardRepositorySuite) TestGetYearlyBalance_Success() {
	year := 2023
	yearlyBalance := []*record.CardYearlyBalance{
		{Year: "2023", TotalBalance: 20000000},
	}

	suite.mockRepo.EXPECT().GetYearlyBalance(year).Return(yearlyBalance, nil)

	result, err := suite.mockRepo.GetYearlyBalance(year)

	suite.NoError(err)
	suite.Equal(yearlyBalance, result)
}

func (suite *CardRepositorySuite) TestGetMonthlyBalancesByCardNumber_Success() {
	req := &requests.MonthYearCardNumberCard{
		CardNumber: "1234-5678-9012-3456",
		Year:       2024,
	}
	monthlyBalances := []*record.CardMonthBalance{
		{Month: "2024-01", TotalBalance: 500000},
		{Month: "2024-02", TotalBalance: 700000},
	}

	suite.mockRepo.EXPECT().GetMonthlyBalancesByCardNumber(req).Return(monthlyBalances, nil)

	result, err := suite.mockRepo.GetMonthlyBalancesByCardNumber(req)

	suite.NoError(err)
	suite.Equal(monthlyBalances, result)
}

func (suite *CardRepositorySuite) TestGetYearlyTopupAmountByCardNumber_Success() {
	req := &requests.MonthYearCardNumberCard{
		CardNumber: "1234-5678-9012-3456",
		Year:       2023,
	}
	yearlyTopup := []*record.CardYearAmount{
		{Year: "2023", TotalAmount: 5000000},
	}

	suite.mockRepo.EXPECT().GetYearlyTopupAmountByCardNumber(req).Return(yearlyTopup, nil)

	result, err := suite.mockRepo.GetYearlyTopupAmountByCardNumber(req)

	suite.NoError(err)
	suite.Equal(yearlyTopup, result)
}

func (suite *CardRepositorySuite) TestGetMonthlyTransferAmountSender_Success() {
	year := 2024
	monthlyTransfer := []*record.CardMonthAmount{
		{Month: "2024-01", TotalAmount: 200000},
		{Month: "2024-02", TotalAmount: 150000},
	}

	suite.mockRepo.EXPECT().GetMonthlyTransferAmountSender(year).Return(monthlyTransfer, nil)

	result, err := suite.mockRepo.GetMonthlyTransferAmountSender(year)

	suite.NoError(err)
	suite.Equal(monthlyTransfer, result)
}

func (suite *CardRepositorySuite) TestGetYearlyTransferAmountReceiver_Success() {
	year := 2024
	yearlyTransfer := []*record.CardYearAmount{
		{Year: "2024", TotalAmount: 1500000},
	}

	suite.mockRepo.EXPECT().GetYearlyTransferAmountReceiver(year).Return(yearlyTransfer, nil)

	result, err := suite.mockRepo.GetYearlyTransferAmountReceiver(year)

	suite.NoError(err)
	suite.Equal(yearlyTransfer, result)
}

func TestCardRepositorySuite(t *testing.T) {
	suite.Run(t, new(CardRepositorySuite))
}
