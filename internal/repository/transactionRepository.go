package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/transaction_errors"
	"context"
	"time"
)

type transactionRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransactionRecordMapping
}

func NewTransactionRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransactionRecordMapping) *transactionRepository {
	return &transactionRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transactionRepository) FindAllTransactions(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	transactions, err := r.db.GetTransactions(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindAllTransactionsFailed
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordAll(transactions), &totalCount, nil
}

func (r *transactionRepository) FindAllTransactionByCardNumber(req *requests.FindAllTransactionCardNumber) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	transactions, err := r.db.GetTransactionsByCardNumber(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindTransactionsByCardNumberFailed
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsByCardNumberRecord(transactions), &totalCount, nil
}

func (r *transactionRepository) FindByActive(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTransactions(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindActiveTransactionsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordActive(res), &totalCount, nil
}

func (r *transactionRepository) FindByTrashed(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTransactions(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindTrashedTransactionsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordTrashed(res), &totalCount, nil
}

func (r *transactionRepository) FindById(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionByID(r.ctx, int32(transaction_id))

	if err != nil {
		return nil, transaction_errors.ErrFindTransactionByIdFailed
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) FindTransactionByMerchantId(merchant_id int) ([]*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionsByMerchantID(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, transaction_errors.ErrFindTransactionByMerchantIdFailed
	}

	return r.mapping.ToTransactionsRecord(res), nil
}

func (r *transactionRepository) GetMonthTransactionStatusSuccess(req *requests.MonthStatusTransaction) ([]*record.TransactionRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusSuccess(r.ctx, db.GetMonthTransactionStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusSuccessFailed
	}

	so := r.mapping.ToTransactionRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusSuccess(year int) ([]*record.TransactionRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransactionStatusSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusSuccessFailed
	}

	so := r.mapping.ToTransactionRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *transactionRepository) GetMonthTransactionStatusFailed(req *requests.MonthStatusTransaction) ([]*record.TransactionRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusFailed(r.ctx, db.GetMonthTransactionStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusFailedFailed
	}

	so := r.mapping.ToTransactionRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusFailed(year int) ([]*record.TransactionRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransactionStatusFailed(r.ctx, int32(year))

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusFailedFailed
	}

	so := r.mapping.ToTransactionRecordsYearStatusFailed(res)

	return so, nil
}

func (r *transactionRepository) GetMonthTransactionStatusSuccessByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*record.TransactionRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusSuccessCardNumber(r.ctx, db.GetMonthTransactionStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusSuccessByCardFailed
	}

	so := r.mapping.ToTransactionRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusSuccessByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*record.TransactionRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransactionStatusSuccessCardNumber(r.ctx, db.GetYearlyTransactionStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusSuccessByCardFailed
	}

	so := r.mapping.ToTransactionRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transactionRepository) GetMonthTransactionStatusFailedByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*record.TransactionRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusFailedCardNumber(r.ctx, db.GetMonthTransactionStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusFailedByCardFailed
	}

	so := r.mapping.ToTransactionRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusFailedByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*record.TransactionRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransactionStatusFailedCardNumber(r.ctx, db.GetYearlyTransactionStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusFailedByCardFailed
	}

	so := r.mapping.ToTransactionRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

func (r *transactionRepository) GetMonthlyPaymentMethods(year int) ([]*record.TransactionMonthMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethods(r.ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyPaymentMethodsFailed
	}

	return r.mapping.ToTransactionMonthlyMethods(res), nil
}

func (r *transactionRepository) GetYearlyPaymentMethods(year int) ([]*record.TransactionYearMethod, error) {
	res, err := r.db.GetYearlyPaymentMethods(r.ctx, year)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyPaymentMethodsFailed
	}

	return r.mapping.ToTransactionYearlyMethods(res), nil
}

func (r *transactionRepository) GetMonthlyAmounts(year int) ([]*record.TransactionMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyAmounts(r.ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountsFailed
	}

	return r.mapping.ToTransactionMonthlyAmounts(res), nil
}

func (r *transactionRepository) GetYearlyAmounts(year int) ([]*record.TransactionYearlyAmount, error) {
	res, err := r.db.GetYearlyAmounts(r.ctx, year)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountsFailed
	}

	return r.mapping.ToTransactionYearlyAmounts(res), nil
}

func (r *transactionRepository) GetMonthlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetMonthlyPaymentMethodsByCardNumber(r.ctx, db.GetMonthlyPaymentMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyPaymentMethodsByCardFailed
	}

	return r.mapping.ToTransactionMonthlyMethodsByCardNumber(res), nil
}

func (r *transactionRepository) GetYearlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyPaymentMethodsByCardNumber(r.ctx, db.GetYearlyPaymentMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyPaymentMethodsByCardFailed
	}

	return r.mapping.ToTransactionYearlyMethodsByCardNumber(res), nil
}

func (r *transactionRepository) GetMonthlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthAmount, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetMonthlyAmountsByCardNumber(r.ctx, db.GetMonthlyAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountsByCardFailed
	}

	return r.mapping.ToTransactionMonthlyAmountsByCardNumber(res), nil
}

func (r *transactionRepository) GetYearlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearlyAmount, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetYearlyAmountsByCardNumber(r.ctx, db.GetYearlyAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountsByCardFailed
	}

	return r.mapping.ToTransactionYearlyAmountsByCardNumber(res), nil
}

func (r *transactionRepository) CreateTransaction(request *requests.CreateTransactionRequest) (*record.TransactionRecord, error) {
	req := db.CreateTransactionParams{
		CardNumber:      request.CardNumber,
		Amount:          int32(request.Amount),
		PaymentMethod:   request.PaymentMethod,
		MerchantID:      int32(*request.MerchantID),
		TransactionTime: request.TransactionTime,
	}

	res, err := r.db.CreateTransaction(r.ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrCreateTransactionFailed
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) UpdateTransaction(request *requests.UpdateTransactionRequest) (*record.TransactionRecord, error) {
	req := db.UpdateTransactionParams{
		TransactionID:   int32(*request.TransactionID),
		CardNumber:      request.CardNumber,
		Amount:          int32(request.Amount),
		PaymentMethod:   request.PaymentMethod,
		MerchantID:      int32(*request.MerchantID),
		TransactionTime: request.TransactionTime,
	}

	res, err := r.db.UpdateTransaction(r.ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionFailed
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) UpdateTransactionStatus(request *requests.UpdateTransactionStatus) (*record.TransactionRecord, error) {
	req := db.UpdateTransactionStatusParams{
		TransactionID: int32(request.TransactionID),
		Status:        request.Status,
	}

	res, err := r.db.UpdateTransactionStatus(r.ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionStatusFailed
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) TrashedTransaction(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.TrashTransaction(r.ctx, int32(transaction_id))
	if err != nil {
		return nil, transaction_errors.ErrTrashedTransactionFailed
	}
	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) RestoreTransaction(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.RestoreTransaction(r.ctx, int32(transaction_id))
	if err != nil {
		return nil, transaction_errors.ErrRestoreTransactionFailed
	}
	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) DeleteTransactionPermanent(transaction_id int) (bool, error) {
	err := r.db.DeleteTransactionPermanently(r.ctx, int32(transaction_id))
	if err != nil {

		return false, transaction_errors.ErrDeleteTransactionPermanentFailed
	}
	return true, nil
}

func (r *transactionRepository) RestoreAllTransaction() (bool, error) {
	err := r.db.RestoreAllTransactions(r.ctx)

	if err != nil {
		return false, transaction_errors.ErrRestoreAllTransactionsFailed
	}

	return true, nil
}

func (r *transactionRepository) DeleteAllTransactionPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTransactions(r.ctx)

	if err != nil {
		return false, transaction_errors.ErrDeleteAllTransactionsPermanentFailed
	}
	return true, nil
}
