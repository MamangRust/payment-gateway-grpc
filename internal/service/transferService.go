package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/service"
	"net/http"

	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/pkg/errors/card_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/transfer_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"

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

func (s *transferService) FindAll(req *requests.FindAllTranfers) ([]*response.TransferResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

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

	transfers, totalRecords, err := s.transferRepository.FindAll(req)

	if err != nil {
		s.logger.Error("Failed to fetch transfer",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, transfer_errors.ErrFailedFindAllTransfers
	}

	so := s.mapping.ToTransfersResponse(transfers)

	s.logger.Debug("Successfully fetched transfer",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transferService) FindById(transferId int) (*response.TransferResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching transfer by ID", zap.Int("transfer_id", transferId))

	transfer, err := s.transferRepository.FindById(transferId)

	if err != nil {
		s.logger.Error("failed to find transfer by ID", zap.Error(err))

		return nil, transfer_errors.ErrTransferNotFound
	}

	so := s.mapping.ToTransferResponse(transfer)

	s.logger.Debug("Successfully fetched transfer", zap.Int("transfer_id", transferId))

	return so, nil
}

func (s *transferService) FindMonthTransferStatusSuccess(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	s.logger.Debug("Fetching monthly Transfer status success", zap.Int("year", year), zap.Int("month", month))

	records, err := s.transferRepository.GetMonthTransferStatusSuccess(req)

	if err != nil {
		s.logger.Error("failed to fetch monthly Transfer status success", zap.Error(err))

		return nil, transfer_errors.ErrFailedFindMonthTransferStatusSuccess
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

		return nil, transfer_errors.ErrFailedFindYearTransferStatusSuccess
	}

	s.logger.Debug("Successfully fetched yearly Transfer status success", zap.Int("year", year))

	so := s.mapping.ToTransferResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *transferService) FindMonthTransferStatusFailed(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	s.logger.Debug("Fetching monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month))

	records, err := s.transferRepository.GetMonthTransferStatusFailed(req)

	if err != nil {
		s.logger.Error("failed to fetch monthly Transfer status Failed", zap.Error(err))

		return nil, transfer_errors.ErrFailedFindMonthTransferStatusFailed
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

		return nil, transfer_errors.ErrFailedFindYearTransferStatusFailed
	}

	s.logger.Debug("Failedfully fetched yearly Transfer status Failed", zap.Int("year", year))

	so := s.mapping.ToTransferResponsesYearStatusFailed(records)

	return so, nil
}

func (s *transferService) FindMonthTransferStatusSuccessByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	s.logger.Debug("Fetching monthly Transfer status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.transferRepository.GetMonthTransferStatusSuccessByCardNumber(req)
	if err != nil {
		s.logger.Error("failed to fetch monthly Transfer status success", zap.Error(err))

		return nil, transfer_errors.ErrFailedFindMonthTransferStatusSuccess
	}

	s.logger.Debug("Successfully fetched monthly Transfer status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToTransferResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *transferService) FindYearlyTransferStatusSuccessByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching yearly Transfer status success", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.transferRepository.GetYearlyTransferStatusSuccessByCardNumber(req)
	if err != nil {
		s.logger.Error("failed to fetch yearly Transfer status success", zap.Error(err))

		return nil, transfer_errors.ErrFailedFindYearTransferStatusSuccessByCard
	}

	s.logger.Debug("Successfully fetched yearly Transfer status success", zap.Int("year", year), zap.String("card_number", card_number))

	so := s.mapping.ToTransferResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *transferService) FindMonthTransferStatusFailedByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	s.logger.Debug("Fetching monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.transferRepository.GetMonthTransferStatusFailedByCardNumber(req)

	if err != nil {
		s.logger.Error("failed to fetch monthly Transfer status Failed", zap.Error(err))

		return nil, transfer_errors.ErrFailedFindMonthTransferStatusFailed
	}

	s.logger.Debug("Failedfully fetched monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToTransferResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *transferService) FindYearlyTransferStatusFailedByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching yearly Transfer status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.transferRepository.GetYearlyTransferStatusFailedByCardNumber(req)
	if err != nil {
		s.logger.Error("failed to fetch yearly Transfer status Failed", zap.Error(err))

		return nil, transfer_errors.ErrFailedFindYearTransferStatusFailedByCard
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
		return nil, transfer_errors.ErrFailedFindMonthlyTransferAmounts
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
		return nil, transfer_errors.ErrFailedFindYearlyTransferAmounts
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	s.logger.Debug("Successfully fetched yearly transfer amounts", zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferService) FindMonthlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching monthly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	amounts, err := s.transferRepository.GetMonthlyTransferAmountsBySenderCardNumber(req)
	if err != nil {
		s.logger.Error("failed to find monthly transfer amounts by sender card number", zap.Error(err))
		return nil, transfer_errors.ErrFailedFindMonthlyTransferAmountsBySenderCard
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

	s.logger.Debug("Successfully fetched monthly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferService) FindMonthlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching monthly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	amounts, err := s.transferRepository.GetMonthlyTransferAmountsByReceiverCardNumber(req)
	if err != nil {
		s.logger.Error("failed to find monthly transfer amounts by receiver card number", zap.Error(err))
		return nil, transfer_errors.ErrFailedFindMonthlyTransferAmountsByReceiverCard
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

	s.logger.Debug("Successfully fetched monthly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferService) FindYearlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching yearly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	amounts, err := s.transferRepository.GetYearlyTransferAmountsBySenderCardNumber(req)
	if err != nil {
		s.logger.Error("failed to find yearly transfer amounts by sender card number", zap.Error(err))
		return nil, transfer_errors.ErrFailedFindYearlyTransferAmountsBySenderCard
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	s.logger.Debug("Successfully fetched yearly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferService) FindYearlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching yearly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	amounts, err := s.transferRepository.GetYearlyTransferAmountsByReceiverCardNumber(req)
	if err != nil {
		s.logger.Error("failed to find yearly transfer amounts by receiver card number", zap.Error(err))
		return nil, transfer_errors.ErrFailedFindYearlyTransferAmountsByReceiverCard
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	s.logger.Debug("Successfully fetched yearly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferService) FindByActive(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

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

	transfers, totalRecords, err := s.transferRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to fetch active transfer",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, transfer_errors.ErrFailedFindActiveTransfers
	}

	so := s.mapping.ToTransfersResponseDeleteAt(transfers)

	s.logger.Debug("Successfully fetched active transfer",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transferService) FindByTrashed(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

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

	transfers, totalRecords, err := s.transferRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to fetch trashed transfer",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, transfer_errors.ErrFailedFindTrashedTransfers
	}

	so := s.mapping.ToTransfersResponseDeleteAt(transfers)

	s.logger.Debug("Successfully fetched trashed transfer",
		zap.Int("totalRecords", *totalRecords),
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

		return nil, transfer_errors.ErrTransferNotFound
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
		return nil, transfer_errors.ErrTransferNotFound
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

	_, err := s.cardRepository.FindCardByCardNumber(request.TransferFrom)
	if err != nil {
		s.logger.Error("failed to find sender card by Number", zap.Error(err))
		return nil, card_errors.ErrCardNotFoundRes
	}

	_, err = s.cardRepository.FindCardByCardNumber(request.TransferTo)
	if err != nil {
		s.logger.Error("failed to find receiver card by number", zap.Error(err))
		return nil, card_errors.ErrCardNotFoundRes
	}

	senderSaldo, err := s.saldoRepository.FindByCardNumber(request.TransferFrom)
	if err != nil {
		s.logger.Error("failed to find sender saldo by card number", zap.Error(err))
		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(request.TransferTo)
	if err != nil {
		s.logger.Error("failed to find receiver saldo by card number", zap.Error(err))
		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	if senderSaldo.TotalBalance < request.TransferAmount {
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance for sender",
			Code:    http.StatusBadRequest,
		}
	}

	senderSaldo.TotalBalance -= request.TransferAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   senderSaldo.CardNumber,
		TotalBalance: senderSaldo.TotalBalance,
	})
	if err != nil {
		s.logger.Error("failed to update sender saldo", zap.Error(err))
		return nil, saldo_errors.ErrFailedUpdateSaldo
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

		return nil, saldo_errors.ErrFailedUpdateSaldo
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

		return nil, transfer_errors.ErrFailedCreateTransfer
	}

	if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
		TransferID: transfer.ID,
		Status:     "success",
	}); err != nil {
		s.logger.Error("failed to update transfer status", zap.Error(err))
		return nil, transfer_errors.ErrFailedUpdateTransfer
	}

	so := s.mapping.ToTransferResponse(transfer)

	s.logger.Debug("successfully create transaction",
		zap.Int("transfer_id", transfer.ID),
	)

	return so, nil
}

func (s *transferService) UpdateTransaction(request *requests.UpdateTransferRequest) (*response.TransferResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting update transaction process",
		zap.Int("transfer_id", *request.TransferID),
	)

	transfer, err := s.transferRepository.FindById(*request.TransferID)
	if err != nil {
		s.logger.Error("Failed to find transfer by ID", zap.Error(err))

		return nil, transfer_errors.ErrTransferNotFound
	}

	amountDifference := request.TransferAmount - transfer.TransferAmount

	senderSaldo, err := s.saldoRepository.FindByCardNumber(transfer.TransferFrom)
	if err != nil {
		s.logger.Error("Failed to find sender's saldo by user ID", zap.Error(err))

		if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update transfer status", zap.Error(err))
		}

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	newSenderBalance := senderSaldo.TotalBalance - amountDifference
	if newSenderBalance < 0 {
		s.logger.Error("Insufficient balance for sender", zap.String("senderID", transfer.TransferFrom))

		if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update transfer status", zap.Error(err))
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance for sender",
			Code:    http.StatusBadRequest,
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
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update transfer status", zap.Error(err))
		}

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(transfer.TransferTo)
	if err != nil {
		s.logger.Error("Failed to find receiver's saldo by user ID", zap.Error(err))

		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(rollbackSenderBalance)
		if rollbackErr != nil {
			s.logger.Error("Failed to rollback sender's saldo after receiver lookup failure", zap.Error(rollbackErr))
		}

		if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update transfer status", zap.Error(err))
		}

		return nil, saldo_errors.ErrFailedSaldoNotFound
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
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update transfer status", zap.Error(err))
		}

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

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
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update transfer status", zap.Error(err))
		}

		return nil, transfer_errors.ErrFailedUpdateTransfer
	}

	if _, err := s.transferRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
		TransferID: *request.TransferID,
		Status:     "success",
	}); err != nil {
		s.logger.Error("Failed to update transfer status", zap.Error(err))
		return nil, transfer_errors.ErrFailedUpdateTransfer
	}

	so := s.mapping.ToTransferResponse(updatedTransfer)

	s.logger.Debug("successfully update transaction",
		zap.Int("transfer_id", *request.TransferID),
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

		return nil, transfer_errors.ErrFailedTrashedTransfer
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

		return nil, transfer_errors.ErrFailedRestoreTransfer
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

		return false, transfer_errors.ErrFailedDeleteTransferPermanent
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
		return false, transfer_errors.ErrFailedRestoreAllTransfers
	}

	s.logger.Debug("Successfully restored all transfers")

	return true, nil
}

func (s *transferService) DeleteAllTransferPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all transfers")

	_, err := s.transferRepository.DeleteAllTransferPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all transfers", zap.Error(err))
		return false, transfer_errors.ErrFailedDeleteAllTransfersPermanent
	}

	s.logger.Debug("Successfully deleted all transfers permanently")
	return true, nil
}
