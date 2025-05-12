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

func TestFindAllSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	page := 1
	pageSize := 10
	search := "example"

	expectedResponse := &pb.ApiResponsePaginationSaldo{
		Status:  "success",
		Message: "Saldo data retrieved successfully",
		Data: []*pb.SaldoResponse{
			{
				SaldoId:      1,
				TotalBalance: 10000,
			},
			{
				SaldoId:      2,
				TotalBalance: 20000,
			},
		},
	}

	mockSaldoClient.EXPECT().
		FindAllSaldo(
			gomock.Any(),
			&pb.FindAllSaldoRequest{
				Page:     int32(page),
				PageSize: int32(pageSize),
				Search:   search,
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo?page=1&page_size=10&search=example", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Saldo data retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, int32(1), resp.Data[0].SaldoId)
	assert.Equal(t, int32(10000), resp.Data[0].TotalBalance)
}

func TestFindAllSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	page := 1
	pageSize := 10
	search := "example"

	mockSaldoClient.EXPECT().
		FindAllSaldo(
			gomock.Any(),
			&pb.FindAllSaldoRequest{
				Page:     int32(page),
				PageSize: int32(pageSize),
				Search:   search,
			},
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve saldo data: ",
		})

	mockLogger.EXPECT().Debug("Failed to retrieve saldo data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo?page=1&page_size=10&search=example", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve saldo data: ", resp.Message)
}

func TestFindAllSaldo_EmptyResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	page := 1
	pageSize := 10
	search := "example"

	expectedResponse := &pb.ApiResponsePaginationSaldo{
		Status:  "success",
		Message: "Saldo data retrieved successfully",
		Data:    []*pb.SaldoResponse{},
	}

	mockSaldoClient.EXPECT().
		FindAllSaldo(
			gomock.Any(),
			&pb.FindAllSaldoRequest{
				Page:     int32(page),
				PageSize: int32(pageSize),
				Search:   search,
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo?page=1&page_size=10&search=example", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Saldo data retrieved successfully", resp.Message)
	assert.Empty(t, resp.Data)
}

func TestFindByIdSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	id := 1

	expectedResponse := &pb.SaldoResponse{
		SaldoId:      1,
		TotalBalance: 10000,
	}
	expect := &pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Saldo data retrieved successfully",
		Data:    expectedResponse,
	}

	mockSaldoClient.EXPECT().
		FindByIdSaldo(
			gomock.Any(),
			&pb.FindByIdSaldoRequest{SaldoId: int32(id)},
		).
		Return(expect, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, int32(1), resp.Data.SaldoId)
	assert.Equal(t, int32(10000), resp.Data.TotalBalance)
}

func TestFindByIdSaldo_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	invalidID := "abc"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/saldo/%s", invalidID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(invalidID)

	mockLogger.EXPECT().Debug("Invalid saldo ID", gomock.Any()).Times(1)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid saldo ID", resp.Message)
}

func TestFindByIdSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	id := 1

	mockSaldoClient.EXPECT().
		FindByIdSaldo(
			gomock.Any(),
			&pb.FindByIdSaldoRequest{SaldoId: int32(id)},
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve saldo data: ",
		})

	mockLogger.EXPECT().Debug("Failed to retrieve saldo data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve saldo data: ", resp.Message)
}

func TestFindByCardNumberSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	cardNumber := "1234567890"
	expectedResponse := &pb.SaldoResponse{
		SaldoId:      1,
		TotalBalance: 10000,
	}

	expect := &pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Saldo data retrieved successfully",
		Data:    expectedResponse,
	}

	mockSaldoClient.EXPECT().
		FindByCardNumber(
			gomock.Any(),
			&pb.FindByCardNumberRequest{CardNumber: cardNumber},
		).
		Return(expect, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo/card-number/1234567890", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindByCardNumber(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, int32(1), resp.Data.SaldoId)
	assert.Equal(t, int32(10000), resp.Data.TotalBalance)
}

func TestFindByCardNumber_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	cardNumber := "1234567890"

	mockSaldoClient.EXPECT().
		FindByCardNumber(
			gomock.Any(),
			&pb.FindByCardNumberRequest{CardNumber: cardNumber},
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve saldo data: ",
		})

	mockLogger.EXPECT().Debug("Failed to retrieve saldo data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo/card-number/1234567890", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindByCardNumber(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve saldo data: ", resp.Message)
}

func TestFindByActiveSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationSaldoDeleteAt{
		Status:  "success",
		Message: "Saldo data retrieved successfully",
		Data: []*pb.SaldoResponseDeleteAt{
			{
				SaldoId:      1,
				TotalBalance: 10000,
			},
			{
				SaldoId:      2,
				TotalBalance: 20000,
			},
		},
	}

	request := &pb.FindAllSaldoRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockSaldoClient.EXPECT().
		FindByActive(gomock.Any(), request).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindByActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Saldo data retrieved successfully", resp.Message)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, int32(1), resp.Data[0].SaldoId)
	assert.Equal(t, int32(10000), resp.Data[0].TotalBalance)
}

func TestFindByActiveSaldo_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllSaldoRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	expectedResponse := &pb.ApiResponsePaginationSaldoDeleteAt{
		Status:  "success",
		Message: "No active saldo data found",
		Data:    []*pb.SaldoResponseDeleteAt{},
	}

	mockSaldoClient.EXPECT().
		FindByActive(gomock.Any(), request).
		Return(expectedResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindByActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No active saldo data found", resp.Message)
	assert.Empty(t, resp.Data)
}

func TestFindByActiveSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllSaldoRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockSaldoClient.EXPECT().
		FindByActive(gomock.Any(), request).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve saldo data: ",
		})

	mockLogger.EXPECT().Debug("Failed to retrieve saldo data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindByActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve saldo data: ", resp.Message)
}

func TestFindByTrashedSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationSaldoDeleteAt{
		Status:  "success",
		Message: "Trashed saldo data retrieved successfully",
		Data: []*pb.SaldoResponseDeleteAt{
			{
				SaldoId:      1,
				TotalBalance: 10000,
			},
			{
				SaldoId:      2,
				TotalBalance: 20000,
			},
		},
	}

	request := &pb.FindAllSaldoRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockSaldoClient.EXPECT().
		FindByTrashed(gomock.Any(), request).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindByTrashed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Trashed saldo data retrieved successfully", resp.Message)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, int32(1), resp.Data[0].SaldoId)
	assert.Equal(t, int32(10000), resp.Data[0].TotalBalance)
}

func TestFindByTrashedSaldo_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllSaldoRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	expectedResponse := &pb.ApiResponsePaginationSaldoDeleteAt{
		Status:  "success",
		Message: "No trashed saldo data found",
		Data:    []*pb.SaldoResponseDeleteAt{},
	}

	mockSaldoClient.EXPECT().
		FindByTrashed(gomock.Any(), request).
		Return(expectedResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindByTrashed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No trashed saldo data found", resp.Message)
	assert.Empty(t, resp.Data)
}

func TestFindByTrashedSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllSaldoRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockSaldoClient.EXPECT().
		FindByTrashed(gomock.Any(), request).
		Return(nil, fmt.Errorf("internal server error"))

	mockLogger.EXPECT().Debug("Failed to retrieve saldo data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/saldo/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.FindByTrashed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve saldo data: ", resp.Message)
}

func TestCreateSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.CreateSaldoRequest{
		CardNumber:   "1234567890",
		TotalBalance: 5000,
	}

	expectedResponse := &pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Saldo created successfully",
		Data: &pb.SaldoResponse{
			SaldoId:      1,
			CardNumber:   "1234567890",
			TotalBalance: 5000,
		},
	}

	mockSaldoClient.EXPECT().
		CreateSaldo(
			gomock.Any(),
			&pb.CreateSaldoRequest{
				CardNumber:   requestBody.CardNumber,
				TotalBalance: int32(requestBody.TotalBalance),
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	bodyJSON, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/saldo", bytes.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Saldo created successfully", resp.Message)
	assert.Equal(t, "1234567890", resp.Data.CardNumber)
	assert.Equal(t, int32(5000), resp.Data.TotalBalance)
}

func TestCreateSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.CreateSaldoRequest{
		CardNumber:   "1234567890",
		TotalBalance: 5000,
	}

	mockSaldoClient.EXPECT().
		CreateSaldo(
			gomock.Any(),
			&pb.CreateSaldoRequest{
				CardNumber:   requestBody.CardNumber,
				TotalBalance: int32(requestBody.TotalBalance),
			},
		).
		Return(nil, fmt.Errorf("internal server error"))

	mockLogger.EXPECT().Debug("Failed to create saldo", gomock.Any()).Times(1)

	e := echo.New()
	bodyJSON, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/saldo", bytes.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to create saldo: ", resp.Message)
}

func TestCreateSaldo_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	invalidRequestBody := requests.CreateSaldoRequest{
		CardNumber:   "",
		TotalBalance: -5000,
	}

	mockLogger.EXPECT().Debug("Validation Error", gomock.Any()).Times(1)

	e := echo.New()
	bodyJSON, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/saldo", bytes.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error")
}

func TestUpdateSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.UpdateSaldoRequest{
		SaldoID:      1,
		CardNumber:   "1234567890",
		TotalBalance: 10000,
	}

	expectedResponse := &pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Saldo updated successfully",
		Data: &pb.SaldoResponse{
			SaldoId:      1,
			CardNumber:   "1234567890",
			TotalBalance: 10000,
		},
	}

	mockSaldoClient.EXPECT().
		UpdateSaldo(
			gomock.Any(),
			&pb.UpdateSaldoRequest{
				SaldoId:      int32(requestBody.SaldoID),
				CardNumber:   requestBody.CardNumber,
				TotalBalance: int32(requestBody.TotalBalance),
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	bodyJSON, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/api/saldo", bytes.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("id", int32(requestBody.SaldoID))

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Saldo updated successfully", resp.Message)
	assert.Equal(t, requestBody.SaldoID, int(resp.Data.SaldoId))
	assert.Equal(t, requestBody.CardNumber, resp.Data.CardNumber)
	assert.Equal(t, requestBody.TotalBalance, int(resp.Data.TotalBalance))
}

func TestUpdateSaldo_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/saldo/update/ab", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("ab")

	mockLogger.EXPECT().Debug("Invalid saldo ID", gomock.Any()).Times(1)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid saldo ID", resp.Message)
}

func TestUpdateSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.UpdateSaldoRequest{
		SaldoID:      1,
		CardNumber:   "1234567890",
		TotalBalance: 10000,
	}

	mockLogger.EXPECT().Debug("Failed to update saldo", gomock.Any()).Times(1)

	mockSaldoClient.EXPECT().
		UpdateSaldo(
			gomock.Any(),
			&pb.UpdateSaldoRequest{
				SaldoId:      int32(requestBody.SaldoID),
				CardNumber:   requestBody.CardNumber,
				TotalBalance: int32(requestBody.TotalBalance),
			},
		).
		Return(nil, errors.New("internal server error"))

	e := echo.New()
	bodyJSON, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/api/saldo", bytes.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("id", int32(requestBody.SaldoID))

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to update saldo: ", resp.Message)
}

func TestUpdateSaldo_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.UpdateSaldoRequest{
		SaldoID:      1,
		CardNumber:   "",
		TotalBalance: 10000,
	}

	mockLogger.EXPECT().Debug("Validation Error", gomock.Any()).Times(1)

	e := echo.New()
	bodyJSON, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/api/saldo", bytes.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("id", int32(requestBody.SaldoID))

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Validation Error: ", resp.Message)
}

func TestTrashSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Saldo trashed successfully",
		Data: &pb.SaldoResponse{
			SaldoId:      1,
			CardNumber:   "1234567890",
			TotalBalance: 10000,
		},
	}

	mockSaldoClient.EXPECT().
		TrashedSaldo(
			gomock.Any(),
			&pb.FindByIdSaldoRequest{
				SaldoId: 1,
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/saldo/trash/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.TrashSaldo(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Saldo trashed successfully", resp.Message)
	assert.Equal(t, int32(1), resp.Data.SaldoId)
	assert.Equal(t, "1234567890", resp.Data.CardNumber)
	assert.Equal(t, int32(10000), resp.Data.TotalBalance)
}

func TestTrashSaldo_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/saldo/trash/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	mockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.TrashSaldo(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Bad Request: Invalid ID")
}

func TestTrashSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockSaldoClient.EXPECT().
		TrashedSaldo(
			gomock.Any(),
			&pb.FindByIdSaldoRequest{
				SaldoId: 1,
			},
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to trashed saldo",
		})

	mockLogger.EXPECT().Debug("Failed to trashed saldo", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/saldo/trash/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.TrashSaldo(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to trashed saldo")
}

func TestRestoreSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Saldo restored successfully",
		Data: &pb.SaldoResponse{
			SaldoId:      1,
			CardNumber:   "1234567890",
			TotalBalance: 10000,
		},
	}

	mockSaldoClient.EXPECT().
		RestoreSaldo(
			gomock.Any(),
			&pb.FindByIdSaldoRequest{
				SaldoId: 1,
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/saldo/restore/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.RestoreSaldo(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Saldo restored successfully", resp.Message)
	assert.Equal(t, int32(1), resp.Data.SaldoId)
	assert.Equal(t, "1234567890", resp.Data.CardNumber)
	assert.Equal(t, int32(10000), resp.Data.TotalBalance)
}

func TestRestoreSaldo_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/saldo/restore/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	mockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.RestoreSaldo(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Bad Request: Invalid ID")
}

func TestRestoreSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockSaldoClient.EXPECT().
		RestoreSaldo(
			gomock.Any(),
			&pb.FindByIdSaldoRequest{
				SaldoId: 1,
			},
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore saldo",
		})

	mockLogger.EXPECT().Debug("Failed to restore saldo", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/saldo/restore/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.RestoreSaldo(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to restore saldo")
}

func TestDeleteSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponseSaldoDelete{
		Status:  "success",
		Message: "Saldo deleted successfully",
	}

	mockSaldoClient.EXPECT().
		DeleteSaldoPermanent(
			gomock.Any(),
			&pb.FindByIdSaldoRequest{
				SaldoId: 1,
			},
		).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/saldo/delete/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.Delete(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Saldo deleted successfully", resp.Message)
}

func TestDeleteSaldo_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/saldo/delete/ss", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("ss")

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

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

func TestDeleteSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoClient := mock_pb.NewMockSaldoServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockSaldoClient.EXPECT().
		DeleteSaldoPermanent(
			gomock.Any(),
			&pb.FindByIdSaldoRequest{
				SaldoId: 1,
			},
		).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete saldo",
		})

	mockLogger.EXPECT().Debug("Failed to delete saldo", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/saldo/delete/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerSaldo(mockSaldoClient, e, mockLogger)

	err := handler.Delete(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to delete saldo")
}
