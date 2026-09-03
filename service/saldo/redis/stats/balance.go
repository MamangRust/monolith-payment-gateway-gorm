package saldostatscache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type saldoStatsBalanceCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewSaldoStatsBalanceCache(store *sharedcachehelpers.CacheStore) SaldoStatsBalanceCache {
	return &saldoStatsBalanceCache{store: store}
}

func (c *saldoStatsBalanceCache) GetMonthlySaldoBalanceCache(ctx context.Context, year int) ([]*models.MonthlySaldoBalanceRow, bool) {
	key := fmt.Sprintf(saldoMonthBalanceCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlySaldoBalanceRow](ctx, c.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *saldoStatsBalanceCache) SetMonthlySaldoBalanceCache(ctx context.Context, year int, data []*models.MonthlySaldoBalanceRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(saldoMonthBalanceCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *saldoStatsBalanceCache) GetYearlySaldoBalanceCache(ctx context.Context, year int) ([]*models.YearlySaldoBalanceRow, bool) {
	key := fmt.Sprintf(saldoYearlyBalanceCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlySaldoBalanceRow](ctx, c.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *saldoStatsBalanceCache) SetYearlySaldoBalanceCache(ctx context.Context, year int, data []*models.YearlySaldoBalanceRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(saldoYearlyBalanceCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}
