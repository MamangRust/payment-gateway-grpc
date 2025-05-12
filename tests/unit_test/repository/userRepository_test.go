package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/tests/utils"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindAllUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

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

	mockRepo.EXPECT().FindAllUsers("", 1, 10).Return(users, 3, nil)

	result, total, err := mockRepo.FindAllUsers("", 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, users, result)
	assert.Equal(t, 3, total)
}

func TestFindAllUsers_Error(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().FindAllUsers("", 1, 10).Return(nil, 0, errors.New("database error"))

	result, total, err := mockRepo.FindAllUsers("", 1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, total)
	assert.EqualError(t, err, "database error")
}

func TestFindAllUsers_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().FindAllUsers("", 1, 10).Return([]*record.UserRecord{}, 0, nil)

	result, total, err := mockRepo.FindAllUsers("", 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, []*record.UserRecord{}, result)
	assert.Equal(t, 0, total)
}

func TestFindUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	user := &record.UserRecord{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "M0Vn2@example.com",
		Password:  "password123",
		CreatedAt: "2024-12-21T09:00:00Z",
		UpdatedAt: "2024-12-21T09:00:00Z",
	}

	mockRepo.EXPECT().FindById(1).Return(user, nil)

	result, err := mockRepo.FindById(1)

	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestFindUserByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().FindById(1).Return(nil, errors.New("database error"))

	result, err := mockRepo.FindById(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")
}

func TestFindUserByEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	user := &record.UserRecord{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "M0Vn2@example.com",
		Password:  "password123",
		CreatedAt: "2024-12-21T09:00:00Z",
		UpdatedAt: "2024-12-21T09:00:00Z",
	}

	mockRepo.EXPECT().FindByEmail("M0Vn2@example.com").Return(user, nil)

	result, err := mockRepo.FindByEmail("M0Vn2@example.com")

	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestFindUserByEmail_Error(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().FindByEmail("M0Vn2@example.com").Return(nil, errors.New("database error"))

	result, err := mockRepo.FindByEmail("M0Vn2@example.com")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")
}

func TestFindByActiveUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	activeUsers := []*record.UserRecord{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			CreatedAt: "2023-01-01",
			UpdatedAt: "2023-12-01",
			DeletedAt: nil,
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			CreatedAt: "2023-02-01",
			UpdatedAt: "2023-12-02",
			DeletedAt: nil,
		},
	}

	search := "user1"
	page := 1
	pageSize := 10
	expected := 2

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(activeUsers, 2, nil)

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "John", result[0].FirstName)
	assert.Equal(t, "jane.smith@example.com", result[1].Email)
}

func TestFindByActiveUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	search := "user1"
	page := 1
	pageSize := 10
	expected := 0

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(nil, 0, errors.New("failed to fetch active users"))

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "failed to fetch active users")
}

func TestFindByTrashedUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	trashedUsers := []*record.UserRecord{
		{
			ID:        3,
			FirstName: "Mark",
			LastName:  "Brown",
			Email:     "mark.brown@example.com",
			CreatedAt: "2024-12-21T09:00:00Z",
			UpdatedAt: "2024-12-21T09:00:00Z",
			DeletedAt: utils.PtrString("2024-12-21T09:00:00Z"),
		},
	}

	search := "user1"
	page := 1
	pageSize := 10
	expected := 1

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(trashedUsers, 1, nil)

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "Mark", result[0].FirstName)
	assert.NotNil(t, result[0].DeletedAt)
}

func TestFindByTrashedUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	search := "user1"
	page := 1
	pageSize := 10
	expected := 0

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, 0, errors.New("failed to fetch trashed users"))

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "failed to fetch trashed users")
}

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	request := requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "strongpassword",
		ConfirmPassword: "strongpassword",
	}

	expectedUser := &record.UserRecord{
		ID:        1,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		CreatedAt: "2023-12-24",
		UpdatedAt: "2023-12-24",
	}

	mockRepo.EXPECT().CreateUser(&request).Return(expectedUser, nil)

	result, err := mockRepo.CreateUser(&request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser, result)
}

func TestCreateUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	request := requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "strongpassword",
		ConfirmPassword: "strongpassword",
	}

	mockRepo.EXPECT().CreateUser(&request).Return(nil, errors.New("failed to create user"))

	result, err := mockRepo.CreateUser(&request)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "failed to create user")
}

func TestCreateUser_ValidationError(t *testing.T) {
	request := requests.CreateUserRequest{
		FirstName:       "",
		LastName:        "Doe",
		Email:           "invalid-email",
		Password:        "pass",
		ConfirmPassword: "password123",
	}

	err := request.Validate()
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "Field validation for 'FirstName' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'Email' failed on the 'email' tag")
	assert.Contains(t, err.Error(), "Field validation for 'Password' failed on the 'min' tag")
	assert.Contains(t, err.Error(), "Field validation for 'ConfirmPassword' failed on the 'eqfield' tag")
}

func TestUpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	request := requests.UpdateUserRequest{
		UserID:          1,
		FirstName:       "Jane",
		LastName:        "Smith",
		Email:           "jane.smith@example.com",
		Password:        "strongpassword",
		ConfirmPassword: "strongpassword",
	}

	expectedUser := &record.UserRecord{
		ID:              request.UserID,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Email:           request.Email,
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
		CreatedAt:       "2023-12-24",
		UpdatedAt:       "2023-12-24",
	}

	mockRepo.EXPECT().UpdateUser(&request).Return(expectedUser, nil)

	result, err := mockRepo.UpdateUser(&request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser, result)
}

func TestUpdateUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	request := requests.UpdateUserRequest{
		UserID:          1,
		FirstName:       "Jane",
		LastName:        "Smith",
		Email:           "jane.smith@example.com",
		Password:        "strongpassword",
		ConfirmPassword: "strongpassword",
	}

	mockRepo.EXPECT().UpdateUser(&request).Return(nil, errors.New("failed to update user"))

	result, err := mockRepo.UpdateUser(&request)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "failed to update user")
}

func TestUpdateUser_ValidationError(t *testing.T) {
	request := requests.UpdateUserRequest{
		UserID:          1,
		FirstName:       "",
		LastName:        "Doe",
		Email:           "invalid-email",
		Password:        "pass",
		ConfirmPassword: "password123",
	}

	err := request.Validate()
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "Field validation for 'FirstName' failed on the 'required' tag")
}

func TestTrashedUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	userID := 1
	trashedUser := &record.UserRecord{
		ID:              userID,
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "strongpassword",
		ConfirmPassword: "strongpassword",
		CreatedAt:       "2024-12-21T09:00:00Z",
		UpdatedAt:       "2024-12-21T09:00:00Z",
		DeletedAt:       utils.PtrString("2024-12-21T09:00:00Z"),
	}

	mockRepo.EXPECT().TrashedUser(userID).Return(trashedUser, nil)

	result, err := mockRepo.TrashedUser(userID)
	assert.NoError(t, err)
	assert.Equal(t, trashedUser, result)
}

func TestTrashedUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	userID := 1

	mockRepo.EXPECT().TrashedUser(userID).Return(nil, errors.New("failed to fetch trashed user"))

	result, err := mockRepo.TrashedUser(userID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "failed to fetch trashed user")
}

func TestRestoreUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	expectedUser := &record.UserRecord{
		ID:              1,
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
		CreatedAt:       "2024-12-21T09:00:00Z",
		UpdatedAt:       "2024-12-21T09:00:00Z",
		DeletedAt:       nil,
	}

	userID := 1

	mockRepo.EXPECT().RestoreUser(userID).Return(expectedUser, nil)

	res, err := mockRepo.RestoreUser(userID)

	assert.NoError(t, err)

	assert.NotNil(t, res)

	assert.Equal(t, expectedUser, res)
}

func TestRestoreUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	userID := 1

	mockRepo.EXPECT().RestoreUser(userID).Return(nil, errors.New("failed to restore user"))

	res, err := mockRepo.RestoreUser(userID)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.EqualError(t, err, "failed to restore user")
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	userID := 1

	mockRepo.EXPECT().DeleteUserPermanent(userID).Return(nil)

	err := mockRepo.DeleteUserPermanent(userID)

	assert.NoError(t, err)
}

func TestDeleteUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	userID := 1

	mockRepo.EXPECT().DeleteUserPermanent(userID).Return(errors.New("failed to delete user"))

	err := mockRepo.DeleteUserPermanent(userID)

	assert.Error(t, err)
}
