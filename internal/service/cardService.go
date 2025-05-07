package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/service"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/pkg/errors/card_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/user_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

type cardService struct {
	cardRepository repository.CardRepository
	userRepository repository.UserRepository
	logger         logger.LoggerInterface
	mapping        responseservice.CardResponseMapper
}

func NewCardService(
	cardRepository repository.CardRepository,
	userRepository repository.UserRepository,
	logger logger.LoggerInterface,
	mapper responseservice.CardResponseMapper,

) *cardService {
	return &cardService{
		cardRepository: cardRepository,
		userRepository: userRepository,
		logger:         logger,
		mapping:        mapper,
	}
}

func (s *cardService) FindAll(req *requests.FindAllCards) ([]*response.CardResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching card records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	cards, totalRecords, err := s.cardRepository.FindAllCards(req)

	if err != nil {
		s.logger.Error("Failed to fetch card",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, card_errors.ErrFailedFindAllCards
	}

	responseData := s.mapping.ToCardsResponse(cards)

	s.logger.Debug("Successfully fetched card records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return responseData, totalRecords, nil
}

func (s *cardService) FindByActive(req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching active card records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.cardRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to fetch active card records",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, card_errors.ErrFailedFindActiveCards
	}

	responseData := s.mapping.ToCardsResponseDeleteAt(res)

	s.logger.Debug("Successfully fetched active card records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return responseData, totalRecords, nil
}

func (s *cardService) FindByTrashed(req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching trashed card records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.cardRepository.FindByTrashed(req)
	if err != nil {
		s.logger.Error("Failed to fetch trashed card records",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, card_errors.ErrFailedFindTrashedCards
	}

	responseData := s.mapping.ToCardsResponseDeleteAt(res)

	s.logger.Debug("Successfully fetched trashed card records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return responseData, totalRecords, nil
}

func (s *cardService) FindById(card_id int) (*response.CardResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching card by ID", zap.Int("card_id", card_id))

	res, err := s.cardRepository.FindById(card_id)

	if err != nil {
		s.logger.Error("Failed to retrieve Card details",
			zap.Error(err),
			zap.Int("card_id", card_id))

		return nil, card_errors.ErrFailedFindById
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully fetched card", zap.Int("card_id", card_id))

	return so, nil
}

func (s *cardService) FindByUserID(user_id int) (*response.CardResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching card by user ID", zap.Int("user_id", user_id))

	res, err := s.cardRepository.FindCardByUserId(user_id)

	if err != nil {
		s.logger.Error("Failed to retrieve Card details by user",
			zap.Error(err),
			zap.Int("user_id", user_id))

		return nil, card_errors.ErrFailedFindByUserID
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully fetched card records by user ID", zap.Int("user_id", user_id))

	return so, nil
}

func (s *cardService) DashboardCard() (*response.DashboardCard, *response.ErrorResponse) {
	s.logger.Debug("Starting DashboardCard service")

	totalBalance, err := s.cardRepository.GetTotalBalances()
	if err != nil {
		s.logger.Error("Failed to retrieve total balance",
			zap.Error(err),
		)
		return nil, card_errors.ErrFailedFindTotalBalances
	}

	totalTopup, err := s.cardRepository.GetTotalTopAmount()
	if err != nil {
		s.logger.Error("Failed to retrieve total top-up amount",
			zap.Error(err),
		)
		return nil, card_errors.ErrFailedFindTotalTopAmount
	}

	totalWithdraw, err := s.cardRepository.GetTotalWithdrawAmount()
	if err != nil {
		s.logger.Error("Failed to retrieve total withdrawal amount",
			zap.Error(err),
		)
		return nil, card_errors.ErrFailedFindTotalWithdrawAmount
	}

	totalTransaction, err := s.cardRepository.GetTotalTransactionAmount()
	if err != nil {
		s.logger.Error("Failed to retrieve total transaction amount",
			zap.Error(err),
		)
		return nil, card_errors.ErrFailedFindTotalTransactionAmount
	}

	totalTransfer, err := s.cardRepository.GetTotalTransferAmount()
	if err != nil {
		s.logger.Error("Failed to retrieve total transfer amount",
			zap.Error(err),
		)
		return nil, card_errors.ErrFailedFindTotalTransferAmount
	}

	s.logger.Debug("Completed DashboardCard service",
		zap.Int("total_balance", int(*totalBalance)),
		zap.Int("total_topup", int(*totalTopup)),
		zap.Int("total_withdraw", int(*totalWithdraw)),
		zap.Int("total_transaction", int(*totalTransaction)),
		zap.Int("total_transfer", int(*totalTransfer)),
	)

	return &response.DashboardCard{
		TotalBalance:     totalBalance,
		TotalTopup:       totalTopup,
		TotalWithdraw:    totalWithdraw,
		TotalTransaction: totalTransaction,
		TotalTransfer:    totalTransfer,
	}, nil
}

func (s *cardService) DashboardCardCardNumber(cardNumber string) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	s.logger.Debug("Starting DashboardCardCardNumber service",
		zap.String("card_number", cardNumber),
	)

	totalBalance, err := s.cardRepository.GetTotalBalanceByCardNumber(cardNumber)
	if err != nil {
		s.logger.Error("Failed to retrieve total balance for card",
			zap.String("card_number", cardNumber),
			zap.Error(err),
		)
		return nil, card_errors.ErrFailedFindTotalBalanceByCard
	}

	totalTopup, err := s.cardRepository.GetTotalTopupAmountByCardNumber(cardNumber)
	if err != nil {
		s.logger.Error("Failed to retrieve total top-up amount for card",
			zap.String("card_number", cardNumber),
			zap.Error(err),
		)
		return nil, card_errors.ErrFailedFindTotalTopupAmountByCard
	}

	totalWithdraw, err := s.cardRepository.GetTotalWithdrawAmountByCardNumber(cardNumber)
	if err != nil {
		s.logger.Error("Failed to retrieve total withdrawal amount for card",
			zap.String("card_number", cardNumber),
			zap.Error(err),
		)
		return nil, card_errors.ErrFailedFindTotalWithdrawAmountByCard
	}

	totalTransaction, err := s.cardRepository.GetTotalTransactionAmountByCardNumber(cardNumber)
	if err != nil {
		s.logger.Error("Failed to retrieve total transaction amount for card",
			zap.String("card_number", cardNumber),
			zap.Error(err),
		)
		return nil, card_errors.ErrFailedFindTotalTransactionAmountByCard
	}

	totalTransferSent, err := s.cardRepository.GetTotalTransferAmountBySender(cardNumber)
	if err != nil {
		s.logger.Error("Failed to retrieve total transfer amount sent by card",
			zap.String("card_number", cardNumber),
			zap.Error(err),
		)
		return nil, card_errors.ErrFailedFindTotalTransferAmountBySender
	}

	totalTransferReceived, err := s.cardRepository.GetTotalTransferAmountByReceiver(cardNumber)
	if err != nil {
		s.logger.Error("Failed to retrieve total transfer amount received by card",
			zap.String("card_number", cardNumber),
			zap.Error(err),
		)
		return nil, card_errors.ErrFailedFindTotalTransferAmountByReceiver
	}

	s.logger.Debug("Completed DashboardCardCardNumber service",
		zap.String("card_number", cardNumber),
		zap.Int("total_balance", int(*totalBalance)),
		zap.Int("total_topup", int(*totalTopup)),
		zap.Int("total_withdraw", int(*totalWithdraw)),
		zap.Int("total_transaction", int(*totalTransaction)),
		zap.Int("total_transfer_sent", int(*totalTransferSent)),
		zap.Int("total_transfer_received", int(*totalTransferReceived)),
	)

	return &response.DashboardCardCardNumber{
		TotalBalance:          totalBalance,
		TotalTopup:            totalTopup,
		TotalWithdraw:         totalWithdraw,
		TotalTransaction:      totalTransaction,
		TotalTransferSend:     totalTransferSent,
		TotalTransferReceiver: totalTransferReceived,
	}, nil
}
func (s *cardService) FindMonthlyBalance(year int) ([]*response.CardResponseMonthBalance, *response.ErrorResponse) {
	s.logger.Debug("FindMonthlyBalance called", zap.Int("year", year))

	res, err := s.cardRepository.GetMonthlyBalance(year)
	if err != nil {
		s.logger.Error("Failed to retrieve monthly balance",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindMonthlyBalance
	}

	so := s.mapping.ToGetMonthlyBalances(res)

	s.logger.Debug("Monthly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindYearlyBalance(year int) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse) {
	s.logger.Debug("FindYearlyBalance called", zap.Int("year", year))

	res, err := s.cardRepository.GetYearlyBalance(year)
	if err != nil {
		s.logger.Error("Failed to retrieve yearly balance",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindYearlyBalance
	}

	so := s.mapping.ToGetYearlyBalances(res)

	s.logger.Debug("Yearly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindMonthlyTopupAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	s.logger.Debug("FindMonthlyTopupAmount called", zap.Int("year", year))

	res, err := s.cardRepository.GetMonthlyTopupAmount(year)

	if err != nil {
		s.logger.Error("Failed to retrieve monthly topup amount",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindMonthlyTopupAmount
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.logger.Debug("Monthly topup amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindYearlyTopupAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	s.logger.Debug("FindYearlyTopupAmount called", zap.Int("year", year))

	res, err := s.cardRepository.GetYearlyTopupAmount(year)

	if err != nil {
		s.logger.Error("Failed to retrieve yearly topup amount",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindYearlyTopupAmount
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.logger.Debug("Yearly topup amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindMonthlyWithdrawAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	s.logger.Debug("FindMonthlyWithdrawAmount called", zap.Int("year", year))

	res, err := s.cardRepository.GetMonthlyWithdrawAmount(year)

	if err != nil {
		s.logger.Error("Failed to retrieve monthly withdraw amount",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindMonthlyWithdrawAmount
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.logger.Debug("Monthly withdraw amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindYearlyWithdrawAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	s.logger.Debug("FindYearlyWithdrawAmount called", zap.Int("year", year))

	if year <= 0 {
		return nil, &response.ErrorResponse{
			Status:  "invalid_request",
			Message: "Year must be a positive number",
			Code:    http.StatusBadRequest,
		}
	}

	res, err := s.cardRepository.GetYearlyWithdrawAmount(year)

	if err != nil {
		s.logger.Error("Failed to retrieve yearly withdraw amount",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindYearlyWithdrawAmount
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.logger.Debug("Yearly withdraw amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindMonthlyTransactionAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	s.logger.Debug("FindMonthlyTransactionAmount called", zap.Int("year", year))

	if year <= 0 {
		return nil, &response.ErrorResponse{
			Status:  "invalid_request",
			Message: "Year must be a positive number",
			Code:    http.StatusBadRequest,
		}
	}

	res, err := s.cardRepository.GetMonthlyTransactionAmount(year)

	if err != nil {
		s.logger.Error("Failed to retrieve monthly transaction amount",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindMonthlyTransactionAmount
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.logger.Debug("Monthly transaction amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindYearlyTransactionAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	s.logger.Debug("FindYearlyTransactionAmount called", zap.Int("year", year))

	res, err := s.cardRepository.GetYearlyTransactionAmount(year)

	if err != nil {
		s.logger.Error("Failed to retrieve yearly transaction amount",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindYearlyTransactionAmount
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.logger.Debug("Yearly transaction amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindMonthlyTransferAmountSender(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	s.logger.Debug("FindMonthlyTransferAmountSender called", zap.Int("year", year))

	res, err := s.cardRepository.GetMonthlyTransferAmountSender(year)
	if err != nil {
		s.logger.Error("Failed to retrieve monthly transfer sender amount",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindMonthlyTransferAmountSender
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.logger.Debug("Monthly transfer sender amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindYearlyTransferAmountSender(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	s.logger.Debug("FindYearlyTransferAmountSender called", zap.Int("year", year))

	res, err := s.cardRepository.GetYearlyTransferAmountSender(year)

	if err != nil {
		s.logger.Error("Failed to retrieve yearly transfer sender amount",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindYearlyTransferAmountSender
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.logger.Debug("Yearly transfer sender amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindMonthlyTransferAmountReceiver(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	s.logger.Debug("FindMonthlyTransferAmountReceiver called", zap.Int("year", year))

	res, err := s.cardRepository.GetMonthlyTransferAmountReceiver(year)

	if err != nil {
		s.logger.Error("Failed to retrieve monthly transfer receiver amount",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindMonthlyTransferAmountReceiver
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.logger.Debug("Monthly transfer receiver amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindYearlyTransferAmountReceiver(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	s.logger.Debug("FindYearlyTransferAmountReceiver called", zap.Int("year", year))

	res, err := s.cardRepository.GetYearlyTransferAmountReceiver(year)

	if err != nil {
		s.logger.Error("Failed to retrieve yearly transfer receiver amount",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindYearlyTransferAmountReceiver
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.logger.Debug("Yearly transfer receiver amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindMonthlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthBalance, *response.ErrorResponse) {
	year := req.Year

	s.logger.Debug("FindMonthlyBalance called", zap.Int("year", year))

	res, err := s.cardRepository.GetMonthlyBalancesByCardNumber(req)

	if err != nil {
		s.logger.Error("Failed to retrieve monthly balance",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindMonthlyBalanceByCard
	}

	so := s.mapping.ToGetMonthlyBalances(res)

	s.logger.Debug("Monthly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindYearlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse) {
	year := req.Year

	s.logger.Debug("FindYearlyBalance called", zap.Int("year", year))

	res, err := s.cardRepository.GetYearlyBalanceByCardNumber(req)

	if err != nil {
		s.logger.Error("Failed to retrieve yearly balance",
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindYearlyBalanceByCard
	}

	so := s.mapping.ToGetYearlyBalances(res)

	s.logger.Debug("Yearly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindMonthlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("FindMonthlyTopupAmountByCardNumber called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	res, err := s.cardRepository.GetMonthlyTopupAmountByCardNumber(req)

	if err != nil {
		s.logger.Error("Failed to retrieve monthly topup amount by card number",
			zap.String("card_number", cardNumber),
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindMonthlyTopupAmountByCard
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.logger.Debug("Monthly topup amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindYearlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("FindYearlyTopupAmountByCardNumber called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	res, err := s.cardRepository.GetYearlyTopupAmountByCardNumber(req)
	if err != nil {
		s.logger.Error("Failed to retrieve yearly topup amount by card number",
			zap.String("card_number", cardNumber),
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindYearlyTopupAmountByCard
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.logger.Debug("Yearly topup amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindMonthlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("FindMonthlyWithdrawAmountByCardNumber called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	res, err := s.cardRepository.GetMonthlyWithdrawAmountByCardNumber(req)

	if err != nil {
		s.logger.Error("Failed to retrieve monthly withdraw amount by card number",
			zap.String("card_number", cardNumber),
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindMonthlyWithdrawAmountByCard
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.logger.Debug("Monthly withdraw amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindYearlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("FindYearlyWithdrawAmountByCardNumber called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	res, err := s.cardRepository.GetYearlyWithdrawAmountByCardNumber(req)
	if err != nil {
		s.logger.Error("Failed to retrieve yearly withdraw amount by card number",
			zap.String("card_number", cardNumber),
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindYearlyWithdrawAmountByCard
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.logger.Debug("Yearly withdraw amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindMonthlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("FindMonthlyTransactionAmountByCardNumber called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	res, err := s.cardRepository.GetMonthlyTransactionAmountByCardNumber(req)

	if err != nil {
		s.logger.Error("Failed to retrieve monthly transaction amount by card number",
			zap.String("card_number", cardNumber),
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindMonthlyTransactionAmountByCard
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.logger.Debug("Monthly transaction amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindYearlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("FindYearlyTransactionAmountByCardNumber called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	res, err := s.cardRepository.GetYearlyTransactionAmountByCardNumber(req)
	if err != nil {
		s.logger.Error("Failed to retrieve yearly transaction amount by card number",
			zap.String("card_number", cardNumber),
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindYearlyTransactionAmountByCard
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.logger.Debug("Yearly transaction amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindMonthlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("FindMonthlyTransferAmountBySender called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	res, err := s.cardRepository.GetMonthlyTransferAmountBySender(req)

	if err != nil {
		s.logger.Error("Failed to retrieve monthly transfer sender amount by card number",
			zap.String("card_number", cardNumber),
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindMonthlyTransferAmountBySender
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.logger.Debug("Monthly transfer sender amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindYearlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("FindYearlyTransferAmountBySender called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	res, err := s.cardRepository.GetYearlyTransferAmountBySender(req)
	if err != nil {
		s.logger.Error("Failed to retrieve yearly transfer sender amount by card number",
			zap.String("card_number", cardNumber),
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindYearlyTransferAmountBySender
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.logger.Debug("Yearly transfer sender amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindMonthlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("FindMonthlyTransferAmountByReceiver called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	res, err := s.cardRepository.GetMonthlyTransferAmountByReceiver(req)

	if err != nil {
		s.logger.Error("Failed to retrieve monthly transfer receiver amount by card number",
			zap.String("card_number", cardNumber),
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindMonthlyTransferAmountByReceiver
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.logger.Debug("Monthly transfer receiver amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindYearlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("FindYearlyTransferAmountByReceiver called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	res, err := s.cardRepository.GetYearlyTransferAmountByReceiver(req)
	if err != nil {
		s.logger.Error("Failed to retrieve yearly transfer receiver amount by card number",
			zap.String("card_number", cardNumber),
			zap.Int("year", year),
			zap.Error(err),
		)

		return nil, card_errors.ErrFailedFindYearlyTransferAmountByReceiver
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.logger.Debug("Yearly transfer receiver amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardService) FindByCardNumber(card_number string) (*response.CardResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching card record by card number", zap.String("card_number", card_number))

	res, err := s.cardRepository.FindCardByCardNumber(card_number)

	if err != nil {
		s.logger.Error("Failed to retrieve find card",
			zap.Error(err),
			zap.String("card_number", card_number))

		return nil, card_errors.ErrCardNotFoundRes
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully fetched card record by card number", zap.String("card_number", card_number))

	return so, nil
}

func (s *cardService) CreateCard(request *requests.CreateCardRequest) (*response.CardResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new card", zap.Any("request", request))

	_, err := s.userRepository.FindById(request.UserID)

	if err != nil {
		s.logger.Error("Failed to retrieve find user",
			zap.Error(err),
			zap.Int("user_id", request.UserID))

		return nil, user_errors.ErrUserNotFoundRes
	}

	res, err := s.cardRepository.CreateCard(request)

	if err != nil {
		s.logger.Error("Failed to create card", zap.Error(err))
		return nil, card_errors.ErrFailedCreateCard
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully created new card", zap.Int("card_id", so.ID))

	return so, nil
}

func (s *cardService) UpdateCard(request *requests.UpdateCardRequest) (*response.CardResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating card", zap.Int("card_id", request.CardID), zap.Any("request", request))

	_, err := s.userRepository.FindById(request.UserID)

	if err != nil {
		s.logger.Error("Failed to retrieve find user",
			zap.Error(err),
			zap.Int("user_id", request.UserID))

		return nil, user_errors.ErrUserNotFoundRes
	}

	res, err := s.cardRepository.UpdateCard(request)

	if err != nil {
		s.logger.Error("Failed to update card", zap.Error(err), zap.Int("cardID", request.CardID))

		return nil, card_errors.ErrFailedUpdateCard
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully updated card", zap.Int("cardID", so.ID))

	return so, nil
}

func (s *cardService) TrashedCard(card_id int) (*response.CardResponse, *response.ErrorResponse) {
	s.logger.Debug("Trashing card", zap.Int("card_id", card_id))

	res, err := s.cardRepository.TrashedCard(card_id)

	if err != nil {
		s.logger.Error("Failed to move card to trash",
			zap.Error(err),
			zap.Int("card_id", card_id))

		return nil, card_errors.ErrFailedTrashCard
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully trashed card", zap.Int("card_id", so.ID))

	return so, nil
}

func (s *cardService) RestoreCard(card_id int) (*response.CardResponse, *response.ErrorResponse) {
	s.logger.Debug("Restoring card", zap.Int("card_id", card_id))

	res, err := s.cardRepository.RestoreCard(card_id)

	if err != nil {
		s.logger.Error("Failed to restore cashier from trash",
			zap.Error(err),
			zap.Int("card_id", card_id))

		return nil, card_errors.ErrFailedRestoreCard
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully restored card", zap.Int("card_id", so.ID))

	return so, nil
}

func (s *cardService) DeleteCardPermanent(card_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting card", zap.Int("card_id", card_id))

	_, err := s.cardRepository.DeleteCardPermanent(card_id)

	if err != nil {
		s.logger.Error("Failed to permanently delete card",
			zap.Error(err),
			zap.Int("card_id", card_id))

		return false, card_errors.ErrFailedDeleteCard
	}

	s.logger.Debug("Successfully deleted card permanently", zap.Int("card_id", card_id))

	return true, nil
}

func (s *cardService) RestoreAllCard() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all cards")

	_, err := s.cardRepository.RestoreAllCard()

	if err != nil {
		s.logger.Error("Failed to restore all trashed cards",
			zap.Error(err))

		return false, card_errors.ErrFailedRestoreAllCards
	}

	s.logger.Debug("Successfully restored all cards")
	return true, nil
}

func (s *cardService) DeleteAllCardPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all cards")

	_, err := s.cardRepository.DeleteAllCardPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed card",
			zap.Error(err))
		return false, card_errors.ErrFailedDeleteAllCards
	}

	s.logger.Debug("Successfully deleted all cards permanently")

	return true, nil
}
