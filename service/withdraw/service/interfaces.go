package service

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// WithdrawQueryService defines query operations for fetching withdraw data.
type WithdrawQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListRow, *int, error)
	FindAllByCardNumber(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*models.WithdrawByCardNumberRow, *int, error)
	FindById(ctx context.Context, withdrawID int) (*models.WithdrawAllFieldsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListWithDeletedRow, *int, error)
}

type WithdrawCommandService interface {
	Create(ctx context.Context, request *requests.CreateWithdrawRequest) (*models.WithdrawAllFieldsRow, error)
	Update(ctx context.Context, request *requests.UpdateWithdrawRequest) (*models.WithdrawAllFieldsRow, error)
	TrashedWithdraw(ctx context.Context, withdraw_id int) (*models.Withdraw, error)
	RestoreWithdraw(ctx context.Context, withdraw_id int) (*models.Withdraw, error)
	DeleteWithdrawPermanent(ctx context.Context, withdraw_id int) (bool, error)

	RestoreAllWithdraw(ctx context.Context) (bool, error)
	DeleteAllWithdrawPermanent(ctx context.Context) (bool, error)
}
