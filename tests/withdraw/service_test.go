package withdraw_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	card_repo_impl "github.com/MamangRust/monolith-payment-gateway-card/repository"
	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	saldo_repo_impl "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	user_repo_impl "github.com/MamangRust/monolith-payment-gateway-user/repository"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/repository"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type testCardRepo struct {
	db *gorm.DB
}

func (r *testCardRepo) FindUserCardByCardNumber(ctx context.Context, card_number string) (*models.CardByEmailRow, error) {
	type cardWithEmail struct {
		models.Card
		Email string
	}
	var result cardWithEmail
	err := r.db.WithContext(ctx).Raw(`
		SELECT c.*, u.email FROM cards c
		JOIN users u ON c.user_id = u.user_id
		WHERE c.card_number = ? AND c.deleted_at IS NULL
	`, card_number).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &models.CardByEmailRow{
		CardID: result.CardID, UserID: result.UserID, CardNumber: result.CardNumber,
		CardType: result.CardType, ExpireDate: result.ExpireDate, Cvv: result.Cvv,
		CardProvider: result.CardProvider, Email: result.Email,
	}, nil
}

type testSaldoRepo struct {
	db *gorm.DB
}

func (r *testSaldoRepo) FindByCardNumber(ctx context.Context, card_number string) (*models.Saldo, error) {
	var result models.Saldo
	err := r.db.WithContext(ctx).Raw(`SELECT saldo_id, card_number, total_balance FROM saldos WHERE card_number = ? AND deleted_at IS NULL`, card_number).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.SaldoID == 0 {
		return nil, fmt.Errorf("saldo not found for card %s", card_number)
	}
	return &result, nil
}

func (r *testSaldoRepo) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*models.UpdateSaldoBalanceRow, error) {
	var result models.Saldo
	err := r.db.WithContext(ctx).Raw(`UPDATE saldos SET total_balance = ?, updated_at = NOW() WHERE card_number = ? AND deleted_at IS NULL RETURNING saldo_id, card_number, total_balance`, request.TotalBalance, request.CardNumber).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &models.UpdateSaldoBalanceRow{
		SaldoID: result.SaldoID, CardNumber: result.CardNumber, TotalBalance: result.TotalBalance,
	}, nil
}

func (r *testSaldoRepo) UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*models.UpdateSaldoWithdrawRow, error) {
	var withdrawTime interface{}
	if request.WithdrawTime != nil {
		withdrawTime = request.WithdrawTime
	}

	var withdrawAmount interface{}
	if request.WithdrawAmount != nil {
		withdrawAmount = *request.WithdrawAmount
	}

	var result models.Saldo
	err := r.db.WithContext(ctx).Raw(`UPDATE saldos SET total_balance = total_balance - ?, withdraw_amount = ?, withdraw_time = ?, updated_at = NOW() WHERE card_number = ? AND deleted_at IS NULL RETURNING saldo_id, card_number, total_balance`, withdrawAmount, withdrawTime, request.CardNumber).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &models.UpdateSaldoWithdrawRow{
		SaldoID: result.SaldoID, CardNumber: result.CardNumber, TotalBalance: result.TotalBalance,
	}, nil
}

type WithdrawServiceTestSuite struct {
	suite.Suite
	ts              *tests.TestSuite
	withdrawService service.Service
	gormDB          *gorm.DB
	withdrawID      int
	cardNumber      string
	userID          int
}

func (s *WithdrawServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(opts)

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}
	s.gormDB = gormDB

	cardRepo := &testCardRepo{db: gormDB}
	saldoRepo := &testSaldoRepo{db: gormDB}
	repos := repository.NewRepositories(gormDB, cardRepo, saldoRepo)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)

	s.withdrawService = service.NewService(&service.Deps{
		Kafka:        nil,
		Repositories: repos,
		Logger:       log,
		Cache:        cacheStore,
	})

	// Seed User
	userRepo := user_repo_impl.NewUserCommandRepository(gormDB)
	user, err := userRepo.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Withdraw",
		LastName:  "Tester",
		Email:     fmt.Sprintf("withdraw.tester.%d@example.com", time.Now().UnixNano()),
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// Seed Card
	cardCmdRepo := card_repo_impl.NewCardCommandRepository(gormDB)
	card, err := cardCmdRepo.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.cardNumber = card.CardNumber

	// Seed Saldo
	saldoCmdRepo := saldo_repo_impl.NewSaldoCommandRepository(gormDB)
	_, err = saldoCmdRepo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   s.cardNumber,
		TotalBalance: 1000000,
	})
	s.Require().NoError(err)
}

func (s *WithdrawServiceTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *WithdrawServiceTestSuite) getSaldoBalance(ctx context.Context, cardNumber string) int32 {
	var saldo models.Saldo
	err := s.gormDB.WithContext(ctx).Raw(`SELECT saldo_id, card_number, total_balance FROM saldos WHERE card_number = ? AND deleted_at IS NULL`, cardNumber).Scan(&saldo).Error
	s.Require().NoError(err)
	return saldo.TotalBalance
}

func (s *WithdrawServiceTestSuite) Test1_WithdrawLifecycle() {
	ctx := context.Background()

	// Create Withdraw
	createReq := &requests.CreateWithdrawRequest{
		CardNumber:     s.cardNumber,
		WithdrawAmount: 100000,
		WithdrawTime:   time.Now(),
	}
	res, err := s.withdrawService.Create(ctx, createReq)
	s.NoError(err)
	s.NotNil(res)
	s.withdrawID = int(res.WithdrawID)
	s.Equal("success", res.Status)

	// FindById
	found, err := s.withdrawService.FindById(ctx, s.withdrawID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.withdrawID), found.WithdrawID)

	// Update Withdraw
	updateReq := &requests.UpdateWithdrawRequest{
		WithdrawID:     &s.withdrawID,
		CardNumber:     s.cardNumber,
		WithdrawAmount: 200000,
		WithdrawTime:   time.Now(),
	}
	updated, err := s.withdrawService.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(200000), updated.WithdrawAmount)
}

func (s *WithdrawServiceTestSuite) Test2_QueryOperations() {
	ctx := context.Background()

	// FindAll
	all, total, err := s.withdrawService.FindAll(ctx, &requests.FindAllWithdraws{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(all)
	s.GreaterOrEqual(*total, 1)

	// FindByActive
	active, totalActive, err := s.withdrawService.FindByActive(ctx, &requests.FindAllWithdraws{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(active)
	s.GreaterOrEqual(*totalActive, 1)

	// FindByCardNumber
	byCard, totalCard, err := s.withdrawService.FindAllByCardNumber(ctx, &requests.FindAllWithdrawCardNumber{
		CardNumber: s.cardNumber,
		Page:       1,
		PageSize:   10,
	})
	s.NoError(err)
	s.NotNil(byCard)
	s.GreaterOrEqual(*totalCard, 1)
}

func (s *WithdrawServiceTestSuite) Test3_TrashAndRestore() {
	ctx := context.Background()
	s.Require().NotZero(s.withdrawID)

	// Trash
	trashed, err := s.withdrawService.TrashedWithdraw(ctx, s.withdrawID)
	s.NoError(err)
	s.NotNil(trashed)

	// FindByTrashed
	trashedList, totalTrashed, err := s.withdrawService.FindByTrashed(ctx, &requests.FindAllWithdraws{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(trashedList)
	s.GreaterOrEqual(*totalTrashed, 1)

	// Restore
	restored, err := s.withdrawService.RestoreWithdraw(ctx, s.withdrawID)
	s.NoError(err)
	s.NotNil(restored)
}

func (s *WithdrawServiceTestSuite) Test4_BulkOperations() {
	ctx := context.Background()

	// Restore All
	ok, err := s.withdrawService.RestoreAllWithdraw(ctx)
	s.NoError(err)
	s.True(ok)

	// Delete All Permanent
	ok, err = s.withdrawService.DeleteAllWithdrawPermanent(ctx)
	s.NoError(err)
	s.True(ok)
}

func (s *WithdrawServiceTestSuite) Test5_WithdrawIdempotentReplay() {
	ctx := context.Background()
	key := fmt.Sprintf("withdraw:create:user:%d:op:%d", s.userID, time.Now().UnixNano())
	amount := 50000

	before := s.getSaldoBalance(ctx, s.cardNumber)

	req := &requests.CreateWithdrawRequest{
		CardNumber:     s.cardNumber,
		WithdrawAmount: amount,
		WithdrawTime:   time.Now(),
		IdempotencyKey: key,
	}

	// A replay with the same key must return the original withdrawal and must
	// NOT debit the balance a second time.
	first, err := s.withdrawService.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(first)

	second, err := s.withdrawService.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(second)
	s.Equal(first.WithdrawID, second.WithdrawID, "replay must return the original withdrawal")

	after := s.getSaldoBalance(ctx, s.cardNumber)
	s.Equal(before-int32(amount), after, "balance must be debited exactly once")
}

func (s *WithdrawServiceTestSuite) TestZZ_NotFound404() {
	ctx := context.Background()
	_, err := s.withdrawService.FindById(ctx, 999999999)
	s.Require().Error(err)
	s.Contains(err.Error(), "not found")
}

func (s *WithdrawServiceTestSuite) Test6_FindByIdNotFound() {
	ctx := context.Background()
	_, err := s.withdrawService.FindById(ctx, 999999999)
	s.Require().Error(err, "a non-existent withdraw must not resolve")
	s.Contains(err.Error(), "not found", "missing withdraw must surface as a 404-style not-found error")
}

func TestWithdrawServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(WithdrawServiceTestSuite))
}
