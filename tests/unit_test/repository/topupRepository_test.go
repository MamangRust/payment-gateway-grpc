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

type TopupRepositorySuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	mockRepo *mocks.MockTopupRepository
}

func (suite *TopupRepositorySuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockTopupRepository(suite.mockCtrl)
}

func (suite *TopupRepositorySuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TopupRepositorySuite) TestFindAllTopups_Success() {
	req := &requests.FindAllTopups{Search: "success", Page: 1, PageSize: 10}
	topups := []*record.TopupRecord{{ID: 1, TopupNo: "TP001"}}
	total := 1

	suite.mockRepo.EXPECT().FindAllTopups(req).Return(topups, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllTopups(req)

	suite.NoError(err)
	suite.Equal(topups, result)
	suite.Equal(1, *totalRes)
}

func (suite *TopupRepositorySuite) TestFindByIdTopup_Success() {
	topupID := 1
	topup := &record.TopupRecord{ID: topupID, TopupNo: "TP001"}

	suite.mockRepo.EXPECT().FindById(topupID).Return(topup, nil)

	result, err := suite.mockRepo.FindById(topupID)

	suite.NoError(err)
	suite.Equal(topup, result)
}

func (suite *TopupRepositorySuite) TestFindByIdTopup_NotFound() {
	topupID := 99
	expectedErr := errors.New("topup not found")

	suite.mockRepo.EXPECT().FindById(topupID).Return(nil, expectedErr)

	result, err := suite.mockRepo.FindById(topupID)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(expectedErr, err)
}

func (suite *TopupRepositorySuite) TestFindAllTopupByCardNumber_Success() {
	req := &requests.FindAllTopupsByCardNumber{CardNumber: "1111-xxxx", Page: 1, PageSize: 10}
	topups := []*record.TopupRecord{{ID: 2, CardNumber: "1111-xxxx"}}
	total := 1

	suite.mockRepo.EXPECT().FindAllTopupByCardNumber(req).Return(topups, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllTopupByCardNumber(req)

	suite.NoError(err)
	suite.Equal(topups, result)
	suite.Equal(1, *totalRes)
}

func (suite *TopupRepositorySuite) TestCreateTopup_Success() {
	req := &requests.CreateTopupRequest{CardNumber: "2222-xxxx", TopupAmount: 100000, TopupMethod: "bank_transfer"}
	createdTopup := &record.TopupRecord{ID: 10, CardNumber: "2222-xxxx", TopupNo: "TP010"}

	suite.mockRepo.EXPECT().CreateTopup(req).Return(createdTopup, nil)

	result, err := suite.mockRepo.CreateTopup(req)

	suite.NoError(err)
	suite.Equal(createdTopup, result)
}

func (suite *TopupRepositorySuite) TestUpdateTopupAmount_Success() {
	req := &requests.UpdateTopupAmount{TopupID: 1, TopupAmount: 150000}
	updatedTopup := &record.TopupRecord{ID: 1, TopupAmount: 150000}

	suite.mockRepo.EXPECT().UpdateTopupAmount(req).Return(updatedTopup, nil)

	result, err := suite.mockRepo.UpdateTopupAmount(req)

	suite.NoError(err)
	suite.Equal(int(150000), result.TopupAmount)
}

func (suite *TopupRepositorySuite) TestUpdateTopupStatus_Success() {
	req := &requests.UpdateTopupStatus{TopupID: 1, Status: "failed"}
	updatedTopup := &record.TopupRecord{ID: 10, CardNumber: "2222-xxxx", TopupNo: "TP010"}

	suite.mockRepo.EXPECT().UpdateTopupStatus(req).Return(updatedTopup, nil)

	result, err := suite.mockRepo.UpdateTopupStatus(req)

	suite.NoError(err)
	suite.Equal(updatedTopup, result)
}

func (suite *TopupRepositorySuite) TestTrashedTopup_Success() {
	topupID := 1
	trashedTime := "2024-01-01T00:00:00Z"
	trashedTopup := &record.TopupRecord{ID: topupID, DeletedAt: &trashedTime}

	suite.mockRepo.EXPECT().TrashedTopup(topupID).Return(trashedTopup, nil)

	result, err := suite.mockRepo.TrashedTopup(topupID)

	suite.NoError(err)
	suite.NotNil(result.DeletedAt)
	suite.Equal(trashedTopup, result)
}

func (suite *TopupRepositorySuite) TestDeleteTopupPermanent_Success() {
	topupID := 1

	suite.mockRepo.EXPECT().DeleteTopupPermanent(topupID).Return(true, nil)

	result, err := suite.mockRepo.DeleteTopupPermanent(topupID)

	suite.NoError(err)
	suite.True(result)
}

func (suite *TopupRepositorySuite) TestRestoreAllTopup_Success() {
	suite.mockRepo.EXPECT().RestoreAllTopup().Return(true, nil)

	result, err := suite.mockRepo.RestoreAllTopup()

	suite.NoError(err)
	suite.True(result)
}

func (suite *TopupRepositorySuite) TestGetMonthTopupStatusSuccess_Success() {
	req := &requests.MonthTopupStatus{Year: 2024, Month: 1}
	monthlySuccess := []*record.TopupRecordMonthStatusSuccess{
		{Year: "2024", Month: "01", TotalSuccess: 100, TotalAmount: 10000000},
	}

	suite.mockRepo.EXPECT().GetMonthTopupStatusSuccess(req).Return(monthlySuccess, nil)

	result, err := suite.mockRepo.GetMonthTopupStatusSuccess(req)

	suite.NoError(err)
	suite.Equal(monthlySuccess, result)
}

func (suite *TopupRepositorySuite) TestGetYearlyTopupStatusFailed_Success() {
	year := 2023
	yearlyFailed := []*record.TopupRecordYearStatusFailed{
		{Year: "2023", TotalFailed: 5, TotalAmount: 250000},
	}

	suite.mockRepo.EXPECT().GetYearlyTopupStatusFailed(year).Return(yearlyFailed, nil)

	result, err := suite.mockRepo.GetYearlyTopupStatusFailed(year)

	suite.NoError(err)
	suite.Equal(yearlyFailed, result)
}

func (suite *TopupRepositorySuite) TestGetMonthTopupStatusSuccessByCardNumber_Success() {
	req := &requests.MonthTopupStatusCardNumber{CardNumber: "1111-xxxx", Year: 2024, Month: 1}
	monthlySuccess := []*record.TopupRecordMonthStatusSuccess{
		{Year: "2024", Month: "01", TotalSuccess: 10, TotalAmount: 1000000},
	}

	suite.mockRepo.EXPECT().GetMonthTopupStatusSuccessByCardNumber(req).Return(monthlySuccess, nil)

	result, err := suite.mockRepo.GetMonthTopupStatusSuccessByCardNumber(req)

	suite.NoError(err)
	suite.Equal(monthlySuccess, result)
}

func (suite *TopupRepositorySuite) TestGetYearlyTopupStatusFailedByCardNumber_Success() {
	req := &requests.YearTopupStatusCardNumber{CardNumber: "2222-xxxx", Year: 2023}
	yearlyFailed := []*record.TopupRecordYearStatusFailed{
		{Year: "2023", TotalFailed: 1, TotalAmount: 50000},
	}

	suite.mockRepo.EXPECT().GetYearlyTopupStatusFailedByCardNumber(req).Return(yearlyFailed, nil)

	result, err := suite.mockRepo.GetYearlyTopupStatusFailedByCardNumber(req)

	suite.NoError(err)
	suite.Equal(yearlyFailed, result)
}

func (suite *TopupRepositorySuite) TestGetMonthlyTopupMethods_Success() {
	year := 2024
	monthlyMethods := []*record.TopupMonthMethod{
		{Month: "2024-01", TopupMethod: "bank_transfer", TotalTopups: 50, TotalAmount: 5000000},
		{Month: "2024-01", TopupMethod: "ewallet", TotalTopups: 50, TotalAmount: 5000000},
	}

	suite.mockRepo.EXPECT().GetMonthlyTopupMethods(year).Return(monthlyMethods, nil)

	result, err := suite.mockRepo.GetMonthlyTopupMethods(year)

	suite.NoError(err)
	suite.Equal(monthlyMethods, result)
}

func (suite *TopupRepositorySuite) TestGetYearlyTopupAmounts_Success() {
	year := 2023
	yearlyAmounts := []*record.TopupYearlyAmount{
		{Year: "2023", TotalAmount: 120000000},
	}

	suite.mockRepo.EXPECT().GetYearlyTopupAmounts(year).Return(yearlyAmounts, nil)

	result, err := suite.mockRepo.GetYearlyTopupAmounts(year)

	suite.NoError(err)
	suite.Equal(yearlyAmounts, result)
}

func (suite *TopupRepositorySuite) TestGetMonthlyTopupMethodsByCardNumber_Success() {
	req := &requests.YearMonthMethod{CardNumber: "1111-xxxx", Year: 2024}
	monthlyMethods := []*record.TopupMonthMethod{
		{Month: "2024-01", TopupMethod: "bank_transfer", TotalTopups: 5, TotalAmount: 500000},
	}

	suite.mockRepo.EXPECT().GetMonthlyTopupMethodsByCardNumber(req).Return(monthlyMethods, nil)

	result, err := suite.mockRepo.GetMonthlyTopupMethodsByCardNumber(req)

	suite.NoError(err)
	suite.Equal(monthlyMethods, result)
}

func (suite *TopupRepositorySuite) TestGetYearlyTopupAmountsByCardNumber_Success() {
	req := &requests.YearMonthMethod{CardNumber: "1111-xxxx", Year: 2023}
	yearlyAmounts := []*record.TopupYearlyAmount{
		{Year: "2023", TotalAmount: 12000000},
	}

	suite.mockRepo.EXPECT().GetYearlyTopupAmountsByCardNumber(req).Return(yearlyAmounts, nil)

	result, err := suite.mockRepo.GetYearlyTopupAmountsByCardNumber(req)

	suite.NoError(err)
	suite.Equal(yearlyAmounts, result)
}

func TestTopupRepositorySuite(t *testing.T) {
	suite.Run(t, new(TopupRepositorySuite))
}
