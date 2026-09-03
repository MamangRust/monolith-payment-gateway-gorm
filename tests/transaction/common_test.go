package transaction_test

import (
	"context"
	"time"

	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	merchant_repo "github.com/MamangRust/monolith-payment-gateway-merchant/repository"
	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type realMerchantRepo struct {
	repo merchant_repo.MerchantQueryRepository
}

func (r *realMerchantRepo) FindByApiKey(ctx context.Context, api_key string) (*models.MerchantAllFieldsRow, error) {
	return r.repo.FindByApiKey(ctx, api_key)
}

type transactionCardRepo struct {
	query   card_repo.CardQueryRepository
	command card_repo.CardCommandRepository
}

func (r *transactionCardRepo) FindCardByUserId(ctx context.Context, user_id int) (*models.CardAllFieldsRow, error) {
	return r.query.FindCardByUserId(ctx, user_id)
}
func (r *transactionCardRepo) FindUserCardByCardNumber(ctx context.Context, card_number string) (*models.CardByEmailRow, error) {
	return r.query.FindUserCardByCardNumber(ctx, card_number)
}
func (r *transactionCardRepo) FindCardByCardNumber(ctx context.Context, card_number string) (*models.CardAllFieldsRow, error) {
	return r.query.FindCardByCardNumber(ctx, card_number)
}
func (r *transactionCardRepo) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*models.CardUpdateRow, error) {
	return r.command.UpdateCard(ctx, request)
}
func (r *transactionCardRepo) UpdateCardOutstandingBalance(ctx context.Context, cardID int, outstandingBalance int) (*models.UpdateOutstandingBalanceRow, error) {
	return &models.UpdateOutstandingBalanceRow{}, nil
}
func (r *transactionCardRepo) AddRewardPoints(ctx context.Context, cardID int, points int) (*models.AddRewardPointsRow, error) {
	return &models.AddRewardPointsRow{}, nil
}

type realCardRepo = transactionCardRepo

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

func (d *dummyCacheMetrics) RecordCacheHit(ctx context.Context, key string)                  {}
func (d *dummyCacheMetrics) RecordCacheMiss(ctx context.Context, key string)                 {}
func (d *dummyCacheMetrics) RecordCacheSet(ctx context.Context, key string, success bool)    {}
func (d *dummyCacheMetrics) RecordCacheDelete(ctx context.Context, key string, success bool) {}
func (d *dummyCacheMetrics) RecordCacheOperationLatency(ctx context.Context, operation string, duration time.Duration) {
}
func (d *dummyCacheMetrics) RecordCacheError(ctx context.Context, operation, key string, err error) {}
