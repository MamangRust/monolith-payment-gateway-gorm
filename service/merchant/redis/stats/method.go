package merchantstatscache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type merchantStatsMethodCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsMethodCache(store *sharedcachehelpers.CacheStore) MerchantStatsMethodCache {
	return &merchantStatsMethodCache{store: store}
}

func (s *merchantStatsMethodCache) GetMonthlyPaymentMethodsMerchantCache(ctx context.Context, year int) ([]*models.MerchantMonthlyPaymentMethodRow, bool) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantMonthlyPaymentMethodRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantStatsMethodCache) SetMonthlyPaymentMethodsMerchantCache(ctx context.Context, year int, data []*models.MerchantMonthlyPaymentMethodRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyPaymentMethodCacheKey, year)

	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *merchantStatsMethodCache) GetYearlyPaymentMethodMerchantCache(ctx context.Context, year int) ([]*models.MerchantYearlyPaymentMethodRow, bool) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantYearlyPaymentMethodRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantStatsMethodCache) SetYearlyPaymentMethodMerchantCache(ctx context.Context, year int, data []*models.MerchantYearlyPaymentMethodRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyPaymentMethodCacheKey, year)

	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
