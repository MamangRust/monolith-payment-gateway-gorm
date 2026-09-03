package cardbillingmencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
)

type CardBillingCache interface {
	GetByCardNumber(ctx context.Context, cardNumber string) ([]*models.BillingCycle, bool)
	SetByCardNumber(ctx context.Context, cardNumber string, data []*models.BillingCycle)
	DeleteByCardNumber(ctx context.Context, cardNumber string)
}
