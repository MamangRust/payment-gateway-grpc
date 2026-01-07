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

type TransferRepositorySuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	mockRepo *mocks.MockTransferRepository
}

func (suite *TransferRepositorySuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockTransferRepository(suite.mockCtrl)
}

func (suite *TransferRepositorySuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TransferRepositorySuite) TestFindAll_Success() {
	req := &requests.FindAllTranfers{Search: "success", Page: 1, PageSize: 10}
	transfers := []*record.TransferRecord{{ID: 1, TransferNo: "TF001"}}
	total := 1

	suite.mockRepo.EXPECT().FindAll(req).Return(transfers, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAll(req)

	suite.NoError(err)
	suite.Equal(transfers, result)
	suite.Equal(1, *totalRes)
}

func (suite *TransferRepositorySuite) TestFindByIdTransfer_Success() {
	id := 1
	transfer := &record.TransferRecord{ID: id, TransferNo: "TF001"}

	suite.mockRepo.EXPECT().FindById(id).Return(transfer, nil)

	result, err := suite.mockRepo.FindById(id)

	suite.NoError(err)
	suite.Equal(transfer, result)
}

func (suite *TransferRepositorySuite) TestFindByIdTransfer_NotFound() {
	id := 99
	expectedErr := errors.New("transfer not found")

	suite.mockRepo.EXPECT().FindById(id).Return(nil, expectedErr)

	result, err := suite.mockRepo.FindById(id)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(expectedErr, err)
}

func (suite *TransferRepositorySuite) TestFindTransferByTransferFrom_Success() {
	transferFrom := "1111-xxxx-xxxx-1111"
	transfers := []*record.TransferRecord{{ID: 2, TransferFrom: transferFrom}}

	suite.mockRepo.EXPECT().FindTransferByTransferFrom(transferFrom).Return(transfers, nil)

	result, err := suite.mockRepo.FindTransferByTransferFrom(transferFrom)

	suite.NoError(err)
	suite.Equal(transfers, result)
}

func (suite *TransferRepositorySuite) TestCreateTransfer_Success() {
	req := &requests.CreateTransferRequest{
		TransferFrom:   "1111-xxxx-xxxx-1111",
		TransferTo:     "2222-xxxx-xxxx-2222",
		TransferAmount: 100000,
	}
	createdTransfer := &record.TransferRecord{ID: 10, TransferNo: "TF010"}

	suite.mockRepo.EXPECT().CreateTransfer(req).Return(createdTransfer, nil)

	result, err := suite.mockRepo.CreateTransfer(req)

	suite.NoError(err)
	suite.Equal(createdTransfer, result)
}

func (suite *TransferRepositorySuite) TestUpdateTransferAmount_Success() {
	req := &requests.UpdateTransferAmountRequest{TransferID: 1, TransferAmount: 150000}
	updatedTransfer := &record.TransferRecord{ID: 1, TransferAmount: 150000}

	suite.mockRepo.EXPECT().UpdateTransferAmount(req).Return(updatedTransfer, nil)

	result, err := suite.mockRepo.UpdateTransferAmount(req)

	suite.NoError(err)
	suite.Equal(int(150000), result.TransferAmount)
}

func (suite *TransferRepositorySuite) TestTrashedTransfer_Success() {
	transferID := 1
	trashedTime := "2024-01-01T00:00:00Z"
	trashedTransfer := &record.TransferRecord{ID: transferID, DeletedAt: &trashedTime}

	suite.mockRepo.EXPECT().TrashedTransfer(transferID).Return(trashedTransfer, nil)

	result, err := suite.mockRepo.TrashedTransfer(transferID)

	suite.NoError(err)
	suite.NotNil(result.DeletedAt)
	suite.Equal(trashedTransfer, result)
}

func (suite *TransferRepositorySuite) TestDeleteTransferPermanent_Success() {
	transferID := 1

	suite.mockRepo.EXPECT().DeleteTransferPermanent(transferID).Return(true, nil)

	result, err := suite.mockRepo.DeleteTransferPermanent(transferID)

	suite.NoError(err)
	suite.True(result)
}

func (suite *TransferRepositorySuite) TestRestoreAllTransfer_Success() {
	suite.mockRepo.EXPECT().RestoreAllTransfer().Return(true, nil)

	result, err := suite.mockRepo.RestoreAllTransfer()

	suite.NoError(err)
	suite.True(result)
}

func (suite *TransferRepositorySuite) TestGetMonthTransferStatusSuccess_Success() {
	req := &requests.MonthStatusTransfer{Year: 2024, Month: 1}
	monthlySuccess := []*record.TransferRecordMonthStatusSuccess{
		{Year: "2024", Month: "01", TotalSuccess: 50, TotalAmount: 5000000},
	}

	suite.mockRepo.EXPECT().GetMonthTransferStatusSuccess(req).Return(monthlySuccess, nil)

	result, err := suite.mockRepo.GetMonthTransferStatusSuccess(req)

	suite.NoError(err)
	suite.Equal(monthlySuccess, result)
}

func (suite *TransferRepositorySuite) TestGetYearlyTransferStatusFailed_Success() {
	year := 2023
	yearlyFailed := []*record.TransferRecordYearStatusFailed{
		{Year: "2023", TotalFailed: 2, TotalAmount: 100000},
	}

	suite.mockRepo.EXPECT().GetYearlyTransferStatusFailed(year).Return(yearlyFailed, nil)

	result, err := suite.mockRepo.GetYearlyTransferStatusFailed(year)

	suite.NoError(err)
	suite.Equal(yearlyFailed, result)
}

func (suite *TransferRepositorySuite) TestGetMonthTransferStatusSuccessByCardNumber_Success() {
	req := &requests.MonthStatusTransferCardNumber{CardNumber: "1111-xxxx", Year: 2024, Month: 1}
	monthlySuccess := []*record.TransferRecordMonthStatusSuccess{
		{Year: "2024", Month: "01", TotalSuccess: 5, TotalAmount: 500000},
	}

	suite.mockRepo.EXPECT().GetMonthTransferStatusSuccessByCardNumber(req).Return(monthlySuccess, nil)

	result, err := suite.mockRepo.GetMonthTransferStatusSuccessByCardNumber(req)

	suite.NoError(err)
	suite.Equal(monthlySuccess, result)
}

func (suite *TransferRepositorySuite) TestGetYearlyTransferStatusFailedByCardNumber_Success() {
	req := &requests.YearStatusTransferCardNumber{CardNumber: "2222-xxxx", Year: 2023}
	yearlyFailed := []*record.TransferRecordYearStatusFailed{
		{Year: "2023", TotalFailed: 1, TotalAmount: 50000},
	}

	suite.mockRepo.EXPECT().GetYearlyTransferStatusFailedByCardNumber(req).Return(yearlyFailed, nil)

	result, err := suite.mockRepo.GetYearlyTransferStatusFailedByCardNumber(req)

	suite.NoError(err)
	suite.Equal(yearlyFailed, result)
}

func (suite *TransferRepositorySuite) TestGetMonthlyTransferAmounts_Success() {
	year := 2024
	monthlyAmounts := []*record.TransferMonthAmount{
		{Month: "2024-01", TotalAmount: 4500000},
		{Month: "2024-02", TotalAmount: 5500000},
	}

	suite.mockRepo.EXPECT().GetMonthlyTransferAmounts(year).Return(monthlyAmounts, nil)

	result, err := suite.mockRepo.GetMonthlyTransferAmounts(year)

	suite.NoError(err)
	suite.Equal(monthlyAmounts, result)
}

func (suite *TransferRepositorySuite) TestGetYearlyTransferAmountsBySenderCardNumber_Success() {
	req := &requests.MonthYearCardNumber{CardNumber: "1111-xxxx", Year: 2023}
	yearlyAmounts := []*record.TransferYearAmount{
		{Year: "2023", TotalAmount: 12000000},
	}

	suite.mockRepo.EXPECT().GetYearlyTransferAmountsBySenderCardNumber(req).Return(yearlyAmounts, nil)

	result, err := suite.mockRepo.GetYearlyTransferAmountsBySenderCardNumber(req)

	suite.NoError(err)
	suite.Equal(yearlyAmounts, result)
}

func (suite *TransferRepositorySuite) TestGetMonthlyTransferAmountsByReceiverCardNumber_Success() {
	req := &requests.MonthYearCardNumber{CardNumber: "2222-xxxx", Year: 2024}
	monthlyAmounts := []*record.TransferMonthAmount{
		{Month: "2024-01", TotalAmount: 200000},
		{Month: "2024-02", TotalAmount: 300000},
	}

	suite.mockRepo.EXPECT().GetMonthlyTransferAmountsByReceiverCardNumber(req).Return(monthlyAmounts, nil)

	result, err := suite.mockRepo.GetMonthlyTransferAmountsByReceiverCardNumber(req)

	suite.NoError(err)
	suite.Equal(monthlyAmounts, result)
}

func TestTransferRepositorySuite(t *testing.T) {
	suite.Run(t, new(TransferRepositorySuite))
}
