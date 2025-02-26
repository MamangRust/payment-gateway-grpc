package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/service"

	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type topupService struct {
	cardRepository  repository.CardRepository
	topupRepository repository.TopupRepository
	saldoRepository repository.SaldoRepository
	logger          logger.LoggerInterface
	mapping         responseservice.TopupResponseMapper
}

func NewTopupService(cardRepository repository.CardRepository,
	topupRepository repository.TopupRepository,
	saldoRepository repository.SaldoRepository,
	logger logger.LoggerInterface, mapping responseservice.TopupResponseMapper) *topupService {
	return &topupService{
		topupRepository: topupRepository,
		saldoRepository: saldoRepository,
		cardRepository:  cardRepository,
		logger:          logger,
		mapping:         mapping,
	}
}

func (s *topupService) FindAll(page int, pageSize int, search string) ([]*response.TopupResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching topup",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	topups, totalRecords, err := s.topupRepository.FindAllTopups(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch topup",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch topups",
		}
	}

	so := s.mapping.ToTopupResponses(topups)

	s.logger.Debug("Successfully fetched topup",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *topupService) FindAllByCardNumber(card_number string, page int, pageSize int, search string) ([]*response.TopupResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching topup",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	topups, totalRecords, err := s.topupRepository.FindAllTopupByCardNumber(card_number, search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch topup",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch topups",
		}
	}

	so := s.mapping.ToTopupResponses(topups)

	s.logger.Debug("Successfully fetched topup",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *topupService) FindById(topupID int) (*response.TopupResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching topup by ID", zap.Int("topup_id", topupID))

	topup, err := s.topupRepository.FindById(topupID)

	if err != nil {
		s.logger.Error("failed to find topup by id", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Topup record not found",
		}
	}

	so := s.mapping.ToTopupResponse(topup)

	s.logger.Debug("Successfully fetched topup", zap.Int("topup_id", topupID))

	return so, nil
}

func (s *topupService) FindMonthTopupStatusSuccess(year int, month int) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly topup status success", zap.Int("year", year), zap.Int("month", month))

	records, err := s.topupRepository.GetMonthTopupStatusSuccess(year, month)
	if err != nil {
		s.logger.Error("failed to fetch monthly topup status success", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup status success",
		}
	}

	s.logger.Debug("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTopupResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *topupService) FindYearlyTopupStatusSuccess(year int) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly topup status success", zap.Int("year", year))

	records, err := s.topupRepository.GetYearlyTopupStatusSuccess(year)
	if err != nil {
		s.logger.Error("failed to fetch yearly topup status success", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup status success",
		}
	}

	s.logger.Debug("Successfully fetched yearly topup status success", zap.Int("year", year))

	so := s.mapping.ToTopupResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *topupService) FindMonthTopupStatusFailed(year int, month int) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly topup status Failed", zap.Int("year", year), zap.Int("month", month))

	records, err := s.topupRepository.GetMonthTopupStatusFailed(year, month)
	if err != nil {
		s.logger.Error("failed to fetch monthly topup status Failed", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup status Failed",
		}
	}

	s.logger.Debug("Failedfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTopupResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *topupService) FindYearlyTopupStatusFailed(year int) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly topup status Failed", zap.Int("year", year))

	records, err := s.topupRepository.GetYearlyTopupStatusFailed(year)
	if err != nil {
		s.logger.Error("failed to fetch yearly topup status Failed", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup status Failed",
		}
	}

	s.logger.Debug("Failedfully fetched yearly topup status Failed", zap.Int("year", year))

	so := s.mapping.ToTopupResponsesYearStatusFailed(records)

	return so, nil
}

func (s *topupService) FindMonthTopupStatusSuccessByCardNumber(card_number string, year int, month int) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly topup status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.topupRepository.GetMonthTopupStatusSuccessByCardNumber(card_number, year, month)
	if err != nil {
		s.logger.Error("failed to fetch monthly topup status success", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup status success",
		}
	}

	s.logger.Debug("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTopupResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *topupService) FindYearlyTopupStatusSuccessByCardNumber(card_number string, year int) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly topup status success", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.topupRepository.GetYearlyTopupStatusSuccessByCardNumber(card_number, year)
	if err != nil {
		s.logger.Error("failed to fetch yearly topup status success", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup status success",
		}
	}

	s.logger.Debug("Successfully fetched yearly topup status success", zap.Int("year", year))

	so := s.mapping.ToTopupResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *topupService) FindMonthTopupStatusFailedByCardNumber(card_number string, year int, month int) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly topup status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.topupRepository.GetMonthTopupStatusFailedByCardNumber(card_number, year, month)
	if err != nil {
		s.logger.Error("failed to fetch monthly topup status Failed", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup status Failed",
		}
	}

	s.logger.Debug("Failedfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTopupResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *topupService) FindYearlyTopupStatusFailedByCardNumber(card_number string, year int) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly topup status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.topupRepository.GetYearlyTopupStatusFailedByCardNumber(card_number, year)
	if err != nil {
		s.logger.Error("failed to fetch yearly topup status Failed", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup status Failed",
		}
	}

	s.logger.Debug("Failedfully fetched yearly topup status Failed", zap.Int("year", year))

	so := s.mapping.ToTopupResponsesYearStatusFailed(records)

	return so, nil
}

func (s *topupService) FindMonthlyTopupMethods(year int) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly topup methods", zap.Int("year", year))

	records, err := s.topupRepository.GetMonthlyTopupMethods(year)
	if err != nil {
		s.logger.Error("Failed to fetch monthly topup methods", zap.Error(err), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup methods",
		}
	}

	responses := s.mapping.ToTopupMonthlyMethodResponses(records)

	s.logger.Debug("Successfully fetched monthly topup methods", zap.Int("year", year))

	return responses, nil
}

func (s *topupService) FindYearlyTopupMethods(year int) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly topup methods", zap.Int("year", year))

	records, err := s.topupRepository.GetYearlyTopupMethods(year)
	if err != nil {
		s.logger.Error("Failed to fetch yearly topup methods", zap.Error(err), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup methods",
		}
	}

	responses := s.mapping.ToTopupYearlyMethodResponses(records)

	s.logger.Debug("Successfully fetched yearly topup methods", zap.Int("year", year))

	return responses, nil
}

func (s *topupService) FindMonthlyTopupAmounts(year int) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly topup amounts", zap.Int("year", year))

	records, err := s.topupRepository.GetMonthlyTopupAmounts(year)
	if err != nil {
		s.logger.Error("Failed to fetch monthly topup amounts", zap.Error(err), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup amounts",
		}
	}

	responses := s.mapping.ToTopupMonthlyAmountResponses(records)

	s.logger.Debug("Successfully fetched monthly topup amounts", zap.Int("year", year))

	return responses, nil
}

func (s *topupService) FindYearlyTopupAmounts(year int) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly topup amounts", zap.Int("year", year))

	records, err := s.topupRepository.GetYearlyTopupAmounts(year)
	if err != nil {
		s.logger.Error("Failed to fetch yearly topup amounts", zap.Error(err), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup amounts",
		}
	}

	responses := s.mapping.ToTopupYearlyAmountResponses(records)

	s.logger.Debug("Successfully fetched yearly topup amounts", zap.Int("year", year))

	return responses, nil
}

func (s *topupService) FindMonthlyTopupMethodsByCardNumber(cardNumber string, year int) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.topupRepository.GetMonthlyTopupMethodsByCardNumber(cardNumber, year)
	if err != nil {
		s.logger.Error("Failed to fetch monthly topup methods by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup methods by card number",
		}
	}

	responses := s.mapping.ToTopupMonthlyMethodResponses(records)

	s.logger.Debug("Successfully fetched monthly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupService) FindYearlyTopupMethodsByCardNumber(cardNumber string, year int) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.topupRepository.GetYearlyTopupMethodsByCardNumber(cardNumber, year)
	if err != nil {
		s.logger.Error("Failed to fetch yearly topup methods by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup methods by card number",
		}
	}

	responses := s.mapping.ToTopupYearlyMethodResponses(records)

	s.logger.Debug("Successfully fetched yearly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupService) FindMonthlyTopupAmountsByCardNumber(cardNumber string, year int) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching monthly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.topupRepository.GetMonthlyTopupAmountsByCardNumber(cardNumber, year)
	if err != nil {
		s.logger.Error("Failed to fetch monthly topup amounts by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch monthly topup amounts by card number",
		}
	}

	responses := s.mapping.ToTopupMonthlyAmountResponses(records)

	s.logger.Debug("Successfully fetched monthly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupService) FindYearlyTopupAmountsByCardNumber(cardNumber string, year int) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching yearly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.topupRepository.GetYearlyTopupAmountsByCardNumber(cardNumber, year)
	if err != nil {
		s.logger.Error("Failed to fetch yearly topup amounts by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch yearly topup amounts by card number",
		}
	}

	responses := s.mapping.ToTopupYearlyAmountResponses(records)

	s.logger.Debug("Successfully fetched yearly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupService) FindByActive(page int, pageSize int, search string) ([]*response.TopupResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching active topup",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	topups, totalRecords, err := s.topupRepository.FindByActive(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch active topup",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find active top-up records",
		}
	}

	so := s.mapping.ToTopupResponsesDeleteAt(topups)

	s.logger.Debug("Successfully fetched active topup",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *topupService) FindByTrashed(page int, pageSize int, search string) ([]*response.TopupResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching trashed topup",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	topups, totalRecords, err := s.topupRepository.FindByTrashed(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch trashed topup",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find trashed top-up records",
		}
	}

	so := s.mapping.ToTopupResponsesDeleteAt(topups)

	s.logger.Debug("Successfully fetched trashed topup",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *topupService) CreateTopup(request *requests.CreateTopupRequest) (*response.TopupResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting CreateTopup process",
		zap.String("cardNumber", request.CardNumber),
		zap.Float64("topupAmount", float64(request.TopupAmount)),
	)

	card, err := s.cardRepository.FindCardByCardNumber(request.CardNumber)

	if err != nil {
		s.logger.Error("failed to find card by number", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Card not found",
		}
	}

	topup, err := s.topupRepository.CreateTopup(request)

	if err != nil {
		s.logger.Error("failed to create topup", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create topup record",
		}
	}

	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)
	if err != nil {
		s.logger.Error("failed to find saldo by user id", zap.Error(err))

		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}

		s.topupRepository.UpdateTopupStatus(&req)
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch user's saldo",
		}
	}

	newBalance := saldo.TotalBalance + request.TopupAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: newBalance,
	})
	if err != nil {
		s.logger.Error("failed to update saldo balance", zap.Error(err))
		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}

		s.topupRepository.UpdateTopupStatus(&req)
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update saldo balance",
		}
	}

	expireDate, err := time.Parse("2006-01-02", card.ExpireDate)
	if err != nil {
		s.logger.Error("failed to parse expire date", zap.Error(err))

		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}

		s.topupRepository.UpdateTopupStatus(&req)
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Invalid expire date format",
		}
	}

	_, err = s.cardRepository.UpdateCard(&requests.UpdateCardRequest{
		CardID:       card.ID,
		UserID:       card.UserID,
		CardType:     card.CardType,
		ExpireDate:   expireDate,
		CVV:          card.CVV,
		CardProvider: card.CardProvider,
	})
	if err != nil {
		s.logger.Error("failed to update card expire date", zap.Error(err))
		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}

		s.topupRepository.UpdateTopupStatus(&req)
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update card expire date",
		}
	}

	req := requests.UpdateTopupStatus{
		TopupID: topup.ID,
		Status:  "success",
	}

	_, err = s.topupRepository.UpdateTopupStatus(&req)
	if err != nil {
		s.logger.Error("failed to update topup status", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update topup status to success",
		}
	}

	so := s.mapping.ToTopupResponse(topup)

	s.logger.Debug("CreateTopup process completed",
		zap.String("cardNumber", request.CardNumber),
		zap.Float64("topupAmount", float64(request.TopupAmount)),
		zap.Float64("newBalance", float64(newBalance)),
	)

	return so, nil
}

func (s *topupService) UpdateTopup(request *requests.UpdateTopupRequest) (*response.TopupResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting UpdateTopup process",
		zap.String("cardNumber", request.CardNumber),
		zap.Int("topupID", request.TopupID),
		zap.Float64("newTopupAmount", float64(request.TopupAmount)),
	)

	_, err := s.cardRepository.FindCardByCardNumber(request.CardNumber)

	if err != nil {
		s.logger.Error("failed to find card by number", zap.Error(err))

		req := requests.UpdateTopupStatus{
			TopupID: request.TopupID,
			Status:  "failed",
		}

		s.topupRepository.UpdateTopupStatus(&req)

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Card not found",
		}
	}

	existingTopup, err := s.topupRepository.FindById(request.TopupID)
	if err != nil || existingTopup == nil {
		s.logger.Error("Failed to find topup by ID", zap.Error(err))

		req := requests.UpdateTopupStatus{
			TopupID: request.TopupID,
			Status:  "failed",
		}

		s.topupRepository.UpdateTopupStatus(&req)
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Topup not found",
		}
	}

	topupDifference := request.TopupAmount - existingTopup.TopupAmount

	_, err = s.topupRepository.UpdateTopup(request)
	if err != nil {
		s.logger.Error("Failed to update topup amount", zap.Error(err))

		req := requests.UpdateTopupStatus{
			TopupID: request.TopupID,
			Status:  "failed",
		}

		s.topupRepository.UpdateTopupStatus(&req)
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to update topup amount: %v", err),
		}
	}

	currentSaldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)
	if err != nil {
		s.logger.Error("Failed to retrieve current saldo", zap.Error(err))

		req := requests.UpdateTopupStatus{
			TopupID: request.TopupID,
			Status:  "failed",
		}

		s.topupRepository.UpdateTopupStatus(&req)

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to retrieve current saldo: %v", err),
		}
	}

	if currentSaldo == nil {
		s.logger.Error("No saldo found for card number", zap.String("card_number", request.CardNumber))

		req := requests.UpdateTopupStatus{
			TopupID: request.TopupID,
			Status:  "failed",
		}

		s.topupRepository.UpdateTopupStatus(&req)

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Card not found",
		}
	}

	newBalance := currentSaldo.TotalBalance + topupDifference
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: newBalance,
	})
	if err != nil {
		s.logger.Error("Failed to update saldo balance", zap.Error(err))

		_, rollbackErr := s.topupRepository.UpdateTopupAmount(&requests.UpdateTopupAmount{
			TopupID:     request.TopupID,
			TopupAmount: existingTopup.TopupAmount,
		})
		if rollbackErr != nil {
			s.logger.Error("Failed to rollback topup update", zap.Error(rollbackErr))
		}

		req := requests.UpdateTopupStatus{
			TopupID: request.TopupID,
			Status:  "failed",
		}

		s.topupRepository.UpdateTopupStatus(&req)
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to update saldo balance: %v", err),
		}
	}

	updatedTopup, err := s.topupRepository.FindById(request.TopupID)

	if err != nil || updatedTopup == nil {
		s.logger.Error("Failed to find updated topup by ID", zap.Error(err))

		req := requests.UpdateTopupStatus{
			TopupID: request.TopupID,
			Status:  "failed",
		}

		s.topupRepository.UpdateTopupStatus(&req)
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Updated topup not found",
		}
	}

	req := requests.UpdateTopupStatus{
		TopupID: request.TopupID,
		Status:  "success",
	}

	_, err = s.topupRepository.UpdateTopupStatus(&req)
	if err != nil {
		s.logger.Error("Failed to update topup status", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update topup status to success",
		}
	}

	so := s.mapping.ToTopupResponse(updatedTopup)

	s.logger.Debug("UpdateTopup process completed",
		zap.String("cardNumber", request.CardNumber),
		zap.Int("topupID", request.TopupID),
		zap.Float64("newTopupAmount", float64(request.TopupAmount)),
		zap.Float64("newBalance", float64(newBalance)),
	)

	return so, nil
}

func (s *topupService) TrashedTopup(topup_id int) (*response.TopupResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting TrashedTopup process",
		zap.Int("topupID", topup_id),
	)

	res, err := s.topupRepository.TrashedTopup(topup_id)

	if err != nil {
		s.logger.Error("Failed to trash topup", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to trash topup: %v", err),
		}
	}

	so := s.mapping.ToTopupResponse(res)

	s.logger.Debug("TrashedTopup process completed",
		zap.Int("topupID", topup_id),
	)

	return so, nil
}

func (s *topupService) RestoreTopup(topup_id int) (*response.TopupResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting RestoreTopup process",
		zap.Int("topupID", topup_id),
	)

	res, err := s.topupRepository.RestoreTopup(topup_id)

	if err != nil {
		s.logger.Error("Failed to restore topup", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to restore topup: %v", err),
		}
	}

	so := s.mapping.ToTopupResponse(res)

	s.logger.Debug("RestoreTopup process completed",
		zap.Int("topupID", topup_id),
	)

	return so, nil
}

func (s *topupService) DeleteTopupPermanent(topup_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Starting DeleteTopupPermanent process",
		zap.Int("topupID", topup_id),
	)

	_, err := s.topupRepository.DeleteTopupPermanent(topup_id)

	if err != nil {
		s.logger.Error("Failed to delete topup permanently", zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to delete topup permanently: %v", err),
		}
	}

	s.logger.Debug("DeleteTopupPermanent process completed",
		zap.Int("topupID", topup_id),
	)

	return true, nil
}

func (s *topupService) RestoreAllTopup() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all topups")

	_, err := s.topupRepository.RestoreAllTopup()

	if err != nil {
		s.logger.Error("Failed to restore all topups", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all topups: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully restored all topups")
	return true, nil
}

func (s *topupService) DeleteAllTopupPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all topups")

	_, err := s.topupRepository.DeleteAllTopupPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all topups", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all topups: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully deleted all topups permanently")
	return true, nil
}
