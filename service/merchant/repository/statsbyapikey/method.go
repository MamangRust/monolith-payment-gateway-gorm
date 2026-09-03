package merchantstatsapikeyrepository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	"gorm.io/gorm"
)

type merchantStatsMethodByApiKeyRepository struct {
	db *gorm.DB
}

func NewMerchantStatsMethodByApiKeyRepository(db *gorm.DB) MerchantStatsMethodByApiKeyRepository {
	return &merchantStatsMethodByApiKeyRepository{db: db}
}

func (r *merchantStatsMethodByApiKeyRepository) GetMonthlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*models.MerchantMonthlyPaymentMethodRow, error) {
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
			SELECT DISTINCT tx.payment_method
			FROM transactions tx
			JOIN merchants mer ON mer.merchant_id = tx.merchant_id AND mer.deleted_at IS NULL
			WHERE mer.api_key = ? AND tx.transaction_time >= ? AND tx.transaction_time < ? AND tx.deleted_at IS NULL
		),
		apikey_data AS (
			SELECT EXTRACT(MONTH FROM tx.transaction_time) AS month_num, tx.payment_method, SUM(tx.amount) AS total_amount
			FROM transactions tx
			JOIN merchants mer ON mer.merchant_id = tx.merchant_id AND mer.deleted_at IS NULL
			WHERE mer.api_key = ? AND tx.transaction_time >= ? AND tx.transaction_time < ? AND tx.deleted_at IS NULL
			GROUP BY EXTRACT(MONTH FROM tx.transaction_time), tx.payment_method
		)
		SELECT TO_CHAR(mo.month, 'Mon') AS month, pm.payment_method, COALESCE(ad.total_amount, 0)::int AS total_amount
		FROM months mo
		CROSS JOIN payment_methods pm
		LEFT JOIN apikey_data ad ON EXTRACT(MONTH FROM mo.month) = ad.month_num AND ad.payment_method = pm.payment_method
		ORDER BY mo.month, pm.payment_method
	`, yearStart, req.Apikey, yearStart, yearEnd, req.Apikey, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodByApikeyFailed.WithInternal(err)
	}
	return res, nil
}

func (r *merchantStatsMethodByApiKeyRepository) GetYearlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*models.MerchantYearlyPaymentMethodRow, error) {
	yearStart := time.Date(req.Year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MerchantYearlyPaymentMethodRow
	err := r.db.WithContext(ctx).Raw(`
		WITH years AS (
			SELECT generate_series(?::int - 4, ?::int) AS year
		),
		yearly_methods AS (
			SELECT DISTINCT tx.payment_method
			FROM transactions tx
			JOIN merchants mer ON mer.merchant_id = tx.merchant_id AND mer.deleted_at IS NULL
			WHERE mer.api_key = ? AND tx.transaction_time >= ? AND tx.transaction_time < ? AND tx.deleted_at IS NULL
		),
		yearly_data AS (
			SELECT EXTRACT(YEAR FROM tx.transaction_time)::text AS year, tx.payment_method, SUM(tx.amount) AS total_amount
			FROM transactions tx
			JOIN merchants mer ON mer.merchant_id = tx.merchant_id AND mer.deleted_at IS NULL
			WHERE mer.api_key = ? AND tx.deleted_at IS NULL
			GROUP BY EXTRACT(YEAR FROM tx.transaction_time), tx.payment_method
		)
		SELECT y.year::text AS year, ym.payment_method, COALESCE(yd.total_amount, 0)::bigint AS total_amount
		FROM years y
		CROSS JOIN yearly_methods ym
		LEFT JOIN yearly_data yd ON y.year = yd.year::int AND yd.payment_method = ym.payment_method
		ORDER BY y.year, ym.payment_method
	`, req.Year, req.Year, req.Apikey, yearStart, yearEnd, req.Apikey).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodByApikeyFailed.WithInternal(err)
	}
	return res, nil
}
