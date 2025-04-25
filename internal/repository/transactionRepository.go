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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no transaction found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to retrieve transaction (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no active transaction found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to find active transaction (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no active transaction found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to find active transaction (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no trashed transaction found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to find trashed transaction (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction not found with ID: %d", transaction_id)
		}
		return nil, fmt.Errorf("failed to get transaction by ID %d: %w", transaction_id, err)
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) FindTransactionByMerchantId(merchant_id int) ([]*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionsByMerchantID(r.ctx, int32(merchant_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction not found with merchant_id: %d", merchant_id)
		}
		return nil, fmt.Errorf("failed to get transaction by merchant_id %d: %w", merchant_id, err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no transaction data for success status found for year %d, month %d", req.Year, req.Month)
		}
		return nil, fmt.Errorf("failed to get monthly transaction status success for year %d, month %d: %w", req.Year, req.Month, err)
	}

	so := r.mapping.ToTransactionRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusSuccess(year int) ([]*record.TransactionRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransactionStatusSuccess(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transaction success found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly transaction status success for year %d: %w", year, err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no month transaction failed found for year %d and month %d", req.Year, req.Month)
		}
		return nil, fmt.Errorf("failed to get month transaction status failed for year %d and month %d: %w", req.Year, req.Month, err)
	}

	so := r.mapping.ToTransactionRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusFailed(year int) ([]*record.TransactionRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransactionStatusFailed(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transaction failed found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly transaction status failed for year %d: %w", year, err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no month transaction success found for year %d, month %d and card_number %s", req.Year, req.Month, req.CardNumber)
		}
		return nil, fmt.Errorf("failed to get month transaction status success for year %d, month %d and card_number %s: %w", req.Year, req.Month, req.CardNumber, err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transaction success found for year %d and card_number %s", req.Year, req.CardNumber)
		}
		return nil, fmt.Errorf("failed to get yearly transaction status success for year %d and card_number %s: %w", req.Year, req.CardNumber, err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no month transaction failed found for year %d, month %d and card_number %s", req.Year, req.Month, req.CardNumber)
		}
		return nil, fmt.Errorf("failed to get month transaction status failed for year %d, month %d and card_number %s: %w", req.Year, req.Month, req.CardNumber, err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transaction failed found for year %d and card_number %s", req.Year, req.CardNumber)
		}
		return nil, fmt.Errorf("failed to get yearly transaction status failed for year %d and card_number %s: %w", req.Year, req.CardNumber, err)
	}

	so := r.mapping.ToTransactionRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

func (r *transactionRepository) GetMonthlyPaymentMethods(year int) ([]*record.TransactionMonthMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethods(r.ctx, yearStart)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transaction payment method found")
		}
		return nil, fmt.Errorf("failed to get monthly transaction payment method: %w", err)
	}

	return r.mapping.ToTransactionMonthlyMethods(res), nil
}

func (r *transactionRepository) GetYearlyPaymentMethods(year int) ([]*record.TransactionYearMethod, error) {
	res, err := r.db.GetYearlyPaymentMethods(r.ctx, year)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transaction payment method found")
		}
		return nil, fmt.Errorf("failed to get yearly transaction payment method: %w", err)
	}

	return r.mapping.ToTransactionYearlyMethods(res), nil
}

func (r *transactionRepository) GetMonthlyAmounts(year int) ([]*record.TransactionMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyAmounts(r.ctx, yearStart)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transaction amount found")
		}
		return nil, fmt.Errorf("failed to get monthly transaction amount: %w", err)
	}

	return r.mapping.ToTransactionMonthlyAmounts(res), nil
}

func (r *transactionRepository) GetYearlyAmounts(year int) ([]*record.TransactionYearlyAmount, error) {
	res, err := r.db.GetYearlyAmounts(r.ctx, year)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transaction amounts found")
		}
		return nil, fmt.Errorf("failed to get yearly transaction amounts: %w", err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transaction payment method card number %s and year %d found", cardNumber, year)
		}
		return nil, fmt.Errorf("failed to get monthly transaction payment method card number %s and year %d: %w", cardNumber, year, err)
	}

	return r.mapping.ToTransactionMonthlyMethodsByCardNumber(res), nil
}

func (r *transactionRepository) GetYearlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearMethod, error) {
	year := req.Year
	cardNumber := req.CardNumber

	res, err := r.db.GetYearlyPaymentMethodsByCardNumber(r.ctx, db.GetYearlyPaymentMethodsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transfer payment method card number %s and year %d found", cardNumber, year)
		}
		return nil, fmt.Errorf("failed to get yearly transfer payment method card number %s and year %d: %w", cardNumber, year, err)
	}

	return r.mapping.ToTransactionYearlyMethodsByCardNumber(res), nil
}

func (r *transactionRepository) GetMonthlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthAmount, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetMonthlyAmountsByCardNumber(r.ctx, db.GetMonthlyAmountsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transaction payment method card number %s and year %d found", cardNumber, year)
		}
		return nil, fmt.Errorf("failed to get monthly transaction payment method card number %s and year %d: %w", cardNumber, year, err)
	}

	return r.mapping.ToTransactionMonthlyAmountsByCardNumber(res), nil
}

func (r *transactionRepository) GetYearlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearlyAmount, error) {
	cardNumber := req.CardNumber
	year := req.Year

	res, err := r.db.GetYearlyAmountsByCardNumber(r.ctx, db.GetYearlyAmountsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Year,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transaction amounts card number %s and year %d found", cardNumber, year)
		}
		return nil, fmt.Errorf("failed to get yearly transaction amounts card number %s and year %d: %w", cardNumber, year, err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid transaction data: %w", err)
		}
		return nil, fmt.Errorf("failed to create transaction: invalid or incomplete transaction data: %w", err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction ID %d not found for update", request.TransactionID)
		}
		return nil, fmt.Errorf("failed to update transaction ID %d: transaction not found or invalid update data", request.TransactionID)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction ID %d not found for update", request.TransactionID)
		}
		return nil, fmt.Errorf("failed to update transaction ID %d: transaction not found or invalid update data", request.TransactionID)
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) TrashedTransaction(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.TrashTransaction(r.ctx, int32(transaction_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction ID %d not found or already trashed", transaction_id)
		}
		return nil, fmt.Errorf("failed to trash transaction ID %d: %w", transaction_id, err)
	}
	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) RestoreTransaction(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.RestoreTransaction(r.ctx, int32(transaction_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction ID %d not found in trash", transaction_id)
		}
		return nil, fmt.Errorf("failed to restore transaction ID %d: %w", transaction_id, err)
	}
	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) DeleteTransactionPermanent(transaction_id int) (bool, error) {
	err := r.db.DeleteTransactionPermanently(r.ctx, int32(transaction_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("transaction ID %d not found or already deleted", transaction_id)
		}
		return false, fmt.Errorf("failed to permanently delete transaction ID %d: %w", transaction_id, err)
	}
	return true, nil
}

func (r *transactionRepository) RestoreAllTransaction() (bool, error) {
	err := r.db.RestoreAllTransactions(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no trashed transaction available to restore")
		}
		return false, fmt.Errorf("failed to restore trashed transaction: %w", err)
	}

	return true, nil
}

func (r *transactionRepository) DeleteAllTransactionPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTransactions(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no trashed transaction available to delete permanently")
		}
		return false, fmt.Errorf("failed to permanently delete transaction: %w", err)
	}
	return true, nil
}
