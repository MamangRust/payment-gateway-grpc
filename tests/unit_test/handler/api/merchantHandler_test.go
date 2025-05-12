package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/handler/api"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	mock_pb "MamangRust/paymentgatewaygrpc/internal/pb/mocks"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFindAllMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	page, pageSize := 1, 10
	search := "merchant"

	expectedResponse := &pb.ApiResponsePaginationMerchant{
		Status:  "success",
		Message: "Merchants retrieved successfully",
		Data: []*pb.MerchantResponse{
			{Id: 1, Name: "Merchant 1"},
			{Id: 2, Name: "Merchant 2"},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant?page=%d&page_size=%d&search=%s", page, pageSize, search), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockMerchantClient.EXPECT().
		FindAllMerchant(gomock.Any(), &pb.FindAllMerchantRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Search:   search,
		}).
		Return(expectedResponse, nil).
		Times(1)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindAll(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Merchants retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, "Merchant 1", resp.Data[0].Name)
	assert.Equal(t, "Merchant 2", resp.Data[1].Name)
}

func TestFindAllMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	page, pageSize := 1, 10
	search := "merchant"

	mockLogger.EXPECT().Debug("Failed to retrieve merchant data", gomock.Any()).Times(1)

	mockMerchantClient.EXPECT().
		FindAllMerchant(gomock.Any(), &pb.FindAllMerchantRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Search:   search,
		}).
		Return(nil, fmt.Errorf("internal server error")).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant?page=%d&page_size=%d&search=%s", page, pageSize, search), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindAll(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve merchant data: ", resp.Message)
}

func TestFindAll_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	page, pageSize := 1, 10
	search := "nonexistent"

	expectedResponse := &pb.ApiResponsePaginationMerchant{
		Status:  "success",
		Message: "No merchants found",
		Data:    []*pb.MerchantResponse{},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant?page=%d&page_size=%d&search=%s", page, pageSize, search), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockMerchantClient.EXPECT().
		FindAllMerchant(gomock.Any(), &pb.FindAllMerchantRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Search:   search,
		}).
		Return(expectedResponse, nil).
		Times(1)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindAll(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No merchants found", resp.Message)
	assert.Len(t, resp.Data, 0)
}

func TestFindByIdMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	id := 1
	expectedResponse := &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Merchant retrieved successfully",
		Data: &pb.MerchantResponse{
			Id:   1,
			Name: "Merchant 1",
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant/%d", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	mockMerchantClient.EXPECT().
		FindByIdMerchant(gomock.Any(), &pb.FindByIdMerchantRequest{
			MerchantId: int32(id),
		}).
		Return(expectedResponse, nil).
		Times(1)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Merchant retrieved successfully", resp.Message)
	assert.Equal(t, int32(1), resp.Data.Id)
	assert.Equal(t, "Merchant 1", resp.Data.Name)
}

func TestFindByIdMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	id := 1

	mockLogger.EXPECT().Debug("Failed to retrieve merchant data", gomock.Any()).Times(1)

	mockMerchantClient.EXPECT().
		FindByIdMerchant(gomock.Any(), &pb.FindByIdMerchantRequest{
			MerchantId: int32(id),
		}).
		Return(nil, fmt.Errorf("internal server error")).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant/%d", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve merchant data: ", resp.Message)
}

func TestFindByIdMerchant_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	invalidID := "abc"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant/%s", invalidID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(invalidID)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid merchant ID", resp.Message)
}

func TestFindByApiKey_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	apiKey := "valid-api-key"
	expectedResponse := &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Merchant retrieved successfully",
		Data: &pb.MerchantResponse{
			Id:   1,
			Name: "Merchant 1",
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant/find?api_key=%s", apiKey), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockMerchantClient.EXPECT().
		FindByApiKey(gomock.Any(), &pb.FindByApiKeyRequest{
			ApiKey: apiKey,
		}).
		Return(expectedResponse, nil).
		Times(1)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindByApiKey(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Merchant retrieved successfully", resp.Message)
	assert.Equal(t, int32(1), resp.Data.Id)
	assert.Equal(t, "Merchant 1", resp.Data.Name)
}

func TestFindByApiKey_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	apiKey := "valid-api-key"

	mockLogger.EXPECT().Debug("Failed to retrieve merchant data", gomock.Any()).Times(1)

	mockMerchantClient.EXPECT().
		FindByApiKey(gomock.Any(), &pb.FindByApiKeyRequest{
			ApiKey: apiKey,
		}).
		Return(nil, fmt.Errorf("internal server error")).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant/find?api_key=%s", apiKey), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindByApiKey(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve merchant data: ", resp.Message)
}

func TestFindByMerchantUserId_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userID := int32(1)

	expectedResponse := []*pb.MerchantResponse{
		{
			Id:   1,
			Name: "Merchant 1",
		},
		{
			Id:   2,
			Name: "Merchant 2",
		},
	}
	mockResponse := &pb.ApiResponsesMerchant{
		Status:  "success",
		Message: "Merchant retrieved successfully",
		Data:    expectedResponse,
	}

	mockMerchantClient.EXPECT().
		FindByMerchantUserId(
			gomock.Any(),
			&pb.FindByMerchantUserIdRequest{UserId: userID},
		).
		Return(mockResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/merchant/merchant-user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)
	err := handler.FindByMerchantUserId(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Merchant retrieved successfully", resp.Message)
	assert.Equal(t, expectedResponse[0].Id, resp.Data[0].Id)
	assert.Equal(t, expectedResponse[0].Name, resp.Data[0].Name)
}

func TestFindByMerchantUserId_InvalidUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	invalidUserID := "not_a_number"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/merchant/merchant-user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", invalidUserID)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)
	err := handler.FindByMerchantUserId(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Invalid merchant ID")
}

func TestFindByMerchantUserId_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userID := int32(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/merchant/merchant-user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)
	err := handler.FindByMerchantUserId(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid merchant ID", resp.Message)

	mockMerchantClient.EXPECT().
		FindByMerchantUserId(
			gomock.Any(),
			&pb.FindByMerchantUserIdRequest{UserId: userID},
		).
		Return(nil, fmt.Errorf("service error"))

	mockLogger.EXPECT().Debug("Failed to retrieve merchant data", gomock.Any()).Times(1)

	e = echo.New()
	req = httptest.NewRequest(http.MethodGet, "/api/merchant/merchant-user", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set("user_id", userID)

	handler = api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)
	err = handler.FindByMerchantUserId(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve merchant data: ", resp.Message)

	mockMerchantClient.EXPECT().
		FindByMerchantUserId(
			gomock.Any(),
			&pb.FindByMerchantUserIdRequest{UserId: userID},
		).
		Return(&pb.ApiResponsesMerchant{
			Status:  "error",
			Message: "Invalid merchant ID",
			Data:    nil,
		}, nil)

	e = echo.New()
	req = httptest.NewRequest(http.MethodGet, "/api/merchant/merchant-user", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set("user_id", userID)

	handler = api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)
	err = handler.FindByMerchantUserId(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var respEmpty pb.ApiResponsesMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &respEmpty)
	assert.NoError(t, err)
	assert.Equal(t, "error", respEmpty.Status)
	assert.Equal(t, "Invalid merchant ID", respEmpty.Message)
	assert.Nil(t, respEmpty.Data)
}

func TestFindByActiveMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := []*pb.MerchantResponseDeleteAt{
		{
			Id:   1,
			Name: "Active Merchant 1",
		},
		{
			Id:   2,
			Name: "Active Merchant 2",
		},
	}

	request := &pb.FindAllMerchantRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockResponse := &pb.ApiResponsePaginationMerchantDeleteAt{
		Status:  "success",
		Message: "Active merchants retrieved successfully",
		Data:    expectedResponse,
	}

	mockMerchantClient.EXPECT().
		FindByActive(gomock.Any(), request).
		Return(mockResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/merchant/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindByActive(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Active merchants retrieved successfully", resp.Message)
	assert.Equal(t, len(expectedResponse), len(resp.Data))
	assert.Equal(t, expectedResponse[0].Id, resp.Data[0].Id)
	assert.Equal(t, expectedResponse[0].Name, resp.Data[0].Name)
}

func TestFindByActiveMerchant_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockResponse := &pb.ApiResponsePaginationMerchantDeleteAt{
		Status:  "success",
		Message: "No active merchants found",
		Data:    []*pb.MerchantResponseDeleteAt{},
	}

	mockMerchantClient.EXPECT().
		FindByActive(gomock.Any(), gomock.Any()).
		Return(mockResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/merchant/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindByActive(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No active merchants found", resp.Message)
	assert.Equal(t, 0, len(resp.Data))
}

func TestFindByActiveMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockMerchantClient.EXPECT().
		FindByActive(gomock.Any(), gomock.Any()).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve merchant data: ",
		}).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve merchant data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/merchant/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindByActive(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve merchant data: ", resp.Message)
}

func TestFindByTrashedMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := []*pb.MerchantResponseDeleteAt{
		{
			Id:   1,
			Name: "Trashed Merchant 1",
		},
		{
			Id:   2,
			Name: "Trashed Merchant 2",
		},
	}
	request := &pb.FindAllMerchantRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockResponse := &pb.ApiResponsePaginationMerchantDeleteAt{
		Status:  "success",
		Message: "Trashed merchants retrieved successfully",
		Data:    expectedResponse,
	}

	mockMerchantClient.EXPECT().
		FindByTrashed(gomock.Any(), request).
		Return(mockResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/merchant/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindByTrashed(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Trashed merchants retrieved successfully", resp.Message)
	assert.Equal(t, len(expectedResponse), len(resp.Data))
	assert.Equal(t, expectedResponse[0].Id, resp.Data[0].Id)
	assert.Equal(t, expectedResponse[0].Name, resp.Data[0].Name)
}

func TestFindByTrashedMerchant_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockResponse := &pb.ApiResponsePaginationMerchantDeleteAt{
		Status:  "success",
		Message: "No trashed merchants found",
		Data:    []*pb.MerchantResponseDeleteAt{},
	}

	request := &pb.FindAllMerchantRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockMerchantClient.EXPECT().
		FindByTrashed(gomock.Any(), request).
		Return(mockResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/merchant/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindByTrashed(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No trashed merchants found", resp.Message)
	assert.Equal(t, 0, len(resp.Data))
}

func TestFindByTrashed_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockMerchantClient.EXPECT().
		FindByTrashed(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("service error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve merchant data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/merchant/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.FindByTrashed(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve merchant data: ", resp.Message)
}

func TestCreateMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.CreateMerchantRequest{
		Name:   "New Merchant",
		UserID: 1,
	}

	expectedResponse := &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Merchant created successfully",
		Data: &pb.MerchantResponse{
			Id:   1,
			Name: "New Merchant",
		},
	}

	mockMerchantClient.EXPECT().
		CreateMerchant(
			gomock.Any(),
			&pb.CreateMerchantRequest{
				Name:   body.Name,
				UserId: int32(body.UserID),
			},
		).
		Return(expectedResponse, nil).
		Times(1)

	e := echo.New()
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/merchant", bytes.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.Create(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Merchant created successfully", resp.Message)
	assert.Equal(t, expectedResponse.Data.Id, resp.Data.Id)
	assert.Equal(t, expectedResponse.Data.Name, resp.Data.Name)
}

func TestCreateMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.CreateMerchantRequest{
		Name:   "New Merchant",
		UserID: 1,
	}

	mockMerchantClient.EXPECT().
		CreateMerchant(
			gomock.Any(),
			&pb.CreateMerchantRequest{
				Name:   body.Name,
				UserId: int32(body.UserID),
			},
		).
		Return(nil, fmt.Errorf("service error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to create merchant", gomock.Any()).Times(1)

	e := echo.New()
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/merchant", bytes.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.Create(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to create merchant:")
}

func TestCreateMerchant_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.CreateMerchantRequest{
		Name:   "",
		UserID: 0,
	}

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).Times(1)

	e := echo.New()
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/merchant", bytes.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.Create(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error:")
	assert.Contains(t, resp.Message, "'Name' failed on the 'required' tag")
	assert.Contains(t, resp.Message, "'UserID' failed on the 'required' tag")
}

func TestUpdateMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.UpdateMerchantRequest{
		MerchantID: 1,
		Name:       "Updated Merchant",
		UserID:     1,
		Status:     "active",
	}

	expectedResponse := &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Merchant updated successfully",
		Data: &pb.MerchantResponse{
			Id:   1,
			Name: "Updated Merchant",
		},
	}

	mockMerchantClient.EXPECT().
		UpdateMerchant(
			gomock.Any(),
			&pb.UpdateMerchantRequest{
				MerchantId: int32(body.MerchantID),
				Name:       body.Name,
				UserId:     int32(body.UserID),
				Status:     body.Status,
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/merchant/update", bytes.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("id", int32(body.MerchantID))

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Merchant updated successfully", resp.Message)
	assert.Equal(t, body.MerchantID, int(resp.Data.Id))
	assert.Equal(t, body.Name, resp.Data.Name)
}

func TestUpdateMerchant_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.UpdateMerchantRequest{
		MerchantID: 0,
		Name:       "Updated Merchant",
		UserID:     1,
		Status:     "active",
	}

	merchantId := "invalid"

	e := echo.New()
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/merchant/update/%s", merchantId), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(merchantId)

	mockLogger.EXPECT().Debug("Invalid merchant ID", gomock.Any()).Times(1)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid merchant ID", resp.Message)

}

func TestUpdateMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.UpdateMerchantRequest{
		MerchantID: 1,
		Name:       "Updated Merchant",
		UserID:     1,
		Status:     "active",
	}

	mockLogger.EXPECT().Debug("Failed to update merchant", gomock.Any()).Times(1)

	mockMerchantClient.EXPECT().
		UpdateMerchant(
			gomock.Any(),
			&pb.UpdateMerchantRequest{
				MerchantId: int32(body.MerchantID),
				Name:       body.Name,
				UserId:     int32(body.UserID),
				Status:     body.Status,
			},
		).
		Return(nil, status.Error(codes.Internal, "internal server error"))

	e := echo.New()
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/merchant/update", bytes.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("id", int32(body.MerchantID))

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to update merchant: ", resp.Message)
}

func TestUpdateMerchant_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.UpdateMerchantRequest{
		MerchantID: 1,
		Name:       "",
		UserID:     1,
		Status:     "active",
	}

	mockLogger.EXPECT().Debug("Validation Error", gomock.Any()).Times(1)

	e := echo.New()
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/merchant/update", bytes.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("id", int32(body.MerchantID))

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Validation Error: ", resp.Message)
}

func TestTrashedMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	merchantID := 1

	expectedResponse := &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Merchant trashed successfully",
		Data: &pb.MerchantResponse{
			Id:   int32(merchantID),
			Name: "Merchant 1",
		},
	}

	mockMerchantClient.EXPECT().
		TrashedMerchant(
			gomock.Any(),
			&pb.FindByIdMerchantRequest{MerchantId: int32(merchantID)},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/merchant/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.TrashedMerchant(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Merchant trashed successfully", resp.Message)
	assert.Equal(t, int32(merchantID), resp.Data.Id)
	assert.Equal(t, "Merchant 1", resp.Data.Name)
}

func TestTrashedMerchant_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/merchant/trashed/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	mockLogger.EXPECT().Debug("Bad Request", gomock.Any()).Times(1)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.TrashedMerchant(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Bad Request: Invalid ID")
}

func TestTrashedMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	merchantID := 1

	mockMerchantClient.EXPECT().
		TrashedMerchant(
			gomock.Any(),
			&pb.FindByIdMerchantRequest{MerchantId: int32(merchantID)},
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to trashed merchant",
		})

	mockLogger.EXPECT().Debug("Failed to trashed merchant", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/merchant/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.TrashedMerchant(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to trashed merchant")
}

func TestRestoreMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	merchantID := 1

	expectedResponse := &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Merchant restored successfully",
		Data: &pb.MerchantResponse{
			Id:     int32(merchantID),
			Name:   "Merchant 1",
			Status: "active",
		},
	}

	mockMerchantClient.EXPECT().
		RestoreMerchant(
			gomock.Any(),
			&pb.FindByIdMerchantRequest{MerchantId: int32(merchantID)},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/merchant/restore/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.RestoreMerchant(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Merchant restored successfully", resp.Message)
	assert.Equal(t, int32(merchantID), resp.Data.Id)
	assert.Equal(t, "active", resp.Data.Status)
}

func TestRestoreMerchant_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/merchant/restore/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	mockLogger.EXPECT().Debug("Bad Request", gomock.Any()).Times(1)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.RestoreMerchant(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Bad Request: Invalid ID")
}

func TestRestoreMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	merchantID := 1

	mockMerchantClient.EXPECT().
		RestoreMerchant(
			gomock.Any(),
			&pb.FindByIdMerchantRequest{MerchantId: int32(merchantID)},
		).
		Return(nil, fmt.Errorf("internal server error"))

	mockLogger.EXPECT().Debug("Failed to restore merchant", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/merchant/restore/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.RestoreMerchant(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to restore merchant:", resp.Message)
}

func TestDeleteMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	merchantID := 1

	expectedResponse := &pb.ApiResponseMerchatDelete{
		Status:  "success",
		Message: "Merchant deleted successfully",
	}

	mockMerchantClient.EXPECT().
		DeleteMerchantPermanent(
			gomock.Any(),
			&pb.FindByIdMerchantRequest{MerchantId: int32(merchantID)},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/merchant/delete/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.Delete(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Merchant deleted successfully", resp.Message)
}

func TestDeleteMerchant_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/merchant/delete/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	// mockLogger.EXPECT().
	// 	Debug("Bad Request: Invalid ID", gomock.Any()).
	// 	Times(1)

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.Delete(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Bad Request: Invalid ID")
}

func TestDeleteMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	merchantID := 1

	mockMerchantClient.EXPECT().
		DeleteMerchantPermanent(
			gomock.Any(),
			&pb.FindByIdMerchantRequest{MerchantId: int32(merchantID)},
		).
		Return(nil, fmt.Errorf("internal server error"))

	mockLogger.EXPECT().Debug("Failed to delete merchant", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/merchant/delete/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger)

	err := handler.Delete(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to delete merchant:", resp.Message)
}
