package cardstatsmencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type cardStatsWithdrawCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsWithdrawCache(store *sharedcachehelpers.CacheStore) CardStatsWithdrawCache {
	return &cardStatsWithdrawCache{store: store}
}

func (c *cardStatsWithdrawCache) GetMonthlyWithdrawCache(ctx context.Context, year int) ([]*models.MonthlyWithdrawRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyWithdrawAmount, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlyWithdrawRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsWithdrawCache) SetMonthlyWithdrawCache(ctx context.Context, year int, data []*models.MonthlyWithdrawRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyWithdrawAmount, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsWithdrawCache) GetYearlyWithdrawCache(ctx context.Context, year int) ([]*models.YearlyWithdrawRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyWithdrawAmount, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlyWithdrawRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsWithdrawCache) SetYearlyWithdrawCache(ctx context.Context, year int, data []*models.YearlyWithdrawRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyWithdrawAmount, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}
