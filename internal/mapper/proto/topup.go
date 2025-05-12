package protomapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type topupProtoMapper struct {
}

func NewTopupProtoMapper() *topupProtoMapper {
	return &topupProtoMapper{}
}

func (t *topupProtoMapper) ToProtoResponsePaginationTopup(pagination *pb.PaginationMeta, status string, message string, s []*response.TopupResponse) *pb.ApiResponsePaginationTopup {
	return &pb.ApiResponsePaginationTopup{
		Status:     status,
		Message:    message,
		Data:       t.mapResponsesTopup(s),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (t *topupProtoMapper) ToProtoResponsePaginationTopupDeleteAt(pagination *pb.PaginationMeta, status string, message string, s []*response.TopupResponseDeleteAt) *pb.ApiResponsePaginationTopupDeleteAt {
	return &pb.ApiResponsePaginationTopupDeleteAt{
		Status:     status,
		Message:    message,
		Data:       t.mapResponsesTopupDeleteAt(s),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (t *topupProtoMapper) ToProtoResponseTopupMonthStatusSuccess(status string, message string, s []*response.TopupResponseMonthStatusSuccess) *pb.ApiResponseTopupMonthStatusSuccess {
	return &pb.ApiResponseTopupMonthStatusSuccess{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTopupMonthStatusSuccess(s),
	}
}

func (t *topupProtoMapper) ToProtoResponseTopupYearStatusSuccess(status string, message string, s []*response.TopupResponseYearStatusSuccess) *pb.ApiResponseTopupYearStatusSuccess {
	return &pb.ApiResponseTopupYearStatusSuccess{
		Status:  status,
		Message: message,
		Data:    t.mapTopupResponsesYearStatusSuccess(s),
	}
}

func (t *topupProtoMapper) ToProtoResponseTopupMonthStatusFailed(status string, message string, s []*response.TopupResponseMonthStatusFailed) *pb.ApiResponseTopupMonthStatusFailed {
	return &pb.ApiResponseTopupMonthStatusFailed{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTopupMonthStatusFailed(s),
	}
}

func (t *topupProtoMapper) ToProtoResponseTopupYearStatusFailed(status string, message string, s []*response.TopupResponseYearStatusFailed) *pb.ApiResponseTopupYearStatusFailed {
	return &pb.ApiResponseTopupYearStatusFailed{
		Status:  status,
		Message: message,
		Data:    t.mapTopupResponsesYearStatusFailed(s),
	}
}

func (t *topupProtoMapper) ToProtoResponseTopupMonthMethod(status string, message string, s []*response.TopupMonthMethodResponse) *pb.ApiResponseTopupMonthMethod {
	return &pb.ApiResponseTopupMonthMethod{
		Status:  status,
		Message: message,
		Data:    t.mapResponseTopupMonthlyMethods(s),
	}
}

func (t *topupProtoMapper) ToProtoResponseTopupYearMethod(status string, message string, s []*response.TopupYearlyMethodResponse) *pb.ApiResponseTopupYearMethod {
	return &pb.ApiResponseTopupYearMethod{
		Status:  status,
		Message: message,
		Data:    t.mapResponseTopupYearlyMethods(s),
	}
}

func (t *topupProtoMapper) ToProtoResponseTopupMonthAmount(status string, message string, s []*response.TopupMonthAmountResponse) *pb.ApiResponseTopupMonthAmount {
	return &pb.ApiResponseTopupMonthAmount{
		Status:  status,
		Message: message,
		Data:    t.mapResponseTopupMonthlyAmounts(s),
	}
}

func (t *topupProtoMapper) ToProtoResponseTopupYearAmount(status string, message string, s []*response.TopupYearlyAmountResponse) *pb.ApiResponseTopupYearAmount {
	return &pb.ApiResponseTopupYearAmount{
		Status:  status,
		Message: message,
		Data:    t.mapResponseTopupYearlyAmounts(s),
	}
}

func (t *topupProtoMapper) ToProtoResponseTopup(status string, message string, s *response.TopupResponse) *pb.ApiResponseTopup {
	return &pb.ApiResponseTopup{
		Status:  status,
		Message: message,
		Data:    t.mapResponseTopup(s),
	}
}

func (t *topupProtoMapper) ToProtoResponseTopupDeletAt(status string, message string, s *response.TopupResponseDeleteAt) *pb.ApiResponseTopupDeleteAt {
	return &pb.ApiResponseTopupDeleteAt{
		Status:  status,
		Message: message,
		Data:    t.mapResponseTopupDeleteAt(s),
	}
}

func (t topupProtoMapper) ToProtoResponseTopupDelete(status string, message string) *pb.ApiResponseTopupDelete {
	return &pb.ApiResponseTopupDelete{
		Status:  status,
		Message: message,
	}
}

func (t topupProtoMapper) ToProtoResponseTopupAll(status string, message string) *pb.ApiResponseTopupAll {
	return &pb.ApiResponseTopupAll{
		Status:  status,
		Message: message,
	}
}

func (t *topupProtoMapper) mapResponseTopup(topup *response.TopupResponse) *pb.TopupResponse {
	return &pb.TopupResponse{
		Id:          int32(topup.ID),
		CardNumber:  topup.CardNumber,
		TopupNo:     topup.TopupNo,
		TopupAmount: int32(topup.TopupAmount),
		TopupMethod: topup.TopupMethod,
		TopupTime:   topup.TopupTime,
		CreatedAt:   topup.CreatedAt,
		UpdatedAt:   topup.UpdatedAt,
	}
}

func (t *topupProtoMapper) mapResponsesTopup(topups []*response.TopupResponse) []*pb.TopupResponse {
	var responses []*pb.TopupResponse

	for _, response := range topups {
		responses = append(responses, t.mapResponseTopup(response))
	}

	return responses
}

func (t *topupProtoMapper) mapResponseTopupDeleteAt(topup *response.TopupResponseDeleteAt) *pb.TopupResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if topup.DeletedAt != nil {
		deletedAt = wrapperspb.String(*topup.DeletedAt)
	}

	return &pb.TopupResponseDeleteAt{
		Id:          int32(topup.ID),
		CardNumber:  topup.CardNumber,
		TopupNo:     topup.TopupNo,
		TopupAmount: int32(topup.TopupAmount),
		TopupMethod: topup.TopupMethod,
		TopupTime:   topup.TopupTime,
		CreatedAt:   topup.CreatedAt,
		UpdatedAt:   topup.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

func (t *topupProtoMapper) mapResponsesTopupDeleteAt(topups []*response.TopupResponseDeleteAt) []*pb.TopupResponseDeleteAt {
	var responses []*pb.TopupResponseDeleteAt

	for _, response := range topups {
		responses = append(responses, t.mapResponseTopupDeleteAt(response))
	}

	return responses
}

func (t *topupProtoMapper) mapResponseTopupMonthStatusSuccess(s *response.TopupResponseMonthStatusSuccess) *pb.TopupMonthStatusSuccessResponse {
	return &pb.TopupMonthStatusSuccessResponse{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int32(s.TotalSuccess),
		TotalAmount:  int32(s.TotalAmount),
	}
}

func (t *topupProtoMapper) mapResponsesTopupMonthStatusSuccess(topups []*response.TopupResponseMonthStatusSuccess) []*pb.TopupMonthStatusSuccessResponse {
	var topupRecords []*pb.TopupMonthStatusSuccessResponse

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.mapResponseTopupMonthStatusSuccess(topup))
	}

	return topupRecords
}

func (t *topupProtoMapper) mapTopupResponseYearStatusSuccess(s *response.TopupResponseYearStatusSuccess) *pb.TopupYearStatusSuccessResponse {
	return &pb.TopupYearStatusSuccessResponse{
		Year:         s.Year,
		TotalSuccess: int32(s.TotalSuccess),
		TotalAmount:  int32(s.TotalAmount),
	}
}

func (t *topupProtoMapper) mapTopupResponsesYearStatusSuccess(topups []*response.TopupResponseYearStatusSuccess) []*pb.TopupYearStatusSuccessResponse {
	var topupRecords []*pb.TopupYearStatusSuccessResponse

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.mapTopupResponseYearStatusSuccess(topup))
	}

	return topupRecords
}

func (t *topupProtoMapper) mapResponseTopupMonthStatusFailed(s *response.TopupResponseMonthStatusFailed) *pb.TopupMonthStatusFailedResponse {
	return &pb.TopupMonthStatusFailedResponse{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int32(s.TotalFailed),
		TotalAmount: int32(s.TotalAmount),
	}
}

func (t *topupProtoMapper) mapResponsesTopupMonthStatusFailed(topups []*response.TopupResponseMonthStatusFailed) []*pb.TopupMonthStatusFailedResponse {
	var topupRecords []*pb.TopupMonthStatusFailedResponse

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.mapResponseTopupMonthStatusFailed(topup))
	}

	return topupRecords
}

func (t *topupProtoMapper) mapTopupResponseYearStatusFailed(s *response.TopupResponseYearStatusFailed) *pb.TopupYearStatusFailedResponse {
	return &pb.TopupYearStatusFailedResponse{
		Year:        s.Year,
		TotalFailed: int32(s.TotalFailed),
		TotalAmount: int32(s.TotalAmount),
	}
}

func (t *topupProtoMapper) mapTopupResponsesYearStatusFailed(topups []*response.TopupResponseYearStatusFailed) []*pb.TopupYearStatusFailedResponse {
	var topupRecords []*pb.TopupYearStatusFailedResponse

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.mapTopupResponseYearStatusFailed(topup))
	}

	return topupRecords
}

func (t *topupProtoMapper) mapResponseTopupMonthlyMethod(s *response.TopupMonthMethodResponse) *pb.TopupMonthMethodResponse {
	return &pb.TopupMonthMethodResponse{
		Month:       s.Month,
		TopupMethod: s.TopupMethod,
		TotalTopups: int32(s.TotalTopups),
		TotalAmount: int32(s.TotalAmount),
	}
}

func (t *topupProtoMapper) mapResponseTopupMonthlyMethods(s []*response.TopupMonthMethodResponse) []*pb.TopupMonthMethodResponse {
	var topupProtos []*pb.TopupMonthMethodResponse
	for _, topup := range s {
		topupProtos = append(topupProtos, t.mapResponseTopupMonthlyMethod(topup))
	}
	return topupProtos
}

func (t *topupProtoMapper) mapResponseTopupYearlyMethod(s *response.TopupYearlyMethodResponse) *pb.TopupYearlyMethodResponse {
	return &pb.TopupYearlyMethodResponse{
		Year:        s.Year,
		TopupMethod: s.TopupMethod,
		TotalTopups: int32(s.TotalTopups),
		TotalAmount: int32(s.TotalAmount),
	}
}

func (t *topupProtoMapper) mapResponseTopupYearlyMethods(s []*response.TopupYearlyMethodResponse) []*pb.TopupYearlyMethodResponse {
	var topupProtos []*pb.TopupYearlyMethodResponse
	for _, topup := range s {
		topupProtos = append(topupProtos, t.mapResponseTopupYearlyMethod(topup))
	}
	return topupProtos
}

func (t *topupProtoMapper) mapResponseTopupMonthlyAmount(s *response.TopupMonthAmountResponse) *pb.TopupMonthAmountResponse {
	return &pb.TopupMonthAmountResponse{
		Month:       s.Month,
		TotalAmount: int32(s.TotalAmount),
	}
}

func (t *topupProtoMapper) mapResponseTopupMonthlyAmounts(s []*response.TopupMonthAmountResponse) []*pb.TopupMonthAmountResponse {
	var topupProtos []*pb.TopupMonthAmountResponse
	for _, topup := range s {
		topupProtos = append(topupProtos, t.mapResponseTopupMonthlyAmount(topup))
	}
	return topupProtos
}

func (t *topupProtoMapper) mapResponseTopupYearlyAmount(s *response.TopupYearlyAmountResponse) *pb.TopupYearlyAmountResponse {
	return &pb.TopupYearlyAmountResponse{
		Year:        s.Year,
		TotalAmount: int32(s.TotalAmount),
	}
}

func (t *topupProtoMapper) mapResponseTopupYearlyAmounts(s []*response.TopupYearlyAmountResponse) []*pb.TopupYearlyAmountResponse {
	var topupProtos []*pb.TopupYearlyAmountResponse
	for _, topup := range s {
		topupProtos = append(topupProtos, t.mapResponseTopupYearlyAmount(topup))
	}
	return topupProtos
}
