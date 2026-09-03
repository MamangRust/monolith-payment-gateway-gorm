package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// UserQueryService handles query operations related to user data.
type UserQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserRow, *int, error)
	FindByID(ctx context.Context, id int) (*models.UserByIDRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserTrashedRow, *int, error)
}

// UserCommandService handles command operations related to user management.
type UserCommandService interface {
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*models.CreateUserRow, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*models.UpdateUserRow, error)
	TrashedUser(ctx context.Context, user_id int) (*models.TrashUserRow, error)
	RestoreUser(ctx context.Context, user_id int) (*models.RestoreUserRow, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)

	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}
