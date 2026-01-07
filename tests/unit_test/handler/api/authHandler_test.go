package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/handler/api"
	mock_apimapper "MamangRust/paymentgatewaygrpc/internal/mapper/response/api/mocks"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	mock_pb "MamangRust/paymentgatewaygrpc/internal/pb/mocks"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	Ctrl           *gomock.Controller
	MockAuthClient *mock_pb.MockAuthServiceClient
	MockLogger     *mock_logger.MockLoggerInterface
	MockMapper     *mock_apimapper.MockAuthResponseMapper
	E              *echo.Echo
	Handler        *api.AuthHandleApi
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockAuthClient = mock_pb.NewMockAuthServiceClient(suite.Ctrl)
	suite.MockLogger = mock_logger.NewMockLoggerInterface(suite.Ctrl)
	suite.MockMapper = mock_apimapper.NewMockAuthResponseMapper(suite.Ctrl)
	suite.E = echo.New()
	suite.Handler = api.NewHandlerAuth(suite.MockAuthClient, suite.E, suite.MockLogger, suite.MockMapper)
}

func (suite *AuthHandlerTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *AuthHandlerTestSuite) TestRegister_Success() {
	requestBody := requests.RegisterRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	grpcRequest := &pb.RegisterRequest{
		Firstname:       requestBody.FirstName,
		Lastname:        requestBody.LastName,
		Email:           requestBody.Email,
		Password:        requestBody.Password,
		ConfirmPassword: requestBody.ConfirmPassword,
	}

	grpcResponse := &pb.ApiResponseRegister{
		Status:  "success",
		Message: "User registered successfully",
		Data: &pb.UserResponse{
			Id:        1,
			Firstname: requestBody.FirstName,
			Lastname:  requestBody.LastName,
			Email:     requestBody.Email,
		},
	}

	expectedApiResponse := &response.ApiResponseRegister{
		Status:  "success",
		Message: "User registered successfully",
		Data: &response.UserResponse{
			ID:        1,
			FirstName: requestBody.FirstName,
			LastName:  requestBody.LastName,
			Email:     requestBody.Email,
		},
	}

	suite.MockAuthClient.EXPECT().RegisterUser(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToResponseRegister(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Register(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseRegister
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
	suite.Equal(1, resp.Data.ID)
}

func (suite *AuthHandlerTestSuite) TestRegister_Failure() {
	requestBody := requests.RegisterRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	grpcError := errors.New("gRPC service unavailable")
	suite.MockAuthClient.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Registration failed", zap.Error(grpcError)).Times(1)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Register(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *AuthHandlerTestSuite) TestRegister_ValidationError() {
	invalidRequestBody := requests.RegisterRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "invalid-email",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	suite.MockLogger.EXPECT().
		Debug("Validation failed", gomock.Any()).
		Times(1)

	err := suite.Handler.Register(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *AuthHandlerTestSuite) TestLogin_Success() {
	requestBody := requests.AuthRequest{
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	grpcRequest := &pb.LoginRequest{
		Email:    requestBody.Email,
		Password: requestBody.Password,
	}

	grpcResponse := &pb.ApiResponseLogin{
		Status: "success",
		Data: &pb.TokenResponse{
			AccessToken:  "access-token-123",
			RefreshToken: "refresh-token-456",
		},
	}

	expectedApiResponse := &response.ApiResponseLogin{
		Status: "success",
		Data: &response.TokenResponse{
			AccessToken:  "access-token-123",
			RefreshToken: "refresh-token-456",
		},
	}

	suite.MockAuthClient.EXPECT().LoginUser(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToResponseLogin(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Login(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseLogin
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
	suite.Equal("access-token-123", resp.Data.AccessToken)
}

func (suite *AuthHandlerTestSuite) TestLogin_Failure() {
	requestBody := requests.AuthRequest{
		Email:    "john.doe@example.com",
		Password: "password123",
	}
	grpcError := errors.New("gRPC service unavailable")
	suite.MockAuthClient.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Login failed", zap.Error(grpcError)).Times(1)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Login(c)
	suite.NoError(err)
	suite.Equal(http.StatusUnauthorized, rec.Code)
}

func (suite *AuthHandlerTestSuite) TestLogin_ValidationError() {
	invalidRequestBody := requests.AuthRequest{
		Email:    "invalid-email",
		Password: "password123",
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	suite.MockLogger.EXPECT().
		Debug("Validation failed", gomock.Any()).
		Times(1)

	err := suite.Handler.Login(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *AuthHandlerTestSuite) TestRefreshToken_Success() {
	refreshToken := "refresh-token-456"

	grpcRequest := &pb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	grpcResponse := &pb.ApiResponseRefreshToken{
		Status: "success",
		Data: &pb.TokenResponse{
			AccessToken:  "new-access-token-789",
			RefreshToken: refreshToken,
		},
	}

	expectedApiResponse := &response.ApiResponseRefreshToken{
		Status: "success",
		Data: &response.TokenResponse{
			AccessToken:  "new-access-token-789",
			RefreshToken: refreshToken,
		},
	}

	suite.MockAuthClient.EXPECT().RefreshToken(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToResponseRefreshToken(grpcResponse).Return(expectedApiResponse)

	requestBody := map[string]string{
		"refresh_token": refreshToken,
	}
	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh-token", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RefreshToken(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseRefreshToken
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
	suite.Equal("new-access-token-789", resp.Data.AccessToken)
}

func (suite *AuthHandlerTestSuite) TestRefreshToken_Failure() {
	refreshToken := "refresh-token-456"
	grpcError := errors.New("gRPC service unavailable")
	suite.MockAuthClient.EXPECT().RefreshToken(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Token refresh failed", zap.Error(grpcError)).Times(1)

	requestBody := map[string]string{
		"refresh_token": refreshToken,
	}
	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh-token", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RefreshToken(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *AuthHandlerTestSuite) TestRefreshToken_ValidationError() {
	requestBody := map[string]string{}
	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh-token", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	suite.MockLogger.EXPECT().
		Debug("Validation failed", gomock.Any()).
		Times(1)

	err := suite.Handler.RefreshToken(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *AuthHandlerTestSuite) TestGetMe_Success() {
	userID := 1
	token := "access-token-123"

	grpcRequest := &pb.GetMeRequest{
		AccessToken: token,
	}

	grpcResponse := &pb.ApiResponseGetMe{
		Status: "success",
		Data: &pb.UserResponse{
			Id:        int32(userID),
			Firstname: "John",
			Lastname:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	expectedApiResponse := &response.ApiResponseGetMe{
		Status: "success",
		Data: &response.UserResponse{
			ID:        userID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	suite.MockAuthClient.EXPECT().GetMe(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToResponseGetMe(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.GetMe(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseGetMe
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
	suite.Equal(userID, resp.Data.ID)
}

func (suite *AuthHandlerTestSuite) TestGetMe_Failure() {
	token := "access-token-123"
	grpcError := errors.New("gRPC service unavailable")
	suite.MockAuthClient.EXPECT().GetMe(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to get user information", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.GetMe(c)
	suite.NoError(err)
	suite.Equal(http.StatusUnauthorized, rec.Code)
}

func (suite *AuthHandlerTestSuite) TestGetMe_MissingToken() {
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	suite.MockLogger.EXPECT().
		Debug("Authorization header is missing or invalid format", gomock.Any()).
		Times(1)

	err := suite.Handler.GetMe(c)

	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}
