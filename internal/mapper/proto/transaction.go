package protomapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type transactionProtoMapper struct{}

func NewTransactionProtoMapper() *transactionProtoMapper {
	return &transactionProtoMapper{}
}

func (m *transactionProtoMapper) ToProtoResponseTransactionMonthStatusSuccess(status string, message string, pbResponse []*response.TransactionResponseMonthStatusSuccess) *pb.ApiResponseTransactionMonthStatusSuccess {
	return &pb.ApiResponseTransactionMonthStatusSuccess{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesTransactionMonthStatusSuccess(pbResponse),
	}
}

func (m *transactionProtoMapper) ToProtoResponseTransactionYearStatusSuccess(status string, message string, pbResponse []*response.TransactionResponseYearStatusSuccess) *pb.ApiResponseTransactionYearStatusSuccess {
	return &pb.ApiResponseTransactionYearStatusSuccess{
		Status:  status,
		Message: message,
		Data:    m.mapTransactionResponsesYearStatusSuccess(pbResponse),
	}
}

func (m *transactionProtoMapper) ToProtoResponseTransactionMonthStatusFailed(status string, message string, pbResponse []*response.TransactionResponseMonthStatusFailed) *pb.ApiResponseTransactionMonthStatusFailed {
	return &pb.ApiResponseTransactionMonthStatusFailed{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesTransactionMonthStatusFailed(pbResponse),
	}
}

func (m *transactionProtoMapper) ToProtoResponseTransactionYearStatusFailed(status string, message string, pbResponse []*response.TransactionResponseYearStatusFailed) *pb.ApiResponseTransactionYearStatusFailed {
	return &pb.ApiResponseTransactionYearStatusFailed{
		Status:  status,
		Message: message,
		Data:    m.mapTransactionResponsesYearStatusFailed(pbResponse),
	}
}

func (m *transactionProtoMapper) ToProtoResponseTransactionMonthMethod(status string, message string, pbResponse []*response.TransactionMonthMethodResponse) *pb.ApiResponseTransactionMonthMethod {
	return &pb.ApiResponseTransactionMonthMethod{
		Status:  status,
		Message: message,
		Data:    m.mapResponseTransactionMonthMethods(pbResponse),
	}
}

func (m *transactionProtoMapper) ToProtoResponseTransactionYearMethod(status string, message string, pbResponse []*response.TransactionYearMethodResponse) *pb.ApiResponseTransactionYearMethod {

	return &pb.ApiResponseTransactionYearMethod{
		Status:  status,
		Message: message,
		Data:    m.mapResponseTransactionYearMethods(pbResponse),
	}
}

func (m *transactionProtoMapper) ToProtoResponseTransactionMonthAmount(status string, message string, pbResponse []*response.TransactionMonthAmountResponse) *pb.ApiResponseTransactionMonthAmount {
	return &pb.ApiResponseTransactionMonthAmount{
		Status:  status,
		Message: message,
		Data:    m.mapResponseTransactionMonthAmounts(pbResponse),
	}
}

func (m *transactionProtoMapper) ToProtoResponseTransactionYearAmount(status string, message string, pbResponse []*response.TransactionYearlyAmountResponse) *pb.ApiResponseTransactionYearAmount {
	return &pb.ApiResponseTransactionYearAmount{
		Status:  status,
		Message: message,
		Data:    m.mapResponseTransactionYearlyAmounts(pbResponse),
	}
}

func (m *transactionProtoMapper) ToProtoResponseTransaction(status string, message string, pbResponse *response.TransactionResponse) *pb.ApiResponseTransaction {
	return &pb.ApiResponseTransaction{
		Status:  status,
		Message: message,
		Data:    m.mapResponseTransaction(pbResponse),
	}
}

func (m *transactionProtoMapper) ToProtoResponseTransactions(status string, message string, pbResponse []*response.TransactionResponse) *pb.ApiResponseTransactions {
	return &pb.ApiResponseTransactions{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesTransaction(pbResponse),
	}
}

func (m *transactionProtoMapper) ToProtoResponseTransactionDelete(status string, message string) *pb.ApiResponseTransactionDelete {
	return &pb.ApiResponseTransactionDelete{
		Status:  status,
		Message: message,
	}
}

func (m *transactionProtoMapper) ToProtoResponseTransactionAll(status string, message string) *pb.ApiResponseTransactionAll {
	return &pb.ApiResponseTransactionAll{
		Status:  status,
		Message: message,
	}
}

func (m *transactionProtoMapper) ToProtoResponsePaginationTransaction(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.TransactionResponse) *pb.ApiResponsePaginationTransaction {

	return &pb.ApiResponsePaginationTransaction{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesTransaction(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *transactionProtoMapper) ToProtoResponsePaginationTransactionDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.TransactionResponseDeleteAt) *pb.ApiResponsePaginationTransactionDeleteAt {

	return &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesTransactionDeleteAt(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *transactionProtoMapper) mapResponseTransaction(transaction *response.TransactionResponse) *pb.TransactionResponse {
	return &pb.TransactionResponse{
		Id:              int32(transaction.ID),
		TransactionNo:   transaction.TransactionNo,
		CardNumber:      transaction.CardNumber,
		Amount:          int32(transaction.Amount),
		PaymentMethod:   transaction.PaymentMethod,
		TransactionTime: transaction.TransactionTime,
		MerchantId:      int32(transaction.MerchantID),
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}

}

func (m *transactionProtoMapper) mapResponsesTransaction(transactions []*response.TransactionResponse) []*pb.TransactionResponse {
	var result []*pb.TransactionResponse
	for _, transaction := range transactions {
		result = append(result, m.mapResponseTransaction(transaction))
	}
	return result
}

func (m *transactionProtoMapper) mapResponseTransactionDeleteAt(transaction *response.TransactionResponseDeleteAt) *pb.TransactionResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if transaction.DeletedAt != nil {
		deletedAt = wrapperspb.String(*transaction.DeletedAt)
	}

	return &pb.TransactionResponseDeleteAt{
		Id:              int32(transaction.ID),
		TransactionNo:   transaction.TransactionNo,
		CardNumber:      transaction.CardNumber,
		Amount:          int32(transaction.Amount),
		PaymentMethod:   transaction.PaymentMethod,
		TransactionTime: transaction.TransactionTime,
		MerchantId:      int32(transaction.MerchantID),
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
		DeletedAt:       deletedAt,
	}

}

func (m *transactionProtoMapper) mapResponsesTransactionDeleteAt(transactions []*response.TransactionResponseDeleteAt) []*pb.TransactionResponseDeleteAt {
	var result []*pb.TransactionResponseDeleteAt

	for _, transaction := range transactions {
		result = append(result, m.mapResponseTransactionDeleteAt(transaction))
	}
	return result
}

func (t *transactionProtoMapper) mapResponseTransactionMonthStatusSuccess(s *response.TransactionResponseMonthStatusSuccess) *pb.TransactionMonthStatusSuccessResponse {
	return &pb.TransactionMonthStatusSuccessResponse{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int32(s.TotalSuccess),
		TotalAmount:  int32(s.TotalAmount),
	}
}

func (t *transactionProtoMapper) mapResponsesTransactionMonthStatusSuccess(Transactions []*response.TransactionResponseMonthStatusSuccess) []*pb.TransactionMonthStatusSuccessResponse {
	var TransactionRecords []*pb.TransactionMonthStatusSuccessResponse

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.mapResponseTransactionMonthStatusSuccess(Transaction))
	}

	return TransactionRecords
}

func (t *transactionProtoMapper) mapTransactionResponseYearStatusSuccess(s *response.TransactionResponseYearStatusSuccess) *pb.TransactionYearStatusSuccessResponse {
	return &pb.TransactionYearStatusSuccessResponse{
		Year:         s.Year,
		TotalSuccess: int32(s.TotalSuccess),
		TotalAmount:  int32(s.TotalAmount),
	}
}

func (t *transactionProtoMapper) mapTransactionResponsesYearStatusSuccess(Transactions []*response.TransactionResponseYearStatusSuccess) []*pb.TransactionYearStatusSuccessResponse {
	var TransactionRecords []*pb.TransactionYearStatusSuccessResponse

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.mapTransactionResponseYearStatusSuccess(Transaction))
	}

	return TransactionRecords
}

func (t *transactionProtoMapper) mapResponseTransactionMonthStatusFailed(s *response.TransactionResponseMonthStatusFailed) *pb.TransactionMonthStatusFailedResponse {
	return &pb.TransactionMonthStatusFailedResponse{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int32(s.TotalFailed),
		TotalAmount: int32(s.TotalAmount),
	}
}

func (t *transactionProtoMapper) mapResponsesTransactionMonthStatusFailed(Transactions []*response.TransactionResponseMonthStatusFailed) []*pb.TransactionMonthStatusFailedResponse {
	var TransactionRecords []*pb.TransactionMonthStatusFailedResponse

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.mapResponseTransactionMonthStatusFailed(Transaction))
	}

	return TransactionRecords
}

func (t *transactionProtoMapper) mapTransactionResponseYearStatusFailed(s *response.TransactionResponseYearStatusFailed) *pb.TransactionYearStatusFailedResponse {
	return &pb.TransactionYearStatusFailedResponse{
		Year:        s.Year,
		TotalFailed: int32(s.TotalFailed),
		TotalAmount: int32(s.TotalAmount),
	}
}

func (t *transactionProtoMapper) mapTransactionResponsesYearStatusFailed(Transactions []*response.TransactionResponseYearStatusFailed) []*pb.TransactionYearStatusFailedResponse {
	var TransactionRecords []*pb.TransactionYearStatusFailedResponse

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.mapTransactionResponseYearStatusFailed(Transaction))
	}

	return TransactionRecords
}

func (m *transactionProtoMapper) mapResponseTransactionMonthMethod(s *response.TransactionMonthMethodResponse) *pb.TransactionMonthMethodResponse {
	return &pb.TransactionMonthMethodResponse{
		Month:             s.Month,
		PaymentMethod:     s.PaymentMethod,
		TotalTransactions: int32(s.TotalTransactions),
		TotalAmount:       int32(s.TotalAmount),
	}
}

func (m *transactionProtoMapper) mapResponseTransactionMonthMethods(s []*response.TransactionMonthMethodResponse) []*pb.TransactionMonthMethodResponse {
	var protoResponses []*pb.TransactionMonthMethodResponse
	for _, transaction := range s {
		protoResponses = append(protoResponses, m.mapResponseTransactionMonthMethod(transaction))
	}
	return protoResponses
}

func (m *transactionProtoMapper) mapResponseTransactionYearMethod(s *response.TransactionYearMethodResponse) *pb.TransactionYearMethodResponse {
	return &pb.TransactionYearMethodResponse{
		Year:              s.Year,
		PaymentMethod:     s.PaymentMethod,
		TotalTransactions: int32(s.TotalTransactions),
		TotalAmount:       int32(s.TotalAmount),
	}
}

func (m *transactionProtoMapper) mapResponseTransactionYearMethods(s []*response.TransactionYearMethodResponse) []*pb.TransactionYearMethodResponse {
	var protoResponses []*pb.TransactionYearMethodResponse
	for _, transaction := range s {
		protoResponses = append(protoResponses, m.mapResponseTransactionYearMethod(transaction))
	}
	return protoResponses
}

func (m *transactionProtoMapper) mapResponseTransactionMonthAmount(s *response.TransactionMonthAmountResponse) *pb.TransactionMonthAmountResponse {
	return &pb.TransactionMonthAmountResponse{
		Month:       s.Month,
		TotalAmount: int32(s.TotalAmount),
	}
}

func (m *transactionProtoMapper) mapResponseTransactionMonthAmounts(s []*response.TransactionMonthAmountResponse) []*pb.TransactionMonthAmountResponse {
	var protoResponses []*pb.TransactionMonthAmountResponse
	for _, transaction := range s {
		protoResponses = append(protoResponses, m.mapResponseTransactionMonthAmount(transaction))
	}
	return protoResponses
}

func (m *transactionProtoMapper) mapResponseTransactionYearlyAmount(s *response.TransactionYearlyAmountResponse) *pb.TransactionYearlyAmountResponse {
	return &pb.TransactionYearlyAmountResponse{
		Year:        s.Year,
		TotalAmount: int32(s.TotalAmount),
	}
}

func (m *transactionProtoMapper) mapResponseTransactionYearlyAmounts(s []*response.TransactionYearlyAmountResponse) []*pb.TransactionYearlyAmountResponse {
	var protoResponses []*pb.TransactionYearlyAmountResponse
	for _, transaction := range s {
		protoResponses = append(protoResponses, m.mapResponseTransactionYearlyAmount(transaction))
	}
	return protoResponses
}
