package merchantstatsbymerchantservice

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type MerchantStatsByMerchantAmountService interface {
	FindMonthlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*models.MerchantMonthlyAmountRow, error)
	FindYearlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*models.MerchantYearlyAmountRow, error)
}

type MerchantStatsByMerchantMethodService interface {
	FindMonthlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*models.MerchantMonthlyPaymentMethodRow, error)
	FindYearlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*models.MerchantYearlyPaymentMethodRow, error)
}

type MerchantStatsByMerchantTotalAmountService interface {
	FindMonthlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*models.MerchantMonthlyTotalAmountRow, error)
	FindYearlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*models.MerchantYearlyTotalAmountRow, error)
}
