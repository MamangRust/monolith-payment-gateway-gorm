package repositorystatsbycard

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardStatsBalanceByCardRepository struct {
	db *gorm.DB
}

func NewCardStatsBalanceByCardRepository(db *gorm.DB) CardStatsBalanceByCardRepository {
	return &cardStatsBalanceByCardRepository{db: db}
}

func (r *cardStatsBalanceByCardRepository) GetMonthlyBalancesByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyBalanceRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MonthlyBalanceRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(created_at, 'Mon') AS month, COALESCE(SUM(total_balance), 0)::int AS total_balance
		FROM saldos WHERE card_number = ? AND created_at >= ? AND created_at < ?
		GROUP BY TO_CHAR(created_at, 'Mon'), EXTRACT(MONTH FROM created_at)
		ORDER BY EXTRACT(MONTH FROM created_at)
	`, req.CardNumber, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetMonthlyBalanceByCardFailed.WithInternal(err)
	}
	return res, nil
}

func (r *cardStatsBalanceByCardRepository) GetYearlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyBalanceRow, error) {
	var res []*models.YearlyBalanceRow
	yearStart := time.Date(req.Year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM created_at)::text AS year, COALESCE(SUM(total_balance), 0)::bigint AS total_balance
		FROM saldos WHERE card_number = ? AND created_at >= ? AND created_at < ?
		GROUP BY EXTRACT(YEAR FROM created_at) ORDER BY EXTRACT(YEAR FROM created_at)
	`, req.CardNumber, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetYearlyBalanceByCardFailed.WithInternal(err)
	}
	return res, nil
}
