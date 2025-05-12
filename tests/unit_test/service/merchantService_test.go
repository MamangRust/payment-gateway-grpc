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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestFindAllMerchants_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	merchants := []*record.MerchantRecord{
		{
			ID:   1,
			Name: "Merchant One",
		},
		{
			ID:   2,
			Name: "Merchant Two",
		},
	}

	expectedResponses := []*response.MerchantResponse{
		{
			ID:   1,
			Name: "Merchant One",
		},
		{
			ID:   2,
			Name: "Merchant Two",
		},
	}

	page := 1
	pageSize := 10
	search := ""
	totalRecords := 2

	mock_merchant_repo.EXPECT().FindAllMerchants(search, page, pageSize).Return(merchants, totalRecords, nil).Times(1)
	mock_mapping.EXPECT().ToMerchantsResponse(merchants).Return(expectedResponses).Times(1)

	result, total, errResp := merchantService.FindAll(page, pageSize, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponses, result)
	assert.Equal(t, totalRecords, total)
}

func TestFindAllMerchants_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, nil)

	page := 1
	pageSize := 10
	search := ""

	mock_merchant_repo.EXPECT().FindAllMerchants(search, page, pageSize).Return(nil, 0, fmt.Errorf("fetch merchant error")).Times(1)
	mock_logger.EXPECT().Error("Failed to fetch merchant records", gomock.Any()).Times(1)

	result, total, errResp := merchantService.FindAll(page, pageSize, search)

	assert.Nil(t, result)
	assert.Equal(t, 0, total)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to fetch merchant records", errResp.Message)
}

func TestFindByIdMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	merchant := &record.MerchantRecord{
		ID:   1,
		Name: "Merchant One",
	}

	expectedResponse := &response.MerchantResponse{
		ID:   1,
		Name: "Merchant One",
	}

	mock_logger.EXPECT().Debug("Finding merchant by ID", zap.Int("merchant_id", 1)).Times(1)
	mock_merchant_repo.EXPECT().FindById(1).Return(merchant, nil).Times(1)
	mock_mapping.EXPECT().ToMerchantResponse(merchant).Return(expectedResponse).Times(1)
	mock_logger.EXPECT().Debug("Successfully found merchant by ID", zap.Int("merchant_id", 1)).Times(1)

	result, errResp := merchantService.FindById(1)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestFindByIdMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, nil)

	mock_logger.EXPECT().Debug("Finding merchant by ID", zap.Int("merchant_id", 1)).Times(1)
	mock_merchant_repo.EXPECT().FindById(1).Return(nil, fmt.Errorf("fetch merchant error")).Times(1)

	mock_logger.EXPECT().Error("Failed to find merchant by ID", gomock.Any(), zap.Int("merchant_id", 1)).Times(1)

	result, errResp := merchantService.FindById(1)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Merchant not found", errResp.Message)
}

func TestFindByActiveMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	merchants := []*record.MerchantRecord{
		{
			ID:   1,
			Name: "Active Merchant One",
		},
		{
			ID:   2,
			Name: "Active Merchant Two",
		},
	}

	expectedResponses := []*response.MerchantResponseDeleteAt{
		{
			ID:   1,
			Name: "Active Merchant One",
		},
		{
			ID:   2,
			Name: "Active Merchant Two",
		},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 2

	mock_logger.EXPECT().Info("Fetching active merchants")
	mock_merchant_repo.EXPECT().FindByActive(search, page, pageSize).Return(merchants, expected, nil).Times(1)
	mock_mapping.EXPECT().ToMerchantsResponseDeleteAt(merchants).Return(expectedResponses).Times(1)

	mock_logger.EXPECT().Info("Successfully fetched active merchants")

	result, totalRecord, errResp := merchantService.FindByActive(pageSize, page, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponses, result)
	assert.Equal(t, expected, totalRecord)
}

func TestFindByActiveMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, nil)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	mock_logger.EXPECT().Info("Fetching active merchants")
	mock_logger.EXPECT().Error("Failed to fetch active merchants", zap.Error(fmt.Errorf("no active merchants"))).Times(1)
	mock_merchant_repo.EXPECT().FindByActive(
		search, page, pageSize,
	).Return(nil, expected, fmt.Errorf("no active merchants")).Times(1)

	result, totalRecord, errResp := merchantService.FindByActive(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "Failed to fetch active merchants", errResp.Message)
}

func TestFindByTrashedMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	merchants := []*record.MerchantRecord{
		{
			ID:   1,
			Name: "Trashed Merchant One",
		},
		{
			ID:   2,
			Name: "Trashed Merchant Two",
		},
	}

	expectedResponses := []*response.MerchantResponseDeleteAt{
		{
			ID:   1,
			Name: "Trashed Merchant One",
		},
		{
			ID:   2,
			Name: "Trashed Merchant Two",
		},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 2

	mock_logger.EXPECT().Info("Fetching trashed merchants")
	mock_merchant_repo.EXPECT().FindByTrashed(search, page, pageSize).Return(merchants, expected, nil).Times(1)
	mock_mapping.EXPECT().ToMerchantsResponseDeleteAt(merchants).Return(expectedResponses).Times(1)

	mock_logger.EXPECT().Info("Successfully fetched trashed merchants")

	result, totalRecord, errResp := merchantService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponses, result)
	assert.Equal(t, expected, totalRecord)
}

func TestFindByTrashedMerchant_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	mock_logger.EXPECT().Info("Fetching trashed merchants")

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	mock_logger.EXPECT().Error("Failed to fetch trashed merchants", zap.Error(fmt.Errorf("failed merchant")))

	mock_merchant_repo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, expected, fmt.Errorf("failed merchant")).Times(1)

	result, totalRecord, errResp := merchantService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to fetch trashed merchants", errResp.Message)
}

func TestFindByApiKeyMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)
	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	api_key := "apikey1"
	merchants := &record.MerchantRecord{
		ID:        1,
		Name:      "Merchant 1",
		ApiKey:    api_key,
		UserID:    101,
		Status:    "active",
		CreatedAt: "2024-01-01T10:00:00Z",
		UpdatedAt: "2024-01-01T10:00:00Z",
	}

	expectedResponse := &response.MerchantResponse{
		ID:        1,
		Name:      "Merchant 1",
		ApiKey:    api_key,
		Status:    "active",
		CreatedAt: "2024-01-01T10:00:00Z",
		UpdatedAt: "2024-01-01T10:00:00Z",
	}

	mock_logger.EXPECT().Debug("Finding merchant by API key", zap.String("api_key", api_key)).Times(1)
	mock_merchant_repo.EXPECT().FindByApiKey(api_key).Return(merchants, nil).Times(1)
	mock_mapping.EXPECT().ToMerchantResponse(merchants).Return(expectedResponse).Times(1)
	mock_logger.EXPECT().Debug("Successfully found merchant by API key", zap.String("api_key", api_key)).Times(1)

	result, errResp := merchantService.FindByApiKey(api_key)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestFindByApiKeyMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)
	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	api_key := "apikey1"
	mock_logger.EXPECT().Debug("Finding merchant by API key", zap.String("api_key", api_key)).Times(1)
	mock_merchant_repo.EXPECT().FindByApiKey(api_key).Return(nil, fmt.Errorf("failed merchant")).Times(1)
	mock_logger.EXPECT().Error("Failed to find merchant by API key", zap.Error(fmt.Errorf("failed merchant"))).Times(1)

	result, errResp := merchantService.FindByApiKey(api_key)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Merchant not found by API key", errResp.Message)
}

func TestFindByMerchantUserId_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	user_id := 101
	merchants := []*record.MerchantRecord{
		{
			ID:        1,
			Name:      "Merchant 1",
			ApiKey:    "apikey1",
			UserID:    user_id,
			Status:    "active",
			CreatedAt: "2024-01-01T10:00:00Z",
			UpdatedAt: "2024-01-01T10:00:00Z",
		},
		{
			ID:        2,
			Name:      "Merchant 2",
			ApiKey:    "apikey2",
			UserID:    user_id,
			Status:    "inactive",
			CreatedAt: "2024-01-02T10:00:00Z",
			UpdatedAt: "2024-01-02T10:00:00Z",
		},
	}

	expectedResponses := []*response.MerchantResponse{
		{
			ID:        1,
			Name:      "Merchant 1",
			ApiKey:    "apikey1",
			Status:    "active",
			CreatedAt: "2024-01-01T10:00:00Z",
			UpdatedAt: "2024-01-01T10:00:00Z",
		},
		{
			ID:        2,
			Name:      "Merchant 2",
			ApiKey:    "apikey2",
			Status:    "inactive",
			CreatedAt: "2024-01-02T10:00:00Z",
			UpdatedAt: "2024-01-02T10:00:00Z",
		},
	}

	mock_logger.EXPECT().Debug("Finding merchant by user ID", zap.Int("user_id", user_id)).Times(1)
	mock_merchant_repo.EXPECT().FindByMerchantUserId(user_id).Return(merchants, nil).Times(1)
	mock_mapping.EXPECT().ToMerchantsResponse(merchants).Return(expectedResponses).Times(1)
	mock_logger.EXPECT().Debug("Successfully found merchant by user ID", zap.Int("user_id", user_id)).Times(1)

	result, errResp := merchantService.FindByMerchantUserId(user_id)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponses, result)
}

func TestFindByMerchantUserId_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)
	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	user_id := 101
	mock_logger.EXPECT().Debug("Finding merchant by user ID", zap.Int("user_id", user_id)).Times(1)
	mock_merchant_repo.EXPECT().FindByMerchantUserId(user_id).Return(nil, fmt.Errorf("failed merchant")).Times(1)
	mock_logger.EXPECT().Error("Failed to find merchant by user ID", zap.Error(fmt.Errorf("failed merchant"))).Times(1)

	result, errResp := merchantService.FindByMerchantUserId(user_id)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Merchant not found by user ID", errResp.Message)
}

func TestCreateMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	userId := 1

	request := requests.CreateMerchantRequest{
		Name:   "Merchant 1",
		UserID: userId,
	}

	merchantRecord := &record.MerchantRecord{
		ID:        1,
		Name:      "Merchant 1",
		ApiKey:    "apikey1",
		UserID:    userId,
		Status:    "active",
		CreatedAt: "2024-01-01T10:00:00Z",
		UpdatedAt: "2024-01-01T10:00:00Z",
	}

	expectedResponse := &response.MerchantResponse{
		ID:        1,
		Name:      "Merchant 1",
		ApiKey:    "apikey1",
		Status:    "active",
		CreatedAt: "2024-01-01T10:00:00Z",
		UpdatedAt: "2024-01-01T10:00:00Z",
	}

	mock_logger.EXPECT().Debug("Creating new merchant", zap.String("merchant_name", request.Name)).Times(1)

	mock_merchant_repo.EXPECT().
		CreateMerchant(gomock.Eq(&request)).
		Return(merchantRecord, nil).
		Times(1)

	mock_mapping.EXPECT().
		ToMerchantResponse(gomock.Eq(merchantRecord)).
		Return(expectedResponse).
		Times(1)

	mock_logger.EXPECT().Debug("Successfully created merchant", zap.Int("merchant_id", merchantRecord.ID)).Times(1)

	result, errResp := merchantService.CreateMerchant(&request)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestCreateMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	request := requests.CreateMerchantRequest{
		Name:   "Merchant 1",
		UserID: 1,
	}

	merchantRecord := &record.MerchantRecord{
		ID:        1,
		Name:      "Merchant 1",
		ApiKey:    "apikey1",
		UserID:    1,
		Status:    "active",
		CreatedAt: "2024-01-01T10:00:00Z",
		UpdatedAt: "2024-01-01T10:00:00Z",
	}

	mock_logger.EXPECT().Debug("Creating new merchant", zap.String("merchant_name", request.Name)).Times(1)

	mock_merchant_repo.EXPECT().
		CreateMerchant(gomock.Eq(&request)).
		Return(merchantRecord, fmt.Errorf("failed to create merchant")).
		Times(1)

	mock_logger.EXPECT().Error("Failed to create merchant", zap.Error(fmt.Errorf("failed to create merchant"))).Times(1)

	result, errResp := merchantService.CreateMerchant(&request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to create merchant", errResp.Message)
}

func TestUpdateMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	requests := requests.UpdateMerchantRequest{
		MerchantID: 1,
		Name:       "Merchant 1",
		UserID:     1,
		Status:     "active",
	}

	merchantRecord := &record.MerchantRecord{
		ID:        1,
		Name:      "Merchant 1",
		ApiKey:    "apikey1",
		UserID:    1,
		Status:    "active",
		CreatedAt: "2024-01-01T10:00:00Z",
		UpdatedAt: "2024-01-01T10:00:00Z",
	}

	expectedResponse := &response.MerchantResponse{
		ID:     1,
		Name:   "Merchant 1",
		ApiKey: "apikey1",
		Status: "active",
	}

	mock_logger.EXPECT().Debug("Updating merchant", zap.Int("merchant_id", merchantRecord.ID)).Times(1)

	mock_merchant_repo.EXPECT().FindById(merchantRecord.ID).Return(merchantRecord, nil).Times(1)

	mock_merchant_repo.EXPECT().
		UpdateMerchant(&requests).
		Return(merchantRecord, nil).
		Times(1)

	mock_mapping.EXPECT().
		ToMerchantResponse(gomock.Eq(merchantRecord)).
		Return(expectedResponse).
		Times(1)

	mock_logger.EXPECT().Debug("Successfully updated merchant", zap.Int("merchant_id", merchantRecord.ID)).Times(1)

	result, errResp := merchantService.UpdateMerchant(&requests)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestUpdateMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	requests := requests.UpdateMerchantRequest{
		MerchantID: 1,
		Name:       "Merchant 1",
		UserID:     1,
		Status:     "active",
	}

	mock_logger.EXPECT().Debug("Updating merchant", zap.Int("merchant_id", requests.MerchantID)).Times(1)

	mock_merchant_repo.EXPECT().FindById(requests.MerchantID).Return(nil, fmt.Errorf("merchant not found")).Times(1)

	mock_logger.EXPECT().Error("Merchant not found for update", zap.Error(fmt.Errorf("merchant not found"))).Times(1)

	result, errResp := merchantService.UpdateMerchant(&requests)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Merchant not found", errResp.Message)
}

func TestTrashedMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	merchantId := 1

	merchantRecord := &record.MerchantRecord{
		ID:        1,
		Name:      "Merchant 1",
		ApiKey:    "apikey1",
		UserID:    1,
		Status:    "active",
		CreatedAt: "2024-01-01T10:00:00Z",
		UpdatedAt: "2024-01-01T10:00:00Z",
	}

	expectedResponse := &response.MerchantResponse{
		ID:     1,
		Name:   "Merchant 1",
		ApiKey: "apikey1",
		Status: "active",
	}

	mock_logger.EXPECT().Debug("Trashing merchant", zap.Int("merchant_id", merchantRecord.ID)).Times(1)

	mock_merchant_repo.EXPECT().TrashedMerchant(merchantId).Return(merchantRecord, nil).Times(1)

	mock_mapping.EXPECT().
		ToMerchantResponse(gomock.Eq(merchantRecord)).
		Return(expectedResponse).
		Times(1)

	mock_logger.EXPECT().Debug("Successfully trashed merchant", zap.Int("merchant_id", merchantId)).Times(1)

	result, errResp := merchantService.TrashedMerchant(merchantId)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestTrashedMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	merchantId := 1
	expectedError := errors.New("Failed to trash merchant")

	mock_logger.EXPECT().Debug("Trashing merchant", zap.Int("merchant_id", merchantId)).Times(1)

	mock_logger.EXPECT().
		Error(
			"Failed to trash merchant",
			gomock.AssignableToTypeOf(zap.Error(expectedError)),
			zap.Int("merchant_id", merchantId),
		).Times(1)

	mock_merchant_repo.EXPECT().TrashedMerchant(merchantId).Return(nil, fmt.Errorf("merchant not found")).Times(1)

	result, errResp := merchantService.TrashedMerchant(merchantId)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to trash merchant", errResp.Message)
}

func TestRestoreMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	merchantId := 1

	merchantRecord := &record.MerchantRecord{
		ID:        1,
		Name:      "Merchant 1",
		ApiKey:    "apikey1",
		UserID:    1,
		Status:    "restored",
		CreatedAt: "2024-01-01T10:00:00Z",
		UpdatedAt: "2024-01-02T10:00:00Z",
	}

	expectedResponse := &response.MerchantResponse{
		ID:     1,
		Name:   "Merchant 1",
		ApiKey: "apikey1",
		Status: "restored",
	}

	mock_logger.EXPECT().Debug("Restoring merchant", zap.Int("merchant_id", merchantId)).Times(1)

	mock_merchant_repo.EXPECT().
		RestoreMerchant(merchantId).
		Return(merchantRecord, nil).
		Times(1)

	mock_mapping.EXPECT().
		ToMerchantResponse(gomock.Eq(merchantRecord)).
		Return(expectedResponse).
		Times(1)

	mock_logger.EXPECT().Debug("Successfully restored merchant", zap.Int("merchant_id", merchantId)).Times(1)

	result, errResp := merchantService.RestoreMerchant(merchantId)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestRestoreMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)
	mock_mapping := mock_responsemapper.NewMockMerchantResponseMapper(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, mock_mapping)

	merchantId := 1
	expectedError := errors.New("Failed to restore merchant")

	mock_logger.EXPECT().Debug("Restoring merchant", zap.Int("merchant_id", merchantId)).Times(1)

	mock_logger.EXPECT().
		Error(
			"Failed to restore merchant",
			gomock.AssignableToTypeOf(zap.Error(expectedError)),
			zap.Int("merchant_id", merchantId),
		).Times(1)

	mock_merchant_repo.EXPECT().
		RestoreMerchant(merchantId).
		Return(nil, fmt.Errorf("merchant not found")).
		Times(1)

	result, errResp := merchantService.RestoreMerchant(merchantId)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to restore merchant", errResp.Message)
}

func TestDeleteMerchantPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, nil)

	merchantId := 1

	mock_logger.EXPECT().Debug("Deleting merchant permanently", zap.Int("merchant_id", merchantId)).Times(1)

	mock_merchant_repo.EXPECT().
		DeleteMerchantPermanent(merchantId).
		Return(nil).
		Times(1)

	mock_logger.EXPECT().Debug("Successfully deleted merchant permanently", zap.Int("merchant_id", merchantId)).Times(1)

	result, errResp := merchantService.DeleteMerchantPermanent(merchantId)

	assert.Nil(t, result)
	assert.Nil(t, errResp)
}

func TestDeleteMerchantPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_merchant_repo := mock_repository.NewMockMerchantRepository(ctrl)
	mock_logger := mock_logger.NewMockLoggerInterface(ctrl)

	merchantService := service.NewMerchantService(mock_merchant_repo, mock_logger, nil)

	merchantId := 1
	expectedError := errors.New("Failed to delete merchant permanently")

	mock_logger.EXPECT().Debug("Deleting merchant permanently", zap.Int("merchant_id", merchantId)).Times(1)

	mock_merchant_repo.EXPECT().
		DeleteMerchantPermanent(merchantId).
		Return(fmt.Errorf("merchant not found")).
		Times(1)

	mock_logger.EXPECT().
		Error(
			"Failed to delete merchant permanently",
			gomock.AssignableToTypeOf(zap.Error(expectedError)),
			zap.Int("merchant_id", merchantId),
		).Times(1)

	result, errResp := merchantService.DeleteMerchantPermanent(merchantId)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to delete merchant permanently", errResp.Message)
}
