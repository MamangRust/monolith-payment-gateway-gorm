package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoRow, *int, error)
	FindById(ctx context.Context, saldo_id int) (*models.SaldoByIDRow, error)
	FindByCardNumber(ctx context.Context, card_number string) (*models.Saldo, error)
	FindByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoTrashedRow, *int, error)
}

type SaldoCommandService interface {
	CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*models.CreateSaldoRow, error)
	CreateSaldoIfNotExists(ctx context.Context, request *requests.CreateSaldoRequest) error
	InvalidateSaldoCache(ctx context.Context)
	UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*models.UpdateSaldoRow, error)
	UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*models.UpdateSaldoBalanceRow, error)
	UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*models.UpdateSaldoWithdrawRow, error)
	TrashSaldo(ctx context.Context, saldo_id int) (*models.Saldo, error)
	RestoreSaldo(ctx context.Context, saldo_id int) (*models.Saldo, error)
	DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, error)

	RestoreAllSaldo(ctx context.Context) (bool, error)
	DeleteAllSaldoPermanent(ctx context.Context) (bool, error)
}
