package saldostatsrepository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
	"gorm.io/gorm"
)

type saldoStatsBalanceRepository struct {
	db *gorm.DB
}

func NewSaldoStatsBalanceRepository(db *gorm.DB) SaldoStatsBalanceRepository {
	return &saldoStatsBalanceRepository{db: db}
}

func (r *saldoStatsBalanceRepository) GetMonthlySaldoBalances(ctx context.Context, year int) ([]*models.MonthlySaldoBalanceRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

	sql := `WITH months AS (
		SELECT generate_series(
			date_trunc('year', ?::timestamp),
			date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day',
			interval '1 month'
		) AS month
	)
	SELECT TO_CHAR(m.month, 'Mon') AS month, COALESCE(SUM(s.total_balance), 0)::int AS total_balance
	FROM months m
	LEFT JOIN saldos s ON EXTRACT(MONTH FROM s.created_at) = EXTRACT(MONTH FROM m.month)
		AND EXTRACT(YEAR FROM s.created_at) = EXTRACT(YEAR FROM m.month)
		AND s.deleted_at IS NULL
	GROUP BY m.month
	ORDER BY m.month`

	var results []*models.MonthlySaldoBalanceRow
	if err := r.db.WithContext(ctx).Raw(sql, yearStart, yearEnd).Scan(&results).Error; err != nil {
		return nil, saldo_errors.ErrGetMonthlySaldoBalancesFailed.WithInternal(err)
	}

	return results, nil
}

func (r *saldoStatsBalanceRepository) GetYearlySaldoBalances(ctx context.Context, year int) ([]*models.YearlySaldoBalanceRow, error) {
	sql := `WITH years AS (
		SELECT generate_series(?::int - 4, ?::int) AS year
	),
	yearly_data AS (
		SELECT EXTRACT(YEAR FROM s.created_at)::text AS year, SUM(s.total_balance) AS total_balance
		FROM saldos s
		WHERE s.deleted_at IS NULL
			AND EXTRACT(YEAR FROM s.created_at) >= ?::int - 4
			AND EXTRACT(YEAR FROM s.created_at) <= ?::int
		GROUP BY EXTRACT(YEAR FROM s.created_at)
	)
	SELECT y.year::text AS y_year, COALESCE(yd.total_balance, 0)::bigint AS total_balance
	FROM years y
	LEFT JOIN yearly_data yd ON y.year = yd.year::int
	ORDER BY y.year`

	var results []*models.YearlySaldoBalanceRow
	if err := r.db.WithContext(ctx).Raw(sql, year, year, year, year).Scan(&results).Error; err != nil {
		return nil, saldo_errors.ErrGetYearlySaldoBalancesFailed.WithInternal(err)
	}

	return results, nil
}
