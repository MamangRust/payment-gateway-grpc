package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	responsemapper "MamangRust/paymentgatewaygrpc/internal/mapper/response"
	"MamangRust/paymentgatewaygrpc/internal/repository"
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
	mapping               responsemapper.TransactionResponseMapper
}

func NewTransactionService(
	merchantRepository repository.MerchantRepository,
	cardRepository repository.CardRepository,
	saldoRepository repository.SaldoRepository,
	transactionRepository repository.TransactionRepository,
	logger logger.LoggerInterface,
	mapping responsemapper.TransactionResponseMapper,
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

func (s *transactionService) FindAll(page int, pageSize int, search string) ([]*response.TransactionResponse, int, *response.ErrorResponse) {
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

	transactions, totalRecords, err := s.transactionRepository.FindAllTransactions(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch transaction",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transactions",
		}
	}

	so := s.mapping.ToTransactionsResponse(transactions)

	s.logger.Debug("Successfully fetched transaction",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionService) FindById(transactionID int) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching transaction by ID", zap.Int("transaction_id", transactionID))

	transaction, err := s.transactionRepository.FindById(transactionID)

	if err != nil {
		s.logger.Error("failed to find transaction", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Transaction not found",
		}
	}

	so := s.mapping.ToTransactionResponse(transaction)

	s.logger.Debug("Successfully fetched transaction", zap.Int("transaction_id", transactionID))

	return so, nil
}

func (s *transactionService) FindByActive(page int, pageSize int, search string) ([]*response.TransactionResponseDeleteAt, int, *response.ErrorResponse) {
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

	transactions, totalRecords, err := s.transactionRepository.FindByActive(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch active transaction",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "No active transaction records found",
		}
	}

	so := s.mapping.ToTransactionsResponseDeleteAt(transactions)

	s.logger.Debug("Successfully fetched active transaction",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionService) FindByTrashed(page int, pageSize int, search string) ([]*response.TransactionResponseDeleteAt, int, *response.ErrorResponse) {
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

	transactions, totalRecords, err := s.transactionRepository.FindByTrashed(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch trashed transaction",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "No trashed transaction records found",
		}
	}

	so := s.mapping.ToTransactionsResponseDeleteAt(transactions)

	s.logger.Debug("Successfully fetched trashed transaction",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionService) FindByCardNumber(card_number string) ([]*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching card record by transaction", zap.String("card_number", card_number))

	res, err := s.transactionRepository.FindByCardNumber(card_number)

	if err != nil {
		s.logger.Error("Failed to fetch transactions by card number", zap.Error(err), zap.String("card_number", card_number))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "No transactions found for the given card number",
		}
	}

	so := s.mapping.ToTransactionsResponse(res)

	s.logger.Debug("Successfully fetched transactions by card number", zap.String("card_number", card_number), zap.Int("record_count", len(res)))

	return so, nil
}

func (s *transactionService) FindTransactionByMerchantId(merchant_id int) ([]*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting FindTransactionByMerchantId process",
		zap.Int("merchantID", merchant_id),
	)

	res, err := s.transactionRepository.FindTransactionByMerchantId(merchant_id)

	if err != nil {
		s.logger.Error("Failed to fetch transaction by merchant ID", zap.Error(err), zap.Int("merchant_id", merchant_id))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "No transaction found for the given merchant ID",
		}
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
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Merchant not found",
		}
	}

	card, err := s.cardRepository.FindCardByCardNumber(request.CardNumber)

	if err != nil {
		s.logger.Error("failed to find card", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Card not found",
		}
	}

	saldo, err := s.saldoRepository.FindByCardNumber(card.CardNumber)

	if err != nil {
		s.logger.Error("failed to find saldo", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Saldo not found",
		}
	}

	if saldo.TotalBalance < request.Amount {
		s.logger.Error("insufficient balance", zap.Int("AvailableBalance", saldo.TotalBalance), zap.Int("TransactionAmount", request.Amount))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance",
		}
	}

	saldo.TotalBalance -= request.Amount

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		s.logger.Error("failed to update saldo", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update saldo",
		}
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
			return nil, &response.ErrorResponse{
				Status:  "error",
				Message: "Failed to update saldo",
			}
		}
	}

	merchantCard, err := s.cardRepository.FindCardByUserId(merchant.UserID)
	if err != nil {
		s.logger.Error("failed to find merchant card", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Merchant card not found",
		}
	}

	merchantSaldo, err := s.saldoRepository.FindByCardNumber(merchantCard.CardNumber)
	if err != nil {
		s.logger.Error("failed to find merchant saldo", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Merchant saldo not found",
		}
	}

	merchantSaldo.TotalBalance += request.Amount

	s.logger.Debug("Updating merchant saldo", zap.Int("NewMerchantBalance",
		merchantSaldo.TotalBalance))

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   merchantCard.CardNumber,
		TotalBalance: merchantSaldo.TotalBalance,
	}); err != nil {
		s.logger.Error("failed to update merchant saldo", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update merchant saldo",
		}
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
		zap.Int("transaction_id", request.TransactionID),
	)

	transaction, err := s.transactionRepository.FindById(request.TransactionID)

	if err != nil {
		s.logger.Error("failed to find transaction", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Transaction not found",
		}
	}

	merchant, err := s.merchantRepository.FindByApiKey(apiKey)

	if err != nil || transaction.MerchantID != merchant.ID {
		s.logger.Error("unauthorized access to transaction", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Unauthorized access to transaction",
		}
	}

	card, err := s.cardRepository.FindCardByCardNumber(transaction.CardNumber)

	if err != nil {
		s.logger.Error("failed to find card", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Card not found",
		}
	}

	saldo, err := s.saldoRepository.FindByCardNumber(card.CardNumber)

	if err != nil {
		s.logger.Error("failed to find saldo", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Saldo not found",
		}
	}

	saldo.TotalBalance += transaction.Amount

	s.logger.Debug("Restoring balance for old transaction amount", zap.Int("RestoredBalance", saldo.TotalBalance))

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		s.logger.Error("failed to restore balance", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore balance",
		}
	}

	if saldo.TotalBalance < request.Amount {
		s.logger.Error("insufficient balance for updated amount", zap.Int("AvailableBalance", saldo.TotalBalance), zap.Int("UpdatedAmount", request.Amount))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance for updated transaction",
		}
	}

	saldo.TotalBalance -= request.Amount

	s.logger.Info("Updating balance for updated transaction amount")

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		s.logger.Error("failed to update balance", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update balance",
		}
	}

	transaction.Amount = request.Amount
	transaction.PaymentMethod = request.PaymentMethod

	layout := "2006-01-02 15:04:05"

	parsedTime, err := time.Parse(layout, transaction.TransactionTime)

	if err != nil {
		s.logger.Error("Failed to parse transaction time", zap.Error(err), zap.String("transaction_time", transaction.TransactionTime))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Invalid transaction time format",
		}
	}

	res, err := s.transactionRepository.UpdateTransaction(&requests.UpdateTransactionRequest{
		TransactionID:   transaction.ID,
		CardNumber:      transaction.CardNumber,
		Amount:          transaction.Amount,
		PaymentMethod:   transaction.PaymentMethod,
		MerchantID:      &transaction.MerchantID,
		TransactionTime: parsedTime,
	})

	if err != nil {
		s.logger.Error("failed to update transaction", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update transaction",
		}
	}

	so := s.mapping.ToTransactionResponse(res)

	s.logger.Debug("UpdateTransaction process completed",
		zap.String("apiKey", apiKey),
		zap.Int("transaction_id", request.TransactionID),
	)

	return so, nil
}

func (s *transactionService) TrashedTransaction(transaction_id int) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting TrashedTransaction process",
		zap.Int("transaction_id", transaction_id),
	)

	res, err := s.transactionRepository.TrashedTransaction(transaction_id)

	if err != nil {
		s.logger.Error("Failed to trash transaction", zap.Error(err), zap.Int("transaction_id", transaction_id))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to trash transaction",
		}
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
		s.logger.Error("Failed to restore transaction", zap.Error(err), zap.Int("transaction_id", transaction_id))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore transaction",
		}
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
		s.logger.Error("Failed to permanently delete transaction", zap.Error(err), zap.Int("transaction_id", transaction_id))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete transaction",
		}
	}

	s.logger.Debug("Successfully permanently deleted transaction", zap.Int("transaction_id", transaction_id))

	return true, nil
}

func (s *transactionService) RestoreAllTransaction() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all transactions")

	_, err := s.transactionRepository.RestoreAllTransaction()
	if err != nil {
		s.logger.Error("Failed to restore all transactions", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all transactions: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully restored all transactions")
	return true, nil
}

func (s *transactionService) DeleteAllTransactionPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all transactions")

	_, err := s.transactionRepository.DeleteAllTransactionPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all transactions", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all transactions: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully deleted all transactions permanently")

	return true, nil
}
