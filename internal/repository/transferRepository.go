package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
	"time"
)

type transferRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransferRecordMapping
}

func NewTransferRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransferRecordMapping) *transferRepository {
	return &transferRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transferRepository) FindAll(search string, page, pageSize int) ([]*record.TransferRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetTransfersParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransfers(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find transfers: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransfersRecordAll(res), totalCount, nil
}

func (r *transferRepository) FindById(id int) (*record.TransferRecord, error) {
	transfer, err := r.db.GetTransferByID(r.ctx, int32(id))

	if err != nil {
		return nil, fmt.Errorf("failed to find by transfer: %w", err)
	}

	return r.mapping.ToTransferRecord(transfer), nil
}

func (r *transferRepository) GetMonthTransferStatusSuccess(year int, month int) ([]*record.TransferRecordMonthStatusSuccess, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccess(r.ctx, db.GetMonthTransferStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get month top-up status success for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToTransferRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *transferRepository) GetYearlyTransferStatusSuccess(year int) ([]*record.TransferRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransferStatusSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status success for year %d: %w", year, err)
	}

	so := r.mapping.ToTransferRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *transferRepository) GetMonthTransferStatusFailed(year int, month int) ([]*record.TransferRecordMonthStatusFailed, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailed(r.ctx, db.GetMonthTransferStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get month top-up status failed for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToTransferRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *transferRepository) GetYearlyTransferStatusFailed(year int) ([]*record.TransferRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransferStatusFailed(r.ctx, int32(year))

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status failed for year %d: %w", year, err)
	}

	so := r.mapping.ToTransferRecordsYearStatusFailed(res)

	return so, nil
}

func (r *transferRepository) GetMonthTransferStatusSuccessByCardNumber(card_number string, year int, month int) ([]*record.TransferRecordMonthStatusSuccess, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccessCardNumber(r.ctx, db.GetMonthTransferStatusSuccessCardNumberParams{
		TransferFrom: card_number,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get month top-up status success for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToTransferRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transferRepository) GetYearlyTransferStatusSuccessByCardNumber(card_number string, year int) ([]*record.TransferRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransferStatusSuccessCardNumber(r.ctx, db.GetYearlyTransferStatusSuccessCardNumberParams{
		TransferFrom: card_number,
		Column2:      int32(year),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status success for year %d: %w", year, err)
	}

	so := r.mapping.ToTransferRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transferRepository) GetMonthTransferStatusFailedByCardNumber(card_number string, year int, month int) ([]*record.TransferRecordMonthStatusFailed, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailedCardNumber(r.ctx, db.GetMonthTransferStatusFailedCardNumberParams{
		TransferFrom: card_number,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get month top-up status failed for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToTransferRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

func (r *transferRepository) GetYearlyTransferStatusFailedByCardNumber(card_number string, year int) ([]*record.TransferRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransferStatusFailedCardNumber(r.ctx, db.GetYearlyTransferStatusFailedCardNumberParams{
		TransferFrom: card_number,
		Column2:      int32(year),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status failed for year %d: %w", year, err)
	}

	so := r.mapping.ToTransferRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

func (r *transferRepository) GetMonthlyTransferAmounts(year int) ([]*record.TransferMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmounts(r.ctx, yearStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly transfer amounts: %w", err)
	}

	return r.mapping.ToTransferMonthAmounts(res), nil
}

func (r *transferRepository) GetYearlyTransferAmounts(year int) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmounts(r.ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly transfer amounts: %w", err)
	}

	return r.mapping.ToTransferYearAmounts(res), nil
}

func (r *transferRepository) GetMonthlyTransferAmountsBySenderCardNumber(cardNumber string, year int) ([]*record.TransferMonthAmount, error) {
	res, err := r.db.GetMonthlyTransferAmountsBySenderCardNumber(r.ctx, db.GetMonthlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: cardNumber,
		Column2:      time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly transfer amounts by sender card number: %w", err)
	}

	return r.mapping.ToTransferMonthAmountsSender(res), nil
}

func (r *transferRepository) GetMonthlyTransferAmountsByReceiverCardNumber(cardNumber string, year int) ([]*record.TransferMonthAmount, error) {
	res, err := r.db.GetMonthlyTransferAmountsByReceiverCardNumber(r.ctx, db.GetMonthlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: cardNumber,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly transfer amounts by receiver card number: %w", err)
	}

	return r.mapping.ToTransferMonthAmountsReceiver(res), nil
}

func (r *transferRepository) GetYearlyTransferAmountsBySenderCardNumber(cardNumber string, year int) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountsBySenderCardNumber(r.ctx, db.GetYearlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: cardNumber,
		Column2:      year,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly transfer amounts by sender card number: %w", err)
	}

	return r.mapping.ToTransferYearAmountsSender(res), nil
}

func (r *transferRepository) GetYearlyTransferAmountsByReceiverCardNumber(cardNumber string, year int) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountsByReceiverCardNumber(r.ctx, db.GetYearlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: cardNumber,
		Column2:    year,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly transfer amounts by receiver card number: %w", err)
	}

	return r.mapping.ToTransferYearAmountsReceiver(res), nil
}

func (r *transferRepository) FindTransferByTransferFrom(transfer_from string) ([]*record.TransferRecord, error) {
	res, err := r.db.GetTransfersBySourceCard(r.ctx, transfer_from)

	if err != nil {
		return nil, fmt.Errorf("failed to find transfer by transfer from: %w", err)
	}

	return r.mapping.ToTransfersRecord(res), nil
}

func (r *transferRepository) FindTransferByTransferTo(transfer_to string) ([]*record.TransferRecord, error) {
	res, err := r.db.GetTransfersByDestinationCard(r.ctx, transfer_to)

	if err != nil {
		return nil, fmt.Errorf("failed to find transfer by transfer to: %w", err)
	}

	return r.mapping.ToTransfersRecord(res), nil
}

func (r *transferRepository) FindByActive(search string, page, pageSize int) ([]*record.TransferRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetActiveTransfersParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTransfers(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find active merchant: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransfersRecordActive(res), totalCount, nil
}

func (r *transferRepository) FindByTrashed(search string, page, pageSize int) ([]*record.TransferRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetTrashedTransfersParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTransfers(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find trashed merchant: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransfersRecordTrashed(res), totalCount, nil
}

func (r *transferRepository) CreateTransfer(request *requests.CreateTransferRequest) (*record.TransferRecord, error) {
	req := db.CreateTransferParams{
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.CreateTransfer(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to create transfer: %w", err)
	}

	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) UpdateTransfer(request *requests.UpdateTransferRequest) (*record.TransferRecord, error) {
	req := db.UpdateTransferParams{
		TransferID:     int32(request.TransferID),
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	err := r.db.UpdateTransfer(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update transfer: %w", err)
	}

	res, err := r.db.GetTransferByID(r.ctx, int32(request.TransferID))

	if err != nil {
		return nil, fmt.Errorf("failed to find transfer: %w", err)
	}

	return r.mapping.ToTransferRecord(res), nil

}

func (r *transferRepository) UpdateTransferAmount(request *requests.UpdateTransferAmountRequest) (*record.TransferRecord, error) {
	req := db.UpdateTransferAmountParams{
		TransferID:     int32(request.TransferID),
		TransferAmount: int32(request.TransferAmount),
	}

	err := r.db.UpdateTransferAmount(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update transfer amount: %w", err)
	}

	res, err := r.db.GetTransferByID(r.ctx, int32(request.TransferID))

	if err != nil {
		return nil, fmt.Errorf("failed to find transfer: %w", err)
	}

	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) UpdateTransferStatus(request *requests.UpdateTransferStatus) (*record.TransferRecord, error) {
	req := db.UpdateTransferStatusParams{
		TransferID: int32(request.TransferID),
		Status:     request.Status,
	}

	err := r.db.UpdateTransferStatus(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update Transfer amount :%w", err)
	}

	res, err := r.db.GetTransferByID(r.ctx, req.TransferID)

	if err != nil {
		return nil, fmt.Errorf("failed to find Transfer: %w", err)
	}

	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) TrashedTransfer(transfer_id int) (*record.TransferRecord, error) {
	err := r.db.TrashTransfer(r.ctx, int32(transfer_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash transfer: %w", err)
	}

	merchant, err := r.db.GetTrashedTransferByID(r.ctx, int32(transfer_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find trashed by id transfer: %w", err)
	}

	return r.mapping.ToTransferRecord(merchant), nil
}

func (r *transferRepository) RestoreTransfer(transfer_id int) (*record.TransferRecord, error) {
	err := r.db.RestoreTransfer(r.ctx, int32(transfer_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore transfer: %w", err)
	}

	transfer, err := r.db.GetTransferByID(r.ctx, int32(transfer_id))

	if err != nil {
		return nil, fmt.Errorf("failed not found transfer :%w", err)
	}

	return r.mapping.ToTransferRecord(transfer), nil
}

func (r *transferRepository) DeleteTransferPermanent(topup_id int) (bool, error) {
	err := r.db.DeleteTransferPermanently(r.ctx, int32(topup_id))
	if err != nil {
		return false, fmt.Errorf("failed to delete transfer: %w", err)
	}
	return true, nil
}

func (r *transferRepository) RestoreAllTransfer() (bool, error) {
	err := r.db.RestoreAllTransfers(r.ctx)
	if err != nil {
		return false, fmt.Errorf("failed to restore all transfers: %w", err)
	}
	return true, nil
}

func (r *transferRepository) DeleteAllTransferPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTransfers(r.ctx)
	if err != nil {
		return false, fmt.Errorf("failed to delete all transfers permanently: %w", err)
	}
	return true, nil
}
