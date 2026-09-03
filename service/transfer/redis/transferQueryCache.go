package mencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type transferCacheResponseAll struct {
	Data         []*models.TransferListRow `json:"data"`
	TotalRecords *int                  `json:"total_records"`
}

type transferCacheResponseActive struct {
	Data         []*models.TransferListRow `json:"data"`
	TotalRecords *int                        `json:"total_records"`
}

type transferCacheResponseTrashed struct {
	Data         []*models.TransferListWithDeletedRow `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

type transferQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransferQueryCache(store *sharedcachehelpers.CacheStore) TransferQueryCache {
	return &transferQueryCache{store: store}
}

func (c *transferQueryCache) GetCachedTransfersCache(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListRow, *int, bool) {
	key := fmt.Sprintf(transferAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[transferCacheResponseAll](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (c *transferQueryCache) SetCachedTransfersCache(ctx context.Context, req *requests.FindAllTransfers, data []*models.TransferListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.TransferListRow{}
	}

	key := fmt.Sprintf(transferAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transferCacheResponseAll{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListRow, *int, bool) {
	key := fmt.Sprintf(transferActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[transferCacheResponseActive](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *transferQueryCache) SetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTransfers, data []*models.TransferListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.TransferListRow{}
	}

	key := fmt.Sprintf(transferActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transferCacheResponseActive{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListWithDeletedRow, *int, bool) {
	key := fmt.Sprintf(transferTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[transferCacheResponseTrashed](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *transferQueryCache) SetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTransfers, data []*models.TransferListWithDeletedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.TransferListWithDeletedRow{}
	}

	key := fmt.Sprintf(transferTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transferCacheResponseTrashed{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferCache(ctx context.Context, id int) (*models.TransferAllFieldsRow, bool) {
	key := fmt.Sprintf(transferByIdCacheKey, id)
	result, found := sharedcachehelpers.GetFromCache[*models.TransferAllFieldsRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *transferQueryCache) SetCachedTransferCache(ctx context.Context, data *models.TransferAllFieldsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferByIdCacheKey, data.TransferID)
	sharedcachehelpers.SetToCache(ctx, c.store, key, data, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferByFrom(ctx context.Context, from string) ([]*models.TransferBySourceCardRow, bool) {
	key := fmt.Sprintf(transferByFromCacheKey, from)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferBySourceCardRow](ctx, c.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *transferQueryCache) SetCachedTransferByFrom(ctx context.Context, from string, data []*models.TransferBySourceCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferByFromCacheKey, from)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferByTo(ctx context.Context, to string) ([]*models.TransferByDestinationCardRow, bool) {
	key := fmt.Sprintf(transferByToCacheKey, to)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferByDestinationCardRow](ctx, c.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *transferQueryCache) SetCachedTransferByTo(ctx context.Context, to string, data []*models.TransferByDestinationCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferByToCacheKey, to)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}
