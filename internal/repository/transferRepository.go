package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/transfer_errors"
	"context"
	"time"
)

type transferRepository struct {
	db *db.Queries
}

func NewTransferRepository(db *db.Queries) TransferRepository {
	return &transferRepository{
		db: db,
	}
}

func (r *transferRepository) FindAll(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetTransfersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransfers(ctx, reqDb)

	if err != nil {
		return nil, transfer_errors.ErrFindAllTransfersFailed
	}

	return res, nil
}

func (r *transferRepository) FindByActive(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetActiveTransfersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTransfers(ctx, reqDb)

	if err != nil {
		return nil, transfer_errors.ErrFindActiveTransfersFailed
	}

	return res, nil
}

func (r *transferRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetTrashedTransfersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTransfers(ctx, reqDb)

	if err != nil {
		return nil, transfer_errors.ErrFindTrashedTransfersFailed
	}

	return res, nil
}

func (r *transferRepository) FindById(ctx context.Context, id int) (*db.GetTransferByIDRow, error) {
	transfer, err := r.db.GetTransferByID(ctx, int32(id))

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByIdFailed
	}

	return transfer, nil
}

func (r *transferRepository) GetMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*db.GetMonthTransferStatusSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccess(ctx, db.GetMonthTransferStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusSuccessFailed
	}

	return res, nil
}

func (r *transferRepository) GetYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransferStatusSuccessRow, error) {
	res, err := r.db.GetYearlyTransferStatusSuccess(ctx, int32(year))

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusSuccessFailed
	}

	return res, nil
}

func (r *transferRepository) GetMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*db.GetMonthTransferStatusFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailed(ctx, db.GetMonthTransferStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusFailedFailed
	}

	return res, nil
}

func (r *transferRepository) GetYearlyTransferStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTransferStatusFailedRow, error) {
	res, err := r.db.GetYearlyTransferStatusFailed(ctx, int32(year))

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusFailedFailed
	}

	return res, nil
}

func (r *transferRepository) GetMonthTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*db.GetMonthTransferStatusSuccessCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccessCardNumber(ctx, db.GetMonthTransferStatusSuccessCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusSuccessByCardFailed
	}

	return res, nil
}

func (r *transferRepository) GetYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*db.GetYearlyTransferStatusSuccessCardNumberRow, error) {
	res, err := r.db.GetYearlyTransferStatusSuccessCardNumber(ctx, db.GetYearlyTransferStatusSuccessCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      int32(req.Year),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusSuccessByCardFailed
	}

	return res, nil
}

func (r *transferRepository) GetMonthTransferStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*db.GetMonthTransferStatusFailedCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailedCardNumber(ctx, db.GetMonthTransferStatusFailedCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusFailedByCardFailed
	}

	return res, nil
}

func (r *transferRepository) GetYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*db.GetYearlyTransferStatusFailedCardNumberRow, error) {
	res, err := r.db.GetYearlyTransferStatusFailedCardNumber(ctx, db.GetYearlyTransferStatusFailedCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      int32(req.Year),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusFailedByCardFailed
	}

	return res, nil
}

func (r *transferRepository) GetMonthlyTransferAmounts(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountsRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmounts(ctx, yearStart)

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsFailed
	}

	return res, nil
}

func (r *transferRepository) GetYearlyTransferAmounts(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountsRow, error) {
	res, err := r.db.GetYearlyTransferAmounts(ctx, year)

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsFailed
	}
	return res, nil
}

func (r *transferRepository) GetMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetMonthlyTransferAmountsBySenderCardNumberRow, error) {
	res, err := r.db.GetMonthlyTransferAmountsBySenderCardNumber(ctx, db.GetMonthlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsBySenderCardFailed
	}

	return res, nil
}

func (r *transferRepository) GetMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetMonthlyTransferAmountsByReceiverCardNumberRow, error) {
	res, err := r.db.GetMonthlyTransferAmountsByReceiverCardNumber(ctx, db.GetMonthlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: req.CardNumber,
		Column2:    time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsByReceiverCardFailed
	}
	return res, nil
}

func (r *transferRepository) GetYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetYearlyTransferAmountsBySenderCardNumberRow, error) {
	res, err := r.db.GetYearlyTransferAmountsBySenderCardNumber(ctx, db.GetYearlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      req.Year,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsBySenderCardFailed
	}

	return res, nil
}

func (r *transferRepository) GetYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetYearlyTransferAmountsByReceiverCardNumberRow, error) {
	res, err := r.db.GetYearlyTransferAmountsByReceiverCardNumber(ctx, db.GetYearlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsByReceiverCardFailed
	}

	return res, nil
}

func (r *transferRepository) FindTransferByTransferFrom(ctx context.Context, transfer_from string) ([]*db.GetTransfersBySourceCardRow, error) {
	res, err := r.db.GetTransfersBySourceCard(ctx, transfer_from)

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByTransferFromFailed
	}

	return res, nil
}

func (r *transferRepository) FindTransferByTransferTo(ctx context.Context, transfer_to string) ([]*db.GetTransfersByDestinationCardRow, error) {
	res, err := r.db.GetTransfersByDestinationCard(ctx, transfer_to)

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByTransferToFailed
	}
	return res, nil
}

func (r *transferRepository) CreateTransfer(ctx context.Context, request *requests.CreateTransferRequest) (*db.CreateTransferRow, error) {
	req := db.CreateTransferParams{
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.CreateTransfer(ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrCreateTransferFailed
	}

	return res, nil
}

func (r *transferRepository) UpdateTransfer(ctx context.Context, request *requests.UpdateTransferRequest) (*db.UpdateTransferRow, error) {
	req := db.UpdateTransferParams{
		TransferID:     int32(*request.TransferID),
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.UpdateTransfer(ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferFailed
	}

	return res, nil

}

func (r *transferRepository) UpdateTransferAmount(ctx context.Context, request *requests.UpdateTransferAmountRequest) (*db.UpdateTransferAmountRow, error) {
	req := db.UpdateTransferAmountParams{
		TransferID:     int32(request.TransferID),
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.UpdateTransferAmount(ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferAmountFailed
	}

	return res, nil
}

func (r *transferRepository) UpdateTransferStatus(ctx context.Context, request *requests.UpdateTransferStatus) (*db.UpdateTransferStatusRow, error) {
	req := db.UpdateTransferStatusParams{
		TransferID: int32(request.TransferID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateTransferStatus(ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferStatusFailed
	}

	return res, nil
}

func (r *transferRepository) TrashedTransfer(ctx context.Context, transfer_id int) (*db.Transfer, error) {
	res, err := r.db.TrashTransfer(ctx, int32(transfer_id))

	if err != nil {
		return nil, transfer_errors.ErrTrashedTransferFailed
	}
	return res, nil
}

func (r *transferRepository) RestoreTransfer(ctx context.Context, transfer_id int) (*db.Transfer, error) {
	res, err := r.db.RestoreTransfer(ctx, int32(transfer_id))
	if err != nil {
		return nil, transfer_errors.ErrRestoreTransferFailed
	}
	return res, nil
}

func (r *transferRepository) DeleteTransferPermanent(ctx context.Context, transfer_id int) (bool, error) {
	err := r.db.DeleteTransferPermanently(ctx, int32(transfer_id))
	if err != nil {
		return false, transfer_errors.ErrDeleteTransferPermanentFailed
	}
	return true, nil
}

func (r *transferRepository) RestoreAllTransfer(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllTransfers(ctx)

	if err != nil {
		return false, transfer_errors.ErrRestoreAllTransfersFailed
	}

	return true, nil
}

func (r *transferRepository) DeleteAllTransferPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentTransfers(ctx)

	if err != nil {
		return false, transfer_errors.ErrDeleteAllTransfersPermanentFailed
	}

	return true, nil
}
