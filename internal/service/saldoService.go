package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/service"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/pkg/errors/card_errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"

	"go.uber.org/zap"
)

type saldoService struct {
	cardRepository  repository.CardRepository
	saldoRepository repository.SaldoRepository
	logger          logger.LoggerInterface
	mapping         responseservice.SaldoResponseMapper
}

func NewSaldoService(saldo repository.SaldoRepository, card repository.CardRepository, logger logger.LoggerInterface, mapping responseservice.SaldoResponseMapper) *saldoService {
	return &saldoService{
		saldoRepository: saldo,
		cardRepository:  card,
		logger:          logger,
		mapping:         mapping,
	}
}

func (s *saldoService) FindAll(req *requests.FindAllSaldos) ([]*response.SaldoResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching saldo",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Debug("Fetching all saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	res, totalRecords, err := s.saldoRepository.FindAllSaldos(req)

	if err != nil {
		s.logger.Error("Failed to fetch saldo",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, saldo_errors.ErrFailedFindAllSaldos
	}

	so := s.mapping.ToSaldoResponses(res)

	s.logger.Debug("Successfully fetched saldo",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return so, totalRecords, nil
}

func (s *saldoService) FindByActive(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching active saldo",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.saldoRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active saldo",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search))

		return nil, nil, saldo_errors.ErrFailedFindActiveSaldos
	}

	so := s.mapping.ToSaldoResponsesDeleteAt(res)

	s.logger.Debug("Successfully fetched active saldo",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *saldoService) FindByTrashed(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching saldo record",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.saldoRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed saldo",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search))

		return nil, nil, saldo_errors.ErrFailedFindTrashedSaldos
	}

	s.logger.Debug("Successfully fetched trashed saldo",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	so := s.mapping.ToSaldoResponsesDeleteAt(res)

	return so, totalRecords, nil
}

func (s *saldoService) FindById(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching saldo record by ID", zap.Int("saldo_id", saldo_id))

	res, err := s.saldoRepository.FindById(saldo_id)

	if err != nil {
		s.logger.Error("Failed to retrieve saldo details",
			zap.Int("saldo_id", saldo_id),
			zap.Error(err))

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	so := s.mapping.ToSaldoResponse(res)

	s.logger.Debug("Successfully fetched saldo", zap.Int("saldo_id", saldo_id))

	return so, nil
}

func (s *saldoService) FindMonthlyTotalSaldoBalance(req *requests.MonthTotalSaldoBalance) ([]*response.SaldoMonthTotalBalanceResponse, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	s.logger.Debug("Fetching monthly total saldo balance", zap.Int("year", year), zap.Int("month", month))

	res, err := s.saldoRepository.GetMonthlyTotalSaldoBalance(req)
	if err != nil {
		s.logger.Error("Failed to fetch monthly total saldo balance", zap.Error(err), zap.Int("year", year), zap.Int("month", month))

		return nil, saldo_errors.ErrFailedFindMonthlyTotalSaldoBalance
	}

	responses := s.mapping.ToSaldoMonthTotalBalanceResponses(res)

	s.logger.Debug("Successfully fetched monthly total saldo balance", zap.Int("year", year), zap.Int("month", month))

	return responses, nil
}

func (s *saldoService) FindYearTotalSaldoBalance(year int) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly total saldo balance", zap.Int("year", year))

	res, err := s.saldoRepository.GetYearTotalSaldoBalance(year)

	if err != nil {
		s.logger.Error("Failed to fetch yearly total saldo balance", zap.Error(err), zap.Int("year", year))

		return nil, saldo_errors.ErrFailedFindYearTotalSaldoBalance
	}

	s.logger.Debug("Successfully fetched yearly total saldo balance", zap.Int("year", year))

	so := s.mapping.ToSaldoYearTotalBalanceResponses(res)

	return so, nil
}

func (s *saldoService) FindMonthlySaldoBalances(year int) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly saldo balances", zap.Int("year", year))

	res, err := s.saldoRepository.GetMonthlySaldoBalances(year)

	if err != nil {
		s.logger.Error("Failed to fetch monthly saldo balances", zap.Error(err), zap.Int("year", year))

		return nil, saldo_errors.ErrFailedFindMonthlySaldoBalances
	}

	responses := s.mapping.ToSaldoMonthBalanceResponses(res)

	s.logger.Debug("Successfully fetched monthly saldo balances", zap.Int("year", year))

	return responses, nil
}

func (s *saldoService) FindYearlySaldoBalances(year int) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly saldo balances", zap.Int("year", year))

	res, err := s.saldoRepository.GetYearlySaldoBalances(year)

	if err != nil {
		s.logger.Error("Failed to fetch yearly saldo balances", zap.Error(err), zap.Int("year", year))

		return nil, saldo_errors.ErrFailedFindYearlySaldoBalances
	}

	responses := s.mapping.ToSaldoYearBalanceResponses(res)

	s.logger.Debug("Successfully fetched yearly saldo balances", zap.Int("year", year))

	return responses, nil
}

func (s *saldoService) FindByCardNumber(card_number string) (*response.SaldoResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching saldo record by card number", zap.String("card_number", card_number))

	res, err := s.saldoRepository.FindByCardNumber(card_number)

	if err != nil {
		s.logger.Error("Failed to retrieve saldo details",
			zap.String("card_number", card_number),
			zap.Error(err))

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	so := s.mapping.ToSaldoResponse(res)

	s.logger.Debug("Successfully fetched saldo by card number", zap.String("card_number", card_number))

	return so, nil
}

func (s *saldoService) CreateSaldo(request *requests.CreateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating saldo record", zap.String("card_number", request.CardNumber))

	_, err := s.cardRepository.FindCardByCardNumber(request.CardNumber)

	if err != nil {
		s.logger.Error("Failed to create new saldo",
			zap.Error(err),
			zap.Any("request", request))

		return nil, card_errors.ErrCardNotFoundRes
	}

	res, err := s.saldoRepository.CreateSaldo(request)

	if err != nil {
		s.logger.Error("Failed to create saldo record",
			zap.Error(err))

		return nil, saldo_errors.ErrFailedCreateSaldo
	}

	so := s.mapping.ToSaldoResponse(res)

	s.logger.Debug("Successfully created saldo record", zap.String("card_number", request.CardNumber))

	return so, nil
}

func (s *saldoService) UpdateSaldo(request *requests.UpdateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating saldo record", zap.String("card_number", request.CardNumber), zap.Float64("amount", float64(request.TotalBalance)))

	_, err := s.cardRepository.FindCardByCardNumber(request.CardNumber)

	if err != nil {
		s.logger.Error("Card not found for card number update",
			zap.String("card_number", request.CardNumber),
			zap.Error(err))

		return nil, card_errors.ErrCardNotFoundRes
	}

	res, err := s.saldoRepository.UpdateSaldo(request)

	if err != nil {
		s.logger.Error("Failed to update saldo", zap.Error(err), zap.String("card_number", request.CardNumber))
		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	so := s.mapping.ToSaldoResponse(res)

	s.logger.Debug("Successfully updated saldo", zap.String("card_number", request.CardNumber), zap.Int("saldo_id", res.ID))

	return so, nil
}

func (s *saldoService) TrashSaldo(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse) {
	s.logger.Debug("Trashing saldo record", zap.Int("saldo_id", saldo_id))

	res, err := s.saldoRepository.TrashedSaldo(saldo_id)

	if err != nil {
		s.logger.Error("Failed to move saldo to trash",
			zap.Int("saldo", saldo_id),
			zap.Error(err))

		return nil, saldo_errors.ErrFailedTrashSaldo
	}
	so := s.mapping.ToSaldoResponse(res)

	s.logger.Debug("Successfully trashed saldo", zap.Int("saldo_id", saldo_id))

	return so, nil
}

func (s *saldoService) RestoreSaldo(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse) {
	s.logger.Debug("Restoring saldo record from trash", zap.Int("saldo_id", saldo_id))

	res, err := s.saldoRepository.RestoreSaldo(saldo_id)

	if err != nil {
		s.logger.Error("Failed to restore saldo from trash",
			zap.Int("saldo_id", saldo_id),
			zap.Error(err))

		return nil, saldo_errors.ErrFailedRestoreSaldo
	}

	so := s.mapping.ToSaldoResponse(res)

	s.logger.Debug("Successfully restored saldo", zap.Int("saldo_id", saldo_id))

	return so, nil
}

func (s *saldoService) DeleteSaldoPermanent(saldo_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting saldo permanently", zap.Int("saldo_id", saldo_id))

	_, err := s.saldoRepository.DeleteSaldoPermanent(saldo_id)

	if err != nil {
		s.logger.Error("Failed to permanently delete saldo",
			zap.Int("saldo_id", saldo_id),
			zap.Error(err))

		return false, saldo_errors.ErrFailedDeleteSaldoPermanent
	}

	s.logger.Debug("Successfully deleted saldo permanently", zap.Int("saldo_id", saldo_id))

	return true, nil
}

func (s *saldoService) RestoreAllSaldo() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all saldo")

	_, err := s.saldoRepository.RestoreAllSaldo()

	if err != nil {
		s.logger.Error("Failed to restore all saldo", zap.Error(err))
		return false, saldo_errors.ErrFailedRestoreAllSaldo
	}

	s.logger.Debug("Successfully restored all saldo")
	return true, nil
}

func (s *saldoService) DeleteAllSaldoPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all saldo")

	_, err := s.saldoRepository.DeleteAllSaldoPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all saldo", zap.Error(err))
		return false, saldo_errors.ErrFailedDeleteAllSaldoPermanent
	}

	s.logger.Debug("Successfully deleted all saldo permanently")
	return true, nil
}
