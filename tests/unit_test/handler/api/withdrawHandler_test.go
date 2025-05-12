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
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFindAllWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationWithdraw{
		Status:  "success",
		Message: "Withdraws retrieved successfully",
		Data: []*pb.WithdrawResponse{
			{
				WithdrawId: 1,
				CardNumber: "1234567890123456",
				CreatedAt:  "2022-01-01T00:00:00Z",
				UpdatedAt:  "2022-01-01T00:00:00Z",
			},
			{
				WithdrawId: 2,
				CardNumber: "9876543210987654",
				CreatedAt:  "2022-01-02T00:00:00Z",
				UpdatedAt:  "2022-01-02T00:00:00Z",
			},
		},
		Pagination: &pb.PaginationMeta{
			CurrentPage: 1,
			PageSize:    2,
			TotalPages:  1,
		},
	}

	mockWithdrawClient.EXPECT().
		FindAllWithdraw(gomock.Any(), &pb.FindAllWithdrawRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Withdraws retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 2)
}

func TestFindAllWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockWithdrawClient.EXPECT().
		FindAllWithdraw(gomock.Any(), &pb.FindAllWithdrawRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}).
		Return(nil, fmt.Errorf("internal server error"))

	mockLogger.EXPECT().Debug("Failed to retrieve withdraw data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve withdraw data: ", resp.Message)
}

func TestFindAllWithdraw_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockWithdrawClient.EXPECT().
		FindAllWithdraw(gomock.Any(), &pb.FindAllWithdrawRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}).
		Return(&pb.ApiResponsePaginationWithdraw{
			Status:     "success",
			Message:    "No withdraws found",
			Data:       []*pb.WithdrawResponse{},
			Pagination: &pb.PaginationMeta{},
		}, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No withdraws found", resp.Message)
	assert.Len(t, resp.Data, 0)
}

func TestFindByIdWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.WithdrawResponse{
		WithdrawId: 1,
		CardNumber: "1234567890123456",
		CreatedAt:  "2022-01-01T00:00:00Z",
		UpdatedAt:  "2022-01-01T00:00:00Z",
	}

	expect := &pb.ApiResponseWithdraw{
		Status:  "success",
		Message: "Withdraw retrieved successfully",
		Data:    expectedResponse,
	}

	mockWithdrawClient.EXPECT().
		FindByIdWithdraw(gomock.Any(), &pb.FindByIdWithdrawRequest{
			WithdrawId: 1,
		}).
		Return(expect, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Withdraw retrieved successfully", resp.Message)

	assert.Equal(t, expectedResponse.WithdrawId, resp.Data.WithdrawId)
	assert.Equal(t, expectedResponse.CardNumber, resp.Data.CardNumber)
	assert.Equal(t, expectedResponse.CreatedAt, resp.Data.CreatedAt)
	assert.Equal(t, expectedResponse.UpdatedAt, resp.Data.UpdatedAt)
}

func TestFindByIdWithdraw_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	mockLogger.EXPECT().Debug("Invalid withdraw ID", gomock.Any()).Times(1)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid withdraw ID", resp.Message)
}

func TestFindByCardNumberWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := []*pb.WithdrawResponse{
		{

			WithdrawId:     1,
			CardNumber:     "1234567890123456",
			WithdrawAmount: 50000,
			CreatedAt:      "2022-01-01T00:00:00Z",
			UpdatedAt:      "2022-01-01T00:00:00Z",
		},
		{

			WithdrawId:     2,
			CardNumber:     "9876543210987654",
			WithdrawAmount: 75000,
			CreatedAt:      "2022-01-02T00:00:00Z",
			UpdatedAt:      "2022-01-02T00:00:00Z",
		},
	}
	expected := &pb.ApiResponsesWithdraw{
		Status:  "success",
		Message: "Withdraw retrieved successfully",
		Data:    expectedResponse,
	}

	mockWithdrawClient.EXPECT().
		FindByCardNumber(gomock.Any(), &pb.FindByCardNumberRequest{
			CardNumber: "1234567890123456",
		}).
		Return(expected, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw?card_number=1234567890123456", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.FindByCardNumber(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Withdraw retrieved successfully", resp.Message)

	assert.Equal(t, expectedResponse[0].WithdrawId, resp.Data[0].WithdrawId)
	assert.Equal(t, expectedResponse[0].CardNumber, resp.Data[0].CardNumber)
	assert.Equal(t, expectedResponse[0].WithdrawAmount, resp.Data[0].WithdrawAmount)
}

func TestFindByCardNumberWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockWithdrawClient.EXPECT().
		FindByCardNumber(gomock.Any(), &pb.FindByCardNumberRequest{
			CardNumber: "1234567890123456",
		}).
		Return(nil, fmt.Errorf("withdraw not found"))

	mockLogger.EXPECT().Debug("Failed to retrieve withdraw data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw?card_number=1234567890123456", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.FindByCardNumber(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve withdraw data: ", resp.Message)
}

func TestFindByActiveWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationWithdrawDeleteAt{
		Status:  "success",
		Message: "Withdraw retrieved successfully",
		Data: []*pb.WithdrawResponseDeleteAt{
			{
				WithdrawId:     1,
				CardNumber:     "1234567890123456",
				WithdrawAmount: 50000,
				CreatedAt:      "2022-01-01T00:00:00Z",
				UpdatedAt:      "2022-01-01T00:00:00Z",
			},
			{
				WithdrawId:     2,
				CardNumber:     "9876543210987654",
				WithdrawAmount: 75000,
				CreatedAt:      "2022-01-02T00:00:00Z",
				UpdatedAt:      "2022-01-02T00:00:00Z",
			},
		},
	}
	request := &pb.FindAllWithdrawRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockWithdrawClient.EXPECT().
		FindByActive(gomock.Any(), request).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.FindByActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Withdraw retrieved successfully", resp.Message)

	assert.Equal(t, expectedResponse.Data[0].WithdrawId, resp.Data[0].WithdrawId)
	assert.Equal(t, expectedResponse.Data[0].CardNumber, resp.Data[0].CardNumber)
	assert.Equal(t, expectedResponse.Data[0].WithdrawAmount, resp.Data[0].WithdrawAmount)
}

func TestFindByActiveWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockWithdrawClient.EXPECT().
		FindByActive(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("internal server error"))

	mockLogger.EXPECT().Debug("Failed to retrieve withdraw data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.FindByActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve withdraw data: ", resp.Message)
}

func TestFindByTrashedWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationWithdrawDeleteAt{
		Status:  "success",
		Message: "Withdraw retrieved successfully",
		Data: []*pb.WithdrawResponseDeleteAt{
			{
				WithdrawId:     1,
				CardNumber:     "1234567890123456",
				WithdrawAmount: 50000,
				CreatedAt:      "2022-01-01T00:00:00Z",
				UpdatedAt:      "2022-01-01T00:00:00Z",
			},
		},
	}

	request := &pb.FindAllWithdrawRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockWithdrawClient.EXPECT().
		FindByTrashed(gomock.Any(), request).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.FindByTrashed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Withdraw retrieved successfully", resp.Message)

	assert.Equal(t, expectedResponse.Data[0].WithdrawId, resp.Data[0].WithdrawId)
	assert.Equal(t, expectedResponse.Data[0].CardNumber, resp.Data[0].CardNumber)
	assert.Equal(t, expectedResponse.Data[0].WithdrawAmount, resp.Data[0].WithdrawAmount)
}

func TestFindByTrashedWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockWithdrawClient.EXPECT().
		FindByTrashed(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("database connection error"))

	mockLogger.EXPECT().Debug("Failed to retrieve withdraw data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.FindByTrashed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve withdraw data: ", resp.Message)
}

func TestCreateWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.CreateWithdrawRequest{
		CardNumber:     "1234567890123456",
		WithdrawAmount: 50000,
		WithdrawTime:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	expectedGRPCRequest := &pb.CreateWithdrawRequest{
		CardNumber:     requestBody.CardNumber,
		WithdrawAmount: int32(requestBody.WithdrawAmount),
		WithdrawTime:   timestamppb.New(requestBody.WithdrawTime),
	}

	expectedGRPCResponse := &pb.WithdrawResponse{
		WithdrawId:     1,
		CardNumber:     "1234567890123456",
		WithdrawAmount: 50000,
		WithdrawTime:   "2022-01-01T00:00:00Z",
	}

	expectedAPIResponse := &pb.ApiResponseWithdraw{
		Status:  "success",
		Message: "Withdraw created successfully",
		Data:    expectedGRPCResponse,
	}

	mockWithdrawClient.EXPECT().
		CreateWithdraw(gomock.Any(), expectedGRPCRequest).
		Return(expectedAPIResponse, nil)

	e := echo.New()
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPost, "/api/withdraw", bytes.NewReader(requestBodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err = handler.Create(c)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Withdraw created successfully", resp.Message)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, expectedGRPCResponse, resp.Data)
}

func TestCreateWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.CreateWithdrawRequest{
		CardNumber:     "1234567890123456",
		WithdrawAmount: 50000,
		WithdrawTime:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	expectedGRPCRequest := &pb.CreateWithdrawRequest{
		CardNumber:     requestBody.CardNumber,
		WithdrawAmount: int32(requestBody.WithdrawAmount),
		WithdrawTime:   timestamppb.New(requestBody.WithdrawTime),
	}

	mockWithdrawClient.EXPECT().
		CreateWithdraw(gomock.Any(), expectedGRPCRequest).
		Return(nil, fmt.Errorf("internal server error"))

	mockLogger.EXPECT().Debug("Failed to create withdraw", gomock.Any()).Times(1)

	e := echo.New()
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPost, "/api/withdraw", bytes.NewReader(requestBodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err = handler.Create(c)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to create withdraw: internal server error", resp.Message)
}

func TestCreateWithdraw_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.CreateWithdrawRequest{
		CardNumber:     "",
		WithdrawAmount: 50000,
		WithdrawTime:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	e := echo.New()
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPost, "/api/withdraw", bytes.NewReader(requestBodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	mockLogger.EXPECT().Debug("Validation Error: Key: 'CreateWithdrawRequest.CardNumber' Error:Field validation for 'CardNumber' failed on the 'required' tag", gomock.Any()).Times(1)

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err = handler.Create(c)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error: ")
}

func TestUpdateWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.UpdateWithdrawRequest{
		WithdrawID:     1,
		CardNumber:     "1234567890123456",
		WithdrawAmount: 50000,
		WithdrawTime:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	withdrawID := 1

	expectedResponse := &pb.WithdrawResponse{
		WithdrawId:     int32(withdrawID),
		CardNumber:     "1234567890123456",
		WithdrawAmount: 50000,
		WithdrawTime:   "2022-01-01T00:00:00Z",
	}

	mockResponse := &pb.ApiResponseWithdraw{
		Status:  "success",
		Message: "Successfully updated withdraw",
		Data:    expectedResponse,
	}

	mockWithdrawClient.EXPECT().
		UpdateWithdraw(gomock.Any(), &pb.UpdateWithdrawRequest{
			WithdrawId:     int32(withdrawID),
			CardNumber:     body.CardNumber,
			WithdrawAmount: int32(body.WithdrawAmount),
			WithdrawTime:   timestamppb.New(body.WithdrawTime),
		}).
		Return(mockResponse, nil).
		Times(1)

	e := echo.New()
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/withdraw/%d", withdrawID), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", withdrawID))

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, mockResponse.Data.WithdrawId, resp.Data.WithdrawId)
	assert.Equal(t, mockResponse.Data.CardNumber, resp.Data.CardNumber)
	assert.Equal(t, mockResponse.Data.WithdrawAmount, resp.Data.WithdrawAmount)
}

func TestUpdateWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.UpdateWithdrawRequest{
		CardNumber:     "1234567890123456",
		WithdrawAmount: 50000,
		WithdrawTime:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	mockLogger.EXPECT().Debug(
		"Validation Error: Key: 'UpdateWithdrawRequest.WithdrawID' Error:Field validation for 'WithdrawID' failed on the 'required' tag",
		gomock.Any()).Times(1)

	mockWithdrawClient.EXPECT().
		UpdateWithdraw(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("service unavailable")).
		Times(0)
	e := echo.New()
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/withdraw/%d", 1), bytes.NewReader(bodyBytes)) // Request path should have the correct ID
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error: Key: 'UpdateWithdrawRequest.WithdrawID' Error:Field validation for 'WithdrawID' failed on the 'required' tag")
}

func TestUpdateWithdraw_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := requests.UpdateWithdrawRequest{
		CardNumber:     "",
		WithdrawID:     1,
		WithdrawAmount: 0,
		WithdrawTime:   time.Time{},
	}
	expectedValidationError := "Validation Error: " +
		"Key: 'UpdateWithdrawRequest.CardNumber' Error:Field validation for 'CardNumber' failed on the 'required' tag\n" +
		"Key: 'UpdateWithdrawRequest.WithdrawAmount' Error:Field validation for 'WithdrawAmount' failed on the 'required' tag\n" +
		"Key: 'UpdateWithdrawRequest.WithdrawTime' Error:Field validation for 'WithdrawTime' failed on the 'required' tag"

	mockLogger.EXPECT().Debug(expectedValidationError, gomock.Any()).Times(1)

	e := echo.New()
	requestBodyBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	httpReq := httptest.NewRequest(http.MethodPut, "/api/withdraw/update/1", bytes.NewReader(requestBodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)
	err = handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	fmt.Println("Error message:", resp.Message)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error:")
}

func TestTrashWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockResponse := &pb.ApiResponseWithdraw{
		Status:  "success",
		Message: "Successfully trashed withdraw",
		Data: &pb.WithdrawResponse{
			WithdrawId:     1,
			CardNumber:     "1234567890123456",
			WithdrawAmount: 50000,
			WithdrawTime:   "2022-01-01T00:00:00Z",
		},
	}

	// Mock the gRPC call
	mockWithdrawClient.EXPECT().
		TrashedWithdraw(gomock.Any(), &pb.FindByIdWithdrawRequest{WithdrawId: 1}).
		Return(mockResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/withdraw/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	// Call the handler
	err := handler.TrashWithdraw(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Successfully trashed withdraw", resp.Message)
}

func TestTrashWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockWithdrawClient.EXPECT().
		TrashedWithdraw(gomock.Any(), &pb.FindByIdWithdrawRequest{WithdrawId: 1}).
		Return(nil, fmt.Errorf("service unavailable")).
		Times(1)

	mockLogger.EXPECT().
		Debug("Failed to trash withdraw", gomock.Any()).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/withdraw/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.TrashWithdraw(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to trash withdraw")
}

func TestTrashWithdraw_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().
		Debug("Invalid withdraw ID", gomock.Any()).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/withdraw/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)

	err := handler.TrashWithdraw(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid withdraw ID", resp.Message)
}

func TestRestoreWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	id := 1

	mockWithdrawClient.EXPECT().
		RestoreWithdraw(gomock.Any(), &pb.FindByIdWithdrawRequest{
			WithdrawId: int32(id),
		}).
		Return(&pb.ApiResponseWithdraw{
			Status:  "success",
			Message: "Withdraw restored successfully",
		}, nil).Times(1)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/withdraw/restore/%d", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)
	err := handler.RestoreWithdraw(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Withdraw restored successfully", resp.Message)
}

func TestRestoreWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	id := 1

	mockWithdrawClient.EXPECT().
		RestoreWithdraw(gomock.Any(), &pb.FindByIdWithdrawRequest{
			WithdrawId: int32(id),
		}).
		Return(nil, fmt.Errorf("gRPC service unavailable")).Times(1)

	mockLogger.EXPECT().Debug("Failed to restore withdraw", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/withdraw/restore/%d", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)
	err := handler.RestoreWithdraw(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to restore withdraw")
}

func TestRestoreWithdraw_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	e := echo.New()
	id := "invalid"

	mockLogger.EXPECT().Debug("Invalid withdraw ID", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/withdraw/restore/%s", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id)

	handler := api.NewHandlerWithdraw(nil, e, mockLogger)
	err := handler.RestoreWithdraw(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid withdraw ID", resp.Message)
}

func TestDeleteWithdrawPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	id := 1

	mockWithdrawClient.EXPECT().
		DeleteWithdrawPermanent(gomock.Any(), &pb.FindByIdWithdrawRequest{
			WithdrawId: int32(id),
		}).
		Return(&pb.ApiResponseWithdrawDelete{
			Status:  "success",
			Message: "Withdraw deleted permanently",
		}, nil).Times(1)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/withdraw/delete-permanent/%d", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)
	err := handler.DeleteWithdrawPermanent(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseWithdrawDelete
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Withdraw deleted permanently", resp.Message)
}

func TestDeleteWithdrawPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawClient := mock_pb.NewMockWithdrawServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	id := 1

	mockWithdrawClient.EXPECT().
		DeleteWithdrawPermanent(gomock.Any(), &pb.FindByIdWithdrawRequest{
			WithdrawId: int32(id),
		}).
		Return(nil, fmt.Errorf("gRPC service unavailable")).Times(1)

	mockLogger.EXPECT().Debug("Failed to delete withdraw permanently", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/withdraw/delete-permanent/%d", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))

	handler := api.NewHandlerWithdraw(mockWithdrawClient, e, mockLogger)
	err := handler.DeleteWithdrawPermanent(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to delete withdraw permanently")
}

func TestDeleteWithdrawPermanent_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	e := echo.New()
	id := "invalid"

	mockLogger.EXPECT().Debug("Invalid withdraw ID", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/withdraw/delete-permanent/%s", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id)

	handler := api.NewHandlerWithdraw(nil, e, mockLogger)
	err := handler.DeleteWithdrawPermanent(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid withdraw ID", resp.Message)
}
