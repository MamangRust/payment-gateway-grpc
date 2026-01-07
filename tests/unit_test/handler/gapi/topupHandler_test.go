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

type TopupHandleGrpcTestSuite struct {
	suite.Suite
	Ctrl             *gomock.Controller
	MockTopupService *mock_service.MockTopupService
	MockProtoMapper  *mock_protomapper.MockTopupProtoMapper
	Handler          gapi.TopupHandleGrpc
}

func (suite *TopupHandleGrpcTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockTopupService = mock_service.NewMockTopupService(suite.Ctrl)
	suite.MockProtoMapper = mock_protomapper.NewMockTopupProtoMapper(suite.Ctrl)
	suite.Handler = gapi.NewTopupHandleGrpc(suite.MockTopupService, suite.MockProtoMapper)
}

func (suite *TopupHandleGrpcTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *TopupHandleGrpcTestSuite) TestFindAllTopup_Success() {
	req := &pb.FindAllTopupRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	totalRecord := 2

	mockTopups := []*response.TopupResponse{
		{ID: 1, CardNumber: "1234567890123456", TopupAmount: 100000},
		{ID: 2, CardNumber: "6543210987654321", TopupAmount: 200000},
	}
	mockProtoTopups := []*pb.TopupResponse{
		{Id: 1, CardNumber: "1234567890123456", TopupAmount: 100000},
		{Id: 2, CardNumber: "6543210987654321", TopupAmount: 200000},
	}

	suite.MockTopupService.EXPECT().
		FindAll(gomock.Eq(&requests.FindAllTopups{Page: 1, PageSize: 10, Search: "test"})).
		Return(mockTopups, &totalRecord, nil)

	suite.MockProtoMapper.EXPECT().
		ToProtoResponsePaginationTopup(gomock.Any(), "success", "Successfully fetch topups", mockTopups).
		Return(&pb.ApiResponsePaginationTopup{
			Status:  "success",
			Message: "Successfully fetch topups",
			Data:    mockProtoTopups,
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   1,
				TotalRecords: 2,
			},
		})

	res, err := suite.Handler.FindAllTopup(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Successfully fetch topups", res.GetMessage())
	suite.Equal(int32(1), res.GetPagination().GetCurrentPage())
	suite.Equal(int32(10), res.GetPagination().GetPageSize())
	suite.Equal(int32(2), res.GetPagination().GetTotalRecords())
	suite.Equal(2, len(res.GetData()))
}

func (suite *TopupHandleGrpcTestSuite) TestFindAllWithdraw_Failure() {
	req := &pb.FindAllTopupRequest{Page: 1, PageSize: 10, Search: "test"}
	serviceError := &response.ErrorResponse{Status: "error", Message: "Failed to fetch topups"}

	totalRecords := 0
	suite.MockTopupService.EXPECT().FindAll(gomock.Any()).Return(nil, &totalRecords, serviceError)

	res, _ := suite.Handler.FindAllTopup(context.Background(), req)

	suite.Nil(res)
}
func (suite *TopupHandleGrpcTestSuite) TestFindAllTopup_Empty() {
	req := &pb.FindAllTopupRequest{
		Page:     1,
		PageSize: 10,
		Search:   "empty",
	}

	mockTopups := []*response.TopupResponse{}
	mockProtoTopups := []*pb.TopupResponse{}
	totalRecords := 0

	suite.MockTopupService.EXPECT().
		FindAll(gomock.Eq(&requests.FindAllTopups{Page: 1, PageSize: 10, Search: "empty"})).
		Return(mockTopups, &totalRecords, nil)

	suite.MockProtoMapper.EXPECT().
		ToProtoResponsePaginationTopup(gomock.Any(), "success", "Successfully fetch topups", mockTopups).
		Return(&pb.ApiResponsePaginationTopup{
			Status:     "success",
			Message:    "Successfully fetch topups",
			Data:       mockProtoTopups,
			Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 0, TotalRecords: 0},
		})

	res, err := suite.Handler.FindAllTopup(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Successfully fetch topups", res.GetMessage())
	suite.Equal(int32(0), res.GetPagination().GetTotalRecords())
	suite.Equal(0, len(res.GetData()))
}

func (suite *TopupHandleGrpcTestSuite) TestFindAllTopupByCardNumber_Success() {
	cardNumber := "1234567890123456"
	req := &pb.FindAllTopupByCardNumberRequest{
		CardNumber: cardNumber,
		Page:       1,
		PageSize:   10,
		Search:     "test",
	}

	totalRecords := 1

	mockTopups := []*response.TopupResponse{
		{ID: 1, CardNumber: cardNumber, TopupAmount: 100000},
	}
	mockProtoTopups := []*pb.TopupResponse{
		{Id: 1, CardNumber: cardNumber, TopupAmount: 100000},
	}

	suite.MockTopupService.EXPECT().
		FindAllByCardNumber(gomock.Eq(&requests.FindAllTopupsByCardNumber{CardNumber: cardNumber, Page: 1, PageSize: 10, Search: "test"})).
		Return(mockTopups, &totalRecords, nil)

	suite.MockProtoMapper.EXPECT().
		ToProtoResponsePaginationTopup(gomock.Any(), "success", "Successfully fetch topups", mockTopups).
		Return(&pb.ApiResponsePaginationTopup{
			Status:  "success",
			Message: "Successfully fetch topups",
			Data:    mockProtoTopups,
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   1,
				TotalRecords: 1,
			},
		})

	res, err := suite.Handler.FindAllTopupByCardNumber(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal(1, len(res.GetData()))
	suite.Equal(cardNumber, res.GetData()[0].GetCardNumber())
}

func (suite *TopupHandleGrpcTestSuite) TestFindByIdTopup_Success() {
	req := &pb.FindByIdTopupRequest{TopupId: 1}
	mockTopup := &response.TopupResponse{ID: 1, CardNumber: "1234567890123456", TopupAmount: 100000}
	mockProtoTopup := &pb.TopupResponse{Id: 1, CardNumber: "1234567890123456", TopupAmount: 100000}

	suite.MockTopupService.EXPECT().FindById(1).Return(mockTopup, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTopup("success", "Successfully fetch topup", mockTopup).Return(&pb.ApiResponseTopup{
		Status:  "success",
		Message: "Successfully fetch topup",
		Data:    mockProtoTopup,
	})

	res, err := suite.Handler.FindByIdTopup(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("1234567890123456", res.GetData().GetCardNumber())
}

func (suite *TopupHandleGrpcTestSuite) TestFindByIdTopup_InvalidId() {
	req := &pb.FindByIdTopupRequest{TopupId: 0}

	res, err := suite.Handler.FindByIdTopup(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid Topup ID")
}

func (suite *TopupHandleGrpcTestSuite) TestFindByActiveTopup_Success() {
	req := &pb.FindAllTopupRequest{Page: 1, PageSize: 10, Search: ""}
	activeTopups := []*response.TopupResponseDeleteAt{
		{ID: 1, CardNumber: "1234567890123456", TopupAmount: 100000},
	}

	totalRecords := 1

	suite.MockTopupService.EXPECT().FindByActive(gomock.Any()).Return(activeTopups, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationTopupDeleteAt(gomock.Any(), "success", "Successfully fetch topups", activeTopups).Return(&pb.ApiResponsePaginationTopupDeleteAt{
		Status:  "success",
		Message: "Successfully fetch topups",
		Data: []*pb.TopupResponseDeleteAt{
			{Id: 1, CardNumber: "1234567890123456", TopupAmount: 100000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByActive(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *TopupHandleGrpcTestSuite) TestFindByTrashedTopup_Success() {
	req := &pb.FindAllTopupRequest{Page: 1, PageSize: 10, Search: ""}
	trashedTopups := []*response.TopupResponseDeleteAt{
		{ID: 1, CardNumber: "1234567890123456", TopupAmount: 100000},
	}

	totalRecords := 1

	suite.MockTopupService.EXPECT().FindByTrashed(gomock.Any()).Return(trashedTopups, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationTopupDeleteAt(gomock.Any(), "success", "Successfully fetch topups", trashedTopups).Return(&pb.ApiResponsePaginationTopupDeleteAt{
		Status:  "success",
		Message: "Successfully fetch topups",
		Data: []*pb.TopupResponseDeleteAt{
			{Id: 1, CardNumber: "1234567890123456", TopupAmount: 100000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByTrashed(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *TopupHandleGrpcTestSuite) TestCreateTopup_Success() {
	req := &pb.CreateTopupRequest{
		CardNumber:  "1234567890123456",
		TopupAmount: 100000,
		TopupMethod: "alfamart",
	}

	mockCreateReq := &requests.CreateTopupRequest{
		CardNumber:  "1234567890123456",
		TopupAmount: 100000,
		TopupMethod: "alfamart",
	}

	mockTopup := &response.TopupResponse{
		ID:          1,
		CardNumber:  "1234567890123456",
		TopupAmount: 100000,
		TopupMethod: "alfamart",
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	mockProtoTopup := &pb.TopupResponse{
		Id:          1,
		CardNumber:  "1234567890123456",
		TopupAmount: 100000,
		TopupMethod: "alfamart",
	}

	suite.MockTopupService.EXPECT().CreateTopup(mockCreateReq).Return(mockTopup, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTopup("success", "Successfully created topup", mockTopup).Return(&pb.ApiResponseTopup{
		Status:  "success",
		Message: "Successfully created topup",
		Data:    mockProtoTopup,
	})

	res, err := suite.Handler.CreateTopup(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal(int32(100000), res.GetData().GetTopupAmount())
}

func (suite *TopupHandleGrpcTestSuite) TestCreateTopup_ValidationError() {
	req := &pb.CreateTopupRequest{
		CardNumber:  "",
		TopupAmount: 10000,
	}

	res, err := suite.Handler.CreateTopup(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid input for create topup")
}

func (suite *TopupHandleGrpcTestSuite) TestUpdateTopup_Success() {
	topupID := 1
	req := &pb.UpdateTopupRequest{
		TopupId:     int32(topupID),
		CardNumber:  "1234567890123456",
		TopupAmount: 150000,
		TopupMethod: "alfamart",
	}

	mockUpdateReq := &requests.UpdateTopupRequest{
		TopupID:     &topupID,
		CardNumber:  "1234567890123456",
		TopupAmount: 150000,
		TopupMethod: "alfamart",
	}

	mockTopup := &response.TopupResponse{
		ID:          1,
		CardNumber:  "1234567890123456",
		TopupAmount: 150000,
		TopupMethod: "alfamart",
	}

	mockProtoTopup := &pb.TopupResponse{Id: 1, CardNumber: "1234567890123456", TopupAmount: 150000, TopupMethod: "alfamart"}

	suite.MockTopupService.EXPECT().UpdateTopup(mockUpdateReq).Return(mockTopup, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTopup("success", "Successfully updated topup", mockTopup).Return(&pb.ApiResponseTopup{
		Status:  "success",
		Message: "Successfully updated topup",
		Data:    mockProtoTopup,
	})

	res, err := suite.Handler.UpdateTopup(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal(int32(150000), res.GetData().GetTopupAmount())
}

func (suite *TopupHandleGrpcTestSuite) TestUpdateTopup_ValidationError() {
	req := &pb.UpdateTopupRequest{
		TopupId:     1,
		CardNumber:  "",
		TopupAmount: 40000,
	}

	res, err := suite.Handler.UpdateTopup(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid input for update topup")
}

func (suite *TopupHandleGrpcTestSuite) TestTrashedTopup_Success() {
	req := &pb.FindByIdTopupRequest{TopupId: 1}
	mockTopup := &response.TopupResponseDeleteAt{ID: 1, CardNumber: "1234567890123456"}

	suite.MockTopupService.EXPECT().TrashedTopup(1).Return(mockTopup, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTopupDeletAt("success", "Successfully trashed topup", mockTopup).Return(&pb.ApiResponseTopupDeleteAt{
		Status:  "success",
		Message: "Successfully trashed topup",
		Data:    &pb.TopupResponseDeleteAt{Id: 1, CardNumber: "1234567890123456"},
	})

	res, err := suite.Handler.TrashedTopup(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *TopupHandleGrpcTestSuite) TestTrashedTopup_InvalidId() {
	req := &pb.FindByIdTopupRequest{TopupId: 0}

	res, err := suite.Handler.TrashedTopup(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
}

func (suite *TopupHandleGrpcTestSuite) TestRestoreTopup_Success() {
	req := &pb.FindByIdTopupRequest{TopupId: 1}
	mockTopup := &response.TopupResponseDeleteAt{ID: 1, CardNumber: "1234567890123456"}

	suite.MockTopupService.EXPECT().RestoreTopup(1).Return(mockTopup, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTopupDeletAt("success", "Successfully restored topup", mockTopup).Return(&pb.ApiResponseTopupDeleteAt{
		Status:  "success",
		Message: "Successfully restored topup",
		Data:    &pb.TopupResponseDeleteAt{Id: 1, CardNumber: "1234567890123456"},
	})

	res, err := suite.Handler.RestoreTopup(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *TopupHandleGrpcTestSuite) TestDeleteTopupPermanent_Success() {
	req := &pb.FindByIdTopupRequest{TopupId: 1}

	suite.MockTopupService.EXPECT().DeleteTopupPermanent(1).Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTopupDelete("success", "Successfully deleted topup permanently").Return(&pb.ApiResponseTopupDelete{
		Status:  "success",
		Message: "Successfully deleted topup permanently",
	})

	res, err := suite.Handler.DeleteTopupPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *TopupHandleGrpcTestSuite) TestRestoreAllTopup_Success() {
	req := &emptypb.Empty{}

	suite.MockTopupService.EXPECT().RestoreAllTopup().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTopupAll("success", "Successfully restore all topup").Return(&pb.ApiResponseTopupAll{
		Status:  "success",
		Message: "Successfully restore all topup",
	})

	res, err := suite.Handler.RestoreAllTopup(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *TopupHandleGrpcTestSuite) TestDeleteAllTopupPermanent_Success() {
	req := &emptypb.Empty{}

	suite.MockTopupService.EXPECT().DeleteAllTopupPermanent().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTopupAll("success", "Successfully delete topup permanent").Return(&pb.ApiResponseTopupAll{
		Status:  "success",
		Message: "Successfully delete topup permanent",
	})

	res, err := suite.Handler.DeleteAllTopupPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}
func TestTopupHandleGrpcSuite(t *testing.T) {
	suite.Run(t, new(TopupHandleGrpcTestSuite))
}
