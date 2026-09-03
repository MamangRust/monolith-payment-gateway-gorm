package repository

import (
	"context"
	"errors"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type billingCycleRepository struct {
	db *gorm.DB
}

func NewBillingCycleRepository(db *gorm.DB) BillingCycleRepository {
	return &billingCycleRepository{db: db}
}

func (r *billingCycleRepository) GetBillingCycleByID(ctx context.Context, billingID int) (*models.BillingCycle, error) {
	var res models.BillingCycle
	err := r.db.WithContext(ctx).Where("billing_id = ?", billingID).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewBadRequestError("billing cycle not found")
		}
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Billing cycle", "get billing cycle")
	}
	return &res, nil
}

func (r *billingCycleRepository) GetBillingCyclesByCardNumber(ctx context.Context, cardNumber string) ([]*models.BillingCycle, error) {
	var res []*models.BillingCycle
	err := r.db.WithContext(ctx).
		Where("card_number = ?", cardNumber).
		Order("created_at DESC").
		Find(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetBillingCyclesFailed.WithInternal(err)
	}
	return res, nil
}
