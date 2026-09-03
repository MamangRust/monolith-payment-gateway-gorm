package repositorystatsbycard

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardStatsTransferByCardRepository struct {
	db *gorm.DB
}

func NewCardStatsTransferByCardRepository(db *gorm.DB) CardStatsTransferByCardRepository {
	return &cardStatsTransferByCardRepository{db: db}
}

func (r *cardStatsTransferByCardRepository) GetMonthlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransferSentRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MonthlyTransferSentRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(transfer_time, 'Mon') AS month, COALESCE(SUM(transfer_amount), 0)::int AS total_sent_amount
		FROM transfers WHERE transfer_from = ? AND transfer_time >= ? AND transfer_time < ?
		GROUP BY TO_CHAR(transfer_time, 'Mon'), EXTRACT(MONTH FROM transfer_time)
		ORDER BY EXTRACT(MONTH FROM transfer_time)
	`, req.CardNumber, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountBySenderFailed
	}
	return res, nil
}

func (r *cardStatsTransferByCardRepository) GetYearlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransferSentRow, error) {
	var res []*models.YearlyTransferSentRow
	yearStart := time.Date(req.Year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM transfer_time)::text AS year, COALESCE(SUM(transfer_amount), 0)::bigint AS total_sent_amount
		FROM transfers WHERE transfer_from = ? AND transfer_time >= ? AND transfer_time < ?
		GROUP BY EXTRACT(YEAR FROM transfer_time) ORDER BY EXTRACT(YEAR FROM transfer_time)
	`, req.CardNumber, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountBySenderFailed
	}
	return res, nil
}

func (r *cardStatsTransferByCardRepository) GetMonthlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransferReceivedRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MonthlyTransferReceivedRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(transfer_time, 'Mon') AS month, COALESCE(SUM(transfer_amount), 0)::int AS total_received_amount
		FROM transfers WHERE transfer_to = ? AND transfer_time >= ? AND transfer_time < ?
		GROUP BY TO_CHAR(transfer_time, 'Mon'), EXTRACT(MONTH FROM transfer_time)
		ORDER BY EXTRACT(MONTH FROM transfer_time)
	`, req.CardNumber, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountByReceiverFailed.WithInternal(err)
	}
	return res, nil
}

func (r *cardStatsTransferByCardRepository) GetYearlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransferReceivedRow, error) {
	var res []*models.YearlyTransferReceivedRow
	yearStart := time.Date(req.Year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM transfer_time)::text AS year, COALESCE(SUM(transfer_amount), 0)::bigint AS total_received_amount
		FROM transfers WHERE transfer_to = ? AND transfer_time >= ? AND transfer_time < ?
		GROUP BY EXTRACT(YEAR FROM transfer_time) ORDER BY EXTRACT(YEAR FROM transfer_time)
	`, req.CardNumber, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountByReceiverFailed.WithInternal(err)
	}
	return res, nil
}
