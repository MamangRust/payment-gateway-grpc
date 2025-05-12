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
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindAllTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationTransfer{
		Status:  "success",
		Message: "Successfully fetch transfers",
		Pagination: &pb.PaginationMeta{
			CurrentPage:  1,
			TotalPages:   1,
			TotalRecords: 1,
		},
		Data: []*pb.TransferResponse{
			{
				Id:             1,
				TransferFrom:   "test",
				TransferTo:     "test",
				TransferAmount: 10000,
				TransferTime:   "2022-01-01 00:00:00",
				CreatedAt:      "2022-01-01 00:00:00",
				UpdatedAt:      "2022-01-01 00:00:00",
			},
		},
	}

	mockTransferClient.EXPECT().FindAllTransfer(gomock.Any(), &pb.FindAllTransferRequest{Page: 1, PageSize: 10, Search: ""}).Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationTransfer
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Successfully fetch transfers", resp.Message)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, int32(1), resp.Data[0].Id)
}

func TestFindAllTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockTransferClient.EXPECT().
		FindAllTransfer(
			gomock.Any(),
			&pb.FindAllTransferRequest{
				Page:     1,
				PageSize: 10,
				Search:   "",
			},
		).
		Return(nil, fmt.Errorf("some internal error"))

	mockLogger.EXPECT().Debug("Failed to retrieve transfer data: ", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve transfer data: ", resp.Message)
}

func TestFindAllTransfer_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockTransferClient.EXPECT().
		FindAllTransfer(gomock.Any(), &pb.FindAllTransferRequest{Page: 1, PageSize: 10, Search: ""}).
		Return(&pb.ApiResponsePaginationTransfer{
			Status:     "success",
			Message:    "No transfers found",
			Data:       []*pb.TransferResponse{},
			Pagination: &pb.PaginationMeta{},
		}, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationTransfer
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No transfers found", resp.Message)
	assert.Len(t, resp.Data, 0)
}

func TestFindByIdTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferID := 1
	expectedGRPCRequest := &pb.FindByIdTransferRequest{
		TransferId: int32(transferID),
	}

	expectedGRPCResponse := &pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: &pb.TransferResponse{
			Id:             1,
			TransferFrom:   "test",
			TransferTo:     "test",
			TransferAmount: 10000,
			TransferTime:   "2022-01-01 00:00:00",
			CreatedAt:      "2022-01-01 00:00:00",
			UpdatedAt:      "2022-01-01 00:00:00",
		},
	}

	mockTransferClient.EXPECT().
		FindByIdTransfer(gomock.Any(), expectedGRPCRequest).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransfer
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transfer retrieved successfully", resp.Message)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, expectedGRPCResponse.Data, resp.Data)
}

func TestFindByIdTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferID := 1
	expectedGRPCRequest := &pb.FindByIdTransferRequest{
		TransferId: int32(transferID),
	}

	mockTransferClient.EXPECT().
		FindByIdTransfer(gomock.Any(), expectedGRPCRequest).
		Return(nil, errors.New("Failed to retrieve transfer data: ")).
		Times(1)

	mockLogger.EXPECT().
		Debug("Failed to retrieve transfer data: ", gomock.Any()).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to retrieve transfer data: ")
}

func TestFindByIdTransfer_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	mockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Bad Request: Invalid ID", resp.Message)
}

func TestFindByTransferByTransferFrom_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferFrom := "test_user"

	expectedGRPCRequest := &pb.FindTransferByTransferFromRequest{
		TransferFrom: transferFrom,
	}

	expectedGRPCResponse :=
		&pb.ApiResponseTransfers{
			Status:  "success",
			Message: "Transfer retrieved successfully",
			Data: []*pb.TransferResponse{
				{
					Id:             1,
					TransferFrom:   transferFrom,
					TransferTo:     "test_to",
					TransferAmount: 10000,
					TransferTime:   "2022-01-01 00:00:00",
				},
			},
		}

	mockTransferClient.EXPECT().
		FindTransferByTransferFrom(gomock.Any(), expectedGRPCRequest).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer/from/test_user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("transfer_from")
	c.SetParamValues("test_user")

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindByTransferByTransferFrom(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransfers
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transfer retrieved successfully", resp.Message)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, expectedGRPCResponse.Data, resp.Data)
}

func TestFindByTransferByTransferFrom_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferFrom := "test_user"
	expectedGRPCRequest := &pb.FindTransferByTransferFromRequest{
		TransferFrom: transferFrom,
	}

	mockTransferClient.EXPECT().
		FindTransferByTransferFrom(gomock.Any(), expectedGRPCRequest).
		Return(nil, errors.New("Failed to retrieve transfer data")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve transfer data: ", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer/from/test_user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("transfer_from")
	c.SetParamValues("test_user")

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindByTransferByTransferFrom(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve transfer data: ", resp.Message)
}

func TestFindByTransferByTransferTo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferTo := "test_to"
	expectedGRPCRequest := &pb.FindTransferByTransferToRequest{
		TransferTo: transferTo,
	}

	expectedGRPCResponse :=
		&pb.ApiResponseTransfers{
			Status:  "success",
			Message: "Transfer retrieved successfully",
			Data: []*pb.TransferResponse{
				{
					Id:             1,
					TransferFrom:   "test_from",
					TransferTo:     transferTo,
					TransferAmount: 10000,
					TransferTime:   "2022-01-01 00:00:00",
				},
			},
		}

	mockTransferClient.EXPECT().
		FindTransferByTransferTo(gomock.Any(), expectedGRPCRequest).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer/to/test_to", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("transfer_to")
	c.SetParamValues("test_to")

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindByTransferByTransferTo(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransfers
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transfer retrieved successfully", resp.Message)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, expectedGRPCResponse.Data, resp.Data)
}

func TestFindByTransferByTransferTo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	transferTo := "test_to"
	expectedGRPCRequest := &pb.FindTransferByTransferToRequest{
		TransferTo: transferTo,
	}

	mockTransferClient.EXPECT().
		FindTransferByTransferTo(gomock.Any(), expectedGRPCRequest).
		Return(nil, errors.New("Failed to retrieve transfer data")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve transfer data: ", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer/to/test_to", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("transfer_to")
	c.SetParamValues("test_to")

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindByTransferByTransferTo(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve transfer data: ", resp.Message)
}

func TestFindByActiveTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedGRPCResponse := &pb.ApiResponsePaginationTransferDeleteAt{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: []*pb.TransferResponseDeleteAt{
			{
				Id:             1,
				TransferFrom:   "test_from",
				TransferTo:     "test_to",
				TransferAmount: 10000,
				TransferTime:   "2022-01-01 00:00:00",
			},
		},
	}

	request := &pb.FindAllTransferRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockTransferClient.EXPECT().
		FindByActiveTransfer(gomock.Any(), request).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindByActiveTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransfers
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transfer retrieved successfully", resp.Message)
}

func TestFindByActiveTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllTransferRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockTransferClient.EXPECT().
		FindByActiveTransfer(gomock.Any(), request).
		Return(nil, errors.New("Failed to retrieve transfer data")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve transfer data: ", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindByActiveTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve transfer data: ", resp.Message)
}

func TestFindByTrashedTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllTransferRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	expectedGRPCResponse := &pb.ApiResponsePaginationTransferDeleteAt{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: []*pb.TransferResponseDeleteAt{
			{
				Id:             1,
				TransferFrom:   "test_from",
				TransferTo:     "test_to",
				TransferAmount: 10000,
				TransferTime:   "2022-01-01 00:00:00",
			},
		},
	}

	mockTransferClient.EXPECT().
		FindByTrashedTransfer(gomock.Any(), request).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindByTrashedTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransfers
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transfer retrieved successfully", resp.Message)
}

func TestFindByTrashedTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllTransferRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockTransferClient.EXPECT().
		FindByTrashedTransfer(gomock.Any(), request).
		Return(nil, errors.New("Failed to retrieve transfer data")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve transfer data: ", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transfer/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.FindByTrashedTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve transfer data: ", resp.Message)
}

func TestCreateTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.CreateTransferRequest{
		TransferFrom:   "test_from",
		TransferTo:     "test_to",
		TransferAmount: 500000,
	}

	expectedGRPCRequest := &pb.CreateTransferRequest{
		TransferFrom:   body.TransferFrom,
		TransferTo:     body.TransferTo,
		TransferAmount: int32(body.TransferAmount),
	}

	expectedGRPCResponse := &pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Transfer created successfully",
		Data: &pb.TransferResponse{
			Id:             1,
			TransferFrom:   "test_from",
			TransferTo:     "test_to",
			TransferAmount: 500000,
			TransferTime:   "2022-01-01 00:00:00",
		},
	}

	mockTransferClient.EXPECT().
		CreateTransfer(gomock.Any(), expectedGRPCRequest).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	requestBodyBytes, err := json.Marshal(body)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/transfer", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err = handler.CreateTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransfer
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transfer created successfully", resp.Message)
}

func TestCreateTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.CreateTransferRequest{
		TransferFrom:   "test_from",
		TransferTo:     "test_to",
		TransferAmount: 500000,
	}

	expectedGRPCRequest := &pb.CreateTransferRequest{
		TransferFrom:   body.TransferFrom,
		TransferTo:     body.TransferTo,
		TransferAmount: int32(body.TransferAmount),
	}

	mockTransferClient.EXPECT().
		CreateTransfer(gomock.Any(), expectedGRPCRequest).
		Return(nil, errors.New("Failed to create transfer")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to create transfer: ", gomock.Any()).Times(1)

	e := echo.New()
	requestBodyBytes, err := json.Marshal(body)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/transfer", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err = handler.CreateTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to create transfer: ", resp.Message)
}

func TestCreateTransfer_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.CreateTransferRequest{
		TransferFrom: "test_from",
		TransferTo:   "test_to",
	}

	e := echo.New()
	requestBodyBytes, err := json.Marshal(body)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/transfer", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockLogger.EXPECT().Debug("Validation Error: ", gomock.Any()).Times(1)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err = handler.CreateTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Validation Error: Key: 'CreateTransferRequest.TransferAmount' Error:Field validation for 'TransferAmount' failed on the 'required' tag", resp.Message)
}

func TestUpdateTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.UpdateTransferRequest{
		TransferID:     1,
		TransferFrom:   "test_from",
		TransferTo:     "test_to",
		TransferAmount: 500000,
	}

	expectedGRPCRequest := &pb.UpdateTransferRequest{
		TransferId:     int32(body.TransferID),
		TransferFrom:   body.TransferFrom,
		TransferTo:     body.TransferTo,
		TransferAmount: int32(body.TransferAmount),
	}

	expectedGRPCResponse := &pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Transfer updated successfully",
		Data: &pb.TransferResponse{
			Id:             1,
			TransferFrom:   "test_from",
			TransferTo:     "test_to",
			TransferAmount: 500000,
			TransferTime:   "2022-01-01 00:00:00",
		},
	}

	mockTransferClient.EXPECT().
		UpdateTransfer(gomock.Any(), expectedGRPCRequest).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	requestBodyBytes, err := json.Marshal(body)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/transfer", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err = handler.UpdateTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransfer
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transfer updated successfully", resp.Message)
}

func TestUpdateTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.UpdateTransferRequest{
		TransferID:     1,
		TransferFrom:   "test_from",
		TransferTo:     "test_to",
		TransferAmount: 500000,
	}

	expectedGRPCRequest := &pb.UpdateTransferRequest{
		TransferId:     int32(body.TransferID),
		TransferFrom:   body.TransferFrom,
		TransferTo:     body.TransferTo,
		TransferAmount: int32(body.TransferAmount),
	}

	mockTransferClient.EXPECT().
		UpdateTransfer(gomock.Any(), expectedGRPCRequest).
		Return(nil, errors.New("Failed to update transfer")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to update transfer: ", gomock.Any()).Times(1)

	e := echo.New()
	requestBodyBytes, err := json.Marshal(body)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/transfer", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err = handler.UpdateTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to update transfer: ", resp.Message)
}

func TestUpdateTransfer_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	body := requests.UpdateTransferRequest{
		TransferFrom: "test_from",
		TransferTo:   "test_to",
	}

	e := echo.New()
	requestBodyBytes, err := json.Marshal(body)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/transfer", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockLogger.EXPECT().Debug("Validation Error: ", gomock.Any()).Times(1)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err = handler.UpdateTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error: ")
}

func TestTrashTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	id := 1
	expectedGRPCResponse := &pb.ApiResponseTransfer{
		Status: "success",
		Data: &pb.TransferResponse{
			Id:             1,
			TransferFrom:   "test_from",
			TransferTo:     "test_to",
			TransferAmount: 500000,
			TransferTime:   "2022-01-01 00:00:00",
		},
		Message: "Transfer trashed successfully",
	}

	mockTransferClient.EXPECT().
		TrashedTransfer(gomock.Any(), &pb.FindByIdTransferRequest{TransferId: int32(id)}).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/transfer/%d/trash", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.TrashTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransfer
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transfer trashed successfully", resp.Message)
}

func TestTrashTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	id := 1
	mockTransferClient.EXPECT().
		TrashedTransfer(gomock.Any(), &pb.FindByIdTransferRequest{TransferId: int32(id)}).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to trash transfer: ", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/transfer/%d/trash", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.TrashTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to trash transfer")
}

func TestTrashTransfer_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	id := "invalid_id"
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/transfer/%s/trash", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id)

	mockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.TrashTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Bad Request: Invalid ID")
}

func TestRestoreTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	id := 1
	expectedGRPCResponse := &pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Transfer restored successfully",
		Data: &pb.TransferResponse{
			Id:             1,
			TransferFrom:   "test_from",
			TransferTo:     "test_to",
			TransferAmount: 500000,
			TransferTime:   "2022-01-01 00:00:00",
		},
	}

	mockTransferClient.EXPECT().
		RestoreTransfer(gomock.Any(), &pb.FindByIdTransferRequest{TransferId: int32(id)}).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/transfer/%d/restore", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.RestoreTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransfer
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transfer restored successfully", resp.Message)

	data := resp.Data

	assert.NoError(t, err)
	assert.Equal(t, 1, int(data.Id))
	assert.Equal(t, "test_from", data.TransferFrom)
	assert.Equal(t, "test_to", data.TransferTo)
	assert.Equal(t, 500000, int(data.TransferAmount))
	assert.Equal(t, "2022-01-01 00:00:00", data.TransferTime)
}

func TestRestoreTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	id := 1
	mockTransferClient.EXPECT().
		RestoreTransfer(gomock.Any(), &pb.FindByIdTransferRequest{TransferId: int32(id)}).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to restore transfer: ", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/transfer/%d/restore", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.RestoreTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to restore transfer:")
}

func TestRestoreTransfer_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	invalidID := "invalid_id"
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/transfer/%s/restore", invalidID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(invalidID)

	mockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.RestoreTransfer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Bad Request: Invalid ID")
}

func TestDeleteTransferPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	id := 1
	expectedGRPCResponse := &pb.ApiResponseTransferDelete{
		Status:  "success",
		Message: "Transfer deleted permanently",
	}

	mockTransferClient.EXPECT().
		DeleteTransferPermanent(gomock.Any(), &pb.FindByIdTransferRequest{TransferId: int32(id)}).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/transfer/%d/delete", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.DeleteTransferPermanent(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransferDelete
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transfer deleted permanently", resp.Message)
}

func TestDeleteTransferPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	id := 1
	mockTransferClient.EXPECT().
		DeleteTransferPermanent(gomock.Any(), &pb.FindByIdTransferRequest{TransferId: int32(id)}).
		Return(nil, errors.New("internal server error")).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/transfer/%d/delete", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.DeleteTransferPermanent(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to delete transfer")
}

func TestDeleteTransferPermanent_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferClient := mock_pb.NewMockTransferServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	invalidID := "invalid_id"
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/transfer/%s/delete", invalidID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(invalidID)

	mockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	handler := api.NewHandlerTransfer(mockTransferClient, e, mockLogger)

	err := handler.DeleteTransferPermanent(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Bad Request: Invalid ID")
}
