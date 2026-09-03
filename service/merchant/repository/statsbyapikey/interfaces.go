package merchantstatsapikeyrepository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type MerchantStatsMethodByApiKeyRepository interface {
	GetMonthlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*models.MerchantMonthlyPaymentMethodRow, error)
	GetYearlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*models.MerchantYearlyPaymentMethodRow, error)
}

type MerchantStatsAmountByApiKeyRepository interface {
	GetMonthlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*models.MerchantMonthlyAmountRow, error)
	GetYearlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*models.MerchantYearlyAmountRow, error)
}

type MerchantStatsTotalAmountByApiKeyRepository interface {
	GetMonthlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*models.MerchantMonthlyTotalAmountRow, error)
	GetYearlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*models.MerchantYearlyTotalAmountRow, error)
}
