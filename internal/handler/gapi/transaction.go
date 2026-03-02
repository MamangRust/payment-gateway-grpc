package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/transaction_errors"
	"context"
	"math"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type transactionHandleGrpc struct {
	pb.UnimplementedTransactionServiceServer
	transactionService service.TransactionService
}

func NewTransactionHandleGrpc(transactionService service.TransactionService) *transactionHandleGrpc {
	return &transactionHandleGrpc{
		transactionService: transactionService,
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

	reqService := requests.FindAllTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := t.transactionService.FindAll(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	transactionResponses := make([]*pb.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionResponses[i] = &pb.TransactionResponse{
			Id:              int32(transaction.TransactionID),
			CardNumber:      transaction.CardNumber,
			TransactionNo:   transaction.TransactionNo.String(),
			Amount:          int32(transaction.Amount),
			PaymentMethod:   transaction.PaymentMethod,
			MerchantId:      int32(transaction.MerchantID),
			TransactionTime: transaction.TransactionTime.Format(time.RFC3339),
			CreatedAt:       transaction.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:       transaction.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return &pb.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transaction records",
		Data:       transactionResponses,
		Pagination: paginationMeta,
	}, nil
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
	reqService := requests.FindAllTransactionCardNumber{
		CardNumber: card_number,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	transactions, totalRecords, err := t.transactionService.FindAllByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	transactionResponses := make([]*pb.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionResponses[i] = &pb.TransactionResponse{
			Id:              int32(transaction.TransactionID),
			CardNumber:      transaction.CardNumber,
			TransactionNo:   transaction.TransactionNo.String(),
			Amount:          int32(transaction.Amount),
			PaymentMethod:   transaction.PaymentMethod,
			MerchantId:      int32(transaction.MerchantID),
			TransactionTime: transaction.TransactionTime.Format(time.RFC3339),
			CreatedAt:       transaction.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:       transaction.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return &pb.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transaction records",
		Data:       transactionResponses,
		Pagination: paginationMeta,
	}, nil
}

func (t *transactionHandleGrpc) FindByIdTransaction(ctx context.Context, req *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(req.GetTransactionId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	transaction, err := t.transactionService.FindById(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Transaction fetched successfully",
		Data: &pb.TransactionResponse{
			Id:              int32(transaction.TransactionID),
			CardNumber:      transaction.CardNumber,
			TransactionNo:   transaction.TransactionNo.String(),
			Amount:          int32(transaction.Amount),
			PaymentMethod:   transaction.PaymentMethod,
			MerchantId:      int32(transaction.MerchantID),
			TransactionTime: transaction.TransactionTime.Format(time.RFC3339),
			CreatedAt:       transaction.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:       transaction.UpdatedAt.Time.Format(time.RFC3339),
		},
	}, nil
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

	reqService := requests.FindAllTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := t.transactionService.FindByActive(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	transactionResponses := make([]*pb.TransactionResponseDeleteAt, len(transactions))
	for i, transaction := range transactions {
		transactionResponses[i] = &pb.TransactionResponseDeleteAt{
			Id:              int32(transaction.TransactionID),
			CardNumber:      transaction.CardNumber,
			TransactionNo:   transaction.TransactionNo.String(),
			Amount:          int32(transaction.Amount),
			PaymentMethod:   transaction.PaymentMethod,
			MerchantId:      int32(transaction.MerchantID),
			TransactionTime: transaction.TransactionTime.Format(time.RFC3339),
			CreatedAt:       transaction.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:       transaction.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:       &wrapperspb.StringValue{Value: transaction.DeletedAt.Time.Format(time.RFC3339)},
		}
	}

	return &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetch transactions",
		Data:       transactionResponses,
		Pagination: paginationMeta,
	}, nil
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

	reqService := requests.FindAllTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := t.transactionService.FindByTrashed(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	transactionResponses := make([]*pb.TransactionResponseDeleteAt, len(transactions))
	for i, transaction := range transactions {
		transactionResponses[i] = &pb.TransactionResponseDeleteAt{
			Id:              int32(transaction.TransactionID),
			CardNumber:      transaction.CardNumber,
			TransactionNo:   transaction.TransactionNo.String(),
			Amount:          int32(transaction.Amount),
			PaymentMethod:   transaction.PaymentMethod,
			MerchantId:      int32(transaction.MerchantID),
			TransactionTime: transaction.TransactionTime.Format(time.RFC3339),
			CreatedAt:       transaction.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:       transaction.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:       &wrapperspb.StringValue{Value: transaction.DeletedAt.Time.Format(time.RFC3339)},
		}
	}

	return &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetch transactions",
		Data:       transactionResponses,
		Pagination: paginationMeta,
	}, nil
}

func (s *transactionHandleGrpc) FindMonthlyTransactionStatusSuccess(ctx context.Context, req *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthStatusTransaction{
		Year:  year,
		Month: month,
	}

	records, err := s.transactionService.FindMonthTransactionStatusSuccess(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionMonthStatusSuccessResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransactionMonthStatusSuccessResponse{
			Year:         record.Year,
			Month:        record.Month,
			TotalSuccess: int32(record.TotalSuccess),
			TotalAmount:  int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched monthly Transaction status success",
		Data:    dataResponses,
	}, nil
}

func (s *transactionHandleGrpc) FindYearlyTransactionStatusSuccess(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearStatusSuccess, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	records, err := s.transactionService.FindYearlyTransactionStatusSuccess(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionYearStatusSuccessResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransactionYearStatusSuccessResponse{
			Year:         record.Year,
			TotalSuccess: int32(record.TotalSuccess),
			TotalAmount:  int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly Transaction status success",
		Data:    dataResponses,
	}, nil
}

func (s *transactionHandleGrpc) FindMonthlyTransactionStatusFailed(ctx context.Context, req *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthStatusTransaction{
		Year:  year,
		Month: month,
	}

	records, err := s.transactionService.FindMonthTransactionStatusFailed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionMonthStatusFailedResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransactionMonthStatusFailedResponse{
			Year:        record.Year,
			Month:       record.Month,
			TotalFailed: int32(record.TotalFailed),
			TotalAmount: int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthStatusFailed{
		Status:  "success",
		Message: "success fetched monthly Transaction status Failed",
		Data:    dataResponses,
	}, nil
}

func (s *transactionHandleGrpc) FindYearlyTransactionStatusFailed(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearStatusFailed, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	records, err := s.transactionService.FindYearlyTransactionStatusFailed(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionYearStatusFailedResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransactionYearStatusFailedResponse{
			Year:        record.Year,
			TotalFailed: int32(record.TotalFailed),
			TotalAmount: int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearStatusFailed{
		Status:  "success",
		Message: "success fetched yearly Transaction status Failed",
		Data:    dataResponses,
	}, nil
}

func (s *transactionHandleGrpc) FindMonthlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTransactionStatusCardNumber) (*pb.ApiResponseTransactionMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthStatusTransactionCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	records, err := s.transactionService.FindMonthTransactionStatusSuccessByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionMonthStatusSuccessResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransactionMonthStatusSuccessResponse{
			Year:         record.Year,
			Month:        record.Month,
			TotalSuccess: int32(record.TotalSuccess),
			TotalAmount:  int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched monthly Transaction status success",
		Data:    dataResponses,
	}, nil
}

func (s *transactionHandleGrpc) FindYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTransactionStatusCardNumber) (*pb.ApiResponseTransactionYearStatusSuccess, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearStatusTransactionCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.transactionService.FindYearlyTransactionStatusSuccessByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionYearStatusSuccessResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransactionYearStatusSuccessResponse{
			Year:         record.Year,
			TotalSuccess: int32(record.TotalSuccess),
			TotalAmount:  int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly Transaction status success",
		Data:    dataResponses,
	}, nil
}

func (s *transactionHandleGrpc) FindMonthlyTransactionStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTransactionStatusCardNumber) (*pb.ApiResponseTransactionMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthStatusTransactionCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	records, err := s.transactionService.FindMonthTransactionStatusFailedByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionMonthStatusFailedResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransactionMonthStatusFailedResponse{
			Year:        record.Year,
			Month:       record.Month,
			TotalFailed: int32(record.TotalFailed),
			TotalAmount: int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthStatusFailed{
		Status:  "success",
		Message: "success fetched monthly Transaction status Failed",
		Data:    dataResponses,
	}, nil
}

func (s *transactionHandleGrpc) FindYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTransactionStatusCardNumber) (*pb.ApiResponseTransactionYearStatusFailed, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearStatusTransactionCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.transactionService.FindYearlyTransactionStatusFailedByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionYearStatusFailedResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransactionYearStatusFailedResponse{
			Year:        record.Year,
			TotalFailed: int32(record.TotalFailed),
			TotalAmount: int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearStatusFailed{
		Status:  "success",
		Message: "success fetched yearly Transaction status Failed",
		Data:    dataResponses,
	}, nil
}

func (t *transactionHandleGrpc) FindMonthlyPaymentMethods(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionMonthMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := t.transactionService.FindMonthlyPaymentMethods(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionMonthMethodResponse, len(methods))
	for i, method := range methods {
		dataResponses[i] = &pb.TransactionMonthMethodResponse{
			Month:             method.Month,
			PaymentMethod:     method.PaymentMethod,
			TotalTransactions: int32(method.TotalTransactions),
			TotalAmount:       int32(method.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthMethod{
		Status:  "success",
		Message: "Successfully fetched monthly payment methods",
		Data:    dataResponses,
	}, nil
}

func (t *transactionHandleGrpc) FindYearlyPaymentMethods(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := t.transactionService.FindYearlyPaymentMethods(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionYearMethodResponse, len(methods))
	for i, method := range methods {
		dataResponses[i] = &pb.TransactionYearMethodResponse{
			Year:              method.Year.Int.String(),
			PaymentMethod:     method.PaymentMethod,
			TotalTransactions: int32(method.TotalTransactions),
			TotalAmount:       int32(method.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearMethod{
		Status:  "success",
		Message: "Successfully fetched yearly payment methods",
		Data:    dataResponses,
	}, nil
}

func (t *transactionHandleGrpc) FindMonthlyAmounts(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionMonthAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	amounts, err := t.transactionService.FindMonthlyAmounts(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionMonthAmountResponse, len(amounts))
	for i, amount := range amounts {
		dataResponses[i] = &pb.TransactionMonthAmountResponse{
			Month:       amount.Month,
			TotalAmount: int32(amount.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthAmount{
		Status:  "success",
		Message: "Successfully fetched monthly amounts",
		Data:    dataResponses,
	}, nil
}

func (t *transactionHandleGrpc) FindYearlyAmounts(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	amounts, err := t.transactionService.FindYearlyAmounts(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionYearlyAmountResponse, len(amounts))
	for i, amount := range amounts {
		dataResponses[i] = &pb.TransactionYearlyAmountResponse{
			Year:        amount.Year.Int.String(),
			TotalAmount: int32(amount.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearAmount{
		Status:  "success",
		Message: "Successfully fetched yearly amounts",
		Data:    dataResponses,
	}, nil
}

func (t *transactionHandleGrpc) FindMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionMonthMethod, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearPaymentMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	methods, err := t.transactionService.FindMonthlyPaymentMethodsByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionMonthMethodResponse, len(methods))
	for i, method := range methods {
		dataResponses[i] = &pb.TransactionMonthMethodResponse{
			Month:             method.Month,
			PaymentMethod:     method.PaymentMethod,
			TotalTransactions: int32(method.TotalTransactions),
			TotalAmount:       int32(method.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthMethod{
		Status:  "success",
		Message: "Successfully fetched monthly payment methods by card number",
		Data:    dataResponses,
	}, nil
}

func (t *transactionHandleGrpc) FindYearlyPaymentMethodsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionYearMethod, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearPaymentMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	methods, err := t.transactionService.FindYearlyPaymentMethodsByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionYearMethodResponse, len(methods))
	for i, method := range methods {
		dataResponses[i] = &pb.TransactionYearMethodResponse{
			Year:              method.Year.Int.String(),
			PaymentMethod:     method.PaymentMethod,
			TotalTransactions: int32(method.TotalTransactions),
			TotalAmount:       int32(method.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearMethod{
		Status:  "success",
		Message: "Successfully fetched yearly payment methods by card number",
		Data:    dataResponses,
	}, nil
}

func (t *transactionHandleGrpc) FindMonthlyAmountsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearPaymentMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := t.transactionService.FindMonthlyAmountsByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionMonthAmountResponse, len(amounts))
	for i, amount := range amounts {
		dataResponses[i] = &pb.TransactionMonthAmountResponse{
			Month:       amount.Month,
			TotalAmount: int32(amount.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthAmount{
		Status:  "success",
		Message: "Successfully fetched monthly amounts by card number",
		Data:    dataResponses,
	}, nil
}

func (t *transactionHandleGrpc) FindYearlyAmountsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearPaymentMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := t.transactionService.FindYearlyAmountsByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionYearlyAmountResponse, len(amounts))
	for i, amount := range amounts {
		dataResponses[i] = &pb.TransactionYearlyAmountResponse{
			Year:        amount.Year.Int.String(),
			TotalAmount: int32(amount.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearAmount{
		Status:  "success",
		Message: "Successfully fetched yearly amounts by card number",
		Data:    dataResponses,
	}, nil
}

func (t *transactionHandleGrpc) FindTransactionByMerchantIdRequest(ctx context.Context, req *pb.FindTransactionByMerchantIdRequest) (*pb.ApiResponseTransactions, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcTransactionInvalidMerchantID
	}

	transactions, err := t.transactionService.FindTransactionByMerchantId(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		dataResponses[i] = &pb.TransactionResponse{
			Id:              int32(transaction.TransactionID),
			CardNumber:      transaction.CardNumber,
			TransactionNo:   transaction.TransactionNo.String(),
			Amount:          int32(transaction.Amount),
			PaymentMethod:   transaction.PaymentMethod,
			MerchantId:      int32(transaction.MerchantID),
			TransactionTime: transaction.TransactionTime.Format(time.RFC3339),
			CreatedAt:       transaction.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:       transaction.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return &pb.ApiResponseTransactions{
		Status:  "success",
		Message: "Successfully fetch transactions",
		Data:    dataResponses,
	}, nil
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

	if err := req.Validate(); err != nil {
		return nil, transaction_errors.ErrGrpcValidateCreateTransactionRequest
	}

	res, err := t.transactionService.Create(ctx, request.ApiKey, &req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully created transaction",
		Data: &pb.TransactionResponse{
			Id:              int32(res.TransactionID),
			CardNumber:      res.CardNumber,
			TransactionNo:   res.TransactionNo.String(),
			Amount:          int32(res.Amount),
			PaymentMethod:   res.PaymentMethod,
			MerchantId:      int32(res.MerchantID),
			TransactionTime: res.TransactionTime.Format(time.RFC3339),
			CreatedAt:       res.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:       res.UpdatedAt.Time.Format(time.RFC3339),
		},
	}, nil
}

func (t *transactionHandleGrpc) UpdateTransaction(ctx context.Context, request *pb.UpdateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(request.GetTransactionId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	transactionTime := request.GetTransactionTime().AsTime()
	merchantID := int(request.GetMerchantId())

	req := requests.UpdateTransactionRequest{
		TransactionID:   &id,
		CardNumber:      request.GetCardNumber(),
		Amount:          int(request.GetAmount()),
		PaymentMethod:   request.GetPaymentMethod(),
		MerchantID:      &merchantID,
		TransactionTime: transactionTime,
	}

	if err := req.Validate(); err != nil {
		return nil, transaction_errors.ErrGrpcValidateUpdateTransactionRequest
	}

	res, err := t.transactionService.Update(ctx, request.ApiKey, &req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully updated transaction",
		Data: &pb.TransactionResponse{
			Id:              int32(res.TransactionID),
			CardNumber:      res.CardNumber,
			TransactionNo:   res.TransactionNo.String(),
			Amount:          int32(res.Amount),
			PaymentMethod:   res.PaymentMethod,
			MerchantId:      int32(res.MerchantID),
			TransactionTime: res.TransactionTime.Format(time.RFC3339),
			CreatedAt:       res.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:       res.UpdatedAt.Time.Format(time.RFC3339),
		},
	}, nil
}

func (t *transactionHandleGrpc) TrashedTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDeleteAt, error) {
	id := int(request.GetTransactionId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	res, err := t.transactionService.TrashedTransaction(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully trashed transaction",
		Data: &pb.TransactionResponseDeleteAt{
			Id:              int32(res.TransactionID),
			CardNumber:      res.CardNumber,
			TransactionNo:   res.TransactionNo.String(),
			Amount:          int32(res.Amount),
			PaymentMethod:   res.PaymentMethod,
			MerchantId:      int32(res.MerchantID),
			TransactionTime: res.TransactionTime.Format(time.RFC3339),
			CreatedAt:       res.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:       res.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:       &wrapperspb.StringValue{Value: res.DeletedAt.Time.Format(time.RFC3339)},
		},
	}, nil
}

func (t *transactionHandleGrpc) RestoreTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDeleteAt, error) {
	id := int(request.GetTransactionId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	res, err := t.transactionService.RestoreTransaction(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully restored transaction",
		Data: &pb.TransactionResponseDeleteAt{
			Id:              int32(res.TransactionID),
			CardNumber:      res.CardNumber,
			TransactionNo:   res.TransactionNo.String(),
			Amount:          int32(res.Amount),
			PaymentMethod:   res.PaymentMethod,
			MerchantId:      int32(res.MerchantID),
			TransactionTime: res.TransactionTime.Format(time.RFC3339),
			CreatedAt:       res.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:       res.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:       &wrapperspb.StringValue{Value: res.DeletedAt.Time.Format(time.RFC3339)},
		},
	}, nil
}

func (t *transactionHandleGrpc) DeleteTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDelete, error) {
	id := int(request.GetTransactionId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	_, err := t.transactionService.DeleteTransactionPermanent(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted transaction",
	}, nil
}

func (t *transactionHandleGrpc) RestoreAllTransaction(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	_, err := t.transactionService.RestoreAllTransaction(ctx)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully restore all transaction",
	}, nil
}

func (t *transactionHandleGrpc) DeleteAllTransactionPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	_, err := t.transactionService.DeleteAllTransactionPermanent(ctx)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully delete transaction permanent",
	}, nil
}
