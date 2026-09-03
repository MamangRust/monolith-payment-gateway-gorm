package withdrawstatsbycardcache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type withdrawStatsByCardStatusCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewWithdrawStatsStatusCache(store *sharedcachehelpers.CacheStore) WithdrawStatsByCardStatusCache {
	return &withdrawStatsByCardStatusCache{store: store}
}

func (w *withdrawStatsByCardStatusCache) GetCachedMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*models.WithdrawMonthlyStatusSuccessByCardRow, bool) {
	key := fmt.Sprintf(monthWithdrawStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*models.WithdrawMonthlyStatusSuccessByCardRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsByCardStatusCache) SetCachedMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber, data []*models.WithdrawMonthlyStatusSuccessByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthWithdrawStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsByCardStatusCache) GetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*models.WithdrawYearlyStatusSuccessByCardRow, bool) {
	key := fmt.Sprintf(yearWithdrawStatusSuccessByCardKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[[]*models.WithdrawYearlyStatusSuccessByCardRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsByCardStatusCache) SetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber, data []*models.WithdrawYearlyStatusSuccessByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawStatusSuccessByCardKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsByCardStatusCache) GetCachedMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*models.WithdrawMonthlyStatusFailedByCardRow, bool) {
	key := fmt.Sprintf(monthWithdrawStatusFailedByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*models.WithdrawMonthlyStatusFailedByCardRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsByCardStatusCache) SetCachedMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber, data []*models.WithdrawMonthlyStatusFailedByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthWithdrawStatusFailedByCardKey, req.CardNumber, req.Month, req.Year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsByCardStatusCache) GetCachedYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*models.WithdrawYearlyStatusFailedByCardRow, bool) {
	key := fmt.Sprintf(yearWithdrawStatusFailedByCardKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[[]*models.WithdrawYearlyStatusFailedByCardRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsByCardStatusCache) SetCachedYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber, data []*models.WithdrawYearlyStatusFailedByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawStatusFailedByCardKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}
