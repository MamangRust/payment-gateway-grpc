package service

import (
	transaction_cache "MamangRust/paymentgatewaygrpc/internal/cache/transaction"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/errorhandler"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/card_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/merchant_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/transaction_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
	"context"
	"errors"
	"fmt"

	"github.com/grafana/pyroscope-go"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transactionService struct {
	merchantRepository    repository.MerchantRepository
	cardRepository        repository.CardRepository
	saldoRepository       repository.SaldoRepository
	transactionRepository repository.TransactionRepository
	cache                 transaction_cache.TransactionMencache
	logger                logger.LoggerInterface
	observability         observability.TraceLoggerObservability
}

type TransactionServiceDeps struct {
	MerchantRepo    repository.MerchantRepository
	CardRepo        repository.CardRepository
	SaldoRepo       repository.SaldoRepository
	TransactionRepo repository.TransactionRepository
	Cache           transaction_cache.TransactionMencache
	Logger          logger.LoggerInterface
	Observability   observability.TraceLoggerObservability
}

func NewTransactionService(deps TransactionServiceDeps) *transactionService {
	return &transactionService{
		merchantRepository:    deps.MerchantRepo,
		cardRepository:        deps.CardRepo,
		saldoRepository:       deps.SaldoRepo,
		transactionRepository: deps.TransactionRepo,
		logger:                deps.Logger,
		cache:                 deps.Cache,
		observability:         deps.Observability,
	}
}

func (s *transactionService) FindAll(ctx context.Context, req *requests.FindAllTransactions) ([]*db.GetTransactionsRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedTransactionsCache(ctx, req); found {
		logSuccess("Successfully fetched card records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	var transactions []*db.GetTransactionsRow
	var err error

	pyroscope.TagWrapper(ctx, pyroscope.Labels(
		"page", fmt.Sprint(page),
		"pageSize", fmt.Sprint(pageSize),
		"search", search,
	), func(ctx context.Context) {
		transactions, err = s.transactionRepository.FindAllTransactions(ctx, req)
	})

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTransactionsRow](
			s.logger,
			transaction_errors.ErrFailedFindAllTransactions,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTransactionsCache(ctx, req, transactions, &totalCount)

	logSuccess("Successfully fetched transaction",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return transactions, &totalCount, nil
}

func (s *transactionService) FindAllByCardNumber(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*db.GetTransactionsByCardNumberRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedTransactionByCardNumberCache(ctx, req); found {
		logSuccess("Successfully retrieved all transaction records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, err := s.transactionRepository.FindAllTransactionByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTransactionsByCardNumberRow](
			s.logger,
			transaction_errors.ErrFailedFindAllByCardNumber,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("card_number", req.CardNumber),
		)
	}

	var totalCount int

	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTransactionByCardNumberCache(ctx, req, transactions, &totalCount)

	logSuccess("Successfully fetched transaction by card number",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("card_number", req.CardNumber))

	return transactions, &totalCount, nil
}

func (s *transactionService) FindById(ctx context.Context, transactionID int) (*db.GetTransactionByIDRow, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transactionID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedTransactionCache(ctx, transactionID); found {
		logSuccess("Successfully fetched transaction from cache", zap.Int("transaction.id", transactionID))
		return data, nil
	}

	transaction, err := s.transactionRepository.FindById(ctx, transactionID)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetTransactionByIDRow](
			s.logger,
			transaction_errors.ErrTransactionNotFound,
			method,
			span,

			zap.Int("transaction_id", transactionID),
		)
	}

	s.cache.SetCachedTransactionCache(ctx, transaction)

	logSuccess("Successfully fetched transaction", zap.Int("transaction_id", transactionID))

	return transaction, nil
}

func (s *transactionService) FindByActive(ctx context.Context, req *requests.FindAllTransactions) ([]*db.GetActiveTransactionsRow, *int, error) {
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
		attribute.Int("page", req.Page),
		attribute.Int("pageSize", req.PageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedTransactionActiveCache(ctx, req); found {
		logSuccess("Successfully fetched active transaction from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, err := s.transactionRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetActiveTransactionsRow](
			s.logger,
			transaction_errors.ErrFailedFindByActiveTransactions,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTransactionActiveCache(ctx, req, transactions, &totalCount)

	logSuccess("Successfully fetched active transaction",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return transactions, &totalCount, nil
}

func (s *transactionService) FindByTrashed(ctx context.Context, req *requests.FindAllTransactions) ([]*db.GetTrashedTransactionsRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedTransactionTrashedCache(ctx, req); found {
		logSuccess("Successfully fetched trashed transaction from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, err := s.transactionRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTrashedTransactionsRow](
			s.logger,
			transaction_errors.ErrFailedFindByTrashedTransactions,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTransactionTrashedCache(ctx, req, transactions, &totalCount)

	logSuccess("Successfully fetched trashed transaction",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return transactions, &totalCount, nil
}

func (s *transactionService) FindTransactionByMerchantId(ctx context.Context, merchant_id int) ([]*db.GetTransactionsByMerchantIDRow, error) {
	const method = "FindTransactionByMerchantId"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", merchant_id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedTransactionByMerchantIdCache(ctx, merchant_id); found {
		logSuccess("Successfully fetched transaction by merchant ID from cache", zap.Int("merchant.id", merchant_id))
		return data, nil
	}

	res, err := s.transactionRepository.FindTransactionByMerchantId(ctx, merchant_id)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetTransactionsByMerchantIDRow](
			s.logger,
			transaction_errors.ErrFailedFindByMerchantID,
			method,
			span,

			zap.Int("merchant_id", merchant_id),
		)
	}
	s.cache.SetCachedTransactionByMerchantIdCache(ctx, merchant_id, res)

	logSuccess("Successfully fetched transaction by merchant ID", zap.Int("merchant_id", merchant_id))

	return res, nil
}

func (s *transactionService) FindMonthTransactionStatusSuccess(ctx context.Context, req *requests.MonthStatusTransaction) ([]*db.GetMonthTransactionStatusSuccessRow, error) {
	const method = "FindMonthTransactionStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetMonthTransactionStatusSuccessCache(ctx, req); found {
		logSuccess("Successfully fetched monthly transaction status success (from cache)", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return dbRows, nil
	}

	s.logger.Debug("Cache miss for monthly transaction status success, fetching from DB", zap.Int("year", req.Year), zap.Int("month", req.Month))

	dbRows, err := s.transactionRepository.GetMonthTransactionStatusSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthTransactionStatusSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthTransactionSuccess,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetMonthTransactionStatusSuccessCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly transaction status success (from DB)", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *transactionService) FindYearlyTransactionStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransactionStatusSuccessRow, error) {
	const method = "FindYearlyTransactionStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Checking cache for yearly transaction status success", zap.Int("year", year))
	dbRows, found := s.cache.GetYearTransactionStatusSuccessCache(ctx, year)
	if found {
		s.logger.Info("Cache hit for yearly transaction status success", zap.Int("year", year))
		status = "ok"
		logSuccess("Successfully fetched yearly transaction status success (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetYearlyTransactionStatusSuccess(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionStatusSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindYearTransactionSuccess,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetYearTransactionStatusSuccessCache(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly transaction status success (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *transactionService) FindMonthTransactionStatusFailed(ctx context.Context, req *requests.MonthStatusTransaction) ([]*db.GetMonthTransactionStatusFailedRow, error) {
	const method = "FindMonthTransactionStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetMonthTransactionStatusFailedCache(ctx, req); found {
		logSuccess("Successfully fetched monthly transaction status failed (from cache)", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetMonthTransactionStatusFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthTransactionStatusFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthTransactionFailed,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetMonthTransactionStatusFailedCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly transaction status failed (from DB)", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *transactionService) FindYearlyTransactionStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTransactionStatusFailedRow, error) {
	const method = "FindYearlyTransactionStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetYearTransactionStatusFailedCache(ctx, year); found {
		s.logger.Info("Cache hit for yearly transaction status failed", zap.Int("year", year))
		status = "ok"
		logSuccess("Successfully fetched yearly transaction status failed (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetYearlyTransactionStatusFailed(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionStatusFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindYearTransactionFailed,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetYearTransactionStatusFailedCache(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly transaction status failed (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *transactionService) FindMonthlyPaymentMethods(ctx context.Context, year int) ([]*db.GetMonthlyPaymentMethodsRow, error) {
	const method = "FindMonthlyPaymentMethods"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetMonthlyPaymentMethodsCache(ctx, year); found {
		logSuccess("Successfully fetched monthly payment methods (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetMonthlyPaymentMethods(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyPaymentMethodsRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyPaymentMethods,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyPaymentMethodsCache(ctx, year, dbRows)

	logSuccess("Successfully fetched monthly payment methods (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *transactionService) FindYearlyPaymentMethods(ctx context.Context, year int) ([]*db.GetYearlyPaymentMethodsRow, error) {
	const method = "FindYearlyPaymentMethods"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetYearlyPaymentMethodsCache(ctx, year); found {
		logSuccess("Successfully fetched yearly payment methods (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetYearlyPaymentMethods(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyPaymentMethodsRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyPaymentMethods,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyPaymentMethodsCache(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly payment methods (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *transactionService) FindMonthlyAmounts(ctx context.Context, year int) ([]*db.GetMonthlyAmountsRow, error) {
	const method = "FindMonthlyAmounts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetMonthlyAmountsCache(ctx, year); found {
		s.logger.Info("Cache hit for monthly amounts", zap.Int("year", year))
		status = "ok"
		logSuccess("Successfully fetched monthly amounts (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetMonthlyAmounts(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountsRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyAmounts,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyAmountsCache(ctx, year, dbRows)

	logSuccess("Successfully fetched monthly amounts (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *transactionService) FindYearlyAmounts(ctx context.Context, year int) ([]*db.GetYearlyAmountsRow, error) {
	const method = "FindYearlyAmounts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetYearlyAmountsCache(ctx, year); found {
		s.logger.Info("Cache hit for yearly amounts", zap.Int("year", year))
		status = "ok"
		logSuccess("Successfully fetched yearly amounts (from cache)", zap.Int("year", year))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetYearlyAmounts(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountsRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyAmounts,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyAmountsCache(ctx, year, dbRows)

	logSuccess("Successfully fetched yearly amounts (from DB)", zap.Int("year", year))

	return dbRows, nil
}

func (s *transactionService) FindMonthTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*db.GetMonthTransactionStatusSuccessCardNumberRow, error) {
	const method = "FindMonthTransactionStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetMonthTransactionStatusSuccessByCardCache(ctx, req); found {
		logSuccess("Successfully fetched monthly transaction status success by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetMonthTransactionStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthTransactionStatusSuccessCardNumberRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthTransactionSuccessByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetMonthTransactionStatusSuccessByCardCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly transaction status success by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *transactionService) FindYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*db.GetYearlyTransactionStatusSuccessCardNumberRow, error) {
	const method = "FindYearlyTransactionStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetYearTransactionStatusSuccessByCardCache(ctx, req); found {
		logSuccess("Successfully fetched yearly transaction status success by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetYearlyTransactionStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionStatusSuccessCardNumberRow](
			s.logger,
			transaction_errors.ErrFailedFindYearTransactionSuccessByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearTransactionStatusSuccessByCardCache(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly transaction status success by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *transactionService) FindMonthTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*db.GetMonthTransactionStatusFailedCardNumberRow, error) {
	const method = "FindMonthTransactionStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetMonthTransactionStatusFailedByCardCache(ctx, req); found {
		logSuccess("Successfully fetched monthly transaction status failed by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetMonthTransactionStatusFailedByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthTransactionStatusFailedCardNumberRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthTransactionFailedByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetMonthTransactionStatusFailedByCardCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly transaction status failed by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return dbRows, nil
}

func (s *transactionService) FindYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*db.GetYearlyTransactionStatusFailedCardNumberRow, error) {
	const method = "FindYearlyTransactionStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetYearTransactionStatusFailedByCardCache(ctx, req); found {
		logSuccess("Successfully fetched yearly transaction status failed by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetYearlyTransactionStatusFailedByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionStatusFailedCardNumberRow](
			s.logger,
			transaction_errors.ErrFailedFindYearTransactionFailedByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearTransactionStatusFailedByCardCache(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly transaction status failed by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *transactionService) FindMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetMonthlyPaymentMethodsByCardNumberRow, error) {
	const method = "FindMonthlyPaymentMethodsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetMonthlyPaymentMethodsByCardCache(ctx, req); found {
		logSuccess("Successfully fetched monthly payment methods by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetMonthlyPaymentMethodsByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyPaymentMethodsByCardNumberRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyPaymentMethodsByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyPaymentMethodsByCardCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly payment methods by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *transactionService) FindYearlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetYearlyPaymentMethodsByCardNumberRow, error) {
	const method = "FindYearlyPaymentMethodsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetYearlyPaymentMethodsByCardCache(ctx, req); found {
		logSuccess("Successfully fetched yearly payment methods by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetYearlyPaymentMethodsByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyPaymentMethodsByCardNumberRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyPaymentMethodsByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyPaymentMethodsByCardCache(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly payment methods by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *transactionService) FindMonthlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetMonthlyAmountsByCardNumberRow, error) {
	const method = "FindMonthlyAmountsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetMonthlyAmountsByCardCache(ctx, req); found {
		logSuccess("Successfully fetched monthly amounts by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetMonthlyAmountsByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountsByCardNumberRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyAmountsByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyAmountsByCardCache(ctx, req, dbRows)

	logSuccess("Successfully fetched monthly amounts by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *transactionService) FindYearlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetYearlyAmountsByCardNumberRow, error) {
	const method = "FindYearlyAmountsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if dbRows, found := s.cache.GetYearlyAmountsByCardCache(ctx, req); found {
		logSuccess("Successfully fetched yearly amounts by card number (from cache)",
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year))
		return dbRows, nil
	}

	dbRows, err := s.transactionRepository.GetYearlyAmountsByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountsByCardNumberRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyAmountsByCard,
			method,
			span,

			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyAmountsByCardCache(ctx, req, dbRows)

	logSuccess("Successfully fetched yearly amounts by card number (from DB)",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year))

	return dbRows, nil
}

func (s *transactionService) Create(ctx context.Context, apiKey string, request *requests.CreateTransactionRequest) (*db.CreateTransactionRow, error) {
	const method = "CreateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("api_key", apiKey),
		attribute.String("card_number", request.CardNumber),
		attribute.Int("amount", request.Amount))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting CreateTransaction process",
		zap.String("apiKey", apiKey),
		zap.Any("request", request),
	)

	merchant, err := s.merchantRepository.FindByApiKey(ctx, apiKey)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			merchant_errors.ErrFailedFindByApiKey,
			method,
			span,

			zap.String("api_key", apiKey),
		)
	}

	card, err := s.cardRepository.FindCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			card_errors.ErrFailedFindByCardNumber,
			method,
			span,

			zap.String("card_number", request.CardNumber),
		)
	}

	saldo, err := s.saldoRepository.FindByCardNumber(ctx, card.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			saldo_errors.ErrFailedSaldoNotFound,
			method,
			span,

			zap.String("card_number", card.CardNumber),
		)
	}

	if int(saldo.TotalBalance) < request.Amount {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			saldo_errors.ErrFailedInsufficientBalance,
			method,
			span,
			zap.Int32("available_balance", saldo.TotalBalance),
			zap.Int("transaction_amount", request.Amount),
		)
	}

	newUserBalance := int(saldo.TotalBalance) - request.Amount
	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: newUserBalance,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			saldo_errors.ErrFailedUpdateSaldo,
			method,
			span,

			zap.String("card_number", card.CardNumber),
		)
	}

	rollbackUserSaldo := func(primaryErr error) error {
		s.logger.Warn("Attempting to rollback user's saldo due to an error", zap.String("card_number", card.CardNumber), zap.Error(primaryErr))
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   card.CardNumber,
			TotalBalance: int(saldo.TotalBalance),
		})
		if rollbackErr != nil {
			s.logger.Error("CRITICAL: Failed to rollback user's saldo after another error", zap.NamedError("primary_error", primaryErr), zap.NamedError("rollback_error", rollbackErr))
		}
		return primaryErr
	}

	merchantID := int(merchant.MerchantID)
	request.MerchantID = &merchantID

	transaction, err := s.transactionRepository.CreateTransaction(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			rollbackUserSaldo(transaction_errors.ErrFailedCreateTransaction),
			method,
			span,
		)
	}

	_, err = s.transactionRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
		TransactionID: int(transaction.TransactionID),
		Status:        "success",
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			transaction_errors.ErrFailedUpdateTransaction,
			method,
			span,

			zap.Int("transaction_id", int(transaction.TransactionID)),
		)
	}

	merchantCard, err := s.cardRepository.FindCardByUserId(ctx, int(merchant.UserID))
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			card_errors.ErrCardNotFoundRes,
			method,
			span,

			zap.Int("merchant_user_id", int(merchant.UserID)),
		)
	}

	merchantSaldo, err := s.saldoRepository.FindByCardNumber(ctx, merchantCard.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			saldo_errors.ErrFailedSaldoNotFound,
			method,
			span,

			zap.String("merchant_card_number", merchantCard.CardNumber),
		)
	}

	newMerchantBalance := int(merchantSaldo.TotalBalance) + request.Amount

	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   merchantCard.CardNumber,
		TotalBalance: newMerchantBalance,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			saldo_errors.ErrFailedUpdateSaldo,
			method,
			span,

			zap.String("merchant_card_number", merchantCard.CardNumber),
		)
	}

	logSuccess("CreateTransaction process completed",
		zap.String("apiKey", apiKey),
		zap.Int("transactionID", int(transaction.TransactionID)),
		zap.Int("debited_amount", request.Amount),
	)

	return transaction, nil
}

func (s *transactionService) Update(ctx context.Context, apiKey string, request *requests.UpdateTransactionRequest) (*db.UpdateTransactionRow, error) {
	const method = "UpdateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("api_key", apiKey),
		attribute.Int("transaction_id", *request.TransactionID),
		attribute.Int("new_amount", request.Amount))

	defer func() {
		end(status)
	}()

	markTransactionAsFailed := func(transactionID int, primaryErr error) error {
		s.logger.Error("Attempting to mark transaction as failed due to an error", zap.Int("transaction_id", transactionID), zap.Error(primaryErr))
		if _, statusErr := s.transactionRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
			TransactionID: transactionID,
			Status:        "failed",
		}); statusErr != nil {
			s.logger.Error("CRITICAL: Failed to update transaction status to 'failed' after another error", zap.NamedError("primary_error", primaryErr), zap.NamedError("status_update_error", statusErr))
		}
		return primaryErr
	}

	originalTransaction, err := s.transactionRepository.FindById(ctx, *request.TransactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			markTransactionAsFailed(*request.TransactionID, transaction_errors.ErrFailedUpdateTransaction),
			method,
			span,

			zap.Int("transaction_id", *request.TransactionID),
		)
	}

	merchant, err := s.merchantRepository.FindByApiKey(ctx, apiKey)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			markTransactionAsFailed(*request.TransactionID, merchant_errors.ErrFailedFindByApiKey),
			method,
			span,

			zap.String("api_key", apiKey),
		)
	}

	if originalTransaction.MerchantID != merchant.MerchantID {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			markTransactionAsFailed(*request.TransactionID, transaction_errors.ErrFailedUpdateTransaction),
			method,
			span,
			zap.Int("merchant_id", int(merchant.MerchantID)),
			zap.Int("transaction_merchant_id", int(originalTransaction.MerchantID)),
		)
	}

	card, err := s.cardRepository.FindCardByCardNumber(ctx, originalTransaction.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			markTransactionAsFailed(*request.TransactionID, card_errors.ErrCardNotFoundRes),
			method,
			span,

			zap.String("card_number", originalTransaction.CardNumber),
		)
	}

	currentSaldo, err := s.saldoRepository.FindByCardNumber(ctx, card.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			markTransactionAsFailed(*request.TransactionID, saldo_errors.ErrFailedSaldoNotFound),
			method,
			span,

			zap.String("card_number", card.CardNumber),
		)
	}

	newFinalBalance := int(currentSaldo.TotalBalance) + int(originalTransaction.Amount) - request.Amount

	if newFinalBalance < 0 {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			markTransactionAsFailed(*request.TransactionID, saldo_errors.ErrFailedInsufficientBalance),
			method,
			span,
			zap.Int("current_balance", int(currentSaldo.TotalBalance)),
			zap.Int("original_amount", int(originalTransaction.Amount)),
			zap.Int("new_amount", request.Amount),
			zap.Int("resulting_balance", newFinalBalance),
		)
	}

	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: newFinalBalance,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			markTransactionAsFailed(*request.TransactionID, saldo_errors.ErrFailedUpdateSaldo),
			method,
			span,

			zap.String("card_number", card.CardNumber),
		)
	}

	rollbackSaldo := func(primaryErr error) error {
		s.logger.Warn("Attempting to rollback saldo due to transaction record update failure", zap.String("card_number", card.CardNumber), zap.Error(primaryErr))
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   card.CardNumber,
			TotalBalance: int(currentSaldo.TotalBalance),
		})
		if rollbackErr != nil {
			s.logger.Error("CRITICAL: Failed to rollback saldo after transaction update failure", zap.NamedError("primary_error", primaryErr), zap.NamedError("rollback_error", rollbackErr))
		}
		return primaryErr
	}

	parsedTime := originalTransaction.TransactionTime

	if originalTransaction.TransactionTime.IsZero() {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			rollbackSaldo(transaction_errors.ErrFailedUpdateTransaction),
			method,
			span,
			zap.Error(errors.New("transaction time is zero")),
			zap.Time("transaction_time", originalTransaction.TransactionTime),
		)
	}

	transactionId := int(originalTransaction.TransactionID)
	merchantId := int(originalTransaction.MerchantID)

	updatedTransaction, err := s.transactionRepository.UpdateTransaction(ctx, &requests.UpdateTransactionRequest{
		TransactionID:   &transactionId,
		CardNumber:      originalTransaction.CardNumber,
		Amount:          request.Amount,
		PaymentMethod:   request.PaymentMethod,
		MerchantID:      &merchantId,
		TransactionTime: parsedTime,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			rollbackSaldo(transaction_errors.ErrFailedUpdateTransaction),
			method,
			span,

			zap.Int("transaction_id", *request.TransactionID),
		)
	}

	_, err = s.transactionRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
		TransactionID: int(updatedTransaction.TransactionID),
		Status:        "success",
	})

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			transaction_errors.ErrFailedUpdateTransaction,
			method,
			span,

			zap.Int("transaction_id", int(updatedTransaction.TransactionID)),
		)
	}

	logSuccess("UpdateTransaction process completed",
		zap.String("apiKey", apiKey),
		zap.Int("transaction_id", int(updatedTransaction.TransactionID)),
		zap.Int("new_amount", int(updatedTransaction.Amount)),
		zap.Int("new_final_balance", newFinalBalance),
	)

	return updatedTransaction, nil
}

func (s *transactionService) TrashedTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	const method = "TrashedTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transaction_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting TrashedTransaction process", zap.Int("transaction_id", transaction_id))

	res, err := s.transactionRepository.TrashedTransaction(ctx, transaction_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedTrashedTransaction,
			method,
			span,

			zap.Int("transaction_id", transaction_id),
		)
	}

	logSuccess("Successfully trashed transaction", zap.Int("transaction_id", transaction_id))

	return res, nil
}

func (s *transactionService) RestoreTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	const method = "RestoreTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transaction_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting RestoreTransaction process", zap.Int("transaction_id", transaction_id))

	res, err := s.transactionRepository.RestoreTransaction(ctx, transaction_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedRestoreTransaction,
			method,
			span,

			zap.Int("transaction_id", transaction_id),
		)
	}

	logSuccess("Successfully restored transaction", zap.Int("transaction_id", transaction_id))

	return res, nil
}

func (s *transactionService) DeleteTransactionPermanent(ctx context.Context, transaction_id int) (bool, error) {
	const method = "DeleteTransactionPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transaction_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting DeleteTransactionPermanent process", zap.Int("transaction_id", transaction_id))

	_, err := s.transactionRepository.DeleteTransactionPermanent(ctx, transaction_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedDeleteTransactionPermanent,
			method,
			span,

			zap.Int("transaction_id", transaction_id),
		)
	}

	logSuccess("Successfully permanently deleted transaction", zap.Int("transaction_id", transaction_id))

	return true, nil
}

func (s *transactionService) RestoreAllTransaction(ctx context.Context) (bool, error) {
	const method = "RestoreAllTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring all transactions")

	_, err := s.transactionRepository.RestoreAllTransaction(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedRestoreAllTransactions,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all transactions")
	return true, nil
}

func (s *transactionService) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllTransactionPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Permanently deleting all transactions")

	_, err := s.transactionRepository.DeleteAllTransactionPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedDeleteAllTransactionsPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all transactions permanently")
	return true, nil
}
