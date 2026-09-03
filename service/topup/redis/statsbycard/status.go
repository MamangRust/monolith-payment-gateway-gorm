package topupstatsbycardcache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type topupStatsStatusByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTopupStatsStatusByCardCache(store *sharedcachehelpers.CacheStore) TopupStatsStatusByCardCache {
	return &topupStatsStatusByCardCache{store: store}
}

func (s *topupStatsStatusByCardCache) GetMonthTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*models.TopupMonthlyStatusRow, bool) {
	key := fmt.Sprintf(monthTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TopupMonthlyStatusRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsStatusByCardCache) SetMonthTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber, data []*models.TopupMonthlyStatusRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *topupStatsStatusByCardCache) GetYearlyTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*models.TopupYearlyStatusRow, bool) {
	key := fmt.Sprintf(yearTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TopupYearlyStatusRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsStatusByCardCache) SetYearlyTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber, data []*models.TopupYearlyStatusRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *topupStatsStatusByCardCache) GetMonthTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*models.TopupMonthlyStatusRow, bool) {
	key := fmt.Sprintf(monthTopupStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TopupMonthlyStatusRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsStatusByCardCache) SetMonthTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber, data []*models.TopupMonthlyStatusRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *topupStatsStatusByCardCache) GetYearlyTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*models.TopupYearlyStatusRow, bool) {
	key := fmt.Sprintf(yearTopupStatusFailedByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TopupYearlyStatusRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsStatusByCardCache) SetYearlyTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber, data []*models.TopupYearlyStatusRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusFailedByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
