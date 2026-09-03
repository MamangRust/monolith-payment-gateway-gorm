package cardstatsbycardmencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type CardStatsBalanceByCardCache interface {
	GetMonthlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyBalanceRow, bool)
	GetYearlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyBalanceRow, bool)

	SetMonthlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.MonthlyBalanceRow)
	SetYearlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.YearlyBalanceRow)
}

type CardStatsTopupByCardCache interface {
	GetMonthlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTopupRow, bool)
	GetYearlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTopupRow, bool)

	SetMonthlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.MonthlyTopupRow)
	SetYearlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.YearlyTopupRow)
}

type CardStatsWithdrawByCardCache interface {
	GetMonthlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyWithdrawRow, bool)
	GetYearlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyWithdrawRow, bool)

	SetMonthlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.MonthlyWithdrawRow)
	SetYearlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.YearlyWithdrawRow)
}

type CardStatsTransactionByCardCache interface {
	GetMonthlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransactionRow, bool)
	GetYearlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransactionRow, bool)

	SetMonthlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.MonthlyTransactionRow)
	SetYearlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.YearlyTransactionRow)
}

type CardStatsTransferByCardCache interface {
	GetMonthlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransferSentRow, bool)
	GetYearlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransferSentRow, bool)

	SetMonthlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.MonthlyTransferSentRow)
	SetYearlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.YearlyTransferSentRow)

	GetMonthlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransferReceivedRow, bool)
	GetYearlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransferReceivedRow, bool)

	SetMonthlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.MonthlyTransferReceivedRow)
	SetYearlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.YearlyTransferReceivedRow)
}
