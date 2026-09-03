package withdrawstatsrepository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type WithdrawStatsStatusRepository interface {
	GetMonthWithdrawStatusSuccess(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*models.WithdrawMonthlyStatusSuccessRow, error)
	GetYearlyWithdrawStatusSuccess(ctx context.Context, year int) ([]*models.WithdrawYearlyStatusSuccessRow, error)
	GetMonthWithdrawStatusFailed(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*models.WithdrawMonthlyStatusFailedRow, error)
	GetYearlyWithdrawStatusFailed(ctx context.Context, year int) ([]*models.WithdrawYearlyStatusFailedRow, error)
}

type WithdrawStatsAmountRepository interface {
	GetMonthlyWithdraws(ctx context.Context, year int) ([]*models.WithdrawMonthlyAmountRow, error)
	GetYearlyWithdraws(ctx context.Context, year int) ([]*models.WithdrawYearlyAmountRow, error)
}
