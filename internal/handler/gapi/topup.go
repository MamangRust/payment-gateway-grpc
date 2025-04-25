package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors_custom"
	"context"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type topupHandleGrpc struct {
	pb.UnimplementedTopupServiceServer
	topupService service.TopupService
	mapping      protomapper.TopupProtoMapper
}

func NewTopupHandleGrpc(topup service.TopupService, mapping protomapper.TopupProtoMapper) *topupHandleGrpc {
	return &topupHandleGrpc{
		topupService: topup,
		mapping:      mapping,
	}
}

func (s *topupHandleGrpc) FindAllTopup(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopup, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTopups{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	topups, totalRecords, err := s.topupService.FindAll(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationTopup(paginationMeta, "success", "Successfully fetch topups", topups)

	return so, nil
}

func (s *topupHandleGrpc) FindAllTopupByCardNumber(ctx context.Context, req *pb.FindAllTopupByCardNumberRequest) (*pb.ApiResponsePaginationTopup, error) {
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

	reqService := requests.FindAllTopupsByCardNumber{
		CardNumber: card_number,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	topups, totalRecords, err := s.topupService.FindAllByCardNumber(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationTopup(paginationMeta, "success", "Successfully fetch topups", topups)

	return so, nil
}

func (s *topupHandleGrpc) FindByIdTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopup, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_request",
				Message: "Valid topup ID is required",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	topup, err := s.topupService.FindById(id)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully fetch topup", topup)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusSuccess(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if month <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid month parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.MonthTopupStatus{
		Year:  year,
		Month: month,
	}

	records, err := s.topupService.FindMonthTopupStatusSuccess(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupMonthStatusSuccess("success", "Successfully fetched monthly topup status success", records)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusSuccess(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusSuccess, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	records, err := s.topupService.FindYearlyTopupStatusSuccess(year)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupYearStatusSuccess("success", "Successfully fetched yearly topup status success", records)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusFailed(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if month <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid month parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.MonthTopupStatus{
		Year:  year,
		Month: month,
	}

	records, err := s.topupService.FindMonthTopupStatusFailed(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupMonthStatusFailed("Successfully", "Successfully fetched monthly topup status Failed", records)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusFailed(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusFailed, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	records, err := s.topupService.FindYearlyTopupStatusFailed(year)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupYearStatusFailed("Successfully", "Successfully fetched yearly topup status Failed", records)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if month <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid month parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if cardNumber == "" {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid card_number parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.MonthTopupStatusCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.topupService.FindMonthTopupStatusSuccessByCardNumber(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupMonthStatusSuccess("success", "Successfully fetched monthly topup status success", records)
	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusSuccess, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if cardNumber == "" {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid card_number parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.YearTopupStatusCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.topupService.FindYearlyTopupStatusSuccessByCardNumber(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupYearStatusSuccess("success", "Successfully fetched yearly topup status success", records)
	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if month <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid month parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if cardNumber == "" {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid card_number parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.MonthTopupStatusCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.topupService.FindMonthTopupStatusFailedByCardNumber(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupMonthStatusFailed("success", "Successfully fetched monthly topup status failed", records)
	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusFailed, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if cardNumber == "" {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid card_number parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.YearTopupStatusCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.topupService.FindYearlyTopupStatusFailedByCardNumber(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupYearStatusFailed("success", "Successfully fetched yearly topup status failed", records)
	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	methods, err := s.topupService.FindMonthlyTopupMethods(year)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupMonthMethod("success", "Successfully fetched monthly topup methods", methods)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	methods, err := s.topupService.FindYearlyTopupMethods(year)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupYearMethod("success", "Successfully fetched yearly topup methods", methods)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	amounts, err := s.topupService.FindMonthlyTopupAmounts(year)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupMonthAmount("success", "Successfully fetched monthly topup amounts", amounts)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	amounts, err := s.topupService.FindYearlyTopupAmounts(year)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupYearAmount("success", "Successfully fetched yearly topup amounts", amounts)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthMethod, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if cardNumber == "" {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid card_number parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	methods, err := s.topupService.FindMonthlyTopupMethodsByCardNumber(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupMonthMethod("success", "Successfully fetched monthly topup methods by card number", methods)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearMethod, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if cardNumber == "" {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid card_number parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	methods, err := s.topupService.FindYearlyTopupMethodsByCardNumber(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupYearMethod("success", "Successfully fetched yearly topup methods by card number", methods)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if cardNumber == "" {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid card_number parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.topupService.FindMonthlyTopupAmountsByCardNumber(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupMonthAmount("success", "Successfully fetched monthly topup amounts by card number", amounts)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if cardNumber == "" {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid card_number parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.topupService.FindYearlyTopupAmountsByCardNumber(&reqService)
	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupYearAmount("success", "Successfully fetched yearly topup amounts by card number", amounts)

	return so, nil
}

func (s *topupHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopupDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTopups{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.topupService.FindByActive(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationTopupDeleteAt(paginationMeta, "success", "Successfully fetch topups", res)

	return so, nil
}

func (s *topupHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopupDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTopups{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.topupService.FindByTrashed(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationTopupDeleteAt(paginationMeta, "success", "Successfully fetch topups", res)

	return so, nil
}

func (s *topupHandleGrpc) CreateTopup(ctx context.Context, req *pb.CreateTopupRequest) (*pb.ApiResponseTopup, error) {
	request := requests.CreateTopupRequest{
		CardNumber:  req.GetCardNumber(),
		TopupAmount: int(req.GetTopupAmount()),
		TopupMethod: req.GetTopupMethod(),
	}

	if err := request.Validate(); err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Unable to create new topup. Please check your input.",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	res, err := s.topupService.CreateTopup(&request)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully created topup", res)

	return so, nil
}

func (s *topupHandleGrpc) UpdateTopup(ctx context.Context, req *pb.UpdateTopupRequest) (*pb.ApiResponseTopup, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Topup ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	request := requests.UpdateTopupRequest{
		TopupID:     &id,
		CardNumber:  req.GetCardNumber(),
		TopupAmount: int(req.GetTopupAmount()),
		TopupMethod: req.GetTopupMethod(),
	}

	if err := request.Validate(); err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Unable to process topup update. Please review your data.",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	res, err := s.topupService.UpdateTopup(&request)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully updated topup", res)

	return so, nil
}

func (s *topupHandleGrpc) TrashedTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopup, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Merchant ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	res, err := s.topupService.TrashedTopup(id)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully trashed topup", res)

	return so, nil
}

func (s *topupHandleGrpc) RestoreTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopup, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Merchant ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	res, err := s.topupService.RestoreTopup(id)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully restored topup", res)

	return so, nil
}

func (s *topupHandleGrpc) DeleteTopupPermanent(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDelete, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Merchant ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	_, err := s.topupService.DeleteTopupPermanent(id)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupDelete("success", "Successfully deleted topup permanently")

	return so, nil
}

func (s *topupHandleGrpc) RestoreAllTopup(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTopupAll, error) {
	_, err := s.topupService.RestoreAllTopup()

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupAll("success", "Successfully restore all topup")

	return so, nil
}

func (s *topupHandleGrpc) DeleteAllTopupPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTopupAll, error) {
	_, err := s.topupService.DeleteAllTopupPermanent()

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseTopupAll("success", "Successfully delete topup permanent")

	return so, nil
}
