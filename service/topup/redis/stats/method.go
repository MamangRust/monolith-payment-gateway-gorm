package topupstatscache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type topupStatsMethodCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTopupStatsMethodCache(store *sharedcachehelpers.CacheStore) TopupStatsMethodCache {
	return &topupStatsMethodCache{store: store}
}

func (c *topupStatsMethodCache) GetMonthlyTopupMethodsCache(ctx context.Context, year int) ([]*models.TopupMonthlyMethodRow, bool) {
	key := fmt.Sprintf(monthTopupMethodCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TopupMonthlyMethodRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *topupStatsMethodCache) SetMonthlyTopupMethodsCache(ctx context.Context, year int, data []*models.TopupMonthlyMethodRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupMethodCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *topupStatsMethodCache) GetYearlyTopupMethodsCache(ctx context.Context, year int) ([]*models.TopupYearlyMethodRow, bool) {
	key := fmt.Sprintf(yearTopupMethodCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TopupYearlyMethodRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *topupStatsMethodCache) SetYearlyTopupMethodsCache(ctx context.Context, year int, data []*models.TopupYearlyMethodRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupMethodCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}
