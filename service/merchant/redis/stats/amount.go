package merchantstatscache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type merchantStatsAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsAmountCache(store *sharedcachehelpers.CacheStore) MerchantStatsAmountCache {
	return &merchantStatsAmountCache{store: store}
}

func (s *merchantStatsAmountCache) GetMonthlyAmountMerchantCache(ctx context.Context, year int) ([]*models.MerchantMonthlyAmountRow, bool) {
	key := fmt.Sprintf(merchantMonthlyAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantMonthlyAmountRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantStatsAmountCache) SetMonthlyAmountMerchantCache(ctx context.Context, year int, data []*models.MerchantMonthlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyAmountCacheKey, year)

	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *merchantStatsAmountCache) GetYearlyAmountMerchantCache(ctx context.Context, year int) ([]*models.MerchantYearlyAmountRow, bool) {
	key := fmt.Sprintf(MerchantYearlyAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantYearlyAmountRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantStatsAmountCache) SetYearlyAmountMerchantCache(ctx context.Context, year int, data []*models.MerchantYearlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(MerchantYearlyAmountCacheKey, year)

	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
