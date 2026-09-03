package service

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransferQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListRow, *int, error)
	FindById(ctx context.Context, transferId int) (*models.TransferAllFieldsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListWithDeletedRow, *int, error)
	FindTransferByTransferFrom(ctx context.Context, transfer_from string) ([]*models.TransferBySourceCardRow, error)
	FindTransferByTransferTo(ctx context.Context, transfer_to string) ([]*models.TransferByDestinationCardRow, error)
}

type TransferCommandService interface {
	CreateTransaction(ctx context.Context, request *requests.CreateTransferRequest) (*models.TransferAllFieldsRow, error)
	UpdateTransaction(ctx context.Context, request *requests.UpdateTransferRequest) (*models.TransferAllFieldsRow, error)
	TrashedTransfer(ctx context.Context, transfer_id int) (*models.Transfer, error)
	RestoreTransfer(ctx context.Context, transfer_id int) (*models.Transfer, error)
	DeleteTransferPermanent(ctx context.Context, transfer_id int) (bool, error)

	RestoreAllTransfer(ctx context.Context) (bool, error)
	DeleteAllTransferPermanent(ctx context.Context) (bool, error)
}
