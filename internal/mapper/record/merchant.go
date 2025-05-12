package recordmapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
)

type merchantRecordMapper struct {
}

func NewMerchantRecordMapper() *merchantRecordMapper {
	return &merchantRecordMapper{}
}

func (m *merchantRecordMapper) ToMerchantRecord(merchant *db.Merchant) *record.MerchantRecord {
	var deletedAt *string

	if merchant.DeletedAt.Valid {
		formatedDeletedAt := merchant.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.MerchantRecord{
		ID:        int(merchant.MerchantID),
		Name:      merchant.Name,
		ApiKey:    merchant.ApiKey,
		UserID:    int(merchant.UserID),
		Status:    merchant.Status,
		CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: deletedAt,
	}
}

func (m *merchantRecordMapper) ToMerchantsRecord(merchants []*db.Merchant) []*record.MerchantRecord {
	var records []*record.MerchantRecord
	for _, merchant := range merchants {
		records = append(records, m.ToMerchantRecord(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantTransactionRecord(merchant *db.FindAllTransactionsRow) *record.MerchantTransactionsRecord {
	var deletedAt *string

	if merchant.DeletedAt.Valid {
		formatedDeletedAt := merchant.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.MerchantTransactionsRecord{
		TransactionID:   merchant.TransactionID,
		CardNumber:      merchant.CardNumber,
		Amount:          merchant.Amount,
		PaymentMethod:   merchant.PaymentMethod,
		MerchantID:      merchant.MerchantID,
		MerchantName:    merchant.MerchantName,
		TransactionTime: merchant.TransactionTime,
		CreatedAt:       merchant.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:       merchant.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:       deletedAt,
	}
}

func (m *merchantRecordMapper) ToMerchantsTransactionRecord(merchants []*db.FindAllTransactionsRow) []*record.MerchantTransactionsRecord {
	var records []*record.MerchantTransactionsRecord
	for _, merchant := range merchants {
		records = append(records, m.ToMerchantTransactionRecord(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantGetAllRecord(merchant *db.GetMerchantsRow) *record.MerchantRecord {
	var deletedAt *string

	if merchant.DeletedAt.Valid {
		formatedDeletedAt := merchant.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.MerchantRecord{
		ID:        int(merchant.MerchantID),
		Name:      merchant.Name,
		ApiKey:    merchant.ApiKey,
		UserID:    int(merchant.UserID),
		Status:    merchant.Status,
		CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: deletedAt,
	}
}

func (m *merchantRecordMapper) ToMerchantsGetAllRecord(merchants []*db.GetMerchantsRow) []*record.MerchantRecord {
	var records []*record.MerchantRecord
	for _, merchant := range merchants {
		records = append(records, m.ToMerchantGetAllRecord(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantMonthlyPaymentMethod(ms *db.GetMonthlyPaymentMethodsMerchantRow) *record.MerchantMonthlyPaymentMethod {
	return &record.MerchantMonthlyPaymentMethod{
		Month:         ms.Month,
		PaymentMethod: ms.PaymentMethod,
		TotalAmount:   int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantMonthlyPaymentMethods(ms []*db.GetMonthlyPaymentMethodsMerchantRow) []*record.MerchantMonthlyPaymentMethod {
	var records []*record.MerchantMonthlyPaymentMethod
	for _, merchant := range ms {
		records = append(records, m.ToMerchantMonthlyPaymentMethod(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantYearlyPaymentMethod(ms *db.GetYearlyPaymentMethodMerchantRow) *record.MerchantYearlyPaymentMethod {
	return &record.MerchantYearlyPaymentMethod{
		Year:          ms.Year,
		PaymentMethod: ms.PaymentMethod,
		TotalAmount:   int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantYearlyPaymentMethods(ms []*db.GetYearlyPaymentMethodMerchantRow) []*record.MerchantYearlyPaymentMethod {
	var records []*record.MerchantYearlyPaymentMethod
	for _, merchant := range ms {
		records = append(records, m.ToMerchantYearlyPaymentMethod(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantMonthlyAmount(ms *db.GetMonthlyAmountMerchantRow) *record.MerchantMonthlyAmount {
	return &record.MerchantMonthlyAmount{
		Month:       ms.Month,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantMonthlyAmounts(ms []*db.GetMonthlyAmountMerchantRow) []*record.MerchantMonthlyAmount {
	var records []*record.MerchantMonthlyAmount
	for _, merchant := range ms {
		records = append(records, m.ToMerchantMonthlyAmount(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantYearlyAmount(ms *db.GetYearlyAmountMerchantRow) *record.MerchantYearlyAmount {
	return &record.MerchantYearlyAmount{
		Year:        ms.Year,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantYearlyAmounts(ms []*db.GetYearlyAmountMerchantRow) []*record.MerchantYearlyAmount {
	var records []*record.MerchantYearlyAmount
	for _, merchant := range ms {
		records = append(records, m.ToMerchantYearlyAmount(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantMonthlyTotalAmount(ms *db.GetMonthlyTotalAmountMerchantRow) *record.MerchantMonthlyTotalAmount {
	return &record.MerchantMonthlyTotalAmount{
		Month:       ms.Month,
		Year:        ms.Year,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantMonthlyTotalAmounts(ms []*db.GetMonthlyTotalAmountMerchantRow) []*record.MerchantMonthlyTotalAmount {
	var records []*record.MerchantMonthlyTotalAmount
	for _, merchant := range ms {
		records = append(records, m.ToMerchantMonthlyTotalAmount(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantYearlyTotalAmount(ms *db.GetYearlyTotalAmountMerchantRow) *record.MerchantYearlyTotalAmount {
	return &record.MerchantYearlyTotalAmount{
		Year:        ms.Year,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantYearlyTotalAmounts(ms []*db.GetYearlyTotalAmountMerchantRow) []*record.MerchantYearlyTotalAmount {
	var records []*record.MerchantYearlyTotalAmount
	for _, merchant := range ms {
		records = append(records, m.ToMerchantYearlyTotalAmount(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantTransactionByMerchantRecord(merchant *db.FindAllTransactionsByMerchantRow) *record.MerchantTransactionsRecord {
	var deletedAt *string

	if merchant.DeletedAt.Valid {
		formatedDeletedAt := merchant.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.MerchantTransactionsRecord{
		TransactionID:   merchant.TransactionID,
		CardNumber:      merchant.CardNumber,
		Amount:          merchant.Amount,
		PaymentMethod:   merchant.PaymentMethod,
		MerchantID:      merchant.MerchantID,
		MerchantName:    merchant.MerchantName,
		TransactionTime: merchant.TransactionTime,
		CreatedAt:       merchant.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:       merchant.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:       deletedAt,
	}
}

func (m *merchantRecordMapper) ToMerchantsTransactionByMerchantRecord(merchants []*db.FindAllTransactionsByMerchantRow) []*record.MerchantTransactionsRecord {
	var records []*record.MerchantTransactionsRecord
	for _, merchant := range merchants {
		records = append(records, m.ToMerchantTransactionByMerchantRecord(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantMonthlyPaymentMethodByMerchant(ms *db.GetMonthlyPaymentMethodByMerchantsRow) *record.MerchantMonthlyPaymentMethod {
	return &record.MerchantMonthlyPaymentMethod{
		Month:         ms.Month,
		PaymentMethod: ms.PaymentMethod,
		TotalAmount:   int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantMonthlyPaymentMethodsByMerchant(ms []*db.GetMonthlyPaymentMethodByMerchantsRow) []*record.MerchantMonthlyPaymentMethod {
	var records []*record.MerchantMonthlyPaymentMethod
	for _, merchant := range ms {
		records = append(records, m.ToMerchantMonthlyPaymentMethodByMerchant(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantYearlyPaymentMethodByMerchant(ms *db.GetYearlyPaymentMethodByMerchantsRow) *record.MerchantYearlyPaymentMethod {
	return &record.MerchantYearlyPaymentMethod{
		Year:          ms.Year,
		PaymentMethod: ms.PaymentMethod,
		TotalAmount:   int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantYearlyPaymentMethodsByMerchant(ms []*db.GetYearlyPaymentMethodByMerchantsRow) []*record.MerchantYearlyPaymentMethod {
	var records []*record.MerchantYearlyPaymentMethod
	for _, merchant := range ms {
		records = append(records, m.ToMerchantYearlyPaymentMethodByMerchant(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantMonthlyAmountByMerchant(ms *db.GetMonthlyAmountByMerchantsRow) *record.MerchantMonthlyAmount {
	return &record.MerchantMonthlyAmount{
		Month:       ms.Month,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantMonthlyAmountsByMerchant(ms []*db.GetMonthlyAmountByMerchantsRow) []*record.MerchantMonthlyAmount {
	var records []*record.MerchantMonthlyAmount
	for _, merchant := range ms {
		records = append(records, m.ToMerchantMonthlyAmountByMerchant(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantYearlyAmountByMerchant(ms *db.GetYearlyAmountByMerchantsRow) *record.MerchantYearlyAmount {
	return &record.MerchantYearlyAmount{
		Year:        ms.Year,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantYearlyAmountsByMerchant(ms []*db.GetYearlyAmountByMerchantsRow) []*record.MerchantYearlyAmount {
	var records []*record.MerchantYearlyAmount
	for _, merchant := range ms {
		records = append(records, m.ToMerchantYearlyAmountByMerchant(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantMonthlyTotalAmountByMerchant(ms *db.GetMonthlyTotalAmountByMerchantRow) *record.MerchantMonthlyTotalAmount {
	return &record.MerchantMonthlyTotalAmount{
		Month:       ms.Month,
		Year:        ms.Year,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantMonthlyTotalAmountsByMerchant(ms []*db.GetMonthlyTotalAmountByMerchantRow) []*record.MerchantMonthlyTotalAmount {
	var records []*record.MerchantMonthlyTotalAmount
	for _, merchant := range ms {
		records = append(records, m.ToMerchantMonthlyTotalAmountByMerchant(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantYearlyTotalAmountByMerchant(ms *db.GetYearlyTotalAmountByMerchantRow) *record.MerchantYearlyTotalAmount {
	return &record.MerchantYearlyTotalAmount{
		Year:        ms.Year,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantYearlyTotalAmountsByMerchant(ms []*db.GetYearlyTotalAmountByMerchantRow) []*record.MerchantYearlyTotalAmount {
	var records []*record.MerchantYearlyTotalAmount
	for _, merchant := range ms {
		records = append(records, m.ToMerchantYearlyTotalAmountByMerchant(merchant))
	}
	return records
}

//

func (m *merchantRecordMapper) ToMerchantTransactionByApikeyRecord(merchant *db.FindAllTransactionsByApikeyRow) *record.MerchantTransactionsRecord {
	var deletedAt *string

	if merchant.DeletedAt.Valid {
		formatedDeletedAt := merchant.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.MerchantTransactionsRecord{
		TransactionID:   merchant.TransactionID,
		CardNumber:      merchant.CardNumber,
		Amount:          merchant.Amount,
		PaymentMethod:   merchant.PaymentMethod,
		MerchantID:      merchant.MerchantID,
		MerchantName:    merchant.MerchantName,
		TransactionTime: merchant.TransactionTime,
		CreatedAt:       merchant.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:       merchant.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:       deletedAt,
	}
}

func (m *merchantRecordMapper) ToMerchantsTransactionByApikeyRecord(merchants []*db.FindAllTransactionsByApikeyRow) []*record.MerchantTransactionsRecord {
	var records []*record.MerchantTransactionsRecord
	for _, merchant := range merchants {
		records = append(records, m.ToMerchantTransactionByApikeyRecord(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantMonthlyPaymentMethodByApikey(ms *db.GetMonthlyPaymentMethodByApikeyRow) *record.MerchantMonthlyPaymentMethod {
	return &record.MerchantMonthlyPaymentMethod{
		Month:         ms.Month,
		PaymentMethod: ms.PaymentMethod,
		TotalAmount:   int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantMonthlyPaymentMethodsByApikey(ms []*db.GetMonthlyPaymentMethodByApikeyRow) []*record.MerchantMonthlyPaymentMethod {
	var records []*record.MerchantMonthlyPaymentMethod
	for _, merchant := range ms {
		records = append(records, m.ToMerchantMonthlyPaymentMethodByApikey(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantYearlyPaymentMethodByApikey(ms *db.GetYearlyPaymentMethodByApikeyRow) *record.MerchantYearlyPaymentMethod {
	return &record.MerchantYearlyPaymentMethod{
		Year:          ms.Year,
		PaymentMethod: ms.PaymentMethod,
		TotalAmount:   int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantYearlyPaymentMethodsByApikey(ms []*db.GetYearlyPaymentMethodByApikeyRow) []*record.MerchantYearlyPaymentMethod {
	var records []*record.MerchantYearlyPaymentMethod
	for _, merchant := range ms {
		records = append(records, m.ToMerchantYearlyPaymentMethodByApikey(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantMonthlyAmountByApikey(ms *db.GetMonthlyAmountByApikeyRow) *record.MerchantMonthlyAmount {
	return &record.MerchantMonthlyAmount{
		Month:       ms.Month,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantMonthlyAmountsByApikey(ms []*db.GetMonthlyAmountByApikeyRow) []*record.MerchantMonthlyAmount {
	var records []*record.MerchantMonthlyAmount
	for _, merchant := range ms {
		records = append(records, m.ToMerchantMonthlyAmountByApikey(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantYearlyAmountByApikey(ms *db.GetYearlyAmountByApikeyRow) *record.MerchantYearlyAmount {
	return &record.MerchantYearlyAmount{
		Year:        ms.Year,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantYearlyAmountsByApikey(ms []*db.GetYearlyAmountByApikeyRow) []*record.MerchantYearlyAmount {
	var records []*record.MerchantYearlyAmount
	for _, merchant := range ms {
		records = append(records, m.ToMerchantYearlyAmountByApikey(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantMonthlyTotalAmountByApikey(ms *db.GetMonthlyTotalAmountByApikeyRow) *record.MerchantMonthlyTotalAmount {
	return &record.MerchantMonthlyTotalAmount{
		Month:       ms.Month,
		Year:        ms.Year,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantMonthlyTotalAmountsByApikey(ms []*db.GetMonthlyTotalAmountByApikeyRow) []*record.MerchantMonthlyTotalAmount {
	var records []*record.MerchantMonthlyTotalAmount
	for _, merchant := range ms {
		records = append(records, m.ToMerchantMonthlyTotalAmountByApikey(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantYearlyTotalAmountByApikey(ms *db.GetYearlyTotalAmountByApikeyRow) *record.MerchantYearlyTotalAmount {
	return &record.MerchantYearlyTotalAmount{
		Year:        ms.Year,
		TotalAmount: int(ms.TotalAmount),
	}
}

func (m *merchantRecordMapper) ToMerchantYearlyTotalAmountsByApikey(ms []*db.GetYearlyTotalAmountByApikeyRow) []*record.MerchantYearlyTotalAmount {
	var records []*record.MerchantYearlyTotalAmount
	for _, merchant := range ms {
		records = append(records, m.ToMerchantYearlyTotalAmountByApikey(merchant))
	}
	return records
}

//

func (m *merchantRecordMapper) ToMerchantActiveRecord(merchant *db.GetActiveMerchantsRow) *record.MerchantRecord {
	var deletedAt *string

	if merchant.DeletedAt.Valid {
		formatedDeletedAt := merchant.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.MerchantRecord{
		ID:        int(merchant.MerchantID),
		Name:      merchant.Name,
		ApiKey:    merchant.ApiKey,
		UserID:    int(merchant.UserID),
		Status:    merchant.Status,
		CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: deletedAt,
	}
}

func (m *merchantRecordMapper) ToMerchantsActiveRecord(merchants []*db.GetActiveMerchantsRow) []*record.MerchantRecord {
	var records []*record.MerchantRecord
	for _, merchant := range merchants {
		records = append(records, m.ToMerchantActiveRecord(merchant))
	}
	return records
}

func (m *merchantRecordMapper) ToMerchantTrashedRecord(merchant *db.GetTrashedMerchantsRow) *record.MerchantRecord {
	var deletedAt *string

	if merchant.DeletedAt.Valid {
		formatedDeletedAt := merchant.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.MerchantRecord{
		ID:        int(merchant.MerchantID),
		Name:      merchant.Name,
		ApiKey:    merchant.ApiKey,
		UserID:    int(merchant.UserID),
		Status:    merchant.Status,
		CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: merchant.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: deletedAt,
	}
}

func (m *merchantRecordMapper) ToMerchantsTrashedRecord(merchants []*db.GetTrashedMerchantsRow) []*record.MerchantRecord {
	var records []*record.MerchantRecord
	for _, merchant := range merchants {
		records = append(records, m.ToMerchantTrashedRecord(merchant))
	}
	return records
}
