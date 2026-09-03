package repository

import (
	"context"
	"errors"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindById(ctx context.Context, user_id int) (*models.UserByIDRow, error) {
	if user_id <= 0 {
		return nil, sharedErrors.NewBadRequestError("user ID must be greater than zero")
	}
	var user models.UserByIDRow
	err := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, created_at, updated_at").
		Where("user_id = ?", user_id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.UserByEmailRow, error) {
	var user models.UserByEmailRow
	err := r.db.WithContext(ctx).Table("users").
		Select("user_id, email").
		Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmailAndVerify(ctx context.Context, email string) (*models.UserByEmailWithPasswordRow, error) {
	var user models.UserByEmailWithPasswordRow
	err := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, created_at, updated_at").
		Where("email = ? AND is_verified = true", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userRepository) FindByVerificationCode(ctx context.Context, verification_code string) (*models.UserByVerificationCodeRow, error) {
	var user models.UserByVerificationCodeRow
	err := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, created_at, updated_at").
		Where("verification_code = ?", verification_code).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, request *requests.RegisterRequest) (*models.CreateUserRow, error) {
	user := models.User{
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: request.VerifiedCode,
		IsVerified:       &request.IsVerified,
	}
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}
	return &models.CreateUserRow{
		UserID:    user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *userRepository) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*models.UserIsVerifiedRow, error) {
	result := r.db.WithContext(ctx).Model(&models.User{}).
		Where("user_id = ?", user_id).
		Update("is_verified", is_verified)
	if result.Error != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(result.Error, "User", "update user verification")
	}
	var user models.UserIsVerifiedRow
	if err := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, created_at, updated_at").
		Where("user_id = ?", user_id).First(&user).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "User", "update user verification")
	}
	return &user, nil
}

func (r *userRepository) UpdateUserPassword(ctx context.Context, user_id int, password string) (*models.UserPasswordRow, error) {
	result := r.db.WithContext(ctx).Model(&models.User{}).
		Where("user_id = ?", user_id).
		Update("password", password)
	if result.Error != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(result.Error, "User", "update user password")
	}
	var user models.UserPasswordRow
	if err := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, created_at, updated_at").
		Where("user_id = ?", user_id).First(&user).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "User", "update user password")
	}
	return &user, nil
}
