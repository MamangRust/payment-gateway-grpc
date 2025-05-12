package apimapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"
)

type cardResponseMapper struct{}

func NewCardResponseMapper() *cardResponseMapper {
	return &cardResponseMapper{}
}

func (s *cardResponseMapper) ToApiResponseCard(card *pb.ApiResponseCard) *response.ApiResponseCard {
	return &response.ApiResponseCard{
		Status:  card.Status,
		Message: card.Message,
		Data:    s.mapCardResponse(card.Data),
	}
}

func (s *cardResponseMapper) ToApiResponsesCard(cards *pb.ApiResponsePaginationCard) *response.ApiResponsePaginationCard {
	return &response.ApiResponsePaginationCard{
		Status:     cards.Status,
		Message:    cards.Message,
		Data:       s.mapCardResponses(cards.Data),
		Pagination: mapPaginationMeta(cards.Pagination),
	}

}

func (s *cardResponseMapper) ToApiResponseCardDeleteAt(card *pb.ApiResponseCardDelete) *response.ApiResponseCardDelete {
	return &response.ApiResponseCardDelete{
		Status:  card.Status,
		Message: card.Message,
	}
}

func (s *cardResponseMapper) ToApiResponseCardAll(card *pb.ApiResponseCardAll) *response.ApiResponseCardAll {
	return &response.ApiResponseCardAll{
		Status:  card.Status,
		Message: card.Message,
	}
}

func (s *cardResponseMapper) ToApiResponsesCardDeletedAt(cards *pb.ApiResponsePaginationCardDeleteAt) *response.ApiResponsePaginationCardDeleteAt {
	return &response.ApiResponsePaginationCardDeleteAt{
		Status:     cards.Status,
		Message:    cards.Message,
		Data:       s.mapCardResponsesDeleteAt(cards.Data),
		Pagination: mapPaginationMeta(cards.Pagination),
	}
}

func (s *cardResponseMapper) ToApiResponseDashboardCard(dash *pb.ApiResponseDashboardCard) *response.ApiResponseDashboardCard {
	return &response.ApiResponseDashboardCard{
		Status:  dash.Status,
		Message: dash.Message,
		Data:    s.mapDashboardCard(dash.Data),
	}
}

func (s *cardResponseMapper) ToApiResponseDashboardCardCardNumber(dash *pb.ApiResponseDashboardCardNumber) *response.ApiResponseDashboardCardNumber {
	return &response.ApiResponseDashboardCardNumber{
		Status:  dash.Status,
		Message: dash.Message,
		Data:    s.mapDashboardCardCardNumber(dash.Data),
	}
}

func (s *cardResponseMapper) ToApiResponseMonthlyBalances(cards *pb.ApiResponseMonthlyBalance) *response.ApiResponseMonthlyBalance {
	return &response.ApiResponseMonthlyBalance{
		Status:  cards.Status,
		Message: cards.Message,
		Data:    s.mapMonthlyBalances(cards.Data),
	}
}

func (s *cardResponseMapper) ToApiResponseYearlyBalances(cards *pb.ApiResponseYearlyBalance) *response.ApiResponseYearlyBalance {
	return &response.ApiResponseYearlyBalance{
		Status:  cards.Status,
		Message: cards.Message,
		Data:    s.mapYearlyBalances(cards.Data),
	}
}

func (s *cardResponseMapper) ToApiResponseMonthlyAmounts(cards *pb.ApiResponseMonthlyAmount) *response.ApiResponseMonthlyAmount {
	return &response.ApiResponseMonthlyAmount{
		Status:  cards.Status,
		Message: cards.Message,
		Data:    s.mapMonthlyAmounts(cards.Data),
	}
}

func (s *cardResponseMapper) ToApiResponseYearlyAmounts(cards *pb.ApiResponseYearlyAmount) *response.ApiResponseYearlyAmount {
	return &response.ApiResponseYearlyAmount{
		Status:  cards.Status,
		Message: cards.Message,
		Data:    s.mapYearlyAmounts(cards.Data),
	}
}

func (s *cardResponseMapper) mapCardResponse(card *pb.CardResponse) *response.CardResponse {
	return &response.CardResponse{
		ID:           int(card.Id),
		UserID:       int(card.UserId),
		CardNumber:   card.CardNumber,
		CardType:     card.CardType,
		ExpireDate:   card.ExpireDate,
		CVV:          card.Cvv,
		CardProvider: card.CardProvider,
		CreatedAt:    card.CreatedAt,
		UpdatedAt:    card.UpdatedAt,
	}
}

func (s *cardResponseMapper) mapCardResponses(cards []*pb.CardResponse) []*response.CardResponse {
	var responseCards []*response.CardResponse

	for _, role := range cards {
		responseCards = append(responseCards, s.mapCardResponse(role))
	}

	return responseCards
}

func (s *cardResponseMapper) mapCardResponseDeleteAt(card *pb.CardResponseDeleteAt) *response.CardResponseDeleteAt {
	var deletedAt string
	if card.DeletedAt != nil {
		deletedAt = card.DeletedAt.Value
	}

	return &response.CardResponseDeleteAt{
		ID:           int(card.Id),
		UserID:       int(card.UserId),
		CardNumber:   card.CardNumber,
		CardType:     card.CardType,
		ExpireDate:   card.ExpireDate,
		CVV:          card.Cvv,
		CardProvider: card.CardProvider,
		CreatedAt:    card.CreatedAt,
		UpdatedAt:    card.UpdatedAt,
		DeletedAt:    &deletedAt,
	}
}

func (s *cardResponseMapper) mapCardResponsesDeleteAt(cards []*pb.CardResponseDeleteAt) []*response.CardResponseDeleteAt {
	var responseCards []*response.CardResponseDeleteAt

	for _, role := range cards {
		responseCards = append(responseCards, s.mapCardResponseDeleteAt(role))
	}

	return responseCards
}

func (s *cardResponseMapper) mapDashboardCard(dash *pb.CardResponseDashboard) *response.DashboardCard {
	return &response.DashboardCard{
		TotalBalance:     &dash.TotalBalance,
		TotalWithdraw:    &dash.TotalWithdraw,
		TotalTopup:       &dash.TotalTopup,
		TotalTransfer:    &dash.TotalTransfer,
		TotalTransaction: &dash.TotalTransaction,
	}
}

func (s *cardResponseMapper) mapDashboardCardCardNumber(dash *pb.CardResponseDashboardCardNumber) *response.DashboardCardCardNumber {
	return &response.DashboardCardCardNumber{
		TotalBalance:          &dash.TotalBalance,
		TotalWithdraw:         &dash.TotalWithdraw,
		TotalTopup:            &dash.TotalTopup,
		TotalTransferSend:     &dash.TotalTransferSend,
		TotalTransferReceiver: &dash.TotalTransferReceiver,
		TotalTransaction:      &dash.TotalTransaction,
	}
}

func (s *cardResponseMapper) mapMonthlyBalance(cards *pb.CardResponseMonthlyBalance) *response.CardResponseMonthBalance {
	return &response.CardResponseMonthBalance{
		Month:        cards.Month,
		TotalBalance: cards.TotalBalance,
	}
}

func (s *cardResponseMapper) mapMonthlyBalances(cards []*pb.CardResponseMonthlyBalance) []*response.CardResponseMonthBalance {
	var responseCards []*response.CardResponseMonthBalance

	for _, role := range cards {
		responseCards = append(responseCards, s.mapMonthlyBalance(role))
	}

	return responseCards
}

func (s *cardResponseMapper) mapYearlyBalance(cards *pb.CardResponseYearlyBalance) *response.CardResponseYearlyBalance {
	return &response.CardResponseYearlyBalance{
		Year:         cards.Year,
		TotalBalance: cards.TotalBalance,
	}
}

func (s *cardResponseMapper) mapYearlyBalances(cards []*pb.CardResponseYearlyBalance) []*response.CardResponseYearlyBalance {
	var responseCards []*response.CardResponseYearlyBalance

	for _, role := range cards {
		responseCards = append(responseCards, s.mapYearlyBalance(role))
	}

	return responseCards
}

func (s *cardResponseMapper) mapMonthlyAmount(cards *pb.CardResponseMonthlyAmount) *response.CardResponseMonthAmount {
	return &response.CardResponseMonthAmount{
		Month:       cards.Month,
		TotalAmount: cards.TotalAmount,
	}
}

func (s *cardResponseMapper) mapMonthlyAmounts(cards []*pb.CardResponseMonthlyAmount) []*response.CardResponseMonthAmount {
	var responseCards []*response.CardResponseMonthAmount

	for _, role := range cards {
		responseCards = append(responseCards, s.mapMonthlyAmount(role))
	}

	return responseCards
}

func (s *cardResponseMapper) mapYearlyAmount(cards *pb.CardResponseYearlyAmount) *response.CardResponseYearAmount {
	return &response.CardResponseYearAmount{
		Year:        cards.Year,
		TotalAmount: cards.TotalAmount,
	}
}

func (s *cardResponseMapper) mapYearlyAmounts(cards []*pb.CardResponseYearlyAmount) []*response.CardResponseYearAmount {
	var responseCards []*response.CardResponseYearAmount

	for _, role := range cards {
		responseCards = append(responseCards, s.mapYearlyAmount(role))
	}

	return responseCards
}
