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

type withdrawHandleGrpc struct {
	pb.UnimplementedWithdrawServiceServer
	withdrawService service.WithdrawService
	mapping         protomapper.WithdrawalProtoMapper
}

func NewWithdrawHandleGrpc(withdraw service.WithdrawService, mapping protomapper.WithdrawalProtoMapper) *withdrawHandleGrpc {
	return &withdrawHandleGrpc{
		withdrawService: withdraw,
		mapping:         mapping,
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

	withdraws, totalRecords, err := w.withdrawService.FindAll(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch withdraws: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
	so := w.mapping.ToProtoResponsePaginationWithdraw(paginationMeta, "success", "withdraw", withdraws)

	return so, nil
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

	withdraws, totalRecords, err := w.withdrawService.FindAllByCardNumber(card_number, page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch withdraws: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}

	so := w.mapping.ToProtoResponsePaginationWithdraw(paginationMeta, "success", "Withdraws fetched successfully", withdraws)

	return so, nil
}

func (w *withdrawHandleGrpc) FindByIdWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	withdraw, err := w.withdrawService.FindById(int(req.GetWithdrawId()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch withdraw: " + err.Message,
		})
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully fetched withdraw", withdraw)

	return so, nil
}

func (s *withdrawHandleGrpc) FindMonthlyWithdrawStatusSuccess(ctx context.Context, req *pb.FindMonthlyWithdrawStatus) (*pb.ApiResponseWithdrawMonthStatusSuccess, error) {
	if req.GetYear() <= 0 || req.GetMonth() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year or month",
		})
	}

	year := req.GetYear()
	month := req.GetMonth()

	records, errResponse := s.withdrawService.FindMonthWithdrawStatusSuccess(int(year), int(month))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Withdraw status success: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseWithdrawMonthStatusSuccess("success", "Successfully fetched withdraw", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusSuccess(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearStatusSuccess, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year",
		})
	}

	year := req.GetYear()

	records, errResponse := s.withdrawService.FindYearlyWithdrawStatusSuccess(int(year))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Withdraw status success: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseWithdrawYearStatusSuccess("success", "Successfully fetched yearly Withdraw status success", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindMonthlyWithdrawStatusFailed(ctx context.Context, req *pb.FindMonthlyWithdrawStatus) (*pb.ApiResponseWithdrawMonthStatusFailed, error) {
	if req.GetYear() <= 0 || req.GetMonth() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year or month",
		})
	}

	year := req.GetYear()
	month := req.GetMonth()

	records, errResponse := s.withdrawService.FindMonthWithdrawStatusFailed(int(year), int(month))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Withdraw status Failed: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseWithdrawMonthStatusFailed("success", "success fetched monthly Withdraw status Failed", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusFailed(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearStatusFailed, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year",
		})
	}

	year := req.GetYear()

	records, errResponse := s.withdrawService.FindYearlyWithdrawStatusFailed(int(year))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Withdraw status Failed: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseWithdrawYearStatusFailed("success", "success fetched yearly Withdraw status Failed", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindMonthlyWithdrawStatusSuccessCardNumber(ctx context.Context, req *pb.FindMonthlyWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawMonthStatusSuccess, error) {
	if req.GetYear() <= 0 || req.GetMonth() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year or month",
		})
	}

	year := req.GetYear()
	month := req.GetMonth()
	cardNumber := req.GetCardNumber()

	records, errResponse := s.withdrawService.FindMonthWithdrawStatusSuccessByCardNumber(cardNumber, int(year), int(month))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Withdraw status success: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseWithdrawMonthStatusSuccess("success", "Successfully fetched withdraw", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusSuccessCardNumber(ctx context.Context, req *pb.FindYearWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawYearStatusSuccess, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year",
		})
	}

	year := req.GetYear()
	cardNumber := req.GetCardNumber()

	records, errResponse := s.withdrawService.FindYearlyWithdrawStatusSuccessByCardNumber(cardNumber, int(year))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Withdraw status success: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseWithdrawYearStatusSuccess("success", "Successfully fetched yearly Withdraw status success", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindMonthlyWithdrawStatusFailedCardNumber(ctx context.Context, req *pb.FindMonthlyWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawMonthStatusFailed, error) {
	if req.GetYear() <= 0 || req.GetMonth() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year or month",
		})
	}

	year := req.GetYear()
	month := req.GetMonth()
	cardNumber := req.GetCardNumber()

	records, errResponse := s.withdrawService.FindMonthWithdrawStatusFailedByCardNumber(cardNumber, int(year), int(month))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Withdraw status failed: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseWithdrawMonthStatusFailed("success", "Successfully fetched monthly Withdraw status failed", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusFailedCardNumber(ctx context.Context, req *pb.FindYearWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawYearStatusFailed, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Bad Request: Invalid year",
		})
	}

	year := req.GetYear()
	cardNumber := req.GetCardNumber()

	records, errResponse := s.withdrawService.FindYearlyWithdrawStatusFailedByCardNumber(cardNumber, int(year))
	if errResponse != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Withdraw status failed: " + errResponse.Message,
		})
	}

	so := s.mapping.ToProtoResponseWithdrawYearStatusFailed("success", "Successfully fetched yearly Withdraw status failed", records)

	return so, nil
}

func (w *withdrawHandleGrpc) FindMonthlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawMonthAmount, error) {
	withdraws, err := w.withdrawService.FindMonthlyWithdraws(int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly withdraws: " + err.Message,
		})
	}

	so := w.mapping.ToProtoResponseWithdrawMonthAmount("success", "Successfully fetched monthly withdraws", withdraws)

	return so, nil
}

func (w *withdrawHandleGrpc) FindYearlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearAmount, error) {
	withdraws, err := w.withdrawService.FindYearlyWithdraws(int(req.GetYear()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly withdraws: " + err.Message,
		})
	}

	so := w.mapping.ToProtoResponseWithdrawYearAmount("success", "Successfully fetched yearly withdraws", withdraws)

	return so, nil
}

func (w *withdrawHandleGrpc) FindMonthlyWithdrawsByCardNumber(ctx context.Context, req *pb.FindYearWithdrawCardNumber) (*pb.ApiResponseWithdrawMonthAmount, error) {
	withdraws, err := w.withdrawService.FindMonthlyWithdrawsByCardNumber(req.GetCardNumber(), int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly withdraws by card number: " + err.Message,
		})
	}

	so := w.mapping.ToProtoResponseWithdrawMonthAmount("success", "Successfully fetched monthly withdraws by card number", withdraws)

	return so, nil
}

func (w *withdrawHandleGrpc) FindYearlyWithdrawsByCardNumber(ctx context.Context, req *pb.FindYearWithdrawCardNumber) (*pb.ApiResponseWithdrawYearAmount, error) {
	withdraws, err := w.withdrawService.FindYearlyWithdrawsByCardNumber(req.GetCardNumber(), int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly withdraws by card number: " + err.Message,
		})
	}

	so := w.mapping.ToProtoResponseWithdrawYearAmount("success", "Successfully fetched yearly withdraws by card number", withdraws)

	return so, nil
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

	res, totalRecords, err := w.withdrawService.FindByActive(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch withdraws: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
	so := w.mapping.ToProtoResponsePaginationWithdrawDeleteAt(paginationMeta, "success", "Successfully fetched withdraws", res)

	return so, nil
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

	res, totalRecords, err := w.withdrawService.FindByTrashed(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch withdraws: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}

	so := w.mapping.ToProtoResponsePaginationWithdrawDeleteAt(paginationMeta, "success", "Successfully fetched withdraws", res)

	return so, nil
}

func (w *withdrawHandleGrpc) CreateWithdraw(ctx context.Context, req *pb.CreateWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	request := &requests.CreateWithdrawRequest{
		CardNumber:     req.CardNumber,
		WithdrawAmount: int(req.WithdrawAmount),
		WithdrawTime:   req.WithdrawTime.AsTime(),
	}

	withdraw, err := w.withdrawService.Create(request)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to create withdraw: " + err.Message,
		})
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully created withdraw", withdraw)

	return so, nil

}

func (w *withdrawHandleGrpc) UpdateWithdraw(ctx context.Context, req *pb.UpdateWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	request := &requests.UpdateWithdrawRequest{
		WithdrawID:     int(req.WithdrawId),
		CardNumber:     req.CardNumber,
		WithdrawAmount: int(req.WithdrawAmount),
		WithdrawTime:   req.WithdrawTime.AsTime(),
	}

	withdraw, err := w.withdrawService.Update(request)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to update withdraw: " + err.Message,
		})
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully updated withdraw", withdraw)

	return so, nil
}

func (w *withdrawHandleGrpc) TrashedWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	if req.WithdrawId <= 0 {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid withdraw id",
		})
	}

	withdraw, err := w.withdrawService.TrashedWithdraw(int(req.WithdrawId))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch withdraw: " + err.Message,
		})
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully trashed withdraw", withdraw)

	return so, nil
}

func (w *withdrawHandleGrpc) RestoreWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	if req.WithdrawId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid withdraw id",
		})
	}

	withdraw, err := w.withdrawService.RestoreWithdraw(int(req.WithdrawId))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch withdraw: " + err.Message,
		})
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully restored withdraw", withdraw)

	return so, nil
}

func (w *withdrawHandleGrpc) DeleteWithdrawPermanent(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdrawDelete, error) {
	if req.WithdrawId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "invalid withdraw id",
		})
	}

	_, err := w.withdrawService.DeleteWithdrawPermanent(int(req.WithdrawId))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch withdraw: " + err.Message,
		})
	}

	so := w.mapping.ToProtoResponseWithdrawDelete("success", "Successfully deleted withdraw permanently")

	return so, nil
}

func (s *withdrawHandleGrpc) RestoreAllWithdraw(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error) {
	_, err := s.withdrawService.RestoreAllWithdraw()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all withdraw: ",
		})
	}

	so := s.mapping.ToProtoResponseWithdrawAll("success", "Successfully restore all withdraw")

	return so, nil
}

func (s *withdrawHandleGrpc) DeleteAllWithdrawPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error) {
	_, err := s.withdrawService.DeleteAllWithdrawPermanent()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete withdraw permanent: ",
		})
	}

	so := s.mapping.ToProtoResponseWithdrawAll("success", "Successfully delete withdraw permanent")

	return so, nil
}
