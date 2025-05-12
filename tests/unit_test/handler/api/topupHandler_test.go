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
)

func TestFindAllTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationTopup{
		Status:  "success",
		Message: "Topup data retrieved successfully",
		Data: []*pb.TopupResponse{
			{
				Id:          1,
				CardNumber:  "1234567890",
				TopupAmount: 10000,
			},
			{
				Id:          2,
				CardNumber:  "0987654321",
				TopupAmount: 20000,
			},
		},
		Pagination: &pb.PaginationMeta{
			CurrentPage: 1,
			PageSize:    2,
			TotalPages:  1,
		},
	}

	mockTopupClient.EXPECT().
		FindAllTopup(
			gomock.Any(),
			&pb.FindAllTopupRequest{
				Page:     1,
				PageSize: 10,
				Search:   "",
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findall?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Topup data retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 2)
}

func TestFindAllTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockTopupClient.EXPECT().
		FindAllTopup(
			gomock.Any(),
			&pb.FindAllTopupRequest{
				Page:     1,
				PageSize: 10,
				Search:   "",
			},
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve topup data",
		})

	mockLogger.EXPECT().Debug("Failed to retrieve topup data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findall?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to retrieve topup data")
}

func TestFindAllTopup_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationTopup{
		Status:     "success",
		Message:    "No topup data found",
		Data:       []*pb.TopupResponse{},
		Pagination: &pb.PaginationMeta{},
	}

	mockTopupClient.EXPECT().
		FindAllTopup(
			gomock.Any(),
			&pb.FindAllTopupRequest{
				Page:     1,
				PageSize: 10,
				Search:   "",
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findall?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No topup data found", resp.Message)
	assert.Empty(t, resp.Data)
}

func TestFindByIdTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponseTopup{
		Status:  "success",
		Message: "Topup data retrieved successfully",
		Data: &pb.TopupResponse{
			Id:          1,
			CardNumber:  "1234567890",
			TopupAmount: 10000,
		},
	}

	mockTopupClient.EXPECT().
		FindByIdTopup(
			gomock.Any(),
			&pb.FindByIdTopupRequest{
				TopupId: 1,
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findbyid/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Topup data retrieved successfully", resp.Message)
	assert.Equal(t, int32(1), resp.Data.Id)
	assert.Equal(t, "1234567890", resp.Data.CardNumber)
	assert.Equal(t, int32(10000), resp.Data.TopupAmount)
}

func TestFindByIdTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockTopupClient.EXPECT().
		FindByIdTopup(
			gomock.Any(),
			&pb.FindByIdTopupRequest{
				TopupId: 1,
			},
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve topup data",
		})

	mockLogger.EXPECT().Debug("Failed to retrieve topup data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findbyid/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to retrieve topup data")
}

func TestFindByIdTopup_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findbyid/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Bad Request: Invalid ID", resp.Message)
}

func TestFindByCardNumberTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponseTopup{
		Status:  "success",
		Message: "Topup data retrieved successfully",
		Data: &pb.TopupResponse{
			Id:          1,
			CardNumber:  "1234567890",
			TopupAmount: 10000,
		},
	}

	mockTopupClient.EXPECT().
		FindByCardNumberTopup(
			gomock.Any(),
			&pb.FindByCardNumberTopupRequest{
				CardNumber: "1234567890",
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findbycardnumber/1234567890", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues("1234567890")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindByCardNumber(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Topup data retrieved successfully", resp.Message)
	assert.Equal(t, int32(1), resp.Data.Id)
	assert.Equal(t, "1234567890", resp.Data.CardNumber)
	assert.Equal(t, int32(10000), resp.Data.TopupAmount)
}

func TestFindByCardNumberTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockTopupClient.EXPECT().
		FindByCardNumberTopup(
			gomock.Any(),
			&pb.FindByCardNumberTopupRequest{
				CardNumber: "1234567890",
			},
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve topup data",
		})

	mockLogger.EXPECT().Debug("Failed to retrieve topup data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findbycardnumber/1234567890", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues("1234567890")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindByCardNumber(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to retrieve topup data")
}

func TestFindByActiveTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationTopupDeleteAt{
		Status:  "success",
		Message: "Topup data retrieved successfully",
		Data: []*pb.TopupResponseDeleteAt{
			{
				Id:          1,
				CardNumber:  "1234567890",
				TopupAmount: 10000,
			},
		},
	}

	request := &pb.FindAllTopupRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockTopupClient.EXPECT().
		FindByActive(
			gomock.Any(),
			request,
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findbyactive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindByActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationTopupDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Topup data retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 1)
}

func TestFindByActiveTopup_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllTopupRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	expectedResponse := &pb.ApiResponsePaginationTopupDeleteAt{
		Status:  "success",
		Message: "No active topup data found",
		Data:    []*pb.TopupResponseDeleteAt{},
	}

	mockTopupClient.EXPECT().
		FindByActive(
			gomock.Any(),
			request,
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findbyactive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindByActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No active topup data found", resp.Message)
	assert.Len(t, resp.Data, 0)
}

func TestFindByActiveTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllTopupRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockTopupClient.EXPECT().
		FindByActive(
			gomock.Any(),
			request,
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve topup data",
		})

	mockLogger.EXPECT().Debug("Failed to retrieve topup data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findbyactive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindByActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to retrieve topup data")
}

func TestFindByTrashedTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllTopupRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	expectedResponse := &pb.ApiResponsePaginationTopupDeleteAt{
		Status:  "success",
		Message: "Topup data retrieved successfully",
		Data: []*pb.TopupResponseDeleteAt{
			{
				Id:          1,
				CardNumber:  "1234567890",
				TopupAmount: 10000,
			},
		},
	}

	mockTopupClient.EXPECT().
		FindByTrashed(
			gomock.Any(),
			request,
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findbytrashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindByTrashed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Topup data retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 1)
}

func TestFindByTrashedTopup_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationTopupDeleteAt{
		Status:  "success",
		Message: "No trashed topup data found",
		Data:    []*pb.TopupResponseDeleteAt{},
	}

	request := &pb.FindAllTopupRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockTopupClient.EXPECT().
		FindByTrashed(
			gomock.Any(),
			request,
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findbytrashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindByTrashed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No trashed topup data found", resp.Message)
	assert.Len(t, resp.Data, 0)
}

func TestFindByTrashedTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllTopupRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockTopupClient.EXPECT().
		FindByTrashed(
			gomock.Any(),
			request,
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve topup data",
		})

	mockLogger.EXPECT().Debug("Failed to retrieve topup data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/topup/findbytrashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.FindByTrashed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to retrieve topup data")
}

func TestCreateTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponseTopup{
		Status:  "success",
		Message: "Topup created successfully",
		Data: &pb.TopupResponse{
			Id:          1,
			CardNumber:  "1234567890",
			TopupAmount: 500000,
		},
	}

	req := &requests.CreateTopupRequest{
		CardNumber:  "1234567890",
		TopupNo:     "TOPUP123",
		TopupAmount: 500000,
		TopupMethod: "mandiri",
	}

	mockTopupClient.EXPECT().
		CreateTopup(
			gomock.Any(),
			&pb.CreateTopupRequest{
				CardNumber:  req.CardNumber,
				TopupNo:     req.TopupNo,
				TopupAmount: int32(req.TopupAmount),
				TopupMethod: req.TopupMethod,
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	bodyJson, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	httpReq := httptest.NewRequest(http.MethodPost, "/api/topup/create", bytes.NewReader(bodyJson))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err = handler.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Topup created successfully", resp.Message)
	assert.Equal(t, "1234567890", resp.Data.CardNumber)
	assert.Equal(t, int32(500000), resp.Data.TopupAmount)
}

func TestCreateTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	req := &requests.CreateTopupRequest{
		CardNumber:  "1234567890",
		TopupNo:     "TOPUP123",
		TopupAmount: 500000,
		TopupMethod: "mandiri",
	}

	mockTopupClient.EXPECT().
		CreateTopup(
			gomock.Any(),
			&pb.CreateTopupRequest{
				CardNumber:  req.CardNumber,
				TopupNo:     req.TopupNo,
				TopupAmount: int32(req.TopupAmount),
				TopupMethod: req.TopupMethod,
			},
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create topup",
		})

	mockLogger.EXPECT().Debug("Failed to create topup", gomock.Any()).Times(1)

	e := echo.New()
	bodyJson, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	httpReq := httptest.NewRequest(http.MethodPost, "/api/topup/create", bytes.NewReader(bodyJson))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err = handler.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to create topup")
}

func TestCreateTopup_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	req := &requests.CreateTopupRequest{
		CardNumber:  "",
		TopupNo:     "",
		TopupAmount: 0,
		TopupMethod: "",
	}

	mockTopupClient.EXPECT().CreateTopup(gomock.Any(), gomock.Any()).Times(0)
	mockLogger.EXPECT().Debug("Validation Error", gomock.Any()).Times(1)

	e := echo.New()
	bodyJson, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	httpReq := httptest.NewRequest(http.MethodPost, "/api/topup/create", bytes.NewReader(bodyJson))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err = handler.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error")
}

func TestUpdateTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponseTopup{
		Status:  "success",
		Message: "Topup updated successfully",
		Data: &pb.TopupResponse{
			Id:          1,
			CardNumber:  "1234567890",
			TopupAmount: 600000,
		},
	}
	req := &requests.UpdateTopupRequest{
		TopupID:     1,
		CardNumber:  "1234567890",
		TopupAmount: 600000,
		TopupMethod: "mandiri",
	}

	mockTopupClient.EXPECT().
		UpdateTopup(
			gomock.Any(),
			&pb.UpdateTopupRequest{
				TopupId:     int32(req.TopupID),
				CardNumber:  req.CardNumber,
				TopupAmount: int32(req.TopupAmount),
				TopupMethod: req.TopupMethod,
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	bodyJson, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/topup/update/1", bytes.NewReader(bodyJson))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Topup updated successfully", resp.Message)
	assert.Equal(t, "1234567890", resp.Data.CardNumber)
}

func TestUpdateTopup_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	invalidID := "abc"

	req := &requests.UpdateTopupRequest{
		TopupID:     1,
		CardNumber:  "1234567890",
		TopupAmount: 600000,
		TopupMethod: "mandiri",
	}

	mockLogger.EXPECT().Debug("Bad Request", gomock.Any()).Times(1)

	e := echo.New()
	bodyJson, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/topup/update/%s", invalidID), bytes.NewReader(bodyJson))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("id")
	c.SetParamValues(invalidID)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Bad Request: Invalid ID")
}

func TestUpdateTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	req := &requests.UpdateTopupRequest{
		TopupID:     1,
		CardNumber:  "1234567890",
		TopupAmount: 600000,
		TopupMethod: "mandiri",
	}

	mockTopupClient.EXPECT().
		UpdateTopup(
			gomock.Any(),
			&pb.UpdateTopupRequest{
				TopupId:     int32(req.TopupID),
				CardNumber:  req.CardNumber,
				TopupAmount: int32(req.TopupAmount),
				TopupMethod: req.TopupMethod,
			},
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update topup",
		})

	mockLogger.EXPECT().Debug("Failed to update topup", gomock.Any()).Times(1)

	e := echo.New()
	bodyJson, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPut, "/api/topup/update/1", bytes.NewReader(bodyJson))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to update topup")
}

func TestUpdateTopup_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	req := &requests.UpdateTopupRequest{
		TopupID:     0,
		CardNumber:  "",
		TopupAmount: 0,
		TopupMethod: "",
	}

	mockLogger.EXPECT().Debug("Bad Request", gomock.Any()).Times(1)

	e := echo.New()
	bodyJson, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPut, "/api/topup/update", bytes.NewReader(bodyJson))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Bad Request")
}

func TestTrashTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponseTopup{
		Status:  "success",
		Message: "Topup trashed successfully",
	}

	mockTopupClient.EXPECT().
		TrashedTopup(
			gomock.Any(),
			&pb.FindByIdTopupRequest{TopupId: 1},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/topup/1/trash", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.TrashTopup(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Topup trashed successfully", resp.Message)
}

func TestTrashTopup_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/topup/trash/ab", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("ab")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.TrashTopup(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Bad Request: Invalid ID", resp.Message)
}

func TestTrashTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockTopupClient.EXPECT().
		TrashedTopup(
			gomock.Any(),
			&pb.FindByIdTopupRequest{TopupId: 1},
		).
		Return(nil, fmt.Errorf("internal server error"))

	mockLogger.EXPECT().Debug("Failed to trashed topup", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/topup/1/trash", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.TrashTopup(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to trashed topup")
}

func TestRestoreTopup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponseTopup{
		Status:  "success",
		Message: "Topup restored successfully",
		Data: &pb.TopupResponse{
			Id:          1,
			CardNumber:  "1234567890",
			TopupAmount: 500000,
		},
	}

	mockTopupClient.EXPECT().
		RestoreTopup(
			gomock.Any(),
			&pb.FindByIdTopupRequest{TopupId: 1},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/topup/1/restore", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.RestoreTopup(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Topup restored successfully", resp.Message)
	assert.Equal(t, int32(1), resp.Data.Id)
	assert.Equal(t, "1234567890", resp.Data.CardNumber)
	assert.Equal(t, int32(500000), resp.Data.TopupAmount)
}

func TestRestoreTopup_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockTopupClient.EXPECT().
		RestoreTopup(
			gomock.Any(),
			&pb.FindByIdTopupRequest{TopupId: 1},
		).
		Return(nil, fmt.Errorf("internal server error"))

	mockLogger.EXPECT().Debug("Failed to restore topup", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/topup/1/restore", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.RestoreTopup(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to restore topup")
}

func TestRestoreTopup_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/topup/invalid/restore", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	mockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.RestoreTopup(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Bad Request: Invalid ID", resp.Message)
}

func TestDeleteTopupPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponseTopupDelete{
		Status:  "success",
		Message: "Topup deleted permanently",
	}

	mockTopupClient.EXPECT().
		DeleteTopupPermanent(
			gomock.Any(),
			&pb.FindByIdTopupRequest{TopupId: 1},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/topup/1/permanent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.DeleteTopupPermanent(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Topup deleted permanently", resp.Message)
}

func TestDeleteTopupPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockTopupClient.EXPECT().
		DeleteTopupPermanent(
			gomock.Any(),
			&pb.FindByIdTopupRequest{TopupId: 1},
		).
		Return(nil, fmt.Errorf("internal server error"))

	mockLogger.EXPECT().Debug("Failed to delete topup", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/topup/permanent/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.DeleteTopupPermanent(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to delete topup")
}

func TestDeleteTopupPermanent_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopupClient := mock_pb.NewMockTopupServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/topup/invalid/permanent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	mockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	handler := api.NewHandlerTopup(mockTopupClient, e, mockLogger)

	err := handler.DeleteTopupPermanent(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Bad Request: Invalid ID", resp.Message)
}
