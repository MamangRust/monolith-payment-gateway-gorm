package repository

import (
	"context"
	"strings"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/repository"
	"gorm.io/gorm"
)

type roleQueryRepository struct {
	db *gorm.DB
}

func NewRoleQueryRepository(db *gorm.DB) RoleQueryRepository {
	return &roleQueryRepository{db: db}
}

func (r *roleQueryRepository) FindAllRoles(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.RoleRow
	query := r.db.WithContext(ctx).Table("roles").
		Select("role_id, role_name, created_at, updated_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("role_name ILIKE ?", "%"+req.Search+"%")
	}

	query = query.Order("created_at ASC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, role_errors.ErrFindAllRoles.WithInternal(err)
	}

	return results, nil
}

func (r *roleQueryRepository) FindById(ctx context.Context, id int) (*models.Role, error) {
	if id <= 0 {
		return nil, sharedErrors.NewBadRequestError("role ID must be greater than zero")
	}

	var role models.Role
	if err := r.db.WithContext(ctx).Where("role_id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &role, nil
}

func (r *roleQueryRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).Where("role_name = ?", name).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &role, nil
}

func (r *roleQueryRepository) FindByUserId(ctx context.Context, userID int) ([]*models.Role, error) {
	var roles []*models.Role
	if err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Order("roles.created_at ASC").
		Find(&roles).Error; err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}

	return roles, nil
}

func (r *roleQueryRepository) FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.RoleActiveRow
	query := r.db.WithContext(ctx).Table("roles").
		Select("role_id, role_name, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		query = query.Where("role_name ILIKE ?", "%"+req.Search+"%")
	}

	query = query.Order("created_at ASC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, role_errors.ErrFindActiveRoles.WithInternal(err)
	}

	return results, nil
}

func (r *roleQueryRepository) FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.RoleTrashedRow
	query := r.db.WithContext(ctx).Table("roles").
		Select("role_id, role_name, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count").
		Where("deleted_at IS NOT NULL")

	if req.Search != "" {
		query = query.Where("role_name ILIKE ?", "%"+strings.TrimSpace(req.Search)+"%")
	}

	query = query.Order("deleted_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, role_errors.ErrFindTrashedRoles.WithInternal(err)
	}

	return results, nil
}
