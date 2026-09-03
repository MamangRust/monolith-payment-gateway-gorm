package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"gorm.io/gorm"
)

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) FindCardByCardNumber(ctx context.Context, card_number string) (*models.Card, error) {
	var card models.Card
	if err := r.db.WithContext(ctx).Where("card_number = ?", card_number).First(&card).Error; err != nil {
		return nil, err
	}
	return &card, nil
}
