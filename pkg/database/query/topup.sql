-- Search Topups with Pagination
-- name: GetTopups :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    topups
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR card_number ILIKE '%' || $1 || '%'
        OR topup_no::TEXT ILIKE '%' || $1 || '%'
        OR topup_method ILIKE '%' || $1 || '%'
        OR status ILIKE '%' || $1 || '%'
    )
ORDER BY
    topup_time DESC
LIMIT $2 OFFSET $3;



-- Get Topup by ID
-- name: GetTopupByID :one
SELECT * FROM topups WHERE topup_id = $1 AND deleted_at IS NULL;


-- Get All Active Topups with Pagination and Search
-- name: GetActiveTopups :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    topups
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR card_number ILIKE '%' || $1 || '%'
        OR topup_no::TEXT ILIKE '%' || $1 || '%'
        OR topup_method ILIKE '%' || $1 || '%'
    )
ORDER BY
    topup_time DESC
LIMIT $2 OFFSET $3;

-- Get Trashed Topups with Pagination and Search
-- name: GetTrashedTopups :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    topups
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR card_number ILIKE '%' || $1 || '%'
        OR topup_no::TEXT ILIKE '%' || $1 || '%'
        OR topup_method ILIKE '%' || $1 || '%'
    )
ORDER BY
    topup_time DESC
LIMIT $2 OFFSET $3;




-- name: GetMonthTopupStatusSuccess :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.topup_time)::integer AS year,
        EXTRACT(MONTH FROM t.topup_time)::integer AS month,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.topup_amount), 0)::integer AS total_amount
    FROM
        topups t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND (
            (t.topup_time >= $1::timestamp AND t.topup_time <= $2::timestamp)
            OR (t.topup_time >= $3::timestamp AND t.topup_time <= $4::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.topup_time),
        EXTRACT(MONTH FROM t.topup_time)
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



-- name: GetYearlyTopupStatusSuccess :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.topup_time)::integer AS year,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.topup_amount), 0)::integer AS total_amount
    FROM
        topups t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND (
            EXTRACT(YEAR FROM t.topup_time) = $1::integer
            OR EXTRACT(YEAR FROM t.topup_time) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.topup_time)
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



-- name: GetMonthTopupStatusFailed :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.topup_time)::integer AS year,
        EXTRACT(MONTH FROM t.topup_time)::integer AS month,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.topup_amount), 0)::integer AS total_amount
    FROM
        topups t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND (
            (t.topup_time >= $1::timestamp AND t.topup_time <= $2::timestamp)
            OR (t.topup_time >= $3::timestamp AND t.topup_time <= $4::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.topup_time),
        EXTRACT(MONTH FROM t.topup_time)
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


-- name: GetYearlyTopupStatusFailed :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.topup_time)::integer AS year,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.topup_amount), 0)::integer AS total_amount
    FROM
        topups t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND (
            EXTRACT(YEAR FROM t.topup_time) = $1::integer
            OR EXTRACT(YEAR FROM t.topup_time) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.topup_time)
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





-- name: GetMonthTopupStatusSuccessCardNumber :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.topup_time)::integer AS year,
        EXTRACT(MONTH FROM t.topup_time)::integer AS month,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.topup_amount), 0)::integer AS total_amount
    FROM
        topups t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND t.card_number = $1
        AND (
            (t.topup_time >= $2::timestamp AND t.topup_time <= $3::timestamp)
            OR (t.topup_time >= $4::timestamp AND t.topup_time <= $5::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.topup_time),
        EXTRACT(MONTH FROM t.topup_time)
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



-- name: GetYearlyTopupStatusSuccessCardNumber :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.topup_time)::integer AS year,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.topup_amount), 0)::integer AS total_amount
    FROM
        topups t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'success'
        AND t.card_number = $1
            AND (
                EXTRACT(YEAR FROM t.topup_time) = $2::integer
                OR EXTRACT(YEAR FROM t.topup_time) = $2::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.topup_time)
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


-- name: GetMonthTopupStatusFailedCardNumber :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.topup_time)::integer AS year,
        EXTRACT(MONTH FROM t.topup_time)::integer AS month,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.topup_amount), 0)::integer AS total_amount
    FROM
        topups t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND t.card_number = $1
        AND (
            (t.topup_time >= $2::timestamp AND t.topup_time <= $3::timestamp)
            OR (t.topup_time >= $4::timestamp AND t.topup_time <= $5::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.topup_time),
        EXTRACT(MONTH FROM t.topup_time)
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
        EXTRACT(YEAR FROM $4::timestamp)::text AS year,
        TO_CHAR($4::timestamp, 'Mon') AS month,
        0 AS total_failed,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $4::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $4::timestamp)::integer
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC,
    TO_DATE(month, 'Mon') DESC;


-- name: GetYearlyTopupStatusFailedCardNumber :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.topup_time)::integer AS year,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.topup_amount), 0)::integer AS total_amount
    FROM
        topups t
    WHERE
        t.deleted_at IS NULL
        AND t.status = 'failed'
        AND t.card_number = $1
            AND (
                EXTRACT(YEAR FROM t.topup_time) = $2::integer
                OR EXTRACT(YEAR FROM t.topup_time) = $2::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.topup_time)
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


-- name: GetMonthlyTopupMethods :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $1::timestamp),
        date_trunc('year', $1::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
),
topup_methods AS (
    SELECT DISTINCT topup_method
    FROM topups
    WHERE deleted_at IS NULL
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    tm.topup_method,
    COALESCE(COUNT(t.topup_id), 0)::int AS total_topups,
    COALESCE(SUM(t.topup_amount), 0)::int AS total_amount
FROM
    months m
CROSS JOIN
    topup_methods tm
LEFT JOIN
    topups t ON EXTRACT(MONTH FROM t.topup_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM t.created_at) = EXTRACT(YEAR FROM m.month)
    AND t.topup_method = tm.topup_method
    AND t.deleted_at IS NULL
GROUP BY
    m.month,
    tm.topup_method
ORDER BY
    m.month,
    tm.topup_method;





-- name: GetYearlyTopupMethods :many
SELECT
    EXTRACT(YEAR FROM t.created_at) AS year,
    t.topup_method,
    COUNT(t.topup_id) AS total_topups,
    SUM(t.topup_amount) AS total_amount
FROM
    topups t
WHERE
    t.deleted_at IS NULL
    AND EXTRACT(YEAR FROM t.topup_time) >= $1 - 4
    AND EXTRACT(YEAR FROM t.topup_time) <= $1
GROUP BY
    EXTRACT(YEAR FROM t.created_at),
    t.topup_method
ORDER BY
    year;



-- name: GetMonthlyTopupAmounts :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $1::timestamp),
        date_trunc('year', $1::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    COALESCE(SUM(t.topup_amount), 0)::int AS total_amount
FROM
    months m
LEFT JOIN
    topups t ON EXTRACT(MONTH FROM t.topup_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM t.topup_time) = EXTRACT(YEAR FROM m.month)
    AND t.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY
    m.month;



-- name: GetYearlyTopupAmounts :many
SELECT
    EXTRACT(YEAR FROM t.topup_time) AS year,
    SUM(t.topup_amount) AS total_amount
FROM
    topups t
WHERE
    t.deleted_at IS NULL
    AND EXTRACT(YEAR FROM t.topup_time) >= $1 - 4
    AND EXTRACT(YEAR FROM t.topup_time) <= $1
GROUP BY
    EXTRACT(YEAR FROM t.topup_time)
ORDER BY
    year;




-- name: GetMonthlyTopupMethodsByCardNumber :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $2::timestamp),
        date_trunc('year', $2::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
),
topup_methods AS (
    SELECT DISTINCT topup_method
    FROM topups
    WHERE deleted_at IS NULL
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    tm.topup_method,
    COALESCE(COUNT(t.topup_id), 0)::int AS total_topups,
    COALESCE(SUM(t.topup_amount), 0)::int AS total_amount
FROM
    months m
CROSS JOIN
    topup_methods tm
LEFT JOIN
    topups t ON EXTRACT(MONTH FROM t.topup_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM t.topup_time) = EXTRACT(YEAR FROM m.month)
    AND t.topup_method = tm.topup_method
    AND t.card_number = $1
    AND t.deleted_at IS NULL
GROUP BY
    m.month,
    tm.topup_method
ORDER BY
    m.month,
    tm.topup_method;





-- name: GetYearlyTopupMethodsByCardNumber :many
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
    AND EXTRACT(YEAR FROM t.created_at) >= $2 - 4
    AND EXTRACT(YEAR FROM t.created_at) <= $2
GROUP BY
    EXTRACT(YEAR FROM t.topup_time),
    t.topup_method
ORDER BY
    year;



-- name: GetMonthlyTopupAmountsByCardNumber :many
WITH months AS (
    SELECT generate_series(
        date_trunc('year', $2::timestamp),
        date_trunc('year', $2::timestamp) + interval '1 year' - interval '1 day',
        interval '1 month'
    ) AS month
)
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    COALESCE(SUM(t.topup_amount), 0)::int AS total_amount
FROM
    months m
LEFT JOIN
    topups t ON EXTRACT(MONTH FROM t.topup_time) = EXTRACT(MONTH FROM m.month)
    AND EXTRACT(YEAR FROM t.topup_time) = EXTRACT(YEAR FROM m.month)
    AND t.card_number = $1
    AND t.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY
    m.month;



-- name: GetYearlyTopupAmountsByCardNumber :many
SELECT
    EXTRACT(YEAR FROM t.topup_time) AS year,
    SUM(t.topup_amount) AS total_amount
FROM
    topups t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
    AND EXTRACT(YEAR FROM t.created_at) >= $2 - 4
    AND EXTRACT(YEAR FROM t.created_at) <= $2
GROUP BY
    EXTRACT(YEAR FROM t.topup_time)
ORDER BY
    year;



-- Create Topup
-- name: CreateTopup :one
INSERT INTO
    topups (
        card_number,
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
        current_timestamp
    ) RETURNING *;


-- Update Topup
-- name: UpdateTopup :exec
UPDATE topups
SET
    card_number = $2,
    topup_amount = $3,
    topup_method = $4,
    topup_time = $5,
    updated_at = current_timestamp
WHERE
    topup_id = $1
    AND deleted_at IS NULL;


-- Update Topup Amount
-- name: UpdateTopupAmount :exec
UPDATE topups
SET
    topup_amount = $2,
    updated_at = current_timestamp
WHERE
    topup_id = $1
    AND deleted_at IS NULL;


-- Update Topup Status
-- name: UpdateTopupStatus :exec
UPDATE topups
SET
    status = $2,
    updated_at = current_timestamp
WHERE
    topup_id = $1
    AND deleted_at IS NULL;


-- Get Topups by Card Number
-- name: GetTopupsByCardNumber :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM
    topups
WHERE
    deleted_at IS NULL
    AND card_number = $1 -- Filter by card_number
    AND (
        $2::TEXT IS NULL
        OR topup_no::TEXT ILIKE '%' || $2 || '%'
        OR topup_method ILIKE '%' || $2 || '%'
        OR status ILIKE '%' || $2 || '%'
    )
ORDER BY
    topup_time DESC
LIMIT $3 OFFSET $4;

-- Get Trashed By Topup ID
-- name: GetTrashedTopupByID :one
SELECT *
FROM topups
WHERE
    topup_id = $1
    AND deleted_at IS NOT NULL;


-- Trash Topup
-- name: TrashTopup :exec
UPDATE topups
SET
    deleted_at = current_timestamp
WHERE
    topup_id = $1
    AND deleted_at IS NULL;

-- Restore Trashed Topup
-- name: RestoreTopup :exec
UPDATE topups
SET
    deleted_at = NULL
WHERE
    topup_id = $1
    AND deleted_at IS NOT NULL;


-- Delete Topup Permanently
-- name: DeleteTopupPermanently :exec
DELETE FROM topups WHERE topup_id = $1;




-- Restore All Trashed Saldos
-- name: RestoreAllTopups :exec
UPDATE topups
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;


-- Delete All Trashed Saldos Permanently
-- name: DeleteAllPermanentTopups :exec
DELETE FROM topups
WHERE
    deleted_at IS NOT NULL;
