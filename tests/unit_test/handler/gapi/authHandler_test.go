package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/handler/gapi"
	mock_protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto/mocks"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	mock_service "MamangRust/paymentgatewaygrpc/internal/service/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandleGrpcTestSuite struct {
	suite.Suite
	Ctrl            *gomock.Controller
	MockAuthService *mock_service.MockAuthService
	MockProtoMapper *mock_protomapper.MockAuthProtoMapper
	Handler         gapi.AuthHandleGrpc
}

func (suite *AuthHandleGrpcTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockAuthService = mock_service.NewMockAuthService(suite.Ctrl)
	suite.MockProtoMapper = mock_protomapper.NewMockAuthProtoMapper(suite.Ctrl)
	suite.Handler = gapi.NewAuthHandleGrpc(suite.MockAuthService, suite.MockProtoMapper)
}

func (suite *AuthHandleGrpcTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *AuthHandleGrpcTestSuite) TestLoginUser_Success() {
	req := &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedServiceReq := &requests.AuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockServiceResponse := &response.TokenResponse{
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
	}

	expectedGrpcResponse := &pb.ApiResponseLogin{
		Status:  "success",
		Message: "Login successfully",
		Data: &pb.TokenResponse{
			AccessToken:  mockServiceResponse.AccessToken,
			RefreshToken: mockServiceResponse.RefreshToken,
		},
	}

	suite.MockAuthService.EXPECT().Login(gomock.Eq(expectedServiceReq)).Return(mockServiceResponse, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseLogin("success", "Login successfully", mockServiceResponse).Return(expectedGrpcResponse)

	resp, err := suite.Handler.LoginUser(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(expectedGrpcResponse.Status, resp.GetStatus())
	suite.Equal(expectedGrpcResponse.Message, resp.GetMessage())
	suite.Equal(expectedGrpcResponse.Data.AccessToken, resp.GetData().GetAccessToken())
	suite.Equal(expectedGrpcResponse.Data.RefreshToken, resp.GetData().GetRefreshToken())
}

func (suite *AuthHandleGrpcTestSuite) TestLoginUser_Failure_InvalidCredentials() {
	req := &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	expectedServiceReq := &requests.AuthRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	serviceError := &response.ErrorResponse{
		Status:  "error",
		Message: "invalid credentials",
		Code:    int(codes.Unauthenticated),
	}

	suite.MockAuthService.EXPECT().Login(gomock.Eq(expectedServiceReq)).Return(nil, serviceError)

	resp, err := suite.Handler.LoginUser(context.Background(), req)

	suite.Error(err)
	suite.Nil(resp)

	statusErr, ok := status.FromError(err)
	suite.True(ok, "Error harus berupa gRPC status error")
	suite.Equal(codes.Unauthenticated, statusErr.Code())
	suite.Contains(statusErr.Message(), serviceError.Message)
}

func (suite *AuthHandleGrpcTestSuite) TestRegisterUser_Success() {
	req := &pb.RegisterRequest{
		Firstname:       "John",
		Lastname:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	expectedServiceReq := &requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockServiceResponse := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	expectedGrpcResponse := &pb.ApiResponseRegister{
		Status:  "success",
		Message: "Registration successfully",
		Data: &pb.UserResponse{
			Id:        int32(mockServiceResponse.ID),
			Firstname: mockServiceResponse.FirstName,
			Lastname:  mockServiceResponse.LastName,
			Email:     mockServiceResponse.Email,
		},
	}

	suite.MockAuthService.EXPECT().Register(gomock.Eq(expectedServiceReq)).Return(mockServiceResponse, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseRegister("success", "Registration successfully", mockServiceResponse).Return(expectedGrpcResponse)

	resp, err := suite.Handler.RegisterUser(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(expectedGrpcResponse.Status, resp.GetStatus())
	suite.Equal(expectedGrpcResponse.Data.Id, resp.GetData().GetId())
}

func (suite *AuthHandleGrpcTestSuite) TestRegisterUser_Failure_EmailExists() {
	req := &pb.RegisterRequest{
		Firstname: "Jane",
		Lastname:  "Doe",
		Email:     "existing@example.com",
		Password:  "password123",
	}

	serviceError := &response.ErrorResponse{
		Status:  "error",
		Message: "email already exists",
		Code:    int(codes.AlreadyExists),
	}

	suite.MockAuthService.EXPECT().Register(gomock.Any()).Return(nil, serviceError)

	resp, err := suite.Handler.RegisterUser(context.Background(), req)

	suite.Error(err)
	suite.Nil(resp)

	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.AlreadyExists, statusErr.Code())
	suite.Contains(statusErr.Message(), serviceError.Message)
}

func (suite *AuthHandleGrpcTestSuite) TestRefreshToken_Success() {
	req := &pb.RefreshTokenRequest{
		RefreshToken: "valid_refresh_token",
	}

	mockServiceResponse := &response.TokenResponse{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
	}

	expectedGrpcResponse := &pb.ApiResponseRefreshToken{
		Status:  "success",
		Message: "Refresh token successfully",
		Data: &pb.TokenResponse{
			AccessToken:  mockServiceResponse.AccessToken,
			RefreshToken: mockServiceResponse.RefreshToken,
		},
	}

	suite.MockAuthService.EXPECT().RefreshToken(req.RefreshToken).Return(mockServiceResponse, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseRefreshToken("success", "Refresh token successfully", mockServiceResponse).Return(expectedGrpcResponse)

	resp, err := suite.Handler.RefreshToken(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(expectedGrpcResponse.Status, resp.GetStatus())
}

func (suite *AuthHandleGrpcTestSuite) TestRefreshToken_Failure_InvalidToken() {
	req := &pb.RefreshTokenRequest{
		RefreshToken: "invalid_token",
	}

	serviceError := &response.ErrorResponse{
		Status:  "error",
		Message: "invalid or expired refresh token",
		Code:    int(codes.Unauthenticated),
	}

	suite.MockAuthService.EXPECT().RefreshToken(req.RefreshToken).Return(nil, serviceError)

	resp, err := suite.Handler.RefreshToken(context.Background(), req)

	suite.Error(err)
	suite.Nil(resp)

	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Unauthenticated, statusErr.Code())
	suite.Contains(statusErr.Message(), serviceError.Message)
}

func (suite *AuthHandleGrpcTestSuite) TestGetMe_Success() {
	req := &pb.GetMeRequest{
		AccessToken: "valid_access_token",
	}

	mockServiceResponse := &response.UserResponse{
		ID:        1,
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
	}

	expectedGrpcResponse := &pb.ApiResponseGetMe{
		Status:  "success",
		Message: "Refresh token successfully",
		Data: &pb.UserResponse{
			Id:        int32(mockServiceResponse.ID),
			Firstname: mockServiceResponse.FirstName,
			Lastname:  mockServiceResponse.LastName,
			Email:     mockServiceResponse.Email,
		},
	}

	suite.MockAuthService.EXPECT().GetMe(req.AccessToken).Return(mockServiceResponse, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseGetMe("success", "Refresh token successfully", mockServiceResponse).Return(expectedGrpcResponse)

	resp, err := suite.Handler.GetMe(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(expectedGrpcResponse.Status, resp.GetStatus())
	suite.Equal(expectedGrpcResponse.Data.Id, resp.GetData().GetId())
}

func (suite *AuthHandleGrpcTestSuite) TestGetMe_Failure_Unauthorized() {
	req := &pb.GetMeRequest{
		AccessToken: "invalid_access_token",
	}

	serviceError := &response.ErrorResponse{
		Status:  "error",
		Message: "invalid token",
		Code:    int(codes.Unauthenticated),
	}

	suite.MockAuthService.EXPECT().GetMe(req.AccessToken).Return(nil, serviceError)

	resp, err := suite.Handler.GetMe(context.Background(), req)

	suite.Error(err)
	suite.Nil(resp)

	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Unauthenticated, statusErr.Code())
	suite.Contains(statusErr.Message(), serviceError.Message)
}

func TestAuthHandleGrpcSuite(t *testing.T) {
	suite.Run(t, new(AuthHandleGrpcTestSuite))
}
