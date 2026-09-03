package saldostatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoStatsTotalBalanceService interface {
	FindMonthlyTotalSaldoBalance(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*models.MonthlyTotalSaldoBalanceRow, error)
	FindYearTotalSaldoBalance(ctx context.Context, year int) ([]*models.YearlyTotalSaldoBalancesRow, error)
}

type SaldoStatsBalanceService interface {
	FindMonthlySaldoBalances(ctx context.Context, year int) ([]*models.MonthlySaldoBalanceRow, error)
	FindYearlySaldoBalances(ctx context.Context, year int) ([]*models.YearlySaldoBalanceRow, error)
}
