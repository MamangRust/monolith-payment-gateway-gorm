package repository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoRepository interface {
	FindByCardNumber(ctx context.Context, card_number string) (*models.Saldo, error)
	UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*models.UpdateSaldoBalanceRow, error)
	UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*models.UpdateSaldoWithdrawRow, error)
}

type WithdrawQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListWithDeletedRow, error)
	FindAllByCardNumber(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*models.WithdrawByCardNumberRow, error)
	FindById(ctx context.Context, id int) (*models.WithdrawAllFieldsRow, error)
}

type WithdrawCommandRepository interface {
	CreateWithdraw(ctx context.Context, request *requests.CreateWithdrawRequest) (*models.WithdrawAllFieldsRow, error)
	CreateWithdrawAtomic(ctx context.Context, request *requests.CreateWithdrawRequest) (*WithdrawAtomicResult, error)
	UpdateWithdraw(ctx context.Context, request *requests.UpdateWithdrawRequest) (*models.WithdrawAllFieldsRow, error)
	UpdateWithdrawAtomic(ctx context.Context, request *requests.UpdateWithdrawRequest) (*models.WithdrawAllFieldsRow, error)
	UpdateSaldoBalanceDelta(ctx context.Context, cardNumber string, delta int) error
	UpdateWithdrawStatus(ctx context.Context, request *requests.UpdateWithdrawStatus) (*models.WithdrawAllFieldsRow, error)
	TrashedWithdraw(ctx context.Context, withdrawID int) (*models.Withdraw, error)
	RestoreWithdraw(ctx context.Context, withdrawID int) (*models.Withdraw, error)
	DeleteWithdrawPermanent(ctx context.Context, withdrawID int) (bool, error)
	RestoreAllWithdraw(ctx context.Context) (bool, error)
	DeleteAllWithdrawPermanent(ctx context.Context) (bool, error)
}

type CardRepository interface {
	FindUserCardByCardNumber(ctx context.Context, card_number string) (*models.CardByEmailRow, error)
}
