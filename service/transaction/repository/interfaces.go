package repository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type MerchantRepository interface {
	FindByApiKey(ctx context.Context, api_key string) (*models.MerchantAllFieldsRow, error)
}

type SaldoRepository interface {
	FindByCardNumber(ctx context.Context, card_number string) (*models.Saldo, error)

	UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*models.UpdateSaldoBalanceRow, error)
}

type CardRepository interface {
	FindCardByUserId(ctx context.Context, user_id int) (*models.CardAllFieldsRow, error)

	FindUserCardByCardNumber(ctx context.Context, card_number string) (*models.CardByEmailRow, error)

	FindCardByCardNumber(ctx context.Context, card_number string) (*models.CardAllFieldsRow, error)

	UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*models.CardUpdateRow, error)

	UpdateCardOutstandingBalance(ctx context.Context, cardID int, outstandingBalance int) (*models.UpdateOutstandingBalanceRow, error)

	AddRewardPoints(ctx context.Context, cardID int, points int) (*models.AddRewardPointsRow, error)
}

type TransactionQueryRepository interface {
	FindAllTransactions(ctx context.Context, req *requests.FindAllTransactions) ([]*models.TransactionListRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllTransactions) ([]*models.TransactionListRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransactions) ([]*models.TransactionListWithDeletedRow, error)
	FindAllTransactionByCardNumber(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*models.TransactionListRow, error)
	FindById(ctx context.Context, transaction_id int) (*models.TransactionAllFieldsRow, error)
	FindTransactionByMerchantId(ctx context.Context, merchant_id int) ([]*models.TransactionListRow, error)
}

type TransactionCommandRepository interface {
	CreateTransactionAtomic(ctx context.Context, params models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error)
	AuthorizeTransactionAtomic(ctx context.Context, params models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error)
	AuthorizeTransactionAtomicIdempotent(ctx context.Context, params models.TransactionAllFieldsRow) (*AuthorizeTransactionResult, error)
	CreateTransactionAtomicIdempotent(ctx context.Context, params models.TransactionAllFieldsRow) (*TransactionAtomicResult, error)
	UpdateSaldoBalanceDelta(ctx context.Context, cardNumber string, delta int) error
	UpdateTransactionAtomic(ctx context.Context, params models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error)
	CaptureTransactionAtomic(ctx context.Context, params models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error)
	VoidTransactionAtomic(ctx context.Context, transactionID int32) (*models.TransactionAllFieldsRow, error)
	RefundTransactionAtomic(ctx context.Context, params models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error)
	CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*models.TransactionAllFieldsRow, error)
	UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*models.TransactionAllFieldsRow, error)
	UpdateTransactionStatus(ctx context.Context, request *requests.UpdateTransactionStatus) (*models.TransactionAllFieldsRow, error)
	TrashedTransaction(ctx context.Context, transaction_id int) (*models.TransactionTrashRestoreRow, error)
	RestoreTransaction(ctx context.Context, topup_id int) (*models.TransactionTrashRestoreRow, error)
	DeleteTransactionPermanent(ctx context.Context, topup_id int) (bool, error)
	RestoreAllTransaction(ctx context.Context) (bool, error)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)
}
