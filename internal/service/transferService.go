package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/service"

	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"fmt"

	"go.uber.org/zap"
)

type transferService struct {
	userRepository     repository.UserRepository
	cardRepository     repository.CardRepository
	saldoRepository    repository.SaldoRepository
	transferRepository repository.TransferRepository
	logger             logger.LoggerInterface
	mapping            responseservice.TransferResponseMapper
}

func NewTransferService(
	userRepository repository.UserRepository,
	cardRepository repository.CardRepository,
	transferRepository repository.TransferRepository,
	saldoRepository repository.SaldoRepository, logger logger.LoggerInterface, mapping responseservice.TransferResponseMapper) *transferService {
	return &transferService{
		userRepository:     userRepository,
		cardRepository:     cardRepository,
		transferRepository: transferRepository,
		saldoRepository:    saldoRepository,
		logger:             logger,
		mapping:            mapping,
	}
}

func (s *transferService) FindAll(page int, pageSize int, search string) ([]*response.TransferResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching transfer",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transfers, totalRecords, err := s.transferRepository.FindAll(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch transfer",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transfers",
		}
	}

	so := s.mapping.ToTransfersResponse(transfers)

	s.logger.Debug("Successfully fetched transfer",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transferService) FindById(transferId int) (*response.TransferResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching transfer by ID", zap.Int("transfer_id", transferId))

	transfer, err := s.transferRepository.FindById(transferId)

	if err != nil {
		s.logger.Error("failed to find transfer by ID", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Transfer not found",
		}
	}

	so := s.mapping.ToTransferResponse(transfer)

	s.logger.Debug("Successfully fetched transfer", zap.Int("transfer_id", transferId))

	return so, nil
}

func (s *transferService) FindMonthTransferStatusSuccess(year int, month int) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly Transfer status success", zap.Int("year", year), zap.Int("month", month))

	records, err := s.transferRepository.GetMonthTransferStatusSuccess(year, month)
	if err != nil {
		s.logger.Error("failed to fetch monthly Transfer status success", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Transfer status success",
		}
	}

	s.logger.Debug("Successfully fetched monthly Transfer status success", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTransferResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *transferService) FindYearlyTransferStatusSuccess(year int) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly Transfer status success", zap.Int("year", year))

	records, err := s.transferRepository.GetYearlyTransferStatusSuccess(year)
	if err != nil {
		s.logger.Error("failed to fetch yearly Transfer status success", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Transfer status success",
		}
	}

	s.logger.Debug("Successfully fetched yearly Transfer status success", zap.Int("year", year))

	so := s.mapping.ToTransferResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *transferService) FindMonthTransferStatusFailed(year int, month int) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month))

	records, err := s.transferRepository.GetMonthTransferStatusFailed(year, month)
	if err != nil {
		s.logger.Error("failed to fetch monthly Transfer status Failed", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Transfer status Failed",
		}
	}

	s.logger.Debug("Failedfully fetched monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTransferResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *transferService) FindYearlyTransferStatusFailed(year int) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly Transfer status Failed", zap.Int("year", year))

	records, err := s.transferRepository.GetYearlyTransferStatusFailed(year)
	if err != nil {
		s.logger.Error("failed to fetch yearly Transfer status Failed", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Transfer status Failed",
		}
	}

	s.logger.Debug("Failedfully fetched yearly Transfer status Failed", zap.Int("year", year))

	so := s.mapping.ToTransferResponsesYearStatusFailed(records)

	return so, nil
}

func (s *transferService) FindMonthTransferStatusSuccessByCardNumber(card_number string, year int, month int) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly Transfer status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.transferRepository.GetMonthTransferStatusSuccessByCardNumber(card_number, year, month)
	if err != nil {
		s.logger.Error("failed to fetch monthly Transfer status success", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Transfer status success",
		}
	}

	s.logger.Debug("Successfully fetched monthly Transfer status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToTransferResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *transferService) FindYearlyTransferStatusSuccessByCardNumber(card_number string, year int) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly Transfer status success", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.transferRepository.GetYearlyTransferStatusSuccessByCardNumber(card_number, year)
	if err != nil {
		s.logger.Error("failed to fetch yearly Transfer status success", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Transfer status success",
		}
	}

	s.logger.Debug("Successfully fetched yearly Transfer status success", zap.Int("year", year), zap.String("card_number", card_number))

	so := s.mapping.ToTransferResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *transferService) FindMonthTransferStatusFailedByCardNumber(card_number string, year int, month int) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.transferRepository.GetMonthTransferStatusFailedByCardNumber(card_number, year, month)
	if err != nil {
		s.logger.Error("failed to fetch monthly Transfer status Failed", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Transfer status Failed",
		}
	}

	s.logger.Debug("Failedfully fetched monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToTransferResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *transferService) FindYearlyTransferStatusFailedByCardNumber(card_number string, year int) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly Transfer status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.transferRepository.GetYearlyTransferStatusFailedByCardNumber(card_number, year)
	if err != nil {
		s.logger.Error("failed to fetch yearly Transfer status Failed", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Transfer status Failed",
		}
	}

	s.logger.Debug("Failedfully fetched yearly Transfer status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	so := s.mapping.ToTransferResponsesYearStatusFailed(records)

	return so, nil
}

func (s *transferService) FindMonthlyTransferAmounts(year int) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly transfer amounts", zap.Int("year", year))

	amounts, err := s.transferRepository.GetMonthlyTransferAmounts(year)
	if err != nil {
		s.logger.Error("failed to find monthly transfer amounts", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly transfer amounts",
		}
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

	s.logger.Debug("Successfully fetched monthly transfer amounts", zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferService) FindYearlyTransferAmounts(year int) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly transfer amounts", zap.Int("year", year))

	amounts, err := s.transferRepository.GetYearlyTransferAmounts(year)
	if err != nil {
		s.logger.Error("failed to find yearly transfer amounts", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly transfer amounts",
		}
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	s.logger.Debug("Successfully fetched yearly transfer amounts", zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferService) FindMonthlyTransferAmountsBySenderCardNumber(cardNumber string, year int) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	amounts, err := s.transferRepository.GetMonthlyTransferAmountsBySenderCardNumber(cardNumber, year)
	if err != nil {
		s.logger.Error("failed to find monthly transfer amounts by sender card number", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly transfer amounts by sender card number",
		}
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

	s.logger.Debug("Successfully fetched monthly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferService) FindMonthlyTransferAmountsByReceiverCardNumber(cardNumber string, year int) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	amounts, err := s.transferRepository.GetMonthlyTransferAmountsByReceiverCardNumber(cardNumber, year)
	if err != nil {
		s.logger.Error("failed to find monthly transfer amounts by receiver card number", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly transfer amounts by receiver card number",
		}
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

	s.logger.Debug("Successfully fetched monthly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferService) FindYearlyTransferAmountsBySenderCardNumber(cardNumber string, year int) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	amounts, err := s.transferRepository.GetYearlyTransferAmountsBySenderCardNumber(cardNumber, year)
	if err != nil {
		s.logger.Error("failed to find yearly transfer amounts by sender card number", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly transfer amounts by sender card number",
		}
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	s.logger.Debug("Successfully fetched yearly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferService) FindYearlyTransferAmountsByReceiverCardNumber(cardNumber string, year int) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	amounts, err := s.transferRepository.GetYearlyTransferAmountsByReceiverCardNumber(cardNumber, year)
	if err != nil {
		s.logger.Error("failed to find yearly transfer amounts by receiver card number", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly transfer amounts by receiver card number",
		}
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	s.logger.Debug("Successfully fetched yearly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferService) FindByActive(page int, pageSize int, search string) ([]*response.TransferResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching active transfer",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transfers, totalRecords, err := s.transferRepository.FindByActive(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch active transfer",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "No active transaction records found",
		}
	}

	so := s.mapping.ToTransfersResponseDeleteAt(transfers)

	s.logger.Debug("Successfully fetched active transfer",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transferService) FindByTrashed(page int, pageSize int, search string) ([]*response.TransferResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching trashed transfer",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transfers, totalRecords, err := s.transferRepository.FindByTrashed(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch trashed transfer",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "No trashed transaction records found",
		}
	}

	so := s.mapping.ToTransfersResponseDeleteAt(transfers)

	s.logger.Debug("Successfully fetched trashed transfer",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transferService) FindTransferByTransferFrom(transfer_from string) ([]*response.TransferResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting fetch transfer by transfer_from",
		zap.String("transfer_from", transfer_from),
	)

	res, err := s.transferRepository.FindTransferByTransferFrom(transfer_from)

	if err != nil {
		s.logger.Error("Failed to fetch transfers by transfer_from", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transfers by transfer_from",
		}
	}

	so := s.mapping.ToTransfersResponse(res)

	s.logger.Debug("Successfully fetched transfer record by transfer_from",
		zap.String("transfer_from", transfer_from),
	)

	return so, nil
}

func (s *transferService) FindTransferByTransferTo(transfer_to string) ([]*response.TransferResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting fetch transfer by transfer_to",
		zap.String("transfer_to", transfer_to),
	)

	res, err := s.transferRepository.FindTransferByTransferTo(transfer_to)

	if err != nil {
		s.logger.Error("Failed to fetch transfers by transfer_to", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transfers by transfer_to",
		}
	}

	so := s.mapping.ToTransfersResponse(res)

	s.logger.Debug("Successfully fetched transfer record by transfer_to",
		zap.String("transfer_to", transfer_to),
	)

	return so, nil
}

func (s *transferService) CreateTransaction(request *requests.CreateTransferRequest) (*response.TransferResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting create transaction process",
		zap.Any("request", request),
	)

	if s.logger == nil {
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Logger not initialized",
		}
	}

	if s.mapping == nil {
		s.logger.Error("Mapping service not initialized")
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Mapping service not initialized",
		}
	}

	_, err := s.cardRepository.FindCardByCardNumber(request.TransferFrom)
	if err != nil {
		s.logger.Error("failed to find sender card by Number", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Sender card not found",
		}
	}

	_, err = s.cardRepository.FindCardByCardNumber(request.TransferTo)
	if err != nil {
		s.logger.Error("failed to find receiver card by number", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Receiver card not found",
		}
	}

	senderSaldo, err := s.saldoRepository.FindByCardNumber(request.TransferFrom)
	if err != nil {
		s.logger.Error("failed to find sender saldo by card number", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find sender saldo",
		}
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(request.TransferTo)
	if err != nil {
		s.logger.Error("failed to find receiver saldo by card number", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find receiver saldo",
		}
	}

	if senderSaldo.TotalBalance < request.TransferAmount {
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance for sender",
		}
	}

	senderSaldo.TotalBalance -= request.TransferAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   senderSaldo.CardNumber,
		TotalBalance: senderSaldo.TotalBalance,
	})
	if err != nil {
		s.logger.Error("failed to update sender saldo", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update sender saldo",
		}
	}

	receiverSaldo.TotalBalance += request.TransferAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   receiverSaldo.CardNumber,
		TotalBalance: receiverSaldo.TotalBalance,
	})
	if err != nil {
		s.logger.Error("failed to update receiver saldo", zap.Error(err))

		senderSaldo.TotalBalance += request.TransferAmount
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
			CardNumber:   senderSaldo.CardNumber,
			TotalBalance: senderSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			s.logger.Error("failed to rollback sender saldo", zap.Error(rollbackErr))
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update receiver saldo",
		}
	}

	transfer, err := s.transferRepository.CreateTransfer(request)
	if err != nil {
		s.logger.Error("failed to create transfer", zap.Error(err))

		senderSaldo.TotalBalance += request.TransferAmount
		receiverSaldo.TotalBalance -= request.TransferAmount

		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
			CardNumber:   senderSaldo.CardNumber,
			TotalBalance: senderSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			s.logger.Error("failed to rollback sender saldo", zap.Error(rollbackErr))
		}

		_, rollbackErr = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
			CardNumber:   receiverSaldo.CardNumber,
			TotalBalance: receiverSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			s.logger.Error("failed to rollback receiver saldo", zap.Error(rollbackErr))
		}

		if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: transfer.ID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("failed to update transfer status", zap.Error(err))
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create transfer",
		}
	}

	if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
		TransferID: transfer.ID,
		Status:     "success",
	}); err != nil {
		s.logger.Error("failed to update transfer status", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update transfer status",
		}
	}

	so := s.mapping.ToTransferResponse(transfer)

	s.logger.Debug("successfully create transaction",
		zap.Int("transfer_id", transfer.ID),
	)

	return so, nil
}

func (s *transferService) UpdateTransaction(request *requests.UpdateTransferRequest) (*response.TransferResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting update transaction process",
		zap.Int("transfer_id", request.TransferID),
	)

	transfer, err := s.transferRepository.FindById(request.TransferID)
	if err != nil {
		s.logger.Error("Failed to find transfer by ID", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Transfer with ID %d not found: %v", request.TransferID, err),
		}
	}

	amountDifference := request.TransferAmount - transfer.TransferAmount

	senderSaldo, err := s.saldoRepository.FindByCardNumber(transfer.TransferFrom)

	if err != nil {
		s.logger.Error("Failed to find sender's saldo by user ID", zap.Error(err))

		if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: request.TransferID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update transfer status", zap.Error(err))
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to find sender's saldo: %v", err),
		}
	}

	newSenderBalance := senderSaldo.TotalBalance - amountDifference
	if newSenderBalance < 0 {
		s.logger.Error("Insufficient balance for sender", zap.String("senderID", transfer.TransferFrom))

		if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: request.TransferID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update transfer status", zap.Error(err))
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance for sender",
		}
	}

	senderSaldo.TotalBalance = newSenderBalance
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   senderSaldo.CardNumber,
		TotalBalance: senderSaldo.TotalBalance,
	})
	if err != nil {
		s.logger.Error("Failed to update sender's saldo", zap.Error(err))

		if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: request.TransferID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update transfer status", zap.Error(err))
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to update sender's saldo: %v", err),
		}
	}

	// Cari saldo penerima
	receiverSaldo, err := s.saldoRepository.FindByCardNumber(transfer.TransferTo)
	if err != nil {
		s.logger.Error("Failed to find receiver's saldo by user ID", zap.Error(err))

		// Rollback saldo pengirim
		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(rollbackSenderBalance)
		if rollbackErr != nil {
			s.logger.Error("Failed to rollback sender's saldo after receiver lookup failure", zap.Error(rollbackErr))
		}

		if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: request.TransferID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update transfer status", zap.Error(err))
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to find receiver's saldo: %v", err),
		}
	}

	newReceiverBalance := receiverSaldo.TotalBalance + amountDifference
	receiverSaldo.TotalBalance = newReceiverBalance
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   receiverSaldo.CardNumber,
		TotalBalance: receiverSaldo.TotalBalance,
	})
	if err != nil {
		s.logger.Error("Failed to update receiver's saldo", zap.Error(err))

		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		rollbackReceiverBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferTo,
			TotalBalance: receiverSaldo.TotalBalance - amountDifference,
		}

		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackSenderBalance); err != nil {
			s.logger.Error("Failed to rollback sender's saldo after receiver update failure", zap.Error(err))
		}
		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackReceiverBalance); err != nil {
			s.logger.Error("Failed to rollback receiver's saldo after sender update failure", zap.Error(err))
		}

		if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: request.TransferID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update transfer status", zap.Error(err))
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update receiver's saldo, rollback attempted but may be incomplete",
		}
	}

	// Update transfer di database
	updatedTransfer, err := s.transferRepository.UpdateTransfer(request)
	if err != nil {
		s.logger.Error("Failed to update transfer", zap.Error(err))

		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		rollbackReceiverBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferTo,
			TotalBalance: receiverSaldo.TotalBalance - amountDifference,
		}

		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackSenderBalance); err != nil {
			s.logger.Error("Failed to rollback sender's saldo after receiver update failure", zap.Error(err))
		}
		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackReceiverBalance); err != nil {
			s.logger.Error("Failed to rollback receiver's saldo after sender update failure", zap.Error(err))
		}

		if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: request.TransferID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update transfer status", zap.Error(err))
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to update transfer: %v", err),
		}
	}

	if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
		TransferID: request.TransferID,
		Status:     "success",
	}); err != nil {
		s.logger.Error("Failed to update transfer status", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update transfer status to success",
		}
	}

	so := s.mapping.ToTransferResponse(updatedTransfer)

	s.logger.Debug("successfully update transaction",
		zap.Int("transfer_id", request.TransferID),
	)

	return so, nil
}

func (s *transferService) TrashedTransfer(transfer_id int) (*response.TransferResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting trashed transfer process",
		zap.Int("transfer_id", transfer_id),
	)

	res, err := s.transferRepository.TrashedTransfer(transfer_id)

	if err != nil {
		s.logger.Error("Failed to trash transfer", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to trash transfer",
		}
	}

	so := s.mapping.ToTransferResponse(res)

	s.logger.Debug("successfully trashed transfer",
		zap.Int("transfer_id", transfer_id),
	)

	return so, nil
}

func (s *transferService) RestoreTransfer(transfer_id int) (*response.TransferResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting restore transfer process",
		zap.Int("transfer_id", transfer_id),
	)

	res, err := s.transferRepository.RestoreTransfer(transfer_id)

	if err != nil {
		s.logger.Error("Failed to restore transfer", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore transfer",
		}
	}

	so := s.mapping.ToTransferResponse(res)

	s.logger.Debug("successfully restore transfer",
		zap.Int("transfer_id", transfer_id),
	)

	return so, nil
}

func (s *transferService) DeleteTransferPermanent(transfer_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Starting delete transfer permanent process",
		zap.Int("transfer_id", transfer_id),
	)

	_, err := s.transferRepository.DeleteTransferPermanent(transfer_id)

	if err != nil {
		s.logger.Error("Failed to permanently delete transfer", zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete transfer",
		}
	}

	s.logger.Debug("successfully delete permanent transfer",
		zap.Int("transfer_id", transfer_id),
	)

	return true, nil
}

func (s *transferService) RestoreAllTransfer() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all transfers")

	_, err := s.transferRepository.RestoreAllTransfer()

	if err != nil {
		s.logger.Error("Failed to restore all transfers", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all transfers: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully restored all transfers")

	return true, nil
}

func (s *transferService) DeleteAllTransferPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all transfers")

	_, err := s.transferRepository.DeleteAllTransferPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all transfers", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all transfers: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully deleted all transfers permanently")
	return true, nil
}
