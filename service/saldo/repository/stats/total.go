package saldostatsrepository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
	"gorm.io/gorm"
)

type saldoStatsTotalBalanceRepository struct {
	db *gorm.DB
}

func NewSaldoStatsTotalBalanceRepository(db *gorm.DB) SaldoStatsTotalSaldoRepository {
	return &saldoStatsTotalBalanceRepository{db: db}
}

func (r *saldoStatsTotalBalanceRepository) GetMonthlyTotalSaldoBalance(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*models.MonthlyTotalSaldoBalanceRow, error) {
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	sql := `WITH monthly_data AS (
		SELECT EXTRACT(YEAR FROM s.created_at)::text AS year,
			EXTRACT(MONTH FROM s.created_at)::integer AS month,
			COALESCE(SUM(s.total_balance), 0) AS total_balance
		FROM saldos s
		WHERE s.deleted_at IS NULL
			AND (
				(s.created_at >= ?::timestamp AND s.created_at <= ?::timestamp)
				OR (s.created_at >= ?::timestamp AND s.created_at <= ?::timestamp)
			)
		GROUP BY EXTRACT(YEAR FROM s.created_at), EXTRACT(MONTH FROM s.created_at)
	),
	formatted_data AS (
		SELECT year::text, TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month, total_balance::integer
		FROM monthly_data
	)
	SELECT year, month, total_balance
	FROM formatted_data
	ORDER BY year DESC, TO_DATE(month, 'Mon') DESC`

	var results []*models.MonthlyTotalSaldoBalanceRow
	if err := r.db.WithContext(ctx).Raw(sql, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth).Scan(&results).Error; err != nil {
		return nil, saldo_errors.ErrGetMonthlyTotalSaldoBalanceFailed.WithInternal(err)
	}

	return results, nil
}

func (r *saldoStatsTotalBalanceRepository) GetYearTotalSaldoBalance(ctx context.Context, year int) ([]*models.YearlyTotalSaldoBalancesRow, error) {
	sql := `WITH yearly_data AS (
		SELECT EXTRACT(YEAR FROM s.created_at)::text AS year,
			COALESCE(SUM(s.total_balance), 0)::integer AS total_balance
		FROM saldos s
		WHERE s.deleted_at IS NULL
			AND (EXTRACT(YEAR FROM s.created_at) = ?::integer
				OR EXTRACT(YEAR FROM s.created_at) = ?::integer - 1)
		GROUP BY EXTRACT(YEAR FROM s.created_at)
	),
	formatted_data AS (
		SELECT year::text, total_balance::integer
		FROM yearly_data
	)
	SELECT year, total_balance
	FROM formatted_data
	ORDER BY year DESC`

	var results []*models.YearlyTotalSaldoBalancesRow
	if err := r.db.WithContext(ctx).Raw(sql, year, year).Scan(&results).Error; err != nil {
		return nil, saldo_errors.ErrGetYearTotalSaldoBalanceFailed.WithInternal(err)
	}

	return results, nil
}
