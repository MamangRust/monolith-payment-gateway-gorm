package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-pkg/date"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-pkg/randomvcc"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// cardSeeder is a struct that represents a seeder for the cards table.
type cardSeeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

// NewCardSeeder creates a new instance of the cardSeeder, which is
// responsible for populating the cards table with fake data.
func NewCardSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *cardSeeder {
	return &cardSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

// Seed populates the cards table with fake data.
//
// It generates 10 cards, 5 of which are active and 5 of which are trashed.
// The cards are seeded with random card numbers, card types, expire dates, and
// CVV numbers. The card provider is randomly selected from a list of providers.
//
// If any errors occur during the seeding process, an error is returned.
func (r *cardSeeder) Seed() error {
	cardTypes := []string{"credit", "debit"}
	cardProviders := []string{"mandiri", "bni", "bri"}

	totalCards := 10
	activeCards := 5
	trashedCards := 5

	cardNumbers := make([]string, totalCards)
	for i := 0; i < totalCards; i++ {
		cardNumber, err := randomvcc.RandomCardNumber()
		if err != nil {
			r.logger.Error("failed to generate card number", zap.Int("index", i), zap.Error(err))
			return fmt.Errorf("failed to generate card number: %w", err)
		}
		cardNumbers[i] = cardNumber
	}

	for i := 0; i < totalCards; i++ {
		card := &models.Card{
			UserID:       int32(i + 1),
			CardNumber:   cardNumbers[i],
			CardType:     cardTypes[i%len(cardTypes)],
			ExpireDate:   date.GenerateExpireDate(),
			Cvv:          fmt.Sprintf("%03d", i%1000),
			CardProvider: cardProviders[i%len(cardProviders)],
		}

		if err := r.db.WithContext(r.ctx).Create(card).Error; err != nil {
			r.logger.Error("failed to seed card", zap.Int("card", i+1), zap.Error(err))
			return fmt.Errorf("failed to seed card %d: %w", i+1, err)
		}

		if i >= activeCards {
			now := time.Now()
			if err := r.db.WithContext(r.ctx).Model(&models.Card{}).
				Where("card_id = ? AND deleted_at IS NULL", card.CardID).
				Update("deleted_at", &now).Error; err != nil {
				r.logger.Error("failed to trash card", zap.Int("card", i+1), zap.Error(err))
				return fmt.Errorf("failed to trash card %d: %w", i+1, err)
			}
		}
	}

	r.logger.Info("card seeded successfully", zap.Int("totalCards", totalCards), zap.Int("activeCards", activeCards), zap.Int("trashedCards", trashedCards))
	return nil
}
