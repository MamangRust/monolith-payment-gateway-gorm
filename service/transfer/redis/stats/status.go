package transferstatscache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type transferStatsStatusCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransferStatsStatusCache(store *sharedcachehelpers.CacheStore) TransferStatsStatusCache {
	return &transferStatsStatusCache{store: store}
}

func (t *transferStatsStatusCache) GetCachedMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*models.TransferMonthlyStatusSuccessRow, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusSuccessKey, req.Month, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferMonthlyStatusSuccessRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsStatusCache) SetCachedMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer, data []*models.TransferMonthlyStatusSuccessRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusSuccessKey, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsStatusCache) GetCachedYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*models.TransferYearlyStatusSuccessRow, bool) {
	key := fmt.Sprintf(transferYearTransferStatusSuccessKey, year)

	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferYearlyStatusSuccessRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsStatusCache) SetCachedYearlyTransferStatusSuccess(ctx context.Context, year int, data []*models.TransferYearlyStatusSuccessRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusSuccessKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsStatusCache) GetCachedMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*models.TransferMonthlyStatusFailedRow, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusFailedKey, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferMonthlyStatusFailedRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsStatusCache) SetCachedMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer, data []*models.TransferMonthlyStatusFailedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusFailedKey, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsStatusCache) GetCachedYearlyTransferStatusFailed(ctx context.Context, year int) ([]*models.TransferYearlyStatusFailedRow, bool) {
	key := fmt.Sprintf(transferYearTransferStatusFailedKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferYearlyStatusFailedRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsStatusCache) SetCachedYearlyTransferStatusFailed(ctx context.Context, year int, data []*models.TransferYearlyStatusFailedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusFailedKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
