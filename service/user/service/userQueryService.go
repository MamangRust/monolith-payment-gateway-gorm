package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errorhandler"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	mencache "github.com/MamangRust/monolith-payment-gateway-user/redis"
	"github.com/MamangRust/monolith-payment-gateway-user/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type UserQueryDeps struct {
	Cache         mencache.UserQueryCache
	Repository    repository.UserQueryRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type userQueryService struct {
	cache               mencache.UserQueryCache
	userQueryRepository repository.UserQueryRepository
	logger              logger.LoggerInterface
	observability       observability.TraceLoggerObservability
}

func NewUserQueryService(params *UserQueryDeps) UserQueryService {
	return &userQueryService{
		cache:               params.Cache,
		userQueryRepository: params.Repository,
		logger:              params.Logger,
		observability:       params.Observability,
	}
}

func (s *userQueryService) FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserRow, *int, error) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedUsersCache(ctx, req); found {
		logSuccess("Successfully retrieved all user records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	users, err := s.userQueryRepository.FindAllUsers(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*models.UserRow](
			s.logger,
			user_errors.ErrFailedFindAll,
			method,
			span,
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int
	if len(users) > 0 {
		totalCount = int(users[0].TotalCount)
	}

	s.cache.SetCachedUsersCache(ctx, req, users, &totalCount)
	logSuccess("Successfully fetched user", zap.Int("totalRecords", totalCount), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return users, &totalCount, nil
}

func (s *userQueryService) FindByID(ctx context.Context, id int) (*models.UserByIDRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedUserCache(ctx, id); found {
		logSuccess("Successfully retrieved user record from cache", zap.Int("user.id", id))
		return data, nil
	}

	user, err := s.userQueryRepository.FindById(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*models.UserByIDRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("user_id", id),
		)
	}

	s.cache.SetCachedUserCache(ctx, user)
	logSuccess("Successfully fetched user", zap.Int("user_id", id))
	return user, nil
}

func (s *userQueryService) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserActiveRow, *int, error) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedUserActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active user records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	users, err := s.userQueryRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*models.UserActiveRow](
			s.logger,
			user_errors.ErrFailedFindActive,
			method,
			span,
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int
	if len(users) > 0 {
		totalCount = int(users[0].TotalCount)
	}

	s.cache.SetCachedUserActiveCache(ctx, req, users, &totalCount)
	logSuccess("Successfully fetched active user", zap.Int("totalRecords", totalCount), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return users, &totalCount, nil
}

func (s *userQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserTrashedRow, *int, error) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedUserTrashedCache(ctx, req); found {
		return data, total, nil
	}

	users, err := s.userQueryRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*models.UserTrashedRow](
			s.logger,
			user_errors.ErrFailedFindTrashed,
			method,
			span,
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int
	if len(users) > 0 {
		totalCount = int(users[0].TotalCount)
	}

	s.cache.SetCachedUserTrashedCache(ctx, req, users, &totalCount)
	logSuccess("Successfully fetched trashed user", zap.Int("totalRecords", totalCount), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return users, &totalCount, nil
}

func (s *userQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
