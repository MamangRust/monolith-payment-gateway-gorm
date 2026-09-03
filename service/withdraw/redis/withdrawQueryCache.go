package mencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type withdrawCachedResponseAll struct {
	Data         []*models.WithdrawListRow `json:"data"`
	TotalRecords *int                  `json:"total_records"`
}

type withdrawCachedResponseByCard struct {
	Data         []*models.WithdrawByCardNumberRow `json:"data"`
	TotalRecords *int                              `json:"total_records"`
}

type withdrawCachedResponseActive struct {
	Data         []*models.WithdrawListRow `json:"data"`
	TotalRecords *int                        `json:"total_records"`
}

type withdrawCachedResponseTrashed struct {
	Data         []*models.WithdrawListWithDeletedRow `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

type withdrawQueryCache struct {
	store *cache.CacheStore
}

func NewWithdrawQueryCache(store *cache.CacheStore) WithdrawQueryCache {
	return &withdrawQueryCache{store: store}
}

func (w *withdrawQueryCache) GetCachedWithdrawsCache(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListRow, *int, bool) {
	key := fmt.Sprintf(withdrawAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[withdrawCachedResponseAll](ctx, w.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (w *withdrawQueryCache) SetCachedWithdrawsCache(ctx context.Context, req *requests.FindAllWithdraws, data []*models.WithdrawListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.WithdrawListRow{}
	}

	key := fmt.Sprintf(withdrawAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &withdrawCachedResponseAll{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, w.store, key, payload, ttlDefault)
}

func (w *withdrawQueryCache) GetCachedWithdrawByCardCache(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*models.WithdrawByCardNumberRow, *int, bool) {
	key := fmt.Sprintf(withdrawByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[withdrawCachedResponseByCard](ctx, w.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (w *withdrawQueryCache) SetCachedWithdrawByCardCache(ctx context.Context, req *requests.FindAllWithdrawCardNumber, data []*models.WithdrawByCardNumberRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.WithdrawByCardNumberRow{}
	}

	key := fmt.Sprintf(withdrawByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)
	payload := &withdrawCachedResponseByCard{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, w.store, key, payload, ttlDefault)
}

func (w *withdrawQueryCache) GetCachedWithdrawActiveCache(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListRow, *int, bool) {
	key := fmt.Sprintf(withdrawActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[withdrawCachedResponseActive](ctx, w.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (w *withdrawQueryCache) SetCachedWithdrawActiveCache(ctx context.Context, req *requests.FindAllWithdraws, data []*models.WithdrawListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.WithdrawListRow{}
	}

	key := fmt.Sprintf(withdrawActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &withdrawCachedResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, w.store, key, payload, ttlDefault)
}

func (w *withdrawQueryCache) GetCachedWithdrawTrashedCache(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListWithDeletedRow, *int, bool) {
	key := fmt.Sprintf(withdrawTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[withdrawCachedResponseTrashed](ctx, w.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (w *withdrawQueryCache) SetCachedWithdrawTrashedCache(ctx context.Context, req *requests.FindAllWithdraws, data []*models.WithdrawListWithDeletedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.WithdrawListWithDeletedRow{}
	}

	key := fmt.Sprintf(withdrawTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &withdrawCachedResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, w.store, key, payload, ttlDefault)
}

func (w *withdrawQueryCache) GetCachedWithdrawCache(ctx context.Context, id int) (*models.WithdrawAllFieldsRow, bool) {
	key := fmt.Sprintf(withdrawByIdCacheKey, id)
	result, found := cache.GetFromCache[*models.WithdrawAllFieldsRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawQueryCache) SetCachedWithdrawCache(ctx context.Context, data *models.WithdrawAllFieldsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(withdrawByIdCacheKey, data.WithdrawID)
	cache.SetToCache(ctx, w.store, key, data, ttlDefault)
}
