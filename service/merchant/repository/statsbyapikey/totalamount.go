package merchantstatsapikeyrepository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	"gorm.io/gorm"
)

type merchantStatsTotalAmountByApiKeyRepository struct {
	db *gorm.DB
}

func NewMerchantStatsTotalAmountByApiKeyRepository(db *gorm.DB) MerchantStatsTotalAmountByApiKeyRepository {
	return &merchantStatsTotalAmountByApiKeyRepository{db: db}
}

func (r *merchantStatsTotalAmountByApiKeyRepository) GetMonthlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*models.MerchantMonthlyTotalAmountRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MerchantMonthlyTotalAmountRow
	yearEndMinus1 := yearEnd.AddDate(0, -11, 0)
	err := r.db.WithContext(ctx).Raw(`
		SELECT CAST(EXTRACT(YEAR FROM d) AS TEXT) AS year, CAST(EXTRACT(MONTH FROM d) AS TEXT) AS month, COALESCE(t.total, 0)::int AS total_amount
		FROM generate_series(?::date, ?::date, '1 month') d
		LEFT JOIN (
			SELECT DATE_TRUNC('month', tx.transaction_time) AS month, SUM(tx.amount) AS total
			FROM transactions tx JOIN merchants m ON m.merchant_id = tx.merchant_id AND m.deleted_at IS NULL
			WHERE m.api_key = ? AND tx.transaction_time >= ? AND tx.transaction_time < ? AND tx.deleted_at IS NULL
			GROUP BY DATE_TRUNC('month', tx.transaction_time)
		) t ON t.month = d
		ORDER BY d
	`, yearStart.Format("2006-01-02"), yearEndMinus1.Format("2006-01-02"), req.Apikey, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountByApikeyFailed.WithInternal(err)
	}
	return res, nil
}

func (r *merchantStatsTotalAmountByApiKeyRepository) GetYearlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*models.MerchantYearlyTotalAmountRow, error) {
	var res []*models.MerchantYearlyTotalAmountRow
	yearStart := time.Date(req.Year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT CAST(y AS TEXT) AS year, COALESCE(t.total, 0)::int AS total_amount
		FROM generate_series(?::int - 4, ?::int, 1) y
		LEFT JOIN (
			SELECT EXTRACT(YEAR FROM tx.transaction_time)::text AS year, SUM(tx.amount) AS total
			FROM transactions tx JOIN merchants m ON m.merchant_id = tx.merchant_id AND m.deleted_at IS NULL
			WHERE m.api_key = ? AND tx.transaction_time >= ? AND tx.transaction_time < ? AND tx.deleted_at IS NULL
			GROUP BY EXTRACT(YEAR FROM tx.transaction_time)
		) t ON t.year = y
		ORDER BY y
	`, req.Year, req.Year, req.Apikey, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountByApikeyFailed.WithInternal(err)
	}
	return res, nil
}
