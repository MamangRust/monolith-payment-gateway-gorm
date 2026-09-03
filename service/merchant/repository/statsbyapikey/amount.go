package merchantstatsapikeyrepository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	"gorm.io/gorm"
)

type merchantStatsAmountByApiKeyRepository struct {
	db *gorm.DB
}

func NewMerchantStatsAmountByApiKeyRepository(db *gorm.DB) MerchantStatsAmountByApiKeyRepository {
	return &merchantStatsAmountByApiKeyRepository{db: db}
}

func (r *merchantStatsAmountByApiKeyRepository) GetMonthlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*models.MerchantMonthlyAmountRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MerchantMonthlyAmountRow
	err := r.db.WithContext(ctx).Raw(`
		WITH months AS (
			SELECT generate_series(
				date_trunc('year', ?::timestamp),
				date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day',
				interval '1 month'
			) AS month
		),
		apikey_data AS (
			SELECT EXTRACT(MONTH FROM t.transaction_time) AS month_num, SUM(t.amount) AS total_amount
			FROM transactions t
			JOIN merchants m ON m.merchant_id = t.merchant_id AND m.deleted_at IS NULL
			WHERE m.api_key = ? AND t.transaction_time >= ? AND t.transaction_time < ? AND t.deleted_at IS NULL
			GROUP BY EXTRACT(MONTH FROM t.transaction_time)
		)
		SELECT TO_CHAR(mo.month, 'Mon') AS month, COALESCE(ad.total_amount, 0)::int AS total_amount
		FROM months mo
		LEFT JOIN apikey_data ad ON EXTRACT(MONTH FROM mo.month) = ad.month_num
		ORDER BY mo.month
	`, yearStart, yearStart, req.Apikey, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountByApikeyFailed.WithInternal(err)
	}
	return res, nil
}

func (r *merchantStatsAmountByApiKeyRepository) GetYearlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*models.MerchantYearlyAmountRow, error) {
	var res []*models.MerchantYearlyAmountRow
	err := r.db.WithContext(ctx).Raw(`
		WITH years AS (
			SELECT generate_series(?::int - 4, ?::int) AS year
		),
		yearly_data AS (
			SELECT EXTRACT(YEAR FROM tx.transaction_time)::text AS year, SUM(tx.amount) AS total_amount
			FROM transactions tx
			JOIN merchants mer ON mer.merchant_id = tx.merchant_id AND mer.deleted_at IS NULL
			WHERE mer.api_key = ? AND tx.deleted_at IS NULL
			GROUP BY EXTRACT(YEAR FROM tx.transaction_time)
		)
		SELECT y.year::text AS year, COALESCE(yd.total_amount, 0)::bigint AS total_amount
		FROM years y
		LEFT JOIN yearly_data yd ON y.year = yd.year::int
		ORDER BY y.year
	`, req.Year, req.Apikey).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountByApikeyFailed.WithInternal(err)
	}
	return res, nil
}
