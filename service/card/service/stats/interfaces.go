package cardstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
)

// CardStatsBalanceService handles balance statistics globally and per specific card number.
type CardStatsBalanceService interface {
	FindMonthlyBalance(ctx context.Context, year int) ([]*models.MonthlyBalanceRow, error)
	FindYearlyBalance(ctx context.Context, year int) ([]*models.YearlyBalanceRow, error)
}

// CardStatsTopupService handles top-up statistics globally and per specific card number.
type CardStatsTopupService interface {
	FindMonthlyTopupAmount(ctx context.Context, year int) ([]*models.MonthlyTopupRow, error)
	FindYearlyTopupAmount(ctx context.Context, year int) ([]*models.YearlyTopupRow, error)
}

// CardStatsWithdrawService handles withdraw statistics globally and per specific card number.
type CardStatsWithdrawService interface {
	FindMonthlyWithdrawAmount(ctx context.Context, year int) ([]*models.MonthlyWithdrawRow, error)
	FindYearlyWithdrawAmount(ctx context.Context, year int) ([]*models.YearlyWithdrawRow, error)
}

// CardStatsTransactionService handles transaction statistics globally and per specific card number.
type CardStatsTransactionService interface {
	FindMonthlyTransactionAmount(ctx context.Context, year int) ([]*models.MonthlyTransactionRow, error)
	FindYearlyTransactionAmount(ctx context.Context, year int) ([]*models.YearlyTransactionRow, error)
}

// CardStatsTransferService handles transfer statistics globally and per specific card number (as sender or receiver).
type CardStatsTransferService interface {
	FindMonthlyTransferAmountSender(ctx context.Context, year int) ([]*models.MonthlyTransferSentRow, error)
	FindYearlyTransferAmountSender(ctx context.Context, year int) ([]*models.YearlyTransferSentRow, error)
	FindMonthlyTransferAmountReceiver(ctx context.Context, year int) ([]*models.MonthlyTransferReceivedRow, error)
	FindYearlyTransferAmountReceiver(ctx context.Context, year int) ([]*models.YearlyTransferReceivedRow, error)
}
