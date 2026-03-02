package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/transfer_errors"
	"context"
	"math"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type transferHandleGrpc struct {
	pb.UnimplementedTransferServiceServer
	transferService service.TransferService
}

func NewTransferHandleGrpc(transferService service.TransferService) *transferHandleGrpc {
	return &transferHandleGrpc{
		transferService: transferService,
	}
}

func (s *transferHandleGrpc) FindAllTransfer(ctx context.Context, request *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransfer, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTranfers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transfers, totalRecords, err := s.transferService.FindAll(ctx, &reqService)

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

	transferResponses := make([]*pb.TransferResponse, len(transfers))
	for i, transfer := range transfers {
		transferResponses[i] = &pb.TransferResponse{
			Id:             int32(transfer.TransferID),
			TransferNo:     transfer.TransferNo.String(),
			TransferFrom:   transfer.TransferFrom,
			TransferTo:     transfer.TransferTo,
			TransferAmount: int32(transfer.TransferAmount),
			TransferTime:   transfer.TransferTime.Format(time.RFC3339),
			CreatedAt:      transfer.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      transfer.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return &pb.ApiResponsePaginationTransfer{
		Status:     "success",
		Message:    "Successfully fetch transfer records",
		Data:       transferResponses,
		Pagination: paginationMeta,
	}, nil
}

func (s *transferHandleGrpc) FindByIdTransfer(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error) {
	id := int(request.GetTransferId())

	if id == 0 {
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	transfer, err := s.transferService.FindById(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Successfully fetch transfer record",
		Data: &pb.TransferResponse{
			Id:             int32(transfer.TransferID),
			TransferNo:     transfer.TransferNo.String(),
			TransferFrom:   transfer.TransferFrom,
			TransferTo:     transfer.TransferTo,
			TransferAmount: int32(transfer.TransferAmount),
			TransferTime:   transfer.TransferTime.Format(time.RFC3339),
			CreatedAt:      transfer.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      transfer.UpdatedAt.Time.Format(time.RFC3339),
		},
	}, nil
}

func (s *transferHandleGrpc) FindByTransferByTransferFrom(ctx context.Context, request *pb.FindTransferByTransferFromRequest) (*pb.ApiResponseTransfers, error) {
	transfer_from := request.GetTransferFrom()

	if transfer_from == "" {
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	transfers, err := s.transferService.FindTransferByTransferFrom(ctx, transfer_from)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	transferResponses := make([]*pb.TransferResponse, len(transfers))
	for i, transfer := range transfers {
		transferResponses[i] = &pb.TransferResponse{
			Id:             int32(transfer.TransferID),
			TransferNo:     transfer.TransferNo.String(),
			TransferFrom:   transfer.TransferFrom,
			TransferTo:     transfer.TransferTo,
			TransferAmount: int32(transfer.TransferAmount),
			TransferTime:   transfer.TransferTime.Format(time.RFC3339),
			CreatedAt:      transfer.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      transfer.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return &pb.ApiResponseTransfers{
		Status:  "success",
		Message: "Successfully fetch transfer records",
		Data:    transferResponses,
	}, nil
}

func (s *transferHandleGrpc) FindByTransferByTransferTo(ctx context.Context, request *pb.FindTransferByTransferToRequest) (*pb.ApiResponseTransfers, error) {
	transfer_to := request.GetTransferTo()

	if transfer_to == "" {
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	transfers, err := s.transferService.FindTransferByTransferTo(ctx, transfer_to)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	transferResponses := make([]*pb.TransferResponse, len(transfers))
	for i, transfer := range transfers {
		transferResponses[i] = &pb.TransferResponse{
			Id:             int32(transfer.TransferID),
			TransferNo:     transfer.TransferNo.String(),
			TransferFrom:   transfer.TransferFrom,
			TransferTo:     transfer.TransferTo,
			TransferAmount: int32(transfer.TransferAmount),
			TransferTime:   transfer.TransferTime.Format(time.RFC3339),
			CreatedAt:      transfer.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      transfer.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return &pb.ApiResponseTransfers{
		Status:  "success",
		Message: "Successfully fetch transfer records",
		Data:    transferResponses,
	}, nil
}

func (s *transferHandleGrpc) FindByActiveTransfer(ctx context.Context, req *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransferDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTranfers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transfers, totalRecords, err := s.transferService.FindByActive(ctx, &reqService)

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

	transferResponses := make([]*pb.TransferResponseDeleteAt, len(transfers))
	for i, transfer := range transfers {
		transferResponses[i] = &pb.TransferResponseDeleteAt{
			Id:             int32(transfer.TransferID),
			TransferNo:     transfer.TransferNo.String(),
			TransferFrom:   transfer.TransferFrom,
			TransferTo:     transfer.TransferTo,
			TransferAmount: int32(transfer.TransferAmount),
			TransferTime:   transfer.TransferTime.Format(time.RFC3339),
			CreatedAt:      transfer.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      transfer.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:      &wrapperspb.StringValue{Value: transfer.DeletedAt.Time.Format(time.RFC3339)},
		}
	}

	return &pb.ApiResponsePaginationTransferDeleteAt{
		Status:     "success",
		Message:    "Successfully fetch transfer records",
		Data:       transferResponses,
		Pagination: paginationMeta,
	}, nil
}

func (s *transferHandleGrpc) FindByTrashedTransfer(ctx context.Context, req *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransferDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTranfers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transfers, totalRecords, err := s.transferService.FindByTrashed(ctx, &reqService)

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

	transferResponses := make([]*pb.TransferResponseDeleteAt, len(transfers))
	for i, transfer := range transfers {
		transferResponses[i] = &pb.TransferResponseDeleteAt{
			Id:             int32(transfer.TransferID),
			TransferNo:     transfer.TransferNo.String(),
			TransferFrom:   transfer.TransferFrom,
			TransferTo:     transfer.TransferTo,
			TransferAmount: int32(transfer.TransferAmount),
			TransferTime:   transfer.TransferTime.Format(time.RFC3339),
			CreatedAt:      transfer.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      transfer.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:      &wrapperspb.StringValue{Value: transfer.DeletedAt.Time.Format(time.RFC3339)},
		}
	}

	return &pb.ApiResponsePaginationTransferDeleteAt{
		Status:     "success",
		Message:    "Successfully fetch transfer records",
		Data:       transferResponses,
		Pagination: paginationMeta,
	}, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferStatusSuccess(ctx context.Context, req *pb.FindMonthlyTransferStatus) (*pb.ApiResponseTransferMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthStatusTransfer{
		Year:  year,
		Month: month,
	}

	records, err := s.transferService.FindMonthTransferStatusSuccess(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferMonthStatusSuccessResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransferMonthStatusSuccessResponse{
			Year:         record.Year,
			Month:        record.Month,
			TotalSuccess: int32(record.TotalSuccess),
			TotalAmount:  int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransferMonthStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched monthly Transfer status success",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindYearlyTransferStatusSuccess(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearStatusSuccess, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	records, err := s.transferService.FindYearlyTransferStatusSuccess(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferYearStatusSuccessResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransferYearStatusSuccessResponse{
			Year:         record.Year,
			TotalSuccess: int32(record.TotalSuccess),
			TotalAmount:  int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransferYearStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly Transfer status success",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferStatusFailed(ctx context.Context, req *pb.FindMonthlyTransferStatus) (*pb.ApiResponseTransferMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthStatusTransfer{
		Year:  year,
		Month: month,
	}

	records, err := s.transferService.FindMonthTransferStatusFailed(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferMonthStatusFailedResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransferMonthStatusFailedResponse{
			Year:        record.Year,
			Month:       record.Month,
			TotalFailed: int32(record.TotalFailed),
			TotalAmount: int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransferMonthStatusFailed{
		Status:  "success",
		Message: "success fetched monthly Transfer status Failed",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindYearlyTransferStatusFailed(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearStatusFailed, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	records, err := s.transferService.FindYearlyTransferStatusFailed(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferYearStatusFailedResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransferYearStatusFailedResponse{
			Year:        record.Year,
			TotalFailed: int32(record.TotalFailed),
			TotalAmount: int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransferYearStatusFailed{
		Status:  "success",
		Message: "success fetched yearly Transfer status Failed",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTransferStatusCardNumber) (*pb.ApiResponseTransferMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthStatusTransferCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.transferService.FindMonthTransferStatusSuccessByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferMonthStatusSuccessResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransferMonthStatusSuccessResponse{
			Year:         record.Year,
			Month:        record.Month,
			TotalSuccess: int32(record.TotalSuccess),
			TotalAmount:  int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransferMonthStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched monthly Transfer status success",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTransferStatusCardNumber) (*pb.ApiResponseTransferYearStatusSuccess, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearStatusTransferCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.transferService.FindYearlyTransferStatusSuccessByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferYearStatusSuccessResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransferYearStatusSuccessResponse{
			Year:         record.Year,
			TotalSuccess: int32(record.TotalSuccess),
			TotalAmount:  int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransferYearStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly Transfer status success",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTransferStatusCardNumber) (*pb.ApiResponseTransferMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthStatusTransferCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.transferService.FindMonthTransferStatusFailedByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferMonthStatusFailedResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransferMonthStatusFailedResponse{
			Year:        record.Year,
			Month:       record.Month,
			TotalFailed: int32(record.TotalFailed),
			TotalAmount: int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransferMonthStatusFailed{
		Status:  "success",
		Message: "success fetched monthly Transfer status Failed",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTransferStatusCardNumber) (*pb.ApiResponseTransferYearStatusFailed, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearStatusTransferCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.transferService.FindYearlyTransferStatusFailedByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferYearStatusFailedResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.TransferYearStatusFailedResponse{
			Year:        record.Year,
			TotalFailed: int32(record.TotalFailed),
			TotalAmount: int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseTransferYearStatusFailed{
		Status:  "success",
		Message: "success fetched yearly Transfer status Failed",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferAmounts(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferMonthAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	amounts, err := s.transferService.FindMonthlyTransferAmounts(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferMonthAmountResponse, len(amounts))
	for i, amount := range amounts {
		dataResponses[i] = &pb.TransferMonthAmountResponse{
			Month:       amount.Month,
			TotalAmount: int32(amount.TotalTransferAmount),
		}
	}

	return &pb.ApiResponseTransferMonthAmount{
		Status:  "success",
		Message: "Successfully fetched monthly transfer amounts",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindYearlyTransferAmounts(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	amounts, err := s.transferService.FindYearlyTransferAmounts(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferYearAmountResponse, len(amounts))
	for i, amount := range amounts {
		dataResponses[i] = &pb.TransferYearAmountResponse{
			Year:        amount.Year.Int.String(),
			TotalAmount: int32(amount.TotalTransferAmount),
		}
	}

	return &pb.ApiResponseTransferYearAmount{
		Status:  "success",
		Message: "Successfully fetched yearly transfer amounts",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.transferService.FindMonthlyTransferAmountsBySenderCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferMonthAmountResponse, len(amounts))
	for i, amount := range amounts {
		dataResponses[i] = &pb.TransferMonthAmountResponse{
			Month:       amount.Month,
			TotalAmount: int32(amount.TotalTransferAmount),
		}
	}

	return &pb.ApiResponseTransferMonthAmount{
		Status:  "success",
		Message: "Successfully fetched monthly transfer amounts by sender card number",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.transferService.FindMonthlyTransferAmountsByReceiverCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferMonthAmountResponse, len(amounts))
	for i, amount := range amounts {
		dataResponses[i] = &pb.TransferMonthAmountResponse{
			Month:       amount.Month,
			TotalAmount: int32(amount.TotalTransferAmount),
		}
	}

	return &pb.ApiResponseTransferMonthAmount{
		Status:  "success",
		Message: "Successfully fetched monthly transfer amounts by receiver card number",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.transferService.FindYearlyTransferAmountsBySenderCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferYearAmountResponse, len(amounts))
	for i, amount := range amounts {
		dataResponses[i] = &pb.TransferYearAmountResponse{
			Year:        amount.Year.Int.String(),
			TotalAmount: int32(amount.TotalTransferAmount),
		}
	}

	return &pb.ApiResponseTransferYearAmount{
		Status:  "success",
		Message: "Successfully fetched yearly transfer amounts by sender card number",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) FindYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.transferService.FindYearlyTransferAmountsByReceiverCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.TransferYearAmountResponse, len(amounts))
	for i, amount := range amounts {
		dataResponses[i] = &pb.TransferYearAmountResponse{
			Year:        amount.Year.Int.String(),
			TotalAmount: int32(amount.TotalTransferAmount),
		}
	}

	return &pb.ApiResponseTransferYearAmount{
		Status:  "success",
		Message: "Successfully fetched yearly transfer amounts by receiver card number",
		Data:    dataResponses,
	}, nil
}

func (s *transferHandleGrpc) CreateTransfer(ctx context.Context, request *pb.CreateTransferRequest) (*pb.ApiResponseTransfer, error) {
	req := requests.CreateTransferRequest{
		TransferFrom:   request.GetTransferFrom(),
		TransferTo:     request.GetTransferTo(),
		TransferAmount: int(request.GetTransferAmount()),
	}

	if err := req.Validate(); err != nil {
		return nil, transfer_errors.ErrGrpcValidateCreateTransferRequest
	}

	transfer, err := s.transferService.CreateTransaction(ctx, &req)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Successfully created transfer",
		Data: &pb.TransferResponse{
			Id:             int32(transfer.TransferID),
			TransferNo:     transfer.TransferNo.String(),
			TransferFrom:   transfer.TransferFrom,
			TransferTo:     transfer.TransferTo,
			TransferAmount: int32(transfer.TransferAmount),
			TransferTime:   transfer.TransferTime.Format(time.RFC3339),
			CreatedAt:      transfer.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      transfer.UpdatedAt.Time.Format(time.RFC3339),
		},
	}, nil
}

func (s *transferHandleGrpc) UpdateTransfer(ctx context.Context, request *pb.UpdateTransferRequest) (*pb.ApiResponseTransfer, error) {
	id := int(request.GetTransferId())

	if id == 0 {
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	req := requests.UpdateTransferRequest{
		TransferID:     &id,
		TransferFrom:   request.GetTransferFrom(),
		TransferTo:     request.GetTransferTo(),
		TransferAmount: int(request.GetTransferAmount()),
	}

	if err := req.Validate(); err != nil {
		return nil, transfer_errors.ErrGrpcValidateUpdateTransferRequest
	}

	transfer, err := s.transferService.UpdateTransaction(ctx, &req)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Successfully updated transfer",
		Data: &pb.TransferResponse{
			Id:             int32(transfer.TransferID),
			TransferNo:     transfer.TransferNo.String(),
			TransferFrom:   transfer.TransferFrom,
			TransferTo:     transfer.TransferTo,
			TransferAmount: int32(transfer.TransferAmount),
			TransferTime:   transfer.TransferTime.Format(time.RFC3339),
			CreatedAt:      transfer.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      transfer.UpdatedAt.Time.Format(time.RFC3339),
		},
	}, nil
}

func (s *transferHandleGrpc) TrashedTransfer(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransferDeleteAt, error) {
	id := int(request.GetTransferId())

	if id == 0 {
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	transfer, err := s.transferService.TrashedTransfer(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransferDeleteAt{
		Status:  "success",
		Message: "Successfully trashed transfer",
		Data: &pb.TransferResponseDeleteAt{
			Id:             int32(transfer.TransferID),
			TransferNo:     transfer.TransferNo.String(),
			TransferFrom:   transfer.TransferFrom,
			TransferTo:     transfer.TransferTo,
			TransferAmount: int32(transfer.TransferAmount),
			TransferTime:   transfer.TransferTime.Format(time.RFC3339),
			CreatedAt:      transfer.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      transfer.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:      &wrapperspb.StringValue{Value: transfer.DeletedAt.Time.Format(time.RFC3339)},
		},
	}, nil
}

func (s *transferHandleGrpc) RestoreTransfer(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransferDeleteAt, error) {
	id := int(request.GetTransferId())

	if id == 0 {
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	transfer, err := s.transferService.RestoreTransfer(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransferDeleteAt{
		Status:  "success",
		Message: "Successfully restored transfer",
		Data: &pb.TransferResponseDeleteAt{
			Id:             int32(transfer.TransferID),
			TransferNo:     transfer.TransferNo.String(),
			TransferFrom:   transfer.TransferFrom,
			TransferTo:     transfer.TransferTo,
			TransferAmount: int32(transfer.TransferAmount),
			TransferTime:   transfer.TransferTime.Format(time.RFC3339),
			CreatedAt:      transfer.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      transfer.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:      &wrapperspb.StringValue{Value: transfer.DeletedAt.Time.Format(time.RFC3339)},
		},
	}, nil
}

func (s *transferHandleGrpc) DeleteTransferPermanent(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransferDelete, error) {
	id := int(request.GetTransferId())

	if id == 0 {
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	_, err := s.transferService.DeleteTransferPermanent(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransferDelete{
		Status:  "success",
		Message: "Successfully deleted transfer permanently",
	}, nil
}

func (s *transferHandleGrpc) RestoreAllTransfer(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransferAll, error) {
	_, err := s.transferService.RestoreAllTransfer(ctx)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransferAll{
		Status:  "success",
		Message: "Successfully restored all transfers",
	}, nil
}

func (s *transferHandleGrpc) DeleteAllTransferPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransferAll, error) {
	_, err := s.transferService.DeleteAllTransferPermanent(ctx)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransferAll{
		Status:  "success",
		Message: "Successfully deleted all transfers permanently",
	}, nil
}
