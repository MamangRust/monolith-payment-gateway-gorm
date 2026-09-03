package transactionbycardrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
)

type transactionStatsByCardMethodRepository struct {
	db *gorm.DB
}

func NewTransactionStatsByCardMethodRepository(db *gorm.DB) TransactionStatsByCardMethodRepository {
	return &transactionStatsByCardMethodRepository{db: db}
}

func (r *transactionStatsByCardMethodRepository) GetMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionMonthlyPaymentMethodByCardRow, error) {
	year := req.Year
	cardNumber := req.CardNumber
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	var results []*models.TransactionMonthlyPaymentMethodByCardRow

	err := r.db.WithContext(ctx).Raw(`
		WITH months AS (
			SELECT generate_series(
				date_trunc('year', ?::timestamp),
				date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day',
				interval '1 month'
			) AS month
		), payment_methods AS (
			SELECT DISTINCT payment_method FROM transactions WHERE deleted_at IS NULL
		)
		SELECT TO_CHAR(m.month, 'Mon') AS month,
			pm.payment_method,
			COALESCE(COUNT(t.transaction_id), 0)::int AS total_count,
			COALESCE(SUM(t.amount), 0)::int AS total_amount
		FROM months m
		CROSS JOIN payment_methods pm
		LEFT JOIN transactions t ON EXTRACT(MONTH FROM t.transaction_time) = EXTRACT(MONTH FROM m.month)
			AND EXTRACT(YEAR FROM t.transaction_time) = EXTRACT(YEAR FROM m.month)
			AND t.payment_method = pm.payment_method
			AND t.card_number = ?
			AND t.deleted_at IS NULL
		GROUP BY m.month, pm.payment_method
		ORDER BY m.month, pm.payment_method
	`, yearStart, yearStart, cardNumber).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyPaymentMethodsByCardFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionStatsByCardMethodRepository) GetYearlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionYearlyPaymentMethodByCardRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	var results []*models.TransactionYearlyPaymentMethodByCardRow

	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM t.transaction_time)::text AS year,
			t.payment_method,
			COUNT(t.transaction_id) AS total_count,
			SUM(t.amount) AS total_amount
		FROM transactions t
		WHERE t.deleted_at IS NULL
			AND t.card_number = ?
			AND EXTRACT(YEAR FROM t.transaction_time) >= ? - 4
			AND EXTRACT(YEAR FROM t.transaction_time) <= ?
		GROUP BY EXTRACT(YEAR FROM t.transaction_time), t.payment_method
		ORDER BY year
	`, cardNumber, year, year).Scan(&results).Error

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyPaymentMethodsByCardFailed.WithInternal(err)
	}

	return results, nil
}
