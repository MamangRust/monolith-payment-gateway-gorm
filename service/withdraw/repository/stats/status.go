package withdrawstatsrepository

import (
	"context"
	"strconv"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
)

type withdrawStatsStatusRepository struct {
	db *gorm.DB
}

func NewWithdrawStatsStatusRepository(db *gorm.DB) WithdrawStatsStatusRepository {
	return &withdrawStatsStatusRepository{db: db}
}

func (r *withdrawStatsStatusRepository) GetMonthWithdrawStatusSuccess(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*models.WithdrawMonthlyStatusSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.WithdrawMonthlyStatusSuccessRow
	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.withdraw_time)::text AS year, EXTRACT(MONTH FROM t.withdraw_time)::integer AS month,
				COUNT(*) AS total_success, COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
			FROM withdraws t WHERE t.deleted_at IS NULL AND t.status = 'success'
				AND ((t.withdraw_time >= ?::timestamp AND t.withdraw_time <= ?::timestamp)
					OR (t.withdraw_time >= ?::timestamp AND t.withdraw_time <= ?::timestamp))
			GROUP BY EXTRACT(YEAR FROM t.withdraw_time), EXTRACT(MONTH FROM t.withdraw_time)
		), formatted_data AS (
			SELECT year::text, TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month, total_success, total_amount FROM monthly_data
			UNION ALL
			SELECT EXTRACT(YEAR FROM ?::timestamp)::text AS year, TO_CHAR(?::timestamp, 'Mon') AS month, 0 AS total_success, 0 AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM monthly_data WHERE year = EXTRACT(YEAR FROM ?::timestamp)::text AND month = EXTRACT(MONTH FROM ?::timestamp)::integer)
			UNION ALL
			SELECT EXTRACT(YEAR FROM ?::timestamp)::text AS year, TO_CHAR(?::timestamp, 'Mon') AS month, 0 AS total_success, 0 AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM monthly_data WHERE year = EXTRACT(YEAR FROM ?::timestamp)::text AND month = EXTRACT(MONTH FROM ?::timestamp)::integer)
		)
		SELECT * FROM formatted_data ORDER BY year DESC, TO_DATE(month, 'Mon') DESC
	`, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth,
		currentDate, currentDate, currentDate, currentDate,
		prevDate, prevDate, prevDate, prevDate).Scan(&results).Error
	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusSuccessFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawStatsStatusRepository) GetYearlyWithdrawStatusSuccess(ctx context.Context, year int) ([]*models.WithdrawYearlyStatusSuccessRow, error) {
	var results []*models.WithdrawYearlyStatusSuccessRow
	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.withdraw_time)::text AS year, COUNT(*) AS total_success, COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
			FROM withdraws t WHERE t.deleted_at IS NULL AND t.status = 'success'
				AND (EXTRACT(YEAR FROM t.withdraw_time) = ? OR EXTRACT(YEAR FROM t.withdraw_time) = ? - 1)
			GROUP BY EXTRACT(YEAR FROM t.withdraw_time)
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
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusSuccessFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawStatsStatusRepository) GetMonthWithdrawStatusFailed(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*models.WithdrawMonthlyStatusFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.WithdrawMonthlyStatusFailedRow
	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.withdraw_time)::text AS year, EXTRACT(MONTH FROM t.withdraw_time)::integer AS month,
				COUNT(*) AS total_failed, COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
			FROM withdraws t WHERE t.deleted_at IS NULL AND t.status = 'failed'
				AND ((t.withdraw_time >= ?::timestamp AND t.withdraw_time <= ?::timestamp)
					OR (t.withdraw_time >= ?::timestamp AND t.withdraw_time <= ?::timestamp))
			GROUP BY EXTRACT(YEAR FROM t.withdraw_time), EXTRACT(MONTH FROM t.withdraw_time)
		), formatted_data AS (
			SELECT year::text, TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month, total_failed, total_amount FROM monthly_data
			UNION ALL
			SELECT EXTRACT(YEAR FROM ?::timestamp)::text AS year, TO_CHAR(?::timestamp, 'Mon') AS month, 0 AS total_failed, 0 AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM monthly_data WHERE year = EXTRACT(YEAR FROM ?::timestamp)::text AND month = EXTRACT(MONTH FROM ?::timestamp)::integer)
			UNION ALL
			SELECT EXTRACT(YEAR FROM ?::timestamp)::text AS year, TO_CHAR(?::timestamp, 'Mon') AS month, 0 AS total_failed, 0 AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM monthly_data WHERE year = EXTRACT(YEAR FROM ?::timestamp)::text AND month = EXTRACT(MONTH FROM ?::timestamp)::integer)
		)
		SELECT * FROM formatted_data ORDER BY year DESC, TO_DATE(month, 'Mon') DESC
	`, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth,
		currentDate, currentDate, currentDate, currentDate,
		prevDate, prevDate, prevDate, prevDate).Scan(&results).Error
	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusFailedFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawStatsStatusRepository) GetYearlyWithdrawStatusFailed(ctx context.Context, year int) ([]*models.WithdrawYearlyStatusFailedRow, error) {
	var results []*models.WithdrawYearlyStatusFailedRow
	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.withdraw_time)::text AS year, COUNT(*) AS total_failed, COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
			FROM withdraws t WHERE t.deleted_at IS NULL AND t.status = 'failed'
				AND (EXTRACT(YEAR FROM t.withdraw_time) = ? OR EXTRACT(YEAR FROM t.withdraw_time) = ? - 1)
			GROUP BY EXTRACT(YEAR FROM t.withdraw_time)
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
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusFailedFailed.WithInternal(err)
	}
	return results, nil
}
