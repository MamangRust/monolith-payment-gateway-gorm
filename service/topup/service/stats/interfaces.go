package topupstatsservice

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TopupStatsAmountService interface {
	FindMonthlyTopupAmounts(ctx context.Context, year int) ([]*models.TopupMonthlyAmountRow, error)
	FindYearlyTopupAmounts(ctx context.Context, year int) ([]*models.TopupYearlyAmountRow, error)
}

type TopupStatsMethodService interface {
	FindMonthlyTopupMethods(ctx context.Context, year int) ([]*models.TopupMonthlyMethodRow, error)
	FindYearlyTopupMethods(ctx context.Context, year int) ([]*models.TopupYearlyMethodRow, error)
}

type TopupStatsStatusService interface {
	FindMonthTopupStatusSuccess(ctx context.Context, req *requests.MonthTopupStatus) ([]*models.TopupMonthlyStatusRow, error)
	FindYearlyTopupStatusSuccess(ctx context.Context, year int) ([]*models.TopupYearlyStatusRow, error)
	FindMonthTopupStatusFailed(ctx context.Context, req *requests.MonthTopupStatus) ([]*models.TopupMonthlyStatusRow, error)
	FindYearlyTopupStatusFailed(ctx context.Context, year int) ([]*models.TopupYearlyStatusRow, error)
}
