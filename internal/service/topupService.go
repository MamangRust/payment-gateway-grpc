package service

import (
	topup_cache "MamangRust/paymentgatewaygrpc/internal/cache/topup"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/errorhandler"
	"context"

	"MamangRust/paymentgatewaygrpc/internal/repository"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/card_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/topup_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type topupService struct {
	cardRepository  repository.CardRepository
	topupRepository repository.TopupRepository
	saldoRepository repository.SaldoRepository
	logger          logger.LoggerInterface
	observability   observability.TraceLoggerObservability
	cache           topup_cache.TopupMencach
}

type TopupServiceDeps struct {
	CardRepo  repository.CardRepository
	TopupRepo repository.TopupRepository
	SaldoRepo repository.SaldoRepository

	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
	Cache         topup_cache.TopupMencach
}

func NewTopupService(deps TopupServiceDeps) *topupService {
	return &topupService{
		cardRepository:  deps.CardRepo,
		topupRepository: deps.TopupRepo,
		saldoRepository: deps.SaldoRepo,
		logger:          deps.Logger,
		observability:   deps.Observability,
		cache:           deps.Cache,
	}
}

func (s *topupService) FindAll(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetTopupsRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedTopupsCache(ctx, req); found {
		logSuccess("Successfully retrieved all topup records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	topups, err := s.topupRepository.FindAllTopups(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTopupsRow](
			s.logger,
			topup_errors.ErrFailedFindAllTopups,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(topups) > 0 {
		totalCount = int(topups[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTopupsCache(ctx, req, topups, &totalCount)

	logSuccess("Successfully fetched topup",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return topups, &totalCount, nil
}

func (s *topupService) FindAllByCardNumber(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*db.GetTopupsByCardNumberRow, *int, error) {
	const method = "FindAllByCardNumber"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search
	card_number := req.CardNumber

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

	if data, total, found := s.cache.GetCacheTopupByCardCache(ctx, req); found {
		logSuccess("Successfully retrieved all topup records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	topups, err := s.topupRepository.FindAllTopupByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTopupsByCardNumberRow](
			s.logger,
			topup_errors.ErrFailedFindAllTopupsByCardNumber,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("card_number", card_number),
		)
	}

	var totalCount int

	if len(topups) > 0 {
		totalCount = int(topups[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCacheTopupByCardCache(ctx, req, topups, &totalCount)

	logSuccess("Successfully fetched topup by card number",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("card_number", card_number))

	return topups, &totalCount, nil
}

func (s *topupService) FindById(ctx context.Context, topupID int) (*db.GetTopupByIDRow, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("topup_id", topupID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedTopupCache(ctx, topupID); found {
		logSuccess("Successfully retrieved topup from cache", zap.Int("topup.id", topupID))
		return data, nil
	}

	topup, err := s.topupRepository.FindById(ctx, topupID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetTopupByIDRow](
			s.logger,
			topup_errors.ErrTopupNotFoundRes,
			method,
			span,

			zap.Int("topup_id", topupID),
		)
	}

	s.cache.SetCachedTopupCache(ctx, topup)

	logSuccess("Successfully fetched topup", zap.Int("topup_id", topupID))

	return topup, nil
}

func (s *topupService) FindByActive(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetActiveTopupsRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedTopupActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved all topup records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	topups, err := s.topupRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetActiveTopupsRow](
			s.logger,
			topup_errors.ErrFailedFindActiveTopups,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(topups) > 0 {
		totalCount = int(topups[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTopupActiveCache(ctx, req, topups, &totalCount)

	logSuccess("Successfully fetched active topup",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return topups, &totalCount, nil
}

func (s *topupService) FindByTrashed(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetTrashedTopupsRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedTopupTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved all topup records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	topups, err := s.topupRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTrashedTopupsRow](
			s.logger,
			topup_errors.ErrFailedFindTrashedTopups,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(topups) > 0 {
		totalCount = int(topups[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTopupTrashedCache(ctx, req, topups, &totalCount)

	logSuccess("Successfully fetched trashed topup",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return topups, &totalCount, nil
}

func (s *topupService) FindMonthTopupStatusSuccess(ctx context.Context, req *requests.MonthTopupStatus) ([]*db.GetMonthTopupStatusSuccessRow, error) {
	const method = "FindMonthTopupStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTopupStatusSuccessCache(ctx, req); found {
		logSuccess("Successfully fetched monthly topup status success from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return data, nil
	}

	s.logger.Debug("Cache miss for monthly topup status success, fetching from DB", zap.Int("year", req.Year), zap.Int("month", req.Month))

	dbRows, err := s.topupRepository.GetMonthTopupStatusSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthTopupStatusSuccessRow](
			s.logger,
			topup_errors.ErrFailedFindMonthTopupStatusSuccess,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetMonthTopupStatusSuccessCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly topup status success (from DB)", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *topupService) FindYearlyTopupStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTopupStatusSuccessRow, error) {
	const method = "FindYearlyTopupStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupStatusSuccessCache(ctx, year); found {
		logSuccess("Successfully fetched yearly topup status success from cache", zap.Int("year", year))
		return data, nil
	}

	dbRows, err := s.topupRepository.GetYearlyTopupStatusSuccess(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTopupStatusSuccessRow](
			s.logger,
			topup_errors.ErrFailedFindYearlyTopupStatusSuccess,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyTopupStatusSuccessCache(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly topup status success (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *topupService) FindMonthTopupStatusFailed(ctx context.Context, req *requests.MonthTopupStatus) ([]*db.GetMonthTopupStatusFailedRow, error) {
	const method = "FindMonthTopupStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTopupStatusFailedCache(ctx, req); found {
		logSuccess("Successfully fetched monthly topup status Failed from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return data, nil
	}

	dbRows, err := s.topupRepository.GetMonthTopupStatusFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthTopupStatusFailedRow](
			s.logger,
			topup_errors.ErrFailedFindMonthTopupStatusFailed,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetMonthTopupStatusFailedCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly topup status failed (from DB)", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *topupService) FindYearlyTopupStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTopupStatusFailedRow, error) {
	const method = "FindYearlyTopupStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupStatusFailedCache(ctx, year); found {
		logSuccess("Successfully fetched yearly topup status Failed from cache", zap.Int("year", year))
		return data, nil
	}

	dbRows, err := s.topupRepository.GetYearlyTopupStatusFailed(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTopupStatusFailedRow](
			s.logger,
			topup_errors.ErrFailedFindYearlyTopupStatusFailed,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyTopupStatusFailedCache(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly topup status failed (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *topupService) FindMonthlyTopupMethods(ctx context.Context, year int) ([]*db.GetMonthlyTopupMethodsRow, error) {
	const method = "FindMonthlyTopupMethods"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTopupMethodsCache(ctx, year); found {
		logSuccess("Successfully fetched monthly topup methods from cache", zap.Int("year", year))
		return data, nil
	}

	dbRows, err := s.topupRepository.GetMonthlyTopupMethods(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTopupMethodsRow](
			s.logger,
			topup_errors.ErrFailedFindMonthlyTopupMethods,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyTopupMethodsCache(ctx, year, dbRows)

	logSuccess("Successfully fetched monthly topup methods (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *topupService) FindYearlyTopupMethods(ctx context.Context, year int) ([]*db.GetYearlyTopupMethodsRow, error) {
	const method = "FindYearlyTopupMethods"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupMethodsCache(ctx, year); found {
		logSuccess("Successfully fetched yearly topup methods from cache", zap.Int("year", year))
		return data, nil
	}

	s.logger.Debug("Cache miss for yearly topup methods, fetching from DB", zap.Int("year", year))

	dbRows, err := s.topupRepository.GetYearlyTopupMethods(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTopupMethodsRow](
			s.logger,
			topup_errors.ErrFailedFindYearlyTopupMethods,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyTopupMethodsCache(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly topup methods (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *topupService) FindMonthlyTopupAmounts(ctx context.Context, year int) ([]*db.GetMonthlyTopupAmountsRow, error) {
	const method = "FindMonthlyTopupAmounts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTopupAmountsCache(ctx, year); found {
		logSuccess("Successfully fetched monthly topup amounts from cache", zap.Int("year", year))
		return data, nil
	}

	dbRows, err := s.topupRepository.GetMonthlyTopupAmounts(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTopupAmountsRow](
			s.logger,
			topup_errors.ErrFailedFindMonthlyTopupAmounts,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyTopupAmountsCache(ctx, year, dbRows)

	logSuccess("Successfully fetched monthly topup amounts (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *topupService) FindYearlyTopupAmounts(ctx context.Context, year int) ([]*db.GetYearlyTopupAmountsRow, error) {
	const method = "FindYearlyTopupAmounts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupAmountsCache(ctx, year); found {
		logSuccess("Successfully fetched yearly topup amounts from cache", zap.Int("year", year))
		return data, nil
	}

	dbRows, err := s.topupRepository.GetYearlyTopupAmounts(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTopupAmountsRow](
			s.logger,
			topup_errors.ErrFailedFindYearlyTopupAmounts,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyTopupAmountsCache(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly topup amounts (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *topupService) FindMonthTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*db.GetMonthTopupStatusSuccessCardNumberRow, error) {
	const method = "FindMonthTopupStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTopupStatusSuccessByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched monthly topup status success", zap.Int("year", req.Year), zap.Int("month", req.Month), zap.String("card_number", req.CardNumber))
		return data, nil
	}

	dbRows, err := s.topupRepository.GetMonthTopupStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthTopupStatusSuccessCardNumberRow](
			s.logger,
			topup_errors.ErrFailedFindMonthTopupStatusSuccessByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetMonthTopupStatusSuccessByCardNumberCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly topup status success by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *topupService) FindYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*db.GetYearlyTopupStatusSuccessCardNumberRow, error) {
	const method = "FindYearlyTopupStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupStatusSuccessByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched yearly topup status success", zap.Int("year", req.Year), zap.String("card_number", req.CardNumber))
		return data, nil
	}

	dbRows, err := s.topupRepository.GetYearlyTopupStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTopupStatusSuccessCardNumberRow](
			s.logger,
			topup_errors.ErrFailedFindYearlyTopupStatusSuccessByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.logger.Debug("Setting cache for yearly topup status success by card number",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	s.cache.SetYearlyTopupStatusSuccessByCardNumberCache(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly topup status success by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *topupService) FindMonthTopupStatusFailedByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*db.GetMonthTopupStatusFailedCardNumberRow, error) {
	const method = "FindMonthTopupStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTopupStatusFailedByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched monthly topup status Failed", zap.Int("year", req.Year), zap.Int("month", req.Month), zap.String("card_number", req.CardNumber))
		return data, nil
	}

	s.logger.Debug("Cache miss for monthly topup status failed by card number, fetching from DB",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	dbRows, err := s.topupRepository.GetMonthTopupStatusFailedByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthTopupStatusFailedCardNumberRow](
			s.logger,
			topup_errors.ErrFailedFindMonthTopupStatusFailedByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetMonthTopupStatusFailedByCardNumberCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly topup status failed by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *topupService) FindYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*db.GetYearlyTopupStatusFailedCardNumberRow, error) {
	const method = "FindYearlyTopupStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupStatusFailedByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched yearly topup status Failed", zap.Int("year", req.Year), zap.String("card_number", req.CardNumber))
		return data, nil
	}

	s.logger.Debug("Cache miss for yearly topup status failed by card number, fetching from DB",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	dbRows, err := s.topupRepository.GetYearlyTopupStatusFailedByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTopupStatusFailedCardNumberRow](
			s.logger,
			topup_errors.ErrFailedFindYearlyTopupStatusFailedByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyTopupStatusFailedByCardNumberCache(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly topup status failed by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *topupService) FindMonthlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetMonthlyTopupMethodsByCardNumberRow, error) {
	const method = "FindMonthlyTopupMethodsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTopupMethodsByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched monthly topup methods by card number", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	s.logger.Debug("Cache miss for monthly topup methods by card number, fetching from DB",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	dbRows, err := s.topupRepository.GetMonthlyTopupMethodsByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTopupMethodsByCardNumberRow](
			s.logger,
			topup_errors.ErrFailedFindMonthlyTopupMethodsByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyTopupMethodsByCardNumberCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly topup methods by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *topupService) FindYearlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetYearlyTopupMethodsByCardNumberRow, error) {
	const method = "FindYearlyTopupMethodsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupMethodsByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched yearly topup methods by card number", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	s.logger.Debug("Cache miss for yearly topup methods by card number, fetching from DB",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	dbRows, err := s.topupRepository.GetYearlyTopupMethodsByCardNumber(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTopupMethodsByCardNumberRow](
			s.logger,
			topup_errors.ErrFailedFindYearlyTopupMethodsByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.logger.Debug("Setting cache for yearly topup methods by card number",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	s.cache.SetYearlyTopupMethodsByCardNumberCache(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly topup methods by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *topupService) FindMonthlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetMonthlyTopupAmountsByCardNumberRow, error) {
	const method = "FindMonthlyTopupAmountsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTopupAmountsByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched monthly topup amounts by card number", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	s.logger.Debug("Cache miss for monthly topup amounts by card number, fetching from DB",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	dbRows, err := s.topupRepository.GetMonthlyTopupAmountsByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTopupAmountsByCardNumberRow](
			s.logger,
			topup_errors.ErrFailedFindMonthlyTopupAmountsByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyTopupAmountsByCardNumberCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly topup amounts by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *topupService) FindYearlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetYearlyTopupAmountsByCardNumberRow, error) {
	const method = "FindYearlyTopupAmountsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupAmountsByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched yearly topup amounts by card number", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	dbRows, err := s.topupRepository.GetYearlyTopupAmountsByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTopupAmountsByCardNumberRow](
			s.logger,
			topup_errors.ErrFailedFindYearlyTopupAmountsByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.logger.Debug("Setting cache for yearly topup amounts by card number",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	s.cache.SetYearlyTopupAmountsByCardNumberCache(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly topup amounts by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *topupService) CreateTopup(ctx context.Context, request *requests.CreateTopupRequest) (*db.CreateTopupRow, error) {
	const method = "CreateTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", request.CardNumber),
		attribute.Float64("topup_amount", float64(request.TopupAmount)))

	defer func() {
		end(status)
	}()

	card, err := s.cardRepository.FindCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTopupRow](
			s.logger,
			card_errors.ErrCardNotFoundRes,
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	topup, err := s.topupRepository.CreateTopup(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTopupRow](
			s.logger,
			topup_errors.ErrFailedCreateTopup,
			method,
			span,
		)
	}

	markTopupAsFailed := func(primaryErr error) error {
		s.logger.Error("Attempting to mark topup as failed due to an error", zap.Int("topup_id", int(topup.TopupID)), zap.Error(primaryErr))
		if _, statusErr := s.topupRepository.UpdateTopupStatus(ctx, &requests.UpdateTopupStatus{
			TopupID: int(topup.TopupID),
			Status:  "failed",
		}); statusErr != nil {
			s.logger.Error("CRITICAL: Failed to update topup status to 'failed' after another error", zap.NamedError("primary_error", primaryErr), zap.NamedError("status_update_error", statusErr))
		}
		return primaryErr
	}

	saldo, err := s.saldoRepository.FindByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTopupRow](
			s.logger,
			markTopupAsFailed(saldo_errors.ErrFailedSaldoNotFound),
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	newBalance := int(saldo.TotalBalance) + request.TopupAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: newBalance,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTopupRow](
			s.logger,
			markTopupAsFailed(topup_errors.ErrFailedUpdateTopup),
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	if !card.ExpireDate.Valid {
		status = "error"
		return errorhandler.HandleError[*db.CreateTopupRow](
			s.logger,
			markTopupAsFailed(topup_errors.ErrFailedUpdateTopup),
			method,
			span,
			zap.String("reason", "expire_date is NULL"),
		)
	}

	expireDate := card.ExpireDate.Time

	_, err = s.cardRepository.UpdateCard(ctx, &requests.UpdateCardRequest{
		CardID:       int(card.CardID),
		UserID:       int(card.UserID),
		CardType:     card.CardType,
		ExpireDate:   expireDate,
		CVV:          card.Cvv,
		CardProvider: card.CardProvider,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTopupRow](
			s.logger,
			markTopupAsFailed(card_errors.ErrFailedUpdateCard),
			method,
			span,

			zap.Int("card_id", int(card.CardID)),
		)
	}

	_, err = s.topupRepository.UpdateTopupStatus(ctx, &requests.UpdateTopupStatus{
		TopupID: int(topup.TopupID),
		Status:  "success",
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTopupRow](
			s.logger,
			topup_errors.ErrFailedUpdateTopup,
			method,
			span,

			zap.Int("topup_id", int(topup.TopupID)),
		)
	}

	logSuccess("CreateTopup process completed",
		zap.String("cardNumber", request.CardNumber),
		zap.Float64("topupAmount", float64(request.TopupAmount)),
		zap.Float64("newBalance", float64(newBalance)),
	)

	return topup, nil
}

func (s *topupService) UpdateTopup(ctx context.Context, request *requests.UpdateTopupRequest) (*db.UpdateTopupRow, error) {
	const method = "UpdateTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", request.CardNumber),
		attribute.Int("topup_id", *request.TopupID),
		attribute.Float64("new_topup_amount", float64(request.TopupAmount)))

	defer func() {
		end(status)
	}()

	markTopupAsFailed := func(primaryErr error) error {
		if _, statusErr := s.topupRepository.UpdateTopupStatus(ctx, &requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}); statusErr != nil {
			s.logger.Error("CRITICAL: Failed to update topup status to 'failed' after another error", zap.NamedError("primary_error", primaryErr), zap.NamedError("status_update_error", statusErr))
		}
		return primaryErr
	}

	_, err := s.cardRepository.FindCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTopupRow](
			s.logger,
			markTopupAsFailed(card_errors.ErrCardNotFoundRes),
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	existingTopup, err := s.topupRepository.FindById(ctx, *request.TopupID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTopupRow](
			s.logger,
			markTopupAsFailed(topup_errors.ErrTopupNotFoundRes),
			method,
			span,

			zap.Int("topup_id", *request.TopupID),
		)
	}

	topupDifference := request.TopupAmount - int(existingTopup.TopupAmount)

	updated, err := s.topupRepository.UpdateTopup(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTopupRow](
			s.logger,
			markTopupAsFailed(topup_errors.ErrFailedUpdateTopup),
			method,
			span,

			zap.Int("topup_id", *request.TopupID),
		)
	}

	currentSaldo, err := s.saldoRepository.FindByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTopupRow](
			s.logger,
			markTopupAsFailed(saldo_errors.ErrFailedSaldoNotFound),
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	newBalance := int(currentSaldo.TotalBalance) + topupDifference

	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: newBalance,
	})
	if err != nil {
		status = "error"
		s.logger.Warn("Attempting to rollback topup amount due to saldo update failure", zap.Int("topup_id", *request.TopupID))
		var rollbackErr error
		if _, rollbackErr = s.topupRepository.UpdateTopupAmount(ctx, &requests.UpdateTopupAmount{
			TopupID:     *request.TopupID,
			TopupAmount: int(existingTopup.TopupAmount),
		}); rollbackErr != nil {
			s.logger.Error("CRITICAL: Failed to rollback topup update", zap.NamedError("primary_error", err), zap.NamedError("rollback_error", rollbackErr))
		}

		return errorhandler.HandleError[*db.UpdateTopupRow](
			s.logger,
			markTopupAsFailed(saldo_errors.ErrFailedUpdateSaldo),
			method,
			span,

			zap.Int("topup_id", *request.TopupID),
		)
	}

	_, err = s.topupRepository.FindById(ctx, *request.TopupID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTopupRow](
			s.logger,
			markTopupAsFailed(topup_errors.ErrTopupNotFoundRes),
			method,
			span,

			zap.Int("topup_id", *request.TopupID),
		)
	}

	_, err = s.topupRepository.UpdateTopupStatus(ctx, &requests.UpdateTopupStatus{
		TopupID: *request.TopupID,
		Status:  "success",
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTopupRow](
			s.logger,
			topup_errors.ErrFailedUpdateTopup,
			method,
			span,

			zap.Int("topup_id", *request.TopupID),
		)
	}

	logSuccess("UpdateTopup process completed",
		zap.String("cardNumber", request.CardNumber),
		zap.Int("topupID", *request.TopupID),
		zap.Float64("newTopupAmount", float64(request.TopupAmount)),
		zap.Float64("newBalance", float64(newBalance)),
	)

	return updated, nil
}

func (s *topupService) TrashedTopup(ctx context.Context, topup_id int) (*db.Topup, error) {
	const method = "TrashedTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("topup_id", topup_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting TrashedTopup process", zap.Int("topup_id", topup_id))

	res, err := s.topupRepository.TrashedTopup(ctx, topup_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Topup](
			s.logger,
			topup_errors.ErrFailedTrashTopup,
			method,
			span,

			zap.Int("topup_id", topup_id),
		)
	}

	logSuccess("TrashedTopup process completed", zap.Int("topup_id", topup_id))

	return res, nil
}

func (s *topupService) RestoreTopup(ctx context.Context, topup_id int) (*db.Topup, error) {
	const method = "RestoreTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("topup_id", topup_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting RestoreTopup process", zap.Int("topup_id", topup_id))

	res, err := s.topupRepository.RestoreTopup(ctx, topup_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Topup](
			s.logger,
			topup_errors.ErrFailedRestoreTopup,
			method,
			span,

			zap.Int("topup_id", topup_id),
		)
	}

	logSuccess("RestoreTopup process completed", zap.Int("topup_id", topup_id))

	return res, nil
}

func (s *topupService) DeleteTopupPermanent(ctx context.Context, topup_id int) (bool, error) {
	const method = "DeleteTopupPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("topup_id", topup_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting DeleteTopupPermanent process", zap.Int("topup_id", topup_id))

	_, err := s.topupRepository.DeleteTopupPermanent(ctx, topup_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			topup_errors.ErrFailedDeleteTopup,
			method,
			span,

			zap.Int("topup_id", topup_id),
		)
	}

	logSuccess("DeleteTopupPermanent process completed", zap.Int("topup_id", topup_id))

	return true, nil
}

func (s *topupService) RestoreAllTopup(ctx context.Context) (bool, error) {
	const method = "RestoreAllTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring all topups")

	_, err := s.topupRepository.RestoreAllTopup(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			topup_errors.ErrFailedRestoreAllTopups,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all topups")
	return true, nil
}

func (s *topupService) DeleteAllTopupPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllTopupPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Permanently deleting all topups")

	_, err := s.topupRepository.DeleteAllTopupPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			topup_errors.ErrFailedDeleteAllTopups,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all topups permanently")
	return true, nil
}
