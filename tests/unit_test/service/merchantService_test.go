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
	"MamangRust/paymentgatewaygrpc/pkg/errors/merchant_errors"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type MerchantServiceSuite struct {
	suite.Suite
	mockMerchantRepo *mocks.MockMerchantRepository
	mockLogger       *mock_logger.MockLoggerInterface
	mockMapper       *mock_responseservice.MockMerchantResponseMapper
	mockCtrl         *gomock.Controller
}

func (suite *MerchantServiceSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockMerchantRepo = mocks.NewMockMerchantRepository(suite.mockCtrl)
	suite.mockLogger = mock_logger.NewMockLoggerInterface(suite.mockCtrl)
	suite.mockMapper = mock_responseservice.NewMockMerchantResponseMapper(suite.mockCtrl)
}

func (suite *MerchantServiceSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *MerchantServiceSuite) TestFindAll_Success() {
	req := &requests.FindAllMerchants{Search: "store", Page: 1, PageSize: 10}
	total := 1
	mockMerchants := []*record.MerchantRecord{{ID: 1, Name: "Test Store"}}
	expectedResponse := []*response.MerchantResponse{{ID: 1, Name: "Test Store"}}

	suite.mockLogger.EXPECT().Debug("Fetching all merchant records", gomock.Any())
	suite.mockMerchantRepo.EXPECT().FindAllMerchants(req).Return(mockMerchants, &total, nil)
	suite.mockMapper.EXPECT().ToMerchantsResponse(mockMerchants).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully all merchant records", gomock.Any())

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *MerchantServiceSuite) TestFindAll_Failure() {
	req := &requests.FindAllMerchants{Search: "store", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching all merchant records", gomock.Any())
	suite.mockMerchantRepo.EXPECT().FindAllMerchants(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch merchants", gomock.Any())

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(merchant_errors.ErrFailedFindAllMerchants, err)
}

func (suite *MerchantServiceSuite) TestFindById_Success() {
	merchantID := 1
	mockMerchant := &record.MerchantRecord{ID: merchantID, Name: "Test Store"}
	expectedResponse := &response.MerchantResponse{ID: merchantID, Name: "Test Store"}

	suite.mockLogger.EXPECT().Debug("Finding merchant by ID", zap.Int("merchant_id", merchantID))
	suite.mockMerchantRepo.EXPECT().FindById(merchantID).Return(mockMerchant, nil)
	suite.mockMapper.EXPECT().ToMerchantResponse(mockMerchant).Return(expectedResponse)

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(merchantID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *MerchantServiceSuite) TestFindById_NotFound() {
	merchantID := 99
	repoError := errors.New("merchant not found")

	suite.mockLogger.EXPECT().Debug("Finding merchant by ID", zap.Int("merchant_id", merchantID))
	suite.mockMerchantRepo.EXPECT().FindById(merchantID).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve merchant details", gomock.Any())

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(merchantID)

	suite.Nil(result)
	suite.Equal(merchant_errors.ErrMerchantNotFoundRes, err)
}

func (suite *MerchantServiceSuite) TestFindByActive_Success() {
	req := &requests.FindAllMerchants{Search: "active", Page: 1, PageSize: 10}
	total := 1
	mockMerchants := []*record.MerchantRecord{{ID: 1, Name: "Active Store", DeletedAt: nil}}
	expectedResponse := []*response.MerchantResponseDeleteAt{{ID: 1, Name: "Active Store", DeletedAt: nil}}

	suite.mockLogger.EXPECT().Debug("Fetching all merchant active", gomock.Any())
	suite.mockMerchantRepo.EXPECT().FindByActive(req).Return(mockMerchants, &total, nil)
	suite.mockMapper.EXPECT().ToMerchantsResponseDeleteAt(mockMerchants).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched active merchants", gomock.Any())

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *MerchantServiceSuite) TestFindByActive_Failure() {
	req := &requests.FindAllMerchants{Search: "active", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching all merchant active", gomock.Any())
	suite.mockMerchantRepo.EXPECT().FindByActive(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve active cashiers", gomock.Any())

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(merchant_errors.ErrFailedFindActiveMerchants, err)
}

func (suite *MerchantServiceSuite) TestFindByTrashed_Success() {
	req := &requests.FindAllMerchants{Page: 1, PageSize: 10}
	total := 1
	trashedTime := "2024-01-01T00:00:00Z"
	mockMerchants := []*record.MerchantRecord{{ID: 2, Name: "Deleted Store", DeletedAt: &trashedTime}}
	expectedResponse := []*response.MerchantResponseDeleteAt{{ID: 2, Name: "Deleted Store", DeletedAt: &trashedTime}}

	suite.mockLogger.EXPECT().Debug("Fetching fetched trashed merchants", gomock.Any())
	suite.mockMerchantRepo.EXPECT().FindByTrashed(req).Return(mockMerchants, &total, nil)
	suite.mockMapper.EXPECT().ToMerchantsResponseDeleteAt(mockMerchants).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched trashed merchants", gomock.Any())

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *MerchantServiceSuite) TestFindByTrashed_Failure() {
	req := &requests.FindAllMerchants{Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching fetched trashed merchants", gomock.Any())
	suite.mockMerchantRepo.EXPECT().FindByTrashed(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch trashed merchants", gomock.Any())

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(merchant_errors.ErrFailedFindTrashedMerchants, err)
}

func (suite *MerchantServiceSuite) TestFindByApiKey_Success() {
	apiKey := "secret-api-key-123"
	mockMerchant := &record.MerchantRecord{ID: 1, ApiKey: apiKey}
	expectedResponse := &response.MerchantResponse{ID: 1, ApiKey: apiKey}

	suite.mockLogger.EXPECT().Debug("Finding merchant by API key", zap.String("api_key", apiKey))
	suite.mockMerchantRepo.EXPECT().FindByApiKey(apiKey).Return(mockMerchant, nil)
	suite.mockMapper.EXPECT().ToMerchantResponse(mockMerchant).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully found merchant by API key", zap.String("api_key", apiKey))

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindByApiKey(apiKey)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *MerchantServiceSuite) TestFindByApiKey_NotFound() {
	apiKey := "notfound-key"
	repoError := errors.New("merchant not found")

	suite.mockLogger.EXPECT().Debug("Finding merchant by API key", zap.String("api_key", apiKey))
	suite.mockMerchantRepo.EXPECT().FindByApiKey(apiKey).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve merchant by api_key", gomock.Any())

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindByApiKey(apiKey)

	suite.Nil(result)
	suite.Equal(merchant_errors.ErrMerchantNotFoundRes, err)
}

func (suite *MerchantServiceSuite) TestFindByMerchantUserId_Success() {
	userID := 123
	mockMerchants := []*record.MerchantRecord{{ID: 1, UserID: userID}}
	expectedResponse := []*response.MerchantResponse{{ID: 1, UserID: userID}}

	suite.mockLogger.EXPECT().Debug("Finding merchant by user ID", zap.Int("user_id", userID))
	suite.mockMerchantRepo.EXPECT().FindByMerchantUserId(userID).Return(mockMerchants, nil)
	suite.mockMapper.EXPECT().ToMerchantsResponse(mockMerchants).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully found merchant by user ID", zap.Int("user_id", userID))

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindByMerchantUserId(userID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *MerchantServiceSuite) TestFindByMerchantUserId_Failure() {
	userID := 999
	repoError := errors.New("merchant not found")

	suite.mockLogger.EXPECT().Debug("Finding merchant by user ID", zap.Int("user_id", userID))
	suite.mockMerchantRepo.EXPECT().FindByMerchantUserId(userID).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve merchant by user_id", gomock.Any())

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindByMerchantUserId(userID)

	suite.Nil(result)
	suite.Equal(merchant_errors.ErrMerchantNotFoundRes, err)
}

func (suite *MerchantServiceSuite) TestTrashedMerchant_Success() {
	merchantID := 1
	trashedMerchant := &record.MerchantRecord{ID: merchantID}
	expectedResponse := &response.MerchantResponseDeleteAt{ID: merchantID}

	suite.mockLogger.EXPECT().Debug("Trashing merchant", zap.Int("merchant_id", merchantID))
	suite.mockMerchantRepo.EXPECT().TrashedMerchant(merchantID).Return(trashedMerchant, nil)
	suite.mockMapper.EXPECT().ToMerchantResponseDeleteAt(trashedMerchant).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully trashed merchant", zap.Int("merchant_id", merchantID))

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashedMerchant(merchantID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *MerchantServiceSuite) TestTrashedMerchant_Failure() {
	merchantID := 1
	dbError := errors.New("failed to trash")

	suite.mockLogger.EXPECT().Debug("Trashing merchant", zap.Int("merchant_id", merchantID))
	suite.mockLogger.EXPECT().Error("Failed to trash merchant", gomock.Any())
	suite.mockMerchantRepo.EXPECT().TrashedMerchant(merchantID).Return(nil, dbError)

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashedMerchant(merchantID)

	suite.Nil(result)
	suite.Equal(merchant_errors.ErrFailedTrashMerchant, err)
}

func (suite *MerchantServiceSuite) TestRestoreMerchant_Success() {
	merchantID := 1
	restoredMerchant := &record.MerchantRecord{ID: merchantID, DeletedAt: nil}
	expectedResponse := &response.MerchantResponseDeleteAt{ID: merchantID, DeletedAt: nil}

	suite.mockLogger.EXPECT().Debug("Restoring merchant", zap.Int("merchant_id", merchantID))
	suite.mockMerchantRepo.EXPECT().RestoreMerchant(merchantID).Return(restoredMerchant, nil)
	suite.mockMapper.EXPECT().ToMerchantResponseDeleteAt(restoredMerchant).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully restored merchant", zap.Int("merchant_id", merchantID))

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreMerchant(merchantID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *MerchantServiceSuite) TestRestoreMerchant_Failure() {
	merchantID := 1
	dbError := errors.New("failed to restore")

	suite.mockLogger.EXPECT().Debug("Restoring merchant", zap.Int("merchant_id", merchantID))
	suite.mockLogger.EXPECT().Error("Failed to restore merchant", gomock.Any())
	suite.mockMerchantRepo.EXPECT().RestoreMerchant(merchantID).Return(nil, dbError)

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreMerchant(merchantID)

	suite.Nil(result)
	suite.Equal(merchant_errors.ErrFailedRestoreMerchant, err)
}

func (suite *MerchantServiceSuite) TestDeleteMerchantPermanent_Success() {
	merchantID := 1

	suite.mockLogger.EXPECT().Debug("Deleting merchant permanently", zap.Int("merchant_id", merchantID))
	suite.mockMerchantRepo.EXPECT().DeleteMerchantPermanent(merchantID).Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted merchant permanently", zap.Int("merchant_id", merchantID))

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteMerchantPermanent(merchantID)

	suite.Nil(err)
	suite.True(result)
}

func (suite *MerchantServiceSuite) TestDeleteMerchantPermanent_Failure() {
	merchantID := 1
	dbError := errors.New("failed to delete permanently")

	suite.mockLogger.EXPECT().Debug("Deleting merchant permanently", zap.Int("merchant_id", merchantID))
	suite.mockLogger.EXPECT().Error("Failed to delete merchant permanently", gomock.Any())
	suite.mockMerchantRepo.EXPECT().DeleteMerchantPermanent(merchantID).Return(false, dbError)

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteMerchantPermanent(merchantID)

	suite.False(result)
	suite.Equal(merchant_errors.ErrFailedDeleteMerchant, err)
}

func (suite *MerchantServiceSuite) TestRestoreAllMerchant_Success() {
	suite.mockLogger.EXPECT().Debug("Restoring all merchants")
	suite.mockMerchantRepo.EXPECT().RestoreAllMerchant().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully restored all merchants")

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllMerchant()

	suite.Nil(err)
	suite.True(result)
}

func (suite *MerchantServiceSuite) TestRestoreAllMerchant_Failure() {
	dbError := errors.New("failed to restore all")

	suite.mockLogger.EXPECT().Debug("Restoring all merchants")
	suite.mockLogger.EXPECT().Error("Failed to restore all merchants", gomock.Any())
	suite.mockMerchantRepo.EXPECT().RestoreAllMerchant().Return(false, dbError)

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllMerchant()

	suite.False(result)
	suite.Equal(merchant_errors.ErrFailedRestoreAllMerchants, err)
}

func (suite *MerchantServiceSuite) TestDeleteAllMerchantPermanent_Success() {
	suite.mockLogger.EXPECT().Debug("Permanently deleting all merchants")
	suite.mockMerchantRepo.EXPECT().DeleteAllMerchantPermanent().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted all merchants permanently")

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllMerchantPermanent()

	suite.Nil(err)
	suite.True(result)
}

func (suite *MerchantServiceSuite) TestDeleteAllMerchantPermanent_Failure() {
	dbError := errors.New("failed to delete all permanently")

	suite.mockLogger.EXPECT().Debug("Permanently deleting all merchants")
	suite.mockLogger.EXPECT().Error("Failed to permanently delete all merchants", gomock.Any())
	suite.mockMerchantRepo.EXPECT().DeleteAllMerchantPermanent().Return(false, dbError)

	svc := service.NewMerchantService(suite.mockMerchantRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllMerchantPermanent()

	suite.False(result)
	suite.Equal(merchant_errors.ErrFailedDeleteAllMerchants, err)
}

func TestMerchantServiceSuite(t *testing.T) {
	suite.Run(t, new(MerchantServiceSuite))
}
