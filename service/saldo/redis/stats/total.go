package saldostatscache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type saldoStatsTotalCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewSaldoStatsTotalCache(store *sharedcachehelpers.CacheStore) SaldoStatsTotalCache {
	return &saldoStatsTotalCache{store: store}
}

func (c *saldoStatsTotalCache) GetMonthlyTotalSaldoBalanceCache(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*models.MonthlyTotalSaldoBalanceRow, bool) {
	key := fmt.Sprintf(saldoMonthTotalBalanceCacheKey, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlyTotalSaldoBalanceRow](ctx, c.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *saldoStatsTotalCache) SetMonthlyTotalSaldoCache(ctx context.Context, req *requests.MonthTotalSaldoBalance, data []*models.MonthlyTotalSaldoBalanceRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(saldoMonthTotalBalanceCacheKey, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *saldoStatsTotalCache) GetYearTotalSaldoBalanceCache(ctx context.Context, year int) ([]*models.YearlyTotalSaldoBalancesRow, bool) {
	key := fmt.Sprintf(saldoYearTotalBalanceCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlyTotalSaldoBalancesRow](ctx, c.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *saldoStatsTotalCache) SetYearTotalSaldoBalanceCache(ctx context.Context, year int, data []*models.YearlyTotalSaldoBalancesRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(saldoYearTotalBalanceCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}
