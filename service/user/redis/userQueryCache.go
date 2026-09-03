package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type userCacheResponseAll struct {
	Data         []*models.UserRow `json:"data"`
	TotalRecords *int              `json:"total_records"`
}

type userCacheResponseActive struct {
	Data         []*models.UserActiveRow `json:"data"`
	TotalRecords *int                    `json:"total_records"`
}

type userCacheResponseTrashed struct {
	Data         []*models.UserTrashedRow `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type userQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewUserQueryCache(store *sharedcachehelpers.CacheStore) UserQueryCache {
	return &userQueryCache{store: store}
}

func (s *userQueryCache) GetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserRow, *int, bool) {
	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcachehelpers.GetFromCache[userCacheResponseAll](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers, data []*models.UserRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.UserRow{}
	}
	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseAll{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserActiveRow, *int, bool) {
	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcachehelpers.GetFromCache[userCacheResponseActive](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers, data []*models.UserActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.UserActiveRow{}
	}
	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseActive{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserTrashedRow, *int, bool) {
	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcachehelpers.GetFromCache[userCacheResponseTrashed](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers, data []*models.UserTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.UserTrashedRow{}
	}
	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseTrashed{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserCache(ctx context.Context, id int) (*models.UserByIDRow, bool) {
	key := fmt.Sprintf(userByIdCacheKey, id)
	result, found := sharedcachehelpers.GetFromCache[*models.UserByIDRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *userQueryCache) SetCachedUserCache(ctx context.Context, data *models.UserByIDRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(userByIdCacheKey, data.UserID)
	sharedcachehelpers.SetToCache(ctx, s.store, key, data, ttlDefault)
}
