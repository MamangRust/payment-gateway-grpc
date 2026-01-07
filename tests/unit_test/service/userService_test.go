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
	"MamangRust/paymentgatewaygrpc/pkg/errors/user_errors"
	mock_hash "MamangRust/paymentgatewaygrpc/pkg/hash/mocks"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type UserServiceSuite struct {
	suite.Suite
	mockRepo   *mocks.MockUserRepository
	mockLogger *mock_logger.MockLoggerInterface
	mockHash   *mock_hash.MockHashPassword
	mockMapper *mock_responseservice.MockUserResponseMapper
	mockCtrl   *gomock.Controller
}

func (suite *UserServiceSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockUserRepository(suite.mockCtrl)
	suite.mockLogger = mock_logger.NewMockLoggerInterface(suite.mockCtrl)
	suite.mockMapper = mock_responseservice.NewMockUserResponseMapper(suite.mockCtrl)
	suite.mockHash = mock_hash.NewMockHashPassword(suite.mockCtrl)
}

func (suite *UserServiceSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *UserServiceSuite) TestFindAll_Success() {
	req := &requests.FindAllUsers{Search: "john", Page: 1, PageSize: 10}
	total := 1
	mockUsers := []*record.UserRecord{{ID: 1, FirstName: "John", LastName: "Doe"}}
	expectedResponse := []*response.UserResponse{{ID: 1, FirstName: "John", LastName: "Doe"}}

	suite.mockLogger.EXPECT().Debug("Fetching users", gomock.Any())
	suite.mockRepo.EXPECT().FindAllUsers(req).Return(mockUsers, &total, nil)
	suite.mockMapper.EXPECT().ToUsersResponse(mockUsers).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched user", gomock.Any())

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, totalRes, err := service.FindAll(req)

	suite.Nil(err)
	suite.NotNil(result)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *UserServiceSuite) TestFindAll_Failure() {
	req := &requests.FindAllUsers{Search: "john", Page: 1, PageSize: 10}
	dbError := errors.New("database connection failed")

	suite.mockLogger.EXPECT().Debug("Fetching users", gomock.Any())
	suite.mockRepo.EXPECT().FindAllUsers(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch user", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, totalRes, err := service.FindAll(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(user_errors.ErrFailedFindAll, err)
}

func (suite *UserServiceSuite) TestFindByID_Success() {
	userID := 1
	mockUser := &record.UserRecord{ID: userID, FirstName: "John", LastName: "Doe"}
	expectedResponse := &response.UserResponse{ID: userID, FirstName: "John", LastName: "Doe"}

	suite.mockLogger.EXPECT().Debug("Fetching user by id", zap.Int("user_id", userID))
	suite.mockRepo.EXPECT().FindById(userID).Return(mockUser, nil)
	suite.mockMapper.EXPECT().ToUserResponse(mockUser).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched user", zap.Int("user_id", userID))

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := service.FindByID(userID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *UserServiceSuite) TestFindByID_NotFound() {
	userID := 99
	repoError := errors.New("user not found")

	suite.mockLogger.EXPECT().Debug("Fetching user by id", zap.Int("user_id", userID))
	suite.mockRepo.EXPECT().FindById(userID).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("failed to find user by ID", gomock.Any())

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := service.FindByID(userID)

	suite.Nil(result)
	suite.Equal(user_errors.ErrUserNotFoundRes, err)
}

func (suite *UserServiceSuite) TestFindByActive_Success() {
	req := &requests.FindAllUsers{Search: "active", Page: 1, PageSize: 10}
	total := 1
	mockUsers := []*record.UserRecord{{ID: 1, FirstName: "Active", LastName: "User", DeletedAt: nil}}
	expectedResponse := []*response.UserResponseDeleteAt{{ID: 1, FirstName: "Active", LastName: "User", DeletedAt: nil}}

	suite.mockLogger.EXPECT().Debug("Fetching active user", gomock.Any())
	suite.mockRepo.EXPECT().FindByActive(req).Return(mockUsers, &total, nil)
	suite.mockMapper.EXPECT().ToUsersResponseDeleteAt(mockUsers).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched active user", gomock.Any())

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, totalRes, err := service.FindByActive(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *UserServiceSuite) TestFindByActive_Failure() {
	req := &requests.FindAllUsers{Search: "active", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching active user", gomock.Any())
	suite.mockRepo.EXPECT().FindByActive(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch active user", gomock.Any())

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, totalRes, err := service.FindByActive(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(user_errors.ErrFailedFindActive, err)
}

func (suite *UserServiceSuite) TestFindByTrashed_Success() {
	req := &requests.FindAllUsers{Page: 1, PageSize: 10}
	total := 1
	trashedTime := "2024-01-01T00:00:00Z"
	mockUsers := []*record.UserRecord{{ID: 2, FirstName: "Trashed", LastName: "User", DeletedAt: &trashedTime}}
	expectedResponse := []*response.UserResponseDeleteAt{{ID: 2, FirstName: "Trashed", LastName: "User", DeletedAt: &trashedTime}}

	suite.mockLogger.EXPECT().Debug("Fetching trashed user", gomock.Any())
	suite.mockRepo.EXPECT().FindByTrashed(req).Return(mockUsers, &total, nil)
	suite.mockMapper.EXPECT().ToUsersResponseDeleteAt(mockUsers).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched trashed user", gomock.Any())

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, totalRes, err := service.FindByTrashed(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *UserServiceSuite) TestFindByTrashed_Failure() {
	req := &requests.FindAllUsers{Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching trashed user", gomock.Any())
	suite.mockRepo.EXPECT().FindByTrashed(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to find trashed users", gomock.Any())

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, totalRes, err := service.FindByTrashed(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(user_errors.ErrFailedFindTrashed, err)
}

func (suite *UserServiceSuite) TestCreateUser_Success() {
	req := &requests.CreateUserRequest{FirstName: "New", LastName: "User", Email: "new@example.com", Password: "password123"}
	hashedPassword := "hashed_password_123"
	createdUser := &record.UserRecord{ID: 10, FirstName: "New", LastName: "User", Email: "new@example.com", Password: hashedPassword}
	expectedResponse := &response.UserResponse{ID: 10, FirstName: "New", LastName: "User", Email: "new@example.com"}

	suite.mockLogger.EXPECT().Debug("Creating new user", zap.String("email", req.Email), zap.Any("request", req))
	suite.mockLogger.EXPECT().Debug("Email is available, proceeding to create user", zap.String("email", req.Email))
	suite.mockLogger.EXPECT().Debug("Successfully created new user", zap.String("email", expectedResponse.Email), zap.Int("user", expectedResponse.ID))

	suite.mockRepo.EXPECT().FindByEmail(req.Email).Return(nil, sql.ErrNoRows)
	suite.mockHash.EXPECT().HashPassword(req.Password).Return(hashedPassword, nil)
	suite.mockRepo.EXPECT().CreateUser(gomock.Any()).DoAndReturn(func(req *requests.CreateUserRequest) (*record.UserRecord, error) {
		suite.Equal(hashedPassword, req.Password)
		return createdUser, nil
	})
	suite.mockMapper.EXPECT().ToUserResponse(createdUser).Return(expectedResponse)

	svc := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := svc.CreateUser(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *UserServiceSuite) TestCreateUser_EmailAlreadyExists() {
	req := &requests.CreateUserRequest{Email: "existing@example.com"}
	existingUser := &record.UserRecord{ID: 1, Email: "existing@example.com"}

	suite.mockLogger.EXPECT().Debug("Creating new user", zap.String("email", req.Email), zap.Any("request", req))
	suite.mockLogger.EXPECT().Error("Email is already in use", zap.String("email", req.Email))

	suite.mockRepo.EXPECT().FindByEmail(req.Email).Return(existingUser, nil)

	svc := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := svc.CreateUser(req)

	suite.Nil(result)
	suite.Equal(user_errors.ErrUserEmailAlready, err)
}

func (suite *UserServiceSuite) TestCreateUser_FindByEmailDbError() {
	req := &requests.CreateUserRequest{Email: "error@example.com"}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Creating new user", zap.String("email", req.Email), zap.Any("request", req))
	suite.mockLogger.EXPECT().Error("Error checking existing email", zap.String("email", req.Email), zap.Error(dbError))

	suite.mockRepo.EXPECT().FindByEmail(req.Email).Return(nil, dbError)

	svc := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := svc.CreateUser(req)

	suite.Nil(result)
	suite.Equal(user_errors.ErrUserEmailAlready, err)
}

func (suite *UserServiceSuite) TestCreateUser_HashPasswordError() {
	req := &requests.CreateUserRequest{Email: "hashfail@example.com", Password: "password"}
	hashError := errors.New("hashing failed")

	suite.mockLogger.EXPECT().Debug("Creating new user", zap.String("email", req.Email), zap.Any("request", req))
	suite.mockLogger.EXPECT().Debug("Email is available, proceeding to create user", zap.String("email", req.Email))
	suite.mockLogger.EXPECT().Error("Failed to hash password", zap.Error(hashError))

	suite.mockRepo.EXPECT().FindByEmail(req.Email).Return(nil, sql.ErrNoRows)
	suite.mockHash.EXPECT().HashPassword(req.Password).Return("", hashError)

	svc := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := svc.CreateUser(req)

	suite.Nil(result)
	suite.Equal(user_errors.ErrUserPassword, err)
}

func (suite *UserServiceSuite) TestCreateUser_CreateUserError() {
	req := &requests.CreateUserRequest{Email: "createfail@example.com", Password: "password"}
	hashedPassword := "hashed_password_123"
	dbError := errors.New("failed to create user")

	suite.mockLogger.EXPECT().Debug("Creating new user", zap.String("email", req.Email), zap.Any("request", req))
	suite.mockLogger.EXPECT().Debug("Email is available, proceeding to create user", zap.String("email", req.Email))
	suite.mockLogger.EXPECT().Error("Failed to create user", zap.Error(dbError))

	suite.mockRepo.EXPECT().FindByEmail(req.Email).Return(nil, sql.ErrNoRows)
	suite.mockHash.EXPECT().HashPassword(req.Password).Return(hashedPassword, nil)
	suite.mockRepo.EXPECT().CreateUser(gomock.Any()).Return(nil, dbError)

	svc := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := svc.CreateUser(req)

	suite.Nil(result)
	suite.Equal(user_errors.ErrFailedCreateUser, err)
}

func (suite *UserServiceSuite) TestUpdateUser_Success() {
	userID := 1
	req := &requests.UpdateUserRequest{UserID: &userID, FirstName: "Updated", LastName: "User", Email: "updated@example.com", Password: "newpass"}
	existingUser := &record.UserRecord{ID: userID, Email: "old@example.com"}
	updatedUser := &record.UserRecord{ID: userID, FirstName: "Updated", LastName: "User", Email: "updated@example.com"}
	expectedResponse := &response.UserResponse{ID: userID, FirstName: "Updated", LastName: "User", Email: "updated@example.com"}
	hashedPassword := "hashed_newpass"

	suite.mockLogger.EXPECT().Debug("Updating user", zap.Int("user_id", userID), zap.Any("request", req))
	suite.mockLogger.EXPECT().Debug("Successfully updated user", zap.Int("user_id", userID))

	suite.mockRepo.EXPECT().FindById(userID).Return(existingUser, nil)
	suite.mockRepo.EXPECT().FindByEmail(req.Email).Return(nil, sql.ErrNoRows)
	suite.mockHash.EXPECT().HashPassword(req.Password).Return(hashedPassword, nil)
	suite.mockRepo.EXPECT().UpdateUser(req).Return(updatedUser, nil)
	suite.mockMapper.EXPECT().ToUserResponse(updatedUser).Return(expectedResponse)

	svc := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := svc.UpdateUser(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *UserServiceSuite) TestUpdateUser_UserNotFound() {
	userID := 99
	req := &requests.UpdateUserRequest{UserID: &userID}
	repoError := errors.New("user not found")

	suite.mockLogger.EXPECT().Debug("Updating user", zap.Int("user_id", userID), zap.Any("request", req))
	suite.mockLogger.EXPECT().Error("Failed to find user by ID", zap.Error(repoError))

	suite.mockRepo.EXPECT().FindById(userID).Return(nil, repoError)

	svc := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := svc.UpdateUser(req)

	suite.Nil(result)
	suite.Equal(user_errors.ErrUserNotFoundRes, err)
}

func (suite *UserServiceSuite) TestUpdateUser_EmailAlreadyExists() {
	userID := 1
	req := &requests.UpdateUserRequest{UserID: &userID, Email: "duplicate@example.com"}
	existingUser := &record.UserRecord{ID: userID, Email: "old@example.com"}
	duplicateUser := &record.UserRecord{ID: 2, Email: "duplicate@example.com"}

	suite.mockLogger.EXPECT().Debug("Updating user", zap.Int("user_id", userID), zap.Any("request", req))

	suite.mockRepo.EXPECT().FindById(userID).Return(existingUser, nil)
	suite.mockRepo.EXPECT().FindByEmail(req.Email).Return(duplicateUser, nil)

	svc := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := svc.UpdateUser(req)

	suite.Nil(result)
	suite.Equal(user_errors.ErrUserEmailAlready, err)
}

func (suite *UserServiceSuite) TestUpdateUser_HashPasswordError() {
	userID := 1
	req := &requests.UpdateUserRequest{UserID: &userID, Password: "newpass"}
	existingUser := &record.UserRecord{ID: userID}
	hashError := errors.New("hashing failed")

	suite.mockLogger.EXPECT().Debug("Updating user", zap.Int("user_id", userID), zap.Any("request", req))
	suite.mockLogger.EXPECT().Error("Failed to hash password", zap.Error(hashError))

	suite.mockRepo.EXPECT().FindById(userID).Return(existingUser, nil)
	suite.mockHash.EXPECT().HashPassword(req.Password).Return("", hashError)

	svc := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := svc.UpdateUser(req)

	suite.Nil(result)
	suite.Equal(user_errors.ErrUserPassword, err)
}

func (suite *UserServiceSuite) TestUpdateUser_UpdateUserError() {
	userID := 1
	req := &requests.UpdateUserRequest{UserID: &userID, FirstName: "Updated"}
	existingUser := &record.UserRecord{ID: userID}
	dbError := errors.New("failed to update user")

	suite.mockLogger.EXPECT().Debug("Updating user", zap.Int("user_id", userID), zap.Any("request", req))
	suite.mockLogger.EXPECT().Error("Failed to update user", zap.Error(dbError))

	suite.mockRepo.EXPECT().FindById(userID).Return(existingUser, nil)
	suite.mockRepo.EXPECT().UpdateUser(req).Return(nil, dbError)

	svc := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := svc.UpdateUser(req)

	suite.Nil(result)
	suite.Equal(user_errors.ErrFailedUpdateUser, err)
}

func (suite *UserServiceSuite) TestTrashedUser_Success() {
	userID := 1
	trashedUser := &record.UserRecord{ID: userID, FirstName: "Trashed", LastName: "User"}
	expectedResponse := &response.UserResponseDeleteAt{ID: userID, FirstName: "Trashed", LastName: "User"}

	suite.mockLogger.EXPECT().Debug("Trashing user", zap.Int("user_id", userID))
	suite.mockRepo.EXPECT().TrashedUser(userID).Return(trashedUser, nil)
	suite.mockMapper.EXPECT().ToUserResponseDeleteAt(trashedUser).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully trashed user", zap.Int("user_id", userID))

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := service.TrashedUser(userID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *UserServiceSuite) TestTrashedUser_Failure() {
	userID := 1
	dbError := errors.New("failed to trash user")

	suite.mockLogger.EXPECT().Debug("Trashing user", zap.Int("user_id", userID))
	suite.mockRepo.EXPECT().TrashedUser(userID).Return(nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to trash user", gomock.Any())

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := service.TrashedUser(userID)

	suite.Nil(result)
	suite.Equal(user_errors.ErrFailedTrashedUser, err)
}

func (suite *UserServiceSuite) TestRestoreUser_Success() {
	userID := 1
	restoredUser := &record.UserRecord{ID: userID, FirstName: "Restored", LastName: "User", DeletedAt: nil}
	expectedResponse := &response.UserResponseDeleteAt{ID: userID, FirstName: "Restored", LastName: "User", DeletedAt: nil}

	suite.mockLogger.EXPECT().Debug("Restoring user", zap.Int("user_id", userID))
	suite.mockRepo.EXPECT().RestoreUser(userID).Return(restoredUser, nil)
	suite.mockMapper.EXPECT().ToUserResponseDeleteAt(restoredUser).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully restored user", zap.Int("user_id", userID))

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := service.RestoreUser(userID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *UserServiceSuite) TestRestoreUser_Failure() {
	userID := 1
	dbError := errors.New("failed to restore user")

	suite.mockLogger.EXPECT().Debug("Restoring user", zap.Int("user_id", userID))
	suite.mockRepo.EXPECT().RestoreUser(userID).Return(nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to restore user", gomock.Any())

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := service.RestoreUser(userID)

	suite.Nil(result)
	suite.Equal(user_errors.ErrFailedRestoreUser, err)
}

func (suite *UserServiceSuite) TestDeleteUserPermanent_Success() {
	userID := 1

	suite.mockLogger.EXPECT().Debug("Deleting user permanently", zap.Int("user_id", userID))
	suite.mockRepo.EXPECT().DeleteUserPermanent(userID).Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted user permanently", zap.Int("user_id", userID))

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := service.DeleteUserPermanent(userID)

	suite.Nil(err)
	suite.True(result)
}

func (suite *UserServiceSuite) TestDeleteUserPermanent_Failure() {
	userID := 1
	dbError := errors.New("failed to delete user permanently")

	suite.mockLogger.EXPECT().Debug("Deleting user permanently", zap.Int("user_id", userID))
	suite.mockRepo.EXPECT().DeleteUserPermanent(userID).Return(false, dbError)
	suite.mockLogger.EXPECT().Error("Failed to delete user permanently", gomock.Any())

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := service.DeleteUserPermanent(userID)

	suite.False(result)
	suite.Equal(user_errors.ErrFailedDeletePermanent, err)
}

func (suite *UserServiceSuite) TestRestoreAllUser_Success() {
	suite.mockLogger.EXPECT().Debug("Restoring all users")
	suite.mockRepo.EXPECT().RestoreAllUser().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully restored all users")

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := service.RestoreAllUser()

	suite.Nil(err)
	suite.True(result)
}

func (suite *UserServiceSuite) TestRestoreAllUser_Failure() {
	dbError := errors.New("failed to restore all users")

	suite.mockLogger.EXPECT().Debug("Restoring all users")
	suite.mockRepo.EXPECT().RestoreAllUser().Return(false, dbError)
	suite.mockLogger.EXPECT().Error("Failed to restore all users", gomock.Any())

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := service.RestoreAllUser()

	suite.False(result)
	suite.Equal(user_errors.ErrFailedRestoreAll, err)
}

func (suite *UserServiceSuite) TestDeleteAllUserPermanent_Success() {
	suite.mockLogger.EXPECT().Debug("Permanently deleting all users")
	suite.mockRepo.EXPECT().DeleteAllUserPermanent().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted all users permanently")

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := service.DeleteAllUserPermanent()

	suite.Nil(err)
	suite.True(result)
}

func (suite *UserServiceSuite) TestDeleteAllUserPermanent_Failure() {
	dbError := errors.New("failed to delete all users permanently")

	suite.mockLogger.EXPECT().Debug("Permanently deleting all users")
	suite.mockRepo.EXPECT().DeleteAllUserPermanent().Return(false, dbError)
	suite.mockLogger.EXPECT().Error("Failed to permanently delete all users", gomock.Any())

	service := service.NewUserService(suite.mockRepo, suite.mockLogger, suite.mockMapper, suite.mockHash)
	result, err := service.DeleteAllUserPermanent()

	suite.False(result)
	suite.Equal(user_errors.ErrFailedDeleteAll, err)
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}
