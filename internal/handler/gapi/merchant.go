package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors/merchant_errors"
	"context"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantHandleGrpc struct {
	pb.UnimplementedMerchantServiceServer
	merchantService service.MerchantService
	mapping         protomapper.MerchantProtoMapper
}

func NewMerchantHandleGrpc(merchantService service.MerchantService, mapping protomapper.MerchantProtoMapper) *merchantHandleGrpc {
	return &merchantHandleGrpc{merchantService: merchantService, mapping: mapping}
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

	merchants, totalRecords, err := s.merchantService.FindAll(&reqService)

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
	so := s.mapping.ToProtoResponsePaginationMerchant(paginationMeta, "success", "Successfully fetched merchant record", merchants)

	return so, nil
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

	merchants, totalRecords, err := s.merchantService.FindAllTransactions(&reqService)

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
	so := s.mapping.ToProtoResponsePaginationMerchantTransaction(paginationMeta, "success", "Successfully fetched merchant record", merchants)

	return so, nil
}

func (s *merchantHandleGrpc) FindByIdMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcFailedFindByIdMerchant
	}

	merchant, err := s.merchantService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully fetched merchant record", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) FindMonthlyPaymentMethodsMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.merchantService.FindMonthlyPaymentMethodsMerchant(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyPaymentMethods("success", "Successfully fetched monthly payment methods for merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindYearlyPaymentMethodMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyPaymentMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.merchantService.FindYearlyPaymentMethodMerchant(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyPaymentMethods("success", "Successfully fetched yearly payment methods for merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindMonthlyAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.merchantService.FindMonthlyAmountMerchant(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Successfully fetched monthly amount for merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindYearlyAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.merchantService.FindYearlyAmountMerchant(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Successfully fetched yearly amount for merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindMonthlyTotalAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyTotalAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.merchantService.FindMonthlyTotalAmountMerchant(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyTotalAmounts("success", "Successfully fetched monthly amount for merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindYearlyTotalAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyTotalAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.merchantService.FindYearlyTotalAmountMerchant(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyTotalAmounts("success", "Successfully fetched yearly amount for merchant", res)

	return so, nil
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

	merchants, totalRecords, err := s.merchantService.FindAllTransactionsByMerchant(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantTransaction(paginationMeta, "success", "Successfully fetched merchant record", merchants)

	return so, nil
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

	res, err := s.merchantService.FindMonthlyPaymentMethodByMerchants(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyPaymentMethods("success", "Successfully fetched monthly payment methods by merchant", res)

	return so, nil
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

	res, err := s.merchantService.FindYearlyPaymentMethodByMerchants(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyPaymentMethods("success", "Successfully fetched yearly payment methods by merchant", res)

	return so, nil
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

	res, err := s.merchantService.FindMonthlyAmountByMerchants(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Successfully fetched monthly amount by merchant", res)

	return so, nil
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
	res, err := s.merchantService.FindYearlyAmountByMerchants(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Successfully fetched yearly amount by merchant", res)

	return so, nil
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

	merchants, totalRecords, err := s.merchantService.FindAllTransactionsByApikey(&reqService)
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

	so := s.mapping.ToProtoResponsePaginationMerchantTransaction(paginationMeta, "success", "Successfully fetched merchant record", merchants)

	return so, nil
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

	res, err := s.merchantService.FindMonthlyPaymentMethodByApikeys(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyPaymentMethods("success", "Successfully fetched monthly payment methods by merchant", res)

	return so, nil
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

	res, err := s.merchantService.FindYearlyPaymentMethodByApikeys(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyPaymentMethods("success", "Successfully fetched yearly payment methods by merchant", res)

	return so, nil
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

	res, err := s.merchantService.FindMonthlyAmountByApikeys(&reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Successfully fetched monthly amount by merchant", res)

	return so, nil
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

	res, err := s.merchantService.FindYearlyAmountByApikeys(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Successfully fetched yearly amount by merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindByApiKey(ctx context.Context, req *pb.FindByApiKeyRequest) (*pb.ApiResponseMerchant, error) {
	api_key := req.GetApiKey()

	if api_key == "" {
		return nil, merchant_errors.ErrGrpcMerchantInvalidApiKey
	}

	merchant, err := s.merchantService.FindByApiKey(api_key)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully fetched merchant record", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) FindByMerchantUserId(ctx context.Context, req *pb.FindByMerchantUserIdRequest) (*pb.ApiResponsesMerchant, error) {
	user_id := req.GetUserId()

	if user_id <= 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidUserID
	}

	res, err := s.merchantService.FindByMerchantUserId(int(user_id))

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchants("success", "Successfully fetched merchant record", res)

	return so, nil
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

	res, totalRecords, err := s.merchantService.FindByActive(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantDeleteAt(paginationMeta, "success", "Successfully fetched merchant record", res)

	return so, nil
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

	res, totalRecords, err := s.merchantService.FindByTrashed(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantDeleteAt(paginationMeta, "success", "Successfully fetched merchant record", res)

	return so, nil
}

func (s *merchantHandleGrpc) CreateMerchant(ctx context.Context, req *pb.CreateMerchantRequest) (*pb.ApiResponseMerchant, error) {
	request := requests.CreateMerchantRequest{
		Name:   req.GetName(),
		UserID: int(req.GetUserId()),
	}

	if err := request.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateCreateMerchant
	}

	merchant, err := s.merchantService.CreateMerchant(&request)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully created merchant", merchant)

	return so, nil

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

	merchant, err := s.merchantService.UpdateMerchant(&request)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully updated merchant", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) TrashedMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	merchant, err := s.merchantService.TrashedMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully trashed merchant", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) RestoreMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	merchant, err := s.merchantService.RestoreMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully restored merchant", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) DeleteMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	_, err := s.merchantService.DeleteMerchantPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant")

	return so, nil
}

func (s *merchantHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantService.RestoreAllMerchant()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAll("success", "Successfully restore all merchant")

	return so, nil
}

func (s *merchantHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantService.DeleteAllMerchantPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAll("success", "Successfully delete all merchant")

	return so, nil
}
