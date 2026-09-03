package merchantstatsapikey

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type merchantStatsMethodByApiKeyCache struct {
	store *cache.CacheStore
}

func NewMerchantStatsMethodByApiKeyCache(store *cache.CacheStore) MerchantStatsMethodByApiKeyCache {
	return &merchantStatsMethodByApiKeyCache{store: store}
}

func (m *merchantStatsMethodByApiKeyCache) GetMonthlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*models.MerchantMonthlyPaymentMethodRow, bool) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	result, found := cache.GetFromCache[[]*models.MerchantMonthlyPaymentMethodRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsMethodByApiKeyCache) SetMonthlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey, data []*models.MerchantMonthlyPaymentMethodRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *merchantStatsMethodByApiKeyCache) GetYearlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*models.MerchantYearlyPaymentMethodRow, bool) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	result, found := cache.GetFromCache[[]*models.MerchantYearlyPaymentMethodRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsMethodByApiKeyCache) SetYearlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey, data []*models.MerchantYearlyPaymentMethodRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}
