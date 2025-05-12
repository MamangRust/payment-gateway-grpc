package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors/transfer_errors"
	"context"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
)

type transferHandleGrpc struct {
	pb.UnimplementedTransferServiceServer
	transferService service.TransferService
	mapping         protomapper.TransferProtoMapper
}

func NewTransferHandleGrpc(transferService service.TransferService,
	mapping protomapper.TransferProtoMapper) *transferHandleGrpc {
	return &transferHandleGrpc{
		transferService: transferService,
		mapping:         mapping,
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

	merchants, totalRecords, err := s.transferService.FindAll(&reqService)

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
	so := s.mapping.ToProtoResponsePaginationTransfer(paginationMeta, "success", "Successfully fetch transfer records", merchants)

	return so, nil
}

func (s *transferHandleGrpc) FindByIdTransfer(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error) {
	id := int(request.GetTransferId())

	if id == 0 {
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	transfer, err := s.transferService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfer("success", "Successfully fetch transfer record", transfer)

	return so, nil
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

	records, err := s.transferService.FindMonthTransferStatusSuccess(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferMonthStatusSuccess("success", "Successfully fetched monthly Transfer status success", records)

	return so, nil
}

func (s *transferHandleGrpc) FindYearlyTransferStatusSuccess(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearStatusSuccess, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	records, err := s.transferService.FindYearlyTransferStatusSuccess(year)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearStatusSuccess("success", "Successfully fetched yearly Transfer status success", records)

	return so, nil
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

	records, err := s.transferService.FindMonthTransferStatusFailed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferMonthStatusFailed("success", "success fetched monthly Transfer status Failed", records)

	return so, nil
}

func (s *transferHandleGrpc) FindYearlyTransferStatusFailed(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearStatusFailed, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	records, err := s.transferService.FindYearlyTransferStatusFailed(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearStatusFailed("success", "success fetched yearly Transfer status Failed", records)

	return so, nil
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

	records, err := s.transferService.FindMonthTransferStatusSuccessByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferMonthStatusSuccess("success", "Successfully fetched monthly Transfer status success", records)

	return so, nil
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

	records, err := s.transferService.FindYearlyTransferStatusSuccessByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearStatusSuccess("success", "Successfully fetched yearly Transfer status success", records)

	return so, nil
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

	records, err := s.transferService.FindMonthTransferStatusFailedByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferMonthStatusFailed("success", "success fetched monthly Transfer status Failed", records)

	return so, nil
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

	records, err := s.transferService.FindYearlyTransferStatusFailedByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearStatusFailed("success", "success fetched yearly Transfer status Failed", records)

	return so, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferAmounts(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferMonthAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	amounts, err := s.transferService.FindMonthlyTransferAmounts(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferMonthAmount("success", "Successfully fetched monthly transfer amounts", amounts)

	return so, nil
}

func (s *transferHandleGrpc) FindYearlyTransferAmounts(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	amounts, err := s.transferService.FindYearlyTransferAmounts(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearAmount("success", "Successfully fetched yearly transfer amounts", amounts)

	return so, nil
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

	amounts, err := s.transferService.FindMonthlyTransferAmountsBySenderCardNumber(&reqService)

	if err != nil {
		return nil, transfer_errors.ErrGrpcFailedFindMonthlyTransferAmountsBySenderCardNumber
	}

	so := s.mapping.ToProtoResponseTransferMonthAmount("success", "Successfully fetched monthly transfer amounts by sender card number", amounts)

	return so, nil
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

	amounts, err := s.transferService.FindMonthlyTransferAmountsByReceiverCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferMonthAmount("success", "Successfully fetched monthly transfer amounts by receiver card number", amounts)

	return so, nil
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

	amounts, err := s.transferService.FindYearlyTransferAmountsBySenderCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearAmount("success", "Successfully fetched yearly transfer amounts by sender card number", amounts)

	return so, nil
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

	amounts, err := s.transferService.FindYearlyTransferAmountsByReceiverCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearAmount("success", "Successfully fetched yearly transfer amounts by receiver card number", amounts)

	return so, nil
}

func (s *transferHandleGrpc) FindByTransferByTransferFrom(ctx context.Context, request *pb.FindTransferByTransferFromRequest) (*pb.ApiResponseTransfers, error) {

	transfer_from := request.GetTransferFrom()

	if transfer_from == "" {
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	merchants, err := s.transferService.FindTransferByTransferFrom(transfer_from)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfers("success", "Successfully fetch transfer records", merchants)

	return so, nil
}

func (s *transferHandleGrpc) FindByTransferByTransferTo(ctx context.Context, request *pb.FindTransferByTransferToRequest) (*pb.ApiResponseTransfers, error) {
	transfer_to := request.GetTransferTo()

	if transfer_to == "" {
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	merchants, err := s.transferService.FindTransferByTransferTo(transfer_to)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfers("success", "Successfully fetch transfer records", merchants)

	return so, nil
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

	res, totalRecords, err := s.transferService.FindByActive(&reqService)

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
	so := s.mapping.ToProtoResponsePaginationTransferDeleteAt(paginationMeta, "success", "Successfully fetch transfer records", res)

	return so, nil
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

	res, totalRecords, err := s.transferService.FindByTrashed(&reqService)

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
	so := s.mapping.ToProtoResponsePaginationTransferDeleteAt(paginationMeta, "success", "Successfully fetch transfer records", res)

	return so, nil
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

	res, err := s.transferService.CreateTransaction(&req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfer("success", "Successfully created transfer", res)

	return so, nil
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

	res, err := s.transferService.UpdateTransaction(&req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfer("success", "Successfully updated transfer", res)

	return so, nil
}

func (s *transferHandleGrpc) TrashedTransfer(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error) {
	id := int(request.GetTransferId())

	if id == 0 {
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	res, err := s.transferService.TrashedTransfer(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfer("success", "Successfully trashed transfer", res)

	return so, nil
}

func (s *transferHandleGrpc) RestoreTransfer(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error) {
	id := int(request.GetTransferId())

	if id == 0 {
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	res, err := s.transferService.RestoreTransfer(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfer("success", "Successfully restored transfer", res)

	return so, nil
}

func (s *transferHandleGrpc) DeleteTransferPermanent(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransferDelete, error) {
	id := int(request.GetTransferId())

	if id == 0 {
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	_, err := s.transferService.DeleteTransferPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferDelete("success", "Successfully restored transfer")

	return so, nil
}

func (s *transferHandleGrpc) RestoreAllTransfer(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransferAll, error) {
	_, err := s.transferService.RestoreAllTransfer()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferAll("success", "Successfully restored transfer")

	return so, nil
}

func (s *transferHandleGrpc) DeleteAllTransferPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransferAll, error) {
	_, err := s.transferService.DeleteAllTransferPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferAll("success", "delete transfer permanent")

	return so, nil
}
