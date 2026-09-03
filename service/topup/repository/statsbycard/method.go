package topupstatsbycardrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
)

type topupStatsByCardMethodRepository struct {
	db *gorm.DB
}

func NewTopupStatsByCardMethodRepository(db *gorm.DB) TopupStatsByCardMethodRepository {
	return &topupStatsByCardMethodRepository{db: db}
}

func (r *topupStatsByCardMethodRepository) GetMonthlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupMonthlyMethodByCardRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	var results []*models.TopupMonthlyMethodByCardRow

	err := r.db.WithContext(ctx).Raw(`
		WITH months AS (
			SELECT generate_series(
				date_trunc('year', ?::timestamp),
				date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day',
				interval '1 month'
			) AS month
		), topup_methods AS (
			SELECT DISTINCT topup_method FROM topups WHERE deleted_at IS NULL
		)
		SELECT TO_CHAR(m.month, 'Mon') AS month,
			tm.topup_method,
			COALESCE(COUNT(t.topup_id), 0)::int AS total_topups,
			COALESCE(SUM(t.topup_amount), 0)::int AS total_amount
		FROM months m
		CROSS JOIN topup_methods tm
		LEFT JOIN topups t ON EXTRACT(MONTH FROM t.topup_time) = EXTRACT(MONTH FROM m.month)
			AND EXTRACT(YEAR FROM t.topup_time) = EXTRACT(YEAR FROM m.month)
			AND t.topup_method = tm.topup_method
			AND t.card_number = ?
			AND t.deleted_at IS NULL
		GROUP BY m.month, tm.topup_method
		ORDER BY m.month, tm.topup_method
	`, yearStart, yearStart, cardNumber).Scan(&results).Error

	if err != nil {
		return nil, topup_errors.ErrGetMonthlyTopupMethodsByCardFailed.WithInternal(err)
	}

	return results, nil
}

func (r *topupStatsByCardMethodRepository) GetYearlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupYearlyMethodByCardRow, error) {
	year := req.Year
	cardNumber := req.CardNumber

	var results []*models.TopupYearlyMethodByCardRow

	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM t.topup_time)::text AS year,
			t.topup_method,
			COUNT(t.topup_id) AS total_topups,
			SUM(t.topup_amount) AS total_amount
		FROM topups t
		WHERE t.deleted_at IS NULL
			AND t.card_number = ?
			AND EXTRACT(YEAR FROM t.topup_time) >= ? - 4
			AND EXTRACT(YEAR FROM t.topup_time) <= ?
		GROUP BY EXTRACT(YEAR FROM t.topup_time), t.topup_method
		ORDER BY year
	`, cardNumber, year, year).Scan(&results).Error

	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupMethodsByCardFailed.WithInternal(err)
	}

	return results, nil
}
