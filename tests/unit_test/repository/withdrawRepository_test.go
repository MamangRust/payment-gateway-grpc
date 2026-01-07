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

type WithdrawRepositorySuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	mockRepo *mocks.MockWithdrawRepository
}

func (suite *WithdrawRepositorySuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockWithdrawRepository(suite.mockCtrl)
}

func (suite *WithdrawRepositorySuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *WithdrawRepositorySuite) TestFindAllWithdraw_Success() {
	req := &requests.FindAllWithdraws{Search: "success", Page: 1, PageSize: 10}
	withdraws := []*record.WithdrawRecord{{ID: 1, WithdrawNo: "WD001"}}
	total := 1

	suite.mockRepo.EXPECT().FindAll(req).Return(withdraws, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAll(req)

	suite.NoError(err)
	suite.Equal(withdraws, result)
	suite.Equal(1, *totalRes)
}

func (suite *WithdrawRepositorySuite) TestFindByIdWithdraw_Success() {
	id := 1
	withdraw := &record.WithdrawRecord{ID: id, WithdrawNo: "WD001"}

	suite.mockRepo.EXPECT().FindById(id).Return(withdraw, nil)

	result, err := suite.mockRepo.FindById(id)

	suite.NoError(err)
	suite.Equal(withdraw, result)
}

func (suite *WithdrawRepositorySuite) TestFindByIdWithdraw_NotFound() {
	id := 99
	expectedErr := errors.New("withdraw not found")

	suite.mockRepo.EXPECT().FindById(id).Return(nil, expectedErr)

	result, err := suite.mockRepo.FindById(id)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(expectedErr, err)
}

func (suite *WithdrawRepositorySuite) TestFindAllByCardNumber_Success() {
	req := &requests.FindAllWithdrawCardNumber{CardNumber: "1111-xxxx", Page: 1, PageSize: 10}
	withdraws := []*record.WithdrawRecord{{ID: 2, CardNumber: "1111-xxxx"}}
	total := 1

	suite.mockRepo.EXPECT().FindAllByCardNumber(req).Return(withdraws, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllByCardNumber(req)

	suite.NoError(err)
	suite.Equal(withdraws, result)
	suite.Equal(1, *totalRes)
}

func (suite *WithdrawRepositorySuite) TestCreateWithdraw_Success() {
	req := &requests.CreateWithdrawRequest{
		CardNumber:     "2222-xxxx",
		WithdrawAmount: 100000,
		WithdrawTime:   time.Now(),
	}
	createdWithdraw := &record.WithdrawRecord{ID: 10, CardNumber: "2222-xxxx", WithdrawNo: "WD010"}

	suite.mockRepo.EXPECT().CreateWithdraw(req).Return(createdWithdraw, nil)

	result, err := suite.mockRepo.CreateWithdraw(req)

	suite.NoError(err)
	suite.Equal(createdWithdraw, result)
}

func (suite *WithdrawRepositorySuite) TestTrashedWithdraw_Success() {
	withdrawID := 1
	trashedTime := "2024-01-01T00:00:00Z"
	trashedWithdraw := &record.WithdrawRecord{ID: withdrawID, DeletedAt: &trashedTime}

	suite.mockRepo.EXPECT().TrashedWithdraw(withdrawID).Return(trashedWithdraw, nil)

	result, err := suite.mockRepo.TrashedWithdraw(withdrawID)

	suite.NoError(err)
	suite.NotNil(result.DeletedAt)
	suite.Equal(trashedWithdraw, result)
}

func (suite *WithdrawRepositorySuite) TestDeleteWithdrawPermanent_Success() {
	withdrawID := 1

	suite.mockRepo.EXPECT().DeleteWithdrawPermanent(withdrawID).Return(true, nil)

	result, err := suite.mockRepo.DeleteWithdrawPermanent(withdrawID)

	suite.NoError(err)
	suite.True(result)
}

func (suite *WithdrawRepositorySuite) TestRestoreAllWithdraw_Success() {
	suite.mockRepo.EXPECT().RestoreAllWithdraw().Return(true, nil)

	result, err := suite.mockRepo.RestoreAllWithdraw()

	suite.NoError(err)
	suite.True(result)
}

func (suite *WithdrawRepositorySuite) TestGetMonthWithdrawStatusSuccess_Success() {
	req := &requests.MonthStatusWithdraw{Year: 2024, Month: 1}
	monthlySuccess := []*record.WithdrawRecordMonthStatusSuccess{
		{Year: "2024", Month: "01", TotalSuccess: 50, TotalAmount: 5000000},
	}

	suite.mockRepo.EXPECT().GetMonthWithdrawStatusSuccess(req).Return(monthlySuccess, nil)

	result, err := suite.mockRepo.GetMonthWithdrawStatusSuccess(req)

	suite.NoError(err)
	suite.Equal(monthlySuccess, result)
}

func (suite *WithdrawRepositorySuite) TestGetYearlyWithdrawStatusFailed_Success() {
	year := 2023
	yearlyFailed := []*record.WithdrawRecordYearStatusFailed{
		{Year: "2023", TotalFailed: 2, TotalAmount: 100000},
	}

	suite.mockRepo.EXPECT().GetYearlyWithdrawStatusFailed(year).Return(yearlyFailed, nil)

	result, err := suite.mockRepo.GetYearlyWithdrawStatusFailed(year)

	suite.NoError(err)
	suite.Equal(yearlyFailed, result)
}

func (suite *WithdrawRepositorySuite) TestGetMonthWithdrawStatusSuccessByCardNumber_Success() {
	req := &requests.MonthStatusWithdrawCardNumber{CardNumber: "1111-xxxx", Year: 2024, Month: 1}
	monthlySuccess := []*record.WithdrawRecordMonthStatusSuccess{
		{Year: "2024", Month: "01", TotalSuccess: 5, TotalAmount: 500000},
	}

	suite.mockRepo.EXPECT().GetMonthWithdrawStatusSuccessByCardNumber(req).Return(monthlySuccess, nil)

	result, err := suite.mockRepo.GetMonthWithdrawStatusSuccessByCardNumber(req)

	suite.NoError(err)
	suite.Equal(monthlySuccess, result)
}

func (suite *WithdrawRepositorySuite) TestGetYearlyWithdrawStatusFailedByCardNumber_Success() {
	req := &requests.YearStatusWithdrawCardNumber{CardNumber: "2222-xxxx", Year: 2023}
	yearlyFailed := []*record.WithdrawRecordYearStatusFailed{
		{Year: "2023", TotalFailed: 1, TotalAmount: 50000},
	}

	suite.mockRepo.EXPECT().GetYearlyWithdrawStatusFailedByCardNumber(req).Return(yearlyFailed, nil)

	result, err := suite.mockRepo.GetYearlyWithdrawStatusFailedByCardNumber(req)

	suite.NoError(err)
	suite.Equal(yearlyFailed, result)
}

func (suite *WithdrawRepositorySuite) TestGetMonthlyWithdraws_Success() {
	year := 2024
	monthlyAmounts := []*record.WithdrawMonthlyAmount{
		{Month: "2024-01", TotalAmount: 4500000},
		{Month: "2024-02", TotalAmount: 5500000},
	}

	suite.mockRepo.EXPECT().GetMonthlyWithdraws(year).Return(monthlyAmounts, nil)

	result, err := suite.mockRepo.GetMonthlyWithdraws(year)

	suite.NoError(err)
	suite.Equal(monthlyAmounts, result)
}

func (suite *WithdrawRepositorySuite) TestGetMonthlyWithdrawsByCardNumber_Success() {
	req := &requests.YearMonthCardNumber{CardNumber: "1111-xxxx", Year: 2024}
	monthlyAmounts := []*record.WithdrawMonthlyAmount{
		{Month: "2024-01", TotalAmount: 500000},
		{Month: "2024-02", TotalAmount: 300000},
	}

	suite.mockRepo.EXPECT().GetMonthlyWithdrawsByCardNumber(req).Return(monthlyAmounts, nil)

	result, err := suite.mockRepo.GetMonthlyWithdrawsByCardNumber(req)

	suite.NoError(err)
	suite.Equal(monthlyAmounts, result)
}

func (suite *WithdrawRepositorySuite) TestGetYearlyWithdrawsByCardNumber_Success() {
	req := &requests.YearMonthCardNumber{CardNumber: "1111-xxxx", Year: 2023}
	yearlyAmounts := []*record.WithdrawYearlyAmount{
		{Year: "2023", TotalAmount: 12000000},
	}

	suite.mockRepo.EXPECT().GetYearlyWithdrawsByCardNumber(req).Return(yearlyAmounts, nil)

	result, err := suite.mockRepo.GetYearlyWithdrawsByCardNumber(req)

	suite.NoError(err)
	suite.Equal(yearlyAmounts, result)
}

func TestWithdrawRepositorySuite(t *testing.T) {
	suite.Run(t, new(WithdrawRepositorySuite))
}
