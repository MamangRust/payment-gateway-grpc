package responseservice

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
)

type cardResponseMapper struct {
}

func NewCardResponseMapper() *cardResponseMapper {
	return &cardResponseMapper{}
}

func (s *cardResponseMapper) ToCardResponse(card *record.CardRecord) *response.CardResponse {
	return &response.CardResponse{
		ID:           card.ID,
		UserID:       card.UserID,
		CardNumber:   card.CardNumber,
		CardType:     card.CardType,
		ExpireDate:   card.ExpireDate,
		CVV:          card.CVV,
		CardProvider: card.CardProvider,
		CreatedAt:    card.CreatedAt,
		UpdatedAt:    card.UpdatedAt,
	}
}

func (s *cardResponseMapper) ToCardsResponse(cards []*record.CardRecord) []*response.CardResponse {
	var response []*response.CardResponse

	for _, card := range cards {
		response = append(response, s.ToCardResponse(card))
	}

	return response
}

func (s *cardResponseMapper) ToCardResponseDeleteAt(card *record.CardRecord) *response.CardResponseDeleteAt {
	return &response.CardResponseDeleteAt{
		ID:           card.ID,
		UserID:       card.UserID,
		CardNumber:   card.CardNumber,
		CardType:     card.CardType,
		ExpireDate:   card.ExpireDate,
		CVV:          card.CVV,
		CardProvider: card.CardProvider,
		CreatedAt:    card.CreatedAt,
		UpdatedAt:    card.UpdatedAt,
		DeletedAt:    card.DeletedAt,
	}
}

func (s *cardResponseMapper) ToCardsResponseDeleteAt(cards []*record.CardRecord) []*response.CardResponseDeleteAt {
	var response []*response.CardResponseDeleteAt

	for _, card := range cards {
		response = append(response, s.ToCardResponseDeleteAt(card))
	}

	return response
}

func (s *cardResponseMapper) ToGetMonthlyBalance(card *record.CardMonthBalance) *response.CardResponseMonthBalance {
	return &response.CardResponseMonthBalance{
		Month:        card.Month,
		TotalBalance: card.TotalBalance,
	}
}

func (s *cardResponseMapper) ToGetMonthlyBalances(cards []*record.CardMonthBalance) []*response.CardResponseMonthBalance {
	var records []*response.CardResponseMonthBalance

	for _, card := range cards {
		records = append(records, s.ToGetMonthlyBalance(card))
	}

	return records
}

func (s *cardResponseMapper) ToGetYearlyBalance(card *record.CardYearlyBalance) *response.CardResponseYearlyBalance {
	return &response.CardResponseYearlyBalance{
		Year:         card.Year,
		TotalBalance: card.TotalBalance,
	}
}

func (s *cardResponseMapper) ToGetYearlyBalances(cards []*record.CardYearlyBalance) []*response.CardResponseYearlyBalance {
	var records []*response.CardResponseYearlyBalance

	for _, card := range cards {
		records = append(records, s.ToGetYearlyBalance(card))
	}

	return records
}

func (s *cardResponseMapper) ToGetMonthlyAmount(card *record.CardMonthAmount) *response.CardResponseMonthAmount {
	return &response.CardResponseMonthAmount{
		Month:       card.Month,
		TotalAmount: card.TotalAmount,
	}
}

func (s *cardResponseMapper) ToGetMonthlyAmounts(cards []*record.CardMonthAmount) []*response.CardResponseMonthAmount {
	var records []*response.CardResponseMonthAmount

	for _, card := range cards {
		records = append(records, s.ToGetMonthlyAmount(card))
	}

	return records
}

func (s *cardResponseMapper) ToGetYearlyAmount(card *record.CardYearAmount) *response.CardResponseYearAmount {
	return &response.CardResponseYearAmount{
		Year:        card.Year,
		TotalAmount: card.TotalAmount,
	}
}

func (s *cardResponseMapper) ToGetYearlyAmounts(cards []*record.CardYearAmount) []*response.CardResponseYearAmount {
	var records []*response.CardResponseYearAmount

	for _, card := range cards {
		records = append(records, s.ToGetYearlyAmount(card))
	}

	return records
}
