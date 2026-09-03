package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	refreshtoken_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/refresh_token_errors/repository"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var res models.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&res).Error
	if err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Refresh Token", "find refresh token by token")
	}
	return &res, nil
}

func (r *refreshTokenRepository) FindByUserId(ctx context.Context, user_id int) (*models.RefreshToken, error) {
	var res models.RefreshToken
	err := r.db.WithContext(ctx).Where("user_id = ?", user_id).First(&res).Error
	if err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Refresh Token", "find refresh token by user ID")
	}
	return &res, nil
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*models.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate.WithInternal(err)
	}
	res := models.RefreshToken{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	}
	if err := r.db.WithContext(ctx).Create(&res).Error; err != nil {
		return nil, refreshtoken_errors.ErrCreateRefreshToken.WithInternal(err)
	}
	return &res, nil
}

func (r *refreshTokenRepository) UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*models.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate.WithInternal(err)
	}
	var res models.RefreshToken
	err = r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("user_id = ?", req.UserId).
		Updates(map[string]interface{}{
			"token":      req.Token,
			"expiration": expirationTime,
		}).First(&res).Error
	if err != nil {
		return nil, refreshtoken_errors.ErrUpdateRefreshToken.WithInternal(err)
	}
	// Re-fetch
	if err := r.db.WithContext(ctx).Where("user_id = ?", req.UserId).First(&res).Error; err != nil {
		return nil, refreshtoken_errors.ErrUpdateRefreshToken.WithInternal(err)
	}
	return &res, nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	err := r.db.WithContext(ctx).Where("token = ?", token).Delete(&models.RefreshToken{}).Error
	if err != nil {
		return refreshtoken_errors.ErrDeleteRefreshToken.WithInternal(err)
	}
	return nil
}

func (r *refreshTokenRepository) DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", user_id).Delete(&models.RefreshToken{}).Error
	if err != nil {
		return refreshtoken_errors.ErrDeleteByUserID.WithInternal(err)
	}
	return nil
}
