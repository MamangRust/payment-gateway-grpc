package recordmapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
)

type saldoRecordMapper struct{}

func NewSaldoRecordMapper() *saldoRecordMapper {
	return &saldoRecordMapper{}
}

func (s *saldoRecordMapper) ToSaldoRecord(saldo *db.Saldo) *record.SaldoRecord {
	var deletedAt *string

	if saldo.DeletedAt.Valid {
		formatedDeletedAt := saldo.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.SaldoRecord{
		ID:             int(saldo.SaldoID),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   int(saldo.TotalBalance),
		WithdrawAmount: int(saldo.WithdrawAmount.Int32),
		WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02"),
		CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:      deletedAt,
	}
}

func (s *saldoRecordMapper) ToSaldosRecord(saldos []*db.Saldo) []*record.SaldoRecord {
	var saldoRecords []*record.SaldoRecord
	for _, saldo := range saldos {
		saldoRecords = append(saldoRecords, s.ToSaldoRecord(saldo))
	}
	return saldoRecords
}

func (s *saldoRecordMapper) ToSaldoRecordAll(saldo *db.GetSaldosRow) *record.SaldoRecord {
	var deletedAt *string

	if saldo.DeletedAt.Valid {
		formatedDeletedAt := saldo.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.SaldoRecord{
		ID:             int(saldo.SaldoID),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   int(saldo.TotalBalance),
		WithdrawAmount: int(saldo.WithdrawAmount.Int32),
		WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02"),
		CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:      deletedAt,
	}
}

func (s *saldoRecordMapper) ToSaldosRecordAll(saldos []*db.GetSaldosRow) []*record.SaldoRecord {
	var saldoRecords []*record.SaldoRecord
	for _, saldo := range saldos {
		saldoRecords = append(saldoRecords, s.ToSaldoRecordAll(saldo))
	}
	return saldoRecords
}

func (s *saldoRecordMapper) ToSaldoRecordActive(saldo *db.GetActiveSaldosRow) *record.SaldoRecord {
	var deletedAt *string

	if saldo.DeletedAt.Valid {
		formatedDeletedAt := saldo.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.SaldoRecord{
		ID:             int(saldo.SaldoID),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   int(saldo.TotalBalance),
		WithdrawAmount: int(saldo.WithdrawAmount.Int32),
		WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02"),
		CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:      deletedAt,
	}
}

func (s *saldoRecordMapper) ToSaldosRecordActive(saldos []*db.GetActiveSaldosRow) []*record.SaldoRecord {
	var saldoRecords []*record.SaldoRecord
	for _, saldo := range saldos {
		saldoRecords = append(saldoRecords, s.ToSaldoRecordActive(saldo))
	}
	return saldoRecords
}

func (s *saldoRecordMapper) ToSaldoRecordTrashed(saldo *db.GetTrashedSaldosRow) *record.SaldoRecord {
	var deletedAt *string

	if saldo.DeletedAt.Valid {
		formatedDeletedAt := saldo.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.SaldoRecord{
		ID:             int(saldo.SaldoID),
		CardNumber:     saldo.CardNumber,
		TotalBalance:   int(saldo.TotalBalance),
		WithdrawAmount: int(saldo.WithdrawAmount.Int32),
		WithdrawTime:   saldo.WithdrawTime.Time.Format("2006-01-02"),
		CreatedAt:      saldo.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:      saldo.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:      deletedAt,
	}
}

func (s *saldoRecordMapper) ToSaldosRecordTrashed(saldos []*db.GetTrashedSaldosRow) []*record.SaldoRecord {
	var saldoRecords []*record.SaldoRecord
	for _, saldo := range saldos {
		saldoRecords = append(saldoRecords, s.ToSaldoRecordTrashed(saldo))
	}
	return saldoRecords
}

func (s *saldoRecordMapper) ToSaldoMonthTotalBalance(ss *db.GetMonthlyTotalSaldoBalanceRow) *record.SaldoMonthTotalBalance {
	totalBalance := 0
	if ss.TotalBalance != 0 {
		totalBalance = int(ss.TotalBalance)
	}

	return &record.SaldoMonthTotalBalance{
		Month:        ss.Month,
		Year:         ss.Year,
		TotalBalance: totalBalance,
	}
}

func (s *saldoRecordMapper) ToSaldoMonthTotalBalances(ss []*db.GetMonthlyTotalSaldoBalanceRow) []*record.SaldoMonthTotalBalance {
	var saldoRecords []*record.SaldoMonthTotalBalance
	for _, saldo := range ss {
		saldoRecords = append(saldoRecords, s.ToSaldoMonthTotalBalance(saldo))
	}
	return saldoRecords
}

func (s *saldoRecordMapper) ToSaldoYearTotalBalance(ss *db.GetYearlyTotalSaldoBalancesRow) *record.SaldoYearTotalBalance {
	totalBalance := 0
	if ss.TotalBalance != 0 {
		totalBalance = int(ss.TotalBalance)
	}

	return &record.SaldoYearTotalBalance{
		Year:         ss.Year,
		TotalBalance: int(totalBalance),
	}
}

func (s *saldoRecordMapper) ToSaldoYearTotalBalances(ss []*db.GetYearlyTotalSaldoBalancesRow) []*record.SaldoYearTotalBalance {
	var saldoRecords []*record.SaldoYearTotalBalance
	for _, saldo := range ss {
		saldoRecords = append(saldoRecords, s.ToSaldoYearTotalBalance(saldo))
	}
	return saldoRecords
}

func (s *saldoRecordMapper) ToSaldoMonthBalance(ss *db.GetMonthlySaldoBalancesRow) *record.SaldoMonthSaldoBalance {
	return &record.SaldoMonthSaldoBalance{
		Month:        ss.Month,
		TotalBalance: int(ss.TotalBalance),
	}
}

func (s *saldoRecordMapper) ToSaldoMonthBalances(ss []*db.GetMonthlySaldoBalancesRow) []*record.SaldoMonthSaldoBalance {
	var saldoRecords []*record.SaldoMonthSaldoBalance
	for _, saldo := range ss {
		saldoRecords = append(saldoRecords, s.ToSaldoMonthBalance(saldo))
	}
	return saldoRecords
}

func (s *saldoRecordMapper) ToSaldoYearSaldoBalance(ss *db.GetYearlySaldoBalancesRow) *record.SaldoYearSaldoBalance {
	return &record.SaldoYearSaldoBalance{
		Year:         ss.Year,
		TotalBalance: int(ss.TotalBalance),
	}
}

func (s *saldoRecordMapper) ToSaldoYearSaldoBalances(ss []*db.GetYearlySaldoBalancesRow) []*record.SaldoYearSaldoBalance {
	var saldoRecords []*record.SaldoYearSaldoBalance
	for _, saldo := range ss {
		saldoRecords = append(saldoRecords, s.ToSaldoYearSaldoBalance(saldo))
	}
	return saldoRecords
}
