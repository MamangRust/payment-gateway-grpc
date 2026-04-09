package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/withdraw_errors"
	"context"
	"math"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type withdrawHandleGrpc struct {
	pb.UnimplementedWithdrawServiceServer
	withdrawService service.WithdrawService
}

func NewWithdrawHandleGrpc(withdraw service.WithdrawService) *withdrawHandleGrpc {
	return &withdrawHandleGrpc{
		withdrawService: withdraw,
	}
}

func (w *withdrawHandleGrpc) FindAllWithdraw(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdraw, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	withdraws, totalRecords, err := w.withdrawService.FindAll(ctx, &reqService)

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

	withdrawResponses := make([]*pb.WithdrawResponse, len(withdraws))
	for i, withdraw := range withdraws {
		withdrawResponses[i] = &pb.WithdrawResponse{
			WithdrawId:     int32(withdraw.WithdrawID),
			WithdrawNo:     withdraw.WithdrawNo.String(),
			CardNumber:     withdraw.CardNumber,
			WithdrawAmount: int32(withdraw.WithdrawAmount),
			WithdrawTime:   withdraw.WithdrawTime.Format(time.RFC3339),
			CreatedAt:      withdraw.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      withdraw.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return &pb.ApiResponsePaginationWithdraw{
		Status:     "success",
		Message:    "withdraw",
		Data:       withdrawResponses,
		Pagination: paginationMeta,
	}, nil
}

func (w *withdrawHandleGrpc) FindAllWithdrawByCardNumber(ctx context.Context, req *pb.FindAllWithdrawByCardNumberRequest) (*pb.ApiResponsePaginationWithdraw, error) {
	card_number := req.GetCardNumber()
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllWithdrawCardNumber{
		CardNumber: card_number,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	withdraws, totalRecords, err := w.withdrawService.FindAllByCardNumber(ctx, &reqService)

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

	withdrawResponses := make([]*pb.WithdrawResponse, len(withdraws))
	for i, withdraw := range withdraws {
		withdrawResponses[i] = &pb.WithdrawResponse{
			WithdrawId:     int32(withdraw.WithdrawID),
			WithdrawNo:     withdraw.WithdrawNo.String(),
			CardNumber:     withdraw.CardNumber,
			WithdrawAmount: int32(withdraw.WithdrawAmount),
			WithdrawTime:   withdraw.WithdrawTime.Format(time.RFC3339),
			CreatedAt:      withdraw.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      withdraw.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return &pb.ApiResponsePaginationWithdraw{
		Status:     "success",
		Message:    "Withdraws fetched successfully",
		Data:       withdrawResponses,
		Pagination: paginationMeta,
	}, nil
}

func (w *withdrawHandleGrpc) FindByIdWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	id := int(req.GetWithdrawId())

	if id == 0 {
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdraw, err := w.withdrawService.FindById(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseWithdraw{
		Status:  "success",
		Message: "Successfully fetched withdraw",
		Data: &pb.WithdrawResponse{
			WithdrawId:     int32(withdraw.WithdrawID),
			WithdrawNo:     withdraw.WithdrawNo.String(),
			CardNumber:     withdraw.CardNumber,
			WithdrawAmount: int32(withdraw.WithdrawAmount),
			WithdrawTime:   withdraw.WithdrawTime.Format(time.RFC3339),
			CreatedAt:      withdraw.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      withdraw.UpdatedAt.Time.Format(time.RFC3339),
		},
	}, nil
}

func (w *withdrawHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdrawDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	withdraws, totalRecords, err := w.withdrawService.FindByActive(ctx, &reqService)

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

	withdrawResponses := make([]*pb.WithdrawResponseDeleteAt, len(withdraws))
	for i, withdraw := range withdraws {
		withdrawResponses[i] = &pb.WithdrawResponseDeleteAt{
			WithdrawId:     int32(withdraw.WithdrawID),
			WithdrawNo:     withdraw.WithdrawNo.String(),
			CardNumber:     withdraw.CardNumber,
			WithdrawAmount: int32(withdraw.WithdrawAmount),
			WithdrawTime:   withdraw.WithdrawTime.Format(time.RFC3339),
			CreatedAt:      withdraw.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      withdraw.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:      &wrapperspb.StringValue{Value: withdraw.DeletedAt.Time.Format(time.RFC3339)},
		}
	}

	return &pb.ApiResponsePaginationWithdrawDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched withdraws",
		Data:       withdrawResponses,
		Pagination: paginationMeta,
	}, nil
}

func (w *withdrawHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdrawDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	withdraws, totalRecords, err := w.withdrawService.FindByTrashed(ctx, &reqService)

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

	withdrawResponses := make([]*pb.WithdrawResponseDeleteAt, len(withdraws))
	for i, withdraw := range withdraws {
		withdrawResponses[i] = &pb.WithdrawResponseDeleteAt{
			WithdrawId:     int32(withdraw.WithdrawID),
			WithdrawNo:     withdraw.WithdrawNo.String(),
			CardNumber:     withdraw.CardNumber,
			WithdrawAmount: int32(withdraw.WithdrawAmount),
			WithdrawTime:   withdraw.WithdrawTime.Format(time.RFC3339),
			CreatedAt:      withdraw.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      withdraw.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:      &wrapperspb.StringValue{Value: withdraw.DeletedAt.Time.Format(time.RFC3339)},
		}
	}

	return &pb.ApiResponsePaginationWithdrawDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched withdraws",
		Data:       withdrawResponses,
		Pagination: paginationMeta,
	}, nil
}

func (s *withdrawHandleGrpc) FindMonthlyWithdrawStatusSuccess(ctx context.Context, req *pb.FindMonthlyWithdrawStatus) (*pb.ApiResponseWithdrawMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthStatusWithdraw{
		Year:  year,
		Month: month,
	}

	records, err := s.withdrawService.FindMonthWithdrawStatusSuccess(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.WithdrawMonthStatusSuccessResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.WithdrawMonthStatusSuccessResponse{
			Year:         record.Year,
			Month:        record.Month,
			TotalSuccess: int32(record.TotalSuccess),
			TotalAmount:  int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseWithdrawMonthStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched withdraw",
		Data:    dataResponses,
	}, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusSuccess(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearStatusSuccess, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	records, err := s.withdrawService.FindYearlyWithdrawStatusSuccess(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.WithdrawYearStatusSuccessResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.WithdrawYearStatusSuccessResponse{
			Year:         record.Year,
			TotalSuccess: int32(record.TotalSuccess),
			TotalAmount:  int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseWithdrawYearStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly Withdraw status success",
		Data:    dataResponses,
	}, nil
}

func (s *withdrawHandleGrpc) FindMonthlyWithdrawStatusFailed(ctx context.Context, req *pb.FindMonthlyWithdrawStatus) (*pb.ApiResponseWithdrawMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthStatusWithdraw{
		Year:  year,
		Month: month,
	}

	records, err := s.withdrawService.FindMonthWithdrawStatusFailed(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.WithdrawMonthStatusFailedResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.WithdrawMonthStatusFailedResponse{
			Year:        record.Year,
			Month:       record.Month,
			TotalFailed: int32(record.TotalFailed),
			TotalAmount: int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseWithdrawMonthStatusFailed{
		Status:  "success",
		Message: "success fetched monthly Withdraw status Failed",
		Data:    dataResponses,
	}, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusFailed(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearStatusFailed, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	records, err := s.withdrawService.FindYearlyWithdrawStatusFailed(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.WithdrawYearStatusFailedResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.WithdrawYearStatusFailedResponse{
			Year:        record.Year,
			TotalFailed: int32(record.TotalFailed),
			TotalAmount: int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseWithdrawYearStatusFailed{
		Status:  "success",
		Message: "success fetched yearly Withdraw status Failed",
		Data:    dataResponses,
	}, nil
}

func (s *withdrawHandleGrpc) FindMonthlyWithdrawStatusSuccessCardNumber(ctx context.Context, req *pb.FindMonthlyWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		return nil, withdraw_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthStatusWithdrawCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.withdrawService.FindMonthWithdrawStatusSuccessByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.WithdrawMonthStatusSuccessResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.WithdrawMonthStatusSuccessResponse{
			Year:         record.Year,
			Month:        record.Month,
			TotalSuccess: int32(record.TotalSuccess),
			TotalAmount:  int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseWithdrawMonthStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched withdraw",
		Data:    dataResponses,
	}, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusSuccessCardNumber(ctx context.Context, req *pb.FindYearWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawYearStatusSuccess, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, withdraw_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearStatusWithdrawCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.withdrawService.FindYearlyWithdrawStatusSuccessByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.WithdrawYearStatusSuccessResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.WithdrawYearStatusSuccessResponse{
			Year:         record.Year,
			TotalSuccess: int32(record.TotalSuccess),
			TotalAmount:  int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseWithdrawYearStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly Withdraw status success",
		Data:    dataResponses,
	}, nil
}

func (s *withdrawHandleGrpc) FindMonthlyWithdrawStatusFailedCardNumber(ctx context.Context, req *pb.FindMonthlyWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		return nil, withdraw_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthStatusWithdrawCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.withdrawService.FindMonthWithdrawStatusFailedByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.WithdrawMonthStatusFailedResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.WithdrawMonthStatusFailedResponse{
			Year:        record.Year,
			Month:       record.Month,
			TotalFailed: int32(record.TotalFailed),
			TotalAmount: int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseWithdrawMonthStatusFailed{
		Status:  "success",
		Message: "Successfully fetched monthly Withdraw status failed",
		Data:    dataResponses,
	}, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusFailedCardNumber(ctx context.Context, req *pb.FindYearWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawYearStatusFailed, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	reqService := requests.YearStatusWithdrawCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.withdrawService.FindYearlyWithdrawStatusFailedByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.WithdrawYearStatusFailedResponse, len(records))
	for i, record := range records {
		dataResponses[i] = &pb.WithdrawYearStatusFailedResponse{
			Year:        record.Year,
			TotalFailed: int32(record.TotalFailed),
			TotalAmount: int32(record.TotalAmount),
		}
	}

	return &pb.ApiResponseWithdrawYearStatusFailed{
		Status:  "success",
		Message: "Successfully fetched yearly Withdraw status failed",
		Data:    dataResponses,
	}, nil
}

func (w *withdrawHandleGrpc) FindMonthlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawMonthAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	withdraws, err := w.withdrawService.FindMonthlyWithdraws(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.WithdrawMonthlyAmountResponse, len(withdraws))
	for i, withdraw := range withdraws {
		dataResponses[i] = &pb.WithdrawMonthlyAmountResponse{
			Month:       withdraw.Month,
			TotalAmount: int32(withdraw.TotalWithdrawAmount),
		}
	}

	return &pb.ApiResponseWithdrawMonthAmount{
		Status:  "success",
		Message: "Successfully fetched monthly withdraws",
		Data:    dataResponses,
	}, nil
}

func (w *withdrawHandleGrpc) FindYearlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	withdraws, err := w.withdrawService.FindYearlyWithdraws(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.WithdrawYearlyAmountResponse, len(withdraws))
	for i, withdraw := range withdraws {
		dataResponses[i] = &pb.WithdrawYearlyAmountResponse{
			Year:        withdraw.Year.Int.String(),
			TotalAmount: int32(withdraw.TotalWithdrawAmount),
		}
	}

	return &pb.ApiResponseWithdrawYearAmount{
		Status:  "success",
		Message: "Successfully fetched yearly withdraws",
		Data:    dataResponses,
	}, nil
}

func (w *withdrawHandleGrpc) FindMonthlyWithdrawsByCardNumber(ctx context.Context, req *pb.FindYearWithdrawCardNumber) (*pb.ApiResponseWithdrawMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, withdraw_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearMonthCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	withdraws, err := w.withdrawService.FindMonthlyWithdrawsByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.WithdrawMonthlyAmountResponse, len(withdraws))
	for i, withdraw := range withdraws {
		dataResponses[i] = &pb.WithdrawMonthlyAmountResponse{
			Month:       withdraw.Month,
			TotalAmount: int32(withdraw.TotalWithdrawAmount),
		}
	}

	return &pb.ApiResponseWithdrawMonthAmount{
		Status:  "success",
		Message: "Successfully fetched monthly withdraws by card number",
		Data:    dataResponses,
	}, nil
}

func (w *withdrawHandleGrpc) FindYearlyWithdrawsByCardNumber(ctx context.Context, req *pb.FindYearWithdrawCardNumber) (*pb.ApiResponseWithdrawYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, withdraw_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearMonthCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	withdraws, err := w.withdrawService.FindYearlyWithdrawsByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	dataResponses := make([]*pb.WithdrawYearlyAmountResponse, len(withdraws))
	for i, withdraw := range withdraws {
		dataResponses[i] = &pb.WithdrawYearlyAmountResponse{
			Year:        withdraw.Year.Int.String(),
			TotalAmount: int32(withdraw.TotalWithdrawAmount),
		}
	}

	return &pb.ApiResponseWithdrawYearAmount{
		Status:  "success",
		Message: "Successfully fetched yearly withdraws by card number",
		Data:    dataResponses,
	}, nil
}

func (w *withdrawHandleGrpc) CreateWithdraw(ctx context.Context, req *pb.CreateWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	if req == nil {
		return nil, errors.ToGrpcError(withdraw_errors.ErrGrpcValidateCreateWithdrawRequest)
	}

	withdrawTime := time.Now()
	if req.WithdrawTime != nil {
		withdrawTime = req.WithdrawTime.AsTime()
	}

	request := requests.CreateWithdrawRequest{
		CardNumber:     req.CardNumber,
		WithdrawAmount: int(req.WithdrawAmount),
		WithdrawTime:   withdrawTime,
	}

	if err := request.Validate(); err != nil {
		return nil, withdraw_errors.ErrGrpcValidateCreateWithdrawRequest
	}

	withdraw, err := w.withdrawService.Create(ctx, &request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	if withdraw == nil {
		return nil, errors.ToGrpcError(withdraw_errors.ErrWithdrawNotFound)
	}

	return &pb.ApiResponseWithdraw{
		Status:  "success",
		Message: "Successfully created withdraw",
		Data: &pb.WithdrawResponse{
			WithdrawId:     int32(withdraw.WithdrawID),
			WithdrawNo:     withdraw.WithdrawNo.String(),
			CardNumber:     withdraw.CardNumber,
			WithdrawAmount: int32(withdraw.WithdrawAmount),
			WithdrawTime:   withdraw.WithdrawTime.Format(time.RFC3339),
			CreatedAt:      withdraw.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      withdraw.UpdatedAt.Time.Format(time.RFC3339),
		},
	}, nil
}

func (w *withdrawHandleGrpc) UpdateWithdraw(ctx context.Context, req *pb.UpdateWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	if req == nil {
		return nil, errors.ToGrpcError(withdraw_errors.ErrGrpcValidateUpdateWithdrawRequest)
	}
	id := int(req.GetWithdrawId())

	if id == 0 {
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdrawTime := time.Now()
	if req.WithdrawTime != nil {
		withdrawTime = req.WithdrawTime.AsTime()
	}

	request := requests.UpdateWithdrawRequest{
		WithdrawID:     &id,
		CardNumber:     req.CardNumber,
		WithdrawAmount: int(req.WithdrawAmount),
		WithdrawTime:   withdrawTime,
	}

	if err := request.Validate(); err != nil {
		return nil, withdraw_errors.ErrGrpcValidateUpdateWithdrawRequest
	}

	withdraw, err := w.withdrawService.Update(ctx, &request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	if withdraw == nil {
		return nil, errors.ToGrpcError(withdraw_errors.ErrWithdrawNotFound)
	}

	return &pb.ApiResponseWithdraw{
		Status:  "success",
		Message: "Successfully updated withdraw",
		Data: &pb.WithdrawResponse{
			WithdrawId:     int32(withdraw.WithdrawID),
			WithdrawNo:     withdraw.WithdrawNo.String(),
			CardNumber:     withdraw.CardNumber,
			WithdrawAmount: int32(withdraw.WithdrawAmount),
			WithdrawTime:   withdraw.WithdrawTime.Format(time.RFC3339),
			CreatedAt:      withdraw.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      withdraw.UpdatedAt.Time.Format(time.RFC3339),
		},
	}, nil
}

func (w *withdrawHandleGrpc) TrashedWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdrawDeleteAt, error) {
	id := int(req.GetWithdrawId())

	if id == 0 {
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdraw, err := w.withdrawService.TrashedWithdraw(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseWithdrawDeleteAt{
		Status:  "success",
		Message: "Successfully trashed withdraw",
		Data: &pb.WithdrawResponseDeleteAt{
			WithdrawId:     int32(withdraw.WithdrawID),
			WithdrawNo:     withdraw.WithdrawNo.String(),
			CardNumber:     withdraw.CardNumber,
			WithdrawAmount: int32(withdraw.WithdrawAmount),
			WithdrawTime:   withdraw.WithdrawTime.Format(time.RFC3339),
			CreatedAt:      withdraw.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      withdraw.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:      &wrapperspb.StringValue{Value: withdraw.DeletedAt.Time.Format(time.RFC3339)},
		},
	}, nil
}

func (w *withdrawHandleGrpc) RestoreWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdrawDeleteAt, error) {
	id := int(req.GetWithdrawId())

	if id == 0 {
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdraw, err := w.withdrawService.RestoreWithdraw(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseWithdrawDeleteAt{
		Status:  "success",
		Message: "Successfully restored withdraw",
		Data: &pb.WithdrawResponseDeleteAt{
			WithdrawId:     int32(withdraw.WithdrawID),
			WithdrawNo:     withdraw.WithdrawNo.String(),
			CardNumber:     withdraw.CardNumber,
			WithdrawAmount: int32(withdraw.WithdrawAmount),
			WithdrawTime:   withdraw.WithdrawTime.Format(time.RFC3339),
			CreatedAt:      withdraw.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      withdraw.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:      &wrapperspb.StringValue{Value: withdraw.DeletedAt.Time.Format(time.RFC3339)},
		},
	}, nil
}

func (w *withdrawHandleGrpc) DeleteWithdrawPermanent(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdrawDelete, error) {
	id := int(req.GetWithdrawId())

	if id == 0 {
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	_, err := w.withdrawService.DeleteWithdrawPermanent(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseWithdrawDelete{
		Status:  "success",
		Message: "Successfully deleted withdraw permanently",
	}, nil
}

func (s *withdrawHandleGrpc) RestoreAllWithdraw(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error) {
	_, err := s.withdrawService.RestoreAllWithdraw(ctx)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseWithdrawAll{
		Status:  "success",
		Message: "Successfully restore all withdraw",
	}, nil
}

func (s *withdrawHandleGrpc) DeleteAllWithdrawPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error) {
	_, err := s.withdrawService.DeleteAllWithdrawPermanent(ctx)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseWithdrawAll{
		Status:  "success",
		Message: "Successfully delete withdraw permanent",
	}, nil
}
