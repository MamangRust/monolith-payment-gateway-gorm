package topupstatsbycardrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
)

type topupStatsByCardStatusRepository struct {
	db *gorm.DB
}

func NewTopupStatsByCardStatusRepository(db *gorm.DB) TopupStatsByCardStatusRepository {
	return &topupStatsByCardStatusRepository{db: db}
}

func (r *topupStatsByCardStatusRepository) GetMonthTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*models.TopupMonthlyStatusRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.TopupMonthlyStatusRow

	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.topup_time)::text AS year,
				EXTRACT(MONTH FROM t.topup_time)::integer AS month,
				COUNT(*) AS total_success,
				COALESCE(SUM(t.topup_amount), 0)::integer AS total_amount
			FROM topups t
			WHERE t.deleted_at IS NULL
				AND t.status = 'success'
				AND t.card_number = ?
				AND ((t.topup_time >= ?::timestamp AND t.topup_time <= ?::timestamp)
					OR (t.topup_time >= ?::timestamp AND t.topup_time <= ?::timestamp))
			GROUP BY EXTRACT(YEAR FROM t.topup_time), EXTRACT(MONTH FROM t.topup_time)
		), formatted_data AS (
			SELECT year::text, TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
				total_success, total_amount
			FROM monthly_data
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
		return nil, topup_errors.ErrGetMonthTopupStatusSuccessByCardFailed.WithInternal(err)
	}

	return results, nil
}

func (r *topupStatsByCardStatusRepository) GetYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*models.TopupYearlyStatusRow, error) {
	var results []*models.TopupYearlyStatusRow

	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.topup_time)::text AS year,
				COUNT(*) AS total_success,
				COALESCE(SUM(t.topup_amount), 0)::integer AS total_amount
			FROM topups t
			WHERE t.deleted_at IS NULL
				AND t.status = 'success'
				AND t.card_number = ?
				AND (EXTRACT(YEAR FROM t.topup_time) = ?
					OR EXTRACT(YEAR FROM t.topup_time) = ? - 1)
			GROUP BY EXTRACT(YEAR FROM t.topup_time)
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
		return nil, topup_errors.ErrGetYearlyTopupStatusSuccessByCardFailed.WithInternal(err)
	}

	return results, nil
}

func (r *topupStatsByCardStatusRepository) GetMonthTopupStatusFailedByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*models.TopupMonthlyStatusRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.TopupMonthlyStatusRow

	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.topup_time)::text AS year,
				EXTRACT(MONTH FROM t.topup_time)::integer AS month,
				COUNT(*) AS total_failed,
				COALESCE(SUM(t.topup_amount), 0)::integer AS total_amount
			FROM topups t
			WHERE t.deleted_at IS NULL
				AND t.status = 'failed'
				AND t.card_number = ?
				AND ((t.topup_time >= ?::timestamp AND t.topup_time <= ?::timestamp)
					OR (t.topup_time >= ?::timestamp AND t.topup_time <= ?::timestamp))
			GROUP BY EXTRACT(YEAR FROM t.topup_time), EXTRACT(MONTH FROM t.topup_time)
		), formatted_data AS (
			SELECT year::text, TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
				total_failed, total_amount
			FROM monthly_data
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
		return nil, topup_errors.ErrGetMonthTopupStatusFailedByCardFailed.WithInternal(err)
	}

	return results, nil
}

func (r *topupStatsByCardStatusRepository) GetYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*models.TopupYearlyStatusRow, error) {
	var results []*models.TopupYearlyStatusRow

	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.topup_time)::text AS year,
				COUNT(*) AS total_failed,
				COALESCE(SUM(t.topup_amount), 0)::integer AS total_amount
			FROM topups t
			WHERE t.deleted_at IS NULL
				AND t.status = 'failed'
				AND t.card_number = ?
				AND (EXTRACT(YEAR FROM t.topup_time) = ?
					OR EXTRACT(YEAR FROM t.topup_time) = ? - 1)
			GROUP BY EXTRACT(YEAR FROM t.topup_time)
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
		return nil, topup_errors.ErrGetYearlyTopupStatusSuccessByCardFailed.WithInternal(err)
	}

	return results, nil
}
