package transfer_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/MamangRust/monolith-payment-gateway-transfer/repository"
	"github.com/MamangRust/monolith-payment-gateway-transfer/service"
	user_repo_impl "github.com/MamangRust/monolith-payment-gateway-user/repository"
	card_repo_impl "github.com/MamangRust/monolith-payment-gateway-card/repository"
	saldo_repo_impl "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
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

func (r *testCardRepo) FindCardByCardNumber(ctx context.Context, card_number string) (*models.CardAllFieldsRow, error) {
	var result models.Card
	err := r.db.WithContext(ctx).Where("card_number = ? AND deleted_at IS NULL", card_number).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &models.CardAllFieldsRow{
		CardID: result.CardID, UserID: result.UserID, CardNumber: result.CardNumber,
		CardType: result.CardType, ExpireDate: result.ExpireDate, Cvv: result.Cvv,
		CardProvider: result.CardProvider, Status: result.Status,
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

type TransferServiceTestSuite struct {
	suite.Suite
	ts              *tests.TestSuite
	transferService service.Service
	dbPool          *pgxpool.Pool
	gormDB          *gorm.DB
	transferID      int
	senderCard      string
	receiverCard    string
	userID          int
}

func (s *TransferServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

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
	repos := repository.NewRepositories(gormDB, saldoRepo, cardRepo)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)

	s.transferService = service.NewService(&service.Deps{
		Kafka:        nil,
		Repositories: repos,
		Logger:       log,
		Cache:        cacheStore,
	})

	// Seed User
	userRepo := user_repo_impl.NewUserCommandRepository(gormDB)
	user, err := userRepo.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Transfer",
		LastName:  "Tester",
		Email:     fmt.Sprintf("transfer.tester.%d@example.com", time.Now().UnixNano()),
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// Seed Sender Card & Saldo
	cardCmdRepo := card_repo_impl.NewCardCommandRepository(gormDB)
	sender, err := cardCmdRepo.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.senderCard = sender.CardNumber

	saldoCmdRepo := saldo_repo_impl.NewSaldoCommandRepository(gormDB)
	_, err = saldoCmdRepo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   s.senderCard,
		TotalBalance: 1000000,
	})
	s.Require().NoError(err)

	// Seed Receiver Card & Saldo
	receiver, err := cardCmdRepo.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "456",
		CardProvider: "mastercard",
	})
	s.Require().NoError(err)
	s.receiverCard = receiver.CardNumber

	_, err = saldoCmdRepo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   s.receiverCard,
		TotalBalance: 0,
	})
	s.Require().NoError(err)
}

func (s *TransferServiceTestSuite) TearDownSuite() {
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *TransferServiceTestSuite) getSaldoBalance(ctx context.Context, cardNumber string) int32 {
	var saldo models.Saldo
	err := s.gormDB.WithContext(ctx).Raw(`SELECT saldo_id, card_number, total_balance FROM saldos WHERE card_number = ? AND deleted_at IS NULL`, cardNumber).Scan(&saldo).Error
	s.Require().NoError(err)
	return saldo.TotalBalance
}

func (s *TransferServiceTestSuite) Test1_TransferLifecycle() {
	ctx := context.Background()

	// Create Transfer
	createReq := &requests.CreateTransferRequest{
		TransferFrom:   s.senderCard,
		TransferTo:     s.receiverCard,
		TransferAmount: 100000,
	}
	res, err := s.transferService.CreateTransaction(ctx, createReq)
	s.NoError(err)
	s.NotNil(res)
	s.transferID = int(res.TransferID)
	s.Equal("success", res.Status)

	// FindById
	found, err := s.transferService.FindById(ctx, s.transferID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.transferID), found.TransferID)

	// Update Transfer
	updateReq := &requests.UpdateTransferRequest{
		TransferID:     &s.transferID,
		TransferFrom:   s.senderCard,
		TransferTo:     s.receiverCard,
		TransferAmount: 200000,
	}
	updated, err := s.transferService.UpdateTransaction(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(200000), updated.TransferAmount)
}

func (s *TransferServiceTestSuite) Test2_QueryOperations() {
	ctx := context.Background()

	// FindAll
	all, total, err := s.transferService.FindAll(ctx, &requests.FindAllTransfers{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(all)
	s.GreaterOrEqual(*total, 1)

	// FindByActive
	active, totalActive, err := s.transferService.FindByActive(ctx, &requests.FindAllTransfers{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(active)
	s.GreaterOrEqual(*totalActive, 1)

	// FindByTransferFrom
	from, err := s.transferService.FindTransferByTransferFrom(ctx, s.senderCard)
	s.NoError(err)
	s.NotNil(from)
	s.GreaterOrEqual(len(from), 1)

	// FindByTransferTo
	to, err := s.transferService.FindTransferByTransferTo(ctx, s.receiverCard)
	s.NoError(err)
	s.NotNil(to)
	s.GreaterOrEqual(len(to), 1)
}

func (s *TransferServiceTestSuite) Test3_TrashAndRestore() {
	ctx := context.Background()
	s.Require().NotZero(s.transferID)

	// Trash
	trashed, err := s.transferService.TrashedTransfer(ctx, s.transferID)
	s.NoError(err)
	s.NotNil(trashed)

	// FindByTrashed
	trashedList, totalTrashed, err := s.transferService.FindByTrashed(ctx, &requests.FindAllTransfers{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(trashedList)
	s.GreaterOrEqual(*totalTrashed, 1)

	// Restore
	restored, err := s.transferService.RestoreTransfer(ctx, s.transferID)
	s.NoError(err)
	s.NotNil(restored)
}

func (s *TransferServiceTestSuite) Test4_BulkOperations() {
	ctx := context.Background()

	// Restore All
	ok, err := s.transferService.RestoreAllTransfer(ctx)
	s.NoError(err)
	s.True(ok)

	// Delete All Permanent
	ok, err = s.transferService.DeleteAllTransferPermanent(ctx)
	s.NoError(err)
	s.True(ok)
}

func (s *TransferServiceTestSuite) Test5_TransferIdempotentReplay() {
	ctx := context.Background()
	key := fmt.Sprintf("transfer:create:user:%d:op:%d", s.userID, time.Now().UnixNano())
	amount := 50000

	senderBefore := s.getSaldoBalance(ctx, s.senderCard)
	receiverBefore := s.getSaldoBalance(ctx, s.receiverCard)

	req := &requests.CreateTransferRequest{
		TransferFrom:   s.senderCard,
		TransferTo:     s.receiverCard,
		TransferAmount: amount,
		IdempotencyKey: key,
	}

	// A replay with the same key must return the original settlement and must
	// NOT move money a second time.
	first, err := s.transferService.CreateTransaction(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(first)

	second, err := s.transferService.CreateTransaction(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(second)
	s.Equal(first.TransferID, second.TransferID, "replay must return the original transfer")

	senderAfter := s.getSaldoBalance(ctx, s.senderCard)
	receiverAfter := s.getSaldoBalance(ctx, s.receiverCard)
	s.Equal(senderBefore-int32(amount), senderAfter, "sender must be debited exactly once")
	s.Equal(receiverBefore+int32(amount), receiverAfter, "receiver must be credited exactly once")
}

func (s *TransferServiceTestSuite) TestZZ_NotFound404() {
	ctx := context.Background()
	_, err := s.transferService.FindById(ctx, 999999999)
	s.Require().Error(err)
	s.Contains(err.Error(), "not found")
}

func (s *TransferServiceTestSuite) Test6_FindByIdNotFound() {
	ctx := context.Background()
	_, err := s.transferService.FindById(ctx, 999999999)
	s.Require().Error(err, "a non-existent transfer must not resolve")
	s.Contains(err.Error(), "not found", "missing transfer must surface as a 404-style not-found error")
}

func TestTransferServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransferServiceTestSuite))
}
