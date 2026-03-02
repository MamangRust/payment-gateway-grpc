package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/saldo_errors"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type saldoRepository struct {
	db *db.Queries
}

func NewSaldoRepository(db *db.Queries) SaldoRepository {
	return &saldoRepository{
		db: db,
	}
}

func (r *saldoRepository) FindAllSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetSaldosRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	saldos, err := r.db.GetSaldos(ctx, reqDb)

	if err != nil {
		return nil, saldo_errors.ErrFindAllSaldosFailed
	}

	return saldos, nil
}

func (r *saldoRepository) FindByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetActiveSaldosRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveSaldos(ctx, reqDb)

	if err != nil {
		return nil, saldo_errors.ErrFindActiveSaldosFailed
	}

	return res, nil
}

func (r *saldoRepository) FindByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetTrashedSaldosRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	saldos, err := r.db.GetTrashedSaldos(ctx, reqDb)

	if err != nil {
		return nil, saldo_errors.ErrFindTrashedSaldosFailed
	}

	return saldos, nil
}

func (r *saldoRepository) FindByCardNumber(ctx context.Context, card_number string) (*db.Saldo, error) {
	res, err := r.db.GetSaldoByCardNumber(ctx, card_number)

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByCardNumberFailed
	}

	return res, nil
}

func (r *saldoRepository) FindById(ctx context.Context, saldo_id int) (*db.GetSaldoByIDRow, error) {
	res, err := r.db.GetSaldoByID(ctx, int32(saldo_id))

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByIdFailed
	}

	return res, nil
}

func (r *saldoRepository) GetMonthlyTotalSaldoBalance(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*db.GetMonthlyTotalSaldoBalanceRow, error) {
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalSaldoBalance(ctx, db.GetMonthlyTotalSaldoBalanceParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, saldo_errors.ErrGetMonthlyTotalSaldoBalanceFailed
	}

	return res, nil
}

func (r *saldoRepository) GetYearTotalSaldoBalance(ctx context.Context, year int) ([]*db.GetYearlyTotalSaldoBalancesRow, error) {
	res, err := r.db.GetYearlyTotalSaldoBalances(ctx, int32(year))

	if err != nil {
		return nil, saldo_errors.ErrGetYearTotalSaldoBalanceFailed
	}

	return res, nil
}

func (r *saldoRepository) GetMonthlySaldoBalances(ctx context.Context, year int) ([]*db.GetMonthlySaldoBalancesRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlySaldoBalances(ctx, yearStart)

	if err != nil {
		return nil, saldo_errors.ErrGetMonthlySaldoBalancesFailed
	}

	return res, nil
}

func (r *saldoRepository) GetYearlySaldoBalances(ctx context.Context, year int) ([]*db.GetYearlySaldoBalancesRow, error) {
	res, err := r.db.GetYearlySaldoBalances(ctx, year)

	if err != nil {
		return nil, saldo_errors.ErrGetYearlySaldoBalancesFailed
	}

	return res, nil
}

func (r *saldoRepository) CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*db.CreateSaldoRow, error) {
	req := db.CreateSaldoParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}
	res, err := r.db.CreateSaldo(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrCreateSaldoFailed
	}

	return res, nil
}

func (r *saldoRepository) UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*db.UpdateSaldoRow, error) {
	req := db.UpdateSaldoParams{
		SaldoID:      int32(*request.SaldoID),
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldo(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoFailed
	}

	return res, nil
}

func (r *saldoRepository) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*db.UpdateSaldoBalanceRow, error) {
	req := db.UpdateSaldoBalanceParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldoBalance(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoBalanceFailed
	}

	return res, nil
}

func (r *saldoRepository) UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*db.UpdateSaldoWithdrawRow, error) {
	var withdrawAmount pgtype.Int4
	if request.WithdrawAmount != nil {
		withdrawAmount = pgtype.Int4{
			Int32: int32(*request.WithdrawAmount),
			Valid: true,
		}
	}

	var withdrawTime pgtype.Timestamp
	if request.WithdrawTime != nil {
		withdrawTime = pgtype.Timestamp{
			Time:  *request.WithdrawTime,
			Valid: true,
		}
	}

	req := db.UpdateSaldoWithdrawParams{
		CardNumber:     request.CardNumber,
		WithdrawAmount: &withdrawAmount.Int32,
		WithdrawTime:   withdrawTime,
	}

	res, err := r.db.UpdateSaldoWithdraw(ctx, req)
	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoWithdrawFailed
	}

	return res, nil
}

func (r *saldoRepository) TrashedSaldo(ctx context.Context, saldo_id int) (*db.Saldo, error) {
	res, err := r.db.TrashSaldo(ctx, int32(saldo_id))
	if err != nil {
		return nil, saldo_errors.ErrTrashSaldoFailed
	}
	return res, nil
}

func (r *saldoRepository) RestoreSaldo(ctx context.Context, saldo_id int) (*db.Saldo, error) {
	res, err := r.db.RestoreSaldo(ctx, int32(saldo_id))
	if err != nil {
		return nil, saldo_errors.ErrRestoreSaldoFailed
	}
	return res, nil
}

func (r *saldoRepository) DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, error) {
	err := r.db.DeleteSaldoPermanently(ctx, int32(saldo_id))
	if err != nil {
		return false, saldo_errors.ErrDeleteSaldoPermanentFailed
	}
	return true, nil
}

func (r *saldoRepository) RestoreAllSaldo(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllSaldos(ctx)

	if err != nil {
		return false, saldo_errors.ErrRestoreAllSaldosFailed
	}

	return true, nil
}

func (r *saldoRepository) DeleteAllSaldoPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentSaldos(ctx)

	if err != nil {
		return false, saldo_errors.ErrDeleteAllSaldosPermanentFailed
	}

	return true, nil
}
