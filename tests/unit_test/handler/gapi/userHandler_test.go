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
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandleGrpcTestSuite struct {
	suite.Suite
	Ctrl            *gomock.Controller
	MockUserService *mock_service.MockUserService
	MockProtoMapper *mock_protomapper.MockUserProtoMapper
	Handler         gapi.UserHandleGrpc
}

func (suite *UserHandleGrpcTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockUserService = mock_service.NewMockUserService(suite.Ctrl)
	suite.MockProtoMapper = mock_protomapper.NewMockUserProtoMapper(suite.Ctrl)
	suite.Handler = gapi.NewUserHandleGrpc(suite.MockUserService, suite.MockProtoMapper)
}

func (suite *UserHandleGrpcTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *UserHandleGrpcTestSuite) TestFindAllUser_Success() {
	req := &pb.FindAllUserRequest{
		Page:     1,
		PageSize: 10,
		Search:   "john",
	}

	mockUsers := []*response.UserResponse{
		{ID: 1, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
		{ID: 2, FirstName: "Jane", LastName: "Doe", Email: "jane.doe@example.com"},
	}
	mockProtoUsers := []*pb.UserResponse{
		{Id: 1, Firstname: "John", Lastname: "Doe", Email: "john.doe@example.com"},
		{Id: 2, Firstname: "Jane", Lastname: "Doe", Email: "jane.doe@example.com"},
	}

	totalRecords := 2
	suite.MockUserService.EXPECT().
		FindAll(gomock.Eq(&requests.FindAllUsers{Page: 1, PageSize: 10, Search: "john"})).
		Return(mockUsers, &totalRecords, nil)

	suite.MockProtoMapper.EXPECT().
		ToProtoResponsePaginationUser(gomock.Any(), "success", "Successfully fetched users", mockUsers).
		Return(&pb.ApiResponsePaginationUser{
			Status:  "success",
			Message: "Successfully fetched users",
			Data:    mockProtoUsers,
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   1,
				TotalRecords: 2,
			},
		})

	res, err := suite.Handler.FindAll(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal(int32(2), res.GetPagination().GetTotalRecords())
	suite.Equal(2, len(res.GetData()))
}

func (suite *UserHandleGrpcTestSuite) TestFindAllUser_Failure() {
	req := &pb.FindAllUserRequest{Page: 1, PageSize: 10, Search: "john"}
	serviceError := &response.ErrorResponse{Status: "error", Message: "Failed to fetch users"}

	totalRecords := 0
	suite.MockUserService.EXPECT().FindAll(gomock.Any()).Return(nil, &totalRecords, serviceError)

	res, _ := suite.Handler.FindAll(context.Background(), req)

	suite.Nil(res)
	// suite.Contains(err.Error(), serviceError.Message)
}

func (suite *UserHandleGrpcTestSuite) TestFindByIdUser_Success() {
	req := &pb.FindByIdUserRequest{Id: 1}
	mockUser := &response.UserResponse{ID: 1, FirstName: "John", Email: "john.doe@example.com"}
	mockProtoUser := &pb.UserResponse{Id: 1, Firstname: "John", Email: "john.doe@example.com"}

	suite.MockUserService.EXPECT().FindByID(1).Return(mockUser, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseUser("success", "Successfully fetched user", mockUser).Return(&pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully fetched user",
		Data:    mockProtoUser,
	})

	res, err := suite.Handler.FindById(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("John", res.GetData().GetFirstname())
}

func (suite *UserHandleGrpcTestSuite) TestFindByIdUser_InvalidId() {
	req := &pb.FindByIdUserRequest{Id: 0}

	res, err := suite.Handler.FindById(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.NotFound, statusErr.Code())
	suite.Contains(statusErr.Message(), "User not found")
}

func (suite *UserHandleGrpcTestSuite) TestFindByActiveUser_Success() {
	req := &pb.FindAllUserRequest{Page: 1, PageSize: 10, Search: ""}
	activeUsers := []*response.UserResponseDeleteAt{
		{ID: 1, FirstName: "John", Email: "john.doe@example.com"},
	}
	totalRecords := 1
	suite.MockUserService.EXPECT().FindByActive(gomock.Any()).Return(activeUsers, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationUserDeleteAt(gomock.Any(), "success", "Successfully fetched active users", activeUsers).Return(&pb.ApiResponsePaginationUserDeleteAt{
		Status:  "success",
		Message: "Successfully fetched active user records",
		Data: []*pb.UserResponseDeleteAt{
			{Id: 1, Firstname: "John", Email: "john.doe@example.com"},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByActive(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *UserHandleGrpcTestSuite) TestFindByTrashedUser_Success() {
	req := &pb.FindAllUserRequest{Page: 1, PageSize: 10, Search: ""}
	trashedUsers := []*response.UserResponseDeleteAt{
		{ID: 1, FirstName: "John", Email: "john.doe@example.com"},
	}
	totalRecords := 1
	suite.MockUserService.EXPECT().FindByTrashed(gomock.Any()).Return(trashedUsers, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationUserDeleteAt(gomock.Any(), "success", "Successfully fetched trashed users", trashedUsers).Return(&pb.ApiResponsePaginationUserDeleteAt{
		Status:  "success",
		Message: "Successfully fetched trashed user records",
		Data: []*pb.UserResponseDeleteAt{
			{Id: 1, Firstname: "John", Email: "john.doe@example.com"},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByTrashed(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *UserHandleGrpcTestSuite) TestCreateUser_Success() {
	req := &pb.CreateUserRequest{
		Firstname:       "John",
		Lastname:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockCreateReq := &requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	mockUser := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	mockProtoUser := &pb.UserResponse{
		Id:        1,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
	}

	suite.MockUserService.EXPECT().CreateUser(mockCreateReq).Return(mockUser, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseUser("success", "Successfully created user", mockUser).Return(&pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully created user",
		Data:    mockProtoUser,
	})

	res, err := suite.Handler.Create(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("John", res.GetData().GetFirstname())
}

func (suite *UserHandleGrpcTestSuite) TestCreateUser_ValidationError() {
	req := &pb.CreateUserRequest{
		Firstname:       "John123",
		Lastname:        "",
		Email:           "not-an-email",
		Password:        "123",
		ConfirmPassword: "456",
	}

	res, err := suite.Handler.Create(context.Background(), req)

	suite.Nil(res)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
}

func (suite *UserHandleGrpcTestSuite) TestUpdateUser_Success() {
	userId := 1
	req := &pb.UpdateUserRequest{
		Id:              int32(userId),
		Firstname:       "John",
		Lastname:        "Smith",
		Email:           "john.smith@example.com",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	mockUpdateReq := &requests.UpdateUserRequest{
		UserID:          &userId,
		FirstName:       "John",
		LastName:        "Smith",
		Email:           "john.smith@example.com",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	mockUser := &response.UserResponse{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john.smith@example.com",
	}

	mockProtoUser := &pb.UserResponse{Id: 1, Firstname: "John", Lastname: "Smith", Email: "john.smith@example.com"}

	suite.MockUserService.EXPECT().UpdateUser(mockUpdateReq).Return(mockUser, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseUser("success", "Successfully updated user", mockUser).Return(&pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully updated user",
		Data:    mockProtoUser,
	})

	res, err := suite.Handler.Update(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Smith", res.GetData().GetLastname())
}

func (suite *UserHandleGrpcTestSuite) TestTrashedUser_Success() {
	req := &pb.FindByIdUserRequest{Id: 1}
	mockUser := &response.UserResponseDeleteAt{ID: 1, FirstName: "John"}

	suite.MockUserService.EXPECT().TrashedUser(1).Return(mockUser, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseUserDeleteAt("success", "Successfully trashed user", mockUser).Return(&pb.ApiResponseUserDeleteAt{
		Status:  "success",
		Message: "Successfully trashed user",
		Data:    &pb.UserResponseDeleteAt{Id: 1, Firstname: "John"},
	})

	res, err := suite.Handler.TrashedUser(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *UserHandleGrpcTestSuite) TestRestoreUser_Success() {
	req := &pb.FindByIdUserRequest{Id: 1}
	mockUser := &response.UserResponseDeleteAt{ID: 1, FirstName: "John"}

	suite.MockUserService.EXPECT().RestoreUser(1).Return(mockUser, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseUserDeleteAt("success", "Successfully restored user", mockUser).Return(&pb.ApiResponseUserDeleteAt{
		Status:  "success",
		Message: "Successfully restored user",
		Data:    &pb.UserResponseDeleteAt{Id: 1, Firstname: "John"},
	})

	res, err := suite.Handler.RestoreUser(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *UserHandleGrpcTestSuite) TestDeleteUserPermanent_Success() {
	req := &pb.FindByIdUserRequest{Id: 1}

	suite.MockUserService.EXPECT().DeleteUserPermanent(1).Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseUserDelete("success", "Successfully deleted user permanently").Return(&pb.ApiResponseUserDelete{
		Status:  "success",
		Message: "Successfully deleted user permanently",
	})

	res, err := suite.Handler.DeleteUserPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *UserHandleGrpcTestSuite) TestRestoreAllUser_Success() {
	req := &emptypb.Empty{}

	suite.MockUserService.EXPECT().RestoreAllUser().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseUserAll("success", "Successfully restored all users").Return(&pb.ApiResponseUserAll{
		Status:  "success",
		Message: "Successfully restored all users",
	})

	res, err := suite.Handler.RestoreAllUser(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *UserHandleGrpcTestSuite) TestDeleteAllUserPermanent_Success() {
	req := &emptypb.Empty{}

	suite.MockUserService.EXPECT().DeleteAllUserPermanent().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseUserAll("success", "Successfully deleted all users permanently").Return(&pb.ApiResponseUserAll{
		Status:  "success",
		Message: "Successfully deleted all users permanently",
	})

	res, err := suite.Handler.DeleteAllUserPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func TestUserHandleGrpcSuite(t *testing.T) {
	suite.Run(t, new(UserHandleGrpcTestSuite))
}
