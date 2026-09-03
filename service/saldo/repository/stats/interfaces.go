package saldostatsrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoStatsBalanceRepository interface {
	GetMonthlySaldoBalances(ctx context.Context, year int) ([]*models.MonthlySaldoBalanceRow, error)
	GetYearlySaldoBalances(ctx context.Context, year int) ([]*models.YearlySaldoBalanceRow, error)
}

type SaldoStatsTotalSaldoRepository interface {
	GetMonthlyTotalSaldoBalance(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*models.MonthlyTotalSaldoBalanceRow, error)
	GetYearTotalSaldoBalance(ctx context.Context, year int) ([]*models.YearlyTotalSaldoBalancesRow, error)
}
