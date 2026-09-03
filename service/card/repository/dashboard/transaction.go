package repositorydashboard

import (
	"context"

	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardDashboardTransactionRepository struct {
	db *gorm.DB
}

func NewCardDashboardTransactionRepository(db *gorm.DB) CardDashboardTransactionRepository {
	return &cardDashboardTransactionRepository{db: db}
}

func (r *cardDashboardTransactionRepository) GetTotalTransactionAmount(ctx context.Context) (*int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, card_errors.ErrGetTotalTransactionsFailed.WithInternal(err)
	}
	return &total, nil
}

func (r *cardDashboardTransactionRepository) GetTotalTransactionAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE card_number = ? AND deleted_at IS NULL", cardNumber).Scan(&total).Error
	if err != nil {
		return nil, card_errors.ErrGetTotalTransactionAmountByCardFailed.WithInternal(err)
	}
	return &total, nil
}
