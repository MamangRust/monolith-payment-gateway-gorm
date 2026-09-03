package service

import (
	"context"
	"time"

	mencache "github.com/MamangRust/monolith-payment-gateway-auth/redis"
	"github.com/MamangRust/monolith-payment-gateway-auth/repository"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	randomstring "github.com/MamangRust/monolith-payment-gateway-pkg/random_string"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-payment-gateway-shared/errorhandler"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"

	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type RegisterServiceDeps struct {
	Cache mencache.RegisterCache

	User repository.UserRepository

	Role repository.RoleRepository

	UserRole repository.UserRoleRepository

	Hash hash.HashPassword

	Logger logger.LoggerInterface

	Observability observability.TraceLoggerObservability
}

type registerService struct {
	mencache mencache.RegisterCache

	user repository.UserRepository

	role repository.RoleRepository

	userRole repository.UserRoleRepository

	hash hash.HashPassword

	logger logger.LoggerInterface

	observability observability.TraceLoggerObservability
}

func NewRegisterService(params *RegisterServiceDeps) *registerService {

	return &registerService{
		mencache:      params.Cache,
		user:          params.User,
		role:          params.Role,
		userRole:      params.UserRole,
		hash:          params.Hash,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *registerService) Register(ctx context.Context, request *requests.RegisterRequest) (*models.CreateUserRow, error) {
	const method = "Register"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	existingUser, err := s.user.FindByEmail(ctx, request.Email)
	if err == nil && existingUser != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.CreateUserRow](
			s.logger,
			user_errors.ErrUserEmailAlready,
			method,
			span,
			zap.String("email", request.Email),
		)
	}

	passwordHash, err := s.hash.HashPassword(request.Password)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.CreateUserRow](s.logger, err, method, span)
	}
	request.Password = passwordHash

	const defaultRoleName = "ROLE_ADMIN"
	role, err := s.role.FindByName(ctx, defaultRoleName)
	if err != nil || role == nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.CreateUserRow](s.logger, err, method, span, zap.String("role_name", defaultRoleName))
	}

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.CreateUserRow](s.logger, err, method, span)
	}
	request.VerifiedCode = random
	request.IsVerified = false

	newUser, err := s.user.CreateUser(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.CreateUserRow](s.logger, err, method, span)
	}

	_, err = s.userRole.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: int(newUser.UserID),
		RoleId: int(role.RoleID),
	})
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.CreateUserRow](s.logger, err, method, span, zap.Int("user.id", int(newUser.UserID)))
	}

	s.mencache.SetVerificationCodeCache(ctx, request.Email, random, 15*time.Minute)

	logSuccess("User registered successfully",
		zap.String("email", request.Email),
		zap.String("first_name", request.FirstName),
		zap.String("last_name", request.LastName),
	)

	return newUser, nil
}
