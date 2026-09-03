package merchantstatsmerchantrepository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	"gorm.io/gorm"
)

type merchantStatsTotalAmountByMerchantRepository struct {
	db *gorm.DB
}

func NewMerchantStatsTotalAmountByMerchantRepository(db *gorm.DB) MerchantStatsTotalAmountByMerchantRepository {
	return &merchantStatsTotalAmountByMerchantRepository{db: db}
}

func (r *merchantStatsTotalAmountByMerchantRepository) GetMonthlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*models.MerchantMonthlyTotalAmountRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MerchantMonthlyTotalAmountRow
	yearEndMinus1 := yearEnd.AddDate(0, -11, 0)
	err := r.db.WithContext(ctx).Raw(`
		SELECT CAST(EXTRACT(YEAR FROM d) AS TEXT) AS year, CAST(EXTRACT(MONTH FROM d) AS TEXT) AS month, COALESCE(t.total, 0)::int AS total_amount
		FROM generate_series(?::date, ?::date, '1 month') d
		LEFT JOIN (
			SELECT DATE_TRUNC('month', transaction_time) AS month, SUM(amount) AS total
			FROM transactions WHERE merchant_id = ? AND transaction_time >= ? AND transaction_time < ? AND deleted_at IS NULL
			GROUP BY DATE_TRUNC('month', transaction_time)
		) t ON t.month = d
		ORDER BY d
	`, yearStart.Format("2006-01-02"), yearEndMinus1.Format("2006-01-02"), req.MerchantID, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountByMerchantsFailed.WithInternal(err)
	}
	return res, nil
}

func (r *merchantStatsTotalAmountByMerchantRepository) GetYearlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*models.MerchantYearlyTotalAmountRow, error) {
	var res []*models.MerchantYearlyTotalAmountRow
	yearStart := time.Date(req.Year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT CAST(y AS TEXT) AS year, COALESCE(t.total, 0)::int AS total_amount
		FROM generate_series(?::int - 4, ?::int, 1) y
		LEFT JOIN (
			SELECT EXTRACT(YEAR FROM transaction_time)::text AS year, SUM(amount) AS total
			FROM transactions WHERE merchant_id = ? AND transaction_time >= ? AND transaction_time < ? AND deleted_at IS NULL
			GROUP BY EXTRACT(YEAR FROM transaction_time)
		) t ON t.year = y
		ORDER BY y
	`, req.Year, req.Year, req.MerchantID, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountByMerchantsFailed.WithInternal(err)
	}
	return res, nil
}
