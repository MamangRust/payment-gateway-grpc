package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type MerchantRepositorySuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	mockRepo *mocks.MockMerchantRepository
}

func (suite *MerchantRepositorySuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockMerchantRepository(suite.mockCtrl)
}

func (suite *MerchantRepositorySuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *MerchantRepositorySuite) TestFindAllMerchants_Success() {
	req := &requests.FindAllMerchants{Search: "store", Page: 1, PageSize: 10}
	merchants := []*record.MerchantRecord{{ID: 1, Name: "Test Store"}}
	total := 1

	suite.mockRepo.EXPECT().FindAllMerchants(req).Return(merchants, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllMerchants(req)

	suite.NoError(err)
	suite.Equal(merchants, result)
	suite.Equal(1, *totalRes)
}

func (suite *MerchantRepositorySuite) TestFindByIdMerchant_Success() {
	merchantID := 1
	merchant := &record.MerchantRecord{ID: merchantID, Name: "Test Store"}

	suite.mockRepo.EXPECT().FindById(merchantID).Return(merchant, nil)

	result, err := suite.mockRepo.FindById(merchantID)

	suite.NoError(err)
	suite.Equal(merchant, result)
}

func (suite *MerchantRepositorySuite) TestFindByIdMerchant_NotFound() {
	merchantID := 99
	expectedErr := errors.New("merchant not found")

	suite.mockRepo.EXPECT().FindById(merchantID).Return(nil, expectedErr)

	result, err := suite.mockRepo.FindById(merchantID)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(expectedErr, err)
}

func (suite *MerchantRepositorySuite) TestFindByApiKey_Success() {
	apiKey := "secret-api-key-123"
	merchant := &record.MerchantRecord{ApiKey: apiKey, Name: "API Store"}

	suite.mockRepo.EXPECT().FindByApiKey(apiKey).Return(merchant, nil)

	result, err := suite.mockRepo.FindByApiKey(apiKey)

	suite.NoError(err)
	suite.Equal(merchant, result)
}

func (suite *MerchantRepositorySuite) TestCreateMerchant_Success() {
	req := &requests.CreateMerchantRequest{Name: "New Merchant", UserID: 123}
	createdMerchant := &record.MerchantRecord{ID: 10, Name: "New Merchant", ApiKey: "new-key-456"}

	suite.mockRepo.EXPECT().CreateMerchant(req).Return(createdMerchant, nil)

	result, err := suite.mockRepo.CreateMerchant(req)

	suite.NoError(err)
	suite.Equal(createdMerchant, result)
}

func (suite *MerchantRepositorySuite) TestUpdateMerchantStatus_Success() {
	req := &requests.UpdateMerchantStatus{MerchantID: 1, Status: "inactive"}
	updatedMerchant := &record.MerchantRecord{ID: 1, Status: "inactive"}

	suite.mockRepo.EXPECT().UpdateMerchantStatus(req).Return(updatedMerchant, nil)

	result, err := suite.mockRepo.UpdateMerchantStatus(req)

	suite.NoError(err)
	suite.Equal("inactive", result.Status)
}

func (suite *MerchantRepositorySuite) TestDeleteMerchantPermanent_Success() {
	merchantID := 1

	suite.mockRepo.EXPECT().DeleteMerchantPermanent(merchantID).Return(true, nil)

	result, err := suite.mockRepo.DeleteMerchantPermanent(merchantID)

	suite.NoError(err)
	suite.True(result)
}

func (suite *MerchantRepositorySuite) TestRestoreAllMerchant_Success() {
	suite.mockRepo.EXPECT().RestoreAllMerchant().Return(true, nil)

	result, err := suite.mockRepo.RestoreAllMerchant()

	suite.NoError(err)
	suite.True(result)
}



func (suite *MerchantRepositorySuite) TestGetMonthlyTotalAmountMerchant_Success() {
	year := 2024
	monthlyTotals := []*record.MerchantMonthlyTotalAmount{
		{Year: "2024", Month: "01", TotalAmount: 10000000},
		{Year: "2024", Month: "02", TotalAmount: 15000000},
	}

	suite.mockRepo.EXPECT().GetMonthlyTotalAmountMerchant(year).Return(monthlyTotals, nil)

	result, err := suite.mockRepo.GetMonthlyTotalAmountMerchant(year)

	suite.NoError(err)
	suite.Equal(monthlyTotals, result)
}

func (suite *MerchantRepositorySuite) TestFindAllTransactionsMerchant_Success() {
	req := &requests.FindAllMerchantTransactions{Search: "success", Page: 1, PageSize: 10}
	transactions := []*record.MerchantTransactionsRecord{{TransactionID: 1, Amount: 50000}}
	total := 1

	suite.mockRepo.EXPECT().FindAllTransactions(req).Return(transactions, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllTransactions(req)

	suite.NoError(err)
	suite.Equal(transactions, result)
	suite.Equal(1, *totalRes)
}

func (suite *MerchantRepositorySuite) TestGetYearlyPaymentMethodMerchant_Success() {
	year := 2023
	yearlyMethods := []*record.MerchantYearlyPaymentMethod{
		{Year: "2023", PaymentMethod: "credit_card", TotalAmount: 500000000},
	}

	suite.mockRepo.EXPECT().GetYearlyPaymentMethodMerchant(year).Return(yearlyMethods, nil)

	result, err := suite.mockRepo.GetYearlyPaymentMethodMerchant(year)

	suite.NoError(err)
	suite.Equal(yearlyMethods, result)
}


func (suite *MerchantRepositorySuite) TestFindAllTransactionsByMerchant_Success() {
	req := &requests.FindAllMerchantTransactionsById{MerchantID: 1, Page: 1, PageSize: 10}
	transactions := []*record.MerchantTransactionsRecord{{TransactionID: 101, MerchantID: 1}}
	total := 1

	suite.mockRepo.EXPECT().FindAllTransactionsByMerchant(req).Return(transactions, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllTransactionsByMerchant(req)

	suite.NoError(err)
	suite.Equal(transactions, result)
	suite.Equal(1, *totalRes)
}

func (suite *MerchantRepositorySuite) TestGetMonthlyPaymentMethodByMerchants_Success() {
	req := &requests.MonthYearPaymentMethodMerchant{MerchantID: 1, Year: 2024}
	monthlyMethods := []*record.MerchantMonthlyPaymentMethod{
		{Month: "2024-01", PaymentMethod: "debit_card", TotalAmount: 2000000},
	}

	suite.mockRepo.EXPECT().GetMonthlyPaymentMethodByMerchants(req).Return(monthlyMethods, nil)

	result, err := suite.mockRepo.GetMonthlyPaymentMethodByMerchants(req)

	suite.NoError(err)
	suite.Equal(monthlyMethods, result)
}

func (suite *MerchantRepositorySuite) TestGetYearlyTotalAmountByMerchants_Success() {
	req := &requests.MonthYearTotalAmountMerchant{MerchantID: 5, Year: 2023}
	yearlyTotals := []*record.MerchantYearlyTotalAmount{
		{Year: "2023", TotalAmount: 120000000},
	}

	suite.mockRepo.EXPECT().GetYearlyTotalAmountByMerchants(req).Return(yearlyTotals, nil)

	result, err := suite.mockRepo.GetYearlyTotalAmountByMerchants(req)

	suite.NoError(err)
	suite.Equal(yearlyTotals, result)
}


func (suite *MerchantRepositorySuite) TestFindAllTransactionsByApikey_Success() {
	req := &requests.FindAllMerchantTransactionsByApiKey{ApiKey: "api-key-xyz", Page: 1, PageSize: 10}
	transactions := []*record.MerchantTransactionsRecord{{TransactionID: 201, MerchantName: "API Store"}}
	total := 1

	suite.mockRepo.EXPECT().FindAllTransactionsByApikey(req).Return(transactions, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllTransactionsByApikey(req)

	suite.NoError(err)
	suite.Equal(transactions, result)
	suite.Equal(1, *totalRes)
}

func (suite *MerchantRepositorySuite) TestGetMonthlyAmountByApikey_Success() {
	req := &requests.MonthYearAmountApiKey{Apikey: "api-key-xyz", Year: 2024}
	monthlyAmounts := []*record.MerchantMonthlyAmount{
		{Month: "2024-01", TotalAmount: 5000000},
		{Month: "2024-02", TotalAmount: 7000000},
	}

	suite.mockRepo.EXPECT().GetMonthlyAmountByApikey(req).Return(monthlyAmounts, nil)

	result, err := suite.mockRepo.GetMonthlyAmountByApikey(req)

	suite.NoError(err)
	suite.Equal(monthlyAmounts, result)
}

func (suite *MerchantRepositorySuite) TestGetYearlyTotalAmountByApikey_Success() {
	req := &requests.MonthYearTotalAmountApiKey{Apikey: "api-key-xyz", Year: 2023}
	yearlyTotal := []*record.MerchantYearlyTotalAmount{
		{Year: "2023", TotalAmount: 150000000},
	}

	suite.mockRepo.EXPECT().GetYearlyTotalAmountByApikey(req).Return(yearlyTotal, nil)

	result, err := suite.mockRepo.GetYearlyTotalAmountByApikey(req)

	suite.NoError(err)
	suite.Equal(yearlyTotal, result)
}

func TestMerchantRepositorySuite(t *testing.T) {
	suite.Run(t, new(MerchantRepositorySuite))
}
