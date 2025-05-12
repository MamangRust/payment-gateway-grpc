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

func TestFindAll_Withdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	search := "card_1234"
	page := 1
	pageSize := 10

	expectedWithdrawRecords := []*record.WithdrawRecord{
		{
			ID:             1,
			CardNumber:     "card_1234",
			WithdrawAmount: 100000,
			WithdrawTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
			DeletedAt:      nil,
		},
		{
			ID:             2,
			CardNumber:     "card_1234",
			WithdrawAmount: 200000,
			WithdrawTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
			DeletedAt:      nil,
		},
	}

	mockRepo.EXPECT().FindAll(search, page, pageSize).Return(expectedWithdrawRecords, 2, nil)

	result, count, err := mockRepo.FindAll(search, page, pageSize)

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, result, 2)
	assert.Equal(t, "card_1234", result[0].CardNumber)
}

func TestFindAll_Withdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	search := "card_1234"
	page := 1
	pageSize := 10

	mockRepo.EXPECT().FindAll(search, page, pageSize).Return(nil, 0, fmt.Errorf("failed to fetch withdraw records"))

	result, count, err := mockRepo.FindAll(search, page, pageSize)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, count)
	assert.Contains(t, err.Error(), "failed to fetch withdraw records")
}

func TestFindAll_Withdraw_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	search := "card_9999"
	page := 1
	pageSize := 10

	mockRepo.EXPECT().FindAll(search, page, pageSize).Return([]*record.WithdrawRecord{}, 0, nil)

	result, count, err := mockRepo.FindAll(search, page, pageSize)

	assert.NoError(t, err)
	assert.Equal(t, 0, count)
	assert.Len(t, result, 0)
}

func TestFindById_Withdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	id := 1

	expectedWithdrawRecord := &record.WithdrawRecord{
		ID:             id,
		CardNumber:     "card_1234",
		WithdrawAmount: 150000,
		WithdrawTime:   time.Now().Format(time.RFC3339),
		CreatedAt:      time.Now().Format(time.RFC3339),
		UpdatedAt:      time.Now().Format(time.RFC3339),
		DeletedAt:      nil,
	}

	mockRepo.EXPECT().FindById(id).Return(expectedWithdrawRecord, nil)

	result, err := mockRepo.FindById(id)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, id, result.ID)
	assert.Equal(t, "card_1234", result.CardNumber)
	assert.Equal(t, 150000, result.WithdrawAmount)
}

func TestFindById_Withdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	id := 1

	mockRepo.EXPECT().FindById(id).Return(nil, fmt.Errorf("failed to fetch withdraw record"))

	result, err := mockRepo.FindById(id)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to fetch withdraw record")
}

func TestFindByActive_Withdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	expectedRecords := []*record.WithdrawRecord{
		{
			ID:             1,
			CardNumber:     "card_1234",
			WithdrawAmount: 150000,
			WithdrawTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
			DeletedAt:      nil,
		},
		{
			ID:             2,
			CardNumber:     "card_5678",
			WithdrawAmount: 200000,
			WithdrawTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
			DeletedAt:      nil,
		},
	}

	search := "user1"
	page := 1
	pageSize := 1
	expected := 2

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(expectedRecords, expected, nil)

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(expectedRecords), len(result))
}

func TestFindByActive_Withdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	search := "user1"
	page := 1
	pageSize := 1
	expected := 0

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(nil, expected, fmt.Errorf("database error"))

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
}

func TestFindByActive_WithdrawRecord_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	search := "user1"
	page := 1
	pageSize := 1
	expected := 0

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return([]*record.WithdrawRecord{}, expected, nil)

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

func TestFindByTrashed_Withdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	expectedRecords := []*record.WithdrawRecord{
		{
			ID:             3,
			CardNumber:     "card_9999",
			WithdrawAmount: 100000,
			WithdrawTime:   time.Now().Format(time.RFC3339),
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
			DeletedAt:      utils.PtrString(time.Now().Format(time.RFC3339)),
		},
	}

	search := "user1"
	page := 1
	pageSize := 1
	expected := 2

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(expectedRecords, expected, nil)

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(expectedRecords), len(result))
}

func TestFindByTrashed_Withdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	search := "user1"
	page := 1
	pageSize := 1
	expected := 0

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, expected, fmt.Errorf("database error"))

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
}

func TestFindByTrashed_Withdraw_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	search := "user1"
	page := 1
	pageSize := 1
	expected := 0

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return([]*record.WithdrawRecord{}, expected, nil)

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

func TestCountWithdrawByActiveDate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	expectedCount := int64(2)
	date := time.Now()

	mockRepo.EXPECT().CountActiveByDate(date).Return((expectedCount), nil)

	count, err := mockRepo.CountActiveByDate(date)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestCountWithdrawByActiveDate_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	expectedError := fmt.Errorf("database error")
	date := time.Now()

	mockRepo.EXPECT().CountActiveByDate(date).Return(int64(0), expectedError)

	count, err := mockRepo.CountActiveByDate(date)

	assert.Error(t, err)
	assert.Equal(t, expectedError.Error(), err.Error())

	assert.Equal(t, int64(0), count)
}

func TestCreateWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	request := requests.CreateWithdrawRequest{
		CardNumber:     "success_card",
		WithdrawAmount: 100000,
		WithdrawTime:   time.Now(),
	}

	expectedRecord := &record.WithdrawRecord{
		ID:             1,
		CardNumber:     request.CardNumber,
		WithdrawAmount: request.WithdrawAmount,
		WithdrawTime:   request.WithdrawTime.Format(time.RFC3339),
		CreatedAt:      time.Now().Format(time.RFC3339),
		UpdatedAt:      time.Now().Format(time.RFC3339),
	}

	mockRepo.EXPECT().CreateWithdraw(&request).Return(expectedRecord, nil)

	record, err := mockRepo.CreateWithdraw(&request)

	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, expectedRecord, record)
}

func TestCreateWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	request := requests.CreateWithdrawRequest{
		CardNumber:     "failure_card",
		WithdrawAmount: 100000,
		WithdrawTime:   time.Now(),
	}

	mockRepo.EXPECT().CreateWithdraw(&request).Return(nil, fmt.Errorf("failed to create withdraw record"))

	record, err := mockRepo.CreateWithdraw(&request)

	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Equal(t, "failed to create withdraw record", err.Error())
}

func TestCreateWithdraw_ValidationError(t *testing.T) {
	request := requests.CreateWithdrawRequest{
		CardNumber:     "",
		WithdrawAmount: 0,
		WithdrawTime:   time.Time{},
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'CardNumber' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'WithdrawAmount' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'WithdrawTime' failed on the 'required' tag")
}

func TestUpdateWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	request := requests.UpdateWithdrawRequest{
		CardNumber:     "updated_card",
		WithdrawID:     1,
		WithdrawAmount: 100000,
		WithdrawTime:   time.Now(),
	}

	expectedRecord := &record.WithdrawRecord{
		ID:             request.WithdrawID,
		CardNumber:     request.CardNumber,
		WithdrawAmount: request.WithdrawAmount,
		WithdrawTime:   request.WithdrawTime.Format(time.RFC3339),
		CreatedAt:      time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		UpdatedAt:      time.Now().Format(time.RFC3339),
	}

	mockRepo.EXPECT().UpdateWithdraw(&request).Return(expectedRecord, nil)

	record, err := mockRepo.UpdateWithdraw(&request)

	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, expectedRecord, record)
}

func TestUpdateWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	request := requests.UpdateWithdrawRequest{
		CardNumber:     "failure_card",
		WithdrawID:     99,
		WithdrawAmount: 100000,
		WithdrawTime:   time.Now(),
	}

	mockRepo.EXPECT().UpdateWithdraw(&request).Return(nil, fmt.Errorf("withdraw record not found"))

	record, err := mockRepo.UpdateWithdraw(&request)

	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Equal(t, "withdraw record not found", err.Error())
}

func TestUpdateWithdraw_ValidationError(t *testing.T) {
	request := requests.UpdateWithdrawRequest{
		CardNumber:     "",
		WithdrawID:     0,
		WithdrawAmount: 0,
		WithdrawTime:   time.Time{},
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'CardNumber' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'WithdrawID' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'WithdrawAmount' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'WithdrawTime' failed on the 'required' tag")
}

func TestTrashedWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	expectedRecord := &record.WithdrawRecord{
		ID:             1,
		CardNumber:     "123456789",
		WithdrawAmount: 75000,
		DeletedAt:      utils.PtrString(time.Now().Format(time.RFC3339)),
	}

	mockRepo.EXPECT().TrashedWithdraw(1).Return(expectedRecord, nil)

	record, err := mockRepo.TrashedWithdraw(1)

	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, expectedRecord, record)
}

func TestTrashedWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	mockRepo.EXPECT().TrashedWithdraw(99).Return(nil, fmt.Errorf("withdraw record not found"))

	record, err := mockRepo.TrashedWithdraw(99)

	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Equal(t, "withdraw record not found", err.Error())
}

func TestRestoreWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	expectedRecord := &record.WithdrawRecord{
		ID:             1,
		CardNumber:     "123456789",
		WithdrawAmount: 75000,
		DeletedAt:      nil,
	}

	mockRepo.EXPECT().RestoreWithdraw(1).Return(expectedRecord, nil)

	record, err := mockRepo.RestoreWithdraw(1)

	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, expectedRecord, record)
}

func TestRestoreWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	mockRepo.EXPECT().RestoreWithdraw(99).Return(nil, fmt.Errorf("withdraw record not found"))

	record, err := mockRepo.RestoreWithdraw(99)

	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Equal(t, "withdraw record not found", err.Error())
}

func TestDeleteWithdrawPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	mockRepo.EXPECT().DeleteWithdrawPermanent(1).Return(nil)

	err := mockRepo.DeleteWithdrawPermanent(1)

	assert.NoError(t, err)
}

func TestDeleteWithdrawPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWithdrawRepository(ctrl)

	mockRepo.EXPECT().DeleteWithdrawPermanent(99).Return(fmt.Errorf("withdraw record not found"))

	err := mockRepo.DeleteWithdrawPermanent(99)

	assert.Error(t, err)
	assert.Equal(t, "withdraw record not found", err.Error())
}
