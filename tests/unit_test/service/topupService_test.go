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

func TestFindAllTopups_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	page := 1
	pageSize := 10
	search := "test"
	totalRecords := 3

	topups := []*record.TopupRecord{
		{
			ID:          1,
			CardNumber:  "1234",
			TopupNo:     "TOPUP-001",
			TopupAmount: 50000,
			TopupMethod: "bank_transfer",
			TopupTime:   "2024-12-25T09:00:00Z",
			CreatedAt:   "2024-12-25T09:00:00Z",
			UpdatedAt:   "2024-12-25T09:30:00Z",
			DeletedAt:   nil,
		},
		{
			ID:          2,
			CardNumber:  "5678",
			TopupNo:     "TOPUP-002",
			TopupAmount: 75000,
			TopupMethod: "credit_card",
			TopupTime:   "2024-12-25T10:00:00Z",
			CreatedAt:   "2024-12-25T10:00:00Z",
			UpdatedAt:   "2024-12-25T10:30:00Z",
			DeletedAt:   nil,
		},
	}

	expectedResponses := []*response.TopupResponse{
		{
			ID:          1,
			CardNumber:  "1234",
			TopupNo:     "TOPUP-001",
			TopupAmount: 50000,
			TopupMethod: "bank_transfer",
			TopupTime:   "2024-12-25T09:00:00Z",
			CreatedAt:   "2024-12-25T09:00:00Z",
			UpdatedAt:   "2024-12-25T09:30:00Z",
		},
		{
			ID:          2,
			CardNumber:  "5678",
			TopupNo:     "TOPUP-002",
			TopupAmount: 75000,
			TopupMethod: "credit_card",
			TopupTime:   "2024-12-25T10:00:00Z",
			CreatedAt:   "2024-12-25T10:00:00Z",
			UpdatedAt:   "2024-12-25T10:30:00Z",
		},
	}

	mock_topup_repo.EXPECT().FindAllTopups(search, page, pageSize).Return(topups, totalRecords, nil).Times(1)
	mock_mapping.EXPECT().ToTopupResponses(topups).Return(expectedResponses).Times(1)

	results, totalPages, errResp := topupService.FindAll(page, pageSize, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponses, results)
	assert.Equal(t, 3, totalPages) // Total pages = ceil(25 / 10) = 3
}

func TestFindAllTopups_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	page := 1
	pageSize := 10
	search := "test"

	mock_topup_repo.EXPECT().FindAllTopups(search, page, pageSize).Return(nil, 0, errors.New("database error")).Times(1)
	mock_logger.EXPECT().Error("failed to fetch topups", zap.Error(errors.New("database error"))).Times(1)

	results, totalPages, errResp := topupService.FindAll(page, pageSize, search)

	assert.Nil(t, results)
	assert.Equal(t, 0, totalPages)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to fetch topups", errResp.Message)
}

func TestFindAllTopups_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	page := 1
	pageSize := 10
	search := "no-records"
	totalRecords := 0

	mock_topup_repo.EXPECT().FindAllTopups(search, page, pageSize).Return(nil, totalRecords, nil).Times(1)
	mock_mapping.EXPECT().ToTopupResponses([]*record.TopupRecord(nil)).Return([]*response.TopupResponse(nil)).Times(1)

	results, totalPages, errResp := topupService.FindAll(page, pageSize, search)

	assert.Nil(t, errResp)
	assert.Empty(t, results)
	assert.Equal(t, 0, totalPages)
}

func TestFindByIdTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	topupID := 1
	topupRecord := &record.TopupRecord{
		ID:          1,
		CardNumber:  "1234",
		TopupNo:     "TOPUP-001",
		TopupAmount: 50000,
		TopupMethod: "bank_transfer",
		TopupTime:   "2024-12-25T09:00:00Z",
		CreatedAt:   "2024-12-25T09:00:00Z",
		UpdatedAt:   "2024-12-25T09:30:00Z",
		DeletedAt:   nil,
	}

	expectedResponse := &response.TopupResponse{
		ID:          1,
		CardNumber:  "1234",
		TopupNo:     "TOPUP-001",
		TopupAmount: 50000,
		TopupMethod: "bank_transfer",
		TopupTime:   "2024-12-25T09:00:00Z",
		CreatedAt:   "2024-12-25T09:00:00Z",
		UpdatedAt:   "2024-12-25T09:30:00Z",
	}

	mock_topup_repo.EXPECT().FindById(topupID).Return(topupRecord, nil).Times(1)

	mock_mapping.EXPECT().ToTopupResponse(topupRecord).Return(expectedResponse).Times(1)

	result, errResp := topupService.FindById(topupID)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestFindByIdTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	topupID := 1
	expectedError := errors.New("Topup not found")

	mock_topup_repo.EXPECT().FindById(topupID).Return(nil, expectedError).Times(1)

	mock_logger.EXPECT().Error("failed to find topup by id", zap.Error(expectedError)).Times(1)

	result, errResp := topupService.FindById(topupID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Topup record not found", errResp.Message)
}

func TestFindByCardNumberTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	cardNumber := "1234"

	topups := []*record.TopupRecord{
		{
			ID:          1,
			CardNumber:  "1234",
			TopupNo:     "TOPUP-001",
			TopupAmount: 50000,
			TopupMethod: "bank_transfer",
			TopupTime:   "2024-12-25T09:00:00Z",
			CreatedAt:   "2024-12-25T09:00:00Z",
			UpdatedAt:   "2024-12-25T09:30:00Z",
			DeletedAt:   nil,
		},
		{
			ID:          2,
			CardNumber:  "5678",
			TopupNo:     "TOPUP-002",
			TopupAmount: 75000,
			TopupMethod: "credit_card",
			TopupTime:   "2024-12-25T10:00:00Z",
			CreatedAt:   "2024-12-25T10:00:00Z",
			UpdatedAt:   "2024-12-25T10:30:00Z",
			DeletedAt:   nil,
		},
	}

	expectedResponse := []*response.TopupResponse{
		{
			ID:          1,
			CardNumber:  "1234",
			TopupNo:     "TOPUP-001",
			TopupAmount: 50000,
			TopupMethod: "bank_transfer",
			TopupTime:   "2024-12-25T09:00:00Z",
			CreatedAt:   "2024-12-25T09:00:00Z",
			UpdatedAt:   "2024-12-25T09:30:00Z",
		},
		{
			ID:          2,
			CardNumber:  "5678",
			TopupNo:     "TOPUP-002",
			TopupAmount: 75000,
			TopupMethod: "credit_card",
			TopupTime:   "2024-12-25T10:00:00Z",
			CreatedAt:   "2024-12-25T10:00:00Z",
			UpdatedAt:   "2024-12-25T10:30:00Z",
		},
	}

	mock_logger.EXPECT().Debug("Finding top-up by card number", zap.String("card_number", cardNumber))

	mock_topup_repo.EXPECT().FindByCardNumber(cardNumber).Return(topups, nil).Times(1)

	mock_mapping.EXPECT().ToTopupResponses(topups).Return(expectedResponse).Times(1)

	mock_logger.EXPECT().Debug("Successfully found top-up by card number", zap.String("card_number", cardNumber))

	result, errResp := topupService.FindByCardNumber(cardNumber)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestFindByCardNumberTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	cardNumber := "1234"
	expectedError := errors.New("Topup record not found")

	mock_logger.EXPECT().Debug("Finding top-up by card number", zap.String("card_number", cardNumber))

	mock_topup_repo.EXPECT().FindByCardNumber(cardNumber).Return(nil, expectedError).Times(1)

	mock_logger.EXPECT().Error("Failed to find top-up by card number", zap.Error(expectedError), zap.String("card_number", cardNumber)).Times(1)

	result, errResp := topupService.FindByCardNumber(cardNumber)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to find top-up by card number", errResp.Message)
}

func TestFindByActiveTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	topups := []*record.TopupRecord{
		{
			ID:          1,
			CardNumber:  "1234",
			TopupNo:     "TOPUP-001",
			TopupAmount: 50000,
			TopupMethod: "bank_transfer",
			TopupTime:   "2024-12-25T09:00:00Z",
			CreatedAt:   "2024-12-25T09:00:00Z",
			UpdatedAt:   "2024-12-25T09:30:00Z",
			DeletedAt:   nil,
		},
		{
			ID:          2,
			CardNumber:  "5678",
			TopupNo:     "TOPUP-002",
			TopupAmount: 75000,
			TopupMethod: "credit_card",
			TopupTime:   "2024-12-25T10:00:00Z",
			CreatedAt:   "2024-12-25T10:00:00Z",
			UpdatedAt:   "2024-12-25T10:30:00Z",
			DeletedAt:   nil,
		},
	}

	expectedResponse := []*response.TopupResponseDeleteAt{
		{
			ID:          1,
			CardNumber:  "1234",
			TopupNo:     "TOPUP-001",
			TopupAmount: 50000,
			TopupMethod: "bank_transfer",
			TopupTime:   "2024-12-25T09:00:00Z",
			CreatedAt:   "2024-12-25T09:00:00Z",
			UpdatedAt:   "2024-12-25T09:30:00Z",
		},
		{
			ID:          2,
			CardNumber:  "5678",
			TopupNo:     "TOPUP-002",
			TopupAmount: 75000,
			TopupMethod: "credit_card",
			TopupTime:   "2024-12-25T10:00:00Z",
			CreatedAt:   "2024-12-25T10:00:00Z",
			UpdatedAt:   "2024-12-25T10:30:00Z",
		},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 2

	mock_logger.EXPECT().Info("Finding active top-up records").Times(1)

	mock_topup_repo.EXPECT().FindByActive(search, page, pageSize).Return(topups, expected, nil).Times(1)

	mock_mapping.EXPECT().ToTopupResponsesDeleteAt(topups).Return(expectedResponse).Times(1)

	mock_logger.EXPECT().Debug("Successfully found active top-up records", zap.Int("count", len(
		expectedResponse,
	)))

	result, totalRecord, errResp := topupService.FindByActive(pageSize, page, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, expectedResponse, result)
}

func TestFindByActiveTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	expectedError := errors.New("Failed to fetch active top-up records")

	mock_logger.EXPECT().Info("Finding active top-up records").Times(1)

	mock_topup_repo.EXPECT().FindByActive(search, page, pageSize).Return(nil, expected, expectedError).Times(1)

	mock_logger.EXPECT().Error("Failed to find active top-up records", zap.Error(expectedError)).Times(1)

	result, totalRecord, errResp := topupService.FindByActive(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to find active top-up records", errResp.Message)
}

func TestFindByTrashedTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	topups := []*record.TopupRecord{
		{
			ID:          1,
			CardNumber:  "1234",
			TopupNo:     "TOPUP-001",
			TopupAmount: 50000,
			TopupMethod: "bank_transfer",
			TopupTime:   "2024-12-25T09:00:00Z",
			CreatedAt:   "2024-12-25T09:00:00Z",
			UpdatedAt:   "2024-12-25T09:30:00Z",
			DeletedAt:   nil,
		},
		{
			ID:          2,
			CardNumber:  "5678",
			TopupNo:     "TOPUP-002",
			TopupAmount: 75000,
			TopupMethod: "credit_card",
			TopupTime:   "2024-12-25T10:00:00Z",
			CreatedAt:   "2024-12-25T10:00:00Z",
			UpdatedAt:   "2024-12-25T10:30:00Z",
			DeletedAt:   nil,
		},
	}

	expectedResponse := []*response.TopupResponseDeleteAt{
		{
			ID:          1,
			CardNumber:  "1234",
			TopupNo:     "TOPUP-001",
			TopupAmount: 50000,
			TopupMethod: "bank_transfer",
			TopupTime:   "2024-12-25T09:00:00Z",
			CreatedAt:   "2024-12-25T09:00:00Z",
			UpdatedAt:   "2024-12-25T09:30:00Z",
		},
		{
			ID:          2,
			CardNumber:  "5678",
			TopupNo:     "TOPUP-002",
			TopupAmount: 75000,
			TopupMethod: "credit_card",
			TopupTime:   "2024-12-25T10:00:00Z",
			CreatedAt:   "2024-12-25T10:00:00Z",
			UpdatedAt:   "2024-12-25T10:30:00Z",
		},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 2

	mock_logger.EXPECT().Info("Finding trashed top-up records").Times(1)

	mock_topup_repo.EXPECT().FindByTrashed(search, page, pageSize).Return(topups, expected, nil).Times(1)

	mock_mapping.EXPECT().ToTopupResponsesDeleteAt(topups).Return(expectedResponse).Times(1)

	mock_logger.EXPECT().Debug("Successfully found trashed top-up records", zap.Int("count", len(topups))).Times(1)

	result, totalRecord, errResp := topupService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, expectedResponse, result)
}

func TestFindByTrashedTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	expectedError := errors.New("Failed to fetch trashed top-up records")

	mock_logger.EXPECT().Info("Finding trashed top-up records").Times(1)

	mock_topup_repo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, expected, expectedError).Times(1)

	mock_logger.EXPECT().Error("Failed to find trashed top-up records", zap.Error(expectedError)).Times(1)

	result, totalRecord, errResp := topupService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to find trashed top-up records", errResp.Message)
}

func TestCreateTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	request := &requests.CreateTopupRequest{
		CardNumber:  "1234",
		TopupAmount: 50000,
		TopupMethod: "bank_transfer",
	}

	card := &record.CardRecord{
		ID:           1,
		UserID:       1,
		CardNumber:   "1234",
		ExpireDate:   "2024-12-31",
		CardType:     "credit",
		CVV:          "123",
		CardProvider: "Visa",
	}
	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(card, nil).Times(1)

	topup := &record.TopupRecord{
		ID:          1,
		CardNumber:  "1234",
		TopupNo:     "TOPUP-001",
		TopupAmount: 50000,
		TopupMethod: "bank_transfer",
		TopupTime:   "2024-12-25T09:00:00Z",
	}
	mock_topup_repo.EXPECT().CreateTopup(request).Return(topup, nil).Times(1)

	saldo := &record.SaldoRecord{
		CardNumber:   "1234",
		TotalBalance: 100000,
	}
	mock_saldo_repo.EXPECT().FindByCardNumber(request.CardNumber).Return(saldo, nil).Times(1)

	mock_saldo_repo.EXPECT().UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: saldo.TotalBalance + request.TopupAmount,
	}).Return(nil, nil).Times(1)

	expireDate, _ := time.Parse("2006-01-02", card.ExpireDate)
	mock_card_repo.EXPECT().UpdateCard(&requests.UpdateCardRequest{
		CardID:       card.ID,
		UserID:       card.UserID,
		CardType:     card.CardType,
		ExpireDate:   expireDate,
		CVV:          card.CVV,
		CardProvider: card.CardProvider,
	}).Return(nil, nil).Times(1)

	expectedResponse := &response.TopupResponse{
		ID:          1,
		CardNumber:  "1234",
		TopupNo:     "TOPUP-001",
		TopupAmount: 50000,
		TopupMethod: "bank_transfer",
		TopupTime:   "2024-12-25T09:00:00Z",
		CreatedAt:   "2024-12-25T09:00:00Z",
		UpdatedAt:   "2024-12-25T09:30:00Z",
	}
	mock_mapping.EXPECT().ToTopupResponse(topup).Return(expectedResponse).Times(1)

	result, errResp := topupService.CreateTopup(request)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestCreateTopup_Failure_CardNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	request := &requests.CreateTopupRequest{
		CardNumber:  "1234",
		TopupAmount: 50000,
		TopupMethod: "bank_transfer",
	}

	mock_logger.EXPECT().Error("failed to find card by number", zap.Error(errors.New("card not found")))

	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(nil, errors.New("card not found")).Times(1)

	result, errResp := topupService.CreateTopup(request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Card not found", errResp.Message)
}

func TestCreateTopup_Failure_TopupCreationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	request := &requests.CreateTopupRequest{
		CardNumber:  "1234",
		TopupAmount: 50000,
		TopupMethod: "bank_transfer",
	}

	card := &record.CardRecord{
		ID:           1,
		UserID:       1,
		CardNumber:   "1234",
		ExpireDate:   "2024-12-31",
		CardType:     "credit",
		CVV:          "123",
		CardProvider: "Visa",
	}

	mock_logger.EXPECT().Error("failed to create topup", zap.Error(errors.New("failed to create topup")))
	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(card, nil).Times(1)

	mock_topup_repo.EXPECT().CreateTopup(request).Return(nil, errors.New("failed to create topup")).Times(1)

	result, errResp := topupService.CreateTopup(request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to create topup record", errResp.Message)
}

func TestCreateTopup_Failure_SaldoUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(
		mock_card_repo,
		mock_topup_repo,
		mock_saldo_repo,
		mock_logger, mock_mapping)

	request := &requests.CreateTopupRequest{
		CardNumber:  "1234",
		TopupAmount: 50000,
		TopupMethod: "bank_transfer",
	}

	card := &record.CardRecord{
		ID:           1,
		UserID:       1,
		CardNumber:   "1234",
		ExpireDate:   "2024-12-31",
		CardType:     "credit",
		CVV:          "123",
		CardProvider: "Visa",
	}
	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(card, nil).Times(1)

	topup := &record.TopupRecord{
		ID:          1,
		CardNumber:  "1234",
		TopupNo:     "TOPUP-001",
		TopupAmount: 50000,
		TopupMethod: "bank_transfer",
		TopupTime:   "2024-12-25T09:00:00Z",
	}

	mock_topup_repo.EXPECT().CreateTopup(request).Return(topup, nil).Times(1)

	saldo := &record.SaldoRecord{
		CardNumber:   "1234",
		TotalBalance: 100000,
	}
	mock_saldo_repo.EXPECT().FindByCardNumber(request.CardNumber).Return(saldo, nil).Times(1)

	mock_logger.EXPECT().
		Error("failed to update saldo balance", zap.Error(errors.New("failed to update saldo"))).
		Times(1)

	mock_saldo_repo.EXPECT().UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: saldo.TotalBalance + request.TopupAmount,
	}).Return(nil, errors.New("failed to update saldo")).Times(1)

	result, errResp := topupService.CreateTopup(request)
	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to update saldo balance", errResp.Message)
}

func TestUpdateTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(mock_card_repo, mock_topup_repo, mock_saldo_repo, mock_logger, mock_mapping)

	request := &requests.UpdateTopupRequest{
		TopupID:     1,
		CardNumber:  "1234",
		TopupAmount: 150000,
	}

	card := &record.CardRecord{
		ID:           1,
		UserID:       1,
		CardNumber:   "1234",
		ExpireDate:   "2024-12-31",
		CardType:     "credit",
		CVV:          "123",
		CardProvider: "Visa",
	}
	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(card, nil).Times(1)

	existingTopup := &record.TopupRecord{
		ID:          1,
		CardNumber:  "1234",
		TopupAmount: 100000,
	}
	mock_topup_repo.EXPECT().FindById(request.TopupID).Return(existingTopup, nil).Times(1)

	mock_topup_repo.EXPECT().UpdateTopup(request).Return(existingTopup, nil).Times(1)

	currentSaldo := &record.SaldoRecord{
		CardNumber:   "1234",
		TotalBalance: 200000,
	}
	mock_saldo_repo.EXPECT().FindByCardNumber(request.CardNumber).Return(currentSaldo, nil).Times(1)

	mock_saldo_repo.EXPECT().UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   "1234",
		TotalBalance: 250000,
	}).Return(nil, nil).Times(1)

	updatedTopup := &record.TopupRecord{
		ID:          1,
		CardNumber:  "1234",
		TopupAmount: 150000,
	}
	mock_topup_repo.EXPECT().FindById(request.TopupID).Return(updatedTopup, nil).Times(1)

	expectedResponse := &response.TopupResponse{
		ID:          1,
		CardNumber:  "1234",
		TopupAmount: 150000,
	}
	mock_mapping.EXPECT().ToTopupResponse(updatedTopup).Return(expectedResponse).Times(1)

	result, errResp := topupService.UpdateTopup(request)

	assert.NotNil(t, result)
	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestUpdateTopup_Failure_TopupNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	topupService := service.NewTopupService(mock_card_repo, mock_topup_repo, nil, mock_logger, nil)

	request := &requests.UpdateTopupRequest{
		TopupID:     1,
		CardNumber:  "1234",
		TopupAmount: 150000,
	}

	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(&record.CardRecord{}, nil).Times(1)
	mock_topup_repo.EXPECT().FindById(request.TopupID).Return(nil, errors.New("topup not found")).Times(1)
	mock_logger.EXPECT().Error("Failed to find topup by ID", zap.Error(errors.New("topup not found"))).Times(1)

	result, errResp := topupService.UpdateTopup(request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Topup not found", errResp.Message)
}

func TestUpdateTopup_Failure_SaldoUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	topupService := service.NewTopupService(mock_card_repo, mock_topup_repo, mock_saldo_repo, mock_logger, nil)

	request := &requests.UpdateTopupRequest{
		TopupID:     1,
		CardNumber:  "1234",
		TopupAmount: 150000,
	}

	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(&record.CardRecord{}, nil).Times(1)

	existingTopup := &record.TopupRecord{
		ID:          1,
		CardNumber:  "1234",
		TopupAmount: 100000,
	}
	mock_topup_repo.EXPECT().FindById(request.TopupID).Return(existingTopup, nil).Times(1)

	mock_topup_repo.EXPECT().UpdateTopup(request).Return(existingTopup, nil).Times(1)

	currentSaldo := &record.SaldoRecord{
		CardNumber:   "1234",
		TotalBalance: 200000,
	}
	mock_saldo_repo.EXPECT().FindByCardNumber(request.CardNumber).Return(currentSaldo, nil).Times(1)

	mock_saldo_repo.EXPECT().UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   "1234",
		TotalBalance: 250000,
	}).Return(nil, errors.New("failed to update saldo")).Times(1)

	mock_logger.EXPECT().Error("Failed to update saldo balance", zap.Error(errors.New("failed to update saldo"))).Times(1)

	mock_topup_repo.EXPECT().UpdateTopupAmount(&requests.UpdateTopupAmount{
		TopupID:     1,
		TopupAmount: 100000,
	}).Return(nil, nil).Times(1)

	result, errResp := topupService.UpdateTopup(request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to update saldo balance: failed to update saldo", errResp.Message)
}

func TestTrashedTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(nil, mock_topup_repo, nil, mock_logger, mock_mapping)

	topupID := 1
	topupRecord := &record.TopupRecord{
		ID:          topupID,
		CardNumber:  "1234",
		TopupAmount: 100000,
		TopupNo:     "TOPUP-001",
	}

	mock_topup_repo.EXPECT().TrashedTopup(topupID).Return(topupRecord, nil).Times(1)

	expectedResponse := &response.TopupResponse{
		ID:          topupID,
		CardNumber:  "1234",
		TopupAmount: 100000,
	}
	mock_mapping.EXPECT().ToTopupResponse(topupRecord).Return(expectedResponse).Times(1)

	result, errResp := topupService.TrashedTopup(topupID)

	assert.NotNil(t, result)
	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestTrashedTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	topupService := service.NewTopupService(nil, mock_topup_repo, nil, mock_logger, nil)

	topupID := 1
	mock_topup_repo.EXPECT().TrashedTopup(topupID).Return(nil, errors.New("failed to trash topup")).Times(1)
	mock_logger.EXPECT().Error("Failed to trash topup", zap.Error(errors.New("failed to trash topup"))).Times(1)

	result, errResp := topupService.TrashedTopup(topupID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to trash topup: failed to trash topup", errResp.Message)
}

func TestRestoreTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTopupResponseMapper(ctrl)

	topupService := service.NewTopupService(nil, mock_topup_repo, nil, mock_logger, mock_mapping)

	topupID := 1
	topupRecord := &record.TopupRecord{
		ID:          topupID,
		CardNumber:  "1234",
		TopupAmount: 100000,
		TopupNo:     "TOPUP-001",
	}

	mock_topup_repo.EXPECT().RestoreTopup(topupID).Return(topupRecord, nil).Times(1)

	expectedResponse := &response.TopupResponse{
		ID:          topupID,
		CardNumber:  "1234",
		TopupAmount: 100000,
	}
	mock_mapping.EXPECT().ToTopupResponse(topupRecord).Return(expectedResponse).Times(1)

	result, errResp := topupService.RestoreTopup(topupID)

	assert.NotNil(t, result)
	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestRestoreTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	topupService := service.NewTopupService(nil, mock_topup_repo, nil, mock_logger, nil)

	topupID := 1
	mock_topup_repo.EXPECT().RestoreTopup(topupID).Return(nil, errors.New("failed to restore topup")).Times(1)
	mock_logger.EXPECT().Error("Failed to restore topup", zap.Error(errors.New("failed to restore topup"))).Times(1)

	result, errResp := topupService.RestoreTopup(topupID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to restore topup: failed to restore topup", errResp.Message)
}

func TestDeleteTopupPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	topupService := service.NewTopupService(nil, mock_topup_repo, nil, mock_logger, nil)

	topupID := 1
	mock_topup_repo.EXPECT().DeleteTopupPermanent(topupID).Return(nil).Times(1)

	result, errResp := topupService.DeleteTopupPermanent(topupID)

	assert.Nil(t, result)
	assert.Nil(t, errResp)
}

func TestDeleteTopupPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_topup_repo := mock_repository.NewMockTopupRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	topupService := service.NewTopupService(nil, mock_topup_repo, nil, mock_logger, nil)

	topupID := 1
	mock_topup_repo.EXPECT().DeleteTopupPermanent(topupID).Return(errors.New("failed to delete topup permanently")).Times(1)
	mock_logger.EXPECT().Error("Failed to delete topup permanently", zap.Error(errors.New("failed to delete topup permanently"))).Times(1)

	result, errResp := topupService.DeleteTopupPermanent(topupID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to delete topup permanently: failed to delete topup permanently", errResp.Message)
}
