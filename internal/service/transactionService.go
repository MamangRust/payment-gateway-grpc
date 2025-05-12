package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/service"
	"net/http"

	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/pkg/errors/card_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/merchant_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/transaction_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"time"

	"go.uber.org/zap"
)

type transactionService struct {
	merchantRepository    repository.MerchantRepository
	cardRepository        repository.CardRepository
	saldoRepository       repository.SaldoRepository
	transactionRepository repository.TransactionRepository
	logger                logger.LoggerInterface
	mapping               responseservice.TransactionResponseMapper
}

func NewTransactionService(
	merchantRepository repository.MerchantRepository,
	cardRepository repository.CardRepository,
	saldoRepository repository.SaldoRepository,
	transactionRepository repository.TransactionRepository,
	logger logger.LoggerInterface,
	mapping responseservice.TransactionResponseMapper,
) *transactionService {
	return &transactionService{
		merchantRepository:    merchantRepository,
		cardRepository:        cardRepository,
		saldoRepository:       saldoRepository,
		transactionRepository: transactionRepository,
		logger:                logger,
		mapping:               mapping,
	}
}

func (s *transactionService) FindAll(req *requests.FindAllTransactions) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching transaction",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionRepository.FindAllTransactions(req)

	if err != nil {
		s.logger.Error("Failed to fetch transaction",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, transaction_errors.ErrFailedFindAllTransactions
	}

	so := s.mapping.ToTransactionsResponse(transactions)

	s.logger.Debug("Successfully fetched transaction",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionService) FindAllByCardNumber(req *requests.FindAllTransactionCardNumber) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching transaction",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionRepository.FindAllTransactionByCardNumber(req)

	if err != nil {
		s.logger.Error("Failed to fetch transaction",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, transaction_errors.ErrFailedFindAllByCardNumber
	}

	so := s.mapping.ToTransactionsResponse(transactions)

	s.logger.Debug("Successfully fetched transaction",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionService) FindById(transactionID int) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching transaction by ID", zap.Int("transaction_id", transactionID))

	transaction, err := s.transactionRepository.FindById(transactionID)

	if err != nil {
		s.logger.Error("failed to find transaction", zap.Error(err))

		return nil, transaction_errors.ErrTransactionNotFound
	}

	so := s.mapping.ToTransactionResponse(transaction)

	s.logger.Debug("Successfully fetched transaction", zap.Int("transaction_id", transactionID))

	return so, nil
}

func (s *transactionService) FindMonthTransactionStatusSuccess(req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse) {
	month := req.Month
	year := req.Year

	s.logger.Debug("Fetching monthly Transaction status success", zap.Int("year", year), zap.Int("month", month))

	records, err := s.transactionRepository.GetMonthTransactionStatusSuccess(req)
	if err != nil {
		s.logger.Error("failed to fetch monthly Transaction status success", zap.Error(err))

		return nil, transaction_errors.ErrFailedFindMonthTransactionSuccess
	}

	s.logger.Debug("Successfully fetched monthly Transaction status success", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTransactionResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *transactionService) FindYearlyTransactionStatusSuccess(year int) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly Transaction status success", zap.Int("year", year))

	records, err := s.transactionRepository.GetYearlyTransactionStatusSuccess(year)
	if err != nil {
		s.logger.Error("failed to fetch yearly Transaction status success", zap.Error(err))

		return nil, transaction_errors.ErrFailedFindYearTransactionSuccess
	}

	s.logger.Debug("Successfully fetched yearly Transaction status success", zap.Int("year", year))

	so := s.mapping.ToTransactionResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *transactionService) FindMonthTransactionStatusFailed(req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	s.logger.Debug("Fetching monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month))

	records, err := s.transactionRepository.GetMonthTransactionStatusFailed(req)
	if err != nil {
		s.logger.Error("failed to fetch monthly Transaction status Failed", zap.Error(err))

		return nil, transaction_errors.ErrFailedFindMonthTransactionFailed
	}

	s.logger.Debug("Failedfully fetched monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTransactionResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *transactionService) FindYearlyTransactionStatusFailed(year int) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly Transaction status Failed", zap.Int("year", year))

	records, err := s.transactionRepository.GetYearlyTransactionStatusFailed(year)
	if err != nil {
		s.logger.Error("failed to fetch yearly Transaction status Failed", zap.Error(err))

		return nil, transaction_errors.ErrFailedFindYearTransactionFailed
	}

	s.logger.Debug("Failedfully fetched yearly Transaction status Failed", zap.Int("year", year))

	so := s.mapping.ToTransactionResponsesYearStatusFailed(records)

	return so, nil
}

func (s *transactionService) FindMonthTransactionStatusSuccessByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	s.logger.Debug("Fetching monthly Transaction status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.transactionRepository.GetMonthTransactionStatusSuccessByCardNumber(req)

	if err != nil {
		s.logger.Error("failed to fetch monthly Transaction status success", zap.Error(err))

		return nil, transaction_errors.ErrFailedFindMonthTransactionSuccessByCard
	}

	s.logger.Debug("Successfully fetched monthly Transaction status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToTransactionResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *transactionService) FindYearlyTransactionStatusSuccessByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching yearly Transaction status success", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.transactionRepository.GetYearlyTransactionStatusSuccessByCardNumber(req)

	if err != nil {
		s.logger.Error("failed to fetch yearly Transaction status success", zap.Error(err))

		return nil, transaction_errors.ErrFailedFindYearTransactionSuccessByCard
	}

	s.logger.Debug("Successfully fetched yearly Transaction status success", zap.Int("year", year), zap.String("card_number", card_number))

	so := s.mapping.ToTransactionResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *transactionService) FindMonthTransactionStatusFailedByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	s.logger.Debug("Fetching monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.transactionRepository.GetMonthTransactionStatusFailedByCardNumber(req)
	if err != nil {
		s.logger.Error("failed to fetch monthly Transaction status Failed", zap.Error(err))

		return nil, transaction_errors.ErrFailedFindMonthTransactionFailedByCard
	}

	s.logger.Debug("Failedfully fetched monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToTransactionResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *transactionService) FindYearlyTransactionStatusFailedByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching yearly Transaction status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.transactionRepository.GetYearlyTransactionStatusFailedByCardNumber(req)

	if err != nil {
		s.logger.Error("failed to fetch yearly Transaction status Failed", zap.Error(err))

		return nil, transaction_errors.ErrFailedFindYearTransactionFailedByCard
	}

	s.logger.Debug("Failedfully fetched yearly Transaction status Failed", zap.Int("year", year))

	so := s.mapping.ToTransactionResponsesYearStatusFailed(records)

	return so, nil
}

func (s *transactionService) FindMonthlyPaymentMethods(year int) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly payment methods", zap.Int("year", year))

	records, err := s.transactionRepository.GetMonthlyPaymentMethods(year)
	if err != nil {
		s.logger.Error("Failed to fetch monthly payment methods", zap.Error(err), zap.Int("year", year))

		return nil, transaction_errors.ErrFailedFindMonthlyPaymentMethods
	}

	responses := s.mapping.ToTransactionMonthlyMethodResponses(records)

	s.logger.Debug("Successfully fetched monthly payment methods", zap.Int("year", year))

	return responses, nil
}

func (s *transactionService) FindYearlyPaymentMethods(year int) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly payment methods", zap.Int("year", year))

	records, err := s.transactionRepository.GetYearlyPaymentMethods(year)
	if err != nil {
		s.logger.Error("Failed to fetch yearly payment methods", zap.Error(err), zap.Int("year", year))

		return nil, transaction_errors.ErrFailedFindYearlyPaymentMethods
	}

	responses := s.mapping.ToTransactionYearlyMethodResponses(records)

	s.logger.Debug("Successfully fetched yearly payment methods", zap.Int("year", year))

	return responses, nil
}

func (s *transactionService) FindMonthlyAmounts(year int) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly amounts", zap.Int("year", year))

	records, err := s.transactionRepository.GetMonthlyAmounts(year)
	if err != nil {
		s.logger.Error("Failed to fetch monthly amounts", zap.Error(err), zap.Int("year", year))

		return nil, transaction_errors.ErrFailedFindMonthlyAmounts
	}

	responses := s.mapping.ToTransactionMonthlyAmountResponses(records)

	s.logger.Debug("Successfully fetched monthly amounts", zap.Int("year", year))

	return responses, nil
}

func (s *transactionService) FindYearlyAmounts(year int) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly amounts", zap.Int("year", year))

	records, err := s.transactionRepository.GetYearlyAmounts(year)
	if err != nil {
		s.logger.Error("Failed to fetch yearly amounts", zap.Error(err), zap.Int("year", year))

		return nil, transaction_errors.ErrFailedFindYearlyAmounts
	}

	responses := s.mapping.ToTransactionYearlyAmountResponses(records)

	s.logger.Debug("Successfully fetched yearly amounts", zap.Int("year", year))

	return responses, nil
}

func (s *transactionService) FindMonthlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching monthly payment methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.transactionRepository.GetMonthlyPaymentMethodsByCardNumber(req)
	if err != nil {
		s.logger.Error("Failed to fetch monthly payment methods by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))

		return nil, transaction_errors.ErrFailedFindMonthlyPaymentMethodsByCard
	}

	responses := s.mapping.ToTransactionMonthlyMethodResponses(records)

	s.logger.Debug("Successfully fetched monthly payment methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionService) FindYearlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching yearly payment methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.transactionRepository.GetYearlyPaymentMethodsByCardNumber(req)

	if err != nil {
		s.logger.Error("Failed to fetch yearly payment methods by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))

		return nil, transaction_errors.ErrFailedFindYearlyPaymentMethodsByCard
	}

	responses := s.mapping.ToTransactionYearlyMethodResponses(records)

	s.logger.Debug("Successfully fetched yearly payment methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionService) FindMonthlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching monthly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.transactionRepository.GetMonthlyAmountsByCardNumber(req)
	if err != nil {
		s.logger.Error("Failed to fetch monthly amounts by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))

		return nil, transaction_errors.ErrFailedFindMonthlyAmountsByCard
	}

	responses := s.mapping.ToTransactionMonthlyAmountResponses(records)

	s.logger.Debug("Successfully fetched monthly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionService) FindYearlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching yearly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.transactionRepository.GetYearlyAmountsByCardNumber(req)
	if err != nil {
		s.logger.Error("Failed to fetch yearly amounts by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))

		return nil, transaction_errors.ErrFailedFindYearlyAmountsByCard
	}

	responses := s.mapping.ToTransactionYearlyAmountResponses(records)

	s.logger.Debug("Successfully fetched yearly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionService) FindByActive(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching active transaction",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to fetch active transaction",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, transaction_errors.ErrFailedFindByActiveTransactions
	}

	so := s.mapping.ToTransactionsResponseDeleteAt(transactions)

	s.logger.Debug("Successfully fetched active transaction",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionService) FindByTrashed(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching trashed transaction",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to fetch trashed transaction",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, transaction_errors.ErrFailedFindByTrashedTransactions
	}

	so := s.mapping.ToTransactionsResponseDeleteAt(transactions)

	s.logger.Debug("Successfully fetched trashed transaction",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionService) FindTransactionByMerchantId(merchant_id int) ([]*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting FindTransactionByMerchantId process",
		zap.Int("merchantID", merchant_id),
	)

	res, err := s.transactionRepository.FindTransactionByMerchantId(merchant_id)

	if err != nil {
		s.logger.Error("Failed to fetch transaction by merchant ID", zap.Error(err), zap.Int("merchant_id", merchant_id))
		return nil, transaction_errors.ErrFailedFindByMerchantID
	}

	so := s.mapping.ToTransactionsResponse(res)

	s.logger.Debug("Successfully fetched transaction by merchant ID", zap.Int("merchant_id", merchant_id))

	return so, nil
}

func (s *transactionService) Create(apiKey string, request *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting CreateTransaction process",
		zap.String("apiKey", apiKey),
		zap.Any("request", request),
	)

	merchant, err := s.merchantRepository.FindByApiKey(apiKey)
	if err != nil {
		s.logger.Error("failed to find merchant", zap.Error(err))
		return nil, merchant_errors.ErrFailedFindByApiKey
	}

	card, err := s.cardRepository.FindCardByCardNumber(request.CardNumber)
	if err != nil {
		s.logger.Error("failed to find card", zap.Error(err))
		return nil, card_errors.ErrFailedFindByCardNumber
	}

	saldo, err := s.saldoRepository.FindByCardNumber(card.CardNumber)
	if err != nil {
		s.logger.Error("failed to find saldo", zap.Error(err))
		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	if saldo.TotalBalance < request.Amount {
		s.logger.Error("insufficient balance", zap.Int("AvailableBalance", saldo.TotalBalance), zap.Int("TransactionAmount", request.Amount))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance",
			Code:    http.StatusBadRequest,
		}
	}

	saldo.TotalBalance -= request.Amount
	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		s.logger.Error("failed to update saldo", zap.Error(err))
		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	request.MerchantID = &merchant.ID

	transaction, err := s.transactionRepository.CreateTransaction(request)
	if err != nil {
		saldo.TotalBalance += request.Amount
		_, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
			CardNumber:   card.CardNumber,
			TotalBalance: saldo.TotalBalance,
		})
		if err != nil {
			s.logger.Error("failed to update saldo", zap.Error(err))
			return nil, saldo_errors.ErrFailedUpdateSaldo
		}

		if _, err := s.transactionRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: transaction.ID,
			Status:        "failed",
		}); err != nil {
			s.logger.Error("failed to update transaction status", zap.Error(err))
		}

		return nil, transaction_errors.ErrFailedCreateTransaction
	}

	if _, err := s.transactionRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
		TransactionID: transaction.ID,
		Status:        "success",
	}); err != nil {
		s.logger.Error("failed to update transaction status", zap.Error(err))
		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	merchantCard, err := s.cardRepository.FindCardByUserId(merchant.UserID)
	if err != nil {
		s.logger.Error("failed to find merchant card", zap.Error(err))
		return nil, card_errors.ErrCardNotFoundRes
	}

	merchantSaldo, err := s.saldoRepository.FindByCardNumber(merchantCard.CardNumber)
	if err != nil {
		s.logger.Error("failed to find merchant saldo", zap.Error(err))
		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	merchantSaldo.TotalBalance += request.Amount
	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   merchantCard.CardNumber,
		TotalBalance: merchantSaldo.TotalBalance,
	}); err != nil {
		s.logger.Error("failed to update merchant saldo", zap.Error(err))
		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	so := s.mapping.ToTransactionResponse(transaction)

	s.logger.Debug("CreateTransaction process completed",
		zap.String("apiKey", apiKey),
		zap.Int("transactionID", transaction.ID),
	)

	return so, nil
}

func (s *transactionService) Update(apiKey string, request *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting UpdateTransaction process",
		zap.String("apiKey", apiKey),
		zap.Int("transaction_id", *request.TransactionID),
	)

	transaction, err := s.transactionRepository.FindById(*request.TransactionID)
	if err != nil {
		s.logger.Error("failed to find transaction", zap.Error(err))

		if _, err := s.transactionRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			s.logger.Error("failed to update transaction status", zap.Error(err))
		}

		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	merchant, err := s.merchantRepository.FindByApiKey(apiKey)
	if err != nil || transaction.MerchantID != merchant.ID {
		s.logger.Error("unauthorized access to transaction", zap.Error(err))

		if _, err := s.transactionRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			s.logger.Error("failed to update transaction status", zap.Error(err))
		}

		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	card, err := s.cardRepository.FindCardByCardNumber(transaction.CardNumber)
	if err != nil {
		s.logger.Error("failed to find card", zap.Error(err))

		if _, err := s.transactionRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			s.logger.Error("failed to update transaction status", zap.Error(err))
		}

		return nil, card_errors.ErrCardNotFoundRes
	}

	saldo, err := s.saldoRepository.FindByCardNumber(card.CardNumber)
	if err != nil {
		s.logger.Error("failed to find saldo", zap.Error(err))

		if _, err := s.transactionRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			s.logger.Error("failed to update transaction status", zap.Error(err))
		}

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	saldo.TotalBalance += transaction.Amount
	s.logger.Debug("Restoring balance for old transaction amount", zap.Int("RestoredBalance", saldo.TotalBalance))

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		s.logger.Error("failed to restore balance", zap.Error(err))

		if _, err := s.transactionRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			s.logger.Error("failed to update transaction status", zap.Error(err))
		}

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	if saldo.TotalBalance < request.Amount {
		s.logger.Error("insufficient balance for updated amount", zap.Int("AvailableBalance", saldo.TotalBalance), zap.Int("UpdatedAmount", request.Amount))

		if _, err := s.transactionRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			s.logger.Error("failed to update transaction status", zap.Error(err))
		}

		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	saldo.TotalBalance -= request.Amount
	s.logger.Info("Updating balance for updated transaction amount")

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		s.logger.Error("failed to update balance", zap.Error(err))

		if _, err := s.transactionRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			s.logger.Error("failed to update transaction status", zap.Error(err))
		}

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	transaction.Amount = request.Amount
	transaction.PaymentMethod = request.PaymentMethod

	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, transaction.TransactionTime)
	if err != nil {
		s.logger.Error("Failed to parse transaction time", zap.Error(err), zap.String("transaction_time", transaction.TransactionTime))

		if _, err := s.transactionRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			s.logger.Error("failed to update transaction status", zap.Error(err))
		}

		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	res, err := s.transactionRepository.UpdateTransaction(&requests.UpdateTransactionRequest{
		TransactionID:   &transaction.ID,
		CardNumber:      transaction.CardNumber,
		Amount:          transaction.Amount,
		PaymentMethod:   transaction.PaymentMethod,
		MerchantID:      &transaction.MerchantID,
		TransactionTime: parsedTime,
	})
	if err != nil {
		s.logger.Error("failed to update transaction", zap.Error(err))

		if _, err := s.transactionRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			s.logger.Error("failed to update transaction status", zap.Error(err))
		}

		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	if _, err := s.transactionRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
		TransactionID: transaction.ID,
		Status:        "success",
	}); err != nil {
		s.logger.Error("failed to update transaction status", zap.Error(err))
		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	so := s.mapping.ToTransactionResponse(res)

	s.logger.Debug("UpdateTransaction process completed",
		zap.String("apiKey", apiKey),
		zap.Int("transaction_id", *request.TransactionID),
	)

	return so, nil
}

func (s *transactionService) TrashedTransaction(transaction_id int) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting TrashedTransaction process",
		zap.Int("transaction_id", transaction_id),
	)

	res, err := s.transactionRepository.TrashedTransaction(transaction_id)

	if err != nil {
		s.logger.Error("Failed to move transaction to trash",
			zap.Int("transaction_id", transaction_id),
			zap.Error(err))

		return nil, transaction_errors.ErrFailedTrashedTransaction
	}

	so := s.mapping.ToTransactionResponse(res)

	s.logger.Debug("Successfully trashed transaction", zap.Int("transaction_id", transaction_id))

	return so, nil
}

func (s *transactionService) RestoreTransaction(transaction_id int) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting RestoreTransaction process",
		zap.Int("transaction_id", transaction_id),
	)

	res, err := s.transactionRepository.RestoreTransaction(transaction_id)

	if err != nil {
		s.logger.Error("Failed to restore transaction from trash",
			zap.Int("transaction_id", transaction_id),
			zap.Error(err))

		return nil, transaction_errors.ErrFailedRestoreTransaction
	}

	so := s.mapping.ToTransactionResponse(res)

	s.logger.Debug("Successfully restored transaction", zap.Int("transaction_id", transaction_id))

	return so, nil
}

func (s *transactionService) DeleteTransactionPermanent(transaction_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Starting DeleteTransactionPermanent process",
		zap.Int("transaction_id", transaction_id),
	)

	_, err := s.transactionRepository.DeleteTransactionPermanent(transaction_id)

	if err != nil {
		s.logger.Error("Failed to permanently delete transaction",
			zap.Int("transaction_id", transaction_id),
			zap.Error(err))

		return false, transaction_errors.ErrFailedDeleteTransactionPermanent
	}

	s.logger.Debug("Successfully permanently deleted transaction", zap.Int("transaction_id", transaction_id))

	return true, nil
}

func (s *transactionService) RestoreAllTransaction() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all transactions")

	_, err := s.transactionRepository.RestoreAllTransaction()
	if err != nil {
		s.logger.Error("Failed to restore all transaction", zap.Error(err))
		return false, transaction_errors.ErrFailedRestoreAllTransactions
	}

	s.logger.Debug("Successfully restored all transactions")
	return true, nil
}

func (s *transactionService) DeleteAllTransactionPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all transactions")

	_, err := s.transactionRepository.DeleteAllTransactionPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all transaction", zap.Error(err))
		return false, transaction_errors.ErrFailedDeleteAllTransactionsPermanent
	}

	s.logger.Debug("Successfully deleted all transactions permanently")

	return true, nil
}
