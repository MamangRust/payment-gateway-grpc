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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestFindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	page := 1
	pageSize := 10
	search := "test"
	totalRecords := 2

	saldoRecords := []*record.SaldoRecord{
		{
			ID:             1,
			CardNumber:     "1234",
			TotalBalance:   1000,
			WithdrawAmount: 500,
			WithdrawTime:   "2024-12-24T10:00:00Z",
			CreatedAt:      "2024-12-24T10:00:00Z",
			UpdatedAt:      "2024-12-24T10:00:00Z",
			DeletedAt:      nil,
		},
		{
			ID:             2,
			CardNumber:     "5678",
			TotalBalance:   2000,
			WithdrawAmount: 700,
			WithdrawTime:   "2024-12-24T11:00:00Z",
			CreatedAt:      "2024-12-24T11:00:00Z",
			UpdatedAt:      "2024-12-24T11:00:00Z",
			DeletedAt:      nil,
		},
	}

	expectedResponse := []*response.SaldoResponse{
		{
			ID:             1,
			CardNumber:     "1234",
			TotalBalance:   1000,
			WithdrawAmount: 500,
			WithdrawTime:   "2024-12-24T10:00:00Z",
			CreatedAt:      "2024-12-24T10:00:00Z",
			UpdatedAt:      "2024-12-24T10:00:00Z",
		},
		{
			ID:             2,
			CardNumber:     "5678",
			TotalBalance:   2000,
			WithdrawAmount: 700,
			WithdrawTime:   "2024-12-24T11:00:00Z",
			CreatedAt:      "2024-12-24T11:00:00Z",
			UpdatedAt:      "2024-12-24T11:00:00Z",
		},
	}

	mock_logger.EXPECT().Debug("Fetching all saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search)).Times(1)

	mock_saldo_repo.EXPECT().FindAllSaldos(search, page, pageSize).Return(saldoRecords, totalRecords, nil).Times(1)

	mock_mapping.EXPECT().ToSaldoResponses(saldoRecords).Return(expectedResponse).Times(1)

	mock_logger.EXPECT().Debug("Successfully fetched saldo records", zap.Int("totalRecords", totalRecords), zap.Int("totalPages", 2)).Times(1)

	result, totalPages, errResp := saldoService.FindAll(page, pageSize, search)

	assert.Nil(t, errResp)
	assert.Equal(t, 2, totalPages)
	assert.Equal(t, expectedResponse, result)
}

func TestFindAll_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	page := 1
	pageSize := 10
	search := "test"

	expectedError := errors.New("database error")

	mock_logger.EXPECT().Debug("Fetching all saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search)).Times(1)

	mock_saldo_repo.EXPECT().FindAllSaldos(search, page, pageSize).Return(nil, 0, expectedError).Times(1)

	mock_logger.EXPECT().Error("Failed to fetch saldo records", zap.Error(expectedError)).Times(1)

	result, totalPages, errResp := saldoService.FindAll(page, pageSize, search)

	assert.Nil(t, result)
	assert.Equal(t, 0, totalPages)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Unable to fetch saldo records", errResp.Message)
}

func TestFindAll_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	page := 1
	pageSize := 10
	search := "not_found"

	saldoRecords := []*record.SaldoRecord{}
	totalRecords := 0

	mock_logger.EXPECT().Debug("Fetching all saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search)).Times(1)

	mock_saldo_repo.EXPECT().FindAllSaldos(search, page, pageSize).Return(saldoRecords, totalRecords, nil).Times(1)

	mock_mapping.EXPECT().ToSaldoResponses(saldoRecords).Return([]*response.SaldoResponse{}).Times(1)

	mock_logger.EXPECT().Debug("Successfully fetched saldo records", zap.Int("totalRecords", totalRecords), zap.Int("totalPages", 0)).Times(1)

	result, totalPages, errResp := saldoService.FindAll(page, pageSize, search)

	assert.Nil(t, errResp)
	assert.Equal(t, 0, totalPages)
	assert.Equal(t, []*response.SaldoResponse{}, result)
}

func TestFindById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	saldoID := 1

	saldoRecord := &record.SaldoRecord{
		ID:             1,
		CardNumber:     "1234",
		TotalBalance:   1000,
		WithdrawAmount: 500,
		WithdrawTime:   "2024-12-24T10:00:00Z",
		CreatedAt:      "2024-12-24T10:00:00Z",
		UpdatedAt:      "2024-12-24T10:00:00Z",
	}

	expectedResponse := &response.SaldoResponse{
		ID:             1,
		CardNumber:     "1234",
		TotalBalance:   1000,
		WithdrawAmount: 500,
		WithdrawTime:   "2024-12-24T10:00:00Z",
		CreatedAt:      "2024-12-24T10:00:00Z",
		UpdatedAt:      "2024-12-24T10:00:00Z",
	}

	mock_logger.EXPECT().Debug("Fetching saldo record by ID", zap.Int("saldo_id", saldoID)).Times(1)

	mock_saldo_repo.EXPECT().FindById(saldoID).Return(saldoRecord, nil).Times(1)

	mock_mapping.EXPECT().ToSaldoResponse(saldoRecord).Return(expectedResponse).Times(1)

	mock_logger.EXPECT().Debug("Successfully fetched saldo by ID", zap.Int("saldo_id", saldoID)).Times(1)

	result, errResp := saldoService.FindById(saldoID)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestFindById_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	saldoID := 1
	expectedError := errors.New("saldo not found")

	mock_logger.EXPECT().Debug("Fetching saldo record by ID", zap.Int("saldo_id", saldoID)).Times(1)

	mock_saldo_repo.EXPECT().FindById(saldoID).Return(nil, expectedError).Times(1)

	mock_logger.EXPECT().Error("Failed to fetch saldo by ID", zap.Error(expectedError), zap.Int("saldo_id", saldoID)).Times(1)

	result, errResp := saldoService.FindById(saldoID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Saldo not found for the given ID", errResp.Message)
}

func TestFindByCardNumberSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	cardNumber := "1234"

	saldoRecord := &record.SaldoRecord{
		ID:             1,
		CardNumber:     cardNumber,
		TotalBalance:   1000,
		WithdrawAmount: 500,
		WithdrawTime:   "2024-12-24T10:00:00Z",
		CreatedAt:      "2024-12-24T10:00:00Z",
		UpdatedAt:      "2024-12-24T10:00:00Z",
	}

	expectedResponse := &response.SaldoResponse{
		ID:             1,
		CardNumber:     cardNumber,
		TotalBalance:   1000,
		WithdrawAmount: 500,
		WithdrawTime:   "2024-12-24T10:00:00Z",
		CreatedAt:      "2024-12-24T10:00:00Z",
		UpdatedAt:      "2024-12-24T10:00:00Z",
	}

	mock_logger.EXPECT().Debug("Fetching saldo record by card number", zap.String("card_number", cardNumber)).Times(1)

	mock_saldo_repo.EXPECT().
		FindByCardNumber(cardNumber).
		Return(saldoRecord, nil).
		Times(1)

	mock_mapping.EXPECT().
		ToSaldoResponse(saldoRecord).
		Return(expectedResponse).
		Times(1)

	mock_logger.EXPECT().Debug("Successfully fetched saldo by card number", zap.String("card_number", cardNumber)).Times(1)

	result, errResp := saldoService.FindByCardNumber(cardNumber)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestFindByCardNumberSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	cardNumber := "1234"
	expectedError := errors.New("saldo not found for the given card number")

	mock_logger.EXPECT().Debug("Fetching saldo record by card number", zap.String("card_number", cardNumber)).Times(1)

	mock_saldo_repo.EXPECT().
		FindByCardNumber(cardNumber).
		Return(nil, expectedError).
		Times(1)

	mock_logger.EXPECT().
		Error(
			"Failed to fetch saldo by card number",
			zap.Error(expectedError),
			zap.String("card_number", cardNumber),
		).Times(1)

	result, errResp := saldoService.FindByCardNumber(cardNumber)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Saldo not found for the given card number", errResp.Message)
}

func TestFindByActiveCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	saldoRecords := []*record.SaldoRecord{
		{
			ID:             1,
			CardNumber:     "1234",
			TotalBalance:   1000,
			WithdrawAmount: 500,
			WithdrawTime:   "2024-12-24T10:00:00Z",
			CreatedAt:      "2024-12-24T10:00:00Z",
			UpdatedAt:      "2024-12-24T10:00:00Z",
		},
	}

	expectedResponses := []*response.SaldoResponseDeleteAt{
		{
			ID:             1,
			CardNumber:     "1234",
			TotalBalance:   1000,
			WithdrawAmount: 500,
			WithdrawTime:   "2024-12-24T10:00:00Z",
			CreatedAt:      "2024-12-24T10:00:00Z",
			UpdatedAt:      "2024-12-24T10:00:00Z",
		},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 1

	mock_saldo_repo.EXPECT().FindByActive(search, page, pageSize).Return(saldoRecords, expected, nil).Times(1)
	mock_mapping.EXPECT().ToSaldoResponsesDeleteAt(saldoRecords).Return(expectedResponses).Times(1)

	result, totalRecord, errResp := saldoService.FindByActive(pageSize, page, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, expectedResponses, result)
}

func TestFindByActive_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	mock_saldo_repo.EXPECT().FindByActive(search, page, pageSize).Return(nil, expected, errors.New("database error")).Times(1)

	result, totalRecord, errResp := saldoService.FindByActive(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "No active saldo records found for the given ID", errResp.Message)
}

func TestFindByActive_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	page := 1
	pageSize := 1
	search := ""
	expected := 1

	mock_saldo_repo.EXPECT().FindByActive(search, page, pageSize).Return([]*record.SaldoRecord{}, expected, nil).Times(1)
	mock_mapping.EXPECT().ToSaldoResponsesDeleteAt([]*record.SaldoRecord{}).Return([]*response.SaldoResponseDeleteAt{}).Times(1)

	result, totalRecord, errResp := saldoService.FindByActive(pageSize, page, search)

	assert.Nil(t, errResp)
	assert.Empty(t, result)
	assert.Equal(t, expected, totalRecord)
}

func TestFindByTrashed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	saldoRecords := []*record.SaldoRecord{
		{
			ID:             1,
			CardNumber:     "5678",
			TotalBalance:   2000,
			WithdrawAmount: 100,
			WithdrawTime:   "2024-12-25T12:00:00Z",
			CreatedAt:      "2024-12-25T12:00:00Z",
			UpdatedAt:      "2024-12-25T12:00:00Z",
		},
	}

	expectedResponses := []*response.SaldoResponseDeleteAt{
		{
			ID:             1,
			CardNumber:     "5678",
			TotalBalance:   2000,
			WithdrawAmount: 100,
			WithdrawTime:   "2024-12-25T12:00:00Z",
			CreatedAt:      "2024-12-25T12:00:00Z",
			UpdatedAt:      "2024-12-25T12:00:00Z",
		},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 1

	mock_logger.EXPECT().Info("Fetching trashed saldo records").Times(1)
	mock_saldo_repo.EXPECT().FindByTrashed(search, page, pageSize).Return(saldoRecords, expected, nil).Times(1)
	mock_mapping.EXPECT().ToSaldoResponsesDeleteAt(saldoRecords).Return(expectedResponses).Times(1)
	mock_logger.EXPECT().Debug("Successfully fetched trashed saldo records", zap.Int("record_count", len(saldoRecords))).Times(1)

	result, totalRecord, errResp := saldoService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, expectedResponses, result)
}

func TestFindByTrashed_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	expectedError := errors.New("database error")

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	mock_logger.EXPECT().Info("Fetching trashed saldo records").Times(1)
	mock_saldo_repo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, expected, expectedError).Times(1)
	mock_logger.EXPECT().Error("Failed to fetch trashed saldo records", zap.Error(expectedError)).Times(1)

	result, totalRecord, errResp := saldoService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "No trashed saldo records found", errResp.Message)
}

func TestFindByTrashed_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)

	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	mock_logger.EXPECT().Info("Fetching trashed saldo records").Times(1)
	mock_saldo_repo.EXPECT().FindByTrashed(search, page, pageSize).Return([]*record.SaldoRecord{}, expected, nil).Times(1)
	mock_mapping.EXPECT().ToSaldoResponsesDeleteAt([]*record.SaldoRecord{}).Return([]*response.SaldoResponseDeleteAt{}).Times(1)
	mock_logger.EXPECT().Debug("Successfully fetched trashed saldo records", zap.Int("record_count", 0)).Times(1)

	result, totalRecord, errResp := saldoService.FindByTrashed(pageSize, page, search)

	assert.Equal(t, expected, totalRecord)
	assert.Nil(t, errResp)
	assert.Empty(t, result)
}

func TestCreateSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	cardNumber := "1234"

	request := &requests.CreateSaldoRequest{
		CardNumber:   cardNumber,
		TotalBalance: 1000,
	}

	cardRecord := &record.CardRecord{
		ID:           1,
		UserID:       1,
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
	}

	saldoRecord := &record.SaldoRecord{
		ID:           1,
		CardNumber:   cardNumber,
		TotalBalance: 1000,
		CreatedAt:    "2024-12-25T12:00:00Z",
		UpdatedAt:    "2024-12-25T12:00:00Z",
	}

	expectedResponse := &response.SaldoResponse{
		ID:           1,
		CardNumber:   cardNumber,
		TotalBalance: 1000,
		CreatedAt:    "2024-12-25T12:00:00Z",
		UpdatedAt:    "2024-12-25T12:00:00Z",
	}

	mock_logger.EXPECT().Debug("Creating saldo record", zap.String("card_number", request.CardNumber)).Times(1)
	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(cardRecord, nil).Times(1)
	mock_saldo_repo.EXPECT().CreateSaldo(request).Return(saldoRecord, nil).Times(1)
	mock_mapping.EXPECT().ToSaldoResponse(saldoRecord).Return(expectedResponse).Times(1)
	mock_logger.EXPECT().Debug("Successfully created saldo record", zap.String("card_number", request.CardNumber)).Times(1)

	result, errResp := saldoService.CreateSaldo(request)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestCreateSaldo_Failure_CardNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	cardNumber := "5678"

	request := &requests.CreateSaldoRequest{
		CardNumber:   cardNumber,
		TotalBalance: 1000,
	}

	mock_logger.EXPECT().Debug("Creating saldo record", zap.String("card_number", request.CardNumber)).Times(1)
	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(nil, errors.New("card not found")).Times(1)
	mock_logger.EXPECT().Error("Card not found for creating saldo", zap.Error(errors.New("card not found")), zap.String("card_number", request.CardNumber)).Times(1)

	result, errResp := saldoService.CreateSaldo(request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Card number not found", errResp.Message)
}

func TestCreateSaldo_Failure_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	cardNumber := "1234"

	request := &requests.CreateSaldoRequest{
		CardNumber:   cardNumber,
		TotalBalance: 1000,
	}

	cardRecord := &record.CardRecord{
		ID:           2,
		UserID:       2,
		CardType:     "Credit",
		ExpireDate:   "2024-06-30",
		CVV:          "456",
		CardProvider: "MasterCard",
	}

	mock_logger.EXPECT().Debug("Creating saldo record", zap.String("card_number", request.CardNumber)).Times(1)
	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(cardRecord, nil).Times(1)
	mock_saldo_repo.EXPECT().CreateSaldo(request).Return(nil, errors.New("database error")).Times(1)
	mock_logger.EXPECT().Error("Failed to create saldo", zap.Error(errors.New("database error"))).Times(1)

	result, errResp := saldoService.CreateSaldo(request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to create saldo record", errResp.Message)
}

func TestUpdateSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	cardNumber := "1234"

	request := &requests.UpdateSaldoRequest{
		SaldoID:      1,
		CardNumber:   cardNumber,
		TotalBalance: 2000,
	}

	cardRecord := &record.CardRecord{
		ID:           1,
		UserID:       1,
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
	}

	updatedSaldoRecord := &record.SaldoRecord{
		ID:           1,
		CardNumber:   cardNumber,
		TotalBalance: 2000,
		CreatedAt:    "2024-12-25T10:00:00Z",
		UpdatedAt:    "2024-12-25T12:00:00Z",
	}

	expectedResponse := &response.SaldoResponse{
		ID:           1,
		CardNumber:   cardNumber,
		TotalBalance: 2000,
		CreatedAt:    "2024-12-25T10:00:00Z",
		UpdatedAt:    "2024-12-25T12:00:00Z",
	}

	mock_logger.EXPECT().Debug("Updating saldo record", zap.String("card_number", request.CardNumber), zap.Float64("amount", float64(request.TotalBalance))).Times(1)
	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(cardRecord, nil).Times(1)
	mock_saldo_repo.EXPECT().UpdateSaldo(request).Return(updatedSaldoRecord, nil).Times(1)
	mock_mapping.EXPECT().ToSaldoResponse(updatedSaldoRecord).Return(expectedResponse).Times(1)
	mock_logger.EXPECT().Debug("Successfully updated saldo", zap.String("card_number", request.CardNumber), zap.Int("saldo_id", updatedSaldoRecord.ID)).Times(1)

	result, errResp := saldoService.UpdateSaldo(request)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestUpdateSaldo_Failure_CardNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	cardNumber := "5678"

	request := &requests.UpdateSaldoRequest{
		SaldoID:      1,
		CardNumber:   cardNumber,
		TotalBalance: 2000,
	}

	mock_logger.EXPECT().Debug("Updating saldo record", zap.String("card_number", request.CardNumber), zap.Float64("amount", float64(request.TotalBalance))).Times(1)
	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(nil, errors.New("card not found")).Times(1)
	mock_logger.EXPECT().Error("Failed to find card by card number", zap.Error(errors.New("card not found")), zap.String("card_number", request.CardNumber)).Times(1)

	result, errResp := saldoService.UpdateSaldo(request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Card number not found", errResp.Message)
}

func TestUpdateSaldo_Failure_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	cardNumber := "1234"

	request := &requests.UpdateSaldoRequest{
		SaldoID:      1,
		CardNumber:   cardNumber,
		TotalBalance: 2000,
	}

	cardRecord := &record.CardRecord{
		ID:         1,
		CardNumber: cardNumber,
		UserID:     1,
	}

	mock_logger.EXPECT().Debug("Updating saldo record", zap.String("card_number", request.CardNumber), zap.Float64("amount", float64(request.TotalBalance))).Times(1)
	mock_card_repo.EXPECT().FindCardByCardNumber(request.CardNumber).Return(cardRecord, nil).Times(1)
	mock_saldo_repo.EXPECT().UpdateSaldo(request).Return(nil, errors.New("database error")).Times(1)
	mock_logger.EXPECT().Error("Failed to update saldo", zap.Error(errors.New("database error")), zap.String("card_number", request.CardNumber)).Times(1)

	result, errResp := saldoService.UpdateSaldo(request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to update saldo", errResp.Message)
}

func TestTrashSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	saldoID := 1

	trashedSaldoRecord := &record.SaldoRecord{
		ID:           saldoID,
		CardNumber:   "1234",
		TotalBalance: 1000,
		CreatedAt:    "2024-12-24T10:00:00Z",
		UpdatedAt:    "2024-12-25T12:00:00Z",
	}

	expectedResponse := &response.SaldoResponse{
		ID:           saldoID,
		CardNumber:   "1234",
		TotalBalance: 1000,
		CreatedAt:    "2024-12-24T10:00:00Z",
		UpdatedAt:    "2024-12-25T12:00:00Z",
	}

	mock_logger.EXPECT().Debug("Trashing saldo record", zap.Int("saldo_id", saldoID)).Times(1)
	mock_saldo_repo.EXPECT().TrashedSaldo(saldoID).Return(trashedSaldoRecord, nil).Times(1)
	mock_mapping.EXPECT().ToSaldoResponse(trashedSaldoRecord).Return(expectedResponse).Times(1)
	mock_logger.EXPECT().Debug("Successfully trashed saldo", zap.Int("saldo_id", saldoID)).Times(1)

	result, errResp := saldoService.TrashSaldo(saldoID)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestRestoreSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	saldoID := 1

	restoredSaldoRecord := &record.SaldoRecord{
		ID:           saldoID,
		CardNumber:   "1234",
		TotalBalance: 1000,
		CreatedAt:    "2024-12-24T10:00:00Z",
		UpdatedAt:    "2024-12-25T12:00:00Z",
	}

	expectedResponse := &response.SaldoResponse{
		ID:           saldoID,
		CardNumber:   "1234",
		TotalBalance: 1000,
		CreatedAt:    "2024-12-24T10:00:00Z",
		UpdatedAt:    "2024-12-25T12:00:00Z",
	}

	mock_logger.EXPECT().Debug("Restoring saldo record from trash", zap.Int("saldo_id", saldoID)).Times(1)
	mock_saldo_repo.EXPECT().RestoreSaldo(saldoID).Return(restoredSaldoRecord, nil).Times(1)
	mock_mapping.EXPECT().ToSaldoResponse(restoredSaldoRecord).Return(expectedResponse).Times(1)
	mock_logger.EXPECT().Debug("Successfully restored saldo", zap.Int("saldo_id", saldoID)).Times(1)

	result, errResp := saldoService.RestoreSaldo(saldoID)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestRestoreSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	saldoID := 1

	mock_logger.EXPECT().Debug("Restoring saldo record from trash", zap.Int("saldo_id", saldoID)).Times(1)
	mock_saldo_repo.EXPECT().RestoreSaldo(saldoID).Return(nil, errors.New("database error")).Times(1)
	mock_logger.EXPECT().Error("Failed to restore saldo", zap.Error(errors.New("database error")), zap.Int("saldo_id", saldoID)).Times(1)

	result, errResp := saldoService.RestoreSaldo(saldoID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to restore saldo from trash", errResp.Message)
}

func TestDeleteSaldoPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	saldoID := 1

	mock_logger.EXPECT().Debug("Deleting saldo permanently", zap.Int("saldo_id", saldoID)).Times(1)
	mock_saldo_repo.EXPECT().DeleteSaldoPermanent(saldoID).Return(nil).Times(1)
	mock_logger.EXPECT().Debug("Successfully deleted saldo permanently", zap.Int("saldo_id", saldoID)).Times(1)

	result, errResp := saldoService.DeleteSaldoPermanent(saldoID)

	assert.Nil(t, result)
	assert.Nil(t, errResp)
}

func TestDeleteSaldoPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_saldo_repo := mock_repository.NewMockSaldoRepository(ctrl)
	mock_mapping := mock_responsemapper.NewMockSaldoResponseMapper(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	saldoService := service.NewSaldoService(mock_saldo_repo, mock_card_repo, mock_logger, mock_mapping)

	saldoID := 1

	mock_logger.EXPECT().Debug("Deleting saldo permanently", zap.Int("saldo_id", saldoID)).Times(1)
	mock_saldo_repo.EXPECT().DeleteSaldoPermanent(saldoID).Return(errors.New("database error")).Times(1)
	mock_logger.EXPECT().Error("Failed to delete saldo permanently", zap.Error(errors.New("database error")), zap.Int("saldo_id", saldoID)).Times(1)

	result, errResp := saldoService.DeleteSaldoPermanent(saldoID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to delete saldo permanently", errResp.Message)
}
