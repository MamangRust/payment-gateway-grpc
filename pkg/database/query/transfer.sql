-- Search Transfers with Pagination
-- name: GetTransfers :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    transfers
WHERE
    deleted_at IS NULL
    AND ($1::TEXT IS NULL OR transfer_from ILIKE '%' || $1 || '%' OR transfer_to ILIKE '%' || $1 || '%')
ORDER BY
    transfer_time DESC
LIMIT $2 OFFSET $3;


-- Get Transfer by ID
-- name: GetTransferByID :one
SELECT *
FROM transfers
WHERE
    transfer_id = $1
    AND deleted_at IS NULL;

-- Get Active Transfers with Search, Pagination, and Total Count
-- name: GetActiveTransfers :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    transfers
WHERE
    deleted_at IS NULL
    AND ($1::TEXT IS NULL OR transfer_from ILIKE '%' || $1 || '%' OR transfer_to ILIKE '%' || $1 || '%')
ORDER BY
    transfer_time DESC
LIMIT $2 OFFSET $3;

-- Get Trashed Transfers with Search, Pagination, and Total Count
-- name: GetTrashedTransfers :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    transfers
WHERE
    deleted_at IS NOT NULL
    AND ($1::TEXT IS NULL OR transfer_from ILIKE '%' || $1 || '%' OR transfer_to ILIKE '%' || $1 || '%')
ORDER BY
    transfer_time DESC
LIMIT $2 OFFSET $3;

-- Get Transfers by Card Number (Source or Destination)
-- name: GetTransfersByCardNumber :many
SELECT *
FROM transfers
WHERE
    deleted_at IS NULL
    AND (
        transfer_from = $1
        OR transfer_to = $1
    )
ORDER BY transfer_time DESC;

-- Get Transfers by Source Card
-- name: GetTransfersBySourceCard :many
SELECT *
FROM transfers
WHERE
    deleted_at IS NULL
    AND transfer_from = $1
ORDER BY transfer_time DESC;

-- Get Transfers by Destination Card
-- name: GetTransfersByDestinationCard :many
SELECT *
FROM transfers
WHERE
    deleted_at IS NULL
    AND transfer_to = $1
ORDER BY transfer_time DESC;

-- Get Trashed By Transfer ID
-- name: GetTrashedTransferByID :one
SELECT *
FROM transfers
WHERE
    transfer_id = $1
    AND deleted_at IS NOT NULL;


-- name: GetMonthTransferStatusSuccess :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transfer_time)::integer AS year,
        EXTRACT(MONTH FROM t.transfer_time)::integer AS month,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
    FROM
        transfers t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND (
            (t.transfer_time >= $1::timestamp AND t.transfer_time <= $2::timestamp)
            OR (t.transfer_time >= $3::timestamp AND t.transfer_time <= $4::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transfer_time),
        EXTRACT(MONTH FROM t.transfer_time)
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



-- name: GetYearlyTransferStatusSuccess :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transfer_time)::integer AS year,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
    FROM
        transfers t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND (
            EXTRACT(YEAR FROM t.transfer_time) = $1::integer
            OR EXTRACT(YEAR FROM t.transfer_time) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transfer_time)
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



-- name: GetMonthTransferStatusFailed :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transfer_time)::integer AS year,
        EXTRACT(MONTH FROM t.transfer_time)::integer AS month,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
    FROM
        transfers t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND (
            (t.transfer_time >= $1::timestamp AND t.transfer_time <= $2::timestamp)
            OR (t.transfer_time >= $3::timestamp AND t.transfer_time <= $4::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transfer_time),
        EXTRACT(MONTH FROM t.transfer_time)
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



-- name: GetYearlyTransferStatusFailed :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transfer_time)::integer AS year,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
    FROM
        transfers t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND (
            EXTRACT(YEAR FROM t.transfer_time) = $1::integer
            OR EXTRACT(YEAR FROM t.transfer_time) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transfer_time)
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



-- name: GetMonthTransferStatusSuccessCardNumber :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transfer_time)::integer AS year,
        EXTRACT(MONTH FROM t.transfer_time)::integer AS month,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
    FROM
        transfers t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND (t.transfer_from = $1 OR t.transfer_to = $1)
        AND (
            (t.transfer_time >= $2::timestamp AND t.transfer_time <= $3::timestamp)
            OR (t.transfer_time >= $4::timestamp AND t.transfer_time <= $5::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transfer_time),
        EXTRACT(MONTH FROM t.transfer_time)
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



-- name: GetYearlyTransferStatusSuccessCardNumber :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transfer_time)::integer AS year,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
    FROM
        transfers t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND (t.transfer_from = $1 OR t.transfer_to = $1)
        AND (
            EXTRACT(YEAR FROM t.transfer_time) = $2::integer
            OR EXTRACT(YEAR FROM t.transfer_time) = $2::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transfer_time)
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



-- name: GetMonthTransferStatusFailedCardNumber :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transfer_time)::integer AS year,
        EXTRACT(MONTH FROM t.transfer_time)::integer AS month,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
    FROM
        transfers t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND (t.transfer_from = $1 OR t.transfer_to = $1)
        AND (
            (t.transfer_time >= $2::timestamp AND t.transfer_time <= $3::timestamp)
            OR (t.transfer_time >= $4::timestamp AND t.transfer_time <= $5::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transfer_time),
        EXTRACT(MONTH FROM t.transfer_time)
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



-- name: GetYearlyTransferStatusFailedCardNumber :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transfer_time)::integer AS year,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
    FROM
        transfers t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND (t.transfer_from = $1 OR t.transfer_to = $1)
        AND (
            EXTRACT(YEAR FROM t.transfer_time) = $2::integer
            OR EXTRACT(YEAR FROM t.transfer_time) = $2::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transfer_time)
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



-- name: GetMonthlyTransferAmounts :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $1::timestamp),
        date_trunc('year', $1::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    COALESCE(SUM(t.transfer_amount), 0)::int AS total_transfer_amount
FROM
    months m
LEFT JOIN
    transfers t ON EXTRACT(MONTH FROM t.transfer_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM t.transfer_time) = EXTRACT(YEAR FROM m.month)
    AND t.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY
    m.month;



-- name: GetYearlyTransferAmounts :many
SELECT
    EXTRACT(YEAR FROM t.created_at) AS year,
    SUM(t.transfer_amount) AS total_transfer_amount
FROM
    transfers t
WHERE
    t.deleted_at IS NULL
    AND EXTRACT(YEAR FROM t.created_at) >= $1 - 4
    AND EXTRACT(YEAR FROM t.created_at) <= $1
GROUP BY
    EXTRACT(YEAR FROM t.created_at)
ORDER BY
    year;



-- name: GetMonthlyTransferAmountsBySenderCardNumber :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $2::timestamp),
        date_trunc('year', $2::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    COALESCE(SUM(t.transfer_amount), 0)::int AS total_transfer_amount
FROM
    months m
LEFT JOIN
    transfers t ON EXTRACT(MONTH FROM t.transfer_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM t.transfer_time) = EXTRACT(YEAR FROM m.month)
    AND t.transfer_from = $1
    AND t.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY
    m.month;




-- name: GetMonthlyTransferAmountsByReceiverCardNumber :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $2::timestamp),
        date_trunc('year', $2::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    COALESCE(SUM(t.transfer_amount), 0)::int AS total_transfer_amount
FROM
    months m
LEFT JOIN
    transfers t ON EXTRACT(MONTH FROM t.transfer_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM t.transfer_time) = EXTRACT(YEAR FROM m.month)
    AND t.transfer_to = $1
    AND t.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY
    m.month;




-- name: GetYearlyTransferAmountsBySenderCardNumber :many
SELECT
    EXTRACT(YEAR FROM t.created_at) AS year,
    SUM(t.transfer_amount) AS total_transfer_amount
FROM
    transfers t
WHERE
    t.deleted_at IS NULL
    AND t.transfer_from = $1
    AND EXTRACT(YEAR FROM t.created_at) >= $2 - 4
    AND EXTRACT(YEAR FROM t.created_at) <= $2
GROUP BY
    EXTRACT(YEAR FROM t.created_at)
ORDER BY
    year;



-- name: GetYearlyTransferAmountsByReceiverCardNumber :many
SELECT
    EXTRACT(YEAR FROM t.created_at) AS year,
    SUM(t.transfer_amount) AS total_transfer_amount
FROM
    transfers t
WHERE
    t.deleted_at IS NULL
    AND t.transfer_to = $1
    AND EXTRACT(YEAR FROM t.created_at) >= $2 - 4
    AND EXTRACT(YEAR FROM t.created_at) <= $2
GROUP BY
    EXTRACT(YEAR FROM t.created_at)
ORDER BY
    year;



-- name: FindAllTransfersByCardNumberAsSender :many
SELECT
    t.transfer_id,
    t.transfer_from,
    t.transfer_to,
    t.transfer_amount,
    t.transfer_time,
    t.created_at,
    t.updated_at,
    t.deleted_at
FROM
    transfers t
WHERE
    t.transfer_from = $1
    AND t.deleted_at IS NULL
ORDER BY
    t.transfer_time DESC;


-- name: FindAllTransfersByCardNumberAsReceiver :many
SELECT
    t.transfer_id,
    t.transfer_from,
    t.transfer_to,
    t.transfer_amount,
    t.transfer_time,
    t.created_at,
    t.updated_at,
    t.deleted_at
FROM
    transfers t
WHERE
    t.transfer_to = $1
    AND t.deleted_at IS NULL
ORDER BY
    t.transfer_time DESC;

-- Create Transfer
-- name: CreateTransfer :one
INSERT INTO
    transfers (
        transfer_from,
        transfer_to,
        transfer_amount,
        transfer_time,
        status,
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


-- Update Transfer
-- name: UpdateTransfer :exec
UPDATE transfers
SET
    transfer_from = $2,
    transfer_to = $3,
    transfer_amount = $4,
    transfer_time = $5,
    updated_at = current_timestamp
WHERE
    transfer_id = $1
    AND deleted_at IS NULL;

-- Update Transfer Amount
-- name: UpdateTransferAmount :exec
UPDATE transfers
SET
    transfer_amount = $2,
    transfer_time = current_timestamp,
    updated_at = current_timestamp
WHERE
    transfer_id = $1
    AND deleted_at IS NULL;


-- Update Transfer Status
-- name: UpdateTransferStatus :exec
UPDATE transfers
SET
    status = $2,
    updated_at = current_timestamp
WHERE
    transfer_id = $1
    AND deleted_at IS NULL;



-- Trash Transfer
-- name: TrashTransfer :exec
UPDATE transfers
SET
    deleted_at = current_timestamp
WHERE
    transfer_id = $1
    AND deleted_at IS NULL;

-- Restore Trashed Transfer
-- name: RestoreTransfer :exec
UPDATE transfers
SET
    deleted_at = NULL
WHERE
    transfer_id = $1
    AND deleted_at IS NOT NULL;

-- Delete Transfer Permanently
-- name: DeleteTransferPermanently :exec
DELETE FROM transfers WHERE transfer_id = $1 AND deleted_at IS NOT NULL;


-- Restore All Trashed Transfers
-- name: RestoreAllTransfers :exec
UPDATE transfers
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;


-- Delete All Trashed Transfers Permanently
-- name: DeleteAllPermanentTransfers :exec
DELETE FROM transfers
WHERE
    deleted_at IS NOT NULL;
