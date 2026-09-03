package transactionstatsrepository

import (
	"context"
	"strconv"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
)

type transactionStatsStatusRepository struct {
	db *gorm.DB
}

func NewTransactionStatsStatusRepository(db *gorm.DB) TransactionStatsStatusRepository {
	return &transactionStatsStatusRepository{db: db}
}

func (r *transactionStatsStatusRepository) GetMonthTransactionStatusSuccess(ctx context.Context, req *requests.MonthStatusTransaction) ([]*models.TransactionMonthlyStatusSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.TransactionMonthlyStatusSuccessRow

	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year,
				EXTRACT(MONTH FROM t.transaction_time)::integer AS month,
				COUNT(*) AS total_success,
				COALESCE(SUM(t.amount), 0)::integer AS total_amount
			FROM transactions t
			WHERE t.deleted_at IS NULL
				AND t.status = 'success'
				AND ((t.transaction_time >= ?::timestamp AND t.transaction_time <= ?::timestamp)
					OR (t.transaction_time >= ?::timestamp AND t.transaction_time <= ?::timestamp))
			GROUP BY EXTRACT(YEAR FROM t.transaction_time), EXTRACT(MONTH FROM t.transaction_time)
		), formatted_data AS (
			SELECT year::text, TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
				total_success, total_amount FROM monthly_data
			UNION ALL
			SELECT EXTRACT(YEAR FROM ?::timestamp)::text AS year,
				TO_CHAR(?::timestamp, 'Mon') AS month, 0 AS total_success, 0 AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM monthly_data
				WHERE year = EXTRACT(YEAR FROM ?::timestamp)::text
				AND month = EXTRACT(MONTH FROM ?::timestamp)::integer)
			UNION ALL
			SELECT EXTRACT(YEAR FROM ?::timestamp)::text AS year,
				TO_CHAR(?::timestamp, 'Mon') AS month, 0 AS total_success, 0 AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM monthly_data
				WHERE year = EXTRACT(YEAR FROM ?::timestamp)::text
				AND month = EXTRACT(MONTH FROM ?::timestamp)::integer)
		)
		SELECT * FROM formatted_data ORDER BY year DESC, TO_DATE(month, 'Mon') DESC
	`, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth,
		currentDate, currentDate, currentDate, currentDate,
		prevDate, prevDate, prevDate, prevDate).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusSuccessFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionStatsStatusRepository) GetYearlyTransactionStatusSuccess(ctx context.Context, year int) ([]*models.TransactionYearlyStatusSuccessRow, error) {
	var results []*models.TransactionYearlyStatusSuccessRow

	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year,
				COUNT(*) AS total_success,
				COALESCE(SUM(t.amount), 0)::integer AS total_amount
			FROM transactions t
			WHERE t.deleted_at IS NULL
				AND t.status = 'success'
				AND (EXTRACT(YEAR FROM t.transaction_time)::text = ?
					OR EXTRACT(YEAR FROM t.transaction_time)::text = (? - 1)::text)
			GROUP BY EXTRACT(YEAR FROM t.transaction_time)
		), formatted_data AS (
			SELECT year::text, total_success::integer, total_amount::integer FROM yearly_data
			UNION ALL
			SELECT ? AS year, 0::integer AS total_success, 0::integer AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM yearly_data WHERE year::integer = ?)
			UNION ALL
			SELECT (? - 1)::text AS year, 0::integer AS total_success, 0::integer AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM yearly_data WHERE year::integer = ? - 1)
		)
		SELECT * FROM formatted_data ORDER BY year DESC
	`, strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year)).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusSuccessFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionStatsStatusRepository) GetMonthTransactionStatusFailed(ctx context.Context, req *requests.MonthStatusTransaction) ([]*models.TransactionMonthlyStatusFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.TransactionMonthlyStatusFailedRow

	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year,
				EXTRACT(MONTH FROM t.transaction_time)::integer AS month,
				COUNT(*) AS total_failed,
				COALESCE(SUM(t.amount), 0)::integer AS total_amount
			FROM transactions t
			WHERE t.deleted_at IS NULL
				AND t.status = 'failed'
				AND ((t.transaction_time >= ?::timestamp AND t.transaction_time <= ?::timestamp)
					OR (t.transaction_time >= ?::timestamp AND t.transaction_time <= ?::timestamp))
			GROUP BY EXTRACT(YEAR FROM t.transaction_time), EXTRACT(MONTH FROM t.transaction_time)
		), formatted_data AS (
			SELECT year::text, TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
				total_failed, total_amount FROM monthly_data
			UNION ALL
			SELECT EXTRACT(YEAR FROM ?::timestamp)::text AS year,
				TO_CHAR(?::timestamp, 'Mon') AS month, 0 AS total_failed, 0 AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM monthly_data
				WHERE year = EXTRACT(YEAR FROM ?::timestamp)::text
				AND month = EXTRACT(MONTH FROM ?::timestamp)::integer)
			UNION ALL
			SELECT EXTRACT(YEAR FROM ?::timestamp)::text AS year,
				TO_CHAR(?::timestamp, 'Mon') AS month, 0 AS total_failed, 0 AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM monthly_data
				WHERE year = EXTRACT(YEAR FROM ?::timestamp)::text
				AND month = EXTRACT(MONTH FROM ?::timestamp)::integer)
		)
		SELECT * FROM formatted_data ORDER BY year DESC, TO_DATE(month, 'Mon') DESC
	`, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth,
		currentDate, currentDate, currentDate, currentDate,
		prevDate, prevDate, prevDate, prevDate).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusFailedFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionStatsStatusRepository) GetYearlyTransactionStatusFailed(ctx context.Context, year int) ([]*models.TransactionYearlyStatusFailedRow, error) {
	var results []*models.TransactionYearlyStatusFailedRow

	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year,
				COUNT(*) AS total_failed,
				COALESCE(SUM(t.amount), 0)::integer AS total_amount
			FROM transactions t
			WHERE t.deleted_at IS NULL
				AND t.status = 'failed'
				AND (EXTRACT(YEAR FROM t.transaction_time)::text = ?
					OR EXTRACT(YEAR FROM t.transaction_time)::text = (? - 1)::text)
			GROUP BY EXTRACT(YEAR FROM t.transaction_time)
		), formatted_data AS (
			SELECT year::text, total_failed::integer, total_amount::integer FROM yearly_data
			UNION ALL
			SELECT ? AS year, 0::integer AS total_failed, 0::integer AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM yearly_data WHERE year::integer = ?)
			UNION ALL
			SELECT (? - 1)::text AS year, 0::integer AS total_failed, 0::integer AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM yearly_data WHERE year::integer = ? - 1)
		)
		SELECT * FROM formatted_data ORDER BY year DESC
	`, strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year)).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusFailedFailed.WithInternal(err)
	}

	return results, nil
}
