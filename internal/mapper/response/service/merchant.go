package responseservice

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
)

type merchantResponseMapper struct{}

func NewMerchantResponseMapper() *merchantResponseMapper {
	return &merchantResponseMapper{}
}

func (s *merchantResponseMapper) ToMerchantResponse(merchant *record.MerchantRecord) *response.MerchantResponse {
	return &response.MerchantResponse{
		ID:        merchant.ID,
		Name:      merchant.Name,
		UserID:    merchant.UserID,
		Status:    merchant.Status,
		ApiKey:    merchant.ApiKey,
		CreatedAt: merchant.CreatedAt,
		UpdatedAt: merchant.UpdatedAt,
	}
}

func (s *merchantResponseMapper) ToMerchantsResponse(merchants []*record.MerchantRecord) []*response.MerchantResponse {
	var response []*response.MerchantResponse
	for _, merchant := range merchants {
		response = append(response, s.ToMerchantResponse(merchant))
	}
	return response
}

func (s *merchantResponseMapper) ToMerchantResponseDeleteAt(merchant *record.MerchantRecord) *response.MerchantResponseDeleteAt {
	return &response.MerchantResponseDeleteAt{
		ID:        merchant.ID,
		Name:      merchant.Name,
		UserID:    merchant.UserID,
		Status:    merchant.Status,
		ApiKey:    merchant.ApiKey,
		CreatedAt: merchant.CreatedAt,
		UpdatedAt: merchant.UpdatedAt,
		DeletedAt: merchant.DeletedAt,
	}
}

func (s *merchantResponseMapper) ToMerchantsResponseDeleteAt(merchants []*record.MerchantRecord) []*response.MerchantResponseDeleteAt {
	var response []*response.MerchantResponseDeleteAt
	for _, merchant := range merchants {
		response = append(response, s.ToMerchantResponseDeleteAt(merchant))
	}
	return response
}

func (m *merchantResponseMapper) ToMerchantTransactionResponse(merchant *record.MerchantTransactionsRecord) *response.MerchantTransactionResponse {

	return &response.MerchantTransactionResponse{
		ID:              int(merchant.TransactionID),
		CardNumber:      merchant.CardNumber,
		Amount:          merchant.Amount,
		PaymentMethod:   merchant.PaymentMethod,
		MerchantID:      merchant.MerchantID,
		MerchantName:    merchant.MerchantName,
		TransactionTime: merchant.TransactionTime.Format("2006-01-02"),
		CreatedAt:       merchant.CreatedAt,
		UpdatedAt:       merchant.UpdatedAt,
	}
}

func (m *merchantResponseMapper) ToMerchantsTransactionResponse(merchants []*record.MerchantTransactionsRecord) []*response.MerchantTransactionResponse {
	var records []*response.MerchantTransactionResponse
	for _, merchant := range merchants {
		records = append(records, m.ToMerchantTransactionResponse(merchant))
	}
	return records
}

func (s *merchantResponseMapper) ToMerchantMonthlyPaymentMethod(ms *record.MerchantMonthlyPaymentMethod) *response.MerchantResponseMonthlyPaymentMethod {
	return &response.MerchantResponseMonthlyPaymentMethod{
		Month:         ms.Month,
		PaymentMethod: ms.PaymentMethod,
		TotalAmount:   ms.TotalAmount,
	}
}

func (s *merchantResponseMapper) ToMerchantMonthlyPaymentMethods(ms []*record.MerchantMonthlyPaymentMethod) []*response.MerchantResponseMonthlyPaymentMethod {
	var response []*response.MerchantResponseMonthlyPaymentMethod
	for _, merchant := range ms {
		response = append(response, s.ToMerchantMonthlyPaymentMethod(merchant))
	}
	return response
}

func (s *merchantResponseMapper) ToMerchantYearlyPaymentMethod(ms *record.MerchantYearlyPaymentMethod) *response.MerchantResponseYearlyPaymentMethod {
	return &response.MerchantResponseYearlyPaymentMethod{
		Year:          ms.Year,
		PaymentMethod: ms.PaymentMethod,
		TotalAmount:   ms.TotalAmount,
	}
}

func (s *merchantResponseMapper) ToMerchantYearlyPaymentMethods(ms []*record.MerchantYearlyPaymentMethod) []*response.MerchantResponseYearlyPaymentMethod {
	var response []*response.MerchantResponseYearlyPaymentMethod
	for _, merchant := range ms {
		response = append(response, s.ToMerchantYearlyPaymentMethod(merchant))
	}
	return response
}

func (s *merchantResponseMapper) ToMerchantMonthlyAmount(ms *record.MerchantMonthlyAmount) *response.MerchantResponseMonthlyAmount {
	return &response.MerchantResponseMonthlyAmount{
		Month:       ms.Month,
		TotalAmount: ms.TotalAmount,
	}
}

func (s *merchantResponseMapper) ToMerchantMonthlyAmounts(ms []*record.MerchantMonthlyAmount) []*response.MerchantResponseMonthlyAmount {
	var response []*response.MerchantResponseMonthlyAmount
	for _, merchant := range ms {
		response = append(response, s.ToMerchantMonthlyAmount(merchant))
	}
	return response
}

func (s *merchantResponseMapper) ToMerchantYearlyAmount(ms *record.MerchantYearlyAmount) *response.MerchantResponseYearlyAmount {
	return &response.MerchantResponseYearlyAmount{
		Year:        ms.Year,
		TotalAmount: ms.TotalAmount,
	}
}

func (s *merchantResponseMapper) ToMerchantYearlyAmounts(ms []*record.MerchantYearlyAmount) []*response.MerchantResponseYearlyAmount {
	var response []*response.MerchantResponseYearlyAmount
	for _, merchant := range ms {
		response = append(response, s.ToMerchantYearlyAmount(merchant))
	}
	return response
}

func (s *merchantResponseMapper) ToMerchantMonthlyTotalAmount(ms *record.MerchantMonthlyTotalAmount) *response.MerchantResponseMonthlyTotalAmount {
	return &response.MerchantResponseMonthlyTotalAmount{
		Month:       ms.Month,
		Year:        ms.Year,
		TotalAmount: ms.TotalAmount,
	}
}

func (s *merchantResponseMapper) ToMerchantMonthlyTotalAmounts(ms []*record.MerchantMonthlyTotalAmount) []*response.MerchantResponseMonthlyTotalAmount {
	var response []*response.MerchantResponseMonthlyTotalAmount
	for _, merchant := range ms {
		response = append(response, s.ToMerchantMonthlyTotalAmount(merchant))
	}
	return response
}

func (s *merchantResponseMapper) ToMerchantYearlyTotalAmount(ms *record.MerchantYearlyTotalAmount) *response.MerchantResponseYearlyTotalAmount {
	return &response.MerchantResponseYearlyTotalAmount{
		Year:        ms.Year,
		TotalAmount: ms.TotalAmount,
	}
}

func (s *merchantResponseMapper) ToMerchantYearlyTotalAmounts(ms []*record.MerchantYearlyTotalAmount) []*response.MerchantResponseYearlyTotalAmount {
	var response []*response.MerchantResponseYearlyTotalAmount
	for _, merchant := range ms {
		response = append(response, s.ToMerchantYearlyTotalAmount(merchant))
	}
	return response
}
