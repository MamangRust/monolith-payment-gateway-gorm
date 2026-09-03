package cardstatsbycardmencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type cardStatsTransferByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardStatsTransferByCardCache(store *sharedcachehelpers.CacheStore) CardStatsTransferByCardCache {
	return &cardStatsTransferByCardCache{store: store}
}

func (c *cardStatsTransferByCardCache) GetMonthlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransferSentRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlySenderByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlyTransferSentRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferByCardCache) SetMonthlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.MonthlyTransferSentRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlySenderByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatsTransferByCardCache) GetYearlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransferSentRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlySenderByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlyTransferSentRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferByCardCache) SetYearlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.YearlyTransferSentRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlySenderByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatsTransferByCardCache) GetMonthlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.MonthlyTransferReceivedRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyReceiverByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MonthlyTransferReceivedRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferByCardCache) SetMonthlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.MonthlyTransferReceivedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyReceiverByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatsTransferByCardCache) GetYearlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*models.YearlyTransferReceivedRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyReceiverByCard, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.YearlyTransferReceivedRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferByCardCache) SetYearlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*models.YearlyTransferReceivedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyReceiverByCard, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}
