package topupstatsbycardrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
)

type topupStatsByCardAmountRepository struct {
	db *gorm.DB
}

func NewTopupStatsByCardAmountRepository(db *gorm.DB) TopupStatsByCardAmountRepository {
	return &topupStatsByCardAmountRepository{db: db}
}

func (r *topupStatsByCardAmountRepository) GetMonthlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupMonthlyAmountRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	var results []*models.TopupMonthlyAmountRow

	err := r.db.WithContext(ctx).Raw(`
		WITH months AS (
			SELECT generate_series(
				date_trunc('year', ?::timestamp),
				date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day',
				interval '1 month'
			) AS month
		)
		SELECT TO_CHAR(m.month, 'Mon') AS month,
			COALESCE(SUM(t.topup_amount), 0)::int AS total_amount
		FROM months m
		LEFT JOIN topups t ON EXTRACT(MONTH FROM t.topup_time) = EXTRACT(MONTH FROM m.month)
			AND EXTRACT(YEAR FROM t.topup_time) = EXTRACT(YEAR FROM m.month)
			AND t.card_number = ?
			AND t.deleted_at IS NULL
		GROUP BY m.month
		ORDER BY m.month
	`, yearStart, yearStart, cardNumber).Scan(&results).Error

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupAmountsByCardFailed.WithInternal(err)
	}

	return results, nil
}

func (r *topupStatsByCardAmountRepository) GetYearlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupYearlyAmountRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	var results []*models.TopupYearlyAmountRow

	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM t.topup_time)::text AS year,
			SUM(t.topup_amount) AS total_amount
		FROM topups t
		WHERE t.deleted_at IS NULL
			AND t.card_number = ?
			AND EXTRACT(YEAR FROM t.topup_time) >= ? - 4
			AND EXTRACT(YEAR FROM t.topup_time) <= ?
		GROUP BY EXTRACT(YEAR FROM t.topup_time)
		ORDER BY year
	`, cardNumber, year, year).Scan(&results).Error

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupAmountsByCardFailed.WithInternal(err)
	}

	return results, nil
}
