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
	"MamangRust/paymentgatewaygrpc/pkg/errors/topup_errors"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type TopupServiceSuite struct {
	suite.Suite
	mockTopupRepo *mocks.MockTopupRepository
	mockCardRepo  *mocks.MockCardRepository
	mockSaldoRepo *mocks.MockSaldoRepository
	mockLogger    *mock_logger.MockLoggerInterface
	mockMapper    *mock_responseservice.MockTopupResponseMapper
	mockCtrl      *gomock.Controller
}

func (suite *TopupServiceSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockTopupRepo = mocks.NewMockTopupRepository(suite.mockCtrl)
	suite.mockCardRepo = mocks.NewMockCardRepository(suite.mockCtrl)
	suite.mockSaldoRepo = mocks.NewMockSaldoRepository(suite.mockCtrl)
	suite.mockLogger = mock_logger.NewMockLoggerInterface(suite.mockCtrl)
	suite.mockMapper = mock_responseservice.NewMockTopupResponseMapper(suite.mockCtrl)
}

func (suite *TopupServiceSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TopupServiceSuite) TestFindAll_Success() {
	req := &requests.FindAllTopups{Search: "success", Page: 1, PageSize: 10}
	total := 1
	mockTopups := []*record.TopupRecord{{ID: 1, TopupNo: "TP001"}}
	expectedResponse := []*response.TopupResponse{{ID: 1, TopupNo: "TP001"}}

	suite.mockLogger.EXPECT().Debug("Fetching topup", gomock.Any())
	suite.mockTopupRepo.EXPECT().FindAllTopups(req).Return(mockTopups, &total, nil)
	suite.mockMapper.EXPECT().ToTopupResponses(mockTopups).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched topup", gomock.Any())

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *TopupServiceSuite) TestFindAll_Failure() {
	req := &requests.FindAllTopups{Search: "success", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching topup", gomock.Any())
	suite.mockTopupRepo.EXPECT().FindAllTopups(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch topup", gomock.Any())

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(topup_errors.ErrFailedFindAllTopups, err)
}

func (suite *TopupServiceSuite) TestFindAllByCardNumber_Success() {
	req := &requests.FindAllTopupsByCardNumber{CardNumber: "1111-xxxx", Page: 1, PageSize: 10}
	total := 1
	mockTopups := []*record.TopupRecord{{ID: 2, CardNumber: "1111-xxxx"}}
	expectedResponse := []*response.TopupResponse{{ID: 2, CardNumber: "1111-xxxx"}}

	suite.mockLogger.EXPECT().Debug("Fetching topup by card number", gomock.Any())
	suite.mockTopupRepo.EXPECT().FindAllTopupByCardNumber(req).Return(mockTopups, &total, nil)
	suite.mockMapper.EXPECT().ToTopupResponses(mockTopups).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched topup", gomock.Any())

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAllByCardNumber(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *TopupServiceSuite) TestFindAllByCardNumber_Failure() {
	req := &requests.FindAllTopupsByCardNumber{CardNumber: "1111-xxxx", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching topup by card number", gomock.Any())
	suite.mockTopupRepo.EXPECT().FindAllTopupByCardNumber(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch topup by card number", gomock.Any())

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAllByCardNumber(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(topup_errors.ErrFailedFindAllTopupsByCardNumber, err)
}

func (suite *TopupServiceSuite) TestFindById_Success() {
	topupID := 1
	mockTopup := &record.TopupRecord{ID: topupID, TopupNo: "TP001"}
	expectedResponse := &response.TopupResponse{ID: topupID, TopupNo: "TP001"}

	suite.mockLogger.EXPECT().Debug("Fetching topup by ID", zap.Int("topup_id", topupID))
	suite.mockTopupRepo.EXPECT().FindById(topupID).Return(mockTopup, nil)
	suite.mockMapper.EXPECT().ToTopupResponse(mockTopup).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched topup", zap.Int("topup_id", topupID))

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(topupID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *TopupServiceSuite) TestFindById_NotFound() {
	topupID := 99
	repoError := errors.New("topup not found")

	suite.mockLogger.EXPECT().Debug("Fetching topup by ID", zap.Int("topup_id", topupID))
	suite.mockTopupRepo.EXPECT().FindById(topupID).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("failed to find topup by id", gomock.Any())

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(topupID)

	suite.Nil(result)
	suite.Equal(topup_errors.ErrTopupNotFoundRes, err)
}

func (suite *TopupServiceSuite) TestFindByActive_Success() {
	req := &requests.FindAllTopups{Search: "active", Page: 1, PageSize: 10}
	total := 1
	mockTopups := []*record.TopupRecord{{ID: 1, TopupNo: "TP001", DeletedAt: nil}}
	expectedResponse := []*response.TopupResponseDeleteAt{{ID: 1, TopupNo: "TP001", DeletedAt: nil}}

	suite.mockLogger.EXPECT().Debug("Fetching active topup", gomock.Any())
	suite.mockTopupRepo.EXPECT().FindByActive(req).Return(mockTopups, &total, nil)
	suite.mockMapper.EXPECT().ToTopupResponsesDeleteAt(mockTopups).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched active topup", gomock.Any())

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *TopupServiceSuite) TestFindByActive_Failure() {
	req := &requests.FindAllTopups{Search: "active", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching active topup", gomock.Any())
	suite.mockTopupRepo.EXPECT().FindByActive(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch active topup", gomock.Any())

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(topup_errors.ErrFailedFindActiveTopups, err)
}

func (suite *TopupServiceSuite) TestFindByTrashed_Success() {
	req := &requests.FindAllTopups{Page: 1, PageSize: 10}
	total := 1
	trashedTime := "2024-01-01T00:00:00Z"
	mockTopups := []*record.TopupRecord{{ID: 2, TopupNo: "TP002", DeletedAt: &trashedTime}}
	expectedResponse := []*response.TopupResponseDeleteAt{{ID: 2, TopupNo: "TP002", DeletedAt: &trashedTime}}

	suite.mockLogger.EXPECT().Debug("Fetching trashed topup", gomock.Any())
	suite.mockTopupRepo.EXPECT().FindByTrashed(req).Return(mockTopups, &total, nil)
	suite.mockMapper.EXPECT().ToTopupResponsesDeleteAt(mockTopups).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched trashed topup", gomock.Any())

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *TopupServiceSuite) TestFindByTrashed_Failure() {
	req := &requests.FindAllTopups{Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching trashed topup", gomock.Any())
	suite.mockTopupRepo.EXPECT().FindByTrashed(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch trashed topup", gomock.Any())

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(topup_errors.ErrFailedFindTrashedTopups, err)
}

func (suite *TopupServiceSuite) TestTrashedTopup_Success() {
	topupID := 1
	trashedTopup := &record.TopupRecord{ID: topupID}
	expectedResponse := &response.TopupResponseDeleteAt{ID: topupID}

	suite.mockLogger.EXPECT().Debug("Starting TrashedTopup process", zap.Int("topup_id", topupID))
	suite.mockTopupRepo.EXPECT().TrashedTopup(topupID).Return(trashedTopup, nil)
	suite.mockMapper.EXPECT().ToTopupResponseDeleteAt(trashedTopup).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("TrashedTopup process completed", zap.Int("topup_id", topupID))

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashedTopup(topupID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *TopupServiceSuite) TestTrashedTopup_Failure() {
	topupID := 1
	dbError := errors.New("failed to trash")

	suite.mockLogger.EXPECT().Debug("Starting TrashedTopup process", zap.Int("topup_id", topupID))
	suite.mockLogger.EXPECT().Error("Failed to move topup to trash", gomock.Any())
	suite.mockTopupRepo.EXPECT().TrashedTopup(topupID).Return(nil, dbError)

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashedTopup(topupID)

	suite.Nil(result)
	suite.Equal(topup_errors.ErrFailedTrashTopup, err)
}

func (suite *TopupServiceSuite) TestRestoreTopup_Success() {
	topupID := 1
	restoredTopup := &record.TopupRecord{ID: topupID, DeletedAt: nil}
	expectedResponse := &response.TopupResponseDeleteAt{ID: topupID, DeletedAt: nil}

	suite.mockLogger.EXPECT().Debug("Starting RestoreTopup process", zap.Int("topup_id", topupID))
	suite.mockTopupRepo.EXPECT().RestoreTopup(topupID).Return(restoredTopup, nil)
	suite.mockMapper.EXPECT().ToTopupResponseDeleteAt(restoredTopup).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("RestoreTopup process completed", zap.Int("topup_id", topupID))

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreTopup(topupID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *TopupServiceSuite) TestRestoreTopup_Failure() {
	topupID := 1
	dbError := errors.New("failed to restore")

	suite.mockLogger.EXPECT().Debug("Starting RestoreTopup process", zap.Int("topup_id", topupID))
	suite.mockLogger.EXPECT().Error("Failed to restore topup from trash", gomock.Any())
	suite.mockTopupRepo.EXPECT().RestoreTopup(topupID).Return(nil, dbError)

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreTopup(topupID)

	suite.Nil(result)
	suite.Equal(topup_errors.ErrFailedRestoreTopup, err)
}

func (suite *TopupServiceSuite) TestDeleteTopupPermanent_Success() {
	topupID := 1

	suite.mockLogger.EXPECT().Debug("Starting DeleteTopupPermanent process", zap.Int("topup_id", topupID))
	suite.mockTopupRepo.EXPECT().DeleteTopupPermanent(topupID).Return(true, nil)
	suite.mockLogger.EXPECT().Debug("DeleteTopupPermanent process completed", zap.Int("topup_id", topupID))

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteTopupPermanent(topupID)

	suite.Nil(err)
	suite.True(result)
}

func (suite *TopupServiceSuite) TestDeleteTopupPermanent_Failure() {
	topupID := 1
	dbError := errors.New("failed to delete permanently")

	suite.mockLogger.EXPECT().Debug("Starting DeleteTopupPermanent process", zap.Int("topup_id", topupID))
	suite.mockLogger.EXPECT().Error("Failed to delete topup permanently", gomock.Any())
	suite.mockTopupRepo.EXPECT().DeleteTopupPermanent(topupID).Return(false, dbError)

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteTopupPermanent(topupID)

	suite.False(result)
	suite.Equal(topup_errors.ErrFailedDeleteTopup, err)
}

func (suite *TopupServiceSuite) TestRestoreAllTopup_Success() {
	suite.mockLogger.EXPECT().Debug("Restoring all topups")
	suite.mockTopupRepo.EXPECT().RestoreAllTopup().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully restored all topups")

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllTopup()

	suite.Nil(err)
	suite.True(result)
}

func (suite *TopupServiceSuite) TestRestoreAllTopup_Failure() {
	dbError := errors.New("failed to restore all")

	suite.mockLogger.EXPECT().Debug("Restoring all topups")
	suite.mockLogger.EXPECT().Error("Failed to restore all topups", gomock.Any())
	suite.mockTopupRepo.EXPECT().RestoreAllTopup().Return(false, dbError)

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllTopup()

	suite.False(result)
	suite.Equal(topup_errors.ErrFailedRestoreAllTopups, err)
}

func (suite *TopupServiceSuite) TestDeleteAllTopupPermanent_Success() {
	suite.mockLogger.EXPECT().Debug("Permanently deleting all topups")
	suite.mockTopupRepo.EXPECT().DeleteAllTopupPermanent().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted all topups permanently")

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllTopupPermanent()

	suite.Nil(err)
	suite.True(result)
}

func (suite *TopupServiceSuite) TestDeleteAllTopupPermanent_Failure() {
	dbError := errors.New("failed to delete all permanently")

	suite.mockLogger.EXPECT().Debug("Permanently deleting all topups")
	suite.mockLogger.EXPECT().Error("Failed to permanently delete all topups", gomock.Any())
	suite.mockTopupRepo.EXPECT().DeleteAllTopupPermanent().Return(false, dbError)

	svc := service.NewTopupService(suite.mockCardRepo, suite.mockTopupRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllTopupPermanent()

	suite.False(result)
	suite.Equal(topup_errors.ErrFailedDeleteAllTopups, err)
}

func TestTopupServiceSuite(t *testing.T) {
	suite.Run(t, new(TopupServiceSuite))
}
