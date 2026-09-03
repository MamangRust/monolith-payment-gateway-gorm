package topupstatsrepository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
	"gorm.io/gorm"
)

type topupStatsMethodRepository struct {
	db *gorm.DB
}

func NewTopupStatsMethodRepository(db *gorm.DB) TopupStatsMethodRepository {
	return &topupStatsMethodRepository{db: db}
}

func (r *topupStatsMethodRepository) GetMonthlyTopupMethods(ctx context.Context, year int) ([]*models.TopupMonthlyMethodRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.TopupMonthlyMethodRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(topup_time, 'Mon') AS month, topup_method, COUNT(*)::int AS total_topups, COALESCE(SUM(topup_amount), 0)::int AS total_amount
		FROM topups WHERE topup_time >= ? AND topup_time < ? AND deleted_at IS NULL
		GROUP BY TO_CHAR(topup_time, 'Mon'), EXTRACT(MONTH FROM topup_time), topup_method
		ORDER BY EXTRACT(MONTH FROM topup_time), topup_method
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupMethodsFailed.WithInternal(err)
	}
	return res, nil
}

func (r *topupStatsMethodRepository) GetYearlyTopupMethods(ctx context.Context, year int) ([]*models.TopupYearlyMethodRow, error) {
	var res []*models.TopupYearlyMethodRow
	yearStart := time.Date(year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM topup_time)::text AS year, topup_method, COUNT(*)::bigint AS total_topups, COALESCE(SUM(topup_amount), 0)::bigint AS total_amount
		FROM topups WHERE topup_time >= ? AND topup_time < ? AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM topup_time), topup_method ORDER BY EXTRACT(YEAR FROM topup_time), topup_method
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupMethodsFailed.WithInternal(err)
	}
	return res, nil
}
