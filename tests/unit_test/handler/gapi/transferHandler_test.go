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

type TransferHandleGrpcTestSuite struct {
	suite.Suite
	Ctrl                *gomock.Controller
	MockTransferService *mock_service.MockTransferService
	MockProtoMapper     *mock_protomapper.MockTransferProtoMapper
	Handler             gapi.TransferHandleGrpc
}

func (suite *TransferHandleGrpcTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockTransferService = mock_service.NewMockTransferService(suite.Ctrl)
	suite.MockProtoMapper = mock_protomapper.NewMockTransferProtoMapper(suite.Ctrl)
	suite.Handler = gapi.NewTransferHandleGrpc(suite.MockTransferService, suite.MockProtoMapper)
}

func (suite *TransferHandleGrpcTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *TransferHandleGrpcTestSuite) TestFindAllTransfer_Success() {
	req := &pb.FindAllTransferRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockTransfers := []*response.TransferResponse{
		{ID: 1, TransferNo: "TR123", TransferFrom: "ACC001", TransferTo: "ACC002", TransferAmount: 100000},
		{ID: 2, TransferNo: "TR124", TransferFrom: "ACC003", TransferTo: "ACC004", TransferAmount: 200000},
	}
	mockProtoTransfers := []*pb.TransferResponse{
		{Id: 1, TransferNo: "TR123", TransferFrom: "ACC001", TransferTo: "ACC002", TransferAmount: 100000},
		{Id: 2, TransferNo: "TR124", TransferFrom: "ACC003", TransferTo: "ACC004", TransferAmount: 200000},
	}

	totalRecords := 2
	suite.MockTransferService.EXPECT().
		FindAll(gomock.Eq(&requests.FindAllTranfers{Page: 1, PageSize: 10, Search: "test"})).
		Return(mockTransfers, &totalRecords, nil)

	suite.MockProtoMapper.EXPECT().
		ToProtoResponsePaginationTransfer(gomock.Any(), "success", "Successfully fetch transfer records", mockTransfers).
		Return(&pb.ApiResponsePaginationTransfer{
			Status:  "success",
			Message: "Successfully fetch transfer records",
			Data:    mockProtoTransfers,
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   1,
				TotalRecords: 2,
			},
		})

	res, err := suite.Handler.FindAllTransfer(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Successfully fetch transfer records", res.GetMessage())
	suite.Equal(int32(2), res.GetPagination().GetTotalRecords())
	suite.Equal(2, len(res.GetData()))
}

func (suite *TransferHandleGrpcTestSuite) TestFindAllTransfer_Failure() {
	req := &pb.FindAllTransferRequest{Page: 1, PageSize: 10, Search: "test"}
	serviceError := &response.ErrorResponse{Status: "error", Message: "Failed to fetch transfers"}

	totalRecords := 0
	suite.MockTransferService.EXPECT().FindAll(gomock.Any()).Return(nil, &totalRecords, serviceError)

	res, _ := suite.Handler.FindAllTransfer(context.Background(), req)

	suite.Nil(res)
}

func (suite *TransferHandleGrpcTestSuite) TestFindByIdTransfer_Success() {
	req := &pb.FindByIdTransferRequest{TransferId: 1}
	mockTransfer := &response.TransferResponse{ID: 1, TransferNo: "TR123", TransferAmount: 100000}
	mockProtoTransfer := &pb.TransferResponse{Id: 1, TransferNo: "TR123", TransferAmount: 100000}

	suite.MockTransferService.EXPECT().FindById(1).Return(mockTransfer, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransfer("success", "Successfully fetch transfer record", mockTransfer).Return(&pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Successfully fetch transfer record",
		Data:    mockProtoTransfer,
	})

	res, err := suite.Handler.FindByIdTransfer(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("TR123", res.GetData().GetTransferNo())
}

func (suite *TransferHandleGrpcTestSuite) TestFindByIdTransfer_InvalidId() {
	req := &pb.FindByIdTransferRequest{TransferId: 0}

	res, err := suite.Handler.FindByIdTransfer(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid Transfer ID")
}

func (suite *TransferHandleGrpcTestSuite) TestFindByTransferByTransferFrom_Success() {
	req := &pb.FindTransferByTransferFromRequest{TransferFrom: "ACC001"}
	mockTransfers := []*response.TransferResponse{
		{ID: 1, TransferNo: "TR123", TransferFrom: "ACC001", TransferTo: "ACC002", TransferAmount: 100000},
	}

	suite.MockTransferService.EXPECT().FindTransferByTransferFrom("ACC001").Return(mockTransfers, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransfers("success", "Successfully fetch transfer records", mockTransfers).Return(&pb.ApiResponseTransfers{
		Status:  "success",
		Message: "Successfully fetch transfer records",
		Data: []*pb.TransferResponse{
			{Id: 1, TransferNo: "TR123", TransferFrom: "ACC001", TransferTo: "ACC002", TransferAmount: 100000},
		},
	})

	res, err := suite.Handler.FindByTransferByTransferFrom(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *TransferHandleGrpcTestSuite) TestFindByTransferByTransferTo_Success() {
	req := &pb.FindTransferByTransferToRequest{TransferTo: "ACC002"}
	mockTransfers := []*response.TransferResponse{
		{ID: 1, TransferNo: "TR123", TransferFrom: "ACC001", TransferTo: "ACC002", TransferAmount: 100000},
	}

	suite.MockTransferService.EXPECT().FindTransferByTransferTo("ACC002").Return(mockTransfers, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransfers("success", "Successfully fetch transfer records", mockTransfers).Return(&pb.ApiResponseTransfers{
		Status:  "success",
		Message: "Successfully fetch transfer records",
		Data: []*pb.TransferResponse{
			{Id: 1, TransferNo: "TR123", TransferFrom: "ACC001", TransferTo: "ACC002", TransferAmount: 100000},
		},
	})

	res, err := suite.Handler.FindByTransferByTransferTo(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *TransferHandleGrpcTestSuite) TestFindByActiveTransfer_Success() {
	req := &pb.FindAllTransferRequest{Page: 1, PageSize: 10, Search: ""}
	activeTransfers := []*response.TransferResponseDeleteAt{
		{ID: 1, TransferNo: "TR123", TransferFrom: "ACC001", TransferTo: "ACC002", TransferAmount: 100000},
	}
	totalRecords := 1
	suite.MockTransferService.EXPECT().FindByActive(gomock.Any()).Return(activeTransfers, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationTransferDeleteAt(gomock.Any(), "success", "Successfully fetch transfer records", activeTransfers).Return(&pb.ApiResponsePaginationTransferDeleteAt{
		Status:  "success",
		Message: "Successfully fetch transfer records",
		Data: []*pb.TransferResponseDeleteAt{
			{Id: 1, TransferNo: "TR123", TransferFrom: "ACC001", TransferTo: "ACC002", TransferAmount: 100000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByActiveTransfer(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *TransferHandleGrpcTestSuite) TestFindByTrashedTransfer_Success() {
	req := &pb.FindAllTransferRequest{Page: 1, PageSize: 10, Search: ""}
	trashedTransfers := []*response.TransferResponseDeleteAt{
		{ID: 1, TransferNo: "TR123", TransferFrom: "ACC001", TransferTo: "ACC002", TransferAmount: 100000},
	}
	totalRecords := 1
	suite.MockTransferService.EXPECT().FindByTrashed(gomock.Any()).Return(trashedTransfers, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationTransferDeleteAt(gomock.Any(), "success", "Successfully fetch transfer records", trashedTransfers).Return(&pb.ApiResponsePaginationTransferDeleteAt{
		Status:  "success",
		Message: "Successfully fetch transfer records",
		Data: []*pb.TransferResponseDeleteAt{
			{Id: 1, TransferNo: "TR123", TransferFrom: "ACC001", TransferTo: "ACC002", TransferAmount: 100000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByTrashedTransfer(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *TransferHandleGrpcTestSuite) TestCreateTransfer_Success() {
	transferTime := time.Now()
	req := &pb.CreateTransferRequest{
		TransferFrom:   "ACC001",
		TransferTo:     "ACC002",
		TransferAmount: 150000,
	}

	mockTransfer := &response.TransferResponse{
		ID:             1,
		TransferNo:     "TR125",
		TransferFrom:   "ACC001",
		TransferTo:     "ACC002",
		TransferAmount: 150000,
		TransferTime:   transferTime.Format(time.RFC3339),
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	mockProtoTransfer := &pb.TransferResponse{
		Id:             1,
		TransferNo:     "TR125",
		TransferFrom:   "ACC001",
		TransferTo:     "ACC002",
		TransferAmount: 150000,
		TransferTime:   transferTime.Format(time.RFC3339),
	}

	suite.MockTransferService.EXPECT().CreateTransaction(gomock.Any()).Return(mockTransfer, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransfer("success", "Successfully created transfer", mockTransfer).Return(&pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Successfully created transfer",
		Data:    mockProtoTransfer,
	})

	res, err := suite.Handler.CreateTransfer(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal(int32(150000), res.GetData().GetTransferAmount())
}

func (suite *TransferHandleGrpcTestSuite) TestCreateTransfer_ValidationError() {
	req := &pb.CreateTransferRequest{
		TransferFrom:   "",
		TransferTo:     "",
		TransferAmount: 40000,
	}

	res, err := suite.Handler.CreateTransfer(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid input for create transfer")
}

func (suite *TransferHandleGrpcTestSuite) TestUpdateTransfer_Success() {
	transferId := 1
	req := &pb.UpdateTransferRequest{
		TransferId:     int32(transferId),
		TransferFrom:   "ACC001",
		TransferTo:     "ACC002",
		TransferAmount: 200000,
	}

	mockTransfer := &response.TransferResponse{
		ID:             1,
		TransferNo:     "TR125",
		TransferAmount: 200000,
	}

	mockProtoTransfer := &pb.TransferResponse{Id: 1, TransferAmount: 200000}

	suite.MockTransferService.EXPECT().UpdateTransaction(gomock.Any()).Return(mockTransfer, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransfer("success", "Successfully updated transfer", mockTransfer).Return(&pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Successfully updated transfer",
		Data:    mockProtoTransfer,
	})

	res, err := suite.Handler.UpdateTransfer(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal(int32(200000), res.GetData().GetTransferAmount())
}

func (suite *TransferHandleGrpcTestSuite) TestTrashedTransfer_Success() {
	req := &pb.FindByIdTransferRequest{TransferId: 1}
	mockTransfer := &response.TransferResponseDeleteAt{ID: 1, TransferNo: "TR125"}

	suite.MockTransferService.EXPECT().TrashedTransfer(1).Return(mockTransfer, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransferDeleteAt("success", "Successfully trashed transfer", mockTransfer).Return(&pb.ApiResponseTransferDeleteAt{
		Status:  "success",
		Message: "Successfully trashed transfer",
		Data:    &pb.TransferResponseDeleteAt{Id: 1, TransferNo: "TR125"},
	})

	res, err := suite.Handler.TrashedTransfer(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *TransferHandleGrpcTestSuite) TestRestoreTransfer_Success() {
	req := &pb.FindByIdTransferRequest{TransferId: 1}
	mockTransfer := &response.TransferResponseDeleteAt{ID: 1, TransferNo: "TR125"}

	suite.MockTransferService.EXPECT().RestoreTransfer(1).Return(mockTransfer, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransferDeleteAt("success", "Successfully restored transfer", mockTransfer).Return(&pb.ApiResponseTransferDeleteAt{
		Status:  "success",
		Message: "Successfully restored transfer",
		Data:    &pb.TransferResponseDeleteAt{Id: 1, TransferNo: "TR125"},
	})

	res, err := suite.Handler.RestoreTransfer(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *TransferHandleGrpcTestSuite) TestDeleteTransferPermanent_Success() {
	req := &pb.FindByIdTransferRequest{TransferId: 1}

	suite.MockTransferService.EXPECT().DeleteTransferPermanent(1).Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransferDelete("success", "Successfully deleted transfer permanently").Return(&pb.ApiResponseTransferDelete{
		Status:  "success",
		Message: "Successfully deleted transfer permanently",
	})

	res, err := suite.Handler.DeleteTransferPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *TransferHandleGrpcTestSuite) TestRestoreAllTransfer_Success() {
	req := &emptypb.Empty{}

	suite.MockTransferService.EXPECT().RestoreAllTransfer().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransferAll("success", "Successfully restored all transfers").Return(&pb.ApiResponseTransferAll{
		Status:  "success",
		Message: "Successfully restored all transfers",
	})

	res, err := suite.Handler.RestoreAllTransfer(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *TransferHandleGrpcTestSuite) TestDeleteAllTransferPermanent_Success() {
	req := &emptypb.Empty{}

	suite.MockTransferService.EXPECT().DeleteAllTransferPermanent().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransferAll("success", "Successfully deleted all transfers permanently").Return(&pb.ApiResponseTransferAll{
		Status:  "success",
		Message: "Successfully deleted all transfers permanently",
	})

	res, err := suite.Handler.DeleteAllTransferPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func TestTransferHandleGrpcSuite(t *testing.T) {
	suite.Run(t, new(TransferHandleGrpcTestSuite))
}
