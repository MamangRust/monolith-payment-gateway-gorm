package saldostatsrepository

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"gorm.io/gorm"
)

type SaldoStatsRepository interface {
	SaldoStatsBalanceRepository
	SaldoStatsTotalSaldoRepository
}

type repository struct {
	SaldoStatsBalanceRepository
	SaldoStatsTotalSaldoRepository
}

func NewSaldoStatsRepository(db *gorm.DB) SaldoStatsRepository {
	return &repository{
		SaldoStatsBalanceRepository:    NewSaldoStatsBalanceRepository(db),
		SaldoStatsTotalSaldoRepository: NewSaldoStatsTotalBalanceRepository(db),
	}
}

func toTotalCount(res []*models.MonthlySaldoBalanceRow) int {
	if len(res) > 0 {
		return len(res)
	}
	return 0
}
