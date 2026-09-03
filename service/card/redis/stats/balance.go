package cardstatsmencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type cardStatsBalanceCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsBalanceCache(store *sharedcachehelpers.CacheStore) CardStatsBalanceCache {
	return &cardStatsBalanceCache{store: store}
}

func (c *cardStatsBalanceCache) GetMonthlyBalanceCache(ctx context.Context, year int) ([]*models.MonthlyBalanceRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyBalance, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlyBalanceRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsBalanceCache) SetMonthlyBalanceCache(ctx context.Context, year int, data []*models.MonthlyBalanceRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyBalance, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsBalanceCache) GetYearlyBalanceCache(ctx context.Context, year int) ([]*models.YearlyBalanceRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyBalance, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlyBalanceRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsBalanceCache) SetYearlyBalanceCache(ctx context.Context, year int, data []*models.YearlyBalanceRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyBalance, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}
