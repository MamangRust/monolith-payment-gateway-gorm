package withdrawstatscache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type WithdrawStatsStatusCache interface {
	GetCachedMonthWithdrawStatusSuccessCache(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*models.WithdrawMonthlyStatusSuccessRow, bool)
	SetCachedMonthWithdrawStatusSuccessCache(ctx context.Context, req *requests.MonthStatusWithdraw, data []*models.WithdrawMonthlyStatusSuccessRow)

	GetCachedYearlyWithdrawStatusSuccessCache(ctx context.Context, year int) ([]*models.WithdrawYearlyStatusSuccessRow, bool)
	SetCachedYearlyWithdrawStatusSuccessCache(ctx context.Context, year int, data []*models.WithdrawYearlyStatusSuccessRow)

	GetCachedMonthWithdrawStatusFailedCache(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*models.WithdrawMonthlyStatusFailedRow, bool)
	SetCachedMonthWithdrawStatusFailedCache(ctx context.Context, req *requests.MonthStatusWithdraw, data []*models.WithdrawMonthlyStatusFailedRow)

	GetCachedYearlyWithdrawStatusFailedCache(ctx context.Context, year int) ([]*models.WithdrawYearlyStatusFailedRow, bool)
	SetCachedYearlyWithdrawStatusFailedCache(ctx context.Context, year int, data []*models.WithdrawYearlyStatusFailedRow)
}

type WithdrawStatsAmountCache interface {
	GetCachedMonthlyWithdraws(ctx context.Context, year int) ([]*models.WithdrawMonthlyAmountRow, bool)
	SetCachedMonthlyWithdraws(ctx context.Context, year int, data []*models.WithdrawMonthlyAmountRow)

	GetCachedYearlyWithdraws(ctx context.Context, year int) ([]*models.WithdrawYearlyAmountRow, bool)
	SetCachedYearlyWithdraws(ctx context.Context, year int, data []*models.WithdrawYearlyAmountRow)
}
