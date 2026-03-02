package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/merchant_errors"
	"context"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantHandleGrpc struct {
	pb.UnimplementedMerchantServiceServer
	merchantService service.MerchantService
}

func NewMerchantHandleGrpc(merchantService service.MerchantService) MerchantHandleGrpc {
	return &merchantHandleGrpc{merchantService: merchantService}
}

func (s *merchantHandleGrpc) FindAllMerchant(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchant, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantService.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchants := make([]*pb.MerchantResponse, len(merchants))
	for i, merchant := range merchants {
		protoMerchants[i] = &pb.MerchantResponse{
			Id:        int32(merchant.MerchantID),
			Name:      merchant.Name,
			ApiKey:    merchant.ApiKey,
			Status:    merchant.Status,
			UserId:    int32(merchant.UserID),
			CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchant{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       protoMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantHandleGrpc) FindByIdMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(req.GetMerchantId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	merchant, err := s.merchantService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchant := &pb.MerchantResponse{
		Id:        int32(merchant.MerchantID),
		Name:      merchant.Name,
		ApiKey:    merchant.ApiKey,
		Status:    merchant.Status,
		UserId:    int32(merchant.UserID),
		CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully fetched merchant record",
		Data:    protoMerchant,
	}, nil
}

func (s *merchantHandleGrpc) FindAllTransactionMerchant(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantTransaction, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchantTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := s.merchantService.FindAllTransactions(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTransactions := make([]*pb.MerchantTransactionResponse, len(transactions))
	for i, txn := range transactions {
		protoTransactions[i] = &pb.MerchantTransactionResponse{
			Id:              int32(txn.TransactionID),
			CardNumber:      txn.CardNumber,
			Amount:          int32(txn.Amount),
			PaymentMethod:   txn.PaymentMethod,
			MerchantId:      int32(txn.MerchantID),
			MerchantName:    txn.MerchantName,
			TransactionTime: txn.TransactionTime.Format("2006-01-02 15:04:05"),
			CreatedAt:       txn.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:       txn.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:       wrapperspb.String(txn.DeletedAt.Time.Format("2006-01-02 15:04:05")),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantTransaction{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       protoTransactions,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantHandleGrpc) FindAllTransactionByMerchant(ctx context.Context, req *pb.FindAllMerchantTransaction) (*pb.ApiResponsePaginationMerchantTransaction, error) {
	merchant_id := int(req.MerchantId)
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchantTransactionsById{
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
		MerchantID: merchant_id,
	}

	transactions, totalRecords, err := s.merchantService.FindAllTransactionsByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTransactions := make([]*pb.MerchantTransactionResponse, len(transactions))
	for i, txn := range transactions {
		protoTransactions[i] = &pb.MerchantTransactionResponse{
			Id:              int32(txn.TransactionID),
			CardNumber:      txn.CardNumber,
			Amount:          int32(txn.Amount),
			PaymentMethod:   txn.PaymentMethod,
			MerchantId:      int32(txn.MerchantID),
			MerchantName:    txn.MerchantName,
			TransactionTime: txn.TransactionTime.Format("2006-01-02 15:04:05"),
			CreatedAt:       txn.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:       txn.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:       wrapperspb.String(txn.DeletedAt.Time.Format("2006-01-02 15:04:05")),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantTransaction{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       protoTransactions,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantHandleGrpc) FindAllTransactionByApikey(ctx context.Context, req *pb.FindAllMerchantApikey) (*pb.ApiResponsePaginationMerchantTransaction, error) {
	api_key := req.GetApiKey()
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchantTransactionsByApiKey{
		ApiKey:   api_key,
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := s.merchantService.FindAllTransactionsByApikey(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTransactions := make([]*pb.MerchantTransactionResponse, len(transactions))
	for i, txn := range transactions {
		protoTransactions[i] = &pb.MerchantTransactionResponse{
			Id:              int32(txn.TransactionID),
			CardNumber:      txn.CardNumber,
			Amount:          int32(txn.Amount),
			PaymentMethod:   txn.PaymentMethod,
			MerchantId:      int32(txn.MerchantID),
			MerchantName:    txn.MerchantName,
			TransactionTime: txn.TransactionTime.Format("2006-01-02 15:04:05"),
			CreatedAt:       txn.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:       txn.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:       wrapperspb.String(txn.DeletedAt.Time.Format("2006-01-02 15:04:05")),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantTransaction{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       protoTransactions,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantHandleGrpc) FindMonthlyPaymentMethodsMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.merchantService.FindMonthlyPaymentMethodsMerchant(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseMonthlyPaymentMethod, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseMonthlyPaymentMethod{
			Month:         item.Month,
			PaymentMethod: item.PaymentMethod,
			TotalAmount:   int64(item.TotalAmount),
		}
	}

	return &pb.ApiResponseMerchantMonthlyPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly payment methods for merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindYearlyPaymentMethodMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyPaymentMethod, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.merchantService.FindYearlyPaymentMethodMerchant(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseYearlyPaymentMethod, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseYearlyPaymentMethod{
			Year:          item.Year.Int.String(),
			PaymentMethod: item.PaymentMethod,
			TotalAmount:   item.TotalAmount,
		}
	}

	return &pb.ApiResponseMerchantYearlyPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched yearly payment methods for merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindMonthlyAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.merchantService.FindMonthlyAmountMerchant(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalAmount),
		}
	}

	return &pb.ApiResponseMerchantMonthlyAmount{
		Status:  "success",
		Message: "Successfully fetched monthly amount for merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindYearlyAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.merchantService.FindYearlyAmountMerchant(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalAmount,
		}
	}

	return &pb.ApiResponseMerchantYearlyAmount{
		Status:  "success",
		Message: "Successfully fetched yearly amount for merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindMonthlyTotalAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyTotalAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.merchantService.FindMonthlyTotalAmountMerchant(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseMonthlyTotalAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseMonthlyTotalAmount{
			Month:       item.Month,
			Year:        item.Year,
			TotalAmount: int64(item.TotalAmount),
		}
	}

	return &pb.ApiResponseMerchantMonthlyTotalAmount{
		Status:  "success",
		Message: "Successfully fetched monthly amount for merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindYearlyTotalAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyTotalAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.merchantService.FindYearlyTotalAmountMerchant(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseYearlyTotalAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseYearlyTotalAmount{
			Year:        item.Year,
			TotalAmount: int64(item.TotalAmount),
		}
	}

	return &pb.ApiResponseMerchantYearlyTotalAmount{
		Status:  "success",
		Message: "Successfully fetched yearly amount for merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindMonthlyPaymentMethodByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error) {
	merchantId := req.GetMerchantId()
	year := req.GetYear()

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if merchantId <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	reqService := requests.MonthYearPaymentMethodMerchant{
		MerchantID: int(req.MerchantId),
		Year:       int(year),
	}

	res, err := s.merchantService.FindMonthlyPaymentMethodByMerchants(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseMonthlyPaymentMethod, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseMonthlyPaymentMethod{
			Month:         item.Month,
			PaymentMethod: item.PaymentMethod,
			TotalAmount:   int64(item.TotalAmount),
		}
	}

	return &pb.ApiResponseMerchantMonthlyPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly payment methods by merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindYearlyPaymentMethodByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantYearlyPaymentMethod, error) {
	merchantId := req.GetMerchantId()
	year := req.GetYear()

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if merchantId <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	reqService := requests.MonthYearPaymentMethodMerchant{
		MerchantID: int(req.MerchantId),
		Year:       int(year),
	}

	res, err := s.merchantService.FindYearlyPaymentMethodByMerchants(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseYearlyPaymentMethod, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseYearlyPaymentMethod{
			Year:          item.Year.Int.String(),
			PaymentMethod: item.PaymentMethod,
			TotalAmount:   item.TotalAmount,
		}
	}

	return &pb.ApiResponseMerchantYearlyPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched yearly payment methods by merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindMonthlyAmountByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantMonthlyAmount, error) {
	merchantId := req.GetMerchantId()
	year := req.GetYear()

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if merchantId <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	reqService := requests.MonthYearAmountMerchant{
		MerchantID: int(req.MerchantId),
		Year:       int(year),
	}

	res, err := s.merchantService.FindMonthlyAmountByMerchants(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalAmount),
		}
	}

	return &pb.ApiResponseMerchantMonthlyAmount{
		Status:  "success",
		Message: "Successfully fetched monthly amount by merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindYearlyAmountByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantYearlyAmount, error) {
	merchantId := req.GetMerchantId()
	year := req.GetYear()

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if merchantId <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	reqService := requests.MonthYearAmountMerchant{
		MerchantID: int(req.MerchantId),
		Year:       int(year),
	}
	res, err := s.merchantService.FindYearlyAmountByMerchants(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalAmount,
		}
	}

	return &pb.ApiResponseMerchantYearlyAmount{
		Status:  "success",
		Message: "Successfully fetched yearly amount by merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindMonthlyPaymentMethodByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error) {
	api_key := req.GetApiKey()
	year := req.GetYear()

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if api_key == "" {
		return nil, merchant_errors.ErrGrpcMerchantInvalidApiKey
	}

	reqService := requests.MonthYearPaymentMethodApiKey{
		Year:   int(year),
		Apikey: api_key,
	}

	res, err := s.merchantService.FindMonthlyPaymentMethodByApikey(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseMonthlyPaymentMethod, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseMonthlyPaymentMethod{
			Month:         item.Month,
			PaymentMethod: item.PaymentMethod,
			TotalAmount:   int64(item.TotalAmount),
		}
	}

	return &pb.ApiResponseMerchantMonthlyPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly payment methods by merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindYearlyPaymentMethodByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantYearlyPaymentMethod, error) {
	api_key := req.GetApiKey()
	year := req.GetYear()

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if api_key == "" {
		return nil, merchant_errors.ErrGrpcMerchantInvalidApiKey
	}

	reqService := requests.MonthYearPaymentMethodApiKey{
		Year:   int(year),
		Apikey: api_key,
	}

	res, err := s.merchantService.FindYearlyPaymentMethodByApikey(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseYearlyPaymentMethod, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseYearlyPaymentMethod{
			Year:          item.Year.Int.String(),
			PaymentMethod: item.PaymentMethod,
			TotalAmount:   item.TotalAmount,
		}
	}

	return &pb.ApiResponseMerchantYearlyPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched yearly payment methods by merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindMonthlyAmountByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantMonthlyAmount, error) {
	api_key := req.GetApiKey()
	year := req.GetYear()

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if api_key == "" {
		return nil, merchant_errors.ErrGrpcMerchantInvalidApiKey
	}

	reqService := requests.MonthYearAmountApiKey{
		Apikey: api_key,
		Year:   int(year),
	}

	res, err := s.merchantService.FindMonthlyAmountByApikey(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalAmount),
		}
	}

	return &pb.ApiResponseMerchantMonthlyAmount{
		Status:  "success",
		Message: "Successfully fetched monthly amount by merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindYearlyAmountByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantYearlyAmount, error) {
	api_key := req.GetApiKey()
	year := req.GetYear()

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if api_key == "" {
		return nil, merchant_errors.ErrGrpcMerchantInvalidApiKey
	}

	reqService := requests.MonthYearAmountApiKey{
		Apikey: api_key,
		Year:   int(year),
	}

	res, err := s.merchantService.FindYearlyAmountByApikey(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.MerchantResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.MerchantResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalAmount,
		}
	}

	return &pb.ApiResponseMerchantYearlyAmount{
		Status:  "success",
		Message: "Successfully fetched yearly amount by merchant",
		Data:    protoData,
	}, nil
}

func (s *merchantHandleGrpc) FindByApiKey(ctx context.Context, req *pb.FindByApiKeyRequest) (*pb.ApiResponseMerchant, error) {
	api_key := req.GetApiKey()

	if api_key == "" {
		return nil, merchant_errors.ErrGrpcMerchantInvalidApiKey
	}

	merchant, err := s.merchantService.FindByApiKey(ctx, api_key)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchant := &pb.MerchantResponse{
		Id:        int32(merchant.MerchantID),
		Name:      merchant.Name,
		ApiKey:    merchant.ApiKey,
		Status:    merchant.Status,
		UserId:    int32(merchant.UserID),
		CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully fetched merchant record",
		Data:    protoMerchant,
	}, nil
}

func (s *merchantHandleGrpc) FindByMerchantUserId(ctx context.Context, req *pb.FindByMerchantUserIdRequest) (*pb.ApiResponsesMerchant, error) {
	user_id := req.GetUserId()

	if user_id <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidUserID
	}

	res, err := s.merchantService.FindByMerchantUserId(ctx, int(user_id))

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchant := make([]*pb.MerchantResponse, 0, len(res))

	for _, merchant := range res {
		protoMerchant = append(protoMerchant, &pb.MerchantResponse{
			Id:        int32(merchant.MerchantID),
			Name:      merchant.Name,
			ApiKey:    merchant.ApiKey,
			Status:    merchant.Status,
			UserId:    int32(merchant.UserID),
			CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.ApiResponsesMerchant{
		Status:  "success",
		Message: "Successfully fetched merchant record",
		Data:    protoMerchant,
	}, nil
}

func (s *merchantHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.merchantService.FindByActive(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchants := make([]*pb.MerchantResponseDeleteAt, len(res))
	for i, merchant := range res {
		protoMerchants[i] = &pb.MerchantResponseDeleteAt{
			Id:        int32(merchant.MerchantID),
			Name:      merchant.Name,
			ApiKey:    merchant.ApiKey,
			Status:    merchant.Status,
			UserId:    int32(merchant.UserID),
			CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt: &wrapperspb.StringValue{Value: merchant.DeletedAt.Time.Format("2006-01-02 15:04:05")},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       protoMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.merchantService.FindByTrashed(ctx, &reqService)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchants := make([]*pb.MerchantResponseDeleteAt, len(res))
	for i, merchant := range res {
		protoMerchants[i] = &pb.MerchantResponseDeleteAt{
			Id:        int32(merchant.MerchantID),
			Name:      merchant.Name,
			ApiKey:    merchant.ApiKey,
			Status:    merchant.Status,
			UserId:    int32(merchant.UserID),
			CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt: &wrapperspb.StringValue{Value: merchant.DeletedAt.Time.Format("2006-01-02 15:04:05")},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       protoMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantHandleGrpc) CreateMerchant(ctx context.Context, req *pb.CreateMerchantRequest) (*pb.ApiResponseMerchant, error) {
	request := requests.CreateMerchantRequest{
		Name:   req.GetName(),
		UserID: int(req.GetUserId()),
	}

	if err := request.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateCreateMerchant
	}

	merchant, err := s.merchantService.CreateMerchant(ctx, &request)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchant := &pb.MerchantResponse{
		Id:        int32(merchant.MerchantID),
		Name:      merchant.Name,
		ApiKey:    merchant.ApiKey,
		Status:    merchant.Status,
		UserId:    int32(merchant.UserID),
		CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully created merchant",
		Data:    protoMerchant,
	}, nil
}

func (s *merchantHandleGrpc) UpdateMerchant(ctx context.Context, req *pb.UpdateMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	request := requests.UpdateMerchantRequest{
		MerchantID: &id,
		Name:       req.GetName(),
		UserID:     int(req.GetUserId()),
		Status:     req.GetStatus(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchant
	}

	merchant, err := s.merchantService.UpdateMerchant(ctx, &request)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchant := &pb.MerchantResponse{
		Id:        int32(merchant.MerchantID),
		Name:      merchant.Name,
		ApiKey:    merchant.ApiKey,
		Status:    merchant.Status,
		UserId:    int32(merchant.UserID),
		CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully updated merchant",
		Data:    protoMerchant,
	}, nil
}

func (s *merchantHandleGrpc) TrashedMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDeleteAt, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	merchant, err := s.merchantService.TrashedMerchant(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchant := &pb.MerchantResponseDeleteAt{
		Id:        int32(merchant.MerchantID),
		Name:      merchant.Name,
		ApiKey:    merchant.ApiKey,
		Status:    merchant.Status,
		UserId:    int32(merchant.UserID),
		CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt: wrapperspb.String(merchant.DeletedAt.Time.Format("2006-01-02 15:04:05")),
	}

	return &pb.ApiResponseMerchantDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant",
		Data:    protoMerchant,
	}, nil
}

func (s *merchantHandleGrpc) RestoreMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDeleteAt, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	merchant, err := s.merchantService.RestoreMerchant(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchant := &pb.MerchantResponseDeleteAt{
		Id:        int32(merchant.MerchantID),
		Name:      merchant.Name,
		ApiKey:    merchant.ApiKey,
		Status:    merchant.Status,
		UserId:    int32(merchant.UserID),
		CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt: wrapperspb.String(merchant.DeletedAt.Time.Format("2006-01-02 15:04:05")),
	}

	return &pb.ApiResponseMerchantDeleteAt{
		Status:  "success",
		Message: "Successfully restored merchant",
		Data:    protoMerchant,
	}, nil
}

func (s *merchantHandleGrpc) DeleteMerchantPermanent(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	_, err := s.merchantService.DeleteMerchantPermanent(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant",
	}, nil
}

func (s *merchantHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantService.RestoreAllMerchant(ctx)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restore all merchant",
	}, nil
}

func (s *merchantHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantService.DeleteAllMerchantPermanent(ctx)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully delete all merchant",
	}, nil
}
