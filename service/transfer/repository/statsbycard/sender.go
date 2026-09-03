package transferstatsbycardrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferStatsByCardAmountSenderRepository struct {
	db *gorm.DB
}

func NewTransferStatsByCardAmountSenderRepository(db *gorm.DB) TransferStatsByCardAmountSenderRepository {
	return &transferStatsByCardAmountSenderRepository{db: db}
}

func (r *transferStatsByCardAmountSenderRepository) GetMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferMonthlyAmountBySenderCardRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	var results []*models.TransferMonthlyAmountBySenderCardRow
	err := r.db.WithContext(ctx).Raw(`
		WITH months AS (
			SELECT generate_series(date_trunc('year', ?::timestamp), date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day', interval '1 month') AS month
		)
		SELECT TO_CHAR(m.month, 'Mon') AS month, COALESCE(SUM(t.transfer_amount), 0)::int AS total_transfer_amount
		FROM months m LEFT JOIN transfers t ON EXTRACT(MONTH FROM t.transfer_time) = EXTRACT(MONTH FROM m.month)
			AND EXTRACT(YEAR FROM t.transfer_time) = EXTRACT(YEAR FROM m.month) AND t.transfer_from = ? AND t.deleted_at IS NULL
		GROUP BY m.month ORDER BY m.month
	`, yearStart, req.CardNumber).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsBySenderCardFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transferStatsByCardAmountSenderRepository) GetYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferYearlyAmountBySenderCardRow, error) {
	var results []*models.TransferYearlyAmountBySenderCardRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM t.transfer_time)::text AS year, SUM(t.transfer_amount) AS total_transfer_amount
		FROM transfers t WHERE t.deleted_at IS NULL AND t.transfer_from = ?
			AND EXTRACT(YEAR FROM t.transfer_time) >= ? - 4 AND EXTRACT(YEAR FROM t.transfer_time) <= ?
		GROUP BY EXTRACT(YEAR FROM t.transfer_time) ORDER BY year
	`, req.CardNumber, req.Year, req.Year).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsBySenderCardFailed.WithInternal(err)
	}
	return results, nil
}
