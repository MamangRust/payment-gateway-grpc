package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/withdraw_errors"
	"context"
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
		return nil, nil, withdraw_errors.ErrFindAllWithdrawsFailed
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
		return nil, nil, withdraw_errors.ErrFindActiveWithdrawsFailed
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
		return nil, nil, withdraw_errors.ErrFindTrashedWithdrawsFailed
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
		return nil, nil, withdraw_errors.ErrFindWithdrawsByCardNumberFailed
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
		return nil, withdraw_errors.ErrFindWithdrawByIdFailed
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
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusSuccessFailed
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusSuccess(year int) ([]*record.WithdrawRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyWithdrawStatusSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusSuccessFailed
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
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusFailedFailed
	}

	so := r.mapping.ToWithdrawRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusFailed(year int) ([]*record.WithdrawRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyWithdrawStatusFailed(r.ctx, int32(year))

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusFailedFailed
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
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusSuccessByCardFailed
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
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusSuccessByCardFailed
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
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusFailedByCardFailed
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
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusFailedByCardFailed
	}

	so := r.mapping.ToWithdrawRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

func (r *withdrawRepository) GetMonthlyWithdraws(year int) ([]*record.WithdrawMonthlyAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdraws(r.ctx, yearStart)

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthlyWithdrawsFailed
	}

	return r.mapping.ToWithdrawsAmountMonthly(res), nil

}

func (r *withdrawRepository) GetYearlyWithdraws(year int) ([]*record.WithdrawYearlyAmount, error) {
	res, err := r.db.GetYearlyWithdraws(r.ctx, year)

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawsFailed
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
		return nil, withdraw_errors.ErrGetMonthlyWithdrawsByCardFailed
	}

	return r.mapping.ToWithdrawsAmountMonthlyByCardNumber(res), nil

}

func (r *withdrawRepository) GetYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*record.WithdrawYearlyAmount, error) {
	res, err := r.db.GetYearlyWithdrawsByCardNumber(r.ctx, db.GetYearlyWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawsByCardFailed
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
		return nil, withdraw_errors.ErrCreateWithdrawFailed
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
		return nil, withdraw_errors.ErrUpdateWithdrawFailed
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
		return nil, withdraw_errors.ErrUpdateWithdrawStatusFailed
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawRepository) TrashedWithdraw(withdraw_id int) (*record.WithdrawRecord, error) {
	res, err := r.db.TrashWithdraw(r.ctx, int32(withdraw_id))

	if err != nil {
		return nil, withdraw_errors.ErrTrashedWithdrawFailed
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawRepository) RestoreWithdraw(withdraw_id int) (*record.WithdrawRecord, error) {
	res, err := r.db.RestoreWithdraw(r.ctx, int32(withdraw_id))

	if err != nil {
		return nil, withdraw_errors.ErrRestoreWithdrawFailed
	}

	return r.mapping.ToWithdrawRecord(res), nil
}

func (r *withdrawRepository) DeleteWithdrawPermanent(withdraw_id int) (bool, error) {
	err := r.db.DeleteWithdrawPermanently(r.ctx, int32(withdraw_id))

	if err != nil {
		return false, withdraw_errors.ErrDeleteWithdrawPermanentFailed
	}

	return true, nil
}

func (r *withdrawRepository) RestoreAllWithdraw() (bool, error) {
	err := r.db.RestoreAllWithdraws(r.ctx)

	if err != nil {
		return false, withdraw_errors.ErrRestoreAllWithdrawsFailed
	}

	return true, nil
}

func (r *withdrawRepository) DeleteAllWithdrawPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentWithdraws(r.ctx)

	if err != nil {
		return false, withdraw_errors.ErrDeleteAllWithdrawsPermanentFailed
	}

	return true, nil
}
