package responseservice

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
)

type transferResponseMapper struct {
}

func NewTransferResponseMapper() *transferResponseMapper {
	return &transferResponseMapper{}
}

func (s *transferResponseMapper) ToTransferResponse(transfer *record.TransferRecord) *response.TransferResponse {
	return &response.TransferResponse{
		ID:             transfer.ID,
		TransferNo:     transfer.TransferNo,
		TransferFrom:   transfer.TransferFrom,
		TransferTo:     transfer.TransferTo,
		TransferAmount: transfer.TransferAmount,
		TransferTime:   transfer.TransferTime,
		CreatedAt:      transfer.CreatedAt,
		UpdatedAt:      transfer.UpdatedAt,
	}
}

func (s *transferResponseMapper) ToTransfersResponse(transfers []*record.TransferRecord) []*response.TransferResponse {
	var responses []*response.TransferResponse

	for _, response := range transfers {
		responses = append(responses, s.ToTransferResponse(response))
	}

	return responses
}

func (s *transferResponseMapper) ToTransferResponseDeleteAt(transfer *record.TransferRecord) *response.TransferResponseDeleteAt {
	return &response.TransferResponseDeleteAt{
		ID:             transfer.ID,
		TransferNo:     transfer.TransferNo,
		TransferFrom:   transfer.TransferFrom,
		TransferTo:     transfer.TransferTo,
		TransferAmount: transfer.TransferAmount,
		TransferTime:   transfer.TransferTime,
		CreatedAt:      transfer.CreatedAt,
		UpdatedAt:      transfer.UpdatedAt,
		DeletedAt:      transfer.DeletedAt,
	}
}

func (s *transferResponseMapper) ToTransfersResponseDeleteAt(transfers []*record.TransferRecord) []*response.TransferResponseDeleteAt {
	var responses []*response.TransferResponseDeleteAt

	for _, response := range transfers {
		responses = append(responses, s.ToTransferResponseDeleteAt(response))
	}

	return responses
}

func (t *transferResponseMapper) ToTransferResponseMonthStatusSuccess(s *record.TransferRecordMonthStatusSuccess) *response.TransferResponseMonthStatusSuccess {
	return &response.TransferResponseMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  s.TotalAmount,
	}
}

func (t *transferResponseMapper) ToTransferResponsesMonthStatusSuccess(Transfers []*record.TransferRecordMonthStatusSuccess) []*response.TransferResponseMonthStatusSuccess {
	var TransferRecords []*response.TransferResponseMonthStatusSuccess

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.ToTransferResponseMonthStatusSuccess(Transfer))
	}

	return TransferRecords
}

func (t *transferResponseMapper) ToTransferResponseYearStatusSuccess(s *record.TransferRecordYearStatusSuccess) *response.TransferResponseYearStatusSuccess {
	return &response.TransferResponseYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  s.TotalAmount,
	}
}

func (t *transferResponseMapper) ToTransferResponsesYearStatusSuccess(Transfers []*record.TransferRecordYearStatusSuccess) []*response.TransferResponseYearStatusSuccess {
	var TransferRecords []*response.TransferResponseYearStatusSuccess

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.ToTransferResponseYearStatusSuccess(Transfer))
	}

	return TransferRecords
}

func (t *transferResponseMapper) ToTransferResponseMonthStatusFailed(s *record.TransferRecordMonthStatusFailed) *response.TransferResponseMonthStatusFailed {
	return &response.TransferResponseMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: s.TotalAmount,
	}
}

func (t *transferResponseMapper) ToTransferResponsesMonthStatusFailed(Transfers []*record.TransferRecordMonthStatusFailed) []*response.TransferResponseMonthStatusFailed {
	var TransferRecords []*response.TransferResponseMonthStatusFailed

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.ToTransferResponseMonthStatusFailed(Transfer))
	}

	return TransferRecords
}

func (t *transferResponseMapper) ToTransferResponseYearStatusFailed(s *record.TransferRecordYearStatusFailed) *response.TransferResponseYearStatusFailed {
	return &response.TransferResponseYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: s.TotalAmount,
	}
}

func (t *transferResponseMapper) ToTransferResponsesYearStatusFailed(Transfers []*record.TransferRecordYearStatusFailed) []*response.TransferResponseYearStatusFailed {
	var TransferRecords []*response.TransferResponseYearStatusFailed

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.ToTransferResponseYearStatusFailed(Transfer))
	}

	return TransferRecords
}

func (t *transferResponseMapper) ToTransferResponseMonthAmount(s *record.TransferMonthAmount) *response.TransferMonthAmountResponse {
	return &response.TransferMonthAmountResponse{
		Month:       s.Month,
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transferResponseMapper) ToTransferResponsesMonthAmount(s []*record.TransferMonthAmount) []*response.TransferMonthAmountResponse {
	var transferResponses []*response.TransferMonthAmountResponse
	for _, transfer := range s {
		transferResponses = append(transferResponses, t.ToTransferResponseMonthAmount(transfer))
	}
	return transferResponses
}

func (t *transferResponseMapper) ToTransferResponseYearAmount(s *record.TransferYearAmount) *response.TransferYearAmountResponse {
	return &response.TransferYearAmountResponse{
		Year:        s.Year,
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transferResponseMapper) ToTransferResponsesYearAmount(s []*record.TransferYearAmount) []*response.TransferYearAmountResponse {
	var transferResponses []*response.TransferYearAmountResponse
	for _, transfer := range s {
		transferResponses = append(transferResponses, t.ToTransferResponseYearAmount(transfer))
	}
	return transferResponses
}
