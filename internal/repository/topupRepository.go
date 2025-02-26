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

type topupRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TopupRecordMapping
}

func NewTopupRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TopupRecordMapping) *topupRepository {
	return &topupRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *topupRepository) FindAllTopups(search string, page, pageSize int) ([]*record.TopupRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetTopupsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTopups(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find topups: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTopupRecordsAll(res), totalCount, nil
}

func (r *topupRepository) FindAllTopupByCardNumber(card_number string, search string, page, pageSize int) ([]*record.TopupRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetTopupsByCardNumberParams{
		CardNumber: card_number,
		Column2:    search,
		Limit:      int32(pageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetTopupsByCardNumber(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find topups: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTopupByCardNumberRecords(res), totalCount, nil
}

func (r *topupRepository) FindById(topup_id int) (*record.TopupRecord, error) {
	res, err := r.db.GetTopupByID(r.ctx, int32(topup_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find topup: %w", err)
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupRepository) GetMonthTopupStatusSuccess(year int, month int) ([]*record.TopupRecordMonthStatusSuccess, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusSuccess(r.ctx, db.GetMonthTopupStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get month top-up status success for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToTopupRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *topupRepository) GetYearlyTopupStatusSuccess(year int) ([]*record.TopupRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTopupStatusSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status success for year %d: %w", year, err)
	}

	so := r.mapping.ToTopupRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *topupRepository) GetMonthTopupStatusFailed(year int, month int) ([]*record.TopupRecordMonthStatusFailed, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusFailed(r.ctx, db.GetMonthTopupStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get month top-up status failed for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToTopupRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *topupRepository) GetYearlyTopupStatusFailed(year int) ([]*record.TopupRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTopupStatusFailed(r.ctx, int32(year))

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status failed for year %d: %w", year, err)
	}

	so := r.mapping.ToTopupRecordsYearStatusFailed(res)

	return so, nil
}

func (r *topupRepository) GetMonthTopupStatusSuccessByCardNumber(card_number string, year int, month int) ([]*record.TopupRecordMonthStatusSuccess, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusSuccessCardNumber(r.ctx, db.GetMonthTopupStatusSuccessCardNumberParams{
		CardNumber: card_number,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get month top-up status success for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToTopupRecordsMonthStatusSuccessByCardNumber(res)

	return so, nil
}

func (r *topupRepository) GetYearlyTopupStatusSuccessByCardNumber(card_number string, year int) ([]*record.TopupRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTopupStatusSuccessCardNumber(r.ctx, db.GetYearlyTopupStatusSuccessCardNumberParams{
		CardNumber: card_number,
		Column2:    int32(year),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status success for year %d: %w", year, err)
	}

	so := r.mapping.ToTopupRecordsYearStatusSuccessByCardNumber(res)

	return so, nil
}

func (r *topupRepository) GetMonthTopupStatusFailedByCardNumber(card_number string, year int, month int) ([]*record.TopupRecordMonthStatusFailed, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusFailedCardNumber(r.ctx, db.GetMonthTopupStatusFailedCardNumberParams{
		CardNumber: card_number,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get month top-up status failed for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToTopupRecordsMonthStatusFailedByCardNumber(res)

	return so, nil
}

func (r *topupRepository) GetYearlyTopupStatusFailedByCardNumber(card_number string, year int) ([]*record.TopupRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTopupStatusFailedCardNumber(r.ctx, db.GetYearlyTopupStatusFailedCardNumberParams{
		CardNumber: card_number,
		Column2:    int32(year),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status failed for year %d: %w", year, err)
	}

	so := r.mapping.ToTopupRecordsYearStatusFailedByCardNumber(res)

	return so, nil
}

func (r *topupRepository) GetMonthlyTopupMethods(year int) ([]*record.TopupMonthMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupMethods(r.ctx, yearStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly topup methods: %w", err)
	}

	return r.mapping.ToTopupMonthlyMethods(res), nil
}

func (r *topupRepository) GetYearlyTopupMethods(year int) ([]*record.TopupYearlyMethod, error) {
	res, err := r.db.GetYearlyTopupMethods(r.ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly topup methods: %w", err)
	}

	return r.mapping.ToTopupYearlyMethods(res), nil
}

func (r *topupRepository) GetMonthlyTopupAmounts(year int) ([]*record.TopupMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmounts(r.ctx, yearStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly topup amounts: %w", err)
	}

	return r.mapping.ToTopupMonthlyAmounts(res), nil
}

func (r *topupRepository) GetYearlyTopupAmounts(year int) ([]*record.TopupYearlyAmount, error) {
	res, err := r.db.GetYearlyTopupAmounts(r.ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly topup amounts: %w", err)
	}

	return r.mapping.ToTopupYearlyAmounts(res), nil
}

func (r *topupRepository) GetMonthlyTopupMethodsByCardNumber(card_number string, year int) ([]*record.TopupMonthMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupMethodsByCardNumber(r.ctx, db.GetMonthlyTopupMethodsByCardNumberParams{
		CardNumber: card_number,
		Column2:    yearStart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly topup methods by card number: %w", err)
	}

	return r.mapping.ToTopupMonthlyMethodsByCardNumber(res), nil
}

func (r *topupRepository) GetYearlyTopupMethodsByCardNumber(card_number string, year int) ([]*record.TopupYearlyMethod, error) {
	res, err := r.db.GetYearlyTopupMethodsByCardNumber(r.ctx, db.GetYearlyTopupMethodsByCardNumberParams{
		CardNumber: card_number,
		Column2:    year,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly topup methods by card number: %w", err)
	}

	return r.mapping.ToTopupYearlyMethodsByCardNumber(res), nil
}

func (r *topupRepository) GetMonthlyTopupAmountsByCardNumber(card_number string, year int) ([]*record.TopupMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmountsByCardNumber(r.ctx, db.GetMonthlyTopupAmountsByCardNumberParams{
		CardNumber: card_number,
		Column2:    yearStart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly topup amounts by card number: %w", err)
	}

	return r.mapping.ToTopupMonthlyAmountsByCardNumber(res), nil
}

func (r *topupRepository) GetYearlyTopupAmountsByCardNumber(card_number string, year int) ([]*record.TopupYearlyAmount, error) {
	res, err := r.db.GetYearlyTopupAmountsByCardNumber(r.ctx, db.GetYearlyTopupAmountsByCardNumberParams{
		CardNumber: card_number,
		Column2:    year,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly topup amounts by card number: %w", err)
	}

	return r.mapping.ToTopupYearlyAmountsByCardNumber(res), nil
}

func (r *topupRepository) FindByActive(search string, page, pageSize int) ([]*record.TopupRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetActiveTopupsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTopups(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find active merchant: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTopupRecordsActive(res), totalCount, nil
}

func (r *topupRepository) FindByTrashed(search string, page, pageSize int) ([]*record.TopupRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetTrashedTopupsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTopups(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find trashed merchant: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTopupRecordsTrashed(res), totalCount, nil
}

func (r *topupRepository) CreateTopup(request *requests.CreateTopupRequest) (*record.TopupRecord, error) {
	req := db.CreateTopupParams{
		CardNumber:  request.CardNumber,
		TopupAmount: int32(request.TopupAmount),
		TopupMethod: request.TopupMethod,
	}

	res, err := r.db.CreateTopup(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to create topup: %w", err)
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupRepository) UpdateTopup(request *requests.UpdateTopupRequest) (*record.TopupRecord, error) {
	req := db.UpdateTopupParams{
		TopupID:     int32(request.TopupID),
		CardNumber:  request.CardNumber,
		TopupAmount: int32(request.TopupAmount),
		TopupMethod: request.TopupMethod,
	}
	err := r.db.UpdateTopup(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update topup: %w", err)
	}

	res, err := r.db.GetTopupByID(r.ctx, req.TopupID)

	if err != nil {
		return nil, fmt.Errorf("failed to find topup: %w", err)
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupRepository) UpdateTopupAmount(request *requests.UpdateTopupAmount) (*record.TopupRecord, error) {
	req := db.UpdateTopupAmountParams{
		TopupID:     int32(request.TopupID),
		TopupAmount: int32(request.TopupAmount),
	}

	err := r.db.UpdateTopupAmount(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update topup amount :%w", err)
	}

	res, err := r.db.GetTopupByID(r.ctx, req.TopupID)

	if err != nil {
		return nil, fmt.Errorf("failed to find topup: %w", err)
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupRepository) UpdateTopupStatus(request *requests.UpdateTopupStatus) (*record.TopupRecord, error) {
	req := db.UpdateTopupStatusParams{
		TopupID: int32(request.TopupID),
		Status:  request.Status,
	}

	err := r.db.UpdateTopupStatus(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update topup amount :%w", err)
	}

	res, err := r.db.GetTopupByID(r.ctx, req.TopupID)

	if err != nil {
		return nil, fmt.Errorf("failed to find topup: %w", err)
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupRepository) TrashedTopup(topup_id int) (*record.TopupRecord, error) {
	err := r.db.TrashTopup(r.ctx, int32(topup_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash topup: %w", err)
	}

	merchant, err := r.db.GetTrashedTopupByID(r.ctx, int32(topup_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find trashed by id topup: %w", err)
	}

	return r.mapping.ToTopupRecord(merchant), nil
}

func (r *topupRepository) RestoreTopup(topup_id int) (*record.TopupRecord, error) {
	err := r.db.RestoreTopup(r.ctx, int32(topup_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore topup: %w", err)
	}

	topup, err := r.db.GetTopupByID(r.ctx, int32(topup_id))

	if err != nil {
		return nil, fmt.Errorf("failed not found topup :%w", err)
	}

	return r.mapping.ToTopupRecord(topup), nil
}

func (r *topupRepository) DeleteTopupPermanent(topup_id int) (bool, error) {
	err := r.db.DeleteTopupPermanently(r.ctx, int32(topup_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete topup permanently: %w", err)
	}

	return true, nil
}

func (r *topupRepository) RestoreAllTopup() (bool, error) {
	err := r.db.RestoreAllTopups(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all topups: %w", err)
	}

	return true, nil
}

func (r *topupRepository) DeleteAllTopupPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTopups(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all topups permanently: %w", err)
	}

	return true, nil
}
