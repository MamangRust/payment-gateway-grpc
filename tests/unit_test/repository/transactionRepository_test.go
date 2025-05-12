package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/tests/utils"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindAllTransactions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

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

	mockRepo.EXPECT().FindAllTransactions("", 1, 10).Return(transactions, 2, nil)

	results, total, err := mockRepo.FindAllTransactions("", 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 2, total)
	assert.Equal(t, transactions, results)
}

func TestFindAllTransactions_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	mockRepo.EXPECT().FindAllTransactions("", 1, 10).Return(nil, 0, fmt.Errorf("database error"))

	results, total, err := mockRepo.FindAllTransactions("", 1, 10)

	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Equal(t, 0, total)
	assert.EqualError(t, err, "database error")
}

func TestFindAllTransactions_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	mockRepo.EXPECT().FindAllTransactions("", 1, 10).Return([]*record.TransactionRecord{}, 0, nil)

	results, total, err := mockRepo.FindAllTransactions("", 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 0, total)
	assert.Empty(t, results)
}

func TestFindByIdTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	expectedTransaction := &record.TransactionRecord{
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

	mockRepo.EXPECT().FindById(1).Return(expectedTransaction, nil)

	result, err := mockRepo.FindById(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTransaction, result)
}

func TestFindByIdTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	mockRepo.EXPECT().FindById(1).Return(nil, fmt.Errorf("transaction not found"))

	result, err := mockRepo.FindById(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "transaction not found")
}

func TestFindByActiveTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	expectedTransactions := []*record.TransactionRecord{
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
			MerchantID:      15,
			TransactionTime: "2024-12-25T12:00:00Z",
			CreatedAt:       "2024-12-25T12:00:00Z",
			UpdatedAt:       "2024-12-25T13:00:00Z",
			DeletedAt:       nil,
		},
	}

	page := 1
	pageSize := 10
	search := ""
	expected := 2

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(expectedTransactions, 2, nil)

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTransactions, result)
}

func TestFindByActiveTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(nil, 0, fmt.Errorf("database error"))

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")
}

func TestFindByActiveTransaction_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return([]*record.TransactionRecord{}, 0, nil)

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestFindByTrashedTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	expectedTransactions := []*record.TransactionRecord{
		{
			ID:              3,
			CardNumber:      "1111",
			Amount:          200000,
			PaymentMethod:   "Cash",
			MerchantID:      5,
			TransactionTime: "2024-12-24T10:00:00Z",
			CreatedAt:       "2024-12-24T10:00:00Z",
			UpdatedAt:       "2024-12-24T11:00:00Z",
			DeletedAt:       utils.PtrString("2024-12-24T12:00:00Z"),
		},
	}

	page := 1
	pageSize := 10
	search := ""
	expected := 1

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(expectedTransactions, 1, nil)

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTransactions, result)
}

func TestFindByTrashedTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, 0, fmt.Errorf("database error"))

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Error(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")
}

func TestFindByTrashedTransaction_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return([]*record.TransactionRecord{}, 0, nil)

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestFindByCardNumberTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	cardNumber := "1234"
	expectedTransactions := []*record.TransactionRecord{
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
	}

	mockRepo.EXPECT().FindByCardNumber(cardNumber).Return(expectedTransactions, nil)

	result, err := mockRepo.FindByCardNumber(cardNumber)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTransactions, result)
}

func TestFindByCardNumberTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	cardNumber := "1234"
	mockRepo.EXPECT().FindByCardNumber(cardNumber).Return(nil, fmt.Errorf("database error"))

	result, err := mockRepo.FindByCardNumber(cardNumber)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")
}

func TestFindByCardNumber_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	cardNumber := "1234"
	mockRepo.EXPECT().FindByCardNumber(cardNumber).Return([]*record.TransactionRecord{}, nil)

	result, err := mockRepo.FindByCardNumber(cardNumber)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestFindTransactionByMerchantId_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	merchantID := 10
	expectedTransactions := []*record.TransactionRecord{
		{
			ID:              2,
			CardNumber:      "5678",
			Amount:          300000,
			PaymentMethod:   "Bank Transfer",
			MerchantID:      10,
			TransactionTime: "2024-12-25T12:00:00Z",
			CreatedAt:       "2024-12-25T12:00:00Z",
			UpdatedAt:       "2024-12-25T13:00:00Z",
			DeletedAt:       nil,
		},
	}

	mockRepo.EXPECT().FindTransactionByMerchantId(merchantID).Return(expectedTransactions, nil)

	result, err := mockRepo.FindTransactionByMerchantId(merchantID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTransactions, result)
}

func TestFindTransactionByMerchantId_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	merchantID := 10
	mockRepo.EXPECT().FindTransactionByMerchantId(merchantID).Return(nil, fmt.Errorf("database error"))

	result, err := mockRepo.FindTransactionByMerchantId(merchantID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")
}

func TestFindTransactionByMerchantId_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	merchantID := 10
	mockRepo.EXPECT().FindTransactionByMerchantId(merchantID).Return([]*record.TransactionRecord{}, nil)

	result, err := mockRepo.FindTransactionByMerchantId(merchantID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestCountTransactionsByDate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	date := "2024-12-25"
	expectedCount := 5

	mockRepo.EXPECT().CountTransactionsByDate(date).Return(expectedCount, nil)

	count, err := mockRepo.CountTransactionsByDate(date)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestCountTransactionsByDate_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	date := "2024-12-25"
	mockRepo.EXPECT().CountTransactionsByDate(date).Return(0, fmt.Errorf("database error"))

	count, err := mockRepo.CountTransactionsByDate(date)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.EqualError(t, err, "database error")
}

func TestCountAllTransactions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	expectedCount := int64(100)
	expectedCountPtr := &expectedCount

	mockRepo.EXPECT().CountAllTransactions().Return(expectedCountPtr, nil)

	count, err := mockRepo.CountAllTransactions()

	assert.NoError(t, err)
	assert.Equal(t, expectedCountPtr, count)
}

func TestCountAllTransactions_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	expectedError := fmt.Errorf("database error")
	mockRepo.EXPECT().CountAllTransactions().Return(nil, expectedError)

	count, err := mockRepo.CountAllTransactions()

	assert.Error(t, err)
	assert.Nil(t, count)
	assert.EqualError(t, err, "database error")
}

func TestCreateTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	now := time.Now()
	request := requests.CreateTransactionRequest{
		CardNumber:      "1234",
		Amount:          100000,
		PaymentMethod:   "Credit Card",
		MerchantID:      utils.PtrInt(1),
		TransactionTime: now,
	}

	expectedTransaction := &record.TransactionRecord{
		ID:              1,
		CardNumber:      "1234",
		Amount:          100000,
		PaymentMethod:   "Credit Card",
		MerchantID:      1,
		TransactionTime: now.Format(time.RFC3339),
		CreatedAt:       now.Format(time.RFC3339),
		UpdatedAt:       now.Format(time.RFC3339),
		DeletedAt:       nil,
	}

	mockRepo.EXPECT().CreateTransaction(&request).Return(expectedTransaction, nil)

	result, err := mockRepo.CreateTransaction(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTransaction, result)
}

func TestCreateTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	now := time.Now()
	request := requests.CreateTransactionRequest{
		CardNumber:      "1234",
		Amount:          100000,
		PaymentMethod:   "Credit Card",
		MerchantID:      utils.PtrInt(1),
		TransactionTime: now,
	}

	mockRepo.EXPECT().CreateTransaction(&request).Return(nil, fmt.Errorf("failed to create transaction"))

	result, err := mockRepo.CreateTransaction(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "failed to create transaction")
}

func TestCreateTransaction_ValidationError(t *testing.T) {
	request := requests.CreateTransactionRequest{
		CardNumber:      "",
		Amount:          0,
		PaymentMethod:   "",
		MerchantID:      nil,
		TransactionTime: time.Now(),
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payment method not found")
	assert.NotContains(t, err.Error(), "Field validation for 'CardNumber' failed on the 'required' tag")
	assert.NotContains(t, err.Error(), "Field validation for 'Amount' failed on the 'required' tag")
	assert.NotContains(t, err.Error(), "Field validation for 'PaymentMethod' failed on the 'required' tag")
	assert.NotContains(t, err.Error(), "Field validation for 'MerchantID' failed on the 'required' tag")
	assert.NotContains(t, err.Error(), "Field validation for 'TransactionTime' failed on the 'required' tag")
}

func TestUpdateTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	now := time.Now()
	request := requests.UpdateTransactionRequest{
		TransactionID:   1,
		CardNumber:      "1234",
		Amount:          150000,
		PaymentMethod:   "Credit Card",
		MerchantID:      utils.PtrInt(1),
		TransactionTime: now,
	}

	updatedTransaction := &record.TransactionRecord{
		ID:              1,
		CardNumber:      "1234",
		Amount:          150000,
		PaymentMethod:   "Credit Card",
		MerchantID:      1,
		TransactionTime: now.Format(time.RFC3339),
		CreatedAt:       "2024-12-24T10:00:00Z",
		UpdatedAt:       time.Now().Format(time.RFC3339),
		DeletedAt:       nil,
	}

	mockRepo.EXPECT().UpdateTransaction(&request).Return(updatedTransaction, nil)

	result, err := mockRepo.UpdateTransaction(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedTransaction, result)
}

func TestUpdateTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	now := time.Now()
	request := requests.UpdateTransactionRequest{
		TransactionID:   1,
		CardNumber:      "1234",
		Amount:          150000,
		PaymentMethod:   "Credit Card",
		MerchantID:      utils.PtrInt(1),
		TransactionTime: now,
	}

	mockRepo.EXPECT().UpdateTransaction(&request).Return(nil, fmt.Errorf("failed to update transaction"))

	result, err := mockRepo.UpdateTransaction(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "failed to update transaction")
}

func TestUpdateTransaction_ValidationError(t *testing.T) {
	request := requests.UpdateTransactionRequest{
		TransactionID:   0,
		CardNumber:      "",
		Amount:          0,
		PaymentMethod:   "",
		MerchantID:      nil,
		TransactionTime: time.Time{},
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payment method not found")
	assert.NotContains(t, err.Error(), "Field validation for 'TransactionID' failed on the 'required' tag")
	assert.NotContains(t, err.Error(), "Field validation for 'CardNumber' failed on the 'required' tag")
	assert.NotContains(t, err.Error(), "Field validation for 'Amount' failed on the 'required' tag")
	assert.NotContains(t, err.Error(), "Field validation for 'PaymentMethod' failed on the 'required' tag")
	assert.NotContains(t, err.Error(), "Field validation for 'MerchantID' failed on the 'required' tag")
	assert.NotContains(t, err.Error(), "Field validation for 'TransactionTime' failed on the 'required' tag")
}

func TestTrashedTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	transactionID := 1

	now := time.Now().Format(time.RFC3339)

	transaction := &record.TransactionRecord{
		ID:              transactionID,
		CardNumber:      "1234",
		Amount:          50000,
		PaymentMethod:   "credit",
		MerchantID:      123,
		TransactionTime: now,
		DeletedAt:       &now,
	}

	mockRepo.EXPECT().TrashedTransaction(transactionID).Return(transaction, nil)

	result, err := mockRepo.TrashedTransaction(transactionID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, transaction.ID, result.ID)
	assert.NotNil(t, result.DeletedAt)
}

func TestTrashedTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)
	transactionID := 999

	mockRepo.EXPECT().TrashedTransaction(transactionID).Return(nil, fmt.Errorf("transaction not found"))

	result, err := mockRepo.TrashedTransaction(transactionID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRestoreTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockTransactionRepository(ctrl)
	transactionID := 1
	transaction := &record.TransactionRecord{
		ID:              transactionID,
		CardNumber:      "1234",
		Amount:          50000,
		PaymentMethod:   "credit",
		MerchantID:      123,
		TransactionTime: time.Now().Format(time.RFC3339),
		DeletedAt:       nil,
	}

	mockRepo.EXPECT().RestoreTransaction(transactionID).Return(transaction, nil)

	result, err := mockRepo.RestoreTransaction(transactionID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.DeletedAt)
}

func TestRestoreTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockTransactionRepository(ctrl)
	transactionID := 999

	mockRepo.EXPECT().RestoreTransaction(transactionID).Return(nil, fmt.Errorf("transaction not found"))

	result, err := mockRepo.RestoreTransaction(transactionID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDeleteTransactionPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	transactionID := 1

	mockRepo.EXPECT().DeleteTransactionPermanent(transactionID).Return(
		nil,
	)

	err := mockRepo.DeleteTransactionPermanent(transactionID)

	assert.NoError(t, err)
}

func TestDeleteTransactionPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockTransactionRepository(ctrl)
	transactionID := 999

	mockRepo.EXPECT().DeleteTransactionPermanent(transactionID).Return(fmt.Errorf("transaction not found"))

	err := mockRepo.DeleteTransactionPermanent(transactionID)
	assert.NotNil(t, err)
	assert.Equal(t, "transaction not found", err.Error())
}
