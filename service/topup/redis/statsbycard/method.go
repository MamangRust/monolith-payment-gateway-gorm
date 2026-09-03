package topupstatsbycardcache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type topupStatsMethodByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTopupStatsMethodByCardCache(store *sharedcachehelpers.CacheStore) TopupStatsMethodByCardCache {
	return &topupStatsMethodByCardCache{store: store}
}

func (s *topupStatsMethodByCardCache) GetMonthlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupMonthlyMethodByCardRow, bool) {
	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TopupMonthlyMethodByCardRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsMethodByCardCache) SetMonthlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*models.TopupMonthlyMethodByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *topupStatsMethodByCardCache) GetYearlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupYearlyMethodByCardRow, bool) {
	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TopupYearlyMethodByCardRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsMethodByCardCache) SetYearlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*models.TopupYearlyMethodByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
