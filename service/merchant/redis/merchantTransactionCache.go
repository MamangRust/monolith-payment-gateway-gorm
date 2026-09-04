package mencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type merchantTransactionAllResponse struct {
	Data         []*models.MerchantTransactionRow `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

type merchantTransactionByMerchantResponse struct {
	Data         []*models.MerchantTransactionRow `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

type merchantTransactionByApikeyResponse struct {
	Data         []*models.MerchantTransactionRow `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

type merchantTransactionCache struct {
	store *cache.CacheStore
}

func NewMerchantTransactionCache(store *cache.CacheStore) MerchantTransactionCache {
	return &merchantTransactionCache{store: store}
}

func (m *merchantTransactionCache) SetCacheAllMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions, data []*models.MerchantTransactionRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.MerchantTransactionRow{}
	}

	key := fmt.Sprintf(merchantTransactionsCacheKey, req.Search, req.Page, req.PageSize)
	payload := &merchantTransactionAllResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantTransactionCache) GetCacheAllMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*models.MerchantTransactionRow, *int, bool) {
	key := fmt.Sprintf(merchantTransactionsCacheKey, req.Search, req.Page, req.PageSize)

	result, found := cache.GetFromCache[merchantTransactionAllResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantTransactionCache) SetCacheMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactionsById, data []*models.MerchantTransactionRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.MerchantTransactionRow{}
	}

	key := fmt.Sprintf(merchantTransactionCacheKey, req.MerchantID, req.Search, req.Page, req.PageSize)
	payload := &merchantTransactionByMerchantResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantTransactionCache) GetCacheMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*models.MerchantTransactionRow, *int, bool) {
	key := fmt.Sprintf(merchantTransactionCacheKey, req.MerchantID, req.Search, req.Page, req.PageSize)

	result, found := cache.GetFromCache[merchantTransactionByMerchantResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantTransactionCache) SetCacheMerchantTransactionApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey, data []*models.MerchantTransactionRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.MerchantTransactionRow{}
	}

	key := fmt.Sprintf(merchantTransactionApikeyCacheKey, req.ApiKey, req.Search, req.Page, req.PageSize)
	payload := &merchantTransactionByApikeyResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantTransactionCache) GetCacheMerchantTransactionApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*models.MerchantTransactionRow, *int, bool) {
	key := fmt.Sprintf(merchantTransactionApikeyCacheKey, req.ApiKey, req.Search, req.Page, req.PageSize)

	result, found := cache.GetFromCache[merchantTransactionByApikeyResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}
