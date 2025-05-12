package recordmapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
)

type topupRecordMapper struct {
}

func NewTopupRecordMapper() *topupRecordMapper {
	return &topupRecordMapper{}
}

func (t *topupRecordMapper) ToTopupRecord(topup *db.Topup) *record.TopupRecord {
	var deleted_at *string

	if topup.DeletedAt.Valid {
		formatedDeletedAt := topup.DeletedAt.Time.Format("2006-01-02")
		deleted_at = &formatedDeletedAt
	}

	return &record.TopupRecord{
		ID:          int(topup.TopupID),
		CardNumber:  topup.CardNumber,
		TopupNo:     topup.TopupNo.String(),
		TopupAmount: int(topup.TopupAmount),
		TopupMethod: topup.TopupMethod,
		TopupTime:   topup.TopupTime.Format("2006-01-02 15:04:05.000"),
		CreatedAt:   topup.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:   topup.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:   deleted_at,
	}
}

func (t *topupRecordMapper) ToTopupRecords(topups []*db.Topup) []*record.TopupRecord {
	var topupRecords []*record.TopupRecord

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupRecord(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupByCardNumberRecord(topup *db.GetTopupsByCardNumberRow) *record.TopupRecord {
	var deleted_at *string

	if topup.DeletedAt.Valid {
		formatedDeletedAt := topup.DeletedAt.Time.Format("2006-01-02")
		deleted_at = &formatedDeletedAt
	}

	return &record.TopupRecord{
		ID:          int(topup.TopupID),
		CardNumber:  topup.CardNumber,
		TopupNo:     topup.TopupNo.String(),
		TopupAmount: int(topup.TopupAmount),
		TopupMethod: topup.TopupMethod,
		TopupTime:   topup.TopupTime.Format("2006-01-02 15:04:05.000"),
		CreatedAt:   topup.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:   topup.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:   deleted_at,
	}
}

func (t *topupRecordMapper) ToTopupByCardNumberRecords(topups []*db.GetTopupsByCardNumberRow) []*record.TopupRecord {
	var topupRecords []*record.TopupRecord

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupByCardNumberRecord(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupRecordAll(topup *db.GetTopupsRow) *record.TopupRecord {
	var deleted_at *string

	if topup.DeletedAt.Valid {
		formatedDeletedAt := topup.DeletedAt.Time.Format("2006-01-02")
		deleted_at = &formatedDeletedAt
	}

	return &record.TopupRecord{
		ID:          int(topup.TopupID),
		CardNumber:  topup.CardNumber,
		TopupNo:     topup.TopupNo.String(),
		TopupAmount: int(topup.TopupAmount),
		TopupMethod: topup.TopupMethod,
		TopupTime:   topup.TopupTime.Format("2006-01-02 15:04:05.000"),
		CreatedAt:   topup.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:   topup.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:   deleted_at,
	}
}

func (t *topupRecordMapper) ToTopupRecordsAll(topups []*db.GetTopupsRow) []*record.TopupRecord {
	var topupRecords []*record.TopupRecord

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupRecordAll(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupRecordActive(topup *db.GetActiveTopupsRow) *record.TopupRecord {
	var deleted_at *string

	if topup.DeletedAt.Valid {
		formatedDeletedAt := topup.DeletedAt.Time.Format("2006-01-02")
		deleted_at = &formatedDeletedAt
	}

	return &record.TopupRecord{
		ID:          int(topup.TopupID),
		CardNumber:  topup.CardNumber,
		TopupNo:     topup.TopupNo.String(),
		TopupAmount: int(topup.TopupAmount),
		TopupMethod: topup.TopupMethod,
		TopupTime:   topup.TopupTime.Format("2006-01-02 15:04:05.000"),
		CreatedAt:   topup.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:   topup.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:   deleted_at,
	}
}

func (t *topupRecordMapper) ToTopupRecordsActive(topups []*db.GetActiveTopupsRow) []*record.TopupRecord {
	var topupRecords []*record.TopupRecord

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupRecordActive(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupRecordTrashed(topup *db.GetTrashedTopupsRow) *record.TopupRecord {
	var deleted_at *string

	if topup.DeletedAt.Valid {
		formatedDeletedAt := topup.DeletedAt.Time.Format("2006-01-02")
		deleted_at = &formatedDeletedAt
	}

	return &record.TopupRecord{
		ID:          int(topup.TopupID),
		CardNumber:  topup.CardNumber,
		TopupNo:     topup.TopupNo.String(),
		TopupAmount: int(topup.TopupAmount),
		TopupMethod: topup.TopupMethod,
		TopupTime:   topup.TopupTime.Format("2006-01-02 15:04:05.000"),
		CreatedAt:   topup.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:   topup.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:   deleted_at,
	}
}

func (t *topupRecordMapper) ToTopupRecordsTrashed(topups []*db.GetTrashedTopupsRow) []*record.TopupRecord {
	var topupRecords []*record.TopupRecord

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupRecordTrashed(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupRecordMonthStatusSuccess(s *db.GetMonthTopupStatusSuccessRow) *record.TopupRecordMonthStatusSuccess {
	return &record.TopupRecordMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupRecordsMonthStatusSuccess(topups []*db.GetMonthTopupStatusSuccessRow) []*record.TopupRecordMonthStatusSuccess {
	var topupRecords []*record.TopupRecordMonthStatusSuccess

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupRecordMonthStatusSuccess(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupRecordYearStatusSuccess(s *db.GetYearlyTopupStatusSuccessRow) *record.TopupRecordYearStatusSuccess {
	return &record.TopupRecordYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupRecordsYearStatusSuccess(topups []*db.GetYearlyTopupStatusSuccessRow) []*record.TopupRecordYearStatusSuccess {
	var topupRecords []*record.TopupRecordYearStatusSuccess

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupRecordYearStatusSuccess(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupRecordMonthStatusFailed(s *db.GetMonthTopupStatusFailedRow) *record.TopupRecordMonthStatusFailed {
	return &record.TopupRecordMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupRecordsMonthStatusFailed(topups []*db.GetMonthTopupStatusFailedRow) []*record.TopupRecordMonthStatusFailed {
	var topupRecords []*record.TopupRecordMonthStatusFailed

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupRecordMonthStatusFailed(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupRecordYearStatusFailed(s *db.GetYearlyTopupStatusFailedRow) *record.TopupRecordYearStatusFailed {
	return &record.TopupRecordYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupRecordsYearStatusFailed(topups []*db.GetYearlyTopupStatusFailedRow) []*record.TopupRecordYearStatusFailed {
	var topupRecords []*record.TopupRecordYearStatusFailed

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupRecordYearStatusFailed(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupRecordMonthStatusSuccessByCardNumber(s *db.GetMonthTopupStatusSuccessCardNumberRow) *record.TopupRecordMonthStatusSuccess {
	return &record.TopupRecordMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupRecordsMonthStatusSuccessByCardNumber(topups []*db.GetMonthTopupStatusSuccessCardNumberRow) []*record.TopupRecordMonthStatusSuccess {
	var topupRecords []*record.TopupRecordMonthStatusSuccess

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupRecordMonthStatusSuccessByCardNumber(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupRecordYearStatusSuccessByCardNumber(s *db.GetYearlyTopupStatusSuccessCardNumberRow) *record.TopupRecordYearStatusSuccess {
	return &record.TopupRecordYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupRecordsYearStatusSuccessByCardNumber(topups []*db.GetYearlyTopupStatusSuccessCardNumberRow) []*record.TopupRecordYearStatusSuccess {
	var topupRecords []*record.TopupRecordYearStatusSuccess

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupRecordYearStatusSuccessByCardNumber(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupRecordMonthStatusFailedByCardNumber(s *db.GetMonthTopupStatusFailedCardNumberRow) *record.TopupRecordMonthStatusFailed {
	return &record.TopupRecordMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupRecordsMonthStatusFailedByCardNumber(topups []*db.GetMonthTopupStatusFailedCardNumberRow) []*record.TopupRecordMonthStatusFailed {
	var topupRecords []*record.TopupRecordMonthStatusFailed

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupRecordMonthStatusFailedByCardNumber(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupRecordYearStatusFailedByCardNumber(s *db.GetYearlyTopupStatusFailedCardNumberRow) *record.TopupRecordYearStatusFailed {
	return &record.TopupRecordYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupRecordsYearStatusFailedByCardNumber(topups []*db.GetYearlyTopupStatusFailedCardNumberRow) []*record.TopupRecordYearStatusFailed {
	var topupRecords []*record.TopupRecordYearStatusFailed

	for _, topup := range topups {
		topupRecords = append(topupRecords, t.ToTopupRecordYearStatusFailedByCardNumber(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupMonthlyMethod(s *db.GetMonthlyTopupMethodsRow) *record.TopupMonthMethod {
	return &record.TopupMonthMethod{
		Month:       s.Month,
		TopupMethod: s.TopupMethod,
		TotalTopups: int(s.TotalTopups),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupMonthlyMethods(s []*db.GetMonthlyTopupMethodsRow) []*record.TopupMonthMethod {
	var topupRecords []*record.TopupMonthMethod

	for _, topup := range s {
		topupRecords = append(topupRecords, t.ToTopupMonthlyMethod(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupYearlyMethod(s *db.GetYearlyTopupMethodsRow) *record.TopupYearlyMethod {
	return &record.TopupYearlyMethod{
		Year:        s.Year,
		TopupMethod: s.TopupMethod,
		TotalTopups: int(s.TotalTopups),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupYearlyMethods(s []*db.GetYearlyTopupMethodsRow) []*record.TopupYearlyMethod {
	var topupRecords []*record.TopupYearlyMethod

	for _, topup := range s {
		topupRecords = append(topupRecords, t.ToTopupYearlyMethod(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupMonthlyAmount(s *db.GetMonthlyTopupAmountsRow) *record.TopupMonthAmount {
	return &record.TopupMonthAmount{
		Month:       s.Month,
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupMonthlyAmounts(s []*db.GetMonthlyTopupAmountsRow) []*record.TopupMonthAmount {
	var topupRecords []*record.TopupMonthAmount

	for _, topup := range s {
		topupRecords = append(topupRecords, t.ToTopupMonthlyAmount(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupYearlyAmount(s *db.GetYearlyTopupAmountsRow) *record.TopupYearlyAmount {
	return &record.TopupYearlyAmount{
		Year:        s.Year,
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupYearlyAmounts(s []*db.GetYearlyTopupAmountsRow) []*record.TopupYearlyAmount {
	var topupRecords []*record.TopupYearlyAmount

	for _, topup := range s {
		topupRecords = append(topupRecords, t.ToTopupYearlyAmount(topup))
	}

	return topupRecords
}

//

func (t *topupRecordMapper) ToTopupMonthlyMethodByCardNumber(s *db.GetMonthlyTopupMethodsByCardNumberRow) *record.TopupMonthMethod {
	return &record.TopupMonthMethod{
		Month:       s.Month,
		TopupMethod: s.TopupMethod,
		TotalTopups: int(s.TotalTopups),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupMonthlyMethodsByCardNumber(s []*db.GetMonthlyTopupMethodsByCardNumberRow) []*record.TopupMonthMethod {
	var topupRecords []*record.TopupMonthMethod

	for _, topup := range s {
		topupRecords = append(topupRecords, t.ToTopupMonthlyMethodByCardNumber(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupYearlyMethodByCardNumber(s *db.GetYearlyTopupMethodsByCardNumberRow) *record.TopupYearlyMethod {
	return &record.TopupYearlyMethod{
		Year:        s.Year,
		TopupMethod: s.TopupMethod,
		TotalTopups: int(s.TotalTopups),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupYearlyMethodsByCardNumber(s []*db.GetYearlyTopupMethodsByCardNumberRow) []*record.TopupYearlyMethod {
	var topupRecords []*record.TopupYearlyMethod

	for _, topup := range s {
		topupRecords = append(topupRecords, t.ToTopupYearlyMethodByCardNumber(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupMonthlyAmountByCardNumber(s *db.GetMonthlyTopupAmountsByCardNumberRow) *record.TopupMonthAmount {
	return &record.TopupMonthAmount{
		Month:       s.Month,
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupMonthlyAmountsByCardNumber(s []*db.GetMonthlyTopupAmountsByCardNumberRow) []*record.TopupMonthAmount {
	var topupRecords []*record.TopupMonthAmount

	for _, topup := range s {
		topupRecords = append(topupRecords, t.ToTopupMonthlyAmountByCardNumber(topup))
	}

	return topupRecords
}

func (t *topupRecordMapper) ToTopupYearlyAmountByCardNumber(s *db.GetYearlyTopupAmountsByCardNumberRow) *record.TopupYearlyAmount {
	return &record.TopupYearlyAmount{
		Year:        s.Year,
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *topupRecordMapper) ToTopupYearlyAmountsByCardNumber(s []*db.GetYearlyTopupAmountsByCardNumberRow) []*record.TopupYearlyAmount {
	var topupRecords []*record.TopupYearlyAmount

	for _, topup := range s {
		topupRecords = append(topupRecords, t.ToTopupYearlyAmountByCardNumber(topup))
	}

	return topupRecords
}
