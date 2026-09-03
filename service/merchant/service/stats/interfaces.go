package merchantstatsservice

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
)

type MerchantStatsAmountService interface {
	FindMonthlyAmountMerchant(ctx context.Context, year int) ([]*models.MerchantMonthlyAmountRow, error)
	FindYearlyAmountMerchant(ctx context.Context, year int) ([]*models.MerchantYearlyAmountRow, error)
}

type MerchantStatsMethodService interface {
	FindMonthlyPaymentMethodsMerchant(ctx context.Context, year int) ([]*models.MerchantMonthlyPaymentMethodRow, error)
	FindYearlyPaymentMethodMerchant(ctx context.Context, year int) ([]*models.MerchantYearlyPaymentMethodRow, error)
}

type MerchantStatsTotalAmountService interface {
	FindMonthlyTotalAmountMerchant(ctx context.Context, year int) ([]*models.MerchantMonthlyTotalAmountRow, error)
	FindYearlyTotalAmountMerchant(ctx context.Context, year int) ([]*models.MerchantYearlyTotalAmountRow, error)
}
