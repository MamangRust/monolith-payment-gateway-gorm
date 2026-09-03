package cardrewardmencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
)

type CardRewardCache interface {
	GetBalance(ctx context.Context, cardNumber string) (int64, bool)
	SetBalance(ctx context.Context, cardNumber string, balance int64)
	DeleteBalance(ctx context.Context, cardNumber string)

	GetHistory(ctx context.Context, cardNumber string) ([]*models.CardReward, bool)
	SetHistory(ctx context.Context, cardNumber string, data []*models.CardReward)
	DeleteHistory(ctx context.Context, cardNumber string)
}
