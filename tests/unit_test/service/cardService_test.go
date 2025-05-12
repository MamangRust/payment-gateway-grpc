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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestFindAllCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_user_repo := mock_repository.NewMockUserRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)

	cardService := service.NewCardService(mock_card_repo, mock_user_repo, mock_logger, mock_mapping)

	cards := []*record.CardRecord{
		{
			ID:           1,
			UserID:       1,
			CardType:     "Debit",
			ExpireDate:   "2025-12-31",
			CVV:          "123",
			CardProvider: "Visa",
		},
		{
			ID:           2,
			UserID:       2,
			CardType:     "Credit",
			ExpireDate:   "2024-06-30",
			CVV:          "456",
			CardProvider: "MasterCard",
		},
	}

	expectedResponses := []*response.CardResponse{
		{
			ID:           1,
			UserID:       1,
			CardType:     "Debit",
			ExpireDate:   "2025-12-31",
			CVV:          "123",
			CardProvider: "Visa",
		},
		{
			ID:           2,
			UserID:       2,
			CardType:     "Credit",
			ExpireDate:   "2024-06-30",
			CVV:          "456",
			CardProvider: "MasterCard",
		},
	}

	page := 1
	pageSize := 10
	search := ""
	totalRecords := 2
	totalPages := (totalRecords + pageSize - 1) / pageSize

	mock_logger.EXPECT().Debug("Fetching all card records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	mock_card_repo.EXPECT().FindAllCards(search, page, pageSize).Return(cards, totalRecords, nil)

	mock_mapping.EXPECT().ToCardsResponse(cards).Return(expectedResponses)

	mock_logger.EXPECT().Debug("Successfully fetched card records",
		zap.Int("totalRecords", totalRecords),
		zap.Int("totalPages", totalPages))

	result, total, errResp := cardService.FindAll(page, pageSize, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponses, result)
	assert.Equal(t, totalPages, total)
}

func TestFindAllCard_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_user_repo := mock_repository.NewMockUserRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)

	cardService := service.NewCardService(mock_card_repo, mock_user_repo, mock_logger, mock_mapping)

	page := 1
	pageSize := 10
	search := ""

	mock_logger.EXPECT().Debug("Fetching all card records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search)).Times(1)

	mock_card_repo.EXPECT().FindAllCards(search, page, pageSize).Return([]*record.CardRecord{}, 0, nil)

	mock_logger.EXPECT().Debug("No card records found",
		zap.String("search", search)).Times(1)

	result, total, errResp := cardService.FindAll(page, pageSize, search)

	assert.Nil(t, result)
	assert.Equal(t, 0, total)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "No card records found", errResp.Message)
}

func TestFindAllCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_user_repo := mock_repository.NewMockUserRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)

	cardService := service.NewCardService(mock_card_repo, mock_user_repo, mock_logger, mock_mapping)

	page := 1
	pageSize := 10
	search := ""
	expectedErr := fmt.Errorf("database error")

	mock_logger.EXPECT().Debug("Fetching all card records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search)).Times(1)

	mock_card_repo.EXPECT().FindAllCards(search, page, pageSize).Return(nil, 0, expectedErr)

	mock_logger.EXPECT().Error("Failed to fetch all card records",
		zap.Error(expectedErr)).Times(1)

	result, total, errResp := cardService.FindAll(page, pageSize, search)

	assert.Nil(t, result)
	assert.Equal(t, 0, total)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to fetch card records", errResp.Message)
}

func TestFindByIdCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)

	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	cardID := 1
	cardRecord := &record.CardRecord{
		ID:           cardID,
		UserID:       1,
		CardNumber:   "1234567890123456",
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-21T09:00:00Z",
	}

	expectedResponse := &response.CardResponse{
		ID:           cardID,
		UserID:       1,
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-21T09:00:00Z",
		CardNumber:   "1234567890123456",
	}

	mock_card_repo.EXPECT().FindById(cardID).Return(cardRecord, nil)
	mock_mapping.EXPECT().ToCardResponse(cardRecord).Return(expectedResponse)
	mock_logger.EXPECT().Debug("Fetching card record by ID", zap.Int("card_id", cardID))
	mock_logger.EXPECT().Debug("Successfully fetched card record", zap.Int("card_id", cardID))

	result, err := cardService.FindById(cardID)

	t.Logf("Result: %+v", result)
	t.Logf("Expected: %+v", expectedResponse)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, result)
}

func TestFindByIdCard_Failure_CardNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)

	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	cardID := 1

	mock_card_repo.EXPECT().FindById(cardID).Return(nil, fmt.Errorf("card not found"))
	mock_logger.EXPECT().Debug("Fetching card record by ID", zap.Int("card_id", cardID))
	mock_logger.EXPECT().Error("Failed to fetch card by ID", zap.Error(fmt.Errorf("card not found")), zap.Int("card_id", cardID))

	result, err := cardService.FindById(cardID)

	t.Logf("Result: %+v", result)
	t.Logf("Error: %+v", err)

	// Assertions
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Card record not found", err.Message)
	assert.Equal(t, "error", err.Status)
}

func TestFindByIdCard_Failure_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)

	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	cardID := 1

	mock_card_repo.EXPECT().FindById(cardID).Return(nil, fmt.Errorf("repository error"))
	mock_logger.EXPECT().Debug("Fetching card record by ID", zap.Int("card_id", cardID))
	mock_logger.EXPECT().Error("Failed to fetch card by ID", zap.Error(fmt.Errorf("repository error")), zap.Int("card_id", cardID))

	result, err := cardService.FindById(cardID)

	t.Logf("Result: %+v", result)
	t.Logf("Error: %+v", err)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Card record not found", err.Message)
	assert.Equal(t, "error", err.Status)
}

func TestFindByUserID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	userID := 1
	cardRecord := &record.CardRecord{
		ID:           1,
		UserID:       userID,
		CardNumber:   "1234567890123456",
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-21T09:00:00Z",
	}

	expectedResponse := &response.CardResponse{
		ID:           1,
		UserID:       userID,
		CardNumber:   "1234567890123456",
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-21T09:00:00Z",
	}

	// Mock expectations
	mock_logger.EXPECT().Debug("Fetching card records by user ID",
		zap.Int("userID", userID)).Times(1)

	mock_card_repo.EXPECT().FindCardByUserId(userID).Return(cardRecord, nil)

	mock_mapping.EXPECT().ToCardResponse(cardRecord).Return(expectedResponse)

	mock_logger.EXPECT().Debug("Successfully fetched card records by user ID",
		zap.Int("userID", userID)).Times(1)

	// Test the method
	result, errResp := cardService.FindByUserID(userID)

	// Log the result and expected output for debugging
	t.Logf("Result: %+v", result)
	t.Logf("Expected: %+v", expectedResponse)

	// Assertions
	assert.Nil(t, errResp)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse.ID, result.ID)
	assert.Equal(t, expectedResponse.UserID, result.UserID)
	assert.Equal(t, expectedResponse.CardNumber, result.CardNumber)
	assert.Equal(t, expectedResponse.CardType, result.CardType)
	assert.Equal(t, expectedResponse.ExpireDate, result.ExpireDate)
	assert.Equal(t, expectedResponse.CVV, result.CVV)
	assert.Equal(t, expectedResponse.CardProvider, result.CardProvider)
	assert.Equal(t, expectedResponse.CreatedAt, result.CreatedAt)
	assert.Equal(t, expectedResponse.UpdatedAt, result.UpdatedAt)
}

func TestFindByUserID_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)

	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	userID := 1

	mock_card_repo.EXPECT().FindCardByUserId(userID).Return(nil, fmt.Errorf("repository error"))
	mock_logger.EXPECT().Debug("Fetching card records by user ID",
		zap.Int("userID", userID)).Times(1)
	mock_logger.EXPECT().Error("Failed to fetch cards by user ID", zap.Error(fmt.Errorf("repository error")), zap.Int("userID", userID))

	result, err := cardService.FindByUserID(userID)

	t.Logf("Result: %+v", result)
	t.Logf("Error: %+v", err)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to fetch cards by user ID", err.Message)
	assert.Equal(t, "error", err.Status)
}

func TestFindByActive_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	cardRecords := []*record.CardRecord{
		{
			ID:           1,
			UserID:       1,
			CardNumber:   "1234567890123456",
			CardType:     "Debit",
			ExpireDate:   "2025-12-31",
			CVV:          "123",
			CardProvider: "Visa",
			CreatedAt:    "2024-12-21T09:00:00Z",
			UpdatedAt:    "2024-12-21T09:00:00Z",
		},
		{
			ID:           2,
			UserID:       2,
			CardNumber:   "9876543210987654",
			CardType:     "Credit",
			ExpireDate:   "2026-11-30",
			CVV:          "456",
			CardProvider: "MasterCard",
			CreatedAt:    "2024-12-21T09:00:00Z",
			UpdatedAt:    "2024-12-21T09:00:00Z",
		},
	}

	expectedResponse := []*response.CardResponseDeleteAt{
		{
			ID:           1,
			UserID:       1,
			CardNumber:   "1234567890123456",
			CardType:     "Debit",
			ExpireDate:   "2025-12-31",
			CVV:          "123",
			CardProvider: "Visa",
			CreatedAt:    "2024-12-21T09:00:00Z",
			UpdatedAt:    "2024-12-21T09:00:00Z",
		},
		{
			ID:           2,
			UserID:       2,
			CardNumber:   "9876543210987654",
			CardType:     "Credit",
			ExpireDate:   "2026-11-30",
			CVV:          "456",
			CardProvider: "MasterCard",
			CreatedAt:    "2024-12-21T09:00:00Z",
			UpdatedAt:    "2024-12-21T09:00:00Z",
		},
	}

	page := 1
	pageSize := 10
	search := ""
	expected := 1

	mock_card_repo.EXPECT().FindByActive(search, page, pageSize).Return(cardRecords, expected, nil)
	mock_mapping.EXPECT().ToCardsResponseDeleteAt(cardRecords).Return(expectedResponse)
	mock_logger.EXPECT().Debug("Successfully fetched active card records").Times(1)

	result, totalRecord, errResp := cardService.FindByActive(page, pageSize, search)

	t.Logf("Result: %+v", result)
	t.Logf("Expected: %+v", expectedResponse)

	assert.Nil(t, errResp)
	assert.NotNil(t, result)
	assert.Equal(t, len(expectedResponse), len(result))
	assert.Equal(t, expected, totalRecord)
	for i := range expectedResponse {
		assert.Equal(t, expectedResponse[i], result[i])
	}
}

func TestFindByActiveCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	page := 1
	pageSize := 10
	search := ""

	expectedError := fmt.Errorf("database error")
	mock_card_repo.EXPECT().FindByActive(search, page, pageSize).Return(nil, 0, expectedError)

	mock_logger.EXPECT().Error("Failed to fetch active cards", gomock.Any()).Times(1)

	result, totalRecord, errResp := cardService.FindByActive(page, pageSize, search)

	assert.Nil(t, result)
	assert.Equal(t, 0, totalRecord)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to fetch active card records", errResp.Message)
}

func TestFindByTrashedCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	cardRecords := []*record.CardRecord{
		{
			ID:           1,
			UserID:       1,
			CardNumber:   "1234567890123456",
			CardType:     "Debit",
			ExpireDate:   "2025-12-31",
			CVV:          "123",
			CardProvider: "Visa",
			CreatedAt:    "2024-12-21T09:00:00Z",
			UpdatedAt:    "2024-12-21T09:00:00Z",
		},
		{
			ID:           2,
			UserID:       2,
			CardNumber:   "9876543210987654",
			CardType:     "Credit",
			ExpireDate:   "2026-11-30",
			CVV:          "456",
			CardProvider: "MasterCard",
			CreatedAt:    "2024-12-21T09:00:00Z",
			UpdatedAt:    "2024-12-21T09:00:00Z",
		},
	}

	expectedResponse := []*response.CardResponseDeleteAt{
		{
			ID:           1,
			UserID:       1,
			CardNumber:   "1234567890123456",
			CardType:     "Debit",
			ExpireDate:   "2025-12-31",
			CVV:          "123",
			CardProvider: "Visa",
			CreatedAt:    "2024-12-21T09:00:00Z",
			UpdatedAt:    "2024-12-21T09:00:00Z",
		},
		{
			ID:           2,
			UserID:       2,
			CardNumber:   "9876543210987654",
			CardType:     "Credit",
			ExpireDate:   "2026-11-30",
			CVV:          "456",
			CardProvider: "MasterCard",
			CreatedAt:    "2024-12-21T09:00:00Z",
			UpdatedAt:    "2024-12-21T09:00:00Z",
		},
	}

	page := 1
	pageSize := 10
	search := ""
	expected := 2

	mock_logger.EXPECT().Info("Fetching trashed card records").Times(1)
	mock_card_repo.EXPECT().FindByTrashed(search, page, pageSize).Return(cardRecords, expected, nil)
	mock_mapping.EXPECT().ToCardsResponseDeleteAt(cardRecords).Return(expectedResponse)
	mock_logger.EXPECT().Info("Successfully fetched trashed card records").Times(1)

	result, totalRecord, errResp := cardService.FindByTrashed(page, pageSize, search)

	t.Logf("Result: %+v", result)
	t.Logf("Error Response: %+v", errResp)

	assert.Nil(t, errResp)
	assert.NotNil(t, result)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, len(expectedResponse), len(result))
	for i := range expectedResponse {
		assert.Equal(t, expectedResponse[i], result[i])
	}
}

func TestFindByTrashedCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	page := 1
	pageSize := 10
	search := ""

	expectedError := fmt.Errorf("database error")
	mock_card_repo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, 0, expectedError)

	mock_logger.EXPECT().Info("Fetching trashed card records").Times(1)
	mock_logger.EXPECT().Error("Failed to fetch trashed cards", gomock.Any()).Times(1)

	result, totalRecord, errResp := cardService.FindByTrashed(page, pageSize, search)

	assert.Nil(t, result)
	assert.Equal(t, 0, totalRecord)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to fetch trashed card records", errResp.Message)
}

func TestFindByCardNumber_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	cardNumber := "1234567890123456"
	cardRecord := &record.CardRecord{
		ID:           1,
		UserID:       1,
		CardNumber:   cardNumber,
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-21T09:00:00Z",
	}
	expectedResponse := &response.CardResponse{
		ID:           1,
		UserID:       1,
		CardNumber:   cardNumber,
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-21T09:00:00Z",
	}

	mock_logger.EXPECT().Debug("Fetching card record by card number", zap.String("card_number", cardNumber)).Times(1)
	mock_card_repo.EXPECT().FindCardByCardNumber(cardNumber).Return(cardRecord, nil)
	mock_mapping.EXPECT().ToCardResponse(cardRecord).Return(expectedResponse)
	mock_logger.EXPECT().Debug("Successfully fetched card record by card number", zap.String("card_number", cardNumber)).Times(1)

	result, errResp := cardService.FindByCardNumber(cardNumber)

	t.Logf("Result: %+v", result)
	t.Logf("Error Response: %+v", errResp)

	assert.Nil(t, errResp)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, result)
}

func TestFindByCardNumber_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	cardNumber := "1234567890123456"
	expectedError := errors.New("card not found")

	mock_logger.EXPECT().Debug("Fetching card record by card number", zap.String("card_number", cardNumber)).Times(1)
	mock_card_repo.EXPECT().FindCardByCardNumber(cardNumber).Return(nil, expectedError)
	mock_logger.EXPECT().Error("Failed to fetch card by card number", zap.Error(expectedError), zap.String("card_number", cardNumber)).Times(1)

	result, errResp := cardService.FindByCardNumber(cardNumber)

	t.Logf("Result: %+v", result)
	t.Logf("Error Response: %+v", errResp)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Card record not found for the given card number", errResp.Message)
}

func TestCreateCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_user_repo := mock_repository.NewMockUserRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, mock_user_repo, mock_logger, mock_mapping)

	createCardRequest := requests.CreateCardRequest{
		UserID:       1,
		CardType:     "Debit",
		ExpireDate:   time.Now(),
		CVV:          "123",
		CardProvider: "Visa",
	}
	userRecord := &record.UserRecord{
		ID:        1,
		FirstName: "",
	}
	cardRecord := &record.CardRecord{
		ID:           1,
		UserID:       1,
		CardType:     createCardRequest.CardType,
		ExpireDate:   createCardRequest.ExpireDate.String(),
		CVV:          createCardRequest.CVV,
		CardProvider: createCardRequest.CardProvider,
		CreatedAt:    "2024-12-25T10:00:00Z",
		UpdatedAt:    "2024-12-25T10:00:00Z",
	}
	expectedResponse := &response.CardResponse{
		ID:           1,
		UserID:       1,
		CardType:     createCardRequest.CardType,
		ExpireDate:   createCardRequest.ExpireDate.String(),
		CVV:          createCardRequest.CVV,
		CardProvider: createCardRequest.CardProvider,
		CreatedAt:    "2024-12-25T10:00:00Z",
		UpdatedAt:    "2024-12-25T10:00:00Z",
	}

	mock_logger.EXPECT().Debug("Creating new card", zap.Int("userID", createCardRequest.UserID)).Times(1)
	mock_user_repo.EXPECT().FindById(createCardRequest.UserID).Return(userRecord, nil)
	mock_card_repo.EXPECT().CreateCard(&createCardRequest).Return(cardRecord, nil)
	mock_mapping.EXPECT().ToCardResponse(cardRecord).Return(expectedResponse)
	mock_logger.EXPECT().Debug("Successfully created new card", zap.Int("cardID", expectedResponse.ID)).Times(1)

	result, errResp := cardService.CreateCard(&createCardRequest)

	t.Logf("Result: %+v", result)
	t.Logf("Error Response: %+v", errResp)

	assert.Nil(t, errResp)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, result)
}

func TestCreateCard_Failure_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_user_repo := mock_repository.NewMockUserRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, mock_user_repo, mock_logger, mock_mapping)

	createCardRequest := requests.CreateCardRequest{
		UserID:       1,
		CardType:     "Debit",
		ExpireDate:   time.Now(),
		CVV:          "123",
		CardProvider: "Visa",
	}
	expectedError := errors.New("user not found")

	mock_logger.EXPECT().Debug("Creating new card", zap.Int("userID", createCardRequest.UserID)).Times(1)
	mock_user_repo.EXPECT().FindById(createCardRequest.UserID).Return(nil, expectedError)
	mock_logger.EXPECT().Error("Failed to find user by ID", zap.Error(expectedError), zap.Int("userID", createCardRequest.UserID)).Times(1)

	result, errResp := cardService.CreateCard(&createCardRequest)

	t.Logf("Result: %+v", result)
	t.Logf("Error Response: %+v", errResp)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "User not found", errResp.Message)
}

func TestCreateCard_Failure_CreationFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_user_repo := mock_repository.NewMockUserRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, mock_user_repo, mock_logger, mock_mapping)

	createCardRequest := requests.CreateCardRequest{
		UserID: 1,

		CardType:     "Debit",
		ExpireDate:   time.Now(),
		CVV:          "123",
		CardProvider: "Visa",
	}
	userRecord := &record.UserRecord{
		ID:        1,
		Email:     "test@example.com",
		Password:  "hashed_password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	expectedError := errors.New("failed to create card")

	mock_logger.EXPECT().Debug("Creating new card", zap.Int("userID", createCardRequest.UserID)).Times(1)
	mock_user_repo.EXPECT().FindById(createCardRequest.UserID).Return(userRecord, nil)
	mock_card_repo.EXPECT().CreateCard(&createCardRequest).Return(nil, expectedError)
	mock_logger.EXPECT().Error("Failed to create card", zap.Error(expectedError)).Times(1)

	result, errResp := cardService.CreateCard(&createCardRequest)

	t.Logf("Result: %+v", result)
	t.Logf("Error Response: %+v", errResp)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to create card", errResp.Message)
}

func TestUpdateCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_user_repo := mock_repository.NewMockUserRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, mock_user_repo, mock_logger, mock_mapping)

	request := requests.UpdateCardRequest{
		UserID:       1,
		CardID:       1,
		CardType:     "Debit",
		ExpireDate:   time.Now(),
		CVV:          "123",
		CardProvider: "Visa",
	}

	updatedCardRecord := record.CardRecord{
		ID:           1,
		UserID:       1,
		CardNumber:   "1234567890123456",
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-25T09:00:00Z",
	}

	expectedResponse := response.CardResponse{
		ID:           1,
		UserID:       1,
		CardNumber:   "1234567890123456",
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-25T09:00:00Z",
	}

	mock_logger.EXPECT().Debug("Updating card", zap.Int("userID", request.UserID), zap.Int("cardID", request.CardID)).Times(1)
	mock_user_repo.EXPECT().FindById(request.UserID).Return(&record.UserRecord{ID: 1, FirstName: "John", LastName: "Doe"}, nil).Times(1)
	mock_card_repo.EXPECT().UpdateCard(&request).Return(&updatedCardRecord, nil).Times(1)
	mock_mapping.EXPECT().ToCardResponse(&updatedCardRecord).Return(&expectedResponse).Times(1)
	mock_logger.EXPECT().Debug("Successfully updated card", zap.Int("cardID", expectedResponse.ID)).Times(1)

	result, errResp := cardService.UpdateCard(&request)

	t.Logf("Result: %+v", result)
	t.Logf("Error Response: %+v", errResp)

	assert.Nil(t, errResp)
	assert.NotNil(t, result)
	assert.Equal(t, &expectedResponse, result)
}

func TestUpdateCard_Failure_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_user_repo := mock_repository.NewMockUserRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, mock_user_repo, mock_logger, mock_mapping)

	request := requests.UpdateCardRequest{
		UserID: 1,
		CardID: 1,

		CardType:     "Debit",
		ExpireDate:   time.Now(),
		CVV:          "123",
		CardProvider: "Visa",
	}

	mock_logger.EXPECT().Debug("Updating card", zap.Int("userID", request.UserID), zap.Int("cardID", request.CardID)).Times(1)
	mock_user_repo.EXPECT().FindById(request.UserID).Return(nil, fmt.Errorf("user not found")).Times(1)
	mock_logger.EXPECT().Error("Failed to find user by ID", gomock.Any(), zap.Int("userID", request.UserID)).Times(1)

	result, errResp := cardService.UpdateCard(&request)

	t.Logf("Result: %+v", result)
	t.Logf("Error Response: %+v", errResp)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "User not found", errResp.Message)
}

func TestUpdateCard_Failure_UpdateFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_user_repo := mock_repository.NewMockUserRepository(ctrl)
	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, mock_user_repo, mock_logger, mock_mapping)

	request := requests.UpdateCardRequest{
		UserID: 1,
		CardID: 1,

		CardType:     "Debit",
		ExpireDate:   time.Now(),
		CVV:          "123",
		CardProvider: "Visa",
	}

	mock_logger.EXPECT().Debug("Updating card", zap.Int("userID", request.UserID), zap.Int("cardID", request.CardID)).Times(1)
	mock_user_repo.EXPECT().FindById(request.UserID).Return(&record.UserRecord{ID: 1, FirstName: "John", LastName: "Doe"}, nil).Times(1)
	mock_card_repo.EXPECT().UpdateCard(&request).Return(nil, fmt.Errorf("update failed")).Times(1)
	mock_logger.EXPECT().Error("Failed to update card", gomock.Any(), zap.Int("cardID", request.CardID)).Times(1)

	result, errResp := cardService.UpdateCard(&request)

	t.Logf("Result: %+v", result)
	t.Logf("Error Response: %+v", errResp)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to update card", errResp.Message)
}

func TestTrashedCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	cardID := 1
	trashedCardRecord := record.CardRecord{
		ID:           cardID,
		UserID:       1,
		CardNumber:   "1234567890123456",
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-25T09:00:00Z",
	}

	expectedResponse := response.CardResponse{
		ID:           cardID,
		UserID:       1,
		CardNumber:   "1234567890123456",
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-25T09:00:00Z",
	}

	mock_logger.EXPECT().Debug("Trashing card", zap.Int("cardID", cardID)).Times(1)
	mock_card_repo.EXPECT().TrashedCard(cardID).Return(&trashedCardRecord, nil).Times(1)
	mock_mapping.EXPECT().ToCardResponse(&trashedCardRecord).Return(&expectedResponse).Times(1)
	mock_logger.EXPECT().Debug("Successfully trashed card", zap.Int("cardID", cardID)).Times(1)

	result, errResp := cardService.TrashedCard(cardID)

	assert.Nil(t, errResp)
	assert.NotNil(t, result)
	assert.Equal(t, &expectedResponse, result)
}

func TestTrashedCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, nil)

	cardID := 1

	mock_logger.EXPECT().Debug("Trashing card", zap.Int("cardID", cardID)).Times(1)
	mock_card_repo.EXPECT().TrashedCard(cardID).Return(nil, fmt.Errorf("trash card error")).Times(1)
	mock_logger.EXPECT().Error("Failed to trash card", gomock.Any(), zap.Int("cardID", cardID)).Times(1)

	result, errResp := cardService.TrashedCard(cardID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to trash card", errResp.Message)
}

func TestRestoreCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockCardResponseMapper(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, mock_mapping)

	cardID := 1
	restoredCardRecord := record.CardRecord{
		ID:           cardID,
		UserID:       1,
		CardNumber:   "1234567890123456",
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-25T09:00:00Z",
	}

	expectedResponse := response.CardResponse{
		ID:           cardID,
		UserID:       1,
		CardNumber:   "1234567890123456",
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-25T09:00:00Z",
	}

	mock_logger.EXPECT().Debug("Restoring card", zap.Int("cardID", cardID)).Times(1)
	mock_card_repo.EXPECT().RestoreCard(cardID).Return(&restoredCardRecord, nil).Times(1)
	mock_mapping.EXPECT().ToCardResponse(&restoredCardRecord).Return(&expectedResponse).Times(1)
	mock_logger.EXPECT().Debug("Successfully restored card", zap.Int("cardID", cardID)).Times(1)

	result, errResp := cardService.RestoreCard(cardID)

	assert.Nil(t, errResp)
	assert.NotNil(t, result)
	assert.Equal(t, &expectedResponse, result)
}

func TestRestoreCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, nil)

	cardID := 1

	mock_logger.EXPECT().Debug("Restoring card", zap.Int("cardID", cardID)).Times(1)
	mock_card_repo.EXPECT().RestoreCard(cardID).Return(nil, fmt.Errorf("restore card error")).Times(1)
	mock_logger.EXPECT().Error("Failed to restore card", gomock.Any(), zap.Int("cardID", cardID)).Times(1)

	result, errResp := cardService.RestoreCard(cardID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to restore card", errResp.Message)
}

func TestDeleteCardPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, nil)

	cardID := 1

	mock_logger.EXPECT().Debug("Permanently deleting card", zap.Int("cardID", cardID)).Times(1)
	mock_card_repo.EXPECT().DeleteCardPermanent(cardID).Return(nil).Times(1)
	mock_logger.EXPECT().Debug("Successfully deleted card permanently", zap.Int("cardID", cardID)).Times(1)

	result, errResp := cardService.DeleteCardPermanent(cardID)

	assert.Nil(t, errResp)
	assert.Nil(t, result)
}

func TestDeleteCardPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_card_repo := mock_repository.NewMockCardRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	cardService := service.NewCardService(mock_card_repo, nil, mock_logger, nil)

	cardID := 1

	mock_logger.EXPECT().Debug("Permanently deleting card", zap.Int("cardID", cardID)).Times(1)
	mock_card_repo.EXPECT().DeleteCardPermanent(cardID).Return(fmt.Errorf("delete card error")).Times(1)
	mock_logger.EXPECT().Error("Failed to permanently delete card", gomock.Any(), zap.Int("cardID", cardID)).Times(1)

	result, errResp := cardService.DeleteCardPermanent(cardID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to permanently delete card: delete card error", errResp.Message)
}
