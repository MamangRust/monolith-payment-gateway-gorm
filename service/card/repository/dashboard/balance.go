package repositorydashboard

import (
	"context"

	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardDashboardBalanceRepository struct {
	db *gorm.DB
}

func NewCardDashboardBalanceRepository(db *gorm.DB) CardDashboardBalanceRepository {
	return &cardDashboardBalanceRepository{db: db}
}

func (r *cardDashboardBalanceRepository) GetTotalBalances(ctx context.Context) (*int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw("SELECT COALESCE(SUM(total_balance), 0) FROM saldos WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, card_errors.ErrGetTotalBalancesFailed.WithInternal(err)
	}
	return &total, nil
}

func (r *cardDashboardBalanceRepository) GetTotalBalanceByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw("SELECT COALESCE(SUM(total_balance), 0) FROM saldos WHERE card_number = ? AND deleted_at IS NULL", cardNumber).Scan(&total).Error
	if err != nil {
		return nil, card_errors.ErrGetTotalBalanceByCardFailed.WithInternal(err)
	}
	return &total, nil
}
