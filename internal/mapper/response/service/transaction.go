package responseservice

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
)

type transactionResponseMapper struct {
}

func NewTransactionResponseMapper() *transactionResponseMapper {
	return &transactionResponseMapper{}
}

func (s *transactionResponseMapper) ToTransactionResponse(transaction *record.TransactionRecord) *response.TransactionResponse {
	return &response.TransactionResponse{
		ID:              transaction.ID,
		TransactionNo:   transaction.TransactionNo,
		CardNumber:      transaction.CardNumber,
		Amount:          transaction.Amount,
		PaymentMethod:   transaction.PaymentMethod,
		MerchantID:      transaction.MerchantID,
		TransactionTime: transaction.TransactionTime,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}
}

func (s *transactionResponseMapper) ToTransactionsResponse(transactions []*record.TransactionRecord) []*response.TransactionResponse {
	responses := make([]*response.TransactionResponse, 0, len(transactions))
	for _, transaction := range transactions {
		responses = append(responses, s.ToTransactionResponse(transaction))
	}
	return responses
}

func (s *transactionResponseMapper) ToTransactionResponseDeleteAt(transaction *record.TransactionRecord) *response.TransactionResponseDeleteAt {
	return &response.TransactionResponseDeleteAt{
		ID:              transaction.ID,
		TransactionNo:   transaction.TransactionNo,
		CardNumber:      transaction.CardNumber,
		Amount:          transaction.Amount,
		PaymentMethod:   transaction.PaymentMethod,
		MerchantID:      transaction.MerchantID,
		TransactionTime: transaction.TransactionTime,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
		DeletedAt:       transaction.DeletedAt,
	}
}

func (s *transactionResponseMapper) ToTransactionsResponseDeleteAt(transactions []*record.TransactionRecord) []*response.TransactionResponseDeleteAt {
	responses := make([]*response.TransactionResponseDeleteAt, 0, len(transactions))

	for _, transaction := range transactions {
		responses = append(responses, s.ToTransactionResponseDeleteAt(transaction))
	}
	return responses
}

func (t *transactionResponseMapper) ToTransactionResponseMonthStatusSuccess(s *record.TransactionRecordMonthStatusSuccess) *response.TransactionResponseMonthStatusSuccess {
	return &response.TransactionResponseMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  s.TotalAmount,
	}
}

func (t *transactionResponseMapper) ToTransactionResponsesMonthStatusSuccess(Transactions []*record.TransactionRecordMonthStatusSuccess) []*response.TransactionResponseMonthStatusSuccess {
	var TransactionRecords []*response.TransactionResponseMonthStatusSuccess

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.ToTransactionResponseMonthStatusSuccess(Transaction))
	}

	return TransactionRecords
}

func (t *transactionResponseMapper) ToTransactionResponseYearStatusSuccess(s *record.TransactionRecordYearStatusSuccess) *response.TransactionResponseYearStatusSuccess {
	return &response.TransactionResponseYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  s.TotalAmount,
	}
}

func (t *transactionResponseMapper) ToTransactionResponsesYearStatusSuccess(Transactions []*record.TransactionRecordYearStatusSuccess) []*response.TransactionResponseYearStatusSuccess {
	var TransactionRecords []*response.TransactionResponseYearStatusSuccess

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.ToTransactionResponseYearStatusSuccess(Transaction))
	}

	return TransactionRecords
}

func (t *transactionResponseMapper) ToTransactionResponseMonthStatusFailed(s *record.TransactionRecordMonthStatusFailed) *response.TransactionResponseMonthStatusFailed {
	return &response.TransactionResponseMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: s.TotalAmount,
	}
}

func (t *transactionResponseMapper) ToTransactionResponsesMonthStatusFailed(Transactions []*record.TransactionRecordMonthStatusFailed) []*response.TransactionResponseMonthStatusFailed {
	var TransactionRecords []*response.TransactionResponseMonthStatusFailed

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.ToTransactionResponseMonthStatusFailed(Transaction))
	}

	return TransactionRecords
}

func (t *transactionResponseMapper) ToTransactionResponseYearStatusFailed(s *record.TransactionRecordYearStatusFailed) *response.TransactionResponseYearStatusFailed {
	return &response.TransactionResponseYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: s.TotalAmount,
	}
}

func (t *transactionResponseMapper) ToTransactionResponsesYearStatusFailed(Transactions []*record.TransactionRecordYearStatusFailed) []*response.TransactionResponseYearStatusFailed {
	var TransactionRecords []*response.TransactionResponseYearStatusFailed

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.ToTransactionResponseYearStatusFailed(Transaction))
	}

	return TransactionRecords
}

func (t *transactionResponseMapper) ToTransactionMonthlyMethodResponse(s *record.TransactionMonthMethod) *response.TransactionMonthMethodResponse {
	return &response.TransactionMonthMethodResponse{
		Month:             s.Month,
		PaymentMethod:     s.PaymentMethod,
		TotalTransactions: int(s.TotalTransactions),
		TotalAmount:       int(s.TotalAmount),
	}
}

func (t *transactionResponseMapper) ToTransactionMonthlyMethodResponses(s []*record.TransactionMonthMethod) []*response.TransactionMonthMethodResponse {
	var transactionResponses []*response.TransactionMonthMethodResponse
	for _, transaction := range s {
		transactionResponses = append(transactionResponses, t.ToTransactionMonthlyMethodResponse(transaction))
	}
	return transactionResponses
}

func (t *transactionResponseMapper) ToTransactionYearlyMethodResponse(s *record.TransactionYearMethod) *response.TransactionYearMethodResponse {
	return &response.TransactionYearMethodResponse{
		Year:              s.Year,
		PaymentMethod:     s.PaymentMethod,
		TotalTransactions: int(s.TotalTransactions),
		TotalAmount:       int(s.TotalAmount),
	}
}

func (t *transactionResponseMapper) ToTransactionYearlyMethodResponses(s []*record.TransactionYearMethod) []*response.TransactionYearMethodResponse {
	var transactionResponses []*response.TransactionYearMethodResponse
	for _, transaction := range s {
		transactionResponses = append(transactionResponses, t.ToTransactionYearlyMethodResponse(transaction))
	}
	return transactionResponses
}

func (t *transactionResponseMapper) ToTransactionMonthlyAmountResponse(s *record.TransactionMonthAmount) *response.TransactionMonthAmountResponse {
	return &response.TransactionMonthAmountResponse{
		Month:       s.Month,
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transactionResponseMapper) ToTransactionMonthlyAmountResponses(s []*record.TransactionMonthAmount) []*response.TransactionMonthAmountResponse {
	var transactionResponses []*response.TransactionMonthAmountResponse
	for _, transaction := range s {
		transactionResponses = append(transactionResponses, t.ToTransactionMonthlyAmountResponse(transaction))
	}
	return transactionResponses
}

func (t *transactionResponseMapper) ToTransactionYearlyAmountResponse(s *record.TransactionYearlyAmount) *response.TransactionYearlyAmountResponse {
	return &response.TransactionYearlyAmountResponse{
		Year:        s.Year,
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transactionResponseMapper) ToTransactionYearlyAmountResponses(s []*record.TransactionYearlyAmount) []*response.TransactionYearlyAmountResponse {
	var transactionResponses []*response.TransactionYearlyAmountResponse
	for _, transaction := range s {
		transactionResponses = append(transactionResponses, t.ToTransactionYearlyAmountResponse(transaction))
	}
	return transactionResponses
}
