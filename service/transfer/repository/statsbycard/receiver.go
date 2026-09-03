package transferstatsbycardrepository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferStatsByCardAmountReceiverRepository struct {
	db *gorm.DB
}

func NewTransferStatsByCardAmountReceiverRepository(db *gorm.DB) TransferStatsByCardAmountReceiverRepository {
	return &transferStatsByCardAmountReceiverRepository{db: db}
}

func (r *transferStatsByCardAmountReceiverRepository) GetMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferMonthlyAmountByReceiverCardRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	var results []*models.TransferMonthlyAmountByReceiverCardRow
	err := r.db.WithContext(ctx).Raw(`
		WITH months AS (
			SELECT generate_series(date_trunc('year', ?::timestamp), date_trunc('year', ?::timestamp) + interval '1 year' - interval '1 day', interval '1 month') AS month
		)
		SELECT TO_CHAR(m.month, 'Mon') AS month, COALESCE(SUM(t.transfer_amount), 0)::int AS total_transfer_amount
		FROM months m LEFT JOIN transfers t ON EXTRACT(MONTH FROM t.transfer_time) = EXTRACT(MONTH FROM m.month)
			AND EXTRACT(YEAR FROM t.transfer_time) = EXTRACT(YEAR FROM m.month) AND t.transfer_to = ? AND t.deleted_at IS NULL
		GROUP BY m.month ORDER BY m.month
	`, yearStart, req.CardNumber).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsByReceiverCardFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transferStatsByCardAmountReceiverRepository) GetYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferYearlyAmountByReceiverCardRow, error) {
	var results []*models.TransferYearlyAmountByReceiverCardRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM t.transfer_time)::text AS year, SUM(t.transfer_amount) AS total_transfer_amount
		FROM transfers t WHERE t.deleted_at IS NULL AND t.transfer_to = ?
			AND EXTRACT(YEAR FROM t.transfer_time) >= ? - 4 AND EXTRACT(YEAR FROM t.transfer_time) <= ?
		GROUP BY EXTRACT(YEAR FROM t.transfer_time) ORDER BY year
	`, req.CardNumber, req.Year, req.Year).Scan(&results).Error
	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsByReceiverCardFailed.WithInternal(err)
	}
	return results, nil
}
