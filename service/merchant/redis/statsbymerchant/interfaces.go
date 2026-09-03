package merchantstatsbymerchant

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type MerchantStatsMethodByMerchantCache interface {
	GetMonthlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*models.MerchantMonthlyPaymentMethodRow, bool)
	SetMonthlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant, data []*models.MerchantMonthlyPaymentMethodRow)

	GetYearlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*models.MerchantYearlyPaymentMethodRow, bool)
	SetYearlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant, data []*models.MerchantYearlyPaymentMethodRow)
}

type MerchantStatsAmountByMerchantCache interface {
	GetMonthlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*models.MerchantMonthlyAmountRow, bool)
	SetMonthlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant, data []*models.MerchantMonthlyAmountRow)

	GetYearlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*models.MerchantYearlyAmountRow, bool)
	SetYearlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant, data []*models.MerchantYearlyAmountRow)
}

type MerchantStatsTotalAmountByMerchantCache interface {
	GetMonthlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*models.MerchantMonthlyTotalAmountRow, bool)
	SetMonthlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant, data []*models.MerchantMonthlyTotalAmountRow)

	GetYearlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*models.MerchantYearlyTotalAmountRow, bool)
	SetYearlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant, data []*models.MerchantYearlyTotalAmountRow)
}
