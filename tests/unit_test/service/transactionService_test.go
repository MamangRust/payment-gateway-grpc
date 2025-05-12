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

func TestFindAllTransactions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_transaction_repo := mock_repository.NewMockTransactionRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTransactionResponseMapper(ctrl)

	transactionService := service.NewTransactionService(nil, nil, nil, mock_transaction_repo, mock_logger, mock_mapping)

	page := 1
	pageSize := 2
	search := "example search"

	transactions := []*record.TransactionRecord{
		{
			ID:              1,
			CardNumber:      "1234",
			Amount:          500000,
			PaymentMethod:   "Credit Card",
			MerchantID:      10,
			TransactionTime: "2024-12-25T10:00:00Z",
			CreatedAt:       "2024-12-25T10:00:00Z",
			UpdatedAt:       "2024-12-25T11:00:00Z",
			DeletedAt:       nil,
		},
		{
			ID:              2,
			CardNumber:      "5678",
			Amount:          300000,
			PaymentMethod:   "Bank Transfer",
			MerchantID:      12,
			TransactionTime: "2024-12-25T12:00:00Z",
			CreatedAt:       "2024-12-25T12:00:00Z",
			UpdatedAt:       "2024-12-25T13:00:00Z",
			DeletedAt:       nil,
		},
	}
	totalRecords := 3

	mock_transaction_repo.EXPECT().FindAllTransactions(search, page, pageSize).Return(transactions, totalRecords, nil).Times(1)

	mappedTransactions := []*response.TransactionResponse{
		{
			ID:              1,
			CardNumber:      "1234",
			Amount:          500000,
			PaymentMethod:   "Credit Card",
			MerchantID:      10,
			TransactionTime: "2024-12-25T10:00:00Z",
			CreatedAt:       "2024-12-25T10:00:00Z",
			UpdatedAt:       "2024-12-25T11:00:00Z",
		},
		{
			ID:              2,
			CardNumber:      "5678",
			Amount:          300000,
			PaymentMethod:   "Bank Transfer",
			MerchantID:      12,
			TransactionTime: "2024-12-25T12:00:00Z",
			CreatedAt:       "2024-12-25T12:00:00Z",
			UpdatedAt:       "2024-12-25T13:00:00Z",
		},
	}
	mock_mapping.EXPECT().ToTransactionsResponse(transactions).Return(mappedTransactions).Times(1)

	result, totalPages, errResp := transactionService.FindAll(page, pageSize, search)

	assert.NotNil(t, result)
	assert.Nil(t, errResp)
	assert.Equal(t, 3, totalPages)
	assert.Equal(t, mappedTransactions, result)
}

func TestFindAllTransactions_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_transaction_repo := mock_repository.NewMockTransactionRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	transactionService := service.NewTransactionService(nil, nil, nil, mock_transaction_repo, mock_logger, nil)

	page := 1
	pageSize := 2
	search := "example search"

	mock_transaction_repo.EXPECT().FindAllTransactions(search, page, pageSize).Return(nil, 0, errors.New("database error")).Times(1)
	mock_logger.EXPECT().Error("failed to fetch transactions", zap.Error(errors.New("database error"))).Times(1)

	result, totalPages, errResp := transactionService.FindAll(page, pageSize, search)

	assert.Nil(t, result)
	assert.Equal(t, 0, totalPages)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to fetch transactions", errResp.Message)
}

func TestFindByIdTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_transaction_repo := mock_repository.NewMockTransactionRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTransactionResponseMapper(ctrl)

	transactionService := service.NewTransactionService(nil, nil, nil, mock_transaction_repo, mock_logger, mock_mapping)

	transactionID := 1

	transaction := &record.TransactionRecord{
		ID:              1,
		CardNumber:      "1234",
		Amount:          500000,
		PaymentMethod:   "Credit Card",
		MerchantID:      10,
		TransactionTime: "2024-12-25T10:00:00Z",
		CreatedAt:       "2024-12-25T10:00:00Z",
		UpdatedAt:       "2024-12-25T11:00:00Z",
		DeletedAt:       nil,
	}

	mock_transaction_repo.EXPECT().FindById(transactionID).Return(transaction, nil).Times(1)

	mappedTransaction := &response.TransactionResponse{
		ID:              1,
		CardNumber:      "1234",
		Amount:          500000,
		PaymentMethod:   "Credit Card",
		MerchantID:      10,
		TransactionTime: "2024-12-25T10:00:00Z",
		CreatedAt:       "2024-12-25T10:00:00Z",
		UpdatedAt:       "2024-12-25T11:00:00Z",
	}
	mock_mapping.EXPECT().ToTransactionResponse(transaction).Return(mappedTransaction).Times(1)

	result, errResp := transactionService.FindById(transactionID)

	assert.NotNil(t, result)
	assert.Nil(t, errResp)
	assert.Equal(t, mappedTransaction, result)
}

func TestFindByIdTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_transaction_repo := mock_repository.NewMockTransactionRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	transactionService := service.NewTransactionService(nil, nil, nil, mock_transaction_repo, mock_logger, nil)

	transactionID := 1

	mock_transaction_repo.EXPECT().FindById(transactionID).Return(nil, errors.New("transaction not found")).Times(1)
	mock_logger.EXPECT().Error("failed to find transaction", zap.Error(errors.New("transaction not found"))).Times(1)

	result, errResp := transactionService.FindById(transactionID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Transaction not found", errResp.Message)
}

func TestFindByActiveTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_transaction_repo := mock_repository.NewMockTransactionRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTransactionResponseMapper(ctrl)

	transactionService := service.NewTransactionService(nil, nil, nil, mock_transaction_repo, mock_logger, mock_mapping)

	transactions := []*record.TransactionRecord{
		{
			ID:              1,
			CardNumber:      "1234",
			Amount:          500000,
			PaymentMethod:   "Credit Card",
			MerchantID:      10,
			TransactionTime: "2024-12-25T10:00:00Z",
			CreatedAt:       "2024-12-25T10:00:00Z",
			UpdatedAt:       "2024-12-25T11:00:00Z",
			DeletedAt:       nil,
		},
		{
			ID:              2,
			CardNumber:      "5678",
			Amount:          300000,
			PaymentMethod:   "Bank Transfer",
			MerchantID:      12,
			TransactionTime: "2024-12-25T12:00:00Z",
			CreatedAt:       "2024-12-25T12:00:00Z",
			UpdatedAt:       "2024-12-25T13:00:00Z",
			DeletedAt:       nil,
		},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 2

	mock_transaction_repo.EXPECT().FindByActive(search, page, pageSize).Return(transactions, expected, nil).Times(1)

	mappedTransactions := []*response.TransactionResponseDeleteAt{
		{
			ID:              1,
			CardNumber:      "1234",
			Amount:          500000,
			PaymentMethod:   "Credit Card",
			MerchantID:      10,
			TransactionTime: "2024-12-25T10:00:00Z",
			CreatedAt:       "2024-12-25T10:00:00Z",
			UpdatedAt:       "2024-12-25T11:00:00Z",
		},
		{
			ID:              2,
			CardNumber:      "5678",
			Amount:          300000,
			PaymentMethod:   "Bank Transfer",
			MerchantID:      12,
			TransactionTime: "2024-12-25T12:00:00Z",
			CreatedAt:       "2024-12-25T12:00:00Z",
			UpdatedAt:       "2024-12-25T13:00:00Z",
		},
	}

	mock_mapping.EXPECT().ToTransactionsResponseDeleteAt(transactions).Return(mappedTransactions).Times(1)
	mock_logger.EXPECT().Debug("Successfully fetched active transaction records", zap.Int("record_count", 2)).Times(1)

	result, totalRecord, errResp := transactionService.FindByActive(pageSize, page, search)

	assert.NotNil(t, result)
	assert.Nil(t, errResp)
	assert.Len(t, result, expected)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, mappedTransactions, result)
}

func TestFindByActiveTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_transaction_repo := mock_repository.NewMockTransactionRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	transactionService := service.NewTransactionService(nil, nil, nil, mock_transaction_repo, mock_logger, nil)

	mock_transaction_repo.EXPECT().FindByActive(search, page, pageSize).Return(nil, expected, errors.New("no active transactions found")).Times(1)
	mock_logger.EXPECT().Error("Failed to fetch active transaction records", zap.Error(errors.New("no active transactions found"))).Times(1)

	result, totalRecord, errResp := transactionService.FindByActive(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "No active transaction records found", errResp.Message)
}

func TestFindByTrashedTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_transaction_repo := mock_repository.NewMockTransactionRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTransactionResponseMapper(ctrl)

	transactionService := service.NewTransactionService(nil, nil, nil, mock_transaction_repo, mock_logger, mock_mapping)

	transactions := []*record.TransactionRecord{
		{
			ID:              1,
			CardNumber:      "1234",
			Amount:          500000,
			PaymentMethod:   "Credit Card",
			MerchantID:      10,
			TransactionTime: "2024-12-25T10:00:00Z",
			CreatedAt:       "2024-12-25T10:00:00Z",
			UpdatedAt:       "2024-12-25T11:00:00Z",
			DeletedAt:       nil,
		},
		{
			ID:              2,
			CardNumber:      "5678",
			Amount:          300000,
			PaymentMethod:   "Bank Transfer",
			MerchantID:      12,
			TransactionTime: "2024-12-25T12:00:00Z",
			CreatedAt:       "2024-12-25T12:00:00Z",
			UpdatedAt:       "2024-12-25T13:00:00Z",
			DeletedAt:       nil,
		},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 2

	mock_transaction_repo.EXPECT().FindByTrashed(search, page, pageSize).Return(transactions, expected, nil).Times(1)

	mappedTransactions := []*response.TransactionResponseDeleteAt{
		{
			ID:              1,
			CardNumber:      "1234",
			Amount:          500000,
			PaymentMethod:   "Credit Card",
			MerchantID:      10,
			TransactionTime: "2024-12-25T10:00:00Z",
			CreatedAt:       "2024-12-25T10:00:00Z",
			UpdatedAt:       "2024-12-25T11:00:00Z",
		},
		{
			ID:              2,
			CardNumber:      "5678",
			Amount:          300000,
			PaymentMethod:   "Bank Transfer",
			MerchantID:      12,
			TransactionTime: "2024-12-25T12:00:00Z",
			CreatedAt:       "2024-12-25T12:00:00Z",
			UpdatedAt:       "2024-12-25T13:00:00Z",
		},
	}

	mock_logger.EXPECT().Info("Fetching trashed transaction records")

	mock_mapping.EXPECT().ToTransactionsResponseDeleteAt(transactions).Return(mappedTransactions).Times(1)
	mock_logger.EXPECT().Debug("Successfully fetched trashed transaction records", zap.Int("record_count", len(transactions))).Times(1)

	result, totalRecord, errResp := transactionService.FindByTrashed(pageSize, page, search)

	assert.NotNil(t, result)
	assert.Nil(t, errResp)
	assert.Len(t, result, 2)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, mappedTransactions, result)
}

func TestFindByTrashedTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_transaction_repo := mock_repository.NewMockTransactionRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	transactionService := service.NewTransactionService(nil, nil, nil, mock_transaction_repo, mock_logger, nil)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	mock_logger.EXPECT().Info("Fetching trashed transaction records")

	mock_transaction_repo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, expected, errors.New("no trashed transactions found")).Times(1)
	mock_logger.EXPECT().Error("Failed to fetch trashed transaction records", zap.Error(errors.New("no trashed transactions found"))).Times(1)

	result, totalRecord, errResp := transactionService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "No trashed transaction records found", errResp.Message)
}

func TestFindByCardNumberTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_transaction_repo := mock_repository.NewMockTransactionRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockTransactionResponseMapper(ctrl)

	transactionService := service.NewTransactionService(nil, nil, nil, mock_transaction_repo, mock_logger, mock_mapping)

	cardNumber := "1234"
	transactions := []*record.TransactionRecord{
		{
			ID:              1,
			CardNumber:      "1234",
			Amount:          500000,
			PaymentMethod:   "Credit Card",
			MerchantID:      10,
			TransactionTime: "2024-12-25T10:00:00Z",
			CreatedAt:       "2024-12-25T10:00:00Z",
			UpdatedAt:       "2024-12-25T11:00:00Z",
			DeletedAt:       nil,
		},
		{
			ID:              2,
			CardNumber:      "5678",
			Amount:          300000,
			PaymentMethod:   "Bank Transfer",
			MerchantID:      12,
			TransactionTime: "2024-12-25T12:00:00Z",
			CreatedAt:       "2024-12-25T12:00:00Z",
			UpdatedAt:       "2024-12-25T13:00:00Z",
			DeletedAt:       nil,
		},
	}

	mock_transaction_repo.EXPECT().FindByCardNumber(cardNumber).Return(transactions, nil).Times(1)

	mappedTransactions := []*response.TransactionResponse{
		{
			ID:              1,
			CardNumber:      "1234",
			Amount:          500000,
			PaymentMethod:   "Credit Card",
			MerchantID:      10,
			TransactionTime: "2024-12-25T10:00:00Z",
			CreatedAt:       "2024-12-25T10:00:00Z",
			UpdatedAt:       "2024-12-25T11:00:00Z",
		},
		{
			ID:              2,
			CardNumber:      "5678",
			Amount:          300000,
			PaymentMethod:   "Bank Transfer",
			MerchantID:      12,
			TransactionTime: "2024-12-25T12:00:00Z",
			CreatedAt:       "2024-12-25T12:00:00Z",
			UpdatedAt:       "2024-12-25T13:00:00Z",
		},
	}

	mock_mapping.EXPECT().ToTransactionsResponse(transactions).Return(mappedTransactions).Times(1)
	mock_logger.EXPECT().Debug("Successfully fetched transactions by card number", zap.String("card_number", cardNumber), zap.Int("record_count", len(transactions))).Times(1)

	result, errResp := transactionService.FindByCardNumber(cardNumber)

	assert.NotNil(t, result)
	assert.Nil(t, errResp)
	assert.Len(t, result, 2)
	assert.Equal(t, mappedTransactions, result)
}

func TestFindByCardNumberTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_transaction_repo := mock_repository.NewMockTransactionRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	transactionService := service.NewTransactionService(nil, nil, nil, mock_transaction_repo, mock_logger, nil)

	cardNumber := "1234"

	mock_transaction_repo.EXPECT().FindByCardNumber(cardNumber).Return(nil, errors.New("no transactions found")).Times(1)
	mock_logger.EXPECT().Error("Failed to fetch transactions by card number", zap.Error(errors.New("no transactions found")), zap.String("card_number", cardNumber)).Times(1)

	result, errResp := transactionService.FindByCardNumber(cardNumber)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "No transactions found for the given card number", errResp.Message)
}

func TestCreateTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantRepo := mock_repository.NewMockMerchantRepository(ctrl)
	mockCardRepo := mock_repository.NewMockCardRepository(ctrl)
	mockSaldoRepo := mock_repository.NewMockSaldoRepository(ctrl)
	mockTransactionRepo := mock_repository.NewMockTransactionRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransactionResponseMapper(ctrl)

	transactionService := service.NewTransactionService(
		mockMerchantRepo,
		mockCardRepo,
		mockSaldoRepo,
		mockTransactionRepo,
		mockLogger,
		mockMapping,
	)

	apiKey := "test-api-key"
	merchantID := 1
	merchantUserID := 2

	merchant := &record.MerchantRecord{
		ID:     merchantID,
		UserID: merchantUserID,
	}

	request := &requests.CreateTransactionRequest{
		CardNumber: "4111111111111111",
		Amount:     1000,
	}

	customerCard := &record.CardRecord{
		CardNumber: request.CardNumber,
	}

	merchantCard := &record.CardRecord{
		CardNumber: "4222222222222222",
	}

	customerSaldo := &record.SaldoRecord{
		CardNumber:   customerCard.CardNumber,
		TotalBalance: 5000,
	}

	merchantSaldo := &record.SaldoRecord{
		CardNumber:   merchantCard.CardNumber,
		TotalBalance: 10000,
	}

	transaction := &record.TransactionRecord{
		ID:         1,
		CardNumber: request.CardNumber,
		Amount:     request.Amount,
	}

	expectedResponse := &response.TransactionResponse{
		ID:     1,
		Amount: request.Amount,
	}

	mockMerchantRepo.EXPECT().
		FindByApiKey(apiKey).
		Return(merchant, nil)

	mockCardRepo.EXPECT().
		FindCardByCardNumber(request.CardNumber).
		Return(customerCard, nil)

	mockSaldoRepo.EXPECT().
		FindByCardNumber(customerCard.CardNumber).
		Return(customerSaldo, nil)

	mockSaldoRepo.EXPECT().
		UpdateSaldoBalance(gomock.Any()).
		DoAndReturn(func(req *requests.UpdateSaldoBalance) (*record.SaldoRecord, error) {
			assert.Equal(t, customerCard.CardNumber, req.CardNumber)
			assert.Equal(t, 4000, req.TotalBalance) // 5000 - 1000
			return customerSaldo, nil
		})

	mockTransactionRepo.EXPECT().
		CreateTransaction(gomock.Any()).
		DoAndReturn(func(req *requests.CreateTransactionRequest) (*record.TransactionRecord, error) {
			assert.Equal(t, merchantID, *req.MerchantID)
			return transaction, nil
		})

	mockCardRepo.EXPECT().
		FindCardByUserId(merchant.UserID).
		Return(merchantCard, nil)

	mockSaldoRepo.EXPECT().
		FindByCardNumber(merchantCard.CardNumber).
		Return(merchantSaldo, nil)

	mockLogger.EXPECT().
		Debug("Updating merchant saldo", gomock.Any())

	mockSaldoRepo.EXPECT().
		UpdateSaldoBalance(gomock.Any()).
		DoAndReturn(func(req *requests.UpdateSaldoBalance) (*record.SaldoRecord, error) {
			assert.Equal(t, merchantCard.CardNumber, req.CardNumber)
			assert.Equal(t, 11000, req.TotalBalance)
			return merchantSaldo, nil
		})

	mockMapping.EXPECT().
		ToTransactionResponse(transaction).
		Return(expectedResponse)

	result, errResp := transactionService.Create(apiKey, request)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestUpdateTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantRepo := mock_repository.NewMockMerchantRepository(ctrl)
	mockCardRepo := mock_repository.NewMockCardRepository(ctrl)
	mockSaldoRepo := mock_repository.NewMockSaldoRepository(ctrl)
	mockTransactionRepo := mock_repository.NewMockTransactionRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransactionResponseMapper(ctrl)

	transactionService := service.NewTransactionService(
		mockMerchantRepo,
		mockCardRepo,
		mockSaldoRepo,
		mockTransactionRepo,
		mockLogger,
		mockMapping,
	)

	apiKey := "test-api-key"
	merchantID := 1
	cardNumber := "1234567890"
	oldAmount := 1000
	newAmount := 500
	layout := "2006-01-02 15:04:05"
	transactionTime := time.Now().Format(layout)
	parsedTime, _ := time.Parse(layout, transactionTime)

	updateRequest := &requests.UpdateTransactionRequest{
		TransactionID:   1,
		Amount:          newAmount,
		PaymentMethod:   "credit_card",
		TransactionTime: parsedTime,
	}

	existingTransaction := &record.TransactionRecord{
		ID:              1,
		MerchantID:      merchantID,
		CardNumber:      cardNumber,
		Amount:          oldAmount,
		PaymentMethod:   "credit_card",
		TransactionTime: transactionTime,
	}

	merchant := &record.MerchantRecord{
		ID:     merchantID,
		ApiKey: apiKey,
	}

	card := &record.CardRecord{
		CardNumber: cardNumber,
	}

	saldo := &record.SaldoRecord{
		CardNumber:   cardNumber,
		TotalBalance: 2000,
	}

	expectedResponse := &response.TransactionResponse{
		ID:              1,
		CardNumber:      cardNumber,
		Amount:          newAmount,
		PaymentMethod:   "credit_card",
		TransactionTime: transactionTime,
	}

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()

	mockTransactionRepo.EXPECT().
		FindById(updateRequest.TransactionID).
		Return(existingTransaction, nil)

	mockMerchantRepo.EXPECT().
		FindByApiKey(apiKey).
		Return(merchant, nil)

	mockCardRepo.EXPECT().
		FindCardByCardNumber(cardNumber).
		Return(card, nil)

	mockSaldoRepo.EXPECT().
		FindByCardNumber(cardNumber).
		Return(saldo, nil)

	mockSaldoRepo.EXPECT().
		UpdateSaldoBalance(gomock.Any()).
		Return(&record.SaldoRecord{
			CardNumber:   cardNumber,
			TotalBalance: saldo.TotalBalance + oldAmount,
		}, nil)

	mockSaldoRepo.EXPECT().
		UpdateSaldoBalance(gomock.Any()).
		Return(&record.SaldoRecord{
			CardNumber:   cardNumber,
			TotalBalance: saldo.TotalBalance - newAmount,
		}, nil)

	mockTransactionRepo.EXPECT().
		UpdateTransaction(gomock.Any()).
		Return(existingTransaction, nil)

	mockMapping.EXPECT().
		ToTransactionResponse(gomock.Any()).
		Return(expectedResponse)

	result, err := transactionService.Update(apiKey, updateRequest)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, result)
}

func TestUpdateTransaction_InsufficientBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantRepo := mock_repository.NewMockMerchantRepository(ctrl)
	mockCardRepo := mock_repository.NewMockCardRepository(ctrl)
	mockSaldoRepo := mock_repository.NewMockSaldoRepository(ctrl)
	mockTransactionRepo := mock_repository.NewMockTransactionRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransactionResponseMapper(ctrl)

	transactionService := service.NewTransactionService(
		mockMerchantRepo,
		mockCardRepo,
		mockSaldoRepo,
		mockTransactionRepo,
		mockLogger,
		mockMapping,
	)

	apiKey := "test-api-key"
	merchantID := 1
	cardNumber := "1234567890"
	oldAmount := 1000
	newAmount := 5000
	layout := "2006-01-02 15:04:05"
	transactionTime := time.Now().Format(layout)
	parsedTime, _ := time.Parse(layout, transactionTime)

	updateRequest := &requests.UpdateTransactionRequest{
		TransactionID:   1,
		Amount:          newAmount,
		PaymentMethod:   "credit_card",
		TransactionTime: parsedTime,
	}

	existingTransaction := &record.TransactionRecord{
		ID:              1,
		MerchantID:      merchantID,
		CardNumber:      cardNumber,
		Amount:          oldAmount,
		PaymentMethod:   "credit_card",
		TransactionTime: transactionTime,
	}

	merchant := &record.MerchantRecord{
		ID:     merchantID,
		ApiKey: apiKey,
	}

	card := &record.CardRecord{
		CardNumber: cardNumber,
	}

	saldo := &record.SaldoRecord{
		CardNumber:   cardNumber,
		TotalBalance: 2000,
	}

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()

	mockTransactionRepo.EXPECT().
		FindById(updateRequest.TransactionID).
		Return(existingTransaction, nil)

	mockMerchantRepo.EXPECT().
		FindByApiKey(apiKey).
		Return(merchant, nil)

	mockCardRepo.EXPECT().
		FindCardByCardNumber(cardNumber).
		Return(card, nil)

	mockSaldoRepo.EXPECT().
		FindByCardNumber(cardNumber).
		Return(saldo, nil)

	mockSaldoRepo.EXPECT().
		UpdateSaldoBalance(gomock.Any()).
		Return(&record.SaldoRecord{
			CardNumber:   cardNumber,
			TotalBalance: saldo.TotalBalance + oldAmount,
		}, nil)

	result, err := transactionService.Update(apiKey, updateRequest)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "error", err.Status)
	assert.Equal(t, "Insufficient balance for updated transaction", err.Message)
}

func TestTrashedTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock_repository.NewMockTransactionRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransactionResponseMapper(ctrl)

	transactionService := service.NewTransactionService(
		nil, nil, nil,
		mockTransactionRepo,
		mockLogger,
		mockMapping,
	)

	transactionID := 1
	expectedRecord := &record.TransactionRecord{
		ID: transactionID,
	}
	expectedResponse := &response.TransactionResponse{
		ID: transactionID,
	}

	mockTransactionRepo.EXPECT().
		TrashedTransaction(transactionID).
		Return(expectedRecord, nil)

	mockMapping.EXPECT().
		ToTransactionResponse(expectedRecord).
		Return(expectedResponse)

	mockLogger.EXPECT().
		Debug("Successfully trashed transaction", zap.Int("transaction_id", transactionID))

	result, errResp := transactionService.TrashedTransaction(transactionID)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestTrashedTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock_repository.NewMockTransactionRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transactionService := service.NewTransactionService(
		nil, nil, nil,
		mockTransactionRepo,
		mockLogger,
		nil,
	)

	transactionID := 1

	mockTransactionRepo.EXPECT().
		TrashedTransaction(transactionID).
		Return(nil, errors.New("database error"))

	mockLogger.EXPECT().
		Error("Failed to trash transaction", gomock.Any(), zap.Int("transaction_id", transactionID))

	result, errResp := transactionService.TrashedTransaction(transactionID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to trash transaction", errResp.Message)
}

func TestRestoreTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock_repository.NewMockTransactionRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockTransactionResponseMapper(ctrl)

	transactionService := service.NewTransactionService(
		nil, nil, nil,
		mockTransactionRepo,
		mockLogger,
		mockMapping,
	)

	transactionID := 1
	expectedRecord := &record.TransactionRecord{
		ID: transactionID,
	}
	expectedResponse := &response.TransactionResponse{
		ID: transactionID,
	}

	mockTransactionRepo.EXPECT().
		RestoreTransaction(transactionID).
		Return(expectedRecord, nil)

	mockMapping.EXPECT().
		ToTransactionResponse(expectedRecord).
		Return(expectedResponse)

	mockLogger.EXPECT().
		Debug("Successfully restored transaction", zap.Int("transaction_id", transactionID))

	result, errResp := transactionService.RestoreTransaction(transactionID)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestRestoreTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock_repository.NewMockTransactionRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transactionService := service.NewTransactionService(
		nil, nil, nil,
		mockTransactionRepo,
		mockLogger,
		nil,
	)

	transactionID := 1

	mockTransactionRepo.EXPECT().
		RestoreTransaction(transactionID).
		Return(nil, errors.New("database error"))

	mockLogger.EXPECT().
		Error("Failed to restore transaction", gomock.Any(), zap.Int("transaction_id", transactionID))

	result, errResp := transactionService.RestoreTransaction(transactionID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to restore transaction", errResp.Message)
}

func TestDeleteTransactionPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock_repository.NewMockTransactionRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transactionService := service.NewTransactionService(
		nil, nil, nil,
		mockTransactionRepo,
		mockLogger,
		nil,
	)

	transactionID := 1

	mockTransactionRepo.EXPECT().
		DeleteTransactionPermanent(transactionID).
		Return(nil)

	mockLogger.EXPECT().
		Debug("Successfully permanently deleted transaction", zap.Int("transaction_id", transactionID))

	result, errResp := transactionService.DeleteTransactionPermanent(transactionID)

	assert.Nil(t, errResp)
	assert.Nil(t, result)
}

func TestDeleteTransactionPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock_repository.NewMockTransactionRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transactionService := service.NewTransactionService(
		nil, nil, nil,
		mockTransactionRepo,
		mockLogger,
		nil,
	)

	transactionID := 1

	mockTransactionRepo.EXPECT().
		DeleteTransactionPermanent(transactionID).
		Return(errors.New("database error"))

	mockLogger.EXPECT().
		Error("Failed to permanently delete transaction", gomock.Any(), zap.Int("transaction_id", transactionID))

	result, errResp := transactionService.DeleteTransactionPermanent(transactionID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to permanently delete transaction", errResp.Message)
}
