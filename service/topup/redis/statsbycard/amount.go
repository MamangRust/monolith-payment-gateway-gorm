package topupstatsbycardcache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type topupStatsAmountByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTopupStatsAmountByCardCache(store *sharedcachehelpers.CacheStore) TopupStatsAmountByCardCache {
	return &topupStatsAmountByCardCache{store: store}
}

func (s *topupStatsAmountByCardCache) GetMonthlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupMonthlyAmountRow, bool) {
	key := fmt.Sprintf(monthTopupAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TopupMonthlyAmountRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsAmountByCardCache) SetMonthlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*models.TopupMonthlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupAmountByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *topupStatsAmountByCardCache) GetYearlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupYearlyAmountRow, bool) {
	key := fmt.Sprintf(yearTopupAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TopupYearlyAmountRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsAmountByCardCache) SetYearlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*models.TopupYearlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupAmountByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
