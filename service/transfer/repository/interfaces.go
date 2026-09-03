package repository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoRepository interface {
	FindByCardNumber(ctx context.Context, card_number string) (*models.Saldo, error)
	UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*models.UpdateSaldoBalanceRow, error)
}

type CardRepository interface {
	FindUserCardByCardNumber(ctx context.Context, card_number string) (*models.CardByEmailRow, error)
	FindCardByCardNumber(ctx context.Context, card_number string) (*models.CardAllFieldsRow, error)
}

type TransferQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListWithDeletedRow, error)
	FindById(ctx context.Context, id int) (*models.TransferAllFieldsRow, error)
	FindTransferByTransferFrom(ctx context.Context, transfer_from string) ([]*models.TransferBySourceCardRow, error)
	FindTransferByTransferTo(ctx context.Context, transfer_to string) ([]*models.TransferByDestinationCardRow, error)
}

type TransferCommandRepository interface {
	CreateTransferAtomic(ctx context.Context, request *requests.CreateTransferRequest) (*TransferAtomicResult, error)
	UpdateTransferSettlementDelta(ctx context.Context, senderCard, receiverCard string, senderDelta, receiverDelta int) error
	UpdateTransferAtomic(ctx context.Context, params models.TransferAllFieldsRow) (*models.TransferAllFieldsRow, error)
	UpdateTransfer(ctx context.Context, request *requests.UpdateTransferRequest) (*models.TransferAllFieldsRow, error)
	UpdateTransferAmount(ctx context.Context, request *requests.UpdateTransferAmountRequest) (*models.TransferAllFieldsRow, error)
	UpdateTransferStatus(ctx context.Context, request *requests.UpdateTransferStatus) (*models.TransferAllFieldsRow, error)
	TrashedTransfer(ctx context.Context, transferID int) (*models.Transfer, error)
	RestoreTransfer(ctx context.Context, transferID int) (*models.Transfer, error)
	DeleteTransferPermanent(ctx context.Context, transferID int) (bool, error)
	RestoreAllTransfer(ctx context.Context) (bool, error)
	DeleteAllTransferPermanent(ctx context.Context) (bool, error)
}
