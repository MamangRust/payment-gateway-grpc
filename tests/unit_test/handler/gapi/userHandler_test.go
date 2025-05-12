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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFindAllUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindAllUserRequest{
		Page:     1,
		PageSize: 10,
		Search:   "John",
	}

	mockUsers := []*response.UserResponse{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane@example.com",
		},
	}
	mockTotalRecords := 20

	mockUserService.EXPECT().FindAll(1, 10, "John").Return(mockUsers, mockTotalRecords, nil)
	mockProtoMapper.EXPECT().ToResponsesUser(mockUsers).Return([]*pb.UserResponse{
		{
			Id:        1,
			Firstname: "John",
			Lastname:  "Doe",
			Email:     "john@example.com",
		},
		{
			Id:        2,
			Firstname: "Jane",
			Lastname:  "Smith",
			Email:     "jane@example.com",
		},
	})

	res, err := userHandler.FindAll(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched users", res.GetMessage())
	assert.Equal(t, int32(20), res.GetPagination().TotalRecords)
	assert.Equal(t, int32(2), int32(len(res.GetData())))
}

func TestFindAllUsers_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindAllUserRequest{
		Page:     1,
		PageSize: 10,
		Search:   "John",
	}

	mockError := &response.ErrorResponse{
		Status:  "error",
		Message: "Database connection failed",
	}

	mockUserService.EXPECT().FindAll(1, 10, "John").Return(nil, 0, mockError)

	res, err := userHandler.FindAll(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to fetch users")
}

func TestFindAllUsers_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindAllUserRequest{
		Page:     1,
		PageSize: 10,
		Search:   "Nonexistent",
	}

	mockUserService.EXPECT().FindAll(1, 10, "Nonexistent").Return([]*response.UserResponse{}, 0, nil)
	mockProtoMapper.EXPECT().ToResponsesUser([]*response.UserResponse{}).Return([]*pb.UserResponse{})

	res, err := userHandler.FindAll(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched users", res.GetMessage())
	assert.Equal(t, int32(0), res.GetPagination().TotalRecords)
	assert.Equal(t, 0, len(res.GetData()))
}

func TestFindByIdUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindByIdUserRequest{
		Id: 1,
	}

	mockUser := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	mockUserService.EXPECT().FindByID(1).Return(mockUser, nil)
	mockProtoMapper.EXPECT().ToResponseUser(mockUser).Return(&pb.UserResponse{
		Id:        1,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john@example.com",
	})

	res, err := userHandler.FindById(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched user", res.GetMessage())
	assert.Equal(t, int32(1), res.GetData().GetId())
	assert.Equal(t, "John", res.GetData().GetFirstname())
	assert.Equal(t, "Doe", res.GetData().GetLastname())
	assert.Equal(t, "john@example.com", res.GetData().GetEmail())
}

func TestFindByIdUser_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindByIdUserRequest{
		Id: 0,
	}

	res, err := userHandler.FindById(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Invalid user id")
}

func TestFindByIdUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindByIdUserRequest{
		Id: 1,
	}

	mockError := &response.ErrorResponse{
		Status:  "error",
		Message: "User not found",
	}

	mockUserService.EXPECT().FindByID(1).Return(nil, mockError)

	res, err := userHandler.FindById(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to fetch user")
}

func TestFindByActive_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	mockUsers := []*response.UserResponseDeleteAt{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane@example.com",
		},
	}

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockUserService.EXPECT().FindByActive(pageSize, page, search).Return(mockUsers, expected, nil)

	mockProtoMapper.EXPECT().ToResponsesUserDeleteAt(mockUsers).Return([]*pb.UserResponseWithDeleteAt{
		{
			Id:        1,
			Firstname: "John",
			Lastname:  "Doe",
			Email:     "john@example.com",
		},
		{
			Id:        2,
			Firstname: "Jane",
			Lastname:  "Smith",
			Email:     "jane@example.com",
		},
	})

	res, err := userHandler.FindByActive(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched active users", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
}

func TestFindByActiveUser_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	mockUsers := []*response.UserResponseDeleteAt{}

	mockProtoResponses := []*pb.UserResponseWithDeleteAt{}

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockUserService.EXPECT().FindByActive(pageSize, page, search).Return(mockUsers, expected, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponsesUserDeleteAt(mockUsers).Return(mockProtoResponses).Times(1)

	res, err := userHandler.FindByActive(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched active users", res.GetMessage())
	assert.Len(t, res.GetData(), 0)
}

func TestFindByActive_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	mockError := &response.ErrorResponse{
		Status:  "error",
		Message: "Database error",
	}

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockUserService.EXPECT().FindByActive(pageSize, page, search).Return(nil, expected, mockError)

	res, err := userHandler.FindByActive(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to fetch active users")
}

func TestFindByTrashedUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	mockUsers := []*response.UserResponseDeleteAt{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane@example.com",
		},
	}

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockUserService.EXPECT().FindByTrashed(pageSize, page, search).Return(mockUsers, expected, nil)
	mockProtoMapper.EXPECT().ToResponsesUserDeleteAt(mockUsers).Return([]*pb.UserResponse{
		{
			Id:        1,
			Firstname: "John",
			Lastname:  "Doe",
			Email:     "john@example.com",
		},
		{
			Id:        2,
			Firstname: "Jane",
			Lastname:  "Smith",
			Email:     "jane@example.com",
		},
	})

	res, err := userHandler.FindByTrashed(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched trashed users", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
}

func TestFindByTrashedUser_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	mockUsers := []*response.UserResponseDeleteAt{}

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockUserService.EXPECT().FindByTrashed(pageSize, page, search).Return(mockUsers, expected, nil)
	mockProtoMapper.EXPECT().ToResponsesUserDeleteAt(mockUsers).Return([]*pb.UserResponseWithDeleteAt{})

	res, err := userHandler.FindByTrashed(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched trashed users", res.GetMessage())
	assert.Len(t, res.GetData(), 0)
}

func TestFindByTrashedUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockUserService.EXPECT().FindByTrashed(pageSize, page, search).Return(nil, expected, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch trashed users: ",
	})

	res, err := userHandler.FindByTrashed(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
	assert.Contains(t, err.Error(), "Failed to fetch trashed users: ")
}

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.CreateUserRequest{
		Firstname:       "John",
		Lastname:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockCreateUserRequest := &requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockUserResponse := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	mockUserService.EXPECT().CreateUser(mockCreateUserRequest).Return(mockUserResponse, nil)
	mockProtoMapper.EXPECT().ToResponseUser(mockUserResponse).Return(&pb.UserResponse{
		Id:        1,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
	})

	res, err := userHandler.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully created user", res.GetMessage())
	assert.NotNil(t, res.GetData())
	assert.Equal(t, int32(1), res.GetData().GetId())
}

func TestCreateUser_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	reqAuth := &requests.CreateUserRequest{
		FirstName: "",
		LastName:  "Doe",
		Email:     "invalid-email",
		Password:  "123",
	}

	req := &pb.CreateUserRequest{
		Firstname: "",
		Lastname:  "Doe",
		Email:     "invalid-email",
		Password:  "123",
	}

	mockUserService.EXPECT().CreateUser(reqAuth).Times(0)

	res, err := userHandler.Create(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to create user")
}

func TestCreateUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.CreateUserRequest{
		Firstname:       "John",
		Lastname:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockCreateUserRequest := &requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockError := &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to save user to database",
	}

	mockUserService.EXPECT().CreateUser(mockCreateUserRequest).Return(nil, mockError)

	res, err := userHandler.Create(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to create user")
}

func TestUpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.UpdateUserRequest{
		Id:              1,
		Firstname:       "John",
		Lastname:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockUpdateUserRequest := &requests.UpdateUserRequest{
		UserID:          1,
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockUserResponse := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	mockUserService.EXPECT().UpdateUser(mockUpdateUserRequest).Return(mockUserResponse, nil)
	mockProtoMapper.EXPECT().ToResponseUser(mockUserResponse).Return(&pb.UserResponse{
		Id:        1,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
	})

	res, err := userHandler.Update(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully updated user", res.GetMessage())
	assert.NotNil(t, res.GetData())
	assert.Equal(t, int32(1), res.GetData().GetId())
}

func TestUpdateUser_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.UpdateUserRequest{
		Id:              0,
		Firstname:       "John",
		Lastname:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockUserService.EXPECT().UpdateUser(gomock.Any()).Times(0)
	mockProtoMapper.EXPECT().ToResponseUser(gomock.Any()).Times(0)

	res, err := userHandler.Update(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, statusErr.Message(), "Invalid user id")
}

func TestUpdateUser_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.UpdateUserRequest{
		Id:              1,
		Firstname:       "",
		Lastname:        "",
		Email:           "invalid-email@gmail.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	res, err := userHandler.Update(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
}

func TestUpdateUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.UpdateUserRequest{
		Id:              1,
		Firstname:       "John",
		Lastname:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockUpdateUserRequest := &requests.UpdateUserRequest{
		UserID:          1,
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockError := &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to update user in database",
	}

	mockUserService.EXPECT().UpdateUser(mockUpdateUserRequest).Return(nil, mockError)

	res, err := userHandler.Update(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to update user")
}

func TestTrashedUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindByIdUserRequest{Id: 1}
	mockUserResponse := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	mockUserService.EXPECT().TrashedUser(1).Return(mockUserResponse, nil)
	mockProtoMapper.EXPECT().ToResponseUser(mockUserResponse).Return(&pb.UserResponse{
		Id:        1,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
	})

	res, err := userHandler.TrashedUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully trashed user", res.GetMessage())
	assert.NotNil(t, res.GetData())
}

func TestTrashedUser_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindByIdUserRequest{Id: 0}

	res, err := userHandler.TrashedUser(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, statusErr.Message(), "Invalid user id")
}

func TestTrashedUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindByIdUserRequest{Id: 1}
	mockError := &response.ErrorResponse{
		Status:  "error",
		Message: "User not found",
	}

	mockUserService.EXPECT().TrashedUser(1).Return(nil, mockError)

	res, err := userHandler.TrashedUser(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to trashed user")
}

func TestRestoreUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindByIdUserRequest{Id: 1}
	mockUserResponse := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	mockUserService.EXPECT().RestoreUser(1).Return(mockUserResponse, nil)
	mockProtoMapper.EXPECT().ToResponseUser(mockUserResponse).Return(&pb.UserResponse{
		Id:        1,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
	})

	res, err := userHandler.RestoreUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully restored user", res.GetMessage())
	assert.NotNil(t, res.GetData())
}

func TestRestoreUser_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindByIdUserRequest{Id: 0}

	res, err := userHandler.RestoreUser(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, statusErr.Message(), "Invalid user id")
}

func TestRestoreUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindByIdUserRequest{Id: 1}
	mockError := &response.ErrorResponse{
		Status:  "error",
		Message: "User not found",
	}

	mockUserService.EXPECT().RestoreUser(1).Return(nil, mockError)

	res, err := userHandler.RestoreUser(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to restore user")
}

func TestDeleteUserPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindByIdUserRequest{Id: 1}

	mockResponse := &pb.ApiResponseUserDelete{
		Status:  "success",
		Message: "Successfully deleted user permanently",
	}

	mockUserService.EXPECT().DeleteUserPermanent(1).Return(mockResponse, nil)

	res, err := userHandler.DeleteUserPermanent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully deleted user permanently", res.GetMessage())
}

func TestDeleteUserPermanent_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindByIdUserRequest{Id: 0}

	res, err := userHandler.DeleteUserPermanent(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, statusErr.Message(), "Invalid user id")
}

func TestDeleteUserPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockUserProtoMapper(ctrl)
	userHandler := gapi.NewUserHandleGrpc(mockUserService, mockProtoMapper)

	req := &pb.FindByIdUserRequest{Id: 1}

	mockUserService.EXPECT().DeleteUserPermanent(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to delete user permanently",
	})

	res, err := userHandler.DeleteUserPermanent(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to delete user permanently")
}
