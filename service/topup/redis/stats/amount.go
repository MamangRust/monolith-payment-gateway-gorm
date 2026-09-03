package topupstatscache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type topupStatsAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTopupStatsAmountCache(store *sharedcachehelpers.CacheStore) TopupStatsAmountCache {
	return &topupStatsAmountCache{store: store}
}

func (c *topupStatsAmountCache) GetMonthlyTopupAmountsCache(ctx context.Context, year int) ([]*models.TopupMonthlyAmountRow, bool) {
	key := fmt.Sprintf(monthTopupAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TopupMonthlyAmountRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupStatsAmountCache) SetMonthlyTopupAmountsCache(ctx context.Context, year int, data []*models.TopupMonthlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupAmountCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *topupStatsAmountCache) GetYearlyTopupAmountsCache(ctx context.Context, year int) ([]*models.TopupYearlyAmountRow, bool) {
	key := fmt.Sprintf(yearTopupAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TopupYearlyAmountRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupStatsAmountCache) SetYearlyTopupAmountsCache(ctx context.Context, year int, data []*models.TopupYearlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupAmountCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}
