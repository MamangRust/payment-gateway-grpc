package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/service"
	"MamangRust/paymentgatewaygrpc/internal/repository"
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

func (s *merchantService) FindAll(page int, pageSize int, search string) ([]*response.MerchantResponse, int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantRepository.FindAllMerchants(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch merchant records", zap.Error(err))
		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch merchant records",
		}
	}

	merchantResponses := s.mapping.ToMerchantsResponse(merchants)

	s.logger.Debug("Successfully all merchant records",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantService) FindAllTransactions(page int, pageSize int, search string) ([]*response.MerchantTransactionResponse, int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantRepository.FindAllTransactions(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch merchant records", zap.Error(err))
		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch merchant records",
		}
	}

	merchantResponses := s.mapping.ToMerchantsTransactionResponse(merchants)

	s.logger.Debug("Successfully all merchant records",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantService) FindById(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Finding merchant by ID", zap.Int("merchant_id", merchant_id))

	res, err := s.merchantRepository.FindById(merchant_id)
	if err != nil {
		s.logger.Error("Failed to find merchant by ID", zap.Error(err), zap.Int("merchant_id", merchant_id))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found",
		}
	}

	so := s.mapping.ToMerchantResponse(res)

	s.logger.Debug("Successfully found merchant by ID", zap.Int("merchant_id", merchant_id))

	return so, nil
}

func (s *merchantService) FindMonthlyPaymentMethodsMerchant(year int) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	s.logger.Debug("Finding monthly payment methods for merchant", zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyPaymentMethodsMerchant(year)

	if err != nil {
		s.logger.Error("Failed to find monthly payment methods for merchant", zap.Error(err), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly payment methods for merchant",
		}
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

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly payment methods for merchant",
		}
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

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly amount for merchant",
		}
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

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly amount for merchant",
		}
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

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly amount for merchant",
		}
	}

	s.logger.Debug("Example", zap.Any("response month", res))

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

	s.logger.Debug("Successfully found monthly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyTotalAmountMerchant(year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding yearly amount for merchant", zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyTotalAmountMerchant(year)
	if err != nil {
		s.logger.Error("Failed to find yearly amount for merchant", zap.Error(err), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly amount for merchant",
		}
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.logger.Debug("Successfully found yearly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindAllTransactionsByMerchant(merchant_id int, page int, pageSize int, search string) ([]*response.MerchantTransactionResponse, int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantRepository.FindAllTransactionsByMerchant(merchant_id, search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch merchant records", zap.Error(err))
		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch merchant records",
		}
	}

	merchantResponses := s.mapping.ToMerchantsTransactionResponse(merchants)

	s.logger.Debug("Successfully all merchant records",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantService) FindMonthlyPaymentMethodByMerchants(merchantID int, year int) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	s.logger.Debug("Finding monthly payment methods by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyPaymentMethodByMerchants(merchantID, year)
	if err != nil {
		s.logger.Error("Failed to find monthly payment methods by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly payment methods by merchant",
		}
	}

	so := s.mapping.ToMerchantMonthlyPaymentMethods(res)

	s.logger.Debug("Successfully found monthly payment methods by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyPaymentMethodByMerchants(merchantID int, year int) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	s.logger.Debug("Finding yearly payment methods by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyPaymentMethodByMerchants(merchantID, year)
	if err != nil {
		s.logger.Error("Failed to find yearly payment methods by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly payment methods by merchant",
		}
	}

	so := s.mapping.ToMerchantYearlyPaymentMethods(res)

	s.logger.Debug("Successfully found yearly payment methods by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindMonthlyAmountByMerchants(merchantID int, year int) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding monthly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyAmountByMerchants(merchantID, year)
	if err != nil {
		s.logger.Error("Failed to find monthly amount by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly amount by merchant",
		}
	}

	so := s.mapping.ToMerchantMonthlyAmounts(res)

	s.logger.Debug("Successfully found monthly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyAmountByMerchants(merchantID int, year int) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding yearly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyAmountByMerchants(merchantID, year)
	if err != nil {
		s.logger.Error("Failed to find yearly amount by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly amount by merchant",
		}
	}

	so := s.mapping.ToMerchantYearlyAmounts(res)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindMonthlyTotalAmountByMerchants(merchantID int, year int) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding monthly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyTotalAmountByMerchants(merchantID, year)
	if err != nil {
		s.logger.Error("Failed to find monthly amount by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly amount by merchant",
		}
	}

	s.logger.Debug("Example", zap.Any("response month", res))

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

	s.logger.Debug("Successfully found monthly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyTotalAmountByMerchants(merchantID int, year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding yearly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyTotalAmountByMerchants(merchantID, year)
	if err != nil {
		s.logger.Error("Failed to find yearly amount by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly amount by merchant",
		}
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindAllTransactionsByApikey(api_key string, page int, pageSize int, search string) ([]*response.MerchantTransactionResponse, int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantRepository.FindAllTransactionsByApikey(api_key, search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch merchant records", zap.Error(err))
		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch merchant records",
		}
	}

	merchantResponses := s.mapping.ToMerchantsTransactionResponse(merchants)

	s.logger.Debug("Successfully all merchant records",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantService) FindMonthlyPaymentMethodByApikeys(api_key string, year int) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	s.logger.Debug("Finding monthly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyPaymentMethodByApikey(api_key, year)
	if err != nil {
		s.logger.Error("Failed to find monthly payment methods by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly payment methods by merchant",
		}
	}

	so := s.mapping.ToMerchantMonthlyPaymentMethods(res)

	s.logger.Debug("Successfully found monthly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyPaymentMethodByApikeys(api_key string, year int) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	s.logger.Debug("Finding yearly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyPaymentMethodByApikey(api_key, year)
	if err != nil {
		s.logger.Error("Failed to find yearly payment methods by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly payment methods by merchant",
		}
	}

	so := s.mapping.ToMerchantYearlyPaymentMethods(res)

	s.logger.Debug("Successfully found yearly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindMonthlyAmountByApikeys(api_key string, year int) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyAmountByApikey(api_key, year)
	if err != nil {
		s.logger.Error("Failed to find monthly amount by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly amount by merchant",
		}
	}

	so := s.mapping.ToMerchantMonthlyAmounts(res)

	s.logger.Debug("Successfully found monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyAmountByApikeys(api_key string, year int) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyAmountByApikey(api_key, year)
	if err != nil {
		s.logger.Error("Failed to find yearly amount by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly amount by merchant",
		}
	}

	so := s.mapping.ToMerchantYearlyAmounts(res)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindMonthlyTotalAmountByApikeys(api_key string, year int) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantRepository.GetMonthlyTotalAmountByApikey(api_key, year)
	if err != nil {
		s.logger.Error("Failed to find monthly amount by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find monthly amount by merchant",
		}
	}

	s.logger.Debug("Example", zap.Any("response month", res))

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

	s.logger.Debug("Successfully found monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantService) FindYearlyTotalAmountByApikeys(api_key string, year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	s.logger.Debug("Finding yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantRepository.GetYearlyTotalAmountByApikey(api_key, year)
	if err != nil {
		s.logger.Error("Failed to find yearly amount by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find yearly amount by merchant",
		}
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

//

func (s *merchantService) FindByActive(page int, pageSize int, search string) ([]*response.MerchantResponseDeleteAt, int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantRepository.FindByActive(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch active merchants", zap.Error(err))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch active merchants",
		}
	}

	so := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.logger.Debug("Successfully fetched active merchants",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *merchantService) FindByTrashed(page int, pageSize int, search string) ([]*response.MerchantResponseDeleteAt, int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantRepository.FindByTrashed(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch trashed merchants",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch trashed merchants",
		}
	}

	so := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.logger.Debug("Successfully fetched trashed merchants",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *merchantService) FindByApiKey(api_key string) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Finding merchant by API key", zap.String("api_key", api_key))

	res, err := s.merchantRepository.FindByApiKey(api_key)
	if err != nil {
		s.logger.Error("Failed to find merchant by API key", zap.Error(err), zap.String("api_key", api_key))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found by API key",
		}
	}

	so := s.mapping.ToMerchantResponse(res)

	s.logger.Debug("Successfully found merchant by API key", zap.String("api_key", api_key))

	return so, nil
}

func (s *merchantService) FindByMerchantUserId(user_id int) ([]*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Finding merchant by user ID", zap.Int("user_id", user_id))

	res, err := s.merchantRepository.FindByMerchantUserId(user_id)
	if err != nil {
		s.logger.Error("Failed to find merchant by user ID", zap.Error(err), zap.Int("user_id", user_id))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found by user ID",
		}
	}

	so := s.mapping.ToMerchantsResponse(res)

	s.logger.Debug("Successfully found merchant by user ID", zap.Int("user_id", user_id))

	return so, nil
}

func (s *merchantService) CreateMerchant(request *requests.CreateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new merchant", zap.String("merchant_name", request.Name))

	res, err := s.merchantRepository.CreateMerchant(request)

	if err != nil {
		s.logger.Error("Failed to create merchant", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create merchant",
		}
	}

	so := s.mapping.ToMerchantResponse(res)

	s.logger.Debug("Successfully created merchant", zap.Int("merchant_id", res.ID))

	return so, nil
}

func (s *merchantService) UpdateMerchant(request *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating merchant", zap.Int("merchant_id", request.MerchantID))

	_, err := s.merchantRepository.FindById(request.MerchantID)

	if err != nil {
		s.logger.Error("Merchant not found for update", zap.Error(err), zap.Int("merchant_id", request.MerchantID))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found",
		}
	}

	res, err := s.merchantRepository.UpdateMerchant(request)

	if err != nil {
		s.logger.Error("Failed to update merchant", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update merchant",
		}
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
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to trash merchant",
		}
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
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore merchant",
		}
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
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete merchant permanently",
		}
	}

	s.logger.Debug("Successfully deleted merchant permanently", zap.Int("merchant_id", merchant_id))

	return true, nil
}

func (s *merchantService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all merchants")

	_, err := s.merchantRepository.RestoreAllMerchant()

	if err != nil {
		s.logger.Error("Failed to restore all merchants", zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all merchants: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully restored all merchants")
	return true, nil
}

func (s *merchantService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all merchants")

	_, err := s.merchantRepository.DeleteAllMerchantPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all merchants", zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all merchants: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully deleted all merchants permanently")
	return true, nil
}
