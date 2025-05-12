package apimapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"
)

type saldoResponse struct {
}

func NewSaldoResponseMapper() *saldoResponse {
	return &saldoResponse{}
}

func (s *saldoResponse) ToApiResponseSaldo(pbResponse *pb.ApiResponseSaldo) *response.ApiResponseSaldo {
	return &response.ApiResponseSaldo{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    *s.mapResponseSaldo(pbResponse.Data),
	}
}

func (s *saldoResponse) ToApiResponsesSaldo(pbResponse *pb.ApiResponsesSaldo) *response.ApiResponsesSaldo {
	return &response.ApiResponsesSaldo{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapResponsesSaldo(pbResponse.Data),
	}
}

func (s *saldoResponse) ToApiResponseSaldoDelete(pbResponse *pb.ApiResponseSaldoDelete) *response.ApiResponseSaldoDelete {
	return &response.ApiResponseSaldoDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *saldoResponse) ToApiResponseSaldoAll(pbResponse *pb.ApiResponseSaldoAll) *response.ApiResponseSaldoAll {
	return &response.ApiResponseSaldoAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *saldoResponse) ToApiResponseMonthTotalSaldo(pbResponse *pb.ApiResponseMonthTotalSaldo) *response.ApiResponseMonthTotalSaldo {
	return &response.ApiResponseMonthTotalSaldo{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapSaldoMonthTotalBalanceResponses(pbResponse.Data),
	}
}

func (s *saldoResponse) ToApiResponseYearTotalSaldo(pbResponse *pb.ApiResponseYearTotalSaldo) *response.ApiResponseYearTotalSaldo {
	return &response.ApiResponseYearTotalSaldo{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapSaldoYearTotalBalanceResponses(pbResponse.Data),
	}
}

func (s *saldoResponse) ToApiResponseMonthSaldoBalances(pbResponse *pb.ApiResponseMonthSaldoBalances) *response.ApiResponseMonthSaldoBalances {
	return &response.ApiResponseMonthSaldoBalances{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapSaldoMonthBalanceResponses(pbResponse.Data),
	}
}

func (s *saldoResponse) ToApiResponseYearSaldoBalances(pbResponse *pb.ApiResponseYearSaldoBalances) *response.ApiResponseYearSaldoBalances {
	return &response.ApiResponseYearSaldoBalances{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapSaldoYearBalanceResponses(pbResponse.Data),
	}
}

func (s *saldoResponse) ToApiResponsePaginationSaldo(pbResponse *pb.ApiResponsePaginationSaldo) *response.ApiResponsePaginationSaldo {
	return &response.ApiResponsePaginationSaldo{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.mapResponsesSaldo(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.Pagination),
	}
}

func (s *saldoResponse) ToApiResponsePaginationSaldoDeleteAt(pbResponse *pb.ApiResponsePaginationSaldoDeleteAt) *response.ApiResponsePaginationSaldoDeleteAt {
	return &response.ApiResponsePaginationSaldoDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.mapResponsesSaldoDeleteAt(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.Pagination),
	}
}

func (s *saldoResponse) mapResponseSaldo(saldo *pb.SaldoResponse) *response.SaldoResponse {
	return &response.SaldoResponse{
		ID:             int(saldo.SaldoId),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   int(saldo.TotalBalance),
		WithdrawTime:   saldo.WithdrawTime,
		WithdrawAmount: int(saldo.WithdrawAmount),
		CreatedAt:      saldo.CreatedAt,
		UpdatedAt:      saldo.UpdatedAt,
	}
}

func (s *saldoResponse) mapResponsesSaldo(saldos []*pb.SaldoResponse) []*response.SaldoResponse {
	var responseSaldos []*response.SaldoResponse

	for _, saldo := range saldos {
		responseSaldos = append(responseSaldos, s.mapResponseSaldo(saldo))
	}

	return responseSaldos
}

func (s *saldoResponse) mapResponseSaldoDeleteAt(saldo *pb.SaldoResponseDeleteAt) *response.SaldoResponseDeleteAt {
	var deletedAt string
	if saldo.DeletedAt != nil {
		deletedAt = saldo.DeletedAt.Value
	}

	return &response.SaldoResponseDeleteAt{
		ID:             int(saldo.SaldoId),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   int(saldo.TotalBalance),
		WithdrawTime:   saldo.WithdrawTime,
		WithdrawAmount: int(saldo.WithdrawAmount),
		CreatedAt:      saldo.CreatedAt,
		UpdatedAt:      saldo.UpdatedAt,
		DeletedAt:      &deletedAt,
	}
}

func (s *saldoResponse) mapResponsesSaldoDeleteAt(saldos []*pb.SaldoResponseDeleteAt) []*response.SaldoResponseDeleteAt {
	var responseSaldos []*response.SaldoResponseDeleteAt

	for _, saldo := range saldos {
		responseSaldos = append(responseSaldos, s.mapResponseSaldoDeleteAt(saldo))
	}

	return responseSaldos
}

func (s *saldoResponse) mapSaldoMonthTotalBalanceResponse(ss *pb.SaldoMonthTotalBalanceResponse) *response.SaldoMonthTotalBalanceResponse {
	totalBalance := 0

	if ss.TotalBalance != 0 {
		totalBalance = int(ss.TotalBalance)
	}

	return &response.SaldoMonthTotalBalanceResponse{
		Month:        ss.Month,
		Year:         ss.Year,
		TotalBalance: totalBalance,
	}
}

func (s *saldoResponse) mapSaldoMonthTotalBalanceResponses(ss []*pb.SaldoMonthTotalBalanceResponse) []*response.SaldoMonthTotalBalanceResponse {
	var saldoProtos []*response.SaldoMonthTotalBalanceResponse
	for _, saldo := range ss {
		saldoProtos = append(saldoProtos, s.mapSaldoMonthTotalBalanceResponse(saldo))
	}
	return saldoProtos
}

func (s *saldoResponse) mapSaldoYearTotalBalanceResponse(ss *pb.SaldoYearTotalBalanceResponse) *response.SaldoYearTotalBalanceResponse {
	totalBalance := 0

	if ss.TotalBalance != 0 {
		totalBalance = int(ss.TotalBalance)
	}

	return &response.SaldoYearTotalBalanceResponse{
		Year:         ss.Year,
		TotalBalance: totalBalance,
	}
}

func (s *saldoResponse) mapSaldoYearTotalBalanceResponses(ss []*pb.SaldoYearTotalBalanceResponse) []*response.SaldoYearTotalBalanceResponse {
	var saldoProtos []*response.SaldoYearTotalBalanceResponse
	for _, saldo := range ss {
		saldoProtos = append(saldoProtos, s.mapSaldoYearTotalBalanceResponse(saldo))
	}
	return saldoProtos
}

func (s *saldoResponse) mapSaldoMonthBalanceResponse(ss *pb.SaldoMonthBalanceResponse) *response.SaldoMonthBalanceResponse {
	return &response.SaldoMonthBalanceResponse{
		Month:        ss.Month,
		TotalBalance: int(ss.TotalBalance),
	}
}

func (s *saldoResponse) mapSaldoMonthBalanceResponses(ss []*pb.SaldoMonthBalanceResponse) []*response.SaldoMonthBalanceResponse {
	var saldoProtos []*response.SaldoMonthBalanceResponse
	for _, saldo := range ss {
		saldoProtos = append(saldoProtos, s.mapSaldoMonthBalanceResponse(saldo))
	}
	return saldoProtos
}

func (s *saldoResponse) mapSaldoYearBalanceResponse(ss *pb.SaldoYearBalanceResponse) *response.SaldoYearBalanceResponse {
	return &response.SaldoYearBalanceResponse{
		Year:         ss.Year,
		TotalBalance: int(ss.TotalBalance),
	}
}

func (s *saldoResponse) mapSaldoYearBalanceResponses(ss []*pb.SaldoYearBalanceResponse) []*response.SaldoYearBalanceResponse {
	var saldoProtos []*response.SaldoYearBalanceResponse
	for _, saldo := range ss {
		saldoProtos = append(saldoProtos, s.mapSaldoYearBalanceResponse(saldo))
	}
	return saldoProtos
}
