package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserRepositorySuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	mockRepo *mocks.MockUserRepository
}

func (suite *UserRepositorySuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockUserRepository(suite.mockCtrl)
}

func (suite *UserRepositorySuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *UserRepositorySuite) TestFindAllUsers_Success() {
	req := &requests.FindAllUsers{Search: "", Page: 1, PageSize: 10}
	users := []*record.UserRecord{{ID: 1, FirstName: "John"}}
	total := 1

	suite.mockRepo.EXPECT().FindAllUsers(req).Return(users, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllUsers(req)

	suite.NoError(err)
	suite.Equal(users, result)
	suite.Equal(1, *totalRes)
}

func (suite *UserRepositorySuite) TestFindByActiveUser_Success() {
	req := &requests.FindAllUsers{Search: "active", Page: 1, PageSize: 10}
	activeUsers := []*record.UserRecord{{ID: 1, FirstName: "Active User"}}
	total := 1

	suite.mockRepo.EXPECT().FindByActive(req).Return(activeUsers, &total, nil)

	result, totalRes, err := suite.mockRepo.FindByActive(req)

	suite.NoError(err)
	suite.Equal(activeUsers, result)
	suite.Equal(1, *totalRes)
}

func (suite *UserRepositorySuite) TestFindByTrashedUser_Success() {
	req := &requests.FindAllUsers{Page: 1, PageSize: 10}
	trashedTime := "2024-01-01T00:00:00Z"
	trashedUsers := []*record.UserRecord{{ID: 2, FirstName: "Deleted User", DeletedAt: &trashedTime}}
	total := 1

	suite.mockRepo.EXPECT().FindByTrashed(req).Return(trashedUsers, &total, nil)

	result, totalRes, err := suite.mockRepo.FindByTrashed(req)

	suite.NoError(err)
	suite.Equal(trashedUsers, result)
	suite.Equal(1, *totalRes)
}

func (suite *UserRepositorySuite) TestFindByIdUser_Success() {
	userID := 1
	user := &record.UserRecord{ID: userID, FirstName: "John"}

	suite.mockRepo.EXPECT().FindById(userID).Return(user, nil)

	result, err := suite.mockRepo.FindById(userID)

	suite.NoError(err)
	suite.Equal(user, result)
}

func (suite *UserRepositorySuite) TestFindByIdUser_NotFound() {
	userID := 99
	expectedErr := errors.New("user not found")

	suite.mockRepo.EXPECT().FindById(userID).Return(nil, expectedErr)

	result, err := suite.mockRepo.FindById(userID)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(expectedErr, err)
}

func (suite *UserRepositorySuite) TestFindByEmail_Success() {
	email := "test@example.com"
	user := &record.UserRecord{ID: 1, Email: email}

	suite.mockRepo.EXPECT().FindByEmail(email).Return(user, nil)

	result, err := suite.mockRepo.FindByEmail(email)

	suite.NoError(err)
	suite.Equal(user, result)
}

func (suite *UserRepositorySuite) TestFindByEmail_NotFound() {
	email := "notfound@example.com"
	expectedErr := errors.New("user not found")

	suite.mockRepo.EXPECT().FindByEmail(email).Return(nil, expectedErr)

	result, err := suite.mockRepo.FindByEmail(email)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(expectedErr, err)
}

func (suite *UserRepositorySuite) TestCreateUser_Success() {
	req := &requests.CreateUserRequest{FirstName: "New", LastName: "User", Email: "new@example.com", Password: "password123", ConfirmPassword: "password123"}
	createdUser := &record.UserRecord{ID: 10, FirstName: "New", LastName: "User", Email: "new@example.com"}

	suite.mockRepo.EXPECT().CreateUser(req).Return(createdUser, nil)

	result, err := suite.mockRepo.CreateUser(req)

	suite.NoError(err)
	suite.Equal(createdUser, result)
}

func (suite *UserRepositorySuite) TestUpdateUser_Success() {
	userID := 1
	req := &requests.UpdateUserRequest{UserID: &userID, FirstName: "Updated", LastName: "User", Email: "updated@example.com", Password: "newpass", ConfirmPassword: "newpass"}
	updatedUser := &record.UserRecord{ID: userID, FirstName: "Updated", LastName: "User", Email: "updated@example.com"}

	suite.mockRepo.EXPECT().UpdateUser(req).Return(updatedUser, nil)

	result, err := suite.mockRepo.UpdateUser(req)

	suite.NoError(err)
	suite.Equal(updatedUser, result)
}

func (suite *UserRepositorySuite) TestTrashedUser_Success() {
	userID := 1
	trashedTime := "2024-01-01T00:00:00Z"
	trashedUser := &record.UserRecord{ID: userID, DeletedAt: &trashedTime}

	suite.mockRepo.EXPECT().TrashedUser(userID).Return(trashedUser, nil)

	result, err := suite.mockRepo.TrashedUser(userID)

	suite.NoError(err)
	suite.NotNil(result.DeletedAt)
	suite.Equal(trashedUser, result)
}

func (suite *UserRepositorySuite) TestRestoreUser_Success() {
	userID := 1
	restoredUser := &record.UserRecord{ID: userID, DeletedAt: nil}

	suite.mockRepo.EXPECT().RestoreUser(userID).Return(restoredUser, nil)

	result, err := suite.mockRepo.RestoreUser(userID)

	suite.NoError(err)
	suite.Nil(result.DeletedAt)
	suite.Equal(restoredUser, result)
}

func (suite *UserRepositorySuite) TestDeleteUserPermanent_Success() {
	userID := 1

	suite.mockRepo.EXPECT().DeleteUserPermanent(userID).Return(true, nil)

	result, err := suite.mockRepo.DeleteUserPermanent(userID)

	suite.NoError(err)
	suite.True(result)
}

func (suite *UserRepositorySuite) TestRestoreAllUser_Success() {
	suite.mockRepo.EXPECT().RestoreAllUser().Return(true, nil)

	result, err := suite.mockRepo.RestoreAllUser()

	suite.NoError(err)
	suite.True(result)
}

func (suite *UserRepositorySuite) TestDeleteAllUserPermanent_Success() {
	suite.mockRepo.EXPECT().DeleteAllUserPermanent().Return(true, nil)

	result, err := suite.mockRepo.DeleteAllUserPermanent()

	suite.NoError(err)
	suite.True(result)
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}
