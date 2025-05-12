package apimapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"
)

type transferResponseMapper struct{}

func NewTransferResponseMapper() *transferResponseMapper {
	return &transferResponseMapper{}
}

func (m *transferResponseMapper) ToApiResponseTransferMonthStatusSuccess(pbResponse *pb.ApiResponseTransferMonthStatusSuccess) *response.ApiResponseTransferMonthStatusSuccess {
	return &response.ApiResponseTransferMonthStatusSuccess{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponsesTransferMonthStatusSuccess(pbResponse.Data),
	}
}

func (m *transferResponseMapper) ToApiResponseTransferYearStatusSuccess(pbResponse *pb.ApiResponseTransferYearStatusSuccess) *response.ApiResponseTransferYearStatusSuccess {
	return &response.ApiResponseTransferYearStatusSuccess{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapTransferResponsesYearStatusSuccess(pbResponse.Data),
	}
}

func (m *transferResponseMapper) ToApiResponseTransferMonthStatusFailed(pbResponse *pb.ApiResponseTransferMonthStatusFailed) *response.ApiResponseTransferMonthStatusFailed {
	return &response.ApiResponseTransferMonthStatusFailed{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponsesTransferMonthStatusFailed(pbResponse.Data),
	}
}

func (m *transferResponseMapper) ToApiResponseTransferYearStatusFailed(pbResponse *pb.ApiResponseTransferYearStatusFailed) *response.ApiResponseTransferYearStatusFailed {
	return &response.ApiResponseTransferYearStatusFailed{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapTransferResponsesYearStatusFailed(pbResponse.Data),
	}
}

func (m *transferResponseMapper) ToApiResponseTransferMonthAmount(pbResponse *pb.ApiResponseTransferMonthAmount) *response.ApiResponseTransferMonthAmount {
	return &response.ApiResponseTransferMonthAmount{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponseTransferMonthAmounts(pbResponse.Data),
	}
}

func (m *transferResponseMapper) ToApiResponseTransferYearAmount(pbResponse *pb.ApiResponseTransferYearAmount) *response.ApiResponseTransferYearAmount {
	return &response.ApiResponseTransferYearAmount{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponseTransferYearAmounts(pbResponse.Data),
	}
}

func (m *transferResponseMapper) ToApiResponseTransfer(pbResponse *pb.ApiResponseTransfer) *response.ApiResponseTransfer {
	return &response.ApiResponseTransfer{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponseTransfer(pbResponse.Data),
	}
}

func (m *transferResponseMapper) ToApiResponseTransfers(pbResponse *pb.ApiResponseTransfers) *response.ApiResponseTransfers {
	return &response.ApiResponseTransfers{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponsesTransfer(pbResponse.Data),
	}
}

func (m *transferResponseMapper) ToApiResponseTransferDelete(pbResponse *pb.ApiResponseTransferDelete) *response.ApiResponseTransferDelete {
	return &response.ApiResponseTransferDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (m *transferResponseMapper) ToApiResponseTransferAll(pbResponse *pb.ApiResponseTransferAll) *response.ApiResponseTransferAll {
	return &response.ApiResponseTransferAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (m *transferResponseMapper) ToApiResponsePaginationTransfer(pbResponse *pb.ApiResponsePaginationTransfer) *response.ApiResponsePaginationTransfer {
	return &response.ApiResponsePaginationTransfer{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.mapResponsesTransfer(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *transferResponseMapper) ToApiResponsePaginationTransferDeleteAt(pbResponse *pb.ApiResponsePaginationTransferDeleteAt) *response.ApiResponsePaginationTransferDeleteAt {
	return &response.ApiResponsePaginationTransferDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.mapResponsesTransferDeleteAt(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.Pagination),
	}
}

func (t *transferResponseMapper) mapResponseTransfer(transfer *pb.TransferResponse) *response.TransferResponse {
	return &response.TransferResponse{
		ID:             int(transfer.Id),
		TransferNo:     transfer.TransferNo,
		TransferFrom:   transfer.TransferFrom,
		TransferTo:     transfer.TransferTo,
		TransferAmount: int(transfer.TransferAmount),
		TransferTime:   transfer.TransferTime,
		CreatedAt:      transfer.CreatedAt,
		UpdatedAt:      transfer.UpdatedAt,
	}
}

func (t *transferResponseMapper) mapResponsesTransfer(transfers []*pb.TransferResponse) []*response.TransferResponse {
	var responses []*response.TransferResponse

	for _, response := range transfers {
		responses = append(responses, t.mapResponseTransfer(response))
	}

	return responses
}

func (t *transferResponseMapper) mapResponseTransferDeleteAt(transfer *pb.TransferResponseDeleteAt) *response.TransferResponseDeleteAt {
	var deletedAt string
	if transfer.DeletedAt != nil {
		deletedAt = transfer.DeletedAt.Value
	}

	return &response.TransferResponseDeleteAt{
		ID:             int(transfer.Id),
		TransferNo:     transfer.TransferNo,
		TransferFrom:   transfer.TransferFrom,
		TransferTo:     transfer.TransferTo,
		TransferAmount: int(transfer.TransferAmount),
		TransferTime:   transfer.TransferTime,
		CreatedAt:      transfer.CreatedAt,
		UpdatedAt:      transfer.UpdatedAt,
		DeletedAt:      &deletedAt,
	}
}

func (t *transferResponseMapper) mapResponsesTransferDeleteAt(transfers []*pb.TransferResponseDeleteAt) []*response.TransferResponseDeleteAt {
	var responses []*response.TransferResponseDeleteAt

	for _, response := range transfers {
		responses = append(responses, t.mapResponseTransferDeleteAt(response))
	}

	return responses
}

func (t *transferResponseMapper) mapResponseTransferMonthStatusSuccess(s *pb.TransferMonthStatusSuccessResponse) *response.TransferResponseMonthStatusSuccess {
	return &response.TransferResponseMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *transferResponseMapper) mapResponsesTransferMonthStatusSuccess(Transfers []*pb.TransferMonthStatusSuccessResponse) []*response.TransferResponseMonthStatusSuccess {
	var TransferRecords []*response.TransferResponseMonthStatusSuccess

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.mapResponseTransferMonthStatusSuccess(Transfer))
	}

	return TransferRecords
}

func (t *transferResponseMapper) mapTransferResponseYearStatusSuccess(s *pb.TransferYearStatusSuccessResponse) *response.TransferResponseYearStatusSuccess {
	return &response.TransferResponseYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *transferResponseMapper) mapTransferResponsesYearStatusSuccess(Transfers []*pb.TransferYearStatusSuccessResponse) []*response.TransferResponseYearStatusSuccess {
	var TransferRecords []*response.TransferResponseYearStatusSuccess

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.mapTransferResponseYearStatusSuccess(Transfer))
	}

	return TransferRecords
}

func (t *transferResponseMapper) mapResponseTransferMonthStatusFailed(s *pb.TransferMonthStatusFailedResponse) *response.TransferResponseMonthStatusFailed {
	return &response.TransferResponseMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transferResponseMapper) mapResponsesTransferMonthStatusFailed(Transfers []*pb.TransferMonthStatusFailedResponse) []*response.TransferResponseMonthStatusFailed {
	var TransferRecords []*response.TransferResponseMonthStatusFailed

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.mapResponseTransferMonthStatusFailed(Transfer))
	}

	return TransferRecords
}

func (t *transferResponseMapper) mapTransferResponseYearStatusFailed(s *pb.TransferYearStatusFailedResponse) *response.TransferResponseYearStatusFailed {
	return &response.TransferResponseYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transferResponseMapper) mapTransferResponsesYearStatusFailed(Transfers []*pb.TransferYearStatusFailedResponse) []*response.TransferResponseYearStatusFailed {
	var TransferRecords []*response.TransferResponseYearStatusFailed

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.mapTransferResponseYearStatusFailed(Transfer))
	}

	return TransferRecords
}

func (m *transferResponseMapper) mapResponseTransferMonthAmount(s *pb.TransferMonthAmountResponse) *response.TransferMonthAmountResponse {
	return &response.TransferMonthAmountResponse{
		Month:       s.Month,
		TotalAmount: int(s.TotalAmount),
	}
}

func (m *transferResponseMapper) mapResponseTransferMonthAmounts(s []*pb.TransferMonthAmountResponse) []*response.TransferMonthAmountResponse {
	var protoResponses []*response.TransferMonthAmountResponse
	for _, transfer := range s {
		protoResponses = append(protoResponses, m.mapResponseTransferMonthAmount(transfer))
	}
	return protoResponses
}

func (m *transferResponseMapper) mapResponseTransferYearAmount(s *pb.TransferYearAmountResponse) *response.TransferYearAmountResponse {
	return &response.TransferYearAmountResponse{
		Year:        s.Year,
		TotalAmount: int(s.TotalAmount),
	}
}

func (m *transferResponseMapper) mapResponseTransferYearAmounts(s []*pb.TransferYearAmountResponse) []*response.TransferYearAmountResponse {
	var protoResponses []*response.TransferYearAmountResponse
	for _, transfer := range s {
		protoResponses = append(protoResponses, m.mapResponseTransferYearAmount(transfer))
	}
	return protoResponses
}
