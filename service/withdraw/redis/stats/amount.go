package withdrawstatscache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type withdrawStatsAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewWithdrawStatsAmountCache(store *sharedcachehelpers.CacheStore) WithdrawStatsAmountCache {
	return &withdrawStatsAmountCache{store: store}
}

func (w *withdrawStatsAmountCache) GetCachedMonthlyWithdraws(ctx context.Context, year int) ([]*models.WithdrawMonthlyAmountRow, bool) {
	key := fmt.Sprintf(montWithdrawAmountKey, year)
	result, found := cache.GetFromCache[[]*models.WithdrawMonthlyAmountRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsAmountCache) SetCachedMonthlyWithdraws(ctx context.Context, year int, data []*models.WithdrawMonthlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(montWithdrawAmountKey, year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsAmountCache) GetCachedYearlyWithdraws(ctx context.Context, year int) ([]*models.WithdrawYearlyAmountRow, bool) {
	key := fmt.Sprintf(yearWithdrawAmountKey, year)
	result, found := cache.GetFromCache[[]*models.WithdrawYearlyAmountRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsAmountCache) SetCachedYearlyWithdraws(ctx context.Context, year int, data []*models.WithdrawYearlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawAmountKey, year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}
