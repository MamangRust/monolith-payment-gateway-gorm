package topupstatsbycardservice

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TopupStatsByCardAmountService interface {
	FindMonthlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupMonthlyAmountRow, error)
	FindYearlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupYearlyAmountRow, error)
}

type TopupStatsByCardMethodService interface {
	FindMonthlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupMonthlyMethodByCardRow, error)
	FindYearlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupYearlyMethodByCardRow, error)
}

type TopupStatsByCardStatusService interface {
	FindMonthTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*models.TopupMonthlyStatusRow, error)
	FindYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*models.TopupYearlyStatusRow, error)
	FindMonthTopupStatusFailedByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*models.TopupMonthlyStatusRow, error)
	FindYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*models.TopupYearlyStatusRow, error)
}
