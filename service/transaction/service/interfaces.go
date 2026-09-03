package service

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// TransactionQueryService handles queries related to transactions.
type TransactionQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllTransactions) ([]*models.TransactionListRow, *int, error)
	FindAllByCardNumber(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*models.TransactionListRow, *int, error)
	FindById(ctx context.Context, transactionID int) (*models.TransactionAllFieldsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllTransactions) ([]*models.TransactionListRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransactions) ([]*models.TransactionListWithDeletedRow, *int, error)
	FindTransactionByMerchantId(ctx context.Context, merchant_id int) ([]*models.TransactionListRow, error)
}

type TransactionCommandService interface {
	Create(ctx context.Context, apiKey string, request *requests.CreateTransactionRequest) (*models.TransactionAllFieldsRow, error)
	Update(ctx context.Context, apiKey string, request *requests.UpdateTransactionRequest) (*models.TransactionAllFieldsRow, error)
	TrashedTransaction(ctx context.Context, transaction_id int) (*models.TransactionTrashRestoreRow, error)
	RestoreTransaction(ctx context.Context, transaction_id int) (*models.TransactionTrashRestoreRow, error)
	DeleteTransactionPermanent(ctx context.Context, transaction_id int) (bool, error)

	RestoreAllTransaction(ctx context.Context) (bool, error)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)

	AuthorizeTransaction(ctx context.Context, apiKey string, request *requests.CreateTransactionRequest) (*models.TransactionAllFieldsRow, error)
	CaptureTransaction(ctx context.Context, apiKey string, transactionID int) (*models.TransactionAllFieldsRow, error)
	VoidTransaction(ctx context.Context, apiKey string, transactionID int) (*models.TransactionAllFieldsRow, error)
	RefundTransaction(ctx context.Context, apiKey string, transactionID int) (*models.TransactionAllFieldsRow, error)
}
