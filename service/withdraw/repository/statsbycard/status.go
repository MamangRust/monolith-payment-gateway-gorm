package withdrawstatsbycardrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
)

type withdrawStatsByCardStatusRepository struct {
	db *gorm.DB
}

type withdrawStatsByCardAmountRepository struct {
	db *gorm.DB
}

type repository struct {
	WithdrawStatsByCardStatusRepository
	WithdrawStatsByCardAmountRepository
}

func NewWithdrawStatsByCardRepository(db *gorm.DB) WithdrawStatsByCardRepository {
	return &repository{
		WithdrawStatsByCardStatusRepository: NewWithdrawStatsByCardStatusRepository(db),
		WithdrawStatsByCardAmountRepository: NewWithdrawStatsByCardAmountRepository(db),
	}
}

func NewWithdrawStatsByCardStatusRepository(db *gorm.DB) WithdrawStatsByCardStatusRepository {
	return &withdrawStatsByCardStatusRepository{db: db}
}

func NewWithdrawStatsByCardAmountRepository(db *gorm.DB) WithdrawStatsByCardAmountRepository {
	return &withdrawStatsByCardAmountRepository{db: db}
}

func (r *withdrawStatsByCardAmountRepository) GetMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*models.WithdrawMonthlyAmountByCardRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	var results []*models.WithdrawMonthlyAmountByCardRow
	err := r.db.WithContext(ctx).Raw(`
		WITH months AS (
			SELECT generate_series(date_trunc('year', ?::timestamp), date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day', interval '1 month') AS month
		)
		SELECT TO_CHAR(m.month, 'Mon') AS month, COALESCE(SUM(w.withdraw_amount), 0)::int AS total_withdraw_amount
		FROM months m LEFT JOIN withdraws w ON EXTRACT(MONTH FROM w.withdraw_time) = EXTRACT(MONTH FROM m.month)
			AND EXTRACT(YEAR FROM w.withdraw_time) = EXTRACT(YEAR FROM m.month) AND w.card_number = ? AND w.deleted_at IS NULL
		GROUP BY m.month ORDER BY m.month
	`, yearStart, yearStart, req.CardNumber).Scan(&results).Error
	if err != nil {
		return nil, withdraw_errors.ErrGetMonthlyWithdrawsByCardFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawStatsByCardAmountRepository) GetYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*models.WithdrawYearlyAmountByCardRow, error) {
	var results []*models.WithdrawYearlyAmountByCardRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM w.withdraw_time)::text AS year, SUM(w.withdraw_amount) AS total_withdraw_amount
		FROM withdraws w WHERE w.deleted_at IS NULL AND w.card_number = ?
			AND EXTRACT(YEAR FROM w.withdraw_time) >= ? - 4 AND EXTRACT(YEAR FROM w.withdraw_time) <= ?
		GROUP BY EXTRACT(YEAR FROM w.withdraw_time) ORDER BY year
	`, req.CardNumber, req.Year, req.Year).Scan(&results).Error
	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawsByCardFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawStatsByCardStatusRepository) GetMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*models.WithdrawMonthlyStatusSuccessByCardRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.WithdrawMonthlyStatusSuccessByCardRow
	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.withdraw_time)::text AS year, EXTRACT(MONTH FROM t.withdraw_time)::integer AS month,
				COUNT(*) AS total_success, COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
			FROM withdraws t WHERE t.deleted_at IS NULL AND t.status = 'success' AND t.card_number = ?
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
	`, req.CardNumber, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth,
		currentDate, currentDate, currentDate, currentDate,
		prevDate, prevDate, prevDate, prevDate).Scan(&results).Error
	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusSuccessByCardFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawStatsByCardStatusRepository) GetYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*models.WithdrawYearlyStatusSuccessByCardRow, error) {
	var results []*models.WithdrawYearlyStatusSuccessByCardRow
	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.withdraw_time)::text AS year, COUNT(*) AS total_success, COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
			FROM withdraws t WHERE t.deleted_at IS NULL AND t.status = 'success' AND t.card_number = ?
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
	`, req.CardNumber, req.Year, req.Year, req.Year, req.Year, req.Year, req.Year).Scan(&results).Error
	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusSuccessByCardFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawStatsByCardStatusRepository) GetMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*models.WithdrawMonthlyStatusFailedByCardRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.WithdrawMonthlyStatusFailedByCardRow
	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.withdraw_time)::text AS year, EXTRACT(MONTH FROM t.withdraw_time)::integer AS month,
				COUNT(*) AS total_failed, COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
			FROM withdraws t WHERE t.deleted_at IS NULL AND t.status = 'failed' AND t.card_number = ?
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
	`, req.CardNumber, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth,
		currentDate, currentDate, currentDate, currentDate,
		prevDate, prevDate, prevDate, prevDate).Scan(&results).Error
	if err != nil {
		return nil, withdraw_errors.ErrGetMonthWithdrawStatusFailedByCardFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawStatsByCardStatusRepository) GetYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*models.WithdrawYearlyStatusFailedByCardRow, error) {
	var results []*models.WithdrawYearlyStatusFailedByCardRow
	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.withdraw_time)::text AS year, COUNT(*) AS total_failed, COALESCE(SUM(t.withdraw_amount), 0)::integer AS total_amount
			FROM withdraws t WHERE t.deleted_at IS NULL AND t.status = 'failed' AND t.card_number = ?
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
	`, req.CardNumber, req.Year, req.Year, req.Year, req.Year, req.Year, req.Year).Scan(&results).Error
	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawStatusFailedByCardFailed.WithInternal(err)
	}
	return results, nil
}
