package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"database/sql"
	"errors"
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

func (r *transferRepository) FindAll(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransfers(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no transfers found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to retrieve transfers (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransfersRecordAll(res), &totalCount, nil
}

func (r *transferRepository) FindByActive(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTransfers(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no active transfers found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to find active transfers (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransfersRecordActive(res), &totalCount, nil
}

func (r *transferRepository) FindByTrashed(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTransfers(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no trashed transfers found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to find trashed transfers (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransfersRecordTrashed(res), &totalCount, nil
}

func (r *transferRepository) FindById(id int) (*record.TransferRecord, error) {
	transfer, err := r.db.GetTransferByID(r.ctx, int32(id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transfer not found with ID: %d", id)
		}
		return nil, fmt.Errorf("failed to find transfer by ID %d: %w", id, err)
	}

	return r.mapping.ToTransferRecord(transfer), nil
}

func (r *transferRepository) GetMonthTransferStatusSuccess(req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no transfer data for success status found for year %d, month %d", req.Year, req.Month)
		}
		return nil, fmt.Errorf("failed to get monthly transfer status success for year %d, month %d: %w", req.Year, req.Month, err)
	}

	so := r.mapping.ToTransferRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *transferRepository) GetYearlyTransferStatusSuccess(year int) ([]*record.TransferRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransferStatusSuccess(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transfer success found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly transfer status success for year %d: %w", year, err)
	}

	so := r.mapping.ToTransferRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *transferRepository) GetMonthTransferStatusFailed(req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no month transfer failed found for year %d and month %d", req.Year, req.Month)
		}
		return nil, fmt.Errorf("failed to get month transfer status failed for year %d and month %d: %w", req.Year, req.Month, err)
	}

	so := r.mapping.ToTransferRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *transferRepository) GetYearlyTransferStatusFailed(year int) ([]*record.TransferRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransferStatusFailed(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transfer failed found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly transfer status failed for year %d: %w", year, err)
	}

	so := r.mapping.ToTransferRecordsYearStatusFailed(res)

	return so, nil
}

func (r *transferRepository) GetMonthTransferStatusSuccessByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccessCardNumber(r.ctx, db.GetMonthTransferStatusSuccessCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no month transfer success found for year %d, month %d and card_number %s", req.Year, req.Month, req.CardNumber)
		}
		return nil, fmt.Errorf("failed to get month transfer status success for year %d, month %d and card_number %s: %w", req.Year, req.Month, req.CardNumber, err)
	}

	so := r.mapping.ToTransferRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transferRepository) GetYearlyTransferStatusSuccessByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransferStatusSuccessCardNumber(r.ctx, db.GetYearlyTransferStatusSuccessCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      int32(req.Year),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transfer success found for year %d and card_number %s", req.Year, req.CardNumber)
		}
		return nil, fmt.Errorf("failed to get yearly transfer status success for year %d and card_number %s: %w", req.Year, req.CardNumber, err)
	}

	so := r.mapping.ToTransferRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transferRepository) GetMonthTransferStatusFailedByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailedCardNumber(r.ctx, db.GetMonthTransferStatusFailedCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no month transfer failed found for year %d, month %d and card_number %s", req.Year, req.Month, req.CardNumber)
		}
		return nil, fmt.Errorf("failed to get month transfer status failed for year %d, month %d and card_number %s: %w", req.Year, req.Month, req.CardNumber, err)
	}

	so := r.mapping.ToTransferRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

func (r *transferRepository) GetYearlyTransferStatusFailedByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransferStatusFailedCardNumber(r.ctx, db.GetYearlyTransferStatusFailedCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      int32(req.Year),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transfer failed found for year %d and card_number %s", req.Year, req.CardNumber)
		}
		return nil, fmt.Errorf("failed to get yearly transfer status failed for year %d and card_number %s: %w", req.Year, req.CardNumber, err)
	}

	so := r.mapping.ToTransferRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

func (r *transferRepository) GetMonthlyTransferAmounts(year int) ([]*record.TransferMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmounts(r.ctx, yearStart)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transfer amounts found")
		}
		return nil, fmt.Errorf("failed to get monthly transfer amounts: %w", err)
	}

	return r.mapping.ToTransferMonthAmounts(res), nil
}

func (r *transferRepository) GetYearlyTransferAmounts(year int) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmounts(r.ctx, year)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transfer amounts found")
		}
		return nil, fmt.Errorf("failed to get yearly transfer amounts: %w", err)
	}
	return r.mapping.ToTransferYearAmounts(res), nil
}

func (r *transferRepository) GetMonthlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error) {
	res, err := r.db.GetMonthlyTransferAmountsBySenderCardNumber(r.ctx, db.GetMonthlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transfer amounts by sender card number %s and year %d found", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get monthly transfer amounts by sender card number %s and year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToTransferMonthAmountsSender(res), nil
}

func (r *transferRepository) GetMonthlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error) {
	res, err := r.db.GetMonthlyTransferAmountsByReceiverCardNumber(r.ctx, db.GetMonthlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: req.CardNumber,
		Column2:    time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transfer amounts by receiver card number %s and year %d found", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get monthly transfer amounts by receiver card number %s and year %d: %w", req.CardNumber, req.Year, err)
	}
	return r.mapping.ToTransferMonthAmountsReceiver(res), nil
}

func (r *transferRepository) GetYearlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountsBySenderCardNumber(r.ctx, db.GetYearlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      req.Year,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transfer amounts by sender card number %s and year %d found", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get yearly transfer amounts by sender card number %s and year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToTransferYearAmountsSender(res), nil
}

func (r *transferRepository) GetYearlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountsByReceiverCardNumber(r.ctx, db.GetYearlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transfer amounts by receiver card number %s and year %d found", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get yearly transfer amounts by receiver card number %s and year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToTransferYearAmountsReceiver(res), nil
}

func (r *transferRepository) FindTransferByTransferFrom(transfer_from string) ([]*record.TransferRecord, error) {
	res, err := r.db.GetTransfersBySourceCard(r.ctx, transfer_from)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no transfer found with transfer_from = %s", transfer_from)
		}
		return nil, fmt.Errorf("failed to find transfer by transfer from %s: %w", transfer_from, err)
	}

	return r.mapping.ToTransfersRecord(res), nil
}

func (r *transferRepository) FindTransferByTransferTo(transfer_to string) ([]*record.TransferRecord, error) {
	res, err := r.db.GetTransfersByDestinationCard(r.ctx, transfer_to)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no transfer found with transfer_to = %s", transfer_to)
		}
		return nil, fmt.Errorf("failed to find transfer by transfer to %s: %w", transfer_to, err)
	}
	return r.mapping.ToTransfersRecord(res), nil
}

func (r *transferRepository) CreateTransfer(request *requests.CreateTransferRequest) (*record.TransferRecord, error) {
	req := db.CreateTransferParams{
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.CreateTransfer(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid transfer data: %w", err)
		}
		return nil, fmt.Errorf("failed to create transfer: invalid or incomplete transfer data: %w", err)
	}

	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) UpdateTransfer(request *requests.UpdateTransferRequest) (*record.TransferRecord, error) {
	req := db.UpdateTransferParams{
		TransferID:     int32(*request.TransferID),
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.UpdateTransfer(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transfer ID %d not found for update", request.TransferID)
		}
		return nil, fmt.Errorf("failed to update transfer ID %d: transfer not found or invalid update data", request.TransferID)
	}

	return r.mapping.ToTransferRecord(res), nil

}

func (r *transferRepository) UpdateTransferAmount(request *requests.UpdateTransferAmountRequest) (*record.TransferRecord, error) {
	req := db.UpdateTransferAmountParams{
		TransferID:     int32(request.TransferID),
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.UpdateTransferAmount(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transfer ID %d not found for update", request.TransferID)
		}
		return nil, fmt.Errorf("failed to update transfer ID %d: transfer not found or invalid update data", request.TransferID)
	}

	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) UpdateTransferStatus(request *requests.UpdateTransferStatus) (*record.TransferRecord, error) {
	req := db.UpdateTransferStatusParams{
		TransferID: int32(request.TransferID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateTransferStatus(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transfer ID %d not found for update", request.TransferID)
		}
		return nil, fmt.Errorf("failed to update transfer ID %d: transfer not found or invalid update data", request.TransferID)
	}

	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) TrashedTransfer(transfer_id int) (*record.TransferRecord, error) {
	res, err := r.db.TrashTransfer(r.ctx, int32(transfer_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transfer ID %d not found or already trashed", transfer_id)
		}
		return nil, fmt.Errorf("failed to move transfer ID %d to trash: %w", transfer_id, err)
	}
	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) RestoreTransfer(transfer_id int) (*record.TransferRecord, error) {
	res, err := r.db.RestoreTransfer(r.ctx, int32(transfer_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transfer ID %d not found in trash", transfer_id)
		}
		return nil, fmt.Errorf("failed to restore transfer ID %d: %w", transfer_id, err)
	}
	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) DeleteTransferPermanent(transfer_id int) (bool, error) {
	err := r.db.DeleteTransferPermanently(r.ctx, int32(transfer_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("transfer ID %d not found", transfer_id)
		}
		return false, fmt.Errorf("failed to permanently delete transfer ID %d: %w", transfer_id, err)
	}
	return true, nil
}

func (r *transferRepository) RestoreAllTransfer() (bool, error) {
	err := r.db.RestoreAllTransfers(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no trashed transfers available to restore")
		}
		return false, fmt.Errorf("failed to restore trashed transfers: %w", err)
	}

	return true, nil
}

func (r *transferRepository) DeleteAllTransferPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTransfers(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no trashed transfers available to delete permanently")
		}
		return false, fmt.Errorf("failed to permanently delete transfers: %w", err)
	}

	return true, nil
}
