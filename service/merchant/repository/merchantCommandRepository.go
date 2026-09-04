package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/api-key"
	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	"gorm.io/gorm"
)

type merchantCommandRepository struct {
	db *gorm.DB
}

func NewMerchantCommandRepository(db *gorm.DB) MerchantCommandRepository {
	return &merchantCommandRepository{db: db}
}

func (r *merchantCommandRepository) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*models.MerchantAllFieldsRow, error) {
	apiKey, err := apikey.GenerateApiKey()
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	merchant := models.Merchant{
		Name:   request.Name,
		ApiKey: apiKey,
		UserID: int32(request.UserID),
		Status: "inactive",
	}
	if err := r.db.WithContext(ctx).Create(&merchant).Error; err != nil {
		return nil, sharedErrors.ErrConstraintOrFailed(err, "Merchant", "create merchant")
	}
	return &models.MerchantAllFieldsRow{
		MerchantID: merchant.MerchantID,
		MerchantNo: merchant.MerchantNo,
		Name:       merchant.Name,
		ApiKey:     merchant.ApiKey,
		UserID:     merchant.UserID,
		Status:     merchant.Status,
		CreatedAt:  merchant.CreatedAt,
		UpdatedAt:  merchant.UpdatedAt,
	}, nil
}

func (r *merchantCommandRepository) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*models.MerchantAllFieldsRow, error) {
	result := r.db.WithContext(ctx).Model(&models.Merchant{}).
		Where("merchant_id = ? AND deleted_at IS NULL", *request.MerchantID).
		Updates(map[string]interface{}{
			"name":    request.Name,
			"user_id": int32(request.UserID),
			"status":  request.Status,
		})
	if result.Error != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(result.Error, "Merchant", "update merchant")
	}
	var merchant models.MerchantAllFieldsRow
	if err := r.db.WithContext(ctx).Table("merchants").Where("merchant_id = ?", *request.MerchantID).Scan(&merchant).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Merchant", "update merchant")
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*models.MerchantAllFieldsRow, error) {
	result := r.db.WithContext(ctx).Model(&models.Merchant{}).
		Where("merchant_id = ? AND deleted_at IS NULL", *request.MerchantID).
		Update("status", request.Status)
	if result.Error != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(result.Error, "Merchant", "update merchant status")
	}
	var merchant models.MerchantAllFieldsRow
	if err := r.db.WithContext(ctx).Table("merchants").Where("merchant_id = ?", *request.MerchantID).Scan(&merchant).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Merchant", "update merchant status")
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) TrashedMerchant(ctx context.Context, merchant_id int) (*models.Merchant, error) {
	result := r.db.WithContext(ctx).Where("merchant_id = ? AND deleted_at IS NULL", merchant_id).Delete(&models.Merchant{})
	if result.Error != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(result.Error, "Merchant", "trash merchant")
	}
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).Unscoped().Where("merchant_id = ?", merchant_id).First(&merchant).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Merchant", "trash merchant")
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) RestoreMerchant(ctx context.Context, merchant_id int) (*models.Merchant, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.Merchant{}).
		Where("merchant_id = ?", merchant_id).Update("deleted_at", nil)
	if result.Error != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(result.Error, "Merchant", "restore merchant")
	}
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).Where("merchant_id = ?", merchant_id).First(&merchant).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Merchant", "restore merchant")
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("merchant_id = ?", merchant_id).Delete(&models.Merchant{})
	if result.Error != nil {
		return false, sharedErrors.ErrNoRowsOrFailed(result.Error, "Merchant", "delete merchant permanently")
	}
	return true, nil
}

func (r *merchantCommandRepository) RestoreAllMerchant(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.Merchant{}).
		Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, merchant_errors.ErrRestoreAllMerchantFailed.WithInternal(result.Error)
	}
	return true, nil
}

func (r *merchantCommandRepository) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.Merchant{})
	if result.Error != nil {
		return false, merchant_errors.ErrDeleteAllMerchantPermanentFailed.WithInternal(result.Error)
	}
	return true, nil
}
