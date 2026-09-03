package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	userrole_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_role_errors/repository"
	"gorm.io/gorm"
)

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*models.UserRole, error) {
	userRole := models.UserRole{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	}
	if err := r.db.WithContext(ctx).Create(&userRole).Error; err != nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}
	return &userRole, nil
}

func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", req.UserId, req.RoleId).
		Delete(&models.UserRole{}).Error
	if err != nil {
		return userrole_errors.ErrRemoveRole.WithInternal(err)
	}
	return nil
}
