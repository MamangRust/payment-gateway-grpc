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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFindAllWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	req := &pb.FindAllWithdrawRequest{
		Page:     1,
		PageSize: 10,
		Search:   "example",
	}

	withdraws := []*response.WithdrawResponse{
		{
			ID:         1,
			CardNumber: "1234",
		},
		{
			ID:         2,
			CardNumber: "5678",
		},
	}
	totalRecords := 2

	mockWithdrawService.EXPECT().FindAll(1, 10, "example").Return(withdraws, totalRecords, nil).Times(1)
	mockMapper.EXPECT().ToResponsesWithdrawal(withdraws).Return([]*pb.WithdrawResponse{
		{
			WithdrawId: 1,
			CardNumber: "1234",
		},
		{
			WithdrawId: 2,
			CardNumber: "5678",
		},
	}).Times(1)

	res, err := mockHandler.FindAllWithdraw(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Withdraws fetched successfully", res.GetMessage())
	assert.Equal(t, int32(2), res.GetPagination().GetTotalRecords())
	assert.Equal(t, 2, len(res.GetData()))
}

func TestFindAllWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.FindAllWithdrawRequest{
		Page:     1,
		PageSize: 10,
		Search:   "example",
	}

	mockWithdrawService.EXPECT().FindAll(1, 10, "example").Return(nil, 0, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch withdraws",
	}).Times(1)

	res, err := mockHandler.FindAllWithdraw(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch withdraws")
}

func TestFindAllWithdraw_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	req := &pb.FindAllWithdrawRequest{
		Page:     1,
		PageSize: 10,
		Search:   "example",
	}

	withdraws := []*response.WithdrawResponse{}
	totalRecords := 0

	mockWithdrawService.EXPECT().FindAll(1, 10, "example").Return(withdraws, totalRecords, nil).Times(1)
	mockMapper.EXPECT().ToResponsesWithdrawal(withdraws).Return([]*pb.WithdrawResponse{}).Times(1)

	res, err := mockHandler.FindAllWithdraw(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Withdraws fetched successfully", res.GetMessage())
	assert.Equal(t, int32(0), res.GetPagination().GetTotalRecords())
	assert.Equal(t, 0, len(res.GetData()))
}

func TestFindByIdWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: 1,
	}

	withdraw := response.WithdrawResponse{
		ID:         1,
		CardNumber: "1234",
	}

	mockWithdrawService.EXPECT().FindById(1).Return(&withdraw, nil).Times(1)
	mockMapper.EXPECT().ToResponseWithdrawal(&withdraw).Return(&pb.WithdrawResponse{
		WithdrawId: 1,
		CardNumber: "1234",
	}).Times(1)

	res, err := mockHandler.FindByIdWithdraw(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched withdraw", res.GetMessage())
	assert.Equal(t, int32(1), res.GetData().GetWithdrawId())
}

func TestFindByIdWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: 1,
	}

	mockWithdrawService.EXPECT().FindById(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch withdraw",
	}).Times(1)

	res, err := mockHandler.FindByIdWithdraw(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch withdraw")
}

func TestFindByCardNumberWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	req := &pb.FindByCardNumberRequest{
		CardNumber: "1234-5678-9012",
	}

	withdraws := []*response.WithdrawResponse{
		{ID: 1, CardNumber: "1234-5678-9012", WithdrawAmount: 1000},
		{ID: 2, CardNumber: "1234-5678-9012", WithdrawAmount: 2000},
	}

	mockWithdrawService.EXPECT().FindByCardNumber("1234-5678-9012").Return(withdraws, nil).Times(1)
	mockMapper.EXPECT().ToResponsesWithdrawal(withdraws).Return([]*pb.WithdrawResponse{
		{WithdrawId: 1, CardNumber: "1234-5678-9012", WithdrawAmount: 1000},
		{WithdrawId: 2, CardNumber: "1234-5678-9012", WithdrawAmount: 2000},
	}).Times(1)

	res, err := mockHandler.FindByCardNumber(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched withdraws", res.GetMessage())
	assert.Equal(t, 2, len(res.GetData()))
}

func TestFindByCardNumberWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.FindByCardNumberRequest{
		CardNumber: "1234-5678-9012",
	}

	mockWithdrawService.EXPECT().FindByCardNumber("1234-5678-9012").Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch withdraws",
	}).Times(1)

	res, err := mockHandler.FindByCardNumber(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch withdraws")
}

func TestFindByActiveWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	withdraws := []*response.WithdrawResponseDeleteAt{
		{ID: 1, CardNumber: "1234-5678-9012", WithdrawAmount: 1000},
		{ID: 2, CardNumber: "1234-5678-9012", WithdrawAmount: 2000},
	}

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockWithdrawService.EXPECT().FindByActive(pageSize, page, search).Return(withdraws, expected, nil).Times(1)
	mockMapper.EXPECT().ToResponsesWithdrawalDeleteAt(withdraws).Return([]*pb.WithdrawResponseDeleteAt{
		{WithdrawId: 1, CardNumber: "1234-5678-9012", WithdrawAmount: 1000},
		{WithdrawId: 2, CardNumber: "1234-5678-9012", WithdrawAmount: 2000},
	}).Times(1)

	res, err := mockHandler.FindByActive(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched withdraws", res.GetMessage())
	assert.Equal(t, 2, len(res.GetData()))
}

func TestFindByActiveWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockWithdrawService.EXPECT().FindByActive(pageSize, page, search).Return(nil, expected, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch withdraws",
	}).Times(1)

	res, err := mockHandler.FindByActive(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch withdraws")
}

func TestFindByActive_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockWithdrawService.EXPECT().FindByActive(pageSize, page, search).Return([]*response.WithdrawResponseDeleteAt{}, expected, nil).Times(1)
	mockMapper.EXPECT().ToResponsesWithdrawalDeleteAt([]*response.WithdrawResponseDeleteAt{}).Return([]*pb.WithdrawResponseDeleteAt{}).Times(1)

	res, err := mockHandler.FindByActive(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched withdraws", res.GetMessage())
	assert.Empty(t, res.GetData())
}

func TestFindByTrashedWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	withdraws := []*response.WithdrawResponseDeleteAt{
		{ID: 1, CardNumber: "1234-5678-9012", WithdrawAmount: 1000},
		{ID: 2, CardNumber: "1234-5678-9012", WithdrawAmount: 2000},
	}

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockWithdrawService.EXPECT().FindByTrashed(pageSize, page, search).Return(withdraws, expected, nil).Times(1)
	mockMapper.EXPECT().ToResponsesWithdrawalDeleteAt(withdraws).Return([]*pb.WithdrawResponseDeleteAt{
		{WithdrawId: 1, CardNumber: "1234-5678-9012", WithdrawAmount: 1000},
		{WithdrawId: 2, CardNumber: "1234-5678-9012", WithdrawAmount: 2000},
	}).Times(1)

	res, err := mockHandler.FindByTrashed(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched withdraws", res.GetMessage())
	assert.Equal(t, 2, len(res.GetData()))
}

func TestFindByTrashedWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockWithdrawService.EXPECT().FindByTrashed(pageSize, page, search).Return(nil, expected, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch withdraws",
	}).Times(1)

	res, err := mockHandler.FindByTrashed(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch withdraws")
}

func TestCreateWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	req := &pb.CreateWithdrawRequest{
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
		WithdrawTime:   timestamppb.Now(),
	}

	createRequest := &requests.CreateWithdrawRequest{
		CardNumber:     req.GetCardNumber(),
		WithdrawAmount: int(req.GetWithdrawAmount()),
		WithdrawTime:   req.GetWithdrawTime().AsTime(),
	}

	withdraw := &response.WithdrawResponse{
		ID:             1,
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
		WithdrawTime:   createRequest.WithdrawTime.String(),
	}

	mockWithdrawService.EXPECT().Create(createRequest).Return(withdraw, nil).Times(1)
	mockMapper.EXPECT().ToResponseWithdrawal(withdraw).Return(&pb.WithdrawResponse{
		WithdrawId:     1,
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
		WithdrawTime:   createRequest.WithdrawTime.String(),
	}).Times(1)

	res, err := mockHandler.CreateWithdraw(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully created withdraw", res.GetMessage())
	assert.Equal(t, int32(1), res.GetData().GetWithdrawId())
	assert.Equal(t, "123456789", res.GetData().GetCardNumber())
}

func TestCreateWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.CreateWithdrawRequest{
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
		WithdrawTime:   timestamppb.Now(),
	}

	createRequest := &requests.CreateWithdrawRequest{
		CardNumber:     req.GetCardNumber(),
		WithdrawAmount: int(req.GetWithdrawAmount()),
		WithdrawTime:   req.GetWithdrawTime().AsTime(),
	}

	mockWithdrawService.EXPECT().Create(createRequest).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to create withdraw",
	}).Times(1)

	res, err := mockHandler.CreateWithdraw(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to create withdraw")
}

func TestCreateWithdraw_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.CreateWithdrawRequest{
		CardNumber:     "",
		WithdrawAmount: 1000,
		WithdrawTime:   timestamppb.Now(),
	}

	createRequest := &requests.CreateWithdrawRequest{
		CardNumber:     req.GetCardNumber(),
		WithdrawAmount: int(req.GetWithdrawAmount()),
		WithdrawTime:   req.GetWithdrawTime().AsTime(),
	}

	mockWithdrawService.EXPECT().Create(createRequest).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "CardNumber is required",
	}).Times(1)

	res, err := mockHandler.CreateWithdraw(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "CardNumber is required")
}

func TestUpdateWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	req := &pb.UpdateWithdrawRequest{
		WithdrawId:     1,
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
		WithdrawTime:   timestamppb.Now(),
	}

	updateRequest := &requests.UpdateWithdrawRequest{
		CardNumber:     req.GetCardNumber(),
		WithdrawID:     int(req.GetWithdrawId()),
		WithdrawAmount: int(req.GetWithdrawAmount()),
		WithdrawTime:   req.GetWithdrawTime().AsTime(),
	}

	withdraw := &response.WithdrawResponse{
		ID:             1,
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
		WithdrawTime:   updateRequest.WithdrawTime.String(),
	}

	mockWithdrawService.EXPECT().Update(updateRequest).Return(withdraw, nil).Times(1)
	mockMapper.EXPECT().ToResponseWithdrawal(withdraw).Return(&pb.WithdrawResponse{
		WithdrawId:     1,
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
		WithdrawTime:   updateRequest.WithdrawTime.String(),
	}).Times(1)

	res, err := mockHandler.UpdateWithdraw(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully updated withdraw", res.GetMessage())
	assert.Equal(t, int32(1), res.GetData().GetWithdrawId())
	assert.Equal(t, "123456789", res.GetData().GetCardNumber())
}

func TestUpdateWithdraw_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	req := &pb.UpdateWithdrawRequest{
		WithdrawId:     0,
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
		WithdrawTime:   timestamppb.Now(),
	}

	updateRequest := &requests.UpdateWithdrawRequest{
		CardNumber:     req.GetCardNumber(),
		WithdrawID:     int(req.GetWithdrawId()),
		WithdrawAmount: int(req.GetWithdrawAmount()),
		WithdrawTime:   req.GetWithdrawTime().AsTime(),
	}

	mockWithdrawService.EXPECT().Update(updateRequest).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Invalid withdraw ID",
	}).Times(1)

	res, err := mockHandler.UpdateWithdraw(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Contains(t, statusErr.Message(), "Invalid withdraw ID")
}
func TestUpdateWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	req := &pb.UpdateWithdrawRequest{
		WithdrawId:     1,
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
		WithdrawTime:   timestamppb.Now(),
	}

	updateRequest := &requests.UpdateWithdrawRequest{
		CardNumber:     req.GetCardNumber(),
		WithdrawID:     int(req.GetWithdrawId()),
		WithdrawAmount: int(req.GetWithdrawAmount()),
		WithdrawTime:   req.GetWithdrawTime().AsTime(),
	}

	mockWithdrawService.EXPECT().Update(updateRequest).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to update withdraw",
	}).Times(1)

	res, err := mockHandler.UpdateWithdraw(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to update withdraw")
}

func TestUpdateWithdraw_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.UpdateWithdrawRequest{
		WithdrawId:     1,
		CardNumber:     "",
		WithdrawAmount: 1000,
		WithdrawTime:   timestamppb.Now(),
	}

	updateRequest := &requests.UpdateWithdrawRequest{
		CardNumber:     req.GetCardNumber(),
		WithdrawID:     int(req.GetWithdrawId()),
		WithdrawAmount: int(req.GetWithdrawAmount()),
		WithdrawTime:   req.GetWithdrawTime().AsTime(),
	}

	mockWithdrawService.EXPECT().Update(updateRequest).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "CardNumber is required",
	}).Times(1)

	res, err := mockHandler.UpdateWithdraw(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "CardNumber is required")
}

func TestTrashedWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: 1,
	}

	withdraw := &response.WithdrawResponse{
		ID:             1,
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
	}

	mockWithdrawService.EXPECT().TrashedWithdraw(1).Return(withdraw, nil).Times(1)
	mockMapper.EXPECT().ToResponseWithdrawal(withdraw).Return(&pb.WithdrawResponse{
		WithdrawId:     1,
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
	}).Times(1)

	res, err := mockHandler.TrashedWithdraw(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully trashed withdraw", res.GetMessage())
	assert.Equal(t, int32(1), res.GetData().GetWithdrawId())
}

func TestTrashedWithdraw_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: 0,
	}

	res, err := mockHandler.TrashedWithdraw(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
	assert.Contains(t, statusErr.Message(), "Invalid withdraw id")
}

func TestTrashedWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: 1,
	}

	mockWithdrawService.EXPECT().TrashedWithdraw(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch withdraw",
	}).Times(1)

	res, err := mockHandler.TrashedWithdraw(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch withdraw")
}

func TestRestoreWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockMapper := mock_protomapper.NewMockWithdrawalProtoMapper(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, mockMapper)

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: 1,
	}

	withdraw := &response.WithdrawResponse{
		ID:             1,
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
	}

	mockWithdrawService.EXPECT().RestoreWithdraw(1).Return(withdraw, nil).Times(1)
	mockMapper.EXPECT().ToResponseWithdrawal(withdraw).Return(&pb.WithdrawResponse{
		WithdrawId:     1,
		CardNumber:     "123456789",
		WithdrawAmount: 1000,
	}).Times(1)

	res, err := mockHandler.RestoreWithdraw(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully restored withdraw", res.GetMessage())
	assert.Equal(t, int32(1), res.GetData().GetWithdrawId())
}

func TestRestoreWithdraw_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: 0,
	}

	res, err := mockHandler.RestoreWithdraw(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)

	statusErr, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, statusErr.Message(), "Invalid withdraw id")

}

func TestRestoreWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: 1,
	}

	mockWithdrawService.EXPECT().RestoreWithdraw(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch withdraw",
	}).Times(1)

	res, err := mockHandler.RestoreWithdraw(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch withdraw")
}

func TestDeleteWithdrawPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: 1,
	}
	mockResponse := &pb.ApiResponseWithdrawDelete{
		Status:  "success",
		Message: "Successfully deleted withdraw permanently",
	}

	mockWithdrawService.EXPECT().DeleteWithdrawPermanent(1).Return(mockResponse, nil).Times(1)

	res, err := mockHandler.DeleteWithdrawPermanent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully deleted withdraw permanently", res.GetMessage())
}

func TestDeleteWithdrawPermanent_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: 0,
	}

	res, err := mockHandler.DeleteWithdrawPermanent(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)

	status, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, status.Code())
	assert.Contains(t, status.Message(), "invalid withdraw id")
}

func TestDeleteWithdrawPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockHandler := gapi.NewWithdrawHandleGrpc(mockWithdrawService, nil)

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: 1,
	}

	mockWithdrawService.EXPECT().DeleteWithdrawPermanent(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch withdraw",
	}).Times(1)

	res, err := mockHandler.DeleteWithdrawPermanent(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch withdraw")
}
