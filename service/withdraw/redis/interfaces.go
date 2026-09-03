package mencache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type WithdrawQueryCache interface {
	GetCachedWithdrawsCache(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListRow, *int, bool)
	SetCachedWithdrawsCache(ctx context.Context, req *requests.FindAllWithdraws, data []*models.WithdrawListRow, total *int)

	GetCachedWithdrawByCardCache(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*models.WithdrawByCardNumberRow, *int, bool)
	SetCachedWithdrawByCardCache(ctx context.Context, req *requests.FindAllWithdrawCardNumber, data []*models.WithdrawByCardNumberRow, total *int)

	GetCachedWithdrawActiveCache(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListRow, *int, bool)
	SetCachedWithdrawActiveCache(ctx context.Context, req *requests.FindAllWithdraws, data []*models.WithdrawListRow, total *int)

	GetCachedWithdrawTrashedCache(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListWithDeletedRow, *int, bool)
	SetCachedWithdrawTrashedCache(ctx context.Context, req *requests.FindAllWithdraws, data []*models.WithdrawListWithDeletedRow, total *int)

	GetCachedWithdrawCache(ctx context.Context, id int) (*models.WithdrawAllFieldsRow, bool)
	SetCachedWithdrawCache(ctx context.Context, data *models.WithdrawAllFieldsRow)
}

type WithdrawCommandCache interface {
	DeleteCachedWithdrawCache(ctx context.Context, id int)
}
