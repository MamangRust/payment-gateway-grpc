package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
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

func (r *transactionRepository) FindAllTransactions(search string, page, pageSize int) ([]*record.TransactionRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetTransactionsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	transactions, err := r.db.GetTransactions(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find transactions: %w", err)
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordAll(transactions), totalCount, nil
}

func (r *transactionRepository) FindAllTransactionByCardNumber(card_number string, search string, page, pageSize int) ([]*record.TransactionRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetTransactionsByCardNumberParams{
		CardNumber: card_number,
		Column2:    search,
		Limit:      int32(pageSize),
		Offset:     int32(offset),
	}

	transactions, err := r.db.GetTransactionsByCardNumber(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find transactions: %w", err)
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsByCardNumberRecord(transactions), totalCount, nil
}

func (r *transactionRepository) FindById(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionByID(r.ctx, int32(transaction_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) FindTransactionByMerchantId(merchant_id int) ([]*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionsByMerchantID(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find transaction by merchant id: %w", err)
	}

	return r.mapping.ToTransactionsRecord(res), nil
}

func (r *transactionRepository) GetMonthTransactionStatusSuccess(year int, month int) ([]*record.TransactionRecordMonthStatusSuccess, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
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
		return nil, fmt.Errorf("failed to get month top-up status success for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToTransactionRecordsMonthStatusSuccess(res)

	return so, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusSuccess(year int) ([]*record.TransactionRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransactionStatusSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status success for year %d: %w", year, err)
	}

	so := r.mapping.ToTransactionRecordsYearStatusSuccess(res)

	return so, nil
}

func (r *transactionRepository) GetMonthTransactionStatusFailed(year int, month int) ([]*record.TransactionRecordMonthStatusFailed, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
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
		return nil, fmt.Errorf("failed to get month top-up status failed for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToTransactionRecordsMonthStatusFailed(res)

	return so, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusFailed(year int) ([]*record.TransactionRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransactionStatusFailed(r.ctx, int32(year))

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status failed for year %d: %w", year, err)
	}

	so := r.mapping.ToTransactionRecordsYearStatusFailed(res)

	return so, nil
}

func (r *transactionRepository) GetMonthTransactionStatusSuccessByCardNumber(card string, year int, month int) ([]*record.TransactionRecordMonthStatusSuccess, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusSuccessCardNumber(r.ctx, db.GetMonthTransactionStatusSuccessCardNumberParams{
		CardNumber: card,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get month top-up status success for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToTransactionRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusSuccessByCardNumber(card_number string, year int) ([]*record.TransactionRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransactionStatusSuccessCardNumber(r.ctx, db.GetYearlyTransactionStatusSuccessCardNumberParams{
		CardNumber: card_number,
		Column2:    int32(year),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status success for year %d: %w", year, err)
	}

	so := r.mapping.ToTransactionRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transactionRepository) GetMonthTransactionStatusFailedByCardNumber(card_number string, year int, month int) ([]*record.TransactionRecordMonthStatusFailed, error) {
	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransactionStatusFailedCardNumber(r.ctx, db.GetMonthTransactionStatusFailedCardNumberParams{
		CardNumber: card_number,
		Column2:    currentDate,
		Column3:    lastDayCurrentMonth,
		Column4:    prevDate,
		Column5:    lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get month top-up status failed for year %d and month %d: %w", year, month, err)
	}

	so := r.mapping.ToTransactionRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

func (r *transactionRepository) GetYearlyTransactionStatusFailedByCardNumber(card_number string, year int) ([]*record.TransactionRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransactionStatusFailedCardNumber(r.ctx, db.GetYearlyTransactionStatusFailedCardNumberParams{
		CardNumber: card_number,
		Column2:    int32(year),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly top-up status failed for year %d: %w", year, err)
	}

	so := r.mapping.ToTransactionRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

func (r *transactionRepository) GetMonthlyPaymentMethods(year int) ([]*record.TransactionMonthMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethods(r.ctx, yearStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly payment methods: %w", err)
	}

	return r.mapping.ToTransactionMonthlyMethods(res), nil
}

func (r *transactionRepository) GetYearlyPaymentMethods(year int) ([]*record.TransactionYearMethod, error) {
	res, err := r.db.GetYearlyPaymentMethods(r.ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly payment methods: %w", err)
	}

	return r.mapping.ToTransactionYearlyMethods(res), nil
}

func (r *transactionRepository) GetMonthlyAmounts(year int) ([]*record.TransactionMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyAmounts(r.ctx, yearStart)

	if err != nil {
		return nil, fmt.Errorf("failed to get monthly amounts: %w", err)
	}

	return r.mapping.ToTransactionMonthlyAmounts(res), nil
}

func (r *transactionRepository) GetYearlyAmounts(year int) ([]*record.TransactionYearlyAmount, error) {
	res, err := r.db.GetYearlyAmounts(r.ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly amounts: %w", err)
	}

	return r.mapping.ToTransactionYearlyAmounts(res), nil
}

func (r *transactionRepository) GetMonthlyPaymentMethodsByCardNumber(card_number string, year int) ([]*record.TransactionMonthMethod, error) {
	res, err := r.db.GetMonthlyPaymentMethodsByCardNumber(r.ctx, db.GetMonthlyPaymentMethodsByCardNumberParams{
		CardNumber: card_number,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly payment methods by card number: %w", err)
	}

	return r.mapping.ToTransactionMonthlyMethodsByCardNumber(res), nil
}

func (r *transactionRepository) GetYearlyPaymentMethodsByCardNumber(card_number string, year int) ([]*record.TransactionYearMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodsByCardNumber(r.ctx, db.GetYearlyPaymentMethodsByCardNumberParams{
		CardNumber: card_number,
		Column2:    year,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly payment methods by card number: %w", err)
	}

	return r.mapping.ToTransactionYearlyMethodsByCardNumber(res), nil
}

func (r *transactionRepository) GetMonthlyAmountsByCardNumber(card_number string, year int) ([]*record.TransactionMonthAmount, error) {
	res, err := r.db.GetMonthlyAmountsByCardNumber(r.ctx, db.GetMonthlyAmountsByCardNumberParams{
		CardNumber: card_number,
		Column2:    time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly amounts by card number: %w", err)
	}

	return r.mapping.ToTransactionMonthlyAmountsByCardNumber(res), nil
}

func (r *transactionRepository) GetYearlyAmountsByCardNumber(card_number string, year int) ([]*record.TransactionYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountsByCardNumber(r.ctx, db.GetYearlyAmountsByCardNumberParams{
		CardNumber: card_number,
		Column2:    year,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly amounts by card number: %w", err)
	}

	return r.mapping.ToTransactionYearlyAmountsByCardNumber(res), nil
}

func (r *transactionRepository) FindByActive(search string, page, pageSize int) ([]*record.TransactionRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetActiveTransactionsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTransactions(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find active transactions: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordActive(res), totalCount, nil
}

func (r *transactionRepository) FindByTrashed(search string, page, pageSize int) ([]*record.TransactionRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetTrashedTransactionsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTransactions(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find trashed transactions: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordTrashed(res), totalCount, nil
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
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) UpdateTransaction(request *requests.UpdateTransactionRequest) (*record.TransactionRecord, error) {
	req := db.UpdateTransactionParams{
		TransactionID:   int32(request.TransactionID),
		CardNumber:      request.CardNumber,
		Amount:          int32(request.Amount),
		PaymentMethod:   request.PaymentMethod,
		MerchantID:      int32(*request.MerchantID),
		TransactionTime: request.TransactionTime,
	}

	err := r.db.UpdateTransaction(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	res, err := r.db.GetTransactionByID(r.ctx, int32(request.TransactionID))

	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) UpdateTransactionStatus(request *requests.UpdateTransactionStatus) (*record.TransactionRecord, error) {
	req := db.UpdateTransactionStatusParams{
		TransactionID: int32(request.TransactionID),
		Status:        request.Status,
	}

	err := r.db.UpdateTransactionStatus(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update Transaction amount :%w", err)
	}

	res, err := r.db.GetTransactionByID(r.ctx, req.TransactionID)

	if err != nil {
		return nil, fmt.Errorf("failed to find Transaction: %w", err)
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) TrashedTransaction(transaction_id int) (*record.TransactionRecord, error) {
	err := r.db.TrashTransaction(r.ctx, int32(transaction_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash transaction: %w", err)
	}

	transaction, err := r.db.GetTrashedTransactionByID(r.ctx, int32(transaction_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find trashed by id transaction: %w", err)
	}

	return r.mapping.ToTransactionRecord(transaction), nil
}

func (r *transactionRepository) RestoreTransaction(topup_id int) (*record.TransactionRecord, error) {
	err := r.db.RestoreTransaction(r.ctx, int32(topup_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore transaction: %w", err)
	}

	topup, err := r.db.GetTransactionByID(r.ctx, int32(topup_id))

	if err != nil {
		return nil, fmt.Errorf("failed not found transaction :%w", err)
	}

	return r.mapping.ToTransactionRecord(topup), nil
}

func (r *transactionRepository) DeleteTransactionPermanent(topup_id int) (bool, error) {
	err := r.db.DeleteTransactionPermanently(r.ctx, int32(topup_id))
	if err != nil {
		return false, fmt.Errorf("failed to delete transaction: %w", err)
	}
	return true, nil
}

func (r *transactionRepository) RestoreAllTransaction() (bool, error) {
	err := r.db.RestoreAllTransactions(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all transactions: %w", err)
	}

	return true, nil
}

func (r *transactionRepository) DeleteAllTransactionPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTransactions(r.ctx)
	if err != nil {
		return false, fmt.Errorf("failed to delete all transactions permanently: %w", err)
	}
	return true, nil
}
