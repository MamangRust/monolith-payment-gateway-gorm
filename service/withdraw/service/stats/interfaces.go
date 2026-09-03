package withdrawstatsservice

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type WithdrawStatsStatusService interface {
	FindMonthWithdrawStatusSuccess(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*models.WithdrawMonthlyStatusSuccessRow, error)
	FindYearlyWithdrawStatusSuccess(ctx context.Context, year int) ([]*models.WithdrawYearlyStatusSuccessRow, error)
	FindMonthWithdrawStatusFailed(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*models.WithdrawMonthlyStatusFailedRow, error)
	FindYearlyWithdrawStatusFailed(ctx context.Context, year int) ([]*models.WithdrawYearlyStatusFailedRow, error)
}

type WithdrawStatsAmountService interface {
	FindMonthlyWithdraws(ctx context.Context, year int) ([]*models.WithdrawMonthlyAmountRow, error)
	FindYearlyWithdraws(ctx context.Context, year int) ([]*models.WithdrawYearlyAmountRow, error)
}
