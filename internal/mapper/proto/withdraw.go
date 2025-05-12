package protomapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type withdrawProtoMapper struct {
}

func NewWithdrawProtoMapper() *withdrawProtoMapper {
	return &withdrawProtoMapper{}
}

func (m *withdrawProtoMapper) ToProtoResponseWithdraw(status string, message string, withdraw *response.WithdrawResponse) *pb.ApiResponseWithdraw {
	return &pb.ApiResponseWithdraw{
		Status:  status,
		Message: message,
		Data:    m.mapResponseWithdrawal(withdraw),
	}
}

func (m *withdrawProtoMapper) ToProtoResponsesWithdraw(status string, message string, pbResponse []*response.WithdrawResponse) *pb.ApiResponsesWithdraw {
	return &pb.ApiResponsesWithdraw{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesWithdrawal(pbResponse),
	}
}

func (m *withdrawProtoMapper) ToProtoResponseWithdrawDelete(status string, message string) *pb.ApiResponseWithdrawDelete {
	return &pb.ApiResponseWithdrawDelete{
		Status:  status,
		Message: message,
	}
}

func (m *withdrawProtoMapper) ToProtoResponseWithdrawAll(status string, message string) *pb.ApiResponseWithdrawAll {
	return &pb.ApiResponseWithdrawAll{
		Status:  status,
		Message: message,
	}
}

func (m *withdrawProtoMapper) ToProtoResponsePaginationWithdraw(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.WithdrawResponse) *pb.ApiResponsePaginationWithdraw {
	return &pb.ApiResponsePaginationWithdraw{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesWithdrawal(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *withdrawProtoMapper) ToProtoResponsePaginationWithdrawDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.WithdrawResponseDeleteAt) *pb.ApiResponsePaginationWithdrawDeleteAt {
	return &pb.ApiResponsePaginationWithdrawDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesWithdrawalDeleteAt(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *withdrawProtoMapper) ToProtoResponseWithdrawMonthStatusSuccess(status string, message string, pbResponse []*response.WithdrawResponseMonthStatusSuccess) *pb.ApiResponseWithdrawMonthStatusSuccess {
	return &pb.ApiResponseWithdrawMonthStatusSuccess{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesWithdrawMonthStatusSuccess(pbResponse),
	}
}

func (m *withdrawProtoMapper) ToProtoResponseWithdrawYearStatusSuccess(status string, message string, pbResponse []*response.WithdrawResponseYearStatusSuccess) *pb.ApiResponseWithdrawYearStatusSuccess {
	return &pb.ApiResponseWithdrawYearStatusSuccess{
		Status:  status,
		Message: message,
		Data:    m.mapWithdrawResponsesYearStatusSuccess(pbResponse),
	}
}

func (m *withdrawProtoMapper) ToProtoResponseWithdrawMonthStatusFailed(status string, message string, pbResponse []*response.WithdrawResponseMonthStatusFailed) *pb.ApiResponseWithdrawMonthStatusFailed {
	return &pb.ApiResponseWithdrawMonthStatusFailed{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesWithdrawMonthStatusFailed(pbResponse),
	}
}

func (m *withdrawProtoMapper) ToProtoResponseWithdrawYearStatusFailed(status string, message string, pbResponse []*response.WithdrawResponseYearStatusFailed) *pb.ApiResponseWithdrawYearStatusFailed {
	return &pb.ApiResponseWithdrawYearStatusFailed{
		Status:  status,
		Message: message,
		Data:    m.mapWithdrawResponsesYearStatusFailed(pbResponse),
	}
}

func (m *withdrawProtoMapper) ToProtoResponseWithdrawMonthAmount(status string, message string, pbResponse []*response.WithdrawMonthlyAmountResponse) *pb.ApiResponseWithdrawMonthAmount {
	return &pb.ApiResponseWithdrawMonthAmount{
		Status:  status,
		Message: message,
		Data:    m.mapResponseWithdrawMonthlyAmounts(pbResponse),
	}
}

func (m *withdrawProtoMapper) ToProtoResponseWithdrawYearAmount(status string, message string, pbResponse []*response.WithdrawYearlyAmountResponse) *pb.ApiResponseWithdrawYearAmount {
	return &pb.ApiResponseWithdrawYearAmount{
		Status:  status,
		Message: message,
		Data:    m.mapResponseWithdrawYearlyAmounts(pbResponse),
	}
}

func (w *withdrawProtoMapper) mapResponseWithdrawal(withdraw *response.WithdrawResponse) *pb.WithdrawResponse {
	return &pb.WithdrawResponse{
		WithdrawId:     int32(withdraw.ID),
		WithdrawNo:     withdraw.WithdrawNo,
		CardNumber:     withdraw.CardNumber,
		WithdrawAmount: int32(withdraw.WithdrawAmount),
		WithdrawTime:   withdraw.WithdrawTime,
		CreatedAt:      withdraw.CreatedAt,
		UpdatedAt:      withdraw.UpdatedAt,
	}
}

func (w *withdrawProtoMapper) mapResponsesWithdrawal(withdraws []*response.WithdrawResponse) []*pb.WithdrawResponse {
	var responseWithdraws []*pb.WithdrawResponse

	for _, withdraw := range withdraws {
		responseWithdraws = append(responseWithdraws, w.mapResponseWithdrawal(withdraw))
	}

	return responseWithdraws
}

func (w *withdrawProtoMapper) mapResponseWithdrawalDeleteAt(withdraw *response.WithdrawResponseDeleteAt) *pb.WithdrawResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if withdraw.DeletedAt != nil {
		deletedAt = wrapperspb.String(*withdraw.DeletedAt)
	}

	return &pb.WithdrawResponseDeleteAt{
		WithdrawId:     int32(withdraw.ID),
		WithdrawNo:     withdraw.WithdrawNo,
		CardNumber:     withdraw.CardNumber,
		WithdrawAmount: int32(withdraw.WithdrawAmount),
		WithdrawTime:   withdraw.WithdrawTime,
		CreatedAt:      withdraw.CreatedAt,
		UpdatedAt:      withdraw.UpdatedAt,
		DeletedAt:      deletedAt,
	}
}

func (w *withdrawProtoMapper) mapResponsesWithdrawalDeleteAt(withdraws []*response.WithdrawResponseDeleteAt) []*pb.WithdrawResponseDeleteAt {
	var responseWithdraws []*pb.WithdrawResponseDeleteAt

	for _, withdraw := range withdraws {
		responseWithdraws = append(responseWithdraws, w.mapResponseWithdrawalDeleteAt(withdraw))
	}

	return responseWithdraws
}

func (t *withdrawProtoMapper) mapResponseWithdrawMonthStatusSuccess(s *response.WithdrawResponseMonthStatusSuccess) *pb.WithdrawMonthStatusSuccessResponse {
	return &pb.WithdrawMonthStatusSuccessResponse{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int32(s.TotalSuccess),
		TotalAmount:  int32(s.TotalAmount),
	}
}

func (t *withdrawProtoMapper) mapResponsesWithdrawMonthStatusSuccess(Withdraws []*response.WithdrawResponseMonthStatusSuccess) []*pb.WithdrawMonthStatusSuccessResponse {
	var WithdrawRecords []*pb.WithdrawMonthStatusSuccessResponse

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.mapResponseWithdrawMonthStatusSuccess(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawProtoMapper) mapWithdrawResponseYearStatusSuccess(s *response.WithdrawResponseYearStatusSuccess) *pb.WithdrawYearStatusSuccessResponse {
	return &pb.WithdrawYearStatusSuccessResponse{
		Year:         s.Year,
		TotalSuccess: int32(s.TotalSuccess),
		TotalAmount:  int32(s.TotalAmount),
	}
}

func (t *withdrawProtoMapper) mapWithdrawResponsesYearStatusSuccess(Withdraws []*response.WithdrawResponseYearStatusSuccess) []*pb.WithdrawYearStatusSuccessResponse {
	var WithdrawRecords []*pb.WithdrawYearStatusSuccessResponse

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.mapWithdrawResponseYearStatusSuccess(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawProtoMapper) mapResponseWithdrawMonthStatusFailed(s *response.WithdrawResponseMonthStatusFailed) *pb.WithdrawMonthStatusFailedResponse {
	return &pb.WithdrawMonthStatusFailedResponse{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int32(s.TotalFailed),
		TotalAmount: int32(s.TotalAmount),
	}
}

func (t *withdrawProtoMapper) mapResponsesWithdrawMonthStatusFailed(Withdraws []*response.WithdrawResponseMonthStatusFailed) []*pb.WithdrawMonthStatusFailedResponse {
	var WithdrawRecords []*pb.WithdrawMonthStatusFailedResponse

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.mapResponseWithdrawMonthStatusFailed(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawProtoMapper) mapWithdrawResponseYearStatusFailed(s *response.WithdrawResponseYearStatusFailed) *pb.WithdrawYearStatusFailedResponse {
	return &pb.WithdrawYearStatusFailedResponse{
		Year:        s.Year,
		TotalFailed: int32(s.TotalFailed),
		TotalAmount: int32(s.TotalAmount),
	}
}

func (t *withdrawProtoMapper) mapWithdrawResponsesYearStatusFailed(Withdraws []*response.WithdrawResponseYearStatusFailed) []*pb.WithdrawYearStatusFailedResponse {
	var WithdrawRecords []*pb.WithdrawYearStatusFailedResponse

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.mapWithdrawResponseYearStatusFailed(Withdraw))
	}

	return WithdrawRecords
}

func (m *withdrawProtoMapper) mapResponseWithdrawMonthlyAmount(s *response.WithdrawMonthlyAmountResponse) *pb.WithdrawMonthlyAmountResponse {
	return &pb.WithdrawMonthlyAmountResponse{
		Month:       s.Month,
		TotalAmount: int32(s.TotalAmount),
	}
}

func (m *withdrawProtoMapper) mapResponseWithdrawMonthlyAmounts(s []*response.WithdrawMonthlyAmountResponse) []*pb.WithdrawMonthlyAmountResponse {
	var protoResponses []*pb.WithdrawMonthlyAmountResponse
	for _, withdraw := range s {
		protoResponses = append(protoResponses, m.mapResponseWithdrawMonthlyAmount(withdraw))
	}
	return protoResponses
}

func (m *withdrawProtoMapper) mapResponseWithdrawYearlyAmount(s *response.WithdrawYearlyAmountResponse) *pb.WithdrawYearlyAmountResponse {
	return &pb.WithdrawYearlyAmountResponse{
		Year:        s.Year,
		TotalAmount: int32(s.TotalAmount),
	}
}

func (m *withdrawProtoMapper) mapResponseWithdrawYearlyAmounts(s []*response.WithdrawYearlyAmountResponse) []*pb.WithdrawYearlyAmountResponse {
	var protoResponses []*pb.WithdrawYearlyAmountResponse
	for _, withdraw := range s {
		protoResponses = append(protoResponses, m.mapResponseWithdrawYearlyAmount(withdraw))
	}
	return protoResponses
}
