package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/transfer_errors"
	"context"
	"time"
)

type transferRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransferRecordMapping
}

func NewTransferRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransferRecordMapping) *transferRepository {
	return &transferRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transferRepository) FindAll(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransfers(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transfer_errors.ErrFindAllTransfersFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransfersRecordAll(res), &totalCount, nil
}

func (r *transferRepository) FindByActive(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTransfers(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transfer_errors.ErrFindActiveTransfersFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransfersRecordActive(res), &totalCount, nil
}

func (r *transferRepository) FindByTrashed(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTransfers(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transfer_errors.ErrFindTrashedTransfersFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransfersRecordTrashed(res), &totalCount, nil
}

func (r *transferRepository) FindById(id int) (*record.TransferRecord, error) {
	transfer, err := r.db.GetTransferByID(r.ctx, int32(id))

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByIdFailed
	}

	return r.mapping.ToTransferRecord(transfer), nil
}

func (r *transferRepository) GetMonthTransferStatusSuccess(req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccess(r.ctx, db.GetMonthTransferStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusSuccessFailed
	}

	so := r.mapping.ToTransferRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *transferRepository) GetYearlyTransferStatusSuccess(year int) ([]*record.TransferRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransferStatusSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusSuccessFailed
	}

	so := r.mapping.ToTransferRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *transferRepository) GetMonthTransferStatusFailed(req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailed(r.ctx, db.GetMonthTransferStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusFailedFailed
	}

	so := r.mapping.ToTransferRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *transferRepository) GetYearlyTransferStatusFailed(year int) ([]*record.TransferRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransferStatusFailed(r.ctx, int32(year))

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusFailedFailed
	}

	so := r.mapping.ToTransferRecordsYearStatusFailed(res)

	return so, nil
}

func (r *transferRepository) GetMonthTransferStatusSuccessByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccessCardNumber(r.ctx, db.GetMonthTransferStatusSuccessCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusSuccessByCardFailed
	}

	so := r.mapping.ToTransferRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transferRepository) GetYearlyTransferStatusSuccessByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransferStatusSuccessCardNumber(r.ctx, db.GetYearlyTransferStatusSuccessCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      int32(req.Year),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusSuccessByCardFailed
	}

	so := r.mapping.ToTransferRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transferRepository) GetMonthTransferStatusFailedByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailedCardNumber(r.ctx, db.GetMonthTransferStatusFailedCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusFailedByCardFailed
	}

	so := r.mapping.ToTransferRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

func (r *transferRepository) GetYearlyTransferStatusFailedByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransferStatusFailedCardNumber(r.ctx, db.GetYearlyTransferStatusFailedCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      int32(req.Year),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusFailedByCardFailed
	}

	so := r.mapping.ToTransferRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

func (r *transferRepository) GetMonthlyTransferAmounts(year int) ([]*record.TransferMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmounts(r.ctx, yearStart)

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsFailed
	}

	return r.mapping.ToTransferMonthAmounts(res), nil
}

func (r *transferRepository) GetYearlyTransferAmounts(year int) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmounts(r.ctx, year)

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsFailed
	}
	return r.mapping.ToTransferYearAmounts(res), nil
}

func (r *transferRepository) GetMonthlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error) {
	res, err := r.db.GetMonthlyTransferAmountsBySenderCardNumber(r.ctx, db.GetMonthlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsBySenderCardFailed
	}

	return r.mapping.ToTransferMonthAmountsSender(res), nil
}

func (r *transferRepository) GetMonthlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error) {
	res, err := r.db.GetMonthlyTransferAmountsByReceiverCardNumber(r.ctx, db.GetMonthlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: req.CardNumber,
		Column2:    time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsByReceiverCardFailed
	}
	return r.mapping.ToTransferMonthAmountsReceiver(res), nil
}

func (r *transferRepository) GetYearlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountsBySenderCardNumber(r.ctx, db.GetYearlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      req.Year,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsBySenderCardFailed
	}

	return r.mapping.ToTransferYearAmountsSender(res), nil
}

func (r *transferRepository) GetYearlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountsByReceiverCardNumber(r.ctx, db.GetYearlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsByReceiverCardFailed
	}

	return r.mapping.ToTransferYearAmountsReceiver(res), nil
}

func (r *transferRepository) FindTransferByTransferFrom(transfer_from string) ([]*record.TransferRecord, error) {
	res, err := r.db.GetTransfersBySourceCard(r.ctx, transfer_from)

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByTransferFromFailed
	}

	return r.mapping.ToTransfersRecord(res), nil
}

func (r *transferRepository) FindTransferByTransferTo(transfer_to string) ([]*record.TransferRecord, error) {
	res, err := r.db.GetTransfersByDestinationCard(r.ctx, transfer_to)

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByTransferToFailed
	}
	return r.mapping.ToTransfersRecord(res), nil
}

func (r *transferRepository) CreateTransfer(request *requests.CreateTransferRequest) (*record.TransferRecord, error) {
	req := db.CreateTransferParams{
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.CreateTransfer(r.ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrCreateTransferFailed
	}

	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) UpdateTransfer(request *requests.UpdateTransferRequest) (*record.TransferRecord, error) {
	req := db.UpdateTransferParams{
		TransferID:     int32(*request.TransferID),
		TransferFrom:   request.TransferFrom,
		TransferTo:     request.TransferTo,
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.UpdateTransfer(r.ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferFailed
	}

	return r.mapping.ToTransferRecord(res), nil

}

func (r *transferRepository) UpdateTransferAmount(request *requests.UpdateTransferAmountRequest) (*record.TransferRecord, error) {
	req := db.UpdateTransferAmountParams{
		TransferID:     int32(request.TransferID),
		TransferAmount: int32(request.TransferAmount),
	}

	res, err := r.db.UpdateTransferAmount(r.ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferAmountFailed
	}

	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) UpdateTransferStatus(request *requests.UpdateTransferStatus) (*record.TransferRecord, error) {
	req := db.UpdateTransferStatusParams{
		TransferID: int32(request.TransferID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateTransferStatus(r.ctx, req)

	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferStatusFailed
	}

	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) TrashedTransfer(transfer_id int) (*record.TransferRecord, error) {
	res, err := r.db.TrashTransfer(r.ctx, int32(transfer_id))

	if err != nil {
		return nil, transfer_errors.ErrTrashedTransferFailed
	}
	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) RestoreTransfer(transfer_id int) (*record.TransferRecord, error) {
	res, err := r.db.RestoreTransfer(r.ctx, int32(transfer_id))
	if err != nil {
		return nil, transfer_errors.ErrRestoreTransferFailed
	}
	return r.mapping.ToTransferRecord(res), nil
}

func (r *transferRepository) DeleteTransferPermanent(transfer_id int) (bool, error) {
	err := r.db.DeleteTransferPermanently(r.ctx, int32(transfer_id))
	if err != nil {
		return false, transfer_errors.ErrDeleteTransferPermanentFailed
	}
	return true, nil
}

func (r *transferRepository) RestoreAllTransfer() (bool, error) {
	err := r.db.RestoreAllTransfers(r.ctx)

	if err != nil {
		return false, transfer_errors.ErrRestoreAllTransfersFailed
	}

	return true, nil
}

func (r *transferRepository) DeleteAllTransferPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTransfers(r.ctx)

	if err != nil {
		return false, transfer_errors.ErrDeleteAllTransfersPermanentFailed
	}

	return true, nil
}
