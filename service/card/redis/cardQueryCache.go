package mencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type cardCachedResponse struct {
	Data         []*models.CardListRow `json:"data"`
	TotalRecords *int                  `json:"total_records"`
}

type cardCachedResponseActive struct {
	Data         []*models.CardListWithDeletedRow `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

type cardCachedResponseTrashed struct {
	Data         []*models.CardListWithDeletedRow `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

type cardQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardQueryCache(store *sharedcachehelpers.CacheStore) CardQueryCache {
	return &cardQueryCache{store: store}
}

func (c *cardQueryCache) GetByIdCache(ctx context.Context, cardID int) (*models.CardAllFieldsRow, bool) {
	key := fmt.Sprintf(cardByIdCacheKey, cardID)

	result, found := sharedcachehelpers.GetFromCache[*models.CardAllFieldsRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardQueryCache) GetByUserIDCache(ctx context.Context, userID int) (*models.CardAllFieldsRow, bool) {
	key := fmt.Sprintf(cardByUserIdCacheKey, userID)

	result, found := sharedcachehelpers.GetFromCache[*models.CardAllFieldsRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardQueryCache) GetByCardNumberCache(ctx context.Context, cardNumber string) (*models.CardAllFieldsRow, bool) {
	key := fmt.Sprintf(cardByCardNumCacheKey, cardNumber)

	result, found := sharedcachehelpers.GetFromCache[*models.CardAllFieldsRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardQueryCache) GetUserCardByCardNumberCache(ctx context.Context, cardNumber string) (*models.CardByEmailRow, bool) {
	key := fmt.Sprintf(cardUserByCardNumCacheKey, cardNumber)

	result, found := sharedcachehelpers.GetFromCache[*models.CardByEmailRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardQueryCache) GetFindAllCache(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListRow, *int, bool) {
	key := fmt.Sprintf(cardAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[cardCachedResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *cardQueryCache) GetByActiveCache(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListWithDeletedRow, *int, bool) {
	key := fmt.Sprintf(cardActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[cardCachedResponseActive](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *cardQueryCache) GetByTrashedCache(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListWithDeletedRow, *int, bool) {
	key := fmt.Sprintf(cardTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[cardCachedResponseTrashed](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *cardQueryCache) SetByIdCache(ctx context.Context, cardID int, data *models.CardAllFieldsRow) {
	key := fmt.Sprintf(cardByIdCacheKey, cardID)
	sharedcachehelpers.SetToCache(ctx, c.store, key, data, ttlDefault)
}

func (c *cardQueryCache) SetByUserIDCache(ctx context.Context, userID int, data *models.CardAllFieldsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cardByUserIdCacheKey, userID)
	sharedcachehelpers.SetToCache(ctx, c.store, key, data, ttlDefault)
}

func (c *cardQueryCache) SetByCardNumberCache(ctx context.Context, cardNumber string, data *models.CardAllFieldsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cardByCardNumCacheKey, cardNumber)
	sharedcachehelpers.SetToCache(ctx, c.store, key, data, ttlDefault)
}

func (c *cardQueryCache) SetFindAllCache(ctx context.Context, req *requests.FindAllCards, data []*models.CardListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.CardListRow{}
	}

	payload := &cardCachedResponse{Data: data, TotalRecords: total}

	key := fmt.Sprintf(cardAllCacheKey, req.Page, req.PageSize, req.Search)
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *cardQueryCache) SetByActiveCache(ctx context.Context, req *requests.FindAllCards, data []*models.CardListWithDeletedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.CardListWithDeletedRow{}
	}

	payload := &cardCachedResponseActive{Data: data, TotalRecords: total}

	key := fmt.Sprintf(cardActiveCacheKey, req.Page, req.PageSize, req.Search)
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *cardQueryCache) SetByTrashedCache(ctx context.Context, req *requests.FindAllCards, data []*models.CardListWithDeletedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.CardListWithDeletedRow{}
	}

	payload := &cardCachedResponseTrashed{Data: data, TotalRecords: total}

	key := fmt.Sprintf(cardTrashedCacheKey, req.Page, req.PageSize, req.Search)
	sharedcachehelpers.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *cardQueryCache) SetUserCardByCardNumberCache(ctx context.Context, cardNumber string, data *models.CardByEmailRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cardUserByCardNumCacheKey, cardNumber)

	sharedcachehelpers.SetToCache(ctx, c.store, key, data, ttlDefault)
}

func (c *cardQueryCache) DeleteByIdCache(ctx context.Context, cardID int) {
	key := fmt.Sprintf(cardByIdCacheKey, cardID)
	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}

func (c *cardQueryCache) DeleteByUserIDCache(ctx context.Context, userID int) {
	key := fmt.Sprintf(cardByUserIdCacheKey, userID)
	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}

func (c *cardQueryCache) DeleteByCardNumberCache(ctx context.Context, cardNumber string) {
	key := fmt.Sprintf(cardByCardNumCacheKey, cardNumber)
	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)

	userCardKey := fmt.Sprintf(cardUserByCardNumCacheKey, cardNumber)
	sharedcachehelpers.DeleteFromCache(ctx, c.store, userCardKey)
}
