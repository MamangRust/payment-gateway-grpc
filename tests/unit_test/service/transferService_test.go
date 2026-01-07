package test

import (
	"errors"
	"testing"

	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	mock_responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/mocks"
	mock_repository "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors/transfer_errors"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type TransferServiceSuite struct {
	suite.Suite
	mockTransferRepo *mock_repository.MockTransferRepository
	mockCardRepo     *mock_repository.MockCardRepository
	mockUserRepo     *mock_repository.MockUserRepository
	mockSaldoRepo    *mock_repository.MockSaldoRepository
	mockLogger       *mock_logger.MockLoggerInterface
	mockMapper       *mock_responseservice.MockTransferResponseMapper
	mockCtrl         *gomock.Controller
}

func (suite *TransferServiceSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockCardRepo = mock_repository.NewMockCardRepository(suite.mockCtrl)
	suite.mockUserRepo = mock_repository.NewMockUserRepository(suite.mockCtrl)
	suite.mockSaldoRepo = mock_repository.NewMockSaldoRepository(suite.mockCtrl)
	suite.mockTransferRepo = mock_repository.NewMockTransferRepository(suite.mockCtrl)
	suite.mockLogger = mock_logger.NewMockLoggerInterface(suite.mockCtrl)
	suite.mockMapper = mock_responseservice.NewMockTransferResponseMapper(suite.mockCtrl)
}

func (suite *TransferServiceSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TransferServiceSuite) TestFindAll_Success() {
	req := &requests.FindAllTranfers{Search: "success", Page: 1, PageSize: 10}
	total := 1
	mockTransfers := []*record.TransferRecord{{ID: 1, TransferNo: "TF001"}}
	expectedResponse := []*response.TransferResponse{{ID: 1, TransferNo: "TF001"}}

	suite.mockLogger.EXPECT().Debug("Fetching transfer", gomock.Any())
	suite.mockTransferRepo.EXPECT().FindAll(req).Return(mockTransfers, &total, nil)
	suite.mockMapper.EXPECT().ToTransfersResponse(mockTransfers).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched transfer", gomock.Any())

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)

	result, totalRes, err := svc.FindAll(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *TransferServiceSuite) TestFindAll_Failure() {
	req := &requests.FindAllTranfers{Search: "success", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching transfer", gomock.Any())
	suite.mockTransferRepo.EXPECT().FindAll(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch transfer", gomock.Any())

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(transfer_errors.ErrFailedFindAllTransfers, err)
}

func (suite *TransferServiceSuite) TestFindById_Success() {
	transferID := 1
	mockTransfer := &record.TransferRecord{ID: transferID, TransferNo: "TF001"}
	expectedResponse := &response.TransferResponse{ID: transferID, TransferNo: "TF001"}

	suite.mockLogger.EXPECT().Debug("Fetching transfer by ID", zap.Int("transfer_id", transferID))
	suite.mockTransferRepo.EXPECT().FindById(transferID).Return(mockTransfer, nil)
	suite.mockMapper.EXPECT().ToTransferResponse(mockTransfer).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched transfer", zap.Int("transfer_id", transferID))

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(transferID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *TransferServiceSuite) TestFindById_NotFound() {
	transferID := 99
	repoError := errors.New("transfer not found")

	suite.mockLogger.EXPECT().Debug("Fetching transfer by ID", zap.Int("transfer_id", transferID))
	suite.mockTransferRepo.EXPECT().FindById(transferID).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("failed to find transfer by ID", gomock.Any())

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(transferID)

	suite.Nil(result)
	suite.Equal(transfer_errors.ErrTransferNotFound, err)
}

func (suite *TransferServiceSuite) TestFindByActive_Success() {
	req := &requests.FindAllTranfers{Search: "active", Page: 1, PageSize: 10}
	total := 1
	mockTransfers := []*record.TransferRecord{{ID: 1, TransferNo: "TF001", DeletedAt: nil}}
	expectedResponse := []*response.TransferResponseDeleteAt{{ID: 1, TransferNo: "TF001", DeletedAt: nil}}

	suite.mockLogger.EXPECT().Debug("Fetching active transfer", gomock.Any())
	suite.mockTransferRepo.EXPECT().FindByActive(req).Return(mockTransfers, &total, nil)
	suite.mockMapper.EXPECT().ToTransfersResponseDeleteAt(mockTransfers).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched active transfer", gomock.Any())

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *TransferServiceSuite) TestFindByActive_Failure() {
	req := &requests.FindAllTranfers{Search: "active", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching active transfer", gomock.Any())
	suite.mockTransferRepo.EXPECT().FindByActive(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch active transfer", gomock.Any())

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(transfer_errors.ErrFailedFindActiveTransfers, err)
}

func (suite *TransferServiceSuite) TestFindByTrashed_Success() {
	req := &requests.FindAllTranfers{Page: 1, PageSize: 10}
	total := 1
	trashedTime := "2024-01-01T00:00:00Z"
	mockTransfers := []*record.TransferRecord{{ID: 2, TransferNo: "TF002", DeletedAt: &trashedTime}}
	expectedResponse := []*response.TransferResponseDeleteAt{{ID: 2, TransferNo: "TF002", DeletedAt: &trashedTime}}

	suite.mockLogger.EXPECT().Debug("Fetching trashed transfer", gomock.Any())
	suite.mockTransferRepo.EXPECT().FindByTrashed(req).Return(mockTransfers, &total, nil)
	suite.mockMapper.EXPECT().ToTransfersResponseDeleteAt(mockTransfers).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched trashed transfer", gomock.Any())

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *TransferServiceSuite) TestFindByTrashed_Failure() {
	req := &requests.FindAllTranfers{Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching trashed transfer", gomock.Any())
	suite.mockTransferRepo.EXPECT().FindByTrashed(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch trashed transfer", gomock.Any())

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(transfer_errors.ErrFailedFindTrashedTransfers, err)
}

func (suite *TransferServiceSuite) TestFindTransferByTransferFrom_Success() {
	transferFrom := "1111-xxxx-xxxx-1111"
	mockTransfers := []*record.TransferRecord{{ID: 1, TransferFrom: transferFrom}}
	expectedResponse := []*response.TransferResponse{{ID: 1, TransferFrom: transferFrom}}

	suite.mockLogger.EXPECT().Debug("Starting fetch transfer by transfer_from", gomock.Any())
	suite.mockTransferRepo.EXPECT().FindTransferByTransferFrom(transferFrom).Return(mockTransfers, nil)

	suite.mockMapper.EXPECT().ToTransfersResponse(mockTransfers).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched transfer record by transfer_from", gomock.Any())

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindTransferByTransferFrom(transferFrom)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *TransferServiceSuite) TestFindTransferByTransferFrom_Failure() {
	transferFrom := "not-found-xxxx"
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Starting fetch transfer by transfer_from", gomock.Any())
	suite.mockTransferRepo.EXPECT().FindTransferByTransferFrom(transferFrom).Return(nil, dbError)

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	suite.mockLogger.EXPECT().Error("Failed to fetch transfers by transfer_from", gomock.Any())

	result, err := svc.FindTransferByTransferFrom(transferFrom)

	suite.Nil(result)
	suite.Equal(transfer_errors.ErrTransferNotFound, err)
}

func (suite *TransferServiceSuite) TestFindTransferByTransferTo_Success() {
	transferTo := "2222-xxxx-xxxx-2222"
	mockTransfers := []*record.TransferRecord{{ID: 2, TransferTo: transferTo}}
	expectedResponse := []*response.TransferResponse{{ID: 2, TransferTo: transferTo}}

	suite.mockLogger.EXPECT().Debug("Starting fetch transfer by transfer_to", gomock.Any())
	suite.mockTransferRepo.EXPECT().FindTransferByTransferTo(transferTo).Return(mockTransfers, nil)
	suite.mockMapper.EXPECT().ToTransfersResponse(mockTransfers).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched transfer record by transfer_to", gomock.Any())

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindTransferByTransferTo(transferTo)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *TransferServiceSuite) TestFindTransferByTransferTo_Failure() {
	transferTo := "not-found-xxxx"
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Starting fetch transfer by transfer_to", gomock.Any())
	suite.mockTransferRepo.EXPECT().FindTransferByTransferTo(transferTo).Return(nil, dbError)

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	suite.mockLogger.EXPECT().Error("Failed to fetch transfers by transfer_to", gomock.Any())
	result, err := svc.FindTransferByTransferTo(transferTo)

	suite.Nil(result)
	suite.Equal(transfer_errors.ErrTransferNotFound, err)
}

func (suite *TransferServiceSuite) TestTrashedTransfer_Success() {
	transferID := 1
	trashedTransfer := &record.TransferRecord{ID: transferID}
	expectedResponse := &response.TransferResponseDeleteAt{ID: transferID}

	suite.mockLogger.EXPECT().Debug("Starting trashed transfer process", zap.Int("transfer_id", transferID))
	suite.mockTransferRepo.EXPECT().TrashedTransfer(transferID).Return(trashedTransfer, nil)
	suite.mockMapper.EXPECT().ToTransferResponseDeleteAt(trashedTransfer).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("successfully trashed transfer", zap.Int("transfer_id", transferID))

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashedTransfer(transferID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *TransferServiceSuite) TestTrashedTransfer_Failure() {
	transferID := 1
	dbError := errors.New("failed to trash")

	suite.mockLogger.EXPECT().Debug("Starting trashed transfer process", zap.Int("transfer_id", transferID))
	suite.mockLogger.EXPECT().Error("Failed to trash transfer", gomock.Any())
	suite.mockTransferRepo.EXPECT().TrashedTransfer(transferID).Return(nil, dbError)

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashedTransfer(transferID)

	suite.Nil(result)
	suite.Equal(transfer_errors.ErrFailedTrashedTransfer, err)
}

func (suite *TransferServiceSuite) TestRestoreTransfer_Success() {
	transferID := 1
	restoredTransfer := &record.TransferRecord{ID: transferID, DeletedAt: nil}
	expectedResponse := &response.TransferResponseDeleteAt{ID: transferID, DeletedAt: nil}

	suite.mockLogger.EXPECT().Debug("Starting restore transfer process", zap.Int("transfer_id", transferID))
	suite.mockTransferRepo.EXPECT().RestoreTransfer(transferID).Return(restoredTransfer, nil)
	suite.mockMapper.EXPECT().ToTransferResponseDeleteAt(restoredTransfer).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("successfully restore transfer", zap.Int("transfer_id", transferID))

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreTransfer(transferID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *TransferServiceSuite) TestRestoreTransfer_Failure() {
	transferID := 1
	dbError := errors.New("failed to restore")

	suite.mockLogger.EXPECT().Debug("Starting restore transfer process", zap.Int("transfer_id", transferID))
	suite.mockLogger.EXPECT().Error("Failed to restore transfer", gomock.Any())
	suite.mockTransferRepo.EXPECT().RestoreTransfer(transferID).Return(nil, dbError)

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreTransfer(transferID)

	suite.Nil(result)
	suite.Equal(transfer_errors.ErrFailedRestoreTransfer, err)
}

func (suite *TransferServiceSuite) TestDeleteTransferPermanent_Success() {
	transferID := 1

	suite.mockLogger.EXPECT().Debug("Starting delete transfer permanent process", zap.Int("transfer_id", transferID))
	suite.mockTransferRepo.EXPECT().DeleteTransferPermanent(transferID).Return(true, nil)
	suite.mockLogger.EXPECT().Debug("successfully delete permanent transfer", zap.Int("transfer_id", transferID))

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteTransferPermanent(transferID)

	suite.Nil(err)
	suite.True(result)
}

func (suite *TransferServiceSuite) TestDeleteTransferPermanent_Failure() {
	transferID := 1
	dbError := errors.New("failed to delete permanently")

	suite.mockLogger.EXPECT().Debug("Starting delete transfer permanent process", zap.Int("transfer_id", transferID))
	suite.mockLogger.EXPECT().Error("Failed to permanently delete transfer", gomock.Any())
	suite.mockTransferRepo.EXPECT().DeleteTransferPermanent(transferID).Return(false, dbError)

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteTransferPermanent(transferID)

	suite.False(result)
	suite.Equal(transfer_errors.ErrFailedDeleteTransferPermanent, err)
}

func (suite *TransferServiceSuite) TestRestoreAllTransfer_Success() {
	suite.mockLogger.EXPECT().Debug("Restoring all transfers")
	suite.mockTransferRepo.EXPECT().RestoreAllTransfer().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully restored all transfers")

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllTransfer()

	suite.Nil(err)
	suite.True(result)
}

func (suite *TransferServiceSuite) TestRestoreAllTransfer_Failure() {
	dbError := errors.New("failed to restore all")

	suite.mockLogger.EXPECT().Debug("Restoring all transfers")
	suite.mockLogger.EXPECT().Error("Failed to restore all transfers", gomock.Any())
	suite.mockTransferRepo.EXPECT().RestoreAllTransfer().Return(false, dbError)

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllTransfer()

	suite.False(result)
	suite.Equal(transfer_errors.ErrFailedRestoreAllTransfers, err)
}

func (suite *TransferServiceSuite) TestDeleteAllTransferPermanent_Success() {
	suite.mockLogger.EXPECT().Debug("Permanently deleting all transfers")
	suite.mockTransferRepo.EXPECT().DeleteAllTransferPermanent().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted all transfers permanently")

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllTransferPermanent()

	suite.Nil(err)
	suite.True(result)
}

func (suite *TransferServiceSuite) TestDeleteAllTransferPermanent_Failure() {
	dbError := errors.New("failed to delete all permanently")

	suite.mockLogger.EXPECT().Debug("Permanently deleting all transfers")
	suite.mockLogger.EXPECT().Error("Failed to permanently delete all transfers", gomock.Any())
	suite.mockTransferRepo.EXPECT().DeleteAllTransferPermanent().Return(false, dbError)

	svc := service.NewTransferService(suite.mockUserRepo, suite.mockCardRepo, suite.mockTransferRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllTransferPermanent()

	suite.False(result)
	suite.Equal(transfer_errors.ErrFailedDeleteAllTransfersPermanent, err)
}

func TestTransferServiceSuite(t *testing.T) {
	suite.Run(t, new(TransferServiceSuite))
}
