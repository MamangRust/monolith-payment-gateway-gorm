package transferstatsbycardcache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type transferStatsByCardAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransferStatsByCardAmountCache(store *sharedcachehelpers.CacheStore) TransferStatsByCardAmountCache {
	return &transferStatsByCardAmountCache{store: store}
}

func (t *transferStatsByCardAmountCache) GetMonthlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferMonthlyAmountBySenderCardRow, bool) {
	key := fmt.Sprintf(transferMonthTransferAmountBySenderCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferMonthlyAmountBySenderCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsByCardAmountCache) SetMonthlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*models.TransferMonthlyAmountBySenderCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferAmountBySenderCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsByCardAmountCache) GetMonthlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferMonthlyAmountByReceiverCardRow, bool) {
	key := fmt.Sprintf(transferMonthTransferAmountByReceiverCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferMonthlyAmountByReceiverCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsByCardAmountCache) SetMonthlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*models.TransferMonthlyAmountByReceiverCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferAmountByReceiverCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsByCardAmountCache) GetYearlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferYearlyAmountBySenderCardRow, bool) {
	key := fmt.Sprintf(transferYearTransferAmountBySenderCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferYearlyAmountBySenderCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsByCardAmountCache) SetYearlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*models.TransferYearlyAmountBySenderCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferAmountBySenderCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsByCardAmountCache) GetYearlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferYearlyAmountByReceiverCardRow, bool) {
	key := fmt.Sprintf(transferYearTransferAmountByReceiverCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*models.TransferYearlyAmountByReceiverCardRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transferStatsByCardAmountCache) SetYearlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*models.TransferYearlyAmountByReceiverCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferAmountByReceiverCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
