-- name: GetMerchants :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM merchants
WHERE deleted_at IS NULL
    AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR api_key ILIKE '%' || $1 || '%' OR status ILIKE '%' || $1 || '%')
ORDER BY merchant_id
LIMIT $2 OFFSET $3;


-- Get Merchant by ID
-- name: GetMerchantByID :one
SELECT *
FROM merchants
WHERE
    merchant_id = $1
    AND deleted_at IS NULL;

-- Get Merchant by API Key
-- name: GetMerchantByApiKey :one
SELECT * FROM merchants WHERE api_key = $1 AND deleted_at IS NULL;

-- Get Merchant by Name
-- name: GetMerchantByName :one
SELECT * FROM merchants WHERE name = $1 AND deleted_at IS NULL;

-- Get Merchants by User ID
-- name: GetMerchantsByUserID :many
SELECT * FROM merchants WHERE user_id = $1 AND deleted_at IS NULL;


-- name: GetActiveMerchants :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM merchants
WHERE deleted_at IS NULL
    AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR api_key ILIKE '%' || $1 || '%' OR status ILIKE '%' || $1 || '%')
ORDER BY merchant_id
LIMIT $2 OFFSET $3;

-- name: GetTrashedMerchants :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM merchants
WHERE deleted_at IS NOT NULL
    AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR api_key ILIKE '%' || $1 || '%' OR status ILIKE '%' || $1 || '%')
ORDER BY merchant_id
LIMIT $2 OFFSET $3;


-- Get Trashed By Merchant ID
-- name: GetTrashedMerchantByID :one
SELECT *
FROM merchants
WHERE
    merchant_id = $1
    AND deleted_at IS NOT NULL;



-- name: GetMonthlyPaymentMethodsMerchant :many
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
LEFT JOIN
    merchants mch ON t.merchant_id = mch.merchant_id
    AND mch.deleted_at IS NULL
GROUP BY
    m.month,
    pm.payment_method
ORDER BY
    m.month,
    pm.payment_method;



-- name: GetYearlyPaymentMethodMerchant :many
WITH last_five_years AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time) AS year,
        t.payment_method,
        SUM(t.amount) AS total_amount
    FROM
        transactions t
    JOIN
        merchants m ON t.merchant_id = m.merchant_id
    WHERE
        t.deleted_at IS NULL AND m.deleted_at IS NULL
        AND EXTRACT(YEAR FROM t.transaction_time) >= $1 - 4
        AND EXTRACT(YEAR FROM t.transaction_time) <= $1
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time),
        t.payment_method
)
SELECT
    year,
    payment_method,
    total_amount
FROM
    last_five_years
ORDER BY
    year;


-- name: GetMonthlyAmountMerchant :many
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
LEFT JOIN
    merchants mch ON t.merchant_id = mch.merchant_id
    AND mch.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY
    m.month;


-- name: GetYearlyAmountMerchant :many
WITH last_five_years AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time) AS year,
        SUM(t.amount) AS total_amount
    FROM
        transactions t
    JOIN
        merchants m ON t.merchant_id = m.merchant_id
    WHERE
        t.deleted_at IS NULL AND m.deleted_at IS NULL
        AND EXTRACT(YEAR FROM t.transaction_time) >= $1 - 4
        AND EXTRACT(YEAR FROM t.transaction_time) <= $1
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time)
)
SELECT
    year,
    total_amount
FROM
    last_five_years
ORDER BY
    year;



-- name: GetMonthlyTotalAmountMerchant :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::text AS year,
        TO_CHAR(t.transaction_time, 'Mon') AS month,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    INNER JOIN
        merchants m ON t.merchant_id = m.merchant_id
    WHERE
        t.deleted_at IS NULL
        AND m.deleted_at IS NULL
        AND (
            t.transaction_time >= date_trunc('month', $1::timestamp) - interval '1 month'
            AND t.transaction_time < date_trunc('month', $1::timestamp) + interval '1 month'
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time),
        TO_CHAR(t.transaction_time, 'Mon')
), missing_months AS (
    SELECT
        EXTRACT(YEAR FROM $1::timestamp)::text AS year,
        TO_CHAR($1::timestamp, 'Mon') AS month,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $1::timestamp)::text
        AND month = TO_CHAR($1::timestamp, 'Mon')
    )
    UNION ALL
    SELECT
        EXTRACT(YEAR FROM date_trunc('month', $1::timestamp) - interval '1 month')::text AS year,
        TO_CHAR(date_trunc('month', $1::timestamp) - interval '1 month', 'Mon') AS month,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM date_trunc('month', $1::timestamp) - interval '1 month')::text
        AND month = TO_CHAR(date_trunc('month', $1::timestamp) - interval '1 month', 'Mon')
    )
)
SELECT year, month, total_amount
FROM (
    SELECT year, month, total_amount FROM monthly_data
    UNION ALL
    SELECT year, month, total_amount FROM missing_months
) combined
ORDER BY
    year DESC,
    TO_DATE(month, 'Mon') DESC;


-- name: GetYearlyTotalAmountMerchant :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    INNER JOIN
        merchants m ON t.merchant_id = m.merchant_id
    WHERE
        t.deleted_at IS NULL
        AND m.deleted_at IS NULL
        AND (
            EXTRACT(YEAR FROM t.transaction_time) = $1::integer
            OR EXTRACT(YEAR FROM t.transaction_time) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        total_amount::integer
    FROM
        yearly_data

    UNION ALL

    SELECT
        $1::text AS year,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer
    )

    UNION ALL

    SELECT
        ($1::integer - 1)::text AS year,
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



-- name: FindAllTransactions :many
SELECT
    t.transaction_id,
    t.card_number,
    t.amount,
    t.payment_method,
    t.merchant_id,
    m.name AS merchant_name,
    t.transaction_time,
    t.created_at,
    t.updated_at,
    t.deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    transactions t
JOIN
    merchants m ON t.merchant_id = m.merchant_id
WHERE
    t.deleted_at IS NULL
    AND ($1::TEXT IS NULL OR t.card_number ILIKE '%' || $1 || '%' OR t.payment_method ILIKE '%' || $1 || '%')
ORDER BY
    t.transaction_time DESC
LIMIT $2 OFFSET $3;



-- name: GetMonthlyPaymentMethodByMerchants :many
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
    AND t.merchant_id = $2
LEFT JOIN
    merchants mch ON t.merchant_id = mch.merchant_id
    AND mch.deleted_at IS NULL
GROUP BY
    m.month,
    pm.payment_method
ORDER BY
    m.month,
    pm.payment_method;


-- name: GetYearlyPaymentMethodByMerchants :many
WITH last_five_years AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time) AS year,
        t.payment_method,
        SUM(t.amount) AS total_amount
    FROM
        transactions t
    JOIN
        merchants m ON t.merchant_id = m.merchant_id
    WHERE
        t.deleted_at IS NULL
        AND m.deleted_at IS NULL
        AND t.merchant_id = $1
        AND EXTRACT(YEAR FROM t.transaction_time) >= $2 - 4
        AND EXTRACT(YEAR FROM t.transaction_time) <= $2
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time),
        t.payment_method
)
SELECT
    year,
    payment_method,
    total_amount
FROM
    last_five_years
ORDER BY
    year;


-- name: GetMonthlyAmountByMerchants :many
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
    AND t.merchant_id = $2
LEFT JOIN
    merchants mch ON t.merchant_id = mch.merchant_id
    AND mch.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY
    m.month;



-- name: GetYearlyAmountByMerchants :many
WITH last_five_years AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time) AS year,
        SUM(t.amount) AS total_amount
    FROM
        transactions t
    JOIN
        merchants m ON t.merchant_id = m.merchant_id
    WHERE
        t.deleted_at IS NULL
        AND m.deleted_at IS NULL
        AND t.merchant_id = $1
        AND EXTRACT(YEAR FROM t.transaction_time) >= $2 - 4
        AND EXTRACT(YEAR FROM t.transaction_time) <= $2
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time)
)
SELECT
    year,
    total_amount
FROM
    last_five_years
ORDER BY
    year;


-- name: GetMonthlyTotalAmountByMerchant :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        EXTRACT(MONTH FROM t.transaction_time)::integer AS month,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    INNER JOIN
        merchants m ON t.merchant_id = m.merchant_id
    WHERE
        t.deleted_at IS NULL
        AND m.deleted_at IS NULL
        AND EXTRACT(YEAR FROM t.transaction_time) = EXTRACT(YEAR FROM $1::timestamp)
        AND t.merchant_id = $2::integer
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time),
        EXTRACT(MONTH FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
        total_amount
    FROM
        monthly_data
    UNION ALL

    SELECT
        EXTRACT(YEAR FROM gs.month)::text AS year,
        TO_CHAR(gs.month, 'Mon') AS month,
        0::integer AS total_amount
    FROM generate_series(
        date_trunc('year', $1::timestamp),
        date_trunc('year', $1::timestamp) + interval '11 month',
        interval '1 month'
    ) AS gs(month)
    WHERE NOT EXISTS (
        SELECT 1 FROM monthly_data md
        WHERE md.year = EXTRACT(YEAR FROM gs.month)::integer
        AND md.month = EXTRACT(MONTH FROM gs.month)::integer
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC,
    TO_DATE(month, 'Mon') DESC;


-- name: GetYearlyTotalAmountByMerchant :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    INNER JOIN
        merchants m ON t.merchant_id = m.merchant_id
    WHERE
        t.deleted_at IS NULL
        AND m.deleted_at IS NULL
        AND EXTRACT(YEAR FROM t.transaction_time) >= $1::integer - 4
        AND EXTRACT(YEAR FROM t.transaction_time) <= $1::integer
        AND t.merchant_id = $2::integer
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        total_amount
    FROM
        yearly_data
    UNION ALL

    SELECT
        y::text AS year,
        0::integer AS total_amount
    FROM generate_series($1::integer - 4, $1::integer) AS y
    WHERE NOT EXISTS (
        SELECT 1 FROM yearly_data yd
        WHERE yd.year = y
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC;


-- name: FindAllTransactionsByMerchant :many
SELECT
    t.transaction_id,
    t.card_number,
    t.amount,
    t.payment_method,
    t.merchant_id,
    m.name AS merchant_name,
    t.transaction_time,
    t.created_at,
    t.updated_at,
    t.deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    transactions t
JOIN
    merchants m ON t.merchant_id = m.merchant_id
WHERE
    t.deleted_at IS NULL
    AND t.merchant_id = $1
    AND ($2::TEXT IS NULL OR t.card_number ILIKE '%' || $2 || '%' OR t.payment_method ILIKE '%' || $2 || '%')
ORDER BY
    t.transaction_time DESC
LIMIT $3 OFFSET $4;


-- name: GetMonthlyPaymentMethodByApikey :many
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
LEFT JOIN
    merchants mch ON t.merchant_id = mch.merchant_id
    AND mch.deleted_at IS NULL
    AND mch.api_key = $2
GROUP BY
    m.month,
    pm.payment_method
ORDER BY
    m.month,
    pm.payment_method;


-- name: GetYearlyPaymentMethodByApikey :many
WITH last_five_years AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time) AS year,
        t.payment_method,
        SUM(t.amount) AS total_amount
    FROM
        transactions t
    JOIN
        merchants m ON t.merchant_id = m.merchant_id
    WHERE
        t.deleted_at IS NULL
        AND m.deleted_at IS NULL
        AND m.api_key = $1
        AND EXTRACT(YEAR FROM t.transaction_time) >= $2 - 4
        AND EXTRACT(YEAR FROM t.transaction_time) <= $2
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time),
        t.payment_method
)
SELECT
    year,
    payment_method,
    total_amount
FROM
    last_five_years
ORDER BY
    year;


-- name: GetMonthlyAmountByApikey :many
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
LEFT JOIN
    merchants mch ON t.merchant_id = mch.merchant_id
    AND mch.deleted_at IS NULL
    AND mch.api_key = $2
GROUP BY
    m.month
ORDER BY
    m.month;


-- name: GetYearlyAmountByApikey :many
WITH last_five_years AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time) AS year,
        SUM(t.amount) AS total_amount
    FROM
        transactions t
    JOIN
        merchants m ON t.merchant_id = m.merchant_id
    WHERE
        t.deleted_at IS NULL
        AND m.deleted_at IS NULL
        AND m.api_key = $1
        AND EXTRACT(YEAR FROM t.transaction_time) >= $2 - 4
        AND EXTRACT(YEAR FROM t.transaction_time) <= $2
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time)
)
SELECT
    year,
    total_amount
FROM
    last_five_years
ORDER BY
    year;



-- name: GetMonthlyTotalAmountByApikey :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        EXTRACT(MONTH FROM t.transaction_time)::integer AS month,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    INNER JOIN
        merchants m ON t.merchant_id = m.merchant_id
    WHERE
        t.deleted_at IS NULL
        AND m.deleted_at IS NULL
        AND EXTRACT(YEAR FROM t.transaction_time) = EXTRACT(YEAR FROM $1::timestamp)
        AND m.api_key = $2
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time),
        EXTRACT(MONTH FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
        total_amount
    FROM
        monthly_data
    UNION ALL

    SELECT
        EXTRACT(YEAR FROM gs.month)::text AS year,
        TO_CHAR(gs.month, 'Mon') AS month,
        0::integer AS total_amount
    FROM generate_series(
        date_trunc('year', $1::timestamp),
        date_trunc('year', $1::timestamp) + interval '11 month',
        interval '1 month'
    ) AS gs(month)
    WHERE NOT EXISTS (
        SELECT 1 FROM monthly_data md
        WHERE md.year = EXTRACT(YEAR FROM gs.month)::integer
        AND md.month = EXTRACT(MONTH FROM gs.month)::integer
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC,
    TO_DATE(month, 'Mon') DESC;


-- name: GetYearlyTotalAmountByApikey :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.transaction_time)::integer AS year,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    INNER JOIN
        merchants m ON t.merchant_id = m.merchant_id
    WHERE
        t.deleted_at IS NULL
        AND m.deleted_at IS NULL
        AND EXTRACT(YEAR FROM t.transaction_time) >= $1::integer - 4
        AND EXTRACT(YEAR FROM t.transaction_time) <= $1::integer
        AND m.api_key = $2
    GROUP BY
        EXTRACT(YEAR FROM t.transaction_time)
), formatted_data AS (
    SELECT
        year::text,
        total_amount
    FROM
        yearly_data
    UNION ALL

    SELECT
        y::text AS year,
        0::integer AS total_amount
    FROM generate_series($1::integer - 4, $1::integer) AS y
    WHERE NOT EXISTS (
        SELECT 1 FROM yearly_data yd
        WHERE yd.year = y
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC;



-- name: FindAllTransactionsByApikey :many
SELECT
    t.transaction_id,
    t.card_number,
    t.amount,
    t.payment_method,
    t.merchant_id,
    m.name AS merchant_name,
    t.transaction_time,
    t.created_at,
    t.updated_at,
    t.deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    transactions t
JOIN
    merchants m ON t.merchant_id = m.merchant_id
WHERE
    t.deleted_at IS NULL
    AND m.api_key = $1
    AND ($2::TEXT IS NULL OR t.card_number ILIKE '%' || $2 || '%' OR t.payment_method ILIKE '%' || $2 || '%')
ORDER BY
    t.transaction_time DESC
LIMIT $3 OFFSET $4;



-- Create Merchant
-- name: CreateMerchant :one
INSERT INTO
    merchants (
        name,
        api_key,
        user_id,
        status,
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



-- Update Merchant
-- name: UpdateMerchant :exec
UPDATE merchants
SET
    name = $2,
    user_id = $3,
    status = $4,
    updated_at = current_timestamp
WHERE
    merchant_id = $1
    AND deleted_at IS NULL;

-- UpdateMerchantStatus
-- name: UpdateMerchantStatus :exec
UPDATE merchants
SET
    status = $2,
    updated_at = current_timestamp
WHERE
    merchant_id = $1
    AND deleted_at IS NULL;



-- Trash Merchant
-- name: TrashMerchant :exec
UPDATE merchants
SET
    deleted_at = current_timestamp
WHERE
    merchant_id = $1
    AND deleted_at IS NULL;

-- Restore Trashed Merchant
-- name: RestoreMerchant :exec
UPDATE merchants
SET
    deleted_at = NULL
WHERE
    merchant_id = $1
    AND deleted_at IS NOT NULL;


-- Delete Merchant Permanently
-- name: DeleteMerchantPermanently :exec
DELETE FROM merchants WHERE merchant_id = $1 AND deleted_at IS NOT NULL;


-- Restore All Trashed Merchants
-- name: RestoreAllMerchants :exec
UPDATE merchants
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- Delete All Trashed Merchants Permanently
-- name: DeleteAllPermanentMerchants :exec
DELETE FROM merchants
WHERE
    deleted_at IS NOT NULL;
