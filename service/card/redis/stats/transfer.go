package cardstatsmencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type cardStatsTransferCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsTransferCache(store *sharedcachehelpers.CacheStore) CardStatsTransferCache {
	return &cardStatsTransferCache{store: store}
}

func (c *cardStatsTransferCache) GetMonthlyTransferSenderCache(ctx context.Context, year int) ([]*models.MonthlyTransferSentRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTransferSender, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlyTransferSentRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferCache) SetMonthlyTransferSenderCache(ctx context.Context, year int, data []*models.MonthlyTransferSentRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTransferSender, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsTransferCache) GetYearlyTransferSenderCache(ctx context.Context, year int) ([]*models.YearlyTransferSentRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTransferSender, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlyTransferSentRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferCache) SetYearlyTransferSenderCache(ctx context.Context, year int, data []*models.YearlyTransferSentRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTransferSender, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsTransferCache) GetMonthlyTransferReceiverCache(ctx context.Context, year int) ([]*models.MonthlyTransferReceivedRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTransferReceiver, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlyTransferReceivedRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferCache) SetMonthlyTransferReceiverCache(ctx context.Context, year int, data []*models.MonthlyTransferReceivedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTransferReceiver, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsTransferCache) GetYearlyTransferReceiverCache(ctx context.Context, year int) ([]*models.YearlyTransferReceivedRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTransferReceiver, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlyTransferReceivedRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferCache) SetYearlyTransferReceiverCache(ctx context.Context, year int, data []*models.YearlyTransferReceivedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTransferReceiver, year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}
