package apimapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"
)

type merchantResponse struct{}

func NewMerchantResponseMapper() *merchantResponse {
	return &merchantResponse{}
}

func (m *merchantResponse) ToApiResponseMerchant(merchants *pb.ApiResponseMerchant) *response.ApiResponseMerchant {
	return &response.ApiResponseMerchant{
		Status:  merchants.Status,
		Message: merchants.Message,
		Data:    *m.mapMerchantResponse(merchants.Data),
	}
}

func (m *merchantResponse) ToApiResponseMerchants(merchants *pb.ApiResponsesMerchant) *response.ApiResponsesMerchant {
	return &response.ApiResponsesMerchant{
		Status:  merchants.Status,
		Message: merchants.Message,
		Data:    m.mapMerchantResponses(merchants.Data),
	}
}

func (m *merchantResponse) ToApiResponsesMerchant(merchants *pb.ApiResponsePaginationMerchant) *response.ApiResponsePaginationMerchant {
	return &response.ApiResponsePaginationMerchant{
		Status:     merchants.Status,
		Message:    merchants.Message,
		Data:       m.mapMerchantResponses(merchants.Data),
		Pagination: mapPaginationMeta(merchants.Pagination),
	}
}

func (m *merchantResponse) ToApiResponsesMerchantDeleteAt(merchants *pb.ApiResponsePaginationMerchantDeleteAt) *response.ApiResponsePaginationMerchantDeleteAt {
	return &response.ApiResponsePaginationMerchantDeleteAt{
		Status:     merchants.Status,
		Message:    merchants.Message,
		Data:       m.mapMerchantResponsesDeleteAt(merchants.Data),
		Pagination: mapPaginationMeta(merchants.Pagination),
	}
}

func (m *merchantResponse) ToApiResponseMerchantsTransactionResponse(merchants *pb.ApiResponsePaginationMerchantTransaction) *response.ApiResponsePaginationMerchantTransaction {

	return &response.ApiResponsePaginationMerchantTransaction{
		Status:     merchants.Status,
		Message:    merchants.Message,
		Data:       m.mapMerchantTransactionResponses(merchants.Data),
		Pagination: mapPaginationMeta(merchants.Pagination),
	}
}

func (m *merchantResponse) ToApiResponseMonthlyPaymentMethods(ms *pb.ApiResponseMerchantMonthlyPaymentMethod) *response.ApiResponseMerchantMonthlyPaymentMethod {
	return &response.ApiResponseMerchantMonthlyPaymentMethod{
		Status:  ms.Status,
		Message: ms.Message,
		Data:    m.mapResponsesMonthlyPaymentMethod(ms.Data),
	}
}

func (m *merchantResponse) ToApiResponseYearlyPaymentMethods(ms *pb.ApiResponseMerchantYearlyPaymentMethod) *response.ApiResponseMerchantYearlyPaymentMethod {
	return &response.ApiResponseMerchantYearlyPaymentMethod{
		Status:  ms.Status,
		Message: ms.Message,
		Data:    m.mapResponsesYearlyPaymentMethod(ms.Data),
	}
}

func (m *merchantResponse) ToApiResponseMonthlyAmounts(ms *pb.ApiResponseMerchantMonthlyAmount) *response.ApiResponseMerchantMonthlyAmount {
	return &response.ApiResponseMerchantMonthlyAmount{
		Status:  ms.Status,
		Message: ms.Message,
		Data:    m.mapResponsesMonthlyAmount(ms.Data),
	}
}

func (m *merchantResponse) ToApiResponseYearlyAmounts(ms *pb.ApiResponseMerchantYearlyAmount) *response.ApiResponseMerchantYearlyAmount {
	return &response.ApiResponseMerchantYearlyAmount{
		Status:  ms.Status,
		Message: ms.Message,
		Data:    m.mapResponsesYearlyAmount(ms.Data),
	}
}

func (m *merchantResponse) ToApiResponseMonthlyTotalAmounts(ms *pb.ApiResponseMerchantMonthlyTotalAmount) *response.ApiResponseMerchantMonthlyTotalAmount {
	return &response.ApiResponseMerchantMonthlyTotalAmount{
		Status:  ms.Status,
		Message: ms.Message,
		Data:    m.mapResponsesMonthlyTotalAmount(ms.Data),
	}
}

func (m *merchantResponse) ToApiResponseYearlyTotalAmounts(ms *pb.ApiResponseMerchantYearlyTotalAmount) *response.ApiResponseMerchantYearlyTotalAmount {
	return &response.ApiResponseMerchantYearlyTotalAmount{
		Status:  ms.Status,
		Message: ms.Message,
		Data:    m.mapResponsesYearlyTotalAmount(ms.Data),
	}
}

func (s *merchantResponse) ToApiResponseMerchantDeleteAt(card *pb.ApiResponseMerchantDelete) *response.ApiResponseMerchantDelete {
	return &response.ApiResponseMerchantDelete{
		Status:  card.Status,
		Message: card.Message,
	}
}

func (s *merchantResponse) ToApiResponseMerchantAll(card *pb.ApiResponseMerchantAll) *response.ApiResponseMerchantAll {
	return &response.ApiResponseMerchantAll{
		Status:  card.Status,
		Message: card.Message,
	}
}

func (m *merchantResponse) mapMerchantResponse(merchant *pb.MerchantResponse) *response.MerchantResponse {
	return &response.MerchantResponse{
		ID:        int(merchant.Id),
		Name:      merchant.Name,
		Status:    merchant.Status,
		ApiKey:    merchant.ApiKey,
		UserID:    int(merchant.UserId),
		CreatedAt: merchant.CreatedAt,
		UpdatedAt: merchant.UpdatedAt,
	}
}

func (m *merchantResponse) mapMerchantResponses(r []*pb.MerchantResponse) []*response.MerchantResponse {
	var responseMerchants []*response.MerchantResponse
	for _, merchant := range r {
		responseMerchants = append(responseMerchants, m.mapMerchantResponse(merchant))
	}

	return responseMerchants
}

func (m *merchantResponse) mapMerchantResponseDeleteAt(merchant *pb.MerchantResponseDeleteAt) *response.MerchantResponseDeleteAt {
	var deletedAt string
	if merchant.DeletedAt != nil {
		deletedAt = merchant.DeletedAt.Value
	}

	return &response.MerchantResponseDeleteAt{
		ID:        int(merchant.Id),
		Name:      merchant.Name,
		Status:    merchant.Status,
		UserID:    int(merchant.UserId),
		ApiKey:    merchant.ApiKey,
		CreatedAt: merchant.CreatedAt,
		UpdatedAt: merchant.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (m *merchantResponse) mapMerchantResponsesDeleteAt(r []*pb.MerchantResponseDeleteAt) []*response.MerchantResponseDeleteAt {
	var responseMerchants []*response.MerchantResponseDeleteAt
	for _, merchant := range r {
		responseMerchants = append(responseMerchants, m.mapMerchantResponseDeleteAt(merchant))
	}

	return responseMerchants
}

func (m *merchantResponse) mapMerchantTransactionResponse(merchant *pb.MerchantTransactionResponse) *response.MerchantTransactionResponse {

	return &response.MerchantTransactionResponse{
		ID:              int(merchant.Id),
		CardNumber:      merchant.CardNumber,
		Amount:          merchant.Amount,
		PaymentMethod:   merchant.PaymentMethod,
		MerchantID:      merchant.MerchantId,
		MerchantName:    merchant.MerchantName,
		TransactionTime: merchant.TransactionTime,
		CreatedAt:       merchant.CreatedAt,
		UpdatedAt:       merchant.UpdatedAt,
	}
}

func (m *merchantResponse) mapMerchantTransactionResponses(r []*pb.MerchantTransactionResponse) []*response.MerchantTransactionResponse {
	var responseMerchants []*response.MerchantTransactionResponse
	for _, merchant := range r {
		responseMerchants = append(responseMerchants, m.mapMerchantTransactionResponse(merchant))
	}

	return responseMerchants
}

func (m *merchantResponse) mapResponseMonthlyPaymentMethod(ms *pb.MerchantResponseMonthlyPaymentMethod) *response.MerchantResponseMonthlyPaymentMethod {
	return &response.MerchantResponseMonthlyPaymentMethod{
		Month:         ms.Month,
		PaymentMethod: ms.PaymentMethod,
		TotalAmount:   int(ms.TotalAmount),
	}
}

func (m *merchantResponse) mapResponsesMonthlyPaymentMethod(r []*pb.MerchantResponseMonthlyPaymentMethod) []*response.MerchantResponseMonthlyPaymentMethod {
	var responseMerchants []*response.MerchantResponseMonthlyPaymentMethod
	for _, merchant := range r {
		responseMerchants = append(responseMerchants, m.mapResponseMonthlyPaymentMethod(merchant))
	}

	return responseMerchants
}

func (m *merchantResponse) mapResponseYearlyPaymentMethod(ms *pb.MerchantResponseYearlyPaymentMethod) *response.MerchantResponseYearlyPaymentMethod {
	return &response.MerchantResponseYearlyPaymentMethod{
		Year:          ms.Year,
		PaymentMethod: ms.PaymentMethod,
		TotalAmount:   int(ms.TotalAmount),
	}
}

func (m *merchantResponse) mapResponsesYearlyPaymentMethod(r []*pb.MerchantResponseYearlyPaymentMethod) []*response.MerchantResponseYearlyPaymentMethod {
	var responseMerchants []*response.MerchantResponseYearlyPaymentMethod
	for _, merchant := range r {
		responseMerchants = append(responseMerchants, m.mapResponseYearlyPaymentMethod(merchant))
	}

	return responseMerchants
}

func (m *merchantResponse) mapResponseMonthlyAmount(ms *pb.MerchantResponseMonthlyAmount) *response.MerchantResponseMonthlyAmount {
	return &response.MerchantResponseMonthlyAmount{
		Month:       ms.Month,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantResponse) mapResponsesMonthlyAmount(r []*pb.MerchantResponseMonthlyAmount) []*response.MerchantResponseMonthlyAmount {
	var responseMerchants []*response.MerchantResponseMonthlyAmount
	for _, merchant := range r {
		responseMerchants = append(responseMerchants, m.mapResponseMonthlyAmount(merchant))
	}

	return responseMerchants
}

func (m *merchantResponse) mapResponseYearlyAmount(ms *pb.MerchantResponseYearlyAmount) *response.MerchantResponseYearlyAmount {
	return &response.MerchantResponseYearlyAmount{
		Year:        ms.Year,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantResponse) mapResponsesYearlyAmount(r []*pb.MerchantResponseYearlyAmount) []*response.MerchantResponseYearlyAmount {
	var responseMerchants []*response.MerchantResponseYearlyAmount
	for _, merchant := range r {
		responseMerchants = append(responseMerchants, m.mapResponseYearlyAmount(merchant))
	}

	return responseMerchants
}

func (m *merchantResponse) mapResponseMonthlyTotalAmount(ms *pb.MerchantResponseMonthlyTotalAmount) *response.MerchantResponseMonthlyTotalAmount {
	return &response.MerchantResponseMonthlyTotalAmount{
		Month:       ms.Month,
		Year:        ms.Year,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantResponse) mapResponsesMonthlyTotalAmount(r []*pb.MerchantResponseMonthlyTotalAmount) []*response.MerchantResponseMonthlyTotalAmount {
	var responseMerchants []*response.MerchantResponseMonthlyTotalAmount
	for _, merchant := range r {
		responseMerchants = append(responseMerchants, m.mapResponseMonthlyTotalAmount(merchant))
	}

	return responseMerchants
}

func (m *merchantResponse) mapResponseYearlyTotalAmount(ms *pb.MerchantResponseYearlyTotalAmount) *response.MerchantResponseYearlyTotalAmount {
	return &response.MerchantResponseYearlyTotalAmount{
		Year:        ms.Year,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantResponse) mapResponsesYearlyTotalAmount(r []*pb.MerchantResponseYearlyTotalAmount) []*response.MerchantResponseYearlyTotalAmount {
	var responseMerchants []*response.MerchantResponseYearlyTotalAmount
	for _, merchant := range r {
		responseMerchants = append(responseMerchants, m.mapResponseYearlyTotalAmount(merchant))
	}

	return responseMerchants
}
