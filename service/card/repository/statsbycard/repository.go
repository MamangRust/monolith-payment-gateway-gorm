package repositorystatsbycard

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"gorm.io/gorm"
)

type CardStatsBalanceByCardRepository interface {
	GetMonthlyBalancesByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyBalanceRow, error)
	GetYearlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyBalanceRow, error)
}

type CardStatsTopupByCardRepository interface {
	GetMonthlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTopupRow, error)
	GetYearlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTopupRow, error)
}

type CardStatsWithdrawByCardRepository interface {
	GetMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyWithdrawRow, error)
	GetYearlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyWithdrawRow, error)
}

type CardStatsTransactionByCardRepository interface {
	GetMonthlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransactionRow, error)
	GetYearlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransactionRow, error)
}

type CardStatsTransferByCardRepository interface {
	GetMonthlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransferSentRow, error)
	GetYearlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransferSentRow, error)
	GetMonthlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransferReceivedRow, error)
	GetYearlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransferReceivedRow, error)
}

type CardStatsByCardRepository interface {
	CardStatsBalanceByCardRepository
	CardStatsTopupByCardRepository
	CardStatsTransactionByCardRepository
	CardStatsTransferByCardRepository
	CardStatsWithdrawByCardRepository
}

type repository struct {
	CardStatsBalanceByCardRepository
	CardStatsTopupByCardRepository
	CardStatsTransactionByCardRepository
	CardStatsTransferByCardRepository
	CardStatsWithdrawByCardRepository
}

func NewCardStatsByCardRepository(db *gorm.DB) CardStatsByCardRepository {
	return &repository{
		CardStatsBalanceByCardRepository:     NewCardStatsBalanceByCardRepository(db),
		CardStatsTopupByCardRepository:       NewCardStatsTopupByCardRepository(db),
		CardStatsTransactionByCardRepository: NewCardStatsTransactionByCardRepository(db),
		CardStatsTransferByCardRepository:    NewCardStatsTransferByCardRepository(db),
		CardStatsWithdrawByCardRepository:    NewCardStatsWithdrawByCardRepository(db),
	}
}
