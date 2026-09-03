package transfer_test

import (
	"context"
	"time"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
)

type realCardRepo struct {
	query card_repo.CardQueryRepository
}
func (r *realCardRepo) FindUserCardByCardNumber(ctx context.Context, card_number string) (*models.CardByEmailRow, error) {
	return r.query.FindUserCardByCardNumber(ctx, card_number)
}
func (r *realCardRepo) FindCardByCardNumber(ctx context.Context, card_number string) (*models.CardAllFieldsRow, error) {
	return r.query.FindCardByCardNumber(ctx, card_number)
}

type realSaldoRepo struct {
	repo saldo_repo.Repositories
}
func (r *realSaldoRepo) FindByCardNumber(ctx context.Context, card_number string) (*models.Saldo, error) {
	return r.repo.FindByCardNumber(ctx, card_number)
}
func (r *realSaldoRepo) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*models.UpdateSaldoBalanceRow, error) {
	return r.repo.UpdateSaldoBalance(ctx, request)
}

type dummyCacheMetrics struct{}
func (d *dummyCacheMetrics) RecordCacheHit(ctx context.Context, key string) {}
func (d *dummyCacheMetrics) RecordCacheMiss(ctx context.Context, key string) {}
func (d *dummyCacheMetrics) RecordCacheSet(ctx context.Context, key string, success bool) {}
func (d *dummyCacheMetrics) RecordCacheDelete(ctx context.Context, key string, success bool) {}
func (d *dummyCacheMetrics) RecordCacheOperationLatency(ctx context.Context, operation string, duration time.Duration) {}
func (d *dummyCacheMetrics) RecordCacheError(ctx context.Context, operation, key string, err error) {}
