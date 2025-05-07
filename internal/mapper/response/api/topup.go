package apimapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"
)

type topupResponseMapper struct {
}

func NewTopupResponseMapper() *topupResponseMapper {
	return &topupResponseMapper{}
}

func mapPaginationMeta(s *pb.PaginationMeta) *response.PaginationMeta {
	return &response.PaginationMeta{
		CurrentPage:  int(s.CurrentPage),
		PageSize:     int(s.PageSize),
		TotalRecords: int(s.TotalRecords),
		TotalPages:   int(s.TotalPages),
	}
}

func (t *topupResponseMapper) ToApiResponseTopup(s *pb.ApiResponseTopup) *response.ApiResponseTopup {
	return &response.ApiResponseTopup{
		Status:  s.Status,
		Message: s.Message,
		Data:    t.mapResponseTopup(s.Data),
	}
}

func (t *topupResponseMapper) ToApiResponseTopupDeleteAt(s *pb.ApiResponseTopupDeleteAt) *response.ApiResponseTopupDeleteAt {
	return &response.ApiResponseTopupDeleteAt{
		Status:  s.Status,
		Message: s.Message,
		Data:    t.mapResponseTopupDeleteAt(s.Data),
	}
}

func (t *topupResponseMapper) ToApiResponseTopupAll(s *pb.ApiResponseTopupAll) *response.ApiResponseTopupAll {
	return &response.ApiResponseTopupAll{
		Status:  s.Status,
		Message: s.Message,
	}
}

func (t *topupResponseMapper) ToApiResponseTopupDelete(s *pb.ApiResponseTopupDelete) *response.ApiResponseTopupDelete {
	return &response.ApiResponseTopupDelete{
		Status:  s.Status,
		Message: s.Message,
	}
}

func (t *topupResponseMapper) ToApiResponsePaginationTopup(s *pb.ApiResponsePaginationTopup) *response.ApiResponsePaginationTopup {
	return &response.ApiResponsePaginationTopup{
		Status:     s.Status,
		Message:    s.Message,
		Data:       t.mapResponsesTopup(s.Data),
		Pagination: mapPaginationMeta(s.Pagination),
	}
}

func (t *topupResponseMapper) ToApiResponsePaginationTopupDeleteAt(s *pb.ApiResponsePaginationTopupDeleteAt) *response.ApiResponsePaginationTopupDeleteAt {
	return &response.ApiResponsePaginationTopupDeleteAt{
		Status:     s.Status,
		Message:    s.Message,
		Data:       t.mapResponsesTopupDeleteAt(s.Data),
		Pagination: mapPaginationMeta(s.Pagination),
	}
}

func (t *topupResponseMapper) ToApiResponseTopupMonthStatusSuccess(s *pb.ApiResponseTopupMonthStatusSuccess) *response.ApiResponseTopupMonthStatusSuccess {
	return &response.ApiResponseTopupMonthStatusSuccess{
		Status:  s.Status,
		Message: s.Message,
		Data:    t.mapResponsesTopupMonthStatusSuccess(s.Data),
	}
}

func (t *topupResponseMapper) ToApiResponseTopupYearStatusSuccess(s *pb.ApiResponseTopupYearStatusSuccess) *response.ApiResponseTopupYearStatusSuccess {
	return &response.ApiResponseTopupYearStatusSuccess{
		Status:  s.Status,
		Message: s.Message,
		Data:    t.mapTopupResponsesYearStatusSuccess(s.Data),
	}
}

func (t *topupResponseMapper) ToApiResponseTopupMonthStatusFailed(s *pb.ApiResponseTopupMonthStatusFailed) *response.ApiResponseTopupMonthStatusFailed {
	return &response.ApiResponseTopupMonthStatusFailed{
		Status:  s.Status,
		Message: s.Message,
		Data:    t.mapResponsesTopupMonthStatusFailed(s.Data),
	}
}

func (t *topupResponseMapper) ToApiResponseTopupYearStatusFailed(s *pb.ApiResponseTopupYearStatusFailed) *response.ApiResponseTopupYearStatusFailed {
	return &response.ApiResponseTopupYearStatusFailed{
		Status:  s.Status,
		Message: s.Message,
		Data:    t.mapTopupResponsesYearStatusFailed(s.Data),
	}
}

func (t *topupResponseMapper) ToApiResponseTopupMonthMethod(s *pb.ApiResponseTopupMonthMethod) *response.ApiResponseTopupMonthMethod {
	return &response.ApiResponseTopupMonthMethod{
		Status:  s.Status,
		Message: s.Message,
		Data:    t.mapResponseTopupMonthlyMethods(s.Data),
	}
}

func (t *topupResponseMapper) ToApiResponseTopupYearMethod(s *pb.ApiResponseTopupYearMethod) *response.ApiResponseTopupYearMethod {
	return &response.ApiResponseTopupYearMethod{
		Status:  s.Status,
		Message: s.Message,
		Data:    t.mapResponseTopupYearlyMethods(s.Data),
	}
}

func (t *topupResponseMapper) ToApiResponseTopupMonthAmount(s *pb.ApiResponseTopupMonthAmount) *response.ApiResponseTopupMonthAmount {
	return &response.ApiResponseTopupMonthAmount{
		Status:  s.Status,
		Message: s.Message,
		Data:    t.mapResponseTopupMonthlyAmounts(s.Data),
	}
}

func (t *topupResponseMapper) ToApiResponseTopupYearAmount(s *pb.ApiResponseTopupYearAmount) *response.ApiResponseTopupYearAmount {
	return &response.ApiResponseTopupYearAmount{
		Status:  s.Status,
		Message: s.Message,
		Data:    t.mapResponseTopupYearlyAmounts(s.Data),
	}
}

func (t *topupResponseMapper) mapResponseTopup(topup *pb.TopupResponse) *response.TopupResponse {
	return &response.TopupResponse{
		ID:          int(topup.Id),
		CardNumber:  topup.CardNumber,
		TopupNo:     topup.TopupNo,
		TopupAmount: int(topup.TopupAmount),
		TopupMethod: topup.TopupMethod,
		TopupTime:   topup.TopupTime,
		CreatedAt:   topup.CreatedAt,
		UpdatedAt:   topup.UpdatedAt,
	}
}

func (t *topupResponseMapper) mapResponsesTopup(topups []*pb.TopupResponse) []*response.TopupResponse {
	var responses []*response.TopupResponse

	for _, topup := range topups {
		responses = append(responses, t.mapResponseTopup(topup))
	}

	return responses
}

func (t *topupResponseMapper) mapResponseTopupDeleteAt(topup *pb.TopupResponseDeleteAt) *response.TopupResponseDeleteAt {
	var deletedAt string
	if topup.DeletedAt != nil {
		deletedAt = topup.DeletedAt.Value
	}

	return &response.TopupResponseDeleteAt{
		ID:          int(topup.Id),
		CardNumber:  topup.CardNumber,
		TopupNo:     topup.TopupNo,
		TopupAmount: int(topup.TopupAmount),
		TopupMethod: topup.TopupMethod,
		TopupTime:   topup.TopupTime,
		CreatedAt:   topup.CreatedAt,
		UpdatedAt:   topup.UpdatedAt,
		DeletedAt:   &deletedAt,
	}
}

func (t *topupResponseMapper) mapResponsesTopupDeleteAt(topups []*pb.TopupResponseDeleteAt) []*response.TopupResponseDeleteAt {
	var responses []*response.TopupResponseDeleteAt

	for _, topup := range topups {
		responses = append(responses, t.mapResponseTopupDeleteAt(topup))
	}

	return responses
}

func (t *topupResponseMapper) mapResponseTopupMonthStatusSuccess(s *pb.TopupMonthStatusSuccessResponse) *response.TopupResponseMonthStatusSuccess {
	return &response.TopupResponseMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *topupResponseMapper) mapResponsesTopupMonthStatusSuccess(topups []*pb.TopupMonthStatusSuccessResponse) []*response.TopupResponseMonthStatusSuccess {
	var topupRecords []*response.TopupResponseMonthStatusSuccess

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.mapResponseTopupMonthStatusSuccess(topup))
	}

	return topupRecords
}

func (t *topupResponseMapper) mapTopupResponseYearStatusSuccess(s *pb.TopupYearStatusSuccessResponse) *response.TopupResponseYearStatusSuccess {
	return &response.TopupResponseYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *topupResponseMapper) mapTopupResponsesYearStatusSuccess(topups []*pb.TopupYearStatusSuccessResponse) []*response.TopupResponseYearStatusSuccess {
	var topupRecords []*response.TopupResponseYearStatusSuccess

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.mapTopupResponseYearStatusSuccess(topup))
	}

	return topupRecords
}

func (t *topupResponseMapper) mapResponseTopupMonthStatusFailed(s *pb.TopupMonthStatusFailedResponse) *response.TopupResponseMonthStatusFailed {
	return &response.TopupResponseMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupResponseMapper) mapResponsesTopupMonthStatusFailed(topups []*pb.TopupMonthStatusFailedResponse) []*response.TopupResponseMonthStatusFailed {
	var topupRecords []*response.TopupResponseMonthStatusFailed

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.mapResponseTopupMonthStatusFailed(topup))
	}

	return topupRecords
}

func (t *topupResponseMapper) mapTopupResponseYearStatusFailed(s *pb.TopupYearStatusFailedResponse) *response.TopupResponseYearStatusFailed {
	return &response.TopupResponseYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupResponseMapper) mapTopupResponsesYearStatusFailed(topups []*pb.TopupYearStatusFailedResponse) []*response.TopupResponseYearStatusFailed {
	var topupRecords []*response.TopupResponseYearStatusFailed

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.mapTopupResponseYearStatusFailed(topup))
	}

	return topupRecords
}

func (t *topupResponseMapper) mapResponseTopupMonthlyMethod(s *pb.TopupMonthMethodResponse) *response.TopupMonthMethodResponse {
	return &response.TopupMonthMethodResponse{
		Month:       s.Month,
		TopupMethod: s.TopupMethod,
		TotalTopups: int(s.TotalTopups),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupResponseMapper) mapResponseTopupMonthlyMethods(s []*pb.TopupMonthMethodResponse) []*response.TopupMonthMethodResponse {
	var topupProtos []*response.TopupMonthMethodResponse
	for _, topup := range s {
		topupProtos = append(topupProtos, t.mapResponseTopupMonthlyMethod(topup))
	}
	return topupProtos
}

func (t *topupResponseMapper) mapResponseTopupYearlyMethod(s *pb.TopupYearlyMethodResponse) *response.TopupYearlyMethodResponse {
	return &response.TopupYearlyMethodResponse{
		Year:        s.Year,
		TopupMethod: s.TopupMethod,
		TotalTopups: int(s.TotalTopups),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupResponseMapper) mapResponseTopupYearlyMethods(s []*pb.TopupYearlyMethodResponse) []*response.TopupYearlyMethodResponse {
	var topupProtos []*response.TopupYearlyMethodResponse
	for _, topup := range s {
		topupProtos = append(topupProtos, t.mapResponseTopupYearlyMethod(topup))
	}
	return topupProtos
}

func (t *topupResponseMapper) mapResponseTopupMonthlyAmount(s *pb.TopupMonthAmountResponse) *response.TopupMonthAmountResponse {
	return &response.TopupMonthAmountResponse{
		Month:       s.Month,
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupResponseMapper) mapResponseTopupMonthlyAmounts(s []*pb.TopupMonthAmountResponse) []*response.TopupMonthAmountResponse {
	var topupProtos []*response.TopupMonthAmountResponse
	for _, topup := range s {
		topupProtos = append(topupProtos, t.mapResponseTopupMonthlyAmount(topup))
	}
	return topupProtos
}

func (t *topupResponseMapper) mapResponseTopupYearlyAmount(s *pb.TopupYearlyAmountResponse) *response.TopupYearlyAmountResponse {
	return &response.TopupYearlyAmountResponse{
		Year:        s.Year,
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupResponseMapper) mapResponseTopupYearlyAmounts(s []*pb.TopupYearlyAmountResponse) []*response.TopupYearlyAmountResponse {
	var topupProtos []*response.TopupYearlyAmountResponse
	for _, topup := range s {
		topupProtos = append(topupProtos, t.mapResponseTopupYearlyAmount(topup))
	}
	return topupProtos
}
