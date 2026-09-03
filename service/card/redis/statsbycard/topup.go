package cardstatsbycardmencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type cardStatsTopupByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsTopupByCardCache(store *sharedcachehelpers.CacheStore) CardStatsTopupByCardCache {
	return &cardStatsTopupByCardCache{store: store}
}

func (c *cardStatsTopupByCardCache) GetMonthlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTopupRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTopupByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlyTopupRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTopupByCardCache) SetMonthlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.MonthlyTopupRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTopupByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatsTopupByCardCache) GetYearlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTopupRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTopupByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlyTopupRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTopupByCardCache) SetYearlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.YearlyTopupRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTopupByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}
