package withdrawstatscache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type withdrawStatsStatusCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewWithdrawStatsStatusCache(store *sharedcachehelpers.CacheStore) WithdrawStatsStatusCache {
	return &withdrawStatsStatusCache{store: store}
}

func (w *withdrawStatsStatusCache) GetCachedMonthWithdrawStatusSuccessCache(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*models.WithdrawMonthlyStatusSuccessRow, bool) {
	key := fmt.Sprintf(montWithdrawStatusSuccessKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*models.WithdrawMonthlyStatusSuccessRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsStatusCache) SetCachedMonthWithdrawStatusSuccessCache(ctx context.Context, req *requests.MonthStatusWithdraw, data []*models.WithdrawMonthlyStatusSuccessRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(montWithdrawStatusSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsStatusCache) GetCachedYearlyWithdrawStatusSuccessCache(ctx context.Context, year int) ([]*models.WithdrawYearlyStatusSuccessRow, bool) {
	key := fmt.Sprintf(yearWithdrawStatusSuccessKey, year)
	result, found := cache.GetFromCache[[]*models.WithdrawYearlyStatusSuccessRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsStatusCache) SetCachedYearlyWithdrawStatusSuccessCache(ctx context.Context, year int, data []*models.WithdrawYearlyStatusSuccessRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawStatusSuccessKey, year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsStatusCache) GetCachedMonthWithdrawStatusFailedCache(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*models.WithdrawMonthlyStatusFailedRow, bool) {
	key := fmt.Sprintf(montWithdrawStatusFailedKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*models.WithdrawMonthlyStatusFailedRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsStatusCache) SetCachedMonthWithdrawStatusFailedCache(ctx context.Context, req *requests.MonthStatusWithdraw, data []*models.WithdrawMonthlyStatusFailedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(montWithdrawStatusFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsStatusCache) GetCachedYearlyWithdrawStatusFailedCache(ctx context.Context, year int) ([]*models.WithdrawYearlyStatusFailedRow, bool) {
	key := fmt.Sprintf(yearWithdrawStatusFailedKey, year)
	result, found := cache.GetFromCache[[]*models.WithdrawYearlyStatusFailedRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsStatusCache) SetCachedYearlyWithdrawStatusFailedCache(ctx context.Context, year int, data []*models.WithdrawYearlyStatusFailedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawStatusFailedKey, year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}
