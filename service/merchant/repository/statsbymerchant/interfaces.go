package merchantstatsmerchantrepository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type MerchantStatsMethodByMerchantRepository interface {
	GetMonthlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*models.MerchantMonthlyPaymentMethodRow, error)
	GetYearlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*models.MerchantYearlyPaymentMethodRow, error)
}

type MerchantStatsAmountByMerchantRepository interface {
	GetMonthlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*models.MerchantMonthlyAmountRow, error)
	GetYearlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*models.MerchantYearlyAmountRow, error)
}

type MerchantStatsTotalAmountByMerchantRepository interface {
	GetMonthlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*models.MerchantMonthlyTotalAmountRow, error)
	GetYearlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*models.MerchantYearlyTotalAmountRow, error)
}
