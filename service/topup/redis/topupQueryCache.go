package mencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type topupCachedResponseAll struct {
	Data  []*models.TopupListRow `json:"data"`
	Total *int               `json:"total_records"`
}

type topupCachedResponseByCard struct {
	Data  []*models.TopupListRow `json:"data"`
	Total *int                           `json:"total_records"`
}

type topupCachedResponseActive struct {
	Data  []*models.TopupListRow `json:"data"`
	Total *int                     `json:"total_records"`
}

type topupCachedResponseTrashed struct {
	Data  []*models.TopupListWithDeletedRow `json:"data"`
	Total *int                      `json:"total_records"`
}

type topupQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTopupQueryCache(store *sharedcachehelpers.CacheStore) TopupQueryCache {
	return &topupQueryCache{store: store}
}

func (c *topupQueryCache) GetCachedTopupsCache(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListRow, *int, bool) {
	key := fmt.Sprintf(topupAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[topupCachedResponseAll](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *topupQueryCache) SetCachedTopupsCache(ctx context.Context, req *requests.FindAllTopups, data []*models.TopupListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.TopupListRow{}
	}

	key := fmt.Sprintf(topupAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &topupCachedResponseAll{Data: data, Total: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *topupQueryCache) GetCacheTopupByCardCache(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*models.TopupListRow, *int, bool) {
	key := fmt.Sprintf(topupByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[topupCachedResponseByCard](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *topupQueryCache) SetCacheTopupByCardCache(ctx context.Context, req *requests.FindAllTopupsByCardNumber, data []*models.TopupListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.TopupListRow{}
	}

	key := fmt.Sprintf(topupByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)
	payload := &topupCachedResponseByCard{Data: data, Total: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *topupQueryCache) GetCachedTopupActiveCache(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListRow, *int, bool) {
	key := fmt.Sprintf(topupActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[topupCachedResponseActive](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *topupQueryCache) SetCachedTopupActiveCache(ctx context.Context, req *requests.FindAllTopups, data []*models.TopupListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.TopupListRow{}
	}

	key := fmt.Sprintf(topupActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &topupCachedResponseActive{Data: data, Total: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *topupQueryCache) GetCachedTopupTrashedCache(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListWithDeletedRow, *int, bool) {
	key := fmt.Sprintf(topupTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[topupCachedResponseTrashed](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *topupQueryCache) SetCachedTopupTrashedCache(ctx context.Context, req *requests.FindAllTopups, data []*models.TopupListWithDeletedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.TopupListWithDeletedRow{}
	}

	key := fmt.Sprintf(topupTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &topupCachedResponseTrashed{Data: data, Total: total}
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *topupQueryCache) GetCachedTopupCache(ctx context.Context, id int) (*models.TopupAllFieldsRow, bool) {
	key := fmt.Sprintf(topupByIdCacheKey, id)

	result, found := sharedcachehelpers.GetFromCache[*models.TopupAllFieldsRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupQueryCache) SetCachedTopupCache(ctx context.Context, data *models.TopupAllFieldsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(topupByIdCacheKey, data.TopupID)
	sharedcachehelpers.SetToCache(ctx, c.store, key, data, ttlDefault)
}
