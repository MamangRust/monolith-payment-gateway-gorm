package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoQueryCache interface {
	GetCachedSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoRow, *int, bool)
	SetCachedSaldos(ctx context.Context, req *requests.FindAllSaldos, data []*models.SaldoRow, totalRecords *int)

	GetCachedSaldoById(ctx context.Context, saldo_id int) (*models.SaldoByIDRow, bool)
	SetCachedSaldoById(ctx context.Context, saldo_id int, data *models.SaldoByIDRow)

	GetCachedSaldoByCardNumber(ctx context.Context, card_number string) (*models.Saldo, bool)
	SetCachedSaldoByCardNumber(ctx context.Context, card_number string, data *models.Saldo)

	GetCachedSaldoByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoActiveRow, *int, bool)
	SetCachedSaldoByActive(ctx context.Context, req *requests.FindAllSaldos, data []*models.SaldoActiveRow, totalRecords *int)

	GetCachedSaldoByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoTrashedRow, *int, bool)
	SetCachedSaldoByTrashed(ctx context.Context, req *requests.FindAllSaldos, data []*models.SaldoTrashedRow, totalRecords *int)
}

type SaldoCommandCache interface {
	DeleteSaldoCache(ctx context.Context, saldo_id int)
	InvalidateSaldoCache(ctx context.Context)
}
