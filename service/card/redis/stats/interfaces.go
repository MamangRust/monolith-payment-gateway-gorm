package cardstatsmencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
)

type CardStatsBalanceCache interface {
	GetMonthlyBalanceCache(ctx context.Context, year int) ([]*models.MonthlyBalanceRow, bool)
	SetMonthlyBalanceCache(ctx context.Context, year int, data []*models.MonthlyBalanceRow)

	GetYearlyBalanceCache(ctx context.Context, year int) ([]*models.YearlyBalanceRow, bool)
	SetYearlyBalanceCache(ctx context.Context, year int, data []*models.YearlyBalanceRow)
}

type CardStatsTopupCache interface {
	GetMonthlyTopupCache(ctx context.Context, year int) ([]*models.MonthlyTopupRow, bool)
	SetMonthlyTopupCache(ctx context.Context, year int, data []*models.MonthlyTopupRow)

	GetYearlyTopupCache(ctx context.Context, year int) ([]*models.YearlyTopupRow, bool)
	SetYearlyTopupCache(ctx context.Context, year int, data []*models.YearlyTopupRow)
}

type CardStatsWithdrawCache interface {
	GetMonthlyWithdrawCache(ctx context.Context, year int) ([]*models.MonthlyWithdrawRow, bool)
	SetMonthlyWithdrawCache(ctx context.Context, year int, data []*models.MonthlyWithdrawRow)

	GetYearlyWithdrawCache(ctx context.Context, year int) ([]*models.YearlyWithdrawRow, bool)
	SetYearlyWithdrawCache(ctx context.Context, year int, data []*models.YearlyWithdrawRow)
}

type CardStatsTransactionCache interface {
	GetMonthlyTransactionCache(ctx context.Context, year int) ([]*models.MonthlyTransactionRow, bool)
	SetMonthlyTransactionCache(ctx context.Context, year int, data []*models.MonthlyTransactionRow)

	GetYearlyTransactionCache(ctx context.Context, year int) ([]*models.YearlyTransactionRow, bool)
	SetYearlyTransactionCache(ctx context.Context, year int, data []*models.YearlyTransactionRow)
}

type CardStatsTransferCache interface {
	GetMonthlyTransferSenderCache(ctx context.Context, year int) ([]*models.MonthlyTransferSentRow, bool)
	SetMonthlyTransferSenderCache(ctx context.Context, year int, data []*models.MonthlyTransferSentRow)

	GetYearlyTransferSenderCache(ctx context.Context, year int) ([]*models.YearlyTransferSentRow, bool)
	SetYearlyTransferSenderCache(ctx context.Context, year int, data []*models.YearlyTransferSentRow)

	GetMonthlyTransferReceiverCache(ctx context.Context, year int) ([]*models.MonthlyTransferReceivedRow, bool)
	SetMonthlyTransferReceiverCache(ctx context.Context, year int, data []*models.MonthlyTransferReceivedRow)

	GetYearlyTransferReceiverCache(ctx context.Context, year int) ([]*models.YearlyTransferReceivedRow, bool)
	SetYearlyTransferReceiverCache(ctx context.Context, year int, data []*models.YearlyTransferReceivedRow)
}
