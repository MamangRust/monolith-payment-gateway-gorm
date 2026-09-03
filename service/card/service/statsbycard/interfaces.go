package cardstatsbycard

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// CardStatsBalanceByCardService provides methods for retrieving card balance statistics by card number.
type CardStatsBalanceByCardService interface {
	FindMonthlyBalancesByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyBalanceRow, error)
	FindYearlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyBalanceRow, error)
}

type CardStatsTopupByCardService interface {
	FindMonthlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTopupRow, error)
	FindYearlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTopupRow, error)
}

type CardStatsWithdrawByCardService interface {
	FindMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyWithdrawRow, error)
	FindYearlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyWithdrawRow, error)
}

type CardStatsTransactionByCardService interface {
	FindMonthlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransactionRow, error)
	FindYearlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransactionRow, error)
}

type CardStatsTransferByCardService interface {
	FindMonthlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransferSentRow, error)
	FindYearlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransferSentRow, error)
	FindMonthlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransferReceivedRow, error)
	FindYearlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransferReceivedRow, error)
}
