package withdrawstatsbycardcache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type withdrawStatsByCardAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewWithdrawStatsAmountCache(store *sharedcachehelpers.CacheStore) WithdrawStatsByCardAmountCache {
	return &withdrawStatsByCardAmountCache{store: store}
}

func (w *withdrawStatsByCardAmountCache) GetCachedMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*models.WithdrawMonthlyAmountByCardRow, bool) {
	key := fmt.Sprintf(monthWithdrawAmountByCardKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[[]*models.WithdrawMonthlyAmountByCardRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsByCardAmountCache) SetCachedMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber, data []*models.WithdrawMonthlyAmountByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthWithdrawAmountByCardKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsByCardAmountCache) GetCachedYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*models.WithdrawYearlyAmountByCardRow, bool) {
	key := fmt.Sprintf(yearWithdrawAmountByCardKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[[]*models.WithdrawYearlyAmountByCardRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsByCardAmountCache) SetCachedYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber, data []*models.WithdrawYearlyAmountByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawAmountByCardKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}
