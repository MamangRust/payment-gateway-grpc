package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors/withdraw_errors"
	"context"
	"math"

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

	reqService := requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	withdraws, totalRecords, err := w.withdrawService.FindAll(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
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

	reqService := requests.FindAllWithdrawCardNumber{
		CardNumber: card_number,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	withdraws, totalRecords, err := w.withdrawService.FindAllByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := w.mapping.ToProtoResponsePaginationWithdraw(paginationMeta, "success", "Withdraws fetched successfully", withdraws)

	return so, nil
}

func (w *withdrawHandleGrpc) FindByIdWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	id := int(req.GetWithdrawId())

	if id == 0 {
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdraw, err := w.withdrawService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully fetched withdraw", withdraw)

	return so, nil
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

	records, err := s.withdrawService.FindMonthWithdrawStatusSuccess(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawMonthStatusSuccess("success", "Successfully fetched withdraw", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusSuccess(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearStatusSuccess, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	records, err := s.withdrawService.FindYearlyWithdrawStatusSuccess(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawYearStatusSuccess("success", "Successfully fetched yearly Withdraw status success", records)

	return so, nil
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

	records, err := s.withdrawService.FindMonthWithdrawStatusFailed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawMonthStatusFailed("success", "success fetched monthly Withdraw status Failed", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusFailed(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearStatusFailed, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	records, err := s.withdrawService.FindYearlyWithdrawStatusFailed(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawYearStatusFailed("success", "success fetched yearly Withdraw status Failed", records)

	return so, nil
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

	records, err := s.withdrawService.FindMonthWithdrawStatusSuccessByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawMonthStatusSuccess("success", "Successfully fetched withdraw", records)

	return so, nil
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

	records, err := s.withdrawService.FindYearlyWithdrawStatusSuccessByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawYearStatusSuccess("success", "Successfully fetched yearly Withdraw status success", records)

	return so, nil
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

	records, err := s.withdrawService.FindMonthWithdrawStatusFailedByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawMonthStatusFailed("success", "Successfully fetched monthly Withdraw status failed", records)

	return so, nil
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

	records, err := s.withdrawService.FindYearlyWithdrawStatusFailedByCardNumber(&reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawYearStatusFailed("success", "Successfully fetched yearly Withdraw status failed", records)

	return so, nil
}

func (w *withdrawHandleGrpc) FindMonthlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawMonthAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	withdraws, err := w.withdrawService.FindMonthlyWithdraws(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdrawMonthAmount("success", "Successfully fetched monthly withdraws", withdraws)

	return so, nil
}

func (w *withdrawHandleGrpc) FindYearlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	withdraws, err := w.withdrawService.FindYearlyWithdraws(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdrawYearAmount("success", "Successfully fetched yearly withdraws", withdraws)

	return so, nil
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

	withdraws, err := w.withdrawService.FindMonthlyWithdrawsByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdrawMonthAmount("success", "Successfully fetched monthly withdraws by card number", withdraws)

	return so, nil
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

	withdraws, err := w.withdrawService.FindYearlyWithdrawsByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
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

	reqService := requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := w.withdrawService.FindByActive(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
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

	reqService := requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := w.withdrawService.FindByTrashed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
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

	if err := request.Validate(); err != nil {
		return nil, withdraw_errors.ErrGrpcValidateCreateWithdrawRequest
	}

	withdraw, err := w.withdrawService.Create(request)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully created withdraw", withdraw)

	return so, nil

}

func (w *withdrawHandleGrpc) UpdateWithdraw(ctx context.Context, req *pb.UpdateWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	id := int(req.GetWithdrawId())

	if id == 0 {
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	request := &requests.UpdateWithdrawRequest{
		WithdrawID:     &id,
		CardNumber:     req.CardNumber,
		WithdrawAmount: int(req.WithdrawAmount),
		WithdrawTime:   req.WithdrawTime.AsTime(),
	}

	if err := request.Validate(); err != nil {
		return nil, withdraw_errors.ErrGrpcValidateUpdateWithdrawRequest
	}

	withdraw, err := w.withdrawService.Update(request)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully updated withdraw", withdraw)

	return so, nil
}

func (w *withdrawHandleGrpc) TrashedWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	id := int(req.GetWithdrawId())

	if id == 0 {
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdraw, err := w.withdrawService.TrashedWithdraw(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully trashed withdraw", withdraw)

	return so, nil
}

func (w *withdrawHandleGrpc) RestoreWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	id := int(req.GetWithdrawId())

	if id == 0 {
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdraw, err := w.withdrawService.RestoreWithdraw(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully restored withdraw", withdraw)

	return so, nil
}

func (w *withdrawHandleGrpc) DeleteWithdrawPermanent(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdrawDelete, error) {
	id := int(req.GetWithdrawId())

	if id == 0 {
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	_, err := w.withdrawService.DeleteWithdrawPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdrawDelete("success", "Successfully deleted withdraw permanently")

	return so, nil
}

func (s *withdrawHandleGrpc) RestoreAllWithdraw(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error) {
	_, err := s.withdrawService.RestoreAllWithdraw()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawAll("success", "Successfully restore all withdraw")

	return so, nil
}

func (s *withdrawHandleGrpc) DeleteAllWithdrawPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error) {
	_, err := s.withdrawService.DeleteAllWithdrawPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawAll("success", "Successfully delete withdraw permanent")

	return so, nil
}
