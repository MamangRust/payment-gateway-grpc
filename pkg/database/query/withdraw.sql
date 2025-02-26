-- Search Withdraws with Pagination
-- name: GetWithdraws :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    withdraws
WHERE
    deleted_at IS NULL
    AND ($1::TEXT IS NULL
        OR card_number ILIKE '%' || $1 || '%'
        OR withdraw_amount::TEXT ILIKE '%' || $1 || '%'
        OR withdraw_time::TEXT ILIKE '%' || $1 || '%'
        OR status ILIKE '%' || $1 || '%'
    )
ORDER BY
    withdraw_time DESC
LIMIT $2 OFFSET $3;



-- Get Withdraw by ID
-- name: GetWithdrawByID :one
SELECT *
FROM withdraws
WHERE
    withdraw_id = $1
    AND deleted_at IS NULL;

-- Get Active Withdraws with Search, Pagination, and Total Count
-- name: GetActiveWithdraws :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    withdraws
WHERE
    deleted_at IS NULL
    AND ($1::TEXT IS NULL
        OR card_number ILIKE '%' || $1 || '%'
        OR withdraw_amount::TEXT ILIKE '%' || $1 || '%'
        OR withdraw_time::TEXT ILIKE '%' || $1 || '%'
        OR status ILIKE '%' || $1 || '%'
    )
ORDER BY
    withdraw_time DESC
LIMIT $2 OFFSET $3;


-- Get Trashed Withdraws with Search, Pagination, and Total Count
-- name: GetTrashedWithdraws :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    withdraws
WHERE
    deleted_at IS NOT NULL
    AND ($1::TEXT IS NULL
        OR card_number ILIKE '%' || $1 || '%'
        OR withdraw_amount::TEXT ILIKE '%' || $1 || '%'
        OR withdraw_time::TEXT ILIKE '%' || $1 || '%'
        OR status ILIKE '%' || $1 || '%'
    )
ORDER BY
    withdraw_time DESC
LIMIT $2 OFFSET $3;


-- Search Withdraws by Card Number with Pagination
-- name: GetWithdrawsByCardNumber :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    withdraws
WHERE
    deleted_at IS NULL
    AND card_number = $1
    AND (
        $2::TEXT IS NULL
        OR CAST(withdraw_amount AS TEXT) ILIKE '%' || $2 || '%'
        OR TO_CHAR(withdraw_time, 'YYYY-MM-DD HH24:MI:SS') ILIKE '%' || $2 || '%'
        OR status ILIKE '%' || $2 || '%'
    )
ORDER BY
    withdraw_time DESC
LIMIT $3 OFFSET $4;


-- Get Trashed By Withdraw ID
-- name: GetTrashedWithdrawByID :one
SELECT *
FROM withdraws
WHERE
    withdraw_id = $1
    AND deleted_at IS NOT NULL;



-- name: FindAllWithdrawsByCardNumber :many
SELECT
    w.withdraw_id,
    w.card_number,
    w.withdraw_amount,
    w.withdraw_time,
    w.created_at,
    w.updated_at,
    w.deleted_at
FROM
    withdraws w
WHERE
    w.card_number = $1
    AND w.deleted_at IS NULL
ORDER BY
    w.withdraw_time DESC;



-- name: GetMonthWithdrawStatusSuccess :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.withdraw_time)::integer AS year,
        EXTRACT(MONTH FROM t.withdraw_time)::integer AS month,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
    FROM
        withdraws t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND (
            (t.withdraw_time >= $1::timestamp AND t.withdraw_time <= $2::timestamp)
            OR (t.withdraw_time >= $3::timestamp AND t.withdraw_time <= $4::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.withdraw_time),
        EXTRACT(MONTH FROM t.withdraw_time)
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



-- name: GetYearlyWithdrawStatusSuccess :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.withdraw_time)::integer AS year,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
    FROM
        withdraws t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND (
            EXTRACT(YEAR FROM t.withdraw_time) = $1::integer
            OR EXTRACT(YEAR FROM t.withdraw_time) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.withdraw_time)
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



-- name: GetMonthWithdrawStatusFailed :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.withdraw_time)::integer AS year,
        EXTRACT(MONTH FROM t.withdraw_time)::integer AS month,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
    FROM
        withdraws t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND (
            (t.withdraw_time >= $1::timestamp AND t.withdraw_time <= $2::timestamp)
            OR (t.withdraw_time >= $3::timestamp AND t.withdraw_time <= $4::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.withdraw_time),
        EXTRACT(MONTH FROM t.withdraw_time)
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



-- name: GetYearlyWithdrawStatusFailed :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.withdraw_time)::integer AS year,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
    FROM
        withdraws t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND (
            EXTRACT(YEAR FROM t.withdraw_time) = $1::integer
            OR EXTRACT(YEAR FROM t.withdraw_time) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.withdraw_time)
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



-- name: GetMonthWithdrawStatusSuccessCardNumber :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.withdraw_time)::integer AS year,
        EXTRACT(MONTH FROM t.withdraw_time)::integer AS month,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
    FROM
        withdraws t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND t.card_number = $1
        AND (
            (t.withdraw_time >= $2::timestamp AND t.withdraw_time <= $3::timestamp)
            OR (t.withdraw_time >= $4::timestamp AND t.withdraw_time <= $5::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.withdraw_time),
        EXTRACT(MONTH FROM t.withdraw_time)
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



-- name: GetYearlyWithdrawStatusSuccessCardNumber :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.withdraw_time)::integer AS year,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
    FROM
        withdraws t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND t.card_number = $1
        AND (
            EXTRACT(YEAR FROM t.withdraw_time) = $2::integer
            OR EXTRACT(YEAR FROM t.withdraw_time) = $2::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.withdraw_time)
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



-- name: GetMonthWithdrawStatusFailedCardNumber :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.withdraw_time)::integer AS year,
        EXTRACT(MONTH FROM t.withdraw_time)::integer AS month,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
    FROM
        withdraws t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND t.card_number = $1
        AND (
            (t.withdraw_time >= $2::timestamp AND t.withdraw_time <= $3::timestamp)
            OR (t.withdraw_time >= $4::timestamp AND t.withdraw_time <= $5::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.withdraw_time),
        EXTRACT(MONTH FROM t.withdraw_time)
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



-- name: GetYearlyWithdrawStatusFailedCardNumber :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.withdraw_time)::integer AS year,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
    FROM
        withdraws t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND t.card_number = $1
        AND (
            EXTRACT(YEAR FROM t.withdraw_time) = $2::integer
            OR EXTRACT(YEAR FROM t.withdraw_time) = $2::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.withdraw_time)
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



-- name: GetMonthlyWithdraws :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $1::timestamp),
        date_trunc('year', $1::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    COALESCE(SUM(w.withdraw_amount), 0)::int AS total_withdraw_amount
FROM
    months m
LEFT JOIN
    withdraws w ON EXTRACT(MONTH FROM w.withdraw_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM w.withdraw_time) = EXTRACT(YEAR FROM m.month)
    AND w.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY
    m.month;



-- name: GetYearlyWithdraws :many
SELECT
    EXTRACT(YEAR FROM w.withdraw_time) AS year,
    SUM(w.withdraw_amount) AS total_withdraw_amount
FROM
    withdraws w
WHERE
    w.deleted_at IS NULL
    AND EXTRACT(YEAR FROM w.withdraw_time) >= $1 - 4
    AND EXTRACT(YEAR FROM w.withdraw_time) <= $1
GROUP BY
    EXTRACT(YEAR FROM w.withdraw_time)
ORDER BY
    year;



-- name: GetMonthlyWithdrawsByCardNumber :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $2::timestamp),
        date_trunc('year', $2::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    COALESCE(SUM(w.withdraw_amount), 0)::int AS total_withdraw_amount
FROM
    months m
LEFT JOIN
    withdraws w ON EXTRACT(MONTH FROM w.withdraw_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM w.withdraw_time) = EXTRACT(YEAR FROM m.month)
    AND w.card_number = $1
    AND w.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY
    m.month;




-- name: GetYearlyWithdrawsByCardNumber :many
SELECT
    EXTRACT(YEAR FROM w.created_at) AS year,
    SUM(w.withdraw_amount) AS total_withdraw_amount
FROM
    withdraws w
WHERE
    w.deleted_at IS NULL
    AND w.card_number = $1
    AND EXTRACT(YEAR FROM w.created_at) >= $2 - 4
    AND EXTRACT(YEAR FROM w.created_at) <= $2
GROUP BY
    EXTRACT(YEAR FROM w.created_at)
ORDER BY
    year;



-- Create Withdraw
-- name: CreateWithdraw :one
INSERT INTO
    withdraws (
        card_number,
        withdraw_amount,
        withdraw_time,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        current_timestamp
    ) RETURNING *;


-- Update Withdraw
-- name: UpdateWithdraw :exec
UPDATE withdraws
SET
    card_number = $2,
    withdraw_amount = $3,
    withdraw_time = $4,
    updated_at = current_timestamp
WHERE
    withdraw_id = $1
    AND deleted_at IS NULL;


-- Update Withdraw Status
-- name: UpdateWithdrawStatus :exec
UPDATE withdraws
SET
    status = $2,
    updated_at = current_timestamp
WHERE
    withdraw_id = $1
    AND deleted_at IS NULL;



-- Trash Withdraw (Soft Delete)
-- name: TrashWithdraw :exec
UPDATE withdraws
SET
    deleted_at = current_timestamp
WHERE
    withdraw_id = $1
    AND deleted_at IS NULL;

-- Restore Withdraw (Undelete)
-- name: RestoreWithdraw :exec
UPDATE withdraws
SET
    deleted_at = NULL
WHERE
    withdraw_id = $1
    AND deleted_at IS NOT NULL;



-- Delete Withdraw Permanently
-- name: DeleteWithdrawPermanently :exec
DELETE FROM withdraws WHERE withdraw_id = $1 AND deleted_at IS NOT NULL;


-- Restore All Trashed Withdraws
-- name: RestoreAllWithdraws :exec
UPDATE withdraws
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;


-- Delete All Trashed Withdraws Permanently
-- name: DeleteAllPermanentWithdraws :exec
DELETE FROM withdraws
WHERE
    deleted_at IS NOT NULL;
