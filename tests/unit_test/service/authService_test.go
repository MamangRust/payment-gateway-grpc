package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	mock_responsemapper "MamangRust/paymentgatewaygrpc/internal/mapper/response/mocks"
	mock_repository "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/internal/service"
	mock_auth "MamangRust/paymentgatewaygrpc/pkg/auth/mocks"
	mock_hash "MamangRust/paymentgatewaygrpc/pkg/hash/mocks"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRefreshToken := mock_repository.NewMockRefreshTokenRepository(ctrl)
	mockUserRole := mock_repository.NewMockUserRoleRepository(ctrl)
	mockRole := mock_repository.NewMockRoleRepository(ctrl)

	mockHash := mock_hash.NewMockHashPassword(ctrl)
	mockToken := mock_auth.NewMockTokenManager(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	authService := service.NewAuthService(mockUserRepo, mockRefreshToken,
		mockRole,
		mockUserRole,
		mockHash, mockToken, mockLogger, mockMapping)

	request := &requests.CreateUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	hashedPassword := "hashed_password123"
	mockHash.EXPECT().HashPassword(request.Password).Return(hashedPassword, nil)

	mockUserRepo.EXPECT().FindByEmail(request.Email).Return(nil, fmt.Errorf("user not found"))

	hashedRequest := *request
	hashedRequest.Password = hashedPassword

	mockUserRepo.EXPECT().CreateUser(&hashedRequest).Return(&record.UserRecord{
		ID:        1,
		Email:     request.Email,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Password:  hashedPassword,
	}, nil)

	mockRole.EXPECT().FindByName("user").Return(&record.RoleRecord{
		Name: "user",
	}, nil)

	mockUserRole.EXPECT().AssignRoleToUser(&requests.CreateUserRoleRequest{
		UserId: 1,
		RoleId: 1,
	})

	expectedResponse := &response.UserResponse{
		ID:        1,
		Email:     request.Email,
		FirstName: request.FirstName,
		LastName:  request.LastName,
	}
	mockMapping.EXPECT().ToUserResponse(gomock.Any()).Return(expectedResponse)

	mockLogger.EXPECT().Debug("User registered successfully", gomock.Any())

	res, err := authService.Register(request)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expectedResponse, res)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRefreshToken := mock_repository.NewMockRefreshTokenRepository(ctrl)
	mockUserRole := mock_repository.NewMockUserRoleRepository(ctrl)
	mockRole := mock_repository.NewMockRoleRepository(ctrl)

	mockHash := mock_hash.NewMockHashPassword(ctrl)
	mockToken := mock_auth.NewMockTokenManager(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	authService := service.NewAuthService(mockUserRepo, mockRefreshToken,
		mockRole,
		mockUserRole,
		mockHash, mockToken, mockLogger, mockMapping)

	request := &requests.CreateUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	mockUserRepo.EXPECT().FindByEmail(request.Email).Return(&record.UserRecord{
		ID:        1,
		Email:     request.Email,
		FirstName: "Existing",
		LastName:  "User",
		Password:  "existing_password",
	}, nil)

	mockLogger.EXPECT().Error("Email already exists", gomock.Any())

	res, err := authService.Register(request)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "Email already exists", err.Message)
	assert.Equal(t, "error", err.Status)
}

func TestRegister_HashPasswordError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRefreshToken := mock_repository.NewMockRefreshTokenRepository(ctrl)
	mockUserRole := mock_repository.NewMockUserRoleRepository(ctrl)
	mockRole := mock_repository.NewMockRoleRepository(ctrl)

	mockHash := mock_hash.NewMockHashPassword(ctrl)
	mockToken := mock_auth.NewMockTokenManager(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	authService := service.NewAuthService(mockUserRepo, mockRefreshToken,
		mockRole,
		mockUserRole,
		mockHash, mockToken, mockLogger, mockMapping)

	request := &requests.CreateUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	mockUserRepo.EXPECT().FindByEmail(request.Email).Return(nil, fmt.Errorf("user not found"))
	mockHash.EXPECT().HashPassword(request.Password).Return("", fmt.Errorf("hash error"))
	mockLogger.EXPECT().Error("Failed to hash password", gomock.Any())

	res, err := authService.Register(request)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to hash password", err.Message)
	assert.Equal(t, "error", err.Status)
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRefreshToken := mock_repository.NewMockRefreshTokenRepository(ctrl)
	mockUserRole := mock_repository.NewMockUserRoleRepository(ctrl)
	mockRole := mock_repository.NewMockRoleRepository(ctrl)

	mockHash := mock_hash.NewMockHashPassword(ctrl)
	mockToken := mock_auth.NewMockTokenManager(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	authService := service.NewAuthService(mockUserRepo, mockRefreshToken,
		mockRole,
		mockUserRole,
		mockHash, mockToken, mockLogger, mockMapping)

	request := &requests.AuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	userRecord := &record.UserRecord{
		ID:        1,
		Email:     "test@example.com",
		Password:  "hashed_password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	expectedToken := "jwt_token_123"

	mockUserRepo.EXPECT().FindByEmail(request.Email).Return(userRecord, nil)
	mockHash.EXPECT().ComparePassword(userRecord.Password, request.Password).Return(nil)
	mockToken.EXPECT().GenerateToken(userRecord.ID, "access").Return(expectedToken, nil)
	mockToken.EXPECT().GenerateToken(userRecord.ID, "refresh").Return(expectedToken, nil)

	mockRefreshToken.EXPECT().DeleteRefreshTokenByUserId(userRecord.ID).Return(nil)
	mockRefreshToken.EXPECT().CreateRefreshToken(gomock.Any()).Return(&record.RefreshTokenRecord{}, nil)

	mockLogger.EXPECT().Debug("User logged in successfully", gomock.Any())

	token, err := authService.Login(request)

	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, expectedToken, token.AccessToken)
	assert.Equal(t, expectedToken, token.RefreshToken)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRefreshToken := mock_repository.NewMockRefreshTokenRepository(ctrl)
	mockUserRole := mock_repository.NewMockUserRoleRepository(ctrl)
	mockRole := mock_repository.NewMockRoleRepository(ctrl)

	mockHash := mock_hash.NewMockHashPassword(ctrl)
	mockToken := mock_auth.NewMockTokenManager(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	authService := service.NewAuthService(mockUserRepo, mockRefreshToken,
		mockRole,
		mockUserRole,
		mockHash, mockToken, mockLogger, mockMapping)

	request := &requests.AuthRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	userRecord := &record.UserRecord{
		ID:        1,
		Email:     "test@example.com",
		Password:  "hashed_password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	mockUserRepo.EXPECT().FindByEmail(request.Email).Return(userRecord, nil)
	mockHash.EXPECT().ComparePassword(userRecord.Password, request.Password).Return(fmt.Errorf("invalid password"))
	mockLogger.EXPECT().Error("Failed to compare password", gomock.Any())

	token, err := authService.Login(request)

	assert.Nil(t, token)
	assert.NotNil(t, err)
	assert.Equal(t, "error", err.Status)
	assert.Equal(t, "Invalid password", err.Message)
}

func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRefreshToken := mock_repository.NewMockRefreshTokenRepository(ctrl)
	mockUserRole := mock_repository.NewMockUserRoleRepository(ctrl)
	mockRole := mock_repository.NewMockRoleRepository(ctrl)

	mockHash := mock_hash.NewMockHashPassword(ctrl)
	mockToken := mock_auth.NewMockTokenManager(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	authService := service.NewAuthService(mockUserRepo, mockRefreshToken,
		mockRole,
		mockUserRole,
		mockHash, mockToken, mockLogger, mockMapping)

	request := &requests.AuthRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	mockUserRepo.EXPECT().FindByEmail(request.Email).Return(nil, fmt.Errorf("user not found"))
	mockLogger.EXPECT().Error("Failed to get user", gomock.Any())

	token, err := authService.Login(request)

	// Assertions
	assert.Nil(t, token)
	assert.NotNil(t, err)
	assert.Equal(t, "error", err.Status)
	assert.Equal(t, "Failed to get user: user not found", err.Message)
}

func TestLogin_TokenGenerationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRefreshToken := mock_repository.NewMockRefreshTokenRepository(ctrl)
	mockUserRole := mock_repository.NewMockUserRoleRepository(ctrl)
	mockRole := mock_repository.NewMockRoleRepository(ctrl)

	mockHash := mock_hash.NewMockHashPassword(ctrl)
	mockToken := mock_auth.NewMockTokenManager(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapping := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	authService := service.NewAuthService(mockUserRepo, mockRefreshToken,
		mockRole,
		mockUserRole,
		mockHash, mockToken, mockLogger, mockMapping)

	request := &requests.AuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	userRecord := &record.UserRecord{
		ID:        1,
		Email:     "test@example.com",
		Password:  "hashed_password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	tokenError := fmt.Errorf("failed to generate token")

	mockUserRepo.EXPECT().FindByEmail(request.Email).Return(userRecord, nil)
	mockHash.EXPECT().ComparePassword(userRecord.Password, request.Password).Return(nil)
	mockToken.EXPECT().GenerateToken(userRecord.ID, "access").Return("", tokenError)
	mockLogger.EXPECT().Error("Failed to generate JWT token", gomock.Any())

	token, err := authService.Login(request)

	assert.Nil(t, token)
	assert.NotNil(t, err)
	assert.Equal(t, "error", err.Status)
	assert.Equal(t, "Failed to generate token: failed to generate token", err.Message)
}
