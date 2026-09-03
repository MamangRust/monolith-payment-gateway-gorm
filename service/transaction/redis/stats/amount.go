package transactionstatscache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type transactionStatsAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransactionStatsAmountCache(store *sharedcachehelpers.CacheStore) TransactionStatsAmountCache {
	return &transactionStatsAmountCache{store: store}
}

func (t *transactionStatsAmountCache) GetMonthlyAmountsCache(ctx context.Context, year int) ([]*models.TransactionMonthlyAmountRow, bool) {
	key := fmt.Sprintf(monthTransactionAmountCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransactionMonthlyAmountRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsAmountCache) SetMonthlyAmountsCache(ctx context.Context, year int, data []*models.TransactionMonthlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionAmountCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsAmountCache) GetYearlyAmountsCache(ctx context.Context, year int) ([]*models.TransactionYearlyAmountRow, bool) {
	key := fmt.Sprintf(yearTransactionAmountCacheKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TransactionYearlyAmountRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsAmountCache) SetYearlyAmountsCache(ctx context.Context, year int, data []*models.TransactionYearlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTransactionAmountCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
