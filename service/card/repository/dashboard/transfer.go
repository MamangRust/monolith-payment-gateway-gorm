package repositorydashboard

import (
	"context"

	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardDashboardTransferRepository struct {
	db *gorm.DB
}

func NewCardDashboardTransferRepository(db *gorm.DB) CardDashboardTransferRepository {
	return &cardDashboardTransferRepository{db: db}
}

func (r *cardDashboardTransferRepository) GetTotalTransferAmount(ctx context.Context) (*int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw("SELECT COALESCE(SUM(transfer_amount), 0) FROM transfers WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, card_errors.ErrGetTotalTransfersFailed.WithInternal(err)
	}
	return &total, nil
}

func (r *cardDashboardTransferRepository) GetTotalTransferAmountBySender(ctx context.Context, senderCardNumber string) (*int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw("SELECT COALESCE(SUM(transfer_amount), 0) FROM transfers WHERE transfer_from = ? AND deleted_at IS NULL", senderCardNumber).Scan(&total).Error
	if err != nil {
		return nil, card_errors.ErrGetTotalTransferSenderByCardFailed.WithInternal(err)
	}
	return &total, nil
}

func (r *cardDashboardTransferRepository) GetTotalTransferAmountByReceiver(ctx context.Context, receiverCardNumber string) (*int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw("SELECT COALESCE(SUM(transfer_amount), 0) FROM transfers WHERE transfer_to = ? AND deleted_at IS NULL", receiverCardNumber).Scan(&total).Error
	if err != nil {
		return nil, card_errors.ErrGetTotalTransferSenderByCardFailed.WithInternal(err)
	}
	return &total, nil
}
