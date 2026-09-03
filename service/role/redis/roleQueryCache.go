package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type roleCachedResponseAll struct {
	Data         []*models.RoleRow `json:"data"`
	TotalRecords *int              `json:"total_records"`
}

type roleCachedResponseActive struct {
	Data         []*models.RoleActiveRow `json:"data"`
	TotalRecords *int                    `json:"total_records"`
}

type roleCachedResponseTrashed struct {
	Data         []*models.RoleTrashedRow `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type roleQueryCache struct {
	store *cache.CacheStore
}

func NewRoleQueryCache(store *cache.CacheStore) RoleQueryCache {
	return &roleQueryCache{store: store}
}

func (r *roleQueryCache) SetCachedRoles(ctx context.Context, req *requests.FindAllRoles, data []*models.RoleRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.RoleRow{}
	}

	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &roleCachedResponseAll{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoles(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleRow, *int, bool) {
	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[roleCachedResponseAll](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *roleQueryCache) SetCachedRoleById(ctx context.Context, id int, data *models.Role) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleByIdCacheKey, id)
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleById(ctx context.Context, id int) (*models.Role, bool) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	result, found := cache.GetFromCache[*models.Role](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (r *roleQueryCache) SetCachedRoleByUserId(ctx context.Context, userId int, data []*models.Role) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleByUserIdCacheKey, userId)
	cache.SetToCache(ctx, r.store, key, &data, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleByUserId(ctx context.Context, userId int) ([]*models.Role, bool) {
	key := fmt.Sprintf(roleByUserIdCacheKey, userId)

	result, found := cache.GetFromCache[[]*models.Role](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (r *roleQueryCache) SetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles, data []*models.RoleActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.RoleActiveRow{}
	}

	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &roleCachedResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleActiveRow, *int, bool) {
	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[roleCachedResponseActive](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *roleQueryCache) SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles, data []*models.RoleTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.RoleTrashedRow{}
	}

	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &roleCachedResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleTrashedRow, *int, bool) {
	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[roleCachedResponseTrashed](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}
