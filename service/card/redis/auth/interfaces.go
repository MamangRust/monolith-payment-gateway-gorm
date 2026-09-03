package cardauthmencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
)

type CardAuthCache interface {
	GetByTxnID(ctx context.Context, txnID string) (*models.CardAuthTransaction, bool)
	SetByTxnID(ctx context.Context, txnID string, data *models.CardAuthTransaction)
	DeleteByTxnID(ctx context.Context, txnID string)
}
