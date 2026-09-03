package repositorydashboard

import (
	"context"

	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardDashboardTopupRepository struct {
	db *gorm.DB
}

func NewCardDashboardTopupRepository(db *gorm.DB) CardDashboardTopupRepository {
	return &cardDashboardTopupRepository{db: db}
}

func (r *cardDashboardTopupRepository) GetTotalTopAmount(ctx context.Context) (*int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw("SELECT COALESCE(SUM(topup_amount), 0) FROM topups WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, card_errors.ErrGetTotalTopupsFailed.WithInternal(err)
	}
	return &total, nil
}

func (r *cardDashboardTopupRepository) GetTotalTopupAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw("SELECT COALESCE(SUM(topup_amount), 0) FROM topups WHERE card_number = ? AND deleted_at IS NULL", cardNumber).Scan(&total).Error
	if err != nil {
		return nil, card_errors.ErrGetTotalTopupAmountByCardFailed.WithInternal(err)
	}
	return &total, nil
}
