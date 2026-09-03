package transactionbycardrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
)

type transactionStatsByCardStatusRepository struct {
	db *gorm.DB
}

func NewTransactionStatsByCardStatusRepository(db *gorm.DB) TransactionStatsByCardStatusRepository {
	return &transactionStatsByCardStatusRepository{db: db}
}

func (r *transactionStatsByCardStatusRepository) GetMonthTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*models.TransactionMonthlyStatusSuccessByCardRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.TransactionMonthlyStatusSuccessByCardRow

	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year,
				EXTRACT(MONTH FROM t.transaction_time)::integer AS month,
				COUNT(*) AS total_success,
				COALESCE(SUM(t.amount), 0)::integer AS total_amount
			FROM transactions t
			WHERE t.deleted_at IS NULL
				AND t.status = 'success'
				AND t.card_number = ?
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
	`, req.CardNumber, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth,
		currentDate, currentDate, currentDate, currentDate,
		prevDate, prevDate, prevDate, prevDate).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusSuccessByCardFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionStatsByCardStatusRepository) GetYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*models.TransactionYearlyStatusSuccessByCardRow, error) {
	var results []*models.TransactionYearlyStatusSuccessByCardRow

	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year,
				COUNT(*) AS total_success,
				COALESCE(SUM(t.amount), 0)::integer AS total_amount
			FROM transactions t
			WHERE t.deleted_at IS NULL
				AND t.status = 'success'
				AND t.card_number = ?
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
	`, req.CardNumber, req.Year, req.Year, req.Year, req.Year, req.Year, req.Year).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusSuccessByCardFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionStatsByCardStatusRepository) GetMonthTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*models.TransactionMonthlyStatusFailedByCardRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.TransactionMonthlyStatusFailedByCardRow

	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year,
				EXTRACT(MONTH FROM t.transaction_time)::integer AS month,
				COUNT(*) AS total_failed,
				COALESCE(SUM(t.amount), 0)::integer AS total_amount
			FROM transactions t
			WHERE t.deleted_at IS NULL
				AND t.status = 'failed'
				AND t.card_number = ?
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
	`, req.CardNumber, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth,
		currentDate, currentDate, currentDate, currentDate,
		prevDate, prevDate, prevDate, prevDate).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetMonthTransactionStatusFailedByCardFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionStatsByCardStatusRepository) GetYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*models.TransactionYearlyStatusFailedByCardRow, error) {
	var results []*models.TransactionYearlyStatusFailedByCardRow

	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year,
				COUNT(*) AS total_failed,
				COALESCE(SUM(t.amount), 0)::integer AS total_amount
			FROM transactions t
			WHERE t.deleted_at IS NULL
				AND t.status = 'failed'
				AND t.card_number = ?
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
	`, req.CardNumber, req.Year, req.Year, req.Year, req.Year, req.Year, req.Year).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionStatusFailedByCardFailed.WithInternal(err)
	}

	return results, nil
}
