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

type SaldoHandleGrpcTestSuite struct {
	suite.Suite
	Ctrl             *gomock.Controller
	MockSaldoService *mock_service.MockSaldoService
	MockProtoMapper  *mock_protomapper.MockSaldoProtoMapper
	Handler          gapi.SaldoHandleGrpc
}

func (suite *SaldoHandleGrpcTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockSaldoService = mock_service.NewMockSaldoService(suite.Ctrl)
	suite.MockProtoMapper = mock_protomapper.NewMockSaldoProtoMapper(suite.Ctrl)
	suite.Handler = gapi.NewSaldoHandleGrpc(suite.MockSaldoService, suite.MockProtoMapper)
}

func (suite *SaldoHandleGrpcTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *SaldoHandleGrpcTestSuite) TestFindAllSaldo_Success() {
	req := &pb.FindAllSaldoRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockSaldos := []*response.SaldoResponse{
		{ID: 1, CardNumber: "1234567890123456", TotalBalance: 100000},
		{ID: 2, CardNumber: "6543210987654321", TotalBalance: 200000},
	}
	mockProtoSaldos := []*pb.SaldoResponse{
		{SaldoId: 1, CardNumber: "1234567890123456", TotalBalance: 100000},
		{SaldoId: 2, CardNumber: "6543210987654321", TotalBalance: 200000},
	}

	totalRecords := 2
	suite.MockSaldoService.EXPECT().
		FindAll(gomock.Eq(&requests.FindAllSaldos{Page: 1, PageSize: 10, Search: "test"})).
		Return(mockSaldos, &totalRecords, nil)

	suite.MockProtoMapper.EXPECT().
		ToProtoResponsePaginationSaldo(gomock.Any(), "success", "Successfully fetched saldo record", mockSaldos).
		Return(&pb.ApiResponsePaginationSaldo{
			Status:  "success",
			Message: "Successfully fetched saldo record",
			Data:    mockProtoSaldos,
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   1,
				TotalRecords: 2,
			},
		})

	res, err := suite.Handler.FindAllSaldo(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Successfully fetched saldo record", res.GetMessage())
	suite.Equal(int32(2), res.GetPagination().GetTotalRecords())
	suite.Equal(2, len(res.GetData()))
}

func (suite *SaldoHandleGrpcTestSuite) TestFindAllSaldo_Failure() {
	req := &pb.FindAllSaldoRequest{Page: 1, PageSize: 10, Search: "test"}
	serviceError := &response.ErrorResponse{Status: "error", Message: "Failed to fetch saldos"}

	totalRecords := 0
	suite.MockSaldoService.EXPECT().FindAll(gomock.Any()).Return(nil, &totalRecords, serviceError)

	res, _ := suite.Handler.FindAllSaldo(context.Background(), req)

	suite.Nil(res)
}

func (suite *SaldoHandleGrpcTestSuite) TestFindByIdSaldo_Success() {
	req := &pb.FindByIdSaldoRequest{SaldoId: 1}
	mockSaldo := &response.SaldoResponse{ID: 1, CardNumber: "1234567890123456", TotalBalance: 100000}
	mockProtoSaldo := &pb.SaldoResponse{SaldoId: 1, CardNumber: "1234567890123456", TotalBalance: 100000}

	suite.MockSaldoService.EXPECT().FindById(1).Return(mockSaldo, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseSaldo("success", "Successfully fetched saldo record", mockSaldo).Return(&pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Successfully fetched saldo record",
		Data:    mockProtoSaldo,
	})

	res, err := suite.Handler.FindByIdSaldo(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("1234567890123456", res.GetData().GetCardNumber())
}

func (suite *SaldoHandleGrpcTestSuite) TestFindByIdSaldo_InvalidId() {
	req := &pb.FindByIdSaldoRequest{SaldoId: 0}

	res, err := suite.Handler.FindByIdSaldo(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid Saldo ID")
}

func (suite *SaldoHandleGrpcTestSuite) TestFindByCardNumber_Success() {
	req := &pb.FindByCardNumberRequest{CardNumber: "1234567890123456"}
	mockSaldo := &response.SaldoResponse{ID: 1, CardNumber: "1234567890123456", TotalBalance: 100000}
	mockProtoSaldo := &pb.SaldoResponse{SaldoId: 1, CardNumber: "1234567890123456", TotalBalance: 100000}

	suite.MockSaldoService.EXPECT().FindByCardNumber("1234567890123456").Return(mockSaldo, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseSaldo("success", "Successfully fetched saldo record", mockSaldo).Return(&pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Successfully fetched saldo record",
		Data:    mockProtoSaldo,
	})

	res, err := suite.Handler.FindByCardNumber(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("1234567890123456", res.GetData().GetCardNumber())
}

func (suite *SaldoHandleGrpcTestSuite) TestFindByActiveSaldo_Success() {
	req := &pb.FindAllSaldoRequest{Page: 1, PageSize: 10, Search: ""}
	activeSaldos := []*response.SaldoResponseDeleteAt{
		{ID: 1, CardNumber: "1234567890123456", TotalBalance: 100000},
	}
	totalRecords := 1
	suite.MockSaldoService.EXPECT().FindByActive(gomock.Any()).Return(activeSaldos, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationSaldoDeleteAt(gomock.Any(), "success", "Successfully fetched saldo record", activeSaldos).Return(&pb.ApiResponsePaginationSaldoDeleteAt{
		Status:  "success",
		Message: "Successfully fetched saldo record",
		Data: []*pb.SaldoResponseDeleteAt{
			{SaldoId: 1, CardNumber: "1234567890123456", TotalBalance: 100000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByActive(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *SaldoHandleGrpcTestSuite) TestFindByTrashedSaldo_Success() {
	req := &pb.FindAllSaldoRequest{Page: 1, PageSize: 10, Search: ""}
	trashedSaldos := []*response.SaldoResponseDeleteAt{
		{ID: 1, CardNumber: "1234567890123456", TotalBalance: 100000},
	}
	totalRecords := 1
	suite.MockSaldoService.EXPECT().FindByTrashed(gomock.Any()).Return(trashedSaldos, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationSaldoDeleteAt(gomock.Any(), "success", "Successfully fetched saldo record", trashedSaldos).Return(&pb.ApiResponsePaginationSaldoDeleteAt{
		Status:  "success",
		Message: "Successfully fetched saldo record",
		Data: []*pb.SaldoResponseDeleteAt{
			{SaldoId: 1, CardNumber: "1234567890123456", TotalBalance: 100000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByTrashed(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *SaldoHandleGrpcTestSuite) TestCreateSaldo_Success() {
	req := &pb.CreateSaldoRequest{
		CardNumber:   "1234567890123456",
		TotalBalance: 150000,
	}

	mockSaldo := &response.SaldoResponse{
		ID:           1,
		CardNumber:   "1234567890123456",
		TotalBalance: 150000,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	mockProtoSaldo := &pb.SaldoResponse{
		SaldoId:      1,
		CardNumber:   "1234567890123456",
		TotalBalance: 150000,
	}

	suite.MockSaldoService.EXPECT().CreateSaldo(gomock.Any()).Return(mockSaldo, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseSaldo("success", "Successfully created saldo record", mockSaldo).Return(&pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Successfully created saldo record",
		Data:    mockProtoSaldo,
	})

	res, err := suite.Handler.CreateSaldo(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal(int32(150000), res.GetData().GetTotalBalance())
}

func (suite *SaldoHandleGrpcTestSuite) TestCreateSaldo_ValidationError() {
	req := &pb.CreateSaldoRequest{
		CardNumber:   "",
		TotalBalance: 0,
	}

	res, err := suite.Handler.CreateSaldo(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid input for create saldo")
}

func (suite *SaldoHandleGrpcTestSuite) TestUpdateSaldo_Success() {
	saldoId := 1
	req := &pb.UpdateSaldoRequest{
		SaldoId:      int32(saldoId),
		CardNumber:   "1234567890123456",
		TotalBalance: 200000,
	}

	mockSaldo := &response.SaldoResponse{
		ID:           1,
		CardNumber:   "1234567890123456",
		TotalBalance: 200000,
	}

	mockProtoSaldo := &pb.SaldoResponse{SaldoId: 1, TotalBalance: 200000}

	suite.MockSaldoService.EXPECT().UpdateSaldo(gomock.Any()).Return(mockSaldo, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseSaldo("success", "Successfully updated saldo record", mockSaldo).Return(&pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Successfully updated saldo record",
		Data:    mockProtoSaldo,
	})

	res, err := suite.Handler.UpdateSaldo(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal(int32(200000), res.GetData().GetTotalBalance())
}

func (suite *SaldoHandleGrpcTestSuite) TestTrashedSaldo_Success() {
	req := &pb.FindByIdSaldoRequest{SaldoId: 1}
	mockSaldo := &response.SaldoResponseDeleteAt{ID: 1, CardNumber: "1234567890123456"}

	suite.MockSaldoService.EXPECT().TrashSaldo(1).Return(mockSaldo, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseSaldoDeleteAt("success", "Successfully trashed saldo record", mockSaldo).Return(&pb.ApiResponseSaldoDeleteAt{
		Status:  "success",
		Message: "Successfully trashed saldo record",
		Data:    &pb.SaldoResponseDeleteAt{SaldoId: 1, CardNumber: "1234567890123456"},
	})

	res, err := suite.Handler.TrashedSaldo(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *SaldoHandleGrpcTestSuite) TestRestoreSaldo_Success() {
	req := &pb.FindByIdSaldoRequest{SaldoId: 1}
	mockSaldo := &response.SaldoResponseDeleteAt{ID: 1, CardNumber: "1234567890123456"}

	suite.MockSaldoService.EXPECT().RestoreSaldo(1).Return(mockSaldo, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseSaldoDeleteAt("success", "Successfully restored saldo record", mockSaldo).Return(&pb.ApiResponseSaldoDeleteAt{
		Status:  "success",
		Message: "Successfully restored saldo record",
		Data:    &pb.SaldoResponseDeleteAt{SaldoId: 1, CardNumber: "1234567890123456"},
	})

	res, err := suite.Handler.RestoreSaldo(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *SaldoHandleGrpcTestSuite) TestDeleteSaldo_Success() {
	req := &pb.FindByIdSaldoRequest{SaldoId: 1}

	suite.MockSaldoService.EXPECT().DeleteSaldoPermanent(1).Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseSaldoDelete("success", "Successfully deleted saldo record").Return(&pb.ApiResponseSaldoDelete{
		Status:  "success",
		Message: "Successfully deleted saldo record",
	})

	res, err := suite.Handler.DeleteSaldo(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *SaldoHandleGrpcTestSuite) TestRestoreAllSaldo_Success() {
	req := &emptypb.Empty{}

	suite.MockSaldoService.EXPECT().RestoreAllSaldo().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseSaldoAll("success", "Successfully restore all saldo").Return(&pb.ApiResponseSaldoAll{
		Status:  "success",
		Message: "Successfully restore all saldo",
	})

	res, err := suite.Handler.RestoreAllSaldo(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *SaldoHandleGrpcTestSuite) TestDeleteAllSaldoPermanent_Success() {
	req := &emptypb.Empty{}

	suite.MockSaldoService.EXPECT().DeleteAllSaldoPermanent().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseSaldoAll("success", "delete saldo permanent").Return(&pb.ApiResponseSaldoAll{
		Status:  "success",
		Message: "delete saldo permanent",
	})

	res, err := suite.Handler.DeleteAllSaldoPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func TestSaldoHandleGrpcSuite(t *testing.T) {
	suite.Run(t, new(SaldoHandleGrpcTestSuite))
}
