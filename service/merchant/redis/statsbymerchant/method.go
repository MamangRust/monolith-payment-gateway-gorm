package merchantstatsbymerchant

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type merchantStatsMethodByMerchant struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantStatsMethodByMerchantCache(store *sharedcachehelpers.CacheStore) MerchantStatsMethodByMerchantCache {
	return &merchantStatsMethodByMerchant{store: store}
}

func (m *merchantStatsMethodByMerchant) GetMonthlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*models.MerchantMonthlyPaymentMethodRow, bool) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantMonthlyPaymentMethodRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsMethodByMerchant) SetMonthlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant, data []*models.MerchantMonthlyPaymentMethodRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *merchantStatsMethodByMerchant) GetYearlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*models.MerchantYearlyPaymentMethodRow, bool) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantYearlyPaymentMethodRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsMethodByMerchant) SetYearlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant, data []*models.MerchantYearlyPaymentMethodRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}
