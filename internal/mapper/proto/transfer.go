package protomapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type transferProtoMapper struct{}

func NewTransferProtoMapper() *transferProtoMapper {
	return &transferProtoMapper{}
}

func (m *transferProtoMapper) ToProtoResponseTransferMonthStatusSuccess(status string, message string, pbResponse []*response.TransferResponseMonthStatusSuccess) *pb.ApiResponseTransferMonthStatusSuccess {
	return &pb.ApiResponseTransferMonthStatusSuccess{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesTransferMonthStatusSuccess(pbResponse),
	}
}

func (m *transferProtoMapper) ToProtoResponseTransferYearStatusSuccess(status string, message string, pbResponse []*response.TransferResponseYearStatusSuccess) *pb.ApiResponseTransferYearStatusSuccess {
	return &pb.ApiResponseTransferYearStatusSuccess{
		Status:  status,
		Message: message,
		Data:    m.mapTransferResponsesYearStatusSuccess(pbResponse),
	}
}

func (m *transferProtoMapper) ToProtoResponseTransferMonthStatusFailed(status string, message string, pbResponse []*response.TransferResponseMonthStatusFailed) *pb.ApiResponseTransferMonthStatusFailed {
	return &pb.ApiResponseTransferMonthStatusFailed{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesTransferMonthStatusFailed(pbResponse),
	}
}

func (m *transferProtoMapper) ToProtoResponseTransferYearStatusFailed(status string, message string, pbResponse []*response.TransferResponseYearStatusFailed) *pb.ApiResponseTransferYearStatusFailed {
	return &pb.ApiResponseTransferYearStatusFailed{
		Status:  status,
		Message: message,
		Data:    m.mapTransferResponsesYearStatusFailed(pbResponse),
	}
}

func (m *transferProtoMapper) ToProtoResponseTransferMonthAmount(status string, message string, pbResponse []*response.TransferMonthAmountResponse) *pb.ApiResponseTransferMonthAmount {
	return &pb.ApiResponseTransferMonthAmount{
		Status:  status,
		Message: message,
		Data:    m.mapResponseTransferMonthAmounts(pbResponse),
	}
}

func (m *transferProtoMapper) ToProtoResponseTransferYearAmount(status string, message string, pbResponse []*response.TransferYearAmountResponse) *pb.ApiResponseTransferYearAmount {
	return &pb.ApiResponseTransferYearAmount{
		Status:  status,
		Message: message,
		Data:    m.mapResponseTransferYearAmounts(pbResponse),
	}
}

func (m *transferProtoMapper) ToProtoResponseTransfer(status string, message string, pbResponse *response.TransferResponse) *pb.ApiResponseTransfer {
	return &pb.ApiResponseTransfer{
		Status:  status,
		Message: message,
		Data:    m.mapResponseTransfer(pbResponse),
	}
}

func (m *transferProtoMapper) ToProtoResponseTransfers(status string, message string, pbResponse []*response.TransferResponse) *pb.ApiResponseTransfers {
	return &pb.ApiResponseTransfers{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesTransfer(pbResponse),
	}
}

func (m *transferProtoMapper) ToProtoResponseTransferDelete(status string, message string) *pb.ApiResponseTransferDelete {
	return &pb.ApiResponseTransferDelete{
		Status:  status,
		Message: message,
	}
}

func (m *transferProtoMapper) ToProtoResponseTransferAll(status string, message string) *pb.ApiResponseTransferAll {
	return &pb.ApiResponseTransferAll{
		Status:  status,
		Message: message,
	}
}

func (m *transferProtoMapper) ToProtoResponsePaginationTransfer(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.TransferResponse) *pb.ApiResponsePaginationTransfer {
	return &pb.ApiResponsePaginationTransfer{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesTransfer(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *transferProtoMapper) ToProtoResponsePaginationTransferDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.TransferResponseDeleteAt) *pb.ApiResponsePaginationTransferDeleteAt {
	return &pb.ApiResponsePaginationTransferDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesTransferDeleteAt(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (t *transferProtoMapper) mapResponseTransfer(transfer *response.TransferResponse) *pb.TransferResponse {
	return &pb.TransferResponse{
		Id:             int32(transfer.ID),
		TransferNo:     transfer.TransferNo,
		TransferFrom:   transfer.TransferFrom,
		TransferTo:     transfer.TransferTo,
		TransferAmount: int32(transfer.TransferAmount),
		TransferTime:   transfer.TransferTime,
		CreatedAt:      transfer.CreatedAt,
		UpdatedAt:      transfer.UpdatedAt,
	}
}

func (t *transferProtoMapper) mapResponsesTransfer(transfers []*response.TransferResponse) []*pb.TransferResponse {
	var responses []*pb.TransferResponse

	for _, response := range transfers {
		responses = append(responses, t.mapResponseTransfer(response))
	}

	return responses
}

func (t *transferProtoMapper) mapResponseTransferDeleteAt(transfer *response.TransferResponseDeleteAt) *pb.TransferResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if transfer.DeletedAt != nil {
		deletedAt = wrapperspb.String(*transfer.DeletedAt)
	}

	return &pb.TransferResponseDeleteAt{
		Id:             int32(transfer.ID),
		TransferNo:     transfer.TransferNo,
		TransferFrom:   transfer.TransferFrom,
		TransferTo:     transfer.TransferTo,
		TransferAmount: int32(transfer.TransferAmount),
		TransferTime:   transfer.TransferTime,
		CreatedAt:      transfer.CreatedAt,
		UpdatedAt:      transfer.UpdatedAt,
		DeletedAt:      deletedAt,
	}
}

func (t *transferProtoMapper) mapResponsesTransferDeleteAt(transfers []*response.TransferResponseDeleteAt) []*pb.TransferResponseDeleteAt {
	var responses []*pb.TransferResponseDeleteAt

	for _, response := range transfers {
		responses = append(responses, t.mapResponseTransferDeleteAt(response))
	}

	return responses
}

func (t *transferProtoMapper) mapResponseTransferMonthStatusSuccess(s *response.TransferResponseMonthStatusSuccess) *pb.TransferMonthStatusSuccessResponse {
	return &pb.TransferMonthStatusSuccessResponse{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int32(s.TotalSuccess),
		TotalAmount:  int32(s.TotalAmount),
	}
}

func (t *transferProtoMapper) mapResponsesTransferMonthStatusSuccess(Transfers []*response.TransferResponseMonthStatusSuccess) []*pb.TransferMonthStatusSuccessResponse {
	var TransferRecords []*pb.TransferMonthStatusSuccessResponse

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.mapResponseTransferMonthStatusSuccess(Transfer))
	}

	return TransferRecords
}

func (t *transferProtoMapper) mapTransferResponseYearStatusSuccess(s *response.TransferResponseYearStatusSuccess) *pb.TransferYearStatusSuccessResponse {
	return &pb.TransferYearStatusSuccessResponse{
		Year:         s.Year,
		TotalSuccess: int32(s.TotalSuccess),
		TotalAmount:  int32(s.TotalAmount),
	}
}

func (t *transferProtoMapper) mapTransferResponsesYearStatusSuccess(Transfers []*response.TransferResponseYearStatusSuccess) []*pb.TransferYearStatusSuccessResponse {
	var TransferRecords []*pb.TransferYearStatusSuccessResponse

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.mapTransferResponseYearStatusSuccess(Transfer))
	}

	return TransferRecords
}

func (t *transferProtoMapper) mapResponseTransferMonthStatusFailed(s *response.TransferResponseMonthStatusFailed) *pb.TransferMonthStatusFailedResponse {
	return &pb.TransferMonthStatusFailedResponse{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int32(s.TotalFailed),
		TotalAmount: int32(s.TotalAmount),
	}
}

func (t *transferProtoMapper) mapResponsesTransferMonthStatusFailed(Transfers []*response.TransferResponseMonthStatusFailed) []*pb.TransferMonthStatusFailedResponse {
	var TransferRecords []*pb.TransferMonthStatusFailedResponse

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.mapResponseTransferMonthStatusFailed(Transfer))
	}

	return TransferRecords
}

func (t *transferProtoMapper) mapTransferResponseYearStatusFailed(s *response.TransferResponseYearStatusFailed) *pb.TransferYearStatusFailedResponse {
	return &pb.TransferYearStatusFailedResponse{
		Year:        s.Year,
		TotalFailed: int32(s.TotalFailed),
		TotalAmount: int32(s.TotalAmount),
	}
}

func (t *transferProtoMapper) mapTransferResponsesYearStatusFailed(Transfers []*response.TransferResponseYearStatusFailed) []*pb.TransferYearStatusFailedResponse {
	var TransferRecords []*pb.TransferYearStatusFailedResponse

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.mapTransferResponseYearStatusFailed(Transfer))
	}

	return TransferRecords
}

func (m *transferProtoMapper) mapResponseTransferMonthAmount(s *response.TransferMonthAmountResponse) *pb.TransferMonthAmountResponse {
	return &pb.TransferMonthAmountResponse{
		Month:       s.Month,
		TotalAmount: int32(s.TotalAmount),
	}
}

func (m *transferProtoMapper) mapResponseTransferMonthAmounts(s []*response.TransferMonthAmountResponse) []*pb.TransferMonthAmountResponse {
	var protoResponses []*pb.TransferMonthAmountResponse
	for _, transfer := range s {
		protoResponses = append(protoResponses, m.mapResponseTransferMonthAmount(transfer))
	}
	return protoResponses
}

func (m *transferProtoMapper) mapResponseTransferYearAmount(s *response.TransferYearAmountResponse) *pb.TransferYearAmountResponse {
	return &pb.TransferYearAmountResponse{
		Year:        s.Year,
		TotalAmount: int32(s.TotalAmount),
	}
}

func (m *transferProtoMapper) mapResponseTransferYearAmounts(s []*response.TransferYearAmountResponse) []*pb.TransferYearAmountResponse {
	var protoResponses []*pb.TransferYearAmountResponse
	for _, transfer := range s {
		protoResponses = append(protoResponses, m.mapResponseTransferYearAmount(transfer))
	}
	return protoResponses
}
