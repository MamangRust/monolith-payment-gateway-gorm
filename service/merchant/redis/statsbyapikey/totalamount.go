package merchantstatsapikey

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type merchantStatsTotalAmountByApiKeyCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsTotalAmountByApiKeyCache(store *sharedcachehelpers.CacheStore) MerchantStatsTotalAmountByApiKeyCache {
	return &merchantStatsTotalAmountByApiKeyCache{store: store}
}

func (m *merchantStatsTotalAmountByApiKeyCache) GetMonthlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*models.MerchantMonthlyTotalAmountRow, bool) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantMonthlyTotalAmountRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsTotalAmountByApiKeyCache) SetMonthlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey, data []*models.MerchantMonthlyTotalAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *merchantStatsTotalAmountByApiKeyCache) GetYearlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*models.MerchantYearlyTotalAmountRow, bool) {
	key := fmt.Sprintf(merchantYearlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantYearlyTotalAmountRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsTotalAmountByApiKeyCache) SetYearlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey, data []*models.MerchantYearlyTotalAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}
