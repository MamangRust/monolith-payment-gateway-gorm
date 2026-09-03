package cardpaymentmencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
)

type CardPaymentCache interface {
	GetPaymentHistory(ctx context.Context, cardNumber string, page, pageSize int) ([]*models.CardPaymentListRow, bool)
	SetPaymentHistory(ctx context.Context, cardNumber string, page, pageSize int, data []*models.CardPaymentListRow)
	DeletePaymentHistory(ctx context.Context, cardNumber string)

	GetPaymentCount(ctx context.Context, cardNumber string) (int, bool)
	SetPaymentCount(ctx context.Context, cardNumber string, count int)
	DeletePaymentCount(ctx context.Context, cardNumber string)
}
