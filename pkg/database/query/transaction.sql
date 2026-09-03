-- GetTransactions: Retrieves paginated transaction records with search capability
-- Purpose: List all transactions for management UI with filtering options
-- Parameters:
--   $1: search_term - Optional text to filter transactions by card number, payment method, or status (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All transaction fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted transactions (deleted_at IS NULL)
--   - Supports partial text matching on multiple fields (case-insensitive)
--   - Orders by transaction_time (newest first)
--   - Provides total_count for pagination calculations
--   - Useful for transaction monitoring and auditing
-- name: GetTransactions :many
SELECT
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM transactions
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR card_number ILIKE '%' || $1 || '%'
        OR payment_method ILIKE '%' || $1 || '%'
        OR status ILIKE '%' || $1 || '%'
    )
ORDER BY transaction_time DESC
LIMIT $2
OFFSET
    $3;

-- GetActiveTransactions: Retrieves paginated active transactions with search
-- Purpose: List all non-deleted transactions with filtering options
-- Parameters:
--   $1: search_term - Optional text to filter by card number or payment method
--   $2: limit - Maximum records to return
--   $3: offset - Records to skip for pagination
-- Returns:
--   All transaction fields plus total_count of matching active records
-- Business Logic:
--   - Only includes active transactions (deleted_at IS NULL)
--   - Filters on card_number and payment_method fields
--   - Orders by transaction_time (newest first)
--   - Provides pagination metadata
--   - Used in transaction management interfaces
-- name: GetActiveTransactions :many
SELECT
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM transactions
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR card_number ILIKE '%' || $1 || '%'
        OR payment_method ILIKE '%' || $1 || '%'
    )
ORDER BY transaction_time DESC
LIMIT $2
OFFSET
    $3;

-- GetTrashedTransactions: Retrieves paginated soft-deleted transactions
-- Purpose: List all deleted transactions for recovery or audit purposes
-- Parameters:
--   $1: search_term - Optional text to filter deleted transactions
--   $2: limit - Maximum records to return
--   $3: offset - Records to skip for pagination
-- Returns:
--   All transaction fields plus total_count of matching deleted records
-- Business Logic:
--   - Only includes soft-deleted transactions (deleted_at IS NOT NULL)
--   - Same filtering capabilities as active transactions
--   - Maintains newest-first ordering
--   - Used in admin interfaces for transaction recovery
-- name: GetTrashedTransactions :many
SELECT
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM transactions
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR card_number ILIKE '%' || $1 || '%'
        OR payment_method ILIKE '%' || $1 || '%'
    )
ORDER BY transaction_time DESC
LIMIT $2
OFFSET
    $3;

-- GetTransactionByID: Retrieves a single transaction by its ID
-- Purpose: Get detailed information about a specific transaction
-- Parameters:
--   $1: transaction_id - The ID of the transaction to retrieve
-- Returns:
--   All fields for the specified transaction or NULL if not found/deleted
-- Business Logic:
--   - Only returns active transactions (deleted_at IS NULL)
--   - Useful for transaction details viewing and verification
-- name: GetTransactionByID :one
SELECT
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    created_at,
    updated_at
FROM transactions
WHERE
    transaction_id = $1
    AND deleted_at IS NULL;

-- GetTransactionsByCardNumber: Retrieves paginated transactions for a specific card
-- Purpose: List all transactions associated with a particular card
-- Parameters:
--   $1: card_number - The card number to filter transactions
--   $2: search_term - Optional text to filter by payment method or status
--   $3: limit - Maximum number of records to return
--   $4: offset - Number of records to skip for pagination
-- Returns:
--   All transaction fields plus total_count of matching records
-- Business Logic:
--   - Only includes active transactions (deleted_at IS NULL)
--   - Strict card number matching combined with optional search filters
--   - Orders by transaction_time (newest first)
--   - Provides pagination support with total_count
--   - Useful for cardholder transaction history
-- name: GetTransactionsByCardNumber :many
SELECT
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM transactions
WHERE
    deleted_at IS NULL
    AND card_number = $1
    AND (
        $2::TEXT IS NULL
        OR payment_method ILIKE '%' || $2 || '%'
        OR status ILIKE '%' || $2 || '%'
    )
ORDER BY transaction_time DESC
LIMIT $3
OFFSET
    $4;

-- GetTransactionsByMerchantID: Retrieves transactions for a specific merchant
-- Purpose: List all transactions associated with a merchant
-- Parameters:
--   $1: merchant_id - The ID of the merchant to filter transactions
-- Returns:
--   All transaction fields for the merchant's transactions
-- Business Logic:
--   - Only includes active transactions (deleted_at IS NULL)
--   - Orders by transaction_time (newest first)
--   - No pagination (assumes manageable number of records per merchant)
--   - Useful for merchant transaction reports
-- name: GetTransactionsByMerchantID :many
SELECT
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    created_at,
    updated_at
FROM transactions
WHERE
    merchant_id = $1
    AND deleted_at IS NULL
ORDER BY transaction_time DESC;

-- GetTrashedTransactionByID: Retrieves a single soft-deleted transaction by ID
-- Purpose: View details of a deleted transaction for recovery or audit
-- Parameters:
--   $1: transaction_id - The ID of the transaction to retrieve
-- Returns:
--   All fields for the specified trashed transaction or NULL if not found/active
-- Business Logic:
--   - Only returns soft-deleted transactions (deleted_at IS NOT NULL)
--   - Used in admin interfaces for transaction recovery
-- name: GetTrashedTransactionByID :one
SELECT *
FROM transactions
WHERE
    transaction_id = $1
    AND deleted_at IS NOT NULL;

-- GetMonthTransactionStatusSuccess: Retrieves monthly success metrics for transactions
-- Purpose: Analyze successful transaction trends across comparison periods
-- Parameters:
--   $1: period1_start - Start date of first comparison period
--   $2: period1_end - End date of first comparison period
--   $3: period2_start - Start date of second comparison period
--   $4: period2_end - End date of second comparison period
-- Returns:
--   year: Year as text (e.g., '2023')
--   month: 3-letter month abbreviation (e.g., 'Jan')
--   total_success: Count of successful transactions
--   total_amount: Sum of successful transaction amounts
-- Business Logic:
--   - Only includes successful transactions (status = 'success')
--   - Covers two customizable time periods for comparison
--   - Zero-fills months with no activity
--   - Formats output for consistent visualization (year as text, month as 'Mon')
--   - Orders by year and month (newest first)
--   - Useful for identifying seasonal transaction patterns and revenue trends
-- name: GetMonthTransactionStatusSuccess :many
WITH
    monthly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transaction_time
            )::integer AS year,
            EXTRACT(
                MONTH
                FROM t.transaction_time
            )::integer AS month,
            COUNT(*) AS total_success,
            COALESCE(SUM(t.amount), 0)::integer AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'success'
            AND (
                (
                    t.transaction_time >= $1::timestamp
                    AND t.transaction_time <= $2::timestamp
                )
                OR (
                    t.transaction_time >= $3::timestamp
                    AND t.transaction_time <= $4::timestamp
                )
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transaction_time
            ),
            EXTRACT(
                MONTH
                FROM t.transaction_time
            )
    ),
    formatted_data AS (
        SELECT
            year::text,
            TO_CHAR(
                TO_DATE(month::text, 'MM'),
                'Mon'
            ) AS month,
            total_success,
            total_amount
        FROM monthly_data
        UNION ALL
        SELECT
            EXTRACT(
                YEAR
                FROM $1::timestamp
            )::text AS year,
            TO_CHAR($1::timestamp, 'Mon') AS month,
            0 AS total_success,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $1::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $1::timestamp
                    )::integer
            )
        UNION ALL
        SELECT
            EXTRACT(
                YEAR
                FROM $3::timestamp
            )::text AS year,
            TO_CHAR($3::timestamp, 'Mon') AS month,
            0 AS total_success,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $3::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $3::timestamp
                    )::integer
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC, TO_DATE(month, 'Mon') DESC;

-- GetYearlyTransactionStatusSuccess: Retrieves yearly success metrics for transactions
-- Purpose: Compare annual successful transaction performance
-- Parameters:
--   $1: current_year - The target year (includes this year and previous)
-- Returns:
--   year: Year as text (e.g., '2023')
--   total_success: Count of successful transactions
--   total_amount: Sum of successful transaction amounts
-- Business Logic:
--   - Only includes successful transactions (status = 'success')
--   - Compares current year with previous year
--   - Zero-fills years with no activity
--   - Orders by year (newest first)
--   - Useful for year-over-year growth analysis and financial reporting
--   - Helps identify annual transaction volume and revenue trends
-- name: GetYearlyTransactionStatusSuccess :many
WITH
    yearly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transaction_time
            )::integer AS year,
            COUNT(*) AS total_success,
            COALESCE(SUM(t.amount), 0)::integer AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'success'
            AND (
                EXTRACT(
                    YEAR
                    FROM t.transaction_time
                ) = $1::integer
                OR EXTRACT(
                    YEAR
                    FROM t.transaction_time
                ) = $1::integer - 1
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transaction_time
            )
    ),
    formatted_data AS (
        SELECT
            year::text,
            total_success::integer,
            total_amount::integer
        FROM yearly_data
        UNION ALL
        SELECT
            $1::text AS year,
            0::integer AS total_success,
            0::integer AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM yearly_data
                WHERE
                    year = $1::integer
            )
        UNION ALL
        SELECT ($1::integer - 1)::text AS year,
            0::integer AS total_success,
            0::integer AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM yearly_data
                WHERE
                    year = $1::integer - 1
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC;

-- GetMonthTransactionStatusFailed: Retrieves monthly failed metrics for transactions
-- Purpose: Analyze failedful transaction trends across comparison periods
-- Parameters:
--   $1: period1_start - Start date of first comparison period
--   $2: period1_end - End date of first comparison period
--   $3: period2_start - Start date of second comparison period
--   $4: period2_end - End date of second comparison period
-- Returns:
--   year: Year as text (e.g., '2023')
--   month: 3-letter month abbreviation (e.g., 'Jan')
--   total_failed: Count of failedful transactions
--   total_amount: Sum of failedful transaction amounts
-- Business Logic:
--   - Only includes failedful transactions (status = 'failed')
--   - Covers two customizable time periods for comparison
--   - Zero-fills months with no activity
--   - Formats output for consistent visualization (year as text, month as 'Mon')
--   - Orders by year and month (newest first)
--   - Useful for identifying seasonal transaction patterns and revenue trends
-- name: GetMonthTransactionStatusFailed :many
WITH
    monthly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transaction_time
            )::integer AS year,
            EXTRACT(
                MONTH
                FROM t.transaction_time
            )::integer AS month,
            COUNT(*) AS total_failed,
            COALESCE(SUM(t.amount), 0)::integer AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'failed'
            AND (
                (
                    t.transaction_time >= $1::timestamp
                    AND t.transaction_time <= $2::timestamp
                )
                OR (
                    t.transaction_time >= $3::timestamp
                    AND t.transaction_time <= $4::timestamp
                )
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transaction_time
            ),
            EXTRACT(
                MONTH
                FROM t.transaction_time
            )
    ),
    formatted_data AS (
        SELECT
            year::text,
            TO_CHAR(
                TO_DATE(month::text, 'MM'),
                'Mon'
            ) AS month,
            total_failed,
            total_amount
        FROM monthly_data
        UNION ALL
        SELECT
            EXTRACT(
                YEAR
                FROM $1::timestamp
            )::text AS year,
            TO_CHAR($1::timestamp, 'Mon') AS month,
            0 AS total_failed,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $1::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $1::timestamp
                    )::integer
            )
        UNION ALL
        SELECT
            EXTRACT(
                YEAR
                FROM $3::timestamp
            )::text AS year,
            TO_CHAR($3::timestamp, 'Mon') AS month,
            0 AS total_failed,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $3::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $3::timestamp
                    )::integer
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC, TO_DATE(month, 'Mon') DESC;

-- GetYearlyTransactionStatusFailed: Retrieves yearly failed metrics for transactions
-- Purpose: Compare annual failedful transaction performance
-- Parameters:
--   $1: current_year - The target year (includes this year and previous)
-- Returns:
--   year: Year as text (e.g., '2023')
--   total_failed: Count of failedful transactions
--   total_amount: Sum of failedful transaction amounts
-- Business Logic:
--   - Only includes failedful transactions (status = 'failed')
--   - Compares current year with previous year
--   - Zero-fills years with no activity
--   - Orders by year (newest first)
--   - Useful for year-over-year growth analysis and financial reporting
--   - Helps identify annual transaction volume and revenue trends
-- name: GetYearlyTransactionStatusFailed :many
WITH
    yearly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transaction_time
            )::integer AS year,
            COUNT(*) AS total_failed,
            COALESCE(SUM(t.amount), 0)::integer AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'failed'
            AND (
                EXTRACT(
                    YEAR
                    FROM t.transaction_time
                ) = $1::integer
                OR EXTRACT(
                    YEAR
                    FROM t.transaction_time
                ) = $1::integer - 1
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transaction_time
            )
    ),
    formatted_data AS (
        SELECT
            year::text,
            total_failed::integer,
            total_amount::integer
        FROM yearly_data
        UNION ALL
        SELECT
            $1::text AS year,
            0::integer AS total_failed,
            0::integer AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM yearly_data
                WHERE
                    year = $1::integer
            )
        UNION ALL
        SELECT ($1::integer - 1)::text AS year,
            0::integer AS total_failed,
            0::integer AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM yearly_data
                WHERE
                    year = $1::integer - 1
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC;

-- GetMonthTransactionStatusSuccessCardNumber: Retrieves monthly success metrics for transactions
-- Purpose: Analyze successful transaction trends across comparison periods
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: period1_start - Start date of first comparison period
--   $3: period1_end - End date of first comparison period
--   $4: period2_start - Start date of second comparison period
--   $5: period2_end - End date of second comparison period
-- Returns:
--   year: Year as text (e.g., '2023')
--   month: 3-letter month abbreviation (e.g., 'Jan')
--   total_success: Count of successful transactions
--   total_amount: Sum of successful transaction amounts
-- Business Logic:
--   - Only includes successful transactions (status = 'success')
--   - Covers two customizable time periods for comparison
--   - Zero-fills months with no activity
--   - Formats output for consistent visualization (year as text, month as 'Mon')
--   - Orders by year and month (newest first)
--   - Useful for identifying seasonal transaction patterns and revenue trends
-- name: GetMonthTransactionStatusSuccessCardNumber :many
WITH
    monthly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transaction_time
            )::integer AS year,
            EXTRACT(
                MONTH
                FROM t.transaction_time
            )::integer AS month,
            COUNT(*) AS total_success,
            COALESCE(SUM(t.amount), 0)::integer AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'success'
            AND t.card_number = $1
            AND (
                (
                    t.transaction_time >= $2::timestamp
                    AND t.transaction_time <= $3::timestamp
                )
                OR (
                    t.transaction_time >= $4::timestamp
                    AND t.transaction_time <= $5::timestamp
                )
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transaction_time
            ),
            EXTRACT(
                MONTH
                FROM t.transaction_time
            )
    ),
    formatted_data AS (
        SELECT
            year::text,
            TO_CHAR(
                TO_DATE(month::text, 'MM'),
                'Mon'
            ) AS month,
            total_success,
            total_amount
        FROM monthly_data
        UNION ALL
        SELECT
            EXTRACT(
                YEAR
                FROM $2::timestamp
            )::text AS year,
            TO_CHAR($2::timestamp, 'Mon') AS month,
            0 AS total_success,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $2::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $2::timestamp
                    )::integer
            )
        UNION ALL
        SELECT
            EXTRACT(
                YEAR
                FROM $3::timestamp
            )::text AS year,
            TO_CHAR($3::timestamp, 'Mon') AS month,
            0 AS total_success,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $3::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $3::timestamp
                    )::integer
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC, TO_DATE(month, 'Mon') DESC;

-- GetYearlyTransactionStatusSuccessCardNumber: Retrieves yearly success metrics for transactions
-- Purpose: Compare annual successful transaction performance
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: current_year - The target year (includes this year and previous)
-- Returns:
--   year: Year as text (e.g., '2023')
--   total_success: Count of successful transactions
--   total_amount: Sum of successful transaction amounts
-- Business Logic:
--   - Only includes successful transactions (status = 'success')
--   - Compares current year with previous year
--   - Zero-fills years with no activity
--   - Orders by year (newest first)
--   - Useful for year-over-year growth analysis and financial reporting
--   - Helps identify annual transaction volume and revenue trends
-- name: GetYearlyTransactionStatusSuccessCardNumber :many
WITH
    yearly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transaction_time
            )::integer AS year,
            COUNT(*) AS total_success,
            COALESCE(SUM(t.amount), 0)::integer AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'success'
            AND t.card_number = $1
            AND (
                EXTRACT(
                    YEAR
                    FROM t.transaction_time
                ) = $2::integer
                OR EXTRACT(
                    YEAR
                    FROM t.transaction_time
                ) = $2::integer - 1
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transaction_time
            )
    ),
    formatted_data AS (
        SELECT
            year::text,
            total_success::integer,
            total_amount::integer
        FROM yearly_data
        UNION ALL
        SELECT
            $2::text AS year,
            0::integer AS total_success,
            0::integer AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM yearly_data
                WHERE
                    year = $2::integer
            )
        UNION ALL
        SELECT ($2::integer - 1)::text AS year,
            0::integer AS total_success,
            0::integer AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM yearly_data
                WHERE
                    year = $2::integer - 1
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC;

-- GetMonthTransactionStatusFailed: Retrieves monthly failed metrics for transactions
-- Purpose: Analyze failedful transaction trends across comparison periods
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: period1_start - Start date of first comparison period
--   $3: period1_end - End date of first comparison period
--   $4: period2_start - Start date of second comparison period
--   $5: period2_end - End date of second comparison period
-- Returns:
--   year: Year as text (e.g., '2023')
--   month: 3-letter month abbreviation (e.g., 'Jan')
--   total_failed: Count of failedful transactions
--   total_amount: Sum of failedful transaction amounts
-- Business Logic:
--   - Only includes failedful transactions (status = 'failed')
--   - Covers two customizable time periods for comparison
--   - Zero-fills months with no activity
--   - Formats output for consistent visualization (year as text, month as 'Mon')
--   - Orders by year and month (newest first)
--   - Useful for identifying seasonal transaction patterns and revenue trends
-- name: GetMonthTransactionStatusFailedCardNumber :many
WITH
    monthly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transaction_time
            )::integer AS year,
            EXTRACT(
                MONTH
                FROM t.transaction_time
            )::integer AS month,
            COUNT(*) AS total_failed,
            COALESCE(SUM(t.amount), 0)::integer AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'failed'
            AND t.card_number = $1
            AND (
                (
                    t.transaction_time >= $2::timestamp
                    AND t.transaction_time <= $3::timestamp
                )
                OR (
                    t.transaction_time >= $4::timestamp
                    AND t.transaction_time <= $5::timestamp
                )
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transaction_time
            ),
            EXTRACT(
                MONTH
                FROM t.transaction_time
            )
    ),
    formatted_data AS (
        SELECT
            year::text,
            TO_CHAR(
                TO_DATE(month::text, 'MM'),
                'Mon'
            ) AS month,
            total_failed,
            total_amount
        FROM monthly_data
        UNION ALL
        SELECT
            EXTRACT(
                YEAR
                FROM $2::timestamp
            )::text AS year,
            TO_CHAR($2::timestamp, 'Mon') AS month,
            0 AS total_failed,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $2::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $2::timestamp
                    )::integer
            )
        UNION ALL
        SELECT
            EXTRACT(
                YEAR
                FROM $3::timestamp
            )::text AS year,
            TO_CHAR($3::timestamp, 'Mon') AS month,
            0 AS total_failed,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $3::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $3::timestamp
                    )::integer
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC, TO_DATE(month, 'Mon') DESC;

-- GetYearlyTransactionStatusFailed: Retrieves yearly failed metrics for transactions
-- Purpose: Compare annual failedful transaction performance
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: current_year - The target year (includes this year and previous)
-- Returns:
--   year: Year as text (e.g., '2023')
--   total_failed: Count of failedful transactions
--   total_amount: Sum of failedful transaction amounts
-- Business Logic:
--   - Only includes failedful transactions (status = 'failed')
--   - Compares current year with previous year
--   - Zero-fills years with no activity
--   - Orders by year (newest first)
--   - Useful for year-over-year growth analysis and financial reporting
--   - Helps identify annual transaction volume and revenue trends
-- name: GetYearlyTransactionStatusFailedCardNumber :many
WITH
    yearly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transaction_time
            )::integer AS year,
            COUNT(*) AS total_failed,
            COALESCE(SUM(t.amount), 0)::integer AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'failed'
            AND t.card_number = $1
            AND (
                EXTRACT(
                    YEAR
                    FROM t.transaction_time
                ) = $2::integer
                OR EXTRACT(
                    YEAR
                    FROM t.transaction_time
                ) = $2::integer - 1
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transaction_time
            )
    ),
    formatted_data AS (
        SELECT
            year::text,
            total_failed::integer,
            total_amount::integer
        FROM yearly_data
        UNION ALL
        SELECT
            $2::text AS year,
            0::integer AS total_failed,
            0::integer AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM yearly_data
                WHERE
                    year = $2::integer
            )
        UNION ALL
        SELECT ($2::integer - 1)::text AS year,
            0::integer AS total_failed,
            0::integer AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM yearly_data
                WHERE
                    year = $2::integer - 1
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC;

-- GetMonthlyPaymentMethods: Retrieves a monthly summary of transaction transactions categorized by payment method
-- Purpose:
--   Useful for visualizing how each payment method is used over time within a given year
-- Parameters:
--   $1: reference_date - Any date within the target year (used to generate monthly range)
-- Returns:
--   - month (e.g., 'Jan', 'Feb')
--   - payment_method (e.g., 'e-wallet', 'bank_transfer')
--   - total_transactions: Number of transactions for the method that month
--   - total_amount: Total amount of transactions for the method that month
-- Business Logic:
--   - Includes all combinations of months and available payment methods (even if 0 data)
--   - Excludes soft-deleted transactions (deleted_at IS NULL)
--   - Uses CROSS JOIN to ensure all months and methods are represented
-- name: GetMonthlyPaymentMethods :many
WITH
    months AS (
        SELECT generate_series(
                date_trunc('year', $1::timestamp), date_trunc('year', $1::timestamp) + interval '1 year' - interval '1 day', interval '1 month'
            ) AS month
    ),
    payment_methods AS (
        SELECT DISTINCT
            payment_method
        FROM transactions
        WHERE
            deleted_at IS NULL
    )
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    pm.payment_method,
    COALESCE(COUNT(t.transaction_id), 0)::int AS total_transactions,
    COALESCE(SUM(t.amount), 0)::int AS total_amount
FROM
    months m
    CROSS JOIN payment_methods pm
    LEFT JOIN transactions t ON EXTRACT(
        MONTH
        FROM t.transaction_time
    ) = EXTRACT(
        MONTH
        FROM m.month
    )
    AND EXTRACT(
        YEAR
        FROM t.transaction_time
    ) = EXTRACT(
        YEAR
        FROM m.month
    )
    AND t.payment_method = pm.payment_method
    AND t.deleted_at IS NULL
GROUP BY
    m.month,
    pm.payment_method
ORDER BY m.month, pm.payment_method;

-- GetYearlyPaymentMethods: Retrieves yearly summary of transaction transactions grouped by payment method over a 5-year span
-- Purpose:
--   Analyze long-term trends of transaction method usage across years
-- Parameters:
--   $1: current_year - The most recent year to include (covers current_year - 4 to current_year)
-- Returns:
--   - year: Year of transaction (e.g., 2020, 2021)
--   - payment_method
--   - total_transactions: Count of transactions per method per year
--   - total_amount: Sum of amounts per method per year
-- Business Logic:
--   - Filters data within a 5-year window
--   - Excludes soft-deleted transactions (deleted_at IS NULL)
-- name: GetYearlyPaymentMethods :many
SELECT
    EXTRACT(
        YEAR
        FROM t.created_at
    ) AS year,
    t.payment_method,
    COUNT(t.transaction_id) AS total_transactions,
    SUM(t.amount) AS total_amount
FROM transactions t
WHERE
    t.deleted_at IS NULL
    AND EXTRACT(
        YEAR
        FROM t.created_at
    ) >= $1 - 4
    AND EXTRACT(
        YEAR
        FROM t.created_at
    ) <= $1
GROUP BY
    EXTRACT(
        YEAR
        FROM t.created_at
    ),
    t.payment_method
ORDER BY year;

-- GetMonthlyAmounts: Retrieves total transaction amount per month for a specific year
-- Purpose:
--   Visualize monthly trends in transaction volume for charting/dashboards
-- Parameters:
--   $1: reference_date - Any date within the target year
-- Returns:
--   - month: 3-letter month abbreviation
--   - total_amount: Sum of transaction amounts in each month
-- Business Logic:
--   - Uses LEFT JOIN to ensure all months are included, even with 0 transactions
--   - Filters out soft-deleted data (deleted_at IS NULL)
-- name: GetMonthlyAmounts :many
WITH
    months AS (
        SELECT generate_series(
                date_trunc('year', $1::timestamp), date_trunc('year', $1::timestamp) + interval '1 year' - interval '1 day', interval '1 month'
            ) AS month
    )
SELECT TO_CHAR(m.month, 'Mon') AS month, COALESCE(SUM(t.amount), 0)::int AS total_amount
FROM
    months m
    LEFT JOIN transactions t ON EXTRACT(
        MONTH
        FROM t.transaction_time
    ) = EXTRACT(
        MONTH
        FROM m.month
    )
    AND EXTRACT(
        YEAR
        FROM t.transaction_time
    ) = EXTRACT(
        YEAR
        FROM m.month
    )
    AND t.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY m.month;

-- GetYearlyAmounts: Retrieves total transaction amount per year over a 5-year span
-- Purpose:
--   Analyze annual growth or decline in transaction volume for trend analysis
-- Parameters:
--   $1: current_year - The most recent year to include (covers current_year - 4 to current_year)
-- Returns:
--   - year: Year of the transaction
--   - total_amount: Total transaction amount for the year
-- Business Logic:
--   - Excludes soft-deleted transactions (deleted_at IS NULL)
-- name: GetYearlyAmounts :many
SELECT EXTRACT(
        YEAR
        FROM t.created_at
    ) AS year, SUM(t.amount) AS total_amount
FROM transactions t
WHERE
    t.deleted_at IS NULL
    AND EXTRACT(
        YEAR
        FROM t.created_at
    ) >= $1 - 4
    AND EXTRACT(
        YEAR
        FROM t.created_at
    ) <= $1
GROUP BY
    EXTRACT(
        YEAR
        FROM t.created_at
    )
ORDER BY year;

-- GetTransactionByCardNumber: Retrieves paginated transactions for a specific card with optional filtering
-- Purpose: View transaction history for a particular card with search capability
-- Parameters:
--   $1: card_number - The card number to filter transactions (exact match)
--   $2: search_term - Optional text to filter by payment method (NULL for no filter)
--   $3: limit - Maximum number of records to return per page
--   $4: offset - Number of records to skip for pagination
-- Returns:
--   All transaction fields plus total_count of matching records
-- Business Logic:
--   - Only returns active transactions (non-deleted records)
--   - Strict matching on card_number combined with optional payment method search
--   - Case-insensitive partial matching on payment_method when search term provided
--   - Orders results by transaction_time (newest transactions first)
--   - Includes pagination metadata via total_count
--   - Useful for cardholder transaction history views and statements
-- name: GetTransactionByCardNumber :many
SELECT *, COUNT(*) OVER () AS total_count
FROM transactions
WHERE
    deleted_at IS NULL
    AND card_number = $1
    AND (
        $2::TEXT IS NULL
        OR payment_method ILIKE '%' || $2 || '%'
    )
ORDER BY transaction_time DESC
LIMIT $3
OFFSET
    $4;

-- GetMonthlyPaymentMethodsByCardNumber: Retrieves a monthly summary of transaction transactions categorized by payment method
-- Purpose:
--   Useful for visualizing how each payment method is used over time within a given year
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: reference_date - Any date within the target year (used to generate monthly range)
-- Returns:
--   - month (e.g., 'Jan', 'Feb')
--   - payment_method (e.g., 'e-wallet', 'bank_transfer')
--   - total_transactions: Number of transactions for the method that month
--   - total_amount: Total amount of transactions for the method that month
-- Business Logic:
--   - Includes all combinations of months and available payment methods (even if 0 data)
--   - Excludes soft-deleted transactions (deleted_at IS NULL)
--   - Uses CROSS JOIN to ensure all months and methods are represent
-- name: GetMonthlyPaymentMethodsByCardNumber :many
WITH
    months AS (
        SELECT generate_series(
                date_trunc('year', $2::timestamp), date_trunc('year', $2::timestamp) + interval '1 year' - interval '1 day', interval '1 month'
            ) AS month
    ),
    payment_methods AS (
        SELECT DISTINCT
            payment_method
        FROM transactions
        WHERE
            deleted_at IS NULL
    )
SELECT
    TO_CHAR(m.month, 'Mon') AS month,
    pm.payment_method,
    COALESCE(COUNT(t.transaction_id), 0)::int AS total_transactions,
    COALESCE(SUM(t.amount), 0)::int AS total_amount
FROM
    months m
    CROSS JOIN payment_methods pm
    LEFT JOIN transactions t ON EXTRACT(
        MONTH
        FROM t.transaction_time
    ) = EXTRACT(
        MONTH
        FROM m.month
    )
    AND EXTRACT(
        YEAR
        FROM t.transaction_time
    ) = EXTRACT(
        YEAR
        FROM m.month
    )
    AND t.payment_method = pm.payment_method
    AND t.card_number = $1
    AND t.deleted_at IS NULL
GROUP BY
    m.month,
    pm.payment_method
ORDER BY m.month, pm.payment_method;

-- GetYearlyPaymentMethodsByCardNumber: Retrieves yearly summary of transaction transactions grouped by payment method over a 5-year span
-- Purpose:
--   Analyze long-term trends of transaction method usage across years
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: current_year - The most recent year to include (covers current_year - 4 to current_year)
-- Returns:
--   - year: Year of transaction (e.g., 2020, 2021)
--   - payment_method
--   - total_transactions: Count of transactions per method per year
--   - total_amount: Sum of amounts per method per year
-- Business Logic:
--   - Filters data within a 5-year window
--   - Excludes soft-deleted transactions (deleted_at IS NULL)
-- name: GetYearlyPaymentMethodsByCardNumber :many
SELECT
    EXTRACT(
        YEAR
        FROM t.created_at
    ) AS year,
    t.payment_method,
    COUNT(t.transaction_id) AS total_transactions,
    SUM(t.amount) AS total_amount
FROM transactions t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
    AND EXTRACT(
        YEAR
        FROM t.created_at
    ) >= $2 - 4
    AND EXTRACT(
        YEAR
        FROM t.created_at
    ) <= $2
GROUP BY
    EXTRACT(
        YEAR
        FROM t.created_at
    ),
    t.payment_method
ORDER BY year;

-- GetMonthlyAmountsByCardNumber: Retrieves total transaction amount per month for a specific year
-- Purpose:
--   Visualize monthly trends in transaction volume for charting/dashboards
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: reference_date - Any date within the target year
-- Returns:
--   - month: 3-letter month abbreviation
--   - total_amount: Sum of transaction amounts in each month
-- Business Logic:
--   - Uses LEFT JOIN to ensure all months are included, even with 0 transactions
--   - Filters out soft-deleted data (deleted_at IS NULL)
-- name: GetMonthlyAmountsByCardNumber :many
WITH
    months AS (
        SELECT generate_series(
                date_trunc('year', $2::timestamp), date_trunc('year', $2::timestamp) + interval '1 year' - interval '1 day', interval '1 month'
            ) AS month
    )
SELECT TO_CHAR(m.month, 'Mon') AS month, COALESCE(SUM(t.amount), 0)::int AS total_amount
FROM
    months m
    LEFT JOIN transactions t ON EXTRACT(
        MONTH
        FROM t.transaction_time
    ) = EXTRACT(
        MONTH
        FROM m.month
    )
    AND EXTRACT(
        YEAR
        FROM t.transaction_time
    ) = EXTRACT(
        YEAR
        FROM m.month
    )
    AND t.card_number = $1
    AND t.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY m.month;

-- GetYearlyAmountsByCardNumber:  Retrieves total transaction amount per year over a 5-year span
-- Purpose:
--   Analyze annual growth or decline in transaction volume for trend analysis
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: current_year - The most recent year to include (covers current_year - 4 to current_year)
-- Returns:
--   - year: Year of the transaction
--   - total_amount: Total transaction amount for the year
-- Business Logic:
--   - Excludes soft-deleted transactions (deleted_at IS NULL)
-- name: GetYearlyAmountsByCardNumber :many
SELECT EXTRACT(
        YEAR
        FROM t.created_at
    ) AS year, SUM(t.amount) AS total_amount
FROM transactions t
WHERE
    t.deleted_at IS NULL
    AND t.card_number = $1
    AND EXTRACT(
        YEAR
        FROM t.created_at
    ) >= $2 - 4
    AND EXTRACT(
        YEAR
        FROM t.created_at
    ) <= $2
GROUP BY
    EXTRACT(
        YEAR
        FROM t.created_at
    )
ORDER BY year;

-- CreateTransactionAtomic: Debits the customer, credits the merchant, and records a successful transaction atomically.
-- Purpose: Keep the two saldo mutations and transaction record in one PostgreSQL statement.
-- Business Logic:
--   - Locks both active saldo rows before validating the available balance.
--   - Rejects missing accounts, identical customer/merchant accounts, and insufficient funds.
--   - Returns no row for a rejected operation; no balance is changed in that case.
-- name: CreateTransactionAtomic :one
WITH active_recipient AS (
    SELECT c.card_number, u.email
    FROM cards c
    JOIN users u ON u.user_id = c.user_id AND u.deleted_at IS NULL
    WHERE c.card_number = sqlc.arg(card_number) AND c.deleted_at IS NULL
), locked_saldos AS MATERIALIZED (
    SELECT s.card_number, s.total_balance
    FROM saldos s
    WHERE s.card_number IN (sqlc.arg(card_number), sqlc.arg(card_number_2))
      AND s.deleted_at IS NULL
    ORDER BY s.card_number
    FOR UPDATE
), eligible AS (
    SELECT 1
    FROM saldos AS user_saldo
    JOIN saldos AS merchant_saldo ON merchant_saldo.card_number = sqlc.arg(card_number_2)
    WHERE user_saldo.card_number = sqlc.arg(card_number)
      AND user_saldo.deleted_at IS NULL
      AND merchant_saldo.deleted_at IS NULL
      AND user_saldo.total_balance >= sqlc.arg(amount)
      AND EXISTS (SELECT 1 FROM locked_saldos)
      AND EXISTS (SELECT 1 FROM active_recipient)
      AND sqlc.arg(amount) > 0
      AND sqlc.arg(card_number) <> sqlc.arg(card_number_2)
), settled_saldos AS (
    UPDATE saldos s
    SET total_balance = s.total_balance + CASE
            WHEN s.card_number = sqlc.arg(card_number) THEN -sqlc.arg(amount)
            ELSE sqlc.arg(amount)
        END,
        updated_at = current_timestamp
    WHERE s.card_number IN (sqlc.arg(card_number), sqlc.arg(card_number_2))
      AND s.deleted_at IS NULL
      AND EXISTS (SELECT 1 FROM eligible)
    RETURNING s.card_number
), inserted_transaction AS (
    INSERT INTO transactions (
        card_number,
        amount,
        payment_method,
        merchant_id,
        transaction_time,
        idempotency_key,
        status,
        created_at,
        updated_at
    )
    SELECT sqlc.arg(card_number), sqlc.arg(amount), sqlc.arg(payment_method), sqlc.arg(merchant_id), sqlc.arg(transaction_time), sqlc.arg(idempotency_key), 'success', current_timestamp, current_timestamp
    WHERE (SELECT COUNT(*) FROM settled_saldos) = 2
      AND EXISTS (SELECT 1 FROM active_recipient)
    RETURNING
        transaction_id,
        transaction_no,
        card_number,
        amount,
        payment_method,
        merchant_id,
        transaction_time,
        status,
        created_at,
        updated_at
), queued_event AS (
    INSERT INTO outbox_events (event_key, topic, message_key, payload)
    SELECT
        'transaction-email:' || it.transaction_id::text,
        'email-service-topic-transaction-create',
        it.transaction_id::text,
        jsonb_build_object(
            'event_id', gen_random_uuid()::text,
            'schema_version', 1,
            'event_type', 'transaction.created',
            'occurred_at', now(),
            'email', COALESCE(u.email, ''),
            'subject', 'Transaction Successful - SanEdge',
            'body', 'Your transaction of ' || it.amount::text || ' has been processed successfully.'
        )
    FROM inserted_transaction it
    JOIN active_recipient u ON u.card_number = it.card_number
    ON CONFLICT (event_key) DO NOTHING
    RETURNING event_key
), queued_cache_event AS (
    INSERT INTO outbox_events (event_key, topic, message_key, payload)
    SELECT
        'transaction-cache-invalidate:' || it.transaction_id::text,
        'saldo-service-topic-invalidate-cache',
        it.transaction_id::text,
        jsonb_build_object(
            'event', 'saldo.cache.invalidate',
            'card_numbers', jsonb_build_array(sqlc.arg(card_number), sqlc.arg(card_number_2))
        )
    FROM inserted_transaction it
    ON CONFLICT (event_key) DO NOTHING
    RETURNING event_key
)
SELECT
    it.transaction_id,
    it.transaction_no,
    it.card_number,
    it.amount,
    it.payment_method,
    it.merchant_id,
    it.transaction_time,
    it.status,
    it.created_at,
    it.updated_at
FROM inserted_transaction it;

-- AuthorizeTransactionAtomic: Places a debit hold or increases a credit outstanding balance
-- and creates an authorized transaction in one PostgreSQL statement.
-- The card row and, for debit cards, the saldo row are locked before the
-- availability check. Any failed precondition rolls back every mutation.
-- name: AuthorizeTransactionAtomic :one
WITH locked_card AS MATERIALIZED (
    SELECT c.card_id, c.card_number, c.card_type, c.outstanding_balance, c.credit_limit
    FROM cards AS c
    WHERE c.card_number = sqlc.arg(card_number)
      AND c.deleted_at IS NULL
      AND c.status = 'active'
    FOR UPDATE
), locked_saldo AS MATERIALIZED (
    SELECT card_number, total_balance
    FROM saldos
    WHERE card_number = sqlc.arg(card_number)
      AND deleted_at IS NULL
    FOR UPDATE
), debit_hold AS (
    UPDATE saldos s
    SET total_balance = s.total_balance - sqlc.arg(amount),
        updated_at = current_timestamp
    FROM locked_card c, locked_saldo ls
    WHERE c.card_type <> 'credit'
      AND sqlc.arg(amount) > 0
      AND ls.total_balance >= sqlc.arg(amount)
      AND s.card_number = ls.card_number
    RETURNING s.card_number
), credit_hold AS (
    UPDATE cards c
    SET outstanding_balance = c.outstanding_balance + sqlc.arg(amount),
        updated_at = current_timestamp
    FROM locked_card lc
    WHERE c.card_id = lc.card_id
      AND lc.card_type = 'credit'
      AND sqlc.arg(amount) > 0
      AND (lc.credit_limit <= 0 OR lc.outstanding_balance + sqlc.arg(amount) <= lc.credit_limit)
    RETURNING c.card_number
), authorized_transaction AS (
    INSERT INTO transactions (
        card_number,
        amount,
        payment_method,
        merchant_id,
        transaction_time,
        idempotency_key,
        status,
        created_at,
        updated_at
    )
    SELECT sqlc.arg(card_number), sqlc.arg(amount), sqlc.arg(payment_method), sqlc.arg(merchant_id), sqlc.arg(transaction_time), sqlc.arg(idempotency_key), 'authorized', current_timestamp, current_timestamp
    WHERE EXISTS (SELECT 1 FROM debit_hold)
       OR EXISTS (SELECT 1 FROM credit_hold)
    RETURNING
        transaction_id,
        transaction_no,
        card_number,
        amount,
        payment_method,
        merchant_id,
        transaction_time,
        status,
        created_at,
        updated_at
), queued_fraud_event AS (
    INSERT INTO outbox_events (event_key, topic, message_key, payload)
    SELECT
        'transaction-fraud:' || at.transaction_id::text,
        'transaction-fraud-check',
        at.transaction_id::text,
        jsonb_build_object(
            'transaction_id', at.transaction_id,
            'card_number', at.card_number,
            'amount', at.amount,
            'card_type', lc.card_type
        )
    FROM authorized_transaction at
    JOIN locked_card lc ON lc.card_number = at.card_number
    ON CONFLICT (event_key) DO NOTHING
    RETURNING event_key
)
SELECT
    at.transaction_id,
    at.transaction_no,
    at.card_number,
    at.amount,
    at.payment_method,
    at.merchant_id,
    at.transaction_time,
    at.status,
    at.created_at,
    at.updated_at
FROM authorized_transaction at;


-- GetTransactionByIdempotencyKey: Returns an existing transaction for a client-supplied idempotency key.
-- Purpose: Support safe replay of transaction create requests.
-- Parameters:
--   $1: idempotency_key - The client-supplied key
-- Business Logic:
--   - Only matches non-empty keys on active records.
--   - Used by the repository to return the original record instead of
--     double-debiting when a request is retried.
-- name: GetTransactionByIdempotencyKey :one
SELECT
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    created_at,
    updated_at
FROM transactions
WHERE idempotency_key = $1
  AND idempotency_key <> '';

-- CreateTransaction: Creates a new transaction record
-- Purpose: Record a financial transaction in the system
-- Parameters:
--   $1: card_number - The card used for the transaction
--   $2: amount - The transaction amount
--   $3: payment_method - Payment method used (e.g., 'credit', 'debit')
--   $4: merchant_id - ID of the merchant where transaction occurred
--   $5: transaction_time - Timestamp of when transaction occurred
-- Returns:
--   The newly created transaction record with all fields
-- Business Logic:
--   - Sets creation and update timestamps automatically
--   - Used for recording purchases, payments, and other financial activities
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
    )
RETURNING
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    created_at,
    updated_at;

-- UpdateTransaction: Modifies an existing transaction's details
-- Purpose: Update transaction information
-- Parameters:
--   $1: transaction_id - ID of transaction to update
--   $2: card_number - Updated card number
--   $3: amount - Updated transaction amount
--   $4: payment_method - Updated payment method
--   $5: merchant_id - Updated merchant ID
--   $6: transaction_time - Updated transaction timestamp
-- Business Logic:
--   - Only updates active transactions (non-deleted)
--   - Automatically updates the modification timestamp
--   - Used for correcting transaction details
-- name: UpdateTransaction :one
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
RETURNING
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    created_at,
    updated_at;

-- UpdateTransactionAtomic: Atomically updates a transaction and settles its
-- card balance in a single SQL statement.
-- Purpose: Race-free update of the transaction amount plus the corresponding
--   balance delta. The delta is computed from the locked pre-update amount, so
--   concurrent updates can never read a stale amount (no lost update, no
--   negative double-apply). Status is set to 'success' only when the balance
--   settlement succeeded. Replaces the old three-step read-modify-write flow.
-- name: UpdateTransactionAtomic :one
WITH locked_transaction AS MATERIALIZED (
    SELECT t.transaction_id, t.transaction_no, t.card_number, t.amount, t.payment_method, t.merchant_id, t.transaction_time, t.status, t.created_at, t.updated_at
    FROM transactions AS t
    WHERE t.transaction_id = sqlc.arg(transaction_id)
      AND t.status IN ('pending', 'failed', 'success', 'compensation_required')
      AND t.deleted_at IS NULL
    FOR UPDATE
), locked_saldos AS MATERIALIZED (
    SELECT s.card_number, s.total_balance
    FROM saldos s
    WHERE s.card_number IN ((SELECT card_number FROM locked_transaction), sqlc.arg(merchant_card_number))
      AND s.deleted_at IS NULL
    ORDER BY s.card_number
    FOR UPDATE
), eligible AS (
    SELECT 1
    FROM locked_transaction lt
    JOIN saldos AS customer_saldo ON customer_saldo.card_number = lt.card_number
    JOIN saldos AS merchant_saldo ON merchant_saldo.card_number = sqlc.arg(merchant_card_number)
    WHERE lt.card_number <> sqlc.arg(merchant_card_number)
      AND customer_saldo.deleted_at IS NULL
      AND merchant_saldo.deleted_at IS NULL
      AND EXISTS (SELECT 1 FROM locked_saldos)
      AND sqlc.arg(amount) > 0
      AND customer_saldo.total_balance + lt.amount - sqlc.arg(amount) >= 0
), settled_saldos AS (
    UPDATE saldos s
    SET total_balance = s.total_balance + CASE
            WHEN s.card_number = lt.card_number THEN lt.amount - sqlc.arg(amount)
            ELSE sqlc.arg(amount) - lt.amount
        END,
        updated_at = current_timestamp
    FROM locked_transaction lt
    WHERE s.card_number IN (lt.card_number, sqlc.arg(merchant_card_number))
      AND s.deleted_at IS NULL
      AND EXISTS (SELECT 1 FROM eligible)
    RETURNING s.card_number
), updated_transaction AS (
    UPDATE transactions t
    SET card_number = sqlc.arg(card_number),
        amount = sqlc.arg(amount),
        payment_method = sqlc.arg(payment_method),
        merchant_id = sqlc.arg(merchant_id),
        transaction_time = sqlc.arg(transaction_time),
        status = 'success',
        updated_at = current_timestamp
    FROM locked_transaction lt
    WHERE t.transaction_id = lt.transaction_id
      AND (SELECT COUNT(*) FROM settled_saldos) = 2
    RETURNING
        t.transaction_id,
        t.transaction_no,
        t.card_number,
        t.amount,
        t.payment_method,
        t.merchant_id,
        t.transaction_time,
        t.status,
        t.created_at,
        t.updated_at
)
SELECT
    ut.transaction_id,
    ut.transaction_no,
    ut.card_number,
    ut.amount,
    ut.payment_method,
    ut.merchant_id,
    ut.transaction_time,
    ut.status,
    ut.created_at,
    ut.updated_at
FROM updated_transaction ut;

-- UpdateTransactionStatus: Changes a transaction's status
-- Purpose: Update transaction processing status
-- Parameters:
--   $1: transaction_id - ID of transaction to update
--   $2: status - New status (e.g., 'success', 'failed', 'pending')
-- Business Logic:
--   - Only updates active transactions
--   - Used to reflect transaction processing outcomes
--   - Important for reconciliation and reporting
-- name: UpdateTransactionStatus :one
UPDATE transactions
SET
    status = $2,
    updated_at = current_timestamp
WHERE
    transaction_id = $1
    AND deleted_at IS NULL
    AND (
        ($2 = 'failed' AND status IN ('pending', 'failed', 'compensation_required'))
        OR ($2 = 'compensation_required' AND status IN ('pending', 'failed', 'success', 'authorized', 'captured', 'compensation_required'))
        OR ($2 = 'success' AND status IN ('pending', 'failed', 'success', 'compensation_required'))
        OR ($2 = 'pending' AND status IN ('failed', 'compensation_required'))
        OR ($2 = 'flagged_fraud' AND status IN ('success', 'authorized', 'captured'))
    )
RETURNING
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    created_at,
    updated_at;

-- CaptureTransactionAtomic: Captures an authorized transaction and credits the
-- merchant in a single SQL statement.
-- Purpose: Only transactions in 'authorized' state can be captured. The row is
--   locked first, so a concurrent double-capture sees the already-'captured'
--   state and affects zero rows (no double merchant credit). The merchant
--   settlement and the status transition commit together or not at all.
-- name: CaptureTransactionAtomic :one
WITH locked_transaction AS MATERIALIZED (
    SELECT t.transaction_id, t.transaction_no, t.card_number, t.amount, t.payment_method, t.merchant_id, t.transaction_time, t.status, t.created_at, t.updated_at
    FROM transactions AS t
    JOIN cards AS customer_card
      ON customer_card.card_number = t.card_number
     AND customer_card.deleted_at IS NULL
    WHERE t.transaction_id = sqlc.arg(transaction_id)
      AND t.status = 'authorized'
      AND t.deleted_at IS NULL
    FOR UPDATE OF t
), settled_merchant AS (
    UPDATE saldos s
    SET total_balance = s.total_balance + lt.amount,
        updated_at = current_timestamp
    FROM locked_transaction lt
    WHERE s.card_number = sqlc.arg(merchant_card_number)
      AND s.deleted_at IS NULL
    RETURNING s.card_number
), reward_ledger AS (
    INSERT INTO card_rewards (
        card_number, txn_id, amount, mcc, points_earned, expires_at, redeemed
    )
    SELECT
        lt.card_number,
        lt.transaction_id::text,
        lt.amount,
        '',
        GREATEST((lt.amount / 10000)::integer, 1),
        current_timestamp + interval '1 year',
        FALSE
    FROM locked_transaction lt
    WHERE EXISTS (SELECT 1 FROM settled_merchant)
    ON CONFLICT (txn_id) WHERE txn_id <> '' DO NOTHING
    RETURNING card_number, points_earned
), earned_rewards AS (
    UPDATE cards c
    SET reward_points = c.reward_points + rl.points_earned,
        updated_at = current_timestamp
    FROM reward_ledger rl
    WHERE c.card_number = rl.card_number
      AND c.deleted_at IS NULL
    RETURNING c.card_number
), captured_transaction AS (

    UPDATE transactions t
    SET status = 'captured',
        updated_at = current_timestamp
    FROM locked_transaction lt
    WHERE t.transaction_id = lt.transaction_id
      AND EXISTS (SELECT 1 FROM settled_merchant)
      AND EXISTS (SELECT 1 FROM earned_rewards)
    RETURNING
        t.transaction_id,
        t.transaction_no,
        t.card_number,
        t.amount,
        t.payment_method,
        t.merchant_id,
        t.transaction_time,
        t.status,
        t.created_at,
        t.updated_at
)
SELECT
    ct.transaction_id,
    ct.transaction_no,
    ct.card_number,
    ct.amount,
    ct.payment_method,
    ct.merchant_id,
    ct.transaction_time,
    ct.status,
    ct.created_at,
    ct.updated_at
FROM captured_transaction ct;

-- VoidTransactionAtomic: Voids a pending/authorized transaction and releases the
-- hold in a single SQL statement.
-- Purpose: Only 'pending' or 'authorized' transactions can be voided. The row
--   is locked first, so a concurrent double-void sees the already-'voided'
--   state and affects zero rows (no double hold release). For debit cards the
--   held amount is credited back to the card saldo; credit cards carry no
--   balance hold so they only transition status. The settlement and the status
--   transition commit together or not at all.
-- name: VoidTransactionAtomic :one
WITH locked_transaction AS MATERIALIZED (
    SELECT t.transaction_id, t.transaction_no, t.card_number, t.amount, t.payment_method, t.merchant_id, t.transaction_time, t.status, t.created_at, t.updated_at, c.card_type
    FROM transactions t
    JOIN cards c ON c.card_number = t.card_number AND c.deleted_at IS NULL
    WHERE t.transaction_id = sqlc.arg(transaction_id)
      AND t.status IN ('pending', 'authorized')
      AND t.deleted_at IS NULL
    FOR UPDATE OF t
), released_hold AS (
    UPDATE saldos s
    SET total_balance = s.total_balance + lt.amount,
        updated_at = current_timestamp
    FROM locked_transaction lt
    WHERE lt.card_type = 'debit'
      AND s.card_number = lt.card_number
      AND s.deleted_at IS NULL
    RETURNING s.card_number
), released_credit_hold AS (
    UPDATE cards c
    SET outstanding_balance = c.outstanding_balance - lt.amount,
        updated_at = current_timestamp
    FROM locked_transaction lt
    WHERE lt.card_type = 'credit'
      AND c.card_number = lt.card_number
      AND c.deleted_at IS NULL
      AND c.outstanding_balance >= lt.amount
    RETURNING c.card_number
), voided_transaction AS (
    UPDATE transactions t
    SET status = 'voided',
        updated_at = current_timestamp
    FROM locked_transaction lt
    WHERE t.transaction_id = lt.transaction_id
      AND ((lt.card_type = 'credit' AND EXISTS (SELECT 1 FROM released_credit_hold))
           OR (lt.card_type <> 'credit' AND EXISTS (SELECT 1 FROM released_hold)))
    RETURNING
        t.transaction_id,
        t.transaction_no,
        t.card_number,
        t.amount,
        t.payment_method,
        t.merchant_id,
        t.transaction_time,
        t.status,
        t.created_at,
        t.updated_at
)
SELECT
    vt.transaction_id,
    vt.transaction_no,
    vt.card_number,
    vt.amount,
    vt.payment_method,
    vt.merchant_id,
    vt.transaction_time,
    vt.status,
    vt.created_at,
    vt.updated_at
FROM voided_transaction vt;

-- RefundTransactionAtomic: Refunds a captured transaction and reverses the
-- settlement in a single SQL statement.
-- Purpose: Only 'captured' transactions can be refunded. The row is locked
--   first, so a concurrent double-refund sees the already-'refunded' state and
--   affects zero rows (no double customer credit / merchant debit). The
--   customer credit (debit card) or outstanding reduction (credit card), the
--   merchant debit with a non-negative guard, and the status transition commit
--   together or not at all.
-- name: RefundTransactionAtomic :one
WITH locked_transaction AS MATERIALIZED (
    SELECT t.transaction_id, t.transaction_no, t.card_number, t.amount, t.payment_method, t.merchant_id, t.transaction_time, t.status, t.created_at, t.updated_at, c.card_type
    FROM transactions t
    JOIN cards c ON c.card_number = t.card_number AND c.deleted_at IS NULL
    WHERE t.transaction_id = sqlc.arg(transaction_id)
      AND t.status = 'captured'
      AND t.deleted_at IS NULL
    FOR UPDATE OF t
), locked_saldos AS MATERIALIZED (
    SELECT s.card_number, s.total_balance
    FROM saldos s
    WHERE s.card_number IN ((SELECT card_number FROM locked_transaction), sqlc.arg(merchant_card_number))
      AND s.deleted_at IS NULL
    ORDER BY s.card_number
    FOR UPDATE
), eligible AS (
    SELECT 1
    FROM locked_transaction lt
    LEFT JOIN locked_saldos customer_saldo ON customer_saldo.card_number = lt.card_number
    JOIN locked_saldos merchant_saldo ON merchant_saldo.card_number = sqlc.arg(merchant_card_number)
    WHERE lt.card_number <> sqlc.arg(merchant_card_number)
      AND merchant_saldo.total_balance >= lt.amount
      AND (lt.card_type = 'credit' OR customer_saldo.card_number IS NOT NULL)
), customer_refunded AS (
    UPDATE saldos s
    SET total_balance = s.total_balance + lt.amount,
        updated_at = current_timestamp
    FROM locked_transaction lt
    WHERE lt.card_type = 'debit'
      AND s.card_number = lt.card_number
      AND s.deleted_at IS NULL
      AND EXISTS (SELECT 1 FROM eligible)
    RETURNING s.card_number
), outstanding_reduced AS (
    UPDATE cards c
    SET outstanding_balance = GREATEST(0, c.outstanding_balance - lt.amount),
        updated_at = current_timestamp
    FROM locked_transaction lt
    WHERE lt.card_type = 'credit'
      AND c.card_number = lt.card_number
      AND c.deleted_at IS NULL
      AND EXISTS (SELECT 1 FROM eligible)
    RETURNING c.card_number
), merchant_debited AS (
    UPDATE saldos s
    SET total_balance = s.total_balance - lt.amount,
        updated_at = current_timestamp
    FROM locked_transaction lt
    WHERE s.card_number = sqlc.arg(merchant_card_number)
      AND s.deleted_at IS NULL
      AND s.total_balance - lt.amount >= 0
      AND EXISTS (SELECT 1 FROM eligible)
    RETURNING s.card_number
), refunded_transaction AS (
    UPDATE transactions t
    SET status = 'refunded',
        updated_at = current_timestamp
    FROM locked_transaction lt
    WHERE t.transaction_id = lt.transaction_id
      AND EXISTS (SELECT 1 FROM merchant_debited)
      AND (EXISTS (SELECT 1 FROM customer_refunded) OR EXISTS (SELECT 1 FROM outstanding_reduced))
    RETURNING
        t.transaction_id,
        t.transaction_no,
        t.card_number,
        t.amount,
        t.payment_method,
        t.merchant_id,
        t.transaction_time,
        t.status,
        t.created_at,
        t.updated_at
)
SELECT
    rt.transaction_id,
    rt.transaction_no,
    rt.card_number,
    rt.amount,
    rt.payment_method,
    rt.merchant_id,
    rt.transaction_time,
    rt.status,
    rt.created_at,
    rt.updated_at
FROM refunded_transaction rt;

-- TrashTransaction: Soft-deletes a transaction record
-- Purpose: Remove transaction from active use without permanent deletion
-- Parameters:
--   $1: transaction_id - ID of transaction to trash
-- Business Logic:
--   - Sets deleted_at timestamp
--   - Preserves data for audit/recovery purposes
--   - Only affects currently active records
-- name: TrashTransaction :one
UPDATE transactions
SET
    deleted_at = current_timestamp
WHERE
    transaction_id = $1
    AND deleted_at IS NULL
RETURNING
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    created_at,
    updated_at,
    deleted_at;

-- RestoreTransaction: Recovers a soft-deleted transaction
-- Purpose: Reactivate a previously trashed transaction
-- Parameters:
--   $1: transaction_id - ID of transaction to restore
-- Business Logic:
--   - Clears the deleted_at timestamp
--   - Only works on currently trashed records
--   - Used for data recovery purposes
-- name: RestoreTransaction :one
UPDATE transactions
SET
    deleted_at = NULL
WHERE
    transaction_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    created_at,
    updated_at,
    deleted_at;

-- DeleteTransactionPermanently: Hard-deletes a trashed transaction
-- Purpose: Permanently remove a transaction from the system
-- Parameters:
--   $1: transaction_id - ID of transaction to delete
-- Business Logic:
--   - Physical deletion from database
--   - Only works on already trashed records
--   - Irreversible operation
--   - Used for data cleanup after retention period
-- name: DeleteTransactionPermanently :exec
DELETE FROM transactions
WHERE
    transaction_id = $1
    AND deleted_at IS NOT NULL;

-- RestoreAllTransactions: Recovers all trashed transactions
-- Purpose: Mass restoration of deleted transactions
-- Business Logic:
--   - Clears deleted_at for all trashed records
--   - Useful for system recovery scenarios
--   - Should be used cautiously in production
-- name: RestoreAllTransactions :exec
UPDATE transactions
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentTransactions: Permanently removes all trashed transactions
-- Purpose: Clean up all soft-deleted transaction records
-- Business Logic:
--   - Irreversible bulk deletion
--   - Only affects records marked as deleted
--   - Frees database space from old records
--   - Typically used during maintenance periods
-- name: DeleteAllPermanentTransactions :exec
DELETE FROM transactions WHERE deleted_at IS NOT NULL;

-- UpdateTransactionStatusNew: Updates transaction status
-- Purpose: Updates the transaction status to a new value
-- Parameters:
--   $1: transaction_id - The ID of the transaction
--   $2: status - New transaction status
-- Returns: Updated transaction record
-- Business Logic:
--   - Wraps the standard status update
--   - Consistent with existing patterns
-- name: UpdateTransactionStatusDirect :one
UPDATE transactions
SET
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    transaction_id = $1
    AND deleted_at IS NULL
RETURNING
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    fraud_score,
    created_at,
    updated_at;

-- UpdateTransactionFraudScore: Updates the fraud score for a transaction
-- Purpose: Set or update the fraud score after fraud checking
-- Parameters:
--   $1: transaction_id - The ID of the transaction
--   $2: fraud_score - The calculated fraud score (0-100)
-- Returns: Updated transaction record with fraud score
-- Business Logic:
--   - Fraud score is calculated asynchronously
--   - Higher scores indicate higher fraud probability
--   - Score >= 80 triggers card suspension
-- name: UpdateTransactionFraudScore :one
UPDATE transactions
SET
    fraud_score = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    transaction_id = $1
    AND deleted_at IS NULL
RETURNING
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    fraud_score,
    created_at,
    updated_at;

-- GetTransactionVelocity: Counts transactions within a time window for velocity check
-- Purpose: Detect rapid-fire transactions (fraud pattern)
-- Parameters:
--   $1: card_number - The card number to check
--   $2: since_time - Start of the time window for velocity check
-- Returns: Count of transactions within the time window
-- Business Logic:
--   - Used for fraud detection: >3 transactions in 10 seconds = suspicious
--   - Only counts active transactions
--   - Excludes soft-deleted records
-- name: GetTransactionVelocity :one
SELECT COUNT(*)::INT AS transaction_count
FROM transactions
WHERE
    deleted_at IS NULL
    AND card_number = $1
    AND created_at >= $2;

-- GetTransactionByIDWithFraudScore: Retrieves a single transaction with fraud score
-- Purpose: Get transaction details including fraud-related fields
-- Parameters:
--   $1: transaction_id - The ID of the transaction to retrieve
-- Returns: All fields including fraud_score for the specified transaction
-- Business Logic:
--   - Includes fraud_score field
--   - Only returns active (non-deleted) transactions
-- name: GetTransactionByIDWithFraudScore :one
SELECT
    transaction_id,
    transaction_no,
    card_number,
    amount,
    payment_method,
    merchant_id,
    transaction_time,
    status,
    fraud_score,
    created_at,
    updated_at
FROM transactions
WHERE
    transaction_id = $1
    AND deleted_at IS NULL;