package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/repository"
	"gorm.io/gorm"
)

type userQueryRepository struct {
	db *gorm.DB
}

func NewUserQueryRepository(db *gorm.DB) UserQueryRepository {
	return &userQueryRepository{db: db}
}

func (r *userQueryRepository) FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.UserRow
	query := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, created_at, updated_at, COUNT(*) OVER () AS total_count").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		search := "%" + strings.TrimSpace(req.Search) + "%"
		query = query.Where("firstname ILIKE ? OR lastname ILIKE ? OR email ILIKE ?", search, search, search)
	}

	query = query.Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, user_errors.ErrFindAllUsers.WithInternal(err)
	}

	return results, nil
}

func (r *userQueryRepository) FindById(ctx context.Context, userID int) (*models.UserByIDRow, error) {
	if userID <= 0 {
		return nil, sharedErrors.NewBadRequestError("user ID must be greater than zero")
	}

	var result models.UserByIDRow
	if err := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, created_at, updated_at").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &result, nil
}

func (r *userQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.UserActiveRow
	query := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		search := "%" + strings.TrimSpace(req.Search) + "%"
		query = query.Where("firstname ILIKE ? OR lastname ILIKE ? OR email ILIKE ?", search, search, search)
	}

	query = query.Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, user_errors.ErrFindActiveUsers.WithInternal(err)
	}

	return results, nil
}

func (r *userQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.UserTrashedRow
	query := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count").
		Where("deleted_at IS NOT NULL")

	if req.Search != "" {
		search := "%" + strings.TrimSpace(req.Search) + "%"
		query = query.Where("firstname ILIKE ? OR lastname ILIKE ? OR email ILIKE ?", search, search, search)
	}

	query = query.Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, user_errors.ErrFindTrashedUsers.WithInternal(err)
	}

	return results, nil
}

func (r *userQueryRepository) FindByEmail(ctx context.Context, email string) (*models.UserByEmailRow, error) {
	var result models.UserByEmailRow
	if err := r.db.WithContext(ctx).Table("users").
		Select("user_id, email").
		Where("email = ? AND deleted_at IS NULL", email).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &result, nil
}
