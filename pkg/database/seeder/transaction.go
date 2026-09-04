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

// transactionSeeder is a struct that represents a seeder for the transactions table.
type transactionSeeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

// NewTransactionSeeder creates a new instance of the transactionSeeder, which is
// responsible for populating the transactions table with fake data.
//
// Args:
// db: a pointer to the database connection
// ctx: a context.Context object
// logger: a logger.LoggerInterface object
//
// Returns:
// a pointer to the transactionSeeder struct
func NewTransactionSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *transactionSeeder {
	return &transactionSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

// Seed populates the transactions table with fake data.
//
// It generates a total of 10 transactions, with 5 of them being active and 5 of them being
// trashed. The transactions are seeded with random payment methods, statuses, and transaction
// times. The card numbers are randomly selected from the cards table. The merchant IDs are
// randomly selected from the merchants table.
//
// If any of the transactions fail to be created, the function returns an error.
//
// Returns:
// an error if any of the transactions fail to be created, otherwise nil
func (r *transactionSeeder) Seed() error {
	total := 10
	active := 5
	trashed := 5

	paymentMethods := []string{"Bank Alpha", "Bank Beta", "Bank Gamma"}
	statusOptions := []string{"pending", "success", "failed"}

	var cards []models.Card
	for i := 1; i <= total; i++ {
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

	if len(cards) < total {
		r.logger.Error("not enough cards for transaction seeding", zap.Int("required", total), zap.Int("available", len(cards)))
		return fmt.Errorf("not enough cards for transaction seeding: required %d, got %d", total, len(cards))
	}

	var merchants []models.Merchant
	if err := r.db.WithContext(r.ctx).Limit(total).Find(&merchants).Error; err != nil {
		r.logger.Error("failed to get merchant list", zap.Error(err))
		return fmt.Errorf("failed to get merchant list: %w", err)
	}

	if len(merchants) < total {
		r.logger.Error("not enough merchants for transaction seeding", zap.Int("required", total), zap.Int("available", len(merchants)))
		return fmt.Errorf("not enough merchants: required %d, got %d", total, len(merchants))
	}

	months := make([]time.Time, 12)
	currentYear := time.Now().Year()
	for i := 0; i < 12; i++ {
		months[i] = time.Date(currentYear, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
	}

	for i := 0; i < total; i++ {
		card := cards[i]
		merchant := merchants[i%len(merchants)]
		paymentMethod := paymentMethods[i%len(paymentMethods)]
		status := statusOptions[i%len(statusOptions)]

		monthIndex := i % 12
		transactionTime := months[monthIndex].Add(time.Duration(rand.Intn(28)) * 24 * time.Hour)

		transaction := &models.Transaction{
			CardNumber:      card.CardNumber,
			Amount:          int32(rand.Intn(1000000-50000) + 50000),
			PaymentMethod:   paymentMethod,
			MerchantID:      merchant.MerchantID,
			TransactionTime: transactionTime,
		}

		if err := r.db.WithContext(r.ctx).Create(transaction).Error; err != nil {
			r.logger.Error("failed to seed transaction", zap.Int("index", i), zap.Error(err))
			return fmt.Errorf("failed to seed transaction %d: %w", i, err)
		}

		if err := r.db.WithContext(r.ctx).Model(&models.Transaction{}).
			Where("transaction_id = ?", transaction.TransactionID).
			Update("status", status).Error; err != nil {
			r.logger.Error("failed to update transaction status", zap.Int("transactionID", int(transaction.TransactionID)), zap.String("status", status), zap.Error(err))
			return fmt.Errorf("failed to update status for transaction ID %d: %w", transaction.TransactionID, err)
		}

		if i >= active {
			now := time.Now()
			if err := r.db.WithContext(r.ctx).Model(&models.Transaction{}).
				Where("transaction_id = ? AND deleted_at IS NULL", transaction.TransactionID).
				Update("deleted_at", &now).Error; err != nil {
				r.logger.Error("failed to trash transaction", zap.Int("index", i), zap.Error(err))
				return fmt.Errorf("failed to trash transaction %d: %w", i, err)
			}
		}
	}

	r.logger.Info("transaction seeded successfully", zap.Int("total", total), zap.Int("active", active), zap.Int("trashed", trashed))
	return nil
}
