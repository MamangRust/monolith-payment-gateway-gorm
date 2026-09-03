package transactionstatsbycarcache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type transactionStatsByCardStatusCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransactionStatsByCardStatusCache(store *sharedcachehelpers.CacheStore) TransactionStatsByCardStatusCache {
	return &transactionStatsByCardStatusCache{store: store}
}

func (t *transactionStatsByCardStatusCache) GetMonthTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*models.TransactionMonthlyStatusSuccessByCardRow, bool) {
	key := fmt.Sprintf(monthTransactionStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransactionMonthlyStatusSuccessByCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByCardStatusCache) SetMonthTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber, data []*models.TransactionMonthlyStatusSuccessByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsByCardStatusCache) GetYearTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*models.TransactionYearlyStatusSuccessByCardRow, bool) {
	key := fmt.Sprintf(yearTransactionStatusSuccessByCardCacheKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransactionYearlyStatusSuccessByCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByCardStatusCache) SetYearTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber, data []*models.TransactionYearlyStatusSuccessByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTransactionStatusSuccessByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsByCardStatusCache) GetMonthTransactionStatusFailedByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*models.TransactionMonthlyStatusFailedByCardRow, bool) {
	key := fmt.Sprintf(monthTransactionStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransactionMonthlyStatusFailedByCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByCardStatusCache) SetMonthTransactionStatusFailedByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber, data []*models.TransactionMonthlyStatusFailedByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsByCardStatusCache) GetYearTransactionStatusFailedByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*models.TransactionYearlyStatusFailedByCardRow, bool) {
	key := fmt.Sprintf(yearTransactionStatusFailedByCardCacheKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransactionYearlyStatusFailedByCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByCardStatusCache) SetYearTransactionStatusFailedByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber, data []*models.TransactionYearlyStatusFailedByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTransactionStatusFailedByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
