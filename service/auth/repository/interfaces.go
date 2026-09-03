package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.UserByEmailRow, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*models.UserByEmailWithPasswordRow, error)
	FindById(ctx context.Context, user_id int) (*models.UserByIDRow, error)
	CreateUser(ctx context.Context, request *requests.RegisterRequest) (*models.CreateUserRow, error)
	UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*models.UserIsVerifiedRow, error)
	UpdateUserPassword(ctx context.Context, user_id int, password string) (*models.UserPasswordRow, error)
	FindByVerificationCode(ctx context.Context, verification_code string) (*models.UserByVerificationCodeRow, error)
}

type ResetTokenRepository interface {
	FindByToken(ctx context.Context, code string) (*models.ResetTokenRow, error)
	CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*models.ResetTokenRow, error)
	DeleteResetToken(ctx context.Context, user_id int) error
}

type RefreshTokenRepository interface {
	FindByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	FindByUserId(ctx context.Context, user_id int) (*models.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*models.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error
}

type UserRoleRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*models.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

type RoleRepository interface {
	FindById(ctx context.Context, id int) (*models.Role, error)
	FindByName(ctx context.Context, name string) (*models.Role, error)
}
