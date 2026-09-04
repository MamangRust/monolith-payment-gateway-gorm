package card_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	repositorystats "github.com/MamangRust/monolith-payment-gateway-card/repository/stats"
	repositorystatsbycard "github.com/MamangRust/monolith-payment-gateway-card/repository/statsbycard"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CardStatsRepositoryTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	gormDB      *gorm.DB
	repo        repositorystats.CardStatsRepository
	repoByCard  repositorystatsbycard.CardStatsByCardRepository
	cardNumber1 string
	cardNumber2 string
	testYear    int
}

func (s *CardStatsRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}
	s.gormDB = gormDB
	s.repo = repositorystats.NewCardStatsRepository(gormDB)
	s.repoByCard = repositorystatsbycard.NewCardStatsByCardRepository(gormDB)
	s.testYear = time.Now().Year()

	// Seed Users and Cards
	ctx := context.Background()
	s.cardNumber1 = "1111222233334444"
	s.cardNumber2 = "5555666677778888"

	// Create a user first
	var userID int32
	err = s.gormDB.WithContext(ctx).Raw("INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('Stats', 'User', 'stats@example.com', 'pass', '123456', true) RETURNING user_id").Scan(&userID).Error
	s.Require().NoError(err)

	// Create Cards
	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", userID, s.cardNumber1).Error
	s.Require().NoError(err)
	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'credit', '456', 'mastercard', '2030-01-01')", userID, s.cardNumber2).Error
	s.Require().NoError(err)

	s.seedHistoricalData()
}

func (s *CardStatsRepositoryTestSuite) seedHistoricalData() {
	// Seed Saldos (Monthly Balance)
	// Jan: 1000 (Card 1)
	// Feb: 2000 (Card 1), 3000 (Card 2)
	s.insertSaldo(s.cardNumber1, 1000, time.Date(s.testYear, 1, 15, 10, 0, 0, 0, time.UTC))
	s.insertSaldo(s.cardNumber1, 2000, time.Date(s.testYear, 2, 15, 10, 0, 0, 0, time.UTC))
	s.insertSaldo(s.cardNumber2, 3000, time.Date(s.testYear, 2, 20, 10, 0, 0, 0, time.UTC))

	// Seed Topups
	// Jan: 500 (Card 1)
	// Feb: 1000 (Card 1), 1500 (Card 2)
	s.insertTopup(s.cardNumber1, 500, time.Date(s.testYear, 1, 10, 10, 0, 0, 0, time.UTC))
	s.insertTopup(s.cardNumber1, 1000, time.Date(s.testYear, 2, 10, 10, 0, 0, 0, time.UTC))
	s.insertTopup(s.cardNumber2, 1500, time.Date(s.testYear, 2, 12, 10, 0, 0, 0, time.UTC))

	// Seed Withdraws
	// Jan: 100 (Card 1)
	// Feb: 200 (Card 1), 300 (Card 2)
	s.insertWithdraw(s.cardNumber1, 100, time.Date(s.testYear, 1, 20, 10, 0, 0, 0, time.UTC))
	s.insertWithdraw(s.cardNumber1, 200, time.Date(s.testYear, 2, 20, 10, 0, 0, 0, time.UTC))
	s.insertWithdraw(s.cardNumber2, 300, time.Date(s.testYear, 2, 25, 10, 0, 0, 0, time.UTC))

	// Seed Transfers
	// Jan: 50 from Card 1 to Card 2
	// Feb: 100 from Card 1 to Card 2
	s.insertTransfer(s.cardNumber1, s.cardNumber2, 50, time.Date(s.testYear, 1, 25, 10, 0, 0, 0, time.UTC))
	s.insertTransfer(s.cardNumber1, s.cardNumber2, 100, time.Date(s.testYear, 2, 28, 10, 0, 0, 0, time.UTC))
}

func (s *CardStatsRepositoryTestSuite) insertSaldo(cardNumber string, amount int, t time.Time) {
	// Saldo keeps a single active snapshot per card (idx_saldos_card_number_active);
	// superseded snapshots are soft-deleted so the stats history still counts them.
	err := s.gormDB.WithContext(context.Background()).Exec(
		"UPDATE saldos SET deleted_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL",
		cardNumber).Error
	s.Require().NoError(err)
	err = s.gormDB.WithContext(context.Background()).Exec(
		"INSERT INTO saldos (card_number, total_balance, created_at, updated_at) VALUES (?, ?, ?, ?)",
		cardNumber, amount, t, t).Error
	s.Require().NoError(err)
}

func (s *CardStatsRepositoryTestSuite) insertTopup(cardNumber string, amount int, t time.Time) {
	err := s.gormDB.WithContext(context.Background()).Exec(
		"INSERT INTO topups (card_number, topup_amount, topup_time, topup_method, status) VALUES (?, ?, ?, 'bank_transfer', 'success')",
		cardNumber, amount, t).Error
	s.Require().NoError(err)
}

func (s *CardStatsRepositoryTestSuite) insertWithdraw(cardNumber string, amount int, t time.Time) {
	err := s.gormDB.WithContext(context.Background()).Exec(
		"INSERT INTO withdraws (card_number, withdraw_amount, withdraw_time, status) VALUES (?, ?, ?, 'success')",
		cardNumber, amount, t).Error
	s.Require().NoError(err)
}

func (s *CardStatsRepositoryTestSuite) insertTransfer(from, to string, amount int, t time.Time) {
	err := s.gormDB.WithContext(context.Background()).Exec(
		"INSERT INTO transfers (transfer_from, transfer_to, transfer_amount, transfer_time, status) VALUES (?, ?, ?, ?, 'success')",
		from, to, amount, t).Error
	s.Require().NoError(err)
}

func (s *CardStatsRepositoryTestSuite) TearDownSuite() {
	if s.ts != nil {
		s.ts.Teardown()
	}
}

// --- Balance Tests ---

func (s *CardStatsRepositoryTestSuite) TestBalanceStats() {
	ctx := context.Background()

	// Global Monthly
	res, err := s.repo.GetMonthlyBalance(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)
	s.Equal(int32(1000), res[0].TotalBalance) // Jan
	s.Equal(int32(5000), res[1].TotalBalance) // Feb (2000 + 3000)

	// By Card Monthly
	req1 := &requests.MonthYearCardNumberCard{CardNumber: s.cardNumber1, Year: s.testYear}
	res1, err := s.repoByCard.GetMonthlyBalancesByCardNumber(ctx, req1)
	s.NoError(err)
	s.Equal(int32(1000), res1[0].TotalBalance) // Jan Card 1
	s.Equal(int32(2000), res1[1].TotalBalance) // Feb Card 1

	// Yearly
	yRes, err := s.repo.GetYearlyBalance(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(yRes)
	found := false
	yearStr := strconv.Itoa(s.testYear)
	for _, r := range yRes {
		// Normalize the numeric value for comparison
		val := fmt.Sprintf("%v", r.Year)
		if fmt.Sprintf("%v", val) == yearStr {
			s.Equal(int32(6000), int32(r.TotalBalance))
			found = true
		}
	}
	s.True(found)
}

// --- Topup Tests ---

func (s *CardStatsRepositoryTestSuite) TestTopupStats() {
	ctx := context.Background()

	// Global Monthly
	res, err := s.repo.GetMonthlyTopupAmount(ctx, s.testYear)
	s.NoError(err)
	s.Equal(int32(500), res[0].TotalTopupAmount)  // Jan
	s.Equal(int32(2500), res[1].TotalTopupAmount) // Feb (1000 + 1500)

	// By Card Monthly
	req1 := &requests.MonthYearCardNumberCard{CardNumber: s.cardNumber1, Year: s.testYear}
	res1, err := s.repoByCard.GetMonthlyTopupAmountByCardNumber(ctx, req1)
	s.NoError(err)
	s.Equal(int32(500), res1[0].TotalTopupAmount)  // Jan Card 1
	s.Equal(int32(1000), res1[1].TotalTopupAmount) // Feb Card 1
}

// --- Withdraw Tests ---

func (s *CardStatsRepositoryTestSuite) TestWithdrawStats() {
	ctx := context.Background()

	// Global Monthly
	res, err := s.repo.GetMonthlyWithdrawAmount(ctx, s.testYear)
	s.NoError(err)
	s.Equal(int32(100), res[0].TotalWithdrawAmount) // Jan
	s.Equal(int32(500), res[1].TotalWithdrawAmount) // Feb (200 + 300)

	// By Card Monthly
	req1 := &requests.MonthYearCardNumberCard{CardNumber: s.cardNumber1, Year: s.testYear}
	res1, err := s.repoByCard.GetMonthlyWithdrawAmountByCardNumber(ctx, req1)
	s.NoError(err)
	s.Equal(int32(100), res1[0].TotalWithdrawAmount) // Jan Card 1
	s.Equal(int32(200), res1[1].TotalWithdrawAmount) // Feb Card 1
}

// --- Transfer Tests ---

func (s *CardStatsRepositoryTestSuite) TestTransferStats() {
	ctx := context.Background()

	// Global Monthly Sender
	resS, err := s.repo.GetMonthlyTransferAmountSender(ctx, s.testYear)
	s.NoError(err)
	s.Equal(int32(50), resS[0].TotalSentAmount)  // Jan
	s.Equal(int32(100), resS[1].TotalSentAmount) // Feb

	// By Card Monthly Sender
	req1 := &requests.MonthYearCardNumberCard{CardNumber: s.cardNumber1, Year: s.testYear}
	res1S, err := s.repoByCard.GetMonthlyTransferAmountBySender(ctx, req1)
	s.NoError(err)
	s.Equal(int32(50), res1S[0].TotalSentAmount)  // Jan Card 1
	s.Equal(int32(100), res1S[1].TotalSentAmount) // Feb Card 1

	// By Card Monthly Receiver (Card 2)
	req2 := &requests.MonthYearCardNumberCard{CardNumber: s.cardNumber2, Year: s.testYear}
	res2R, err := s.repoByCard.GetMonthlyTransferAmountByReceiver(ctx, req2)
	s.NoError(err)
	s.Equal(int32(50), res2R[0].TotalReceivedAmount)  // Jan Card 2
	s.Equal(int32(100), res2R[1].TotalReceivedAmount) // Feb Card 2
}

func TestCardStatsRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CardStatsRepositoryTestSuite))
}
