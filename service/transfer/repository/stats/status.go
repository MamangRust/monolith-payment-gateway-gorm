package transferstatsrepository

import (
	"context"
	"strconv"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferStatsStatusRepository struct {
	db *gorm.DB
}

func NewTransferStatsStatusRepository(db *gorm.DB) TransferStatsStatusRepository {
	return &transferStatsStatusRepository{db: db}
}

func (r *transferStatsStatusRepository) GetMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*models.TransferMonthlyStatusSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.TransferMonthlyStatusSuccessRow
	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.transfer_time)::text AS year, EXTRACT(MONTH FROM t.transfer_time)::integer AS month,
				COUNT(*) AS total_success, COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
			FROM transfers t WHERE t.deleted_at IS NULL AND t.status = 'success'
				AND ((t.transfer_time >= ?::timestamp AND t.transfer_time <= ?::timestamp)
					OR (t.transfer_time >= ?::timestamp AND t.transfer_time <= ?::timestamp))
			GROUP BY EXTRACT(YEAR FROM t.transfer_time), EXTRACT(MONTH FROM t.transfer_time)
		), formatted_data AS (
			SELECT year::text, TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month, total_success, total_amount FROM monthly_data
			UNION ALL
			SELECT EXTRACT(YEAR FROM ?::timestamp)::text AS year, TO_CHAR(?::timestamp, 'Mon') AS month, 0 AS total_success, 0 AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM monthly_data WHERE year = EXTRACT(YEAR FROM ?::timestamp)::text AND month = EXTRACT(MONTH FROM ?::timestamp)::integer)
			UNION ALL
			SELECT EXTRACT(YEAR FROM ?::timestamp)::text AS year, TO_CHAR(?::timestamp, 'Mon') AS month, 0 AS total_success, 0 AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM monthly_data WHERE year = EXTRACT(YEAR FROM ?::timestamp)::text AND month = EXTRACT(MONTH FROM ?::timestamp)::integer)
		)
		SELECT * FROM formatted_data ORDER BY year DESC, TO_DATE(month, 'Mon') DESC
	`, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth,
		currentDate, currentDate, currentDate, currentDate,
		prevDate, prevDate, prevDate, prevDate).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusSuccessFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transferStatsStatusRepository) GetYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*models.TransferYearlyStatusSuccessRow, error) {
	var results []*models.TransferYearlyStatusSuccessRow
	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.transfer_time)::text AS year, COUNT(*) AS total_success, COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
			FROM transfers t WHERE t.deleted_at IS NULL AND t.status = 'success'
				AND (EXTRACT(YEAR FROM t.transfer_time) = ? OR EXTRACT(YEAR FROM t.transfer_time) = ? - 1)
			GROUP BY EXTRACT(YEAR FROM t.transfer_time)
		), formatted_data AS (
			SELECT year::text, total_success::integer, total_amount::integer FROM yearly_data
			UNION ALL
			SELECT ? AS year, 0::integer AS total_success, 0::integer AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM yearly_data WHERE year::integer = ?)
			UNION ALL
			SELECT (? - 1)::text AS year, 0::integer AS total_success, 0::integer AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM yearly_data WHERE year::integer = ? - 1)
		)
		SELECT * FROM formatted_data ORDER BY year DESC
	`, strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year)).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusSuccessFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transferStatsStatusRepository) GetMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*models.TransferMonthlyStatusFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.TransferMonthlyStatusFailedRow
	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.transfer_time)::text AS year, EXTRACT(MONTH FROM t.transfer_time)::integer AS month,
				COUNT(*) AS total_failed, COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
			FROM transfers t WHERE t.deleted_at IS NULL AND t.status = 'failed'
				AND ((t.transfer_time >= ?::timestamp AND t.transfer_time <= ?::timestamp)
					OR (t.transfer_time >= ?::timestamp AND t.transfer_time <= ?::timestamp))
			GROUP BY EXTRACT(YEAR FROM t.transfer_time), EXTRACT(MONTH FROM t.transfer_time)
		), formatted_data AS (
			SELECT year::text, TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month, total_failed, total_amount FROM monthly_data
			UNION ALL
			SELECT EXTRACT(YEAR FROM ?::timestamp)::text AS year, TO_CHAR(?::timestamp, 'Mon') AS month, 0 AS total_failed, 0 AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM monthly_data WHERE year = EXTRACT(YEAR FROM ?::timestamp)::text AND month = EXTRACT(MONTH FROM ?::timestamp)::integer)
			UNION ALL
			SELECT EXTRACT(YEAR FROM ?::timestamp)::text AS year, TO_CHAR(?::timestamp, 'Mon') AS month, 0 AS total_failed, 0 AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM monthly_data WHERE year = EXTRACT(YEAR FROM ?::timestamp)::text AND month = EXTRACT(MONTH FROM ?::timestamp)::integer)
		)
		SELECT * FROM formatted_data ORDER BY year DESC, TO_DATE(month, 'Mon') DESC
	`, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth,
		currentDate, currentDate, currentDate, currentDate,
		prevDate, prevDate, prevDate, prevDate).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusFailedFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transferStatsStatusRepository) GetYearlyTransferStatusFailed(ctx context.Context, year int) ([]*models.TransferYearlyStatusFailedRow, error) {
	var results []*models.TransferYearlyStatusFailedRow
	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.transfer_time)::text AS year, COUNT(*) AS total_failed, COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
			FROM transfers t WHERE t.deleted_at IS NULL AND t.status = 'failed'
				AND (EXTRACT(YEAR FROM t.transfer_time) = ? OR EXTRACT(YEAR FROM t.transfer_time) = ? - 1)
			GROUP BY EXTRACT(YEAR FROM t.transfer_time)
		), formatted_data AS (
			SELECT year::text, total_failed::integer, total_amount::integer FROM yearly_data
			UNION ALL
			SELECT ? AS year, 0::integer AS total_failed, 0::integer AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM yearly_data WHERE year::integer = ?)
			UNION ALL
			SELECT (? - 1)::text AS year, 0::integer AS total_failed, 0::integer AS total_amount
			WHERE NOT EXISTS (SELECT 1 FROM yearly_data WHERE year::integer = ? - 1)
		)
		SELECT * FROM formatted_data ORDER BY year DESC
	`, strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year), strconv.Itoa(year)).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusFailedFailed.WithInternal(err)
	}
	return results, nil
}
