package topupstatsrepository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
	"gorm.io/gorm"
)

type topupStatsAmountRepository struct {
	db *gorm.DB
}

func NewTopupStatsAmountRepository(db *gorm.DB) TopupStatsAmountRepository {
	return &topupStatsAmountRepository{db: db}
}

func (r *topupStatsAmountRepository) GetMonthlyTopupAmounts(ctx context.Context, year int) ([]*models.TopupMonthlyAmountRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.TopupMonthlyAmountRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(topup_time, 'Mon') AS month, COALESCE(SUM(topup_amount), 0)::int AS total_amount
		FROM topups WHERE topup_time >= ? AND topup_time < ? AND deleted_at IS NULL
		GROUP BY TO_CHAR(topup_time, 'Mon'), EXTRACT(MONTH FROM topup_time)
		ORDER BY EXTRACT(MONTH FROM topup_time)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupAmountsFailed.WithInternal(err)
	}
	return res, nil
}

func (r *topupStatsAmountRepository) GetYearlyTopupAmounts(ctx context.Context, year int) ([]*models.TopupYearlyAmountRow, error) {
	var res []*models.TopupYearlyAmountRow
	yearStart := time.Date(year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM topup_time)::text AS year, COALESCE(SUM(topup_amount), 0)::bigint AS total_amount
		FROM topups WHERE topup_time >= ? AND topup_time < ? AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM topup_time) ORDER BY EXTRACT(YEAR FROM topup_time)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupAmountsFailed.WithInternal(err)
	}
	return res, nil
}
