package topupstatscache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TopupStatsStatusCache interface {
	GetMonthTopupStatusSuccessCache(ctx context.Context, req *requests.MonthTopupStatus) ([]*models.TopupMonthlyStatusRow, bool)
	SetMonthTopupStatusSuccessCache(ctx context.Context, req *requests.MonthTopupStatus, data []*models.TopupMonthlyStatusRow)

	GetYearlyTopupStatusSuccessCache(ctx context.Context, year int) ([]*models.TopupYearlyStatusRow, bool)
	SetYearlyTopupStatusSuccessCache(ctx context.Context, year int, data []*models.TopupYearlyStatusRow)

	GetMonthTopupStatusFailedCache(ctx context.Context, req *requests.MonthTopupStatus) ([]*models.TopupMonthlyStatusRow, bool)
	SetMonthTopupStatusFailedCache(ctx context.Context, req *requests.MonthTopupStatus, data []*models.TopupMonthlyStatusRow)

	GetYearlyTopupStatusFailedCache(ctx context.Context, year int) ([]*models.TopupYearlyStatusRow, bool)
	SetYearlyTopupStatusFailedCache(ctx context.Context, year int, data []*models.TopupYearlyStatusRow)
}

type TopupStatsMethodCache interface {
	GetMonthlyTopupMethodsCache(ctx context.Context, year int) ([]*models.TopupMonthlyMethodRow, bool)
	SetMonthlyTopupMethodsCache(ctx context.Context, year int, data []*models.TopupMonthlyMethodRow)

	GetYearlyTopupMethodsCache(ctx context.Context, year int) ([]*models.TopupYearlyMethodRow, bool)
	SetYearlyTopupMethodsCache(ctx context.Context, year int, data []*models.TopupYearlyMethodRow)
}

type TopupStatsAmountCache interface {
	GetMonthlyTopupAmountsCache(ctx context.Context, year int) ([]*models.TopupMonthlyAmountRow, bool)
	SetMonthlyTopupAmountsCache(ctx context.Context, year int, data []*models.TopupMonthlyAmountRow)

	GetYearlyTopupAmountsCache(ctx context.Context, year int) ([]*models.TopupYearlyAmountRow, bool)
	SetYearlyTopupAmountsCache(ctx context.Context, year int, data []*models.TopupYearlyAmountRow)
}
