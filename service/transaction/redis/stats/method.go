package transactionstatscache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type transactionStatsMethodCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransactionStatsMethodCache(store *sharedcachehelpers.CacheStore) TransactionStatsMethodCache {
	return &transactionStatsMethodCache{store: store}
}

func (t *transactionStatsMethodCache) GetMonthlyPaymentMethodsCache(ctx context.Context, year int) ([]*models.TransactionMonthlyPaymentMethodRow, bool) {
	key := fmt.Sprintf(monthTransactionMethodCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransactionMonthlyPaymentMethodRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsMethodCache) SetMonthlyPaymentMethodsCache(ctx context.Context, year int, data []*models.TransactionMonthlyPaymentMethodRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionMethodCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsMethodCache) GetYearlyPaymentMethodsCache(ctx context.Context, year int) ([]*models.TransactionYearlyPaymentMethodRow, bool) {
	key := fmt.Sprintf(yearTransactionMethodCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransactionYearlyPaymentMethodRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsMethodCache) SetYearlyPaymentMethodsCache(ctx context.Context, year int, data []*models.TransactionYearlyPaymentMethodRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTransactionMethodCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
