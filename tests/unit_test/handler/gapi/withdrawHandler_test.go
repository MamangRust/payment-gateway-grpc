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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WithdrawHandleGrpcTestSuite struct {
	suite.Suite
	Ctrl                *gomock.Controller
	MockWithdrawService *mock_service.MockWithdrawService
	MockProtoMapper     *mock_protomapper.MockWithdrawalProtoMapper
	Handler             gapi.WithdrawHandleGrpc
}

func (suite *WithdrawHandleGrpcTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockWithdrawService = mock_service.NewMockWithdrawService(suite.Ctrl)
	suite.MockProtoMapper = mock_protomapper.NewMockWithdrawalProtoMapper(suite.Ctrl)
	suite.Handler = gapi.NewWithdrawHandleGrpc(suite.MockWithdrawService, suite.MockProtoMapper)
}

func (suite *WithdrawHandleGrpcTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *WithdrawHandleGrpcTestSuite) TestFindAllWithdraw_Success() {
	req := &pb.FindAllWithdrawRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockWithdraws := []*response.WithdrawResponse{
		{ID: 1, WithdrawNo: "WD123", CardNumber: "1234567890123456", WithdrawAmount: 100000},
		{ID: 2, WithdrawNo: "WD124", CardNumber: "6543210987654321", WithdrawAmount: 200000},
	}
	mockProtoWithdraws := []*pb.WithdrawResponse{
		{WithdrawId: 1, WithdrawNo: "WD123", CardNumber: "1234567890123456", WithdrawAmount: 100000},
		{WithdrawId: 2, WithdrawNo: "WD124", CardNumber: "6543210987654321", WithdrawAmount: 200000},
	}

	totalRecords := 2
	suite.MockWithdrawService.EXPECT().
		FindAll(gomock.Eq(&requests.FindAllWithdraws{Page: 1, PageSize: 10, Search: "test"})).
		Return(mockWithdraws, &totalRecords, nil)

	suite.MockProtoMapper.EXPECT().
		ToProtoResponsePaginationWithdraw(gomock.Any(), "success", "withdraw", mockWithdraws).
		Return(&pb.ApiResponsePaginationWithdraw{
			Status:  "success",
			Message: "Successfully fetched withdraw records",
			Data:    mockProtoWithdraws,
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   1,
				TotalRecords: 2,
			},
		})

	res, err := suite.Handler.FindAllWithdraw(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Successfully fetched withdraw records", res.GetMessage())
	suite.Equal(int32(2), res.GetPagination().GetTotalRecords())
	suite.Equal(2, len(res.GetData()))
}

func (suite *WithdrawHandleGrpcTestSuite) TestFindAllWithdraw_Failure() {
	req := &pb.FindAllWithdrawRequest{Page: 1, PageSize: 10, Search: "test"}
	serviceError := &response.ErrorResponse{Status: "error", Message: "Failed to fetch withdraws"}

	totalRecords := 0
	suite.MockWithdrawService.EXPECT().FindAll(gomock.Any()).Return(nil, &totalRecords, serviceError)

	res, _ := suite.Handler.FindAllWithdraw(context.Background(), req)

	suite.Nil(res)
}

func (suite *WithdrawHandleGrpcTestSuite) TestFindByIdWithdraw_Success() {
	req := &pb.FindByIdWithdrawRequest{WithdrawId: 1}
	mockWithdraw := &response.WithdrawResponse{ID: 1, WithdrawNo: "WD123", WithdrawAmount: 100000}
	mockProtoWithdraw := &pb.WithdrawResponse{WithdrawId: 1, WithdrawNo: "WD123", WithdrawAmount: 100000}

	suite.MockWithdrawService.EXPECT().FindById(1).Return(mockWithdraw, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseWithdraw("success", "Successfully fetched withdraw", mockWithdraw).Return(&pb.ApiResponseWithdraw{
		Status:  "success",
		Message: "Successfully fetched withdraw",
		Data:    mockProtoWithdraw,
	})

	res, err := suite.Handler.FindByIdWithdraw(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("WD123", res.GetData().GetWithdrawNo())
}

func (suite *WithdrawHandleGrpcTestSuite) TestFindByIdWithdraw_InvalidId() {
	req := &pb.FindByIdWithdrawRequest{WithdrawId: 0}

	res, err := suite.Handler.FindByIdWithdraw(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid Withdraw ID")
}

func (suite *WithdrawHandleGrpcTestSuite) TestFindByActiveWithdraw_Success() {
	req := &pb.FindAllWithdrawRequest{Page: 1, PageSize: 10, Search: ""}
	activeWithdraws := []*response.WithdrawResponseDeleteAt{
		{ID: 1, WithdrawNo: "WD123", WithdrawAmount: 100000},
	}
	totalRecords := 1
	suite.MockWithdrawService.EXPECT().FindByActive(gomock.Any()).Return(activeWithdraws, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationWithdrawDeleteAt(gomock.Any(), "success", "Successfully fetched withdraws", activeWithdraws).Return(&pb.ApiResponsePaginationWithdrawDeleteAt{
		Status:  "success",
		Message: "Successfully fetched active withdraw records",
		Data: []*pb.WithdrawResponseDeleteAt{
			{WithdrawId: 1, WithdrawNo: "WD123", WithdrawAmount: 100000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByActive(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *WithdrawHandleGrpcTestSuite) TestFindByTrashedWithdraw_Success() {
	req := &pb.FindAllWithdrawRequest{Page: 1, PageSize: 10, Search: ""}
	trashedWithdraws := []*response.WithdrawResponseDeleteAt{
		{ID: 1, WithdrawNo: "WD123", WithdrawAmount: 100000},
	}
	totalRecords := 1
	suite.MockWithdrawService.EXPECT().FindByTrashed(gomock.Any()).Return(trashedWithdraws, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationWithdrawDeleteAt(gomock.Any(), "success", "Successfully fetched withdraws", trashedWithdraws).Return(&pb.ApiResponsePaginationWithdrawDeleteAt{
		Status:  "success",
		Message: "Successfully fetched trashed withdraw records",
		Data: []*pb.WithdrawResponseDeleteAt{
			{WithdrawId: 1, WithdrawNo: "WD123", WithdrawAmount: 100000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByTrashed(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *WithdrawHandleGrpcTestSuite) TestCreateWithdraw_Success() {
	withdrawTime := time.Now()
	req := &pb.CreateWithdrawRequest{
		CardNumber:     "1234567890123456",
		WithdrawAmount: 150000,
		WithdrawTime:   timestamppb.New(withdrawTime),
	}

	mockWithdraw := &response.WithdrawResponse{
		ID:             1,
		WithdrawNo:     "WD125",
		CardNumber:     "1234567890123456",
		WithdrawAmount: 150000,
		WithdrawTime:   withdrawTime.Format(time.RFC3339),
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	mockProtoWithdraw := &pb.WithdrawResponse{
		WithdrawId:     1,
		WithdrawNo:     "WD125",
		CardNumber:     "1234567890123456",
		WithdrawAmount: 150000,
		WithdrawTime:   withdrawTime.Format(time.RFC3339),
	}

	suite.MockWithdrawService.EXPECT().Create(gomock.Any()).Return(mockWithdraw, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseWithdraw("success", "Successfully created withdraw", mockWithdraw).Return(&pb.ApiResponseWithdraw{
		Status:  "success",
		Message: "Successfully created withdraw",
		Data:    mockProtoWithdraw,
	})

	res, err := suite.Handler.CreateWithdraw(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal(int32(150000), res.GetData().GetWithdrawAmount())
}

func (suite *WithdrawHandleGrpcTestSuite) TestCreateWithdraw_ValidationError() {
	req := &pb.CreateWithdrawRequest{
		CardNumber:     "",
		WithdrawAmount: 40000,
		WithdrawTime:   nil,
	}

	res, err := suite.Handler.CreateWithdraw(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid input for create withdraw")
}

func (suite *WithdrawHandleGrpcTestSuite) TestUpdateWithdraw_Success() {
	withdrawId := 1
	withdrawTime := time.Now()
	req := &pb.UpdateWithdrawRequest{
		WithdrawId:     int32(withdrawId),
		CardNumber:     "1234567890123456",
		WithdrawAmount: 200000,
		WithdrawTime:   timestamppb.New(withdrawTime),
	}

	mockWithdraw := &response.WithdrawResponse{
		ID:             1,
		WithdrawNo:     "WD125",
		WithdrawAmount: 200000,
	}

	mockProtoWithdraw := &pb.WithdrawResponse{WithdrawId: 1, WithdrawAmount: 200000}

	suite.MockWithdrawService.EXPECT().Update(gomock.Any()).Return(mockWithdraw, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseWithdraw("success", "Successfully updated withdraw", mockWithdraw).Return(&pb.ApiResponseWithdraw{
		Status:  "success",
		Message: "Successfully updated withdraw",
		Data:    mockProtoWithdraw,
	})

	res, err := suite.Handler.UpdateWithdraw(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal(int32(200000), res.GetData().GetWithdrawAmount())
}

func (suite *WithdrawHandleGrpcTestSuite) TestTrashedWithdraw_Success() {
	req := &pb.FindByIdWithdrawRequest{WithdrawId: 1}
	mockWithdraw := &response.WithdrawResponseDeleteAt{ID: 1, WithdrawNo: "WD125"}

	suite.MockWithdrawService.EXPECT().TrashedWithdraw(1).Return(mockWithdraw, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseWithdrawDeleteAt("success", "Successfully trashed withdraw", mockWithdraw).Return(&pb.ApiResponseWithdrawDeleteAt{
		Status:  "success",
		Data:    &pb.WithdrawResponseDeleteAt{WithdrawId: 1, WithdrawNo: "WD125"},
		Message: "Successfully trashed withdraw",
	})

	res, err := suite.Handler.TrashedWithdraw(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *WithdrawHandleGrpcTestSuite) TestRestoreWithdraw_Success() {
	req := &pb.FindByIdWithdrawRequest{WithdrawId: 1}
	mockWithdraw := &response.WithdrawResponseDeleteAt{ID: 1, WithdrawNo: "WD125"}

	suite.MockWithdrawService.EXPECT().RestoreWithdraw(1).Return(mockWithdraw, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseWithdrawDeleteAt("success", "Successfully restored withdraw", mockWithdraw).Return(&pb.ApiResponseWithdrawDeleteAt{
		Status:  "success",
		Message: "Successfully restored withdraw",
		Data:    &pb.WithdrawResponseDeleteAt{WithdrawId: 1, WithdrawNo: "WD125"},
	})

	res, err := suite.Handler.RestoreWithdraw(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *WithdrawHandleGrpcTestSuite) TestDeleteWithdrawPermanent_Success() {
	req := &pb.FindByIdWithdrawRequest{WithdrawId: 1}

	suite.MockWithdrawService.EXPECT().DeleteWithdrawPermanent(1).Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseWithdrawDelete("success", "Successfully deleted withdraw permanently").Return(&pb.ApiResponseWithdrawDelete{
		Status:  "success",
		Message: "Successfully deleted withdraw permanently",
	})

	res, err := suite.Handler.DeleteWithdrawPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *WithdrawHandleGrpcTestSuite) TestRestoreAllWithdraw_Success() {
	req := &emptypb.Empty{}

	suite.MockWithdrawService.EXPECT().RestoreAllWithdraw().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseWithdrawAll("success", "Successfully restore all withdraw").Return(&pb.ApiResponseWithdrawAll{
		Status:  "success",
		Message: "Successfully restore all withdraw",
	})

	res, err := suite.Handler.RestoreAllWithdraw(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *WithdrawHandleGrpcTestSuite) TestDeleteAllWithdrawPermanent_Success() {
	req := &emptypb.Empty{}

	suite.MockWithdrawService.EXPECT().DeleteAllWithdrawPermanent().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseWithdrawAll("success", "Successfully delete withdraw permanent").Return(&pb.ApiResponseWithdrawAll{
		Status:  "success",
		Message: "Successfully delete withdraw permanent",
	})

	res, err := suite.Handler.DeleteAllWithdrawPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func TestWithdrawHandleGrpcSuite(t *testing.T) {
	suite.Run(t, new(WithdrawHandleGrpcTestSuite))
}
