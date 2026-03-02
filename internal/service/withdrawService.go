package service

import (
	withdraw_cache "MamangRust/paymentgatewaygrpc/internal/cache/withdraw"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/errorhandler"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/withdraw_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type withdrawService struct {
	userRepository     repository.UserRepository
	saldoRepository    repository.SaldoRepository
	withdrawRepository repository.WithdrawRepository
	logger             logger.LoggerInterface
	cache              withdraw_cache.WithdrawMencache
	observability      observability.TraceLoggerObservability
}

type WithdrawServiceDeps struct {
	UserRepo      repository.UserRepository
	SaldoRepo     repository.SaldoRepository
	WithdrawRepo  repository.WithdrawRepository
	Cache         withdraw_cache.WithdrawMencache
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewWithdrawService(deps WithdrawServiceDeps) *withdrawService {
	return &withdrawService{
		userRepository:     deps.UserRepo,
		saldoRepository:    deps.SaldoRepo,
		withdrawRepository: deps.WithdrawRepo,
		logger:             deps.Logger,
		observability:      deps.Observability,
		cache:              deps.Cache,
	}
}

func (s *withdrawService) FindAll(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetWithdrawsRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedWithdrawsCache(ctx, req); found {
		logSuccess("Successfully retrieved all withdraw records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	withdraws, err := s.withdrawRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetWithdrawsRow](
			s.logger,
			withdraw_errors.ErrFailedFindAllWithdraws,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(withdraws) > 0 {
		totalCount = int(withdraws[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedWithdrawsCache(ctx, req, withdraws, &totalCount)

	logSuccess("Successfully fetched withdraw",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdraws, &totalCount, nil
}

func (s *withdrawService) FindAllByCardNumber(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*db.GetWithdrawsByCardNumberRow, *int, error) {
	const method = "FindAllByCardNumber"

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
		attribute.String("search", search),
		attribute.String("card_number", req.CardNumber))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedWithdrawByCardCache(ctx, req); found {
		logSuccess("Successfully retrieved all withdraw records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	withdraws, err := s.withdrawRepository.FindAllByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetWithdrawsByCardNumberRow](
			s.logger,
			withdraw_errors.ErrFailedFindAllWithdrawsByCard,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("card_number", req.CardNumber),
		)
	}

	var totalCount int

	if len(withdraws) > 0 {
		totalCount = int(withdraws[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedWithdrawByCardCache(ctx, req, withdraws, &totalCount)

	logSuccess("Successfully fetched withdraw",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdraws, &totalCount, nil
}

func (s *withdrawService) FindByActive(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetActiveWithdrawsRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedWithdrawActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved all withdraw records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	withdraws, err := s.withdrawRepository.FindByActive(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetActiveWithdrawsRow](
			s.logger,
			withdraw_errors.ErrFailedFindActiveWithdraws,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(withdraws) > 0 {
		totalCount = int(withdraws[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedWithdrawActiveCache(ctx, req, withdraws, &totalCount)

	logSuccess("Successfully fetched active withdraw",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdraws, &totalCount, nil
}

func (s *withdrawService) FindByTrashed(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetTrashedWithdrawsRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedWithdrawTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved all withdraw records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	withdraws, err := s.withdrawRepository.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTrashedWithdrawsRow](
			s.logger,
			withdraw_errors.ErrFailedFindTrashedWithdraws,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(withdraws) > 0 {
		totalCount = int(withdraws[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedWithdrawTrashedCache(ctx, req, withdraws, &totalCount)

	logSuccess("Successfully fetched trashed withdraw",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdraws, &totalCount, nil
}

func (s *withdrawService) FindById(ctx context.Context, withdrawID int) (*db.GetWithdrawByIDRow, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("withdraw_id", withdrawID))

	defer func() {
		end(status)
	}()

	withdraw, err := s.withdrawRepository.FindById(ctx, withdrawID)

	if data, found := s.cache.GetCachedWithdrawCache(ctx, withdrawID); found {
		logSuccess("Successfully retrieved withdraw from cache", zap.Int("withdraw_id", withdrawID))
		return data, nil
	}

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetWithdrawByIDRow](
			s.logger,
			withdraw_errors.ErrWithdrawNotFound,
			method,
			span,

			zap.Int("withdraw_id", withdrawID),
		)
	}

	s.cache.SetCachedWithdrawCache(ctx, withdraw)

	logSuccess("Successfully fetched withdraw", zap.Int("withdraw_id", withdrawID))

	return withdraw, nil
}

func (s *withdrawService) FindMonthWithdrawStatusSuccess(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*db.GetMonthWithdrawStatusSuccessRow, error) {
	const method = "FindMonthWithdrawStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Checking cache for monthly withdraw status success", zap.Int("year", req.Year), zap.Int("month", req.Month))

	if dbRows, found := s.cache.GetCachedMonthWithdrawStatusSuccessCache(ctx, req); found {
		s.logger.Info("Cache hit for monthly withdraw status success", zap.Int("year", req.Year), zap.Int("month", req.Month))
		status = "ok"
		logSuccess("Successfully fetched monthly withdraw status success (from cache)", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return dbRows, nil
	}

	dbRows, err := s.withdrawRepository.GetMonthWithdrawStatusSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthWithdrawStatusSuccessRow](
			s.logger,
			withdraw_errors.ErrFailedFindMonthWithdrawStatusSuccess,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetCachedMonthWithdrawStatusSuccessCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly withdraw status success (from DB)", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *withdrawService) FindYearlyWithdrawStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyWithdrawStatusSuccessRow, error) {
	const method = "FindYearlyWithdrawStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedYearlyWithdrawStatusSuccessCache(ctx, year); found {
		logSuccess("Successfully fetched yearly withdraw status success (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.withdrawRepository.GetYearlyWithdrawStatusSuccess(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyWithdrawStatusSuccessRow](
			s.logger,
			withdraw_errors.ErrFailedFindYearWithdrawStatusSuccess,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedYearlyWithdrawStatusSuccessCache(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly withdraw status success (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *withdrawService) FindMonthWithdrawStatusFailed(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*db.GetMonthWithdrawStatusFailedRow, error) {
	const method = "FindMonthWithdrawStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedMonthWithdrawStatusFailedCache(ctx, req); found {
		logSuccess("Successfully fetched monthly withdraw status failed (from cache)", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return dbRows, nil
	}

	dbRows, err := s.withdrawRepository.GetMonthWithdrawStatusFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthWithdrawStatusFailedRow](
			s.logger,
			withdraw_errors.ErrFailedFindMonthWithdrawStatusFailed,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetCachedMonthWithdrawStatusFailedCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly withdraw status failed (from DB)", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *withdrawService) FindYearlyWithdrawStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyWithdrawStatusFailedRow, error) {
	const method = "FindYearlyWithdrawStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedYearlyWithdrawStatusFailedCache(ctx, year); found {
		logSuccess("Successfully fetched yearly withdraw status failed (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.withdrawRepository.GetYearlyWithdrawStatusFailed(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyWithdrawStatusFailedRow](
			s.logger,
			withdraw_errors.ErrFailedFindYearWithdrawStatusFailed,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedYearlyWithdrawStatusFailedCache(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly withdraw status failed (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *withdrawService) FindMonthlyWithdraws(ctx context.Context, year int) ([]*db.GetMonthlyWithdrawsRow, error) {
	const method = "FindMonthlyWithdraws"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedMonthlyWithdraws(ctx, year); found {
		logSuccess("Successfully fetched monthly withdraws (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.withdrawRepository.GetMonthlyWithdraws(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyWithdrawsRow](
			s.logger,
			withdraw_errors.ErrFailedFindMonthlyWithdraws,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedMonthlyWithdraws(ctx, year, dbRows)

	logSuccess("Successfully fetched monthly withdraws (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *withdrawService) FindYearlyWithdraws(ctx context.Context, year int) ([]*db.GetYearlyWithdrawsRow, error) {
	const method = "FindYearlyWithdraws"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedYearlyWithdraws(ctx, year); found {
		logSuccess("Successfully fetched yearly withdraws (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.withdrawRepository.GetYearlyWithdraws(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyWithdrawsRow](
			s.logger,
			withdraw_errors.ErrFailedFindYearlyWithdraws,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedYearlyWithdraws(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly withdraws (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *withdrawService) FindMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*db.GetMonthWithdrawStatusSuccessCardNumberRow, error) {
	const method = "FindMonthWithdrawStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedMonthWithdrawStatusSuccessByCardNumber(ctx, req); found {
		logSuccess("Successfully fetched monthly withdraw status success by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month))
		return dbRows, nil
	}

	dbRows, err := s.withdrawRepository.GetMonthWithdrawStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthWithdrawStatusSuccessCardNumberRow](
			s.logger,
			withdraw_errors.ErrFailedFindMonthWithdrawStatusSuccess,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetCachedMonthWithdrawStatusSuccessByCardNumber(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly withdraw status success by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *withdrawService) FindYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*db.GetYearlyWithdrawStatusSuccessCardNumberRow, error) {
	const method = "FindYearlyWithdrawStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx, req); found {
		logSuccess("Successfully fetched yearly withdraw status success by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.withdrawRepository.GetYearlyWithdrawStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyWithdrawStatusSuccessCardNumberRow](
			s.logger,
			withdraw_errors.ErrFailedFindYearWithdrawStatusSuccess,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly withdraw status success by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *withdrawService) FindMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*db.GetMonthWithdrawStatusFailedCardNumberRow, error) {
	const method = "FindMonthWithdrawStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedMonthWithdrawStatusFailedByCardNumber(ctx, req); found {
		logSuccess("Successfully fetched monthly withdraw status failed by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month))
		return dbRows, nil
	}

	dbRows, err := s.withdrawRepository.GetMonthWithdrawStatusFailedByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthWithdrawStatusFailedCardNumberRow](
			s.logger,
			withdraw_errors.ErrFailedFindMonthWithdrawStatusFailed,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetCachedMonthWithdrawStatusFailedByCardNumber(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly withdraw status failed by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *withdrawService) FindYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*db.GetYearlyWithdrawStatusFailedCardNumberRow, error) {
	const method = "FindYearlyWithdrawStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedYearlyWithdrawStatusFailedByCardNumber(ctx, req); found {
		logSuccess("Successfully fetched yearly withdraw status failed by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.withdrawRepository.GetYearlyWithdrawStatusFailedByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyWithdrawStatusFailedCardNumberRow](
			s.logger,
			withdraw_errors.ErrFailedFindYearWithdrawStatusFailed,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetCachedYearlyWithdrawStatusFailedByCardNumber(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly withdraw status failed by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *withdrawService) FindMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetMonthlyWithdrawsByCardNumberRow, error) {
	const method = "FindMonthlyWithdrawsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedMonthlyWithdrawsByCardNumber(ctx, req); found {
		logSuccess("Successfully fetched monthly withdraws by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.withdrawRepository.GetMonthlyWithdrawsByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyWithdrawsByCardNumberRow](
			s.logger,
			withdraw_errors.ErrFailedFindMonthlyWithdraws,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetCachedMonthlyWithdrawsByCardNumber(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly withdraws by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *withdrawService) FindYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetYearlyWithdrawsByCardNumberRow, error) {
	const method = "FindYearlyWithdrawsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedYearlyWithdrawsByCardNumber(ctx, req); found {
		logSuccess("Successfully fetched yearly withdraws by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.withdrawRepository.GetYearlyWithdrawsByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyWithdrawsByCardNumberRow](
			s.logger,
			withdraw_errors.ErrFailedFindYearlyWithdraws,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetCachedYearlyWithdrawsByCardNumber(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly withdraws by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *withdrawService) Create(ctx context.Context, request *requests.CreateWithdrawRequest) (*db.CreateWithdrawRow, error) {
	const method = "CreateWithdraw"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", request.CardNumber),
		attribute.Int("withdraw_amount", request.WithdrawAmount))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Creating new withdraw", zap.Any("request", request))

	saldo, err := s.saldoRepository.FindByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateWithdrawRow](
			s.logger,
			saldo_errors.ErrFailedSaldoNotFound,
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	if saldo == nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateWithdrawRow](
			s.logger,
			saldo_errors.ErrFailedSaldoNotFound,
			method,
			span,
			zap.String("card_number", request.CardNumber),
		)
	}

	if int(saldo.TotalBalance) < request.WithdrawAmount {
		status = "error"
		return errorhandler.HandleError[*db.CreateWithdrawRow](
			s.logger,
			saldo_errors.ErrFailedInsufficientBalance,
			method,
			span,
			zap.String("card_number", request.CardNumber),
			zap.Int("requested_amount", request.WithdrawAmount),
			zap.Int32("current_balance", saldo.TotalBalance),
		)
	}

	newTotalBalance := int(saldo.TotalBalance) - request.WithdrawAmount
	updateData := &requests.UpdateSaldoWithdraw{
		CardNumber:     request.CardNumber,
		TotalBalance:   newTotalBalance,
		WithdrawAmount: &request.WithdrawAmount,
		WithdrawTime:   &request.WithdrawTime,
	}

	_, err = s.saldoRepository.UpdateSaldoWithdraw(ctx, updateData)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateWithdrawRow](
			s.logger,
			saldo_errors.ErrFailedUpdateSaldo,
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	withdrawRecord, err := s.withdrawRepository.CreateWithdraw(ctx, request)
	if err != nil {
		status = "error"
		s.logger.Warn("Attempting to rollback saldo due to withdraw creation failure", zap.String("card_number", request.CardNumber))
		rollbackData := &requests.UpdateSaldoWithdraw{
			CardNumber:     request.CardNumber,
			TotalBalance:   int(saldo.TotalBalance),
			WithdrawAmount: &request.WithdrawAmount,
			WithdrawTime:   &request.WithdrawTime,
		}
		var rollbackErr error
		if _, rollbackErr = s.saldoRepository.UpdateSaldoWithdraw(ctx, rollbackData); rollbackErr != nil {
			s.logger.Error("Failed to rollback saldo after withdraw creation failure", zap.Error(rollbackErr))
		}

		errorFields := []zap.Field{

			zap.String("card_number", request.CardNumber),
		}
		if rollbackErr != nil {
			errorFields = append(errorFields, zap.NamedError("rollback_error", rollbackErr))
		}

		return errorhandler.HandleError[*db.CreateWithdrawRow](
			s.logger,
			withdraw_errors.ErrFailedCreateWithdraw,
			method,
			span,
			errorFields...,
		)
	}

	if _, err := s.withdrawRepository.UpdateWithdrawStatus(ctx, &requests.UpdateWithdrawStatus{
		WithdrawID: int(withdrawRecord.WithdrawID),
		Status:     "success",
	}); err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateWithdrawRow](
			s.logger,
			withdraw_errors.ErrFailedUpdateWithdraw,
			method,
			span,

			zap.Int("withdraw_id", int(withdrawRecord.WithdrawID)),
		)
	}

	logSuccess("Successfully created withdraw", zap.Int("withdraw_id", int(withdrawRecord.WithdrawID)))

	return withdrawRecord, nil
}

func (s *withdrawService) Update(ctx context.Context, request *requests.UpdateWithdrawRequest) (*db.UpdateWithdrawRow, error) {
	const method = "UpdateWithdraw"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("withdraw_id", *request.WithdrawID),
		attribute.String("card_number", request.CardNumber))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Updating withdraw", zap.Int("withdraw_id", *request.WithdrawID), zap.Any("request", request))

	_, err := s.withdrawRepository.FindById(ctx, *request.WithdrawID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateWithdrawRow](
			s.logger,
			withdraw_errors.ErrWithdrawNotFound,
			method,
			span,

			zap.Int("withdraw_id", *request.WithdrawID),
		)
	}

	saldo, err := s.saldoRepository.FindByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateWithdrawRow](
			s.logger,
			saldo_errors.ErrFailedSaldoNotFound,
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	if int(saldo.TotalBalance) < request.WithdrawAmount {
		status = "error"
		return errorhandler.HandleError[*db.UpdateWithdrawRow](
			s.logger,
			saldo_errors.ErrFailedInsufficientBalance,
			method,
			span,
			zap.String("card_number", request.CardNumber),
			zap.Int("requested_amount", request.WithdrawAmount),
			zap.Int32("current_balance", saldo.TotalBalance),
		)
	}

	newTotalBalance := int(saldo.TotalBalance) - request.WithdrawAmount
	updateSaldoData := &requests.UpdateSaldoWithdraw{
		CardNumber:     saldo.CardNumber,
		TotalBalance:   newTotalBalance,
		WithdrawAmount: &request.WithdrawAmount,
		WithdrawTime:   &request.WithdrawTime,
	}
	_, err = s.saldoRepository.UpdateSaldoWithdraw(ctx, updateSaldoData)
	if err != nil {
		status = "error"
		if _, statusErr := s.withdrawRepository.UpdateWithdrawStatus(ctx, &requests.UpdateWithdrawStatus{
			WithdrawID: *request.WithdrawID,
			Status:     "failed",
		}); statusErr != nil {
			s.logger.Error("Failed to update withdraw status to 'failed'", zap.Error(statusErr))
		}
		return errorhandler.HandleError[*db.UpdateWithdrawRow](
			s.logger,
			saldo_errors.ErrFailedUpdateSaldo,
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	updatedWithdraw, err := s.withdrawRepository.UpdateWithdraw(ctx, request)
	if err != nil {
		status = "error"
		s.logger.Warn("Attempting to rollback saldo due to withdraw update failure", zap.String("card_number", request.CardNumber))
		rollbackData := &requests.UpdateSaldoBalance{
			CardNumber:   saldo.CardNumber,
			TotalBalance: int(saldo.TotalBalance),
		}
		var rollbackErr error
		if _, rollbackErr = s.saldoRepository.UpdateSaldoBalance(ctx, rollbackData); rollbackErr != nil {
			s.logger.Error("Failed to rollback saldo after withdraw update failure", zap.Error(rollbackErr))
		}

		if _, statusErr := s.withdrawRepository.UpdateWithdrawStatus(ctx, &requests.UpdateWithdrawStatus{
			WithdrawID: *request.WithdrawID,
			Status:     "failed",
		}); statusErr != nil {
			s.logger.Error("Failed to update withdraw status to 'failed'", zap.Error(statusErr))
		}

		errorFields := []zap.Field{

			zap.Int("withdraw_id", *request.WithdrawID),
		}
		if rollbackErr != nil {
			errorFields = append(errorFields, zap.NamedError("rollback_error", rollbackErr))
		}

		return errorhandler.HandleError[*db.UpdateWithdrawRow](
			s.logger,
			withdraw_errors.ErrFailedUpdateWithdraw,
			method,
			span,
			errorFields...,
		)
	}

	if _, err := s.withdrawRepository.UpdateWithdrawStatus(ctx, &requests.UpdateWithdrawStatus{
		WithdrawID: int(updatedWithdraw.WithdrawID),
		Status:     "success",
	}); err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateWithdrawRow](
			s.logger,
			withdraw_errors.ErrFailedUpdateWithdraw,
			method,
			span,

			zap.Int("withdraw_id", int(updatedWithdraw.WithdrawID)),
		)
	}

	logSuccess("Successfully updated withdraw", zap.Int("withdraw_id", int(updatedWithdraw.WithdrawID)))

	return updatedWithdraw, nil
}

func (s *withdrawService) TrashedWithdraw(ctx context.Context, withdraw_id int) (*db.Withdraw, error) {
	const method = "TrashedWithdraw"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("withdraw_id", withdraw_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Trashing withdraw", zap.Int("withdraw_id", withdraw_id))

	res, err := s.withdrawRepository.TrashedWithdraw(ctx, withdraw_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Withdraw](
			s.logger,
			withdraw_errors.ErrFailedTrashedWithdraw,
			method,
			span,

			zap.Int("withdraw_id", withdraw_id),
		)
	}

	logSuccess("Successfully trashed withdraw", zap.Int("withdraw_id", withdraw_id))

	return res, nil
}

func (s *withdrawService) RestoreWithdraw(ctx context.Context, withdraw_id int) (*db.Withdraw, error) {
	const method = "RestoreWithdraw"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("withdraw_id", withdraw_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring withdraw", zap.Int("withdraw_id", withdraw_id))

	res, err := s.withdrawRepository.RestoreWithdraw(ctx, withdraw_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Withdraw](
			s.logger,
			withdraw_errors.ErrFailedRestoreWithdraw,
			method,
			span,

			zap.Int("withdraw_id", withdraw_id),
		)
	}

	logSuccess("Successfully restored withdraw", zap.Int("withdraw_id", withdraw_id))

	return res, nil
}

func (s *withdrawService) DeleteWithdrawPermanent(ctx context.Context, withdraw_id int) (bool, error) {
	const method = "DeleteWithdrawPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("withdraw_id", withdraw_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Deleting withdraw permanently", zap.Int("withdraw_id", withdraw_id))

	_, err := s.withdrawRepository.DeleteWithdrawPermanent(ctx, withdraw_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			withdraw_errors.ErrFailedDeleteWithdrawPermanent,
			method,
			span,

			zap.Int("withdraw_id", withdraw_id),
		)
	}

	logSuccess("Successfully deleted withdraw permanently", zap.Int("withdraw_id", withdraw_id))

	return true, nil
}

func (s *withdrawService) RestoreAllWithdraw(ctx context.Context) (bool, error) {
	const method = "RestoreAllWithdraw"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring all withdraws")

	_, err := s.withdrawRepository.RestoreAllWithdraw(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			withdraw_errors.ErrFailedRestoreAllWithdraw,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all withdraws")
	return true, nil
}

func (s *withdrawService) DeleteAllWithdrawPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllWithdrawPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Permanently deleting all withdraws")

	_, err := s.withdrawRepository.DeleteAllWithdrawPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			withdraw_errors.ErrFailedDeleteAllWithdrawPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all withdraws permanently")
	return true, nil
}
