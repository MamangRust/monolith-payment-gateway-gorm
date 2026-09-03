package transferstatsbycardcache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type transferStatsByCardStatusCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransferStatsByCardStatusCache(store *sharedcachehelpers.CacheStore) TransferStatsByCardStatusCache {
	return &transferStatsByCardStatusCache{store: store}
}

func (t *transferStatsByCardStatusCache) GetMonthTransferStatusSuccessByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*models.TransferMonthlyStatusSuccessByCardRow, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferMonthlyStatusSuccessByCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsByCardStatusCache) SetMonthTransferStatusSuccessByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber, data []*models.TransferMonthlyStatusSuccessByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsByCardStatusCache) GetYearlyTransferStatusSuccessByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*models.TransferYearlyStatusSuccessByCardRow, bool) {
	key := fmt.Sprintf(transferYearTransferStatusSuccessByCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferYearlyStatusSuccessByCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsByCardStatusCache) SetYearlyTransferStatusSuccessByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber, data []*models.TransferYearlyStatusSuccessByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusSuccessByCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsByCardStatusCache) GetMonthTransferStatusFailedByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*models.TransferMonthlyStatusFailedByCardRow, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusFailedByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferMonthlyStatusFailedByCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsByCardStatusCache) SetMonthTransferStatusFailedByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber, data []*models.TransferMonthlyStatusFailedByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusFailedByCardKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsByCardStatusCache) GetYearlyTransferStatusFailedByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*models.TransferYearlyStatusFailedByCardRow, bool) {
	key := fmt.Sprintf(transferYearTransferStatusFailedByCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferYearlyStatusFailedByCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transferStatsByCardStatusCache) SetYearlyTransferStatusFailedByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber, data []*models.TransferYearlyStatusFailedByCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusFailedByCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
