package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoQueryRepository interface {
	FindAllSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoTrashedRow, error)
	FindById(ctx context.Context, saldo_id int) (*models.SaldoByIDRow, error)
	FindByCardNumber(ctx context.Context, card_number string) (*models.Saldo, error)
}

type SaldoCommandRepository interface {
	CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*models.CreateSaldoRow, error)
	CreateSaldoIfNotExists(ctx context.Context, request *requests.CreateSaldoRequest) error
	UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*models.UpdateSaldoRow, error)
	UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*models.UpdateSaldoBalanceRow, error)
	UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*models.UpdateSaldoWithdrawRow, error)
	TrashedSaldo(ctx context.Context, saldoID int) (*models.Saldo, error)
	RestoreSaldo(ctx context.Context, saldoID int) (*models.Saldo, error)
	DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, error)

	RestoreAllSaldo(ctx context.Context) (bool, error)
	DeleteAllSaldoPermanent(ctx context.Context) (bool, error)
}

type CardRepository interface {
	FindCardByCardNumber(ctx context.Context, card_number string) (*models.Card, error)
}
