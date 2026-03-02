package service

import (
	merchant_cache "MamangRust/paymentgatewaygrpc/internal/cache/merchant"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/errorhandler"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/merchant_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/user_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantService struct {
	userRepository     repository.UserRepository
	merchantRepository repository.MerchantRepository
	logger             logger.LoggerInterface
	observability      observability.TraceLoggerObservability
	cache              merchant_cache.MerchantMencache
}

type MerchantServiceDeps struct {
	UserRepo      repository.UserRepository
	MerchantRepo  repository.MerchantRepository
	Cache         merchant_cache.MerchantMencache
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewMerchantService(deps MerchantServiceDeps) *merchantService {
	return &merchantService{
		userRepository:     deps.UserRepo,
		merchantRepository: deps.MerchantRepo,
		logger:             deps.Logger,
		observability:      deps.Observability,
		cache:              deps.Cache,
	}
}

func (s *merchantService) FindAll(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsRow, *int, error) {
	const method = "FindAll"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchants(ctx, req); found {
		logSuccess("Successfully retrieved all merchant records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	merchants, err := s.merchantRepository.FindAllMerchants(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsRow](
			s.logger,
			merchant_errors.ErrFailedFindAllMerchants,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedMerchants(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched all merchant records",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantService) FindById(ctx context.Context, merchant_id int) (*db.GetMerchantByIDRow, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", merchant_id))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetCachedMerchant(ctx, merchant_id); found {
		logSuccess("Successfully retrieved merchant from cache", zap.Int("merchant.id", merchant_id))
		return cachedMerchant, nil
	}

	res, err := s.merchantRepository.FindById(ctx, merchant_id)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantByIDRow](
			s.logger,
			merchant_errors.ErrMerchantNotFoundRes,
			method,
			span,

			zap.Int("merchant_id", merchant_id),
		)
	}

	s.cache.SetCachedMerchant(ctx, res)

	logSuccess("Successfully found merchant by ID", zap.Int("merchant_id", merchant_id))

	return res, nil
}

func (s *merchantService) FindByApiKey(ctx context.Context, api_key string) (*db.GetMerchantByApiKeyRow, error) {
	const method = "FindByApiKey"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetCachedMerchantByApiKey(ctx, api_key); found {
		logSuccess("Successfully found merchant by API key from cache", zap.String("api_key", api_key))
		return cachedMerchant, nil
	}

	res, err := s.merchantRepository.FindByApiKey(ctx, api_key)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantByApiKeyRow](
			s.logger,
			merchant_errors.ErrMerchantNotFoundRes,
			method,
			span,
		)
	}

	s.cache.SetCachedMerchantByApiKey(ctx, api_key, res)

	logSuccess("Successfully found merchant by Api key")

	return res, nil
}

func (s *merchantService) FindByMerchantUserId(ctx context.Context, user_id int) ([]*db.GetMerchantsByUserIDRow, error) {
	const method = "FindByMerchantUserId"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	if cachedMerchants, found := s.cache.GetCachedMerchantsByUserId(ctx, user_id); found {
		logSuccess("Successfully found merchants by user ID from cache", zap.Int("user.id", user_id), zap.Int("count", len(cachedMerchants)))

		return cachedMerchants, nil
	}

	res, err := s.merchantRepository.FindByMerchantUserId(ctx, user_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMerchantsByUserIDRow](
			s.logger,
			merchant_errors.ErrMerchantNotFoundRes,
			method,
			span,
		)
	}

	s.cache.SetCachedMerchantsByUserId(ctx, user_id, res)

	logSuccess("Successfully found merchant by Api key")

	return res, nil
}

func (s *merchantService) FindAllTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*db.FindAllTransactionsRow, *int, error) {
	const method = "FindAllTransactions"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCacheAllMerchantTransactions(ctx, req); found {
		logSuccess("Successfully retrieved all merchant records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	merchants, err := s.merchantRepository.FindAllTransactions(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.FindAllTransactionsRow](
			s.logger,
			merchant_errors.ErrFailedFindAllTransactions,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCacheAllMerchantTransactions(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched all merchant transaction records",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantService) FindAllTransactionsByMerchant(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*db.FindAllTransactionsByMerchantRow, *int, error) {
	const method = "FindAllTransactionsByMerchant"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.Int("merchant_id", req.MerchantID))

	defer func() {
		end(status)
	}()

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	if data, total, found := s.cache.GetCacheMerchantTransactions(ctx, req); found {
		logSuccess("Successfully retrieved all merchant records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	merchants, err := s.merchantRepository.FindAllTransactionsByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.FindAllTransactionsByMerchantRow](
			s.logger,
			merchant_errors.ErrFailedFindAllTransactionsByMerchant,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.Int("merchant_id", req.MerchantID),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCacheMerchantTransactions(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched merchant transactions",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.Int("merchant_id", req.MerchantID))

	return merchants, &totalCount, nil
}

func (s *merchantService) FindAllTransactionsByApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*db.FindAllTransactionsByApikeyRow, *int, error) {
	const method = "FindAllTransactionsByApikey"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.String("api_key", req.ApiKey))

	defer func() {
		end(status)
	}()

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	if data, total, found := s.cache.GetCacheMerchantTransactionApikey(ctx, req); found {
		logSuccess("Successfully retrieved all merchant records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	merchants, err := s.merchantRepository.FindAllTransactionsByApikey(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.FindAllTransactionsByApikeyRow](
			s.logger,
			merchant_errors.ErrFailedFindAllTransactionsByApikey,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("api_key", req.ApiKey),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCacheMerchantTransactionApikey(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched merchant transactions by API key",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("api_key", req.ApiKey))

	return merchants, &totalCount, nil
}

func (s *merchantService) FindMonthlyPaymentMethodsMerchant(ctx context.Context, year int) ([]*db.GetMonthlyPaymentMethodsMerchantRow, error) {
	const method = "FindMonthlyPaymentMethodsMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetMonthlyPaymentMethodsMerchantCache(ctx, year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	s.logger.Debug("Cache miss for monthly payment methods for merchant, fetching from DB", zap.Int("year", year))

	dbRows, err := s.merchantRepository.GetMonthlyPaymentMethodsMerchant(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyPaymentMethodsMerchantRow](
			s.logger,
			merchant_errors.ErrFailedFindMonthlyPaymentMethodsMerchant,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyPaymentMethodsMerchantCache(ctx, year, dbRows)

	logSuccess("Successfully found monthly payment methods for merchant (from DB)", zap.Int("year", year))
	return dbRows, nil
}

func (s *merchantService) FindYearlyPaymentMethodMerchant(ctx context.Context, year int) ([]*db.GetYearlyPaymentMethodMerchantRow, error) {
	const method = "FindYearlyPaymentMethodMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetYearlyPaymentMethodMerchantCache(ctx, year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	s.logger.Debug("Cache miss for yearly payment methods for merchant, fetching from DB", zap.Int("year", year))

	dbRows, err := s.merchantRepository.GetYearlyPaymentMethodMerchant(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyPaymentMethodMerchantRow](
			s.logger,
			merchant_errors.ErrFailedFindYearlyPaymentMethodMerchant,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyPaymentMethodMerchantCache(ctx, year, dbRows)

	logSuccess("Successfully found yearly payment methods for merchant (from DB)", zap.Int("year", year))
	return dbRows, nil
}

func (s *merchantService) FindMonthlyAmountMerchant(ctx context.Context, year int) ([]*db.GetMonthlyAmountMerchantRow, error) {
	const method = "FindMonthlyAmountMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetMonthlyAmountMerchantCache(ctx, year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	s.logger.Debug("Cache miss for monthly amount for merchant, fetching from DB", zap.Int("year", year))

	dbRows, err := s.merchantRepository.GetMonthlyAmountMerchant(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountMerchantRow](
			s.logger,
			merchant_errors.ErrFailedFindMonthlyAmountMerchant,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyAmountMerchantCache(ctx, year, dbRows)

	logSuccess("Successfully found monthly amount for merchant (from DB)", zap.Int("year", year))
	return dbRows, nil
}

func (s *merchantService) FindYearlyAmountMerchant(ctx context.Context, year int) ([]*db.GetYearlyAmountMerchantRow, error) {
	const method = "FindYearlyAmountMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetYearlyAmountMerchantCache(ctx, year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetYearlyAmountMerchant(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountMerchantRow](
			s.logger,
			merchant_errors.ErrFailedFindYearlyAmountMerchant,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyAmountMerchantCache(ctx, year, dbRows)

	logSuccess("Successfully found yearly amount for merchant (from DB)", zap.Int("year", year))
	return dbRows, nil
}

func (s *merchantService) FindMonthlyTotalAmountMerchant(ctx context.Context, year int) ([]*db.GetMonthlyTotalAmountMerchantRow, error) {
	const method = "FindMonthlyTotalAmountMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetMonthlyTotalAmountMerchantCache(ctx, year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	s.logger.Debug("Cache miss for monthly total amount for merchant, fetching from DB", zap.Int("year", year))

	dbRows, err := s.merchantRepository.GetMonthlyTotalAmountMerchant(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTotalAmountMerchantRow](
			s.logger,
			merchant_errors.ErrFailedFindMonthlyTotalAmountMerchant,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyTotalAmountMerchantCache(ctx, year, dbRows)

	logSuccess("Successfully found monthly total amount for merchant (from DB)", zap.Int("year", year))
	return dbRows, nil
}

func (s *merchantService) FindYearlyTotalAmountMerchant(ctx context.Context, year int) ([]*db.GetYearlyTotalAmountMerchantRow, error) {
	const method = "FindYearlyTotalAmountMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetYearlyTotalAmountMerchantCache(ctx, year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetYearlyTotalAmountMerchant(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTotalAmountMerchantRow](
			s.logger,
			merchant_errors.ErrFailedFindYearlyTotalAmountMerchant,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyTotalAmountMerchantCache(ctx, year, dbRows)

	logSuccess("Successfully found yearly total amount for merchant (from DB)", zap.Int("year", year))
	return dbRows, nil
}

func (s *merchantService) FindMonthlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*db.GetMonthlyPaymentMethodByMerchantsRow, error) {
	const method = "FindMonthlyPaymentMethodByMerchants"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchant_id", req.MerchantID))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetMonthlyPaymentMethodByMerchantsCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant_id", req.MerchantID))
		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetMonthlyPaymentMethodByMerchants(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyPaymentMethodByMerchantsRow](
			s.logger,
			merchant_errors.ErrFailedFindMonthlyPaymentMethodByMerchants,
			method,
			span,

			zap.Int("merchant_id", req.MerchantID),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyPaymentMethodByMerchantsCache(ctx, req, dbRows)

	logSuccess("Successfully found monthly payment methods by merchant (from DB)", zap.Int("merchant_id", req.MerchantID), zap.Int("year", req.Year))
	return dbRows, nil
}

func (s *merchantService) FindYearlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*db.GetYearlyPaymentMethodByMerchantsRow, error) {
	const method = "FindYearlyPaymentMethodByMerchants"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchant_id", req.MerchantID))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetYearlyPaymentMethodByMerchantsCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant_id", req.MerchantID))
		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetYearlyPaymentMethodByMerchants(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyPaymentMethodByMerchantsRow](
			s.logger,
			merchant_errors.ErrFailedFindYearlyPaymentMethodByMerchants,
			method,
			span,

			zap.Int("merchant_id", req.MerchantID),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyPaymentMethodByMerchantsCache(ctx, req, dbRows)

	logSuccess("Successfully found yearly payment methods by merchant (from DB)", zap.Int("merchant_id", req.MerchantID), zap.Int("year", req.Year))
	return dbRows, nil
}

func (s *merchantService) FindMonthlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*db.GetMonthlyAmountByMerchantsRow, error) {
	const method = "FindMonthlyAmountByMerchants"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchant_id", req.MerchantID))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetMonthlyAmountByMerchantsCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", req.MerchantID))

		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetMonthlyAmountByMerchants(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountByMerchantsRow](
			s.logger,
			merchant_errors.ErrFailedFindMonthlyAmountByMerchants,
			method,
			span,

			zap.Int("merchant_id", req.MerchantID),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyAmountByMerchantsCache(ctx, req, dbRows)

	logSuccess("Successfully found monthly amount by merchant (from DB)", zap.Int("merchant_id", req.MerchantID), zap.Int("year", req.Year))
	return dbRows, nil
}

func (s *merchantService) FindYearlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*db.GetYearlyAmountByMerchantsRow, error) {
	const method = "FindYearlyAmountByMerchants"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchant_id", req.MerchantID))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetYearlyAmountByMerchantsCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", req.MerchantID))

		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetYearlyAmountByMerchants(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountByMerchantsRow](
			s.logger,
			merchant_errors.ErrFailedFindYearlyAmountByMerchants,
			method,
			span,

			zap.Int("merchant_id", req.MerchantID),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyAmountByMerchantsCache(ctx, req, dbRows)

	logSuccess("Successfully found yearly amount by merchant (from DB)", zap.Int("merchant_id", req.MerchantID), zap.Int("year", req.Year))
	return dbRows, nil
}

func (s *merchantService) FindMonthlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*db.GetMonthlyTotalAmountByMerchantRow, error) {
	const method = "FindMonthlyTotalAmountByMerchants"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchant_id", req.MerchantID))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetMonthlyTotalAmountByMerchantsCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", req.MerchantID), zap.Int("year", req.Year))
		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetMonthlyTotalAmountByMerchants(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTotalAmountByMerchantRow](
			s.logger,
			merchant_errors.ErrFailedFindMonthlyTotalAmountByMerchants,
			method,
			span,

			zap.Int("merchant_id", req.MerchantID),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyTotalAmountByMerchantsCache(ctx, req, dbRows)

	logSuccess("Successfully found monthly total amount by merchant (from DB)", zap.Int("merchant_id", req.MerchantID), zap.Int("year", req.Year))
	return dbRows, nil
}

func (s *merchantService) FindYearlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*db.GetYearlyTotalAmountByMerchantRow, error) {
	const method = "FindYearlyTotalAmountByMerchants"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchant_id", req.MerchantID))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetYearlyTotalAmountByMerchantsCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", req.MerchantID), zap.Int("year", req.Year))
		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetYearlyTotalAmountByMerchants(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTotalAmountByMerchantRow](
			s.logger,
			merchant_errors.ErrFailedFindYearlyTotalAmountByMerchants,
			method,
			span,

			zap.Int("merchant_id", req.MerchantID),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyTotalAmountByMerchantsCache(ctx, req, dbRows)

	logSuccess("Successfully found yearly total amount by merchant (from DB)", zap.Int("merchant_id", req.MerchantID), zap.Int("year", req.Year))
	return dbRows, nil
}

func (s *merchantService) FindMonthlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*db.GetMonthlyPaymentMethodByApikeyRow, error) {
	const method = "FindMonthlyPaymentMethodByApikeys"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("api_key", req.Apikey),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetMonthlyPaymentMethodByApikeysCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.String("api_key", req.Apikey), zap.Int("year", req.Year))

		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetMonthlyPaymentMethodByApikey(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyPaymentMethodByApikeyRow](
			s.logger,
			merchant_errors.ErrFailedFindMonthlyPaymentMethodByApikeys,
			method,
			span,

			zap.String("api_key", req.Apikey),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyPaymentMethodByApikeysCache(ctx, req, dbRows)

	logSuccess("Successfully found monthly payment methods by API key (from DB)", zap.String("api_key", req.Apikey), zap.Int("year", req.Year))
	return dbRows, nil
}

func (s *merchantService) FindYearlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*db.GetYearlyPaymentMethodByApikeyRow, error) {
	const method = "FindYearlyPaymentMethodByApikeys"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("api_key", req.Apikey),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetYearlyPaymentMethodByApikeysCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.String("api_key", req.Apikey), zap.Int("year", req.Year))
		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetYearlyPaymentMethodByApikey(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyPaymentMethodByApikeyRow](
			s.logger,
			merchant_errors.ErrFailedFindYearlyPaymentMethodByApikeys,
			method,
			span,

			zap.String("api_key", req.Apikey),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyPaymentMethodByApikeysCache(ctx, req, dbRows)

	logSuccess("Successfully found yearly payment methods by API key (from DB)", zap.String("api_key", req.Apikey), zap.Int("year", req.Year))
	return dbRows, nil
}

func (s *merchantService) FindMonthlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*db.GetMonthlyAmountByApikeyRow, error) {
	const method = "FindMonthlyAmountByApikeys"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("api_key", req.Apikey),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetMonthlyAmountByApikeysCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.String("api_key", req.Apikey), zap.Int("year", req.Year))

		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetMonthlyAmountByApikey(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountByApikeyRow](
			s.logger,
			merchant_errors.ErrFailedFindMonthlyAmountByApikeys,
			method,
			span,

			zap.String("api_key", req.Apikey),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyAmountByApikeysCache(ctx, req, dbRows)

	logSuccess("Successfully found monthly amount by API key (from DB)", zap.String("api_key", req.Apikey), zap.Int("year", req.Year))
	return dbRows, nil
}

func (s *merchantService) FindYearlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*db.GetYearlyAmountByApikeyRow, error) {
	const method = "FindYearlyAmountByApikeys"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("api_key", req.Apikey),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetYearlyAmountByApikeysCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.String("api_key", req.Apikey), zap.Int("year", req.Year))
		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetYearlyAmountByApikey(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountByApikeyRow](
			s.logger,
			merchant_errors.ErrFailedFindYearlyAmountByApikeys,
			method,
			span,

			zap.String("api_key", req.Apikey),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyAmountByApikeysCache(ctx, req, dbRows)

	logSuccess("Successfully found yearly amount by API key (from DB)", zap.String("api_key", req.Apikey), zap.Int("year", req.Year))
	return dbRows, nil
}

func (s *merchantService) FindMonthlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*db.GetMonthlyTotalAmountByApikeyRow, error) {
	const method = "FindMonthlyTotalAmountByApikeys"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("api_key", req.Apikey),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetMonthlyTotalAmountByApikeysCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.String("api_key", req.Apikey), zap.Int("year", req.Year))
		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetMonthlyTotalAmountByApikey(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTotalAmountByApikeyRow](
			s.logger,
			merchant_errors.ErrFailedFindMonthlyTotalAmountByApikeys,
			method,
			span,

			zap.String("api_key", req.Apikey),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyTotalAmountByApikeysCache(ctx, req, dbRows)

	logSuccess("Successfully found monthly total amount by API key (from DB)", zap.String("api_key", req.Apikey), zap.Int("year", req.Year))
	return dbRows, nil
}

func (s *merchantService) FindYearlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*db.GetYearlyTotalAmountByApikeyRow, error) {
	const method = "FindYearlyTotalAmountByApikeys"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("api_key", req.Apikey),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.cache.GetYearlyTotalAmountByApikeysCache(ctx, req); found {
		logSuccess("Successfully fetched merchant from cache", zap.String("api_key", req.Apikey), zap.Int("year", req.Year))
		return cachedMerchant, nil
	}

	dbRows, err := s.merchantRepository.GetYearlyTotalAmountByApikey(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTotalAmountByApikeyRow](
			s.logger,
			merchant_errors.ErrFailedFindYearlyTotalAmountByApikeys,
			method,
			span,

			zap.String("api_key", req.Apikey),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyTotalAmountByApikeysCache(ctx, req, dbRows)

	logSuccess("Successfully found yearly total amount by API key (from DB)", zap.String("api_key", req.Apikey), zap.Int("year", req.Year))
	return dbRows, nil
}

func (s *merchantService) FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetActiveMerchantsRow, *int, error) {
	const method = "FindByActive"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantActive(ctx, req); found {
		logSuccess("Successfully fetched active merchants from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetActiveMerchantsRow](
			s.logger,
			merchant_errors.ErrFailedFindActiveMerchants,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedMerchantActive(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched active merchants",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetTrashedMerchantsRow, *int, error) {
	const method = "FindByTrashed"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed merchants from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTrashedMerchantsRow](
			s.logger,
			merchant_errors.ErrFailedFindTrashedMerchants,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedMerchantTrashed(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched trashed merchants",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantService) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*db.CreateMerchantRow, error) {
	const method = "CreateMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("merchant_name", request.Name),
		attribute.Int("user_id", request.UserID))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Creating new merchant", zap.String("merchant_name", request.Name))

	_, err := s.userRepository.FindById(ctx, request.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,

			zap.Int("user_id", request.UserID),
		)
	}

	res, err := s.merchantRepository.CreateMerchant(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantRow](
			s.logger,
			merchant_errors.ErrFailedCreateMerchant,
			method,
			span,

			zap.Any("request", request),
		)
	}

	logSuccess("Successfully created merchant", zap.Int("merchant_id", int(res.MerchantID)))

	return res, nil
}

func (s *merchantService) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.UpdateMerchantRow, error) {
	const method = "UpdateMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", *request.MerchantID),
		attribute.Int("user_id", request.UserID))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Updating merchant", zap.Int("merchant_id", *request.MerchantID))

	_, err := s.userRepository.FindById(ctx, request.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,

			zap.Int("user_id", request.UserID),
		)
	}

	res, err := s.merchantRepository.UpdateMerchant(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantRow](
			s.logger,
			merchant_errors.ErrFailedUpdateMerchant,
			method,
			span,

			zap.Int("merchant_id", *request.MerchantID),
		)
	}

	s.cache.DeleteCachedMerchant(ctx, *request.MerchantID)

	logSuccess("Successfully updated merchant", zap.Int("merchant_id", int(res.MerchantID)))

	return res, nil
}

func (s *merchantService) TrashedMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	const method = "TrashedMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", merchant_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Trashing merchant", zap.Int("merchant_id", merchant_id))

	res, err := s.merchantRepository.TrashedMerchant(ctx, merchant_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedTrashMerchant,
			method,
			span,

			zap.Int("merchant_id", merchant_id),
		)
	}

	logSuccess("Successfully trashed merchant", zap.Int("merchant_id", merchant_id))

	return res, nil
}

func (s *merchantService) RestoreMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	const method = "RestoreMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", merchant_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring merchant", zap.Int("merchant_id", merchant_id))

	res, err := s.merchantRepository.RestoreMerchant(ctx, merchant_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedRestoreMerchant,
			method,
			span,

			zap.Int("merchant_id", merchant_id),
		)
	}

	logSuccess("Successfully restored merchant", zap.Int("merchant_id", merchant_id))

	return res, nil
}

func (s *merchantService) DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error) {
	const method = "DeleteMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", merchant_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Deleting merchant permanently", zap.Int("merchant_id", merchant_id))

	_, err := s.merchantRepository.DeleteMerchantPermanent(ctx, merchant_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_errors.ErrFailedDeleteMerchant,
			method,
			span,

			zap.Int("merchant_id", merchant_id),
		)
	}

	logSuccess("Successfully deleted merchant permanently", zap.Int("merchant_id", merchant_id))

	return true, nil
}

func (s *merchantService) RestoreAllMerchant(ctx context.Context) (bool, error) {
	const method = "RestoreAllMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring all merchants")

	_, err := s.merchantRepository.RestoreAllMerchant(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_errors.ErrFailedRestoreAllMerchants,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all merchants")
	return true, nil
}

func (s *merchantService) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Permanently deleting all merchants")

	_, err := s.merchantRepository.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_errors.ErrFailedDeleteAllMerchants,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all merchants permanently")
	return true, nil
}
