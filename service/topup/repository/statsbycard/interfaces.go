package topupstatsbycardrepository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TopupStatsByCardAmountRepository interface {
	GetMonthlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupMonthlyAmountRow, error)
	GetYearlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupYearlyAmountRow, error)
}

type TopupStatsByCardStatusRepository interface {
	GetMonthTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*models.TopupMonthlyStatusRow, error)
	GetYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*models.TopupYearlyStatusRow, error)

	GetMonthTopupStatusFailedByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*models.TopupMonthlyStatusRow, error)
	GetYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*models.TopupYearlyStatusRow, error)
}

type TopupStatsByCardMethodRepository interface {
	GetMonthlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupMonthlyMethodByCardRow, error)
	GetYearlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupYearlyMethodByCardRow, error)
}
