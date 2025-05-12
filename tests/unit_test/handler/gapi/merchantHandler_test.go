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

func TestFindAllMerchants_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	search := ""
	pageSize := 1
	page := 1

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockMerchants := []*response.MerchantResponse{
		{
			ID:   1,
			Name: "Merchant A",
		},
		{
			ID:   2,
			Name: "Merchant B",
		},
	}

	mockProtoMerchants := []*pb.MerchantResponse{
		{
			Id:   1,
			Name: "Merchant A",
		},
		{
			Id:   2,
			Name: "Merchant B",
		},
	}

	mockMerchantService.EXPECT().
		FindAll(pageSize, page, search).
		Return(mockMerchants, 2, nil).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponsesMerchant(mockMerchants).
		Return(mockProtoMerchants).
		Times(1)

	response, err := merchantHandler.FindAllMerchant(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "success", response.GetStatus())
	assert.Equal(t, "Successfully fetched merchant record", response.GetMessage())
	assert.NotNil(t, response.GetData())
	assert.Equal(t, int32(1), response.GetPagination().GetCurrentPage())
	assert.Equal(t, int32(1), response.GetPagination().GetPageSize())
	assert.Equal(t, int32(2), response.GetPagination().GetTotalRecords())
	assert.Equal(t, int32(2), response.GetPagination().GetTotalPages())
}

func TestFindAllMerchants_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, nil)

	req := &pb.FindAllMerchantRequest{
		Page:     1,
		PageSize: 10,
		Search:   "electronics",
	}

	mockMerchantService.EXPECT().
		FindAll(1, 10, "electronics").
		Return(nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch card records",
		}).
		Times(1)

	response, err := merchantHandler.FindAllMerchant(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
	assert.Contains(t, statusErr.Message(), "Failed to fetch card records")
}

func TestFindAllMerchants_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindAllMerchantRequest{
		Page:     1,
		PageSize: 10,
		Search:   "nonexistent",
	}

	mockMerchants := []*response.MerchantResponse{}
	mockProtoMerchants := []*pb.MerchantResponse{}

	mockMerchantService.EXPECT().
		FindAll(1, 10, "nonexistent").
		Return(mockMerchants, 0, nil).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponsesMerchant(mockMerchants).
		Return(mockProtoMerchants).
		Times(1)

	response, err := merchantHandler.FindAllMerchant(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "success", response.GetStatus())
	assert.Equal(t, "Successfully fetched merchant record", response.GetMessage())
	assert.Empty(t, response.GetData())
	assert.Equal(t, int32(1), response.GetPagination().GetCurrentPage())
	assert.Equal(t, int32(10), response.GetPagination().GetPageSize())
	assert.Equal(t, int32(0), response.GetPagination().GetTotalRecords())
	assert.Equal(t, int32(0), response.GetPagination().GetTotalPages())
}

func TestFindByIdMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByIdMerchantRequest{
		MerchantId: 1,
	}

	mockMerchant := &response.MerchantResponse{
		ID:     1,
		Name:   "Merchant One",
		ApiKey: "api-key-123",
	}

	mockProtoMerchant := &pb.MerchantResponse{
		Id:     1,
		Name:   "Merchant One",
		ApiKey: "api-key-123",
	}

	mockMerchantService.EXPECT().
		FindById(1).
		Return(mockMerchant, nil).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponseMerchant(mockMerchant).
		Return(mockProtoMerchant).
		Times(1)

	response, err := merchantHandler.FindByIdMerchant(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Successfully fetched merchant record", response.Message)
	assert.Equal(t, mockProtoMerchant, response.Data)
}

func TestFindByIdMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByIdMerchantRequest{
		MerchantId: 999,
	}

	mockMerchantService.EXPECT().
		FindById(999).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: ",
		}).
		Times(1)

	response, err := merchantHandler.FindByIdMerchant(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Merchant not found: ")
}

func TestFindByIdMerchant_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByIdMerchantRequest{
		MerchantId: 0,
	}

	response, err := merchantHandler.FindByIdMerchant(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Invalid merchant ID")
}

func TestFindByApiKey_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByApiKeyRequest{
		ApiKey: "api-key-123",
	}

	mockMerchant := &response.MerchantResponse{
		ID:     1,
		Name:   "Merchant One",
		ApiKey: "api-key-123",
	}

	mockProtoMerchant := &pb.MerchantResponse{
		Id:     1,
		Name:   "Merchant One",
		ApiKey: "api-key-123",
	}

	mockMerchantService.EXPECT().
		FindByApiKey("api-key-123").
		Return(mockMerchant, nil).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponseMerchant(mockMerchant).
		Return(mockProtoMerchant).
		Times(1)

	response, err := merchantHandler.FindByApiKey(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Successfully fetched merchant record", response.Message)
	assert.Equal(t, mockProtoMerchant, response.Data)
}

func TestFindByApiKey_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByApiKeyRequest{
		ApiKey: "invalid-api-key",
	}

	mockMerchantService.EXPECT().
		FindByApiKey("invalid-api-key").
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: ",
		}).
		Times(1)

	response, err := merchantHandler.FindByApiKey(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Merchant not found: ")
}

func TestFindByMerchantUserId_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByMerchantUserIdRequest{
		UserId: 1,
	}

	mockMerchants := []*response.MerchantResponse{
		{
			ID:     1,
			Name:   "Merchant One",
			ApiKey: "api-key-123",
		},
		{
			ID:     2,
			Name:   "Merchant Two",
			ApiKey: "api-key-456",
		},
	}

	mockProtoMerchants := []*pb.MerchantResponse{
		{
			Id:     1,
			Name:   "Merchant One",
			ApiKey: "api-key-123",
		},
		{
			Id:     2,
			Name:   "Merchant Two",
			ApiKey: "api-key-456",
		},
	}

	mockMerchantService.EXPECT().
		FindByMerchantUserId(1).
		Return(mockMerchants, nil).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponsesMerchant(mockMerchants).
		Return(mockProtoMerchants).
		Times(1)

	response, err := merchantHandler.FindByMerchantUserId(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Successfully fetched merchant record", response.Message)
	assert.Equal(t, mockProtoMerchants, response.Data)

}

func TestFindByMerchantUserId_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByMerchantUserIdRequest{
		UserId: 1,
	}

	mockMerchantService.EXPECT().
		FindByMerchantUserId(1).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: ",
		}).
		Times(1)

	response, err := merchantHandler.FindByMerchantUserId(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Merchant not found: ")
}

func TestFindByMerchantUserId_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByMerchantUserIdRequest{
		UserId: 0,
	}

	mockMerchantService.EXPECT().
		FindByMerchantUserId(0).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: ",
		}).
		Times(1)

	response, err := merchantHandler.FindByMerchantUserId(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Merchant not found: ")
}

func TestFindByActiveMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	mockMerchants := []*response.MerchantResponseDeleteAt{
		{
			ID:     1,
			Name:   "Merchant One",
			ApiKey: "api-key-123",
		},
		{
			ID:     2,
			Name:   "Merchant Two",
			ApiKey: "api-key-456",
		},
	}

	mockProtoMerchants := []*pb.MerchantResponseDeleteAt{
		{
			Id:     1,
			Name:   "Merchant One",
			ApiKey: "api-key-123",
		},
		{
			Id:     2,
			Name:   "Merchant Two",
			ApiKey: "api-key-456",
		},
	}

	search := ""
	pageSize := 1
	page := 1
	expected := 2

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockMerchantService.EXPECT().
		FindByActive(pageSize, page, search).
		Return(mockMerchants, expected, nil).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponsesMerchantDeleteAt(mockMerchants).
		Return(mockProtoMerchants).
		Times(1)

	res, err := merchantHandler.FindByActive(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully fetched merchant record", res.Message)
	assert.Equal(t, mockProtoMerchants, res.Data)

}

func TestFindByActiveMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	search := ""
	pageSize := 1
	page := 1
	expected := 0

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockMerchantService.EXPECT().
		FindByActive(pageSize, page, search).
		Return(nil, expected, &response.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: ",
		}).
		Times(1)

	res, err := merchantHandler.FindByActive(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Merchant not found: ")
}

func TestCreateMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.CreateMerchantRequest{
		Name:   "Merchant One",
		UserId: 1,
	}

	mockRequest := &requests.CreateMerchantRequest{
		Name:   "Merchant One",
		UserID: 1,
	}

	mockMerchant := &response.MerchantResponse{
		ID:     1,
		Name:   "Merchant One",
		ApiKey: "api-key-123",
	}

	mockProtoMerchant := &pb.MerchantResponse{
		Id:     1,
		Name:   "Merchant One",
		ApiKey: "api-key-123",
	}

	mockMerchantService.EXPECT().
		CreateMerchant(mockRequest).
		Return(mockMerchant, nil).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponseMerchant(mockMerchant).
		Return(mockProtoMerchant).
		Times(1)

	res, err := merchantHandler.CreateMerchant(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully created merchant", res.Message)
	assert.Equal(t, mockProtoMerchant, res.Data)
}

func TestCreateMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.CreateMerchantRequest{
		Name:   "Merchant One",
		UserId: 1,
	}

	mockRequest := &requests.CreateMerchantRequest{
		Name:   "Merchant One",
		UserID: 1,
	}

	mockMerchantService.EXPECT().
		CreateMerchant(mockRequest).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create merchant: ",
		}).
		Times(1)

	res, err := merchantHandler.CreateMerchant(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to create merchant: ")
}

func TestCreateMerchant_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.CreateMerchantRequest{
		Name:   "",
		UserId: 1,
	}

	res, err := merchantHandler.CreateMerchant(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
}

func TestUpdateMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	mockRequest := &requests.UpdateMerchantRequest{
		MerchantID: 1,
		Name:       "Merchant One",
		UserID:     1,
		Status:     "active",
	}

	req := &pb.UpdateMerchantRequest{
		MerchantId: 1,
		Name:       "Merchant One",
		UserId:     1,
		Status:     "active",
	}

	mockResponse := &response.MerchantResponse{
		ID:     1,
		Name:   "Merchant One",
		ApiKey: "api-key-123",
	}

	mockMerchantService.EXPECT().UpdateMerchant(mockRequest).Return(mockResponse, nil)

	mockProtoMapper.EXPECT().ToResponseMerchant(mockResponse).Return(
		&pb.MerchantResponse{
			Id:     1,
			Name:   "Merchant One",
			ApiKey: "api-key-123",
		},
	)

	res, err := merchantHandler.UpdateMerchant(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully updated merchant", res.GetMessage())

}

func TestUpdateMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.UpdateMerchantRequest{
		MerchantId: 1,
		Name:       "Merchant One",
		UserId:     1,
		Status:     "active",
	}
	mockError := &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to update merchant",
	}

	mockRequest := &requests.UpdateMerchantRequest{
		MerchantID: 1,
		Name:       "Merchant One",
		UserID:     1,
		Status:     "active",
	}

	mockMerchantService.EXPECT().UpdateMerchant(mockRequest).Return(nil, mockError)

	res, err := merchantHandler.UpdateMerchant(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to update merchant")
}

func TestUpdateMerchant_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.UpdateMerchantRequest{
		MerchantId: 1,
		Name:       "",
		UserId:     1,
		Status:     "active",
	}

	res, err := merchantHandler.UpdateMerchant(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
}

func TestTrashedMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByIdMerchantRequest{MerchantId: 1}
	mockMerchant := &response.MerchantResponse{
		ID:     1,
		Name:   "Test Merchant",
		ApiKey: "api-key-123",
	}

	mockMerchantService.EXPECT().TrashedMerchant(1).Return(mockMerchant, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponseMerchant(mockMerchant).Return(&pb.MerchantResponse{
		Id:   1,
		Name: "Test Merchant",
	}).Times(1)

	res, err := merchantHandler.TrashedMerchant(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully trashed merchant", res.Message)
	assert.Equal(t, int32(1), res.Data.Id)
	assert.Equal(t, "Test Merchant", res.Data.Name)
}

func TestTrashedMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByIdMerchantRequest{MerchantId: 1}

	mockMerchantService.EXPECT().TrashedMerchant(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Merchant not found",
	}).Times(1)

	res, err := merchantHandler.TrashedMerchant(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Merchant not found")
}

func TestTrashedMerchant_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByIdMerchantRequest{MerchantId: 0}

	res, err := merchantHandler.TrashedMerchant(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "merchant id is required")
}

func TestRestoreMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByIdMerchantRequest{MerchantId: 1}
	mockMerchant := &response.MerchantResponse{ID: 1, Name: "Test Merchant"}

	mockMerchantService.EXPECT().RestoreMerchant(1).Return(mockMerchant, nil).Times(1)
	mockProtoMapper.EXPECT().ToResponseMerchant(mockMerchant).Return(&pb.MerchantResponse{
		Id:   1,
		Name: "Test Merchant",
	}).Times(1)

	res, err := merchantHandler.RestoreMerchant(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "Successfully restored merchant", res.Message)
	assert.Equal(t, int32(1), res.Data.Id)
	assert.Equal(t, "Test Merchant", res.Data.Name)
}

func TestRestoreMerchant_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByIdMerchantRequest{MerchantId: 1}

	mockMerchantService.EXPECT().RestoreMerchant(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Merchant not found",
	}).Times(1)

	res, err := merchantHandler.RestoreMerchant(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Merchant not found")
}

func TestDeleteMerchant_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByIdMerchantRequest{MerchantId: 1}

	mockResponse := &pb.ApiResponseMerchatDelete{
		Status:  "success",
		Message: "Successfully deleted merchant",
	}

	mockMerchantService.EXPECT().DeleteMerchantPermanent(1).Return(mockResponse, nil).Times(1)

	res, err := merchantHandler.DeleteMerchant(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully deleted merchant", res.Message)
}

func TestDeleteMerchant_Failure_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMerchantService := mock_service.NewMockMerchantService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockMerchantProtoMapper(ctrl)
	merchantHandler := gapi.NewMerchantHandleGrpc(mockMerchantService, mockProtoMapper)

	req := &pb.FindByIdMerchantRequest{MerchantId: 1}

	mockMerchantService.EXPECT().DeleteMerchantPermanent(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Merchant not found",
	}).Times(1)

	res, err := merchantHandler.DeleteMerchant(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Merchant not found")
}
