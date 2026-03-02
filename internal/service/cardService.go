package service

import (
	card_cache "MamangRust/paymentgatewaygrpc/internal/cache/card"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/errorhandler"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/card_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/user_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type cardService struct {
	cardRepository repository.CardRepository
	userRepository repository.UserRepository
	logger         logger.LoggerInterface
	observability  observability.TraceLoggerObservability
	cache          card_cache.CardMencache
}

type CardServiceDeps struct {
	CardRepo      repository.CardRepository
	UserRepo      repository.UserRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
	cache         card_cache.CardMencache
}

func NewCardService(deps CardServiceDeps) CardService {
	return &cardService{
		cardRepository: deps.CardRepo,
		userRepository: deps.UserRepo,
		logger:         deps.Logger,
		observability:  deps.Observability,
		cache:          deps.cache,
	}
}

func (s *cardService) FindAll(ctx context.Context, req *requests.FindAllCards) ([]*db.GetCardsRow, *int, error) {
	const method = "FindAll"

	page := req.Page
	pageSize := req.PageSize

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetFindAllCache(ctx, req); found {
		logSuccess("Successfully fetched card records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	cards, err := s.cardRepository.FindAllCards(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetCardsRow](
			s.logger,
			card_errors.ErrFailedFindAllCards,
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(cards) > 0 {
		totalCount = int(cards[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetFindAllCache(ctx, req, cards, &totalCount)

	logSuccess("Successfully fetched card records",
		zap.Int("totalRecords", totalCount),

		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return cards, &totalCount, nil
}

func (s *cardService) FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*db.GetActiveCardsWithCountRow, *int, error) {
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

	if data, total, found := s.cache.GetByActiveCache(ctx, req); found {
		logSuccess("Successfully fetched active card records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, err := s.cardRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetActiveCardsWithCountRow](
			s.logger,
			card_errors.ErrFailedFindActiveCards,
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

	s.cache.SetByActiveCache(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched active card records",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return res, &totalCount, nil
}

func (s *cardService) FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*db.GetTrashedCardsWithCountRow, *int, error) {
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

	if data, total, found := s.cache.GetByTrashedCache(ctx, req); found {
		logSuccess("Successfully fetched trashed card records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, err := s.cardRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTrashedCardsWithCountRow](
			s.logger,
			card_errors.ErrFailedFindTrashedCards,
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

	s.cache.SetByTrashedCache(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched trashed card records",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return res, &totalCount, nil
}

func (s *cardService) FindById(ctx context.Context, card_id int) (*db.GetCardByIDRow, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("card_id", card_id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetByIdCache(ctx, card_id); found {
		logSuccess("Successfully fetched card from cache", zap.Int("card.id", card_id))
		return data, nil
	}

	res, err := s.cardRepository.FindById(ctx, card_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetCardByIDRow](
			s.logger,
			card_errors.ErrFailedFindById,
			method,
			span,

			zap.Int("card_id", card_id),
		)
	}

	s.cache.SetByIdCache(ctx, card_id, res)

	logSuccess("Successfully fetched card", zap.Int("card_id", card_id))

	return res, nil
}

func (s *cardService) FindByUserID(ctx context.Context, user_id int) (*db.GetCardByUserIDRow, error) {
	const method = "FindByUserID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", user_id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetByUserIDCache(ctx, user_id); found {
		logSuccess("Successfully fetched card records by user ID from cache", zap.Int("user.id", user_id))
		return data, nil
	}

	res, err := s.cardRepository.FindCardByUserId(ctx, user_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetCardByUserIDRow](
			s.logger,
			card_errors.ErrFailedFindByUserID,
			method,
			span,

			zap.Int("user_id", user_id),
		)
	}

	s.cache.SetByUserIDCache(ctx, user_id, res)

	logSuccess("Successfully fetched card records by user ID", zap.Int("user_id", user_id))

	return res, nil
}

func (s *cardService) DashboardCard(ctx context.Context) (*response.DashboardCard, error) {
	const method = "DashboardCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetDashboardCardCache(ctx); found {
		s.logger.Debug("DashboardCard cache hit")
		return data, nil
	}

	totalBalance, err := s.cardRepository.GetTotalBalances(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.DashboardCard](
			s.logger,
			card_errors.ErrFailedFindTotalBalances,
			method,
			span,
		)
	}

	totalTopup, err := s.cardRepository.GetTotalTopAmount(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.DashboardCard](
			s.logger,
			card_errors.ErrFailedFindTotalTopAmount,
			method,
			span,
		)
	}

	totalWithdraw, err := s.cardRepository.GetTotalWithdrawAmount(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.DashboardCard](
			s.logger,
			card_errors.ErrFailedFindTotalWithdrawAmount,
			method,
			span,
		)
	}

	totalTransaction, err := s.cardRepository.GetTotalTransactionAmount(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.DashboardCard](
			s.logger,
			card_errors.ErrFailedFindTotalTransactionAmount,
			method,
			span,
		)
	}

	totalTransfer, err := s.cardRepository.GetTotalTransferAmount(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.DashboardCard](
			s.logger,
			card_errors.ErrFailedFindTotalTransferAmount,
			method,
			span,
		)
	}

	result := &response.DashboardCard{
		TotalBalance:     totalBalance,
		TotalTopup:       totalTopup,
		TotalWithdraw:    totalWithdraw,
		TotalTransaction: totalTransaction,
		TotalTransfer:    totalTransfer,
	}

	s.cache.SetDashboardCardCache(ctx, result)

	logSuccess("Completed DashboardCard service",
		zap.Int64("total_balance", *totalBalance),
		zap.Int64("total_topup", *totalTopup),
		zap.Int64("total_withdraw", *totalWithdraw),
		zap.Int64("total_transaction", *totalTransaction),
		zap.Int64("total_transfer", *totalTransfer),
	)

	return result, nil
}

func (s *cardService) DashboardCardCardNumber(ctx context.Context, cardNumber string) (*response.DashboardCardCardNumber, error) {
	const method = "DashboardCardCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetDashboardCardCardNumberCache(ctx, cardNumber); found {
		s.logger.Debug("DashboardCardCardNumber cache hit", zap.String("card_number", cardNumber))
		return data, nil
	}

	totalBalance, err := s.cardRepository.GetTotalBalanceByCardNumber(ctx, cardNumber)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.DashboardCardCardNumber](
			s.logger,
			card_errors.ErrFailedFindTotalBalanceByCard,
			method,
			span,
			zap.String("card_number", cardNumber),
		)
	}

	totalTopup, err := s.cardRepository.GetTotalTopupAmountByCardNumber(ctx, cardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.DashboardCardCardNumber](
			s.logger,
			card_errors.ErrFailedFindTotalTopupAmountByCard,
			method,
			span,
			zap.String("card_number", cardNumber),
		)
	}

	totalWithdraw, err := s.cardRepository.GetTotalWithdrawAmountByCardNumber(ctx, cardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.DashboardCardCardNumber](
			s.logger,
			card_errors.ErrFailedFindTotalWithdrawAmountByCard,
			method,
			span,
			zap.String("card_number", cardNumber),
		)
	}

	totalTransaction, err := s.cardRepository.GetTotalTransactionAmountByCardNumber(ctx, cardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.DashboardCardCardNumber](
			s.logger,
			card_errors.ErrFailedFindTotalTransactionAmountByCard,
			method,
			span,
			zap.String("card_number", cardNumber),
		)
	}

	totalTransferSent, err := s.cardRepository.GetTotalTransferAmountBySender(ctx, cardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.DashboardCardCardNumber](
			s.logger,
			card_errors.ErrFailedFindTotalTransferAmountBySender,
			method,
			span,
			zap.String("card_number", cardNumber),
		)
	}

	totalTransferReceived, err := s.cardRepository.GetTotalTransferAmountByReceiver(ctx, cardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.DashboardCardCardNumber](
			s.logger,
			card_errors.ErrFailedFindTotalTransferAmountByReceiver,
			method,
			span,
			zap.String("card_number", cardNumber),
		)
	}

	result := &response.DashboardCardCardNumber{
		TotalBalance:          totalBalance,
		TotalTopup:            totalTopup,
		TotalWithdraw:         totalWithdraw,
		TotalTransaction:      totalTransaction,
		TotalTransferSend:     totalTransferSent,
		TotalTransferReceiver: totalTransferReceived,
	}

	s.cache.SetDashboardCardCardNumberCache(ctx, cardNumber, result)

	logSuccess("Completed DashboardCardCardNumber service",
		zap.String("card_number", cardNumber),
		zap.Int64("total_balance", *totalBalance),
		zap.Int64("total_topup", *totalTopup),
		zap.Int64("total_withdraw", *totalWithdraw),
		zap.Int64("total_transaction", *totalTransaction),
		zap.Int64("total_transfer_sent", *totalTransferSent),
		zap.Int64("total_transfer_received", *totalTransferReceived),
	)

	return result, nil
}

func (s *cardService) FindMonthlyBalance(ctx context.Context, year int) ([]*db.GetMonthlyBalancesRow, error) {
	const method = "FindMonthlyBalance"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyBalanceCache(ctx, year); found {
		logSuccess("Monthly balance cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardRepository.GetMonthlyBalance(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyBalancesRow](
			s.logger,
			card_errors.ErrFailedFindMonthlyBalance,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyBalanceCache(ctx, year, res)

	logSuccess("Monthly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindYearlyBalance(ctx context.Context, year int) ([]*db.GetYearlyBalancesRow, error) {
	const method = "FindYearlyBalance"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyBalanceCache(ctx, year); found {
		logSuccess("Yearly balance cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardRepository.GetYearlyBalance(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyBalancesRow](
			s.logger,
			card_errors.ErrFailedFindYearlyBalance,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyBalanceCache(ctx, year, res)

	logSuccess("Yearly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindMonthlyTopupAmount(ctx context.Context, year int) ([]*db.GetMonthlyTopupAmountRow, error) {
	const method = "FindMonthlyTopupAmount"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTopupCache(ctx, year); found {
		logSuccess("Monthly topup amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardRepository.GetMonthlyTopupAmount(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTopupAmountRow](
			s.logger,
			card_errors.ErrFailedFindMonthlyTopupAmount,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyTopupCache(ctx, year, res)

	logSuccess("Monthly topup amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindYearlyTopupAmount(ctx context.Context, year int) ([]*db.GetYearlyTopupAmountRow, error) {
	const method = "FindYearlyTopupAmount"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupCache(ctx, year); found {
		logSuccess("Yearly topup amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardRepository.GetYearlyTopupAmount(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTopupAmountRow](
			s.logger,
			card_errors.ErrFailedFindYearlyTopupAmount,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyTopupCache(ctx, year, res)

	logSuccess("Yearly topup amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindMonthlyWithdrawAmount(ctx context.Context, year int) ([]*db.GetMonthlyWithdrawAmountRow, error) {
	const method = "FindMonthlyWithdrawAmount"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyWithdrawCache(ctx, year); found {
		logSuccess("Monthly withdraw amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardRepository.GetMonthlyWithdrawAmount(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyWithdrawAmountRow](
			s.logger,
			card_errors.ErrFailedFindMonthlyWithdrawAmount,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyWithdrawCache(ctx, year, res)

	logSuccess("Monthly withdraw amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindYearlyWithdrawAmount(ctx context.Context, year int) ([]*db.GetYearlyWithdrawAmountRow, error) {
	const method = "FindYearlyWithdrawAmount"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyWithdrawCache(ctx, year); found {
		logSuccess("Yearly withdraw amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardRepository.GetYearlyWithdrawAmount(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyWithdrawAmountRow](
			s.logger,
			card_errors.ErrFailedFindYearlyWithdrawAmount,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyWithdrawCache(ctx, year, res)

	logSuccess("Yearly withdraw amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindMonthlyTransactionAmount(ctx context.Context, year int) ([]*db.GetMonthlyTransactionAmountRow, error) {
	const method = "FindMonthlyTransactionAmount"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransactionCache(ctx, year); found {
		logSuccess("Monthly transaction amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardRepository.GetMonthlyTransactionAmount(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransactionAmountRow](
			s.logger,
			card_errors.ErrFailedFindMonthlyTransactionAmount,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyTransactionCache(ctx, year, res)

	logSuccess("Monthly transaction amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindYearlyTransactionAmount(ctx context.Context, year int) ([]*db.GetYearlyTransactionAmountRow, error) {
	const method = "FindYearlyTransactionAmount"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransactionCache(ctx, year); found {
		logSuccess("Yearly transaction amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardRepository.GetYearlyTransactionAmount(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionAmountRow](
			s.logger,
			card_errors.ErrFailedFindYearlyTransactionAmount,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyTransactionCache(ctx, year, res)

	logSuccess("Yearly transaction amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindMonthlyTransferAmountSender(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountSenderRow, error) {
	const method = "FindMonthlyTransferAmountSender"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransferSenderCache(ctx, year); found {
		logSuccess("Monthly transfer sender amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardRepository.GetMonthlyTransferAmountSender(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransferAmountSenderRow](
			s.logger,
			card_errors.ErrFailedFindMonthlyTransferAmountSender,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyTransferSenderCache(ctx, year, res)

	logSuccess("Monthly transfer sender amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindYearlyTransferAmountSender(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountSenderRow, error) {
	const method = "FindYearlyTransferAmountSender"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransferSenderCache(ctx, year); found {
		logSuccess("Yearly transfer sender amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardRepository.GetYearlyTransferAmountSender(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransferAmountSenderRow](
			s.logger,
			card_errors.ErrFailedFindYearlyTransferAmountSender,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyTransferSenderCache(ctx, year, res)

	logSuccess("Yearly transfer sender amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindMonthlyTransferAmountReceiver(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountReceiverRow, error) {
	const method = "FindMonthlyTransferAmountReceiver"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransferReceiverCache(ctx, year); found {
		logSuccess("Monthly transfer receiver amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardRepository.GetMonthlyTransferAmountReceiver(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransferAmountReceiverRow](
			s.logger,
			card_errors.ErrFailedFindMonthlyTransferAmountReceiver,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyTransferReceiverCache(ctx, year, res)

	logSuccess("Monthly transfer receiver amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindYearlyTransferAmountReceiver(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountReceiverRow, error) {
	const method = "FindYearlyTransferAmountReceiver"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransferReceiverCache(ctx, year); found {
		logSuccess("Yearly transfer receiver amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardRepository.GetYearlyTransferAmountReceiver(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransferAmountReceiverRow](
			s.logger,
			card_errors.ErrFailedFindYearlyTransferAmountReceiver,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyTransferReceiverCache(ctx, year, res)

	logSuccess("Yearly transfer receiver amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindMonthlyBalancesByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyBalancesByCardNumberRow, error) {
	const method = "FindMonthlyBalancesByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.String("card_number", req.CardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyBalanceByNumberCache(ctx, req); found {
		logSuccess("Cache hit for monthly balance card", zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardRepository.GetMonthlyBalancesByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyBalancesByCardNumberRow](
			s.logger,
			card_errors.ErrFailedFindMonthlyBalanceByCard,
			method,
			span,
			zap.Int("year", req.Year),
			zap.String("card_number", req.CardNumber),
		)
	}

	s.cache.SetMonthlyBalanceByNumberCache(ctx, req, res)

	logSuccess("Monthly balance retrieved successfully",
		zap.Int("year", req.Year),
		zap.String("card_number", req.CardNumber),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindYearlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyBalancesByCardNumberRow, error) {
	const method = "FindYearlyBalanceByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.String("card_number", req.CardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyBalanceByNumberCache(ctx, req); found {
		logSuccess("Cache hit for yearly balance card", zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardRepository.GetYearlyBalanceByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyBalancesByCardNumberRow](
			s.logger,
			card_errors.ErrFailedFindYearlyBalanceByCard,
			method,
			span,
			zap.Int("year", req.Year),
			zap.String("card_number", req.CardNumber),
		)
	}

	s.cache.SetYearlyBalanceByNumberCache(ctx, req, res)

	logSuccess("Yearly balance retrieved successfully",
		zap.Int("year", req.Year),
		zap.String("card_number", req.CardNumber),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindMonthlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTopupAmountByCardNumberRow, error) {
	const method = "FindMonthlyTopupAmountByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTopupByNumberCache(ctx, req); found {
		logSuccess("Cache hit for monthly topup amount card", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardRepository.GetMonthlyTopupAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTopupAmountByCardNumberRow](
			s.logger,
			card_errors.ErrFailedFindMonthlyTopupAmountByCard,
			method,
			span,
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyTopupByNumberCache(ctx, req, res)

	logSuccess("Monthly topup amount by card number retrieved successfully",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindYearlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTopupAmountByCardNumberRow, error) {
	const method = "FindYearlyTopupAmountByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupByNumberCache(ctx, req); found {
		logSuccess("Cache hit for yearly topup amount card", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardRepository.GetYearlyTopupAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTopupAmountByCardNumberRow](
			s.logger,
			card_errors.ErrFailedFindYearlyTopupAmountByCard,
			method,
			span,
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyTopupByNumberCache(ctx, req, res)

	logSuccess("Yearly topup amount by card number retrieved successfully",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyWithdrawAmountByCardNumberRow, error) {
	const method = "FindMonthlyWithdrawAmountByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyWithdrawByNumberCache(ctx, req); found {
		logSuccess("Cache hit for monthly withdraw amount card", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardRepository.GetMonthlyWithdrawAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyWithdrawAmountByCardNumberRow](
			s.logger,
			card_errors.ErrFailedFindMonthlyWithdrawAmountByCard,
			method,
			span,
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyWithdrawByNumberCache(ctx, req, res)

	logSuccess("Monthly withdraw amount by card number retrieved successfully",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindYearlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyWithdrawAmountByCardNumberRow, error) {
	const method = "FindYearlyWithdrawAmountByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyWithdrawByNumberCache(ctx, req); found {
		logSuccess("Cache hit for yearly withdraw amount card", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardRepository.GetYearlyWithdrawAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyWithdrawAmountByCardNumberRow](
			s.logger,
			card_errors.ErrFailedFindYearlyWithdrawAmountByCard,
			method,
			span,
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyWithdrawByNumberCache(ctx, req, res)

	logSuccess("Yearly withdraw amount by card number retrieved successfully",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindMonthlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransactionAmountByCardNumberRow, error) {
	const method = "FindMonthlyTransactionAmountByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransactionByNumberCache(ctx, req); found {
		logSuccess("Cache hit for monthly transaction amount card", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardRepository.GetMonthlyTransactionAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransactionAmountByCardNumberRow](
			s.logger,
			card_errors.ErrFailedFindMonthlyTransactionAmountByCard,
			method,
			span,
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyTransactionByNumberCache(ctx, req, res)

	logSuccess("Monthly transaction amount by card number retrieved successfully",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindYearlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransactionAmountByCardNumberRow, error) {
	const method = "FindYearlyTransactionAmountByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransactionByNumberCache(ctx, req); found {
		logSuccess("Cache hit for yearly transaction amount card", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardRepository.GetYearlyTransactionAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionAmountByCardNumberRow](
			s.logger,
			card_errors.ErrFailedFindYearlyTransactionAmountByCard,
			method,
			span,
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyTransactionByNumberCache(ctx, req, res)

	logSuccess("Yearly transaction amount by card number retrieved successfully",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindMonthlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransferAmountBySenderRow, error) {
	const method = "FindMonthlyTransferAmountBySender"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransferBySenderCache(ctx, req); found {
		logSuccess("Cache hit for monthly transfer sender amount card", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardRepository.GetMonthlyTransferAmountBySender(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransferAmountBySenderRow](
			s.logger,
			card_errors.ErrFailedFindMonthlyTransferAmountBySender,
			method,
			span,
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyTransferBySenderCache(ctx, req, res)

	logSuccess("Monthly transfer sender amount by card number retrieved successfully",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindYearlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransferAmountBySenderRow, error) {
	const method = "FindYearlyTransferAmountBySender"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransferBySenderCache(ctx, req); found {
		logSuccess("Cache hit for yearly transfer sender amount card", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardRepository.GetYearlyTransferAmountBySender(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransferAmountBySenderRow](
			s.logger,
			card_errors.ErrFailedFindYearlyTransferAmountBySender,
			method,
			span,
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyTransferBySenderCache(ctx, req, res)

	logSuccess("Yearly transfer sender amount by card number retrieved successfully",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindMonthlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransferAmountByReceiverRow, error) {
	const method = "FindMonthlyTransferAmountByReceiver"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransferByReceiverCache(ctx, req); found {
		logSuccess("Cache hit for monthly transfer receiver amount card", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardRepository.GetMonthlyTransferAmountByReceiver(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransferAmountByReceiverRow](
			s.logger,
			card_errors.ErrFailedFindMonthlyTransferAmountByReceiver,
			method,
			span,
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetMonthlyTransferByReceiverCache(ctx, req, res)

	logSuccess("Monthly transfer receiver amount by card number retrieved successfully",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindYearlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransferAmountByReceiverRow, error) {
	const method = "FindYearlyTransferAmountByReceiver"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransferByReceiverCache(ctx, req); found {
		logSuccess("Cache hit for yearly transfer receiver amount card", zap.String("card_number", req.CardNumber), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardRepository.GetYearlyTransferAmountByReceiver(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransferAmountByReceiverRow](
			s.logger,
			card_errors.ErrFailedFindYearlyTransferAmountByReceiver,
			method,
			span,
			zap.String("card_number", req.CardNumber),
			zap.Int("year", req.Year),
		)
	}

	s.cache.SetYearlyTransferByReceiverCache(ctx, req, res)

	logSuccess("Yearly transfer receiver amount by card number retrieved successfully",
		zap.String("card_number", req.CardNumber),
		zap.Int("year", req.Year),
		zap.Int("result_count", len(res)),
	)

	return res, nil
}

func (s *cardService) FindByCardNumber(ctx context.Context, card_number string) (*db.GetCardByCardNumberRow, error) {
	const method = "FindByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetByCardNumberCache(ctx, card_number); found {
		logSuccess("Successfully fetched card record by card number from cache", zap.String("card_number", card_number))
		return data, nil
	}

	res, err := s.cardRepository.FindCardByCardNumber(ctx, card_number)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetCardByCardNumberRow](
			s.logger,
			card_errors.ErrCardNotFoundRes,
			method,
			span,
			zap.String("card_number", card_number),
		)
	}

	s.cache.SetByCardNumberCache(ctx, card_number, res)

	logSuccess("Successfully fetched card record by card number", zap.String("card_number", card_number))

	return res, nil
}

func (s *cardService) CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*db.CreateCardRow, error) {
	const method = "CreateCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", request.UserID))

	defer func() {
		end(status)
	}()

	_, err := s.userRepository.FindById(ctx, request.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateCardRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,

			zap.Int("user_id", request.UserID),
		)
	}

	res, err := s.cardRepository.CreateCard(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateCardRow](
			s.logger,
			card_errors.ErrFailedCreateCard,
			method,
			span,
		)
	}

	logSuccess("Successfully created new card", zap.Int("card_id", int(res.CardID)))

	return res, nil
}

func (s *cardService) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*db.UpdateCardRow, error) {
	const method = "UpdateCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("card_id", request.CardID),
		attribute.Int("user_id", request.UserID))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Updating card", zap.Int("card_id", request.CardID), zap.Any("request", request))

	_, err := s.userRepository.FindById(ctx, request.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateCardRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,

			zap.Int("user_id", request.UserID),
		)
	}

	res, err := s.cardRepository.UpdateCard(ctx, request)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateCardRow](
			s.logger,
			card_errors.ErrFailedUpdateCard,
			method,
			span,

			zap.Int("cardID", request.CardID),
		)
	}

	s.cache.DeleteCardCommandCache(ctx, request.CardID)

	logSuccess("Successfully updated card", zap.Int("cardID", int(res.CardID)))

	return res, nil
}

func (s *cardService) TrashedCard(ctx context.Context, card_id int) (*db.Card, error) {
	const method = "TrashedCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("card_id", card_id))

	defer func() {
		end(status)
	}()

	res, err := s.cardRepository.TrashedCard(ctx, card_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Card](
			s.logger,
			card_errors.ErrFailedTrashCard,
			method,
			span,

			zap.Int("card_id", card_id),
		)
	}

	s.cache.DeleteCardCommandCache(ctx, card_id)

	logSuccess("Successfully trashed card", zap.Int("card_id", int(res.CardID)))

	return res, nil
}

func (s *cardService) RestoreCard(ctx context.Context, card_id int) (*db.Card, error) {
	const method = "RestoreCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("card_id", card_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring card", zap.Int("card_id", card_id))

	res, err := s.cardRepository.RestoreCard(ctx, card_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Card](
			s.logger,
			card_errors.ErrFailedRestoreCard,
			method,
			span,

			zap.Int("card_id", card_id),
		)
	}

	s.cache.DeleteCardCommandCache(ctx, card_id)

	logSuccess("Successfully restored card", zap.Int("card_id", int(res.CardID)))

	return res, nil
}

func (s *cardService) DeleteCardPermanent(ctx context.Context, card_id int) (bool, error) {
	const method = "DeleteCardPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("card_id", card_id))

	defer func() {
		end(status)
	}()

	_, err := s.cardRepository.DeleteCardPermanent(ctx, card_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			card_errors.ErrFailedDeleteCard,
			method,
			span,

			zap.Int("card_id", card_id),
		)
	}

	s.cache.DeleteCardCommandCache(ctx, card_id)

	logSuccess("Successfully deleted card permanently", zap.Int("card_id", card_id))

	return true, nil
}

func (s *cardService) RestoreAllCard(ctx context.Context) (bool, error) {
	const method = "RestoreAllCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.cardRepository.RestoreAllCard(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			card_errors.ErrFailedRestoreAllCards,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all cards")

	return true, nil
}

func (s *cardService) DeleteAllCardPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllCardPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.cardRepository.DeleteAllCardPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			card_errors.ErrFailedDeleteAllCards,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all cards permanently")

	return true, nil
}
