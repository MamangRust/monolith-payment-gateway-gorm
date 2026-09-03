package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userCommandRepository struct {
	db *gorm.DB
}

func NewUserCommandRepository(db *gorm.DB) UserCommandRepository {
	return &userCommandRepository{db: db}
}

func (r *userCommandRepository) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*models.CreateUserRow, error) {
	verifyCode := uuid.New().String()
	verified := false

	user := &models.User{
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: verifyCode,
		IsVerified:       &verified,
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, sharedErrors.ErrConstraintOrFailed(err, "User", "create user")
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

func (r *userCommandRepository) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*models.UpdateUserRow, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", *request.UserID).
		First(&user).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "User", "update user")
	}

	user.Firstname = request.FirstName
	user.Lastname = request.LastName
	user.Email = request.Email
	user.Password = request.Password
	user.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "User", "update user")
	}

	return &models.UpdateUserRow{
		UserID:    user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *userCommandRepository) TrashedUser(ctx context.Context, userID int) (*models.TrashUserRow, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		First(&user).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "User", "trash user")
	}

	now := time.Now()
	user.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	user.UpdatedAt = now

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "User", "trash user")
	}

	return &models.TrashUserRow{
		UserID:    user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: &user.DeletedAt.Time,
	}, nil
}

func (r *userCommandRepository) RestoreUser(ctx context.Context, userID int) (*models.RestoreUserRow, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Unscoped().
		Where("user_id = ? AND deleted_at IS NOT NULL", userID).
		First(&user).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "User", "restore user")
	}

	user.DeletedAt = gorm.DeletedAt{Valid: false}
	user.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Unscoped().Save(&user).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "User", "restore user")
	}

	return &models.RestoreUserRow{
		UserID:    user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: nil,
	}, nil
}

func (r *userCommandRepository) DeleteUserPermanent(ctx context.Context, userID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("user_id = ? AND deleted_at IS NOT NULL", userID).
		Delete(&models.User{})

	if result.Error != nil {
		return false, sharedErrors.ErrNoRowsOrFailed(result.Error, "User", "delete user")
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *userCommandRepository) RestoreAllUser(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Model(&models.User{}).
		Where("deleted_at IS NOT NULL").
		Update("deleted_at", nil)

	if result.Error != nil {
		return false, user_errors.ErrRestoreAllUsers.WithInternal(result.Error)
	}

	return true, nil
}

func (r *userCommandRepository) DeleteAllUserPermanent(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("deleted_at IS NOT NULL").
		Delete(&models.User{})

	if result.Error != nil {
		return false, user_errors.ErrDeleteAllUsers.WithInternal(result.Error)
	}

	return true, nil
}
