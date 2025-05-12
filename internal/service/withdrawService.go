package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/service"
	"net/http"

	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/withdraw_errors"
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

func (s *withdrawService) FindAll(req *requests.FindAllWithdraws) ([]*response.WithdrawResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

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

	withdraws, totalRecords, err := s.withdrawRepository.FindAll(req)

	if err != nil {
		s.logger.Error("Failed to fetch withdraw",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, withdraw_errors.ErrFailedFindAllWithdraws
	}

	withdrawResponse := s.mapping.ToWithdrawsResponse(withdraws)

	s.logger.Debug("Successfully fetched withdraw",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdrawResponse, totalRecords, nil
}

func (s *withdrawService) FindAllByCardNumber(req *requests.FindAllWithdrawCardNumber) ([]*response.WithdrawResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

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

	withdraws, totalRecords, err := s.withdrawRepository.FindAllByCardNumber(req)

	if err != nil {
		s.logger.Error("Failed to fetch withdraw",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, withdraw_errors.ErrFailedFindAllWithdrawsByCard
	}

	withdrawResponse := s.mapping.ToWithdrawsResponse(withdraws)

	s.logger.Debug("Successfully fetched withdraw",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdrawResponse, totalRecords, nil
}

func (s *withdrawService) FindById(withdrawID int) (*response.WithdrawResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching withdraw by ID", zap.Int("withdraw_id", withdrawID))

	withdraw, err := s.withdrawRepository.FindById(withdrawID)

	if err != nil {
		s.logger.Error("failed to find withdraw by id", zap.Error(err))
		return nil, withdraw_errors.ErrWithdrawNotFound
	}
	so := s.mapping.ToWithdrawResponse(withdraw)

	s.logger.Debug("Successfully fetched withdraw", zap.Int("withdraw_id", withdrawID))

	return so, nil
}

func (s *withdrawService) FindMonthWithdrawStatusSuccess(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	s.logger.Debug("Fetching monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month))

	records, err := s.withdrawRepository.GetMonthWithdrawStatusSuccess(req)

	if err != nil {
		s.logger.Error("failed to fetch monthly Withdraw status success", zap.Error(err))

		return nil, withdraw_errors.ErrFailedFindMonthWithdrawStatusSuccess
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

		return nil, withdraw_errors.ErrFailedFindYearWithdrawStatusSuccess
	}

	s.logger.Debug("Successfully fetched yearly Withdraw status success", zap.Int("year", year))

	so := s.mapping.ToWithdrawResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *withdrawService) FindMonthWithdrawStatusFailed(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	s.logger.Debug("Fetching monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month))

	records, err := s.withdrawRepository.GetMonthWithdrawStatusFailed(req)

	if err != nil {
		s.logger.Error("failed to fetch monthly Withdraw status Failed", zap.Error(err))

		return nil, withdraw_errors.ErrFailedFindMonthWithdrawStatusFailed
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

		return nil, withdraw_errors.ErrFailedFindYearWithdrawStatusFailed
	}

	s.logger.Debug("Failedfully fetched yearly Withdraw status Failed", zap.Int("year", year))

	so := s.mapping.ToWithdrawResponsesYearStatusFailed(records)

	return so, nil
}

func (s *withdrawService) FindMonthWithdrawStatusSuccessByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse) {
	year := req.Year
	card_number := req.CardNumber
	month := req.Month

	s.logger.Debug("Fetching monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.withdrawRepository.GetMonthWithdrawStatusSuccessByCardNumber(req)

	if err != nil {
		s.logger.Error("failed to fetch monthly Withdraw status success", zap.Error(err))

		return nil, withdraw_errors.ErrFailedFindMonthWithdrawStatusSuccess
	}

	s.logger.Debug("Successfully fetched monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToWithdrawResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *withdrawService) FindYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {
	year := req.Year
	card_number := req.CardNumber

	s.logger.Debug("Fetching yearly Withdraw status success", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.withdrawRepository.GetYearlyWithdrawStatusSuccessByCardNumber(req)
	if err != nil {
		s.logger.Error("failed to fetch yearly Withdraw status success", zap.Error(err))

		return nil, withdraw_errors.ErrFailedFindYearWithdrawStatusSuccess
	}

	s.logger.Debug("Successfully fetched yearly Withdraw status success", zap.Int("year", year), zap.String("card_number", card_number))

	so := s.mapping.ToWithdrawResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *withdrawService) FindMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	card_number := req.CardNumber
	month := req.Month

	s.logger.Debug("Fetching monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.withdrawRepository.GetMonthWithdrawStatusFailedByCardNumber(req)

	if err != nil {
		s.logger.Error("failed to fetch monthly Withdraw status Failed", zap.Error(err))

		return nil, withdraw_errors.ErrFailedFindMonthWithdrawStatusFailed
	}

	s.logger.Debug("Failedfully fetched monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToWithdrawResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *withdrawService) FindYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	year := req.Year
	card_number := req.CardNumber

	s.logger.Debug("Fetching yearly Withdraw status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.withdrawRepository.GetYearlyWithdrawStatusFailedByCardNumber(req)
	if err != nil {
		s.logger.Error("failed to fetch yearly Withdraw status Failed", zap.Error(err))

		return nil, withdraw_errors.ErrFailedFindYearWithdrawStatusFailed
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
		return nil, withdraw_errors.ErrFailedFindMonthlyWithdraws
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
		return nil, withdraw_errors.ErrFailedFindYearlyWithdraws
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountYearlyResponses(withdraws)

	s.logger.Debug("Successfully fetched yearly withdraws", zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawService) FindMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching monthly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	withdraws, err := s.withdrawRepository.GetMonthlyWithdrawsByCardNumber(req)
	if err != nil {
		s.logger.Error("failed to find monthly withdraws by card number", zap.Error(err))
		return nil, withdraw_errors.ErrFailedFindMonthlyWithdraws
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountMonthlyResponses(withdraws)

	s.logger.Debug("Successfully fetched monthly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawService) FindYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching yearly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	withdraws, err := s.withdrawRepository.GetYearlyWithdrawsByCardNumber(req)
	if err != nil {
		s.logger.Error("failed to find yearly withdraws by card number", zap.Error(err))
		return nil, withdraw_errors.ErrFailedFindYearlyWithdraws
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountYearlyResponses(withdraws)

	s.logger.Debug("Successfully fetched yearly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawService) FindByActive(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

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

	withdraws, totalRecords, err := s.withdrawRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to fetch active withdraw",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, withdraw_errors.ErrFailedFindActiveWithdraws
	}

	withdrawResponses := s.mapping.ToWithdrawsResponseDeleteAt(withdraws)

	s.logger.Debug("Successfully fetched active withdraw",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdrawResponses, totalRecords, nil
}

func (s *withdrawService) FindByTrashed(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

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

	withdraws, totalRecords, err := s.withdrawRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to fetch trashed withdraw",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, withdraw_errors.ErrFailedFindTrashedWithdraws
	}

	withdrawResponses := s.mapping.ToWithdrawsResponseDeleteAt(withdraws)

	s.logger.Debug("Successfully fetched trashed withdraw",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdrawResponses, totalRecords, nil
}

func (s *withdrawService) Create(request *requests.CreateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new withdraw", zap.Any("request", request))

	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)

	if err != nil {
		s.logger.Error("Failed to find saldo by user ID", zap.Error(err))
		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	if saldo == nil {
		s.logger.Error("Saldo not found for user", zap.String("cardNumber", request.CardNumber))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Saldo not found for the specified user ID.",
			Code:    http.StatusNotFound,
		}
	}
	if saldo.TotalBalance < request.WithdrawAmount {
		s.logger.Error("Insufficient balance for user", zap.String("cardNumber", request.CardNumber), zap.Int("requested", request.WithdrawAmount))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance for withdrawal.",
			Code:    http.StatusBadRequest,
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
		return nil, saldo_errors.ErrFailedUpdateSaldo
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
		return nil, withdraw_errors.ErrFailedCreateWithdraw
	}
	if _, err := s.withdrawRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
		WithdrawID: withdrawRecord.ID,
		Status:     "success",
	}); err != nil {
		s.logger.Error("Failed to update withdraw status", zap.Error(err))
		return nil, withdraw_errors.ErrFailedUpdateWithdraw
	}
	so := s.mapping.ToWithdrawResponse(withdrawRecord)
	s.logger.Debug("Successfully created withdraw", zap.Int("withdraw_id", withdrawRecord.ID))
	return so, nil
}

func (s *withdrawService) Update(request *requests.UpdateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating withdraw", zap.Int("withdraw_id", *request.WithdrawID), zap.Any("request", request))
	_, err := s.withdrawRepository.FindById(*request.WithdrawID)
	if err != nil {
		s.logger.Error("Failed to find withdraw record by ID", zap.Error(err))
		return nil, withdraw_errors.ErrWithdrawNotFound
	}
	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)
	if err != nil {
		s.logger.Error("Failed to fetch saldo by user ID", zap.Error(err))
		return nil, saldo_errors.ErrFailedSaldoNotFound
	}
	if saldo.TotalBalance < request.WithdrawAmount {
		s.logger.Error("Insufficient balance for user", zap.String("cardNumber", request.CardNumber))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance for withdrawal update.",
			Code:    http.StatusBadRequest,
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
			WithdrawID: *request.WithdrawID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update withdraw status", zap.Error(err))
		}
		return nil, withdraw_errors.ErrFailedUpdateWithdraw
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
			WithdrawID: *request.WithdrawID,
			Status:     "failed",
		}); err != nil {
			s.logger.Error("Failed to update withdraw status", zap.Error(err))
		}
		return nil, withdraw_errors.ErrFailedUpdateWithdraw
	}
	if _, err := s.withdrawRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
		WithdrawID: updatedWithdraw.ID,
		Status:     "success",
	}); err != nil {
		s.logger.Error("Failed to update withdraw status", zap.Error(err))
		return nil, withdraw_errors.ErrFailedUpdateWithdraw
	}
	so := s.mapping.ToWithdrawResponse(updatedWithdraw)
	s.logger.Debug("Successfully updated withdraw", zap.Int("withdraw_id", so.ID))
	return so, nil
}

func (s *withdrawService) TrashedWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse) {
	s.logger.Debug("Trashing withdraw", zap.Int("withdraw_id", withdraw_id))

	res, err := s.withdrawRepository.TrashedWithdraw(withdraw_id)

	if err != nil {
		s.logger.Error("Failed to move withdraw to trash",
			zap.Int("withdraw_id", withdraw_id),
			zap.Error(err))

		return nil, withdraw_errors.ErrFailedTrashedWithdraw
	}

	withdrawResponse := s.mapping.ToWithdrawResponse(res)

	s.logger.Debug("Successfully trashed withdraw", zap.Int("withdraw_id", withdraw_id))

	return withdrawResponse, nil
}

func (s *withdrawService) RestoreWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse) {
	s.logger.Debug("Restoring withdraw", zap.Int("withdraw_id", withdraw_id))

	res, err := s.withdrawRepository.RestoreWithdraw(withdraw_id)

	if err != nil {
		s.logger.Error("Failed to restore withdraw from trash",
			zap.Int("withdraw_id", withdraw_id),
			zap.Error(err))

		return nil, withdraw_errors.ErrFailedRestoreWithdraw
	}

	withdrawResponse := s.mapping.ToWithdrawResponse(res)

	s.logger.Debug("Successfully restored withdraw", zap.Int("withdraw_id", withdraw_id))

	return withdrawResponse, nil
}

func (s *withdrawService) DeleteWithdrawPermanent(withdraw_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting withdraw permanently", zap.Int("withdraw_id", withdraw_id))

	_, err := s.withdrawRepository.DeleteWithdrawPermanent(withdraw_id)

	if err != nil {
		s.logger.Error("Failed to permanently delete withdraw",
			zap.Int("withdraw_id", withdraw_id),
			zap.Error(err))

		return false, withdraw_errors.ErrFailedDeleteWithdrawPermanent
	}

	s.logger.Debug("Successfully deleted withdraw permanently", zap.Int("withdraw_id", withdraw_id))

	return true, nil
}

func (s *withdrawService) RestoreAllWithdraw() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all withdraws")

	_, err := s.withdrawRepository.RestoreAllWithdraw()

	if err != nil {
		s.logger.Error("Failed to restore all withdraws", zap.Error(err))
		return false, withdraw_errors.ErrFailedRestoreAllWithdraw
	}

	s.logger.Debug("Successfully restored all withdraws")
	return true, nil
}

func (s *withdrawService) DeleteAllWithdrawPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all withdraws")

	_, err := s.withdrawRepository.DeleteAllWithdrawPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all withdraws", zap.Error(err))
		return false, withdraw_errors.ErrFailedDeleteAllWithdrawPermanent
	}

	s.logger.Debug("Successfully deleted all withdraws permanently")
	return true, nil
}
