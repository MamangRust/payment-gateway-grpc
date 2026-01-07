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

type SaldoRepositorySuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	mockRepo *mocks.MockSaldoRepository
}

func (suite *SaldoRepositorySuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockSaldoRepository(suite.mockCtrl)
}

func (suite *SaldoRepositorySuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *SaldoRepositorySuite) TestFindAllSaldos_Success() {
	req := &requests.FindAllSaldos{Search: "1234", Page: 1, PageSize: 10}
	saldos := []*record.SaldoRecord{{ID: 1, CardNumber: "1111-xxxx-xxxx-1111"}}
	total := 1

	suite.mockRepo.EXPECT().FindAllSaldos(req).Return(saldos, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllSaldos(req)

	suite.NoError(err)
	suite.Equal(saldos, result)
	suite.Equal(1, *totalRes)
}

func (suite *SaldoRepositorySuite) TestFindByIdSaldo_Success() {
	saldoID := 1
	saldo := &record.SaldoRecord{ID: saldoID, CardNumber: "1111-xxxx-xxxx-1111"}

	suite.mockRepo.EXPECT().FindById(saldoID).Return(saldo, nil)

	result, err := suite.mockRepo.FindById(saldoID)

	suite.NoError(err)
	suite.Equal(saldo, result)
}

func (suite *SaldoRepositorySuite) TestFindByIdSaldo_NotFound() {
	saldoID := 99
	expectedErr := errors.New("saldo not found")

	suite.mockRepo.EXPECT().FindById(saldoID).Return(nil, expectedErr)

	result, err := suite.mockRepo.FindById(saldoID)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(expectedErr, err)
}

func (suite *SaldoRepositorySuite) TestFindByCardNumberSaldo_Success() {
	cardNumber := "1234-5678-9012-3456"
	saldo := &record.SaldoRecord{CardNumber: cardNumber, TotalBalance: 1000000}

	suite.mockRepo.EXPECT().FindByCardNumber(cardNumber).Return(saldo, nil)

	result, err := suite.mockRepo.FindByCardNumber(cardNumber)

	suite.NoError(err)
	suite.Equal(saldo, result)
}

func (suite *SaldoRepositorySuite) TestCreateSaldo_Success() {
	req := &requests.CreateSaldoRequest{CardNumber: "2222-xxxx-xxxx-2222", TotalBalance: 500000}
	createdSaldo := &record.SaldoRecord{ID: 10, CardNumber: "2222-xxxx-xxxx-2222", TotalBalance: 500000}

	suite.mockRepo.EXPECT().CreateSaldo(req).Return(createdSaldo, nil)

	result, err := suite.mockRepo.CreateSaldo(req)

	suite.NoError(err)
	suite.Equal(createdSaldo, result)
}

func (suite *SaldoRepositorySuite) TestUpdateSaldoBalance_Success() {
	req := &requests.UpdateSaldoBalance{CardNumber: "1111-xxxx-xxxx-1111", TotalBalance: 1500000}
	updatedSaldo := &record.SaldoRecord{CardNumber: "1111-xxxx-xxxx-1111", TotalBalance: 1500000}

	suite.mockRepo.EXPECT().UpdateSaldoBalance(req).Return(updatedSaldo, nil)

	result, err := suite.mockRepo.UpdateSaldoBalance(req)

	suite.NoError(err)
	suite.Equal(int(1500000), result.TotalBalance)
}

func (suite *SaldoRepositorySuite) TestUpdateSaldoWithdraw_Success() {
	withdrawAmount := 100000
	withdrawTime := time.Now()
	req := &requests.UpdateSaldoWithdraw{
		CardNumber:     "1111-xxxx-xxxx-1111",
		TotalBalance:   900000,
		WithdrawAmount: &withdrawAmount,
		WithdrawTime:   &withdrawTime,
	}
	updatedSaldo := &record.SaldoRecord{CardNumber: "1111-xxxx-xxxx-1111", TotalBalance: 900000, WithdrawAmount: withdrawAmount}

	suite.mockRepo.EXPECT().UpdateSaldoWithdraw(req).Return(updatedSaldo, nil)

	result, err := suite.mockRepo.UpdateSaldoWithdraw(req)

	suite.NoError(err)
	suite.Equal(int(900000), result.TotalBalance)
	suite.Equal(withdrawAmount, result.WithdrawAmount)
}

func (suite *SaldoRepositorySuite) TestTrashedSaldo_Success() {
	saldoID := 1
	trashedTime := "2024-01-01T00:00:00Z"
	trashedSaldo := &record.SaldoRecord{ID: saldoID, DeletedAt: &trashedTime}

	suite.mockRepo.EXPECT().TrashedSaldo(saldoID).Return(trashedSaldo, nil)

	result, err := suite.mockRepo.TrashedSaldo(saldoID)

	suite.NoError(err)
	suite.NotNil(result.DeletedAt)
	suite.Equal(trashedSaldo, result)
}

func (suite *SaldoRepositorySuite) TestDeleteSaldoPermanent_Success() {
	saldoID := 1

	suite.mockRepo.EXPECT().DeleteSaldoPermanent(saldoID).Return(true, nil)

	result, err := suite.mockRepo.DeleteSaldoPermanent(saldoID)

	suite.NoError(err)
	suite.True(result)
}

func (suite *SaldoRepositorySuite) TestRestoreAllSaldo_Success() {
	suite.mockRepo.EXPECT().RestoreAllSaldo().Return(true, nil)

	result, err := suite.mockRepo.RestoreAllSaldo()

	suite.NoError(err)
	suite.True(result)
}

func (suite *SaldoRepositorySuite) TestGetMonthlyTotalSaldoBalance_Success() {
	req := &requests.MonthTotalSaldoBalance{Year: 2024, Month: 1}
	monthlyTotals := []*record.SaldoMonthTotalBalance{
		{Year: "2024", Month: "01", TotalBalance: 50000000},
	}

	suite.mockRepo.EXPECT().GetMonthlyTotalSaldoBalance(req).Return(monthlyTotals, nil)

	result, err := suite.mockRepo.GetMonthlyTotalSaldoBalance(req)

	suite.NoError(err)
	suite.Equal(monthlyTotals, result)
}

func (suite *SaldoRepositorySuite) TestGetYearTotalSaldoBalance_Success() {
	year := 2023
	yearlyTotals := []*record.SaldoYearTotalBalance{
		{Year: "2023", TotalBalance: 600000000},
	}

	suite.mockRepo.EXPECT().GetYearTotalSaldoBalance(year).Return(yearlyTotals, nil)

	result, err := suite.mockRepo.GetYearTotalSaldoBalance(year)

	suite.NoError(err)
	suite.Equal(yearlyTotals, result)
}

func (suite *SaldoRepositorySuite) TestGetMonthlySaldoBalances_Success() {
	year := 2024
	monthlyBalances := []*record.SaldoMonthSaldoBalance{
		{Month: "2024-01", TotalBalance: 45000000},
		{Month: "2024-02", TotalBalance: 55000000},
	}

	suite.mockRepo.EXPECT().GetMonthlySaldoBalances(year).Return(monthlyBalances, nil)

	result, err := suite.mockRepo.GetMonthlySaldoBalances(year)

	suite.NoError(err)
	suite.Equal(monthlyBalances, result)
}

func (suite *SaldoRepositorySuite) TestGetYearlySaldoBalances_Success() {
	year := 2023
	yearlyBalances := []*record.SaldoYearSaldoBalance{
		{Year: "2023", TotalBalance: 550000000},
	}

	suite.mockRepo.EXPECT().GetYearlySaldoBalances(year).Return(yearlyBalances, nil)

	result, err := suite.mockRepo.GetYearlySaldoBalances(year)

	suite.NoError(err)
	suite.Equal(yearlyBalances, result)
}

func TestSaldoRepositorySuite(t *testing.T) {
	suite.Run(t, new(SaldoRepositorySuite))
}
