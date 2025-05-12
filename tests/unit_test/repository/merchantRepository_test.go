package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/tests/utils"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindAllMerchants_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	merchants := []*record.MerchantRecord{
		{
			ID:        1,
			Name:      "Merchant 1",
			ApiKey:    "apikey1",
			UserID:    101,
			Status:    "active",
			CreatedAt: "2024-01-01T10:00:00Z",
			UpdatedAt: "2024-01-01T10:00:00Z",
		},
		{
			ID:        2,
			Name:      "Merchant 2",
			ApiKey:    "apikey2",
			UserID:    102,
			Status:    "inactive",
			CreatedAt: "2024-02-01T10:00:00Z",
			UpdatedAt: "2024-02-01T10:00:00Z",
		},
	}

	mockRepo.EXPECT().FindAllMerchants("", 1, 10).Return(merchants, 2, nil)

	result, total, err := mockRepo.FindAllMerchants("", 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, merchants, result)
	assert.Equal(t, 2, total)
}

func TestFindAllMerchants_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	mockRepo.EXPECT().FindAllMerchants("", 1, 10).Return(nil, 0, errors.New("database error"))

	result, total, err := mockRepo.FindAllMerchants("", 1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, total)
}

func TestFindAllMerchants_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	mockRepo.EXPECT().FindAllMerchants("", 1, 10).Return([]*record.MerchantRecord{}, 0, nil)

	result, total, err := mockRepo.FindAllMerchants("", 1, 10)

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 0, total)
}

func TestFindByIdMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	expectedMerchant := &record.MerchantRecord{
		ID:        1,
		Name:      "Merchant 1",
		ApiKey:    "apikey1",
		UserID:    101,
		Status:    "active",
		CreatedAt: "2024-01-01T10:00:00Z",
		UpdatedAt: "2024-01-01T10:00:00Z",
	}

	merchantID := 1

	mockRepo.EXPECT().FindById(merchantID).Return(expectedMerchant, nil)

	result, err := mockRepo.FindById(merchantID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMerchant, result)
}
func TestFindById_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	merchantID := 999

	mockRepo.EXPECT().FindById(merchantID).Return(nil, errors.New("merchant not found"))

	result, err := mockRepo.FindById(merchantID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "merchant not found")
}

func TestFindByActiveMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	expectedMerchants := []*record.MerchantRecord{
		{
			ID:        1,
			Name:      "Merchant 1",
			ApiKey:    "apikey1",
			UserID:    101,
			Status:    "active",
			CreatedAt: "2024-01-01T10:00:00Z",
			UpdatedAt: "2024-01-01T10:00:00Z",
		},
		{
			ID:        2,
			Name:      "Merchant 2",
			ApiKey:    "apikey2",
			UserID:    102,
			Status:    "active",
			CreatedAt: "2024-01-02T10:00:00Z",
			UpdatedAt: "2024-01-02T10:00:00Z",
		},
	}

	page := 1
	pageSize := 10
	search := ""
	expected := 2

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(expectedMerchants, 2, nil)

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMerchants, result)
}

func TestFindByActiveMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(nil, 0, fmt.Errorf("failed to fetch active merchants"))

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expected, totalRecord)
	assert.Contains(t, err.Error(), "failed to fetch active merchants")
}

func TestFindByTrashedMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	expectedMerchants := []*record.MerchantRecord{
		{
			ID:        3,
			Name:      "Merchant 3",
			ApiKey:    "apikey3",
			UserID:    103,
			Status:    "trashed",
			CreatedAt: "2024-01-03T10:00:00Z",
			UpdatedAt: "2024-01-03T10:00:00Z",
			DeletedAt: utils.PtrString("2024-01-03T10:00:00Z"),
		},
	}

	page := 1
	pageSize := 10
	search := ""
	expected := 1

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(expectedMerchants, 1, nil)

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, expectedMerchants, result)
}

func TestFindByTrashedMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, 0, fmt.Errorf("failed to fetch trashed merchants"))

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Error(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to fetch trashed merchants")
}

func TestFindByName_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	name := "Merchant 1"
	expectedMerchant := &record.MerchantRecord{
		ID:        1,
		Name:      name,
		ApiKey:    "apikey123",
		UserID:    101,
		Status:    "active",
		CreatedAt: "2024-01-01T10:00:00Z",
		UpdatedAt: "2024-01-01T10:00:00Z",
	}

	mockRepo.EXPECT().FindByName(name).Return(expectedMerchant, nil)

	result, err := mockRepo.FindByName(name)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMerchant, result)
}

func TestFindByName_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	name := "Nonexistent Merchant"

	mockRepo.EXPECT().FindByName(name).Return(nil, fmt.Errorf("merchant not found"))

	result, err := mockRepo.FindByName(name)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "merchant not found")
}

func TestFindByApiKey_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	apiKey := "apikey123"
	expectedMerchant := &record.MerchantRecord{
		ID:        1,
		Name:      "Merchant 1",
		ApiKey:    apiKey,
		UserID:    101,
		Status:    "active",
		CreatedAt: "2024-01-01T10:00:00Z",
		UpdatedAt: "2024-01-01T10:00:00Z",
	}

	mockRepo.EXPECT().FindByApiKey(apiKey).Return(expectedMerchant, nil)

	result, err := mockRepo.FindByApiKey(apiKey)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMerchant, result)
}

func TestFindByApiKey_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	apiKey := "nonexistent-api-key"

	mockRepo.EXPECT().FindByApiKey(apiKey).Return(nil, fmt.Errorf("merchant not found"))

	result, err := mockRepo.FindByApiKey(apiKey)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "merchant not found")
}

func TestFindByMerchantUserId_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	userID := 101
	expectedMerchants := []*record.MerchantRecord{
		{
			ID:        1,
			Name:      "Merchant 1",
			ApiKey:    "apikey123",
			UserID:    userID,
			Status:    "active",
			CreatedAt: "2024-01-01T10:00:00Z",
			UpdatedAt: "2024-01-01T10:00:00Z",
		},
		{
			ID:        2,
			Name:      "Merchant 2",
			ApiKey:    "apikey456",
			UserID:    userID,
			Status:    "inactive",
			CreatedAt: "2024-01-02T10:00:00Z",
			UpdatedAt: "2024-01-02T10:00:00Z",
		},
	}

	mockRepo.EXPECT().FindByMerchantUserId(userID).Return(expectedMerchants, nil)

	result, err := mockRepo.FindByMerchantUserId(userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMerchants, result)
}

func TestFindByMerchantUserId_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	userID := 999

	mockRepo.EXPECT().FindByMerchantUserId(userID).Return(nil, fmt.Errorf("merchants not found for user ID %d", userID))

	result, err := mockRepo.FindByMerchantUserId(userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "merchants not found")
}

func TestCreateMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	request := requests.CreateMerchantRequest{
		Name:   "Merchant 1",
		UserID: 101,
	}

	expectedMerchant := &record.MerchantRecord{
		ID:        1,
		Name:      request.Name,
		UserID:    request.UserID,
		ApiKey:    "generatedapikey",
		Status:    "active",
		CreatedAt: "2024-12-24T10:00:00Z",
		UpdatedAt: "2024-12-24T10:00:00Z",
	}

	mockRepo.EXPECT().CreateMerchant(&request).Return(expectedMerchant, nil)

	result, err := mockRepo.CreateMerchant(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMerchant, result)
}

func TestCreateMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	request := requests.CreateMerchantRequest{
		Name:   "Merchant 1",
		UserID: 101,
	}

	mockRepo.EXPECT().CreateMerchant(&request).Return(nil, fmt.Errorf("failed to create merchant"))

	result, err := mockRepo.CreateMerchant(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create merchant")
}

func TestCreateMerchant_ValidationError(t *testing.T) {
	request := requests.CreateMerchantRequest{
		Name:   "",
		UserID: 0,
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'Name' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'UserID' failed on the 'required' tag")
}

func TestUpdateMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	request := requests.UpdateMerchantRequest{
		MerchantID: 1,
		Name:       "Updated Merchant",
		UserID:     101,
	}

	expectedMerchant := &record.MerchantRecord{
		ID:        request.MerchantID,
		Name:      request.Name,
		UserID:    request.UserID,
		ApiKey:    "existingapikey",
		Status:    "active",
		CreatedAt: "2024-12-23T10:00:00Z",
		UpdatedAt: "2024-12-24T10:00:00Z",
	}

	mockRepo.EXPECT().UpdateMerchant(&request).Return(expectedMerchant, nil)

	result, err := mockRepo.UpdateMerchant(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMerchant, result)
}

func TestUpdateMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	request := requests.UpdateMerchantRequest{
		MerchantID: 1,
		Name:       "Updated Merchant",
		UserID:     101,
	}

	mockRepo.EXPECT().UpdateMerchant(&request).Return(nil, fmt.Errorf("failed to update merchant"))

	result, err := mockRepo.UpdateMerchant(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update merchant")
}

func TestUpdateMerchant_ValidationError(t *testing.T) {
	request := requests.UpdateMerchantRequest{
		MerchantID: 0,
		Name:       "",
		UserID:     0,
		Status:     "",
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'MerchantID' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'Name' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'UserID' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'Status' failed on the 'required' tag")
}

func TestTrashedMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	merchantID := 1
	expectedMerchant := &record.MerchantRecord{
		ID:        merchantID,
		Name:      "Merchant 1",
		UserID:    101,
		ApiKey:    "apikey123",
		Status:    "trashed",
		CreatedAt: "2024-12-23T10:00:00Z",
		UpdatedAt: "2024-12-24T10:00:00Z",
		DeletedAt: utils.PtrString("2024-12-24T10:00:00Z"),
	}

	mockRepo.EXPECT().TrashedMerchant(merchantID).Return(expectedMerchant, nil)

	result, err := mockRepo.TrashedMerchant(merchantID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMerchant, result)
}

func TestTrashedMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	merchantID := 1

	mockRepo.EXPECT().TrashedMerchant(merchantID).Return(nil, fmt.Errorf("failed to trash merchant"))

	result, err := mockRepo.TrashedMerchant(merchantID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to trash merchant")
}

func TestRestoreMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	merchantID := 1
	expectedMerchant := &record.MerchantRecord{
		ID:        merchantID,
		Name:      "Merchant 1",
		UserID:    101,
		ApiKey:    "apikey123",
		Status:    "active",
		CreatedAt: "2024-12-23T10:00:00Z",
		UpdatedAt: "2024-12-24T12:00:00Z",
		DeletedAt: nil,
	}

	mockRepo.EXPECT().RestoreMerchant(merchantID).Return(expectedMerchant, nil)

	result, err := mockRepo.RestoreMerchant(merchantID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMerchant, result)
}

func TestRestoreMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	merchantID := 1

	mockRepo.EXPECT().RestoreMerchant(merchantID).Return(nil, fmt.Errorf("failed to restore merchant"))

	result, err := mockRepo.RestoreMerchant(merchantID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to restore merchant")
}

func TestDeleteMerchantPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	merchantID := 1

	mockRepo.EXPECT().DeleteMerchantPermanent(merchantID).Return(nil)

	err := mockRepo.DeleteMerchantPermanent(merchantID)

	assert.NoError(t, err)
}

func TestDeleteMerchantPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMerchantRepository(ctrl)

	merchantID := 1

	mockRepo.EXPECT().DeleteMerchantPermanent(merchantID).Return(fmt.Errorf("failed to delete merchant"))

	err := mockRepo.DeleteMerchantPermanent(merchantID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete merchant")
}
