package merchantstatscache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
)

type MerchantStatsMethodCache interface {
	GetMonthlyPaymentMethodsMerchantCache(ctx context.Context, year int) ([]*models.MerchantMonthlyPaymentMethodRow, bool)
	SetMonthlyPaymentMethodsMerchantCache(ctx context.Context, year int, data []*models.MerchantMonthlyPaymentMethodRow)

	GetYearlyPaymentMethodMerchantCache(ctx context.Context, year int) ([]*models.MerchantYearlyPaymentMethodRow, bool)
	SetYearlyPaymentMethodMerchantCache(ctx context.Context, year int, data []*models.MerchantYearlyPaymentMethodRow)
}

type MerchantStatsAmountCache interface {
	GetMonthlyAmountMerchantCache(ctx context.Context, year int) ([]*models.MerchantMonthlyAmountRow, bool)
	SetMonthlyAmountMerchantCache(ctx context.Context, year int, data []*models.MerchantMonthlyAmountRow)

	GetYearlyAmountMerchantCache(ctx context.Context, year int) ([]*models.MerchantYearlyAmountRow, bool)
	SetYearlyAmountMerchantCache(ctx context.Context, year int, data []*models.MerchantYearlyAmountRow)
}

type MerchantStatsTotalAmountCache interface {
	GetMonthlyTotalAmountMerchantCache(ctx context.Context, year int) ([]*models.MerchantMonthlyTotalAmountRow, bool)
	SetMonthlyTotalAmountMerchantCache(ctx context.Context, year int, data []*models.MerchantMonthlyTotalAmountRow)

	GetYearlyTotalAmountMerchantCache(ctx context.Context, year int) ([]*models.MerchantYearlyTotalAmountRow, bool)
	SetYearlyTotalAmountMerchantCache(ctx context.Context, year int, data []*models.MerchantYearlyTotalAmountRow)
}
