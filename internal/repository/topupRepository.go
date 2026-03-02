package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/topup_errors"
	"context"
	"time"
)

type topupRepository struct {
	db *db.Queries
}

func NewTopupRepository(db *db.Queries) TopupRepository {
	return &topupRepository{
		db: db,
	}
}

func (r *topupRepository) FindAllTopups(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetTopupsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTopups(ctx, reqDb)

	if err != nil {
		return nil, topup_errors.ErrFindAllTopupsFailed
	}

	return res, nil
}

func (r *topupRepository) FindByActive(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetActiveTopupsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTopups(ctx, reqDb)

	if err != nil {
		return nil, topup_errors.ErrFindTopupsByActiveFailed
	}

	return res, nil
}

func (r *topupRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetTrashedTopupsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTopups(ctx, reqDb)

	if err != nil {
		return nil, topup_errors.ErrFindTopupsByTrashedFailed
	}

	return res, nil
}

func (r *topupRepository) FindAllTopupByCardNumber(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*db.GetTopupsByCardNumberRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTopupsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetTopupsByCardNumber(ctx, reqDb)

	if err != nil {
		return nil, topup_errors.ErrFindTopupsByCardNumberFailed
	}

	return res, nil
}

func (r *topupRepository) FindById(ctx context.Context, topup_id int) (*db.GetTopupByIDRow, error) {
	res, err := r.db.GetTopupByID(ctx, int32(topup_id))
	if err != nil {
		return nil, topup_errors.ErrFindTopupByIdFailed
	}
	return res, nil
}

func (r *topupRepository) GetMonthTopupStatusSuccess(ctx context.Context, req *requests.MonthTopupStatus) ([]*db.GetMonthTopupStatusSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusSuccess(ctx, db.GetMonthTopupStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusSuccessFailed
	}

	return res, nil
}

func (r *topupRepository) GetYearlyTopupStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTopupStatusSuccessRow, error) {
	res, err := r.db.GetYearlyTopupStatusSuccess(ctx, int32(year))

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusSuccessFailed
	}

	return res, nil
}

func (r *topupRepository) GetMonthTopupStatusFailed(ctx context.Context, req *requests.MonthTopupStatus) ([]*db.GetMonthTopupStatusFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusFailed(ctx, db.GetMonthTopupStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusFailedFailed
	}

	return res, nil
}

func (r *topupRepository) GetYearlyTopupStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTopupStatusFailedRow, error) {
	res, err := r.db.GetYearlyTopupStatusFailed(ctx, int32(year))

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusFailedFailed
	}

	return res, nil
}

func (r *topupRepository) GetMonthTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*db.GetMonthTopupStatusSuccessCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusSuccessCardNumber(ctx, db.GetMonthTopupStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusSuccessByCardFailed
	}

	return res, nil
}

func (r *topupRepository) GetYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*db.GetYearlyTopupStatusSuccessCardNumberRow, error) {
	res, err := r.db.GetYearlyTopupStatusSuccessCardNumber(ctx, db.GetYearlyTopupStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusSuccessByCardFailed
	}

	return res, nil
}

func (r *topupRepository) GetMonthTopupStatusFailedByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*db.GetMonthTopupStatusFailedCardNumberRow, error) {
	cardNumber := req.CardNumber
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTopupStatusFailedCardNumber(ctx, db.GetMonthTopupStatusFailedCardNumberParams{
		CardNumber: cardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusFailedByCardFailed
	}

	return res, nil
}

func (r *topupRepository) GetYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*db.GetYearlyTopupStatusFailedCardNumberRow, error) {
	res, err := r.db.GetYearlyTopupStatusFailedCardNumber(ctx, db.GetYearlyTopupStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusSuccessByCardFailed
	}

	return res, nil
}

func (r *topupRepository) GetMonthlyTopupMethods(ctx context.Context, year int) ([]*db.GetMonthlyTopupMethodsRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupMethods(ctx, yearStart)

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupMethodsFailed
	}

	return res, nil
}

func (r *topupRepository) GetYearlyTopupMethods(ctx context.Context, year int) ([]*db.GetYearlyTopupMethodsRow, error) {
	res, err := r.db.GetYearlyTopupMethods(ctx, year)

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupMethodsFailed
	}

	return res, nil
}

func (r *topupRepository) GetMonthlyTopupAmounts(ctx context.Context, year int) ([]*db.GetMonthlyTopupAmountsRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmounts(ctx, yearStart)

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupAmountsFailed
	}

	return res, nil
}

func (r *topupRepository) GetYearlyTopupAmounts(ctx context.Context, year int) ([]*db.GetYearlyTopupAmountsRow, error) {
	res, err := r.db.GetYearlyTopupAmounts(ctx, year)

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupAmountsFailed
	}

	return res, nil
}

func (r *topupRepository) GetMonthlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetMonthlyTopupMethodsByCardNumberRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupMethodsByCardNumber(ctx, db.GetMonthlyTopupMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupMethodsByCardFailed
	}

	return res, nil
}

func (r *topupRepository) GetYearlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetYearlyTopupMethodsByCardNumberRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyTopupMethodsByCardNumber(ctx, db.GetYearlyTopupMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupMethodsByCardFailed
	}

	return res, nil
}

func (r *topupRepository) GetMonthlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetMonthlyTopupAmountsByCardNumberRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmountsByCardNumber(ctx, db.GetMonthlyTopupAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    yearStart,
	})

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupAmountsByCardFailed
	}

	return res, nil
}

func (r *topupRepository) GetYearlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetYearlyTopupAmountsByCardNumberRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyTopupAmountsByCardNumber(ctx, db.GetYearlyTopupAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupAmountsByCardFailed
	}

	return res, nil
}

func (r *topupRepository) CreateTopup(ctx context.Context, request *requests.CreateTopupRequest) (*db.CreateTopupRow, error) {
	req := db.CreateTopupParams{
		CardNumber:  request.CardNumber,
		TopupAmount: int32(request.TopupAmount),
		TopupMethod: request.TopupMethod,
	}

	res, err := r.db.CreateTopup(ctx, req)

	if err != nil {
		return nil, topup_errors.ErrCreateTopupFailed
	}

	return res, nil
}

func (r *topupRepository) UpdateTopup(ctx context.Context, request *requests.UpdateTopupRequest) (*db.UpdateTopupRow, error) {
	req := db.UpdateTopupParams{
		TopupID:     int32(*request.TopupID),
		CardNumber:  request.CardNumber,
		TopupAmount: int32(request.TopupAmount),
		TopupMethod: request.TopupMethod,
	}

	res, err := r.db.UpdateTopup(ctx, req)

	if err != nil {
		return nil, topup_errors.ErrUpdateTopupFailed
	}

	return res, nil
}

func (r *topupRepository) UpdateTopupAmount(ctx context.Context, request *requests.UpdateTopupAmount) (*db.UpdateTopupAmountRow, error) {
	req := db.UpdateTopupAmountParams{
		TopupID:     int32(request.TopupID),
		TopupAmount: int32(request.TopupAmount),
	}

	res, err := r.db.UpdateTopupAmount(ctx, req)

	if err != nil {
		return nil, topup_errors.ErrUpdateTopupAmountFailed
	}

	return res, nil
}

func (r *topupRepository) UpdateTopupStatus(ctx context.Context, request *requests.UpdateTopupStatus) (*db.UpdateTopupStatusRow, error) {
	req := db.UpdateTopupStatusParams{
		TopupID: int32(request.TopupID),
		Status:  request.Status,
	}

	res, err := r.db.UpdateTopupStatus(ctx, req)

	if err != nil {
		return nil, topup_errors.ErrUpdateTopupStatusFailed
	}

	return res, nil
}

func (r *topupRepository) TrashedTopup(ctx context.Context, topup_id int) (*db.Topup, error) {
	res, err := r.db.TrashTopup(ctx, int32(topup_id))
	if err != nil {
		return nil, topup_errors.ErrTrashedTopupFailed
	}
	return res, nil
}

func (r *topupRepository) RestoreTopup(ctx context.Context, topup_id int) (*db.Topup, error) {
	res, err := r.db.RestoreTopup(ctx, int32(topup_id))
	if err != nil {
		return nil, topup_errors.ErrRestoreTopupFailed
	}
	return res, nil
}

func (r *topupRepository) DeleteTopupPermanent(ctx context.Context, topup_id int) (bool, error) {
	err := r.db.DeleteTopupPermanently(ctx, int32(topup_id))
	if err != nil {
		return false, topup_errors.ErrDeleteTopupPermanentFailed
	}
	return true, nil
}

func (r *topupRepository) RestoreAllTopup(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllTopups(ctx)

	if err != nil {
		return false, topup_errors.ErrRestoreAllTopupFailed
	}

	return true, nil
}

func (r *topupRepository) DeleteAllTopupPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentTopups(ctx)

	if err != nil {
		return false, topup_errors.ErrDeleteAllTopupPermanentFailed
	}

	return true, nil
}
