package service

import (
	saldo_cache "MamangRust/paymentgatewaygrpc/internal/cache/saldo"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/errorhandler"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/card_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type saldoService struct {
	cardRepository  repository.CardRepository
	saldoRepository repository.SaldoRepository
	logger          logger.LoggerInterface
	cache           saldo_cache.SaldoMencache
	observability   observability.TraceLoggerObservability
}

type SaldoServiceDeps struct {
	SaldoRepo     repository.SaldoRepository
	CardRepo      repository.CardRepository
	Cache         saldo_cache.SaldoMencache
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewSaldoService(deps SaldoServiceDeps) *saldoService {
	return &saldoService{
		saldoRepository: deps.SaldoRepo,
		cardRepository:  deps.CardRepo,
		logger:          deps.Logger,
		cache:           deps.Cache,
		observability:   deps.Observability,
	}
}

func (s *saldoService) FindAll(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetSaldosRow, *int, error) {
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

	res, err := s.saldoRepository.FindAllSaldos(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetSaldosRow](
			s.logger,
			saldo_errors.ErrFailedFindAllSaldos,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	logSuccess("Successfully fetched saldo",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return res, &totalCount, nil
}

func (s *saldoService) FindByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetActiveSaldosRow, *int, error) {
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

	res, err := s.saldoRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetActiveSaldosRow](
			s.logger,
			saldo_errors.ErrFailedFindActiveSaldos,
			method,
			span,

			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	logSuccess("Successfully fetched active saldo",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return res, &totalCount, nil
}

func (s *saldoService) FindByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetTrashedSaldosRow, *int, error) {
	const method = "FindByTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", req.Page),
		attribute.Int("pageSize", req.PageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, err := s.saldoRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTrashedSaldosRow](
			s.logger,
			saldo_errors.ErrFailedFindTrashedSaldos,
			method,
			span,

			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	logSuccess("Successfully fetched trashed saldo",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return res, &totalCount, nil
}

func (s *saldoService) FindById(ctx context.Context, saldo_id int) (*db.GetSaldoByIDRow, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("saldo_id", saldo_id))

	defer func() {
		end(status)
	}()

	res, err := s.saldoRepository.FindById(ctx, saldo_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetSaldoByIDRow](
			s.logger,
			saldo_errors.ErrFailedSaldoNotFound,
			method,
			span,

			zap.Int("saldo_id", saldo_id),
		)
	}

	logSuccess("Successfully fetched saldo", zap.Int("saldo_id", saldo_id))

	return res, nil
}

func (s *saldoService) FindByCardNumber(ctx context.Context, card_number string) (*db.Saldo, error) {
	const method = "FindByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	res, err := s.saldoRepository.FindByCardNumber(ctx, card_number)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Saldo](
			s.logger,
			saldo_errors.ErrFailedSaldoNotFound,
			method,
			span,

			zap.String("card_number", card_number),
		)
	}

	logSuccess("Successfully fetched saldo by card number", zap.String("card_number", card_number))

	return res, nil
}

func (s *saldoService) FindMonthlyTotalSaldoBalance(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*db.GetMonthlyTotalSaldoBalanceRow, error) {
	const method = "FindMonthlyTotalSaldoBalance"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if cache, found := s.cache.GetMonthlyTotalSaldoBalanceCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total saldo balance from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return cache, nil
	}

	dbRows, err := s.saldoRepository.GetMonthlyTotalSaldoBalance(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTotalSaldoBalanceRow](
			s.logger,
			saldo_errors.ErrFailedFindMonthlyTotalSaldoBalance,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.logger.Debug("Setting cache for monthly total saldo balance", zap.Int("year", req.Year), zap.Int("month", req.Month))
	s.cache.SetMonthlyTotalSaldoCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly total saldo balance (from DB)", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *saldoService) FindYearTotalSaldoBalance(ctx context.Context, year int) ([]*db.GetYearlyTotalSaldoBalancesRow, error) {
	const method = "FindYearTotalSaldoBalance"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cache, found := s.cache.GetYearTotalSaldoBalanceCache(ctx, year); found {
		logSuccess("Successfully fetched yearly total saldo balance from cache", zap.Int("year", year))
		return cache, nil
	}

	s.logger.Debug("Cache miss for yearly total saldo balance, fetching from DB", zap.Int("year", year))

	dbRows, err := s.saldoRepository.GetYearTotalSaldoBalance(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTotalSaldoBalancesRow](
			s.logger,
			saldo_errors.ErrFailedFindYearTotalSaldoBalance,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.logger.Debug("Setting cache for yearly total saldo balance", zap.Int("year", year))

	s.cache.SetYearTotalSaldoBalanceCache(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly total saldo balance (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *saldoService) FindMonthlySaldoBalances(ctx context.Context, year int) ([]*db.GetMonthlySaldoBalancesRow, error) {
	const method = "FindMonthlySaldoBalances"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cache, found := s.cache.GetMonthlySaldoBalanceCache(ctx, year); found {
		logSuccess("Successfully fetched monthly saldo balances from cache", zap.Int("year", year))
		return cache, nil
	}

	s.logger.Debug("Cache miss for monthly saldo balances, fetching from DB", zap.Int("year", year))

	dbRows, err := s.saldoRepository.GetMonthlySaldoBalances(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlySaldoBalancesRow](
			s.logger,
			saldo_errors.ErrFailedFindMonthlySaldoBalances,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.logger.Debug("Setting cache for monthly saldo balances", zap.Int("year", year))

	s.cache.SetMonthlySaldoBalanceCache(ctx, year, dbRows)

	logSuccess("Successfully fetched monthly saldo balances (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *saldoService) FindYearlySaldoBalances(ctx context.Context, year int) ([]*db.GetYearlySaldoBalancesRow, error) {
	const method = "FindYearlySaldoBalances"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cache, found := s.cache.GetYearlySaldoBalanceCache(ctx, year); found {
		logSuccess("Successfully fetched yearly saldo balances from cache", zap.Int("year", year))
		return cache, nil
	}

	s.logger.Debug("Cache miss for yearly saldo balances, fetching from DB", zap.Int("year", year))

	dbRows, err := s.saldoRepository.GetYearlySaldoBalances(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlySaldoBalancesRow](
			s.logger,
			saldo_errors.ErrFailedFindYearlySaldoBalances,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.logger.Debug("Setting cache for yearly saldo balances", zap.Int("year", year))
	s.cache.SetYearlySaldoBalanceCache(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly saldo balances (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *saldoService) CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*db.CreateSaldoRow, error) {
	const method = "CreateSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", request.CardNumber))

	defer func() {
		end(status)
	}()

	_, err := s.cardRepository.FindCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateSaldoRow](
			s.logger,
			card_errors.ErrCardNotFoundRes,
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	res, err := s.saldoRepository.CreateSaldo(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateSaldoRow](
			s.logger,
			saldo_errors.ErrFailedCreateSaldo,
			method,
			span,
		)
	}

	logSuccess("Successfully created saldo record", zap.String("card_number", request.CardNumber))

	return res, nil
}

func (s *saldoService) UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*db.UpdateSaldoRow, error) {
	const method = "UpdateSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", request.CardNumber),
		attribute.Float64("amount", float64(request.TotalBalance)))

	defer func() {
		end(status)
	}()

	_, err := s.cardRepository.FindCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateSaldoRow](
			s.logger,
			card_errors.ErrCardNotFoundRes,
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	res, err := s.saldoRepository.UpdateSaldo(ctx, request)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateSaldoRow](
			s.logger,
			saldo_errors.ErrFailedUpdateSaldo,
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	logSuccess("Successfully updated saldo", zap.String("card_number", request.CardNumber), zap.Int("saldo_id", int(res.SaldoID)))

	return res, nil
}

func (s *saldoService) TrashSaldo(ctx context.Context, saldo_id int) (*db.Saldo, error) {
	const method = "TrashSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("saldo_id", saldo_id))

	defer func() {
		end(status)
	}()

	res, err := s.saldoRepository.TrashedSaldo(ctx, saldo_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Saldo](
			s.logger,
			saldo_errors.ErrFailedTrashSaldo,
			method,
			span,

			zap.Int("saldo_id", saldo_id),
		)
	}

	logSuccess("Successfully trashed saldo", zap.Int("saldo_id", saldo_id))

	return res, nil
}

func (s *saldoService) RestoreSaldo(ctx context.Context, saldo_id int) (*db.Saldo, error) {
	const method = "RestoreSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("saldo_id", saldo_id))

	defer func() {
		end(status)
	}()

	res, err := s.saldoRepository.RestoreSaldo(ctx, saldo_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Saldo](
			s.logger,
			saldo_errors.ErrFailedRestoreSaldo,
			method,
			span,

			zap.Int("saldo_id", saldo_id),
		)
	}

	logSuccess("Successfully restored saldo", zap.Int("saldo_id", saldo_id))

	return res, nil
}

func (s *saldoService) DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, error) {
	const method = "DeleteSaldoPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("saldo_id", saldo_id))

	defer func() {
		end(status)
	}()

	_, err := s.saldoRepository.DeleteSaldoPermanent(ctx, saldo_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			saldo_errors.ErrFailedDeleteSaldoPermanent,
			method,
			span,

			zap.Int("saldo_id", saldo_id),
		)
	}

	logSuccess("Successfully deleted saldo permanently", zap.Int("saldo_id", saldo_id))

	return true, nil
}

func (s *saldoService) RestoreAllSaldo(ctx context.Context) (bool, error) {
	const method = "RestoreAllSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.saldoRepository.RestoreAllSaldo(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			saldo_errors.ErrFailedRestoreAllSaldo,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all saldo")
	return true, nil
}

func (s *saldoService) DeleteAllSaldoPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllSaldoPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.saldoRepository.DeleteAllSaldoPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			saldo_errors.ErrFailedDeleteAllSaldoPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all saldo permanently")
	return true, nil
}
