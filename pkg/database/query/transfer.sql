-- GetTransfers: Retrieves paginated transfer records with search capability
-- Purpose: List all active transfers for management UI with filtering options
-- Parameters:
--   $1: search_term - Optional text to filter transfers by source or destination account (NULL for no filter)
--   $2: limit - Maximum number of records to return per page
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All transfer fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted transfers (deleted_at IS NULL)
--   - Supports partial text matching on transfer_from and transfer_to fields (case-insensitive)
--   - Orders by transfer_time (newest transfers first)
--   - Provides total_count for pagination calculations
--   - Useful for transfer monitoring and auditing
-- name: GetTransfers :many
SELECT
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM transfers
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR transfer_from ILIKE '%' || $1 || '%'
        OR transfer_to ILIKE '%' || $1 || '%'
    )
ORDER BY transfer_time DESC
LIMIT $2
OFFSET
    $3;

-- GetTransferByID: Retrieves a single transfer by its ID
-- Purpose: Get detailed information about a specific transfer
-- Parameters:
--   $1: transfer_id - The ID of the transfer to retrieve
-- Returns:
--   All fields for the specified transfer or NULL if not found/deleted
-- Business Logic:
--   - Only returns active transfers (deleted_at IS NULL)
--   - Useful for transfer verification and detailed viewing
-- name: GetTransferByID :one
SELECT
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at
FROM transfers
WHERE
    transfer_id = $1
    AND deleted_at IS NULL;

-- GetActiveTransfers: Retrieves paginated active transfers with search
-- Purpose: List all non-deleted transfers with filtering options
-- Parameters:
--   $1: search_term - Optional text to filter by source or destination account
--   $2: limit - Maximum records to return per page
--   $3: offset - Records to skip for pagination
-- Returns:
--   All transfer fields plus total_count of matching active records
-- Business Logic:
--   - Only includes active transfers (deleted_at IS NULL)
--   - Filters on transfer_from and transfer_to fields
--   - Orders by transfer_time (newest first)
--   - Provides pagination metadata
--   - Used in transfer management interfaces
-- name: GetActiveTransfers :many
SELECT
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM transfers
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR transfer_from ILIKE '%' || $1 || '%'
        OR transfer_to ILIKE '%' || $1 || '%'
    )
ORDER BY transfer_time DESC
LIMIT $2
OFFSET
    $3;

-- GetTrashedTransfers: Retrieves paginated soft-deleted transfers
-- Purpose: List all deleted transfers for recovery or audit purposes
-- Parameters:
--   $1: search_term - Optional text to filter deleted transfers
--   $2: limit - Maximum records to return per page
--   $3: offset - Records to skip for pagination
-- Returns:
--   All transfer fields plus total_count of matching deleted records
-- Business Logic:
--   - Only includes soft-deleted transfers (deleted_at IS NOT NULL)
--   - Same filtering capabilities as active transfers
--   - Maintains newest-first ordering
--   - Used in admin interfaces for transfer recovery
-- name: GetTrashedTransfers :many
SELECT
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM transfers
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR transfer_from ILIKE '%' || $1 || '%'
        OR transfer_to ILIKE '%' || $1 || '%'
    )
ORDER BY transfer_time DESC
LIMIT $2
OFFSET
    $3;

-- GetTransfersByCardNumber:  Retrieves all transfers where the given card number is either the sender or the receiver
-- Purpose:
--   Useful for displaying all transfer history related to a specific card
-- Parameters:
--   $1: card_number - Card number to search for in both transfer_from and transfer_to
-- Returns:
--   All transfer columns for matched records
-- Business Logic:
--   - Excludes soft-deleted records (deleted_at IS NULL)
--   - Sorted by most recent transfer first (DESC)
-- name: GetTransfersByCardNumber :many
SELECT
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at
FROM transfers
WHERE
    deleted_at IS NULL
    AND (
        transfer_from = $1
        OR transfer_to = $1
    )
ORDER BY transfer_time DESC;

-- GetTransfersBySourceCard: Retrieves all transfers where the specified card is the source (transfer_from)
-- Purpose:
--   Track outgoing transfer history for auditing or user activity
-- Parameters:
--   $1: card_number - The source card number
-- Returns:
--   All transfer columns for matched records
-- Business Logic:
--   - Excludes soft-deleted records (deleted_at IS NULL)
--   - Sorted by most recent transfer first (DESC)
-- name: GetTransfersBySourceCard :many
SELECT
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at
FROM transfers
WHERE
    deleted_at IS NULL
    AND transfer_from = $1
ORDER BY transfer_time DESC;

-- GetTransfersByDestinationCard: Retrieves all transfers where the specified card is the destination (transfer_to)
-- Purpose:
--   Track incoming transfers for a user or card
-- Parameters:
--   $1: card_number - The destination card number
-- Returns:
--   All transfer columns for matched records
-- Business Logic:
--   - Excludes soft-deleted records (deleted_at IS NULL)
--   - Sorted by most recent transfer first (DESC)
-- name: GetTransfersByDestinationCard :many
SELECT
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at
FROM transfers
WHERE
    deleted_at IS NULL
    AND transfer_to = $1
ORDER BY transfer_time DESC;

-- GetTrashedTransferByID: Retrieves a single soft-deleted transfer by its ID
-- Purpose:
--   Used for viewing trashed data or restoring transfers
-- Parameters:
--   $1: transfer_id - ID of the transfer
-- Returns:
--   The transfer row if it exists and is soft-deleted
-- Business Logic:
--   - Includes only soft-deleted records (deleted_at IS NOT NULL)
-- name: GetTrashedTransferByID :one
SELECT *
FROM transfers
WHERE
    transfer_id = $1
    AND deleted_at IS NOT NULL;

-- GetMonthTransferStatusSuccess: Retrieves monthly success metrics for fund transfers
-- Purpose: Analyze successful transfer trends across comparison periods
-- Parameters:
--   $1: period1_start - Start date of first comparison period (timestamp)
--   $2: period1_end - End date of first comparison period (timestamp)
--   $3: period2_start - Start date of second comparison period (timestamp)
--   $4: period2_end - End date of second comparison period (timestamp)
-- Returns:
--   year: Year as text (e.g., '2023')
--   month: 3-letter month abbreviation (e.g., 'Jan')
--   total_success: Count of successful transfers
--   total_amount: Sum of successful transfer amounts
-- Business Logic:
--   - Only includes successful transfers (status = 'success')
--   - Covers two customizable time periods for comparison
--   - Zero-fills months with no transfer activity
--   - Formats output for consistent visualization (year as text, month as 'Mon')
--   - Orders by year and month (newest first)
--   - Useful for identifying seasonal transfer patterns and cash flow analysis
-- name: GetMonthTransferStatusSuccess :many
WITH
    monthly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transfer_time
            )::integer AS year,
            EXTRACT(
                MONTH
                FROM t.transfer_time
            )::integer AS month,
            COUNT(*) AS total_success,
            COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
        FROM transfers t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'success'
            AND (
                (
                    t.transfer_time >= $1::timestamp
                    AND t.transfer_time <= $2::timestamp
                )
                OR (
                    t.transfer_time >= $3::timestamp
                    AND t.transfer_time <= $4::timestamp
                )
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transfer_time
            ),
            EXTRACT(
                MONTH
                FROM t.transfer_time
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

-- GetYearlyTransferStatusSuccess: Retrieves yearly success metrics for fund transfers
-- Purpose: Compare annual successful transfer performance year-over-year
-- Parameters:
--   $1: current_year - The target year (includes this year and previous)
-- Returns:
--   year: Year as text (e.g., '2023')
--   total_success: Count of successful transfers
--   total_amount: Sum of successful transfer amounts
-- Business Logic:
--   - Only includes successful transfers (status = 'success')
--   - Compares current year with previous year
--   - Zero-fills years with no transfer activity
--   - Orders by year (newest first)
--   - Useful for year-over-year growth analysis and financial reporting
--   - Helps identify annual transfer volume and money movement trends
-- name: GetYearlyTransferStatusSuccess :many
WITH
    yearly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transfer_time
            )::integer AS year,
            COUNT(*) AS total_success,
            COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
        FROM transfers t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'success'
            AND (
                EXTRACT(
                    YEAR
                    FROM t.transfer_time
                ) = $1::integer
                OR EXTRACT(
                    YEAR
                    FROM t.transfer_time
                ) = $1::integer - 1
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transfer_time
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

-- GetMonthTransferStatusFailed: Retrieves monthly failed metrics for fund transfers
-- Purpose: Analyze failedful transfer trends across comparison periods
-- Parameters:
--   $1: period1_start - Start date of first comparison period (timestamp)
--   $2: period1_end - End date of first comparison period (timestamp)
--   $3: period2_start - Start date of second comparison period (timestamp)
--   $4: period2_end - End date of second comparison period (timestamp)
-- Returns:
--   year: Year as text (e.g., '2023')
--   month: 3-letter month abbreviation (e.g., 'Jan')
--   total_failed: Count of failedful transfers
--   total_amount: Sum of failedful transfer amounts
-- Business Logic:
--   - Only includes failedful transfers (status = 'failed')
--   - Covers two customizable time periods for comparison
--   - Zero-fills months with no transfer activity
--   - Formats output for consistent visualization (year as text, month as 'Mon')
--   - Orders by year and month (newest first)
--   - Useful for identifying seasonal transfer patterns and cash flow analysis
-- name: GetMonthTransferStatusFailed :many
WITH
    monthly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transfer_time
            )::integer AS year,
            EXTRACT(
                MONTH
                FROM t.transfer_time
            )::integer AS month,
            COUNT(*) AS total_failed,
            COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
        FROM transfers t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'failed'
            AND (
                (
                    t.transfer_time >= $1::timestamp
                    AND t.transfer_time <= $2::timestamp
                )
                OR (
                    t.transfer_time >= $3::timestamp
                    AND t.transfer_time <= $4::timestamp
                )
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transfer_time
            ),
            EXTRACT(
                MONTH
                FROM t.transfer_time
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

-- GetYearlyTransferStatusFailed: Retrieves yearly failed metrics for fund transfers
-- Purpose: Compare annual failedful transfer performance year-over-year
-- Parameters:
--   $1: current_year - The target year (includes this year and previous)
-- Returns:
--   year: Year as text (e.g., '2023')
--   total_failed: Count of failedful transfers
--   total_amount: Sum of failedful transfer amounts
-- Business Logic:
--   - Only includes failedful transfers (status = 'failed')
--   - Compares current year with previous year
--   - Zero-fills years with no transfer activity
--   - Orders by year (newest first)
--   - Useful for year-over-year growth analysis and financial reporting
--   - Helps identify annual transfer volume and money movement trends
-- name: GetYearlyTransferStatusFailed :many
WITH
    yearly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transfer_time
            )::integer AS year,
            COUNT(*) AS total_failed,
            COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
        FROM transfers t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'failed'
            AND (
                EXTRACT(
                    YEAR
                    FROM t.transfer_time
                ) = $1::integer
                OR EXTRACT(
                    YEAR
                    FROM t.transfer_time
                ) = $1::integer - 1
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transfer_time
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

-- GetMonthTransferStatusSuccessCardNumber: Retrieves monthly success metrics for fund transfers
-- Purpose: Analyze successful transfer trends across comparison periods
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: period1_start - Start date of first comparison period (timestamp)
--   $3: period1_end - End date of first comparison period (timestamp)
--   $4: period2_start - Start date of second comparison period (timestamp)
--   $5: period2_end - End date of second comparison period (timestamp)
-- Returns:
--   year: Year as text (e.g., '2023')
--   month: 3-letter month abbreviation (e.g., 'Jan')
--   total_success: Count of successful transfers
--   total_amount: Sum of successful transfer amounts
-- Business Logic:
--   - Only includes successful transfers (status = 'success')
--   - Covers two customizable time periods for comparison
--   - Zero-fills months with no transfer activity
--   - Formats output for consistent visualization (year as text, month as 'Mon')
--   - Orders by year and month (newest first)
--   - Useful for identifying seasonal transfer patterns and cash flow analysis
-- name: GetMonthTransferStatusSuccessCardNumber :many
WITH
    monthly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transfer_time
            )::integer AS year,
            EXTRACT(
                MONTH
                FROM t.transfer_time
            )::integer AS month,
            COUNT(*) AS total_success,
            COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
        FROM transfers t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'success'
            AND (
                t.transfer_from = $1
                OR t.transfer_to = $1
            )
            AND (
                (
                    t.transfer_time >= $2::timestamp
                    AND t.transfer_time <= $3::timestamp
                )
                OR (
                    t.transfer_time >= $4::timestamp
                    AND t.transfer_time <= $5::timestamp
                )
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transfer_time
            ),
            EXTRACT(
                MONTH
                FROM t.transfer_time
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
                FROM $4::timestamp
            )::text AS year,
            TO_CHAR($4::timestamp, 'Mon') AS month,
            0 AS total_success,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $4::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $4::timestamp
                    )::integer
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC, TO_DATE(month, 'Mon') DESC;

-- GetYearlyTransferStatusSuccessCardNumber: Retrieves yearly success metrics for fund transfers
-- Purpose: Compare annual successful transfer performance year-over-year
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: current_year - The target year (includes this year and previous)
-- Returns:
--   year: Year as text (e.g., '2023')
--   total_success: Count of successful transfers
--   total_amount: Sum of successful transfer amounts
-- Business Logic:
--   - Only includes successful transfers (status = 'success')
--   - Compares current year with previous year
--   - Zero-fills years with no transfer activity
--   - Orders by year (newest first)
--   - Useful for year-over-year growth analysis and financial reporting
--   - Helps identify annual transfer volume and money movement trends
-- name: GetYearlyTransferStatusSuccessCardNumber :many
WITH
    yearly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transfer_time
            )::integer AS year,
            COUNT(*) AS total_success,
            COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
        FROM transfers t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'success'
            AND (
                t.transfer_from = $1
                OR t.transfer_to = $1
            )
            AND (
                EXTRACT(
                    YEAR
                    FROM t.transfer_time
                ) = $2::integer
                OR EXTRACT(
                    YEAR
                    FROM t.transfer_time
                ) = $2::integer - 1
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transfer_time
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

-- GetMonthTransferStatusFailedCardNumber: Retrieves monthly failed metrics for fund transfers
-- Purpose: Analyze failedful transfer trends across comparison periods
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: period1_start - Start date of first comparison period (timestamp)
--   $3: period1_end - End date of first comparison period (timestamp)
--   $4: period2_start - Start date of second comparison period (timestamp)
--   $5: period2_end - End date of second comparison period (timestamp)
-- Returns:
--   year: Year as text (e.g., '2023')
--   month: 3-letter month abbreviation (e.g., 'Jan')
--   total_failed: Count of failedful transfers
--   total_amount: Sum of failedful transfer amounts
-- Business Logic:
--   - Only includes failedful transfers (status = 'failed')
--   - Covers two customizable time periods for comparison
--   - Zero-fills months with no transfer activity
--   - Formats output for consistent visualization (year as text, month as 'Mon')
--   - Orders by year and month (newest first)
--   - Useful for identifying seasonal transfer patterns and cash flow analysis
-- name: GetMonthTransferStatusFailedCardNumber :many
WITH
    monthly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transfer_time
            )::integer AS year,
            EXTRACT(
                MONTH
                FROM t.transfer_time
            )::integer AS month,
            COUNT(*) AS total_failed,
            COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
        FROM transfers t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'failed'
            AND (
                t.transfer_from = $1
                OR t.transfer_to = $1
            )
            AND (
                (
                    t.transfer_time >= $2::timestamp
                    AND t.transfer_time <= $3::timestamp
                )
                OR (
                    t.transfer_time >= $4::timestamp
                    AND t.transfer_time <= $5::timestamp
                )
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transfer_time
            ),
            EXTRACT(
                MONTH
                FROM t.transfer_time
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
                FROM $4::timestamp
            )::text AS year,
            TO_CHAR($4::timestamp, 'Mon') AS month,
            0 AS total_failed,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $4::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $4::timestamp
                    )::integer
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC, TO_DATE(month, 'Mon') DESC;

-- GetYearlyTransferStatusFailedCardNumber: Retrieves yearly failed metrics for fund transfers
-- Purpose: Compare annual failedful transfer performance year-over-year
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: current_year - The target year (includes this year and previous)
-- Returns:
--   year: Year as text (e.g., '2023')
--   total_failed: Count of failedful transfers
--   total_amount: Sum of failedful transfer amounts
-- Business Logic:
--   - Only includes failedful transfers (status = 'failed')
--   - Compares current year with previous year
--   - Zero-fills years with no transfer activity
--   - Orders by year (newest first)
--   - Useful for year-over-year growth analysis and financial reporting
--   - Helps identify annual transfer volume and money movement trends
-- name: GetYearlyTransferStatusFailedCardNumber :many
WITH
    yearly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.transfer_time
            )::integer AS year,
            COUNT(*) AS total_failed,
            COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
        FROM transfers t
        WHERE
            t.deleted_at IS NULL
            AND t.status = 'failed'
            AND (
                t.transfer_from = $1
                OR t.transfer_to = $1
            )
            AND (
                EXTRACT(
                    YEAR
                    FROM t.transfer_time
                ) = $2::integer
                OR EXTRACT(
                    YEAR
                    FROM t.transfer_time
                ) = $2::integer - 1
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.transfer_time
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

-- GetMonthlyTransferAmounts: Retrieves monthly transfer amounts
-- Purpose: Track total transfer amounts over each month of the selected year
-- Parameters:
--   $1: reference_date - Any date within the target year (used to define monthly range)
-- Returns:
--   month: 3-letter month abbreviation (e.g., 'Jan')
--   total_transfer_amount: Sum of transfer amounts in that month
-- Business Logic:
--   - Generates complete monthly series for the target year
--   - Includes all transfers regardless of method
--   - Filters out soft-deleted records (deleted_at IS NULL)
--   - Zero-fills months with no transfer activity
--   - Useful for visualizing monthly cash flow patterns
-- name: GetMonthlyTransferAmounts :many
WITH
    months AS (
        SELECT generate_series(
                date_trunc('year', $1::timestamp), date_trunc('year', $1::timestamp) + interval '1 year' - interval '1 day', interval '1 month'
            ) AS month
    )
SELECT TO_CHAR(m.month, 'Mon') AS month, COALESCE(SUM(t.transfer_amount), 0)::int AS total_transfer_amount
FROM months m
    LEFT JOIN transfers t ON EXTRACT(
        MONTH
        FROM t.transfer_time
    ) = EXTRACT(
        MONTH
        FROM m.month
    )
    AND EXTRACT(
        YEAR
        FROM t.transfer_time
    ) = EXTRACT(
        YEAR
        FROM m.month
    )
    AND t.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY m.month;

-- GetYearlyTransferAmounts: Retrieves yearly transfer amounts
-- Purpose: Analyze total transfer amounts over a 5-year period
-- Parameters:
--   $1: current_year - The final year to include (e.g., 2024), includes 5-year span (current_year - 4)
-- Returns:
--   year: 4-digit year
--   total_transfer_amount: Sum of transfer amounts in that year
-- Business Logic:
--   - Covers a 5-year window (current_year - 4 to current_year)
--   - Filters out soft-deleted records (deleted_at IS NULL)
--   - Groups by calendar year
--   - Useful for identifying long-term money movement trends
-- name: GetYearlyTransferAmounts :many
SELECT EXTRACT(
        YEAR
        FROM t.created_at
    ) AS year, SUM(t.transfer_amount) AS total_transfer_amount
FROM transfers t
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

-- GetMonthlyTransferAmountsBySenderCardNumber: Retrieves monthly transfer amounts
-- Purpose: Track total transfer amounts over each month of the selected year
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: reference_date - Any date within the target year (used to define monthly range)
-- Returns:
--   month: 3-letter month abbreviation (e.g., 'Jan')
--   total_transfer_amount: Sum of transfer amounts in that month
-- Business Logic:
--   - Generates complete monthly series for the target year
--   - Includes all transfers regardless of method
--   - Filters out soft-deleted records (deleted_at IS NULL)
--   - Zero-fills months with no transfer activity
--   - Useful for visualizing monthly cash flow patterns
-- name: GetMonthlyTransferAmountsBySenderCardNumber :many
WITH
    months AS (
        SELECT generate_series(
                date_trunc('year', $2::timestamp), date_trunc('year', $2::timestamp) + interval '1 year' - interval '1 day', interval '1 month'
            ) AS month
    )
SELECT TO_CHAR(m.month, 'Mon') AS month, COALESCE(SUM(t.transfer_amount), 0)::int AS total_transfer_amount
FROM
    months m
    LEFT JOIN transfers t ON EXTRACT(
        MONTH
        FROM t.transfer_time
    ) = EXTRACT(
        MONTH
        FROM m.month
    )
    AND EXTRACT(
        YEAR
        FROM t.transfer_time
    ) = EXTRACT(
        YEAR
        FROM m.month
    )
    AND t.transfer_from = $1
    AND t.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY m.month;

-- GetMonthlyTransferAmountsByReceiverCardNumber: Retrieves monthly transfer amounts
-- Purpose: Track total transfer amounts over each month of the selected year
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: reference_date - Any date within the target year (used to define monthly range)
-- Returns:
--   month: 3-letter month abbreviation (e.g., 'Jan')
--   total_transfer_amount: Sum of transfer amounts in that month
-- Business Logic:
--   - Generates complete monthly series for the target year
--   - Includes all transfers regardless of method
--   - Filters out soft-deleted records (deleted_at IS NULL)
--   - Zero-fills months with no transfer activity
--   - Useful for visualizing monthly cash flow patterns
-- name: GetMonthlyTransferAmountsByReceiverCardNumber :many
WITH
    months AS (
        SELECT generate_series(
                date_trunc('year', $2::timestamp), date_trunc('year', $2::timestamp) + interval '1 year' - interval '1 day', interval '1 month'
            ) AS month
    )
SELECT TO_CHAR(m.month, 'Mon') AS month, COALESCE(SUM(t.transfer_amount), 0)::int AS total_transfer_amount
FROM
    months m
    LEFT JOIN transfers t ON EXTRACT(
        MONTH
        FROM t.transfer_time
    ) = EXTRACT(
        MONTH
        FROM m.month
    )
    AND EXTRACT(
        YEAR
        FROM t.transfer_time
    ) = EXTRACT(
        YEAR
        FROM m.month
    )
    AND t.transfer_to = $1
    AND t.deleted_at IS NULL
GROUP BY
    m.month
ORDER BY m.month;

-- GetYearlyTransferAmountsBySenderCardNumber: Retrieves yearly transfer amounts
-- Purpose: Analyze total transfer amounts over a 5-year period
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: current_year - The final year to include (e.g., 2024), includes 5-year span (current_year - 4)
-- Returns:
--   year: 4-digit year
--   total_transfer_amount: Sum of transfer amounts in that year
-- Business Logic:
--   - Covers a 5-year window (current_year - 4 to current_year)
--   - Filters out soft-deleted records (deleted_at IS NULL)
--   - Groups by calendar year
--   - Useful for identifying long-term money movement trends
-- name: GetYearlyTransferAmountsBySenderCardNumber :many
SELECT EXTRACT(
        YEAR
        FROM t.created_at
    ) AS year, SUM(t.transfer_amount) AS total_transfer_amount
FROM transfers t
WHERE
    t.deleted_at IS NULL
    AND t.transfer_from = $1
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

-- GetYearlyTransferAmountsByReceiverCardNumber: Retrieves yearly transfer amounts
-- Purpose: Analyze total transfer amounts over a 5-year period
-- Parameters:
--   $1: card_number  - filter by card_number
--   $2: current_year - The final year to include (e.g., 2024), includes 5-year span (current_year - 4)
-- Returns:
--   year: 4-digit year
--   total_transfer_amount: Sum of transfer amounts in that year
-- Business Logic:
--   - Covers a 5-year window (current_year - 4 to current_year)
--   - Filters out soft-deleted records (deleted_at IS NULL)
--   - Groups by calendar year
--   - Useful for identifying long-term money movement trends
-- name: GetYearlyTransferAmountsByReceiverCardNumber :many
SELECT EXTRACT(
        YEAR
        FROM t.created_at
    ) AS year, SUM(t.transfer_amount) AS total_transfer_amount
FROM transfers t
WHERE
    t.deleted_at IS NULL
    AND t.transfer_to = $1
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

-- FindAllTransfersByCardNumberAsSender: Retrieves all transfers where the card was the sender
-- Purpose: View outgoing transfer history for a specific card
-- Parameters:
--   $1: card_number - The card number that initiated the transfers
-- Returns:
--   All transfer fields for outgoing transfers (transfer_from = card_number)
-- Business Logic:
--   - Only includes active transfers (non-deleted)
--   - Orders by transfer_time (newest first)
--   - Useful for tracking money sent by a cardholder
-- name: FindAllTransfersByCardNumberAsSender :many
SELECT t.transfer_id, t.transfer_from, t.transfer_to, t.transfer_amount, t.transfer_time, t.created_at, t.updated_at, t.deleted_at
FROM transfers t
WHERE
    t.transfer_from = $1
    AND t.deleted_at IS NULL
ORDER BY t.transfer_time DESC;

-- FindAllTransfersByCardNumberAsReceiver: Retrieves all transfers where the card was the receiver
-- Purpose: View incoming transfer history for a specific card
-- Parameters:
--   $1: card_number - The card number that received the transfers
-- Returns:
--   All transfer fields for incoming transfers (transfer_to = card_number)
-- Business Logic:
--   - Only includes active transfers (non-deleted)
--   - Orders by transfer_time (newest first)
--   - Useful for tracking money received by a cardholder
-- name: FindAllTransfersByCardNumberAsReceiver :many
SELECT t.transfer_id, t.transfer_from, t.transfer_to, t.transfer_amount, t.transfer_time, t.created_at, t.updated_at, t.deleted_at
FROM transfers t
WHERE
    t.transfer_to = $1
    AND t.deleted_at IS NULL
ORDER BY t.transfer_time DESC;

-- CreateTransfer: Records a new transfer transaction
-- Purpose: Create a transfer between accounts/cards
-- Parameters:
--   $1: transfer_from - Source account/card number
--   $2: transfer_to - Destination account/card number
--   $3: transfer_amount - Amount transferred
--   $4: transfer_time - When the transfer occurred
--   $5: status - Initial status of the transfer
-- Returns:
--   The newly created transfer record
-- Business Logic:
--   - Sets creation and update timestamps automatically
--   - Used for recording money movements between accounts
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
    )
RETURNING
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at;

-- CreateTransferAtomic: Debits the sender, credits the receiver, and records a successful transfer atomically.
-- Purpose: Keep the two saldo mutations and transfer record in one PostgreSQL statement.
-- Business Logic:
--   - Locks both active saldo rows before validating the available balance.
--   - Rejects missing accounts, identical sender/receiver accounts, and insufficient funds.
--   - Returns no row for a rejected operation; no balance is changed in that case.
-- name: CreateTransferAtomic :one
WITH active_recipient AS (
    SELECT c.card_number, u.email
    FROM cards c
    JOIN users u ON u.user_id = c.user_id AND u.deleted_at IS NULL
    WHERE c.card_number = sqlc.arg(transfer_from) AND c.deleted_at IS NULL
), locked_saldos AS MATERIALIZED (
    SELECT s.card_number, s.total_balance
    FROM saldos s
    WHERE s.card_number IN (sqlc.arg(transfer_from), sqlc.arg(transfer_to))
      AND s.deleted_at IS NULL
    ORDER BY s.card_number
    FOR UPDATE
), eligible AS (
    SELECT 1
    FROM saldos AS sender_saldo
    JOIN saldos AS receiver_saldo ON receiver_saldo.card_number = sqlc.arg(transfer_to)
    WHERE sender_saldo.card_number = sqlc.arg(transfer_from)
      AND sender_saldo.deleted_at IS NULL
      AND receiver_saldo.deleted_at IS NULL
      AND sender_saldo.total_balance >= sqlc.arg(transfer_amount)
      AND EXISTS (SELECT 1 FROM locked_saldos)
      AND EXISTS (SELECT 1 FROM active_recipient)
      AND sqlc.arg(transfer_amount) > 0
      AND sqlc.arg(transfer_from) <> sqlc.arg(transfer_to)
), settled_saldos AS (
    UPDATE saldos s
    SET total_balance = s.total_balance + CASE
            WHEN s.card_number = sqlc.arg(transfer_from) THEN -sqlc.arg(transfer_amount)
            ELSE sqlc.arg(transfer_amount)
        END,
        updated_at = current_timestamp
    WHERE s.card_number IN (sqlc.arg(transfer_from), sqlc.arg(transfer_to))
      AND s.deleted_at IS NULL
      AND EXISTS (SELECT 1 FROM eligible)
    RETURNING s.card_number
), inserted_transfer AS (
    INSERT INTO transfers (
        transfer_from,
        transfer_to,
        transfer_amount,
        transfer_time,
        idempotency_key,
        status,
        created_at,
        updated_at
    )
    SELECT sqlc.arg(transfer_from), sqlc.arg(transfer_to), sqlc.arg(transfer_amount), current_timestamp, sqlc.arg(idempotency_key), 'success', current_timestamp, current_timestamp
    WHERE (SELECT COUNT(*) FROM settled_saldos) = 2
      AND EXISTS (SELECT 1 FROM active_recipient)
    RETURNING
        transfer_id,
        transfer_no,
        transfer_from,
        transfer_to,
        transfer_amount,
        transfer_time,
        status,
        created_at,
        updated_at
), queued_event AS (
    INSERT INTO outbox_events (event_key, topic, message_key, payload)
    SELECT
        'transfer-email:' || it.transfer_id::text,
        'email-service-topic-transfer-create',
        it.transfer_id::text,
        jsonb_build_object(
            'event_id', gen_random_uuid()::text,
            'schema_version', 1,
            'event_type', 'transfer.created',
            'occurred_at', now(),
            'email', COALESCE(u.email, ''),
            'subject', 'Transfer Successful - SanEdge',
            'body', 'Your transfer of ' || it.transfer_amount::text || ' has been processed successfully.'
        )
    FROM inserted_transfer it
    JOIN active_recipient u ON u.card_number = it.transfer_from
    ON CONFLICT (event_key) DO NOTHING
    RETURNING event_key
)
SELECT
    it.transfer_id,
    it.transfer_no,
    it.transfer_from,
    it.transfer_to,
    it.transfer_amount,
    it.transfer_time,
    it.status,
    it.created_at,
    it.updated_at
FROM inserted_transfer it;

-- GetTransferByIdempotencyKey: Returns an existing transfer for a client-supplied idempotency key.
-- Purpose: Support safe replay of transfer create requests.
-- Parameters:
--   $1: idempotency_key - The client-supplied key
-- Business Logic:
--   - Only matches non-empty keys on active records.
--   - Used by the repository to return the original record instead of
--     double-settling when a request is retried.
-- name: GetTransferByIdempotencyKey :one
SELECT
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at
FROM transfers
WHERE idempotency_key = $1
  AND idempotency_key <> '';

-- UpdateTransferSettlementDelta: Applies sender/receiver balance deltas atomically.
-- Purpose: Race-free balance adjustment for transfer amount updates.
-- Parameters:
--   $1: transfer_from - Sender card number
--   $2: transfer_to - Receiver card number
--   $3: sender_delta - Signed delta for the sender (negative debits)
--   $4: receiver_delta - Signed delta for the receiver
-- Business Logic:
--   - Locks both active saldo rows (ordered) before adjusting.
--   - Rejects when either account is missing or the sender's resulting
--     balance would be negative (returns 0; balance unchanged).
--   - Returns the number of saldo rows actually updated (2 = success).
-- name: UpdateTransferSettlementDelta :one
WITH locked AS MATERIALIZED (
    SELECT s.card_number
    FROM saldos s
    WHERE s.card_number IN ($1, $2)
      AND s.deleted_at IS NULL
    ORDER BY s.card_number
    FOR UPDATE
), eligible AS (
    SELECT 1
    FROM locked
    JOIN saldos sender ON sender.card_number = $1 AND sender.deleted_at IS NULL
    WHERE (SELECT COUNT(*) FROM locked) = 2
      AND sender.total_balance + $3 >= 0
), updated AS (
    UPDATE saldos s
    SET total_balance = s.total_balance + CASE
            WHEN s.card_number = $1 THEN $3
            ELSE $4
        END,
        updated_at = current_timestamp
    WHERE s.card_number IN ($1, $2)
      AND s.deleted_at IS NULL
      AND EXISTS (SELECT 1 FROM eligible)
    RETURNING s.card_number
)
SELECT COUNT(*)::int AS updated_count FROM updated;

-- UpdateTransferAtomic: Updates transfer details and both settlement balances atomically.
-- The transfer row and both saldo rows are locked before the amount delta is
-- computed, so a partial update cannot leave balances and ledger out of sync.
-- name: UpdateTransferAtomic :one
WITH locked_transfer AS MATERIALIZED (
    SELECT transfers.transfer_id, transfers.transfer_from, transfers.transfer_to, transfers.transfer_amount
    FROM transfers
    WHERE transfers.transfer_id = sqlc.arg(transfer_id)
      AND transfers.deleted_at IS NULL
      AND transfers.status IN ('pending', 'failed', 'success', 'compensation_required')
    FOR UPDATE
), locked_saldos AS MATERIALIZED (
    SELECT s.card_number, s.total_balance
    FROM saldos s, locked_transfer t
    WHERE s.card_number IN (t.transfer_from, t.transfer_to)
      AND s.deleted_at IS NULL
    ORDER BY s.card_number
    FOR UPDATE
), eligible AS (
    SELECT 1
    FROM locked_transfer t
    JOIN saldos AS sender_saldo ON sender_saldo.card_number = t.transfer_from
    JOIN saldos AS receiver_saldo ON receiver_saldo.card_number = t.transfer_to
    WHERE t.transfer_from = sqlc.arg(transfer_from)
      AND t.transfer_to = sqlc.arg(transfer_to)
      AND t.transfer_from <> t.transfer_to
      AND sender_saldo.deleted_at IS NULL
      AND receiver_saldo.deleted_at IS NULL
      AND EXISTS (SELECT 1 FROM locked_saldos)
      AND sender_saldo.total_balance + t.transfer_amount - sqlc.arg(transfer_amount) >= 0
      AND sqlc.arg(transfer_amount) > 0
), settled_saldos AS (
    UPDATE saldos s
    SET total_balance = s.total_balance + CASE
        WHEN s.card_number = t.transfer_from THEN t.transfer_amount - sqlc.arg(transfer_amount)
        ELSE sqlc.arg(transfer_amount) - t.transfer_amount
    END,
    updated_at = current_timestamp
    FROM locked_transfer t
    WHERE s.card_number IN (t.transfer_from, t.transfer_to)
      AND s.deleted_at IS NULL
      AND EXISTS (SELECT 1 FROM eligible)
    RETURNING s.card_number
), updated_transfer AS (
    UPDATE transfers t
    SET transfer_from = sqlc.arg(transfer_from),
        transfer_to = sqlc.arg(transfer_to),
        transfer_amount = sqlc.arg(transfer_amount),
        transfer_time = sqlc.arg(transfer_time),
        status = 'success',
        updated_at = current_timestamp
    WHERE t.transfer_id = (SELECT transfer_id FROM locked_transfer)
      AND (SELECT COUNT(*) FROM settled_saldos) = 2
    RETURNING transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at
)
SELECT transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at
FROM updated_transfer;

-- UpdateTransfer: Modifies transfer details
-- Purpose: Update all fields of an existing transfer
-- Parameters:
--   $1: transfer_id - ID of transfer to update
--   $2: transfer_from - Updated source account
--   $3: transfer_to - Updated destination account
--   $4: transfer_amount - Updated amount
--   $5: transfer_time - Updated timestamp
-- Business Logic:
--   - Only updates active transfers (non-deleted)
--   - Updates modification timestamp automatically
--   - Used for correcting transfer details
-- name: UpdateTransfer :one
UPDATE transfers
SET
    transfer_from = $2,
    transfer_to = $3,
    transfer_amount = $4,
    transfer_time = $5,
    updated_at = current_timestamp
WHERE
    transfer_id = $1
    AND deleted_at IS NULL
RETURNING
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at;

-- UpdateTransferAmount: Changes only the transfer amount
-- Purpose: Adjust the amount of a transfer
-- Parameters:
--   $1: transfer_id - ID of transfer to update
--   $2: transfer_amount - New transfer amount
-- Business Logic:
--   - Only updates active transfers
--   - Updates both amount and transfer timestamp
--   - Used for amount corrections
-- name: UpdateTransferAmount :one
UPDATE transfers
SET
    transfer_amount = $2,
    transfer_time = current_timestamp,
    updated_at = current_timestamp
WHERE
    transfer_id = $1
    AND deleted_at IS NULL
RETURNING
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at;

-- UpdateTransferStatus: Changes transfer status
-- Purpose: Update processing status of a transfer
-- Parameters:
--   $1: transfer_id - ID of transfer to update
--   $2: status - New status (e.g., 'completed', 'failed')
-- Business Logic:
--   - Only updates active transfers
--   - Updates modification timestamp
--   - Used to reflect transfer processing outcomes
-- name: UpdateTransferStatus :one
UPDATE transfers
SET
    status = $2,
    updated_at = current_timestamp
WHERE
    transfer_id = $1
    AND deleted_at IS NULL
    AND (
        ($2 = 'failed' AND status IN ('pending', 'failed', 'compensation_required'))
        OR ($2 = 'compensation_required' AND status IN ('pending', 'failed', 'success', 'compensation_required'))
        OR ($2 = 'success' AND status IN ('pending', 'failed', 'success', 'compensation_required'))
        OR ($2 = 'pending' AND status IN ('failed', 'compensation_required'))
    )
RETURNING
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at;

-- TrashTransfer: Soft-deletes a transfer
-- Purpose: Remove transfer from active view without permanent deletion
-- Parameters:
--   $1: transfer_id - ID of transfer to trash
-- Business Logic:
--   - Sets deleted_at timestamp
--   - Only affects active transfers
--   - Preserves data for audit/recovery
-- name: TrashTransfer :one
UPDATE transfers
SET
    deleted_at = current_timestamp
WHERE
    transfer_id = $1
    AND deleted_at IS NULL
RETURNING
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at,
    deleted_at;

-- RestoreTransfer: Recovers a soft-deleted transfer
-- Purpose: Reactivate a previously deleted transfer
-- Parameters:
--   $1: transfer_id - ID of transfer to restore
-- Business Logic:
--   - Clears deleted_at timestamp
--   - Only works on trashed transfers
-- name: RestoreTransfer :one
UPDATE transfers
SET
    deleted_at = NULL
WHERE
    transfer_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    transfer_id,
    transfer_no,
    transfer_from,
    transfer_to,
    transfer_amount,
    transfer_time,
    status,
    created_at,
    updated_at,
    deleted_at;

-- DeleteTransferPermanently: Hard-deletes a transfer
-- Purpose: Permanently remove a transfer record
-- Parameters:
--   $1: transfer_id - ID of transfer to delete
-- Business Logic:
--   - Physical deletion from database
--   - Only works on already trashed transfers
--   - Irreversible operation
-- name: DeleteTransferPermanently :exec
DELETE FROM transfers
WHERE
    transfer_id = $1
    AND deleted_at IS NOT NULL;

-- RestoreAllTransfers: Recovers all trashed transfers
-- Purpose: Mass restoration of deleted transfers
-- Business Logic:
--   - Clears deleted_at for all trashed transfers
--   - Useful for system recovery
--   - Use with caution in production
-- name: RestoreAllTransfers :exec
UPDATE transfers
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentTransfers: Permanently removes all trashed transfers
-- Purpose: Clean up all soft-deleted transfers
-- Business Logic:
--   - Irreversible bulk deletion
--   - Only affects trashed records
--   - Frees database space
-- name: DeleteAllPermanentTransfers :exec
DELETE FROM transfers WHERE deleted_at IS NOT NULL;