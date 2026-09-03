package cardstatsbycardmencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type cardStatsTransactionByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsTransactionByCardCache(store *sharedcachehelpers.CacheStore) CardStatsTransactionByCardCache {
	return &cardStatsTransactionByCardCache{store: store}
}

func (c *cardStatsTransactionByCardCache) GetMonthlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransactionRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTxnByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlyTransactionRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransactionByCardCache) SetMonthlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.MonthlyTransactionRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTxnByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatsTransactionByCardCache) GetYearlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransactionRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTxnByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlyTransactionRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransactionByCardCache) SetYearlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.YearlyTransactionRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTxnByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}
