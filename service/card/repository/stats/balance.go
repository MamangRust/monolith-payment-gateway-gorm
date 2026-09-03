package repositorystats

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type CardStatsBalanceRepository interface {
	GetMonthlyBalance(ctx context.Context, year int) ([]*models.MonthlyBalanceRow, error)
	GetYearlyBalance(ctx context.Context, year int) ([]*models.YearlyBalanceRow, error)
}

type CardStatsTopupRepository interface {
	GetMonthlyTopupAmount(ctx context.Context, year int) ([]*models.MonthlyTopupRow, error)
	GetYearlyTopupAmount(ctx context.Context, year int) ([]*models.YearlyTopupRow, error)
}

type CardStatsWithdrawRepository interface {
	GetMonthlyWithdrawAmount(ctx context.Context, year int) ([]*models.MonthlyWithdrawRow, error)
	GetYearlyWithdrawAmount(ctx context.Context, year int) ([]*models.YearlyWithdrawRow, error)
}

type CardStatsTransactionRepository interface {
	GetMonthlyTransactionAmount(ctx context.Context, year int) ([]*models.MonthlyTransactionRow, error)
	GetYearlyTransactionAmount(ctx context.Context, year int) ([]*models.YearlyTransactionRow, error)
}

type CardStatsTransferRepository interface {
	GetMonthlyTransferAmountSender(ctx context.Context, year int) ([]*models.MonthlyTransferSentRow, error)
	GetYearlyTransferAmountSender(ctx context.Context, year int) ([]*models.YearlyTransferSentRow, error)
	GetMonthlyTransferAmountReceiver(ctx context.Context, year int) ([]*models.MonthlyTransferReceivedRow, error)
	GetYearlyTransferAmountReceiver(ctx context.Context, year int) ([]*models.YearlyTransferReceivedRow, error)
}

type cardStatsBalanceRepository struct {
	db *gorm.DB
}

func NewCardStatsBalanceRepository(db *gorm.DB) CardStatsBalanceRepository {
	return &cardStatsBalanceRepository{db: db}
}

func (r *cardStatsBalanceRepository) GetMonthlyBalance(ctx context.Context, year int) ([]*models.MonthlyBalanceRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	var res []*models.MonthlyBalanceRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(created_at, 'Mon') AS month, COALESCE(SUM(total_balance), 0)::int AS total_balance
		FROM saldos WHERE created_at >= ? AND created_at < ?
		GROUP BY TO_CHAR(created_at, 'Mon'), EXTRACT(MONTH FROM created_at)
		ORDER BY EXTRACT(MONTH FROM created_at)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetMonthlyBalanceFailed.WithInternal(err)
	}
	return res, nil
}

func (r *cardStatsBalanceRepository) GetYearlyBalance(ctx context.Context, year int) ([]*models.YearlyBalanceRow, error) {
	var res []*models.YearlyBalanceRow
	yearStart := time.Date(year-4, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM created_at)::text AS year, COALESCE(SUM(total_balance), 0)::bigint AS total_balance
		FROM saldos WHERE created_at >= ? AND created_at < ?
		GROUP BY EXTRACT(YEAR FROM created_at) ORDER BY EXTRACT(YEAR FROM created_at)
	`, yearStart, yearEnd).Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetYearlyBalanceFailed.WithInternal(err)
	}
	return res, nil
}
