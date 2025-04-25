package apimapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"
)

type transactionResponseMapper struct {
}

func NewTransactionResponseMapper() *transactionResponseMapper {
	return &transactionResponseMapper{}
}

func (m *transactionResponseMapper) ToApiResponseTransactionMonthStatusSuccess(pbResponse *pb.ApiResponseTransactionMonthStatusSuccess) *response.ApiResponseTransactionMonthStatusSuccess {

	return &response.ApiResponseTransactionMonthStatusSuccess{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesTransactionMonthStatusSuccess(pbResponse.Data),
	}
}

func (m *transactionResponseMapper) ToApiResponseTransactionYearStatusSuccess(pbResponse *pb.ApiResponseTransactionYearStatusSuccess) *response.ApiResponseTransactionYearStatusSuccess {

	return &response.ApiResponseTransactionYearStatusSuccess{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapTransactionResponsesYearStatusSuccess(pbResponse.Data),
	}
}

func (m *transactionResponseMapper) ToApiResponseTransactionMonthStatusFailed(pbResponse *pb.ApiResponseTransactionMonthStatusFailed) *response.ApiResponseTransactionMonthStatusFailed {

	return &response.ApiResponseTransactionMonthStatusFailed{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponsesTransactionMonthStatusFailed(pbResponse.Data),
	}
}

func (m *transactionResponseMapper) ToApiResponseTransactionYearStatusFailed(pbResponse *pb.ApiResponseTransactionYearStatusFailed) *response.ApiResponseTransactionYearStatusFailed {

	return &response.ApiResponseTransactionYearStatusFailed{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapTransactionResponsesYearStatusFailed(pbResponse.Data),
	}
}

func (m *transactionResponseMapper) ToApiResponseTransactionMonthMethod(pbResponse *pb.ApiResponseTransactionMonthMethod) *response.ApiResponseTransactionMonthMethod {

	return &response.ApiResponseTransactionMonthMethod{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponseTransactionMonthMethods(pbResponse.Data),
	}
}

func (m *transactionResponseMapper) ToApiResponseTransactionYearMethod(pbResponse *pb.ApiResponseTransactionYearMethod) *response.ApiResponseTransactionYearMethod {

	return &response.ApiResponseTransactionYearMethod{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponseTransactionYearMethods(pbResponse.Data),
	}
}

func (m *transactionResponseMapper) ToApiResponseTransactionMonthAmount(pbResponse *pb.ApiResponseTransactionMonthAmount) *response.ApiResponseTransactionMonthAmount {

	return &response.ApiResponseTransactionMonthAmount{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponseTransactionMonthAmounts(pbResponse.Data),
	}
}

func (m *transactionResponseMapper) ToApiResponseTransactionYearAmount(pbResponse *pb.ApiResponseTransactionYearAmount) *response.ApiResponseTransactionYearAmount {

	return &response.ApiResponseTransactionYearAmount{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponseTransactionYearlyAmounts(pbResponse.Data),
	}
}

func (m *transactionResponseMapper) ToApiResponseTransaction(pbResponse *pb.ApiResponseTransaction) *response.ApiResponseTransaction {
	return &response.ApiResponseTransaction{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponseTransaction(pbResponse.Data),
	}
}

func (m *transactionResponseMapper) ToApiResponseTransactions(pbResponse *pb.ApiResponseTransactions) *response.ApiResponseTransactions {

	return &response.ApiResponseTransactions{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.mapResponsesTransaction(pbResponse.Data),
	}
}

func (m *transactionResponseMapper) ToApiResponseTransactionDelete(pbResponse *pb.ApiResponseTransactionDelete) *response.ApiResponseTransactionDelete {
	return &response.ApiResponseTransactionDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (m *transactionResponseMapper) ToApiResponseTransactionAll(pbResponse *pb.ApiResponseTransactionAll) *response.ApiResponseTransactionAll {
	return &response.ApiResponseTransactionAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (m *transactionResponseMapper) ToApiResponsePaginationTransaction(pbResponse *pb.ApiResponsePaginationTransaction) *response.ApiResponsePaginationTransaction {

	return &response.ApiResponsePaginationTransaction{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.mapResponsesTransaction(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *transactionResponseMapper) ToApiResponsePaginationTransactionDeleteAt(pbResponse *pb.ApiResponsePaginationTransactionDeleteAt) *response.ApiResponsePaginationTransactionDeleteAt {

	return &response.ApiResponsePaginationTransactionDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesTransactionDeleteAt(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *transactionResponseMapper) mapResponseTransaction(transaction *pb.TransactionResponse) *response.TransactionResponse {
	return &response.TransactionResponse{
		ID:              int(transaction.Id),
		TransactionNo:   transaction.TransactionNo,
		CardNumber:      transaction.CardNumber,
		Amount:          int(transaction.Amount),
		PaymentMethod:   transaction.PaymentMethod,
		TransactionTime: transaction.TransactionTime,
		MerchantID:      int(transaction.MerchantId),
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}
}

func (m *transactionResponseMapper) mapResponsesTransaction(transactions []*pb.TransactionResponse) []*response.TransactionResponse {
	var result []*response.TransactionResponse
	for _, transaction := range transactions {
		result = append(result, m.mapResponseTransaction(transaction))
	}
	return result
}

func (m *transactionResponseMapper) mapResponseTransactionDeleteAt(transaction *pb.TransactionResponseDeleteAt) *response.TransactionResponseDeleteAt {
	var deletedAt string
	if transaction.DeletedAt != nil {
		deletedAt = transaction.DeletedAt.Value
	}

	return &response.TransactionResponseDeleteAt{
		ID:              int(transaction.Id),
		TransactionNo:   transaction.TransactionNo,
		CardNumber:      transaction.CardNumber,
		Amount:          int(transaction.Amount),
		PaymentMethod:   transaction.PaymentMethod,
		TransactionTime: transaction.TransactionTime,
		MerchantID:      int(transaction.MerchantId),
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
		DeletedAt:       &deletedAt,
	}
}

func (m *transactionResponseMapper) ToResponsesTransactionDeleteAt(transactions []*pb.TransactionResponseDeleteAt) []*response.TransactionResponseDeleteAt {
	var result []*response.TransactionResponseDeleteAt
	for _, transaction := range transactions {
		result = append(result, m.mapResponseTransactionDeleteAt(transaction))
	}
	return result
}

func (m *transactionResponseMapper) mapResponseTransactionMonthStatusSuccess(s *pb.TransactionMonthStatusSuccessResponse) *response.TransactionResponseMonthStatusSuccess {
	return &response.TransactionResponseMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (m *transactionResponseMapper) ToResponsesTransactionMonthStatusSuccess(transactions []*pb.TransactionMonthStatusSuccessResponse) []*response.TransactionResponseMonthStatusSuccess {
	var transactionRecords []*response.TransactionResponseMonthStatusSuccess
	for _, transaction := range transactions {
		transactionRecords = append(transactionRecords, m.mapResponseTransactionMonthStatusSuccess(transaction))
	}
	return transactionRecords
}

func (m *transactionResponseMapper) mapTransactionResponseYearStatusSuccess(s *pb.TransactionYearStatusSuccessResponse) *response.TransactionResponseYearStatusSuccess {
	return &response.TransactionResponseYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (m *transactionResponseMapper) mapTransactionResponsesYearStatusSuccess(transactions []*pb.TransactionYearStatusSuccessResponse) []*response.TransactionResponseYearStatusSuccess {
	var transactionRecords []*response.TransactionResponseYearStatusSuccess
	for _, transaction := range transactions {
		transactionRecords = append(transactionRecords, m.mapTransactionResponseYearStatusSuccess(transaction))
	}
	return transactionRecords
}

func (m *transactionResponseMapper) mapResponseTransactionMonthStatusFailed(s *pb.TransactionMonthStatusFailedResponse) *response.TransactionResponseMonthStatusFailed {
	return &response.TransactionResponseMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (m *transactionResponseMapper) mapResponsesTransactionMonthStatusFailed(transactions []*pb.TransactionMonthStatusFailedResponse) []*response.TransactionResponseMonthStatusFailed {
	var transactionRecords []*response.TransactionResponseMonthStatusFailed
	for _, transaction := range transactions {
		transactionRecords = append(transactionRecords, m.mapResponseTransactionMonthStatusFailed(transaction))
	}
	return transactionRecords
}

func (m *transactionResponseMapper) mapTransactionResponseYearStatusFailed(s *pb.TransactionYearStatusFailedResponse) *response.TransactionResponseYearStatusFailed {
	return &response.TransactionResponseYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (m *transactionResponseMapper) mapTransactionResponsesYearStatusFailed(transactions []*pb.TransactionYearStatusFailedResponse) []*response.TransactionResponseYearStatusFailed {
	var transactionRecords []*response.TransactionResponseYearStatusFailed
	for _, transaction := range transactions {
		transactionRecords = append(transactionRecords, m.mapTransactionResponseYearStatusFailed(transaction))
	}
	return transactionRecords
}

func (m *transactionResponseMapper) mapResponseTransactionMonthMethod(s *pb.TransactionMonthMethodResponse) *response.TransactionMonthMethodResponse {
	return &response.TransactionMonthMethodResponse{
		Month:             s.Month,
		PaymentMethod:     s.PaymentMethod,
		TotalTransactions: int(s.TotalTransactions),
		TotalAmount:       int(s.TotalAmount),
	}
}

func (m *transactionResponseMapper) mapResponseTransactionMonthMethods(s []*pb.TransactionMonthMethodResponse) []*response.TransactionMonthMethodResponse {
	var responses []*response.TransactionMonthMethodResponse
	for _, transaction := range s {
		responses = append(responses, m.mapResponseTransactionMonthMethod(transaction))
	}
	return responses
}

func (m *transactionResponseMapper) mapResponseTransactionYearMethod(s *pb.TransactionYearMethodResponse) *response.TransactionYearMethodResponse {
	return &response.TransactionYearMethodResponse{
		Year:              s.Year,
		PaymentMethod:     s.PaymentMethod,
		TotalTransactions: int(s.TotalTransactions),
		TotalAmount:       int(s.TotalAmount),
	}
}

func (m *transactionResponseMapper) mapResponseTransactionYearMethods(s []*pb.TransactionYearMethodResponse) []*response.TransactionYearMethodResponse {
	var responses []*response.TransactionYearMethodResponse
	for _, transaction := range s {
		responses = append(responses, m.mapResponseTransactionYearMethod(transaction))
	}
	return responses
}

func (m *transactionResponseMapper) mapResponseTransactionMonthAmount(s *pb.TransactionMonthAmountResponse) *response.TransactionMonthAmountResponse {
	return &response.TransactionMonthAmountResponse{
		Month:       s.Month,
		TotalAmount: int(s.TotalAmount),
	}
}

func (m *transactionResponseMapper) mapResponseTransactionMonthAmounts(s []*pb.TransactionMonthAmountResponse) []*response.TransactionMonthAmountResponse {
	var responses []*response.TransactionMonthAmountResponse
	for _, transaction := range s {
		responses = append(responses, m.mapResponseTransactionMonthAmount(transaction))
	}
	return responses
}

func (m *transactionResponseMapper) mapResponseTransactionYearlyAmount(s *pb.TransactionYearlyAmountResponse) *response.TransactionYearlyAmountResponse {
	return &response.TransactionYearlyAmountResponse{
		Year:        s.Year,
		TotalAmount: int(s.TotalAmount),
	}
}

func (m *transactionResponseMapper) mapResponseTransactionYearlyAmounts(s []*pb.TransactionYearlyAmountResponse) []*response.TransactionYearlyAmountResponse {
	var responses []*response.TransactionYearlyAmountResponse
	for _, transaction := range s {
		responses = append(responses, m.mapResponseTransactionYearlyAmount(transaction))
	}
	return responses
}
