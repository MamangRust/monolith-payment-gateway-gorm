package saldo_test

import (
	"context"
	"testing"
	"time"

	repository "github.com/MamangRust/monolith-payment-gateway-saldo/repository/stats"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SaldoStatsRepositoryTestSuite struct {
	suite.Suite
	ts       *tests.TestSuite
	gormDB   *gorm.DB
	repo     repository.SaldoStatsRepository
	testYear int
}

func (s *SaldoStatsRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}
	s.gormDB = gormDB
	s.repo = repository.NewSaldoStatsRepository(gormDB)
	s.testYear = time.Now().Year()
}

func (s *SaldoStatsRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *SaldoStatsRepositoryTestSuite) TestBalanceStats() {
	ctx := context.Background()

	// Seed data
	var userID int32
	err := s.gormDB.WithContext(ctx).Raw("INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('Saldo', 'Stats', 'saldo_stats_balance@example.com', 'pass', '123', true) RETURNING user_id").Scan(&userID).Error
	s.Require().NoError(err)

	cardNumber := "1111222233334444"
	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", userID, cardNumber).Error
	s.Require().NoError(err)

	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO saldos (card_number, total_balance, created_at, updated_at) VALUES (?, ?, ?, ?)",
		cardNumber, 50000, time.Date(s.testYear, 1, 15, 10, 0, 0, 0, time.UTC), time.Date(s.testYear, 1, 15, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)

	// Monthly Global
	res, err := s.repo.GetMonthlySaldoBalances(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)

	// Yearly Global (5 years)
	yearlyRes, err := s.repo.GetYearlySaldoBalances(ctx, s.testYear)
	s.NoError(err)
	s.Len(yearlyRes, 5)

}

func (s *SaldoStatsRepositoryTestSuite) TestTotalStats() {
	ctx := context.Background()

	// Seed data for current month
	var userID int32
	err := s.gormDB.WithContext(ctx).Raw("INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('Saldo', 'Total', 'saldo_total_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&userID).Error
	s.Require().NoError(err)
	cardNumber := "2222333344445555"
	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", userID, cardNumber).Error
	s.Require().NoError(err)
	now := time.Now()
	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO saldos (card_number, total_balance, created_at, updated_at) VALUES (?, ?, ?, ?)", cardNumber, 50000, now, now).Error
	s.Require().NoError(err)

	// Monthly Total Balance (at least current month)
	req := &requests.MonthTotalSaldoBalance{
		Year:  s.testYear,
		Month: int(time.Now().Month()),
	}
	res, err := s.repo.GetMonthlyTotalSaldoBalance(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)

	// Yearly Total Balance
	yearlyRes, err := s.repo.GetYearTotalSaldoBalance(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(yearlyRes)
}

func TestSaldoStatsRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SaldoStatsRepositoryTestSuite))
}
