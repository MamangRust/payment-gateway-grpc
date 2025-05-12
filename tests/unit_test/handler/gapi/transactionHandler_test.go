package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/handler/gapi"
	mock_protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto/mocks"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	mock_service "MamangRust/paymentgatewaygrpc/internal/service/mocks"
	"MamangRust/paymentgatewaygrpc/tests/utils"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFindAllTransactions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindAllTransactionRequest{Page: 1, PageSize: 10, Search: "test"}

	transactions := []*response.TransactionResponse{
		{
			ID:         1,
			CardNumber: "1234",
		},
		{
			ID:         2,
			CardNumber: "1234",
		},
	}

	mockTransaction := []*pb.TransactionResponse{
		{
			Id:         1,
			CardNumber: "1234",
		},
		{
			Id:         2,
			CardNumber: "1234",
		},
	}

	totalRecords := 2

	mockTransactionService.EXPECT().FindAll(1, 10, "test").Return(transactions, totalRecords, nil).Times(1)

	mockTransactionMapper.EXPECT().ToResponsesTransaction(transactions).Return(mockTransaction).Times(1)

	res, err := mockHandler.FindAllTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Transactions fetched successfully", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
	assert.Equal(t, int32(1), res.GetPagination().GetTotalPages())
	assert.Equal(t, int32(2), res.GetPagination().GetTotalRecords())
}

func TestFindAllTransactions_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindAllTransactionRequest{Page: 1, PageSize: 10, Search: "test"}

	mockTransactionService.EXPECT().FindAll(1, 10, "test").Return(nil, 0, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch transactions: database error",
	}).Times(1)

	res, err := mockHandler.FindAllTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch transactions: database error")
}

func TestFindAllTransactions_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindAllTransactionRequest{Page: 1, PageSize: 10, Search: "test"}
	mockTransaction := []*response.TransactionResponse{}
	mockProto := []*pb.TransactionResponse{}

	mockTransactionService.EXPECT().FindAll(1, 10, "test").Return(mockTransaction, 0, nil).Times(1)

	mockTransactionMapper.EXPECT().ToResponsesTransaction(mockTransaction).Return(mockProto).Times(1)

	res, err := mockHandler.FindAllTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Transactions fetched successfully", res.GetMessage())
	assert.Len(t, res.GetData(), 0)
	assert.Equal(t, int32(0), res.GetPagination().GetTotalPages())
	assert.Equal(t, int32(0), res.GetPagination().GetTotalRecords())
}

func TestFindTransactionById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindByIdTransactionRequest{TransactionId: 1}

	transaction := &response.TransactionResponse{
		ID:         1,
		CardNumber: "1234",
	}

	mockTransaction := &pb.TransactionResponse{
		Id:         1,
		CardNumber: "1234",
	}

	mockTransactionService.EXPECT().FindById(1).Return(transaction, nil).Times(1)
	mockTransactionMapper.EXPECT().ToResponseTransaction(transaction).Return(mockTransaction).Times(1)

	res, err := mockHandler.FindTransactionById(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int32(1), res.GetId())
	assert.Equal(t, "1234", res.GetCardNumber())
}

func TestFindTransactionById_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindByIdTransactionRequest{TransactionId: -1}

	res, err := mockHandler.FindTransactionById(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Bad Request: Invalid ID")
}

func TestFindTransactionById_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindByIdTransactionRequest{TransactionId: 1}

	mockTransactionService.EXPECT().FindById(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch transaction: transaction not found",
	}).Times(1)

	res, err := mockHandler.FindTransactionById(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch transaction: transaction not found")
}

func TestFindByCardNumberTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindByCardNumberTransactionRequest{CardNumber: "1234"}

	transactions := []*response.TransactionResponse{
		{
			ID:         1,
			CardNumber: "1234",
		},
		{
			ID:         2,
			CardNumber: "1234",
		},
	}

	mockTransaction := []*pb.TransactionResponse{
		{
			Id:         1,
			CardNumber: "1234",
		},
		{
			Id:         2,
			CardNumber: "1234",
		},
	}

	mockTransactionService.EXPECT().FindByCardNumber("1234").Return(transactions, nil).Times(1)
	mockTransactionMapper.EXPECT().ToResponsesTransaction(transactions).Return(mockTransaction).Times(1)

	res, err := mockHandler.FindByCardNumberTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transactions", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
}

func TestFindByCardNumberTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindByCardNumberTransactionRequest{CardNumber: "1234"}

	mockTransactionService.EXPECT().FindByCardNumber("1234").Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch transactions: transactions not found",
	}).Times(1)

	res, err := mockHandler.FindByCardNumberTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch transactions: transactions not found")
}

func TestFindTransactionByMerchantIdRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindTransactionByMerchantIdRequest{MerchantId: 1}

	transactions := []*response.TransactionResponse{
		{
			ID:         1,
			CardNumber: "1234",
		},
		{
			ID:         2,
			CardNumber: "5678",
		},
	}

	mockTransaction := []*pb.TransactionResponse{
		{
			Id:         1,
			CardNumber: "1234",
		},
		{
			Id:         2,
			CardNumber: "5678",
		},
	}

	mockTransactionService.EXPECT().FindTransactionByMerchantId(1).Return(transactions, nil).Times(1)
	mockTransactionMapper.EXPECT().ToResponsesTransaction(transactions).Return(mockTransaction).Times(1)

	res, err := mockHandler.FindTransactionByMerchantIdRequest(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transactions", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
	assert.Equal(t, int32(1), res.GetData()[0].GetId())
	assert.Equal(t, int32(2), res.GetData()[1].GetId())
}

func TestFindTransactionByMerchantIdRequest_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindTransactionByMerchantIdRequest{MerchantId: 0}

	res, err := mockHandler.FindTransactionByMerchantIdRequest(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Bad Request: Invalid ID")
}

func TestFindTransactionByMerchantIdRequest_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindTransactionByMerchantIdRequest{MerchantId: 1}

	mockTransactionService.EXPECT().FindTransactionByMerchantId(1).Return([]*response.TransactionResponse{}, nil).Times(1)
	mockTransactionMapper.EXPECT().ToResponsesTransaction([]*response.TransactionResponse{}).Return([]*pb.TransactionResponse{}).Times(1)

	res, err := mockHandler.FindTransactionByMerchantIdRequest(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transactions", res.GetMessage())
	assert.Len(t, res.GetData(), 0)
}

func TestFindTransactionByMerchantIdRequest_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindTransactionByMerchantIdRequest{MerchantId: 1}

	mockTransactionService.EXPECT().FindTransactionByMerchantId(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch transactions",
	}).Times(1)

	res, err := mockHandler.FindTransactionByMerchantIdRequest(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
	assert.Contains(t, err.Error(), "Failed to fetch transactions")
}

func TestFindByActiveTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	transactions := []*response.TransactionResponseDeleteAt{
		{ID: 1, CardNumber: "1234"},
		{ID: 2, CardNumber: "5678"},
	}
	mockTransaction := []*pb.TransactionResponseDeleteAt{
		{Id: 1, CardNumber: "1234"},
		{Id: 2, CardNumber: "5678"},
	}

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockTransactionService.EXPECT().FindByActive(pageSize, page, search).Return(transactions, expected, nil).Times(1)
	mockTransactionMapper.EXPECT().ToResponsesTransactionDeleteAt(transactions).Return(mockTransaction).Times(1)

	res, err := mockHandler.FindByActiveTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transactions", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
}

func TestFindByActiveTransaction_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockTransactionService.EXPECT().FindByActive(pageSize, page, search).Return([]*response.TransactionResponseDeleteAt{}, expected, nil).Times(1)
	mockTransactionMapper.EXPECT().ToResponsesTransactionDeleteAt([]*response.TransactionResponseDeleteAt{}).Return([]*pb.TransactionResponseDeleteAt{}).Times(1)

	res, err := mockHandler.FindByActiveTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transactions", res.GetMessage())
	assert.Len(t, res.GetData(), 0)
}

func TestFindByActiveTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, nil)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockTransactionService.EXPECT().FindByActive(pageSize, page, search).Return(nil, expected, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch transactions: database error",
	}).Times(1)

	res, err := mockHandler.FindByActiveTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch transactions: database error")
}

func TestFindByTrashedTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	transactions := []*response.TransactionResponseDeleteAt{
		{ID: 3, CardNumber: "9012"},
		{ID: 4, CardNumber: "3456"},
	}
	mockTransaction := []*pb.TransactionResponseDeleteAt{
		{Id: 3, CardNumber: "9012"},
		{Id: 4, CardNumber: "3456"},
	}

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockTransactionService.EXPECT().FindByTrashed(pageSize, page, search).Return(transactions, expected, nil).Times(1)
	mockTransactionMapper.EXPECT().ToResponsesTransactionDeleteAt(transactions).Return(mockTransaction).Times(1)

	res, err := mockHandler.FindByTrashedTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transactions", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
}

func TestFindByTrashedTransaction_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockTransactionService.EXPECT().FindByTrashed(pageSize, page, search).Return([]*response.TransactionResponseDeleteAt{}, expected, nil).Times(1)
	mockTransactionMapper.EXPECT().ToResponsesTransactionDeleteAt([]*response.TransactionResponseDeleteAt{}).Return([]*pb.TransactionResponseDeleteAt{}).Times(1)

	res, err := mockHandler.FindByTrashedTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetch transactions", res.GetMessage())
	assert.Len(t, res.GetData(), 0)
}

func TestFindByTrashedTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, nil)

	search := ""
	pageSize := 1
	page := 1
	expected := 1

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockTransactionService.EXPECT().FindByTrashed(pageSize, page, search).Return(nil, expected, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch transactions: database error",
	}).Times(1)

	res, err := mockHandler.FindByTrashedTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to fetch transactions: database error")
}

func TestCreateTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.CreateTransactionRequest{
		ApiKey:          "test-api-key",
		CardNumber:      "1234-5678-9012-3456",
		Amount:          1000,
		PaymentMethod:   "credit_card",
		MerchantId:      1,
		TransactionTime: timestamppb.Now(),
	}

	serviceReq := &requests.CreateTransactionRequest{
		CardNumber:      "1234-5678-9012-3456",
		Amount:          1000,
		PaymentMethod:   "credit_card",
		MerchantID:      utils.PtrInt(1),
		TransactionTime: req.GetTransactionTime().AsTime(),
	}

	serviceRes := &response.TransactionResponse{
		ID:              1,
		CardNumber:      "1234-5678-9012-3456",
		Amount:          1000,
		PaymentMethod:   "credit_card",
		TransactionTime: req.GetTransactionTime().String(),
	}

	pbRes := &pb.TransactionResponse{
		Id:              1,
		CardNumber:      "1234-5678-9012-3456",
		Amount:          1000,
		PaymentMethod:   "credit_card",
		TransactionTime: req.GetTransactionTime().String(),
	}

	mockTransactionService.EXPECT().Create("test-api-key", serviceReq).Return(serviceRes, nil).Times(1)
	mockTransactionMapper.EXPECT().ToResponseTransaction(serviceRes).Return(pbRes).Times(1)

	res, err := mockHandler.CreateTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully created transaction", res.GetMessage())
	assert.Equal(t, pbRes, res.GetData())
}

func TestCreateTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, nil)

	req := &pb.CreateTransactionRequest{
		ApiKey:          "test-api-key",
		CardNumber:      "1234-5678-9012-3456",
		Amount:          1000,
		PaymentMethod:   "credit_card",
		MerchantId:      1,
		TransactionTime: timestamppb.Now(),
	}

	serviceReq := &requests.CreateTransactionRequest{
		CardNumber:      "1234-5678-9012-3456",
		Amount:          1000,
		PaymentMethod:   "credit_card",
		MerchantID:      utils.PtrInt(1),
		TransactionTime: req.GetTransactionTime().AsTime(),
	}

	mockTransactionService.EXPECT().Create("test-api-key", serviceReq).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to create transaction: database error",
	}).Times(1)

	res, err := mockHandler.CreateTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to create transaction: database error")
}

func TestCreateTransaction_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, nil)

	req := &pb.CreateTransactionRequest{
		ApiKey:          "test-api-key",
		CardNumber:      "",
		Amount:          -1000,
		PaymentMethod:   "",
		MerchantId:      0,
		TransactionTime: timestamppb.Now(),
	}

	serviceReq := &requests.CreateTransactionRequest{
		CardNumber:      "",
		Amount:          -1000,
		PaymentMethod:   "",
		MerchantID:      utils.PtrInt(0),
		TransactionTime: req.GetTransactionTime().AsTime(),
	}

	mockTransactionService.EXPECT().Create("test-api-key", serviceReq).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Validation failed",
	}).Times(1)

	res, err := mockHandler.CreateTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Validation failed")
}

func TestUpdateTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.UpdateTransactionRequest{
		ApiKey:          "test-api-key",
		TransactionId:   1,
		CardNumber:      "1234-5678-9012-3456",
		Amount:          2000,
		PaymentMethod:   "debit_card",
		MerchantId:      2,
		TransactionTime: timestamppb.Now(),
	}

	serviceReq := &requests.UpdateTransactionRequest{
		TransactionID:   1,
		CardNumber:      "1234-5678-9012-3456",
		Amount:          2000,
		PaymentMethod:   "debit_card",
		MerchantID:      utils.PtrInt(2),
		TransactionTime: req.GetTransactionTime().AsTime(),
	}

	serviceRes := &response.TransactionResponse{
		ID:              1,
		CardNumber:      "1234-5678-9012-3456",
		Amount:          2000,
		PaymentMethod:   "debit_card",
		TransactionTime: req.GetTransactionTime().String(),
	}

	pbRes := &pb.TransactionResponse{
		Id:              1,
		CardNumber:      "1234-5678-9012-3456",
		Amount:          2000,
		PaymentMethod:   "debit_card",
		TransactionTime: req.GetTransactionTime().String(),
	}

	mockTransactionService.EXPECT().Update("test-api-key", serviceReq).Return(serviceRes, nil).Times(1)
	mockTransactionMapper.EXPECT().ToResponseTransaction(serviceRes).Return(pbRes).Times(1)

	res, err := mockHandler.UpdateTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully updated transaction", res.GetMessage())
	assert.Equal(t, pbRes, res.GetData())
}

func TestUpdateTransaction_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.UpdateTransactionRequest{
		ApiKey:          "test-api-key",
		TransactionId:   0,
		CardNumber:      "1234-5678-9012-3456",
		Amount:          2000,
		PaymentMethod:   "debit_card",
		MerchantId:      2,
		TransactionTime: timestamppb.Now(),
	}

	res, err := mockHandler.UpdateTransaction(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Bad Request: Invalid ID")
}

func TestUpdateTransaction_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, nil)

	req := &pb.UpdateTransactionRequest{
		ApiKey:          "test-api-key",
		TransactionId:   1,
		CardNumber:      "1234-5678-9012-3456",
		Amount:          2000,
		PaymentMethod:   "debit_card",
		MerchantId:      2,
		TransactionTime: timestamppb.Now(),
	}

	serviceReq := &requests.UpdateTransactionRequest{
		TransactionID:   1,
		CardNumber:      "1234-5678-9012-3456",
		Amount:          2000,
		PaymentMethod:   "debit_card",
		MerchantID:      utils.PtrInt(2),
		TransactionTime: req.GetTransactionTime().AsTime(),
	}

	mockTransactionService.EXPECT().Update("test-api-key", serviceReq).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to update transaction: database error",
	}).Times(1)

	res, err := mockHandler.UpdateTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to update transaction: database error")
}

func TestUpdateTransaction_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, nil)

	req := &pb.UpdateTransactionRequest{
		TransactionId:   1,
		ApiKey:          "test-api-key",
		CardNumber:      "",
		Amount:          -1000,
		PaymentMethod:   "",
		MerchantId:      0,
		TransactionTime: timestamppb.Now(),
	}

	serviceReq := &requests.UpdateTransactionRequest{
		TransactionID:   1,
		CardNumber:      "",
		Amount:          -1000,
		PaymentMethod:   "",
		MerchantID:      utils.PtrInt(0),
		TransactionTime: req.GetTransactionTime().AsTime(),
	}

	mockTransactionService.EXPECT().Update("test-api-key", serviceReq).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Validation failed",
	}).Times(1)

	res, err := mockHandler.UpdateTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Validation failed")
}

func TestTrashedTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindByIdTransactionRequest{TransactionId: 1}

	transaction := &response.TransactionResponse{
		ID:         1,
		CardNumber: "1234",
	}

	mockTransactionService.EXPECT().TrashedTransaction(1).Return(transaction, nil).Times(1)
	mockTransactionMapper.EXPECT().ToResponseTransaction(transaction).Return(&pb.TransactionResponse{
		Id:         1,
		CardNumber: "1234",
	}).Times(1)

	res, err := mockHandler.TrashedTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully trashed transaction", res.GetMessage())
}

func TestTrashedTransaction_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, nil)

	req := &pb.FindByIdTransactionRequest{TransactionId: 0}

	res, err := mockHandler.TrashedTransaction(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Bad Request: Invalid ID")
}

func TestTrashedTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindByIdTransactionRequest{TransactionId: 1}

	mockTransactionService.EXPECT().TrashedTransaction(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to trash transaction",
	}).Times(1)

	res, err := mockHandler.TrashedTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to trash transaction")
}

func TestRestoreTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindByIdTransactionRequest{TransactionId: 1}

	transaction := &response.TransactionResponse{
		ID:         1,
		CardNumber: "1234",
	}

	mockTransactionService.EXPECT().RestoreTransaction(1).Return(transaction, nil).Times(1)
	mockTransactionMapper.EXPECT().ToResponseTransaction(transaction).Return(&pb.TransactionResponse{
		Id:         1,
		CardNumber: "1234",
	}).Times(1)

	res, err := mockHandler.RestoreTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully restored transaction", res.GetMessage())
}

func TestRestoreTransaction_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, nil)

	req := &pb.FindByIdTransactionRequest{TransactionId: 0}

	res, err := mockHandler.RestoreTransaction(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Bad Request: Invalid ID")
}

func TestRestoreTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindByIdTransactionRequest{TransactionId: 1}

	mockTransactionService.EXPECT().RestoreTransaction(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to restore transaction",
	}).Times(1)

	res, err := mockHandler.RestoreTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to restore transaction")
}

func TestDeleteTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindByIdTransactionRequest{TransactionId: 1}

	mockResponse := &pb.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted transaction",
	}

	mockTransactionService.EXPECT().DeleteTransactionPermanent(1).Return(mockResponse, nil).Times(1)

	res, err := mockHandler.DeleteTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully deleted transaction", res.GetMessage())
}

func TestDeleteTransaction_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, nil)

	req := &pb.FindByIdTransactionRequest{TransactionId: 0}

	res, err := mockHandler.DeleteTransaction(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Bad Request: Invalid ID")
}

func TestDeleteTransaction_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransactionService(ctrl)
	mockTransactionMapper := mock_protomapper.NewMockTransactionProtoMapper(ctrl)
	mockHandler := gapi.NewTransactionHandleGrpc(mockTransactionService, mockTransactionMapper)

	req := &pb.FindByIdTransactionRequest{TransactionId: 1}

	mockTransactionService.EXPECT().DeleteTransactionPermanent(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to delete transaction permanently",
	}).Times(1)

	res, err := mockHandler.DeleteTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "Failed to delete transaction permanently")
}
