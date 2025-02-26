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

	topups, totalRecords, err := s.topupService.FindAll(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch topups: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
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

	topups, totalRecords, err := s.topupService.FindAllByCardNumber(card_number, page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch topups: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationTopup(paginationMeta, "success", "Successfully fetch topups", topups)

	return so, nil
}

func (s *topupHandleGrpc) FindByIdTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopup, error) {
	if req.GetTopupId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid ID",
		})
	}

	id := req.GetTopupId()

	topup, err := s.topupService.FindById(int(id))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch topup: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully fetch topup", topup)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusSuccess(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusSuccess, error) {
	if req.GetYear() <= 0 || req.GetMonth() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year or month",
		})
	}

	year := req.GetYear()
	month := req.GetMonth()

	records, errResponse := s.topupService.FindMonthTopupStatusSuccess(int(year), int(month))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup status success: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupMonthStatusSuccess("success", "Successfully fetched monthly topup status success", records)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusSuccess(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusSuccess, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year",
		})
	}

	year := req.GetYear()

	records, errResponse := s.topupService.FindYearlyTopupStatusSuccess(int(year))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup status success: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupYearStatusSuccess("success", "Successfully fetched yearly topup status success", records)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusFailed(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusFailed, error) {
	if req.GetYear() <= 0 || req.GetMonth() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year or month",
		})
	}

	year := req.GetYear()
	month := req.GetMonth()

	records, errResponse := s.topupService.FindMonthTopupStatusFailed(int(year), int(month))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup status Failed: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupMonthStatusFailed("Successfully", "Successfully fetched monthly topup status Failed", records)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusFailed(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusFailed, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year",
		})
	}

	year := req.GetYear()

	records, errResponse := s.topupService.FindYearlyTopupStatusFailed(int(year))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup status Failed: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupYearStatusFailed("Successfully", "Successfully fetched yearly topup status Failed", records)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusSuccess, error) {
	if req.GetYear() <= 0 || req.GetMonth() <= 0 || req.GetCardNumber() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year, month, or card number",
		})
	}

	year := req.GetYear()
	month := req.GetMonth()
	cardNumber := req.GetCardNumber()

	records, errResponse := s.topupService.FindMonthTopupStatusSuccessByCardNumber(cardNumber, int(year), int(month))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup status success: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupMonthStatusSuccess("success", "Successfully fetched monthly topup status success", records)
	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusSuccess, error) {
	if req.GetYear() <= 0 || req.GetCardNumber() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year or card number",
		})
	}

	year := req.GetYear()
	cardNumber := req.GetCardNumber()

	records, errResponse := s.topupService.FindYearlyTopupStatusSuccessByCardNumber(cardNumber, int(year))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup status success: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupYearStatusSuccess("success", "Successfully fetched yearly topup status success", records)
	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusFailed, error) {
	if req.GetYear() <= 0 || req.GetMonth() <= 0 || req.GetCardNumber() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year, month, or card number",
		})
	}

	year := req.GetYear()
	month := req.GetMonth()
	cardNumber := req.GetCardNumber()

	records, errResponse := s.topupService.FindMonthTopupStatusFailedByCardNumber(cardNumber, int(year), int(month))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup status failed: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupMonthStatusFailed("success", "Successfully fetched monthly topup status failed", records)
	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusFailed, error) {
	if req.GetYear() <= 0 || req.GetCardNumber() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year or card number",
		})
	}

	year := req.GetYear()
	cardNumber := req.GetCardNumber()

	records, errResponse := s.topupService.FindYearlyTopupStatusFailedByCardNumber(cardNumber, int(year))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup status failed: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupYearStatusFailed("success", "Successfully fetched yearly topup status failed", records)
	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthMethod, error) {
	methods, err := s.topupService.FindMonthlyTopupMethods(int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup methods: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupMonthMethod("success", "Successfully fetched monthly topup methods", methods)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearMethod, error) {
	methods, err := s.topupService.FindYearlyTopupMethods(int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup methods: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupYearMethod("success", "Successfully fetched yearly topup methods", methods)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthAmount, error) {
	amounts, err := s.topupService.FindMonthlyTopupAmounts(int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup amounts: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupMonthAmount("success", "Successfully fetched monthly topup amounts", amounts)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearAmount, error) {
	amounts, err := s.topupService.FindYearlyTopupAmounts(int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup amounts: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupYearAmount("success", "Successfully fetched yearly topup amounts", amounts)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthMethod, error) {
	methods, err := s.topupService.FindMonthlyTopupMethodsByCardNumber(req.GetCardNumber(), int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup methods by card number: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupMonthMethod("success", "Successfully fetched monthly topup methods by card number", methods)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearMethod, error) {
	methods, err := s.topupService.FindYearlyTopupMethodsByCardNumber(req.GetCardNumber(), int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup methods by card number: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupYearMethod("success", "Successfully fetched yearly topup methods by card number", methods)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthAmount, error) {
	amounts, err := s.topupService.FindMonthlyTopupAmountsByCardNumber(req.GetCardNumber(), int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup amounts by card number: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupMonthAmount("success", "Successfully fetched monthly topup amounts by card number", amounts)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearAmount, error) {
	amounts, err := s.topupService.FindYearlyTopupAmountsByCardNumber(req.GetCardNumber(), int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup amounts by card number: " + err.Message,
		})
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

	res, totalRecords, err := s.topupService.FindByActive(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch topups: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
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

	res, totalRecords, err := s.topupService.FindByTrashed(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch topups: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
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

	res, err := s.topupService.CreateTopup(&request)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to create topup: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully created topup", res)

	return so, nil
}

func (s *topupHandleGrpc) UpdateTopup(ctx context.Context, req *pb.UpdateTopupRequest) (*pb.ApiResponseTopup, error) {
	if req.GetTopupId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid ID",
		})
	}

	request := requests.UpdateTopupRequest{
		TopupID:     int(req.GetTopupId()),
		CardNumber:  req.GetCardNumber(),
		TopupAmount: int(req.GetTopupAmount()),
		TopupMethod: req.GetTopupMethod(),
	}

	res, err := s.topupService.UpdateTopup(&request)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to update topup: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully updated topup", res)

	return so, nil
}

func (s *topupHandleGrpc) TrashedTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopup, error) {
	if req.GetTopupId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid ID",
		})
	}

	res, err := s.topupService.TrashedTopup(int(req.GetTopupId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to trash topup: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully trashed topup", res)

	return so, nil
}

func (s *topupHandleGrpc) RestoreTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopup, error) {
	if req.GetTopupId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid ID",
		})
	}

	res, err := s.topupService.RestoreTopup(int(req.GetTopupId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore topup: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully restored topup", res)

	return so, nil
}

func (s *topupHandleGrpc) DeleteTopupPermanent(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDelete, error) {
	if req.GetTopupId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid ID",
		})
	}

	_, err := s.topupService.DeleteTopupPermanent(int(req.GetTopupId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete topup permanently: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseTopupDelete("success", "Successfully deleted topup permanently")

	return so, nil
}

func (s *topupHandleGrpc) RestoreAllTopup(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTopupAll, error) {
	_, err := s.topupService.RestoreAllTopup()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all topup: ",
		})
	}

	so := s.mapping.ToProtoResponseTopupAll("success", "Successfully restore all topup")

	return so, nil
}

func (s *topupHandleGrpc) DeleteAllTopupPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTopupAll, error) {
	_, err := s.topupService.DeleteAllTopupPermanent()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete topup permanent: ",
		})
	}

	so := s.mapping.ToProtoResponseTopupAll("success", "Successfully delete topup permanent")

	return so, nil
}
