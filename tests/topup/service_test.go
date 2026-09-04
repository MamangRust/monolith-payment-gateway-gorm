package topup_test

import (
	"context"
	"fmt"
	"sync"
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
	"github.com/MamangRust/monolith-payment-gateway-topup/repository"
	"github.com/MamangRust/monolith-payment-gateway-topup/service"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"
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
		CardNumber: result.CardNumber, Email: result.Email,
	}, nil
}

func (r *testCardRepo) FindCardByCardNumber(ctx context.Context, card_number string) (*models.CardAllFieldsRow, error) {
	var result models.Card
	err := r.db.WithContext(ctx).Where("card_number = ? AND deleted_at IS NULL", card_number).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &models.CardAllFieldsRow{
		CardID: result.CardID, UserID: result.UserID, CardNumber: result.CardNumber,
	}, nil
}

func (r *testCardRepo) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*models.CardUpdateRow, error) {
	err := r.db.WithContext(ctx).Model(&models.Card{}).Where("id = ?", request.CardID).Updates(map[string]interface{}{
		"card_type":     request.CardType,
		"expire_date":   request.ExpireDate,
		"cvv":           request.CVV,
		"card_provider": request.CardProvider,
	}).Error
	if err != nil {
		return nil, err
	}
	var result models.Card
	err = r.db.WithContext(ctx).Where("id = ?", request.CardID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &models.CardUpdateRow{
		CardID: result.CardID, CardNumber: result.CardNumber, CardType: result.CardType,
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

type TopupServiceTestSuite struct {
	suite.Suite
	ts           *tests.TestSuite
	topupService service.Service
	gormDB       *gorm.DB
	topupID      int
	cardNumber   string
	userID       int
}

func (s *TopupServiceTestSuite) SetupSuite() {
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

	s.topupService = service.NewService(&service.Deps{
		Kafka:        nil,
		Repositories: repos,
		Logger:       log,
		Cache:        cacheStore,
	})

	// Seed User
	userRepo := user_repo.NewUserCommandRepository(gormDB)
	user, err := userRepo.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Topup",
		LastName:  "Tester",
		Email:     fmt.Sprintf("topup.tester.%d@example.com", time.Now().UnixNano()),
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

func (s *TopupServiceTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

// helper
func (s *TopupServiceTestSuite) getSaldoBalance(ctx context.Context, cardNumber string) int32 {
	var saldo models.Saldo
	err := s.gormDB.WithContext(ctx).Raw(`SELECT saldo_id, card_number, total_balance FROM saldos WHERE card_number = ? AND deleted_at IS NULL`, cardNumber).Scan(&saldo).Error
	s.Require().NoError(err)
	return saldo.TotalBalance
}

// helper
func (s *TopupServiceTestSuite) getTopupAmount(ctx context.Context, topupID int) int32 {
	var t struct {
		TopupAmount int32
	}
	err := s.gormDB.WithContext(ctx).Raw(`SELECT topup_amount FROM topups WHERE topup_id = ? AND deleted_at IS NULL`, topupID).Scan(&t).Error
	s.Require().NoError(err)
	return t.TopupAmount
}

func (s *TopupServiceTestSuite) Test1_TopupLifecycle() {
	ctx := context.Background()

	// Create Topup
	createReq := &requests.CreateTopupRequest{
		CardNumber:  s.cardNumber,
		TopupAmount: 100000,
		TopupMethod: "gopay",
	}
	res, err := s.topupService.CreateTopup(ctx, createReq)
	s.NoError(err)
	s.NotNil(res)
	s.topupID = int(res.TopupID)
	s.Equal("success", res.Status)

	// FindById
	found, err := s.topupService.FindById(ctx, s.topupID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.topupID), found.TopupID)

	// Update Topup
	updateReq := &requests.UpdateTopupRequest{
		TopupID:     &s.topupID,
		CardNumber:  s.cardNumber,
		TopupAmount: 200000,
		TopupMethod: "dana",
	}
	updated, err := s.topupService.UpdateTopup(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("success", updated.Status)
}

func (s *TopupServiceTestSuite) Test2_QueryOperations() {
	ctx := context.Background()

	// FindAll
	all, total, err := s.topupService.FindAll(ctx, &requests.FindAllTopups{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(all)
	s.GreaterOrEqual(*total, 1)

	// FindByActive
	active, totalActive, err := s.topupService.FindByActive(ctx, &requests.FindAllTopups{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(active)
	s.GreaterOrEqual(*totalActive, 1)

	// FindByCardNumber
	byCard, _, err := s.topupService.FindAllByCardNumber(ctx, &requests.FindAllTopupsByCardNumber{
		CardNumber: s.cardNumber,
		Page:       1,
		PageSize:   10,
	})
	s.NoError(err)
	s.NotNil(byCard)
	s.GreaterOrEqual(len(byCard), 1)
}

func (s *TopupServiceTestSuite) Test3_TrashAndRestore() {
	ctx := context.Background()
	s.Require().NotZero(s.topupID)

	// Trash
	trashed, err := s.topupService.TrashedTopup(ctx, s.topupID)
	s.NoError(err)
	s.NotNil(trashed)

	// FindByTrashed
	trashedList, totalTrashed, err := s.topupService.FindByTrashed(ctx, &requests.FindAllTopups{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(trashedList)
	s.GreaterOrEqual(*totalTrashed, 1)

	// Restore
	restored, err := s.topupService.RestoreTopup(ctx, s.topupID)
	s.NoError(err)
	s.NotNil(restored)
}

func (s *TopupServiceTestSuite) Test4_BulkOperations() {
	ctx := context.Background()

	// Restore All
	ok, err := s.topupService.RestoreAllTopup(ctx)
	s.NoError(err)
	s.True(ok)

	// Delete All Permanent
	ok, err = s.topupService.DeleteAllTopupPermanent(ctx)
	s.NoError(err)
	s.True(ok)
}

func (s *TopupServiceTestSuite) Test5_TopupIdempotentReplay() {
	ctx := context.Background()
	key := fmt.Sprintf("topup:create:user:%d:op:%d", s.userID, time.Now().UnixNano())

	req := &requests.CreateTopupRequest{
		CardNumber:     s.cardNumber,
		TopupAmount:    150000,
		TopupMethod:    "gopay",
		IdempotencyKey: key,
	}

	first, err := s.topupService.CreateTopup(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(first)
	s.Equal("success", first.Status)

	balanceAfterCreate := s.getSaldoBalance(ctx, s.cardNumber)

	// Replay with the same key: the original record is returned and the
	// balance must not be credited a second time.
	second, err := s.topupService.CreateTopup(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(second)
	s.Equal(first.TopupID, second.TopupID, "replay must return the original record")

	balanceAfterReplay := s.getSaldoBalance(ctx, s.cardNumber)
	s.Equal(balanceAfterCreate, balanceAfterReplay, "replay must not credit the balance twice")
}

func (s *TopupServiceTestSuite) Test6_TopupConcurrentIdempotentReplay() {
	ctx := context.Background()
	key := fmt.Sprintf("topup:create:user:%d:concurrent:%d", s.userID, time.Now().UnixNano())
	amount := 120000

	before := s.getSaldoBalance(ctx, s.cardNumber)

	// Two concurrent creates with the SAME idempotency key: both must succeed
	// (the loser hits the unique index and falls back to returning the winner's
	// record — replay, not 409) and only one record/credit may exist.
	requestsReady := make(chan struct{})
	results := make(chan error, 2)
	ids := make(chan int32, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			<-requestsReady
			res, err := s.topupService.CreateTopup(ctx, &requests.CreateTopupRequest{
				CardNumber:     s.cardNumber,
				TopupAmount:    amount,
				TopupMethod:    "gopay",
				IdempotencyKey: key,
			})
			results <- err
			if err == nil {
				ids <- res.TopupID
			}
		}()
	}
	close(requestsReady)
	wg.Wait()
	close(results)
	close(ids)

	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	s.Equal(2, successes, "concurrent duplicate with the same key must both succeed (replay, not 409)")

	var firstID int32
	first := true
	for id := range ids {
		if first {
			firstID = id
			first = false
		} else {
			s.Equal(firstID, id, "both concurrent requests must return the same record id")
		}
	}

	var rowCount int64
	err := s.gormDB.WithContext(ctx).Model(&models.Topup{}).Where("idempotency_key = ?", key).Count(&rowCount).Error
	s.Require().NoError(err)
	s.Equal(int64(1), rowCount, "exactly one record may exist for the key")

	after := s.getSaldoBalance(ctx, s.cardNumber)
	s.Equal(before+int32(amount), after, "balance must be credited exactly once")
}

func (s *TopupServiceTestSuite) Test7_TopupConcurrentUpdatesApplyAllDeltas() {
	ctx := context.Background()

	// Seed a topup to update concurrently.
	key := fmt.Sprintf("topup:create:user:%d:update:%d", s.userID, time.Now().UnixNano())
	res, err := s.topupService.CreateTopup(ctx, &requests.CreateTopupRequest{
		CardNumber:     s.cardNumber,
		TopupAmount:    100000,
		TopupMethod:    "gopay",
		IdempotencyKey: key,
	})
	s.Require().NoError(err)
	topupID := int(res.TopupID)

	before := s.getSaldoBalance(ctx, s.cardNumber)

	// Two concurrent updates with different target amounts.
	requestsReady := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	for _, newAmount := range []int{150000, 200000} {
		newAmount := newAmount
		go func() {
			defer wg.Done()
			<-requestsReady
			_, err := s.topupService.UpdateTopup(ctx, &requests.UpdateTopupRequest{
				TopupID:     &topupID,
				CardNumber:  s.cardNumber,
				TopupAmount: newAmount,
				TopupMethod: "gopay",
			})
			results <- err
		}()
	}
	close(requestsReady)
	wg.Wait()
	close(results)

	for err := range results {
		s.NoError(err, "concurrent updates must not fail")
	}

	finalAmount := s.getTopupAmount(ctx, topupID)
	after := s.getSaldoBalance(ctx, s.cardNumber)
	expected := before - 100000 + finalAmount
	s.Equal(expected, after, "balance must be consistent with the final topup amount")
}

func (s *TopupServiceTestSuite) TestZZ_NotFound404() {
	ctx := context.Background()
	_, err := s.topupService.FindById(ctx, 999999999)
	s.Require().Error(err)
	s.Contains(err.Error(), "not found")
}

func (s *TopupServiceTestSuite) Test8_FindByIdNotFound() {
	ctx := context.Background()
	_, err := s.topupService.FindById(ctx, 999999999)
	s.Require().Error(err, "a non-existent topup must not resolve")
	s.Contains(err.Error(), "not found", "missing topup must surface as a 404-style not-found error")
}

func TestTopupServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TopupServiceTestSuite))
}
