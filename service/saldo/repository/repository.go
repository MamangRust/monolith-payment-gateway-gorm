package repository

import (
	saldostatsrepository "github.com/MamangRust/monolith-payment-gateway-saldo/repository/stats"
	"gorm.io/gorm"
)

// Repositories is a struct containing all saldo repositories.
type Repositories interface {
	SaldoQueryRepository
	SaldoCommandRepository
	saldostatsrepository.SaldoStatsRepository
	CardRepository
}

type repositories struct {
	SaldoQueryRepository
	SaldoCommandRepository
	saldostatsrepository.SaldoStatsRepository
	CardRepository
}

func NewRepositories(db *gorm.DB) Repositories {
	return &repositories{
		SaldoQueryRepository:   NewSaldoQueryRepository(db),
		SaldoCommandRepository: NewSaldoCommandRepository(db),
		SaldoStatsRepository:   saldostatsrepository.NewSaldoStatsRepository(db),
		CardRepository:         NewCardRepository(db),
	}
}
