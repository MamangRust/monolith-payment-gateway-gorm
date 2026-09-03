package repositorystats

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardStaticsTransferRepository struct {
	db *gorm.DB
}

func NewCardStatsTransferRepository(db *gorm.DB) CardStatsTransferRepository {
	return &cardStaticsTransferRepository{db: db}
}

func (r *cardStaticsTransferRepository) GetMonthlyTransferAmountSender(ctx context.Context, year int) ([]*models.MonthlyTransferSentRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MonthlyTransferSentRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(transfer_time, 'Mon') AS month, COALESCE(SUM(transfer_amount), 0)::int AS total_sent_amount
		FROM transfers WHERE transfer_time >= ? AND transfer_time < ?
		GROUP BY TO_CHAR(transfer_time, 'Mon'), EXTRACT(MONTH FROM transfer_time)
		ORDER BY EXTRACT(MONTH FROM transfer_time)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountSenderFailed
	}
	return res, nil
}

func (r *cardStaticsTransferRepository) GetYearlyTransferAmountSender(ctx context.Context, year int) ([]*models.YearlyTransferSentRow, error) {
	var res []*models.YearlyTransferSentRow
	yearStart := time.Date(year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM transfer_time)::text AS year, COALESCE(SUM(transfer_amount), 0)::bigint AS total_sent_amount
		FROM transfers WHERE transfer_time >= ? AND transfer_time < ?
		GROUP BY EXTRACT(YEAR FROM transfer_time) ORDER BY EXTRACT(YEAR FROM transfer_time)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountSenderFailed
	}
	return res, nil
}

func (r *cardStaticsTransferRepository) GetMonthlyTransferAmountReceiver(ctx context.Context, year int) ([]*models.MonthlyTransferReceivedRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MonthlyTransferReceivedRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(transfer_time, 'Mon') AS month, COALESCE(SUM(transfer_amount), 0)::int AS total_received_amount
		FROM transfers WHERE transfer_time >= ? AND transfer_time < ?
		GROUP BY TO_CHAR(transfer_time, 'Mon'), EXTRACT(MONTH FROM transfer_time)
		ORDER BY EXTRACT(MONTH FROM transfer_time)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountReceiverFailed.WithInternal(err)
	}
	return res, nil
}

func (r *cardStaticsTransferRepository) GetYearlyTransferAmountReceiver(ctx context.Context, year int) ([]*models.YearlyTransferReceivedRow, error) {
	var res []*models.YearlyTransferReceivedRow
	yearStart := time.Date(year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM transfer_time)::text AS year, COALESCE(SUM(transfer_amount), 0)::bigint AS total_received_amount
		FROM transfers WHERE transfer_time >= ? AND transfer_time < ?
		GROUP BY EXTRACT(YEAR FROM transfer_time) ORDER BY EXTRACT(YEAR FROM transfer_time)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountReceiverFailed.WithInternal(err)
	}
	return res, nil
}
