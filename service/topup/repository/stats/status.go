package topupstatsrepository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
	"gorm.io/gorm"
)

type topupStatsStatusRepository struct {
	db *gorm.DB
}

func NewTopupStatsStatusRepository(db *gorm.DB) TopupStatsStatusRepository {
	return &topupStatsStatusRepository{db: db}
}

func (r *topupStatsStatusRepository) GetMonthTopupStatusSuccess(ctx context.Context, req *requests.MonthTopupStatus) ([]*models.TopupMonthlyStatusRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var res []*models.TopupMonthlyStatusRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT CAST(EXTRACT(YEAR FROM p.period_start) AS TEXT) AS year,
			CAST(EXTRACT(MONTH FROM p.period_start) AS TEXT) AS month,
			COALESCE(SUM(CASE WHEN t.status = 'success' THEN 1 ELSE 0 END), 0)::bigint AS total_success,
			COALESCE(SUM(CASE WHEN t.status = 'success' THEN t.topup_amount ELSE 0 END), 0)::int AS total_amount
		FROM (SELECT ?::date AS period_start, ?::date AS period_end
			UNION ALL SELECT ?::date, ?::date) p
		LEFT JOIN topups t ON t.topup_time >= p.period_start AND t.topup_time <= p.period_end AND t.deleted_at IS NULL
		GROUP BY p.period_start, p.period_start ORDER BY p.period_start
	`, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth).Scan(&res).Error
	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusSuccessFailed.WithInternal(err)
	}
	return res, nil
}

func (r *topupStatsStatusRepository) GetYearlyTopupStatusSuccess(ctx context.Context, year int) ([]*models.TopupYearlyStatusRow, error) {
	var res []*models.TopupYearlyStatusRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT CAST(y AS TEXT) AS year, COALESCE(SUM(CASE WHEN t.status = 'success' THEN 1 ELSE 0 END), 0)::int AS total_success,
			0::int AS total_failed, COALESCE(SUM(CASE WHEN t.status = 'success' THEN t.topup_amount ELSE 0 END), 0)::int AS total_amount
		FROM generate_series(?::int - 1, ?::int, 1) y
		LEFT JOIN topups t ON EXTRACT(YEAR FROM t.topup_time) = y AND t.deleted_at IS NULL
		GROUP BY y ORDER BY y
	`, year, year).Scan(&res).Error
	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusSuccessFailed.WithInternal(err)
	}
	return res, nil
}

func (r *topupStatsStatusRepository) GetMonthTopupStatusFailed(ctx context.Context, req *requests.MonthTopupStatus) ([]*models.TopupMonthlyStatusRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var res []*models.TopupMonthlyStatusRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT CAST(EXTRACT(YEAR FROM p.period_start) AS TEXT) AS year,
			CAST(EXTRACT(MONTH FROM p.period_start) AS TEXT) AS month,
			0::bigint AS total_success,
			COALESCE(SUM(CASE WHEN t.status = 'failed' THEN 1 ELSE 0 END), 0)::bigint AS total_failed,
			COALESCE(SUM(CASE WHEN t.status = 'failed' THEN t.topup_amount ELSE 0 END), 0)::int AS total_amount
		FROM (SELECT ?::date AS period_start, ?::date AS period_end
			UNION ALL SELECT ?::date, ?::date) p
		LEFT JOIN topups t ON t.topup_time >= p.period_start AND t.topup_time <= p.period_end AND t.deleted_at IS NULL
		GROUP BY p.period_start, p.period_start ORDER BY p.period_start
	`, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth).Scan(&res).Error
	if err != nil {
		return nil, topup_errors.ErrGetMonthTopupStatusFailedFailed.WithInternal(err)
	}
	return res, nil
}

func (r *topupStatsStatusRepository) GetYearlyTopupStatusFailed(ctx context.Context, year int) ([]*models.TopupYearlyStatusRow, error) {
	var res []*models.TopupYearlyStatusRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT CAST(y AS TEXT) AS year, 0::int AS total_success,
			COALESCE(SUM(CASE WHEN t.status = 'failed' THEN 1 ELSE 0 END), 0)::int AS total_failed,
			COALESCE(SUM(CASE WHEN t.status = 'failed' THEN t.topup_amount ELSE 0 END), 0)::int AS total_amount
		FROM generate_series(?::int - 1, ?::int, 1) y
		LEFT JOIN topups t ON EXTRACT(YEAR FROM t.topup_time) = y AND t.deleted_at IS NULL
		GROUP BY y ORDER BY y
	`, year, year).Scan(&res).Error
	if err != nil {
		return nil, topup_errors.ErrGetYearlyTopupStatusFailedFailed.WithInternal(err)
	}
	return res, nil
}
