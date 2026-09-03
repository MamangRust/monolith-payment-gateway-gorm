package transferstatscache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type transferStatsAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransferStatsAmountCache(store *sharedcachehelpers.CacheStore) TransferStatsAmountCache {
	return &transferStatsAmountCache{store: store}
}

func (t *transferStatsAmountCache) GetCachedMonthTransferAmounts(ctx context.Context, year int) ([]*models.TransferMonthlyAmountRow, bool) {
	key := fmt.Sprintf(transferMonthTransferAmountKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferMonthlyAmountRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsAmountCache) SetCachedMonthTransferAmounts(ctx context.Context, year int, data []*models.TransferMonthlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferAmountKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsAmountCache) GetCachedYearlyTransferAmounts(ctx context.Context, year int) ([]*models.TransferYearlyAmountRow, bool) {
	key := fmt.Sprintf(transferYearTransferAmountKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferYearlyAmountRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsAmountCache) SetCachedYearlyTransferAmounts(ctx context.Context, year int, data []*models.TransferYearlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferAmountKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
