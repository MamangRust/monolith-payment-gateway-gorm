package transferstatsrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferStatsAmountRepository struct {
	db *gorm.DB
}

func NewTransferStatsAmountRepository(db *gorm.DB) TransferStatsAmountRepository {
	return &transferStatsAmountRepository{db: db}
}

func (r *transferStatsAmountRepository) GetMonthlyTransferAmounts(ctx context.Context, year int) ([]*models.TransferMonthlyAmountRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var results []*models.TransferMonthlyAmountRow
	err := r.db.WithContext(ctx).Raw(`
		WITH months AS (
			SELECT generate_series(date_trunc('year', ?::timestamp), date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day', interval '1 month') AS month
		)
		SELECT TO_CHAR(m.month, 'Mon') AS month, COALESCE(SUM(t.transfer_amount), 0)::int AS total_transfer_amount
		FROM months m LEFT JOIN transfers t ON EXTRACT(MONTH FROM t.transfer_time) = EXTRACT(MONTH FROM m.month)
			AND EXTRACT(YEAR FROM t.transfer_time) = EXTRACT(YEAR FROM m.month) AND t.deleted_at IS NULL
		GROUP BY m.month ORDER BY m.month
	`, yearStart, yearEnd).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transferStatsAmountRepository) GetYearlyTransferAmounts(ctx context.Context, year int) ([]*models.TransferYearlyAmountRow, error) {
	var results []*models.TransferYearlyAmountRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM t.transfer_time)::text AS year, SUM(t.transfer_amount) AS total_transfer_amount
		FROM transfers t WHERE t.deleted_at IS NULL
			AND EXTRACT(YEAR FROM t.transfer_time) >= ? - 4 AND EXTRACT(YEAR FROM t.transfer_time) <= ?
		GROUP BY EXTRACT(YEAR FROM t.transfer_time) ORDER BY year
	`, year, year).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsFailed.WithInternal(err)
	}
	return results, nil
}
