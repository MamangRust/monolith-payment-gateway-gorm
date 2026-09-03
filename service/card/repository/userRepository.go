package repository

import (
	"context"
	"errors"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindById(ctx context.Context, user_id int) (*models.UserByIDRow, error) {
	var user models.UserByIDRow
	err := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, created_at, updated_at").
		Where("user_id = ?", user_id).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUserNotFound
	}

	return &user, nil
}
