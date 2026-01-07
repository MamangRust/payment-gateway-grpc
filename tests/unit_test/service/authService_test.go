package test

import (
	"database/sql"
	"errors"
	"testing"

	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	mock_responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/mocks"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/auth"
	mock_auth "MamangRust/paymentgatewaygrpc/pkg/auth/mocks"
	refreshtoken_errors "MamangRust/paymentgatewaygrpc/pkg/errors/refresh_token_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/user_errors"
	mock_hash "MamangRust/paymentgatewaygrpc/pkg/hash/mocks"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type AuthServiceSuite struct {
	suite.Suite
	mockUserRepo         *mocks.MockUserRepository
	mockRefreshTokenRepo *mocks.MockRefreshTokenRepository
	mockUserRoleRepo     *mocks.MockUserRoleRepository
	mockRoleRepo         *mocks.MockRoleRepository
	mockHash             *mock_hash.MockHashPassword
	mockToken            *mock_auth.MockTokenManager
	mockLogger           *mock_logger.MockLoggerInterface
	mockMapper           *mock_responseservice.MockUserResponseMapper
	mockCtrl             *gomock.Controller
}

func (suite *AuthServiceSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockUserRepo = mocks.NewMockUserRepository(suite.mockCtrl)
	suite.mockRefreshTokenRepo = mocks.NewMockRefreshTokenRepository(suite.mockCtrl)
	suite.mockUserRoleRepo = mocks.NewMockUserRoleRepository(suite.mockCtrl)
	suite.mockRoleRepo = mocks.NewMockRoleRepository(suite.mockCtrl)
	suite.mockHash = mock_hash.NewMockHashPassword(suite.mockCtrl)
	suite.mockToken = mock_auth.NewMockTokenManager(suite.mockCtrl)
	suite.mockLogger = mock_logger.NewMockLoggerInterface(suite.mockCtrl)
	suite.mockMapper = mock_responseservice.NewMockUserResponseMapper(suite.mockCtrl)
}

func (suite *AuthServiceSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *AuthServiceSuite) TestRegister_Success() {
	req := &requests.CreateUserRequest{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Password: "password123"}
	hashedPassword := "hashed_password_123"
	newUser := &record.UserRecord{ID: 1, Email: req.Email, Password: hashedPassword}
	defaultRole := &record.RoleRecord{ID: 2, Name: "ROLE_ADMIN"}
	expectedResponse := &response.UserResponse{ID: 1, Email: req.Email}

	suite.mockLogger.EXPECT().Debug("Starting user registration", gomock.Any())

	suite.mockUserRepo.EXPECT().FindByEmail(req.Email).Return(nil, sql.ErrNoRows)

	suite.mockHash.EXPECT().HashPassword(req.Password).Return(hashedPassword, nil)

	suite.mockRoleRepo.EXPECT().FindByName("ROLE_ADMIN").Return(defaultRole, nil)

	suite.mockUserRepo.EXPECT().CreateUser(gomock.Any()).DoAndReturn(func(req *requests.CreateUserRequest) (*record.UserRecord, error) {
		suite.Equal(hashedPassword, req.Password)
		return newUser, nil
	})

	suite.mockUserRoleRepo.EXPECT().AssignRoleToUser(&requests.CreateUserRoleRequest{UserId: newUser.ID, RoleId: defaultRole.ID}).Return(nil, nil)
	suite.mockMapper.EXPECT().ToUserResponse(newUser).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("User registered successfully", gomock.Any())

	svc := service.NewAuthService(suite.mockUserRepo, suite.mockRefreshTokenRepo, suite.mockRoleRepo, suite.mockUserRoleRepo, suite.mockHash, suite.mockToken, suite.mockLogger, suite.mockMapper)
	result, err := svc.Register(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *AuthServiceSuite) TestRegister_EmailAlreadyExists() {
	req := &requests.CreateUserRequest{Email: "existing@example.com"}
	existingUser := &record.UserRecord{ID: 1, Email: req.Email}

	suite.mockLogger.EXPECT().Debug("Starting user registration", gomock.Any())
	suite.mockUserRepo.EXPECT().FindByEmail(req.Email).Return(existingUser, nil)
	suite.mockLogger.EXPECT().Debug("Email already exists", gomock.Any())

	svc := service.NewAuthService(suite.mockUserRepo, suite.mockRefreshTokenRepo, suite.mockRoleRepo, suite.mockUserRoleRepo, suite.mockHash, suite.mockToken, suite.mockLogger, suite.mockMapper)
	result, err := svc.Register(req)

	suite.Nil(result)
	suite.Equal(user_errors.ErrUserEmailAlready, err)
}

func (suite *AuthServiceSuite) TestRegister_FailedCreateUser() {
	req := &requests.CreateUserRequest{Email: "fail@example.com"}
	hashedPassword := "hashed_password_123"
	defaultRole := &record.RoleRecord{ID: 2, Name: "ROLE_ADMIN"}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Starting user registration", gomock.Any())
	suite.mockUserRepo.EXPECT().FindByEmail(req.Email).Return(nil, sql.ErrNoRows)
	suite.mockHash.EXPECT().HashPassword(req.Password).Return(hashedPassword, nil)
	suite.mockRoleRepo.EXPECT().FindByName("ROLE_ADMIN").Return(defaultRole, nil)
	suite.mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to create user", gomock.Any())

	svc := service.NewAuthService(suite.mockUserRepo, suite.mockRefreshTokenRepo, suite.mockRoleRepo, suite.mockUserRoleRepo, suite.mockHash, suite.mockToken, suite.mockLogger, suite.mockMapper)
	result, err := svc.Register(req)

	suite.Nil(result)
	suite.Equal(user_errors.ErrFailedCreateUser, err)
}

func (suite *AuthServiceSuite) TestLogin_Success() {
	req := &requests.AuthRequest{Email: "test@example.com", Password: "password"}
	user := &record.UserRecord{ID: 1, Email: req.Email, Password: "hashed_password"}
	accessToken := "access_token_123"
	refreshToken := "refresh_token_123"

	suite.mockLogger.EXPECT().Debug("Starting login process", gomock.Any())
	suite.mockUserRepo.EXPECT().FindByEmail(req.Email).Return(user, nil)
	suite.mockHash.EXPECT().ComparePassword(user.Password, req.Password).Return(nil)

	suite.mockLogger.EXPECT().Debug("Creating access token", zap.Int("userID", user.ID))
	suite.mockToken.EXPECT().GenerateToken(user.ID, "access").Return(accessToken, nil)
	suite.mockLogger.EXPECT().Debug("Access token created successfully", zap.Int("userID", user.ID))

	suite.mockLogger.EXPECT().Debug("Creating refresh token", zap.Int("userID", user.ID))
	suite.mockToken.EXPECT().GenerateToken(user.ID, "refresh").Return(refreshToken, nil)
	suite.mockRefreshTokenRepo.EXPECT().DeleteRefreshTokenByUserId(user.ID).Return(sql.ErrNoRows)
	suite.mockRefreshTokenRepo.EXPECT().CreateRefreshToken(gomock.Any()).Return(nil, nil)
	suite.mockLogger.EXPECT().Debug("Refresh token created successfully", zap.Int("userID", user.ID))

	suite.mockLogger.EXPECT().Debug("User logged in successfully", zap.String("email", req.Email))

	svc := service.NewAuthService(suite.mockUserRepo, suite.mockRefreshTokenRepo, suite.mockRoleRepo, suite.mockUserRoleRepo, suite.mockHash, suite.mockToken, suite.mockLogger, suite.mockMapper)
	result, err := svc.Login(req)

	suite.Nil(err)
	suite.Equal(&response.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}, result)
}

func (suite *AuthServiceSuite) TestLogin_UserNotFound() {
	req := &requests.AuthRequest{Email: "notfound@example.com"}
	repoError := errors.New("user not found")

	suite.mockLogger.EXPECT().Debug("Starting login process", gomock.Any())
	suite.mockUserRepo.EXPECT().FindByEmail(req.Email).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("Failed to get user", gomock.Any())

	svc := service.NewAuthService(suite.mockUserRepo, suite.mockRefreshTokenRepo, suite.mockRoleRepo, suite.mockUserRoleRepo, suite.mockHash, suite.mockToken, suite.mockLogger, suite.mockMapper)
	result, err := svc.Login(req)

	suite.Nil(result)
	suite.Equal(user_errors.ErrUserNotFoundRes, err)
}

func (suite *AuthServiceSuite) TestLogin_PasswordMismatch() {
	req := &requests.AuthRequest{Email: "test@example.com", Password: "wrongpassword"}
	user := &record.UserRecord{ID: 1, Email: req.Email, Password: "hashed_password"}

	suite.mockLogger.EXPECT().Debug("Starting login process", gomock.Any())
	suite.mockUserRepo.EXPECT().FindByEmail(req.Email).Return(user, nil)
	suite.mockHash.EXPECT().ComparePassword(user.Password, req.Password).Return(errors.New("password mismatch"))
	suite.mockLogger.EXPECT().Error("Failed to compare password", gomock.Any())

	svc := service.NewAuthService(suite.mockUserRepo, suite.mockRefreshTokenRepo, suite.mockRoleRepo, suite.mockUserRoleRepo, suite.mockHash, suite.mockToken, suite.mockLogger, suite.mockMapper)
	result, err := svc.Login(req)

	suite.Nil(result)
	suite.Equal(user_errors.ErrUserPassword, err)
}

func (suite *AuthServiceSuite) TestRefreshToken_Success() {
	token := "old_refresh_token"
	userIdStr := "1"
	userId := 1
	newAccessToken := "new_access_token"
	newRefreshToken := "new_refresh_token"

	suite.mockLogger.EXPECT().Debug("Refreshing token", zap.String("token", token))
	suite.mockToken.EXPECT().ValidateToken(token).Return(userIdStr, nil)

	suite.mockLogger.EXPECT().Debug("Creating access token", zap.Int("userID", userId))
	suite.mockToken.EXPECT().GenerateToken(userId, "access").Return(newAccessToken, nil)
	suite.mockLogger.EXPECT().Debug("Access token created successfully", zap.Int("userID", userId))

	suite.mockLogger.EXPECT().Debug("Creating refresh token", zap.Int("userID", userId))
	suite.mockToken.EXPECT().GenerateToken(userId, "refresh").Return(newRefreshToken, nil)
	suite.mockRefreshTokenRepo.EXPECT().DeleteRefreshTokenByUserId(userId).Return(sql.ErrNoRows)
	suite.mockRefreshTokenRepo.EXPECT().CreateRefreshToken(gomock.Any()).Return(nil, nil)
	suite.mockLogger.EXPECT().Debug("Refresh token created successfully", zap.Int("userID", userId))

	suite.mockRefreshTokenRepo.EXPECT().UpdateRefreshToken(gomock.Any()).Return(nil, nil)
	suite.mockLogger.EXPECT().Debug("Refresh token refreshed successfully")

	svc := service.NewAuthService(suite.mockUserRepo, suite.mockRefreshTokenRepo, suite.mockRoleRepo, suite.mockUserRoleRepo, suite.mockHash, suite.mockToken, suite.mockLogger, suite.mockMapper)
	result, err := svc.RefreshToken(token)

	suite.Nil(err)
	suite.Equal(&response.TokenResponse{AccessToken: newAccessToken, RefreshToken: newRefreshToken}, result)
}

func (suite *AuthServiceSuite) TestRefreshToken_TokenExpired() {
	token := "expired_token"
	userIdStr := "1"

	suite.mockLogger.EXPECT().Debug("Refreshing token", gomock.Any())
	suite.mockToken.EXPECT().ValidateToken(token).Return(userIdStr, auth.ErrTokenExpired)
	suite.mockRefreshTokenRepo.EXPECT().DeleteRefreshToken(token).Return(nil)
	suite.mockLogger.EXPECT().Error("Refresh token has expired", gomock.Any())

	svc := service.NewAuthService(suite.mockUserRepo, suite.mockRefreshTokenRepo, suite.mockRoleRepo, suite.mockUserRoleRepo, suite.mockHash, suite.mockToken, suite.mockLogger, suite.mockMapper)
	result, err := svc.RefreshToken(token)

	suite.Nil(result)
	suite.Equal(refreshtoken_errors.ErrFailedExpire, err)
}

func (suite *AuthServiceSuite) TestRefreshToken_InvalidToken() {
	token := "invalid_token"

	suite.mockLogger.EXPECT().Debug("Refreshing token", gomock.Any())
	suite.mockToken.EXPECT().ValidateToken(token).Return("", errors.New("invalid token"))
	suite.mockLogger.EXPECT().Error("Invalid refresh token", gomock.Any())

	svc := service.NewAuthService(suite.mockUserRepo, suite.mockRefreshTokenRepo, suite.mockRoleRepo, suite.mockUserRoleRepo, suite.mockHash, suite.mockToken, suite.mockLogger, suite.mockMapper)
	result, err := svc.RefreshToken(token)

	suite.Nil(result)
	suite.Equal(refreshtoken_errors.ErrRefreshTokenNotFound, err)
}

func (suite *AuthServiceSuite) TestGetMe_Success() {
	token := "valid_token"
	userIdStr := "1"
	userId := 1
	user := &record.UserRecord{ID: userId, Email: "test@example.com"}
	expectedResponse := &response.UserResponse{ID: userId, Email: "test@example.com"}

	suite.mockLogger.EXPECT().Debug("Fetching user details", gomock.Any())
	suite.mockToken.EXPECT().ValidateToken(token).Return(userIdStr, nil)
	suite.mockUserRepo.EXPECT().FindById(userId).Return(user, nil)
	suite.mockMapper.EXPECT().ToUserResponse(user).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("User details fetched successfully", gomock.Any())

	svc := service.NewAuthService(suite.mockUserRepo, suite.mockRefreshTokenRepo, suite.mockRoleRepo, suite.mockUserRoleRepo, suite.mockHash, suite.mockToken, suite.mockLogger, suite.mockMapper)
	result, err := svc.GetMe(token)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *AuthServiceSuite) TestGetMe_InvalidToken() {
	token := "invalid_token"

	suite.mockLogger.EXPECT().Debug("Fetching user details", gomock.Any())
	suite.mockToken.EXPECT().ValidateToken(token).Return("", errors.New("invalid token"))
	suite.mockLogger.EXPECT().Error("Invalid access token", gomock.Any())

	svc := service.NewAuthService(suite.mockUserRepo, suite.mockRefreshTokenRepo, suite.mockRoleRepo, suite.mockUserRoleRepo, suite.mockHash, suite.mockToken, suite.mockLogger, suite.mockMapper)
	result, err := svc.GetMe(token)

	suite.Nil(result)
	suite.Equal(refreshtoken_errors.ErrFailedInValidToken, err)
}

func (suite *AuthServiceSuite) TestGetMe_UserNotFound() {
	token := "valid_token_for_nonexistent_user"
	userIdStr := "999"
	userId := 999

	suite.mockLogger.EXPECT().Debug("Fetching user details", gomock.Any())
	suite.mockToken.EXPECT().ValidateToken(token).Return(userIdStr, nil)
	suite.mockUserRepo.EXPECT().FindById(userId).Return(nil, errors.New("user not found"))
	suite.mockLogger.EXPECT().Error("Failed to find user by ID", gomock.Any())

	svc := service.NewAuthService(suite.mockUserRepo, suite.mockRefreshTokenRepo, suite.mockRoleRepo, suite.mockUserRoleRepo, suite.mockHash, suite.mockToken, suite.mockLogger, suite.mockMapper)
	result, err := svc.GetMe(token)

	suite.Nil(result)
	suite.Equal(user_errors.ErrUserNotFoundRes, err)
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceSuite))
}
