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

// withdrawSeeder is a struct that represents a seeder for the withdraws table.
type withdrawSeeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

// NewWithdrawSeeder creates a new instance of the withdrawSeeder, which is
// responsible for populating the withdraws table with fake data.
//
// Args:
// db: a pointer to the database connection
// ctx: a context.Context object
// logger: a logger.LoggerInterface object
//
// Returns:
// a pointer to the withdrawSeeder struct
func NewWithdrawSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *withdrawSeeder {
	return &withdrawSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

// Seed populates the withdraws table with fake data.
//
// It generates a total of 10 withdraws, with 5 of them being active and 5 of them being
// trashed. The withdraws are seeded with random card numbers, withdraw amounts, and
// withdraw times. The status of each withdraw is randomly set to one of the following:
// pending, success, or failed. The card numbers are randomly selected from the cards
// table.
//
// If any errors occur during the seeding process, an error is returned.
//
// Returns:
// an error if any of the withdraws fail to be created or updated, otherwise nil
func (r *withdrawSeeder) Seed() error {
	total := 10
	active := 5
	trashed := total - active

	var cards []models.Card
	for i := 1; i <= total; i++ {
		var card models.Card
		err := r.db.WithContext(r.ctx).Where("user_id = ?", int32(i)).First(&card).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				r.logger.Debug("failed to get card for user", zap.Int("userID", i), zap.Error(err))
			}
			continue
		}
		cards = append(cards, card)
	}

	if len(cards) < total {
		r.logger.Error("not enough cards for withdraw seeding", zap.Int("required", total), zap.Int("available", len(cards)))
		return fmt.Errorf("not enough cards: required %d, got %d", total, len(cards))
	}

	statusOptions := []string{"pending", "success", "failed"}

	months := make([]time.Time, 12)
	currentYear := time.Now().Year()
	for i := 0; i < 12; i++ {
		months[i] = time.Date(currentYear, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
	}

	for i := 0; i < total; i++ {
		card := cards[i]
		status := statusOptions[rand.Intn(len(statusOptions))]

		monthIndex := i % 12
		withdrawTime := months[monthIndex].Add(time.Duration(rand.Intn(28)) * 24 * time.Hour)

		withdraw := &models.Withdraw{
			CardNumber:     card.CardNumber,
			WithdrawAmount: int32(rand.Intn(1000000) + 50000),
			WithdrawTime:   withdrawTime,
		}

		if err := r.db.WithContext(r.ctx).Create(withdraw).Error; err != nil {
			r.logger.Error("failed to seed withdraw", zap.Int("index", i), zap.Error(err))
			return fmt.Errorf("failed to create withdraw %d: %w", i, err)
		}

		if err := r.db.WithContext(r.ctx).Model(&models.Withdraw{}).
			Where("withdraw_id = ?", withdraw.WithdrawID).
			Update("status", status).Error; err != nil {
			r.logger.Error("failed to update withdraw status", zap.Int("withdraw.id", int(withdraw.WithdrawID)), zap.Error(err))
			return fmt.Errorf("failed to update status: %w", err)
		}

		if i >= active {
			now := time.Now()
			if err := r.db.WithContext(r.ctx).Model(&models.Withdraw{}).
				Where("withdraw_id = ? AND deleted_at IS NULL", withdraw.WithdrawID).
				Update("deleted_at", &now).Error; err != nil {
				r.logger.Error("failed to trash withdraw", zap.Int("withdraw.id", int(withdraw.WithdrawID)), zap.Error(err))
				return fmt.Errorf("failed to trash withdraw %d: %w", i, err)
			}
		}
	}

	r.logger.Info("withdraw seeding completed",
		zap.Int("total", total),
		zap.Int("active", active),
		zap.Int("trashed", trashed))

	return nil
}
