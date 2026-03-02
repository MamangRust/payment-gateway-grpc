package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	"context"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type saldoHandleGrpc struct {
	pb.UnimplementedSaldoServiceServer
	saldoService service.SaldoService
}

func NewSaldoHandleGrpc(saldo service.SaldoService) *saldoHandleGrpc {
	return &saldoHandleGrpc{
		saldoService: saldo,
	}
}

func (s *saldoHandleGrpc) FindAllSaldo(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldo, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllSaldos{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.saldoService.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSaldos := make([]*pb.SaldoResponse, len(res))
	for i, saldo := range res {
		protoSaldos[i] = &pb.SaldoResponse{
			SaldoId:        int32(saldo.SaldoID),
			CardNumber:     saldo.CardNumber,
			TotalBalance:   saldo.TotalBalance,
			WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02 15:04:05"),
			WithdrawAmount: *saldo.WithdrawAmount,
			CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationSaldo{
		Status:     "success",
		Message:    "Successfully fetched saldo record",
		Data:       protoSaldos,
		Pagination: paginationMeta,
	}, nil
}

func (s *saldoHandleGrpc) FindByIdSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldo, error) {
	id := int(req.GetSaldoId())
	if id == 0 {
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	saldo, err := s.saldoService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSaldo := &pb.SaldoResponse{
		SaldoId:        int32(saldo.SaldoID),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   saldo.TotalBalance,
		WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02 15:04:05"),
		WithdrawAmount: *saldo.WithdrawAmount,
		CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Successfully fetched saldo record",
		Data:    protoSaldo,
	}, nil
}

func (s *saldoHandleGrpc) FindByCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponseSaldo, error) {
	cardNumber := req.GetCardNumber()
	if cardNumber == "" {
		return nil, saldo_errors.ErrGrpcSaldoInvalidCardNumber
	}

	saldo, err := s.saldoService.FindByCardNumber(ctx, cardNumber)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSaldo := &pb.SaldoResponse{
		SaldoId:        int32(saldo.SaldoID),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   saldo.TotalBalance,
		WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02 15:04:05"),
		WithdrawAmount: *saldo.WithdrawAmount,
		CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Successfully fetched saldo record",
		Data:    protoSaldo,
	}, nil
}

func (s *saldoHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldoDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllSaldos{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.saldoService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSaldos := make([]*pb.SaldoResponseDeleteAt, len(res))
	for i, saldo := range res {
		protoSaldos[i] = &pb.SaldoResponseDeleteAt{
			SaldoId:        int32(saldo.SaldoID),
			CardNumber:     saldo.CardNumber,
			TotalBalance:   saldo.TotalBalance,
			WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02 15:04:05"),
			WithdrawAmount: *saldo.WithdrawAmount,
			CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:      wrapperspb.String(saldo.DeletedAt.Time.Format("2006-01-02 15:04:05")),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationSaldoDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched saldo record",
		Data:       protoSaldos,
		Pagination: paginationMeta,
	}, nil
}

func (s *saldoHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldoDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllSaldos{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.saldoService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSaldos := make([]*pb.SaldoResponseDeleteAt, len(res))
	for i, saldo := range res {
		protoSaldos[i] = &pb.SaldoResponseDeleteAt{
			SaldoId:        int32(saldo.SaldoID),
			CardNumber:     saldo.CardNumber,
			TotalBalance:   saldo.TotalBalance,
			WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02 15:04:05"),
			WithdrawAmount: *saldo.WithdrawAmount,
			CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:      wrapperspb.String(saldo.DeletedAt.Time.Format("2006-01-02 15:04:05")),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationSaldoDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched saldo record",
		Data:       protoSaldos,
		Pagination: paginationMeta,
	}, nil
}

func (s *saldoHandleGrpc) FindMonthlyTotalSaldoBalance(ctx context.Context, req *pb.FindMonthlySaldoTotalBalance) (*pb.ApiResponseMonthTotalSaldo, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, saldo_errors.ErrGrpcSaldoInvalidYear
	}
	if month <= 0 {
		return nil, saldo_errors.ErrGrpcSaldoInvalidMonth
	}

	reqService := requests.MonthTotalSaldoBalance{
		Year:  year,
		Month: month,
	}

	res, err := s.saldoService.FindMonthlyTotalSaldoBalance(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.SaldoMonthTotalBalanceResponse, len(res))
	for i, item := range res {
		protoData[i] = &pb.SaldoMonthTotalBalanceResponse{
			Month:        item.Month,
			Year:         item.Year,
			TotalBalance: int32(item.TotalBalance),
		}
	}

	return &pb.ApiResponseMonthTotalSaldo{
		Status:  "success",
		Message: "Successfully fetched monthly total saldo balance",
		Data:    protoData,
	}, nil
}

func (s *saldoHandleGrpc) FindYearTotalSaldoBalance(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseYearTotalSaldo, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, saldo_errors.ErrGrpcSaldoInvalidYear
	}

	res, err := s.saldoService.FindYearTotalSaldoBalance(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.SaldoYearTotalBalanceResponse, len(res))
	for i, item := range res {
		protoData[i] = &pb.SaldoYearTotalBalanceResponse{
			Year:         item.Year,
			TotalBalance: int32(item.TotalBalance),
		}
	}

	return &pb.ApiResponseYearTotalSaldo{
		Status:  "success",
		Message: "Successfully fetched yearly total saldo balance",
		Data:    protoData,
	}, nil
}

func (s *saldoHandleGrpc) FindMonthlySaldoBalances(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseMonthSaldoBalances, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, saldo_errors.ErrGrpcSaldoInvalidYear
	}

	res, err := s.saldoService.FindMonthlySaldoBalances(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.SaldoMonthBalanceResponse, len(res))
	for i, item := range res {
		protoData[i] = &pb.SaldoMonthBalanceResponse{
			Month:        item.Month,
			TotalBalance: int32(item.TotalBalance),
		}
	}

	return &pb.ApiResponseMonthSaldoBalances{
		Status:  "success",
		Message: "Successfully fetched monthly saldo balances",
		Data:    protoData,
	}, nil
}

func (s *saldoHandleGrpc) FindYearlySaldoBalances(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseYearSaldoBalances, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, saldo_errors.ErrGrpcSaldoInvalidYear
	}

	res, err := s.saldoService.FindYearlySaldoBalances(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.SaldoYearBalanceResponse, len(res))
	for i, item := range res {
		protoData[i] = &pb.SaldoYearBalanceResponse{
			Year:         item.Year.Int.String(),
			TotalBalance: int32(item.TotalBalance),
		}
	}

	return &pb.ApiResponseYearSaldoBalances{
		Status:  "success",
		Message: "Successfully fetched yearly saldo balances",
		Data:    protoData,
	}, nil
}

func (s *saldoHandleGrpc) CreateSaldo(ctx context.Context, req *pb.CreateSaldoRequest) (*pb.ApiResponseSaldo, error) {
	request := requests.CreateSaldoRequest{
		CardNumber:   req.GetCardNumber(),
		TotalBalance: int(req.GetTotalBalance()),
	}

	if err := request.Validate(); err != nil {
		return nil, saldo_errors.ErrGrpcValidateCreateSaldo
	}

	saldo, err := s.saldoService.CreateSaldo(ctx, &request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSaldo := &pb.SaldoResponse{
		SaldoId:        int32(saldo.SaldoID),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   int32(saldo.TotalBalance),
		WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02 15:04:05"),
		WithdrawAmount: int32(*saldo.WithdrawAmount),
		CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Successfully created saldo record",
		Data:    protoSaldo,
	}, nil
}

func (s *saldoHandleGrpc) UpdateSaldo(ctx context.Context, req *pb.UpdateSaldoRequest) (*pb.ApiResponseSaldo, error) {
	id := int(req.GetSaldoId())

	if id == 0 {
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	request := requests.UpdateSaldoRequest{
		SaldoID:      &id,
		CardNumber:   req.GetCardNumber(),
		TotalBalance: int(req.GetTotalBalance()),
	}

	if err := request.Validate(); err != nil {
		return nil, saldo_errors.ErrGrpcValidateUpdateSaldo
	}

	saldo, err := s.saldoService.UpdateSaldo(ctx, &request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSaldo := &pb.SaldoResponse{
		SaldoId:        int32(saldo.SaldoID),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   int32(saldo.TotalBalance),
		WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02 15:04:05"),
		WithdrawAmount: int32(*saldo.WithdrawAmount),
		CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseSaldo{
		Status:  "success",
		Message: "Successfully updated saldo record",
		Data:    protoSaldo,
	}, nil
}

func (s *saldoHandleGrpc) TrashedSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldoDeleteAt, error) {
	id := int(req.GetSaldoId())

	if id == 0 {
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	saldo, err := s.saldoService.TrashSaldo(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSaldo := &pb.SaldoResponseDeleteAt{
		SaldoId:        int32(saldo.SaldoID),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   int32(saldo.TotalBalance),
		WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02 15:04:05"),
		WithdrawAmount: int32(*saldo.WithdrawAmount),
		CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      wrapperspb.String(saldo.DeletedAt.Time.Format("2006-01-02 15:04:05")),
	}

	return &pb.ApiResponseSaldoDeleteAt{
		Status:  "success",
		Message: "Successfully trashed saldo record",
		Data:    protoSaldo,
	}, nil
}

func (s *saldoHandleGrpc) RestoreSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldoDeleteAt, error) {
	id := int(req.GetSaldoId())

	if id == 0 {
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	saldo, err := s.saldoService.RestoreSaldo(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSaldo := &pb.SaldoResponseDeleteAt{
		SaldoId:        int32(saldo.SaldoID),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   int32(saldo.TotalBalance),
		WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02 15:04:05"),
		WithdrawAmount: int32(*saldo.WithdrawAmount),
		CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      wrapperspb.String(saldo.DeletedAt.Time.Format("2006-01-02 15:04:05")),
	}

	return &pb.ApiResponseSaldoDeleteAt{
		Status:  "success",
		Message: "Successfully restored saldo record",
		Data:    protoSaldo,
	}, nil
}

func (s *saldoHandleGrpc) DeleteSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldoDelete, error) {
	id := int(req.GetSaldoId())

	if id == 0 {
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	_, err := s.saldoService.DeleteSaldoPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseSaldoDelete{
		Status:  "success",
		Message: "Successfully deleted saldo record",
	}, nil
}

func (s *saldoHandleGrpc) RestoreAllSaldo(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSaldoAll, error) {
	_, err := s.saldoService.RestoreAllSaldo(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseSaldoAll{
		Status:  "success",
		Message: "Successfully restore all saldo",
	}, nil
}

func (s *saldoHandleGrpc) DeleteAllSaldoPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSaldoAll, error) {
	_, err := s.saldoService.DeleteAllSaldoPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseSaldoAll{
		Status:  "success",
		Message: "delete saldo permanent",
	}, nil
}
