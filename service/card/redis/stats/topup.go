package cardstatsmencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type cardStatsTopupCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsTopupCache(store *sharedcachehelpers.CacheStore) CardStatsTopupCache {
	return &cardStatsTopupCache{store: store}
}

func (c *cardStatsTopupCache) GetMonthlyTopupCache(ctx context.Context, year int) ([]*models.MonthlyTopupRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTopupAmount, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlyTopupRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTopupCache) SetMonthlyTopupCache(ctx context.Context, year int, data []*models.MonthlyTopupRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTopupAmount, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsTopupCache) GetYearlyTopupCache(ctx context.Context, year int) ([]*models.YearlyTopupRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTopupAmount, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlyTopupRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTopupCache) SetYearlyTopupCache(ctx context.Context, year int, data []*models.YearlyTopupRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTopupAmount, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}
