package saldo_test

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	tests "github.com/MamangRust/monolith-payment-gateway-test"

	"github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-saldo/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type SaldoStatsServiceTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	svc         service.Service
	userID      int32
	cardNumber  string
	testYear    int
}

func (s *SaldoStatsServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}
	repos := repository.NewRepositories(gormDB)

	zapLog := zap.NewNop()
	myLogger := &logger.Logger{Log: zapLog}

	redisOption, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(redisOption)
	cacheStore := cache.NewCacheStore(redisClient, myLogger, &dummyCacheMetrics{})

	s.svc = service.NewService(&service.Deps{
		Repositories: repos,
		Logger:       myLogger,
		Cache:        cacheStore,
	})

	s.testYear = time.Now().Year()

	ctx := context.Background()
	err = s.dbPool.QueryRow(ctx, "INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('SaldoService', 'Stats', 'saldo_svc_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&s.userID)
	s.Require().NoError(err)

	s.cardNumber = "4444555566667777"
	_, err = s.dbPool.Exec(ctx, "INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES ($1, $2, 'debit', '123', 'visa', '2030-01-01')", s.userID, s.cardNumber)
	s.Require().NoError(err)

	// Seed Saldo with balance history
	_, err = s.dbPool.Exec(ctx, "INSERT INTO saldos (card_number, total_balance) VALUES ($1, $2)", s.cardNumber, 1000000)
	s.Require().NoError(err)
}

func (s *SaldoStatsServiceTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	s.ts.Teardown()
}

func (s *SaldoStatsServiceTestSuite) TestSaldoStatsService() {
	ctx := context.Background()

	// Monthly balances
	resMonthly, err := s.svc.FindMonthlySaldoBalances(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(resMonthly)

	// Yearly balances
	resYearly, err := s.svc.FindYearlySaldoBalances(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(resYearly)
}


func TestSaldoStatsServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SaldoStatsServiceTestSuite))
}
