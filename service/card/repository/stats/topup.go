package repositorystats

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardStatsTopupRepository struct {
	db *gorm.DB
}

func NewCardStatsTopupRepository(db *gorm.DB) CardStatsTopupRepository {
	return &cardStatsTopupRepository{db: db}
}

func (r *cardStatsTopupRepository) GetMonthlyTopupAmount(ctx context.Context, year int) ([]*models.MonthlyTopupRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MonthlyTopupRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(topup_time, 'Mon') AS month, COALESCE(SUM(topup_amount), 0)::int AS total_topup_amount
		FROM topups WHERE topup_time >= ? AND topup_time < ?
		GROUP BY TO_CHAR(topup_time, 'Mon'), EXTRACT(MONTH FROM topup_time)
		ORDER BY EXTRACT(MONTH FROM topup_time)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetMonthlyTopupAmountFailed.WithInternal(err)
	}
	return res, nil
}

func (r *cardStatsTopupRepository) GetYearlyTopupAmount(ctx context.Context, year int) ([]*models.YearlyTopupRow, error) {
	var res []*models.YearlyTopupRow
	yearStart := time.Date(year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM topup_time)::text AS year, COALESCE(SUM(topup_amount), 0)::bigint AS total_topup_amount
		FROM topups WHERE topup_time >= ? AND topup_time < ?
		GROUP BY EXTRACT(YEAR FROM topup_time) ORDER BY EXTRACT(YEAR FROM topup_time)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetYearlyTopupAmountFailed.WithInternal(err)
	}
	return res, nil
}
