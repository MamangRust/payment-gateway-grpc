package recordmapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
)

type transactionRecordMapper struct{}

func NewTransactionRecordMapper() *transactionRecordMapper {
	return &transactionRecordMapper{}
}

func (s *transactionRecordMapper) ToTransactionRecord(transaction *db.Transaction) *record.TransactionRecord {
	var deletedAt *string

	if transaction.DeletedAt.Valid {
		formatedDeletedAt := transaction.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.TransactionRecord{
		ID:              int(transaction.TransactionID),
		TransactionNo:   transaction.TransactionNo.String(),
		CardNumber:      transaction.CardNumber,
		Amount:          int(transaction.Amount),
		PaymentMethod:   transaction.PaymentMethod,
		MerchantID:      int(transaction.MerchantID),
		TransactionTime: transaction.TransactionTime.Format("2006-01-02 15:04:05"),
		CreatedAt:       transaction.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:       transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:       deletedAt,
	}
}

func (s *transactionRecordMapper) ToTransactionsRecord(transactions []*db.Transaction) []*record.TransactionRecord {
	var transactionRecords []*record.TransactionRecord
	for _, transaction := range transactions {
		transactionRecords = append(transactionRecords, s.ToTransactionRecord(transaction))
	}
	return transactionRecords
}

func (s *transactionRecordMapper) ToTransactionByCardNumberRecord(transaction *db.GetTransactionsByCardNumberRow) *record.TransactionRecord {
	var deletedAt *string

	if transaction.DeletedAt.Valid {
		formatedDeletedAt := transaction.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.TransactionRecord{
		ID:              int(transaction.TransactionID),
		TransactionNo:   transaction.TransactionNo.String(),
		CardNumber:      transaction.CardNumber,
		Amount:          int(transaction.Amount),
		PaymentMethod:   transaction.PaymentMethod,
		MerchantID:      int(transaction.MerchantID),
		TransactionTime: transaction.TransactionTime.Format("2006-01-02 15:04:05"),
		CreatedAt:       transaction.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:       transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:       deletedAt,
	}
}

func (s *transactionRecordMapper) ToTransactionsByCardNumberRecord(transactions []*db.GetTransactionsByCardNumberRow) []*record.TransactionRecord {
	var transactionRecords []*record.TransactionRecord
	for _, transaction := range transactions {
		transactionRecords = append(transactionRecords, s.ToTransactionByCardNumberRecord(transaction))
	}
	return transactionRecords
}

func (s *transactionRecordMapper) ToTransactionRecordAll(transaction *db.GetTransactionsRow) *record.TransactionRecord {
	var deletedAt *string

	if transaction.DeletedAt.Valid {
		formatedDeletedAt := transaction.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.TransactionRecord{
		ID:              int(transaction.TransactionID),
		TransactionNo:   transaction.TransactionNo.String(),
		CardNumber:      transaction.CardNumber,
		Amount:          int(transaction.Amount),
		PaymentMethod:   transaction.PaymentMethod,
		MerchantID:      int(transaction.MerchantID),
		TransactionTime: transaction.TransactionTime.Format("2006-01-02 15:04:05"),
		CreatedAt:       transaction.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:       transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:       deletedAt,
	}
}

func (s *transactionRecordMapper) ToTransactionsRecordAll(transactions []*db.GetTransactionsRow) []*record.TransactionRecord {
	var transactionRecords []*record.TransactionRecord
	for _, transaction := range transactions {
		transactionRecords = append(transactionRecords, s.ToTransactionRecordAll(transaction))
	}
	return transactionRecords
}

func (s *transactionRecordMapper) ToTransactionRecordActive(transaction *db.GetActiveTransactionsRow) *record.TransactionRecord {
	var deletedAt *string

	if transaction.DeletedAt.Valid {
		formatedDeletedAt := transaction.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.TransactionRecord{
		ID:              int(transaction.TransactionID),
		TransactionNo:   transaction.TransactionNo.String(),
		CardNumber:      transaction.CardNumber,
		Amount:          int(transaction.Amount),
		PaymentMethod:   transaction.PaymentMethod,
		MerchantID:      int(transaction.MerchantID),
		TransactionTime: transaction.TransactionTime.Format("2006-01-02 15:04:05"),
		CreatedAt:       transaction.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:       transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:       deletedAt,
	}
}

func (s *transactionRecordMapper) ToTransactionsRecordActive(transactions []*db.GetActiveTransactionsRow) []*record.TransactionRecord {
	var transactionRecords []*record.TransactionRecord
	for _, transaction := range transactions {
		transactionRecords = append(transactionRecords, s.ToTransactionRecordActive(transaction))
	}
	return transactionRecords
}

func (s *transactionRecordMapper) ToTransactionRecordTrashed(transaction *db.GetTrashedTransactionsRow) *record.TransactionRecord {
	var deletedAt *string

	if transaction.DeletedAt.Valid {
		formatedDeletedAt := transaction.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.TransactionRecord{
		ID:              int(transaction.TransactionID),
		TransactionNo:   transaction.TransactionNo.String(),
		CardNumber:      transaction.CardNumber,
		Amount:          int(transaction.Amount),
		PaymentMethod:   transaction.PaymentMethod,
		MerchantID:      int(transaction.MerchantID),
		TransactionTime: transaction.TransactionTime.Format("2006-01-02 15:04:05"),
		CreatedAt:       transaction.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:       transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:       deletedAt,
	}
}

func (s *transactionRecordMapper) ToTransactionsRecordTrashed(transactions []*db.GetTrashedTransactionsRow) []*record.TransactionRecord {
	var transactionRecords []*record.TransactionRecord
	for _, transaction := range transactions {
		transactionRecords = append(transactionRecords, s.ToTransactionRecordTrashed(transaction))
	}
	return transactionRecords
}

func (t *transactionRecordMapper) ToTransactionRecordMonthStatusSuccess(s *db.GetMonthTransactionStatusSuccessRow) *record.TransactionRecordMonthStatusSuccess {
	return &record.TransactionRecordMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *transactionRecordMapper) ToTransactionRecordsMonthStatusSuccess(Transactions []*db.GetMonthTransactionStatusSuccessRow) []*record.TransactionRecordMonthStatusSuccess {
	var TransactionRecords []*record.TransactionRecordMonthStatusSuccess

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.ToTransactionRecordMonthStatusSuccess(Transaction))
	}

	return TransactionRecords
}

func (t *transactionRecordMapper) ToTransactionRecordYearStatusSuccess(s *db.GetYearlyTransactionStatusSuccessRow) *record.TransactionRecordYearStatusSuccess {
	return &record.TransactionRecordYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *transactionRecordMapper) ToTransactionRecordsYearStatusSuccess(Transactions []*db.GetYearlyTransactionStatusSuccessRow) []*record.TransactionRecordYearStatusSuccess {
	var TransactionRecords []*record.TransactionRecordYearStatusSuccess

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.ToTransactionRecordYearStatusSuccess(Transaction))
	}

	return TransactionRecords
}

func (t *transactionRecordMapper) ToTransactionRecordMonthStatusFailed(s *db.GetMonthTransactionStatusFailedRow) *record.TransactionRecordMonthStatusFailed {
	return &record.TransactionRecordMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transactionRecordMapper) ToTransactionRecordsMonthStatusFailed(Transactions []*db.GetMonthTransactionStatusFailedRow) []*record.TransactionRecordMonthStatusFailed {
	var TransactionRecords []*record.TransactionRecordMonthStatusFailed

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.ToTransactionRecordMonthStatusFailed(Transaction))
	}

	return TransactionRecords
}

func (t *transactionRecordMapper) ToTransactionRecordYearStatusFailed(s *db.GetYearlyTransactionStatusFailedRow) *record.TransactionRecordYearStatusFailed {
	return &record.TransactionRecordYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transactionRecordMapper) ToTransactionRecordsYearStatusFailed(Transactions []*db.GetYearlyTransactionStatusFailedRow) []*record.TransactionRecordYearStatusFailed {
	var TransactionRecords []*record.TransactionRecordYearStatusFailed

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.ToTransactionRecordYearStatusFailed(Transaction))
	}

	return TransactionRecords
}

func (t *transactionRecordMapper) ToTransactionRecordMonthStatusSuccessCardNumber(s *db.GetMonthTransactionStatusSuccessCardNumberRow) *record.TransactionRecordMonthStatusSuccess {
	return &record.TransactionRecordMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *transactionRecordMapper) ToTransactionRecordsMonthStatusSuccessCardNumber(Transactions []*db.GetMonthTransactionStatusSuccessCardNumberRow) []*record.TransactionRecordMonthStatusSuccess {
	var TransactionRecords []*record.TransactionRecordMonthStatusSuccess

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.ToTransactionRecordMonthStatusSuccessCardNumber(Transaction))
	}

	return TransactionRecords
}

func (t *transactionRecordMapper) ToTransactionRecordYearStatusSuccessCardNumber(s *db.GetYearlyTransactionStatusSuccessCardNumberRow) *record.TransactionRecordYearStatusSuccess {
	return &record.TransactionRecordYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *transactionRecordMapper) ToTransactionRecordsYearStatusSuccessCardNumber(Transactions []*db.GetYearlyTransactionStatusSuccessCardNumberRow) []*record.TransactionRecordYearStatusSuccess {
	var TransactionRecords []*record.TransactionRecordYearStatusSuccess

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.ToTransactionRecordYearStatusSuccessCardNumber(Transaction))
	}

	return TransactionRecords
}

func (t *transactionRecordMapper) ToTransactionRecordMonthStatusFailedCardNumber(s *db.GetMonthTransactionStatusFailedCardNumberRow) *record.TransactionRecordMonthStatusFailed {
	return &record.TransactionRecordMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transactionRecordMapper) ToTransactionRecordsMonthStatusFailedCardNumber(Transactions []*db.GetMonthTransactionStatusFailedCardNumberRow) []*record.TransactionRecordMonthStatusFailed {
	var TransactionRecords []*record.TransactionRecordMonthStatusFailed

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.ToTransactionRecordMonthStatusFailedCardNumber(Transaction))
	}

	return TransactionRecords
}

func (t *transactionRecordMapper) ToTransactionRecordYearStatusFailedCardNumber(s *db.GetYearlyTransactionStatusFailedCardNumberRow) *record.TransactionRecordYearStatusFailed {
	return &record.TransactionRecordYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transactionRecordMapper) ToTransactionRecordsYearStatusFailedCardNumber(Transactions []*db.GetYearlyTransactionStatusFailedCardNumberRow) []*record.TransactionRecordYearStatusFailed {
	var TransactionRecords []*record.TransactionRecordYearStatusFailed

	for _, Transaction := range Transactions {
		TransactionRecords = append(TransactionRecords, t.ToTransactionRecordYearStatusFailedCardNumber(Transaction))
	}

	return TransactionRecords
}

func (s *transactionRecordMapper) ToTransactionMonthlyMethod(ss *db.GetMonthlyPaymentMethodsRow) *record.TransactionMonthMethod {
	return &record.TransactionMonthMethod{
		Month:             ss.Month,
		PaymentMethod:     ss.PaymentMethod,
		TotalTransactions: int(ss.TotalTransactions),
		TotalAmount:       int(ss.TotalAmount),
	}
}

func (s *transactionRecordMapper) ToTransactionMonthlyMethods(ss []*db.GetMonthlyPaymentMethodsRow) []*record.TransactionMonthMethod {
	var transactionRecords []*record.TransactionMonthMethod
	for _, transaction := range ss {
		transactionRecords = append(transactionRecords, s.ToTransactionMonthlyMethod(transaction))
	}
	return transactionRecords
}

func (s *transactionRecordMapper) ToTransactionYearlyMethod(ss *db.GetYearlyPaymentMethodsRow) *record.TransactionYearMethod {
	return &record.TransactionYearMethod{
		Year:              ss.Year,
		PaymentMethod:     ss.PaymentMethod,
		TotalTransactions: int(ss.TotalTransactions),
		TotalAmount:       int(ss.TotalAmount),
	}
}

func (s *transactionRecordMapper) ToTransactionYearlyMethods(ss []*db.GetYearlyPaymentMethodsRow) []*record.TransactionYearMethod {
	var transactionRecords []*record.TransactionYearMethod
	for _, transaction := range ss {
		transactionRecords = append(transactionRecords, s.ToTransactionYearlyMethod(transaction))
	}
	return transactionRecords
}

//

func (s *transactionRecordMapper) ToTransactionMonthlyAmount(ss *db.GetMonthlyAmountsRow) *record.TransactionMonthAmount {
	return &record.TransactionMonthAmount{
		Month:       ss.Month,
		TotalAmount: int(ss.TotalAmount),
	}
}

func (s *transactionRecordMapper) ToTransactionMonthlyAmounts(ss []*db.GetMonthlyAmountsRow) []*record.TransactionMonthAmount {
	var transactionRecords []*record.TransactionMonthAmount
	for _, transaction := range ss {
		transactionRecords = append(transactionRecords, s.ToTransactionMonthlyAmount(transaction))
	}
	return transactionRecords
}

func (s *transactionRecordMapper) ToTransactionYearlyAmount(ss *db.GetYearlyAmountsRow) *record.TransactionYearlyAmount {
	return &record.TransactionYearlyAmount{
		Year:        ss.Year,
		TotalAmount: int(ss.TotalAmount),
	}
}

func (s *transactionRecordMapper) ToTransactionYearlyAmounts(ss []*db.GetYearlyAmountsRow) []*record.TransactionYearlyAmount {
	var transactionRecords []*record.TransactionYearlyAmount
	for _, transaction := range ss {
		transactionRecords = append(transactionRecords, s.ToTransactionYearlyAmount(transaction))
	}
	return transactionRecords
}

/////

func (s *transactionRecordMapper) ToTransactionMonthlyMethodByCardNumber(ss *db.GetMonthlyPaymentMethodsByCardNumberRow) *record.TransactionMonthMethod {
	return &record.TransactionMonthMethod{
		Month:             ss.Month,
		PaymentMethod:     ss.PaymentMethod,
		TotalTransactions: int(ss.TotalTransactions),
		TotalAmount:       int(ss.TotalAmount),
	}
}

func (s *transactionRecordMapper) ToTransactionMonthlyMethodsByCardNumber(ss []*db.GetMonthlyPaymentMethodsByCardNumberRow) []*record.TransactionMonthMethod {
	var transactionRecords []*record.TransactionMonthMethod
	for _, transaction := range ss {
		transactionRecords = append(transactionRecords, s.ToTransactionMonthlyMethodByCardNumber(transaction))
	}
	return transactionRecords
}

func (s *transactionRecordMapper) ToTransactionYearlyMethodByCardNumber(ss *db.GetYearlyPaymentMethodsByCardNumberRow) *record.TransactionYearMethod {
	return &record.TransactionYearMethod{
		Year:              ss.Year,
		PaymentMethod:     ss.PaymentMethod,
		TotalTransactions: int(ss.TotalTransactions),
		TotalAmount:       int(ss.TotalAmount),
	}
}

func (s *transactionRecordMapper) ToTransactionYearlyMethodsByCardNumber(ss []*db.GetYearlyPaymentMethodsByCardNumberRow) []*record.TransactionYearMethod {
	var transactionRecords []*record.TransactionYearMethod
	for _, transaction := range ss {
		transactionRecords = append(transactionRecords, s.ToTransactionYearlyMethodByCardNumber(transaction))
	}
	return transactionRecords
}

//

func (s *transactionRecordMapper) ToTransactionMonthlyAmountByCardNumber(ss *db.GetMonthlyAmountsByCardNumberRow) *record.TransactionMonthAmount {
	return &record.TransactionMonthAmount{
		Month:       ss.Month,
		TotalAmount: int(ss.TotalAmount),
	}

}

func (s *transactionRecordMapper) ToTransactionMonthlyAmountsByCardNumber(ss []*db.GetMonthlyAmountsByCardNumberRow) []*record.TransactionMonthAmount {
	var transactionRecords []*record.TransactionMonthAmount
	for _, transaction := range ss {
		transactionRecords = append(transactionRecords, s.ToTransactionMonthlyAmountByCardNumber(transaction))
	}
	return transactionRecords

}

func (s *transactionRecordMapper) ToTransactionYearlyAmountByCardNumber(ss *db.GetYearlyAmountsByCardNumberRow) *record.TransactionYearlyAmount {
	return &record.TransactionYearlyAmount{
		Year:        ss.Year,
		TotalAmount: int(ss.TotalAmount),
	}
}

func (s *transactionRecordMapper) ToTransactionYearlyAmountsByCardNumber(ss []*db.GetYearlyAmountsByCardNumberRow) []*record.TransactionYearlyAmount {
	var transactionRecords []*record.TransactionYearlyAmount
	for _, transaction := range ss {
		transactionRecords = append(transactionRecords, s.ToTransactionYearlyAmountByCardNumber(transaction))
	}
	return transactionRecords

}
