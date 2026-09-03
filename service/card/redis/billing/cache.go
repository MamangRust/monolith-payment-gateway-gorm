package cardbillingmencache

import (
	"context"
	"fmt"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type cardBillingCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewCardBillingCache(store *sharedcachehelpers.CacheStore) CardBillingCache {
	return &cardBillingCache{store: store}
}

func (c *cardBillingCache) GetByCardNumber(ctx context.Context, cardNumber string) ([]*models.BillingCycle, bool) {
	key := fmt.Sprintf(billingByCardCacheKey, cardNumber)
	result, found := sharedcachehelpers.GetFromCache[[]*models.BillingCycle](ctx, c.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *cardBillingCache) SetByCardNumber(ctx context.Context, cardNumber string, data []*models.BillingCycle) {
	key := fmt.Sprintf(billingByCardCacheKey, cardNumber)
	sharedcachehelpers.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *cardBillingCache) DeleteByCardNumber(ctx context.Context, cardNumber string) {
	key := fmt.Sprintf(billingByCardCacheKey, cardNumber)
	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}
