package saldostatscache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoStatsTotalCache interface {
	GetMonthlyTotalSaldoBalanceCache(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*models.MonthlyTotalSaldoBalanceRow, bool)
	SetMonthlyTotalSaldoCache(ctx context.Context, req *requests.MonthTotalSaldoBalance, data []*models.MonthlyTotalSaldoBalanceRow)

	GetYearTotalSaldoBalanceCache(ctx context.Context, year int) ([]*models.YearlyTotalSaldoBalancesRow, bool)
	SetYearTotalSaldoBalanceCache(ctx context.Context, year int, data []*models.YearlyTotalSaldoBalancesRow)
}

type SaldoStatsBalanceCache interface {
	GetMonthlySaldoBalanceCache(ctx context.Context, year int) ([]*models.MonthlySaldoBalanceRow, bool)
	SetMonthlySaldoBalanceCache(ctx context.Context, year int, data []*models.MonthlySaldoBalanceRow)

	GetYearlySaldoBalanceCache(ctx context.Context, year int) ([]*models.YearlySaldoBalanceRow, bool)
	SetYearlySaldoBalanceCache(ctx context.Context, year int, data []*models.YearlySaldoBalanceRow)
}
