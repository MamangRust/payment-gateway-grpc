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

type MerchantHandleGrpcTestSuite struct {
	suite.Suite
	Ctrl                *gomock.Controller
	MockMerchantService *mock_service.MockMerchantService
	MockProtoMapper     *mock_protomapper.MockMerchantProtoMapper
	Handler             gapi.MerchantHandleGrpc
}

func (suite *MerchantHandleGrpcTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockMerchantService = mock_service.NewMockMerchantService(suite.Ctrl)
	suite.MockProtoMapper = mock_protomapper.NewMockMerchantProtoMapper(suite.Ctrl)
	suite.Handler = gapi.NewMerchantHandleGrpc(suite.MockMerchantService, suite.MockProtoMapper)
}

func (suite *MerchantHandleGrpcTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *MerchantHandleGrpcTestSuite) TestFindAllMerchant_Success() {
	req := &pb.FindAllMerchantRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockMerchants := []*response.MerchantResponse{
		{ID: 1, Name: "Merchant 1", UserID: 1, ApiKey: "API123", Status: "active"},
		{ID: 2, Name: "Merchant 2", UserID: 2, ApiKey: "API456", Status: "active"},
	}
	mockProtoMerchants := []*pb.MerchantResponse{
		{Id: 1, Name: "Merchant 1", UserId: 1, ApiKey: "API123", Status: "active"},
		{Id: 2, Name: "Merchant 2", UserId: 2, ApiKey: "API456", Status: "active"},
	}

	totalRecords := 2
	suite.MockMerchantService.EXPECT().
		FindAll(gomock.Eq(&requests.FindAllMerchants{Page: 1, PageSize: 10, Search: "test"})).
		Return(mockMerchants, &totalRecords, nil)

	suite.MockProtoMapper.EXPECT().
		ToProtoResponsePaginationMerchant(gomock.Any(), "success", "Successfully fetched merchant record", mockMerchants).
		Return(&pb.ApiResponsePaginationMerchant{
			Status:  "success",
			Message: "Successfully fetched merchant record",
			Data:    mockProtoMerchants,
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   1,
				TotalRecords: 2,
			},
		})

	res, err := suite.Handler.FindAllMerchant(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Successfully fetched merchant record", res.GetMessage())
	suite.Equal(int32(2), res.GetPagination().GetTotalRecords())
	suite.Equal(2, len(res.GetData()))
}

func (suite *MerchantHandleGrpcTestSuite) TestFindAllMerchant_Failure() {
	req := &pb.FindAllMerchantRequest{Page: 1, PageSize: 10, Search: "test"}
	serviceError := &response.ErrorResponse{Status: "error", Message: "Failed to fetch merchants"}

	totalRecords := 0
	suite.MockMerchantService.EXPECT().FindAll(gomock.Any()).Return(nil, &totalRecords, serviceError)

	res, _ := suite.Handler.FindAllMerchant(context.Background(), req)

	suite.Nil(res)
}

func (suite *MerchantHandleGrpcTestSuite) TestFindByIdMerchant_Success() {
	req := &pb.FindByIdMerchantRequest{MerchantId: 1}
	mockMerchant := &response.MerchantResponse{ID: 1, Name: "Merchant 1", UserID: 1, ApiKey: "API123", Status: "active"}
	mockProtoMerchant := &pb.MerchantResponse{Id: 1, Name: "Merchant 1", UserId: 1, ApiKey: "API123", Status: "active"}

	suite.MockMerchantService.EXPECT().FindById(1).Return(mockMerchant, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseMerchant("success", "Successfully fetched merchant record", mockMerchant).Return(&pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully fetched merchant record",
		Data:    mockProtoMerchant,
	})

	res, err := suite.Handler.FindByIdMerchant(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Merchant 1", res.GetData().GetName())
}

func (suite *MerchantHandleGrpcTestSuite) TestFindByIdMerchant_InvalidId() {
	req := &pb.FindByIdMerchantRequest{MerchantId: 0}

	res, err := suite.Handler.FindByIdMerchant(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Failed to fetch merchant by ID")
}

func (suite *MerchantHandleGrpcTestSuite) TestFindByUserIdMerchant_Success() {
	req := &pb.FindByMerchantUserIdRequest{UserId: 1}
	mockMerchants := []*response.MerchantResponse{
		{
			ID:     1,
			Name:   "Merchant 1",
			UserID: 1,
			ApiKey: "API123",
			Status: "active",
		},
	}

	mockProtoMerchants := []*pb.MerchantResponse{
		{
			Id:     1,
			Name:   "Merchant 1",
			UserId: 1,
			ApiKey: "API123",
			Status: "active",
		},
	}

	suite.MockMerchantService.EXPECT().FindByMerchantUserId(1).Return(mockMerchants, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseMerchants("success", "Successfully fetched merchant record", mockMerchants).Return(&pb.ApiResponsesMerchant{
		Status:  "success",
		Message: "Successfully fetched merchant record",
		Data:    mockProtoMerchants,
	})

	res, err := suite.Handler.FindByMerchantUserId(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *MerchantHandleGrpcTestSuite) TestFindByApiKeyMerchant_Success() {
	req := &pb.FindByApiKeyRequest{ApiKey: "API123"}
	mockMerchant := &response.MerchantResponse{ID: 1, Name: "Merchant 1", UserID: 1, ApiKey: "API123", Status: "active"}
	mockProtoMerchant := &pb.MerchantResponse{Id: 1, Name: "Merchant 1", UserId: 1, ApiKey: "API123", Status: "active"}

	suite.MockMerchantService.EXPECT().FindByApiKey("API123").Return(mockMerchant, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseMerchant("success", "Successfully fetched merchant record", mockMerchant).Return(&pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully fetched merchant record",
		Data:    mockProtoMerchant,
	})

	res, err := suite.Handler.FindByApiKey(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Merchant 1", res.GetData().GetName())
}

func (suite *MerchantHandleGrpcTestSuite) TestFindByActiveMerchant_Success() {
	req := &pb.FindAllMerchantRequest{Page: 1, PageSize: 10, Search: ""}
	activeMerchants := []*response.MerchantResponseDeleteAt{
		{ID: 1, Name: "Merchant 1", UserID: 1, ApiKey: "API123", Status: "active"},
	}
	totalRecords := 1
	suite.MockMerchantService.EXPECT().FindByActive(gomock.Any()).Return(activeMerchants, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationMerchantDeleteAt(gomock.Any(), "success", "Successfully fetched merchant record", activeMerchants).Return(&pb.ApiResponsePaginationMerchantDeleteAt{
		Status:  "success",
		Message: "Successfully fetched merchant record",
		Data: []*pb.MerchantResponseDeleteAt{
			{Id: 1, Name: "Merchant 1", UserId: 1, ApiKey: "API123", Status: "active"},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByActive(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *MerchantHandleGrpcTestSuite) TestFindByTrashedMerchant_Success() {
	req := &pb.FindAllMerchantRequest{Page: 1, PageSize: 10, Search: ""}
	trashedMerchants := []*response.MerchantResponseDeleteAt{
		{ID: 1, Name: "Merchant 1", UserID: 1, ApiKey: "API123", Status: "active"},
	}
	totalRecords := 1
	suite.MockMerchantService.EXPECT().FindByTrashed(gomock.Any()).Return(trashedMerchants, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationMerchantDeleteAt(gomock.Any(), "success", "Successfully fetched merchant record", trashedMerchants).Return(&pb.ApiResponsePaginationMerchantDeleteAt{
		Status:  "success",
		Message: "Successfully fetched merchant record",
		Data: []*pb.MerchantResponseDeleteAt{
			{Id: 1, Name: "Merchant 1", UserId: 1, ApiKey: "API123", Status: "active"},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByTrashed(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *MerchantHandleGrpcTestSuite) TestCreateMerchant_Success() {
	req := &pb.CreateMerchantRequest{
		Name:   "New Merchant",
		UserId: 1,
	}

	mockMerchant := &response.MerchantResponse{
		ID:        1,
		Name:      "New Merchant",
		UserID:    1,
		ApiKey:    "API123",
		Status:    "active",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	mockProtoMerchant := &pb.MerchantResponse{
		Id:     1,
		Name:   "New Merchant",
		UserId: 1,
		ApiKey: "API123",
		Status: "active",
	}

	suite.MockMerchantService.EXPECT().CreateMerchant(gomock.Any()).Return(mockMerchant, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseMerchant("success", "Successfully created merchant", mockMerchant).Return(&pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully created merchant",
		Data:    mockProtoMerchant,
	})

	res, err := suite.Handler.CreateMerchant(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("New Merchant", res.GetData().GetName())
}

func (suite *MerchantHandleGrpcTestSuite) TestCreateMerchant_ValidationError() {
	req := &pb.CreateMerchantRequest{
		Name:   "",
		UserId: 0,
	}

	res, err := suite.Handler.CreateMerchant(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid input for create merchant")
}

func (suite *MerchantHandleGrpcTestSuite) TestUpdateMerchant_Success() {
	merchantId := 1
	req := &pb.UpdateMerchantRequest{
		MerchantId: int32(merchantId),
		Name:       "Updated Merchant",
		UserId:     1,
		Status:     "active",
	}

	mockMerchant := &response.MerchantResponse{
		ID:     1,
		Name:   "Updated Merchant",
		UserID: 1,
		Status: "active",
	}

	mockProtoMerchant := &pb.MerchantResponse{Id: 1, Name: "Updated Merchant", Status: "active"}

	suite.MockMerchantService.EXPECT().UpdateMerchant(gomock.Any()).Return(mockMerchant, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseMerchant("success", "Successfully updated merchant", mockMerchant).Return(&pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully updated merchant",
		Data:    mockProtoMerchant,
	})

	res, err := suite.Handler.UpdateMerchant(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Updated Merchant", res.GetData().GetName())
}

func (suite *MerchantHandleGrpcTestSuite) TestTrashedMerchant_Success() {
	req := &pb.FindByIdMerchantRequest{MerchantId: 1}
	mockMerchant := &response.MerchantResponseDeleteAt{ID: 1, Name: "Merchant 1"}

	suite.MockMerchantService.EXPECT().TrashedMerchant(1).Return(mockMerchant, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseMerchantDeleteAt("success", "Successfully trashed merchant", mockMerchant).Return(&pb.ApiResponseMerchantDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant",
		Data:    &pb.MerchantResponseDeleteAt{Id: 1, Name: "Merchant 1"},
	})

	res, err := suite.Handler.TrashedMerchant(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *MerchantHandleGrpcTestSuite) TestRestoreMerchant_Success() {
	req := &pb.FindByIdMerchantRequest{MerchantId: 1}
	mockMerchant := &response.MerchantResponseDeleteAt{ID: 1, Name: "Merchant 1"}

	suite.MockMerchantService.EXPECT().RestoreMerchant(1).Return(mockMerchant, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseMerchantDeleteAt("success", "Successfully restored merchant", mockMerchant).Return(&pb.ApiResponseMerchantDeleteAt{
		Status:  "success",
		Message: "Successfully restored merchant",
		Data:    &pb.MerchantResponseDeleteAt{Id: 1, Name: "Merchant 1"},
	})

	res, err := suite.Handler.RestoreMerchant(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *MerchantHandleGrpcTestSuite) TestDeleteMerchantPermanent_Success() {
	req := &pb.FindByIdMerchantRequest{MerchantId: 1}

	suite.MockMerchantService.EXPECT().DeleteMerchantPermanent(1).Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseMerchantDelete("success", "Successfully deleted merchant").Return(&pb.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant",
	})

	res, err := suite.Handler.DeleteMerchantPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *MerchantHandleGrpcTestSuite) TestRestoreAllMerchant_Success() {
	req := &emptypb.Empty{}

	suite.MockMerchantService.EXPECT().RestoreAllMerchant().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseMerchantAll("success", "Successfully restore all merchant").Return(&pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restore all merchant",
	})

	res, err := suite.Handler.RestoreAllMerchant(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *MerchantHandleGrpcTestSuite) TestDeleteAllMerchantPermanent_Success() {
	req := &emptypb.Empty{}

	suite.MockMerchantService.EXPECT().DeleteAllMerchantPermanent().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseMerchantAll("success", "Successfully delete all merchant").Return(&pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully delete all merchant",
	})

	res, err := suite.Handler.DeleteAllMerchantPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func TestMerchantHandleGrpcSuite(t *testing.T) {
	suite.Run(t, new(MerchantHandleGrpcTestSuite))
}
