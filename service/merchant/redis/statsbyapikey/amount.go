package merchantstatsapikey

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type merchantStatsAmountByApiKeyCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsAmountByApiKeyCache(store *sharedcachehelpers.CacheStore) MerchantStatsAmountByApiKeyCache {
	return &merchantStatsAmountByApiKeyCache{store: store}
}

func (m *merchantStatsAmountByApiKeyCache) GetMonthlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*models.MerchantMonthlyAmountRow, bool) {
	key := fmt.Sprintf(merchantMonthlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantMonthlyAmountRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsAmountByApiKeyCache) SetMonthlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey, data []*models.MerchantMonthlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *merchantStatsAmountByApiKeyCache) GetYearlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*models.MerchantYearlyAmountRow, bool) {
	key := fmt.Sprintf(merchantYearlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantYearlyAmountRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsAmountByApiKeyCache) SetYearlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey, data []*models.MerchantYearlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}
