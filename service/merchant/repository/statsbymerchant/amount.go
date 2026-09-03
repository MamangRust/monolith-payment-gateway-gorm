package merchantstatsmerchantrepository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	"gorm.io/gorm"
)

type merchantStatsAmountByMerchantRepository struct {
	db *gorm.DB
}

func NewMerchantStatsAmountByMerchantRepository(db *gorm.DB) MerchantStatsAmountByMerchantRepository {
	return &merchantStatsAmountByMerchantRepository{db: db}
}

func (r *merchantStatsAmountByMerchantRepository) GetMonthlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*models.MerchantMonthlyAmountRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
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
			AND t.merchant_id = ?
			AND t.deleted_at IS NULL
		GROUP BY m.month
		ORDER BY m.month
	`, yearStart, yearStart, req.MerchantID).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountByMerchantsFailed.WithInternal(err)
	}
	return res, nil
}

func (r *merchantStatsAmountByMerchantRepository) GetYearlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*models.MerchantYearlyAmountRow, error) {
	var res []*models.MerchantYearlyAmountRow
	err := r.db.WithContext(ctx).Raw(`
		WITH years AS (
			SELECT generate_series(?::int - 4, ?::int) AS year
		),
		yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year, SUM(t.amount) AS total_amount
			FROM transactions t
			WHERE t.merchant_id = ? AND t.deleted_at IS NULL
			GROUP BY EXTRACT(YEAR FROM t.transaction_time)
		)
		SELECT y.year::text AS year, COALESCE(yd.total_amount, 0)::bigint AS total_amount
		FROM years y
		LEFT JOIN yearly_data yd ON y.year = yd.year::int
		ORDER BY y.year
	`, req.Year, req.MerchantID).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountByMerchantsFailed.WithInternal(err)
	}
	return res, nil
}
