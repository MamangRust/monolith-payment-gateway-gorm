package repositorystatsbycard

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardStatsWithdrawByCardRepository struct {
	db *gorm.DB
}

func NewCardStatsWithdrawByCardRepository(db *gorm.DB) CardStatsWithdrawByCardRepository {
	return &cardStatsWithdrawByCardRepository{db: db}
}

func (r *cardStatsWithdrawByCardRepository) GetMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyWithdrawRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MonthlyWithdrawRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(withdraw_time, 'Mon') AS month, COALESCE(SUM(withdraw_amount), 0)::int AS total_withdraw_amount
		FROM withdraws WHERE card_number = ? AND withdraw_time >= ? AND withdraw_time < ?
		GROUP BY TO_CHAR(withdraw_time, 'Mon'), EXTRACT(MONTH FROM withdraw_time)
		ORDER BY EXTRACT(MONTH FROM withdraw_time)
	`, req.CardNumber, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetMonthlyWithdrawAmountByCardFailed
	}
	return res, nil
}

func (r *cardStatsWithdrawByCardRepository) GetYearlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyWithdrawRow, error) {
	var res []*models.YearlyWithdrawRow
	yearStart := time.Date(req.Year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM withdraw_time)::text AS year, COALESCE(SUM(withdraw_amount), 0)::bigint AS total_withdraw_amount
		FROM withdraws WHERE card_number = ? AND withdraw_time >= ? AND withdraw_time < ?
		GROUP BY EXTRACT(YEAR FROM withdraw_time) ORDER BY EXTRACT(YEAR FROM withdraw_time)
	`, req.CardNumber, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetYearlyWithdrawAmountByCardFailed
	}
	return res, nil
}
