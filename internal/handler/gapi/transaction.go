package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"context"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionHandleGrpc struct {
	pb.UnimplementedTransactionServiceServer
	transactionService service.TransactionService
	mapping            protomapper.TransactionProtoMapper
}

func NewTransactionHandleGrpc(transactionService service.TransactionService, mapping protomapper.TransactionProtoMapper) *transactionHandleGrpc {
	return &transactionHandleGrpc{
		transactionService: transactionService,
		mapping:            mapping,
	}
}

func (t *transactionHandleGrpc) FindAllTransaction(ctx context.Context, request *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransaction, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := t.transactionService.FindAll(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transactions: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
	so := t.mapping.ToProtoResponsePaginationTransaction(paginationMeta, "", "", transactions)

	return so, nil
}

func (t *transactionHandleGrpc) FindAllTransactionByCardNumber(ctx context.Context, request *pb.FindAllTransactionCardNumberRequest) (*pb.ApiResponsePaginationTransaction, error) {
	card_number := request.GetCardNumber()
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := t.transactionService.FindAllByCardNumber(card_number, page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transactions: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
	so := t.mapping.ToProtoResponsePaginationTransaction(paginationMeta, "", "", transactions)

	return so, nil
}

func (t *transactionHandleGrpc) FindByIdTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	if request.GetTransactionId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid ID",
		})
	}

	id := request.GetTransactionId()

	transaction, err := t.transactionService.FindById(int(id))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transaction: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransaction("success", "Transaction fetched successfully", transaction)

	return so, nil
}

func (s *transactionHandleGrpc) FindMonthlyTransactionStatusSuccess(ctx context.Context, req *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthStatusSuccess, error) {
	if req.GetYear() <= 0 || req.GetMonth() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year or month",
		})
	}

	year := req.GetYear()
	month := req.GetMonth()

	records, errResponse := s.transactionService.FindMonthTransactionStatusSuccess(int(year), int(month))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Transaction status success: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTransactionMonthStatusSuccess("success", "Successfully fetched monthly Transaction status success", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindYearlyTransactionStatusSuccess(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearStatusSuccess, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year",
		})
	}

	year := req.GetYear()

	records, errResponse := s.transactionService.FindYearlyTransactionStatusSuccess(int(year))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Transaction status success: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTransactionYearStatusSuccess("success", "Successfully fetched yearly Transaction status success", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindMonthlyTransactionStatusFailed(ctx context.Context, req *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthStatusFailed, error) {
	if req.GetYear() <= 0 || req.GetMonth() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year or month",
		})
	}

	year := req.GetYear()
	month := req.GetMonth()

	records, errResponse := s.transactionService.FindMonthTransactionStatusFailed(int(year), int(month))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Transaction status Failed: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTransactionMonthStatusFailed("success", "success fetched monthly Transaction status Failed", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindYearlyTransactionStatusFailed(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearStatusFailed, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year",
		})
	}

	year := req.GetYear()

	records, errResponse := s.transactionService.FindYearlyTransactionStatusFailed(int(year))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Transaction status Failed: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTransactionYearStatusFailed("success", "success fetched yearly Transaction status Failed", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindMonthlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTransactionStatusCardNumber) (*pb.ApiResponseTransactionMonthStatusSuccess, error) {
	if req.GetYear() <= 0 || req.GetMonth() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year or month",
		})
	}

	year := req.GetYear()
	month := req.GetMonth()
	cardNumber := req.GetCardNumber()

	records, errResponse := s.transactionService.FindMonthTransactionStatusSuccessByCardNumber(cardNumber, int(year), int(month))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Transaction status success: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTransactionMonthStatusSuccess("success", "Successfully fetched monthly Transaction status success", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTransactionStatusCardNumber) (*pb.ApiResponseTransactionYearStatusSuccess, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year",
		})
	}

	year := req.GetYear()
	cardNumber := req.GetCardNumber()

	records, errResponse := s.transactionService.FindYearlyTransactionStatusSuccessByCardNumber(cardNumber, int(year))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Transaction status success: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTransactionYearStatusSuccess("success", "Successfully fetched yearly Transaction status success", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindMonthlyTransactionStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTransactionStatusCardNumber) (*pb.ApiResponseTransactionMonthStatusFailed, error) {
	if req.GetYear() <= 0 || req.GetMonth() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year or month",
		})
	}

	year := req.GetYear()
	month := req.GetMonth()
	cardNumber := req.GetCardNumber()

	records, errResponse := s.transactionService.FindMonthTransactionStatusFailedByCardNumber(cardNumber, int(year), int(month))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Transaction status Failed: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTransactionMonthStatusFailed("success", "success fetched monthly Transaction status Failed", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTransactionStatusCardNumber) (*pb.ApiResponseTransactionYearStatusFailed, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year",
		})
	}

	year := req.GetYear()
	cardNumber := req.GetCardNumber()

	records, errResponse := s.transactionService.FindYearlyTransactionStatusFailedByCardNumber(cardNumber, int(year))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Transaction status Failed: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTransactionYearStatusFailed("success", "success fetched yearly Transaction status Failed", records)

	return so, nil
}

func (t *transactionHandleGrpc) FindMonthlyPaymentMethods(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionMonthMethod, error) {
	methods, err := t.transactionService.FindMonthlyPaymentMethods(int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly payment methods: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransactionMonthMethod("success", "Successfully fetched monthly payment methods", methods)

	return so, nil
}

func (t *transactionHandleGrpc) FindYearlyPaymentMethods(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearMethod, error) {
	methods, err := t.transactionService.FindYearlyPaymentMethods(int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly payment methods: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransactionYearMethod("success", "Successfully fetched yearly payment methods", methods)

	return so, nil
}

func (t *transactionHandleGrpc) FindMonthlyAmounts(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionMonthAmount, error) {
	amounts, err := t.transactionService.FindMonthlyAmounts(int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly amounts: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransactionMonthAmount("success", "Successfully fetched monthly amounts", amounts)

	return so, nil
}

func (t *transactionHandleGrpc) FindYearlyAmounts(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearAmount, error) {
	amounts, err := t.transactionService.FindYearlyAmounts(int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly amounts: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransactionYearAmount("success", "Successfully fetched yearly amounts", amounts)

	return so, nil
}

func (t *transactionHandleGrpc) FindMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionMonthMethod, error) {
	methods, err := t.transactionService.FindMonthlyPaymentMethodsByCardNumber(req.GetCardNumber(), int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly payment methods by card number: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransactionMonthMethod("success", "Successfully fetched monthly payment methods by card number", methods)

	return so, nil
}

func (t *transactionHandleGrpc) FindYearlyPaymentMethodsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionYearMethod, error) {
	methods, err := t.transactionService.FindYearlyPaymentMethodsByCardNumber(req.GetCardNumber(), int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly payment methods by card number: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransactionYearMethod("success", "Successfully fetched yearly payment methods by card number", methods)

	return so, nil
}

func (t *transactionHandleGrpc) FindMonthlyAmountsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionMonthAmount, error) {
	amounts, err := t.transactionService.FindMonthlyAmountsByCardNumber(req.GetCardNumber(), int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly amounts by card number: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransactionMonthAmount("success", "Successfully fetched monthly amounts by card number", amounts)

	return so, nil
}

func (t *transactionHandleGrpc) FindYearlyAmountsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionYearAmount, error) {
	amounts, err := t.transactionService.FindYearlyAmountsByCardNumber(req.GetCardNumber(), int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly amounts by card number: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransactionYearAmount("success", "Successfully fetched yearly amounts by card number", amounts)

	return so, nil
}

func (t *transactionHandleGrpc) FindTransactionByMerchantIdRequest(ctx context.Context, request *pb.FindTransactionByMerchantIdRequest) (*pb.ApiResponseTransactions, error) {
	if request.GetMerchantId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid ID",
		})
	}

	merchantId := request.GetMerchantId()

	transactions, err := t.transactionService.FindTransactionByMerchantId(int(merchantId))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transactions: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransactions("success", "Successfully fetch transactions", transactions)

	return so, nil
}

func (t *transactionHandleGrpc) FindByActiveTransaction(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := t.transactionService.FindByActive(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transactions: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
	so := t.mapping.ToProtoResponsePaginationTransactionDeleteAt(paginationMeta, "success", "Successfully fetch transactions", transactions)

	return so, nil
}

func (t *transactionHandleGrpc) FindByTrashedTransaction(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := t.transactionService.FindByTrashed(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transactions: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
	so := t.mapping.ToProtoResponsePaginationTransactionDeleteAt(paginationMeta, "success", "Successfully fetch transactions", transactions)

	return so, nil
}

func (t *transactionHandleGrpc) CreateTransaction(ctx context.Context, request *pb.CreateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	transactionTime := request.GetTransactionTime().AsTime()
	merchantID := int(request.GetMerchantId())

	req := requests.CreateTransactionRequest{
		CardNumber:      request.GetCardNumber(),
		Amount:          int(request.GetAmount()),
		PaymentMethod:   request.GetPaymentMethod(),
		MerchantID:      &merchantID,
		TransactionTime: transactionTime,
	}

	res, err := t.transactionService.Create(request.ApiKey, &req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to create transaction: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransaction("success", "Successfully created transaction", res)

	return so, nil
}

func (t *transactionHandleGrpc) UpdateTransaction(ctx context.Context, request *pb.UpdateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	if request.GetTransactionId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid ID",
		})
	}

	transactionTime := request.GetTransactionTime().AsTime()
	merchantID := int(request.GetMerchantId())

	req := requests.UpdateTransactionRequest{
		TransactionID:   int(request.GetTransactionId()),
		CardNumber:      request.GetCardNumber(),
		Amount:          int(request.GetAmount()),
		PaymentMethod:   request.GetPaymentMethod(),
		MerchantID:      &merchantID,
		TransactionTime: transactionTime,
	}

	res, err := t.transactionService.Update(request.ApiKey, &req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to update transaction: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransaction("success", "Successfully updated transaction", res)

	return so, nil
}

func (t *transactionHandleGrpc) TrashedTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	if request.GetTransactionId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid ID",
		})
	}

	res, err := t.transactionService.TrashedTransaction(int(request.GetTransactionId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transaction: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransaction("success", "Successfully trashed transaction", res)

	return so, nil
}

func (t *transactionHandleGrpc) RestoreTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	if request.GetTransactionId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid ID",
		})
	}

	res, err := t.transactionService.RestoreTransaction(int(request.GetTransactionId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transaction: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransaction("success", "Successfully restored transaction", res)

	return so, nil
}

func (t *transactionHandleGrpc) DeleteTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDelete, error) {
	if request.GetTransactionId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid ID",
		})
	}

	_, err := t.transactionService.DeleteTransactionPermanent(int(request.GetTransactionId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transaction: " + err.Message,
		})
	}

	so := t.mapping.ToProtoResponseTransactionDelete("success", "Successfully deleted transaction")

	return so, nil

}

func (t *transactionHandleGrpc) RestoreAllTransaction(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	_, err := t.transactionService.RestoreAllTransaction()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all transaction: ",
		})
	}

	so := t.mapping.ToProtoResponseTransactionAll("success", "Successfully restore all transaction")

	return so, nil
}

func (t *transactionHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	_, err := t.transactionService.DeleteAllTransactionPermanent()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete transaction permanent: ",
		})
	}

	so := t.mapping.ToProtoResponseTransactionAll("success", "Successfully delete transaction permanent")

	return so, nil
}
