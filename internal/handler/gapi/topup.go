package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/topup_errors"
	"context"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type topupHandleGrpc struct {
	pb.UnimplementedTopupServiceServer
	topupService service.TopupService
}

func NewTopupHandleGrpc(topup service.TopupService) *topupHandleGrpc {
	return &topupHandleGrpc{
		topupService: topup,
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

	topups, totalRecords, err := s.topupService.FindAll(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTopups := make([]*pb.TopupResponse, len(topups))
	for i, topup := range topups {
		protoTopups[i] = &pb.TopupResponse{
			Id:          int32(topup.TopupID),
			CardNumber:  topup.CardNumber,
			TopupNo:     topup.TopupNo.String(),
			TopupAmount: int32(topup.TopupAmount),
			TopupMethod: topup.TopupMethod,
			TopupTime:   topup.TopupTime.Format("2006-01-02 15:04:05"),
			CreatedAt:   topup.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:   topup.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationTopup{
		Status:     "success",
		Message:    "Successfully fetch topups",
		Data:       protoTopups,
		Pagination: paginationMeta,
	}, nil
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

	topups, totalRecords, err := s.topupService.FindAllByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTopups := make([]*pb.TopupResponse, len(topups))
	for i, topup := range topups {
		protoTopups[i] = &pb.TopupResponse{
			Id:          int32(topup.TopupID),
			CardNumber:  topup.CardNumber,
			TopupNo:     topup.TopupNo.String(),
			TopupAmount: int32(topup.TopupAmount),
			TopupMethod: topup.TopupMethod,
			TopupTime:   topup.TopupTime.Format("2006-01-02 15:04:05"),
			CreatedAt:   topup.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:   topup.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationTopup{
		Status:     "success",
		Message:    "Successfully fetch topups",
		Data:       protoTopups,
		Pagination: paginationMeta,
	}, nil
}

func (s *topupHandleGrpc) FindByIdTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopup, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	topup, err := s.topupService.FindById(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTopup := &pb.TopupResponse{
		Id:          int32(topup.TopupID),
		CardNumber:  topup.CardNumber,
		TopupNo:     topup.TopupNo.String(),
		TopupAmount: int32(topup.TopupAmount),
		TopupMethod: topup.TopupMethod,
		TopupTime:   topup.TopupTime.Format("2006-01-02 15:04:05"),
		CreatedAt:   topup.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:   topup.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseTopup{
		Status:  "success",
		Message: "Successfully fetch topup",
		Data:    protoTopup,
	}, nil
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

	res, totalRecords, err := s.topupService.FindByActive(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTopups := make([]*pb.TopupResponseDeleteAt, len(res))
	for i, topup := range res {
		protoTopups[i] = &pb.TopupResponseDeleteAt{
			Id:          int32(topup.TopupID),
			CardNumber:  topup.CardNumber,
			TopupNo:     topup.TopupNo.String(),
			TopupAmount: int32(topup.TopupAmount),
			TopupMethod: topup.TopupMethod,
			TopupTime:   topup.TopupTime.Format("2006-01-02 15:04:05"),
			CreatedAt:   topup.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:   topup.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:   wrapperspb.String(topup.DeletedAt.Time.Format("2006-01-02 15:04:05")),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationTopupDeleteAt{
		Status:     "success",
		Message:    "Successfully fetch topups",
		Data:       protoTopups,
		Pagination: paginationMeta,
	}, nil
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

	res, totalRecords, err := s.topupService.FindByTrashed(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTopups := make([]*pb.TopupResponseDeleteAt, len(res))
	for i, topup := range res {
		protoTopups[i] = &pb.TopupResponseDeleteAt{
			Id:          int32(topup.TopupID),
			CardNumber:  topup.CardNumber,
			TopupNo:     topup.TopupNo.String(),
			TopupAmount: int32(topup.TopupAmount),
			TopupMethod: topup.TopupMethod,
			TopupTime:   topup.TopupTime.Format("2006-01-02 15:04:05"),
			CreatedAt:   topup.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:   topup.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:   wrapperspb.String(topup.DeletedAt.Time.Format("2006-01-02 15:04:05")),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationTopupDeleteAt{
		Status:     "success",
		Message:    "Successfully fetch topups",
		Data:       protoTopups,
		Pagination: paginationMeta,
	}, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusSuccess(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}
	if month <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidMonth
	}

	reqService := requests.MonthTopupStatus{
		Year:  year,
		Month: month,
	}

	records, err := s.topupService.FindMonthTopupStatusSuccess(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupMonthStatusSuccessResponse, len(records))
	for i, item := range records {
		protoData[i] = &pb.TopupMonthStatusSuccessResponse{
			Year:         item.Year,
			Month:        item.Month,
			TotalSuccess: int32(item.TotalSuccess),
			TotalAmount:  item.TotalAmount,
		}
	}

	return &pb.ApiResponseTopupMonthStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched monthly topup status success",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusSuccess(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusSuccess, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	records, err := s.topupService.FindYearlyTopupStatusSuccess(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupYearStatusSuccessResponse, len(records))
	for i, item := range records {
		protoData[i] = &pb.TopupYearStatusSuccessResponse{
			Year:         item.Year,
			TotalSuccess: item.TotalSuccess,
			TotalAmount:  item.TotalAmount,
		}
	}

	return &pb.ApiResponseTopupYearStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly topup status success",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusFailed(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}
	if month <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidMonth
	}

	reqService := requests.MonthTopupStatus{
		Year:  year,
		Month: month,
	}

	records, err := s.topupService.FindMonthTopupStatusFailed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupMonthStatusFailedResponse, len(records))
	for i, item := range records {
		protoData[i] = &pb.TopupMonthStatusFailedResponse{
			Year:        item.Year,
			Month:       item.Month,
			TotalFailed: int32(item.TotalFailed),
			TotalAmount: item.TotalAmount,
		}
	}

	return &pb.ApiResponseTopupMonthStatusFailed{
		Status:  "success",
		Message: "Successfully fetched monthly topup status failed",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusFailed(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusFailed, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	records, err := s.topupService.FindYearlyTopupStatusFailed(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupYearStatusFailedResponse, len(records))
	for i, item := range records {
		protoData[i] = &pb.TopupYearStatusFailedResponse{
			Year:        item.Year,
			TotalFailed: item.TotalFailed,
			TotalAmount: item.TotalAmount,
		}
	}

	return &pb.ApiResponseTopupYearStatusFailed{
		Status:  "success",
		Message: "Successfully fetched yearly topup status failed",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}
	if month <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidMonth
	}
	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthTopupStatusCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.topupService.FindMonthTopupStatusSuccessByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupMonthStatusSuccessResponse, len(records))
	for i, item := range records {
		protoData[i] = &pb.TopupMonthStatusSuccessResponse{
			Year:         item.Year,
			Month:        item.Month,
			TotalSuccess: int32(item.TotalSuccess),
			TotalAmount:  item.TotalAmount,
		}
	}

	return &pb.ApiResponseTopupMonthStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched monthly topup status success",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusSuccess, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}
	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearTopupStatusCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.topupService.FindYearlyTopupStatusSuccessByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupYearStatusSuccessResponse, len(records))
	for i, item := range records {
		protoData[i] = &pb.TopupYearStatusSuccessResponse{
			Year:         item.Year,
			TotalSuccess: item.TotalSuccess,
			TotalAmount:  item.TotalAmount,
		}
	}

	return &pb.ApiResponseTopupYearStatusSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly topup status success",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}
	if month <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidMonth
	}
	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthTopupStatusCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.topupService.FindMonthTopupStatusFailedByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupMonthStatusFailedResponse, len(records))
	for i, item := range records {
		protoData[i] = &pb.TopupMonthStatusFailedResponse{
			Year:        item.Year,
			Month:       item.Month,
			TotalFailed: int32(item.TotalFailed),
			TotalAmount: item.TotalAmount,
		}
	}

	return &pb.ApiResponseTopupMonthStatusFailed{
		Status:  "success",
		Message: "Successfully fetched monthly topup status failed",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusFailed, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}
	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearTopupStatusCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.topupService.FindYearlyTopupStatusFailedByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupYearStatusFailedResponse, len(records))
	for i, item := range records {
		protoData[i] = &pb.TopupYearStatusFailedResponse{
			Year:        item.Year,
			TotalFailed: item.TotalFailed,
			TotalAmount: item.TotalAmount,
		}
	}

	return &pb.ApiResponseTopupYearStatusFailed{
		Status:  "success",
		Message: "Successfully fetched yearly topup status failed",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	methods, err := s.topupService.FindMonthlyTopupMethods(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupMonthMethodResponse, len(methods))
	for i, item := range methods {
		protoData[i] = &pb.TopupMonthMethodResponse{
			Month:       item.Month,
			TopupMethod: item.TopupMethod,
			TotalTopups: item.TotalTopups,
			TotalAmount: item.TotalAmount,
		}
	}

	return &pb.ApiResponseTopupMonthMethod{
		Status:  "success",
		Message: "Successfully fetched monthly topup methods",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindYearlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	methods, err := s.topupService.FindYearlyTopupMethods(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupYearlyMethodResponse, len(methods))
	for i, item := range methods {
		protoData[i] = &pb.TopupYearlyMethodResponse{
			Year:        item.Year.Int.String(),
			TopupMethod: item.TopupMethod,
			TotalTopups: int32(item.TotalTopups),
			TotalAmount: int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTopupYearMethod{
		Status:  "success",
		Message: "Successfully fetched yearly topup methods",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	amounts, err := s.topupService.FindMonthlyTopupAmounts(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupMonthAmountResponse, len(amounts))
	for i, item := range amounts {
		protoData[i] = &pb.TopupMonthAmountResponse{
			Month:       item.Month,
			TotalAmount: item.TotalAmount,
		}
	}

	return &pb.ApiResponseTopupMonthAmount{
		Status:  "success",
		Message: "Successfully fetched monthly topup amounts",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindYearlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	amounts, err := s.topupService.FindYearlyTopupAmounts(ctx, year)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupYearlyAmountResponse, len(amounts))
	for i, item := range amounts {
		protoData[i] = &pb.TopupYearlyAmountResponse{
			Year:        item.Year.Int.String(),
			TotalAmount: int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTopupYearAmount{
		Status:  "success",
		Message: "Successfully fetched yearly topup amounts",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthMethod, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	methods, err := s.topupService.FindMonthlyTopupMethodsByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupMonthMethodResponse, len(methods))
	for i, item := range methods {
		protoData[i] = &pb.TopupMonthMethodResponse{
			Month:       item.Month,
			TopupMethod: item.TopupMethod,
			TotalTopups: item.TotalTopups,
			TotalAmount: item.TotalAmount,
		}
	}

	return &pb.ApiResponseTopupMonthMethod{
		Status:  "success",
		Message: "Successfully fetched monthly topup methods by card number",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindYearlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearMethod, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	methods, err := s.topupService.FindYearlyTopupMethodsByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupYearlyMethodResponse, len(methods))
	for i, item := range methods {
		protoData[i] = &pb.TopupYearlyMethodResponse{
			Year:        item.Year.Int.String(),
			TopupMethod: item.TopupMethod,
			TotalTopups: int32(item.TotalTopups),
			TotalAmount: int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTopupYearMethod{
		Status:  "success",
		Message: "Successfully fetched yearly topup methods by card number",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.topupService.FindMonthlyTopupAmountsByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupMonthAmountResponse, len(amounts))
	for i, item := range amounts {
		protoData[i] = &pb.TopupMonthAmountResponse{
			Month:       item.Month,
			TotalAmount: item.TotalAmount,
		}
	}

	return &pb.ApiResponseTopupMonthAmount{
		Status:  "success",
		Message: "Successfully fetched monthly topup amounts by card number",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) FindYearlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.topupService.FindYearlyTopupAmountsByCardNumber(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TopupYearlyAmountResponse, len(amounts))
	for i, item := range amounts {
		protoData[i] = &pb.TopupYearlyAmountResponse{
			Year:        item.Year.Int.String(),
			TotalAmount: int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTopupYearAmount{
		Status:  "success",
		Message: "Successfully fetched yearly topup amounts by card number",
		Data:    protoData,
	}, nil
}

func (s *topupHandleGrpc) CreateTopup(ctx context.Context, req *pb.CreateTopupRequest) (*pb.ApiResponseTopup, error) {
	request := requests.CreateTopupRequest{
		CardNumber:  req.GetCardNumber(),
		TopupAmount: int(req.GetTopupAmount()),
		TopupMethod: req.GetTopupMethod(),
	}

	if err := request.Validate(); err != nil {
		return nil, topup_errors.ErrGrpcValidateCreateTopup
	}

	res, err := s.topupService.CreateTopup(ctx, &request)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTopup := &pb.TopupResponse{
		Id:          int32(res.TopupID),
		CardNumber:  res.CardNumber,
		TopupNo:     res.TopupNo.String(),
		TopupAmount: int32(res.TopupAmount),
		TopupMethod: res.TopupMethod,
		TopupTime:   res.TopupTime.Format("2006-01-02 15:04:05"),
		CreatedAt:   res.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:   res.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseTopup{
		Status:  "success",
		Message: "Successfully created topup",
		Data:    protoTopup,
	}, nil
}

func (s *topupHandleGrpc) UpdateTopup(ctx context.Context, req *pb.UpdateTopupRequest) (*pb.ApiResponseTopup, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	request := requests.UpdateTopupRequest{
		TopupID:     &id,
		CardNumber:  req.GetCardNumber(),
		TopupAmount: int(req.GetTopupAmount()),
		TopupMethod: req.GetTopupMethod(),
	}

	if err := request.Validate(); err != nil {
		return nil, topup_errors.ErrGrpcValidateUpdateTopup
	}

	res, err := s.topupService.UpdateTopup(ctx, &request)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTopup := &pb.TopupResponse{
		Id:          int32(res.TopupID),
		CardNumber:  res.CardNumber,
		TopupNo:     res.TopupNo.String(),
		TopupAmount: int32(res.TopupAmount),
		TopupMethod: res.TopupMethod,
		TopupTime:   res.TopupTime.Format("2006-01-02 15:04:05"),
		CreatedAt:   res.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:   res.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseTopup{
		Status:  "success",
		Message: "Successfully updated topup",
		Data:    protoTopup,
	}, nil
}

func (s *topupHandleGrpc) TrashedTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDeleteAt, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	res, err := s.topupService.TrashedTopup(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTopup := &pb.TopupResponseDeleteAt{
		Id:          int32(res.TopupID),
		CardNumber:  res.CardNumber,
		TopupNo:     res.TopupNo.String(),
		TopupAmount: int32(res.TopupAmount),
		TopupMethod: res.TopupMethod,
		TopupTime:   res.TopupTime.Format("2006-01-02 15:04:05"),
		CreatedAt:   res.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:   res.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:   wrapperspb.String(res.DeletedAt.Time.Format("2006-01-02 15:04:05")),
	}

	return &pb.ApiResponseTopupDeleteAt{
		Status:  "success",
		Message: "Successfully trashed topup",
		Data:    protoTopup,
	}, nil
}

func (s *topupHandleGrpc) RestoreTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDeleteAt, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	res, err := s.topupService.RestoreTopup(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTopup := &pb.TopupResponseDeleteAt{
		Id:          int32(res.TopupID),
		CardNumber:  res.CardNumber,
		TopupNo:     res.TopupNo.String(),
		TopupAmount: int32(res.TopupAmount),
		TopupMethod: res.TopupMethod,
		TopupTime:   res.TopupTime.Format("2006-01-02 15:04:05"),
		CreatedAt:   res.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:   res.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:   wrapperspb.String(res.DeletedAt.Time.Format("2006-01-02 15:04:05")),
	}

	return &pb.ApiResponseTopupDeleteAt{
		Status:  "success",
		Message: "Successfully restored topup",
		Data:    protoTopup,
	}, nil
}

func (s *topupHandleGrpc) DeleteTopupPermanent(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDelete, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	_, err := s.topupService.DeleteTopupPermanent(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTopupDelete{
		Status:  "success",
		Message: "Successfully deleted topup permanently",
	}, nil
}

func (s *topupHandleGrpc) RestoreAllTopup(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTopupAll, error) {
	_, err := s.topupService.RestoreAllTopup(ctx)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTopupAll{
		Status:  "success",
		Message: "Successfully restore all topup",
	}, nil
}

func (s *topupHandleGrpc) DeleteAllTopupPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTopupAll, error) {
	_, err := s.topupService.DeleteAllTopupPermanent(ctx)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTopupAll{
		Status:  "success",
		Message: "Successfully delete topup permanent",
	}, nil
}
