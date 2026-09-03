package mencache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransferQueryCache interface {
	GetCachedTransfersCache(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListRow, *int, bool)
	SetCachedTransfersCache(ctx context.Context, req *requests.FindAllTransfers, data []*models.TransferListRow, total *int)

	GetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListRow, *int, bool)
	SetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTransfers, data []*models.TransferListRow, total *int)

	GetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListWithDeletedRow, *int, bool)
	SetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTransfers, data []*models.TransferListWithDeletedRow, total *int)

	GetCachedTransferCache(ctx context.Context, id int) (*models.TransferAllFieldsRow, bool)
	SetCachedTransferCache(ctx context.Context, data *models.TransferAllFieldsRow)

	GetCachedTransferByFrom(ctx context.Context, from string) ([]*models.TransferBySourceCardRow, bool)
	SetCachedTransferByFrom(ctx context.Context, from string, data []*models.TransferBySourceCardRow)

	GetCachedTransferByTo(ctx context.Context, to string) ([]*models.TransferByDestinationCardRow, bool)
	SetCachedTransferByTo(ctx context.Context, to string, data []*models.TransferByDestinationCardRow)
}

type TransferCommandCache interface {
	DeleteTransferCache(ctx context.Context, id int)
}
