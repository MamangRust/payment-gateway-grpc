// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: topup.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const countAllTopups = `-- name: CountAllTopups :one
SELECT COUNT(*) FROM topups WHERE deleted_at IS NULL
`

// Count All Topups
func (q *Queries) CountAllTopups(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countAllTopups)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countTopups = `-- name: CountTopups :one
SELECT COUNT(*)
FROM topups
WHERE deleted_at IS NULL
    AND ($1::TEXT IS NULL OR
        card_number ILIKE '%' || $1 || '%' OR
        topup_method ILIKE '%' || $1 || '%' OR
        topup_status ILIKE '%' || $1 || '%')
`

func (q *Queries) CountTopups(ctx context.Context, dollar_1 string) (int64, error) {
	row := q.db.QueryRowContext(ctx, countTopups, dollar_1)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countTopupsByDate = `-- name: CountTopupsByDate :one
SELECT COUNT(*)
FROM topups
WHERE deleted_at IS NULL
  AND topup_time::DATE = $1::DATE
`

// Count Topups by Date
func (q *Queries) CountTopupsByDate(ctx context.Context, dollar_1 time.Time) (int64, error) {
	row := q.db.QueryRowContext(ctx, countTopupsByDate, dollar_1)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createTopup = `-- name: CreateTopup :one
INSERT INTO
    topups (
        card_number,
        topup_no,
        topup_amount,
        topup_method,
        topup_time,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        current_timestamp,
        current_timestamp,
        current_timestamp
    ) RETURNING topup_id, card_number, topup_no, topup_amount, topup_method, topup_time, created_at, updated_at, deleted_at
`

type CreateTopupParams struct {
	CardNumber  string `json:"card_number"`
	TopupNo     string `json:"topup_no"`
	TopupAmount int32  `json:"topup_amount"`
	TopupMethod string `json:"topup_method"`
}

// Create Topup
func (q *Queries) CreateTopup(ctx context.Context, arg CreateTopupParams) (*Topup, error) {
	row := q.db.QueryRowContext(ctx, createTopup,
		arg.CardNumber,
		arg.TopupNo,
		arg.TopupAmount,
		arg.TopupMethod,
	)
	var i Topup
	err := row.Scan(
		&i.TopupID,
		&i.CardNumber,
		&i.TopupNo,
		&i.TopupAmount,
		&i.TopupMethod,
		&i.TopupTime,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const deleteAllPermanentTopups = `-- name: DeleteAllPermanentTopups :exec
DELETE FROM topups
WHERE
    deleted_at IS NOT NULL
`

// Delete All Trashed Saldos Permanently
func (q *Queries) DeleteAllPermanentTopups(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllPermanentTopups)
	return err
}

const deleteTopupPermanently = `-- name: DeleteTopupPermanently :exec
DELETE FROM topups WHERE topup_id = $1
`

// Delete Topup Permanently
func (q *Queries) DeleteTopupPermanently(ctx context.Context, topupID int32) error {
	_, err := q.db.ExecContext(ctx, deleteTopupPermanently, topupID)
	return err
}

const getActiveTopups = `-- name: GetActiveTopups :many
SELECT
    topup_id, card_number, topup_no, topup_amount, topup_method, topup_time, created_at, updated_at, deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    topups
WHERE
    deleted_at IS NULL
    AND ($1::TEXT IS NULL OR card_number ILIKE '%' || $1 || '%' OR topup_no ILIKE '%' || $1 || '%' OR topup_method ILIKE '%' || $1 || '%')
ORDER BY
    topup_time DESC
LIMIT $2 OFFSET $3
`

type GetActiveTopupsParams struct {
	Column1 string `json:"column_1"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

type GetActiveTopupsRow struct {
	TopupID     int32        `json:"topup_id"`
	CardNumber  string       `json:"card_number"`
	TopupNo     string       `json:"topup_no"`
	TopupAmount int32        `json:"topup_amount"`
	TopupMethod string       `json:"topup_method"`
	TopupTime   time.Time    `json:"topup_time"`
	CreatedAt   sql.NullTime `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
	TotalCount  int64        `json:"total_count"`
}

// Get All Active Topups with Pagination and Search
func (q *Queries) GetActiveTopups(ctx context.Context, arg GetActiveTopupsParams) ([]*GetActiveTopupsRow, error) {
	rows, err := q.db.QueryContext(ctx, getActiveTopups, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetActiveTopupsRow
	for rows.Next() {
		var i GetActiveTopupsRow
		if err := rows.Scan(
			&i.TopupID,
			&i.CardNumber,
			&i.TopupNo,
			&i.TopupAmount,
			&i.TopupMethod,
			&i.TopupTime,
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

const getMonthlyTopupAmounts = `-- name: GetMonthlyTopupAmounts :many
SELECT
    TO_CHAR(t.topup_time, 'Mon') AS month,
    SUM(t.topup_amount) AS total_amount
FROM
    topups t
WHERE
    t.deleted_at IS NULL
    AND EXTRACT(YEAR FROM t.topup_time) = $1
GROUP BY
    TO_CHAR(t.topup_time, 'Mon'),
    EXTRACT(MONTH FROM t.topup_time)
ORDER BY
    EXTRACT(MONTH FROM t.topup_time)
`

type GetMonthlyTopupAmountsRow struct {
	Month       string `json:"month"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) GetMonthlyTopupAmounts(ctx context.Context, topupTime time.Time) ([]*GetMonthlyTopupAmountsRow, error) {
	rows, err := q.db.QueryContext(ctx, getMonthlyTopupAmounts, topupTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetMonthlyTopupAmountsRow
	for rows.Next() {
		var i GetMonthlyTopupAmountsRow
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

const getMonthlyTopupAmountsByCardNumber = `-- name: GetMonthlyTopupAmountsByCardNumber :many
SELECT
    TO_CHAR(t.topup_time, 'Mon') AS month,
    SUM(t.topup_amount) AS total_amount
FROM
    topups t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
    AND EXTRACT(YEAR FROM t.topup_time) = $2
GROUP BY
    TO_CHAR(t.topup_time, 'Mon'),
    EXTRACT(MONTH FROM t.topup_time)
ORDER BY
    EXTRACT(MONTH FROM t.topup_time)
`

type GetMonthlyTopupAmountsByCardNumberParams struct {
	CardNumber string    `json:"card_number"`
	TopupTime  time.Time `json:"topup_time"`
}

type GetMonthlyTopupAmountsByCardNumberRow struct {
	Month       string `json:"month"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) GetMonthlyTopupAmountsByCardNumber(ctx context.Context, arg GetMonthlyTopupAmountsByCardNumberParams) ([]*GetMonthlyTopupAmountsByCardNumberRow, error) {
	rows, err := q.db.QueryContext(ctx, getMonthlyTopupAmountsByCardNumber, arg.CardNumber, arg.TopupTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetMonthlyTopupAmountsByCardNumberRow
	for rows.Next() {
		var i GetMonthlyTopupAmountsByCardNumberRow
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

const getMonthlyTopupMethods = `-- name: GetMonthlyTopupMethods :many
SELECT
    TO_CHAR(t.topup_time, 'Mon') AS month,
    t.topup_method,
    COUNT(t.topup_id) AS total_topups,
    SUM(t.topup_amount) AS total_amount
FROM
    topups t
WHERE
    t.deleted_at IS NULL
    AND EXTRACT(YEAR FROM t.topup_time) = $1
GROUP BY
    TO_CHAR(t.topup_time, 'Mon'),
    EXTRACT(MONTH FROM t.topup_time),
    t.topup_method
ORDER BY
    EXTRACT(MONTH FROM t.topup_time)
`

type GetMonthlyTopupMethodsRow struct {
	Month       string `json:"month"`
	TopupMethod string `json:"topup_method"`
	TotalTopups int64  `json:"total_topups"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) GetMonthlyTopupMethods(ctx context.Context, topupTime time.Time) ([]*GetMonthlyTopupMethodsRow, error) {
	rows, err := q.db.QueryContext(ctx, getMonthlyTopupMethods, topupTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetMonthlyTopupMethodsRow
	for rows.Next() {
		var i GetMonthlyTopupMethodsRow
		if err := rows.Scan(
			&i.Month,
			&i.TopupMethod,
			&i.TotalTopups,
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

const getMonthlyTopupMethodsByCardNumber = `-- name: GetMonthlyTopupMethodsByCardNumber :many
SELECT
    TO_CHAR(t.topup_time, 'Mon') AS month,
    t.topup_method,
    COUNT(t.topup_id) AS total_topups,
    SUM(t.topup_amount) AS total_amount
FROM
    topups t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
    AND EXTRACT(YEAR FROM t.topup_time) = $2
GROUP BY
    TO_CHAR(t.topup_time, 'Mon'),
    EXTRACT(MONTH FROM t.topup_time),
    t.topup_method
ORDER BY
    EXTRACT(MONTH FROM t.topup_time)
`

type GetMonthlyTopupMethodsByCardNumberParams struct {
	CardNumber string    `json:"card_number"`
	TopupTime  time.Time `json:"topup_time"`
}

type GetMonthlyTopupMethodsByCardNumberRow struct {
	Month       string `json:"month"`
	TopupMethod string `json:"topup_method"`
	TotalTopups int64  `json:"total_topups"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) GetMonthlyTopupMethodsByCardNumber(ctx context.Context, arg GetMonthlyTopupMethodsByCardNumberParams) ([]*GetMonthlyTopupMethodsByCardNumberRow, error) {
	rows, err := q.db.QueryContext(ctx, getMonthlyTopupMethodsByCardNumber, arg.CardNumber, arg.TopupTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetMonthlyTopupMethodsByCardNumberRow
	for rows.Next() {
		var i GetMonthlyTopupMethodsByCardNumberRow
		if err := rows.Scan(
			&i.Month,
			&i.TopupMethod,
			&i.TotalTopups,
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

const getTopupByID = `-- name: GetTopupByID :one
SELECT topup_id, card_number, topup_no, topup_amount, topup_method, topup_time, created_at, updated_at, deleted_at FROM topups WHERE topup_id = $1 AND deleted_at IS NULL
`

// Get Topup by ID
func (q *Queries) GetTopupByID(ctx context.Context, topupID int32) (*Topup, error) {
	row := q.db.QueryRowContext(ctx, getTopupByID, topupID)
	var i Topup
	err := row.Scan(
		&i.TopupID,
		&i.CardNumber,
		&i.TopupNo,
		&i.TopupAmount,
		&i.TopupMethod,
		&i.TopupTime,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const getTopups = `-- name: GetTopups :many
SELECT
    topup_id, card_number, topup_no, topup_amount, topup_method, topup_time, created_at, updated_at, deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    topups
WHERE
    deleted_at IS NULL
    AND ($1::TEXT IS NULL OR card_number ILIKE '%' || $1 || '%' OR topup_no ILIKE '%' || $1 || '%' OR topup_method ILIKE '%' || $1 || '%')
ORDER BY
    topup_time DESC
LIMIT $2 OFFSET $3
`

type GetTopupsParams struct {
	Column1 string `json:"column_1"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

type GetTopupsRow struct {
	TopupID     int32        `json:"topup_id"`
	CardNumber  string       `json:"card_number"`
	TopupNo     string       `json:"topup_no"`
	TopupAmount int32        `json:"topup_amount"`
	TopupMethod string       `json:"topup_method"`
	TopupTime   time.Time    `json:"topup_time"`
	CreatedAt   sql.NullTime `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
	TotalCount  int64        `json:"total_count"`
}

// Search Topups with Pagination
func (q *Queries) GetTopups(ctx context.Context, arg GetTopupsParams) ([]*GetTopupsRow, error) {
	rows, err := q.db.QueryContext(ctx, getTopups, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetTopupsRow
	for rows.Next() {
		var i GetTopupsRow
		if err := rows.Scan(
			&i.TopupID,
			&i.CardNumber,
			&i.TopupNo,
			&i.TopupAmount,
			&i.TopupMethod,
			&i.TopupTime,
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

const getTopupsByCardNumber = `-- name: GetTopupsByCardNumber :many
SELECT topup_id, card_number, topup_no, topup_amount, topup_method, topup_time, created_at, updated_at, deleted_at
FROM topups
WHERE
    deleted_at IS NULL
    AND card_number = $1
ORDER BY topup_time DESC
`

// Get Topups by Card Number
func (q *Queries) GetTopupsByCardNumber(ctx context.Context, cardNumber string) ([]*Topup, error) {
	rows, err := q.db.QueryContext(ctx, getTopupsByCardNumber, cardNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Topup
	for rows.Next() {
		var i Topup
		if err := rows.Scan(
			&i.TopupID,
			&i.CardNumber,
			&i.TopupNo,
			&i.TopupAmount,
			&i.TopupMethod,
			&i.TopupTime,
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

const getTrashedTopupByID = `-- name: GetTrashedTopupByID :one
SELECT topup_id, card_number, topup_no, topup_amount, topup_method, topup_time, created_at, updated_at, deleted_at
FROM topups
WHERE
    topup_id = $1
    AND deleted_at IS NOT NULL
`

// Get Trashed By Topup ID
func (q *Queries) GetTrashedTopupByID(ctx context.Context, topupID int32) (*Topup, error) {
	row := q.db.QueryRowContext(ctx, getTrashedTopupByID, topupID)
	var i Topup
	err := row.Scan(
		&i.TopupID,
		&i.CardNumber,
		&i.TopupNo,
		&i.TopupAmount,
		&i.TopupMethod,
		&i.TopupTime,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const getTrashedTopups = `-- name: GetTrashedTopups :many
SELECT
    topup_id, card_number, topup_no, topup_amount, topup_method, topup_time, created_at, updated_at, deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    topups
WHERE
    deleted_at IS NOT NULL
    AND ($1::TEXT IS NULL OR card_number ILIKE '%' || $1 || '%' OR topup_no ILIKE '%' || $1 || '%' OR topup_method ILIKE '%' || $1 || '%')
ORDER BY
    topup_time DESC
LIMIT $2 OFFSET $3
`

type GetTrashedTopupsParams struct {
	Column1 string `json:"column_1"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

type GetTrashedTopupsRow struct {
	TopupID     int32        `json:"topup_id"`
	CardNumber  string       `json:"card_number"`
	TopupNo     string       `json:"topup_no"`
	TopupAmount int32        `json:"topup_amount"`
	TopupMethod string       `json:"topup_method"`
	TopupTime   time.Time    `json:"topup_time"`
	CreatedAt   sql.NullTime `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
	TotalCount  int64        `json:"total_count"`
}

// Get Trashed Topups with Pagination and Search
func (q *Queries) GetTrashedTopups(ctx context.Context, arg GetTrashedTopupsParams) ([]*GetTrashedTopupsRow, error) {
	rows, err := q.db.QueryContext(ctx, getTrashedTopups, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetTrashedTopupsRow
	for rows.Next() {
		var i GetTrashedTopupsRow
		if err := rows.Scan(
			&i.TopupID,
			&i.CardNumber,
			&i.TopupNo,
			&i.TopupAmount,
			&i.TopupMethod,
			&i.TopupTime,
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

const getYearlyTopupAmounts = `-- name: GetYearlyTopupAmounts :many
SELECT
    EXTRACT(YEAR FROM t.topup_time) AS year,
    SUM(t.topup_amount) AS total_amount
FROM
    topups t
WHERE
    t.deleted_at IS NULL
GROUP BY
    EXTRACT(YEAR FROM t.topup_time)
ORDER BY
    year
`

type GetYearlyTopupAmountsRow struct {
	Year        string `json:"year"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) GetYearlyTopupAmounts(ctx context.Context) ([]*GetYearlyTopupAmountsRow, error) {
	rows, err := q.db.QueryContext(ctx, getYearlyTopupAmounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetYearlyTopupAmountsRow
	for rows.Next() {
		var i GetYearlyTopupAmountsRow
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

const getYearlyTopupAmountsByCardNumber = `-- name: GetYearlyTopupAmountsByCardNumber :many
SELECT
    EXTRACT(YEAR FROM t.topup_time) AS year,
    SUM(t.topup_amount) AS total_amount
FROM
    topups t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
GROUP BY
    EXTRACT(YEAR FROM t.topup_time)
ORDER BY
    year
`

type GetYearlyTopupAmountsByCardNumberRow struct {
	Year        string `json:"year"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) GetYearlyTopupAmountsByCardNumber(ctx context.Context, cardNumber string) ([]*GetYearlyTopupAmountsByCardNumberRow, error) {
	rows, err := q.db.QueryContext(ctx, getYearlyTopupAmountsByCardNumber, cardNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetYearlyTopupAmountsByCardNumberRow
	for rows.Next() {
		var i GetYearlyTopupAmountsByCardNumberRow
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

const getYearlyTopupMethods = `-- name: GetYearlyTopupMethods :many
SELECT
    EXTRACT(YEAR FROM t.topup_time) AS year,
    t.topup_method,
    COUNT(t.topup_id) AS total_topups,
    SUM(t.topup_amount) AS total_amount
FROM
    topups t
WHERE
    t.deleted_at IS NULL
GROUP BY
    EXTRACT(YEAR FROM t.topup_time),
    t.topup_method
ORDER BY
    year
`

type GetYearlyTopupMethodsRow struct {
	Year        string `json:"year"`
	TopupMethod string `json:"topup_method"`
	TotalTopups int64  `json:"total_topups"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) GetYearlyTopupMethods(ctx context.Context) ([]*GetYearlyTopupMethodsRow, error) {
	rows, err := q.db.QueryContext(ctx, getYearlyTopupMethods)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetYearlyTopupMethodsRow
	for rows.Next() {
		var i GetYearlyTopupMethodsRow
		if err := rows.Scan(
			&i.Year,
			&i.TopupMethod,
			&i.TotalTopups,
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

const getYearlyTopupMethodsByCardNumber = `-- name: GetYearlyTopupMethodsByCardNumber :many
SELECT
    EXTRACT(YEAR FROM t.topup_time) AS year,
    t.topup_method,
    COUNT(t.topup_id) AS total_topups,
    SUM(t.topup_amount) AS total_amount
FROM
    topups t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
GROUP BY
    EXTRACT(YEAR FROM t.topup_time),
    t.topup_method
ORDER BY
    year
`

type GetYearlyTopupMethodsByCardNumberRow struct {
	Year        string `json:"year"`
	TopupMethod string `json:"topup_method"`
	TotalTopups int64  `json:"total_topups"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) GetYearlyTopupMethodsByCardNumber(ctx context.Context, cardNumber string) ([]*GetYearlyTopupMethodsByCardNumberRow, error) {
	rows, err := q.db.QueryContext(ctx, getYearlyTopupMethodsByCardNumber, cardNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetYearlyTopupMethodsByCardNumberRow
	for rows.Next() {
		var i GetYearlyTopupMethodsByCardNumberRow
		if err := rows.Scan(
			&i.Year,
			&i.TopupMethod,
			&i.TotalTopups,
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

const restoreAllTopups = `-- name: RestoreAllTopups :exec
UPDATE topups
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL
`

// Restore All Trashed Saldos
func (q *Queries) RestoreAllTopups(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, restoreAllTopups)
	return err
}

const restoreTopup = `-- name: RestoreTopup :exec
UPDATE topups
SET
    deleted_at = NULL
WHERE
    topup_id = $1
    AND deleted_at IS NOT NULL
`

// Restore Trashed Topup
func (q *Queries) RestoreTopup(ctx context.Context, topupID int32) error {
	_, err := q.db.ExecContext(ctx, restoreTopup, topupID)
	return err
}

const topup_CountAll = `-- name: Topup_CountAll :one
SELECT COUNT(*)
FROM topups
WHERE deleted_at IS NULL
`

func (q *Queries) Topup_CountAll(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, topup_CountAll)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const trashTopup = `-- name: TrashTopup :exec
UPDATE topups
SET
    deleted_at = current_timestamp
WHERE
    topup_id = $1
    AND deleted_at IS NULL
`

// Trash Topup
func (q *Queries) TrashTopup(ctx context.Context, topupID int32) error {
	_, err := q.db.ExecContext(ctx, trashTopup, topupID)
	return err
}

const updateTopup = `-- name: UpdateTopup :exec
UPDATE topups
SET
    card_number = $2,
    topup_amount = $3,
    topup_method = $4,
    topup_time = current_timestamp,
    updated_at = current_timestamp
WHERE
    topup_id = $1
    AND deleted_at IS NULL
`

type UpdateTopupParams struct {
	TopupID     int32  `json:"topup_id"`
	CardNumber  string `json:"card_number"`
	TopupAmount int32  `json:"topup_amount"`
	TopupMethod string `json:"topup_method"`
}

// Update Topup
func (q *Queries) UpdateTopup(ctx context.Context, arg UpdateTopupParams) error {
	_, err := q.db.ExecContext(ctx, updateTopup,
		arg.TopupID,
		arg.CardNumber,
		arg.TopupAmount,
		arg.TopupMethod,
	)
	return err
}

const updateTopupAmount = `-- name: UpdateTopupAmount :exec
UPDATE topups
SET
    topup_amount = $2,
    updated_at = current_timestamp
WHERE
    topup_id = $1
    AND deleted_at IS NULL
`

type UpdateTopupAmountParams struct {
	TopupID     int32 `json:"topup_id"`
	TopupAmount int32 `json:"topup_amount"`
}

// Update Topup Amount
func (q *Queries) UpdateTopupAmount(ctx context.Context, arg UpdateTopupAmountParams) error {
	_, err := q.db.ExecContext(ctx, updateTopupAmount, arg.TopupID, arg.TopupAmount)
	return err
}
