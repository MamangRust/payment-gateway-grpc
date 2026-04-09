package service

import (
	transfer_cache "MamangRust/paymentgatewaygrpc/internal/cache/transfer"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/errorhandler"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/card_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/transfer_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
	"context"
	"fmt"

	"github.com/grafana/pyroscope-go"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transferService struct {
	userRepository     repository.UserRepository
	cardRepository     repository.CardRepository
	saldoRepository    repository.SaldoRepository
	transferRepository repository.TransferRepository
	logger             logger.LoggerInterface
	cache              transfer_cache.TransferMencache
	observability      observability.TraceLoggerObservability
}

type TransferServiceDeps struct {
	UserRepo      repository.UserRepository
	CardRepo      repository.CardRepository
	SaldoRepo     repository.SaldoRepository
	TransferRepo  repository.TransferRepository
	Cache         transfer_cache.TransferMencache
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewTransferService(deps TransferServiceDeps) *transferService {
	return &transferService{
		userRepository:     deps.UserRepo,
		cardRepository:     deps.CardRepo,
		saldoRepository:    deps.SaldoRepo,
		transferRepository: deps.TransferRepo,
		logger:             deps.Logger,
		cache:              deps.Cache,
		observability:      deps.Observability,
	}
}

func (s *transferService) FindAll(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetTransfersRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedTransfersCache(ctx, req); found {
		logSuccess("Successfully retrieved all transfer records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	var transfers []*db.GetTransfersRow
	var err error

	pyroscope.TagWrapper(ctx, pyroscope.Labels(
		"page", fmt.Sprint(page),
		"pageSize", fmt.Sprint(pageSize),
		"search", search,
	), func(ctx context.Context) {
		transfers, err = s.transferRepository.FindAll(ctx, req)
	})

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTransfersRow](
			s.logger,
			transfer_errors.ErrFailedFindAllTransfers,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(transfers) > 0 {
		totalCount = int(transfers[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTransfersCache(ctx, req, transfers, &totalCount)

	logSuccess("Successfully fetched transfer",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return transfers, &totalCount, nil
}

func (s *transferService) FindById(ctx context.Context, transferId int) (*db.GetTransferByIDRow, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transfer_id", transferId))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedTransferCache(ctx, transferId); found {
		logSuccess("Successfully fetched transfer from cache", zap.Int("transfer.id", transferId))
		return data, nil
	}

	var transfer *db.GetTransferByIDRow
	var err error

	pyroscope.TagWrapper(ctx, pyroscope.Labels("transfer_id", fmt.Sprint(transferId)), func(ctx context.Context) {
		transfer, err = s.transferRepository.FindById(ctx, transferId)
	})

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetTransferByIDRow](
			s.logger,
			transfer_errors.ErrTransferNotFound,
			method,
			span,

			zap.Int("transfer_id", transferId),
		)
	}
	s.cache.SetCachedTransferCache(ctx, transfer)

	logSuccess("Successfully fetched transfer", zap.Int("transfer_id", transferId))

	return transfer, nil
}

func (s *transferService) FindByActive(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetActiveTransfersRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedTransferActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active transfer records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transfers, err := s.transferRepository.FindByActive(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetActiveTransfersRow](
			s.logger,
			transfer_errors.ErrFailedFindActiveTransfers,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(transfers) > 0 {
		totalCount = int(transfers[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTransferActiveCache(ctx, req, transfers, &totalCount)

	logSuccess("Successfully fetched active transfer",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return transfers, &totalCount, nil
}

func (s *transferService) FindByTrashed(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetTrashedTransfersRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedTransferTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed transfer records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transfers, err := s.transferRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTrashedTransfersRow](
			s.logger,
			transfer_errors.ErrFailedFindTrashedTransfers,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(transfers) > 0 {
		totalCount = int(transfers[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTransferTrashedCache(ctx, req, transfers, &totalCount)

	logSuccess("Successfully fetched trashed transfer",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return transfers, &totalCount, nil
}

func (s *transferService) FindTransferByTransferFrom(ctx context.Context, transfer_from string) ([]*db.GetTransfersBySourceCardRow, error) {
	const method = "FindTransferByTransferFrom"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("transfer_from", transfer_from))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedTransferByFrom(ctx, transfer_from); found {
		logSuccess("Successfully fetched transfer from cache", zap.String("transfer_from", transfer_from))
		return data, nil
	}

	res, err := s.transferRepository.FindTransferByTransferFrom(ctx, transfer_from)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetTransfersBySourceCardRow](
			s.logger,
			transfer_errors.ErrTransferNotFound,
			method,
			span,

			zap.String("transfer_from", transfer_from),
		)
	}

	s.cache.SetCachedTransferByFrom(ctx, transfer_from, res)

	logSuccess("Successfully fetched transfer record by transfer_from", zap.String("transfer_from", transfer_from))

	return res, nil
}

func (s *transferService) FindTransferByTransferTo(ctx context.Context, transfer_to string) ([]*db.GetTransfersByDestinationCardRow, error) {
	const method = "FindTransferByTransferTo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("transfer_to", transfer_to))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedTransferByTo(ctx, transfer_to); found {
		logSuccess("Successfully fetched transfer from cache", zap.String("transfer_to", transfer_to))
		return data, nil
	}

	res, err := s.transferRepository.FindTransferByTransferTo(ctx, transfer_to)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetTransfersByDestinationCardRow](
			s.logger,
			transfer_errors.ErrTransferNotFound,
			method,
			span,

			zap.String("transfer_to", transfer_to),
		)
	}

	s.cache.SetCachedTransferByTo(ctx, transfer_to, res)

	logSuccess("Successfully fetched transfer record by transfer_to", zap.String("transfer_to", transfer_to))

	return res, nil
}

func (s *transferService) FindMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*db.GetMonthTransferStatusSuccessRow, error) {
	const method = "FindMonthTransferStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedMonthTransferStatusSuccess(ctx, req); found {
		logSuccess("Successfully fetched monthly transfer status success (from cache)", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetMonthTransferStatusSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthTransferStatusSuccessRow](
			s.logger,
			transfer_errors.ErrFailedFindMonthTransferStatusSuccess,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetCachedMonthTransferStatusSuccess(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly transfer status success (from DB)", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *transferService) FindYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransferStatusSuccessRow, error) {
	const method = "FindYearlyTransferStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedYearlyTransferStatusSuccess(ctx, year); found {
		logSuccess("Successfully fetched yearly transfer status success (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetYearlyTransferStatusSuccess(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransferStatusSuccessRow](
			s.logger,
			transfer_errors.ErrFailedFindYearTransferStatusSuccess,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedYearlyTransferStatusSuccess(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly transfer status success (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *transferService) FindMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*db.GetMonthTransferStatusFailedRow, error) {
	const method = "FindMonthTransferStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedMonthTransferStatusFailed(ctx, req); found {
		logSuccess("Successfully fetched monthly transfer status failed (from cache)", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetMonthTransferStatusFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthTransferStatusFailedRow](
			s.logger,
			transfer_errors.ErrFailedFindMonthTransferStatusFailed,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetCachedMonthTransferStatusFailed(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly transfer status failed (from DB)", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *transferService) FindYearlyTransferStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTransferStatusFailedRow, error) {
	const method = "FindYearlyTransferStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedYearlyTransferStatusFailed(ctx, year); found {
		logSuccess("Successfully fetched yearly transfer status failed (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetYearlyTransferStatusFailed(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransferStatusFailedRow](
			s.logger,
			transfer_errors.ErrFailedFindYearTransferStatusFailed,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedYearlyTransferStatusFailed(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly transfer status failed (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *transferService) FindMonthlyTransferAmounts(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountsRow, error) {
	const method = "FindMonthlyTransferAmounts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedMonthTransferAmounts(ctx, year); found {
		logSuccess("Successfully fetched monthly transfer amounts (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetMonthlyTransferAmounts(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransferAmountsRow](
			s.logger,
			transfer_errors.ErrFailedFindMonthlyTransferAmounts,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedMonthTransferAmounts(ctx, year, dbRows)

	logSuccess("Successfully fetched monthly transfer amounts (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *transferService) FindYearlyTransferAmounts(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountsRow, error) {
	const method = "FindYearlyTransferAmounts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetCachedYearlyTransferAmounts(ctx, year); found {
		logSuccess("Successfully fetched yearly transfer amounts (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetYearlyTransferAmounts(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransferAmountsRow](
			s.logger,
			transfer_errors.ErrFailedFindYearlyTransferAmounts,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedYearlyTransferAmounts(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly transfer amounts (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *transferService) FindMonthTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*db.GetMonthTransferStatusSuccessCardNumberRow, error) {
	const method = "FindMonthTransferStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetMonthTransferStatusSuccessByCard(ctx, req); found {
		logSuccess("Successfully fetched monthly transfer status success by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetMonthTransferStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthTransferStatusSuccessCardNumberRow](
			s.logger,
			transfer_errors.ErrFailedFindMonthTransferStatusSuccess,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetMonthTransferStatusSuccessByCard(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly transfer status success by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *transferService) FindYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*db.GetYearlyTransferStatusSuccessCardNumberRow, error) {
	const method = "FindYearlyTransferStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetYearlyTransferStatusSuccessByCard(ctx, req); found {
		logSuccess("Successfully fetched yearly transfer status success by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetYearlyTransferStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransferStatusSuccessCardNumberRow](
			s.logger,
			transfer_errors.ErrFailedFindYearTransferStatusSuccessByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyTransferStatusSuccessByCard(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly transfer status success by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *transferService) FindMonthTransferStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*db.GetMonthTransferStatusFailedCardNumberRow, error) {
	const method = "FindMonthTransferStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetMonthTransferStatusFailedByCard(ctx, req); found {
		logSuccess("Successfully fetched monthly transfer status failed by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetMonthTransferStatusFailedByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthTransferStatusFailedCardNumberRow](
			s.logger,
			transfer_errors.ErrFailedFindMonthTransferStatusFailed,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetMonthTransferStatusFailedByCard(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly transfer status failed by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *transferService) FindYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*db.GetYearlyTransferStatusFailedCardNumberRow, error) {
	const method = "FindYearlyTransferStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetYearlyTransferStatusFailedByCard(ctx, req); found {
		logSuccess("Successfully fetched yearly transfer status failed by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetYearlyTransferStatusFailedByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransferStatusFailedCardNumberRow](
			s.logger,
			transfer_errors.ErrFailedFindYearTransferStatusFailedByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyTransferStatusFailedByCard(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly transfer status failed by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *transferService) FindMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetMonthlyTransferAmountsBySenderCardNumberRow, error) {
	const method = "FindMonthlyTransferAmountsBySenderCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetMonthlyTransferAmountsBySenderCard(ctx, req); found {
		logSuccess("Successfully fetched monthly transfer amounts by sender card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetMonthlyTransferAmountsBySenderCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransferAmountsBySenderCardNumberRow](
			s.logger,
			transfer_errors.ErrFailedFindMonthlyTransferAmountsBySenderCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyTransferAmountsBySenderCard(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly transfer amounts by sender card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *transferService) FindMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetMonthlyTransferAmountsByReceiverCardNumberRow, error) {
	const method = "FindMonthlyTransferAmountsByReceiverCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetMonthlyTransferAmountsByReceiverCard(ctx, req); found {
		logSuccess("Successfully fetched monthly transfer amounts by receiver card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetMonthlyTransferAmountsByReceiverCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransferAmountsByReceiverCardNumberRow](
			s.logger,
			transfer_errors.ErrFailedFindMonthlyTransferAmountsByReceiverCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyTransferAmountsByReceiverCard(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly transfer amounts by receiver card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *transferService) FindYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetYearlyTransferAmountsBySenderCardNumberRow, error) {
	const method = "FindYearlyTransferAmountsBySenderCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetYearlyTransferAmountsBySenderCard(ctx, req); found {
		logSuccess("Successfully fetched yearly transfer amounts by sender card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetYearlyTransferAmountsBySenderCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransferAmountsBySenderCardNumberRow](
			s.logger,
			transfer_errors.ErrFailedFindYearlyTransferAmountsBySenderCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyTransferAmountsBySenderCard(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly transfer amounts by sender card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *transferService) FindYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetYearlyTransferAmountsByReceiverCardNumberRow, error) {
	const method = "FindYearlyTransferAmountsByReceiverCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetYearlyTransferAmountsByReceiverCard(ctx, req); found {
		logSuccess("Successfully fetched yearly transfer amounts by receiver card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.transferRepository.GetYearlyTransferAmountsByReceiverCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransferAmountsByReceiverCardNumberRow](
			s.logger,
			transfer_errors.ErrFailedFindYearlyTransferAmountsByReceiverCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyTransferAmountsByReceiverCard(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly transfer amounts by receiver card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *transferService) CreateTransaction(ctx context.Context, request *requests.CreateTransferRequest) (*db.CreateTransferRow, error) {
	const method = "CreateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("transfer_from", request.TransferFrom),
		attribute.String("transfer_to", request.TransferTo),
		attribute.Int("transfer_amount", request.TransferAmount))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting create transaction process", zap.Any("request", request))

	_, err := s.cardRepository.FindCardByCardNumber(ctx, request.TransferFrom)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransferRow](
			s.logger,
			card_errors.ErrCardNotFoundRes,
			method,
			span,

			zap.String("card_number", request.TransferFrom),
			zap.String("role", "sender"),
		)
	}

	_, err = s.cardRepository.FindCardByCardNumber(ctx, request.TransferTo)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransferRow](
			s.logger,
			card_errors.ErrCardNotFoundRes,
			method,
			span,

			zap.String("card_number", request.TransferTo),
			zap.String("role", "receiver"),
		)
	}

	senderSaldo, err := s.saldoRepository.FindByCardNumber(ctx, request.TransferFrom)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransferRow](
			s.logger,
			saldo_errors.ErrFailedSaldoNotFound,
			method,
			span,

			zap.String("card_number", request.TransferFrom),
			zap.String("role", "sender"),
		)
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(ctx, request.TransferTo)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransferRow](
			s.logger,
			saldo_errors.ErrFailedSaldoNotFound,
			method,
			span,

			zap.String("card_number", request.TransferTo),
			zap.String("role", "receiver"),
		)
	}

	if int(senderSaldo.TotalBalance) < request.TransferAmount {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransferRow](
			s.logger,
			saldo_errors.ErrFailedInsufficientBalance,
			method,
			span,
			zap.Int("available_balance", int(senderSaldo.TotalBalance)),
			zap.Int("transfer_amount", request.TransferAmount),
		)
	}

	rollbackSenderSaldo := func(primaryErr error) error {
		s.logger.Warn("Attempting to rollback sender's saldo due to an error", zap.String("card_number", senderSaldo.CardNumber), zap.Error(primaryErr))
		senderSaldo.TotalBalance += int32(request.TransferAmount)
		if _, rollbackErr := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   senderSaldo.CardNumber,
			TotalBalance: int(senderSaldo.TotalBalance),
		}); rollbackErr != nil {
			s.logger.Error("CRITICAL: Failed to rollback sender's saldo after another error", zap.NamedError("primary_error", primaryErr), zap.NamedError("rollback_error", rollbackErr))
		}
		return primaryErr
	}

	senderSaldo.TotalBalance -= int32(request.TransferAmount)

	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   senderSaldo.CardNumber,
		TotalBalance: int(senderSaldo.TotalBalance),
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransferRow](
			s.logger,
			saldo_errors.ErrFailedUpdateSaldo,
			method,
			span,

			zap.String("card_number", senderSaldo.CardNumber),
			zap.String("role", "sender"),
		)
	}

	receiverSaldo.TotalBalance += int32(request.TransferAmount)

	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   receiverSaldo.CardNumber,
		TotalBalance: int(receiverSaldo.TotalBalance),
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransferRow](
			s.logger,
			rollbackSenderSaldo(saldo_errors.ErrFailedUpdateSaldo),
			method,
			span,

			zap.String("card_number", receiverSaldo.CardNumber),
			zap.String("role", "receiver"),
		)
	}

	rollbackFullTransfer := func(primaryErr error) error {
		s.logger.Warn("Attempting to rollback full transfer due to an error", zap.Error(primaryErr))

		senderSaldo.TotalBalance += int32(request.TransferAmount)

		if _, rollbackErr := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   senderSaldo.CardNumber,
			TotalBalance: int(senderSaldo.TotalBalance),
		}); rollbackErr != nil {
			s.logger.Error("CRITICAL: Failed to rollback sender's saldo during full transfer rollback", zap.NamedError("primary_error", primaryErr), zap.NamedError("rollback_error", rollbackErr))
		}

		receiverSaldo.TotalBalance -= int32(request.TransferAmount)

		if _, rollbackErr := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   receiverSaldo.CardNumber,
			TotalBalance: int(receiverSaldo.TotalBalance),
		}); rollbackErr != nil {
			s.logger.Error("CRITICAL: Failed to rollback receiver's saldo during full transfer rollback", zap.NamedError("primary_error", primaryErr), zap.NamedError("rollback_error", rollbackErr))
		}
		return primaryErr
	}

	transfer, err := s.transferRepository.CreateTransfer(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransferRow](
			s.logger,
			rollbackFullTransfer(transfer_errors.ErrFailedCreateTransfer),
			method,
			span,
		)
	}

	_, err = s.transferRepository.UpdateTransferStatus(ctx, &requests.UpdateTransferStatus{
		TransferID: int(transfer.TransferID),
		Status:     "success",
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransferRow](
			s.logger,
			transfer_errors.ErrFailedUpdateTransfer,
			method,
			span,

			zap.Int("transfer_id", int(transfer.TransferID)),
		)
	}

	logSuccess("Successfully created transaction",
		zap.Int("transfer_id", int(transfer.TransferID)),
		zap.String("transfer_from", request.TransferFrom),
		zap.String("transfer_to", request.TransferTo),
		zap.Int("transfer_amount", request.TransferAmount),
	)

	return transfer, nil
}

func (s *transferService) UpdateTransaction(ctx context.Context, request *requests.UpdateTransferRequest) (*db.UpdateTransferRow, error) {
	const method = "UpdateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transfer_id", *request.TransferID),
		attribute.Int("new_amount", request.TransferAmount))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting update transaction process", zap.Int("transfer_id", *request.TransferID))

	originalTransfer, err := s.transferRepository.FindById(ctx, *request.TransferID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransferRow](
			s.logger,
			transfer_errors.ErrTransferNotFound,
			method,
			span,

			zap.Int("transfer_id", *request.TransferID),
		)
	}

	amountDifference := request.TransferAmount - int(originalTransfer.TransferAmount)

	senderSaldo, err := s.saldoRepository.FindByCardNumber(ctx, originalTransfer.TransferFrom)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransferRow](
			s.logger,
			saldo_errors.ErrFailedSaldoNotFound,
			method,
			span,

			zap.String("card_number", originalTransfer.TransferFrom),
			zap.String("role", "sender"),
		)
	}

	if int(senderSaldo.TotalBalance)+int(originalTransfer.TransferAmount) < (request.TransferAmount) {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransferRow](
			s.logger,
			saldo_errors.ErrFailedInsufficientBalance,
			method,
			span,
			zap.Int("current_balance", int(senderSaldo.TotalBalance)),
			zap.Int("original_amount", int(originalTransfer.TransferAmount)),
			zap.Int("new_amount", request.TransferAmount),
		)
	}

	newSenderBalance := int(senderSaldo.TotalBalance) - amountDifference
	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   senderSaldo.CardNumber,
		TotalBalance: newSenderBalance,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransferRow](
			s.logger,
			saldo_errors.ErrFailedUpdateSaldo,
			method,
			span,

			zap.String("card_number", senderSaldo.CardNumber),
			zap.String("role", "sender"),
		)
	}

	rollbackSenderSaldo := func(primaryErr error) error {
		s.logger.Warn("Attempting to rollback sender's saldo due to an error", zap.String("card_number", senderSaldo.CardNumber), zap.Error(primaryErr))
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   senderSaldo.CardNumber,
			TotalBalance: int(senderSaldo.TotalBalance),
		})
		if rollbackErr != nil {
			s.logger.Error("CRITICAL: Failed to rollback sender's saldo after another error", zap.NamedError("primary_error", primaryErr), zap.NamedError("rollback_error", rollbackErr))
		}
		return primaryErr
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(ctx, originalTransfer.TransferTo)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransferRow](
			s.logger,
			rollbackSenderSaldo(saldo_errors.ErrFailedSaldoNotFound),
			method,
			span,

			zap.String("card_number", originalTransfer.TransferTo),
			zap.String("role", "receiver"),
		)
	}

	newReceiverBalance := int(receiverSaldo.TotalBalance) + amountDifference
	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   receiverSaldo.CardNumber,
		TotalBalance: newReceiverBalance,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransferRow](
			s.logger,
			rollbackSenderSaldo(saldo_errors.ErrFailedUpdateSaldo),
			method,
			span,

			zap.String("card_number", receiverSaldo.CardNumber),
			zap.String("role", "receiver"),
		)
	}

	rollbackFullAdjustment := func(primaryErr error) error {
		s.logger.Warn("Attempting to rollback full transfer adjustment due to an error", zap.Error(primaryErr))

		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   senderSaldo.CardNumber,
			TotalBalance: int(senderSaldo.TotalBalance),
		})
		if rollbackErr != nil {
			s.logger.Error("CRITICAL: Failed to rollback sender's saldo during full rollback", zap.NamedError("primary_error", primaryErr), zap.NamedError("rollback_error", rollbackErr))
		}

		_, rollbackErr = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   receiverSaldo.CardNumber,
			TotalBalance: int(receiverSaldo.TotalBalance),
		})
		if rollbackErr != nil {
			s.logger.Error("CRITICAL: Failed to rollback receiver's saldo during full rollback", zap.NamedError("primary_error", primaryErr), zap.NamedError("rollback_error", rollbackErr))
		}
		return primaryErr
	}

	updatedTransfer, err := s.transferRepository.UpdateTransfer(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransferRow](
			s.logger,
			rollbackFullAdjustment(transfer_errors.ErrFailedUpdateTransfer),
			method,
			span,

			zap.Int("transfer_id", *request.TransferID),
		)
	}

	_, err = s.transferRepository.UpdateTransferStatus(ctx, &requests.UpdateTransferStatus{
		TransferID: int(updatedTransfer.TransferID),
		Status:     "success",
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransferRow](
			s.logger,
			transfer_errors.ErrFailedUpdateTransfer,
			method,
			span,

			zap.Int("transfer_id", int(updatedTransfer.TransferID)),
		)
	}

	logSuccess("Successfully updated transaction",
		zap.Int("transfer_id", int(updatedTransfer.TransferID)),
		zap.Int("new_amount", int(updatedTransfer.TransferAmount)),
		zap.Int("amount_difference", amountDifference),
	)

	return updatedTransfer, nil
}

func (s *transferService) TrashedTransfer(ctx context.Context, transfer_id int) (*db.Transfer, error) {
	const method = "TrashedTransfer"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transfer_id", transfer_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting trashed transfer process", zap.Int("transfer_id", transfer_id))

	res, err := s.transferRepository.TrashedTransfer(ctx, transfer_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Transfer](
			s.logger,
			transfer_errors.ErrFailedTrashedTransfer,
			method,
			span,

			zap.Int("transfer_id", transfer_id),
		)
	}

	logSuccess("Successfully trashed transfer", zap.Int("transfer_id", transfer_id))

	return res, nil
}

func (s *transferService) RestoreTransfer(ctx context.Context, transfer_id int) (*db.Transfer, error) {
	const method = "RestoreTransfer"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transfer_id", transfer_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting restore transfer process", zap.Int("transfer_id", transfer_id))

	res, err := s.transferRepository.RestoreTransfer(ctx, transfer_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Transfer](
			s.logger,
			transfer_errors.ErrFailedRestoreTransfer,
			method,
			span,

			zap.Int("transfer_id", transfer_id),
		)
	}

	logSuccess("Successfully restored transfer", zap.Int("transfer_id", transfer_id))

	return res, nil
}

func (s *transferService) DeleteTransferPermanent(ctx context.Context, transfer_id int) (bool, error) {
	const method = "DeleteTransferPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transfer_id", transfer_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting delete transfer permanent process", zap.Int("transfer_id", transfer_id))

	_, err := s.transferRepository.DeleteTransferPermanent(ctx, transfer_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			transfer_errors.ErrFailedDeleteTransferPermanent,
			method,
			span,

			zap.Int("transfer_id", transfer_id),
		)
	}

	logSuccess("Successfully deleted transfer permanently", zap.Int("transfer_id", transfer_id))

	return true, nil
}

func (s *transferService) RestoreAllTransfer(ctx context.Context) (bool, error) {
	const method = "RestoreAllTransfer"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring all transfers")

	_, err := s.transferRepository.RestoreAllTransfer(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			transfer_errors.ErrFailedRestoreAllTransfers,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all transfers")
	return true, nil
}

func (s *transferService) DeleteAllTransferPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllTransferPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Permanently deleting all transfers")

	_, err := s.transferRepository.DeleteAllTransferPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			transfer_errors.ErrFailedDeleteAllTransfersPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all transfers permanently")
	return true, nil
}
