package merchantstatsapikey

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type MerchantStatsMethodByApiKeyCache interface {
	GetMonthlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*models.MerchantMonthlyPaymentMethodRow, bool)
	SetMonthlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey, data []*models.MerchantMonthlyPaymentMethodRow)

	GetYearlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*models.MerchantYearlyPaymentMethodRow, bool)
	SetYearlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey, data []*models.MerchantYearlyPaymentMethodRow)
}

type MerchantStatsAmountByApiKeyCache interface {
	GetMonthlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*models.MerchantMonthlyAmountRow, bool)
	SetMonthlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey, data []*models.MerchantMonthlyAmountRow)

	GetYearlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*models.MerchantYearlyAmountRow, bool)
	SetYearlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey, data []*models.MerchantYearlyAmountRow)
}

type MerchantStatsTotalAmountByApiKeyCache interface {
	GetMonthlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*models.MerchantMonthlyTotalAmountRow, bool)
	SetMonthlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey, data []*models.MerchantMonthlyTotalAmountRow)

	GetYearlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*models.MerchantYearlyTotalAmountRow, bool)
	SetYearlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey, data []*models.MerchantYearlyTotalAmountRow)
}
