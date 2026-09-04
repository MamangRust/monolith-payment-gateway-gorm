package transaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/MamangRust/monolith-payment-gateway-transaction/repository"
	"github.com/MamangRust/monolith-payment-gateway-transaction/service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	merchant_repo "github.com/MamangRust/monolith-payment-gateway-merchant/repository"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
)

type TransactionStatsServiceTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	svc        service.Service
	userID     int32
	cardNumber string
	merchantID int32
	testYear   int
	testMonth  int
}

func (s *TransactionStatsServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}

	// Real repository implementations for full integration
	realSaldo := &realSaldoRepo{repo: saldo_repo.NewRepositories(gormDB)}
	realCard := &realCardRepo{
		query:   card_repo.NewCardQueryRepository(gormDB),
		command: card_repo.NewCardCommandRepository(gormDB),
	}
	realMerchant := &realMerchantRepo{repo: merchant_repo.NewMerchantQueryRepository(gormDB)}

	repos := repository.NewRepositories(gormDB, realSaldo, realCard, realMerchant)

	zapLog := zap.NewNop()
	myLogger := &logger.Logger{Log: zapLog}

	redisOption, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(redisOption)
	cacheStore := cache.NewCacheStore(redisClient, myLogger, &dummyCacheMetrics{})

	s.svc = service.NewService(&service.Deps{
		Kafka:        nil,
		Repositories: repos,
		Logger:       myLogger,
		Cache:        cacheStore,
	})

	s.testYear = time.Now().Year()
	s.testMonth = int(time.Now().Month())

	ctx := context.Background()
	err = gormDB.WithContext(ctx).Raw("INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('TransactionService', 'Stats', 'transaction_svc_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&s.userID).Error
	s.Require().NoError(err)

	s.cardNumber = "1111222233334444"
	err = gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", s.userID, s.cardNumber).Error
	s.Require().NoError(err)

	// Merchant needs to be active for transactions
	err = gormDB.WithContext(ctx).Raw("INSERT INTO merchants (name, api_key, user_id, status) VALUES ('Test Merchant', 'test_key', ?, 'active') RETURNING merchant_id", s.userID).Scan(&s.merchantID).Error
	s.Require().NoError(err)

	// Seed Transactions
	err = gormDB.WithContext(ctx).Exec("INSERT INTO transactions (card_number, merchant_id, amount, payment_method, transaction_time, status) VALUES (?, ?, ?, 'bank_transfer', ?, 'success')",
		s.cardNumber, s.merchantID, 100000, time.Date(s.testYear, time.Month(s.testMonth), 10, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)

	err = gormDB.WithContext(ctx).Exec("INSERT INTO transactions (card_number, merchant_id, amount, payment_method, transaction_time, status) VALUES (?, ?, ?, 'bank_transfer', ?, 'failed')",
		s.cardNumber, s.merchantID, 50000, time.Date(s.testYear, time.Month(s.testMonth), 11, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)
}

func (s *TransactionStatsServiceTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *TransactionStatsServiceTestSuite) TestTransactionStatsService() {
	ctx := context.Background()

	// Global Monthly Success
	reqMonth := &requests.MonthStatusTransaction{Year: s.testYear, Month: s.testMonth}
	resSuccess, err := s.svc.FindMonthTransactionStatusSuccess(ctx, reqMonth)
	s.NoError(err)
	s.NotEmpty(resSuccess)

	// Global Monthly Failed
	resFailed, err := s.svc.FindMonthTransactionStatusFailed(ctx, reqMonth)
	s.NoError(err)
	s.NotEmpty(resFailed)

	// Yearly Success
	resYearSuccess, err := s.svc.FindYearlyTransactionStatusSuccess(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(resYearSuccess)

}

func TestTransactionStatsServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsServiceTestSuite))
}
