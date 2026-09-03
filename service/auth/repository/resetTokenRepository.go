package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	refresh_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/refresh_token_errors/repository"
	"gorm.io/gorm"
)

type resetTokenRepository struct {
	db *gorm.DB
}

func NewResetTokenRepository(db *gorm.DB) ResetTokenRepository {
	return &resetTokenRepository{db: db}
}

func (r *resetTokenRepository) FindByToken(ctx context.Context, code string) (*models.ResetTokenRow, error) {
	var token models.ResetToken
	err := r.db.WithContext(ctx).Where("token = ?", code).First(&token).Error
	if err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Reset Token", "find reset token")
	}
	return &models.ResetTokenRow{
		UserID:     token.UserID,
		Token:      token.Token,
		ExpiryDate: token.ExpiryDate,
	}, nil
}

func (r *resetTokenRepository) CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*models.ResetTokenRow, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	token := models.ResetToken{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	}
	if err := r.db.WithContext(ctx).Create(&token).Error; err != nil {
		return nil, refresh_errors.ErrCreateRefreshToken.WithInternal(err)
	}
	return &models.ResetTokenRow{
		UserID:     token.UserID,
		Token:      token.Token,
		ExpiryDate: token.ExpiryDate,
	}, nil
}

func (r *resetTokenRepository) DeleteResetToken(ctx context.Context, user_id int) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", user_id).Delete(&models.ResetToken{}).Error
	if err != nil {
		return refresh_errors.ErrDeleteByUserID.WithInternal(err)
	}
	return nil
}
