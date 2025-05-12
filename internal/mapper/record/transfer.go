package recordmapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
)

type transferRecordMapper struct {
}

func NewTransferRecordMapper() *transferRecordMapper {
	return &transferRecordMapper{}
}

func (t *transferRecordMapper) ToTransferRecord(transfer *db.Transfer) *record.TransferRecord {
	var deletedAt *string

	return &record.TransferRecord{
		ID:             int(transfer.TransferID),
		TransferNo:     transfer.TransferNo.String(),
		TransferFrom:   transfer.TransferFrom,
		TransferTo:     transfer.TransferTo,
		TransferAmount: int(transfer.TransferAmount),
		TransferTime:   transfer.TransferTime.String(),
		CreatedAt:      transfer.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      transfer.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *transferRecordMapper) ToTransfersRecord(transfers []*db.Transfer) []*record.TransferRecord {
	var transferRecords []*record.TransferRecord

	for _, transfer := range transfers {
		transferRecords = append(transferRecords, s.ToTransferRecord(transfer))
	}

	return transferRecords
}

func (t *transferRecordMapper) ToTransferRecordAll(transfer *db.GetTransfersRow) *record.TransferRecord {
	var deletedAt *string

	return &record.TransferRecord{
		ID:             int(transfer.TransferID),
		TransferNo:     transfer.TransferNo.String(),
		TransferFrom:   transfer.TransferFrom,
		TransferTo:     transfer.TransferTo,
		TransferAmount: int(transfer.TransferAmount),
		TransferTime:   transfer.TransferTime.String(),
		CreatedAt:      transfer.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      transfer.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *transferRecordMapper) ToTransfersRecordAll(transfers []*db.GetTransfersRow) []*record.TransferRecord {
	var transferRecords []*record.TransferRecord

	for _, transfer := range transfers {
		transferRecords = append(transferRecords, s.ToTransferRecordAll(transfer))
	}

	return transferRecords
}

func (t *transferRecordMapper) ToTransferRecordActive(transfer *db.GetActiveTransfersRow) *record.TransferRecord {
	var deletedAt *string

	return &record.TransferRecord{
		ID:             int(transfer.TransferID),
		TransferNo:     transfer.TransferNo.String(),
		TransferFrom:   transfer.TransferFrom,
		TransferTo:     transfer.TransferTo,
		TransferAmount: int(transfer.TransferAmount),
		TransferTime:   transfer.TransferTime.String(),
		CreatedAt:      transfer.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      transfer.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *transferRecordMapper) ToTransfersRecordActive(transfers []*db.GetActiveTransfersRow) []*record.TransferRecord {
	var transferRecords []*record.TransferRecord

	for _, transfer := range transfers {
		transferRecords = append(transferRecords, s.ToTransferRecordActive(transfer))
	}

	return transferRecords
}

func (t *transferRecordMapper) ToTransferRecordTrashed(transfer *db.GetTrashedTransfersRow) *record.TransferRecord {
	var deletedAt *string

	if transfer.DeletedAt.Valid {
		formatedDeletedAt := transfer.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.TransferRecord{
		ID:             int(transfer.TransferID),
		TransferNo:     transfer.TransferNo.String(),
		TransferFrom:   transfer.TransferFrom,
		TransferTo:     transfer.TransferTo,
		TransferAmount: int(transfer.TransferAmount),
		TransferTime:   transfer.TransferTime.String(),
		CreatedAt:      transfer.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      transfer.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *transferRecordMapper) ToTransfersRecordTrashed(transfers []*db.GetTrashedTransfersRow) []*record.TransferRecord {
	var transferRecords []*record.TransferRecord

	for _, transfer := range transfers {
		transferRecords = append(transferRecords, s.ToTransferRecordTrashed(transfer))
	}

	return transferRecords
}

func (t *transferRecordMapper) ToTransferRecordMonthStatusSuccess(s *db.GetMonthTransferStatusSuccessRow) *record.TransferRecordMonthStatusSuccess {
	return &record.TransferRecordMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *transferRecordMapper) ToTransferRecordsMonthStatusSuccess(Transfers []*db.GetMonthTransferStatusSuccessRow) []*record.TransferRecordMonthStatusSuccess {
	var TransferRecords []*record.TransferRecordMonthStatusSuccess

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.ToTransferRecordMonthStatusSuccess(Transfer))
	}

	return TransferRecords
}

func (t *transferRecordMapper) ToTransferRecordYearStatusSuccess(s *db.GetYearlyTransferStatusSuccessRow) *record.TransferRecordYearStatusSuccess {
	return &record.TransferRecordYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *transferRecordMapper) ToTransferRecordsYearStatusSuccess(Transfers []*db.GetYearlyTransferStatusSuccessRow) []*record.TransferRecordYearStatusSuccess {
	var TransferRecords []*record.TransferRecordYearStatusSuccess

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.ToTransferRecordYearStatusSuccess(Transfer))
	}

	return TransferRecords
}

func (t *transferRecordMapper) ToTransferRecordMonthStatusFailed(s *db.GetMonthTransferStatusFailedRow) *record.TransferRecordMonthStatusFailed {
	return &record.TransferRecordMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transferRecordMapper) ToTransferRecordsMonthStatusFailed(Transfers []*db.GetMonthTransferStatusFailedRow) []*record.TransferRecordMonthStatusFailed {
	var TransferRecords []*record.TransferRecordMonthStatusFailed

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.ToTransferRecordMonthStatusFailed(Transfer))
	}

	return TransferRecords
}

func (t *transferRecordMapper) ToTransferRecordYearStatusFailed(s *db.GetYearlyTransferStatusFailedRow) *record.TransferRecordYearStatusFailed {
	return &record.TransferRecordYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transferRecordMapper) ToTransferRecordsYearStatusFailed(Transfers []*db.GetYearlyTransferStatusFailedRow) []*record.TransferRecordYearStatusFailed {
	var TransferRecords []*record.TransferRecordYearStatusFailed

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.ToTransferRecordYearStatusFailed(Transfer))
	}

	return TransferRecords
}

func (t *transferRecordMapper) ToTransferRecordMonthStatusSuccessCardNumber(s *db.GetMonthTransferStatusSuccessCardNumberRow) *record.TransferRecordMonthStatusSuccess {
	return &record.TransferRecordMonthStatusSuccess{
		Year:         s.Year,
		Month:        s.Month,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *transferRecordMapper) ToTransferRecordsMonthStatusSuccessCardNumber(Transfers []*db.GetMonthTransferStatusSuccessCardNumberRow) []*record.TransferRecordMonthStatusSuccess {
	var TransferRecords []*record.TransferRecordMonthStatusSuccess

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.ToTransferRecordMonthStatusSuccessCardNumber(Transfer))
	}

	return TransferRecords
}

func (t *transferRecordMapper) ToTransferRecordYearStatusSuccessCardNumber(s *db.GetYearlyTransferStatusSuccessCardNumberRow) *record.TransferRecordYearStatusSuccess {
	return &record.TransferRecordYearStatusSuccess{
		Year:         s.Year,
		TotalSuccess: int(s.TotalSuccess),
		TotalAmount:  int(s.TotalAmount),
	}
}

func (t *transferRecordMapper) ToTransferRecordsYearStatusSuccessCardNumber(Transfers []*db.GetYearlyTransferStatusSuccessCardNumberRow) []*record.TransferRecordYearStatusSuccess {
	var TransferRecords []*record.TransferRecordYearStatusSuccess

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.ToTransferRecordYearStatusSuccessCardNumber(Transfer))
	}

	return TransferRecords
}

func (t *transferRecordMapper) ToTransferRecordMonthStatusFailedCardNumber(s *db.GetMonthTransferStatusFailedCardNumberRow) *record.TransferRecordMonthStatusFailed {
	return &record.TransferRecordMonthStatusFailed{
		Year:        s.Year,
		Month:       s.Month,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transferRecordMapper) ToTransferRecordsMonthStatusFailedCardNumber(Transfers []*db.GetMonthTransferStatusFailedCardNumberRow) []*record.TransferRecordMonthStatusFailed {
	var TransferRecords []*record.TransferRecordMonthStatusFailed

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.ToTransferRecordMonthStatusFailedCardNumber(Transfer))
	}

	return TransferRecords
}

func (t *transferRecordMapper) ToTransferRecordYearStatusFailedCardNumber(s *db.GetYearlyTransferStatusFailedCardNumberRow) *record.TransferRecordYearStatusFailed {
	return &record.TransferRecordYearStatusFailed{
		Year:        s.Year,
		TotalFailed: int(s.TotalFailed),
		TotalAmount: int(s.TotalAmount),
	}
}

func (t *transferRecordMapper) ToTransferRecordsYearStatusFailedCardNumber(Transfers []*db.GetYearlyTransferStatusFailedCardNumberRow) []*record.TransferRecordYearStatusFailed {
	var TransferRecords []*record.TransferRecordYearStatusFailed

	for _, Transfer := range Transfers {
		TransferRecords = append(TransferRecords, t.ToTransferRecordYearStatusFailedCardNumber(Transfer))
	}

	return TransferRecords
}

func (s *transferRecordMapper) ToTransferMonthAmount(ss *db.GetMonthlyTransferAmountsRow) *record.TransferMonthAmount {
	return &record.TransferMonthAmount{
		Month:       ss.Month,
		TotalAmount: int(ss.TotalTransferAmount),
	}
}

func (s *transferRecordMapper) ToTransferMonthAmounts(ss []*db.GetMonthlyTransferAmountsRow) []*record.TransferMonthAmount {
	var transferRecords []*record.TransferMonthAmount

	for _, transfer := range ss {
		transferRecords = append(transferRecords, s.ToTransferMonthAmount(transfer))
	}

	return transferRecords
}

func (s *transferRecordMapper) ToTransferYearAmount(ss *db.GetYearlyTransferAmountsRow) *record.TransferYearAmount {
	return &record.TransferYearAmount{
		Year:        ss.Year,
		TotalAmount: int(ss.TotalTransferAmount),
	}
}

func (s *transferRecordMapper) ToTransferYearAmounts(ss []*db.GetYearlyTransferAmountsRow) []*record.TransferYearAmount {
	var transferRecords []*record.TransferYearAmount

	for _, transfer := range ss {
		transferRecords = append(transferRecords, s.ToTransferYearAmount(transfer))
	}

	return transferRecords
}

func (s *transferRecordMapper) ToTransferMonthAmountSender(ss *db.GetMonthlyTransferAmountsBySenderCardNumberRow) *record.TransferMonthAmount {
	return &record.TransferMonthAmount{
		Month:       ss.Month,
		TotalAmount: int(ss.TotalTransferAmount),
	}
}

func (s *transferRecordMapper) ToTransferMonthAmountsSender(ss []*db.GetMonthlyTransferAmountsBySenderCardNumberRow) []*record.TransferMonthAmount {
	var transferRecords []*record.TransferMonthAmount

	for _, transfer := range ss {
		transferRecords = append(transferRecords, s.ToTransferMonthAmountSender(transfer))
	}

	return transferRecords
}

func (s *transferRecordMapper) ToTransferYearAmountSender(ss *db.GetYearlyTransferAmountsBySenderCardNumberRow) *record.TransferYearAmount {
	return &record.TransferYearAmount{
		Year:        ss.Year,
		TotalAmount: int(ss.TotalTransferAmount),
	}
}

func (s *transferRecordMapper) ToTransferYearAmountsSender(ss []*db.GetYearlyTransferAmountsBySenderCardNumberRow) []*record.TransferYearAmount {
	var transferRecords []*record.TransferYearAmount

	for _, transfer := range ss {
		transferRecords = append(transferRecords, s.ToTransferYearAmountSender(transfer))
	}

	return transferRecords
}

func (s *transferRecordMapper) ToTransferMonthAmountReceiver(ss *db.GetMonthlyTransferAmountsByReceiverCardNumberRow) *record.TransferMonthAmount {
	return &record.TransferMonthAmount{
		Month:       ss.Month,
		TotalAmount: int(ss.TotalTransferAmount),
	}
}

func (s *transferRecordMapper) ToTransferMonthAmountsReceiver(ss []*db.GetMonthlyTransferAmountsByReceiverCardNumberRow) []*record.TransferMonthAmount {
	var transferRecords []*record.TransferMonthAmount

	for _, transfer := range ss {
		transferRecords = append(transferRecords, s.ToTransferMonthAmountReceiver(transfer))
	}

	return transferRecords
}

func (s *transferRecordMapper) ToTransferYearAmountReceiver(ss *db.GetYearlyTransferAmountsByReceiverCardNumberRow) *record.TransferYearAmount {
	return &record.TransferYearAmount{
		Year:        ss.Year,
		TotalAmount: int(ss.TotalTransferAmount),
	}
}

func (s *transferRecordMapper) ToTransferYearAmountsReceiver(ss []*db.GetYearlyTransferAmountsByReceiverCardNumberRow) []*record.TransferYearAmount {
	var transferRecords []*record.TransferYearAmount

	for _, transfer := range ss {
		transferRecords = append(transferRecords, s.ToTransferYearAmountReceiver(transfer))
	}

	return transferRecords
}
