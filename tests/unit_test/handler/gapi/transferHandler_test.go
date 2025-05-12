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

func TestFindAllTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.FindAllTransferRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	transfers := []*response.TransferResponse{
		{
			ID:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
		},
		{
			ID:             2,
			TransferFrom:   "user3",
			TransferTo:     "user4",
			TransferAmount: 2000,
		},
	}

	mockTransferService.EXPECT().FindAll(1, 10, "test").Return(transfers, 2, nil).Times(1)
	mockTransferMapper.EXPECT().ToResponsesTransfer(transfers).Return([]*pb.TransferResponse{
		{
			Id:             1,
			TransferFrom:   "user1",
			TransferTo:     "user2",
			TransferAmount: 1000,
		},
		{
			Id:             2,
			TransferFrom:   "user3",
			TransferTo:     "user4",
			TransferAmount: 2000,
		},
	}).Times(1)

	res, err := mockHandler.FindAllTransfer(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transfer records", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
	assert.Equal(t, int32(1), res.GetPagination().GetTotalPages())
	assert.Equal(t, int32(2), res.GetPagination().GetTotalRecords())
}

func TestFindAllTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.FindAllTransferRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockTransferService.EXPECT().FindAll(1, 10, "test").Return(nil, 0, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch transfer records",
	}).Times(1)

	res, err := mockHandler.FindAllTransfer(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch transfer records")
}

func TestFindAllTransfer_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.FindAllTransferRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockTransferService.EXPECT().FindAll(1, 10, "test").Return([]*response.TransferResponse{}, 0, nil).Times(1)
	mockTransferMapper.EXPECT().ToResponsesTransfer([]*response.TransferResponse{}).Return([]*pb.TransferResponse{}).Times(1)

	res, err := mockHandler.FindAllTransfer(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transfer records", res.GetMessage())
	assert.Len(t, res.GetData(), 0)
	assert.Equal(t, int32(0), res.GetPagination().GetTotalPages())
	assert.Equal(t, int32(0), res.GetPagination().GetTotalRecords())
}

func TestFindTransferById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.FindByIdTransferRequest{
		TransferId: 1,
	}

	transfer := &response.TransferResponse{
		ID:             1,
		TransferFrom:   "user1",
		TransferTo:     "user2",
		TransferAmount: 1000,
	}

	mockTransferService.EXPECT().FindById(1).Return(transfer, nil).Times(1)
	mockTransferMapper.EXPECT().ToResponseTransfer(transfer).Return(&pb.TransferResponse{
		Id:             1,
		TransferFrom:   "user1",
		TransferTo:     "user2",
		TransferAmount: 1000,
	}).Times(1)

	res, err := mockHandler.FindTransferById(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transfer record", res.GetMessage())
	assert.Equal(t, int32(1), res.GetData().GetId())
	assert.Equal(t, int32(1000), res.GetData().GetTransferAmount())
}

func TestFindTransferById_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.FindByIdTransferRequest{
		TransferId: 0,
	}

	res, err := mockHandler.FindTransferById(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Invalid transfer id")
}

func TestFindTransferById_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.FindByIdTransferRequest{
		TransferId: 1,
	}

	mockTransferService.EXPECT().FindById(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch transfer record",
	}).Times(1)

	res, err := mockHandler.FindTransferById(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
	assert.Contains(t, err.Error(), "Failed to fetch transfer record")
}

func TestFindByTransferByTransferFrom_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.FindTransferByTransferFromRequest{
		TransferFrom: "sourceAccount",
	}

	transfer := []*response.TransferResponse{
		{
			ID:             1,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount",
			TransferAmount: 1000,
		},
		{
			ID:             2,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount2",
			TransferAmount: 2000,
		},
	}

	mockTransferService.EXPECT().FindTransferByTransferFrom("sourceAccount").Return(transfer, nil).Times(1)
	mockTransferMapper.EXPECT().ToResponsesTransfer(transfer).Return([]*pb.TransferResponse{
		{
			Id:             1,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount",
			TransferAmount: 1000,
		},
		{
			Id:             2,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount2",
			TransferAmount: 2000,
		},
	}).Times(1)

	res, err := mockHandler.FindByTransferByTransferFrom(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transfer records", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
}

func TestFindByTransferByTransferFrom_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.FindTransferByTransferFromRequest{
		TransferFrom: "sourceAccount",
	}

	mockTransferService.EXPECT().FindTransferByTransferFrom("sourceAccount").Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch transfer records",
	}).Times(1)

	res, err := mockHandler.FindByTransferByTransferFrom(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch transfer records")
}

func TestFindByTransferByTransferTo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.FindTransferByTransferToRequest{
		TransferTo: "destinationAccount",
	}

	transfer := []*response.TransferResponse{
		{
			ID:             1,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount",
			TransferAmount: 1000,
		},
		{
			ID:             2,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount2",
			TransferAmount: 2000,
		},
	}

	mockTransferService.EXPECT().FindTransferByTransferTo("destinationAccount").Return(transfer, nil).Times(1)
	mockTransferMapper.EXPECT().ToResponsesTransfer(transfer).Return([]*pb.TransferResponse{
		{
			Id:             1,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount",
			TransferAmount: 1000,
		},
		{
			Id:             2,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount2",
			TransferAmount: 2000,
		},
	}).Times(1)

	res, err := mockHandler.FindByTransferByTransferTo(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transfer records", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
}

func TestFindByTransferByTransferTo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.FindTransferByTransferToRequest{
		TransferTo: "destinationAccount",
	}

	mockTransferService.EXPECT().FindTransferByTransferTo("destinationAccount").Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch transfer records",
	}).Times(1)

	res, err := mockHandler.FindByTransferByTransferTo(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch transfer records")
}

func TestFindByActiveTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllTransferRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	transfer := []*response.TransferResponseDeleteAt{
		{
			ID:             1,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount",
			TransferAmount: 1000,
		},
		{
			ID:             2,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount2",
			TransferAmount: 2000,
		},
	}

	mockTransferService.EXPECT().FindByActive(pageSize, page, search).Return(transfer, expected, nil).Times(1)
	mockTransferMapper.EXPECT().ToResponsesTransferDeleteAt(transfer).Return([]*pb.TransferResponseDeleteAt{
		{
			Id:             1,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount",
			TransferAmount: 1000,
		},
		{
			Id:             2,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount2",
			TransferAmount: 2000,
		},
	}).Times(1)

	res, err := mockHandler.FindByActiveTransfer(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transfer records", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
}

func TestFindByActiveTransfer_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllTransferRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockTransferService.EXPECT().FindByActive(pageSize, page, search).Return([]*response.TransferResponseDeleteAt{}, expected, nil).Times(1)
	mockTransferMapper.EXPECT().ToResponsesTransferDeleteAt([]*response.TransferResponseDeleteAt{}).Return([]*pb.TransferResponseDeleteAt{}).Times(1)

	res, err := mockHandler.FindByActiveTransfer(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transfer records", res.GetMessage())
	assert.Len(t, res.GetData(), 0)
}

func TestFindByActiveTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllTransferRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockTransferService.EXPECT().FindByActive(pageSize, page, search).Return(nil, expected, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch transfer records",
	}).Times(1)

	res, err := mockHandler.FindByActiveTransfer(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch transfer records")
}

func TestFindByTrashedTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	transfer := []*response.TransferResponseDeleteAt{
		{
			ID:             1,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount",
			TransferAmount: 1000,
		},
		{
			ID:             2,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount2",
			TransferAmount: 2000,
		},
	}

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllTransferRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockTransferService.EXPECT().FindByTrashed(pageSize, page, search).Return(transfer, expected, nil).Times(1)
	mockTransferMapper.EXPECT().ToResponsesTransferDeleteAt(transfer).Return([]*pb.TransferResponseDeleteAt{
		{
			Id:             1,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount",
			TransferAmount: 1000,
		},
		{
			Id:             2,
			TransferFrom:   "sourceAccount",
			TransferTo:     "destinationAccount2",
			TransferAmount: 2000,
		},
	}).Times(1)

	res, err := mockHandler.FindByTrashedTransfer(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transfer records", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
}

func TestFindByTrashedTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllTransferRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockTransferService.EXPECT().FindByTrashed(pageSize, page, search).Return(nil, expected, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch transfer records",
	}).Times(1)

	res, err := mockHandler.FindByTrashedTransfer(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch transfer records")
}

func TestFindByTrashedTransfer_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllTransferRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockTransferService.EXPECT().FindByTrashed(pageSize, page, search).Return([]*response.TransferResponseDeleteAt{}, expected, nil).Times(1)
	mockTransferMapper.EXPECT().ToResponsesTransferDeleteAt([]*response.TransferResponseDeleteAt{}).Return([]*pb.TransferResponseDeleteAt{}).Times(1)

	res, err := mockHandler.FindByTrashedTransfer(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transfer records", res.GetMessage())
	assert.Len(t, res.GetData(), 0)
}

func TestCreateTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.CreateTransferRequest{
		TransferFrom:   "sourceAccount",
		TransferTo:     "destinationAccount",
		TransferAmount: 1000,
	}

	createReq := &requests.CreateTransferRequest{
		TransferFrom:   req.GetTransferFrom(),
		TransferTo:     req.GetTransferTo(),
		TransferAmount: int(req.GetTransferAmount()),
	}

	expectedResponse := &response.TransferResponse{
		ID:             1,
		TransferFrom:   "sourceAccount",
		TransferTo:     "destinationAccount",
		TransferAmount: 1000,
	}

	mockTransferService.EXPECT().CreateTransaction(createReq).Return(expectedResponse, nil).Times(1)
	mockTransferMapper.EXPECT().ToResponseTransfer(expectedResponse).Return(&pb.TransferResponse{
		Id:             1,
		TransferFrom:   "sourceAccount",
		TransferTo:     "destinationAccount",
		TransferAmount: 1000,
	}).Times(1)

	res, err := mockHandler.CreateTransfer(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully created transfer", res.GetMessage())
	assert.Equal(t, int32(1), res.GetData().GetId())
}

func TestCreateTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.CreateTransferRequest{
		TransferFrom:   "sourceAccount",
		TransferTo:     "destinationAccount",
		TransferAmount: 1000,
	}

	createReq := &requests.CreateTransferRequest{
		TransferFrom:   req.GetTransferFrom(),
		TransferTo:     req.GetTransferTo(),
		TransferAmount: int(req.GetTransferAmount()),
	}

	mockTransferService.EXPECT().CreateTransaction(createReq).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to create transfer",
	}).Times(1)

	res, err := mockHandler.CreateTransfer(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to create transfer")
}

func TestCreateTransfer_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.CreateTransferRequest{
		TransferFrom:   "sourceAccount",
		TransferTo:     "destinationAccount",
		TransferAmount: -1000,
	}

	request := &requests.CreateTransferRequest{
		TransferFrom:   req.GetTransferFrom(),
		TransferTo:     req.GetTransferTo(),
		TransferAmount: int(req.GetTransferAmount()),
	}

	mockTransferService.EXPECT().CreateTransaction(request).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "validation error",
	}).Times(1)

	res, err := mockHandler.CreateTransfer(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "validation error")
}

func TestUpdateTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockMapper)

	req := &pb.UpdateTransferRequest{
		TransferId:     1,
		TransferFrom:   "AccountA",
		TransferTo:     "AccountB",
		TransferAmount: 1000,
	}

	updateReq := &requests.UpdateTransferRequest{
		TransferID:     1,
		TransferFrom:   "AccountA",
		TransferTo:     "AccountB",
		TransferAmount: 1000,
	}

	transfer := &response.TransferResponse{
		ID:             1,
		TransferFrom:   "AccountA",
		TransferTo:     "AccountB",
		TransferAmount: 1000,
	}

	mockTransferService.EXPECT().UpdateTransaction(updateReq).Return(transfer, nil).Times(1)
	mockMapper.EXPECT().ToResponseTransfer(transfer).Return(&pb.TransferResponse{
		Id:             1,
		TransferFrom:   "AccountA",
		TransferTo:     "AccountB",
		TransferAmount: 1000,
	}).Times(1)

	res, err := mockHandler.UpdateTransfer(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully updated transfer", res.GetMessage())
	assert.NotNil(t, res.GetData())
	assert.Equal(t, int32(1), res.GetData().GetId())
}

func TestUpdateTransfer_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockMapper)

	req := &pb.UpdateTransferRequest{
		TransferId:     -1,
		TransferFrom:   "AccountA",
		TransferTo:     "AccountB",
		TransferAmount: 1000,
	}

	res, err := mockHandler.UpdateTransfer(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)

	statusErr, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Transfer ID is required")
}

func TestUpdateTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, nil)

	req := &pb.UpdateTransferRequest{
		TransferId:     2,
		TransferFrom:   "AccountA",
		TransferTo:     "AccountB",
		TransferAmount: 2000,
	}

	updateReq := &requests.UpdateTransferRequest{
		TransferID:     2,
		TransferFrom:   "AccountA",
		TransferTo:     "AccountB",
		TransferAmount: 2000,
	}

	mockTransferService.EXPECT().UpdateTransaction(updateReq).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to update transfer",
	}).Times(1)

	res, err := mockHandler.UpdateTransfer(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to update transfer")
}

func TestUpdateTransfer_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockTransferMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockTransferMapper)

	req := &pb.UpdateTransferRequest{
		TransferId:     1,
		TransferFrom:   "sourceAccount",
		TransferTo:     "destinationAccount",
		TransferAmount: -1000,
	}

	request := &requests.UpdateTransferRequest{
		TransferID:     1,
		TransferFrom:   req.GetTransferFrom(),
		TransferTo:     req.GetTransferTo(),
		TransferAmount: int(req.GetTransferAmount()),
	}

	mockTransferService.EXPECT().UpdateTransaction(request).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "validation error",
	}).Times(1)

	res, err := mockHandler.UpdateTransfer(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "validation error")
}

func TestTrashedTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockMapper)

	req := &pb.FindByIdTransferRequest{TransferId: 1}
	expectedResponse := &response.TransferResponse{
		ID:             1,
		TransferFrom:   "AccountA",
		TransferTo:     "AccountB",
		TransferAmount: 1000,
	}

	mockTransferService.EXPECT().TrashedTransfer(1).Return(expectedResponse, nil).Times(1)
	mockMapper.EXPECT().ToResponseTransfer(expectedResponse).Return(&pb.TransferResponse{
		Id:             1,
		TransferFrom:   "AccountA",
		TransferTo:     "AccountB",
		TransferAmount: 1000,
	}).Times(1)

	res, err := mockHandler.TrashedTransfer(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully trashed transfer", res.GetMessage())
	assert.Equal(t, int32(1), res.GetData().GetId())
}

func TestTrashedTransfer_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockMapper)

	req := &pb.FindByIdTransferRequest{TransferId: 0}

	res, err := mockHandler.TrashedTransfer(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Transfer ID is required")
}

func TestTrashedTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, nil)

	req := &pb.FindByIdTransferRequest{TransferId: 1}

	mockTransferService.EXPECT().TrashedTransfer(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to trash transfer",
	}).Times(1)

	res, err := mockHandler.TrashedTransfer(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to trash transfer")
}

func TestRestoreTransfer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockMapper := mock_protomapper.NewMockTransferProtoMapper(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, mockMapper)

	req := &pb.FindByIdTransferRequest{TransferId: 1}
	expectedResponse := &response.TransferResponse{
		ID:             1,
		TransferFrom:   "AccountA",
		TransferTo:     "AccountB",
		TransferAmount: 1000,
	}

	mockTransferService.EXPECT().RestoreTransfer(1).Return(expectedResponse, nil).Times(1)
	mockMapper.EXPECT().ToResponseTransfer(expectedResponse).Return(&pb.TransferResponse{
		Id:             1,
		TransferFrom:   "AccountA",
		TransferTo:     "AccountB",
		TransferAmount: 1000,
	}).Times(1)

	res, err := mockHandler.RestoreTransfer(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully restored transfer", res.GetMessage())
	assert.Equal(t, int32(1), res.GetData().GetId())
}

func TestRestoreTransfer_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)

	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, nil)

	req := &pb.FindByIdTransferRequest{TransferId: 0}

	res, err := mockHandler.RestoreTransfer(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Transfer ID is required")
}

func TestRestoreTransfer_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, nil)

	req := &pb.FindByIdTransferRequest{TransferId: 1}

	mockTransferService.EXPECT().RestoreTransfer(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to restore transfer",
	}).Times(1)

	res, err := mockHandler.RestoreTransfer(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to restore transfer")
}

func TestDeleteTransferPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, nil)

	req := &pb.FindByIdTransferRequest{TransferId: 1}

	mockResponse := &pb.ApiResponseTransferDelete{
		Status:  "success",
		Message: "Successfully deleted transfer",
	}

	mockTransferService.EXPECT().DeleteTransferPermanent(1).Return(mockResponse, nil).Times(1)

	res, err := mockHandler.DeleteTransferPermanent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully deleted transfer", res.GetMessage())
}

func TestDeleteTransferPermanent_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)

	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, nil)

	req := &pb.FindByIdTransferRequest{TransferId: 0}

	res, err := mockHandler.DeleteTransferPermanent(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Transfer ID is required")
}

func TestDeleteTransferPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransferService := mock_service.NewMockTransferService(ctrl)
	mockHandler := gapi.NewTransferHandleGrpc(mockTransferService, nil)

	req := &pb.FindByIdTransferRequest{TransferId: 1}

	mockTransferService.EXPECT().DeleteTransferPermanent(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to delete transfer",
	}).Times(1)

	res, err := mockHandler.DeleteTransferPermanent(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to delete transfer")
}
