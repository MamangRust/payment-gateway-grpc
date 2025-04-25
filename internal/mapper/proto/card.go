package protomapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type cardProtoMapper struct{}

func NewCardProtoMapper() *cardProtoMapper {
	return &cardProtoMapper{}
}

func (s *cardProtoMapper) ToProtoResponseCard(status string, message string, card *response.CardResponse) *pb.ApiResponseCard {
	return &pb.ApiResponseCard{
		Status:  status,
		Message: message,
		Data:    s.mapCardResponse(card),
	}
}

func (s *cardProtoMapper) ToProtoResponsePaginationCard(pagination *pb.PaginationMeta, status string, message string, cards []*response.CardResponse) *pb.ApiResponsePaginationCard {
	return &pb.ApiResponsePaginationCard{
		Status:     status,
		Message:    message,
		Data:       s.mapCardResponses(cards),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (s *cardProtoMapper) ToProtoResponseCardDeleteAt(status string, message string) *pb.ApiResponseCardDelete {
	return &pb.ApiResponseCardDelete{
		Status:  status,
		Message: message,
	}
}

func (s *cardProtoMapper) ToProtoResponseCardAll(status string, message string) *pb.ApiResponseCardAll {
	return &pb.ApiResponseCardAll{
		Status:  status,
		Message: message,
	}
}

func (s *cardProtoMapper) ToProtoResponsePaginationCardDeletedAt(pagination *pb.PaginationMeta, status string, message string, cards []*response.CardResponseDeleteAt) *pb.ApiResponsePaginationCardDeleteAt {
	return &pb.ApiResponsePaginationCardDeleteAt{
		Status:     status,
		Message:    message,
		Data:       s.mapCardResponsesDeleteAt(cards),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (s *cardProtoMapper) ToProtoResponseDashboardCard(status string, message string, dash *response.DashboardCard) *pb.ApiResponseDashboardCard {
	return &pb.ApiResponseDashboardCard{
		Status:  status,
		Message: message,
		Data:    s.mapDashboardCard(dash),
	}
}

func (s *cardProtoMapper) ToProtoResponseDashboardCardCardNumber(status string, message string, dash *response.DashboardCardCardNumber) *pb.ApiResponseDashboardCardNumber {
	return &pb.ApiResponseDashboardCardNumber{
		Status:  status,
		Message: message,
		Data:    s.mapDashboardCardCardNumber(dash),
	}
}

func (s *cardProtoMapper) ToProtoResponseMonthlyBalances(status string, message string, cards []*response.CardResponseMonthBalance) *pb.ApiResponseMonthlyBalance {

	return &pb.ApiResponseMonthlyBalance{
		Status:  status,
		Message: message,
		Data:    s.mapMonthlyBalances(cards),
	}
}

func (s *cardProtoMapper) ToProtoResponseYearlyBalances(status string, message string, cards []*response.CardResponseYearlyBalance) *pb.ApiResponseYearlyBalance {

	return &pb.ApiResponseYearlyBalance{
		Status:  status,
		Message: message,
		Data:    s.mapYearlyBalances(cards),
	}
}

func (s *cardProtoMapper) ToProtoResponseMonthlyAmounts(status string, message string, cards []*response.CardResponseMonthAmount) *pb.ApiResponseMonthlyAmount {

	return &pb.ApiResponseMonthlyAmount{
		Status:  status,
		Message: message,
		Data:    s.mapMonthlyAmounts(cards),
	}
}

func (s *cardProtoMapper) ToProtoResponseYearlyAmounts(status string, message string, cards []*response.CardResponseYearAmount) *pb.ApiResponseYearlyAmount {
	return &pb.ApiResponseYearlyAmount{
		Status:  status,
		Message: message,
		Data:    s.mapYearlyAmounts(cards),
	}
}

func (s *cardProtoMapper) mapCardResponse(card *response.CardResponse) *pb.CardResponse {
	return &pb.CardResponse{
		Id:           int32(card.ID),
		UserId:       int32(card.UserID),
		CardNumber:   card.CardNumber,
		CardType:     card.CardType,
		ExpireDate:   card.ExpireDate,
		Cvv:          card.CVV,
		CardProvider: card.CardProvider,
		CreatedAt:    card.CreatedAt,
		UpdatedAt:    card.UpdatedAt,
	}
}

func (s *cardProtoMapper) mapCardResponses(roles []*response.CardResponse) []*pb.CardResponse {
	var responseRoles []*pb.CardResponse

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapCardResponse(role))
	}

	return responseRoles
}

func (s *cardProtoMapper) mapCardResponseDeleteAt(card *response.CardResponseDeleteAt) *pb.CardResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if card.DeletedAt != nil {
		deletedAt = wrapperspb.String(*card.DeletedAt)
	}

	return &pb.CardResponseDeleteAt{
		Id:           int32(card.ID),
		UserId:       int32(card.UserID),
		CardNumber:   card.CardNumber,
		CardType:     card.CardType,
		ExpireDate:   card.ExpireDate,
		Cvv:          card.CVV,
		CardProvider: card.CardProvider,
		CreatedAt:    card.CreatedAt,
		UpdatedAt:    card.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (s *cardProtoMapper) mapCardResponsesDeleteAt(roles []*response.CardResponseDeleteAt) []*pb.CardResponseDeleteAt {
	var responseRoles []*pb.CardResponseDeleteAt

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapCardResponseDeleteAt(role))
	}

	return responseRoles
}

func (s *cardProtoMapper) mapDashboardCard(dash *response.DashboardCard) *pb.CardResponseDashboard {
	return &pb.CardResponseDashboard{
		TotalBalance:     *dash.TotalBalance,
		TotalWithdraw:    *dash.TotalWithdraw,
		TotalTopup:       *dash.TotalTopup,
		TotalTransfer:    *dash.TotalTransfer,
		TotalTransaction: *dash.TotalTransaction,
	}
}

func (s *cardProtoMapper) mapDashboardCardCardNumber(dash *response.DashboardCardCardNumber) *pb.CardResponseDashboardCardNumber {
	return &pb.CardResponseDashboardCardNumber{
		TotalBalance:          *dash.TotalBalance,
		TotalWithdraw:         *dash.TotalWithdraw,
		TotalTopup:            *dash.TotalTopup,
		TotalTransferSend:     *dash.TotalTransferSend,
		TotalTransferReceiver: *dash.TotalTransferReceiver,
		TotalTransaction:      *dash.TotalTransaction,
	}
}

func (s *cardProtoMapper) mapMonthlyBalance(card *response.CardResponseMonthBalance) *pb.CardResponseMonthlyBalance {
	return &pb.CardResponseMonthlyBalance{
		Month:        card.Month,
		TotalBalance: card.TotalBalance,
	}
}

func (s *cardProtoMapper) mapMonthlyBalances(roles []*response.CardResponseMonthBalance) []*pb.CardResponseMonthlyBalance {
	var responseRoles []*pb.CardResponseMonthlyBalance

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapMonthlyBalance(role))
	}

	return responseRoles
}

func (s *cardProtoMapper) mapYearlyBalance(card *response.CardResponseYearlyBalance) *pb.CardResponseYearlyBalance {
	return &pb.CardResponseYearlyBalance{
		Year:         card.Year,
		TotalBalance: card.TotalBalance,
	}
}

func (s *cardProtoMapper) mapYearlyBalances(roles []*response.CardResponseYearlyBalance) []*pb.CardResponseYearlyBalance {
	var responseRoles []*pb.CardResponseYearlyBalance

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapYearlyBalance(role))
	}

	return responseRoles
}

func (s *cardProtoMapper) mapMonthlyAmount(card *response.CardResponseMonthAmount) *pb.CardResponseMonthlyAmount {
	return &pb.CardResponseMonthlyAmount{
		Month:       card.Month,
		TotalAmount: card.TotalAmount,
	}
}

func (s *cardProtoMapper) mapMonthlyAmounts(roles []*response.CardResponseMonthAmount) []*pb.CardResponseMonthlyAmount {
	var responseRoles []*pb.CardResponseMonthlyAmount

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapMonthlyAmount(role))
	}

	return responseRoles
}

func (s *cardProtoMapper) mapYearlyAmount(card *response.CardResponseYearAmount) *pb.CardResponseYearlyAmount {
	return &pb.CardResponseYearlyAmount{
		Year:        card.Year,
		TotalAmount: card.TotalAmount,
	}
}

func (s *cardProtoMapper) mapYearlyAmounts(roles []*response.CardResponseYearAmount) []*pb.CardResponseYearlyAmount {
	var responseRoles []*pb.CardResponseYearlyAmount

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapYearlyAmount(role))
	}

	return responseRoles
}
