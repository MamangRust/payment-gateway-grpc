package recordmapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
)

type withdrawRecordMapper struct{}

func NewWithdrawRecordMapper() *withdrawRecordMapper {
	return &withdrawRecordMapper{}
}

func (s *withdrawRecordMapper) ToWithdrawRecord(withdraw *db.Withdraw) *record.WithdrawRecord {
	var deletedAt *string

	if withdraw.DeletedAt.Valid {
		formatedDeletedAt := withdraw.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.WithdrawRecord{
		ID:             int(withdraw.WithdrawID),
		WithdrawNo:     withdraw.WithdrawNo.String(),
		CardNumber:     withdraw.CardNumber,
		WithdrawAmount: int(withdraw.WithdrawAmount),
		WithdrawTime:   withdraw.WithdrawTime.String(),
		CreatedAt:      withdraw.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      withdraw.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *withdrawRecordMapper) ToWithdrawsRecord(withdraws []*db.Withdraw) []*record.WithdrawRecord {
	var withdrawRecords []*record.WithdrawRecord

	for _, withdraw := range withdraws {
		withdrawRecords = append(withdrawRecords, s.ToWithdrawRecord(withdraw))
	}

	return withdrawRecords
}

func (s *withdrawRecordMapper) ToWithdrawByCardNumberRecord(withdraw *db.GetWithdrawsByCardNumberRow) *record.WithdrawRecord {
	var deletedAt *string

	if withdraw.DeletedAt.Valid {
		formatedDeletedAt := withdraw.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.WithdrawRecord{
		ID:             int(withdraw.WithdrawID),
		WithdrawNo:     withdraw.WithdrawNo.String(),
		CardNumber:     withdraw.CardNumber,
		WithdrawAmount: int(withdraw.WithdrawAmount),
		WithdrawTime:   withdraw.WithdrawTime.String(),
		CreatedAt:      withdraw.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      withdraw.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *withdrawRecordMapper) ToWithdrawsByCardNumberRecord(withdraws []*db.GetWithdrawsByCardNumberRow) []*record.WithdrawRecord {
	var withdrawRecords []*record.WithdrawRecord

	for _, withdraw := range withdraws {
		withdrawRecords = append(withdrawRecords, s.ToWithdrawByCardNumberRecord(withdraw))
	}

	return withdrawRecords
}

func (s *withdrawRecordMapper) ToWithdrawRecordAll(withdraw *db.GetWithdrawsRow) *record.WithdrawRecord {
	var deletedAt *string

	if withdraw.DeletedAt.Valid {
		formatedDeletedAt := withdraw.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.WithdrawRecord{
		ID:             int(withdraw.WithdrawID),
		WithdrawNo:     withdraw.WithdrawNo.String(),
		CardNumber:     withdraw.CardNumber,
		WithdrawAmount: int(withdraw.WithdrawAmount),
		WithdrawTime:   withdraw.WithdrawTime.String(),
		CreatedAt:      withdraw.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      withdraw.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *withdrawRecordMapper) ToWithdrawsRecordALl(withdraws []*db.GetWithdrawsRow) []*record.WithdrawRecord {
	var withdrawRecords []*record.WithdrawRecord

	for _, withdraw := range withdraws {
		withdrawRecords = append(withdrawRecords, s.ToWithdrawRecordAll(withdraw))
	}

	return withdrawRecords
}

func (s *withdrawRecordMapper) ToWithdrawRecordActive(withdraw *db.GetActiveWithdrawsRow) *record.WithdrawRecord {
	var deletedAt *string

	if withdraw.DeletedAt.Valid {
		formatedDeletedAt := withdraw.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.WithdrawRecord{
		ID:             int(withdraw.WithdrawID),
		WithdrawNo:     withdraw.WithdrawNo.String(),
		CardNumber:     withdraw.CardNumber,
		WithdrawAmount: int(withdraw.WithdrawAmount),
		WithdrawTime:   withdraw.WithdrawTime.String(),
		CreatedAt:      withdraw.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      withdraw.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *withdrawRecordMapper) ToWithdrawsRecordActive(withdraws []*db.GetActiveWithdrawsRow) []*record.WithdrawRecord {
	var withdrawRecords []*record.WithdrawRecord

	for _, withdraw := range withdraws {
		withdrawRecords = append(withdrawRecords, s.ToWithdrawRecordActive(withdraw))
	}

	return withdrawRecords
}

func (s *withdrawRecordMapper) ToWithdrawRecordTrashed(withdraw *db.GetTrashedWithdrawsRow) *record.WithdrawRecord {
	var deletedAt *string

	if withdraw.DeletedAt.Valid {
		formatedDeletedAt := withdraw.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.WithdrawRecord{
		ID:             int(withdraw.WithdrawID),
		WithdrawNo:     withdraw.WithdrawNo.String(),
		CardNumber:     withdraw.CardNumber,
		WithdrawAmount: int(withdraw.WithdrawAmount),
		WithdrawTime:   withdraw.WithdrawTime.String(),
		CreatedAt:      withdraw.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      withdraw.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *withdrawRecordMapper) ToWithdrawsRecordTrashed(withdraws []*db.GetTrashedWithdrawsRow) []*record.WithdrawRecord {
	var withdrawRecords []*record.WithdrawRecord

	for _, withdraw := range withdraws {
		withdrawRecords = append(withdrawRecords, s.ToWithdrawRecordTrashed(withdraw))
	}

	return withdrawRecords
}

func (t *withdrawRecordMapper) ToWithdrawRecordMonthStatusSuccess(s *db.GetMonthWithdrawStatusSuccessRow) *record.WithdrawRecordMonthStatusSuccess {
	return &record.WithdrawRecordMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *withdrawRecordMapper) ToWithdrawRecordsMonthStatusSuccess(Withdraws []*db.GetMonthWithdrawStatusSuccessRow) []*record.WithdrawRecordMonthStatusSuccess {
	var WithdrawRecords []*record.WithdrawRecordMonthStatusSuccess

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.ToWithdrawRecordMonthStatusSuccess(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawRecordMapper) ToWithdrawRecordYearStatusSuccess(s *db.GetYearlyWithdrawStatusSuccessRow) *record.WithdrawRecordYearStatusSuccess {
	return &record.WithdrawRecordYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *withdrawRecordMapper) ToWithdrawRecordsYearStatusSuccess(Withdraws []*db.GetYearlyWithdrawStatusSuccessRow) []*record.WithdrawRecordYearStatusSuccess {
	var WithdrawRecords []*record.WithdrawRecordYearStatusSuccess

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.ToWithdrawRecordYearStatusSuccess(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawRecordMapper) ToWithdrawRecordMonthStatusFailed(s *db.GetMonthWithdrawStatusFailedRow) *record.WithdrawRecordMonthStatusFailed {
	return &record.WithdrawRecordMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *withdrawRecordMapper) ToWithdrawRecordsMonthStatusFailed(Withdraws []*db.GetMonthWithdrawStatusFailedRow) []*record.WithdrawRecordMonthStatusFailed {
	var WithdrawRecords []*record.WithdrawRecordMonthStatusFailed

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.ToWithdrawRecordMonthStatusFailed(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawRecordMapper) ToWithdrawRecordYearStatusFailed(s *db.GetYearlyWithdrawStatusFailedRow) *record.WithdrawRecordYearStatusFailed {
	return &record.WithdrawRecordYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *withdrawRecordMapper) ToWithdrawRecordsYearStatusFailed(Withdraws []*db.GetYearlyWithdrawStatusFailedRow) []*record.WithdrawRecordYearStatusFailed {
	var WithdrawRecords []*record.WithdrawRecordYearStatusFailed

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.ToWithdrawRecordYearStatusFailed(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawRecordMapper) ToWithdrawRecordMonthStatusSuccessCardNumber(s *db.GetMonthWithdrawStatusSuccessCardNumberRow) *record.WithdrawRecordMonthStatusSuccess {
	return &record.WithdrawRecordMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *withdrawRecordMapper) ToWithdrawRecordsMonthStatusSuccessCardNumber(Withdraws []*db.GetMonthWithdrawStatusSuccessCardNumberRow) []*record.WithdrawRecordMonthStatusSuccess {
	var WithdrawRecords []*record.WithdrawRecordMonthStatusSuccess

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.ToWithdrawRecordMonthStatusSuccessCardNumber(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawRecordMapper) ToWithdrawRecordYearStatusSuccessCardNumber(s *db.GetYearlyWithdrawStatusSuccessCardNumberRow) *record.WithdrawRecordYearStatusSuccess {
	return &record.WithdrawRecordYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *withdrawRecordMapper) ToWithdrawRecordsYearStatusSuccessCardNumber(Withdraws []*db.GetYearlyWithdrawStatusSuccessCardNumberRow) []*record.WithdrawRecordYearStatusSuccess {
	var WithdrawRecords []*record.WithdrawRecordYearStatusSuccess

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.ToWithdrawRecordYearStatusSuccessCardNumber(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawRecordMapper) ToWithdrawRecordMonthStatusFailedCardNumber(s *db.GetMonthWithdrawStatusFailedCardNumberRow) *record.WithdrawRecordMonthStatusFailed {
	return &record.WithdrawRecordMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *withdrawRecordMapper) ToWithdrawRecordsMonthStatusFailedCardNumber(Withdraws []*db.GetMonthWithdrawStatusFailedCardNumberRow) []*record.WithdrawRecordMonthStatusFailed {
	var WithdrawRecords []*record.WithdrawRecordMonthStatusFailed

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.ToWithdrawRecordMonthStatusFailedCardNumber(Withdraw))
	}

	return WithdrawRecords
}

func (t *withdrawRecordMapper) ToWithdrawRecordYearStatusFailedCardNumber(s *db.GetYearlyWithdrawStatusFailedCardNumberRow) *record.WithdrawRecordYearStatusFailed {
	return &record.WithdrawRecordYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *withdrawRecordMapper) ToWithdrawRecordsYearStatusFailedCardNumber(Withdraws []*db.GetYearlyWithdrawStatusFailedCardNumberRow) []*record.WithdrawRecordYearStatusFailed {
	var WithdrawRecords []*record.WithdrawRecordYearStatusFailed

	for _, Withdraw := range Withdraws {
		WithdrawRecords = append(WithdrawRecords, t.ToWithdrawRecordYearStatusFailedCardNumber(Withdraw))
	}

	return WithdrawRecords
}

func (r *withdrawRecordMapper) ToWithdrawAmountMonthly(ss *db.GetMonthlyWithdrawsRow) *record.WithdrawMonthlyAmount {
	return &record.WithdrawMonthlyAmount{
		Month:       ss.Month,
		TotalAmount: int(ss.TotalWithdrawAmount),
	}
}

func (s *withdrawRecordMapper) ToWithdrawsAmountMonthly(ss []*db.GetMonthlyWithdrawsRow) []*record.WithdrawMonthlyAmount {
	var withdrawRecords []*record.WithdrawMonthlyAmount

	for _, withdraw := range ss {
		withdrawRecords = append(withdrawRecords, s.ToWithdrawAmountMonthly(withdraw))
	}

	return withdrawRecords
}

func (r *withdrawRecordMapper) ToWithdrawAmountYearly(ss *db.GetYearlyWithdrawsRow) *record.WithdrawYearlyAmount {
	return &record.WithdrawYearlyAmount{
		Year:        ss.Year,
		TotalAmount: int(ss.TotalWithdrawAmount),
	}
}

func (s *withdrawRecordMapper) ToWithdrawsAmountYearly(ss []*db.GetYearlyWithdrawsRow) []*record.WithdrawYearlyAmount {
	var withdrawRecords []*record.WithdrawYearlyAmount

	for _, withdraw := range ss {
		withdrawRecords = append(withdrawRecords, s.ToWithdrawAmountYearly(withdraw))
	}

	return withdrawRecords
}

func (r *withdrawRecordMapper) ToWithdrawAmountMonthlyByCardNumber(ss *db.GetMonthlyWithdrawsByCardNumberRow) *record.WithdrawMonthlyAmount {
	return &record.WithdrawMonthlyAmount{
		Month:       ss.Month,
		TotalAmount: int(ss.TotalWithdrawAmount),
	}
}

func (s *withdrawRecordMapper) ToWithdrawsAmountMonthlyByCardNumber(ss []*db.GetMonthlyWithdrawsByCardNumberRow) []*record.WithdrawMonthlyAmount {
	var withdrawRecords []*record.WithdrawMonthlyAmount

	for _, withdraw := range ss {
		withdrawRecords = append(withdrawRecords, s.ToWithdrawAmountMonthlyByCardNumber(withdraw))
	}

	return withdrawRecords
}

func (r *withdrawRecordMapper) ToWithdrawAmountYearlyByCardNumber(ss *db.GetYearlyWithdrawsByCardNumberRow) *record.WithdrawYearlyAmount {
	return &record.WithdrawYearlyAmount{
		Year:        ss.Year,
		TotalAmount: int(ss.TotalWithdrawAmount),
	}
}

func (s *withdrawRecordMapper) ToWithdrawsAmountYearlyByCardNumber(ss []*db.GetYearlyWithdrawsByCardNumberRow) []*record.WithdrawYearlyAmount {
	var withdrawRecords []*record.WithdrawYearlyAmount

	for _, withdraw := range ss {
		withdrawRecords = append(withdrawRecords, s.ToWithdrawAmountYearlyByCardNumber(withdraw))
	}

	return withdrawRecords
}
