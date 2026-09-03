package cardstatsbycardmencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type cardStatsWithdrawByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsWithdrawByCardCache(store *sharedcachehelpers.CacheStore) CardStatsWithdrawByCardCache {
	return &cardStatsWithdrawByCardCache{store: store}
}

func (c *cardStatsWithdrawByCardCache) GetMonthlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyWithdrawRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyWithdrawByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlyWithdrawRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsWithdrawByCardCache) SetMonthlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.MonthlyWithdrawRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyWithdrawByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatsWithdrawByCardCache) GetYearlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyWithdrawRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyWithdrawByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlyWithdrawRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsWithdrawByCardCache) SetYearlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.YearlyWithdrawRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyWithdrawByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}
