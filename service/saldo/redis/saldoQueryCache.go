package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type saldoCachedResponseAll struct {
	Data         []*models.SaldoRow `json:"data"`
	TotalRecords *int               `json:"total_records"`
}

type saldoCachedResponseActive struct {
	Data         []*models.SaldoActiveRow `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type saldoCachedResponseTrashed struct {
	Data         []*models.SaldoTrashedRow `json:"data"`
	TotalRecords *int                      `json:"total_records"`
}

type saldoQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewSaldoQueryCache(store *sharedcachehelpers.CacheStore) SaldoQueryCache {
	return &saldoQueryCache{store: store}
}

func (s *saldoQueryCache) GetCachedSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoRow, *int, bool) {
	key := fmt.Sprintf(saldoAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcachehelpers.GetFromCache[saldoCachedResponseAll](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *saldoQueryCache) SetCachedSaldos(ctx context.Context, req *requests.FindAllSaldos, data []*models.SaldoRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.SaldoRow{}
	}
	key := fmt.Sprintf(saldoAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &saldoCachedResponseAll{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *saldoQueryCache) GetCachedSaldoByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoActiveRow, *int, bool) {
	key := fmt.Sprintf(saldoActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcachehelpers.GetFromCache[saldoCachedResponseActive](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *saldoQueryCache) SetCachedSaldoByActive(ctx context.Context, req *requests.FindAllSaldos, data []*models.SaldoActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.SaldoActiveRow{}
	}
	key := fmt.Sprintf(saldoActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &saldoCachedResponseActive{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *saldoQueryCache) GetCachedSaldoByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoTrashedRow, *int, bool) {
	key := fmt.Sprintf(saldoTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcachehelpers.GetFromCache[saldoCachedResponseTrashed](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *saldoQueryCache) SetCachedSaldoByTrashed(ctx context.Context, req *requests.FindAllSaldos, data []*models.SaldoTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.SaldoTrashedRow{}
	}
	key := fmt.Sprintf(saldoTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &saldoCachedResponseTrashed{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *saldoQueryCache) GetCachedSaldoById(ctx context.Context, saldo_id int) (*models.SaldoByIDRow, bool) {
	key := fmt.Sprintf(saldoByIdCacheKey, saldo_id)
	result, found := sharedcachehelpers.GetFromCache[*models.SaldoByIDRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *saldoQueryCache) SetCachedSaldoById(ctx context.Context, saldo_id int, data *models.SaldoByIDRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(saldoByIdCacheKey, saldo_id)
	sharedcachehelpers.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *saldoQueryCache) GetCachedSaldoByCardNumber(ctx context.Context, card_number string) (*models.Saldo, bool) {
	key := fmt.Sprintf(saldoByCardNumberKey, card_number)
	result, found := sharedcachehelpers.GetFromCache[*models.Saldo](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *saldoQueryCache) SetCachedSaldoByCardNumber(ctx context.Context, card_number string, data *models.Saldo) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(saldoByCardNumberKey, card_number)
	sharedcachehelpers.SetToCache(ctx, s.store, key, data, ttlDefault)
}
