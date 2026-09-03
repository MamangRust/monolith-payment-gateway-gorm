package repository

import (
	"gorm.io/gorm"

	withdrawstatsrepository "github.com/MamangRust/monolith-payment-gateway-withdraw/repository/stats"
	withdrawstatsbycardrepository "github.com/MamangRust/monolith-payment-gateway-withdraw/repository/statsbycard"
)

type Repositories interface {
	CardRepository
	SaldoRepository
	WithdrawQueryRepository
	WithdrawCommandRepository
	withdrawstatsrepository.WithdrawStatsRepository
	withdrawstatsbycardrepository.WithdrawStatsByCardRepository
}

type repositories struct {
	CardRepository
	SaldoRepository
	WithdrawQueryRepository
	WithdrawCommandRepository
	withdrawstatsrepository.WithdrawStatsRepository
	withdrawstatsbycardrepository.WithdrawStatsByCardRepository
}

func NewRepositories(
	db *gorm.DB,
	card CardRepository,
	saldo SaldoRepository,
	dailyLimit ...int64,
) Repositories {
	return &repositories{
		CardRepository:                card,
		SaldoRepository:               saldo,
		WithdrawQueryRepository:       NewWithdrawQueryRepository(db),
		WithdrawCommandRepository:     NewWithdrawCommandRepository(db, dailyLimit...),
		WithdrawStatsRepository:       withdrawstatsrepository.NewWithdrawStatsRepository(db),
		WithdrawStatsByCardRepository: withdrawstatsbycardrepository.NewWithdrawStatsByCardRepository(db),
	}
}
