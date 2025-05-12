package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/handler/api"
	"MamangRust/paymentgatewaygrpc/internal/middlewares"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	mock_pb "MamangRust/paymentgatewaygrpc/internal/pb/mocks"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"
	"MamangRust/paymentgatewaygrpc/tests/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFindAllTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationTransaction{
		Status:  "success",
		Message: "Transactions retrieved successfully",
		Data: []*pb.TransactionResponse{
			{
				Id:         1,
				CardNumber: "1234567890123456",
			},
			{
				Id:         2,
				CardNumber: "1234567890123457",
			},
		},
		Pagination: &pb.PaginationMeta{
			CurrentPage: 1,
			PageSize:    2,
			TotalPages:  1,
		},
	}

	mockTransactionClient.EXPECT().
		FindAllTransaction(
			gomock.Any(),
			&pb.FindAllTransactionRequest{
				Page:     1,
				PageSize: 10,
				Search:   "",
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transactions?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationTransaction
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transactions retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 2)
}

func TestFindAllTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockTransactionClient.EXPECT().
		FindAllTransaction(
			gomock.Any(),
			&pb.FindAllTransactionRequest{
				Page:     1,
				PageSize: 10,
				Search:   "",
			},
		).
		Return(nil, fmt.Errorf("some internal error"))

	mockLogger.EXPECT().Debug("Failed to retrieve transaction data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transactions?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve transaction data: ", resp.Message)
}

func TestFindAllTransaction_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationTransaction{
		Status:  "success",
		Message: "No transactions found",
		Data:    []*pb.TransactionResponse{},
		Pagination: &pb.PaginationMeta{
			CurrentPage: 1,
			PageSize:    10,
			TotalPages:  1,
		},
	}

	mockTransactionClient.EXPECT().
		FindAllTransaction(
			gomock.Any(),
			&pb.FindAllTransactionRequest{
				Page:     1,
				PageSize: 10,
				Search:   "",
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transactions?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationTransaction
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No transactions found", resp.Message)
	assert.Len(t, resp.Data, 0)
}

func TestFindById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	id := 1
	expectedGRPCResponse := &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Transaction retrieved successfully",
		Data: &pb.TransactionResponse{
			Id:         1,
			CardNumber: "1234567890123456",
		},
	}

	mockTransactionClient.EXPECT().
		FindByIdTransaction(gomock.Any(), &pb.FindByIdTransactionRequest{TransactionId: int32(id)}).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/transaction/%d", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransaction
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transaction retrieved successfully", resp.Message)
	assert.Equal(t, int32(1), resp.Data.Id)
}

func TestFindById_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	id := 1
	mockTransactionClient.EXPECT().
		FindByIdTransaction(gomock.Any(), &pb.FindByIdTransactionRequest{TransactionId: int32(id)}).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve transaction data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/transaction/%d", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to retrieve transaction data")
}

func TestFindById_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	invalidID := "abc"
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/transaction/%s", invalidID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(invalidID)

	mockLogger.EXPECT().Debug("Invalid transaction ID", gomock.Any()).Times(1)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Bad Request: Invalid ID")
}

func TestFindByCardNumberTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	cardNumber := "1234567890123456"

	expectedGRPCResponse := &pb.ApiResponseTransactions{
		Status:  "success",
		Message: "Transaction retrieved successfully",
		Data: []*pb.TransactionResponse{
			{
				Id:         1,
				CardNumber: cardNumber,
			},
			{
				Id:         2,
				CardNumber: "1234567890123457",
			},
		},
	}

	mockTransactionClient.EXPECT().
		FindByCardNumberTransaction(gomock.Any(), &pb.FindByCardNumberTransactionRequest{CardNumber: cardNumber}).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transaction?card_number="+cardNumber, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindByCardNumber(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransactions
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transaction retrieved successfully", resp.Message)
	assert.Equal(t, cardNumber, resp.Data[0].CardNumber)
}

func TestFindByCardNumberTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	cardNumber := "1234567890123456"

	mockTransactionClient.EXPECT().
		FindByCardNumberTransaction(gomock.Any(), &pb.FindByCardNumberTransactionRequest{CardNumber: cardNumber}).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve transaction data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transaction?card_number="+cardNumber, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindByCardNumber(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to retrieve transaction data")
}

func TestFindByTransactionMerchantId_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	merchantId := 12345
	expectedGRPCResponse := &pb.ApiResponseTransactions{
		Status:  "success",
		Message: "Transaction retrieved successfully",
		Data: []*pb.TransactionResponse{
			{
				Id:         1,
				CardNumber: "1234567890123456",
			},
			{
				Id:         2,
				CardNumber: "1234567890123457",
			},
		},
	}

	mockTransactionClient.EXPECT().
		FindTransactionByMerchantId(gomock.Any(), &pb.FindTransactionByMerchantIdRequest{MerchantId: int32(merchantId)}).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transaction?merchant_id=12345", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindByTransactionMerchantId(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransactions
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transaction retrieved successfully", resp.Message)
	assert.Equal(t, "1234567890123456", resp.Data[0].CardNumber)
}

func TestFindByTransactionMerchantId_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	merchantId := 12345

	mockTransactionClient.EXPECT().
		FindTransactionByMerchantId(gomock.Any(), &pb.FindTransactionByMerchantIdRequest{MerchantId: int32(merchantId)}).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve transaction data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transaction?merchant_id=12345", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindByTransactionMerchantId(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to retrieve transaction data")
}

func TestFindByTransactionMerchantId_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transaction?merchant_id=invalid_id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindByTransactionMerchantId(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Bad Request: Invalid ID", resp.Message)
}

func TestFindByActiveTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	expectedGRPCResponse := &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:  "success",
		Message: "Transaction retrieved successfully",
		Data: []*pb.TransactionResponseDeleteAt{
			{
				Id:         1,
				CardNumber: "1234567890123456",
			},
			{
				Id:         2,
				CardNumber: "1234567890123457",
			},
		},
	}
	request := &pb.FindAllTransactionRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockTransactionClient.EXPECT().
		FindByActiveTransaction(gomock.Any(), request).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transaction/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindByActiveTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransactions
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transaction retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, int32(1), resp.Data[0].Id)
	assert.Equal(t, "1234567890123456", resp.Data[0].CardNumber)
}

func TestFindByActiveTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	request := &pb.FindAllTransactionRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockTransactionClient.EXPECT().
		FindByActiveTransaction(gomock.Any(), request).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve transaction data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transaction/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindByActiveTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to retrieve transaction data")
}

func TestFindByActiveTransaction_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	request := &pb.FindAllTransactionRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	expectedGRPCResponse := &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:  "success",
		Message: "Transaction retrieved successfully",
		Data:    []*pb.TransactionResponseDeleteAt{},
	}

	mockTransactionClient.EXPECT().
		FindByActiveTransaction(gomock.Any(), request).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transaction/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindByActiveTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransactions
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transaction retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 0)
}

func TestFindByTrashedTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	request := &pb.FindAllTransactionRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	expectedGRPCResponse := &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:  "success",
		Message: "Trashed transactions retrieved successfully",
		Data: []*pb.TransactionResponseDeleteAt{
			{
				Id:         1,
				CardNumber: "1234567890123456",
			},
			{
				Id:         2,
				CardNumber: "1234567890123457",
			},
		},
	}

	mockTransactionClient.EXPECT().
		FindByTrashedTransaction(gomock.Any(), request).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transaction/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindByTrashedTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransactions
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Trashed transactions retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, int32(1), resp.Data[0].Id)
	assert.Equal(t, "1234567890123456", resp.Data[0].CardNumber)
}

func TestFindByTrashedTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	request := &pb.FindAllTransactionRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockTransactionClient.EXPECT().
		FindByTrashedTransaction(gomock.Any(), request).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve transaction data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transaction/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindByTrashedTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to retrieve transaction data")
}

func TestFindByTrashedTransaction_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	expectedGRPCResponse := &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:  "success",
		Message: "Trashed transactions retrieved successfully",
		Data:    []*pb.TransactionResponseDeleteAt{},
	}
	request := &pb.FindAllTransactionRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockTransactionClient.EXPECT().
		FindByTrashedTransaction(gomock.Any(), request).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transaction/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	err := handler.FindByTrashedTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransactions
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Trashed transactions retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 0)
}

func TestCreateTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	body := requests.CreateTransactionRequest{
		CardNumber:      "1234567890123456",
		Amount:          1000000,
		PaymentMethod:   "mandiri",
		MerchantID:      utils.PtrInt(1),
		TransactionTime: time.Now(),
	}

	mockMerchantClient.EXPECT().
		FindByApiKey(gomock.Any(), &pb.FindByApiKeyRequest{
			ApiKey: "test-api-key",
		}).
		Return(&pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant found successfully",
			Data: &pb.MerchantResponse{
				Id:   1,
				Name: "test-merchant",
			},
		}, nil).
		Times(1)

	expectedGRPCRequest := &pb.CreateTransactionRequest{
		CardNumber:      body.CardNumber,
		Amount:          int32(body.Amount),
		PaymentMethod:   body.PaymentMethod,
		MerchantId:      int32(*body.MerchantID),
		TransactionTime: timestamppb.New(body.TransactionTime),
		ApiKey:          "test-api-key",
	}

	expectedGRPCResponse := &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Transaction created successfully",
		Data: &pb.TransactionResponse{
			Id:         1,
			CardNumber: "1234567890123456",
			Amount:     1000,
		},
	}

	mockTransactionClient.EXPECT().
		CreateTransaction(gomock.Any(), expectedGRPCRequest).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/transaction", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "test-api-key")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	middleware := middlewares.ApiKeyMiddleware(mockMerchantClient)
	h := middleware(handler.Create)

	err = h(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransaction
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transaction created successfully", resp.Message)
}

func TestCreateTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	body := requests.CreateTransactionRequest{
		CardNumber:      "1234567890123456",
		Amount:          1000000,
		PaymentMethod:   "mandiri",
		MerchantID:      utils.PtrInt(1),
		TransactionTime: time.Now(),
	}

	mockMerchantClient.EXPECT().
		FindByApiKey(gomock.Any(), &pb.FindByApiKeyRequest{
			ApiKey: "test-api-key",
		}).
		Return(&pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant found successfully",
			Data: &pb.MerchantResponse{
				Id:   1,
				Name: "test-merchant",
			},
		}, nil).
		Times(1)

	expectedGRPCRequest := &pb.CreateTransactionRequest{
		CardNumber:      body.CardNumber,
		Amount:          int32(body.Amount),
		PaymentMethod:   body.PaymentMethod,
		MerchantId:      int32(*body.MerchantID),
		TransactionTime: timestamppb.New(body.TransactionTime),
		ApiKey:          "test-api-key",
	}

	mockTransactionClient.EXPECT().
		CreateTransaction(gomock.Any(), expectedGRPCRequest).
		Return(nil, errors.New("gRPC error")).
		Times(1)

	e := echo.New()
	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/transaction", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "test-api-key")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	middleware := middlewares.ApiKeyMiddleware(mockMerchantClient)
	h := middleware(handler.Create)

	err = h(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to create transaction: ", resp.Message)
}

func TestCreateTransaction_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	mockMerchantClient.EXPECT().
		FindByApiKey(gomock.Any(), &pb.FindByApiKeyRequest{
			ApiKey: "test-api-key",
		}).
		Return(&pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant found successfully",
			Data: &pb.MerchantResponse{
				Id:   1,
				Name: "test-merchant",
			},
		}, nil).
		Times(1)

	body := requests.CreateTransactionRequest{
		CardNumber:    "",
		Amount:        0,
		PaymentMethod: "",
	}

	e := echo.New()
	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/transaction", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "test-api-key")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	middleware := middlewares.ApiKeyMiddleware(mockMerchantClient)
	h := middleware(handler.Create)

	err = h(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error")
}

func TestUpdateTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	body := requests.UpdateTransactionRequest{
		TransactionID:   1,
		CardNumber:      "1234567890123456",
		Amount:          1000000,
		PaymentMethod:   "mandiri",
		MerchantID:      utils.PtrInt(1),
		TransactionTime: time.Now(),
	}

	mockMerchantClient.EXPECT().
		FindByApiKey(gomock.Any(), &pb.FindByApiKeyRequest{
			ApiKey: "test-api-key",
		}).
		Return(&pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant found successfully",
			Data: &pb.MerchantResponse{
				Id:   1,
				Name: "test-merchant",
			},
		}, nil).
		Times(1)

	expectedGRPCRequest := &pb.UpdateTransactionRequest{
		TransactionId:   int32(body.TransactionID),
		CardNumber:      body.CardNumber,
		Amount:          int32(body.Amount),
		PaymentMethod:   body.PaymentMethod,
		MerchantId:      int32(*body.MerchantID),
		TransactionTime: timestamppb.New(body.TransactionTime),
		ApiKey:          "test-api-key",
	}

	expectedGRPCResponse := &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Transaction updated successfully",
		Data: &pb.TransactionResponse{
			Id:         1,
			CardNumber: "1234567890123456",
			Amount:     1000000,
		},
	}

	mockTransactionClient.EXPECT().
		UpdateTransaction(gomock.Any(), expectedGRPCRequest).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	e.POST("/api/transaction/update/:id", func(c echo.Context) error {
		handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)
		return handler.Update(c)
	})

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/transaction/update/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "test-api-key")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/api/transaction/update/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	middleware := middlewares.ApiKeyMiddleware(mockMerchantClient)
	h := middleware(func(c echo.Context) error {
		handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)
		return handler.Update(c)
	})

	err = h(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransaction
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transaction updated successfully", resp.Message)
}

func TestUpdateTransaction_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	mockMerchantClient.EXPECT().
		FindByApiKey(gomock.Any(), &pb.FindByApiKeyRequest{
			ApiKey: "test-api-key",
		}).
		Return(&pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant found successfully",
			Data: &pb.MerchantResponse{
				Id:   1,
				Name: "test-merchant",
			},
		}, nil).
		Times(1)

	body := requests.UpdateTransactionRequest{
		TransactionID:   0,
		CardNumber:      "1234567890123456",
		Amount:          1000000,
		PaymentMethod:   "mandiri",
		MerchantID:      utils.PtrInt(1),
		TransactionTime: time.Now(),
	}

	e := echo.New()
	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/transaction/update/0", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "test-api-key")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)

	middleware := middlewares.ApiKeyMiddleware(mockMerchantClient)
	h := middleware(handler.Update)

	err = h(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Bad Request: Invalid ID")
}

func TestUpdateTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	body := requests.UpdateTransactionRequest{
		TransactionID:   1,
		CardNumber:      "1234567890123456",
		Amount:          1000000,
		PaymentMethod:   "mandiri",
		MerchantID:      utils.PtrInt(1),
		TransactionTime: time.Now(),
	}

	mockMerchantClient.EXPECT().
		FindByApiKey(gomock.Any(), &pb.FindByApiKeyRequest{
			ApiKey: "test-api-key",
		}).
		Return(&pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant found successfully",
			Data: &pb.MerchantResponse{
				Id:   1,
				Name: "test-merchant",
			},
		}, nil).
		Times(1)

	expectedGRPCRequest := &pb.UpdateTransactionRequest{
		TransactionId:   int32(body.TransactionID),
		CardNumber:      body.CardNumber,
		Amount:          int32(body.Amount),
		PaymentMethod:   body.PaymentMethod,
		MerchantId:      int32(*body.MerchantID),
		TransactionTime: timestamppb.New(body.TransactionTime),
		ApiKey:          "test-api-key",
	}

	mockTransactionClient.EXPECT().
		UpdateTransaction(gomock.Any(), expectedGRPCRequest).
		Return(nil, status.Error(codes.Internal, "internal server error")).
		Times(1)

	e := echo.New()
	e.POST("/api/transaction/update/:id", func(c echo.Context) error {
		handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)
		return handler.Update(c)
	})

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/transaction/update/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "test-api-key")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/api/transaction/update/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	middleware := middlewares.ApiKeyMiddleware(mockMerchantClient)
	h := middleware(func(c echo.Context) error {
		handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)
		return handler.Update(c)
	})

	err = h(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to update transaction: ", resp.Message)
}

func TestUpdateTransaction_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	body := requests.UpdateTransactionRequest{
		TransactionID:   1,
		CardNumber:      "",
		Amount:          1000000,
		PaymentMethod:   "mandiri",
		MerchantID:      utils.PtrInt(1),
		TransactionTime: time.Now(),
	}

	mockMerchantClient.EXPECT().FindByApiKey(gomock.Any(), &pb.FindByApiKeyRequest{
		ApiKey: "test-api-key",
	}).Return(&pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Merchant found successfully",
		Data: &pb.MerchantResponse{
			Id:   1,
			Name: "test-merchant",
		},
	}, nil)

	e := echo.New()
	e.POST("/api/transaction/update/:id", func(c echo.Context) error {
		handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)
		return handler.Update(c)
	})

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/transaction/update/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "test-api-key")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/api/transaction/update/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	middleware := middlewares.ApiKeyMiddleware(mockMerchantClient)
	h := middleware(func(c echo.Context) error {
		handler := api.NewHandlerTransaction(mockTransactionClient, mockMerchantClient, e, mockLogger)
		return handler.Update(c)
	})

	err = h(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Validation Error: Key: 'UpdateTransactionRequest.CardNumber' Error:Field validation for 'CardNumber' failed on the 'required' tag", resp.Message)
}

func TestTrashedTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	expectedGRPCResponse := &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Transaction trashed successfully",
		Data: &pb.TransactionResponse{
			Id:         1,
			CardNumber: "1234567890123456",
		},
	}

	mockTransactionClient.EXPECT().
		TrashedTransaction(gomock.Any(), &pb.FindByIdTransactionRequest{
			TransactionId: 1,
		}).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/transaction/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTransaction(mockTransactionClient, nil, e, mockLogger)

	err := handler.TrashedTransaction(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransaction
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transaction trashed successfully", resp.Message)
}

func TestTrashedTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	mockTransactionClient.EXPECT().
		TrashedTransaction(gomock.Any(), &pb.FindByIdTransactionRequest{
			TransactionId: 1,
		}).
		Return(nil, errors.New("gRPC error")).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/transaction/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTransaction(mockTransactionClient, nil, e, mockLogger)

	err := handler.TrashedTransaction(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to trashed transaction:", resp.Message)
}

func TestTrashedTransaction_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/transaction/trashed/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	handler := api.NewHandlerTransaction(mockTransactionClient, nil, e, mockLogger)

	err := handler.TrashedTransaction(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Bad Request: Invalid ID", resp.Message)
}

func TestRestoreTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	expectedGRPCResponse := &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Transaction restored successfully",
		Data: &pb.TransactionResponse{
			Id:         1,
			CardNumber: "1234567890123456",
		},
	}

	mockTransactionClient.EXPECT().
		RestoreTransaction(gomock.Any(), &pb.FindByIdTransactionRequest{
			TransactionId: 1,
		}).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/transaction/restore/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTransaction(mockTransactionClient, nil, e, mockLogger)

	err := handler.RestoreTransaction(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransaction
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transaction restored successfully", resp.Message)
}

func TestRestoreTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	mockTransactionClient.EXPECT().
		RestoreTransaction(gomock.Any(), &pb.FindByIdTransactionRequest{
			TransactionId: 1,
		}).
		Return(nil, errors.New("gRPC error")).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/transaction/restore/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTransaction(mockTransactionClient, nil, e, mockLogger)

	err := handler.RestoreTransaction(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to restore transaction:", resp.Message)
}

func TestRestoreTransaction_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/transaction/restore/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	handler := api.NewHandlerTransaction(mockTransactionClient, nil, e, mockLogger)

	err := handler.RestoreTransaction(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Bad Request: Invalid ID", resp.Message)
}

func TestDeletePermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	expectedGRPCResponse := &pb.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Transaction deleted permanently",
	}

	mockTransactionClient.EXPECT().
		DeleteTransactionPermanent(gomock.Any(), &pb.FindByIdTransactionRequest{
			TransactionId: 1,
		}).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/transaction/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTransaction(mockTransactionClient, nil, e, mockLogger)

	err := handler.DeletePermanent(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseTransactionDelete
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Transaction deleted permanently", resp.Message)
}

func TestDeletePermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	mockTransactionClient.EXPECT().
		DeleteTransactionPermanent(gomock.Any(), &pb.FindByIdTransactionRequest{
			TransactionId: 1,
		}).
		Return(nil, errors.New("gRPC error")).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/transaction/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerTransaction(mockTransactionClient, nil, e, mockLogger)
	err := handler.DeletePermanent(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to delete transaction:", resp.Message)
}

func TestDeletePermanent_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionClient := mock_pb.NewMockTransactionServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/transaction/permanent/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	handler := api.NewHandlerTransaction(mockTransactionClient, nil, e, mockLogger)

	err := handler.DeletePermanent(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Bad Request: Invalid ID", resp.Message)
}
