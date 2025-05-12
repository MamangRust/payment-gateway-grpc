package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	mock_responsemapper "MamangRust/paymentgatewaygrpc/internal/mapper/response/mocks"
	mock_repository "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/internal/service"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestWithdrawService_FindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockWithdrawResponseMapper(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		mockMapper,
	)

	page := 1
	pageSize := 1
	search := "user1"
	expectedWithdraws := []*response.WithdrawResponse{
		{
			ID:             1,
			CardNumber:     "card_1234",
			WithdrawAmount: 100000,
			WithdrawTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
		{
			ID:             2,
			CardNumber:     "card_1234",
			WithdrawAmount: 200000,
			WithdrawTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
	}
	totalRecords := 2

	mockWithdrawRepo.EXPECT().
		FindAll(search, page, pageSize).
		Return([]*record.WithdrawRecord{
			{
				ID:             1,
				CardNumber:     "card_1234",
				WithdrawAmount: 100000,
				WithdrawTime:   time.Now().Format(time.RFC3339),
				CreatedAt:      time.Now().Format(time.RFC3339),
				UpdatedAt:      time.Now().Format(time.RFC3339),
			},
			{
				ID:             2,
				CardNumber:     "card_1234",
				WithdrawAmount: 200000,
				WithdrawTime:   time.Now().Format(time.RFC3339),
				CreatedAt:      time.Now().Format(time.RFC3339),
				UpdatedAt:      time.Now().Format(time.RFC3339),
			},
		}, totalRecords, nil)

	mockMapper.EXPECT().
		ToWithdrawsResponse(gomock.Any()).
		Return(expectedWithdraws)

	result, totalPages, errResp := withdrawService.FindAll(page, pageSize, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedWithdraws, result)
	assert.Equal(t, 2, totalPages)
}

func TestWithdrawService_FindAll_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		nil,
	)

	page := 1
	pageSize := 10
	search := "user1"

	mockWithdrawRepo.EXPECT().
		FindAll(search, page, pageSize).
		Return(nil, 0, errors.New("failed to fetch withdraws"))

	mockLogger.EXPECT().
		Error("failed to fetch withdraws", gomock.Any())

	result, totalPages, errResp := withdrawService.FindAll(page, pageSize, search)

	assert.Nil(t, result)
	assert.Equal(t, 0, totalPages)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to fetch withdraws", errResp.Message)
}

func TestWithdrawService_FindById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockWithdrawResponseMapper(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		mockMapper,
	)

	withdrawID := 1

	mockWithdrawRecord := &record.WithdrawRecord{
		ID:             withdrawID,
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 5000,
		WithdrawTime:   "2024-12-25T14:30:00Z",
		CreatedAt:      "2024-12-25T14:00:00Z",
		UpdatedAt:      "2024-12-25T14:15:00Z",
		DeletedAt:      nil,
	}

	mockWithdrawResponse := &response.WithdrawResponse{
		ID:             withdrawID,
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 5000,
		WithdrawTime:   "2024-12-25T14:30:00Z",
		CreatedAt:      "2024-12-25T14:00:00Z",
		UpdatedAt:      "2024-12-25T14:15:00Z",
	}

	mockWithdrawRepo.EXPECT().FindById(withdrawID).Return(mockWithdrawRecord, nil)
	mockMapper.EXPECT().ToWithdrawResponse(mockWithdrawRecord).Return(mockWithdrawResponse)

	result, err := withdrawService.FindById(withdrawID)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockWithdrawResponse, result)
}

func TestWithdrawService_FindById_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		nil,
	)

	withdrawID := 1

	mockWithdrawRepo.EXPECT().FindById(withdrawID).Return(nil, errors.New("record not found"))

	mockLogger.EXPECT().Error("failed to find withdraw by id", gomock.Any())

	result, err := withdrawService.FindById(withdrawID)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to fetch withdraw record by ID.", err.Message)
}

func TestWithdrawService_FindByCardNumber_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockWithdrawResponseMapper(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		mockMapper,
	)

	cardNumber := "1234-5678-9012-3456"
	mockWithdrawRecords := []*record.WithdrawRecord{
		{ID: 1, CardNumber: cardNumber, WithdrawAmount: 5000},
		{ID: 2, CardNumber: cardNumber, WithdrawAmount: 10000},
	}
	mockWithdrawResponses := []*response.WithdrawResponse{
		{ID: 1, CardNumber: cardNumber, WithdrawAmount: 5000},
		{ID: 2, CardNumber: cardNumber, WithdrawAmount: 10000},
	}

	mockWithdrawRepo.EXPECT().FindByCardNumber(cardNumber).Return(mockWithdrawRecords, nil)
	mockMapper.EXPECT().ToWithdrawsResponse(mockWithdrawRecords).Return(mockWithdrawResponses)

	result, err := withdrawService.FindByCardNumber(cardNumber)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockWithdrawResponses, result)
}

func TestWithdrawService_FindByCardNumber_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		nil,
	)

	cardNumber := "1234-5678-9012-3456"

	mockWithdrawRepo.EXPECT().FindByCardNumber(cardNumber).Return(nil, errors.New("database error"))
	mockLogger.EXPECT().Error("Failed to fetch withdraw records by card number", gomock.Any(), gomock.Any())

	result, err := withdrawService.FindByCardNumber(cardNumber)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to fetch withdraw records for the given card number", err.Message)
}

func TestWithdrawService_FindByCardNumber_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockMapper := mock_responsemapper.NewMockWithdrawResponseMapper(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		nil,
		mockMapper,
	)

	cardNumber := "1234-5678-9012-3456"
	mockWithdrawRecords := []*record.WithdrawRecord{}

	mockWithdrawRepo.EXPECT().FindByCardNumber(cardNumber).Return(mockWithdrawRecords, nil)
	mockMapper.EXPECT().ToWithdrawsResponse(mockWithdrawRecords).Return([]*response.WithdrawResponse{})

	result, err := withdrawService.FindByCardNumber(cardNumber)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

func TestWithdrawService_FindByActive_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockWithdrawResponseMapper(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		mockMapper,
	)

	mockWithdrawRecords := []*record.WithdrawRecord{
		{ID: 1, CardNumber: "1234-5678-9012-3456", WithdrawAmount: 5000},
		{ID: 2, CardNumber: "9876-5432-1098-7654", WithdrawAmount: 10000},
	}
	mockWithdrawResponses := []*response.WithdrawResponseDeleteAt{
		{ID: 1, CardNumber: "1234-5678-9012-3456", WithdrawAmount: 5000},
		{ID: 2, CardNumber: "9876-5432-1098-7654", WithdrawAmount: 10000},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 2

	mockWithdrawRepo.EXPECT().FindByActive(search, page, pageSize).Return(mockWithdrawRecords, expected, nil)
	mockMapper.EXPECT().ToWithdrawsResponseDeleteAt(mockWithdrawRecords).Return(mockWithdrawResponses)

	result, totalRecord, err := withdrawService.FindByActive(pageSize, page, search)

	assert.Nil(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.NotNil(t, result)
	assert.Equal(t, mockWithdrawResponses, result)
}

func TestWithdrawService_FindByActive_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		nil,
	)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	mockLogger.EXPECT().Error("Failed to fetch active withdraw records", gomock.Any())

	mockWithdrawRepo.EXPECT().FindByActive(search, page, pageSize).Return(nil, expected, errors.New("database error"))

	result, totalRecord, err := withdrawService.FindByActive(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "Failed to fetch active withdraw records", err.Message)
}

func TestWithdrawService_FindByActive_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockMapper := mock_responsemapper.NewMockWithdrawResponseMapper(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		nil,
		mockMapper,
	)

	page := 1
	pageSize := 1
	search := "user1"
	expected := 0

	mockWithdrawRecords := []*record.WithdrawRecord{}

	mockWithdrawRepo.EXPECT().FindByActive(search, page, pageSize).Return(mockWithdrawRecords, expected, nil)
	mockMapper.EXPECT().ToWithdrawsResponseDeleteAt(mockWithdrawRecords).Return([]*response.WithdrawResponseDeleteAt{})

	result, totalRecord, err := withdrawService.FindByActive(pageSize, page, search)

	assert.Nil(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

func TestWithdrawService_FindByTrashed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockWithdrawResponseMapper(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		mockMapper,
	)

	page := 1
	pageSize := 1
	search := "user1"
	expected := 2

	mockWithdrawRecords := []*record.WithdrawRecord{
		{ID: 1, CardNumber: "1234-5678-9012-3456", WithdrawAmount: 5000},
		{ID: 2, CardNumber: "9876-5432-1098-7654", WithdrawAmount: 10000},
	}
	mockWithdrawResponses := []*response.WithdrawResponseDeleteAt{
		{ID: 1, CardNumber: "1234-5678-9012-3456", WithdrawAmount: 5000},
		{ID: 2, CardNumber: "9876-5432-1098-7654", WithdrawAmount: 10000},
	}

	mockWithdrawRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(mockWithdrawRecords, expected, nil)
	mockMapper.EXPECT().ToWithdrawsResponseDeleteAt(mockWithdrawRecords).Return(mockWithdrawResponses)

	result, totalRecord, err := withdrawService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, mockWithdrawResponses, result)
}

func TestWithdrawService_FindByTrashed_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		nil,
	)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	mockLogger.EXPECT().Error("Failed to fetch trashed withdraw records", gomock.Any())

	mockWithdrawRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, expected, errors.New("database error"))

	result, totalRecord, err := withdrawService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "Failed to fetch trashed withdraw records", err.Message)
}

func TestWithdrawService_FindByTrashed_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockMapper := mock_responsemapper.NewMockWithdrawResponseMapper(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		nil,
		mockMapper,
	)
	page := 1
	pageSize := 1
	search := "user1"
	expected := 2

	mockWithdrawRecords := []*record.WithdrawRecord{}

	mockWithdrawRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(mockWithdrawRecords, expected, nil)
	mockMapper.EXPECT().ToWithdrawsResponseDeleteAt(mockWithdrawRecords).Return([]*response.WithdrawResponseDeleteAt{})

	result, totalRecord, err := withdrawService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

func TestWithdrawService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoRepo := mock_repository.NewMockSaldoRepository(ctrl)
	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockWithdrawResponseMapper(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		mockSaldoRepo,
		mockLogger,
		mockMapper,
	)

	request := &requests.CreateWithdrawRequest{
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 1000,
		WithdrawTime:   time.Now(),
	}

	mockSaldo := &record.SaldoRecord{
		CardNumber:   "1234-5678-9012-3456",
		TotalBalance: 5000,
	}

	mockWithdrawRecord := &record.WithdrawRecord{
		ID:             1,
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 1000,
		WithdrawTime:   "2024-12-26T10:00:00",
	}

	mockWithdrawResponse := &response.WithdrawResponse{
		ID:             1,
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 1000,
		WithdrawTime:   "2024-12-26T10:00:00",
	}

	mockSaldoRepo.EXPECT().FindByCardNumber(request.CardNumber).Return(mockSaldo, nil)
	mockSaldoRepo.EXPECT().UpdateSaldoWithdraw(gomock.Any()).Return(nil, nil)
	mockWithdrawRepo.EXPECT().CreateWithdraw(request).Return(mockWithdrawRecord, nil)
	mockMapper.EXPECT().ToWithdrawResponse(mockWithdrawRecord).Return(mockWithdrawResponse)

	result, err := withdrawService.Create(request)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockWithdrawResponse, result)
}

func TestWithdrawService_Create_Failure_SaldoNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoRepo := mock_repository.NewMockSaldoRepository(ctrl)
	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		mockSaldoRepo,
		mockLogger,
		nil,
	)

	request := &requests.CreateWithdrawRequest{
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 1000,
		WithdrawTime:   time.Now(),
	}

	mockLogger.EXPECT().Error("Failed to find saldo by user ID", zap.Error(errors.New("saldo not found")))

	mockSaldoRepo.EXPECT().FindByCardNumber(request.CardNumber).Return(nil, errors.New("saldo not found"))

	result, err := withdrawService.Create(request)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to fetch saldo for the user.", err.Message)
}

func TestWithdrawService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoRepo := mock_repository.NewMockSaldoRepository(ctrl)
	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockWithdrawResponseMapper(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		mockSaldoRepo,
		mockLogger,
		mockMapper,
	)

	request := &requests.UpdateWithdrawRequest{
		WithdrawID:     1,
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 1000,
		WithdrawTime:   time.Now(),
	}

	mockSaldo := &record.SaldoRecord{
		CardNumber:   "1234-5678-9012-3456",
		TotalBalance: 5000,
	}

	mockUpdatedWithdraw := &record.WithdrawRecord{
		ID:             1,
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 1000,
		WithdrawTime:   "2024-12-26T10:00:00",
	}

	mockWithdrawResponse := &response.WithdrawResponse{
		ID:             1,
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 1000,
		WithdrawTime:   "2024-12-26T10:00:00",
	}

	mockWithdrawRepo.EXPECT().FindById(request.WithdrawID).Return(mockUpdatedWithdraw, nil)
	mockSaldoRepo.EXPECT().FindByCardNumber(request.CardNumber).Return(mockSaldo, nil)
	mockSaldoRepo.EXPECT().UpdateSaldoWithdraw(gomock.Any()).Return(nil, nil)
	mockWithdrawRepo.EXPECT().UpdateWithdraw(request).Return(mockUpdatedWithdraw, nil)
	mockMapper.EXPECT().ToWithdrawResponse(mockUpdatedWithdraw).Return(mockWithdrawResponse)

	result, err := withdrawService.Update(request)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockWithdrawResponse, result)
}

func TestWithdrawService_Update_Failure_WithdrawNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoRepo := mock_repository.NewMockSaldoRepository(ctrl)
	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		mockSaldoRepo,
		mockLogger,
		nil,
	)

	request := &requests.UpdateWithdrawRequest{
		WithdrawID:     1,
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 1000,
		WithdrawTime:   time.Now(),
	}

	mockLogger.EXPECT().Error("Failed to find withdraw record by ID", zap.Error(errors.New("withdraw record not found")))

	mockWithdrawRepo.EXPECT().FindById(request.WithdrawID).Return(nil, errors.New("withdraw record not found"))

	result, err := withdrawService.Update(request)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Withdraw record not found.", err.Message)
}

func TestWithdrawService_TrashedWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockWithdrawResponseMapper(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		mockMapper,
	)

	withdrawID := 1
	mockWithdrawRecord := &record.WithdrawRecord{
		ID:             withdrawID,
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 1000,
		WithdrawTime:   "2024-12-26T10:00:00",
	}

	mockWithdrawResponse := &response.WithdrawResponse{
		ID:             withdrawID,
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 1000,
		WithdrawTime:   "2024-12-26T10:00:00",
	}

	mockLogger.EXPECT().Debug("Trashing withdraw", zap.Int("withdraw_id", withdrawID)).Times(1)

	mockWithdrawRepo.EXPECT().TrashedWithdraw(withdrawID).Return(mockWithdrawRecord, nil)
	mockMapper.EXPECT().ToWithdrawResponse(mockWithdrawRecord).Return(mockWithdrawResponse)

	mockLogger.EXPECT().Debug("Successfully trashed withdraw", zap.Int("withdraw_id", withdrawID)).Times(1)

	result, err := withdrawService.TrashedWithdraw(withdrawID)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockWithdrawResponse, result)
}

func TestWithdrawService_TrashedWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		nil,
	)

	withdrawID := 1

	mockLogger.EXPECT().Debug("Trashing withdraw", zap.Int("withdraw_id", withdrawID)).Times(1)

	mockWithdrawRepo.EXPECT().TrashedWithdraw(withdrawID).Return(nil, errors.New("failed to trash withdraw"))
	mockLogger.EXPECT().Error("Failed to trash withdraw", zap.Error(errors.New("failed to trash withdraw")), zap.Int("withdraw_id", withdrawID)).Times(1)

	result, err := withdrawService.TrashedWithdraw(withdrawID)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to trash withdraw", err.Message)
}

func TestWithdrawService_RestoreWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockWithdrawResponseMapper(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		mockMapper,
	)

	withdrawID := 1
	mockWithdrawRecord := &record.WithdrawRecord{
		ID:             withdrawID,
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 1000,
		WithdrawTime:   "2024-12-26T10:00:00",
	}

	mockWithdrawResponse := &response.WithdrawResponse{
		ID:             withdrawID,
		CardNumber:     "1234-5678-9012-3456",
		WithdrawAmount: 1000,
		WithdrawTime:   "2024-12-26T10:00:00",
	}

	mockLogger.EXPECT().Debug("Restoring withdraw", zap.Int("withdraw_id", withdrawID)).Times(1)

	mockWithdrawRepo.EXPECT().RestoreWithdraw(withdrawID).Return(mockWithdrawRecord, nil)
	mockMapper.EXPECT().ToWithdrawResponse(mockWithdrawRecord).Return(mockWithdrawResponse)

	mockLogger.EXPECT().Debug("Successfully restored withdraw", zap.Int("withdraw_id", withdrawID)).Times(1)

	result, err := withdrawService.RestoreWithdraw(withdrawID)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockWithdrawResponse, result)
}

func TestWithdrawService_RestoreWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		nil,
	)

	withdrawID := 1

	mockLogger.EXPECT().Debug("Restoring withdraw", zap.Int("withdraw_id", withdrawID)).Times(1)

	mockWithdrawRepo.EXPECT().RestoreWithdraw(withdrawID).Return(nil, errors.New("failed to restore withdraw"))

	mockLogger.EXPECT().Error("Failed to restore withdraw", zap.Error(errors.New("failed to restore withdraw")), zap.Int("withdraw_id", withdrawID)).Times(1)

	result, err := withdrawService.RestoreWithdraw(withdrawID)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to restore withdraw", err.Message)
}

func TestWithdrawService_DeleteWithdrawPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		nil,
	)

	withdrawID := 1

	mockLogger.EXPECT().Debug("Deleting withdraw permanently", zap.Int("withdraw_id", withdrawID)).Times(1)

	mockWithdrawRepo.EXPECT().DeleteWithdrawPermanent(withdrawID).Return(nil)

	mockLogger.EXPECT().Debug("Successfully deleted withdraw permanently", zap.Int("withdraw_id", withdrawID)).Times(1)

	result, err := withdrawService.DeleteWithdrawPermanent(withdrawID)

	assert.Nil(t, err)
	assert.Nil(t, result)
}

func TestWithdrawService_DeleteWithdrawPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawRepo := mock_repository.NewMockWithdrawRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	withdrawService := service.NewWithdrawService(
		nil,
		mockWithdrawRepo,
		nil,
		mockLogger,
		nil,
	)

	withdrawID := 1

	mockLogger.EXPECT().Debug("Deleting withdraw permanently", zap.Int("withdraw_id", withdrawID)).Times(1)

	mockLogger.EXPECT().Error("Failed to delete withdraw permanently", zap.Error(errors.New("failed to delete withdraw permanently")), zap.Int("withdraw_id", withdrawID)).Times(1)

	mockWithdrawRepo.EXPECT().DeleteWithdrawPermanent(withdrawID).Return(errors.New("failed to delete withdraw permanently"))

	result, err := withdrawService.DeleteWithdrawPermanent(withdrawID)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to delete withdraw permanently", err.Message)
}
