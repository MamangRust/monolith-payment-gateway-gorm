package merchantstatsbymerchant

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type merchantStatsAmountByMerchant struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsAmountByMerchantCache(store *sharedcachehelpers.CacheStore) MerchantStatsAmountByMerchantCache {
	return &merchantStatsAmountByMerchant{store: store}
}

func (m *merchantStatsAmountByMerchant) GetMonthlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*models.MerchantMonthlyAmountRow, bool) {
	key := fmt.Sprintf(merchantMonthlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantMonthlyAmountRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsAmountByMerchant) SetMonthlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant, data []*models.MerchantMonthlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *merchantStatsAmountByMerchant) GetYearlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*models.MerchantYearlyAmountRow, bool) {
	key := fmt.Sprintf(merchantYearlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantYearlyAmountRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsAmountByMerchant) SetYearlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant, data []*models.MerchantYearlyAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}
