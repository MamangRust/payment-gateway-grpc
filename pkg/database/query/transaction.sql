-- Search Transactions with Pagination
-- name: GetTransactions :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    transactions
WHERE
    deleted_at IS NULL
    AND ($1::TEXT IS NULL
        OR card_number ILIKE '%' || $1 || '%'
        OR payment_method ILIKE '%' || $1 || '%'
        OR status ILIKE '%' || $1 || '%'
    )
ORDER BY
    transaction_time DESC
LIMIT $2 OFFSET $3;

-- Get Transaction by ID
-- name: GetTransactionByID :one
SELECT *
FROM transactions
WHERE
    transaction_id = $1
    AND deleted_at IS NULL;

-- Get Transactions by Card Number
-- name: GetTransactionsByCardNumber :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    transactions
WHERE
    deleted_at IS NULL
    AND card_number = $1
    AND (
        $2::TEXT IS NULL
        OR payment_method ILIKE '%' || $2 || '%'
        OR status ILIKE '%' || $2 || '%'
    )
ORDER BY
    transaction_time DESC
LIMIT $3 OFFSET $4;

-- Get Transactions by Merchant ID
-- name: GetTransactionsByMerchantID :many
SELECT *
FROM transactions
WHERE
    merchant_id = $1
    AND deleted_at IS NULL
ORDER BY transaction_time DESC;

-- Get Trashed By Transaction ID
-- name: GetTrashedTransactionByID :one
SELECT *
FROM transactions
WHERE
    transaction_id = $1
    AND deleted_at IS NOT NULL;


-- name: GetMonthTransactionStatusSuccess :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        EXTRACT(MONTH FROM t.transaction_time)::integer AS month,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND (
            (t.transaction_time >= $1::timestamp AND t.transaction_time <= $2::timestamp)
            OR (t.transaction_time >= $3::timestamp AND t.transaction_time <= $4::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time),
        EXTRACT(MONTH FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
        total_success,
        total_amount
    FROM
        monthly_data

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $1::timestamp)::text AS year,
        TO_CHAR($1::timestamp, 'Mon') AS month,
        0 AS total_success,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $1::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $1::timestamp)::integer
    )

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $3::timestamp)::text AS year,
        TO_CHAR($3::timestamp, 'Mon') AS month,
        0 AS total_success,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $3::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $3::timestamp)::integer
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC,
    TO_DATE(month, 'Mon') DESC;



-- name: GetYearlyTransactionStatusSuccess :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND (
            EXTRACT(YEAR FROM t.transaction_time) = $1::integer
            OR EXTRACT(YEAR FROM t.transaction_time) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        total_success::integer,
        total_amount::integer
    FROM
        yearly_data

    UNION ALL

    SELECT
        $1::text AS year,
        0::integer AS total_success,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer
    )

    UNION ALL

    SELECT
        ($1::integer - 1)::text AS year,
        0::integer AS total_success,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer - 1
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC;



-- name: GetMonthTransactionStatusFailed :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        EXTRACT(MONTH FROM t.transaction_time)::integer AS month,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND (
            (t.transaction_time >= $1::timestamp AND t.transaction_time <= $2::timestamp)
            OR (t.transaction_time >= $3::timestamp AND t.transaction_time <= $4::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time),
        EXTRACT(MONTH FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
        total_failed,
        total_amount
    FROM
        monthly_data

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $1::timestamp)::text AS year,
        TO_CHAR($1::timestamp, 'Mon') AS month,
        0 AS total_failed,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $1::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $1::timestamp)::integer
    )

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $3::timestamp)::text AS year,
        TO_CHAR($3::timestamp, 'Mon') AS month,
        0 AS total_failed,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $3::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $3::timestamp)::integer
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC,
    TO_DATE(month, 'Mon') DESC;


-- name: GetYearlyTransactionStatusFailed :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND (
            EXTRACT(YEAR FROM t.transaction_time) = $1::integer
            OR EXTRACT(YEAR FROM t.transaction_time) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        total_failed::integer,
        total_amount::integer
    FROM
        yearly_data

    UNION ALL

    SELECT
        $1::text AS year,
        0::integer AS total_failed,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer
    )

    UNION ALL

    SELECT
        ($1::integer - 1)::text AS year,
        0::integer AS total_failed,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer - 1
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC;


-- name: GetMonthTransactionStatusSuccessCardNumber :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        EXTRACT(MONTH FROM t.transaction_time)::integer AS month,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND t.card_number = $1
        AND (
            (t.transaction_time >= $2::timestamp AND t.transaction_time <= $3::timestamp)
            OR (t.transaction_time >= $4::timestamp AND t.transaction_time <= $5::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time),
        EXTRACT(MONTH FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
        total_success,
        total_amount
    FROM
        monthly_data

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $2::timestamp)::text AS year,
        TO_CHAR($2::timestamp, 'Mon') AS month,
        0 AS total_success,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $2::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $2::timestamp)::integer
    )

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $3::timestamp)::text AS year,
        TO_CHAR($3::timestamp, 'Mon') AS month,
        0 AS total_success,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $3::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $3::timestamp)::integer
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC,
    TO_DATE(month, 'Mon') DESC;



-- name: GetYearlyTransactionStatusSuccessCardNumber :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND t.card_number = $1
        AND (
            EXTRACT(YEAR FROM t.transaction_time) = $2::integer
            OR EXTRACT(YEAR FROM t.transaction_time) = $2::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        total_success::integer,
        total_amount::integer
    FROM
        yearly_data

    UNION ALL

    SELECT
        $2::text AS year,
        0::integer AS total_success,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $2::integer
    )

    UNION ALL

    SELECT
        ($2::integer - 1)::text AS year,
        0::integer AS total_success,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $2::integer - 1
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC;



-- name: GetMonthTransactionStatusFailedCardNumber :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        EXTRACT(MONTH FROM t.transaction_time)::integer AS month,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND t.card_number = $1
        AND (
            (t.transaction_time >= $2::timestamp AND t.transaction_time <= $3::timestamp)
            OR (t.transaction_time >= $4::timestamp AND t.transaction_time <= $5::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time),
        EXTRACT(MONTH FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
        total_failed,
        total_amount
    FROM
        monthly_data

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $2::timestamp)::text AS year,
        TO_CHAR($2::timestamp, 'Mon') AS month,
        0 AS total_failed,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $2::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $2::timestamp)::integer
    )

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $3::timestamp)::text AS year,
        TO_CHAR($3::timestamp, 'Mon') AS month,
        0 AS total_failed,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $3::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $3::timestamp)::integer
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC,
    TO_DATE(month, 'Mon') DESC;


-- name: GetYearlyTransactionStatusFailedCardNumber :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND t.card_number = $1
        AND (
            EXTRACT(YEAR FROM t.transaction_time) = $2::integer
            OR EXTRACT(YEAR FROM t.transaction_time) = $2::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        total_failed::integer,
        total_amount::integer
    FROM
        yearly_data

    UNION ALL

    SELECT
        $2::text AS year,
        0::integer AS total_failed,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $2::integer
    )

    UNION ALL

    SELECT
        ($2::integer - 1)::text AS year,
        0::integer AS total_failed,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $2::integer - 1
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC;


-- name: GetMonthlyPaymentMethods :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $1::timestamp),
        date_trunc('year', $1::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
),
payment_methods AS (
    SELECT DISTINCT payment_method
    FROM transactions
    WHERE deleted_at IS NULL
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    pm.payment_method,
    COALESCE(COUNT(t.transaction_id), 0)::int AS total_transactions,
    COALESCE(SUM(t.amount), 0)::int AS total_amount
FROM
    months m
CROSS JOIN
    payment_methods pm
LEFT JOIN
    transactions t ON EXTRACT(MONTH FROM t.transaction_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM t.transaction_time) = EXTRACT(YEAR FROM m.month)
    AND t.payment_method = pm.payment_method
    AND t.deleted_at IS NULL
GROUP BY
    m.month,
    pm.payment_method
ORDER BY
    m.month,
    pm.payment_method;


-- name: GetYearlyPaymentMethods :many
SELECT
    EXTRACT(YEAR FROM t.created_at) AS year,
    t.payment_method,
    COUNT(t.transaction_id) AS total_transactions,
    SUM(t.amount) AS total_amount
FROM
    transactions t
WHERE
    t.deleted_at IS NULL
    AND EXTRACT(YEAR FROM t.created_at) >= $1 - 4
    AND EXTRACT(YEAR FROM t.created_at) <= $1
GROUP BY
    EXTRACT(YEAR FROM t.created_at),
    t.payment_method
ORDER BY
    year;


-- name: GetMonthlyAmounts :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $1::timestamp),
        date_trunc('year', $1::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    COALESCE(SUM(t.amount), 0)::int AS total_amount
FROM
    months m
LEFT JOIN
    transactions t ON EXTRACT(MONTH FROM t.transaction_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM t.transaction_time) = EXTRACT(YEAR FROM m.month)
    AND t.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY
    m.month;


-- name: GetYearlyAmounts :many
SELECT
    EXTRACT(YEAR FROM t.created_at) AS year,
    SUM(t.amount) AS total_amount
FROM
    transactions t
WHERE
    t.deleted_at IS NULL
    AND EXTRACT(YEAR FROM t.created_at) >= $1 - 4
    AND EXTRACT(YEAR FROM t.created_at) <= $1
GROUP BY
    EXTRACT(YEAR FROM t.created_at)
ORDER BY
    year;



-- name: GetTransactionByCardNumber :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    transactions
WHERE
    deleted_at IS NULL
    AND card_number = $1
    AND ($2::TEXT IS NULL OR payment_method ILIKE '%' || $2 || '%')
ORDER BY
    transaction_time DESC
LIMIT $3 OFFSET $4;


-- name: GetMonthlyPaymentMethodsByCardNumber :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $2::timestamp),
        date_trunc('year', $2::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
),
payment_methods AS (
    SELECT DISTINCT payment_method
    FROM transactions
    WHERE deleted_at IS NULL
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    pm.payment_method,
    COALESCE(COUNT(t.transaction_id), 0)::int AS total_transactions,
    COALESCE(SUM(t.amount), 0)::int AS total_amount
FROM
    months m
CROSS JOIN
    payment_methods pm
LEFT JOIN
    transactions t ON EXTRACT(MONTH FROM t.transaction_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM t.transaction_time) = EXTRACT(YEAR FROM m.month)
    AND t.payment_method = pm.payment_method
    AND t.card_number = $1
    AND t.deleted_at IS NULL
GROUP BY
    m.month,
    pm.payment_method
ORDER BY
    m.month,
    pm.payment_method;




-- name: GetYearlyPaymentMethodsByCardNumber :many
SELECT
    EXTRACT(YEAR FROM t.created_at) AS year,
    t.payment_method,
    COUNT(t.transaction_id) AS total_transactions,
    SUM(t.amount) AS total_amount
FROM
    transactions t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
    AND EXTRACT(YEAR FROM t.created_at) >= $2 - 4
    AND EXTRACT(YEAR FROM t.created_at) <= $2
GROUP BY
    EXTRACT(YEAR FROM t.created_at),
    t.payment_method
ORDER BY
    year;



-- name: GetMonthlyAmountsByCardNumber :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $2::timestamp),
        date_trunc('year', $2::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    COALESCE(SUM(t.amount), 0)::int AS total_amount
FROM
    months m
LEFT JOIN
    transactions t ON EXTRACT(MONTH FROM t.transaction_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM t.transaction_time) = EXTRACT(YEAR FROM m.month)
    AND t.card_number = $1
    AND t.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY
    m.month;



-- name: GetYearlyAmountsByCardNumber :many
SELECT
    EXTRACT(YEAR FROM t.created_at) AS year,
    SUM(t.amount) AS total_amount
FROM
    transactions t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
    AND EXTRACT(YEAR FROM t.created_at) >= $2 - 4
    AND EXTRACT(YEAR FROM t.created_at) <= $2
GROUP BY
    EXTRACT(YEAR FROM t.created_at)
ORDER BY
    year;


-- Get Active Transactions with Pagination, Search, and Count
-- name: GetActiveTransactions :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    transactions
WHERE
    deleted_at IS NULL
    AND ($1::TEXT IS NULL OR card_number ILIKE '%' || $1 || '%' OR payment_method ILIKE '%' || $1 || '%')
ORDER BY
    transaction_time DESC
LIMIT $2 OFFSET $3;


-- Get Trashed Transactions with Pagination, Search, and Count
-- name: GetTrashedTransactions :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    transactions
WHERE
    deleted_at IS NOT NULL
    AND ($1::TEXT IS NULL OR card_number ILIKE '%' || $1 || '%' OR payment_method ILIKE '%' || $1 || '%')
ORDER BY
    transaction_time DESC
LIMIT $2 OFFSET $3;

-- Create Transaction
-- name: CreateTransaction :one
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
    ) RETURNING *;



-- Update Transaction
-- name: UpdateTransaction :exec
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
    AND deleted_at IS NULL;


-- Update Transaction Status
-- name: UpdateTransactionStatus :exec
UPDATE transactions
SET
    status = $2,
    updated_at = current_timestamp
WHERE
    transaction_id = $1
    AND deleted_at IS NULL;


-- Trash Transaction
-- name: TrashTransaction :exec
UPDATE transactions
SET
    deleted_at = current_timestamp
WHERE
    transaction_id = $1
    AND deleted_at IS NULL;


-- Restore Trashed Transaction
-- name: RestoreTransaction :exec
UPDATE transactions
SET
    deleted_at = NULL
WHERE
    transaction_id = $1
    AND deleted_at IS NOT NULL;


-- Delete Transaction Permanently
-- name: DeleteTransactionPermanently :exec
DELETE FROM transactions WHERE transaction_id = $1 AND deleted_at IS NOT NULL;


-- Restore All Trashed Transactions
-- name: RestoreAllTransactions :exec
UPDATE transactions
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;


-- Delete All Trashed Transactions Permanently
-- name: DeleteAllPermanentTransactions :exec
DELETE FROM transactions
WHERE
    deleted_at IS NOT NULL;
