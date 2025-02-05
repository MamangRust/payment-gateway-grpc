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

	merchants, totalRecords, err := s.merchantService.FindAll(page, pageSize, search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch card records: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
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

	merchants, totalRecords, err := s.merchantService.FindAllTransactions(page, pageSize, search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch card records: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationMerchantTransaction(paginationMeta, "success", "Successfully fetched merchant record", merchants)

	return so, nil
}

func (s *merchantHandleGrpc) FindByIdMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	if req.GetMerchantId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant ID",
		})
	}

	merchant, err := s.merchantService.FindById(int(req.GetMerchantId()))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: ",
		})
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully fetched merchant record", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) FindMonthlyPaymentMethodsMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid year",
		})
	}

	res, err := s.merchantService.FindMonthlyPaymentMethodsMerchant(int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly payment methods for merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMonthlyPaymentMethods("success", "Successfully fetched monthly payment methods for merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindYearlyPaymentMethodMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyPaymentMethod, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid year",
		})
	}

	res, err := s.merchantService.FindYearlyPaymentMethodMerchant(int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly payment methods for merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseYearlyPaymentMethods("success", "Successfully fetched yearly payment methods for merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindMonthlyAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyAmount, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid year",
		})
	}

	res, err := s.merchantService.FindMonthlyAmountMerchant(int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly amount for merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Successfully fetched monthly amount for merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindYearlyAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyAmount, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid year",
		})
	}

	res, err := s.merchantService.FindYearlyAmountMerchant(int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly amount for merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Successfully fetched yearly amount for merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindMonthlyTotalAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyTotalAmount, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid year",
		})
	}

	res, err := s.merchantService.FindMonthlyTotalAmountMerchant(int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly amount for merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMonthlyTotalAmounts("success", "Successfully fetched monthly amount for merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindYearlyTotalAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyTotalAmount, error) {
	if req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid year",
		})
	}

	res, err := s.merchantService.FindYearlyTotalAmountMerchant(int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly amount for merchant: " + err.Message,
		})
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

	merchants, totalRecords, err := s.merchantService.FindAllTransactionsByMerchant(merchant_id, page, pageSize, search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch card records: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationMerchantTransaction(paginationMeta, "success", "Successfully fetched merchant record", merchants)

	return so, nil
}

func (s *merchantHandleGrpc) FindMonthlyPaymentMethodByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error) {
	if req.GetMerchantId() <= 0 || req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant ID or year",
		})
	}

	res, err := s.merchantService.FindMonthlyPaymentMethodByMerchants(int(req.GetMerchantId()), int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly payment methods by merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMonthlyPaymentMethods("success", "Successfully fetched monthly payment methods by merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindYearlyPaymentMethodByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantYearlyPaymentMethod, error) {
	if req.GetMerchantId() <= 0 || req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant ID or year",
		})
	}

	res, err := s.merchantService.FindYearlyPaymentMethodByMerchants(int(req.GetMerchantId()), int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly payment methods by merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseYearlyPaymentMethods("success", "Successfully fetched yearly payment methods by merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindMonthlyAmountByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantMonthlyAmount, error) {
	if req.GetMerchantId() <= 0 || req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant ID or year",
		})
	}

	res, err := s.merchantService.FindMonthlyAmountByMerchants(int(req.GetMerchantId()), int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly amount by merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Successfully fetched monthly amount by merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindYearlyAmountByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantYearlyAmount, error) {
	if req.GetMerchantId() <= 0 || req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant ID or year",
		})
	}

	res, err := s.merchantService.FindYearlyAmountByMerchants(int(req.GetMerchantId()), int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly amount by merchant: " + err.Message,
		})
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

	merchants, totalRecords, err := s.merchantService.FindAllTransactionsByApikey(api_key, page, pageSize, search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch card records: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationMerchantTransaction(paginationMeta, "success", "Successfully fetched merchant record", merchants)

	return so, nil
}

func (s *merchantHandleGrpc) FindMonthlyPaymentMethodByApikeys(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error) {
	if req.GetApiKey() == "" || req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant ID or year",
		})
	}

	res, err := s.merchantService.FindMonthlyPaymentMethodByApikeys(req.GetApiKey(), int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly payment methods by merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMonthlyPaymentMethods("success", "Successfully fetched monthly payment methods by merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindYearlyPaymentMethodByApikeys(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantYearlyPaymentMethod, error) {
	if req.GetApiKey() == "" || req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant ID or year",
		})
	}

	res, err := s.merchantService.FindYearlyPaymentMethodByApikeys(req.GetApiKey(), int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly payment methods by merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseYearlyPaymentMethods("success", "Successfully fetched yearly payment methods by merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindMonthlyAmountByApikeys(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantMonthlyAmount, error) {
	if req.GetApiKey() == "" || req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant ID or year",
		})
	}

	res, err := s.merchantService.FindMonthlyAmountByApikeys(req.GetApiKey(), int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly amount by merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Successfully fetched monthly amount by merchant", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindYearlyAmountByApikeys(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantYearlyAmount, error) {
	if req.GetApiKey() == "" || req.GetYear() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant ID or year",
		})
	}

	res, err := s.merchantService.FindYearlyAmountByApikeys(req.GetApiKey(), int(req.GetYear()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly amount by merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Successfully fetched yearly amount by merchant", res)

	return so, nil
}

//

func (s *merchantHandleGrpc) FindByApiKey(ctx context.Context, req *pb.FindByApiKeyRequest) (*pb.ApiResponseMerchant, error) {
	merchant, err := s.merchantService.FindByApiKey(req.ApiKey)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: ",
		})
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully fetched merchant record", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) FindByMerchantUserId(ctx context.Context, req *pb.FindByMerchantUserIdRequest) (*pb.ApiResponsesMerchant, error) {
	res, err := s.merchantService.FindByMerchantUserId(int(req.GetUserId()))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: " + err.Message,
		})
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

	res, totalRecords, err := s.merchantService.FindByActive(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
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

	res, totalRecords, err := s.merchantService.FindByTrashed(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
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
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to create merchant: " + err.Error(),
		})
	}

	merchant, err := s.merchantService.CreateMerchant(&request)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to create merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully created merchant", merchant)

	return so, nil

}

func (s *merchantHandleGrpc) UpdateMerchant(ctx context.Context, req *pb.UpdateMerchantRequest) (*pb.ApiResponseMerchant, error) {
	request := requests.UpdateMerchantRequest{
		MerchantID: int(req.GetMerchantId()),
		Name:       req.GetName(),
		UserID:     int(req.GetUserId()),
		Status:     req.GetStatus(),
	}

	if err := request.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to update merchant: " + err.Error(),
		})
	}

	merchant, err := s.merchantService.UpdateMerchant(&request)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to update merchant: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully updated merchant", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) TrashedMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	if req.GetMerchantId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to trashed merchant: merchant id is required",
		})
	}

	merchant, err := s.merchantService.TrashedMerchant(int(req.GetMerchantId()))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully trashed merchant", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) RestoreMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	merchant, err := s.merchantService.RestoreMerchant(int(req.GetMerchantId()))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully restored merchant", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) DeleteMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	_, err := s.merchantService.DeleteMerchantPermanent(int(req.GetMerchantId()))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant")

	return so, nil
}

func (s *merchantHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantService.RestoreAllMerchant()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all merchant: ",
		})
	}

	so := s.mapping.ToProtoResponseMerchantAll("success", "Successfully restore all merchant")

	return so, nil
}

func (s *merchantHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantService.DeleteAllMerchantPermanent()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete merchant permanent: ",
		})
	}

	so := s.mapping.ToProtoResponseMerchantAll("success", "Successfully delete all merchant")

	return so, nil
}
