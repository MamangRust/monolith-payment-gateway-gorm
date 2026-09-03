package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type UserQueryRepository interface {
	FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserTrashedRow, error)
	FindById(ctx context.Context, user_id int) (*models.UserByIDRow, error)
	FindByEmail(ctx context.Context, email string) (*models.UserByEmailRow, error)
}

type UserCommandRepository interface {
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*models.CreateUserRow, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*models.UpdateUserRow, error)
	TrashedUser(ctx context.Context, user_id int) (*models.TrashUserRow, error)
	RestoreUser(ctx context.Context, user_id int) (*models.RestoreUserRow, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)
	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}

type RoleRepository interface {
	FindById(ctx context.Context, role_id int) (*models.Role, error)
	FindByName(ctx context.Context, name string) (*models.Role, error)
}
