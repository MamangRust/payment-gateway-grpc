package test

import (
	"errors"
	"testing"

	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	mock_responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/mocks"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type SaldoServiceSuite struct {
	suite.Suite
	mockCardRepo  *mocks.MockCardRepository
	mockSaldoRepo *mocks.MockSaldoRepository
	mockLogger    *mock_logger.MockLoggerInterface
	mockMapper    *mock_responseservice.MockSaldoResponseMapper
	mockCtrl      *gomock.Controller
}

func (suite *SaldoServiceSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockCardRepo = mocks.NewMockCardRepository(suite.mockCtrl)
	suite.mockSaldoRepo = mocks.NewMockSaldoRepository(suite.mockCtrl)
	suite.mockLogger = mock_logger.NewMockLoggerInterface(suite.mockCtrl)
	suite.mockMapper = mock_responseservice.NewMockSaldoResponseMapper(suite.mockCtrl)
}

func (suite *SaldoServiceSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *SaldoServiceSuite) TestFindAll_Success() {
	req := &requests.FindAllSaldos{Search: "1234", Page: 1, PageSize: 10}
	total := 1
	mockSaldos := []*record.SaldoRecord{{ID: 1, CardNumber: "1111-xxxx-xxxx-1111"}}
	expectedResponse := []*response.SaldoResponse{{ID: 1, CardNumber: "1111-xxxx-xxxx-1111"}}

	suite.mockLogger.EXPECT().Debug("Fetching saldo", gomock.Any())
	suite.mockSaldoRepo.EXPECT().FindAllSaldos(req).Return(mockSaldos, &total, nil)
	suite.mockMapper.EXPECT().ToSaldoResponses(mockSaldos).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched saldo", gomock.Any())

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *SaldoServiceSuite) TestFindAll_Failure() {
	req := &requests.FindAllSaldos{Search: "1234", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching saldo", gomock.Any())
	suite.mockSaldoRepo.EXPECT().FindAllSaldos(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch saldo", gomock.Any())

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(saldo_errors.ErrFailedFindAllSaldos, err)
}

func (suite *SaldoServiceSuite) TestFindById_Success() {
	saldoID := 1
	mockSaldo := &record.SaldoRecord{ID: saldoID, CardNumber: "1111-xxxx-xxxx-1111"}
	expectedResponse := &response.SaldoResponse{ID: saldoID, CardNumber: "1111-xxxx-xxxx-1111"}

	suite.mockLogger.EXPECT().Debug("Fetching saldo record by ID", zap.Int("saldo_id", saldoID))
	suite.mockSaldoRepo.EXPECT().FindById(saldoID).Return(mockSaldo, nil)
	suite.mockMapper.EXPECT().ToSaldoResponse(mockSaldo).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched saldo", zap.Int("saldo_id", saldoID))

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(saldoID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *SaldoServiceSuite) TestFindById_NotFound() {
	saldoID := 99
	repoError := errors.New("saldo not found")

	suite.mockLogger.EXPECT().Debug("Fetching saldo record by ID", zap.Int("saldo_id", saldoID))
	suite.mockSaldoRepo.EXPECT().FindById(saldoID).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve saldo details", gomock.Any())

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(saldoID)

	suite.Nil(result)
	suite.Equal(saldo_errors.ErrFailedSaldoNotFound, err)
}

func (suite *SaldoServiceSuite) TestFindByCardNumber_Success() {
	cardNumber := "1234-5678-9012-3456"
	mockSaldo := &record.SaldoRecord{CardNumber: cardNumber, TotalBalance: 1000000}
	expectedResponse := &response.SaldoResponse{CardNumber: cardNumber, TotalBalance: 1000000}

	suite.mockLogger.EXPECT().Debug("Fetching saldo record by card number", zap.String("card_number", cardNumber))
	suite.mockSaldoRepo.EXPECT().FindByCardNumber(cardNumber).Return(mockSaldo, nil)
	suite.mockMapper.EXPECT().ToSaldoResponse(mockSaldo).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched saldo by card number", zap.String("card_number", cardNumber))

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindByCardNumber(cardNumber)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *SaldoServiceSuite) TestFindByCardNumber_NotFound() {
	cardNumber := "not-found-xxxx"
	repoError := errors.New("saldo not found")

	suite.mockLogger.EXPECT().Debug("Fetching saldo record by card number", zap.String("card_number", cardNumber))
	suite.mockSaldoRepo.EXPECT().FindByCardNumber(cardNumber).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve saldo details", gomock.Any())

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindByCardNumber(cardNumber)

	suite.Nil(result)
	suite.Equal(saldo_errors.ErrFailedSaldoNotFound, err)
}

func (suite *SaldoServiceSuite) TestFindByActive_Success() {
	req := &requests.FindAllSaldos{Search: "active", Page: 1, PageSize: 10}
	total := 1
	mockSaldos := []*record.SaldoRecord{{ID: 1, DeletedAt: nil}}
	expectedResponse := []*response.SaldoResponseDeleteAt{{ID: 1, DeletedAt: nil}}

	suite.mockLogger.EXPECT().Debug("Fetching active saldo", gomock.Any())
	suite.mockSaldoRepo.EXPECT().FindByActive(req).Return(mockSaldos, &total, nil)
	suite.mockMapper.EXPECT().ToSaldoResponsesDeleteAt(mockSaldos).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched active saldo", gomock.Any())

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *SaldoServiceSuite) TestFindByActive_Failure() {
	req := &requests.FindAllSaldos{Search: "active", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching active saldo", gomock.Any())
	suite.mockSaldoRepo.EXPECT().FindByActive(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve active saldo", gomock.Any())

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(saldo_errors.ErrFailedFindActiveSaldos, err)
}

func (suite *SaldoServiceSuite) TestFindByTrashed_Success() {
	req := &requests.FindAllSaldos{Page: 1, PageSize: 10}
	total := 1
	trashedTime := "2024-01-01T00:00:00Z"
	mockSaldos := []*record.SaldoRecord{{ID: 2, DeletedAt: &trashedTime}}
	expectedResponse := []*response.SaldoResponseDeleteAt{{ID: 2, DeletedAt: &trashedTime}}

	suite.mockLogger.EXPECT().Debug("Fetching saldo record", gomock.Any())
	suite.mockSaldoRepo.EXPECT().FindByTrashed(req).Return(mockSaldos, &total, nil)
	suite.mockMapper.EXPECT().ToSaldoResponsesDeleteAt(mockSaldos).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched trashed saldo", gomock.Any())

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *SaldoServiceSuite) TestFindByTrashed_Failure() {
	req := &requests.FindAllSaldos{Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching saldo record", gomock.Any())
	suite.mockSaldoRepo.EXPECT().FindByTrashed(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve trashed saldo", gomock.Any())

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(saldo_errors.ErrFailedFindTrashedSaldos, err)
}

func (suite *SaldoServiceSuite) TestTrashedSaldo_Success() {
	saldoID := 1
	trashedSaldo := &record.SaldoRecord{ID: saldoID}
	expectedResponse := &response.SaldoResponseDeleteAt{ID: saldoID}

	suite.mockLogger.EXPECT().Debug("Trashing saldo record", zap.Int("saldo_id", saldoID))
	suite.mockSaldoRepo.EXPECT().TrashedSaldo(saldoID).Return(trashedSaldo, nil)
	suite.mockMapper.EXPECT().ToSaldoResponseDeleteAt(trashedSaldo).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully trashed saldo", zap.Int("saldo_id", saldoID))

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashSaldo(saldoID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *SaldoServiceSuite) TestTrashedSaldo_Failure() {
	saldoID := 1
	dbError := errors.New("failed to trash")

	suite.mockLogger.EXPECT().Debug("Trashing saldo record", zap.Int("saldo_id", saldoID))
	suite.mockLogger.EXPECT().Error("Failed to move saldo to trash", gomock.Any())
	suite.mockSaldoRepo.EXPECT().TrashedSaldo(saldoID).Return(nil, dbError)

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashSaldo(saldoID)

	suite.Nil(result)
	suite.Equal(saldo_errors.ErrFailedTrashSaldo, err)
}

func (suite *SaldoServiceSuite) TestRestoreSaldo_Success() {
	saldoID := 1
	restoredSaldo := &record.SaldoRecord{ID: saldoID, DeletedAt: nil}
	expectedResponse := &response.SaldoResponseDeleteAt{ID: saldoID, DeletedAt: nil}

	suite.mockLogger.EXPECT().Debug("Restoring saldo record from trash", zap.Int("saldo_id", saldoID))
	suite.mockSaldoRepo.EXPECT().RestoreSaldo(saldoID).Return(restoredSaldo, nil)
	suite.mockMapper.EXPECT().ToSaldoResponseDeleteAt(restoredSaldo).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully restored saldo", zap.Int("saldo_id", saldoID))

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreSaldo(saldoID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *SaldoServiceSuite) TestRestoreSaldo_Failure() {
	saldoID := 1
	dbError := errors.New("failed to restore")

	suite.mockLogger.EXPECT().Debug("Restoring saldo record from trash", zap.Int("saldo_id", saldoID))
	suite.mockLogger.EXPECT().Error("Failed to restore saldo from trash", gomock.Any())
	suite.mockSaldoRepo.EXPECT().RestoreSaldo(saldoID).Return(nil, dbError)

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreSaldo(saldoID)

	suite.Nil(result)
	suite.Equal(saldo_errors.ErrFailedRestoreSaldo, err)
}

func (suite *SaldoServiceSuite) TestDeleteSaldoPermanent_Success() {
	saldoID := 1

	suite.mockLogger.EXPECT().Debug("Deleting saldo permanently", zap.Int("saldo_id", saldoID))
	suite.mockSaldoRepo.EXPECT().DeleteSaldoPermanent(saldoID).Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted saldo permanently", zap.Int("saldo_id", saldoID))

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteSaldoPermanent(saldoID)

	suite.Nil(err)
	suite.True(result)
}

func (suite *SaldoServiceSuite) TestDeleteSaldoPermanent_Failure() {
	saldoID := 1
	dbError := errors.New("failed to delete permanently")

	suite.mockLogger.EXPECT().Debug("Deleting saldo permanently", zap.Int("saldo_id", saldoID))
	suite.mockLogger.EXPECT().Error("Failed to permanently delete saldo", gomock.Any())
	suite.mockSaldoRepo.EXPECT().DeleteSaldoPermanent(saldoID).Return(false, dbError)

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteSaldoPermanent(saldoID)

	suite.False(result)
	suite.Equal(saldo_errors.ErrFailedDeleteSaldoPermanent, err)
}

func (suite *SaldoServiceSuite) TestRestoreAllSaldo_Success() {
	suite.mockLogger.EXPECT().Debug("Restoring all saldo")
	suite.mockSaldoRepo.EXPECT().RestoreAllSaldo().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully restored all saldo")

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllSaldo()

	suite.Nil(err)
	suite.True(result)
}

func (suite *SaldoServiceSuite) TestRestoreAllSaldo_Failure() {
	dbError := errors.New("failed to restore all")

	suite.mockLogger.EXPECT().Debug("Restoring all saldo")
	suite.mockLogger.EXPECT().Error("Failed to restore all saldo", gomock.Any())
	suite.mockSaldoRepo.EXPECT().RestoreAllSaldo().Return(false, dbError)

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllSaldo()

	suite.False(result)
	suite.Equal(saldo_errors.ErrFailedRestoreAllSaldo, err)
}

func (suite *SaldoServiceSuite) TestDeleteAllSaldoPermanent_Success() {
	suite.mockLogger.EXPECT().Debug("Permanently deleting all saldo")
	suite.mockSaldoRepo.EXPECT().DeleteAllSaldoPermanent().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted all saldo permanently")

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllSaldoPermanent()

	suite.Nil(err)
	suite.True(result)
}

func (suite *SaldoServiceSuite) TestDeleteAllSaldoPermanent_Failure() {
	dbError := errors.New("failed to delete all permanently")

	suite.mockLogger.EXPECT().Debug("Permanently deleting all saldo")
	suite.mockLogger.EXPECT().Error("Failed to permanently delete all saldo", gomock.Any())
	suite.mockSaldoRepo.EXPECT().DeleteAllSaldoPermanent().Return(false, dbError)

	svc := service.NewSaldoService(suite.mockSaldoRepo, suite.mockCardRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllSaldoPermanent()

	suite.False(result)
	suite.Equal(saldo_errors.ErrFailedDeleteAllSaldoPermanent, err)
}

func TestSaldoServiceSuite(t *testing.T) {
	suite.Run(t, new(SaldoServiceSuite))
}
