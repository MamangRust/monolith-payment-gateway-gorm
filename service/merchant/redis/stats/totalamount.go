package merchantstatscache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type merchantStatsTotalAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsTotalAmountCache(store *sharedcachehelpers.CacheStore) MerchantStatsTotalAmountCache {
	return &merchantStatsTotalAmountCache{store: store}
}

func (s *merchantStatsTotalAmountCache) GetMonthlyTotalAmountMerchantCache(ctx context.Context, year int) ([]*models.MerchantMonthlyTotalAmountRow, bool) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantMonthlyTotalAmountRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantStatsTotalAmountCache) SetMonthlyTotalAmountMerchantCache(ctx context.Context, year int, data []*models.MerchantMonthlyTotalAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyTotalAmountCacheKey, year)

	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *merchantStatsTotalAmountCache) GetYearlyTotalAmountMerchantCache(ctx context.Context, year int) ([]*models.MerchantYearlyTotalAmountRow, bool) {
	key := fmt.Sprintf(merchantYearlyTotalAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantYearlyTotalAmountRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantStatsTotalAmountCache) SetYearlyTotalAmountMerchantCache(ctx context.Context, year int, data []*models.MerchantYearlyTotalAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyTotalAmountCacheKey, year)

	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
