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

type TransactionRepositorySuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	mockRepo *mocks.MockTransactionRepository
}

func (suite *TransactionRepositorySuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockTransactionRepository(suite.mockCtrl)
}

func (suite *TransactionRepositorySuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TransactionRepositorySuite) TestFindAllTransactions_Success() {
	req := &requests.FindAllTransactions{Search: "success", Page: 1, PageSize: 10}
	transactions := []*record.TransactionRecord{{ID: 1, TransactionNo: "TX001"}}
	total := 1

	suite.mockRepo.EXPECT().FindAllTransactions(req).Return(transactions, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllTransactions(req)

	suite.NoError(err)
	suite.Equal(transactions, result)
	suite.Equal(1, *totalRes)
}

func (suite *TransactionRepositorySuite) TestFindByIdTransaction_Success() {
	transactionID := 1
	transaction := &record.TransactionRecord{ID: transactionID, TransactionNo: "TX001"}

	suite.mockRepo.EXPECT().FindById(transactionID).Return(transaction, nil)

	result, err := suite.mockRepo.FindById(transactionID)

	suite.NoError(err)
	suite.Equal(transaction, result)
}

func (suite *TransactionRepositorySuite) TestFindByIdTransaction_NotFound() {
	transactionID := 99
	expectedErr := errors.New("transaction not found")

	suite.mockRepo.EXPECT().FindById(transactionID).Return(nil, expectedErr)

	result, err := suite.mockRepo.FindById(transactionID)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(expectedErr, err)
}

func (suite *TransactionRepositorySuite) TestFindAllTransactionByCardNumber_Success() {
	req := &requests.FindAllTransactionCardNumber{CardNumber: "1111-xxxx", Page: 1, PageSize: 10}
	transactions := []*record.TransactionRecord{{ID: 2, CardNumber: "1111-xxxx"}}
	total := 1

	suite.mockRepo.EXPECT().FindAllTransactionByCardNumber(req).Return(transactions, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllTransactionByCardNumber(req)

	suite.NoError(err)
	suite.Equal(transactions, result)
	suite.Equal(1, *totalRes)
}

func (suite *TransactionRepositorySuite) TestFindTransactionByMerchantId_Success() {
	merchantID := 5
	transactions := []*record.TransactionRecord{{ID: 3, MerchantID: 5}}

	suite.mockRepo.EXPECT().FindTransactionByMerchantId(merchantID).Return(transactions, nil)

	result, err := suite.mockRepo.FindTransactionByMerchantId(merchantID)

	suite.NoError(err)
	suite.Equal(transactions, result)
}

func (suite *TransactionRepositorySuite) TestCreateTransaction_Success() {
	merchantID := 1
	req := &requests.CreateTransactionRequest{
		CardNumber:      "2222-xxxx",
		Amount:          100000,
		PaymentMethod:   "credit_card",
		MerchantID:      &merchantID,
		TransactionTime: time.Now(),
	}
	createdTransaction := &record.TransactionRecord{ID: 10, CardNumber: "2222-xxxx", TransactionNo: "TX010"}

	suite.mockRepo.EXPECT().CreateTransaction(req).Return(createdTransaction, nil)

	result, err := suite.mockRepo.CreateTransaction(req)

	suite.NoError(err)
	suite.Equal(createdTransaction, result)
}

func (suite *TransactionRepositorySuite) TestUpdateTransactionStatus_Success() {
	req := &requests.UpdateTransactionStatus{TransactionID: 1, Status: "failed"}
	updatedTransaction := &record.TransactionRecord{ID: 10, CardNumber: "2222-xxxx", TransactionNo: "TX010"}

	suite.mockRepo.EXPECT().UpdateTransactionStatus(req).Return(updatedTransaction, nil)

	result, err := suite.mockRepo.UpdateTransactionStatus(req)

	suite.NoError(err)

	suite.Equal(updatedTransaction, result)
}

func (suite *TransactionRepositorySuite) TestTrashedTransaction_Success() {
	transactionID := 1
	trashedTime := "2024-01-01T00:00:00Z"
	trashedTransaction := &record.TransactionRecord{ID: transactionID, DeletedAt: &trashedTime}

	suite.mockRepo.EXPECT().TrashedTransaction(transactionID).Return(trashedTransaction, nil)

	result, err := suite.mockRepo.TrashedTransaction(transactionID)

	suite.NoError(err)
	suite.NotNil(result.DeletedAt)
	suite.Equal(trashedTransaction, result)
}

func (suite *TransactionRepositorySuite) TestDeleteTransactionPermanent_Success() {
	transactionID := 1

	suite.mockRepo.EXPECT().DeleteTransactionPermanent(transactionID).Return(true, nil)

	result, err := suite.mockRepo.DeleteTransactionPermanent(transactionID)

	suite.NoError(err)
	suite.True(result)
}

func (suite *TransactionRepositorySuite) TestRestoreAllTransaction_Success() {
	suite.mockRepo.EXPECT().RestoreAllTransaction().Return(true, nil)

	result, err := suite.mockRepo.RestoreAllTransaction()

	suite.NoError(err)
	suite.True(result)
}

func (suite *TransactionRepositorySuite) TestGetMonthTransactionStatusSuccess_Success() {
	req := &requests.MonthStatusTransaction{Year: 2024, Month: 1}
	monthlySuccess := []*record.TransactionRecordMonthStatusSuccess{
		{Year: "2024", Month: "01", TotalSuccess: 100, TotalAmount: 10000000},
	}

	suite.mockRepo.EXPECT().GetMonthTransactionStatusSuccess(req).Return(monthlySuccess, nil)

	result, err := suite.mockRepo.GetMonthTransactionStatusSuccess(req)

	suite.NoError(err)
	suite.Equal(monthlySuccess, result)
}

func (suite *TransactionRepositorySuite) TestGetYearlyTransactionStatusFailed_Success() {
	year := 2023
	yearlyFailed := []*record.TransactionRecordYearStatusFailed{
		{Year: "2023", TotalFailed: 5, TotalAmount: 250000},
	}

	suite.mockRepo.EXPECT().GetYearlyTransactionStatusFailed(year).Return(yearlyFailed, nil)

	result, err := suite.mockRepo.GetYearlyTransactionStatusFailed(year)

	suite.NoError(err)
	suite.Equal(yearlyFailed, result)
}

func (suite *TransactionRepositorySuite) TestGetMonthTransactionStatusSuccessByCardNumber_Success() {
	req := &requests.MonthStatusTransactionCardNumber{CardNumber: "1111-xxxx", Year: 2024, Month: 1}
	monthlySuccess := []*record.TransactionRecordMonthStatusSuccess{
		{Year: "2024", Month: "01", TotalSuccess: 10, TotalAmount: 1000000},
	}

	suite.mockRepo.EXPECT().GetMonthTransactionStatusSuccessByCardNumber(req).Return(monthlySuccess, nil)

	result, err := suite.mockRepo.GetMonthTransactionStatusSuccessByCardNumber(req)

	suite.NoError(err)
	suite.Equal(monthlySuccess, result)
}

func (suite *TransactionRepositorySuite) TestGetYearlyTransactionStatusFailedByCardNumber_Success() {
	req := &requests.YearStatusTransactionCardNumber{CardNumber: "2222-xxxx", Year: 2023}
	yearlyFailed := []*record.TransactionRecordYearStatusFailed{
		{Year: "2023", TotalFailed: 1, TotalAmount: 50000},
	}

	suite.mockRepo.EXPECT().GetYearlyTransactionStatusFailedByCardNumber(req).Return(yearlyFailed, nil)

	result, err := suite.mockRepo.GetYearlyTransactionStatusFailedByCardNumber(req)

	suite.NoError(err)
	suite.Equal(yearlyFailed, result)
}

func (suite *TransactionRepositorySuite) TestGetMonthlyPaymentMethods_Success() {
	year := 2024
	monthlyMethods := []*record.TransactionMonthMethod{
		{Month: "2024-01", PaymentMethod: "credit_card", TotalTransactions: 50, TotalAmount: 5000000},
		{Month: "2024-01", PaymentMethod: "debit_card", TotalTransactions: 50, TotalAmount: 5000000},
	}

	suite.mockRepo.EXPECT().GetMonthlyPaymentMethods(year).Return(monthlyMethods, nil)

	result, err := suite.mockRepo.GetMonthlyPaymentMethods(year)

	suite.NoError(err)
	suite.Equal(monthlyMethods, result)
}

func (suite *TransactionRepositorySuite) TestGetYearlyAmounts_Success() {
	year := 2023
	yearlyAmounts := []*record.TransactionYearlyAmount{
		{Year: "2023", TotalAmount: 120000000},
	}

	suite.mockRepo.EXPECT().GetYearlyAmounts(year).Return(yearlyAmounts, nil)

	result, err := suite.mockRepo.GetYearlyAmounts(year)

	suite.NoError(err)
	suite.Equal(yearlyAmounts, result)
}

func (suite *TransactionRepositorySuite) TestGetMonthlyPaymentMethodsByCardNumber_Success() {
	req := &requests.MonthYearPaymentMethod{CardNumber: "1111-xxxx", Year: 2024}
	monthlyMethods := []*record.TransactionMonthMethod{
		{Month: "2024-01", PaymentMethod: "credit_card", TotalTransactions: 5, TotalAmount: 500000},
	}

	suite.mockRepo.EXPECT().GetMonthlyPaymentMethodsByCardNumber(req).Return(monthlyMethods, nil)

	result, err := suite.mockRepo.GetMonthlyPaymentMethodsByCardNumber(req)

	suite.NoError(err)
	suite.Equal(monthlyMethods, result)
}

func (suite *TransactionRepositorySuite) TestGetYearlyAmountsByCardNumber_Success() {
	req := &requests.MonthYearPaymentMethod{CardNumber: "1111-xxxx", Year: 2023}
	yearlyAmounts := []*record.TransactionYearlyAmount{
		{Year: "2023", TotalAmount: 12000000},
	}

	suite.mockRepo.EXPECT().GetYearlyAmountsByCardNumber(req).Return(yearlyAmounts, nil)

	result, err := suite.mockRepo.GetYearlyAmountsByCardNumber(req)

	suite.NoError(err)
	suite.Equal(yearlyAmounts, result)
}

func TestTransactionRepositorySuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositorySuite))
}
