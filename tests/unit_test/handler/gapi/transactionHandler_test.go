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

type TransactionHandleGrpcTestSuite struct {
	suite.Suite
	Ctrl                   *gomock.Controller
	MockTransactionService *mock_service.MockTransactionService
	MockProtoMapper        *mock_protomapper.MockTransactionProtoMapper
	Handler                gapi.TransactionHandleGrpc
}

func (suite *TransactionHandleGrpcTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockTransactionService = mock_service.NewMockTransactionService(suite.Ctrl)
	suite.MockProtoMapper = mock_protomapper.NewMockTransactionProtoMapper(suite.Ctrl)
	suite.Handler = gapi.NewTransactionHandleGrpc(suite.MockTransactionService, suite.MockProtoMapper)
}

func (suite *TransactionHandleGrpcTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *TransactionHandleGrpcTestSuite) TestFindAllTransaction_Success() {
	req := &pb.FindAllTransactionRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockTransactions := []*response.TransactionResponse{
		{ID: 1, TransactionNo: "TX123", CardNumber: "1234567890123456", Amount: 100000, PaymentMethod: "alfamart", MerchantID: 1},
		{ID: 2, TransactionNo: "TX124", CardNumber: "6543210987654321", Amount: 200000, PaymentMethod: "Debit Card", MerchantID: 2},
	}
	mockProtoTransactions := []*pb.TransactionResponse{
		{Id: 1, TransactionNo: "TX123", CardNumber: "1234567890123456", Amount: 100000, PaymentMethod: "alfamart", MerchantId: 1},
		{Id: 2, TransactionNo: "TX124", CardNumber: "6543210987654321", Amount: 200000, PaymentMethod: "Debit Card", MerchantId: 2},
	}

	totalRecords := 2
	suite.MockTransactionService.EXPECT().
		FindAll(gomock.Eq(&requests.FindAllTransactions{Page: 1, PageSize: 10, Search: "test"})).
		Return(mockTransactions, &totalRecords, nil)

	suite.MockProtoMapper.EXPECT().
		ToProtoResponsePaginationTransaction(gomock.Any(), "success", "Successfully fetched transaction records", mockTransactions).
		Return(&pb.ApiResponsePaginationTransaction{
			Status:  "success",
			Message: "Successfully fetched transaction records",
			Data:    mockProtoTransactions,
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   1,
				TotalRecords: 2,
			},
		})

	res, err := suite.Handler.FindAllTransaction(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Successfully fetched transaction records", res.GetMessage())
	suite.Equal(int32(2), res.GetPagination().GetTotalRecords())
	suite.Equal(2, len(res.GetData()))
}

func (suite *TransactionHandleGrpcTestSuite) TestFindAllTransaction_Failure() {
	req := &pb.FindAllTransactionRequest{Page: 1, PageSize: 10, Search: "test"}
	serviceError := &response.ErrorResponse{Status: "error", Message: "Failed to fetch transactions"}

	totalRecords := 0
	suite.MockTransactionService.EXPECT().FindAll(gomock.Any()).Return(nil, &totalRecords, serviceError)

	res, _ := suite.Handler.FindAllTransaction(context.Background(), req)

	suite.Nil(res)
}

func (suite *TransactionHandleGrpcTestSuite) TestFindByIdTransaction_Success() {
	req := &pb.FindByIdTransactionRequest{TransactionId: 1}
	mockTransaction := &response.TransactionResponse{ID: 1, TransactionNo: "TX123", Amount: 100000}
	mockProtoTransaction := &pb.TransactionResponse{Id: 1, TransactionNo: "TX123", Amount: 100000}

	suite.MockTransactionService.EXPECT().FindById(1).Return(mockTransaction, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransaction("success", "Transaction fetched successfully", mockTransaction).Return(&pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Transaction fetched successfully",
		Data:    mockProtoTransaction,
	})

	res, err := suite.Handler.FindByIdTransaction(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("TX123", res.GetData().GetTransactionNo())
}

func (suite *TransactionHandleGrpcTestSuite) TestFindByIdTransaction_InvalidId() {
	req := &pb.FindByIdTransactionRequest{TransactionId: 0}

	res, err := suite.Handler.FindByIdTransaction(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid Transaction ID")
}

func (suite *TransactionHandleGrpcTestSuite) TestFindTransactionByMerchantIdRequest_Success() {
	req := &pb.FindTransactionByMerchantIdRequest{MerchantId: 1}
	mockTransactions := []*response.TransactionResponse{
		{ID: 1, TransactionNo: "TX123", CardNumber: "1234567890123456", Amount: 100000, PaymentMethod: "alfamart", MerchantID: 1},
	}

	suite.MockTransactionService.EXPECT().FindTransactionByMerchantId(1).Return(mockTransactions, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransactions("success", "Successfully fetch transactions", mockTransactions).Return(&pb.ApiResponseTransactions{
		Status:  "success",
		Message: "Successfully fetch transactions",
		Data: []*pb.TransactionResponse{
			{Id: 1, TransactionNo: "TX123", CardNumber: "1234567890123456", Amount: 100000, PaymentMethod: "alfamart", MerchantId: 1},
		},
	})

	res, err := suite.Handler.FindTransactionByMerchantIdRequest(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *TransactionHandleGrpcTestSuite) TestFindByActiveTransaction_Success() {
	req := &pb.FindAllTransactionRequest{Page: 1, PageSize: 10, Search: ""}
	activeTransactions := []*response.TransactionResponseDeleteAt{
		{ID: 1, TransactionNo: "TX123", CardNumber: "1234567890123456", Amount: 100000, PaymentMethod: "alfamart", MerchantID: 1},
	}
	totalRecords := 1
	suite.MockTransactionService.EXPECT().FindByActive(gomock.Any()).Return(activeTransactions, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationTransactionDeleteAt(gomock.Any(), "success", "Successfully fetch transactions", activeTransactions).Return(&pb.ApiResponsePaginationTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully fetch transactions",
		Data: []*pb.TransactionResponseDeleteAt{
			{Id: 1, TransactionNo: "TX123", CardNumber: "1234567890123456", Amount: 100000, PaymentMethod: "alfamart", MerchantId: 1},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByActiveTransaction(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *TransactionHandleGrpcTestSuite) TestFindByTrashedTransaction_Success() {
	req := &pb.FindAllTransactionRequest{Page: 1, PageSize: 10, Search: ""}
	trashedTransactions := []*response.TransactionResponseDeleteAt{
		{ID: 1, TransactionNo: "TX123", CardNumber: "1234567890123456", Amount: 100000, PaymentMethod: "alfamart", MerchantID: 1},
	}
	totalRecords := 1
	suite.MockTransactionService.EXPECT().FindByTrashed(gomock.Any()).Return(trashedTransactions, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationTransactionDeleteAt(gomock.Any(), "success", "Successfully fetch transactions", trashedTransactions).Return(&pb.ApiResponsePaginationTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully fetch transactions",
		Data: []*pb.TransactionResponseDeleteAt{
			{Id: 1, TransactionNo: "TX123", CardNumber: "1234567890123456", Amount: 100000, PaymentMethod: "alfamart", MerchantId: 1},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByTrashedTransaction(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *TransactionHandleGrpcTestSuite) TestCreateTransaction_Success() {
	transactionTime := time.Now()
	req := &pb.CreateTransactionRequest{
		ApiKey:          "hello",
		CardNumber:      "1234567890123456",
		Amount:          150000,
		PaymentMethod:   "alfamart",
		MerchantId:      1,
		TransactionTime: timestamppb.New(transactionTime),
	}

	mockTransaction := &response.TransactionResponse{
		ID:              1,
		TransactionNo:   "TX125",
		CardNumber:      "1234567890123456",
		Amount:          150000,
		PaymentMethod:   "alfamart",
		MerchantID:      1,
		TransactionTime: transactionTime.Format(time.RFC3339),
	}

	mockProtoTransaction := &pb.TransactionResponse{
		Id:              1,
		TransactionNo:   "TX125",
		CardNumber:      "1234567890123456",
		Amount:          150000,
		PaymentMethod:   "alfamart",
		MerchantId:      1,
		TransactionTime: transactionTime.Format(time.RFC3339),
	}

	suite.MockTransactionService.EXPECT().Create("hello", gomock.Any()).Return(mockTransaction, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransaction("success", "Successfully created transaction", mockTransaction).Return(&pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully created transaction",
		Data:    mockProtoTransaction,
	})

	res, err := suite.Handler.CreateTransaction(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal(int32(150000), res.GetData().GetAmount())
}

func (suite *TransactionHandleGrpcTestSuite) TestCreateTransaction_ValidationError() {
	req := &pb.CreateTransactionRequest{
		CardNumber:      "",
		Amount:          40000,
		PaymentMethod:   "",
		MerchantId:      0,
		TransactionTime: nil,
	}

	res, err := suite.Handler.CreateTransaction(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid input for create transaction")
}

func (suite *TransactionHandleGrpcTestSuite) TestUpdateTransaction_Success() {
	transactionId := 1
	transactionTime := time.Now()
	req := &pb.UpdateTransactionRequest{
		ApiKey:          "hello",
		TransactionId:   int32(transactionId),
		CardNumber:      "1234567890123456",
		Amount:          200000,
		PaymentMethod:   "alfamart",
		MerchantId:      1,
		TransactionTime: timestamppb.New(transactionTime),
	}

	mockTransaction := &response.TransactionResponse{
		ID:            1,
		TransactionNo: "TX125",
		Amount:        200000,
	}

	mockProtoTransaction := &pb.TransactionResponse{Id: 1, Amount: 200000}

	suite.MockTransactionService.EXPECT().Update("hello", gomock.Any()).Return(mockTransaction, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransaction("success", "Successfully updated transaction", mockTransaction).Return(&pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully updated transaction",
		Data:    mockProtoTransaction,
	})

	res, err := suite.Handler.UpdateTransaction(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal(int32(200000), res.GetData().GetAmount())
}

func (suite *TransactionHandleGrpcTestSuite) TestTrashedTransaction_Success() {
	req := &pb.FindByIdTransactionRequest{TransactionId: 1}
	mockTransaction := &response.TransactionResponseDeleteAt{ID: 1, TransactionNo: "TX125"}

	suite.MockTransactionService.EXPECT().TrashedTransaction(1).Return(mockTransaction, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransactionDeleteAt("success", "Successfully trashed transaction", mockTransaction).Return(&pb.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully trashed transaction",
		Data:    &pb.TransactionResponseDeleteAt{Id: 1, TransactionNo: "TX125"},
	})

	res, err := suite.Handler.TrashedTransaction(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *TransactionHandleGrpcTestSuite) TestRestoreTransaction_Success() {
	req := &pb.FindByIdTransactionRequest{TransactionId: 1}
	mockTransaction := &response.TransactionResponseDeleteAt{ID: 1, TransactionNo: "TX125"}

	suite.MockTransactionService.EXPECT().RestoreTransaction(1).Return(mockTransaction, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransactionDeleteAt("success", "Successfully restored transaction", mockTransaction).Return(&pb.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully restored transaction",
		Data:    &pb.TransactionResponseDeleteAt{Id: 1, TransactionNo: "TX125"},
	})

	res, err := suite.Handler.RestoreTransaction(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *TransactionHandleGrpcTestSuite) TestDeleteTransaction_Success() {
	req := &pb.FindByIdTransactionRequest{TransactionId: 1}

	suite.MockTransactionService.EXPECT().DeleteTransactionPermanent(1).Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransactionDelete("success", "Successfully deleted transaction").Return(&pb.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted transaction",
	})

	res, err := suite.Handler.DeleteTransaction(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *TransactionHandleGrpcTestSuite) TestRestoreAllTransaction_Success() {
	req := &emptypb.Empty{}

	suite.MockTransactionService.EXPECT().RestoreAllTransaction().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransactionAll("success", "Successfully restore all transaction").Return(&pb.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully restore all transaction",
	})

	res, err := suite.Handler.RestoreAllTransaction(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *TransactionHandleGrpcTestSuite) TestDeleteAllTransactionPermanent_Success() {
	req := &emptypb.Empty{}

	suite.MockTransactionService.EXPECT().DeleteAllTransactionPermanent().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseTransactionAll("success", "Successfully delete transaction permanent").Return(&pb.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully delete transaction permanent",
	})

	res, err := suite.Handler.DeleteAllTransactionPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func TestTransactionHandleGrpcSuite(t *testing.T) {
	suite.Run(t, new(TransactionHandleGrpcTestSuite))
}
