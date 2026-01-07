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
	"MamangRust/paymentgatewaygrpc/pkg/errors/withdraw_errors"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WithdrawServiceSuite struct {
	suite.Suite
	mockWithdrawRepo *mocks.MockWithdrawRepository
	mockSaldoRepo    *mocks.MockSaldoRepository
	mockUserRepo     *mocks.MockUserRepository
	mockLogger       *mock_logger.MockLoggerInterface
	mockMapper       *mock_responseservice.MockWithdrawResponseMapper
	mockCtrl         *gomock.Controller
}

func (suite *WithdrawServiceSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockWithdrawRepo = mocks.NewMockWithdrawRepository(suite.mockCtrl)
	suite.mockSaldoRepo = mocks.NewMockSaldoRepository(suite.mockCtrl)
	suite.mockUserRepo = mocks.NewMockUserRepository(suite.mockCtrl)
	suite.mockLogger = mock_logger.NewMockLoggerInterface(suite.mockCtrl)
	suite.mockMapper = mock_responseservice.NewMockWithdrawResponseMapper(suite.mockCtrl)
}

func (suite *WithdrawServiceSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *WithdrawServiceSuite) TestFindAll_Success() {
	req := &requests.FindAllWithdraws{Search: "success", Page: 1, PageSize: 10}
	total := 1
	mockWithdraws := []*record.WithdrawRecord{{ID: 1, WithdrawNo: "WD001"}}
	expectedResponse := []*response.WithdrawResponse{{ID: 1, WithdrawNo: "WD001"}}

	suite.mockLogger.EXPECT().Debug("Fetching withdraw", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().FindAll(req).Return(mockWithdraws, &total, nil)
	suite.mockMapper.EXPECT().ToWithdrawsResponse(mockWithdraws).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched withdraw", gomock.Any())

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *WithdrawServiceSuite) TestFindAll_Failure() {
	req := &requests.FindAllWithdraws{Search: "success", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching withdraw", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().FindAll(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch withdraw", gomock.Any())

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(withdraw_errors.ErrFailedFindAllWithdraws, err)
}

func (suite *WithdrawServiceSuite) TestFindAllByCardNumber_Success() {
	req := &requests.FindAllWithdrawCardNumber{CardNumber: "1111-xxxx", Page: 1, PageSize: 10}
	total := 1
	mockWithdraws := []*record.WithdrawRecord{{ID: 2, CardNumber: "1111-xxxx"}}
	expectedResponse := []*response.WithdrawResponse{{ID: 2, CardNumber: "1111-xxxx"}}

	suite.mockLogger.EXPECT().Debug("Fetching withdraw", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().FindAllByCardNumber(req).Return(mockWithdraws, &total, nil)
	suite.mockMapper.EXPECT().ToWithdrawsResponse(mockWithdraws).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched withdraw", gomock.Any())

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAllByCardNumber(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *WithdrawServiceSuite) TestFindAllByCardNumber_Failure() {
	req := &requests.FindAllWithdrawCardNumber{CardNumber: "1111-xxxx", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching withdraw", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().FindAllByCardNumber(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch withdraw", gomock.Any())

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAllByCardNumber(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(withdraw_errors.ErrFailedFindAllWithdrawsByCard, err)
}

func (suite *WithdrawServiceSuite) TestFindById_Success() {
	withdrawID := 1
	mockWithdraw := &record.WithdrawRecord{ID: withdrawID, WithdrawNo: "WD001"}
	expectedResponse := &response.WithdrawResponse{ID: withdrawID, WithdrawNo: "WD001"}

	suite.mockLogger.EXPECT().Debug("Fetching withdraw by ID", zap.Int("withdraw_id", withdrawID))
	suite.mockWithdrawRepo.EXPECT().FindById(withdrawID).Return(mockWithdraw, nil)
	suite.mockMapper.EXPECT().ToWithdrawResponse(mockWithdraw).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched withdraw", zap.Int("withdraw_id", withdrawID))

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(withdrawID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *WithdrawServiceSuite) TestFindById_NotFound() {
	withdrawID := 99
	repoError := errors.New("withdraw not found")

	suite.mockLogger.EXPECT().Debug("Fetching withdraw by ID", zap.Int("withdraw_id", withdrawID))
	suite.mockWithdrawRepo.EXPECT().FindById(withdrawID).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("failed to find withdraw by id", zap.Error(repoError))

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(withdrawID)

	suite.Nil(result)
	suite.Equal(withdraw_errors.ErrWithdrawNotFound, err)
}

func (suite *WithdrawServiceSuite) TestFindByActive_Success() {
	req := &requests.FindAllWithdraws{Search: "active", Page: 1, PageSize: 10}
	total := 1
	mockWithdraws := []*record.WithdrawRecord{{ID: 1, WithdrawNo: "WD001", DeletedAt: nil}}
	expectedResponse := []*response.WithdrawResponseDeleteAt{{ID: 1, WithdrawNo: "WD001", DeletedAt: nil}}

	suite.mockLogger.EXPECT().Debug("Fetching active withdraw", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().FindByActive(req).Return(mockWithdraws, &total, nil)
	suite.mockMapper.EXPECT().ToWithdrawsResponseDeleteAt(mockWithdraws).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched active withdraw", gomock.Any())

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *WithdrawServiceSuite) TestFindByActive_Failure() {
	req := &requests.FindAllWithdraws{Search: "active", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching active withdraw", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().FindByActive(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch active withdraw", gomock.Any())

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(withdraw_errors.ErrFailedFindActiveWithdraws, err)
}

func (suite *WithdrawServiceSuite) TestFindByTrashed_Success() {
	req := &requests.FindAllWithdraws{Page: 1, PageSize: 10}
	total := 1
	trashedTime := "2024-01-01T00:00:00Z"
	mockWithdraws := []*record.WithdrawRecord{{ID: 2, WithdrawNo: "WD002", DeletedAt: &trashedTime}}
	expectedResponse := []*response.WithdrawResponseDeleteAt{{ID: 2, WithdrawNo: "WD002", DeletedAt: &trashedTime}}

	suite.mockLogger.EXPECT().Debug("Fetching trashed withdraw", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().FindByTrashed(req).Return(mockWithdraws, &total, nil)
	suite.mockMapper.EXPECT().ToWithdrawsResponseDeleteAt(mockWithdraws).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched trashed withdraw", gomock.Any())

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *WithdrawServiceSuite) TestFindByTrashed_Failure() {
	req := &requests.FindAllWithdraws{Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching trashed withdraw", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().FindByTrashed(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch trashed withdraw", gomock.Any())

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(withdraw_errors.ErrFailedFindTrashedWithdraws, err)
}

func (suite *WithdrawServiceSuite) TestTrashedWithdraw_Success() {
	withdrawID := 1
	trashedWithdraw := &record.WithdrawRecord{ID: withdrawID}
	expectedResponse := &response.WithdrawResponseDeleteAt{ID: withdrawID}

	suite.mockLogger.EXPECT().Debug("Trashing withdraw", zap.Int("withdraw_id", withdrawID))
	suite.mockWithdrawRepo.EXPECT().TrashedWithdraw(withdrawID).Return(trashedWithdraw, nil)
	suite.mockMapper.EXPECT().ToWithdrawResponseDeleteAt(trashedWithdraw).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully trashed withdraw", zap.Int("withdraw_id", withdrawID))

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashedWithdraw(withdrawID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *WithdrawServiceSuite) TestTrashedWithdraw_Failure() {
	withdrawID := 1
	dbError := errors.New("failed to trash")

	suite.mockLogger.EXPECT().Debug("Trashing withdraw", zap.Int("withdraw_id", withdrawID))
	suite.mockLogger.EXPECT().Error("Failed to move withdraw to trash", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().TrashedWithdraw(withdrawID).Return(nil, dbError)

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashedWithdraw(withdrawID)

	suite.Nil(result)
	suite.Equal(withdraw_errors.ErrFailedTrashedWithdraw, err)
}

func (suite *WithdrawServiceSuite) TestRestoreWithdraw_Success() {
	withdrawID := 1
	restoredWithdraw := &record.WithdrawRecord{ID: withdrawID, DeletedAt: nil}
	expectedResponse := &response.WithdrawResponseDeleteAt{ID: withdrawID, DeletedAt: nil}

	suite.mockLogger.EXPECT().Debug("Restoring withdraw", zap.Int("withdraw_id", withdrawID))
	suite.mockWithdrawRepo.EXPECT().RestoreWithdraw(withdrawID).Return(restoredWithdraw, nil)
	suite.mockMapper.EXPECT().ToWithdrawResponseDeleteAt(restoredWithdraw).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully restored withdraw", zap.Int("withdraw_id", withdrawID))

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreWithdraw(withdrawID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *WithdrawServiceSuite) TestRestoreWithdraw_Failure() {
	withdrawID := 1
	dbError := errors.New("failed to restore")

	suite.mockLogger.EXPECT().Debug("Restoring withdraw", zap.Int("withdraw_id", withdrawID))
	suite.mockLogger.EXPECT().Error("Failed to restore withdraw from trash", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().RestoreWithdraw(withdrawID).Return(nil, dbError)

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreWithdraw(withdrawID)

	suite.Nil(result)
	suite.Equal(withdraw_errors.ErrFailedRestoreWithdraw, err)
}

func (suite *WithdrawServiceSuite) TestDeleteWithdrawPermanent_Success() {
	withdrawID := 1

	suite.mockLogger.EXPECT().Debug("Deleting withdraw permanently", zap.Int("withdraw_id", withdrawID))
	suite.mockWithdrawRepo.EXPECT().DeleteWithdrawPermanent(withdrawID).Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted withdraw permanently", zap.Int("withdraw_id", withdrawID))

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteWithdrawPermanent(withdrawID)

	suite.Nil(err)
	suite.True(result)
}

func (suite *WithdrawServiceSuite) TestDeleteWithdrawPermanent_Failure() {
	withdrawID := 1
	dbError := errors.New("failed to delete permanently")

	suite.mockLogger.EXPECT().Debug("Deleting withdraw permanently", zap.Int("withdraw_id", withdrawID))
	suite.mockLogger.EXPECT().Error("Failed to permanently delete withdraw", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().DeleteWithdrawPermanent(withdrawID).Return(false, dbError)

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteWithdrawPermanent(withdrawID)

	suite.False(result)
	suite.Equal(withdraw_errors.ErrFailedDeleteWithdrawPermanent, err)
}

func (suite *WithdrawServiceSuite) TestRestoreAllWithdraw_Success() {
	suite.mockLogger.EXPECT().Debug("Restoring all withdraws")
	suite.mockWithdrawRepo.EXPECT().RestoreAllWithdraw().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully restored all withdraws")

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllWithdraw()

	suite.Nil(err)
	suite.True(result)
}

func (suite *WithdrawServiceSuite) TestRestoreAllWithdraw_Failure() {
	dbError := errors.New("failed to restore all")

	suite.mockLogger.EXPECT().Debug("Restoring all withdraws")
	suite.mockLogger.EXPECT().Error("Failed to restore all withdraws", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().RestoreAllWithdraw().Return(false, dbError)

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllWithdraw()

	suite.False(result)
	suite.Equal(withdraw_errors.ErrFailedRestoreAllWithdraw, err)
}

func (suite *WithdrawServiceSuite) TestDeleteAllWithdrawPermanent_Success() {
	suite.mockLogger.EXPECT().Debug("Permanently deleting all withdraws")
	suite.mockWithdrawRepo.EXPECT().DeleteAllWithdrawPermanent().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted all withdraws permanently")

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllWithdrawPermanent()

	suite.Nil(err)
	suite.True(result)
}

func (suite *WithdrawServiceSuite) TestDeleteAllWithdrawPermanent_Failure() {
	dbError := errors.New("failed to delete all permanently")

	suite.mockLogger.EXPECT().Debug("Permanently deleting all withdraws")
	suite.mockLogger.EXPECT().Error("Failed to permanently delete all withdraws", gomock.Any())
	suite.mockWithdrawRepo.EXPECT().DeleteAllWithdrawPermanent().Return(false, dbError)

	svc := service.NewWithdrawService(suite.mockUserRepo, suite.mockWithdrawRepo, suite.mockSaldoRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllWithdrawPermanent()

	suite.False(result)
	suite.Equal(withdraw_errors.ErrFailedDeleteAllWithdrawPermanent, err)
}

func TestWithdrawServiceSuite(t *testing.T) {
	suite.Run(t, new(WithdrawServiceSuite))
}
