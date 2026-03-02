package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/transaction_errors"
	"context"
	"time"
)

type transactionRepository struct {
	db *db.Queries
}

func NewTransactionRepository(db *db.Queries) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) FindAllTransactions(ctx context.Context, req *requests.FindAllTransactions) ([]*db.GetTransactionsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	transactions, err := r.db.GetTransactions(ctx, reqDb)

	if err != nil {
		return nil, transaction_errors.ErrFindAllTransactionsFailed
	}

	return transactions, nil
}

func (r *transactionRepository) FindAllTransactionByCardNumber(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*db.GetTransactionsByCardNumberRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	transactions, err := r.db.GetTransactionsByCardNumber(ctx, reqDb)

	if err != nil {
		return nil, transaction_errors.ErrFindTransactionsByCardNumberFailed
	}

	return transactions, nil
}

func (r *transactionRepository) FindByActive(ctx context.Context, req *requests.FindAllTransactions) ([]*db.GetActiveTransactionsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTransactions(ctx, reqDb)

	if err != nil {
		return nil, transaction_errors.ErrFindActiveTransactionsFailed
	}

	return res, nil
}

func (r *transactionRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTransactions) ([]*db.GetTrashedTransactionsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTransactions(ctx, reqDb)

	if err != nil {
		return nil, transaction_errors.ErrFindTrashedTransactionsFailed
	}

	return res, nil
}

func (r *transactionRepository) FindById(ctx context.Context, transaction_id int) (*db.GetTransactionByIDRow, error) {
	res, err := r.db.GetTransactionByID(ctx, int32(transaction_id))

	if err != nil {
		return nil, transaction_errors.ErrFindTransactionByIdFailed
	}

	return res, nil
}

func (r *transactionRepository) FindTransactionByMerchantId(ctx context.Context, merchant_id int) ([]*db.GetTransactionsByMerchantIDRow, error) {
	res, err := r.db.GetTransactionsByMerchantID(ctx, int32(merchant_id))

	if err != nil {
		return nil, transaction_errors.ErrFindTransactionByMerchantIdFailed
	}

	return res, nil
}

func (r *transactionRepository) GetMonthTransactionStatusSuccess(ctx context.Context, req *requests.MonthStatusTransaction) ([]*db.GetMonthTransactionStatusSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusSuccess(ctx, db.GetMonthTransactionStatusSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusSuccessFailed
	}

	return res, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransactionStatusSuccessRow, error) {
	res, err := r.db.GetYearlyTransactionStatusSuccess(ctx, int32(year))

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusSuccessFailed
	}

	return res, nil
}

func (r *transactionRepository) GetMonthTransactionStatusFailed(ctx context.Context, req *requests.MonthStatusTransaction) ([]*db.GetMonthTransactionStatusFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusFailed(ctx, db.GetMonthTransactionStatusFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusFailedFailed
	}

	return res, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTransactionStatusFailedRow, error) {
	res, err := r.db.GetYearlyTransactionStatusFailed(ctx, int32(year))

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusFailedFailed
	}

	return res, nil
}

func (r *transactionRepository) GetMonthTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*db.GetMonthTransactionStatusSuccessCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusSuccessCardNumber(ctx, db.GetMonthTransactionStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusSuccessByCardFailed
	}

	return res, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*db.GetYearlyTransactionStatusSuccessCardNumberRow, error) {
	res, err := r.db.GetYearlyTransactionStatusSuccessCardNumber(ctx, db.GetYearlyTransactionStatusSuccessCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusSuccessByCardFailed
	}

	return res, nil
}

func (r *transactionRepository) GetMonthTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*db.GetMonthTransactionStatusFailedCardNumberRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusFailedCardNumber(ctx, db.GetMonthTransactionStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusFailedByCardFailed
	}

	return res, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*db.GetYearlyTransactionStatusFailedCardNumberRow, error) {
	res, err := r.db.GetYearlyTransactionStatusFailedCardNumber(ctx, db.GetYearlyTransactionStatusFailedCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    int32(req.Year),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusFailedByCardFailed
	}

	return res, nil
}

func (r *transactionRepository) GetMonthlyPaymentMethods(ctx context.Context, year int) ([]*db.GetMonthlyPaymentMethodsRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethods(ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyPaymentMethodsFailed
	}

	return res, nil
}

func (r *transactionRepository) GetYearlyPaymentMethods(ctx context.Context, year int) ([]*db.GetYearlyPaymentMethodsRow, error) {
	res, err := r.db.GetYearlyPaymentMethods(ctx, year)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyPaymentMethodsFailed
	}

	return res, nil
}

func (r *transactionRepository) GetMonthlyAmounts(ctx context.Context, year int) ([]*db.GetMonthlyAmountsRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyAmounts(ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountsFailed
	}

	return res, nil
}

func (r *transactionRepository) GetYearlyAmounts(ctx context.Context, year int) ([]*db.GetYearlyAmountsRow, error) {
	res, err := r.db.GetYearlyAmounts(ctx, year)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountsFailed
	}

	return res, nil
}

func (r *transactionRepository) GetMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetMonthlyPaymentMethodsByCardNumberRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetMonthlyPaymentMethodsByCardNumber(ctx, db.GetMonthlyPaymentMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyPaymentMethodsByCardFailed
	}

	return res, nil
}

func (r *transactionRepository) GetYearlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetYearlyPaymentMethodsByCardNumberRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyPaymentMethodsByCardNumber(ctx, db.GetYearlyPaymentMethodsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyPaymentMethodsByCardFailed
	}

	return res, nil
}

func (r *transactionRepository) GetMonthlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetMonthlyAmountsByCardNumberRow, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetMonthlyAmountsByCardNumber(ctx, db.GetMonthlyAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountsByCardFailed
	}

	return res, nil
}

func (r *transactionRepository) GetYearlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetYearlyAmountsByCardNumberRow, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetYearlyAmountsByCardNumber(ctx, db.GetYearlyAmountsByCardNumberParams{
		CardNumber: cardNumber,
		Column2:    year,
	})
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountsByCardFailed
	}

	return res, nil
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*db.CreateTransactionRow, error) {
	req := db.CreateTransactionParams{
		CardNumber:      request.CardNumber,
		Amount:          int32(request.Amount),
		PaymentMethod:   request.PaymentMethod,
		MerchantID:      int32(*request.MerchantID),
		TransactionTime: request.TransactionTime,
	}

	res, err := r.db.CreateTransaction(ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrCreateTransactionFailed
	}

	return res, nil
}

func (r *transactionRepository) UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*db.UpdateTransactionRow, error) {
	req := db.UpdateTransactionParams{
		TransactionID:   int32(*request.TransactionID),
		CardNumber:      request.CardNumber,
		Amount:          int32(request.Amount),
		PaymentMethod:   request.PaymentMethod,
		MerchantID:      int32(*request.MerchantID),
		TransactionTime: request.TransactionTime,
	}

	res, err := r.db.UpdateTransaction(ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionFailed
	}

	return res, nil
}

func (r *transactionRepository) UpdateTransactionStatus(ctx context.Context, request *requests.UpdateTransactionStatus) (*db.UpdateTransactionStatusRow, error) {
	req := db.UpdateTransactionStatusParams{
		TransactionID: int32(request.TransactionID),
		Status:        request.Status,
	}

	res, err := r.db.UpdateTransactionStatus(ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionStatusFailed
	}

	return res, nil
}

func (r *transactionRepository) TrashedTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	res, err := r.db.TrashTransaction(ctx, int32(transaction_id))
	if err != nil {
		return nil, transaction_errors.ErrTrashedTransactionFailed
	}
	return res, nil
}

func (r *transactionRepository) RestoreTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	res, err := r.db.RestoreTransaction(ctx, int32(transaction_id))
	if err != nil {
		return nil, transaction_errors.ErrRestoreTransactionFailed
	}
	return res, nil
}

func (r *transactionRepository) DeleteTransactionPermanent(ctx context.Context, transaction_id int) (bool, error) {
	err := r.db.DeleteTransactionPermanently(ctx, int32(transaction_id))
	if err != nil {

		return false, transaction_errors.ErrDeleteTransactionPermanentFailed
	}
	return true, nil
}

func (r *transactionRepository) RestoreAllTransaction(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllTransactions(ctx)

	if err != nil {
		return false, transaction_errors.ErrRestoreAllTransactionsFailed
	}

	return true, nil
}

func (r *transactionRepository) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentTransactions(ctx)

	if err != nil {
		return false, transaction_errors.ErrDeleteAllTransactionsPermanentFailed
	}
	return true, nil
}
