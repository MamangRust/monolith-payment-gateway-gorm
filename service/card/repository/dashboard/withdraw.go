package repositorydashboard

import (
	"context"

	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardDashboardWithdrawRepository struct {
	db *gorm.DB
}

func NewCardDashboardWithdrawRepository(db *gorm.DB) CardDashboardWithdrawRepository {
	return &cardDashboardWithdrawRepository{db: db}
}

func (r *cardDashboardWithdrawRepository) GetTotalWithdrawAmount(ctx context.Context) (*int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw("SELECT COALESCE(SUM(withdraw_amount), 0) FROM withdraws WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, card_errors.ErrGetTotalWithdrawsFailed.WithInternal(err)
	}
	return &total, nil
}

func (r *cardDashboardWithdrawRepository) GetTotalWithdrawAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw("SELECT COALESCE(SUM(withdraw_amount), 0) FROM withdraws WHERE card_number = ? AND deleted_at IS NULL", cardNumber).Scan(&total).Error
	if err != nil {
		return nil, card_errors.ErrGetTotalWithdrawAmountByCardFailed.WithInternal(err)
	}
	return &total, nil
}
