package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type RoleQueryCache interface {
	SetCachedRoles(ctx context.Context, req *requests.FindAllRoles, data []*models.RoleRow, total *int)
	GetCachedRoles(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleRow, *int, bool)

	GetCachedRoleById(ctx context.Context, id int) (*models.Role, bool)
	SetCachedRoleById(ctx context.Context, id int, data *models.Role)

	GetCachedRoleByUserId(ctx context.Context, userId int) ([]*models.Role, bool)
	SetCachedRoleByUserId(ctx context.Context, userId int, data []*models.Role)

	GetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleActiveRow, *int, bool)
	SetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles, data []*models.RoleActiveRow, total *int)

	GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleTrashedRow, *int, bool)
	SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles, data []*models.RoleTrashedRow, total *int)
}

type RoleCommandCache interface {
	DeleteCachedRole(ctx context.Context, id int)
}
