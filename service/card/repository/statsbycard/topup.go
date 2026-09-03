package repositorystatsbycard

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardStatsTopupByCardRepository struct {
	db *gorm.DB
}

func NewCardStatsTopupByCardRepository(db *gorm.DB) CardStatsTopupByCardRepository {
	return &cardStatsTopupByCardRepository{db: db}
}

func (r *cardStatsTopupByCardRepository) GetMonthlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTopupRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MonthlyTopupRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(topup_time, 'Mon') AS month, COALESCE(SUM(topup_amount), 0)::int AS total_topup_amount
		FROM topups WHERE card_number = ? AND topup_time >= ? AND topup_time < ?
		GROUP BY TO_CHAR(topup_time, 'Mon'), EXTRACT(MONTH FROM topup_time)
		ORDER BY EXTRACT(MONTH FROM topup_time)
	`, req.CardNumber, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetMonthlyTopupAmountByCardFailed.WithInternal(err)
	}
	return res, nil
}

func (r *cardStatsTopupByCardRepository) GetYearlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTopupRow, error) {
	var res []*models.YearlyTopupRow
	yearStart := time.Date(req.Year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM topup_time)::text AS year, COALESCE(SUM(topup_amount), 0)::bigint AS total_topup_amount
		FROM topups WHERE card_number = ? AND topup_time >= ? AND topup_time < ?
		GROUP BY EXTRACT(YEAR FROM topup_time) ORDER BY EXTRACT(YEAR FROM topup_time)
	`, req.CardNumber, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetYearlyTopupAmountByCardFailed.WithInternal(err)
	}
	return res, nil
}
