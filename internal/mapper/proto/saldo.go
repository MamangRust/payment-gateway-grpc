package protomapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type saldoProtoMapper struct {
}

func NewSaldoProtoMapper() *saldoProtoMapper {
	return &saldoProtoMapper{}
}

func (s *saldoProtoMapper) ToProtoResponseSaldo(status string, message string, pbResponse *response.SaldoResponse) *pb.ApiResponseSaldo {
	return &pb.ApiResponseSaldo{
		Status:  status,
		Message: message,
		Data:    s.mapResponseSaldo(pbResponse),
	}
}

func (s *saldoProtoMapper) ToProtoResponsesSaldo(status string, message string, pbResponse []*response.SaldoResponse) *pb.ApiResponsesSaldo {
	return &pb.ApiResponsesSaldo{
		Status:  status,
		Message: message,
		Data:    s.mapResponsesSaldo(pbResponse),
	}
}

func (s *saldoProtoMapper) ToProtoResponseSaldoDelete(status string, message string) *pb.ApiResponseSaldoDelete {
	return &pb.ApiResponseSaldoDelete{
		Status:  status,
		Message: message,
	}
}

func (s *saldoProtoMapper) ToProtoResponseSaldoAll(status string, message string) *pb.ApiResponseSaldoAll {
	return &pb.ApiResponseSaldoAll{
		Status:  status,
		Message: message,
	}
}

func (s *saldoProtoMapper) ToProtoResponseMonthTotalSaldo(status string, message string, pbResponse []*response.SaldoMonthTotalBalanceResponse) *pb.ApiResponseMonthTotalSaldo {
	return &pb.ApiResponseMonthTotalSaldo{
		Status:  status,
		Message: message,
		Data:    s.mapSaldoMonthTotalBalanceResponses(pbResponse),
	}
}

func (s *saldoProtoMapper) ToProtoResponseYearTotalSaldo(status string, message string, pbResponse []*response.SaldoYearTotalBalanceResponse) *pb.ApiResponseYearTotalSaldo {
	return &pb.ApiResponseYearTotalSaldo{
		Status:  status,
		Message: message,
		Data:    s.mapSaldoYearTotalBalanceResponses(pbResponse),
	}
}

func (s *saldoProtoMapper) ToProtoResponseMonthSaldoBalances(status string, message string, pbResponse []*response.SaldoMonthBalanceResponse) *pb.ApiResponseMonthSaldoBalances {
	return &pb.ApiResponseMonthSaldoBalances{
		Status:  status,
		Message: message,
		Data:    s.mapSaldoMonthBalanceResponses(pbResponse),
	}
}

func (s *saldoProtoMapper) ToProtoResponseYearSaldoBalances(status string, message string, pbResponse []*response.SaldoYearBalanceResponse) *pb.ApiResponseYearSaldoBalances {
	return &pb.ApiResponseYearSaldoBalances{
		Status:  status,
		Message: message,
		Data:    s.mapSaldoYearBalanceResponses(pbResponse),
	}
}

func (s *saldoProtoMapper) ToProtoResponsePaginationSaldo(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.SaldoResponse) *pb.ApiResponsePaginationSaldo {
	return &pb.ApiResponsePaginationSaldo{
		Status:     status,
		Message:    message,
		Data:       s.mapResponsesSaldo(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (s *saldoProtoMapper) ToProtoResponsePaginationSaldoDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.SaldoResponseDeleteAt) *pb.ApiResponsePaginationSaldoDeleteAt {
	return &pb.ApiResponsePaginationSaldoDeleteAt{
		Status:     status,
		Message:    message,
		Data:       s.mapResponsesSaldoDeleteAt(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (s *saldoProtoMapper) mapResponseSaldo(saldo *response.SaldoResponse) *pb.SaldoResponse {
	return &pb.SaldoResponse{
		SaldoId:        int32(saldo.ID),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   int32(saldo.TotalBalance),
		WithdrawTime:   saldo.WithdrawTime,
		WithdrawAmount: int32(saldo.WithdrawAmount),
		CreatedAt:      saldo.CreatedAt,
		UpdatedAt:      saldo.UpdatedAt,
	}
}

func (s *saldoProtoMapper) mapResponsesSaldo(saldos []*response.SaldoResponse) []*pb.SaldoResponse {
	var responseSaldos []*pb.SaldoResponse

	for _, saldo := range saldos {
		responseSaldos = append(responseSaldos, s.mapResponseSaldo(saldo))
	}

	return responseSaldos
}

func (s *saldoProtoMapper) mapResponseSaldoDeleteAt(saldo *response.SaldoResponseDeleteAt) *pb.SaldoResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if saldo.DeletedAt != nil {
		deletedAt = wrapperspb.String(*saldo.DeletedAt)
	}

	return &pb.SaldoResponseDeleteAt{
		SaldoId:        int32(saldo.ID),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   int32(saldo.TotalBalance),
		WithdrawTime:   saldo.WithdrawTime,
		WithdrawAmount: int32(saldo.WithdrawAmount),
		CreatedAt:      saldo.CreatedAt,
		UpdatedAt:      saldo.UpdatedAt,
		DeletedAt:      deletedAt,
	}
}

func (s *saldoProtoMapper) mapResponsesSaldoDeleteAt(saldos []*response.SaldoResponseDeleteAt) []*pb.SaldoResponseDeleteAt {
	var responseSaldos []*pb.SaldoResponseDeleteAt

	for _, saldo := range saldos {
		responseSaldos = append(responseSaldos, s.mapResponseSaldoDeleteAt(saldo))
	}

	return responseSaldos
}

func (s *saldoProtoMapper) mapSaldoMonthTotalBalanceResponse(ss *response.SaldoMonthTotalBalanceResponse) *pb.SaldoMonthTotalBalanceResponse {
	totalBalance := 0

	if ss.TotalBalance != 0 {
		totalBalance = ss.TotalBalance
	}

	return &pb.SaldoMonthTotalBalanceResponse{
		Month:        ss.Month,
		Year:         ss.Year,
		TotalBalance: int32(totalBalance),
	}
}

func (s *saldoProtoMapper) mapSaldoMonthTotalBalanceResponses(ss []*response.SaldoMonthTotalBalanceResponse) []*pb.SaldoMonthTotalBalanceResponse {
	var saldoProtos []*pb.SaldoMonthTotalBalanceResponse
	for _, saldo := range ss {
		saldoProtos = append(saldoProtos, s.mapSaldoMonthTotalBalanceResponse(saldo))
	}
	return saldoProtos
}

func (s *saldoProtoMapper) mapSaldoYearTotalBalanceResponse(ss *response.SaldoYearTotalBalanceResponse) *pb.SaldoYearTotalBalanceResponse {
	totalBalance := 0

	if ss.TotalBalance != 0 {
		totalBalance = ss.TotalBalance
	}

	return &pb.SaldoYearTotalBalanceResponse{
		Year:         ss.Year,
		TotalBalance: int32(totalBalance),
	}
}

func (s *saldoProtoMapper) mapSaldoYearTotalBalanceResponses(ss []*response.SaldoYearTotalBalanceResponse) []*pb.SaldoYearTotalBalanceResponse {
	var saldoProtos []*pb.SaldoYearTotalBalanceResponse
	for _, saldo := range ss {
		saldoProtos = append(saldoProtos, s.mapSaldoYearTotalBalanceResponse(saldo))
	}
	return saldoProtos
}

func (s *saldoProtoMapper) mapSaldoMonthBalanceResponse(ss *response.SaldoMonthBalanceResponse) *pb.SaldoMonthBalanceResponse {
	return &pb.SaldoMonthBalanceResponse{
		Month:        ss.Month,
		TotalBalance: int32(ss.TotalBalance),
	}
}

func (s *saldoProtoMapper) mapSaldoMonthBalanceResponses(ss []*response.SaldoMonthBalanceResponse) []*pb.SaldoMonthBalanceResponse {
	var saldoProtos []*pb.SaldoMonthBalanceResponse
	for _, saldo := range ss {
		saldoProtos = append(saldoProtos, s.mapSaldoMonthBalanceResponse(saldo))
	}
	return saldoProtos
}

func (s *saldoProtoMapper) mapSaldoYearBalanceResponse(ss *response.SaldoYearBalanceResponse) *pb.SaldoYearBalanceResponse {
	return &pb.SaldoYearBalanceResponse{
		Year:         ss.Year,
		TotalBalance: int32(ss.TotalBalance),
	}
}

func (s *saldoProtoMapper) mapSaldoYearBalanceResponses(ss []*response.SaldoYearBalanceResponse) []*pb.SaldoYearBalanceResponse {
	var saldoProtos []*pb.SaldoYearBalanceResponse
	for _, saldo := range ss {
		saldoProtos = append(saldoProtos, s.mapSaldoYearBalanceResponse(saldo))
	}
	return saldoProtos
}
