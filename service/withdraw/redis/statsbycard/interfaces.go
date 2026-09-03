package withdrawstatsbycardcache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type WithdrawStatsByCardStatusCache interface {
	GetCachedMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*models.WithdrawMonthlyStatusSuccessByCardRow, bool)
	SetCachedMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber, data []*models.WithdrawMonthlyStatusSuccessByCardRow)

	GetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*models.WithdrawYearlyStatusSuccessByCardRow, bool)
	SetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber, data []*models.WithdrawYearlyStatusSuccessByCardRow)

	GetCachedMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*models.WithdrawMonthlyStatusFailedByCardRow, bool)
	SetCachedMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber, data []*models.WithdrawMonthlyStatusFailedByCardRow)

	GetCachedYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*models.WithdrawYearlyStatusFailedByCardRow, bool)
	SetCachedYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber, data []*models.WithdrawYearlyStatusFailedByCardRow)
}

type WithdrawStatsByCardAmountCache interface {
	GetCachedMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*models.WithdrawMonthlyAmountByCardRow, bool)
	SetCachedMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber, data []*models.WithdrawMonthlyAmountByCardRow)

	GetCachedYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*models.WithdrawYearlyAmountByCardRow, bool)
	SetCachedYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber, data []*models.WithdrawYearlyAmountByCardRow)
}
