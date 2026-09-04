package seeder

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// userSeeder is a struct that represents a seeder for the users table.
type userSeeder struct {
	db     *gorm.DB
	hash   hash.HashPassword
	ctx    context.Context
	logger logger.LoggerInterface
}

// NewUserSeeder creates a new instance of the userSeeder, which is
// responsible for populating the users table with fake data.
//
// Args:
// db: a pointer to the database connection
// ctx: a context.Context object
// hash: a hash.HashPassword object
// logger: a logger.LoggerInterface object
//
// Returns:
// a pointer to the userSeeder struct
func NewUserSeeder(db *gorm.DB, ctx context.Context, hash hash.HashPassword, logger logger.LoggerInterface) *userSeeder {
	return &userSeeder{
		db:     db,
		hash:   hash,
		ctx:    ctx,
		logger: logger,
	}
}

// Seed populates the users table with fake data.
//
// It creates a total of 15 users, with 10 of them being active and 5 of them being
// trashed. The active users have a status of "active", while the trashed users have a
// status of "deactive". The user names are generated randomly from the combination of
// the adjectives and nouns provided. The API key is generated using the GenerateApiKey
// function from the api-key package.
//
// If any errors occur during the seeding process, an error is returned.
//
// Returns:
// an error if any of the users fail to be created, otherwise nil
func (r *userSeeder) Seed() error {
	totalUsers := 15
	activeUsers := 10
	trashedUsers := totalUsers - activeUsers

	randomRoles := []string{
		"Super Admin",
		"Admin",
		"Merchant Admin",
		"Merchant Operator",
		"Finance",
		"Compliance",
		"Auditor",
		"Support",
		"Viewer",
		"User",
	}

	for i := 1; i <= totalUsers; i++ {
		email := fmt.Sprintf("user_%s@example.com", uuid.NewString())

		hash, err := r.hash.HashPassword("password")
		if err != nil {
			r.logger.Error("failed to hash password", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to hash password for user %d: %w", i, err)
		}

		isVerified := true

		user := &models.User{
			Firstname:        fmt.Sprintf("User%d", i),
			Lastname:         fmt.Sprintf("Last%d", i),
			Email:            email,
			Password:         hash,
			VerificationCode: uuid.NewString(),
			IsVerified:       &isVerified,
		}

		if err := r.db.WithContext(r.ctx).Create(user).Error; err != nil {
			r.logger.Error("failed to create user", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to create user %d: %w", i, err)
		}

		randomRole := randomRoles[rand.Intn(len(randomRoles))]
		var role models.Role
		if err := r.db.WithContext(r.ctx).Where("role_name = ?", randomRole).First(&role).Error; err != nil {
			r.logger.Error("failed to get role", zap.String("role", randomRole), zap.Error(err))
			return fmt.Errorf("failed to get role %s: %w", randomRole, err)
		}

		if err := r.db.WithContext(r.ctx).Create(&models.UserRole{
			RoleID: role.RoleID,
			UserID: user.UserID,
		}).Error; err != nil {
			r.logger.Error("failed to assign role to user", zap.Int("userID", int(user.UserID)), zap.String("role", randomRole), zap.Error(err))
			return fmt.Errorf("failed to assign role %s to user %d: %w", randomRole, user.UserID, err)
		}

		if i > activeUsers {
			now := time.Now()
			if err := r.db.WithContext(r.ctx).Model(&models.User{}).
				Where("user_id = ? AND deleted_at IS NULL", user.UserID).
				Update("deleted_at", &now).Error; err != nil {
				r.logger.Error("failed to trash user", zap.Int("user", i), zap.Error(err))
				return fmt.Errorf("failed to trash user %d: %w", i, err)
			}
		}
	}

	r.logger.Info("user seeded successfully", zap.Int("totalUsers", totalUsers), zap.Int("activeUsers", activeUsers), zap.Int("trashedUsers", trashedUsers))

	return nil
}
