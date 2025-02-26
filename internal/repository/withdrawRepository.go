package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
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

func (r *withdrawRepository) FindAll(search string, page, pageSize int) ([]*record.WithdrawRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetWithdrawsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	withdraw, err := r.db.GetWithdraws(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get withdraws: %w", err)
	}

	var totalCount int
	if len(withdraw) > 0 {
		totalCount = int(withdraw[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToWithdrawsRecordALl(withdraw), totalCount, nil

}

func (r *withdrawRepository) FindAllByCardNumber(card_number string, search string, page, pageSize int) ([]*record.WithdrawRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetWithdrawsByCardNumberParams{
		CardNumber: card_number,
		Column2:    search,
		Limit:      int32(pageSize),
		Offset:     int32(offset),
	}

	withdraw, err := r.db.GetWithdrawsByCardNumber(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get withdraws by card number: %w", err)
	}

	var totalCount int
	if len(withdraw) > 0 {
		totalCount = int(withdraw[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToWithdrawsByCardNumberRecord(withdraw), totalCount, nil

}

func (r *withdrawRepository) FindById(id int) (*record.WithdrawRecord, error) {
	withdraw, err := r.db.GetWithdrawByID(r.ctx, int32(id))

	if err != nil {
		return nil, errors.New("failed get withdraw")
	}

	return r.mapping.ToWithdrawRecord(withdraw), nil
}

func (r *withdrawRepository) GetMonthWithdrawStatusSuccess(year int, month int) ([]*record.WithdrawRecordMonthStatusSuccess, error) {
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
		return nil, fmt.Errorf("failed to get month top-up status success for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusSuccess(year int) ([]*record.WithdrawRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyWithdrawStatusSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status success for year %d: %w", year, err)
	}

	so := r.mapping.ToWithdrawRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *withdrawRepository) GetMonthWithdrawStatusFailed(year int, month int) ([]*record.WithdrawRecordMonthStatusFailed, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
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
		return nil, fmt.Errorf("failed to get month top-up status failed for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusFailed(year int) ([]*record.WithdrawRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyWithdrawStatusFailed(r.ctx, int32(year))

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status failed for year %d: %w", year, err)
	}

	so := r.mapping.ToWithdrawRecordsYearStatusFailed(res)

	return so, nil
}

func (r *withdrawRepository) GetMonthWithdrawStatusSuccessByCardNumber(card_number string, year int, month int) ([]*record.WithdrawRecordMonthStatusSuccess, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusSuccessCardNumber(r.ctx, db.GetMonthWithdrawStatusSuccessCardNumberParams{
		CardNumber: card_number,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get month top-up status success for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusSuccessByCardNumber(card_number string, year int) ([]*record.WithdrawRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyWithdrawStatusSuccessCardNumber(r.ctx, db.GetYearlyWithdrawStatusSuccessCardNumberParams{
		CardNumber: card_number,
		Column2:    int32(year),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status success for year %d: %w", year, err)
	}

	so := r.mapping.ToWithdrawRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

func (r *withdrawRepository) GetMonthWithdrawStatusFailedByCardNumber(card_number string, year int, month int) ([]*record.WithdrawRecordMonthStatusFailed, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusFailedCardNumber(r.ctx, db.GetMonthWithdrawStatusFailedCardNumberParams{
		CardNumber: card_number,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get month top-up status failed for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusFailedByCardNumber(card_number string, year int) ([]*record.WithdrawRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyWithdrawStatusFailedCardNumber(r.ctx, db.GetYearlyWithdrawStatusFailedCardNumberParams{
		CardNumber: card_number,
		Column2:    int32(year),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status failed for year %d: %w", year, err)
	}

	so := r.mapping.ToWithdrawRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

func (r *withdrawRepository) GetMonthlyWithdraws(year int) ([]*record.WithdrawMonthlyAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdraws(r.ctx, yearStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly withdrawals for year %d: %w", year, err)
	}

	return r.mapping.ToWithdrawsAmountMonthly(res), nil

}

func (r *withdrawRepository) GetYearlyWithdraws(year int) ([]*record.WithdrawYearlyAmount, error) {
	res, err := r.db.GetYearlyWithdraws(r.ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly withdrawals for year %d: %w", year, err)
	}

	return r.mapping.ToWithdrawsAmountYearly(res), nil

}

func (r *withdrawRepository) GetMonthlyWithdrawsByCardNumber(cardNumber string, year int) ([]*record.WithdrawMonthlyAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawsByCardNumber(r.ctx, db.GetMonthlyWithdrawsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    yearStart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly withdrawals for card number %s and year %d: %w", cardNumber, year, err)
	}

	return r.mapping.ToWithdrawsAmountMonthlyByCardNumber(res), nil

}

func (r *withdrawRepository) GetYearlyWithdrawsByCardNumber(cardNumber string, year int) ([]*record.WithdrawYearlyAmount, error) {

	res, err := r.db.GetYearlyWithdrawsByCardNumber(r.ctx, db.GetYearlyWithdrawsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly withdrawals for card number %s and year %d: %w", cardNumber, year, err)
	}

	return r.mapping.ToWithdrawsAmountYearlyByCardNumber(res), nil
}

func (r *withdrawRepository) FindByActive(search string, page, pageSize int) ([]*record.WithdrawRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetActiveWithdrawsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveWithdraws(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find active withdraw: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToWithdrawsRecordActive(res), totalCount, nil
}

func (r *withdrawRepository) FindByTrashed(search string, page, pageSize int) ([]*record.WithdrawRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetTrashedWithdrawsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedWithdraws(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find trashed merchant: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToWithdrawsRecordTrashed(res), totalCount, nil
}

func (r *withdrawRepository) CreateWithdraw(request *requests.CreateWithdrawRequest) (*record.WithdrawRecord, error) {
	req := db.CreateWithdrawParams{
		CardNumber:     request.CardNumber,
		WithdrawAmount: int32(request.WithdrawAmount),
		WithdrawTime:   request.WithdrawTime,
	}

	res, err := r.db.CreateWithdraw(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update withdraw :%w", err)
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawRepository) UpdateWithdraw(request *requests.UpdateWithdrawRequest) (*record.WithdrawRecord, error) {
	req := db.UpdateWithdrawParams{
		WithdrawID:     int32(request.WithdrawID),
		CardNumber:     request.CardNumber,
		WithdrawAmount: int32(request.WithdrawAmount),
		WithdrawTime:   request.WithdrawTime,
	}

	err := r.db.UpdateWithdraw(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update withdraw: %w", err)
	}

	res, err := r.db.GetWithdrawByID(r.ctx, req.WithdrawID)

	if err != nil {
		return nil, fmt.Errorf("failed to find withdraw: %w", err)
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawRepository) UpdateWithdrawStatus(request *requests.UpdateWithdrawStatus) (*record.WithdrawRecord, error) {
	req := db.UpdateWithdrawStatusParams{
		WithdrawID: int32(request.WithdrawID),
		Status:     request.Status,
	}

	err := r.db.UpdateWithdrawStatus(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update Withdraw amount :%w", err)
	}

	res, err := r.db.GetWithdrawByID(r.ctx, req.WithdrawID)

	if err != nil {
		return nil, fmt.Errorf("failed to find Withdraw: %w", err)
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawRepository) TrashedWithdraw(WithdrawID int) (*record.WithdrawRecord, error) {
	err := r.db.TrashWithdraw(r.ctx, int32(WithdrawID))

	if err != nil {
		return nil, fmt.Errorf("failed to trash withdraw: %w", err)
	}

	merchant, err := r.db.GetTrashedWithdrawByID(r.ctx, int32(WithdrawID))

	if err != nil {
		return nil, fmt.Errorf("failed to find trashed by id topup: %w", err)
	}

	return r.mapping.ToWithdrawRecord(merchant), nil
}

func (r *withdrawRepository) RestoreWithdraw(WithdrawID int) (*record.WithdrawRecord, error) {
	err := r.db.RestoreWithdraw(r.ctx, int32(WithdrawID))

	if err != nil {
		return nil, fmt.Errorf("failed to restore withdraw: %w", err)
	}

	withdraw, err := r.db.GetWithdrawByID(r.ctx, int32(WithdrawID))

	if err != nil {
		return nil, fmt.Errorf("failed not found withdraw :%w", err)
	}

	return r.mapping.ToWithdrawRecord(withdraw), nil
}

func (r *withdrawRepository) DeleteWithdrawPermanent(WithdrawID int) (bool, error) {
	err := r.db.DeleteWithdrawPermanently(r.ctx, int32(WithdrawID))

	if err != nil {
		return false, fmt.Errorf("failed to delete withdraw: %w", err)
	}

	return true, nil
}

func (r *withdrawRepository) RestoreAllWithdraw() (bool, error) {
	err := r.db.RestoreAllWithdraws(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all withdraws: %w", err)
	}

	return true, nil
}

func (r *withdrawRepository) DeleteAllWithdrawPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentWithdraws(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all withdraws permanently: %w", err)
	}

	return true, nil
}
