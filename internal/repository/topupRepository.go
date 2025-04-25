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

func (r *topupRepository) FindAllTopups(req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTopups(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no topups found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to retrieve topups (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTopupRecordsAll(res), &totalCount, nil
}

func (r *topupRepository) FindByActive(req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTopups(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no active topups found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to find active topups (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTopupRecordsActive(res), &totalCount, nil
}

func (r *topupRepository) FindByTrashed(req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTopups(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no trashed topups found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to find trashed topups (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTopupRecordsTrashed(res), &totalCount, nil
}

func (r *topupRepository) FindAllTopupByCardNumber(req *requests.FindAllTopupsByCardNumber) ([]*record.TopupRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTopupsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetTopupsByCardNumber(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no topups found matching the criteria (page %d, size %d, search '%s', card_number '%s')", req.Page, req.PageSize, req.Search, req.CardNumber)
		}
		return nil, nil, fmt.Errorf("failed to retrieve topups (page %d, size %d, search '%s', card_number '%s'): %w", req.Page, req.PageSize, req.Search, req.CardNumber, err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTopupByCardNumberRecords(res), &totalCount, nil
}

func (r *topupRepository) FindById(topup_id int) (*record.TopupRecord, error) {
	res, err := r.db.GetTopupByID(r.ctx, int32(topup_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("topup not found with ID: %d", topup_id)
		}
		return nil, fmt.Errorf("failed to get topup by ID %d: %w", topup_id, err)
	}
	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupRepository) GetMonthTopupStatusSuccess(req *requests.MonthTopupStatus) ([]*record.TopupRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no successful topup records found for year %d and month %d", req.Year, req.Month)
		}
		return nil, fmt.Errorf("failed to get successful monthly topup status for year %d and month %d: %w", req.Year, req.Month, err)
	}

	so := r.mapping.ToTopupRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *topupRepository) GetYearlyTopupStatusSuccess(year int) ([]*record.TopupRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTopupStatusSuccess(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no successful yearly topup records found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get successful yearly topup status for year %d: %w", year, err)
	}
	so := r.mapping.ToTopupRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *topupRepository) GetMonthTopupStatusFailed(req *requests.MonthTopupStatus) ([]*record.TopupRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no failed topup records found for year %d and month %d", req.Year, req.Month)
		}
		return nil, fmt.Errorf("failed to get failed monthly topup status for year %d and month %d: %w", req.Year, req.Month, err)
	}

	so := r.mapping.ToTopupRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *topupRepository) GetYearlyTopupStatusFailed(year int) ([]*record.TopupRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTopupStatusFailed(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no failed yearly topup records found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get failed yearly topup status for year %d: %w", year, err)
	}

	so := r.mapping.ToTopupRecordsYearStatusFailed(res)

	return so, nil
}

func (r *topupRepository) GetMonthTopupStatusSuccessByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*record.TopupRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusSuccessCardNumber(r.ctx, db.GetMonthTopupStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no successful monthly topup records for card number %s in %d-%02d", req.CardNumber, req.Year, req.Month)
		}
		return nil, fmt.Errorf("failed to get successful monthly topup by card number %s for %d-%02d: %w", req.CardNumber, req.Year, req.Month, err)
	}

	so := r.mapping.ToTopupRecordsMonthStatusSuccessByCardNumber(res)

	return so, nil
}

func (r *topupRepository) GetYearlyTopupStatusSuccessByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*record.TopupRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTopupStatusSuccessCardNumber(r.ctx, db.GetYearlyTopupStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no successful yearly topup records for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get successful yearly topup by card number %s for year %d: %w", req.CardNumber, req.Year, err)
	}

	so := r.mapping.ToTopupRecordsYearStatusSuccessByCardNumber(res)

	return so, nil
}

func (r *topupRepository) GetMonthTopupStatusFailedByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*record.TopupRecordMonthStatusFailed, error) {
	cardNumber := req.CardNumber
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusFailedCardNumber(r.ctx, db.GetMonthTopupStatusFailedCardNumberParams{
		CardNumber: cardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no failed monthly topup records for card number %s in %d-%02d", req.CardNumber, req.Year, req.Month)
		}
		return nil, fmt.Errorf("failed to get failed monthly topup by card number %s for %d-%02d: %w", req.CardNumber, req.Year, req.Month, err)
	}

	so := r.mapping.ToTopupRecordsMonthStatusFailedByCardNumber(res)

	return so, nil
}

func (r *topupRepository) GetYearlyTopupStatusFailedByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*record.TopupRecordYearStatusFailed, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetYearlyTopupStatusFailedCardNumber(r.ctx, db.GetYearlyTopupStatusFailedCardNumberParams{
		CardNumber: cardNumber,
		Column2:    int32(year),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no failed yearly topup records for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get failed yearly topup by card number %s for year %d: %w", req.CardNumber, req.Year, err)
	}

	so := r.mapping.ToTopupRecordsYearStatusFailedByCardNumber(res)

	return so, nil
}

func (r *topupRepository) GetMonthlyTopupMethods(year int) ([]*record.TopupMonthMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupMethods(r.ctx, yearStart)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly topup method data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get monthly topup methods for year %d: %w", year, err)
	}

	return r.mapping.ToTopupMonthlyMethods(res), nil
}

func (r *topupRepository) GetYearlyTopupMethods(year int) ([]*record.TopupYearlyMethod, error) {
	res, err := r.db.GetYearlyTopupMethods(r.ctx, year)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly topup method data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly topup methods for year %d: %w", year, err)
	}

	return r.mapping.ToTopupYearlyMethods(res), nil
}

func (r *topupRepository) GetMonthlyTopupAmounts(year int) ([]*record.TopupMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmounts(r.ctx, yearStart)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly topup amount data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get monthly topup amounts for year %d: %w", year, err)
	}

	return r.mapping.ToTopupMonthlyAmounts(res), nil
}

func (r *topupRepository) GetYearlyTopupAmounts(year int) ([]*record.TopupYearlyAmount, error) {
	res, err := r.db.GetYearlyTopupAmounts(r.ctx, year)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly topup amount data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly topup amounts for year %d: %w", year, err)
	}

	return r.mapping.ToTopupYearlyAmounts(res), nil
}

func (r *topupRepository) GetMonthlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupMonthMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupMethodsByCardNumber(r.ctx, db.GetMonthlyTopupMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly topup method data for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get monthly topup methods by card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToTopupMonthlyMethodsByCardNumber(res), nil
}

func (r *topupRepository) GetYearlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupYearlyMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyTopupMethodsByCardNumber(r.ctx, db.GetYearlyTopupMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly topup method data for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get yearly topup methods by card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToTopupYearlyMethodsByCardNumber(res), nil
}

func (r *topupRepository) GetMonthlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupMonthAmount, error) {
	year := req.Year
	cardNumber := req.CardNumber

	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmountsByCardNumber(r.ctx, db.GetMonthlyTopupAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly topup amount data for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get monthly topup amounts by card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToTopupMonthlyAmountsByCardNumber(res), nil
}

func (r *topupRepository) GetYearlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupYearlyAmount, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyTopupAmountsByCardNumber(r.ctx, db.GetYearlyTopupAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly topup amount data for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get yearly topup amounts by card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToTopupYearlyAmountsByCardNumber(res), nil
}

func (r *topupRepository) CreateTopup(request *requests.CreateTopupRequest) (*record.TopupRecord, error) {
	req := db.CreateTopupParams{
		CardNumber:  request.CardNumber,
		TopupAmount: int32(request.TopupAmount),
		TopupMethod: request.TopupMethod,
	}

	res, err := r.db.CreateTopup(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid topup data: %w", err)
		}
		return nil, fmt.Errorf("failed to create topup: invalid or incomplete topup data: %w", err)
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupRepository) UpdateTopup(request *requests.UpdateTopupRequest) (*record.TopupRecord, error) {
	req := db.UpdateTopupParams{
		TopupID:     int32(*request.TopupID),
		CardNumber:  request.CardNumber,
		TopupAmount: int32(request.TopupAmount),
		TopupMethod: request.TopupMethod,
	}

	res, err := r.db.UpdateTopup(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("topup ID %d not found for update", request.TopupID)
		}
		return nil, fmt.Errorf("failed to update topup ID %d: topup not found or invalid update data", request.TopupID)
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupRepository) UpdateTopupAmount(request *requests.UpdateTopupAmount) (*record.TopupRecord, error) {
	req := db.UpdateTopupAmountParams{
		TopupID:     int32(request.TopupID),
		TopupAmount: int32(request.TopupAmount),
	}

	res, err := r.db.UpdateTopupAmount(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("topup ID %d not found for update", request.TopupID)
		}
		return nil, fmt.Errorf("failed to update topup ID %d: topup not found or invalid update data", request.TopupID)
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupRepository) UpdateTopupStatus(request *requests.UpdateTopupStatus) (*record.TopupRecord, error) {
	req := db.UpdateTopupStatusParams{
		TopupID: int32(request.TopupID),
		Status:  request.Status,
	}

	res, err := r.db.UpdateTopupStatus(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("topup ID %d not found for update", request.TopupID)
		}
		return nil, fmt.Errorf("failed to update topup ID %d: topup not found or invalid update data", request.TopupID)
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupRepository) TrashedTopup(topup_id int) (*record.TopupRecord, error) {
	res, err := r.db.TrashTopup(r.ctx, int32(topup_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("topup ID %d not found or already trashed", topup_id)
		}
		return nil, fmt.Errorf("failed to trash topup ID %d: %w", topup_id, err)
	}
	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupRepository) RestoreTopup(topup_id int) (*record.TopupRecord, error) {
	res, err := r.db.RestoreTopup(r.ctx, int32(topup_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("topup ID %d not found in trash", topup_id)
		}
		return nil, fmt.Errorf("failed to restore topup ID %d: %w", topup_id, err)
	}
	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupRepository) DeleteTopupPermanent(topup_id int) (bool, error) {
	err := r.db.DeleteTopupPermanently(r.ctx, int32(topup_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("topup ID %d not found or already deleted", topup_id)
		}
		return false, fmt.Errorf("failed to permanently delete topup ID %d: %w", topup_id, err)
	}
	return true, nil
}

func (r *topupRepository) RestoreAllTopup() (bool, error) {
	err := r.db.RestoreAllTopups(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no trashed topup available to restore")
		}
		return false, fmt.Errorf("failed to restore trashed topup: %w", err)
	}

	return true, nil
}

func (r *topupRepository) DeleteAllTopupPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTopups(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no trashed topup available to delete permanently")
		}
		return false, fmt.Errorf("failed to permanently delete topup: %w", err)
	}

	return true, nil
}
