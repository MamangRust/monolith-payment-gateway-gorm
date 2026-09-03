package withdrawstatsrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
)

type withdrawStatsAmountRepository struct {
	db *gorm.DB
}

func NewWithdrawStatsAmountRepository(db *gorm.DB) WithdrawStatsAmountRepository {
	return &withdrawStatsAmountRepository{db: db}
}

func (r *withdrawStatsAmountRepository) GetMonthlyWithdraws(ctx context.Context, year int) ([]*models.WithdrawMonthlyAmountRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var results []*models.WithdrawMonthlyAmountRow
	err := r.db.WithContext(ctx).Raw(`
		WITH months AS (
			SELECT generate_series(date_trunc('year', ?::timestamp), date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day', interval '1 month') AS month
		)
		SELECT TO_CHAR(m.month, 'Mon') AS month, COALESCE(SUM(w.withdraw_amount), 0)::int AS total_withdraw_amount
		FROM months m LEFT JOIN withdraws w ON EXTRACT(MONTH FROM w.withdraw_time) = EXTRACT(MONTH FROM m.month)
			AND EXTRACT(YEAR FROM w.withdraw_time) = EXTRACT(YEAR FROM m.month) AND w.deleted_at IS NULL
		GROUP BY m.month ORDER BY m.month
	`, yearStart, yearEnd).Scan(&results).Error
	if err != nil {
		return nil, withdraw_errors.ErrGetMonthlyWithdrawsFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawStatsAmountRepository) GetYearlyWithdraws(ctx context.Context, year int) ([]*models.WithdrawYearlyAmountRow, error) {
	var results []*models.WithdrawYearlyAmountRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM w.withdraw_time)::text AS year, SUM(w.withdraw_amount) AS total_withdraw_amount
		FROM withdraws w WHERE w.deleted_at IS NULL
			AND EXTRACT(YEAR FROM w.withdraw_time) >= ? - 4 AND EXTRACT(YEAR FROM w.withdraw_time) <= ?
		GROUP BY EXTRACT(YEAR FROM w.withdraw_time) ORDER BY year
	`, year, year).Scan(&results).Error
	if err != nil {
		return nil, withdraw_errors.ErrGetYearlyWithdrawsFailed.WithInternal(err)
	}
	return results, nil
}
