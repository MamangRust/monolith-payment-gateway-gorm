package repositorystatsbycard

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardStatsTransactionByCardRepository struct {
	db *gorm.DB
}

func NewCardStatsTransactionByCardRepository(db *gorm.DB) CardStatsTransactionByCardRepository {
	return &cardStatsTransactionByCardRepository{db: db}
}

func (r *cardStatsTransactionByCardRepository) GetMonthlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransactionRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MonthlyTransactionRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(transaction_time, 'Mon') AS month, COALESCE(SUM(amount), 0)::int AS total_transaction_amount
		FROM transactions WHERE card_number = ? AND transaction_time >= ? AND transaction_time < ?
		GROUP BY TO_CHAR(transaction_time, 'Mon'), EXTRACT(MONTH FROM transaction_time)
		ORDER BY EXTRACT(MONTH FROM transaction_time)
	`, req.CardNumber, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransactionAmountByCardFailed.WithInternal(err)
	}
	return res, nil
}

func (r *cardStatsTransactionByCardRepository) GetYearlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransactionRow, error) {
	var res []*models.YearlyTransactionRow
	yearStart := time.Date(req.Year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM transaction_time)::text AS year, COALESCE(SUM(amount), 0)::bigint AS total_transaction_amount
		FROM transactions WHERE card_number = ? AND transaction_time >= ? AND transaction_time < ?
		GROUP BY EXTRACT(YEAR FROM transaction_time) ORDER BY EXTRACT(YEAR FROM transaction_time)
	`, req.CardNumber, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetYearlyTransactionAmountByCardFailed.WithInternal(err)
	}
	return res, nil
}
