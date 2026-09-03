package cardstatsmencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type cardStatsTransactionCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsTransactionCache(store *sharedcachehelpers.CacheStore) CardStatsTransactionCache {
	return &cardStatsTransactionCache{store: store}
}

func (c *cardStatsTransactionCache) GetMonthlyTransactionCache(ctx context.Context, year int) ([]*models.MonthlyTransactionRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTransactionAmount, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlyTransactionRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransactionCache) SetMonthlyTransactionCache(ctx context.Context, year int, data []*models.MonthlyTransactionRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTransactionAmount, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsTransactionCache) GetYearlyTransactionCache(ctx context.Context, year int) ([]*models.YearlyTransactionRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTransactionAmount, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlyTransactionRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransactionCache) SetYearlyTransactionCache(ctx context.Context, year int, data []*models.YearlyTransactionRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTransactionAmount, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}
