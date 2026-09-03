package merchantstatsmerchantrepository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	"gorm.io/gorm"
)

type merchantStatsMethodByMerchantRepository struct {
	db *gorm.DB
}

func NewMerchantStatsMethodByMerchantRepository(db *gorm.DB) MerchantStatsMethodByMerchantRepository {
	return &merchantStatsMethodByMerchantRepository{db: db}
}

func (r *merchantStatsMethodByMerchantRepository) GetMonthlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*models.MerchantMonthlyPaymentMethodRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MerchantMonthlyPaymentMethodRow
	err := r.db.WithContext(ctx).Raw(`
		WITH months AS (
			SELECT generate_series(
				date_trunc('year', ?::timestamp),
				date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day',
				interval '1 month'
			) AS month
		),
		payment_methods AS (
			SELECT DISTINCT payment_method
			FROM transactions WHERE merchant_id = ? AND transaction_time >= ? AND transaction_time < ? AND deleted_at IS NULL
		)
		SELECT TO_CHAR(m.month, 'Mon') AS month, pm.payment_method, COALESCE(SUM(t.amount), 0)::int AS total_amount
		FROM months m
		CROSS JOIN payment_methods pm
		LEFT JOIN transactions t ON EXTRACT(MONTH FROM t.transaction_time) = EXTRACT(MONTH FROM m.month)
			AND EXTRACT(YEAR FROM t.transaction_time) = EXTRACT(YEAR FROM m.month)
			AND t.merchant_id = ?
			AND t.payment_method = pm.payment_method
			AND t.deleted_at IS NULL
		GROUP BY m.month, pm.payment_method
		ORDER BY m.month, pm.payment_method
	`, yearStart, yearStart, req.MerchantID, yearStart, yearEnd, req.MerchantID).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodByMerchantsFailed.WithInternal(err)
	}
	return res, nil
}

func (r *merchantStatsMethodByMerchantRepository) GetYearlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*models.MerchantYearlyPaymentMethodRow, error) {
	yearStart := time.Date(req.Year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MerchantYearlyPaymentMethodRow
	err := r.db.WithContext(ctx).Raw(`
		WITH years AS (
			SELECT generate_series(?::int - 4, ?::int) AS year
		),
		yearly_methods AS (
			SELECT DISTINCT payment_method
			FROM transactions WHERE merchant_id = ? AND transaction_time >= ? AND transaction_time < ? AND deleted_at IS NULL
		),
		yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year, t.payment_method, SUM(t.amount) AS total_amount
			FROM transactions t
			WHERE t.merchant_id = ? AND t.deleted_at IS NULL
			GROUP BY EXTRACT(YEAR FROM t.transaction_time), t.payment_method
		)
		SELECT y.year::text AS year, ym.payment_method, COALESCE(yd.total_amount, 0)::bigint AS total_amount
		FROM years y
		CROSS JOIN yearly_methods ym
		LEFT JOIN yearly_data yd ON y.year = yd.year::int AND yd.payment_method = ym.payment_method
		ORDER BY y.year, ym.payment_method
	`, req.Year, req.Year, req.MerchantID, yearStart, yearEnd, req.MerchantID).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodByMerchantsFailed.WithInternal(err)
	}
	return res, nil
}
