package responseservice

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
)

type saldoResponseMapper struct {
}

func NewSaldoResponseMapper() *saldoResponseMapper {
	return &saldoResponseMapper{}
}

func (s *saldoResponseMapper) ToSaldoResponse(saldo *record.SaldoRecord) *response.SaldoResponse {
	return &response.SaldoResponse{
		ID:             saldo.ID,
		CardNumber:     saldo.CardNumber,
		TotalBalance:   saldo.TotalBalance,
		WithdrawAmount: saldo.WithdrawAmount,
		WithdrawTime:   saldo.WithdrawTime,
		CreatedAt:      saldo.CreatedAt,
		UpdatedAt:      saldo.UpdatedAt,
	}
}

func (s *saldoResponseMapper) ToSaldoResponses(saldos []*record.SaldoRecord) []*response.SaldoResponse {
	var responses []*response.SaldoResponse

	for _, response := range saldos {
		responses = append(responses, s.ToSaldoResponse(response))
	}

	return responses
}

func (s *saldoResponseMapper) ToSaldoResponseDeleteAt(saldo *record.SaldoRecord) *response.SaldoResponseDeleteAt {
	return &response.SaldoResponseDeleteAt{
		ID:             saldo.ID,
		CardNumber:     saldo.CardNumber,
		TotalBalance:   saldo.TotalBalance,
		WithdrawAmount: saldo.WithdrawAmount,
		WithdrawTime:   saldo.WithdrawTime,
		CreatedAt:      saldo.CreatedAt,
		UpdatedAt:      saldo.UpdatedAt,
		DeletedAt:      saldo.DeletedAt,
	}
}

func (s *saldoResponseMapper) ToSaldoResponsesDeleteAt(saldos []*record.SaldoRecord) []*response.SaldoResponseDeleteAt {
	var responses []*response.SaldoResponseDeleteAt

	for _, response := range saldos {
		responses = append(responses, s.ToSaldoResponseDeleteAt(response))
	}

	return responses
}

func (s *saldoResponseMapper) ToSaldoMonthTotalBalanceResponse(ss *record.SaldoMonthTotalBalance) *response.SaldoMonthTotalBalanceResponse {
	totalBalance := 0

	if ss.TotalBalance != 0 {
		totalBalance = ss.TotalBalance
	}

	return &response.SaldoMonthTotalBalanceResponse{
		Month:        ss.Month,
		Year:         ss.Year,
		TotalBalance: totalBalance,
	}
}

func (s *saldoResponseMapper) ToSaldoMonthTotalBalanceResponses(ss []*record.SaldoMonthTotalBalance) []*response.SaldoMonthTotalBalanceResponse {
	var saldoResponses []*response.SaldoMonthTotalBalanceResponse
	for _, saldo := range ss {
		saldoResponses = append(saldoResponses, s.ToSaldoMonthTotalBalanceResponse(saldo))
	}
	return saldoResponses
}

func (s *saldoResponseMapper) ToSaldoYearTotalBalanceResponse(ss *record.SaldoYearTotalBalance) *response.SaldoYearTotalBalanceResponse {
	return &response.SaldoYearTotalBalanceResponse{
		Year:         ss.Year,
		TotalBalance: ss.TotalBalance,
	}
}

func (s *saldoResponseMapper) ToSaldoYearTotalBalanceResponses(ss []*record.SaldoYearTotalBalance) []*response.SaldoYearTotalBalanceResponse {
	var saldoResponses []*response.SaldoYearTotalBalanceResponse
	for _, saldo := range ss {
		saldoResponses = append(saldoResponses, s.ToSaldoYearTotalBalanceResponse(saldo))
	}
	return saldoResponses
}

func (s *saldoResponseMapper) ToSaldoMonthBalanceResponse(ss *record.SaldoMonthSaldoBalance) *response.SaldoMonthBalanceResponse {
	return &response.SaldoMonthBalanceResponse{
		Month:        ss.Month,
		TotalBalance: ss.TotalBalance,
	}
}

func (s *saldoResponseMapper) ToSaldoMonthBalanceResponses(ss []*record.SaldoMonthSaldoBalance) []*response.SaldoMonthBalanceResponse {
	var saldoResponses []*response.SaldoMonthBalanceResponse
	for _, saldo := range ss {
		saldoResponses = append(saldoResponses, s.ToSaldoMonthBalanceResponse(saldo))
	}
	return saldoResponses
}

func (s *saldoResponseMapper) ToSaldoYearBalanceResponse(ss *record.SaldoYearSaldoBalance) *response.SaldoYearBalanceResponse {
	return &response.SaldoYearBalanceResponse{
		Year:         ss.Year,
		TotalBalance: ss.TotalBalance,
	}
}

func (s *saldoResponseMapper) ToSaldoYearBalanceResponses(ss []*record.SaldoYearSaldoBalance) []*response.SaldoYearBalanceResponse {
	var saldoResponses []*response.SaldoYearBalanceResponse
	for _, saldo := range ss {
		saldoResponses = append(saldoResponses, s.ToSaldoYearBalanceResponse(saldo))
	}
	return saldoResponses
}
