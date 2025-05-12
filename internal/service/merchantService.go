package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/service"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/pkg/errors/merchant_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"

	"go.uber.org/zap"
)

type merchantService struct {
	merchantRepository repository.MerchantRepository
	logger             logger.LoggerInterface
	mapping            responseservice.MerchantResponseMapper
}

func NewMerchantService(
	merchantRepository repository.MerchantRepository,
	logger logger.LoggerInterface,
	mapping responseservice.MerchantResponseMapper,
) *merchantService {
	return &merchantService{
		merchantRepository: merchantRepository,
		logger:             logger,
		mapping:            mapping,
	}
}

func (s *merchantService) FindAll(req *requests.FindAllMerchants) ([]*response.MerchantResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all merchant records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantRepository.FindAllMerchants(req)

	if err != nil {
		s.logger.Error("Failed to fetch merchants",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, merchant_errors.ErrFailedFindAllMerchants
	}

	merchantResponses := s.mapping.ToMerchantsResponse(merchants)

	s.logger.Debug("Successfully all merchant records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantService) FindAllTransactions(req *requests.FindAllMerchantTransactions) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all merchant records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantRepository.FindAllTransactions(req)

	if err != nil {
		s.logger.Error("Failed to fetch merchant records", zap.Error(err))
		return nil, nil, merchant_errors.ErrFailedFindAllTransactions
	}

	merchantResponses := s.mapping.ToMerchantsTransactionResponse(merchants)

	s.logger.Debug("Successfully all merchant records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantService) FindById(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Finding merchant by ID", zap.Int("merchant_id", merchant_id))

	res, err := s.merchantRepository.FindById(merchant_id)

	if err != nil {
		s.logger.Error("Failed to retrieve merchant details",
			zap.Error(err),
			zap.Int("merchant_id", merchant_id))

		return nil, merchant_errors.ErrMerchantNotFoundRes
	}

	so := s.mapping.ToMerchantResponse(res)

	return so, nil
}

func (s *merchantService) FindMonthlyPaymentMethodsMerchant(year int) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	s.logger.Debug("Finding monthly payment methods for merchant", zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyPaymentMethodsMerchant(year)

	if err != nil {
		s.logger.Error("Failed to find monthly payment methods for merchant", zap.Error(err), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindMonthlyPaymentMethodsMerchant
	}

	so := s.mapping.ToMerchantMonthlyPaymentMethods(res)

	s.logger.Debug("Successfully found monthly payment methods for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyPaymentMethodMerchant(year int) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	s.logger.Debug("Finding yearly payment methods for merchant", zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyPaymentMethodMerchant(year)

	if err != nil {
		s.logger.Error("Failed to find yearly payment methods for merchant", zap.Error(err), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindYearlyPaymentMethodMerchant
	}

	so := s.mapping.ToMerchantYearlyPaymentMethods(res)

	s.logger.Debug("Successfully found yearly payment methods for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindMonthlyAmountMerchant(year int) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding monthly amount for merchant", zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyAmountMerchant(year)

	if err != nil {
		s.logger.Error("Failed to find monthly amount for merchant", zap.Error(err), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindMonthlyAmountMerchant
	}

	so := s.mapping.ToMerchantMonthlyAmounts(res)

	s.logger.Debug("Successfully found monthly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyAmountMerchant(year int) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding yearly amount for merchant", zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyAmountMerchant(year)

	if err != nil {
		s.logger.Error("Failed to find yearly amount for merchant", zap.Error(err), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindYearlyAmountMerchant
	}

	so := s.mapping.ToMerchantYearlyAmounts(res)

	s.logger.Debug("Successfully found yearly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindMonthlyTotalAmountMerchant(year int) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding monthly amount for merchant", zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyTotalAmountMerchant(year)

	if err != nil {
		s.logger.Error("Failed to find monthly amount for merchant", zap.Error(err), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindMonthlyTotalAmountMerchant
	}

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

	s.logger.Debug("Successfully found monthly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyTotalAmountMerchant(year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding yearly amount for merchant", zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyTotalAmountMerchant(year)

	if err != nil {
		s.logger.Error("Failed to find yearly amount for merchant", zap.Error(err), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindYearlyTotalAmountMerchant
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.logger.Debug("Successfully found yearly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindAllTransactionsByMerchant(req *requests.FindAllMerchantTransactionsById) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all merchant records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantRepository.FindAllTransactionsByMerchant(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active merchant",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, merchant_errors.ErrFailedFindAllTransactionsByMerchant
	}

	merchantResponses := s.mapping.ToMerchantsTransactionResponse(merchants)

	s.logger.Debug("Successfully fetched active merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantService) FindMonthlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	s.logger.Debug("Finding monthly payment methods by merchant", zap.Int("merchant_id", merchantID), zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyPaymentMethodByMerchants(req)

	if err != nil {
		s.logger.Error("Failed to find monthly payment methods by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindMonthlyPaymentMethodByMerchants
	}

	so := s.mapping.ToMerchantMonthlyPaymentMethods(res)

	s.logger.Debug("Successfully found monthly payment methods by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	s.logger.Debug("Finding yearly payment methods by merchant", zap.Int("merchant_id", merchantID), zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyPaymentMethodByMerchants(req)

	if err != nil {
		s.logger.Error("Failed to find yearly payment methods by merchant", zap.Error(err), zap.Int("merchant_id", merchantID), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindYearlyPaymentMethodByMerchants
	}

	so := s.mapping.ToMerchantYearlyPaymentMethods(res)

	s.logger.Debug("Successfully found yearly payment methods by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindMonthlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	s.logger.Debug("Finding monthly amount by merchant", zap.Int("merchant_id", merchantID), zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyAmountByMerchants(req)
	if err != nil {
		s.logger.Error("Failed to find monthly amount by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindMonthlyAmountByMerchants
	}

	so := s.mapping.ToMerchantMonthlyAmounts(res)

	s.logger.Debug("Successfully found monthly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	s.logger.Debug("Finding yearly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyAmountByMerchants(req)

	if err != nil {
		s.logger.Error("Failed to find yearly amount by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindYearlyAmountByMerchants
	}

	so := s.mapping.ToMerchantYearlyAmounts(res)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindMonthlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	s.logger.Debug("Finding monthly total amount by merchant", zap.Int("merchant_id", merchantID), zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyTotalAmountByMerchants(req)

	if err != nil {
		s.logger.Error("Failed to find monthly total amount by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindMonthlyTotalAmountByMerchants
	}

	s.logger.Debug("Example", zap.Any("response month", res))

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

	s.logger.Debug("Successfully found monthly total amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	s.logger.Debug("Finding yearly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyTotalAmountByMerchants(req)
	if err != nil {
		s.logger.Error("Failed to find yearly amount by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindYearlyTotalAmountByMerchants
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindAllTransactionsByApikey(req *requests.FindAllMerchantTransactionsByApiKey) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all transaction merchant records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantRepository.FindAllTransactionsByApikey(req)

	if err != nil {
		s.logger.Error("Failed to retrieve transaction merchant",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, merchant_errors.ErrFailedFindAllTransactionsByApikey
	}

	merchantResponses := s.mapping.ToMerchantsTransactionResponse(merchants)

	s.logger.Debug("Successfully all transaction merchant records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantService) FindMonthlyPaymentMethodByApikeys(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding monthly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyPaymentMethodByApikey(req)

	if err != nil {
		s.logger.Error("Failed to find monthly payment methods by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindMonthlyPaymentMethodByApikeys
	}

	so := s.mapping.ToMerchantMonthlyPaymentMethods(res)

	s.logger.Debug("Successfully found monthly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyPaymentMethodByApikeys(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding yearly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyPaymentMethodByApikey(req)

	if err != nil {
		s.logger.Error("Failed to find yearly payment methods by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindYearlyPaymentMethodByApikeys
	}

	so := s.mapping.ToMerchantYearlyPaymentMethods(res)

	s.logger.Debug("Successfully found yearly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindMonthlyAmountByApikeys(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyAmountByApikey(req)

	if err != nil {
		s.logger.Error("Failed to find monthly amount by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindMonthlyAmountByApikeys
	}

	so := s.mapping.ToMerchantMonthlyAmounts(res)

	s.logger.Debug("Successfully found monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyAmountByApikeys(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyAmountByApikey(req)

	if err != nil {
		s.logger.Error("Failed to find yearly amount by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindYearlyAmountByApikeys
	}

	so := s.mapping.ToMerchantYearlyAmounts(res)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindMonthlyTotalAmountByApikeys(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyTotalAmountByApikey(req)

	if err != nil {
		s.logger.Error("Failed to find monthly amount by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindMonthlyTotalAmountByApikeys
	}

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

	s.logger.Debug("Successfully found monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyTotalAmountByApikeys(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyTotalAmountByApikey(req)

	if err != nil {
		s.logger.Error("Failed to find yearly amount by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindYearlyTotalAmountByApikeys
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindByActive(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all merchant active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active cashiers",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, merchant_errors.ErrFailedFindActiveMerchants
	}

	so := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.logger.Debug("Successfully fetched active merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *merchantService) FindByTrashed(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching fetched trashed merchants",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to fetch trashed merchants",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, merchant_errors.ErrFailedFindTrashedMerchants
	}

	so := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.logger.Debug("Successfully fetched trashed merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *merchantService) FindByApiKey(api_key string) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Finding merchant by API key", zap.String("api_key", api_key))

	res, err := s.merchantRepository.FindByApiKey(api_key)

	if err != nil {
		s.logger.Error("Failed to retrieve merchant by api_key",
			zap.Error(err),
			zap.String("api_key", api_key))

		return nil, merchant_errors.ErrMerchantNotFoundRes
	}

	so := s.mapping.ToMerchantResponse(res)

	s.logger.Debug("Successfully found merchant by API key", zap.String("api_key", api_key))

	return so, nil
}

func (s *merchantService) FindByMerchantUserId(user_id int) ([]*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Finding merchant by user ID", zap.Int("user_id", user_id))

	res, err := s.merchantRepository.FindByMerchantUserId(user_id)

	if err != nil {
		s.logger.Error("Failed to retrieve merchant by user_id",
			zap.Error(err),
			zap.Int("user_id", user_id))

		return nil, merchant_errors.ErrMerchantNotFoundRes
	}

	so := s.mapping.ToMerchantsResponse(res)

	s.logger.Debug("Successfully found merchant by user ID", zap.Int("user_id", user_id))

	return so, nil
}

func (s *merchantService) CreateMerchant(request *requests.CreateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new merchant", zap.String("merchant_name", request.Name))

	res, err := s.merchantRepository.CreateMerchant(request)

	if err != nil {
		s.logger.Error("Failed to create new merchant",
			zap.Error(err),
			zap.Any("request", request))

		return nil, merchant_errors.ErrFailedCreateMerchant
	}

	so := s.mapping.ToMerchantResponse(res)

	s.logger.Debug("Successfully created merchant", zap.Int("merchant_id", res.ID))

	return so, nil
}

func (s *merchantService) UpdateMerchant(request *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating merchant", zap.Int("merchant_id", *request.MerchantID))

	res, err := s.merchantRepository.UpdateMerchant(request)

	if err != nil {
		s.logger.Error("Failed to update merchant",
			zap.Error(err),
			zap.Int("merchant", *request.MerchantID))

		return nil, merchant_errors.ErrFailedUpdateMerchant
	}

	so := s.mapping.ToMerchantResponse(res)

	s.logger.Debug("Successfully updated merchant", zap.Int("merchant_id", res.ID))

	return so, nil
}

func (s *merchantService) TrashedMerchant(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Trashing merchant", zap.Int("merchant_id", merchant_id))

	res, err := s.merchantRepository.TrashedMerchant(merchant_id)

	if err != nil {
		s.logger.Error("Failed to trash merchant", zap.Error(err), zap.Int("merchant_id", merchant_id))
		return nil, merchant_errors.ErrFailedTrashMerchant
	}

	s.logger.Debug("Successfully trashed merchant", zap.Int("merchant_id", merchant_id))

	so := s.mapping.ToMerchantResponse(res)

	return so, nil
}

func (s *merchantService) RestoreMerchant(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Restoring merchant", zap.Int("merchant_id", merchant_id))

	res, err := s.merchantRepository.RestoreMerchant(merchant_id)

	if err != nil {
		s.logger.Error("Failed to restore merchant", zap.Error(err), zap.Int("merchant_id", merchant_id))
		return nil, merchant_errors.ErrFailedRestoreMerchant
	}
	s.logger.Debug("Successfully restored merchant", zap.Int("merchant_id", merchant_id))

	so := s.mapping.ToMerchantResponse(res)

	return so, nil
}

func (s *merchantService) DeleteMerchantPermanent(merchant_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting merchant permanently", zap.Int("merchant_id", merchant_id))

	_, err := s.merchantRepository.DeleteMerchantPermanent(merchant_id)

	if err != nil {
		s.logger.Error("Failed to delete merchant permanently", zap.Error(err), zap.Int("merchant_id", merchant_id))

		return false, merchant_errors.ErrFailedDeleteMerchant
	}

	s.logger.Debug("Successfully deleted merchant permanently", zap.Int("merchant_id", merchant_id))

	return true, nil
}

func (s *merchantService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all merchants")

	_, err := s.merchantRepository.RestoreAllMerchant()

	if err != nil {
		s.logger.Error("Failed to restore all merchants", zap.Error(err))

		return false, merchant_errors.ErrFailedRestoreAllMerchants
	}

	s.logger.Debug("Successfully restored all merchants")
	return true, nil
}

func (s *merchantService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all merchants")

	_, err := s.merchantRepository.DeleteAllMerchantPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all merchants", zap.Error(err))

		return false, merchant_errors.ErrFailedDeleteAllMerchants
	}

	s.logger.Debug("Successfully deleted all merchants permanently")
	return true, nil
}
