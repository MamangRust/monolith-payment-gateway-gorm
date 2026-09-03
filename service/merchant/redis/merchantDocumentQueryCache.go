package mencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type merchantDocumentQueryCachedResponse struct {
	Data         []*models.MerchantDocumentListRow `json:"data"`
	TotalRecords *int                          `json:"total_records"`
}

type merchantDocumentQueryCachedResponseActive struct {
	Data         []*models.MerchantDocumentListWithDeletedRow `json:"data"`
	TotalRecords *int                                `json:"total_records"`
}

type merchantDocumentQueryCachedResponseTrashed struct {
	Data         []*models.MerchantDocumentListWithDeletedRow `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantDocumentQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewMerchantDocumentQueryCache(store *sharedcachehelpers.CacheStore) MerchantDocumentQueryCache {
	return &merchantDocumentQueryCache{store: store}
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListRow, *int, bool) {
	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[merchantDocumentQueryCachedResponse](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*models.MerchantDocumentListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.MerchantDocumentListRow{}
	}

	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantDocumentQueryCachedResponse{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListWithDeletedRow, *int, bool) {
	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[merchantDocumentQueryCachedResponseActive](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*models.MerchantDocumentListWithDeletedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.MerchantDocumentListWithDeletedRow{}
	}

	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantDocumentQueryCachedResponseActive{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListWithDeletedRow, *int, bool) {
	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[merchantDocumentQueryCachedResponseTrashed](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*models.MerchantDocumentListWithDeletedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.MerchantDocumentListWithDeletedRow{}
	}

	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantDocumentQueryCachedResponseTrashed{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocument(ctx context.Context, id int) (*models.MerchantDocumentAllFieldsRow, bool) {
	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)

	result, found := sharedcachehelpers.GetFromCache[models.MerchantDocumentAllFieldsRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocument(ctx context.Context, id int, data *models.MerchantDocumentAllFieldsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)
	sharedcachehelpers.SetToCache(ctx, s.store, key, data, ttlDefault)
}
