package mencache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransactionQueryCache interface {
	GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransactions) ([]*models.TransactionListRow, *int, bool)
	SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransactions, data []*models.TransactionListRow, total *int)

	GetCachedTransactionByCardNumberCache(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*models.TransactionListRow, *int, bool)
	SetCachedTransactionByCardNumberCache(ctx context.Context, req *requests.FindAllTransactionCardNumber, data []*models.TransactionListRow, total *int)

	GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransactions) ([]*models.TransactionListRow, *int, bool)
	SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransactions, data []*models.TransactionListRow, total *int)

	GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransactions) ([]*models.TransactionListWithDeletedRow, *int, bool)
	SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransactions, data []*models.TransactionListWithDeletedRow, total *int)

	GetCachedTransactionByMerchantIdCache(ctx context.Context, merchant_id int) ([]*models.TransactionListRow, bool)
	SetCachedTransactionByMerchantIdCache(ctx context.Context, merchant_id int, data []*models.TransactionListRow)

	GetCachedTransactionCache(ctx context.Context, id int) (*models.TransactionAllFieldsRow, bool)
	SetCachedTransactionCache(ctx context.Context, data *models.TransactionAllFieldsRow)
}

type TransactionCommandCache interface {
	DeleteTransactionCache(ctx context.Context, id int)
	InvalidateTransactionCache(ctx context.Context)
}
