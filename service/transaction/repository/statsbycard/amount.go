package transactionbycardrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
)

type transactionStatsByCardAmountRepository struct {
	db *gorm.DB
}

func NewTransactionStatsByCardAmountRepository(db *gorm.DB) TransactionStatsByCardAmountRepository {
	return &transactionStatsByCardAmountRepository{db: db}
}

func (r *transactionStatsByCardAmountRepository) GetMonthlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionMonthlyAmountByCardRow, error) {
	cardNumber := req.CardNumber
	year := req.Year
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	var results []*models.TransactionMonthlyAmountByCardRow

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
			AND t.card_number = ?
			AND t.deleted_at IS NULL
		GROUP BY m.month
		ORDER BY m.month
	`, yearStart, yearStart, cardNumber).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountsByCardFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionStatsByCardAmountRepository) GetYearlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionYearlyAmountByCardRow, error) {
	cardNumber := req.CardNumber
	year := req.Year

	var results []*models.TransactionYearlyAmountByCardRow

	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year,
			SUM(t.amount) AS total_amount
		FROM transactions t
		WHERE t.deleted_at IS NULL
			AND t.card_number = ?
			AND EXTRACT(YEAR FROM t.transaction_time) >= ? - 4
			AND EXTRACT(YEAR FROM t.transaction_time) <= ?
		GROUP BY EXTRACT(YEAR FROM t.transaction_time)
		ORDER BY year
	`, cardNumber, year, year).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountsByCardFailed.WithInternal(err)
	}

	return results, nil
}
