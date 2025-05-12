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

func TestFindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockTransferRepository(ctrl)
	search := "user1"
	page := 1
	pageSize := 10
	transferRecords := []*record.TransferRecord{
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

	mockRepo.EXPECT().FindAll(search, page, pageSize).Return(transferRecords, 2, nil)

	result, totalRecords, err := mockRepo.FindAll(search, page, pageSize)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, totalRecords)
}

func TestFindAll_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockTransferRepository(ctrl)
	search := "user1"
	page := 1
	pageSize := 10

	mockRepo.EXPECT().FindAll(search, page, pageSize).Return(nil, 0, fmt.Errorf("database error"))

	result, totalRecords, err := mockRepo.FindAll(search, page, pageSize)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, totalRecords)
}

func TestFindAll_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockTransferRepository(ctrl)
	search := "user999"
	page := 1
	pageSize := 10

	mockRepo.EXPECT().FindAll(search, page, pageSize).Return(nil, 0, nil)

	result, totalRecords, err := mockRepo.FindAll(search, page, pageSize)

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 0, totalRecords)
}

func TestFindByIdTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockTransferRepository(ctrl)
	transferRecord := &record.TransferRecord{
		ID:             1,
		TransferFrom:   "user1",
		TransferTo:     "user2",
		TransferAmount: 1000,
		TransferTime:   time.Now().Format(time.RFC3339),
		CreatedAt:      time.Now().Format(time.RFC3339),
		UpdatedAt:      time.Now().Format(time.RFC3339),
	}

	mockRepo.EXPECT().FindById(1).Return(transferRecord, nil)

	result, err := mockRepo.FindById(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "user1", result.TransferFrom)
	assert.Equal(t, "user2", result.TransferTo)
}

func TestFindByIdTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	mockRepo.EXPECT().FindById(999).Return(nil, fmt.Errorf("transfer record with ID 999 not found"))

	result, err := mockRepo.FindById(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestFindByActiveTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	activeTransfers := []*record.TransferRecord{
		{
			ID:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
			DeletedAt:      nil,
		},
		{
			ID:             2,
			TransferFrom:   "user3",
			TransferTo:     "user4",
			TransferAmount: 2000,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
			DeletedAt:      nil,
		},
	}

	search := "user1"
	page := 1
	pageSize := 10
	expected := 2

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(activeTransfers, 2, nil)

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.NoError(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, "user1", result[0].TransferFrom)
}

func TestFindByActiveTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	search := "user1"
	page := 1
	pageSize := 10
	expected := 0

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(nil, 0, fmt.Errorf("failed to retrieve active transfer records"))

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to retrieve active transfer records")
}

func TestFindByTrashedTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	trashedTransfers := []*record.TransferRecord{
		{
			ID:             3,
			TransferFrom:   "user5",
			TransferTo:     "user6",
			TransferAmount: 1500,
			TransferTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
			DeletedAt:      utils.PtrString(time.Now().Format(time.RFC3339)),
		},
	}

	search := "user1"
	page := 1
	pageSize := 10
	expected := 1

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(trashedTransfers, 1, nil)

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, 3, result[0].ID)
	assert.NotNil(t, result[0].DeletedAt)
}

func TestFindByTrashedTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	search := "user1"
	page := 1
	pageSize := 10
	expected := 0

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, 0, fmt.Errorf("failed to retrieve trashed transfer records"))

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to retrieve trashed transfer records")
}

func TestFindTransferByTransferFrom_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

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
			DeletedAt:      nil,
		},
	}

	mockRepo.EXPECT().FindTransferByTransferFrom(transferFrom).Return(transfers, nil)

	result, err := mockRepo.FindTransferByTransferFrom(transferFrom)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "user1", result[0].TransferFrom)
}

func TestFindTransferByTransferFrom_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	transferFrom := "user1"

	mockRepo.EXPECT().FindTransferByTransferFrom(transferFrom).Return(nil, fmt.Errorf("failed to retrieve transfers from %s", transferFrom))

	result, err := mockRepo.FindTransferByTransferFrom(transferFrom)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to retrieve transfers from user1")
}

func TestFindTransferByTransferFrom_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	transferFrom := "user1"

	mockRepo.EXPECT().FindTransferByTransferFrom(transferFrom).Return([]*record.TransferRecord{}, nil)

	result, err := mockRepo.FindTransferByTransferFrom(transferFrom)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestFindTransferByTransferTo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

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
			DeletedAt:      nil,
		},
	}

	mockRepo.EXPECT().FindTransferByTransferTo(transferTo).Return(transfers, nil)

	result, err := mockRepo.FindTransferByTransferTo(transferTo)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "user2", result[0].TransferTo)
}

func TestFindTransferByTransferTo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	transferTo := "user2"

	mockRepo.EXPECT().FindTransferByTransferTo(transferTo).Return(nil, fmt.Errorf("failed to retrieve transfers to %s", transferTo))

	result, err := mockRepo.FindTransferByTransferTo(transferTo)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to retrieve transfers to user2")
}

func TestFindTransferByTransferTo_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	transferTo := "user2"

	mockRepo.EXPECT().FindTransferByTransferTo(transferTo).Return([]*record.TransferRecord{}, nil)

	result, err := mockRepo.FindTransferByTransferTo(transferTo)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestCountTransfersByDate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	date := "2024-12-25"
	expectedCount := 5

	mockRepo.EXPECT().CountTransfersByDate(date).Return(expectedCount, nil)

	count, err := mockRepo.CountTransfersByDate(date)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestCountTransfersByDate_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	date := "2024-12-25"

	mockRepo.EXPECT().CountTransfersByDate(date).Return(0, fmt.Errorf("failed to count transfers for date %s", date))

	count, err := mockRepo.CountTransfersByDate(date)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Contains(t, err.Error(), "failed to count transfers for date 2024-12-25")
}

func TestCountAllTransfers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransferRepository(ctrl)

	expectedCount := int64(50)
	expectedCountPtr := &expectedCount

	mockRepo.EXPECT().CountAllTransfers().Return(expectedCountPtr, nil)

	count, err := mockRepo.CountAllTransfers()

	assert.NoError(t, err)
	assert.Equal(t, expectedCountPtr, count)
}

func TestCountAllTransfers_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransferRepository(ctrl)

	expectedError := fmt.Errorf("failed to count all transfers")
	mockRepo.EXPECT().CountAllTransfers().Return(nil, expectedError)

	count, err := mockRepo.CountAllTransfers()

	assert.Error(t, err)
	assert.Nil(t, count)
	assert.Contains(t, err.Error(), "failed to count all transfers")
}

func TestCreateTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	request := requests.CreateTransferRequest{
		TransferFrom:   "Account1",
		TransferTo:     "Account2",
		TransferAmount: 100000,
	}

	expectedTransfer := &record.TransferRecord{
		ID:             1,
		TransferFrom:   "Account1",
		TransferTo:     "Account2",
		TransferAmount: 100000,
		TransferTime:   time.Now().Format(time.RFC3339),
	}

	mockRepo.EXPECT().CreateTransfer(&request).Return(expectedTransfer, nil)

	result, err := mockRepo.CreateTransfer(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTransfer.TransferFrom, result.TransferFrom)
	assert.Equal(t, expectedTransfer.TransferTo, result.TransferTo)
	assert.Equal(t, expectedTransfer.TransferAmount, result.TransferAmount)
}

func TestCreateTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	request := requests.CreateTransferRequest{
		TransferFrom:   "Account1",
		TransferTo:     "Account2",
		TransferAmount: 100000,
	}

	mockRepo.EXPECT().CreateTransfer(&request).Return(nil, fmt.Errorf("failed to create transfer"))

	result, err := mockRepo.CreateTransfer(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create transfer")
}

func TestCreateTransfer_ValidationError(t *testing.T) {
	request := requests.CreateTransferRequest{
		TransferFrom:   "",
		TransferTo:     "",
		TransferAmount: 0,
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'TransferFrom' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TransferTo' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TransferAmount' failed on the 'required' tag")

}

func TestUpdateTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	request := requests.UpdateTransferRequest{
		TransferID:     1,
		TransferFrom:   "Account1",
		TransferTo:     "Account2",
		TransferAmount: 150000,
	}

	expectedTransfer := &record.TransferRecord{
		ID:             1,
		TransferFrom:   "Account1",
		TransferTo:     "Account2",
		TransferAmount: 150000,
		TransferTime:   time.Now().Format(time.RFC3339),
	}

	mockRepo.EXPECT().UpdateTransfer(&request).Return(expectedTransfer, nil)

	result, err := mockRepo.UpdateTransfer(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTransfer.TransferFrom, result.TransferFrom)
	assert.Equal(t, expectedTransfer.TransferTo, result.TransferTo)
	assert.Equal(t, expectedTransfer.TransferAmount, result.TransferAmount)
}

func TestUpdateTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	request := requests.UpdateTransferRequest{
		TransferID:     1,
		TransferFrom:   "Account1",
		TransferTo:     "Account2",
		TransferAmount: 150000,
	}

	mockRepo.EXPECT().UpdateTransfer(&request).Return(nil, fmt.Errorf("failed to update transfer"))

	result, err := mockRepo.UpdateTransfer(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update transfer")
}

func TestUpdateTransfer_ValidationError(t *testing.T) {
	request := requests.UpdateTransferRequest{
		TransferID:     0,
		TransferFrom:   "",
		TransferTo:     "",
		TransferAmount: 0,
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'TransferID' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TransferFrom' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TransferTo' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TransferAmount' failed on the 'required' tag")
}

func TestUpdateTransferAmount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	request := requests.UpdateTransferAmountRequest{
		TransferID:     1,
		TransferAmount: 200000,
	}

	expectedTransfer := &record.TransferRecord{
		ID:             1,
		TransferFrom:   "Account1",
		TransferTo:     "Account2",
		TransferAmount: 200000,
		TransferTime:   time.Now().Format(time.RFC3339),
	}

	mockRepo.EXPECT().UpdateTransferAmount(&request).Return(expectedTransfer, nil)

	result, err := mockRepo.UpdateTransferAmount(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTransfer.TransferAmount, result.TransferAmount)
}

func TestUpdateTransferAmount_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	request := requests.UpdateTransferAmountRequest{
		TransferID:     1,
		TransferAmount: 200000,
	}

	mockRepo.EXPECT().UpdateTransferAmount(&request).Return(nil, fmt.Errorf("failed to update transfer amount"))

	result, err := mockRepo.UpdateTransferAmount(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update transfer amount")
}

func TestUpdateTransferAmount_ValidationError(t *testing.T) {
	request := requests.UpdateTransferAmountRequest{
		TransferID:     0,
		TransferAmount: -50000,
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'TransferID' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TransferAmount' failed on the 'gt' tag")
}

func TestTrashedTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	transferID := 1

	trashedTransfer := &record.TransferRecord{
		ID:             1,
		TransferFrom:   "Account1",
		TransferTo:     "Account2",
		TransferAmount: 100000,
		TransferTime:   time.Now().Format(time.RFC3339),
		DeletedAt:      new(string),
	}

	mockRepo.EXPECT().TrashedTransfer(transferID).Return(trashedTransfer, nil)

	result, err := mockRepo.TrashedTransfer(transferID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, trashedTransfer.ID, result.ID)
	assert.NotNil(t, result.DeletedAt)
}

func TestTrashedTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	transferID := 1

	mockRepo.EXPECT().TrashedTransfer(transferID).Return(nil, fmt.Errorf("failed to trash transfer"))

	result, err := mockRepo.TrashedTransfer(transferID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to trash transfer")
}

func TestRestoreTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	transferID := 1

	restoredTransfer := &record.TransferRecord{
		ID:             1,
		TransferFrom:   "Account1",
		TransferTo:     "Account2",
		TransferAmount: 100000,
		TransferTime:   time.Now().Format(time.RFC3339),
		DeletedAt:      nil,
	}

	mockRepo.EXPECT().RestoreTransfer(transferID).Return(restoredTransfer, nil)

	result, err := mockRepo.RestoreTransfer(transferID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.DeletedAt)
}

func TestRestoreTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	transferID := 1

	mockRepo.EXPECT().RestoreTransfer(transferID).Return(nil, fmt.Errorf("failed to restore transfer"))

	result, err := mockRepo.RestoreTransfer(transferID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to restore transfer")
}

func TestDeleteTransferPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	transferID := 1

	mockRepo.EXPECT().DeleteTransferPermanent(transferID).Return(nil)

	err := mockRepo.DeleteTransferPermanent(transferID)

	assert.NoError(t, err)
}

func TestDeleteTransferPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTransferRepository(ctrl)

	transferID := 1

	mockRepo.EXPECT().DeleteTransferPermanent(transferID).Return(fmt.Errorf("failed to delete transfer permanently"))

	err := mockRepo.DeleteTransferPermanent(transferID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete transfer permanently")
}
