package repositorystats

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardStatsTransactionRepository struct {
	db *gorm.DB
}

func NewCardStatsTransactionRepository(db *gorm.DB) CardStatsTransactionRepository {
	return &cardStatsTransactionRepository{db: db}
}

func (r *cardStatsTransactionRepository) GetMonthlyTransactionAmount(ctx context.Context, year int) ([]*models.MonthlyTransactionRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MonthlyTransactionRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(transaction_time, 'Mon') AS month, COALESCE(SUM(amount), 0)::int AS total_transaction_amount
		FROM transactions WHERE transaction_time >= ? AND transaction_time < ?
		GROUP BY TO_CHAR(transaction_time, 'Mon'), EXTRACT(MONTH FROM transaction_time)
		ORDER BY EXTRACT(MONTH FROM transaction_time)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransactionAmountFailed.WithInternal(err)
	}
	return res, nil
}

func (r *cardStatsTransactionRepository) GetYearlyTransactionAmount(ctx context.Context, year int) ([]*models.YearlyTransactionRow, error) {
	var res []*models.YearlyTransactionRow
	yearStart := time.Date(year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM transaction_time)::text AS year, COALESCE(SUM(amount), 0)::bigint AS total_transaction_amount
		FROM transactions WHERE transaction_time >= ? AND transaction_time < ?
		GROUP BY EXTRACT(YEAR FROM transaction_time) ORDER BY EXTRACT(YEAR FROM transaction_time)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetYearlyTransactionAmountFailed.WithInternal(err)
	}
	return res, nil
}
