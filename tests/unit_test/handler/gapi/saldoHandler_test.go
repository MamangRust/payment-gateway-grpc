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

func TestFindAllSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindAllSaldoRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockSaldo := []*response.SaldoResponse{
		{
			ID:           1,
			CardNumber:   "1234",
			TotalBalance: 10000,
		},
		{
			ID:           2,
			CardNumber:   "5678",
			TotalBalance: 20000,
		},
	}
	mockResponseSaldo := []*pb.SaldoResponse{
		{
			SaldoId:      1,
			CardNumber:   "1234",
			TotalBalance: 10000,
		},
		{
			SaldoId:      2,
			CardNumber:   "5678",
			TotalBalance: 20000,
		},
	}

	mockSaldoService.EXPECT().FindAll(1, 10, "test").Return(mockSaldo, 2, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponsesSaldo(mockSaldo).Return(mockResponseSaldo).Times(1)

	res, err := saldoHandler.FindAllSaldo(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully fetched saldo record", res.Message)
	assert.Equal(t, 2, int(res.Pagination.TotalRecords))
	assert.Equal(t, int32(1), res.Pagination.TotalPages)
	assert.Equal(t, mockResponseSaldo, res.Data)
}

func TestFindAllSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindAllSaldoRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockSaldoService.EXPECT().FindAll(1, 10, "test").Return(nil, 0, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch saldo records",
	}).Times(1)

	res, err := saldoHandler.FindAllSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to fetch saldo records")
}

func TestFindAllSaldo_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindAllSaldoRequest{
		Page:     1,
		PageSize: 10,
		Search:   "notfound",
	}

	mockSaldoService.EXPECT().FindAll(1, 10, "notfound").Return([]*response.SaldoResponse{}, 0, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponsesSaldo([]*response.SaldoResponse{}).Return([]*pb.SaldoResponse{}).Times(1)

	res, err := saldoHandler.FindAllSaldo(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully fetched saldo record", res.Message)
	assert.Equal(t, int32(0), res.Pagination.TotalRecords)
	assert.Equal(t, int32(0), res.Pagination.TotalPages)
	assert.Empty(t, res.Data)
}

func TestFindByIdSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByIdSaldoRequest{
		SaldoId: 1,
	}

	mockSaldo := &response.SaldoResponse{
		ID:           1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockSaldoPb := &pb.SaldoResponse{
		SaldoId:      1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockSaldoService.EXPECT().FindById(1).Return(mockSaldo, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponseSaldo(mockSaldo).Return(mockSaldoPb).Times(1)

	res, err := saldoHandler.FindByIdSaldo(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully fetched saldo record", res.Message)
	assert.Equal(t, mockSaldoPb, res.Data)
}

func TestFindByIdSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByIdSaldoRequest{
		SaldoId: 1,
	}

	mockSaldoService.EXPECT().FindById(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Saldo not found",
	}).Times(1)

	res, err := saldoHandler.FindByIdSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to fetch saldo record: Saldo not found")
}

func TestFindByIdSaldo_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByIdSaldoRequest{
		SaldoId: -1,
	}

	res, err := saldoHandler.FindByIdSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Invalid ID")
}

func TestFindByCardNumberSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByCardNumberRequest{
		CardNumber: "1234",
	}

	mockSaldo := &response.SaldoResponse{
		ID:           1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockSaldoPb := &pb.SaldoResponse{
		SaldoId:      1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockSaldoService.EXPECT().FindByCardNumber("1234").Return(mockSaldo, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponseSaldo(mockSaldo).Return(mockSaldoPb).Times(1)

	res, err := saldoHandler.FindByCardNumber(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully fetched saldo record", res.Message)
	assert.Equal(t, mockSaldoPb, res.Data)
}

func TestFindByCardNumberSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByCardNumberRequest{
		CardNumber: "1234",
	}

	mockSaldoService.EXPECT().FindByCardNumber("1234").Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Saldo not found",
	}).Times(1)

	res, err := saldoHandler.FindByCardNumber(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to fetch saldo record: Saldo not found")
}

func TestFindByActiveSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	mockSaldoResponses := []*response.SaldoResponseDeleteAt{
		{
			ID:           1,
			CardNumber:   "1234",
			TotalBalance: 10000,
		},
		{
			ID:           2,
			CardNumber:   "5678",
			TotalBalance: 5000,
		},
	}

	mockSaldoPbResponses := []*pb.SaldoResponseDeleteAt{
		{
			SaldoId:      1,
			CardNumber:   "1234",
			TotalBalance: 10000,
		},
		{
			SaldoId:      2,
			CardNumber:   "5678",
			TotalBalance: 5000,
		},
	}

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllSaldoRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockSaldoService.EXPECT().FindByActive(pageSize, page, search).Return(mockSaldoResponses, expected, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponsesSaldoDeleteAt(mockSaldoResponses).Return(mockSaldoPbResponses).Times(1)

	res, err := saldoHandler.FindByActive(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully fetched saldo record", res.Message)
	assert.Equal(t, mockSaldoPbResponses, res.Data)
}

func TestFindByActiveSaldo_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	mockSaldoResponses := []*response.SaldoResponseDeleteAt{}

	search := ""
	pageSize := 1
	page := 1
	expected := 0

	req := &pb.FindAllSaldoRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockSaldoService.EXPECT().FindByActive(pageSize, page, search).Return(mockSaldoResponses, expected, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponsesSaldoDeleteAt(mockSaldoResponses).Return([]*pb.SaldoResponseDeleteAt{}).Times(1)

	res, err := saldoHandler.FindByActive(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully fetched saldo record", res.Message)
	assert.Empty(t, res.Data)
}

func TestFindByActiveSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	search := ""
	pageSize := 1
	page := 1
	expected := 0

	req := &pb.FindAllSaldoRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockSaldoService.EXPECT().FindByActive(pageSize, page, search).Return(nil, expected, &response.ErrorResponse{
		Status:  "error",
		Message: "No active saldo found",
	}).Times(1)

	res, err := saldoHandler.FindByActive(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Saldo not found: No active saldo found")
}

func TestFindByTrashed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	mockSaldoResponses := []*response.SaldoResponseDeleteAt{
		{
			ID:           3,
			CardNumber:   "9999",
			TotalBalance: 0,
		},
	}

	mockSaldoPbResponses := []*pb.SaldoResponseDeleteAt{
		{
			SaldoId:      3,
			CardNumber:   "9999",
			TotalBalance: 0,
		},
	}

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllSaldoRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockSaldoService.EXPECT().FindByTrashed(pageSize, page, search).Return(mockSaldoResponses, expected, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponsesSaldoDeleteAt(mockSaldoResponses).Return(mockSaldoPbResponses).Times(1)

	res, err := saldoHandler.FindByTrashed(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully fetched saldo record", res.Message)
	assert.Equal(t, mockSaldoPbResponses, res.Data)
}

func TestFindByTrashedSaldo_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	mockSaldoResponses := []*response.SaldoResponseDeleteAt{}

	search := ""
	pageSize := 1
	page := 1
	expected := 0

	req := &pb.FindAllSaldoRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockSaldoService.EXPECT().FindByTrashed(pageSize, page, search).Return(mockSaldoResponses, expected, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponsesSaldoDeleteAt(mockSaldoResponses).Return([]*pb.SaldoResponseDeleteAt{}).Times(1)

	res, err := saldoHandler.FindByTrashed(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully fetched saldo record", res.Message)
	assert.Empty(t, res.Data)
}

func TestFindByTrashedSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	search := ""
	pageSize := 1
	page := 1
	expected := 0

	req := &pb.FindAllSaldoRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockSaldoService.EXPECT().FindByTrashed(pageSize, page, search).Return(nil, expected, &response.ErrorResponse{
		Status:  "error",
		Message: "No trashed saldo found",
	}).Times(1)

	res, err := saldoHandler.FindByTrashed(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Saldo not found: No trashed saldo found")
}

func TestCreateSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.CreateSaldoRequest{
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockRequest := &requests.CreateSaldoRequest{
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockResponse := &response.SaldoResponse{
		ID:           1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockPbResponse := &pb.SaldoResponse{
		SaldoId:      1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockSaldoService.EXPECT().CreateSaldo(mockRequest).Return(mockResponse, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponseSaldo(mockResponse).Return(mockPbResponse).Times(1)

	res, err := saldoHandler.CreateSaldo(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully created saldo record", res.Message)
	assert.Equal(t, mockPbResponse, res.Data)
}

func TestCreateSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.CreateSaldoRequest{
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockRequest := &requests.CreateSaldoRequest{
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockSaldoService.EXPECT().CreateSaldo(mockRequest).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to create saldo",
	}).Times(1)

	res, err := saldoHandler.CreateSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to create saldo")
}

func TestCreateSaldo_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.CreateSaldoRequest{
		CardNumber:   "",
		TotalBalance: -500,
	}

	res, err := saldoHandler.CreateSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())

}

func TestUpdateSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.UpdateSaldoRequest{
		SaldoId:      1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockRequest := &requests.UpdateSaldoRequest{
		SaldoID:      1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockResponse := &response.SaldoResponse{
		ID:           1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockPbResponse := &pb.SaldoResponse{
		SaldoId:      1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockSaldoService.EXPECT().UpdateSaldo(mockRequest).Return(mockResponse, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponseSaldo(mockResponse).Return(mockPbResponse).Times(1)

	res, err := saldoHandler.UpdateSaldo(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully updated saldo record", res.GetMessage())

}

func TestUpdateSaldo_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.UpdateSaldoRequest{
		SaldoId:      0,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	res, err := saldoHandler.UpdateSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Invalid ID")
}

func TestUpdateSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.UpdateSaldoRequest{
		SaldoId:      1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockRequest := &requests.UpdateSaldoRequest{
		SaldoID:      1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockSaldoService.EXPECT().UpdateSaldo(mockRequest).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to update saldo",
	}).Times(1)

	res, err := saldoHandler.UpdateSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to update saldo")

}

func TestUpdateSaldo_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.UpdateSaldoRequest{
		SaldoId:      1,
		CardNumber:   "",
		TotalBalance: -500,
	}

	res, err := saldoHandler.UpdateSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
}

func TestTrashSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByIdSaldoRequest{SaldoId: 1}

	mockResponse := &response.SaldoResponse{
		ID:           1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockPbResponse := &pb.SaldoResponse{
		SaldoId:      1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockSaldoService.EXPECT().TrashSaldo(1).Return(mockResponse, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponseSaldo(mockResponse).Return(mockPbResponse).Times(1)

	res, err := saldoHandler.TrashedSaldo(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully trashed saldo record", res.GetMessage())
	assert.Equal(t, mockPbResponse, res.GetData())
}

func TestTrashSaldo_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByIdSaldoRequest{SaldoId: 0}

	res, err := saldoHandler.TrashedSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Invalid ID")
}

func TestTrashSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByIdSaldoRequest{SaldoId: 1}

	mockSaldoService.EXPECT().TrashSaldo(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Saldo not found",
	}).Times(1)

	res, err := saldoHandler.TrashedSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to trash saldo record")
}

func TestRestoreSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByIdSaldoRequest{SaldoId: 1}

	mockResponse := &response.SaldoResponse{
		ID:           1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockPbResponse := &pb.SaldoResponse{
		SaldoId:      1,
		CardNumber:   "1234",
		TotalBalance: 10000,
	}

	mockSaldoService.EXPECT().RestoreSaldo(1).Return(mockResponse, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponseSaldo(mockResponse).Return(mockPbResponse).Times(1)

	res, err := saldoHandler.RestoreSaldo(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully restored saldo record", res.GetMessage())
	assert.Equal(t, mockPbResponse, res.GetData())
}

func TestRestoreSaldo_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByIdSaldoRequest{SaldoId: 0}

	res, err := saldoHandler.RestoreSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Invalid ID")
}

func TestRestoreSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByIdSaldoRequest{SaldoId: 1}

	mockSaldoService.EXPECT().RestoreSaldo(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Saldo not found",
	}).Times(1)

	res, err := saldoHandler.RestoreSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to restore saldo record")
}

func TestDeleteSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByIdSaldoRequest{SaldoId: 1}

	mockResponse := &pb.ApiResponseSaldoDelete{
		Status:  "success",
		Message: "Successfully deleted saldo record",
	}

	mockSaldoService.EXPECT().DeleteSaldoPermanent(1).Return(mockResponse, nil).Times(1)

	res, err := saldoHandler.DeleteSaldo(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully deleted saldo record", res.GetMessage())
}

func TestDeleteSaldo_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByIdSaldoRequest{SaldoId: 0}

	res, err := saldoHandler.DeleteSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Invalid ID")
}

func TestDeleteSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaldoService := mock_service.NewMockSaldoService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockSaldoProtoMapper(ctrl)
	saldoHandler := gapi.NewSaldoHandleGrpc(mockSaldoService, mockProtoMapper)

	req := &pb.FindByIdSaldoRequest{SaldoId: 1}

	mockSaldoService.EXPECT().DeleteSaldoPermanent(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to delete saldo record",
	}).Times(1)

	res, err := saldoHandler.DeleteSaldo(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to delete saldo record")
}
