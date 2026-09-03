package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type RoleQueryRepository interface {
	FindAllRoles(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleRow, error)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleActiveRow, error)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleTrashedRow, error)
	FindById(ctx context.Context, role_id int) (*models.Role, error)
	FindByName(ctx context.Context, name string) (*models.Role, error)
	FindByUserId(ctx context.Context, user_id int) ([]*models.Role, error)
}

type RoleCommandRepository interface {
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*models.Role, error)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*models.Role, error)
	TrashedRole(ctx context.Context, role_id int) (*models.Role, error)
	RestoreRole(ctx context.Context, role_id int) (*models.Role, error)
	DeleteRolePermanent(ctx context.Context, role_id int) (bool, error)
	RestoreAllRole(ctx context.Context) (bool, error)
	DeleteAllRolePermanent(ctx context.Context) (bool, error)
}
