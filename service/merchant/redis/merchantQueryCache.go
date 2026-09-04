package mencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type merchantCachedResponseAll struct {
	Data         []*models.MerchantListRow `json:"data"`
	TotalRecords *int                      `json:"total_records"`
}

type merchantCachedResponseActive struct {
	Data         []*models.MerchantListWithDeletedRow `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantCachedResponseTrashed struct {
	Data         []*models.MerchantListWithDeletedRow `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantQueryCache(store *sharedcachehelpers.CacheStore) MerchantQueryCache {
	return &merchantQueryCache{store: store}
}

func (m *merchantQueryCache) GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListRow, *int, bool) {
	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[merchantCachedResponseAll](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantQueryCache) SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants, data []*models.MerchantListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.MerchantListRow{}
	}

	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantCachedResponseAll{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListWithDeletedRow, *int, bool) {
	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[merchantCachedResponseActive](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantQueryCache) SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants, data []*models.MerchantListWithDeletedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.MerchantListWithDeletedRow{}
	}

	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantCachedResponseActive{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListWithDeletedRow, *int, bool) {
	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[merchantCachedResponseTrashed](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantQueryCache) SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants, data []*models.MerchantListWithDeletedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.MerchantListWithDeletedRow{}
	}

	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantCachedResponseTrashed{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchant(ctx context.Context, id int) (*models.MerchantAllFieldsRow, bool) {
	key := fmt.Sprintf(merchantByIdCacheKey, id)

	result, found := sharedcachehelpers.GetFromCache[*models.MerchantAllFieldsRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantQueryCache) SetCachedMerchant(ctx context.Context, data *models.MerchantAllFieldsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantByIdCacheKey, data.MerchantID)

	sharedcachehelpers.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantsByUserId(ctx context.Context, userId int) ([]*models.MerchantListByUserRow, bool) {
	key := fmt.Sprintf(merchantByUserIdCacheKey, userId)

	result, found := sharedcachehelpers.GetFromCache[[]*models.MerchantListByUserRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantQueryCache) SetCachedMerchantsByUserId(ctx context.Context, userId int, data []*models.MerchantListByUserRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantByUserIdCacheKey, userId)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantByApiKey(ctx context.Context, apiKey string) (*models.MerchantAllFieldsRow, bool) {
	key := fmt.Sprintf(merchantByApiKeyCacheKey, apiKey)

	result, found := sharedcachehelpers.GetFromCache[*models.MerchantAllFieldsRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantQueryCache) SetCachedMerchantByApiKey(ctx context.Context, apiKey string, data *models.MerchantAllFieldsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantByApiKeyCacheKey, apiKey)

	sharedcachehelpers.SetToCache(ctx, m.store, key, data, ttlDefault)
}
