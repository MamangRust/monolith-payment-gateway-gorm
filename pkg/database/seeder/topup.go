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

// topupSeeder is a struct that represents a seeder for the topups table.
type topupSeeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

// NewTopupSeeder creates a new instance of the topupSeeder, which is
// responsible for populating the topups table with fake data.
//
// Args:
// db: a pointer to the database connection
// ctx: a context.Context object
// logger: a logger.LoggerInterface object
//
// Returns:
// a pointer to the topupSeeder struct
func NewTopupSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *topupSeeder {
	return &topupSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

// Seed populates the topups table with fake data.
//
// It creates a total of 10 topups, with 5 of them being active and 5 of them being
// trashed. The topups are seeded with a random card number, topup amount, topup method,
// and topup time. The status of the topup is randomly set to one of the following:
// pending, success, or failed.
//
// If any errors occur during the seeding process, an error is returned.
//
// Returns:
// an error if any of the topups fail to be created, otherwise nil
func (r *topupSeeder) Seed() error {
	totalTopups := 10
	activeTopups := 5
	trashedTopups := 5

	var cards []models.Card

	for i := 1; i <= totalTopups; i++ {
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

	if len(cards) < totalTopups {
		r.logger.Error("not enough cards found to seed topups", zap.Int("found", len(cards)))
		return fmt.Errorf("not enough cards to seed topups")
	}

	topupMethods := []string{"Bank Alpha", "Bank Beta", "Bank Gamma"}
	statusOptions := []string{"pending", "success", "failed"}

	months := make([]time.Time, 12)
	currentYear := time.Now().Year()
	for i := 0; i < 12; i++ {
		months[i] = time.Date(currentYear, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
	}

	for i := 0; i < totalTopups; i++ {
		card := cards[i]
		cardNumber := card.CardNumber

		monthIndex := i % 12
		topupTime := months[monthIndex].Add(time.Duration(rand.Intn(28)) * 24 * time.Hour)

		topup := &models.Topup{
			CardNumber:  cardNumber,
			TopupAmount: int32(rand.Intn(10000000) + 1000000),
			TopupMethod: topupMethods[rand.Intn(len(topupMethods))],
			TopupTime:   topupTime,
		}

		if err := r.db.WithContext(r.ctx).Create(topup).Error; err != nil {
			r.logger.Error("failed to seed topup", zap.String("card", cardNumber), zap.Error(err))
			return fmt.Errorf("failed to seed topup for card %s: %w", cardNumber, err)
		}

		status := statusOptions[rand.Intn(len(statusOptions))]
		if err := r.db.WithContext(r.ctx).Model(&models.Topup{}).
			Where("topup_id = ?", topup.TopupID).
			Update("status", status).Error; err != nil {
			r.logger.Error("failed to update topup status", zap.Int("topupID", int(topup.TopupID)), zap.Error(err))
			return fmt.Errorf("failed to update status for topup %d: %w", topup.TopupID, err)
		}

		if i >= activeTopups {
			now := time.Now()
			if err := r.db.WithContext(r.ctx).Model(&models.Topup{}).
				Where("topup_id = ? AND deleted_at IS NULL", topup.TopupID).
				Update("deleted_at", &now).Error; err != nil {
				r.logger.Error("failed to trash topup", zap.Int("topup", i+1), zap.String("card", cardNumber), zap.Error(err))
				return fmt.Errorf("failed to trash topup %d for card %s: %w", i+1, cardNumber, err)
			}
		}
	}

	r.logger.Info("topup seeded successfully", zap.Int("totalTopups", totalTopups), zap.Int("activeTopups", activeTopups), zap.Int("trashedTopups", trashedTopups))
	return nil
}
