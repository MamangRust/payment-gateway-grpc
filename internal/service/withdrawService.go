package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/service"

	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/pkg/logger"

	"go.uber.org/zap"
)

type withdrawService struct {
	userRepository     repository.UserRepository
	saldoRepository    repository.SaldoRepository
	withdrawRepository repository.WithdrawRepository
	logger             logger.LoggerInterface
	mapping            responseservice.WithdrawResponseMapper
}

func NewWithdrawService(
	userRepository repository.UserRepository,
	withdrawRepository repository.WithdrawRepository, saldoRepository repository.SaldoRepository, logger logger.LoggerInterface, mapping responseservice.WithdrawResponseMapper) *withdrawService {
	return &withdrawService{
		userRepository:     userRepository,
		saldoRepository:    saldoRepository,
		withdrawRepository: withdrawRepository,
		logger:             logger,
		mapping:            mapping,
	}
}

func (s *withdrawService) FindAll(page int, pageSize int, search string) ([]*response.WithdrawResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching withdraw",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	withdraws, totalRecords, err := s.withdrawRepository.FindAll(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch withdraw",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch withdraws",
		}
	}

	withdrawResponse := s.mapping.ToWithdrawsResponse(withdraws)

	s.logger.Debug("Successfully fetched withdraw",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdrawResponse, totalRecords, nil
}

func (s *withdrawService) FindAllByCardNumber(card_number string, page int, pageSize int, search string) ([]*response.WithdrawResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching withdraw",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	withdraws, totalRecords, err := s.withdrawRepository.FindAllByCardNumber(card_number, search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch withdraw",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch withdraws",
		}
	}

	withdrawResponse := s.mapping.ToWithdrawsResponse(withdraws)

	s.logger.Debug("Successfully fetched withdraw",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdrawResponse, totalRecords, nil
}

func (s *withdrawService) FindById(withdrawID int) (*response.WithdrawResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching withdraw by ID", zap.Int("withdraw_id", withdrawID))

	withdraw, err := s.withdrawRepository.FindById(withdrawID)

	if err != nil {
		s.logger.Error("failed to find withdraw by id", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch withdraw record by ID.",
		}
	}
	so := s.mapping.ToWithdrawResponse(withdraw)

	s.logger.Debug("Successfully fetched withdraw", zap.Int("withdraw_id", withdrawID))

	return so, nil
}

func (s *withdrawService) FindMonthWithdrawStatusSuccess(year int, month int) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month))

	records, err := s.withdrawRepository.GetMonthWithdrawStatusSuccess(year, month)
	if err != nil {
		s.logger.Error("failed to fetch monthly Withdraw status success", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Withdraw status success",
		}
	}

	s.logger.Debug("Successfully fetched monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToWithdrawResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *withdrawService) FindYearlyWithdrawStatusSuccess(year int) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly Withdraw status success", zap.Int("year", year))

	records, err := s.withdrawRepository.GetYearlyWithdrawStatusSuccess(year)
	if err != nil {
		s.logger.Error("failed to fetch yearly Withdraw status success", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Withdraw status success",
		}
	}

	s.logger.Debug("Successfully fetched yearly Withdraw status success", zap.Int("year", year))

	so := s.mapping.ToWithdrawResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *withdrawService) FindMonthWithdrawStatusFailed(year int, month int) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month))

	records, err := s.withdrawRepository.GetMonthWithdrawStatusFailed(year, month)
	if err != nil {
		s.logger.Error("failed to fetch monthly Withdraw status Failed", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Withdraw status Failed",
		}
	}

	s.logger.Debug("Failedfully fetched monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToWithdrawResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *withdrawService) FindYearlyWithdrawStatusFailed(year int) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly Withdraw status Failed", zap.Int("year", year))

	records, err := s.withdrawRepository.GetYearlyWithdrawStatusFailed(year)
	if err != nil {
		s.logger.Error("failed to fetch yearly Withdraw status Failed", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Withdraw status Failed",
		}
	}

	s.logger.Debug("Failedfully fetched yearly Withdraw status Failed", zap.Int("year", year))

	so := s.mapping.ToWithdrawResponsesYearStatusFailed(records)

	return so, nil
}

func (s *withdrawService) FindMonthWithdrawStatusSuccessByCardNumber(card_number string, year int, month int) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.withdrawRepository.GetMonthWithdrawStatusSuccessByCardNumber(card_number, year, month)
	if err != nil {
		s.logger.Error("failed to fetch monthly Withdraw status success", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Withdraw status success",
		}
	}

	s.logger.Debug("Successfully fetched monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToWithdrawResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *withdrawService) FindYearlyWithdrawStatusSuccessByCardNumber(card_number string, year int) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly Withdraw status success", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.withdrawRepository.GetYearlyWithdrawStatusSuccessByCardNumber(card_number, year)
	if err != nil {
		s.logger.Error("failed to fetch yearly Withdraw status success", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Withdraw status success",
		}
	}

	s.logger.Debug("Successfully fetched yearly Withdraw status success", zap.Int("year", year), zap.String("card_number", card_number))

	so := s.mapping.ToWithdrawResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *withdrawService) FindMonthWithdrawStatusFailedByCardNumber(card_number string, year int, month int) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.withdrawRepository.GetMonthWithdrawStatusFailedByCardNumber(card_number, year, month)
	if err != nil {
		s.logger.Error("failed to fetch monthly Withdraw status Failed", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly Withdraw status Failed",
		}
	}

	s.logger.Debug("Failedfully fetched monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToWithdrawResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *withdrawService) FindYearlyWithdrawStatusFailedByCardNumber(card_number string, year int) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly Withdraw status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.withdrawRepository.GetYearlyWithdrawStatusFailedByCardNumber(card_number, year)
	if err != nil {
		s.logger.Error("failed to fetch yearly Withdraw status Failed", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly Withdraw status Failed",
		}
	}

	s.logger.Debug("Failedfully fetched yearly Withdraw status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	so := s.mapping.ToWithdrawResponsesYearStatusFailed(records)

	return so, nil
}

func (s *withdrawService) FindMonthlyWithdraws(year int) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly withdraws", zap.Int("year", year))

	withdraws, err := s.withdrawRepository.GetMonthlyWithdraws(year)
	if err != nil {
		s.logger.Error("failed to find monthly withdraws", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly withdraws",
		}
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountMonthlyResponses(withdraws)

	s.logger.Debug("Successfully fetched monthly withdraws", zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawService) FindYearlyWithdraws(year int) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly withdraws", zap.Int("year", year))

	withdraws, err := s.withdrawRepository.GetYearlyWithdraws(year)
	if err != nil {
		s.logger.Error("failed to find yearly withdraws", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly withdraws",
		}
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountYearlyResponses(withdraws)

	s.logger.Debug("Successfully fetched yearly withdraws", zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawService) FindMonthlyWithdrawsByCardNumber(cardNumber string, year int) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	withdraws, err := s.withdrawRepository.GetMonthlyWithdrawsByCardNumber(cardNumber, year)
	if err != nil {
		s.logger.Error("failed to find monthly withdraws by card number", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly withdraws by card number",
		}
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountMonthlyResponses(withdraws)

	s.logger.Debug("Successfully fetched monthly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawService) FindYearlyWithdrawsByCardNumber(cardNumber string, year int) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	withdraws, err := s.withdrawRepository.GetYearlyWithdrawsByCardNumber(cardNumber, year)
	if err != nil {
		s.logger.Error("failed to find yearly withdraws by card number", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly withdraws by card number",
		}
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountYearlyResponses(withdraws)

	s.logger.Debug("Successfully fetched yearly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawService) FindByActive(page int, pageSize int, search string) ([]*response.WithdrawResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching active withdraw",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	withdraws, totalRecords, err := s.withdrawRepository.FindByActive(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch active withdraw",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch active withdraw records",
		}
	}

	withdrawResponses := s.mapping.ToWithdrawsResponseDeleteAt(withdraws)

	s.logger.Debug("Successfully fetched active withdraw",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdrawResponses, totalRecords, nil
}

func (s *withdrawService) FindByTrashed(page int, pageSize int, search string) ([]*response.WithdrawResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching trashed withdraw",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	withdraws, totalRecords, err := s.withdrawRepository.FindByTrashed(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch trashed withdraw",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch trashed withdraw records",
		}
	}

	withdrawResponses := s.mapping.ToWithdrawsResponseDeleteAt(withdraws)

	s.logger.Debug("Successfully fetched trashed withdraw",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdrawResponses, totalRecords, nil
}

func (s *withdrawService) Create(request *requests.CreateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new withdraw", zap.Any("request", request))

	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)

	if err != nil {
		s.logger.Error("Failed to find saldo by user ID", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch saldo for the user.",
		}
	}

	if saldo == nil {
		s.logger.Error("Saldo not found for user", zap.String("cardNumber", request.CardNumber))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Saldo not found for the specified user ID.",
		}
	}

	if saldo.TotalBalance < request.WithdrawAmount {
		s.logger.Error("Insufficient balance for user", zap.String("cardNumber", request.CardNumber), zap.Int("requested", request.WithdrawAmount))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance for withdrawal.",
		}
	}

	newTotalBalance := saldo.TotalBalance - request.WithdrawAmount

	updateData := &requests.UpdateSaldoWithdraw{
		CardNumber:     request.CardNumber,
		TotalBalance:   newTotalBalance,
		WithdrawAmount: &request.WithdrawAmount,
		WithdrawTime:   &request.WithdrawTime,
	}

	_, err = s.saldoRepository.UpdateSaldoWithdraw(updateData)
	if err != nil {
		s.logger.Error("Failed to update saldo after withdrawal", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update saldo after withdrawal.",
		}
	}

	withdrawRecord, err := s.withdrawRepository.CreateWithdraw(request)
	if err != nil {
		s.logger.Error("Failed to create withdraw record", zap.Error(err))

		rollbackData := &requests.UpdateSaldoWithdraw{
			CardNumber:     request.CardNumber,
			TotalBalance:   saldo.TotalBalance,
			WithdrawAmount: &request.WithdrawAmount,
			WithdrawTime:   &request.WithdrawTime,
		}
		if _, rollbackErr := s.saldoRepository.UpdateSaldoWithdraw(rollbackData); rollbackErr != nil {
			s.logger.Error("Failed to rollback saldo after withdraw creation failure", zap.Error(rollbackErr))
		}

		if _, err := s.withdrawRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
			WithdrawID: withdrawRecord.ID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update withdraw status", zap.Error(err))
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create withdraw record.",
		}
	}

	if _, err := s.withdrawRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
		WithdrawID: withdrawRecord.ID,
		Status:     "success",
	}); err != nil {
		s.logger.Error("Failed to update withdraw status", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update withdraw status to success.",
		}
	}

	so := s.mapping.ToWithdrawResponse(withdrawRecord)

	s.logger.Debug("Successfully created withdraw", zap.Int("withdraw_id", withdrawRecord.ID))

	return so, nil
}

func (s *withdrawService) Update(request *requests.UpdateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating withdraw", zap.Int("withdraw_id", request.WithdrawID), zap.Any("request", request))

	_, err := s.withdrawRepository.FindById(request.WithdrawID)
	if err != nil {
		s.logger.Error("Failed to find withdraw record by ID", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Withdraw record not found.",
		}
	}

	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)
	if err != nil {
		s.logger.Error("Failed to fetch saldo by user ID", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch saldo for the user.",
		}
	}

	if saldo.TotalBalance < request.WithdrawAmount {
		s.logger.Error("Insufficient balance for user", zap.String("cardNumber", request.CardNumber))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance for withdrawal update.",
		}
	}

	newTotalBalance := saldo.TotalBalance - request.WithdrawAmount
	updateSaldoData := &requests.UpdateSaldoWithdraw{
		CardNumber:     saldo.CardNumber,
		TotalBalance:   newTotalBalance,
		WithdrawAmount: &request.WithdrawAmount,
		WithdrawTime:   &request.WithdrawTime,
	}

	_, err = s.saldoRepository.UpdateSaldoWithdraw(updateSaldoData)
	if err != nil {
		s.logger.Error("Failed to update saldo balance", zap.Error(err))

		if _, err := s.withdrawRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
			WithdrawID: request.WithdrawID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update withdraw status", zap.Error(err))
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update saldo balance.",
		}
	}

	updatedWithdraw, err := s.withdrawRepository.UpdateWithdraw(request)
	if err != nil {
		s.logger.Error("Failed to update withdraw record", zap.Error(err))

		rollbackData := &requests.UpdateSaldoBalance{
			CardNumber:   saldo.CardNumber,
			TotalBalance: saldo.TotalBalance,
		}
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(rollbackData)
		if rollbackErr != nil {
			s.logger.Error("Failed to rollback saldo after withdraw update failure", zap.Error(rollbackErr))
		}

		if _, err := s.withdrawRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
			WithdrawID: request.WithdrawID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update withdraw status", zap.Error(err))
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update withdraw record.",
		}
	}

	if _, err := s.withdrawRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
		WithdrawID: updatedWithdraw.ID,
		Status:     "success",
	}); err != nil {
		s.logger.Error("Failed to update withdraw status", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update withdraw status to success.",
		}
	}

	so := s.mapping.ToWithdrawResponse(updatedWithdraw)

	s.logger.Debug("Successfully updated withdraw", zap.Int("withdraw_id", so.ID))

	return so, nil
}

func (s *withdrawService) TrashedWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse) {
	s.logger.Debug("Trashing withdraw", zap.Int("withdraw_id", withdraw_id))

	res, err := s.withdrawRepository.TrashedWithdraw(withdraw_id)

	if err != nil {
		s.logger.Error("Failed to trash withdraw", zap.Error(err), zap.Int("withdraw_id", withdraw_id))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to trash withdraw",
		}
	}

	withdrawResponse := s.mapping.ToWithdrawResponse(res)

	s.logger.Debug("Successfully trashed withdraw", zap.Int("withdraw_id", withdraw_id))

	return withdrawResponse, nil
}

func (s *withdrawService) RestoreWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse) {
	s.logger.Debug("Restoring withdraw", zap.Int("withdraw_id", withdraw_id))

	res, err := s.withdrawRepository.RestoreWithdraw(withdraw_id)

	if err != nil {
		s.logger.Error("Failed to restore withdraw", zap.Error(err), zap.Int("withdraw_id", withdraw_id))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore withdraw",
		}
	}

	withdrawResponse := s.mapping.ToWithdrawResponse(res)

	s.logger.Debug("Successfully restored withdraw", zap.Int("withdraw_id", withdraw_id))

	return withdrawResponse, nil
}

func (s *withdrawService) DeleteWithdrawPermanent(withdraw_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting withdraw permanently", zap.Int("withdraw_id", withdraw_id))

	_, err := s.withdrawRepository.DeleteWithdrawPermanent(withdraw_id)

	if err != nil {
		s.logger.Error("Failed to delete withdraw permanently", zap.Error(err), zap.Int("withdraw_id", withdraw_id))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete withdraw permanently",
		}
	}

	s.logger.Debug("Successfully deleted withdraw permanently", zap.Int("withdraw_id", withdraw_id))

	return true, nil
}

func (s *withdrawService) RestoreAllWithdraw() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all withdraws")

	_, err := s.withdrawRepository.RestoreAllWithdraw()

	if err != nil {
		s.logger.Error("Failed to restore all withdraws", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all withdraws: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully restored all withdraws")
	return true, nil
}

func (s *withdrawService) DeleteAllWithdrawPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all withdraws")

	_, err := s.withdrawRepository.DeleteAllWithdrawPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all withdraws", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all withdraws: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully deleted all withdraws permanently")
	return true, nil
}
