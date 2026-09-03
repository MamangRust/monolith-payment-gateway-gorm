package topupstatsrepository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TopupStatsAmountRepository interface {
	GetMonthlyTopupAmounts(ctx context.Context, year int) ([]*models.TopupMonthlyAmountRow, error)
	GetYearlyTopupAmounts(ctx context.Context, year int) ([]*models.TopupYearlyAmountRow, error)
}

type TopupStatsStatusRepository interface {
	GetMonthTopupStatusSuccess(ctx context.Context, req *requests.MonthTopupStatus) ([]*models.TopupMonthlyStatusRow, error)
	GetYearlyTopupStatusSuccess(ctx context.Context, year int) ([]*models.TopupYearlyStatusRow, error)

	GetMonthTopupStatusFailed(ctx context.Context, req *requests.MonthTopupStatus) ([]*models.TopupMonthlyStatusRow, error)
	GetYearlyTopupStatusFailed(ctx context.Context, year int) ([]*models.TopupYearlyStatusRow, error)
}

type TopupStatsMethodRepository interface {
	GetMonthlyTopupMethods(ctx context.Context, year int) ([]*models.TopupMonthlyMethodRow, error)
	GetYearlyTopupMethods(ctx context.Context, year int) ([]*models.TopupYearlyMethodRow, error)
}
