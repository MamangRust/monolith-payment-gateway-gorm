package merchantstatsbyapikeyservice

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type MerchantStatsByApiKeyAmountService interface {
	FindMonthlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*models.MerchantMonthlyAmountRow, error)
	FindYearlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*models.MerchantYearlyAmountRow, error)
}

type MerchantStatsByApiKeyMethodService interface {
	FindMonthlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*models.MerchantMonthlyPaymentMethodRow, error)
	FindYearlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*models.MerchantYearlyPaymentMethodRow, error)
}

type MerchantStatsByApiKeyTotalAmountService interface {
	FindMonthlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*models.MerchantMonthlyTotalAmountRow, error)
	FindYearlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*models.MerchantYearlyTotalAmountRow, error)
}
