package merchantstatsrepository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	"gorm.io/gorm"
)

type merchantStatsAmountRepository struct {
	db *gorm.DB
}

func NewMerchantStatsAmountRepository(db *gorm.DB) MerchantStatsAmountRepository {
	return &merchantStatsAmountRepository{db: db}
}

func (r *merchantStatsAmountRepository) GetMonthlyAmountMerchant(ctx context.Context, year int) ([]*models.MerchantMonthlyAmountRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MerchantMonthlyAmountRow
	err := r.db.WithContext(ctx).Raw(`
		WITH months AS (
			SELECT generate_series(
				date_trunc('year', ?::timestamp),
				date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day',
				interval '1 month'
			) AS month
		)
		SELECT TO_CHAR(m.month, 'Mon') AS month, COALESCE(SUM(t.amount), 0)::int AS total_amount
		FROM months m
		LEFT JOIN transactions t ON EXTRACT(MONTH FROM t.transaction_time) = EXTRACT(MONTH FROM m.month)
			AND EXTRACT(YEAR FROM t.transaction_time) = EXTRACT(YEAR FROM m.month)
			AND t.deleted_at IS NULL
		GROUP BY m.month
		ORDER BY m.month
	`, yearStart, yearStart).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountMerchantFailed.WithInternal(err)
	}
	return res, nil
}

func (r *merchantStatsAmountRepository) GetYearlyAmountMerchant(ctx context.Context, year int) ([]*models.MerchantYearlyAmountRow, error) {
	var res []*models.MerchantYearlyAmountRow
	err := r.db.WithContext(ctx).Raw(`
		WITH years AS (
			SELECT generate_series(?::int - 4, ?::int) AS year
		),
		yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year, SUM(t.amount) AS total_amount
			FROM transactions t
			WHERE t.deleted_at IS NULL
			GROUP BY EXTRACT(YEAR FROM t.transaction_time)
		)
		SELECT y.year::text AS year, COALESCE(yd.total_amount, 0)::bigint AS total_amount
		FROM years y
		LEFT JOIN yearly_data yd ON y.year = yd.year::int
		ORDER BY y.year
	`, year, year).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountMerchantFailed.WithInternal(err)
	}
	return res, nil
}
