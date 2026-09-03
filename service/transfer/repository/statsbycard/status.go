package transferstatsbycardrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferStatsByCardStatusRepository struct {
	db *gorm.DB
}

func NewTransferStatsByCardStatusRepository(db *gorm.DB) TransferStatsByCardStatusRepository {
	return &transferStatsByCardStatusRepository{db: db}
}

func (r *transferStatsByCardStatusRepository) GetMonthTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*models.TransferMonthlyStatusSuccessByCardRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.TransferMonthlyStatusSuccessByCardRow
	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.transfer_time)::text AS year, EXTRACT(MONTH FROM t.transfer_time)::integer AS month,
				COUNT(*) AS total_success, COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
			FROM transfers t WHERE t.deleted_at IS NULL AND t.status = 'success'
				AND (t.transfer_from = ? OR t.transfer_to = ?)
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
	`, req.CardNumber, req.CardNumber, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth,
		currentDate, currentDate, currentDate, currentDate,
		prevDate, prevDate, prevDate, prevDate).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusSuccessByCardFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transferStatsByCardStatusRepository) GetYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*models.TransferYearlyStatusSuccessByCardRow, error) {
	var results []*models.TransferYearlyStatusSuccessByCardRow
	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.transfer_time)::text AS year, COUNT(*) AS total_success, COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
			FROM transfers t WHERE t.deleted_at IS NULL AND t.status = 'success'
				AND (t.transfer_from = ? OR t.transfer_to = ?)
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
	`, req.CardNumber, req.CardNumber, req.Year, req.Year, req.Year, req.Year, req.Year, req.Year).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusSuccessByCardFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transferStatsByCardStatusRepository) GetMonthTransferStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*models.TransferMonthlyStatusFailedByCardRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var results []*models.TransferMonthlyStatusFailedByCardRow
	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_data AS (
			SELECT EXTRACT(YEAR FROM t.transfer_time)::text AS year, EXTRACT(MONTH FROM t.transfer_time)::integer AS month,
				COUNT(*) AS total_failed, COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
			FROM transfers t WHERE t.deleted_at IS NULL AND t.status = 'failed'
				AND (t.transfer_from = ? OR t.transfer_to = ?)
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
	`, req.CardNumber, req.CardNumber, currentDate, lastDayCurrentMonth, prevDate, lastDayPrevMonth,
		currentDate, currentDate, currentDate, currentDate,
		prevDate, prevDate, prevDate, prevDate).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusFailedByCardFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transferStatsByCardStatusRepository) GetYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*models.TransferYearlyStatusFailedByCardRow, error) {
	var results []*models.TransferYearlyStatusFailedByCardRow
	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT EXTRACT(YEAR FROM t.transfer_time)::text AS year, COUNT(*) AS total_failed, COALESCE(SUM(t.transfer_amount), 0)::integer AS total_amount
			FROM transfers t WHERE t.deleted_at IS NULL AND t.status = 'failed'
				AND (t.transfer_from = ? OR t.transfer_to = ?)
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
	`, req.CardNumber, req.CardNumber, req.Year, req.Year, req.Year, req.Year, req.Year, req.Year).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusFailedByCardFailed.WithInternal(err)
	}
	return results, nil
}
