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

type RoleHandleGrpcTestSuite struct {
	suite.Suite
	Ctrl            *gomock.Controller
	MockRoleService *mock_service.MockRoleService
	MockProtoMapper *mock_protomapper.MockRoleProtoMapper
	Handler         gapi.RoleHandleGrpc
}

func (suite *RoleHandleGrpcTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockRoleService = mock_service.NewMockRoleService(suite.Ctrl)
	suite.MockProtoMapper = mock_protomapper.NewMockRoleProtoMapper(suite.Ctrl)
	suite.Handler = gapi.NewRoleHandleGrpc(suite.MockRoleService, suite.MockProtoMapper)
}

func (suite *RoleHandleGrpcTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *RoleHandleGrpcTestSuite) TestFindAllRole_Success() {
	req := &pb.FindAllRoleRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockRoles := []*response.RoleResponse{
		{ID: 1, Name: "Admin"},
		{ID: 2, Name: "User"},
	}
	mockProtoRoles := []*pb.RoleResponse{
		{Id: 1, Name: "Admin"},
		{Id: 2, Name: "User"},
	}

	totalRecords := 2
	suite.MockRoleService.EXPECT().
		FindAll(gomock.Eq(&requests.FindAllRoles{Page: 1, PageSize: 10, Search: "test"})).
		Return(mockRoles, &totalRecords, nil)

	suite.MockProtoMapper.EXPECT().
		ToProtoResponsePaginationRole(gomock.Any(), "success", "Successfully fetched role records", mockRoles).
		Return(&pb.ApiResponsePaginationRole{
			Status:  "success",
			Message: "Successfully fetched role records",
			Data:    mockProtoRoles,
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   1,
				TotalRecords: 2,
			},
		})

	res, err := suite.Handler.FindAllRole(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Successfully fetched role records", res.GetMessage())
	suite.Equal(int32(2), res.GetPagination().GetTotalRecords())
	suite.Equal(2, len(res.GetData()))
}

func (suite *RoleHandleGrpcTestSuite) TestFindAllRole_Failure() {
	req := &pb.FindAllRoleRequest{Page: 1, PageSize: 10, Search: "test"}
	serviceError := &response.ErrorResponse{Status: "error", Message: "Failed to fetch roles"}

	totalRecords := 0
	suite.MockRoleService.EXPECT().FindAll(gomock.Any()).Return(nil, &totalRecords, serviceError)

	res, _ := suite.Handler.FindAllRole(context.Background(), req)

	suite.Nil(res)
}

func (suite *RoleHandleGrpcTestSuite) TestFindByIdRole_Success() {
	req := &pb.FindByIdRoleRequest{RoleId: 1}
	mockRole := &response.RoleResponse{ID: 1, Name: "Admin"}
	mockProtoRole := &pb.RoleResponse{Id: 1, Name: "Admin"}

	suite.MockRoleService.EXPECT().FindById(1).Return(mockRole, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseRole("success", "Successfully fetched role", mockRole).Return(&pb.ApiResponseRole{
		Status:  "success",
		Message: "Successfully fetched role",
		Data:    mockProtoRole,
	})

	res, err := suite.Handler.FindByIdRole(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Admin", res.GetData().GetName())
}

func (suite *RoleHandleGrpcTestSuite) TestFindByIdRole_InvalidId() {
	req := &pb.FindByIdRoleRequest{RoleId: 0}

	res, err := suite.Handler.FindByIdRole(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.NotFound, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid Role ID")
}

func (suite *RoleHandleGrpcTestSuite) TestFindByUserId_Success() {
	req := &pb.FindByIdUserRoleRequest{UserId: 1}
	mockRoles := []*response.RoleResponse{
		{ID: 1, Name: "Admin"},
		{ID: 2, Name: "User"},
	}
	mockProtoRoles := []*pb.RoleResponse{
		{Id: 1, Name: "Admin"},
		{Id: 2, Name: "User"},
	}

	suite.MockRoleService.EXPECT().FindByUserId(1).Return(mockRoles, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsesRole("success", "Successfully fetched role by user id", mockRoles).Return(&pb.ApiResponsesRole{
		Status:  "success",
		Message: "Successfully fetched role by user id",
		Data:    mockProtoRoles,
	})

	res, err := suite.Handler.FindByUserId(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 2)
}

func (suite *RoleHandleGrpcTestSuite) TestFindByActive_Success() {
	req := &pb.FindAllRoleRequest{Page: 1, PageSize: 10, Search: ""}
	activeRoles := []*response.RoleResponseDeleteAt{
		{ID: 1, Name: "Admin"},
	}
	totalRecords := 1
	suite.MockRoleService.EXPECT().FindByActiveRole(gomock.Any()).Return(activeRoles, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationRoleDeleteAt(gomock.Any(), "success", "Successfully fetched active roles", activeRoles).Return(&pb.ApiResponsePaginationRoleDeleteAt{
		Status:  "success",
		Message: "Successfully fetched active roles",
		Data: []*pb.RoleResponseDeleteAt{
			{Id: 1, Name: "Admin"},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByActive(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *RoleHandleGrpcTestSuite) TestFindByTrashed_Success() {
	req := &pb.FindAllRoleRequest{Page: 1, PageSize: 10, Search: ""}
	trashedRoles := []*response.RoleResponseDeleteAt{
		{ID: 1, Name: "Admin"},
	}
	totalRecords := 1
	suite.MockRoleService.EXPECT().FindByTrashedRole(gomock.Any()).Return(trashedRoles, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationRoleDeleteAt(gomock.Any(), "success", "Successfully fetched trashed roles", trashedRoles).Return(&pb.ApiResponsePaginationRoleDeleteAt{
		Status:  "success",
		Message: "Successfully fetched trashed roles",
		Data: []*pb.RoleResponseDeleteAt{
			{Id: 1, Name: "Admin"},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByTrashed(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *RoleHandleGrpcTestSuite) TestCreateRole_Success() {
	req := &pb.CreateRoleRequest{
		Name: "Manager",
	}

	mockRole := &response.RoleResponse{
		ID:        1,
		Name:      "Manager",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	mockProtoRole := &pb.RoleResponse{
		Id:   1,
		Name: "Manager",
	}

	suite.MockRoleService.EXPECT().CreateRole(gomock.Any()).Return(mockRole, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseRole("success", "Successfully created role", mockRole).Return(&pb.ApiResponseRole{
		Status:  "success",
		Message: "Successfully created role",
		Data:    mockProtoRole,
	})

	res, err := suite.Handler.CreateRole(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Manager", res.GetData().GetName())
}

func (suite *RoleHandleGrpcTestSuite) TestCreateRole_ValidationError() {
	req := &pb.CreateRoleRequest{
		Name: "",
	}

	res, err := suite.Handler.CreateRole(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.Internal, statusErr.Code())
}

func (suite *RoleHandleGrpcTestSuite) TestUpdateRole_Success() {
	roleId := 1
	req := &pb.UpdateRoleRequest{
		Id:   int32(roleId),
		Name: "Updated Role",
	}

	mockRole := &response.RoleResponse{
		ID:   1,
		Name: "Updated Role",
	}

	mockProtoRole := &pb.RoleResponse{Id: 1, Name: "Updated Role"}

	suite.MockRoleService.EXPECT().UpdateRole(gomock.Any()).Return(mockRole, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseRole("success", "Successfully updated role", mockRole).Return(&pb.ApiResponseRole{
		Status:  "success",
		Message: "Successfully updated role",
		Data:    mockProtoRole,
	})

	res, err := suite.Handler.UpdateRole(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Updated Role", res.GetData().GetName())
}

func (suite *RoleHandleGrpcTestSuite) TestTrashedRole_Success() {
	req := &pb.FindByIdRoleRequest{RoleId: 1}
	mockRole := &response.RoleResponseDeleteAt{ID: 1, Name: "Admin"}

	suite.MockRoleService.EXPECT().TrashedRole(1).Return(mockRole, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseRoleDeleteAt("success", "Successfully trashed role", mockRole).Return(&pb.ApiResponseRoleDeleteAt{
		Status:  "success",
		Message: "Successfully trashed role",
		Data:    &pb.RoleResponseDeleteAt{Id: 1, Name: "Admin"},
	})

	res, err := suite.Handler.TrashedRole(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *RoleHandleGrpcTestSuite) TestRestoreRole_Success() {
	req := &pb.FindByIdRoleRequest{RoleId: 1}
	mockRole := &response.RoleResponseDeleteAt{ID: 1, Name: "Admin"}

	suite.MockRoleService.EXPECT().RestoreRole(1).Return(mockRole, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseRoleDeleteAt("success", "Successfully restored role", mockRole).Return(&pb.ApiResponseRoleDeleteAt{
		Status:  "success",
		Message: "Successfully restored role",
		Data:    &pb.RoleResponseDeleteAt{Id: 1, Name: "Admin"},
	})

	res, err := suite.Handler.RestoreRole(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *RoleHandleGrpcTestSuite) TestDeleteRolePermanent_Success() {
	req := &pb.FindByIdRoleRequest{RoleId: 1}

	suite.MockRoleService.EXPECT().DeleteRolePermanent(1).Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseRoleDelete("success", "Successfully deleted role permanently").Return(&pb.ApiResponseRoleDelete{
		Status:  "success",
		Message: "Successfully deleted role permanently",
	})

	res, err := suite.Handler.DeleteRolePermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *RoleHandleGrpcTestSuite) TestRestoreAllRole_Success() {
	req := &emptypb.Empty{}

	suite.MockRoleService.EXPECT().RestoreAllRole().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseRoleAll("success", "Successfully restored all roles").Return(&pb.ApiResponseRoleAll{
		Status:  "success",
		Message: "Successfully restored all roles",
	})

	res, err := suite.Handler.RestoreAllRole(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *RoleHandleGrpcTestSuite) TestDeleteAllRolePermanent_Success() {
	req := &emptypb.Empty{}

	suite.MockRoleService.EXPECT().DeleteAllRolePermanent().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseRoleAll("success", "Successfully deleted all roles").Return(&pb.ApiResponseRoleAll{
		Status:  "success",
		Message: "Successfully deleted all roles",
	})

	res, err := suite.Handler.DeleteAllRolePermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func TestRoleHandleGrpcSuite(t *testing.T) {
	suite.Run(t, new(RoleHandleGrpcTestSuite))
}
