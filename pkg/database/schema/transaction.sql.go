// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: transaction.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const countAllTransactions = `-- name: CountAllTransactions :one
SELECT COUNT(*) FROM transactions WHERE deleted_at IS NULL
`

// Count All Transactions
func (q *Queries) CountAllTransactions(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countAllTransactions)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countTransactions = `-- name: CountTransactions :one
SELECT COUNT(*)
FROM transactions
WHERE deleted_at IS NULL
    AND ($1::TEXT IS NULL OR
        card_number ILIKE '%' || $1 || '%' OR
        payment_method ILIKE '%' || $1 || '%' OR
        CAST(transaction_time AS TEXT) ILIKE '%' || $1 || '%')
`

func (q *Queries) CountTransactions(ctx context.Context, dollar_1 string) (int64, error) {
	row := q.db.QueryRowContext(ctx, countTransactions, dollar_1)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countTransactionsByDate = `-- name: CountTransactionsByDate :one
SELECT COUNT(*)
FROM transactions
WHERE deleted_at IS NULL
  AND transaction_time::DATE = $1::DATE
`

// Count Transactions by Date
func (q *Queries) CountTransactionsByDate(ctx context.Context, dollar_1 time.Time) (int64, error) {
	row := q.db.QueryRowContext(ctx, countTransactionsByDate, dollar_1)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO
    transactions (
        card_number,
        amount,
        payment_method,
        merchant_id,
        transaction_time,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        current_timestamp,
        current_timestamp
    ) RETURNING transaction_id, card_number, amount, payment_method, merchant_id, transaction_time, created_at, updated_at, deleted_at
`

type CreateTransactionParams struct {
	CardNumber      string    `json:"card_number"`
	Amount          int32     `json:"amount"`
	PaymentMethod   string    `json:"payment_method"`
	MerchantID      int32     `json:"merchant_id"`
	TransactionTime time.Time `json:"transaction_time"`
}

// Create Transaction
func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (*Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.CardNumber,
		arg.Amount,
		arg.PaymentMethod,
		arg.MerchantID,
		arg.TransactionTime,
	)
	var i Transaction
	err := row.Scan(
		&i.TransactionID,
		&i.CardNumber,
		&i.Amount,
		&i.PaymentMethod,
		&i.MerchantID,
		&i.TransactionTime,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const deleteAllPermanentTransactions = `-- name: DeleteAllPermanentTransactions :exec
DELETE FROM transactions
WHERE
    deleted_at IS NOT NULL
`

// Delete All Trashed Transactions Permanently
func (q *Queries) DeleteAllPermanentTransactions(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllPermanentTransactions)
	return err
}

const deleteTransactionPermanently = `-- name: DeleteTransactionPermanently :exec
DELETE FROM transactions WHERE transaction_id = $1 AND deleted_at IS NOT NULL
`

// Delete Transaction Permanently
func (q *Queries) DeleteTransactionPermanently(ctx context.Context, transactionID int32) error {
	_, err := q.db.ExecContext(ctx, deleteTransactionPermanently, transactionID)
	return err
}

const getActiveTransactions = `-- name: GetActiveTransactions :many
SELECT
    transaction_id, card_number, amount, payment_method, merchant_id, transaction_time, created_at, updated_at, deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    transactions
WHERE
    deleted_at IS NULL
    AND ($1::TEXT IS NULL OR card_number ILIKE '%' || $1 || '%' OR payment_method ILIKE '%' || $1 || '%')
ORDER BY
    transaction_time DESC
LIMIT $2 OFFSET $3
`

type GetActiveTransactionsParams struct {
	Column1 string `json:"column_1"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

type GetActiveTransactionsRow struct {
	TransactionID   int32        `json:"transaction_id"`
	CardNumber      string       `json:"card_number"`
	Amount          int32        `json:"amount"`
	PaymentMethod   string       `json:"payment_method"`
	MerchantID      int32        `json:"merchant_id"`
	TransactionTime time.Time    `json:"transaction_time"`
	CreatedAt       sql.NullTime `json:"created_at"`
	UpdatedAt       sql.NullTime `json:"updated_at"`
	DeletedAt       sql.NullTime `json:"deleted_at"`
	TotalCount      int64        `json:"total_count"`
}

// Get Active Transactions with Pagination, Search, and Count
func (q *Queries) GetActiveTransactions(ctx context.Context, arg GetActiveTransactionsParams) ([]*GetActiveTransactionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getActiveTransactions, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetActiveTransactionsRow
	for rows.Next() {
		var i GetActiveTransactionsRow
		if err := rows.Scan(
			&i.TransactionID,
			&i.CardNumber,
			&i.Amount,
			&i.PaymentMethod,
			&i.MerchantID,
			&i.TransactionTime,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMonthlyAmounts = `-- name: GetMonthlyAmounts :many
SELECT
    TO_CHAR(t.transaction_time, 'Mon') AS month,
    SUM(t.amount) AS total_amount
FROM
    transactions t
WHERE
    t.deleted_at IS NULL
    AND EXTRACT(YEAR FROM t.transaction_time) = $1
GROUP BY
    TO_CHAR(t.transaction_time, 'Mon'),
    EXTRACT(MONTH FROM t.transaction_time)
ORDER BY
    EXTRACT(MONTH FROM t.transaction_time)
`

type GetMonthlyAmountsRow struct {
	Month       string `json:"month"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) GetMonthlyAmounts(ctx context.Context, transactionTime time.Time) ([]*GetMonthlyAmountsRow, error) {
	rows, err := q.db.QueryContext(ctx, getMonthlyAmounts, transactionTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetMonthlyAmountsRow
	for rows.Next() {
		var i GetMonthlyAmountsRow
		if err := rows.Scan(&i.Month, &i.TotalAmount); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMonthlyAmountyByCardNumber = `-- name: GetMonthlyAmountyByCardNumber :many
SELECT
    TO_CHAR(t.transaction_time, 'Mon') AS month,
    SUM(t.amount) AS total_amount
FROM
    transactions t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
    AND EXTRACT(YEAR FROM t.transaction_time) = $2
GROUP BY
    TO_CHAR(t.transaction_time, 'Mon'),
    EXTRACT(MONTH FROM t.transaction_time)
ORDER BY
    EXTRACT(MONTH FROM t.transaction_time)
`

type GetMonthlyAmountyByCardNumberParams struct {
	CardNumber      string    `json:"card_number"`
	TransactionTime time.Time `json:"transaction_time"`
}

type GetMonthlyAmountyByCardNumberRow struct {
	Month       string `json:"month"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) GetMonthlyAmountyByCardNumber(ctx context.Context, arg GetMonthlyAmountyByCardNumberParams) ([]*GetMonthlyAmountyByCardNumberRow, error) {
	rows, err := q.db.QueryContext(ctx, getMonthlyAmountyByCardNumber, arg.CardNumber, arg.TransactionTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetMonthlyAmountyByCardNumberRow
	for rows.Next() {
		var i GetMonthlyAmountyByCardNumberRow
		if err := rows.Scan(&i.Month, &i.TotalAmount); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMonthlyPaymentMethods = `-- name: GetMonthlyPaymentMethods :many
SELECT
    TO_CHAR(t.transaction_time, 'Mon') AS month,
    t.payment_method,
    COUNT(t.transaction_id) AS total_transactions,
    SUM(t.amount) AS total_amount
FROM
    transactions t
WHERE
    t.deleted_at IS NULL
    AND EXTRACT(YEAR FROM t.transaction_time) = $1
GROUP BY
    TO_CHAR(t.transaction_time, 'Mon'),
    EXTRACT(MONTH FROM t.transaction_time),
    t.payment_method
ORDER BY
    EXTRACT(MONTH FROM t.transaction_time)
`

type GetMonthlyPaymentMethodsRow struct {
	Month             string `json:"month"`
	PaymentMethod     string `json:"payment_method"`
	TotalTransactions int64  `json:"total_transactions"`
	TotalAmount       int64  `json:"total_amount"`
}

func (q *Queries) GetMonthlyPaymentMethods(ctx context.Context, transactionTime time.Time) ([]*GetMonthlyPaymentMethodsRow, error) {
	rows, err := q.db.QueryContext(ctx, getMonthlyPaymentMethods, transactionTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetMonthlyPaymentMethodsRow
	for rows.Next() {
		var i GetMonthlyPaymentMethodsRow
		if err := rows.Scan(
			&i.Month,
			&i.PaymentMethod,
			&i.TotalTransactions,
			&i.TotalAmount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMonthlyPaymentMethodsByCardNumber = `-- name: GetMonthlyPaymentMethodsByCardNumber :many
SELECT
    TO_CHAR(t.transaction_time, 'Mon') AS month,
    t.payment_method,
    COUNT(t.transaction_id) AS total_transactions,
    SUM(t.amount) AS total_amount
FROM
    transactions t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
    AND EXTRACT(YEAR FROM t.transaction_time) = $2
GROUP BY
    TO_CHAR(t.transaction_time, 'Mon'),
    EXTRACT(MONTH FROM t.transaction_time),
    t.payment_method
ORDER BY
    EXTRACT(MONTH FROM t.transaction_time)
`

type GetMonthlyPaymentMethodsByCardNumberParams struct {
	CardNumber      string    `json:"card_number"`
	TransactionTime time.Time `json:"transaction_time"`
}

type GetMonthlyPaymentMethodsByCardNumberRow struct {
	Month             string `json:"month"`
	PaymentMethod     string `json:"payment_method"`
	TotalTransactions int64  `json:"total_transactions"`
	TotalAmount       int64  `json:"total_amount"`
}

func (q *Queries) GetMonthlyPaymentMethodsByCardNumber(ctx context.Context, arg GetMonthlyPaymentMethodsByCardNumberParams) ([]*GetMonthlyPaymentMethodsByCardNumberRow, error) {
	rows, err := q.db.QueryContext(ctx, getMonthlyPaymentMethodsByCardNumber, arg.CardNumber, arg.TransactionTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetMonthlyPaymentMethodsByCardNumberRow
	for rows.Next() {
		var i GetMonthlyPaymentMethodsByCardNumberRow
		if err := rows.Scan(
			&i.Month,
			&i.PaymentMethod,
			&i.TotalTransactions,
			&i.TotalAmount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTransactionByCardNumber = `-- name: GetTransactionByCardNumber :many
SELECT
    transaction_id, card_number, amount, payment_method, merchant_id, transaction_time, created_at, updated_at, deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    transactions
WHERE
    deleted_at IS NULL
    AND card_number = $1
    AND ($2::TEXT IS NULL OR payment_method ILIKE '%' || $2 || '%')
ORDER BY
    transaction_time DESC
LIMIT $3 OFFSET $4
`

type GetTransactionByCardNumberParams struct {
	CardNumber string `json:"card_number"`
	Column2    string `json:"column_2"`
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
}

type GetTransactionByCardNumberRow struct {
	TransactionID   int32        `json:"transaction_id"`
	CardNumber      string       `json:"card_number"`
	Amount          int32        `json:"amount"`
	PaymentMethod   string       `json:"payment_method"`
	MerchantID      int32        `json:"merchant_id"`
	TransactionTime time.Time    `json:"transaction_time"`
	CreatedAt       sql.NullTime `json:"created_at"`
	UpdatedAt       sql.NullTime `json:"updated_at"`
	DeletedAt       sql.NullTime `json:"deleted_at"`
	TotalCount      int64        `json:"total_count"`
}

func (q *Queries) GetTransactionByCardNumber(ctx context.Context, arg GetTransactionByCardNumberParams) ([]*GetTransactionByCardNumberRow, error) {
	rows, err := q.db.QueryContext(ctx, getTransactionByCardNumber,
		arg.CardNumber,
		arg.Column2,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetTransactionByCardNumberRow
	for rows.Next() {
		var i GetTransactionByCardNumberRow
		if err := rows.Scan(
			&i.TransactionID,
			&i.CardNumber,
			&i.Amount,
			&i.PaymentMethod,
			&i.MerchantID,
			&i.TransactionTime,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTransactionByID = `-- name: GetTransactionByID :one
SELECT transaction_id, card_number, amount, payment_method, merchant_id, transaction_time, created_at, updated_at, deleted_at
FROM transactions
WHERE
    transaction_id = $1
    AND deleted_at IS NULL
`

// Get Transaction by ID
func (q *Queries) GetTransactionByID(ctx context.Context, transactionID int32) (*Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransactionByID, transactionID)
	var i Transaction
	err := row.Scan(
		&i.TransactionID,
		&i.CardNumber,
		&i.Amount,
		&i.PaymentMethod,
		&i.MerchantID,
		&i.TransactionTime,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const getTransactions = `-- name: GetTransactions :many
SELECT
    transaction_id, card_number, amount, payment_method, merchant_id, transaction_time, created_at, updated_at, deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    transactions
WHERE
    deleted_at IS NULL
    AND ($1::TEXT IS NULL OR card_number ILIKE '%' || $1 || '%' OR payment_method ILIKE '%' || $1 || '%')
ORDER BY
    transaction_time DESC
LIMIT $2 OFFSET $3
`

type GetTransactionsParams struct {
	Column1 string `json:"column_1"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

type GetTransactionsRow struct {
	TransactionID   int32        `json:"transaction_id"`
	CardNumber      string       `json:"card_number"`
	Amount          int32        `json:"amount"`
	PaymentMethod   string       `json:"payment_method"`
	MerchantID      int32        `json:"merchant_id"`
	TransactionTime time.Time    `json:"transaction_time"`
	CreatedAt       sql.NullTime `json:"created_at"`
	UpdatedAt       sql.NullTime `json:"updated_at"`
	DeletedAt       sql.NullTime `json:"deleted_at"`
	TotalCount      int64        `json:"total_count"`
}

// Search Transactions with Pagination
func (q *Queries) GetTransactions(ctx context.Context, arg GetTransactionsParams) ([]*GetTransactionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getTransactions, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetTransactionsRow
	for rows.Next() {
		var i GetTransactionsRow
		if err := rows.Scan(
			&i.TransactionID,
			&i.CardNumber,
			&i.Amount,
			&i.PaymentMethod,
			&i.MerchantID,
			&i.TransactionTime,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTransactionsByCardNumber = `-- name: GetTransactionsByCardNumber :many
SELECT transaction_id, card_number, amount, payment_method, merchant_id, transaction_time, created_at, updated_at, deleted_at
FROM transactions
WHERE
    card_number = $1
    AND deleted_at IS NULL
ORDER BY transaction_time DESC
`

// Get Transactions by Card Number
func (q *Queries) GetTransactionsByCardNumber(ctx context.Context, cardNumber string) ([]*Transaction, error) {
	rows, err := q.db.QueryContext(ctx, getTransactionsByCardNumber, cardNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.TransactionID,
			&i.CardNumber,
			&i.Amount,
			&i.PaymentMethod,
			&i.MerchantID,
			&i.TransactionTime,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTransactionsByMerchantID = `-- name: GetTransactionsByMerchantID :many
SELECT transaction_id, card_number, amount, payment_method, merchant_id, transaction_time, created_at, updated_at, deleted_at
FROM transactions
WHERE
    merchant_id = $1
    AND deleted_at IS NULL
ORDER BY transaction_time DESC
`

// Get Transactions by Merchant ID
func (q *Queries) GetTransactionsByMerchantID(ctx context.Context, merchantID int32) ([]*Transaction, error) {
	rows, err := q.db.QueryContext(ctx, getTransactionsByMerchantID, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.TransactionID,
			&i.CardNumber,
			&i.Amount,
			&i.PaymentMethod,
			&i.MerchantID,
			&i.TransactionTime,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTrashedTransactionByID = `-- name: GetTrashedTransactionByID :one
SELECT transaction_id, card_number, amount, payment_method, merchant_id, transaction_time, created_at, updated_at, deleted_at
FROM transactions
WHERE
    transaction_id = $1
    AND deleted_at IS NOT NULL
`

// Get Trashed By Transaction ID
func (q *Queries) GetTrashedTransactionByID(ctx context.Context, transactionID int32) (*Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTrashedTransactionByID, transactionID)
	var i Transaction
	err := row.Scan(
		&i.TransactionID,
		&i.CardNumber,
		&i.Amount,
		&i.PaymentMethod,
		&i.MerchantID,
		&i.TransactionTime,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const getTrashedTransactions = `-- name: GetTrashedTransactions :many
SELECT
    transaction_id, card_number, amount, payment_method, merchant_id, transaction_time, created_at, updated_at, deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    transactions
WHERE
    deleted_at IS NOT NULL
    AND ($1::TEXT IS NULL OR card_number ILIKE '%' || $1 || '%' OR payment_method ILIKE '%' || $1 || '%')
ORDER BY
    transaction_time DESC
LIMIT $2 OFFSET $3
`

type GetTrashedTransactionsParams struct {
	Column1 string `json:"column_1"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

type GetTrashedTransactionsRow struct {
	TransactionID   int32        `json:"transaction_id"`
	CardNumber      string       `json:"card_number"`
	Amount          int32        `json:"amount"`
	PaymentMethod   string       `json:"payment_method"`
	MerchantID      int32        `json:"merchant_id"`
	TransactionTime time.Time    `json:"transaction_time"`
	CreatedAt       sql.NullTime `json:"created_at"`
	UpdatedAt       sql.NullTime `json:"updated_at"`
	DeletedAt       sql.NullTime `json:"deleted_at"`
	TotalCount      int64        `json:"total_count"`
}

// Get Trashed Transactions with Pagination, Search, and Count
func (q *Queries) GetTrashedTransactions(ctx context.Context, arg GetTrashedTransactionsParams) ([]*GetTrashedTransactionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getTrashedTransactions, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetTrashedTransactionsRow
	for rows.Next() {
		var i GetTrashedTransactionsRow
		if err := rows.Scan(
			&i.TransactionID,
			&i.CardNumber,
			&i.Amount,
			&i.PaymentMethod,
			&i.MerchantID,
			&i.TransactionTime,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getYearlyAmounts = `-- name: GetYearlyAmounts :many
SELECT
    EXTRACT(YEAR FROM t.transaction_time) AS year,
    SUM(t.amount) AS total_amount
FROM
    transactions t
WHERE
    t.deleted_at IS NULL
GROUP BY
    EXTRACT(YEAR FROM t.transaction_time)
ORDER BY
    year
`

type GetYearlyAmountsRow struct {
	Year        string `json:"year"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) GetYearlyAmounts(ctx context.Context) ([]*GetYearlyAmountsRow, error) {
	rows, err := q.db.QueryContext(ctx, getYearlyAmounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetYearlyAmountsRow
	for rows.Next() {
		var i GetYearlyAmountsRow
		if err := rows.Scan(&i.Year, &i.TotalAmount); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getYearlyAmountsByCardNumber = `-- name: GetYearlyAmountsByCardNumber :many
SELECT
    EXTRACT(YEAR FROM t.transaction_time) AS year,
    SUM(t.amount) AS total_amount
FROM
    transactions t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
GROUP BY
    EXTRACT(YEAR FROM t.transaction_time)
ORDER BY
    year
`

type GetYearlyAmountsByCardNumberRow struct {
	Year        string `json:"year"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) GetYearlyAmountsByCardNumber(ctx context.Context, cardNumber string) ([]*GetYearlyAmountsByCardNumberRow, error) {
	rows, err := q.db.QueryContext(ctx, getYearlyAmountsByCardNumber, cardNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetYearlyAmountsByCardNumberRow
	for rows.Next() {
		var i GetYearlyAmountsByCardNumberRow
		if err := rows.Scan(&i.Year, &i.TotalAmount); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getYearlyPaymentMethods = `-- name: GetYearlyPaymentMethods :many
SELECT
    EXTRACT(YEAR FROM t.transaction_time) AS year,
    t.payment_method,
    COUNT(t.transaction_id) AS total_transactions,
    SUM(t.amount) AS total_amount
FROM
    transactions t
WHERE
    t.deleted_at IS NULL
GROUP BY
    EXTRACT(YEAR FROM t.transaction_time),
    t.payment_method
ORDER BY
    year
`

type GetYearlyPaymentMethodsRow struct {
	Year              string `json:"year"`
	PaymentMethod     string `json:"payment_method"`
	TotalTransactions int64  `json:"total_transactions"`
	TotalAmount       int64  `json:"total_amount"`
}

func (q *Queries) GetYearlyPaymentMethods(ctx context.Context) ([]*GetYearlyPaymentMethodsRow, error) {
	rows, err := q.db.QueryContext(ctx, getYearlyPaymentMethods)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetYearlyPaymentMethodsRow
	for rows.Next() {
		var i GetYearlyPaymentMethodsRow
		if err := rows.Scan(
			&i.Year,
			&i.PaymentMethod,
			&i.TotalTransactions,
			&i.TotalAmount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getYearlyPaymentMethodsByCardNumber = `-- name: GetYearlyPaymentMethodsByCardNumber :many
SELECT
    EXTRACT(YEAR FROM t.transaction_time) AS year,
    t.payment_method,
    COUNT(t.transaction_id) AS total_transactions,
    SUM(t.amount) AS total_amount
FROM
    transactions t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
GROUP BY
    EXTRACT(YEAR FROM t.transaction_time),
    t.payment_method
ORDER BY
    year
`

type GetYearlyPaymentMethodsByCardNumberRow struct {
	Year              string `json:"year"`
	PaymentMethod     string `json:"payment_method"`
	TotalTransactions int64  `json:"total_transactions"`
	TotalAmount       int64  `json:"total_amount"`
}

func (q *Queries) GetYearlyPaymentMethodsByCardNumber(ctx context.Context, cardNumber string) ([]*GetYearlyPaymentMethodsByCardNumberRow, error) {
	rows, err := q.db.QueryContext(ctx, getYearlyPaymentMethodsByCardNumber, cardNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetYearlyPaymentMethodsByCardNumberRow
	for rows.Next() {
		var i GetYearlyPaymentMethodsByCardNumberRow
		if err := rows.Scan(
			&i.Year,
			&i.PaymentMethod,
			&i.TotalTransactions,
			&i.TotalAmount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const restoreAllTransactions = `-- name: RestoreAllTransactions :exec
UPDATE transactions
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL
`

// Restore All Trashed Transactions
func (q *Queries) RestoreAllTransactions(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, restoreAllTransactions)
	return err
}

const restoreTransaction = `-- name: RestoreTransaction :exec
UPDATE transactions
SET
    deleted_at = NULL
WHERE
    transaction_id = $1
    AND deleted_at IS NOT NULL
`

// Restore Trashed Transaction
func (q *Queries) RestoreTransaction(ctx context.Context, transactionID int32) error {
	_, err := q.db.ExecContext(ctx, restoreTransaction, transactionID)
	return err
}

const transaction_CountAll = `-- name: Transaction_CountAll :one
SELECT COUNT(*)
FROM transactions
WHERE deleted_at IS NULL
`

func (q *Queries) Transaction_CountAll(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, transaction_CountAll)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const trashTransaction = `-- name: TrashTransaction :exec
UPDATE transactions
SET
    deleted_at = current_timestamp
WHERE
    transaction_id = $1
    AND deleted_at IS NULL
`

// Trash Transaction
func (q *Queries) TrashTransaction(ctx context.Context, transactionID int32) error {
	_, err := q.db.ExecContext(ctx, trashTransaction, transactionID)
	return err
}

const updateTransaction = `-- name: UpdateTransaction :exec
UPDATE transactions
SET
    card_number = $2,
    amount = $3,
    payment_method = $4,
    merchant_id = $5,
    transaction_time = $6,
    updated_at = current_timestamp
WHERE
    transaction_id = $1
    AND deleted_at IS NULL
`

type UpdateTransactionParams struct {
	TransactionID   int32     `json:"transaction_id"`
	CardNumber      string    `json:"card_number"`
	Amount          int32     `json:"amount"`
	PaymentMethod   string    `json:"payment_method"`
	MerchantID      int32     `json:"merchant_id"`
	TransactionTime time.Time `json:"transaction_time"`
}

// Update Transaction
func (q *Queries) UpdateTransaction(ctx context.Context, arg UpdateTransactionParams) error {
	_, err := q.db.ExecContext(ctx, updateTransaction,
		arg.TransactionID,
		arg.CardNumber,
		arg.Amount,
		arg.PaymentMethod,
		arg.MerchantID,
		arg.TransactionTime,
	)
	return err
}
