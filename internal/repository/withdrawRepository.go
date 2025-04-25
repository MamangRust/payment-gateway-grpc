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

type withdrawRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.WithdrawRecordMapping
}

func NewWithdrawRepository(db *db.Queries, ctx context.Context, mapping recordmapper.WithdrawRecordMapping) *withdrawRepository {
	return &withdrawRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *withdrawRepository) FindAll(req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	withdraw, err := r.db.GetWithdraws(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no withdraws found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}

		return nil, nil, fmt.Errorf("failed to retrieve withdraws: invalid pagination (page %d, size %d) or search criteria '%s'", req.Page, req.PageSize, req.Search)
	}

	var totalCount int
	if len(withdraw) > 0 {
		totalCount = int(withdraw[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToWithdrawsRecordALl(withdraw), &totalCount, nil

}

func (r *withdrawRepository) FindByActive(req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveWithdraws(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no active withdraws found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}

		return nil, nil, fmt.Errorf("failed to find active withdraws: invalid parameters (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToWithdrawsRecordActive(res), &totalCount, nil
}

func (r *withdrawRepository) FindByTrashed(req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedWithdraws(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no trashed withdraws found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}

		return nil, nil, fmt.Errorf("failed to find trashed withdraws: invalid parameters (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToWithdrawsRecordTrashed(res), &totalCount, nil
}

func (r *withdrawRepository) FindAllByCardNumber(req *requests.FindAllWithdrawCardNumber) ([]*record.WithdrawRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	withdraw, err := r.db.GetWithdrawsByCardNumber(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no withdraws found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}

		return nil, nil, fmt.Errorf("failed to retrieve withdraws: invalid pagination (page %d, size %d) or search criteria '%s'", req.Page, req.PageSize, req.Search)
	}
	var totalCount int
	if len(withdraw) > 0 {
		totalCount = int(withdraw[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToWithdrawsByCardNumberRecord(withdraw), &totalCount, nil

}

func (r *withdrawRepository) FindById(id int) (*record.WithdrawRecord, error) {
	withdraw, err := r.db.GetWithdrawByID(r.ctx, int32(id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("withdraw not found with ID: %d", id)
		}
		return nil, fmt.Errorf("failed to find withdraw with ID %d: %w", id, err)
	}

	return r.mapping.ToWithdrawRecord(withdraw), nil
}

func (r *withdrawRepository) GetMonthWithdrawStatusSuccess(req *requests.MonthStatusWithdraw) ([]*record.WithdrawRecordMonthStatusSuccess, error) {
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusSuccess(r.ctx, db.GetMonthWithdrawStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no month withdraw success data found for year %d and month %d", year, month)
		}
		return nil, fmt.Errorf("failed to get month withdraw status success for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusSuccess(year int) ([]*record.WithdrawRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyWithdrawStatusSuccess(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly withdraw success data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly withdraw status success for year %d: %w", year, err)
	}

	so := r.mapping.ToWithdrawRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *withdrawRepository) GetMonthWithdrawStatusFailed(req *requests.MonthStatusWithdraw) ([]*record.WithdrawRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusFailed(r.ctx, db.GetMonthWithdrawStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no month withdraw failed data found for year %d and month %d", req.Year, req.Month)
		}
		return nil, fmt.Errorf("failed to get month withdraw status failed for year %d and month %d: %w", req.Year, req.Month, err)
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusFailed(year int) ([]*record.WithdrawRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyWithdrawStatusFailed(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly withdraw failed data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly withdraw status failed for year %d: %w", year, err)
	}

	so := r.mapping.ToWithdrawRecordsYearStatusFailed(res)

	return so, nil
}

func (r *withdrawRepository) GetMonthWithdrawStatusSuccessByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*record.WithdrawRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusSuccessCardNumber(r.ctx, db.GetMonthWithdrawStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no month withdraw success data found for year %d, month %d and card_number %s", req.Year, req.Month, req.CardNumber)
		}
		return nil, fmt.Errorf("failed to get month withdraw status success for year %d, month %d and card_number %s: %w", req.Year, req.Month, req.CardNumber, err)
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*record.WithdrawRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyWithdrawStatusSuccessCardNumber(r.ctx, db.GetYearlyWithdrawStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly withdraw success data found for year %d and card_number %s", req.Year, req.CardNumber)
		}
		return nil, fmt.Errorf("failed to get yearly withdraw status success for year %d and card_number %s: %w", req.Year, req.CardNumber, err)
	}

	so := r.mapping.ToWithdrawRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

func (r *withdrawRepository) GetMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*record.WithdrawRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusFailedCardNumber(r.ctx, db.GetMonthWithdrawStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no month withdraw failed data found for year %d, month %d and card_number %s", req.Year, req.Month, req.CardNumber)
		}
		return nil, fmt.Errorf("failed to get month withdraw status failed for year %d, month %d and card_number %s: %w", req.Year, req.Month, req.CardNumber, err)
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*record.WithdrawRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyWithdrawStatusFailedCardNumber(r.ctx, db.GetYearlyWithdrawStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly withdraw failed data found for year %d and card_number %s", req.Year, req.CardNumber)
		}
		return nil, fmt.Errorf("failed to get yearly withdraw status failed for year %d and card_number %s: %w", req.Year, req.CardNumber, err)
	}

	so := r.mapping.ToWithdrawRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

func (r *withdrawRepository) GetMonthlyWithdraws(year int) ([]*record.WithdrawMonthlyAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdraws(r.ctx, yearStart)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly withdrawals found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get monthly withdrawals for year %d: %w", year, err)
	}

	return r.mapping.ToWithdrawsAmountMonthly(res), nil

}

func (r *withdrawRepository) GetYearlyWithdraws(year int) ([]*record.WithdrawYearlyAmount, error) {
	res, err := r.db.GetYearlyWithdraws(r.ctx, year)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly withdrawals found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly withdrawals for year %d: %w", year, err)
	}

	return r.mapping.ToWithdrawsAmountYearly(res), nil

}

func (r *withdrawRepository) GetMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*record.WithdrawMonthlyAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawsByCardNumber(r.ctx, db.GetMonthlyWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly withdrawals found for card number %s and year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get monthly withdrawals for card number %s and year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToWithdrawsAmountMonthlyByCardNumber(res), nil

}

func (r *withdrawRepository) GetYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*record.WithdrawYearlyAmount, error) {
	res, err := r.db.GetYearlyWithdrawsByCardNumber(r.ctx, db.GetYearlyWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly withdrawals found for card number %s and year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get yearly withdrawals for card number %s and year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToWithdrawsAmountYearlyByCardNumber(res), nil
}

func (r *withdrawRepository) CreateWithdraw(request *requests.CreateWithdrawRequest) (*record.WithdrawRecord, error) {
	req := db.CreateWithdrawParams{
		CardNumber:     request.CardNumber,
		WithdrawAmount: int32(request.WithdrawAmount),
		WithdrawTime:   request.WithdrawTime,
	}

	res, err := r.db.CreateWithdraw(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to create withdraw: invalid or incomplete withdraw data")
		}
		return nil, fmt.Errorf("failed to create withdraw: invalid or incomplete withdraw data")
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawRepository) UpdateWithdraw(request *requests.UpdateWithdrawRequest) (*record.WithdrawRecord, error) {
	req := db.UpdateWithdrawParams{
		WithdrawID:     int32(*request.WithdrawID),
		CardNumber:     request.CardNumber,
		WithdrawAmount: int32(request.WithdrawAmount),
		WithdrawTime:   request.WithdrawTime,
	}

	res, err := r.db.UpdateWithdraw(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to update withdraw ID %d: withdraw not found or invalid update data", request.WithdrawID)
		}

		return nil, fmt.Errorf("failed to update withdraw ID %d: withdraw not found or invalid update data", request.WithdrawID)
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawRepository) UpdateWithdrawStatus(request *requests.UpdateWithdrawStatus) (*record.WithdrawRecord, error) {
	req := db.UpdateWithdrawStatusParams{
		WithdrawID: int32(request.WithdrawID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateWithdrawStatus(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to update withdraw status ID %d: withdraw not found or invalid update data", request.WithdrawID)
		}

		return nil, fmt.Errorf("failed to update withdraw status ID %d: withdraw not found or invalid update data", request.WithdrawID)
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawRepository) TrashedWithdraw(withdraw_id int) (*record.WithdrawRecord, error) {
	res, err := r.db.TrashWithdraw(r.ctx, int32(withdraw_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to move withdraw ID %d to trash: withdraw not found or already trashed", withdraw_id)
		}

		return nil, fmt.Errorf("failed to move withdraw ID %d to trash: withdraw not found or already trashed", withdraw_id)
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawRepository) RestoreWithdraw(withdraw_id int) (*record.WithdrawRecord, error) {
	res, err := r.db.RestoreWithdraw(r.ctx, int32(withdraw_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to restore withdraw ID %d: withdraw not found in trash", withdraw_id)
		}

		return nil, fmt.Errorf("failed to restore withdraw ID %d: withdraw not found in trash", withdraw_id)
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawRepository) DeleteWithdrawPermanent(withdraw_id int) (bool, error) {
	err := r.db.DeleteWithdrawPermanently(r.ctx, int32(withdraw_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("failed to permanently delete withdraw ID %d: usewithdrawr not found", withdraw_id)
		}

		return false, fmt.Errorf("failed to permanently delete withdraw ID %d: withdraw not found", withdraw_id)
	}

	return true, nil
}

func (r *withdrawRepository) RestoreAllWithdraw() (bool, error) {
	err := r.db.RestoreAllWithdraws(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no trashed withdraw available to restore")
		}

		return false, fmt.Errorf("no trashed withdraw available to restore")
	}

	return true, nil
}

func (r *withdrawRepository) DeleteAllWithdrawPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentWithdraws(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("cannot permanently delete all withdraw: operation not allowed")
		}

		return false, fmt.Errorf("cannot permanently delete all withdraw: operation not allowed")
	}

	return true, nil
}
