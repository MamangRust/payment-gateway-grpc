package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindAllTopups_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

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

	mockRepo.EXPECT().FindAllTopups("", 1, 10).Return(topups, 2, nil)

	result, total, err := mockRepo.FindAllTopups("", 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, total)
	assert.Equal(t, topups, result)
}

func TestFindAllTopups_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	mockRepo.EXPECT().FindAllTopups("", 1, 10).Return(nil, 0, fmt.Errorf("database error"))

	result, total, err := mockRepo.FindAllTopups("", 1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, total)
	assert.Contains(t, err.Error(), "database error")
}

func TestFindAllTopups_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	mockRepo.EXPECT().FindAllTopups("", 1, 10).Return([]*record.TopupRecord{}, 0, nil)

	result, total, err := mockRepo.FindAllTopups("", 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, total)
	assert.Empty(t, result)
}

func TestFindByIdTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	expectedTopup := &record.TopupRecord{
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

	mockRepo.EXPECT().FindById(1).Return(expectedTopup, nil)

	result, err := mockRepo.FindById(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTopup, result)
}

func TestFindByIdTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	mockRepo.EXPECT().FindById(1).Return(nil, fmt.Errorf("database error"))

	result, err := mockRepo.FindById(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
}

func TestFindByCardNumber_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	expectedTopups := []*record.TopupRecord{
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
	}

	mockRepo.EXPECT().FindByCardNumber("1234").Return(expectedTopups, nil)

	result, err := mockRepo.FindByCardNumber("1234")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTopups, result)
}

func TestFindByCardNumber_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	mockRepo.EXPECT().FindByCardNumber("1234").Return(nil, fmt.Errorf("database error"))

	result, err := mockRepo.FindByCardNumber("1234")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
}

func TestFindByActiveTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	expectedTopups := []*record.TopupRecord{
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
	}
	page := 1
	pageSize := 10
	search := ""
	expected := 1

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(expectedTopups, 1, nil)

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.NoError(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTopups, result)
}

func TestFindByActiveTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(nil, 0, fmt.Errorf("database error"))

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
}

func TestFindByTrashedTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	expectedTopups := []*record.TopupRecord{
		{
			ID:          2,
			CardNumber:  "5678",
			TopupNo:     "TOPUP-002",
			TopupAmount: 60000,
			TopupMethod: "credit_card",
			TopupTime:   "2024-12-25T10:00:00Z",
			CreatedAt:   "2024-12-25T10:00:00Z",
			UpdatedAt:   "2024-12-25T10:30:00Z",
			DeletedAt:   nil,
		},
	}

	page := 1
	pageSize := 10
	search := ""
	expected := 1

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(expectedTopups, 1, nil)

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTopups, result)
}
func TestFindByTrashedTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, 0, fmt.Errorf("database error"))

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
}

func TestCountTopupsByDate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	expectedCount := 10
	date := "2024-12-25"

	mockRepo.EXPECT().CountTopupsByDate(date).Return(expectedCount, nil)

	result, err := mockRepo.CountTopupsByDate(date)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)
}

func TestCountTopupsByDate_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	date := "2024-12-25"

	mockRepo.EXPECT().CountTopupsByDate(date).Return(0, fmt.Errorf("database error"))

	result, err := mockRepo.CountTopupsByDate(date)

	assert.Error(t, err)
	assert.Equal(t, 0, result)
	assert.Contains(t, err.Error(), "database error")
}

func TestCountAllTopups_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	expectedCount := int64(100)
	expectedCountPtr := &expectedCount

	mockRepo.EXPECT().CountAllTopups().Return(expectedCountPtr, nil)

	result, err := mockRepo.CountAllTopups()

	assert.NoError(t, err)
	assert.Equal(t, expectedCountPtr, result)
}

func TestCountAllTopups_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	expectedError := fmt.Errorf("database error")
	mockRepo.EXPECT().CountAllTopups().Return(nil, expectedError)

	result, err := mockRepo.CountAllTopups()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
}

func TestCreateTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	request := requests.CreateTopupRequest{
		CardNumber:  "1234",
		TopupNo:     "TOPUP001",
		TopupAmount: 100000,
		TopupMethod: "Bank Transfer",
	}

	expectedTopup := &record.TopupRecord{
		ID:          1,
		CardNumber:  "1234",
		TopupNo:     "TOPUP001",
		TopupAmount: 100000,
		TopupMethod: "Bank Transfer",
		TopupTime:   "2024-12-25T10:00:00Z",
		CreatedAt:   "2024-12-25T10:00:00Z",
		UpdatedAt:   "2024-12-25T10:00:00Z",
	}

	mockRepo.EXPECT().CreateTopup(&request).Return(expectedTopup, nil)

	result, err := mockRepo.CreateTopup(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTopup, result)
}

func TestCreateTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	request := requests.CreateTopupRequest{
		CardNumber:  "1234",
		TopupNo:     "TOPUP001",
		TopupAmount: 100000,
		TopupMethod: "Bank Transfer",
	}

	mockRepo.EXPECT().CreateTopup(&request).Return(nil, fmt.Errorf("database error"))

	result, err := mockRepo.CreateTopup(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
}

func TestCreateTopup_ValidationError(t *testing.T) {
	request := requests.CreateTopupRequest{
		CardNumber:  "",
		TopupNo:     "",
		TopupAmount: 0,
		TopupMethod: "",
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'CardNumber' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TopupNo' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TopupAmount' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TopupMethod' failed on the 'required' tag")
}

func TestUpdateTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	request := requests.UpdateTopupRequest{
		CardNumber:  "1234",
		TopupID:     1,
		TopupAmount: 200000,
		TopupMethod: "Bank Transfer",
	}

	expectedTopup := &record.TopupRecord{
		ID:          1,
		CardNumber:  "1234",
		TopupNo:     "TOPUP001",
		TopupAmount: 200000,
		TopupMethod: "Bank Transfer",
		TopupTime:   "2024-12-25T10:00:00Z",
		CreatedAt:   "2024-12-25T10:00:00Z",
		UpdatedAt:   "2024-12-25T10:10:00Z",
	}

	mockRepo.EXPECT().UpdateTopup(&request).Return(expectedTopup, nil)

	result, err := mockRepo.UpdateTopup(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTopup, result)
}

func TestUpdateTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	request := requests.UpdateTopupRequest{
		CardNumber:  "1234",
		TopupID:     1,
		TopupAmount: 200000,
		TopupMethod: "Bank Transfer",
	}

	mockRepo.EXPECT().UpdateTopup(&request).Return(nil, fmt.Errorf("failed to update topup"))

	result, err := mockRepo.UpdateTopup(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update topup")
}

func TestUpdateTopup_ValidationError(t *testing.T) {
	request := requests.UpdateTopupRequest{
		CardNumber:  "",
		TopupID:     0,
		TopupAmount: 0,
		TopupMethod: "",
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'CardNumber' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TopupID' failed on the 'required' tag")

	assert.Contains(t, err.Error(), "Field validation for 'TopupAmount' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TopupMethod' failed on the 'required' tag")
}

func TestUpdateTopupAmount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	request := requests.UpdateTopupAmount{
		TopupID:     1,
		TopupAmount: 150000,
	}

	expectedTopup := &record.TopupRecord{
		ID:          1,
		CardNumber:  "1234",
		TopupNo:     "TOPUP001",
		TopupAmount: 150000,
		TopupMethod: "Bank Transfer",
		TopupTime:   "2024-12-25T12:00:00Z",
		CreatedAt:   "2024-12-25T10:00:00Z",
		UpdatedAt:   "2024-12-25T12:00:00Z",
	}

	mockRepo.EXPECT().UpdateTopupAmount(&request).Return(expectedTopup, nil)

	result, err := mockRepo.UpdateTopupAmount(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTopup.TopupAmount, result.TopupAmount)
	assert.Equal(t, expectedTopup.ID, result.ID)
}

func TestUpdateTopupAmount_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	request := requests.UpdateTopupAmount{
		TopupID:     1,
		TopupAmount: 150000,
	}

	mockRepo.EXPECT().UpdateTopupAmount(&request).Return(nil, fmt.Errorf("failed to update topup amount"))

	result, err := mockRepo.UpdateTopupAmount(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update topup amount")
}

func TestUpdateTopupAmount_ValidationError(t *testing.T) {
	request := requests.UpdateTopupAmount{
		TopupID:     0,
		TopupAmount: 0,
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'TopupID' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TopupAmount' failed on the 'required' tag")
}

func TestTrashedTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	expectedTopup := &record.TopupRecord{
		ID:          1,
		CardNumber:  "1234",
		TopupNo:     "TOPUP001",
		TopupAmount: 150000,
		TopupMethod: "Bank Transfer",
		TopupTime:   "2024-12-25T12:00:00Z",
		CreatedAt:   "2024-12-25T10:00:00Z",
		UpdatedAt:   "2024-12-25T12:00:00Z",
		DeletedAt:   new(string),
	}

	mockRepo.EXPECT().TrashedTopup(1).Return(expectedTopup, nil)

	result, err := mockRepo.TrashedTopup(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTopup.ID, result.ID)
	assert.NotNil(t, result.DeletedAt)
}

func TestTrashedTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	mockRepo.EXPECT().TrashedTopup(1).Return(nil, fmt.Errorf("failed to trash topup"))

	result, err := mockRepo.TrashedTopup(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to trash topup")
}

func TestRestoreTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	expectedTopup := &record.TopupRecord{
		ID:          1,
		CardNumber:  "1234",
		TopupNo:     "TOPUP001",
		TopupAmount: 150000,
		TopupMethod: "Bank Transfer",
		TopupTime:   "2024-12-25T12:00:00Z",
		CreatedAt:   "2024-12-25T10:00:00Z",
		UpdatedAt:   "2024-12-25T12:00:00Z",
		DeletedAt:   nil,
	}

	mockRepo.EXPECT().RestoreTopup(1).Return(expectedTopup, nil)

	result, err := mockRepo.RestoreTopup(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTopup.ID, result.ID)
	assert.Nil(t, result.DeletedAt)
}

func TestRestoreTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	mockRepo.EXPECT().RestoreTopup(1).Return(nil, fmt.Errorf("failed to restore topup"))

	result, err := mockRepo.RestoreTopup(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to restore topup")
}

func TestDeleteTopupPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	mockRepo.EXPECT().DeleteTopupPermanent(1).Return(nil)

	err := mockRepo.DeleteTopupPermanent(1)

	assert.NoError(t, err)
}

func TestDeleteTopupPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopupRepository(ctrl)

	mockRepo.EXPECT().DeleteTopupPermanent(1).Return(fmt.Errorf("failed to delete topup permanently"))

	err := mockRepo.DeleteTopupPermanent(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete topup permanently")
}
