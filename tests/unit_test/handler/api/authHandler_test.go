package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/handler/api"
	mock_apimapper "MamangRust/paymentgatewaygrpc/internal/mapper/response/api/mocks"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	mock_pb "MamangRust/paymentgatewaygrpc/internal/pb/mocks"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestHandleHello(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockClient := mock_pb.NewMockAuthServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockAuthResponseMapper(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/auth/hello", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerAuth(mockClient, e, mockLogger, mockMapper)

	err := handler.HandleHello(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Hello", rec.Body.String())
}

func TestHandleRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_pb.NewMockAuthServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockAuthResponseMapper(ctrl)

	requestBody := requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	expectedRequest := &pb.RegisterRequest{
		Firstname:       "John",
		Lastname:        "Doe",
		Email:           "john@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	expectedResponse := &pb.ApiResponseRegister{
		Status:  "success",
		Message: "User registered successfully",
	}

	mockClient.EXPECT().
		RegisterUser(context.Background(), expectedRequest).
		Return(expectedResponse, nil)

	e := echo.New()
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerAuth(mockClient, e, mockLogger, mockMapper)
	err := handler.Register(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response pb.ApiResponseRegister
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Status, response.Status)
	assert.Equal(t, expectedResponse.Message, response.Message)
}

func TestHandleRegister_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_pb.NewMockAuthServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockAuthResponseMapper(ctrl)

	requestBody := requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	expectedRequest := &pb.RegisterRequest{
		Firstname:       "John",
		Lastname:        "Doe",
		Email:           "john@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	expectedResponse := &pb.ApiResponseRegister{
		Status:  "error",
		Message: "User registration failed",
	}

	mockClient.EXPECT().
		RegisterUser(context.Background(), expectedRequest).
		Return(expectedResponse, nil)

	e := echo.New()
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerAuth(mockClient, e, mockLogger, mockMapper)
	err := handler.Register(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response pb.ApiResponseRegister
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Status, response.Status)
	assert.Equal(t, expectedResponse.Message, response.Message)
}

func TestHandleRegister_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_pb.NewMockAuthServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockAuthResponseMapper(ctrl)

	mockLogger.EXPECT().Debug(
		"Validation Error",
		gomock.AssignableToTypeOf(zap.Field{}),
	).Times(1)

	requestBody := `{}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer([]byte(requestBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerAuth(mockClient, e, mockLogger, mockMapper)

	err := handler.Register(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error")
}

func TestHandleLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_pb.NewMockAuthServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockAuthResponseMapper(ctrl)

	mockResponse := &pb.ApiResponseLogin{
		Status:  "success",
		Message: "Login successful",
		Data: &pb.TokenResponse{
			AccessToken:  "mockToken123",
			RefreshToken: "mockRefreshToken123",
		},
	}
	mockClient.EXPECT().LoginUser(context.Background(), gomock.Any()).Return(mockResponse, nil)

	requestBody := `{"email":"test@example.com","password":"password123"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer([]byte(requestBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerAuth(mockClient, e, mockLogger, mockMapper)

	err := handler.Login(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response pb.ApiResponseLogin

	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Login successful", response.Message)

	assert.NotNil(t, response.Data)
	assert.Equal(t, "mockToken123", response.Data.AccessToken)
	assert.Equal(t, "mockRefreshToken123", response.Data.RefreshToken)
}

func TestHandleLogin_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_pb.NewMockAuthServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockAuthResponseMapper(ctrl)

	mockClient.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Internal Server Error: ",
	})
	mockLogger.EXPECT().Debug("Failed to login user", gomock.Any()).Times(1)

	requestBody := `{"email":"test@example.com","password":"password123"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer([]byte(requestBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerAuth(mockClient, e, mockLogger, mockMapper)

	err := handler.Login(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Internal Server Error: ", response.Message)
}

func TestHandleLogin_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_pb.NewMockAuthServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockAuthResponseMapper(ctrl)

	mockLogger.EXPECT().Debug(gomock.Eq("Validation Error"), gomock.Any()).Times(1)

	requestBody := `{}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer([]byte(requestBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerAuth(mockClient, e, mockLogger, mockMapper)

	err := handler.Login(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Message, "Validation Error")
}
