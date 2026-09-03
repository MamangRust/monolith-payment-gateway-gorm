package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// RoleQueryService is an interface for querying role records
type RoleQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleRow, *int, error)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleActiveRow, *int, error)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleTrashedRow, *int, error)
	FindById(ctx context.Context, role_id int) (*models.Role, error)
	FindByUserId(ctx context.Context, id int) ([]*models.Role, error)
}

// RoleCommandService is an interface for creating, updating, and deleting role records
type RoleCommandService interface {
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*models.Role, error)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*models.Role, error)
	TrashedRole(ctx context.Context, role_id int) (*models.Role, error)
	RestoreRole(ctx context.Context, role_id int) (*models.Role, error)
	DeleteRolePermanent(ctx context.Context, role_id int) (bool, error)

	RestoreAllRole(ctx context.Context) (bool, error)
	DeleteAllRolePermanent(ctx context.Context) (bool, error)
}
