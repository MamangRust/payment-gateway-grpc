package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	"context"
	"database/sql"
	"time"
)

type saldoRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.SaldoRecordMapping
}

func NewSaldoRepository(db *db.Queries, ctx context.Context, mapping recordmapper.SaldoRecordMapping) *saldoRepository {
	return &saldoRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *saldoRepository) FindAllSaldos(req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	saldos, err := r.db.GetSaldos(r.ctx, reqDb)

	if err != nil {
		return nil, nil, saldo_errors.ErrFindAllSaldosFailed
	}

	var totalCount int
	if len(saldos) > 0 {
		totalCount = int(saldos[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSaldosRecordAll(saldos), &totalCount, nil
}

func (r *saldoRepository) FindByActive(req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveSaldos(r.ctx, reqDb)

	if err != nil {
		return nil, nil, saldo_errors.ErrFindActiveSaldosFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSaldosRecordActive(res), &totalCount, nil

}

func (r *saldoRepository) FindByTrashed(req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	saldos, err := r.db.GetTrashedSaldos(r.ctx, reqDb)

	if err != nil {
		return nil, nil, saldo_errors.ErrFindTrashedSaldosFailed
	}

	var totalCount int
	if len(saldos) > 0 {
		totalCount = int(saldos[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSaldosRecordTrashed(saldos), &totalCount, nil
}

func (r *saldoRepository) FindByCardNumber(card_number string) (*record.SaldoRecord, error) {
	res, err := r.db.GetSaldoByCardNumber(r.ctx, card_number)

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByCardNumberFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoRepository) FindById(saldo_id int) (*record.SaldoRecord, error) {
	res, err := r.db.GetSaldoByID(r.ctx, int32(saldo_id))

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByIdFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoRepository) GetMonthlyTotalSaldoBalance(req *requests.MonthTotalSaldoBalance) ([]*record.SaldoMonthTotalBalance, error) {
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalSaldoBalance(r.ctx, db.GetMonthlyTotalSaldoBalanceParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, saldo_errors.ErrGetMonthlyTotalSaldoBalanceFailed
	}

	so := r.mapping.ToSaldoMonthTotalBalances(res)
	return so, nil
}

func (r *saldoRepository) GetYearTotalSaldoBalance(year int) ([]*record.SaldoYearTotalBalance, error) {
	res, err := r.db.GetYearlyTotalSaldoBalances(r.ctx, int32(year))

	if err != nil {
		return nil, saldo_errors.ErrGetYearTotalSaldoBalanceFailed
	}

	so := r.mapping.ToSaldoYearTotalBalances(res)

	return so, nil
}

func (r *saldoRepository) GetMonthlySaldoBalances(year int) ([]*record.SaldoMonthSaldoBalance, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlySaldoBalances(r.ctx, yearStart)

	if err != nil {
		return nil, saldo_errors.ErrGetMonthlySaldoBalancesFailed
	}

	so := r.mapping.ToSaldoMonthBalances(res)

	return so, nil
}

func (r *saldoRepository) GetYearlySaldoBalances(year int) ([]*record.SaldoYearSaldoBalance, error) {
	res, err := r.db.GetYearlySaldoBalances(r.ctx, year)

	if err != nil {
		return nil, saldo_errors.ErrGetYearlySaldoBalancesFailed
	}

	so := r.mapping.ToSaldoYearSaldoBalances(res)

	return so, nil
}

func (r *saldoRepository) CreateSaldo(request *requests.CreateSaldoRequest) (*record.SaldoRecord, error) {
	req := db.CreateSaldoParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}
	res, err := r.db.CreateSaldo(r.ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrCreateSaldoFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoRepository) UpdateSaldo(request *requests.UpdateSaldoRequest) (*record.SaldoRecord, error) {
	req := db.UpdateSaldoParams{
		SaldoID:      int32(*request.SaldoID),
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldo(r.ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoRepository) UpdateSaldoBalance(request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error) {
	req := db.UpdateSaldoBalanceParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldoBalance(r.ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoBalanceFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoRepository) UpdateSaldoWithdraw(request *requests.UpdateSaldoWithdraw) (*record.SaldoRecord, error) {
	withdrawAmount := sql.NullInt32{
		Int32: int32(*request.WithdrawAmount),
		Valid: request.WithdrawAmount != nil,
	}
	var withdrawTime sql.NullTime
	if request.WithdrawTime != nil {
		withdrawTime = sql.NullTime{
			Time:  *request.WithdrawTime,
			Valid: true,
		}
	}

	req := db.UpdateSaldoWithdrawParams{
		CardNumber:     request.CardNumber,
		WithdrawAmount: withdrawAmount,
		WithdrawTime:   withdrawTime,
	}

	res, err := r.db.UpdateSaldoWithdraw(r.ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoWithdrawFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoRepository) TrashedSaldo(saldo_id int) (*record.SaldoRecord, error) {
	res, err := r.db.TrashSaldo(r.ctx, int32(saldo_id))
	if err != nil {
		return nil, saldo_errors.ErrTrashSaldoFailed
	}
	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoRepository) RestoreSaldo(saldo_id int) (*record.SaldoRecord, error) {
	res, err := r.db.RestoreSaldo(r.ctx, int32(saldo_id))
	if err != nil {
		return nil, saldo_errors.ErrRestoreSaldoFailed
	}
	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoRepository) DeleteSaldoPermanent(saldo_id int) (bool, error) {
	err := r.db.DeleteSaldoPermanently(r.ctx, int32(saldo_id))
	if err != nil {
		return false, saldo_errors.ErrDeleteSaldoPermanentFailed
	}
	return true, nil
}

func (r *saldoRepository) RestoreAllSaldo() (bool, error) {
	err := r.db.RestoreAllSaldos(r.ctx)

	if err != nil {
		return false, saldo_errors.ErrRestoreAllSaldosFailed
	}

	return true, nil
}

func (r *saldoRepository) DeleteAllSaldoPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentSaldos(r.ctx)

	if err != nil {
		return false, saldo_errors.ErrDeleteAllSaldosPermanentFailed
	}

	return true, nil
}
