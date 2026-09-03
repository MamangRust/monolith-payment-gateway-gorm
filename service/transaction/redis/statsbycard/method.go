package transactionstatsbycarcache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type transactionStatsByCardMethodCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransactionStatsByCardMethodCache(store *sharedcachehelpers.CacheStore) TransactionStatsByCardMethodCache {
	return &transactionStatsByCardMethodCache{store: store}
}

func (t *transactionStatsByCardMethodCache) GetMonthlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionMonthlyPaymentMethodByCardRow, bool) {
	key := fmt.Sprintf(monthTransactionMethodByCardCacheKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransactionMonthlyPaymentMethodByCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByCardMethodCache) SetMonthlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*models.TransactionMonthlyPaymentMethodByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionMethodByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsByCardMethodCache) GetYearlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionYearlyPaymentMethodByCardRow, bool) {
	key := fmt.Sprintf(yearTransactionMethodByCardCacheKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransactionYearlyPaymentMethodByCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByCardMethodCache) SetYearlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*models.TransactionYearlyPaymentMethodByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTransactionMethodByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
