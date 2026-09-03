package withdrawstatsbycardservice

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type WithdrawStatsByCardStatusService interface {
	FindMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*models.WithdrawMonthlyStatusSuccessByCardRow, error)
	FindYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*models.WithdrawYearlyStatusSuccessByCardRow, error)
	FindMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*models.WithdrawMonthlyStatusFailedByCardRow, error)
	FindYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*models.WithdrawYearlyStatusFailedByCardRow, error)
}

type WithdrawStatsByCardAmountService interface {
	FindMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*models.WithdrawMonthlyAmountByCardRow, error)
	FindYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*models.WithdrawYearlyAmountByCardRow, error)
}
