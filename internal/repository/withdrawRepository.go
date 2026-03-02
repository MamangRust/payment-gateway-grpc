package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/withdraw_errors"
	"context"
	"time"
)

type withdrawRepository struct {
	db *db.Queries
}

func NewWithdrawRepository(db *db.Queries) WithdrawRepository {
	return &withdrawRepository{
		db: db,
	}
}

func (r *withdrawRepository) FindAll(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetWithdrawsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	withdraw, err := r.db.GetWithdraws(ctx, reqDb)

	if err != nil {
		return nil, withdraw_errors.ErrFindAllWithdrawsFailed
	}

	return withdraw, nil

}

func (r *withdrawRepository) FindByActive(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetActiveWithdrawsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveWithdraws(ctx, reqDb)

	if err != nil {
		return nil, withdraw_errors.ErrFindActiveWithdrawsFailed
	}

	return res, nil
}

func (r *withdrawRepository) FindByTrashed(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetTrashedWithdrawsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedWithdraws(ctx, reqDb)

	if err != nil {
		return nil, withdraw_errors.ErrFindTrashedWithdrawsFailed
	}

	return res, nil
}

func (r *withdrawRepository) FindAllByCardNumber(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*db.GetWithdrawsByCardNumberRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	withdraw, err := r.db.GetWithdrawsByCardNumber(ctx, reqDb)

	if err != nil {
		return nil, withdraw_errors.ErrFindWithdrawsByCardNumberFailed
	}

	return withdraw, nil

}

func (r *withdrawRepository) FindById(ctx context.Context, id int) (*db.GetWithdrawByIDRow, error) {
	withdraw, err := r.db.GetWithdrawByID(ctx, int32(id))

	if err != nil {
		return nil, withdraw_errors.ErrFindWithdrawByIdFailed
	}

	return withdraw, nil
}

func (r *withdrawRepository) GetMonthWithdrawStatusSuccess(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*db.GetMonthWithdrawStatusSuccessRow, error) {
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusSuccess(ctx, db.GetMonthWithdrawStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusSuccessFailed
	}

	return res, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyWithdrawStatusSuccessRow, error) {
	res, err := r.db.GetYearlyWithdrawStatusSuccess(ctx, int32(year))

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusSuccessFailed
	}

	return res, nil
}

func (r *withdrawRepository) GetMonthWithdrawStatusFailed(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*db.GetMonthWithdrawStatusFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusFailed(ctx, db.GetMonthWithdrawStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusFailedFailed
	}

	return res, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyWithdrawStatusFailedRow, error) {
	res, err := r.db.GetYearlyWithdrawStatusFailed(ctx, int32(year))

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusFailedFailed
	}

	return res, nil
}

func (r *withdrawRepository) GetMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*db.GetMonthWithdrawStatusSuccessCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusSuccessCardNumber(ctx, db.GetMonthWithdrawStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusSuccessByCardFailed
	}

	return res, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*db.GetYearlyWithdrawStatusSuccessCardNumberRow, error) {
	res, err := r.db.GetYearlyWithdrawStatusSuccessCardNumber(ctx, db.GetYearlyWithdrawStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusSuccessByCardFailed
	}

	return res, nil
}

func (r *withdrawRepository) GetMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*db.GetMonthWithdrawStatusFailedCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthWithdrawStatusFailedCardNumber(ctx, db.GetMonthWithdrawStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusFailedByCardFailed
	}

	return res, nil
}

func (r *withdrawRepository) GetYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*db.GetYearlyWithdrawStatusFailedCardNumberRow, error) {
	res, err := r.db.GetYearlyWithdrawStatusFailedCardNumber(ctx, db.GetYearlyWithdrawStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusFailedByCardFailed
	}

	return res, nil
}

func (r *withdrawRepository) GetMonthlyWithdraws(ctx context.Context, year int) ([]*db.GetMonthlyWithdrawsRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdraws(ctx, yearStart)

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthlyWithdrawsFailed
	}

	return res, nil

}

func (r *withdrawRepository) GetYearlyWithdraws(ctx context.Context, year int) ([]*db.GetYearlyWithdrawsRow, error) {
	res, err := r.db.GetYearlyWithdraws(ctx, year)

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawsFailed
	}

	return res, nil

}

func (r *withdrawRepository) GetMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetMonthlyWithdrawsByCardNumberRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawsByCardNumber(ctx, db.GetMonthlyWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetMonthlyWithdrawsByCardFailed
	}

	return res, nil

}

func (r *withdrawRepository) GetYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetYearlyWithdrawsByCardNumberRow, error) {
	res, err := r.db.GetYearlyWithdrawsByCardNumber(ctx, db.GetYearlyWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawsByCardFailed
	}

	return res, nil
}

func (r *withdrawRepository) CreateWithdraw(ctx context.Context, request *requests.CreateWithdrawRequest) (*db.CreateWithdrawRow, error) {
	req := db.CreateWithdrawParams{
		CardNumber:     request.CardNumber,
		WithdrawAmount: int32(request.WithdrawAmount),
		WithdrawTime:   request.WithdrawTime,
	}

	res, err := r.db.CreateWithdraw(ctx, req)

	if err != nil {
		return nil, withdraw_errors.ErrCreateWithdrawFailed
	}

	return res, nil
}

func (r *withdrawRepository) UpdateWithdraw(ctx context.Context, request *requests.UpdateWithdrawRequest) (*db.UpdateWithdrawRow, error) {
	req := db.UpdateWithdrawParams{
		WithdrawID:     int32(*request.WithdrawID),
		CardNumber:     request.CardNumber,
		WithdrawAmount: int32(request.WithdrawAmount),
		WithdrawTime:   request.WithdrawTime,
	}

	res, err := r.db.UpdateWithdraw(ctx, req)

	if err != nil {
		return nil, withdraw_errors.ErrUpdateWithdrawFailed
	}

	return res, nil
}

func (r *withdrawRepository) UpdateWithdrawStatus(ctx context.Context, request *requests.UpdateWithdrawStatus) (*db.UpdateWithdrawStatusRow, error) {
	req := db.UpdateWithdrawStatusParams{
		WithdrawID: int32(request.WithdrawID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateWithdrawStatus(ctx, req)

	if err != nil {
		return nil, withdraw_errors.ErrUpdateWithdrawStatusFailed
	}

	return res, nil
}

func (r *withdrawRepository) TrashedWithdraw(ctx context.Context, withdraw_id int) (*db.Withdraw, error) {
	res, err := r.db.TrashWithdraw(ctx, int32(withdraw_id))

	if err != nil {
		return nil, withdraw_errors.ErrTrashedWithdrawFailed
	}

	return res, nil
}

func (r *withdrawRepository) RestoreWithdraw(ctx context.Context, withdraw_id int) (*db.Withdraw, error) {
	res, err := r.db.RestoreWithdraw(ctx, int32(withdraw_id))

	if err != nil {
		return nil, withdraw_errors.ErrRestoreWithdrawFailed
	}

	return res, nil
}

func (r *withdrawRepository) DeleteWithdrawPermanent(ctx context.Context, withdraw_id int) (bool, error) {
	err := r.db.DeleteWithdrawPermanently(ctx, int32(withdraw_id))

	if err != nil {
		return false, withdraw_errors.ErrDeleteWithdrawPermanentFailed
	}

	return true, nil
}

func (r *withdrawRepository) RestoreAllWithdraw(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllWithdraws(ctx)

	if err != nil {
		return false, withdraw_errors.ErrRestoreAllWithdrawsFailed
	}

	return true, nil
}

func (r *withdrawRepository) DeleteAllWithdrawPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentWithdraws(ctx)

	if err != nil {
		return false, withdraw_errors.ErrDeleteAllWithdrawsPermanentFailed
	}

	return true, nil
}
