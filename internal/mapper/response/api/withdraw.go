package apimapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"
)

type withdrawResponseMapper struct {
}

func NewWithdrawResponseMapper() *withdrawResponseMapper {
	return &withdrawResponseMapper{}
}

func (m *withdrawResponseMapper) ToApiResponseWithdraw(pbResponse *pb.ApiResponseWithdraw) *response.ApiResponseWithdraw {
	return &response.ApiResponseWithdraw{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponseWithdrawal(pbResponse.Data),
	}
}

func (m *withdrawResponseMapper) ToApiResponsesWithdraw(pbResponse *pb.ApiResponsesWithdraw) *response.ApiResponsesWithdraw {
	return &response.ApiResponsesWithdraw{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponsesWithdrawal(pbResponse.Data),
	}
}

func (m *withdrawResponseMapper) ToApiResponseWithdrawDelete(pbResponse *pb.ApiResponseWithdrawDelete) *response.ApiResponseWithdrawDelete {
	return &response.ApiResponseWithdrawDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (m *withdrawResponseMapper) ToApiResponseWithdrawAll(pbResponse *pb.ApiResponseWithdrawAll) *response.ApiResponseWithdrawAll {
	return &response.ApiResponseWithdrawAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (m *withdrawResponseMapper) ToApiResponsePaginationWithdraw(pbResponse *pb.ApiResponsePaginationWithdraw) *response.ApiResponsePaginationWithdraw {
	return &response.ApiResponsePaginationWithdraw{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.mapResponsesWithdrawal(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *withdrawResponseMapper) ToApiResponsePaginationWithdrawDeleteAt(pbResponse *pb.ApiResponsePaginationWithdrawDeleteAt) *response.ApiResponsePaginationWithdrawDeleteAt {
	return &response.ApiResponsePaginationWithdrawDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.mapResponsesWithdrawalDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *withdrawResponseMapper) ToApiResponseWithdrawMonthStatusSuccess(pbResponse *pb.ApiResponseWithdrawMonthStatusSuccess) *response.ApiResponseWithdrawMonthStatusSuccess {
	return &response.ApiResponseWithdrawMonthStatusSuccess{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponsesWithdrawMonthStatusSuccess(pbResponse.Data),
	}
}

func (m *withdrawResponseMapper) ToApiResponseWithdrawYearStatusSuccess(pbResponse *pb.ApiResponseWithdrawYearStatusSuccess) *response.ApiResponseWithdrawYearStatusSuccess {
	return &response.ApiResponseWithdrawYearStatusSuccess{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapWithdrawResponsesYearStatusSuccess(pbResponse.Data),
	}
}

func (m *withdrawResponseMapper) ToApiResponseWithdrawMonthStatusFailed(pbResponse *pb.ApiResponseWithdrawMonthStatusFailed) *response.ApiResponseWithdrawMonthStatusFailed {
	return &response.ApiResponseWithdrawMonthStatusFailed{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponsesWithdrawMonthStatusFailed(pbResponse.Data),
	}
}

func (m *withdrawResponseMapper) ToApiResponseWithdrawYearStatusFailed(pbResponse *pb.ApiResponseWithdrawYearStatusFailed) *response.ApiResponseWithdrawYearStatusFailed {
	return &response.ApiResponseWithdrawYearStatusFailed{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapWithdrawResponsesYearStatusFailed(pbResponse.Data),
	}
}

func (m *withdrawResponseMapper) ToApiResponseWithdrawMonthAmount(pbResponse *pb.ApiResponseWithdrawMonthAmount) *response.ApiResponseWithdrawMonthAmount {
	return &response.ApiResponseWithdrawMonthAmount{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponseWithdrawMonthlyAmounts(pbResponse.Data),
	}
}

func (m *withdrawResponseMapper) ToApiResponseWithdrawYearAmount(pbResponse *pb.ApiResponseWithdrawYearAmount) *response.ApiResponseWithdrawYearAmount {
	return &response.ApiResponseWithdrawYearAmount{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponseWithdrawYearlyAmounts(pbResponse.Data),
	}
}

func (w *withdrawResponseMapper) mapResponseWithdrawal(withdraw *pb.WithdrawResponse) *response.WithdrawResponse {
	return &response.WithdrawResponse{
		ID:             int(withdraw.WithdrawId),
		WithdrawNo:     withdraw.WithdrawNo,
		CardNumber:     withdraw.CardNumber,
		WithdrawAmount: int(withdraw.WithdrawAmount),
		WithdrawTime:   withdraw.WithdrawTime,
		CreatedAt:      withdraw.CreatedAt,
		UpdatedAt:      withdraw.UpdatedAt,
	}
}

func (w *withdrawResponseMapper) mapResponsesWithdrawal(withdraws []*pb.WithdrawResponse) []*response.WithdrawResponse {
	var responseWithdraws []*response.WithdrawResponse

	for _, withdraw := range withdraws {
		responseWithdraws = append(responseWithdraws, w.mapResponseWithdrawal(withdraw))
	}

	return responseWithdraws
}

func (w *withdrawResponseMapper) mapResponseWithdrawalDeleteAt(withdraw *pb.WithdrawResponseDeleteAt) *response.WithdrawResponseDeleteAt {
	var deletedAt string
	if withdraw.DeletedAt != nil {
		deletedAt = withdraw.DeletedAt.Value
	}

	return &response.WithdrawResponseDeleteAt{
		ID:             int(withdraw.WithdrawId),
		WithdrawNo:     withdraw.WithdrawNo,
		CardNumber:     withdraw.CardNumber,
		WithdrawAmount: int(withdraw.WithdrawAmount),
		WithdrawTime:   withdraw.WithdrawTime,
		CreatedAt:      withdraw.CreatedAt,
		UpdatedAt:      withdraw.UpdatedAt,
		DeletedAt:      &deletedAt,
	}
}

func (w *withdrawResponseMapper) mapResponsesWithdrawalDeleteAt(withdraws []*pb.WithdrawResponseDeleteAt) []*response.WithdrawResponseDeleteAt {
	var responseWithdraws []*response.WithdrawResponseDeleteAt

	for _, withdraw := range withdraws {
		responseWithdraws = append(responseWithdraws, w.mapResponseWithdrawalDeleteAt(withdraw))
	}

	return responseWithdraws
}

func (t *withdrawResponseMapper) mapResponseWithdrawMonthStatusSuccess(s *pb.WithdrawMonthStatusSuccessResponse) *response.WithdrawResponseMonthStatusSuccess {
	return &response.WithdrawResponseMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *withdrawResponseMapper) mapResponsesWithdrawMonthStatusSuccess(Withdraws []*pb.WithdrawMonthStatusSuccessResponse) []*response.WithdrawResponseMonthStatusSuccess {
	var WithdrawRecords []*response.WithdrawResponseMonthStatusSuccess

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.mapResponseWithdrawMonthStatusSuccess(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawResponseMapper) mapWithdrawResponseYearStatusSuccess(s *pb.WithdrawYearStatusSuccessResponse) *response.WithdrawResponseYearStatusSuccess {
	return &response.WithdrawResponseYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *withdrawResponseMapper) mapWithdrawResponsesYearStatusSuccess(Withdraws []*pb.WithdrawYearStatusSuccessResponse) []*response.WithdrawResponseYearStatusSuccess {
	var WithdrawRecords []*response.WithdrawResponseYearStatusSuccess

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.mapWithdrawResponseYearStatusSuccess(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawResponseMapper) mapResponseWithdrawMonthStatusFailed(s *pb.WithdrawMonthStatusFailedResponse) *response.WithdrawResponseMonthStatusFailed {
	return &response.WithdrawResponseMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *withdrawResponseMapper) mapResponsesWithdrawMonthStatusFailed(Withdraws []*pb.WithdrawMonthStatusFailedResponse) []*response.WithdrawResponseMonthStatusFailed {
	var WithdrawRecords []*response.WithdrawResponseMonthStatusFailed

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.mapResponseWithdrawMonthStatusFailed(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawResponseMapper) mapWithdrawResponseYearStatusFailed(s *pb.WithdrawYearStatusFailedResponse) *response.WithdrawResponseYearStatusFailed {
	return &response.WithdrawResponseYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *withdrawResponseMapper) mapWithdrawResponsesYearStatusFailed(Withdraws []*pb.WithdrawYearStatusFailedResponse) []*response.WithdrawResponseYearStatusFailed {
	var WithdrawRecords []*response.WithdrawResponseYearStatusFailed

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.mapWithdrawResponseYearStatusFailed(Withdraw))
	}

	return WithdrawRecords
}

func (m *withdrawResponseMapper) mapResponseWithdrawMonthlyAmount(s *pb.WithdrawMonthlyAmountResponse) *response.WithdrawMonthlyAmountResponse {
	return &response.WithdrawMonthlyAmountResponse{
		Month:       s.Month,
		TotalAmount: int(s.TotalAmount),
	}
}

func (m *withdrawResponseMapper) mapResponseWithdrawMonthlyAmounts(s []*pb.WithdrawMonthlyAmountResponse) []*response.WithdrawMonthlyAmountResponse {
	var protoResponses []*response.WithdrawMonthlyAmountResponse
	for _, withdraw := range s {
		protoResponses = append(protoResponses, m.mapResponseWithdrawMonthlyAmount(withdraw))
	}
	return protoResponses
}

func (m *withdrawResponseMapper) mapResponseWithdrawYearlyAmount(s *pb.WithdrawYearlyAmountResponse) *response.WithdrawYearlyAmountResponse {
	return &response.WithdrawYearlyAmountResponse{
		Year:        s.Year,
		TotalAmount: int(s.TotalAmount),
	}
}

func (m *withdrawResponseMapper) mapResponseWithdrawYearlyAmounts(s []*pb.WithdrawYearlyAmountResponse) []*response.WithdrawYearlyAmountResponse {
	var protoResponses []*response.WithdrawYearlyAmountResponse
	for _, withdraw := range s {
		protoResponses = append(protoResponses, m.mapResponseWithdrawYearlyAmount(withdraw))
	}
	return protoResponses
}
