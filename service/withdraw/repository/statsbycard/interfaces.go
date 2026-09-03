package withdrawstatsbycardrepository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type WithdrawStatsByCardRepository interface {
	WithdrawStatsByCardStatusRepository
	WithdrawStatsByCardAmountRepository
}

type WithdrawStatsByCardStatusRepository interface {
	GetMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*models.WithdrawMonthlyStatusSuccessByCardRow, error)
	GetYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*models.WithdrawYearlyStatusSuccessByCardRow, error)
	GetMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*models.WithdrawMonthlyStatusFailedByCardRow, error)
	GetYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*models.WithdrawYearlyStatusFailedByCardRow, error)
}

type WithdrawStatsByCardAmountRepository interface {
	GetMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*models.WithdrawMonthlyAmountByCardRow, error)
	GetYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*models.WithdrawYearlyAmountByCardRow, error)
}
