package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
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

func TestFindAllTransfers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransferResponseMapper(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		mockMapping,
	)

	page := 1
	pageSize := 10
	search := "test"
	totalRecords := 2

	transfers := []*record.TransferRecord{
		{
			ID:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
		{
			ID:             2,
			TransferFrom:   "user1",
			TransferTo:     "user3",
			TransferAmount: 500,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
	}
	expectedResponses := []*response.TransferResponse{
		{
			ID:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
		{
			ID:             2,
			TransferFrom:   "user1",
			TransferTo:     "user3",
			TransferAmount: 500,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
	}

	mockTransferRepo.EXPECT().
		FindAll(search, page, pageSize).
		Return(transfers, totalRecords, nil)

	mockMapping.EXPECT().
		ToTransfersResponse(transfers).
		Return(expectedResponses)

	result, totalPages, errResp := transferService.FindAll(page, pageSize, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponses, result)
	assert.Equal(t, totalPages, totalRecords)
}

func TestFindAllTransfers_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		nil,
	)

	page := 1
	pageSize := 10
	search := "test"

	mockTransferRepo.EXPECT().
		FindAll(search, page, pageSize).
		Return(nil, 0, errors.New("database error"))

	mockLogger.EXPECT().
		Error("failed to fetch transfers", gomock.Any())

	result, totalPages, errResp := transferService.FindAll(page, pageSize, search)

	assert.Nil(t, result)
	assert.Equal(t, 0, totalPages)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to fetch transfers", errResp.Message)
}

func TestFindByIdTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransferResponseMapper(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		mockMapping,
	)

	transferId := 1
	transferRecord := &record.TransferRecord{
		ID:             1,
		TransferFrom:   "user1",
		TransferTo:     "user2",
		TransferAmount: 1000,
		TransferTime:   time.Now().Format(time.RFC3339),
		CreatedAt:      time.Now().Format(time.RFC3339),
		UpdatedAt:      time.Now().Format(time.RFC3339),
	}

	expectedResponse := &response.TransferResponse{
		ID:             1,
		TransferFrom:   "user1",
		TransferTo:     "user2",
		TransferAmount: 1000,
		TransferTime:   time.Now().Format(time.RFC3339),
		CreatedAt:      time.Now().Format(time.RFC3339),
		UpdatedAt:      time.Now().Format(time.RFC3339),
	}

	mockTransferRepo.EXPECT().
		FindById(transferId).
		Return(transferRecord, nil)

	mockMapping.EXPECT().
		ToTransferResponse(transferRecord).
		Return(expectedResponse)

	result, errResp := transferService.FindById(transferId)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestFindByIdTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		nil,
	)

	transferId := 1

	mockTransferRepo.EXPECT().
		FindById(transferId).
		Return(nil, errors.New("record not found"))

	mockLogger.EXPECT().
		Error("failed to find transfer by ID", gomock.Any())

	result, errResp := transferService.FindById(transferId)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Transfer not found", errResp.Message)
}

func TestFindByActiveTransfers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransferResponseMapper(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		mockMapping,
	)

	transfers := []*record.TransferRecord{
		{
			ID:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
		{
			ID:             2,
			TransferFrom:   "user1",
			TransferTo:     "user3",
			TransferAmount: 500,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
	}
	expectedResponses := []*response.TransferResponseDeleteAt{
		{
			ID:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
		{
			ID:             2,
			TransferFrom:   "user1",
			TransferTo:     "user3",
			TransferAmount: 500,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 2

	mockTransferRepo.EXPECT().
		FindByActive(search, page, pageSize).
		Return(transfers, expected, nil)

	mockMapping.EXPECT().
		ToTransfersResponseDeleteAt(transfers).
		Return(expectedResponses)

	mockLogger.EXPECT().
		Debug("Successfully fetched active transaction records", zap.Int("record_count", len(transfers)))

	result, totalRecord, errResp := transferService.FindByActive(pageSize, page, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, expectedResponses, result)
}

func TestFindByActiveTransfers_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		nil,
	)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	mockTransferRepo.EXPECT().
		FindByActive(search, page, pageSize).
		Return(nil, expected, errors.New("no active records found"))

	mockLogger.EXPECT().
		Error("Failed to fetch active transaction records", gomock.Any())

	result, totalRecord, errResp := transferService.FindByActive(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "No active transaction records found", errResp.Message)
}

func TestFindByTrashedTransfers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransferResponseMapper(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		mockMapping,
	)

	transfers := []*record.TransferRecord{
		{
			ID:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
		{
			ID:             2,
			TransferFrom:   "user1",
			TransferTo:     "user3",
			TransferAmount: 500,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
	}
	expectedResponses := []*response.TransferResponseDeleteAt{
		{
			ID:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
		{
			ID:             2,
			TransferFrom:   "user1",
			TransferTo:     "user3",
			TransferAmount: 500,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 2

	mockLogger.EXPECT().
		Info("Fetching trashed transaction records")

	mockTransferRepo.EXPECT().
		FindByTrashed(search, page, pageSize).
		Return(transfers, expected, nil)

	mockMapping.EXPECT().
		ToTransfersResponseDeleteAt(transfers).
		Return(expectedResponses)

	mockLogger.EXPECT().
		Debug("Successfully fetched trashed transaction records", zap.Int("record_count", len(transfers)))

	result, totalRecord, errResp := transferService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, expectedResponses, result)
}

func TestFindByTrashedTransfers_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		nil,
	)

	page := 1
	pageSize := 1
	search := "test"
	expected := 0

	mockLogger.EXPECT().
		Info("Fetching trashed transaction records")

	mockTransferRepo.EXPECT().
		FindByTrashed(search, page, pageSize).
		Return(nil, expected, errors.New("no trashed records found"))

	mockLogger.EXPECT().
		Error("Failed to fetch trashed transaction records", gomock.Any())

	result, totalRecord, errResp := transferService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "No trashed transaction records found", errResp.Message)
}

func TestFindTransferByTransferFrom_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransferResponseMapper(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		mockMapping,
	)

	transferFrom := "user1"

	transfers := []*record.TransferRecord{
		{
			ID:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
		{
			ID:             2,
			TransferFrom:   "user1",
			TransferTo:     "user3",
			TransferAmount: 500,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
	}

	expectedResponses := []*response.TransferResponse{
		{
			ID:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
			TransferTime:   transfers[0].TransferTime,
			CreatedAt:      transfers[0].CreatedAt,
			UpdatedAt:      transfers[0].UpdatedAt,
		},
		{
			ID:             2,
			TransferFrom:   "user1",
			TransferTo:     "user3",
			TransferAmount: 500,
			TransferTime:   transfers[1].TransferTime,
			CreatedAt:      transfers[1].CreatedAt,
			UpdatedAt:      transfers[1].UpdatedAt,
		},
	}

	mockTransferRepo.EXPECT().
		FindTransferByTransferFrom(transferFrom).
		Return(transfers, nil)

	mockMapping.EXPECT().
		ToTransfersResponse(transfers).
		Return(expectedResponses)

	result, errResp := transferService.FindTransferByTransferFrom(transferFrom)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponses, result)
}

func TestFindTransferByTransferFrom_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		nil,
	)

	transferFrom := "user1"

	mockTransferRepo.EXPECT().
		FindTransferByTransferFrom(transferFrom).
		Return(nil, errors.New("failed to fetch transfers"))

	mockLogger.EXPECT().
		Error("Failed to fetch transfers by transfer_from", gomock.Any())

	result, errResp := transferService.FindTransferByTransferFrom(transferFrom)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to fetch transfers by transfer_from", errResp.Message)
}

func TestFindTransferByTransferTo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransferResponseMapper(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		mockMapping,
	)

	transferTo := "user2"

	transfers := []*record.TransferRecord{
		{
			ID:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		},
	}

	expectedResponses := []*response.TransferResponse{
		{
			ID:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
			TransferTime:   transfers[0].TransferTime,
			CreatedAt:      transfers[0].CreatedAt,
			UpdatedAt:      transfers[0].UpdatedAt,
		},
	}

	mockTransferRepo.EXPECT().
		FindTransferByTransferTo(transferTo).
		Return(transfers, nil)

	mockMapping.EXPECT().
		ToTransfersResponse(transfers).
		Return(expectedResponses)

	result, errResp := transferService.FindTransferByTransferTo(transferTo)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponses, result)
}

func TestFindTransferByTransferTo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		nil,
	)

	transferTo := "user2"

	mockTransferRepo.EXPECT().
		FindTransferByTransferTo(transferTo).
		Return(nil, errors.New("failed to fetch transfers"))

	mockLogger.EXPECT().
		Error("Failed to fetch transfers by transfer_to", gomock.Any())

	result, errResp := transferService.FindTransferByTransferTo(transferTo)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to fetch transfers by transfer_to", errResp.Message)
}

func TestTrashedTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransferResponseMapper(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		mockMapping,
	)

	transferID := 1
	transferRecord := &record.TransferRecord{
		ID:             transferID,
		TransferFrom:   "user1",
		TransferTo:     "user2",
		TransferAmount: 1000,
		TransferTime:   time.Now().Format(time.RFC3339),
	}

	expectedResponse := &response.TransferResponse{
		ID:             transferID,
		TransferFrom:   "user1",
		TransferTo:     "user2",
		TransferAmount: 1000,
		TransferTime:   transferRecord.TransferTime,
	}

	mockTransferRepo.EXPECT().
		TrashedTransfer(transferID).
		Return(transferRecord, nil)

	mockMapping.EXPECT().
		ToTransferResponse(transferRecord).
		Return(expectedResponse)

	result, errResp := transferService.TrashedTransfer(transferID)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestTrashedTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		nil,
	)

	transferID := 1

	mockTransferRepo.EXPECT().
		TrashedTransfer(transferID).
		Return(nil, errors.New("failed to trash transfer"))

	mockLogger.EXPECT().
		Error("Failed to trash transfer", gomock.Any())

	result, errResp := transferService.TrashedTransfer(transferID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to trash transfer", errResp.Message)
}

func TestRestoreTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransferResponseMapper(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		mockMapping,
	)

	transferID := 1
	transferRecord := &record.TransferRecord{
		ID:             transferID,
		TransferFrom:   "user1",
		TransferTo:     "user2",
		TransferAmount: 1000,
		TransferTime:   time.Now().Format(time.RFC3339),
	}

	expectedResponse := &response.TransferResponse{
		ID:             transferID,
		TransferFrom:   "user1",
		TransferTo:     "user2",
		TransferAmount: 1000,
		TransferTime:   transferRecord.TransferTime,
	}

	mockTransferRepo.EXPECT().
		RestoreTransfer(transferID).
		Return(transferRecord, nil)

	mockMapping.EXPECT().
		ToTransferResponse(transferRecord).
		Return(expectedResponse)

	result, errResp := transferService.RestoreTransfer(transferID)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestRestoreTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		nil,
	)

	transferID := 1

	mockTransferRepo.EXPECT().
		RestoreTransfer(transferID).
		Return(nil, errors.New("failed to restore transfer"))

	mockLogger.EXPECT().
		Error("Failed to restore transfer", gomock.Any())

	result, errResp := transferService.RestoreTransfer(transferID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to restore transfer", errResp.Message)
}

func TestDeleteTransferPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		nil,
	)

	transferID := 1

	mockTransferRepo.EXPECT().
		DeleteTransferPermanent(transferID).
		Return(nil)

	result, errResp := transferService.DeleteTransferPermanent(transferID)

	assert.Nil(t, errResp)
	assert.Nil(t, result)
}

func TestDeleteTransferPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferRepo := mock_repository.NewMockTransferRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferService := service.NewTransferService(
		nil, nil,
		mockTransferRepo,
		nil,
		mockLogger,
		nil,
	)

	transferID := 1

	mockTransferRepo.EXPECT().
		DeleteTransferPermanent(transferID).
		Return(errors.New("failed to delete transfer"))

	mockLogger.EXPECT().
		Error("Failed to permanently delete transfer", gomock.Any())

	result, errResp := transferService.DeleteTransferPermanent(transferID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to permanently delete transfer", errResp.Message)
}
