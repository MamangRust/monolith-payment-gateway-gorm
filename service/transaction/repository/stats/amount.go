package transactionstatsrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
)

type transactionStatsAmountRepository struct {
	db *gorm.DB
}

func NewTransactionStatsAmountRepository(db *gorm.DB) TransactionStatsAmountRepository {
	return &transactionStatsAmountRepository{db: db}
}

func (r *transactionStatsAmountRepository) GetMonthlyAmounts(ctx context.Context, year int) ([]*models.TransactionMonthlyAmountRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

	var results []*models.TransactionMonthlyAmountRow

	err := r.db.WithContext(ctx).Raw(`
		WITH months AS (
			SELECT generate_series(
				date_trunc('year', ?::timestamp),
				date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day',
				interval '1 month'
			) AS month
		)
		SELECT TO_CHAR(m.month, 'Mon') AS month,
			COALESCE(SUM(t.amount), 0)::int AS total_amount
		FROM months m
		LEFT JOIN transactions t ON EXTRACT(MONTH FROM t.transaction_time) = EXTRACT(MONTH FROM m.month)
			AND EXTRACT(YEAR FROM t.transaction_time) = EXTRACT(YEAR FROM m.month)
			AND t.deleted_at IS NULL
		GROUP BY m.month
		ORDER BY m.month
	`, yearStart, yearEnd).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountsFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionStatsAmountRepository) GetYearlyAmounts(ctx context.Context, year int) ([]*models.TransactionYearlyAmountRow, error) {
	var results []*models.TransactionYearlyAmountRow

	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year,
			SUM(t.amount) AS total_amount
		FROM transactions t
		WHERE t.deleted_at IS NULL
			AND EXTRACT(YEAR FROM t.transaction_time) >= ? - 4
			AND EXTRACT(YEAR FROM t.transaction_time) <= ?
		GROUP BY EXTRACT(YEAR FROM t.transaction_time)
		ORDER BY year
	`, year, year).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountsFailed.WithInternal(err)
	}

	return results, nil
}
