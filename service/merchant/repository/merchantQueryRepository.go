package repository

import (
	"fmt"
	"context"
	"errors"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	"gorm.io/gorm"
)

type merchantQueryRepository struct {
	db *gorm.DB
}

func NewMerchantQueryRepository(db *gorm.DB) MerchantQueryRepository {
	return &merchantQueryRepository{db: db}
}

func (r *merchantQueryRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var totalCount int64
	countQuery := r.db.WithContext(ctx).Table("merchants").Where("deleted_at IS NULL")
	if req.Search != "" {
		countQuery = countQuery.Where("(name ILIKE ? OR api_key ILIKE ? OR status ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, merchant_errors.ErrFindAllMerchantsFailed.WithInternal(err)
	}

	var merchants []*models.MerchantListRow
	query := r.db.WithContext(ctx).Table("merchants").
		Select(fmt.Sprintf("merchants.*, %d as total_count", totalCount)).
		Where("merchants.deleted_at IS NULL")
	if req.Search != "" {
		query = query.Where("(merchants.name ILIKE ? OR merchants.api_key ILIKE ? OR merchants.status ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := query.Order("merchants.merchant_id DESC").Limit(req.PageSize).Offset(offset).Scan(&merchants).Error; err != nil {
		return nil, merchant_errors.ErrFindAllMerchantsFailed.WithInternal(err)
	}
	return merchants, nil
}

func (r *merchantQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListWithDeletedRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var totalCount int64
	countQuery := r.db.WithContext(ctx).Table("merchants").Where("deleted_at IS NULL")
	if req.Search != "" {
		countQuery = countQuery.Where("(name ILIKE ? OR api_key ILIKE ? OR status ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, merchant_errors.ErrFindActiveMerchantsFailed.WithInternal(err)
	}

	var merchants []*models.MerchantListWithDeletedRow
	query := r.db.WithContext(ctx).Table("merchants").
		Select(fmt.Sprintf("merchants.*, %d as total_count", totalCount)).
		Where("merchants.deleted_at IS NULL")
	if req.Search != "" {
		query = query.Where("(merchants.name ILIKE ? OR merchants.api_key ILIKE ? OR merchants.status ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := query.Order("merchants.merchant_id DESC").Limit(req.PageSize).Offset(offset).Scan(&merchants).Error; err != nil {
		return nil, merchant_errors.ErrFindActiveMerchantsFailed.WithInternal(err)
	}
	return merchants, nil
}

func (r *merchantQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListWithDeletedRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var totalCount int64
	countQuery := r.db.WithContext(ctx).Table("merchants").Where("deleted_at IS NOT NULL")
	if req.Search != "" {
		countQuery = countQuery.Where("(name ILIKE ? OR api_key ILIKE ? OR status ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, merchant_errors.ErrFindTrashedMerchantsFailed.WithInternal(err)
	}

	var merchants []*models.MerchantListWithDeletedRow
	query := r.db.WithContext(ctx).Table("merchants").
		Select(fmt.Sprintf("merchants.*, %d as total_count", totalCount)).
		Where("merchants.deleted_at IS NOT NULL")
	if req.Search != "" {
		query = query.Where("(merchants.name ILIKE ? OR merchants.api_key ILIKE ? OR merchants.status ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := query.Order("merchants.merchant_id DESC").Limit(req.PageSize).Offset(offset).Scan(&merchants).Error; err != nil {
		return nil, merchant_errors.ErrFindTrashedMerchantsFailed.WithInternal(err)
	}
	return merchants, nil
}

func (r *merchantQueryRepository) FindByMerchantId(ctx context.Context, merchant_id int) (*models.MerchantAllFieldsRow, error) {
	if merchant_id <= 0 {
		return nil, sharedErrors.NewBadRequestError("merchant ID must be greater than zero")
	}
	var merchant models.MerchantAllFieldsRow
	err := r.db.WithContext(ctx).Table("merchants").
		Where("merchant_id = ? AND deleted_at IS NULL", merchant_id).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantQueryRepository) FindByApiKey(ctx context.Context, api_key string) (*models.MerchantAllFieldsRow, error) {
	var merchant models.MerchantAllFieldsRow
	err := r.db.WithContext(ctx).Table("merchants").
		Where("api_key = ? AND deleted_at IS NULL", api_key).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantQueryRepository) FindByName(ctx context.Context, name string) (*models.MerchantAllFieldsRow, error) {
	var merchant models.MerchantAllFieldsRow
	err := r.db.WithContext(ctx).Table("merchants").
		Where("name = ? AND deleted_at IS NULL", name).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantQueryRepository) FindByMerchantUserId(ctx context.Context, user_id int) ([]*models.MerchantListByUserRow, error) {
	var merchants []*models.MerchantListByUserRow
	err := r.db.WithContext(ctx).Table("merchants").
		Where("user_id = ? AND deleted_at IS NULL", user_id).
		Order("merchant_id DESC").Scan(&merchants).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return merchants, nil
}
