package seeder

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// saldoSeeder is a struct that represents a seeder for the saldos table.
type saldoSeeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

// NewSaldoSeeder creates a new instance of the saldoSeeder, which is
// responsible for populating the saldos table with fake data.
//
// Args:
// db: a pointer to the database connection
// ctx: a context.Context object
// logger: a logger.LoggerInterface object
//
// Returns:
// a pointer to the saldoSeeder struct
func NewSaldoSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *saldoSeeder {
	return &saldoSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

// Seed populates the saldos table with fake data.
//
// It generates 10 saldos, 5 of which are active and 5 of which are trashed.
// The saldos are seeded with random totalBalance values, and the cardNumber
// is randomly selected from the cards table.
//
// If any errors occur during the seeding process, an error is returned.
//
// Returns:
// an error if any of the saldos fail to be created, otherwise nil
func (r *saldoSeeder) Seed() error {
	totalSaldos := 10
	activeSaldos := 5
	trashedSaldos := 5

	var cards []models.Card
	for i := 1; i <= totalSaldos; i++ {
		var card models.Card
		err := r.db.WithContext(r.ctx).Where("user_id = ?", int32(i)).First(&card).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				r.logger.Error("no card found for user", zap.Int("userID", i))
				continue
			}
			r.logger.Error("failed to get card for user", zap.Int("userID", i), zap.Error(err))
			return fmt.Errorf("failed to get card for user %d: %w", i, err)
		}
		cards = append(cards, card)
	}

	if len(cards) < totalSaldos {
		r.logger.Error("not enough cards to seed saldo", zap.Int("required", totalSaldos), zap.Int("available", len(cards)))
		return fmt.Errorf("not enough cards to seed saldo: required %d, got %d", totalSaldos, len(cards))
	}

	for i, card := range cards {
		saldo := &models.Saldo{
			CardNumber:   card.CardNumber,
			TotalBalance: int32(rand.Intn(9_000_000) + 1_000_000),
		}

		if err := r.db.WithContext(r.ctx).Create(saldo).Error; err != nil {
			r.logger.Error("failed to seed saldo", zap.Int("index", i), zap.String("card", card.CardNumber), zap.Error(err))
			return fmt.Errorf("failed to seed saldo for card %s: %w", card.CardNumber, err)
		}

		if i >= activeSaldos {
			now := time.Now()
			if err := r.db.WithContext(r.ctx).Model(&models.Saldo{}).
				Where("saldo_id = ? AND deleted_at IS NULL", saldo.SaldoID).
				Update("deleted_at", &now).Error; err != nil {
				r.logger.Error("failed to trash saldo", zap.Int("index", i), zap.String("card", card.CardNumber), zap.Error(err))
				return fmt.Errorf("failed to trash saldo %d for card %s: %w", i+1, card.CardNumber, err)
			}
		}
	}

	r.logger.Info("saldo seeded successfully",
		zap.Int("totalSaldos", totalSaldos),
		zap.Int("activeSaldos", activeSaldos),
		zap.Int("trashedSaldos", trashedSaldos))

	return nil
}
