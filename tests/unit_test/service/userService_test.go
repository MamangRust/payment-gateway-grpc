package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	mock_responsemapper "MamangRust/paymentgatewaygrpc/internal/mapper/response/mocks"
	mock_repository "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/internal/service"
	mock_hash "MamangRust/paymentgatewaygrpc/pkg/hash/mocks"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestUserService_FindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		mockMapper,
		nil,
	)

	page := 1
	pageSize := 10
	search := "John"
	totalRecords := 3

	users := []*record.UserRecord{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "M0Vn2@example.com",
			Password:  "password123",
			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},

		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "0KdXb@example.com",
			Password:  "password123",
			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},
		{
			ID:        3,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "M0Vn2@example.com",
			Password:  "password123",
			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},
	}

	expectedResponse := []*response.UserResponse{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "M0Vn2@example.com",

			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},

		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "0KdXb@example.com",

			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},
		{
			ID:        3,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "M0Vn2@example.com",

			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},
	}

	mockUserRepo.EXPECT().
		FindAllUsers(search, page, pageSize).
		Return(users, totalRecords, nil)

	mockMapper.EXPECT().
		ToUsersResponse(users).
		Return(expectedResponse)

	results, totalPages, errResp := userService.FindAll(page, pageSize, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, results)
	assert.Equal(t, totalPages, totalRecords)
}

func TestUserService_FindAll_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		nil,
		nil,
	)

	page := 1
	pageSize := 10
	search := "Jane"

	mockUserRepo.EXPECT().
		FindAllUsers(search, page, pageSize).
		Return(nil, 0, errors.New("database error"))

	mockLogger.EXPECT().
		Error("failed to fetch users", gomock.Any())

	results, totalPages, errResp := userService.FindAll(page, pageSize, search)

	assert.Nil(t, results)
	assert.Equal(t, 0, totalPages)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to fetch users", errResp.Message)
}

func TestUserService_FindByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		mockMapper,
		nil,
	)

	user := &record.UserRecord{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "M0Vn2@example.com",
		Password:  "password123",
		CreatedAt: "2024-12-21T09:00:00Z",
		UpdatedAt: "2024-12-21T09:00:00Z",
	}

	expectedResponse := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "M0Vn2@example.com",
		CreatedAt: "2024-12-21T09:00:00Z",
		UpdatedAt: "2024-12-21T09:00:00Z",
	}

	mockUserRepo.EXPECT().
		FindById(1).
		Return(user, nil)

	mockMapper.EXPECT().
		ToUserResponse(user).
		Return(expectedResponse)

	result, errResp := userService.FindByID(1)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedResponse, result)
}

func TestUserService_FindByID_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		nil,
		nil,
	)

	mockUserRepo.EXPECT().
		FindById(999).
		Return(nil, errors.New("user not found"))

	mockLogger.EXPECT().
		Error("failed to find user by ID", gomock.Any())

	result, errResp := userService.FindByID(999)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "User not found", errResp.Message)
}

func TestUserService_FindByActive_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		mockMapper,
		nil,
	)

	users := []*record.UserRecord{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "M0Vn2@example.com",
			Password:  "password123",
			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "0KdXb@example.com",
			Password:  "password123",
			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},
	}

	expectedResponse := []*response.UserResponseDeleteAt{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "M0Vn2@example.com",
			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "0KdXb@example.com",
			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},
	}

	page := 1
	pageSize := 1
	search := "John"
	expected := 2

	mockUserRepo.EXPECT().
		FindByActive(search, page, pageSize).
		Return(users, expected, nil)

	mockMapper.EXPECT().
		ToUsersResponseDeleteAt(users).
		Return(expectedResponse)

	result, totalRecord, errResp := userService.FindByActive(pageSize, page, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, expectedResponse, result)
}

func TestUserService_FindByActive_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		nil,
		nil,
	)

	page := 1
	pageSize := 1
	search := "John"
	expected := 0

	mockUserRepo.EXPECT().
		FindByActive(search, page, pageSize).
		Return(nil, expected, errors.New("database error"))

	mockLogger.EXPECT().
		Error("Failed to find active users", gomock.Any())

	result, totalRecord, errResp := userService.FindByActive(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, totalRecord, expected)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to find active users", errResp.Message)
}

func TestUserService_FindByTrashed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		mockMapper,
		nil,
	)

	users := []*record.UserRecord{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "M0Vn2@example.com",
			Password:  "password123",
			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "0KdXb@example.com",
			Password:  "password123",
			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},
	}

	expectedResponse := []*response.UserResponseDeleteAt{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "M0Vn2@example.com",
			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "0KdXb@example.com",
			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
		},
	}

	page := 1
	pageSize := 1
	search := "John"
	expected := 2

	mockUserRepo.EXPECT().
		FindByTrashed(search, page, pageSize).
		Return(users, expected, nil)

	mockMapper.EXPECT().
		ToUsersResponseDeleteAt(users).
		Return(expectedResponse)

	result, totalRecord, errResp := userService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, expectedResponse, result)
}

func TestUserService_FindByTrashed_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		nil,
		nil,
	)

	page := 1
	pageSize := 1
	search := "John"
	expected := 0

	mockUserRepo.EXPECT().
		FindByTrashed(search, page, pageSize).
		Return(nil, expected, errors.New("database error"))

	mockLogger.EXPECT().
		Error("Failed to find trashed users", gomock.Any())

	result, totalRecord, errResp := userService.FindByTrashed(pageSize, page, search)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to find trashed users", errResp.Message)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockUserResponseMapper(ctrl)
	mockHashing := mock_hash.NewMockHashPassword(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		mockMapper,
		mockHashing,
	)

	request := &requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	hashedPassword := "hashedpassword123"

	expectedUser := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	mockUserRepo.EXPECT().
		FindByEmail(request.Email).
		Return(nil, nil)

	mockHashing.EXPECT().
		HashPassword(request.Password).
		Return(hashedPassword, nil)

	mockUserRepo.EXPECT().
		CreateUser(request).
		Return(&record.UserRecord{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  hashedPassword,
		}, nil)

	mockMapper.EXPECT().
		ToUserResponse(gomock.Any()).
		Return(expectedUser)

	result, errResp := userService.CreateUser(request)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedUser, result)
}

func TestUserService_CreateUser_EmailAlreadyInUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		nil,
		nil,
	)

	request := &requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	existingUser := &record.UserRecord{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	mockUserRepo.EXPECT().
		FindByEmail(request.Email).
		Return(existingUser, nil)

	mockLogger.EXPECT().
		Error("Email is already in use", gomock.Any())

	result, errResp := userService.CreateUser(request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Email is already in use", errResp.Message)
}

func TestUserService_CreateUser_PasswordHashingFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockUserResponseMapper(ctrl)
	mockHashing := mock_hash.NewMockHashPassword(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		mockMapper,
		mockHashing,
	)

	request := &requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockUserRepo.EXPECT().
		FindByEmail(request.Email).
		Return(nil, nil)

	mockHashing.EXPECT().
		HashPassword(request.Password).
		Return("", errors.New("hashing failed"))

	mockLogger.EXPECT().
		Error("Failed to hash password", gomock.Any())

	result, errResp := userService.CreateUser(request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to hash password", errResp.Message)
}

func TestUserService_UpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockUserResponseMapper(ctrl)
	mockHashing := mock_hash.NewMockHashPassword(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		mockMapper,
		mockHashing,
	)

	request := &requests.UpdateUserRequest{
		UserID:          1,
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@newemail.com",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	existingUser := &record.UserRecord{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	hashedPassword := "hashednewpassword123"

	expectedUser := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@newemail.com",
	}

	mockUserRepo.EXPECT().
		FindById(request.UserID).
		Return(existingUser, nil)

	mockUserRepo.EXPECT().
		FindByEmail(request.Email).
		Return(nil, nil)

	mockHashing.EXPECT().
		HashPassword(request.Password).
		Return(hashedPassword, nil)

	mockUserRepo.EXPECT().
		UpdateUser(request).
		Return(&record.UserRecord{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@newemail.com",
			Password:  hashedPassword,
		}, nil)

	mockMapper.EXPECT().
		ToUserResponse(gomock.Any()).
		Return(expectedUser)

	result, errResp := userService.UpdateUser(request)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedUser, result)
}

func TestUserService_UpdateUser_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		nil,
		nil,
	)

	request := &requests.UpdateUserRequest{
		UserID:          999,
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@newemail.com",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	mockUserRepo.EXPECT().
		FindById(request.UserID).
		Return(nil, errors.New("user not found"))

	mockLogger.EXPECT().
		Error("Failed to find user by ID", gomock.Any())

	result, errResp := userService.UpdateUser(request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "User not found", errResp.Message)
}

func TestUserService_UpdateUser_EmailAlreadyInUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		nil,
		nil,
	)

	request := &requests.UpdateUserRequest{
		UserID:          1,
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@newemail.com",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	existingUser := &record.UserRecord{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	duplicateUser := &record.UserRecord{
		ID:        2,
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "john.doe@newemail.com",
	}

	mockUserRepo.EXPECT().
		FindById(request.UserID).
		Return(existingUser, nil)

	mockUserRepo.EXPECT().
		FindByEmail(request.Email).
		Return(duplicateUser, nil)

	result, errResp := userService.UpdateUser(request)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Email is already in use", errResp.Message)
}

func TestUserService_TrashedUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		mockMapper,
		nil,
	)

	userID := 1
	expectedUser := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	mockLogger.EXPECT().Debug("Trashing user", zap.Int("user_id", userID)).Times(1)

	mockUserRepo.EXPECT().
		TrashedUser(userID).
		Return(&record.UserRecord{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}, nil)

	mockMapper.EXPECT().
		ToUserResponse(gomock.Any()).
		Return(expectedUser)

	mockLogger.EXPECT().Debug("Successfully trashed user", zap.Int("user_id", userID)).Times(1)

	result, errResp := userService.TrashedUser(userID)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedUser, result)
}

func TestUserService_TrashedUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		nil,
		nil,
	)

	userID := 1

	mockLogger.EXPECT().Debug("Trashing user", zap.Int("user_id", userID)).Times(1)

	mockUserRepo.EXPECT().
		TrashedUser(userID).
		Return(nil, errors.New("trash user failed"))

	mockLogger.EXPECT().
		Error("Failed to trash user", zap.Error(errors.New("trash user failed")), zap.Int("user_id", userID))

	result, errResp := userService.TrashedUser(userID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to trash user", errResp.Message)
}

func TestUserService_RestoreUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_responsemapper.NewMockUserResponseMapper(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		mockMapper,
		nil,
	)

	userID := 1
	expectedUser := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	mockLogger.EXPECT().Debug("Restoring user", zap.Int("user_id", userID)).Times(1)

	mockUserRepo.EXPECT().
		RestoreUser(userID).
		Return(&record.UserRecord{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}, nil)

	mockMapper.EXPECT().
		ToUserResponse(gomock.Any()).
		Return(expectedUser)

	mockLogger.EXPECT().Debug("Successfully restored user", zap.Int("user_id", userID)).Times(1)

	result, errResp := userService.RestoreUser(userID)

	assert.Nil(t, errResp)
	assert.Equal(t, expectedUser, result)
}

func TestUserService_RestoreUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		nil,
		nil,
	)

	userID := 1

	mockLogger.EXPECT().Debug("Restoring user", zap.Int("user_id", userID)).Times(1)

	mockUserRepo.EXPECT().
		RestoreUser(userID).
		Return(nil, errors.New("restore user failed"))

	mockLogger.EXPECT().
		Error("Failed to restore user", gomock.Any(), gomock.Any())

	result, errResp := userService.RestoreUser(userID)

	assert.Nil(t, result)
	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to restore user", errResp.Message)
}

func TestUserService_DeleteUserPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		nil,
		nil,
	)

	userID := 1

	mockLogger.EXPECT().Debug("Deleting user permanently", zap.Int("user_id", userID)).Times(1)

	mockUserRepo.EXPECT().
		DeleteUserPermanent(userID).
		Return(nil)
	mockLogger.EXPECT().Debug("Successfully deleted user permanently", zap.Int("user_id", userID)).Times(1)

	result, errResp := userService.DeleteUserPermanent(userID)

	assert.Nil(t, result)
	assert.Nil(t, errResp)
}

func TestUserService_DeleteUserPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userService := service.NewUserService(
		mockUserRepo,
		mockLogger,
		nil,
		nil,
	)

	userID := 1

	mockLogger.EXPECT().
		Debug("Deleting user permanently", gomock.Any()).
		Times(1)

	mockUserRepo.EXPECT().
		DeleteUserPermanent(userID).
		Return(errors.New("Failed to delete user permanently"))

	mockLogger.EXPECT().
		Error("Failed to delete user permanently",
			gomock.Any(),
			zap.Int("user_id", userID),
		)

	_, errResp := userService.DeleteUserPermanent(userID)

	assert.NotNil(t, errResp)
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "Failed to delete user permanently", errResp.Message)
}
