package merchantstatsrepository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
)

type MerchantStatsMethodRepository interface {
	GetMonthlyPaymentMethodsMerchant(ctx context.Context, year int) ([]*models.MerchantMonthlyPaymentMethodRow, error)
	GetYearlyPaymentMethodMerchant(ctx context.Context, year int) ([]*models.MerchantYearlyPaymentMethodRow, error)
}

type MerchantStatsAmountRepository interface {
	GetMonthlyAmountMerchant(ctx context.Context, year int) ([]*models.MerchantMonthlyAmountRow, error)
	GetYearlyAmountMerchant(ctx context.Context, year int) ([]*models.MerchantYearlyAmountRow, error)
}

type MerchantStatsTotalAmountRepository interface {
	GetMonthlyTotalAmountMerchant(ctx context.Context, year int) ([]*models.MerchantMonthlyTotalAmountRow, error)
	GetYearlyTotalAmountMerchant(ctx context.Context, year int) ([]*models.MerchantYearlyTotalAmountRow, error)
}
