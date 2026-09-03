package repositorystats

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardStatsWithdrawRepository struct {
	db *gorm.DB
}

func NewCardStatsWithdrawRepository(db *gorm.DB) CardStatsWithdrawRepository {
	return &cardStatsWithdrawRepository{db: db}
}

func (r *cardStatsWithdrawRepository) GetMonthlyWithdrawAmount(ctx context.Context, year int) ([]*models.MonthlyWithdrawRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MonthlyWithdrawRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(withdraw_time, 'Mon') AS month, COALESCE(SUM(withdraw_amount), 0)::int AS total_withdraw_amount
		FROM withdraws WHERE withdraw_time >= ? AND withdraw_time < ?
		GROUP BY TO_CHAR(withdraw_time, 'Mon'), EXTRACT(MONTH FROM withdraw_time)
		ORDER BY EXTRACT(MONTH FROM withdraw_time)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetMonthlyWithdrawAmountFailed
	}
	return res, nil
}

func (r *cardStatsWithdrawRepository) GetYearlyWithdrawAmount(ctx context.Context, year int) ([]*models.YearlyWithdrawRow, error) {
	var res []*models.YearlyWithdrawRow
	yearStart := time.Date(year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM withdraw_time)::text AS year, COALESCE(SUM(withdraw_amount), 0)::bigint AS total_withdraw_amount
		FROM withdraws WHERE withdraw_time >= ? AND withdraw_time < ?
		GROUP BY EXTRACT(YEAR FROM withdraw_time) ORDER BY EXTRACT(YEAR FROM withdraw_time)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetYearlyWithdrawAmountFailed
	}
	return res, nil
}
