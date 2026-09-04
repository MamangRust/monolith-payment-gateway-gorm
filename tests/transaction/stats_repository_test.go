package transaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	repository "github.com/MamangRust/monolith-payment-gateway-transaction/repository/stats"
	statsbycard_repository "github.com/MamangRust/monolith-payment-gateway-transaction/repository/statsbycard"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TransactionStatsRepositoryTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	repo       repository.TransactionStatsStatusRepository
	repoByCard statsbycard_repository.TransactionStatsByCardStatusRepository
	gormDB     *gorm.DB
	testYear   int
	testMonth  int
}

func (s *TransactionStatsRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}
	s.gormDB = gormDB
	s.repo = repository.NewTransactionStatsStatusRepository(gormDB)
	s.repoByCard = statsbycard_repository.NewTransactionStatsByCardStatusRepository(gormDB)
	s.testYear = time.Now().Year()
	s.testMonth = int(time.Now().Month())
}

func (s *TransactionStatsRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *TransactionStatsRepositoryTestSuite) TestTransactionStatusStats() {
	ctx := context.Background()

	// Seed data
	var userID int32
	err := s.gormDB.WithContext(ctx).Raw("INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('Transaction', 'Stats', 'transaction_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&userID).Error
	s.Require().NoError(err)

	cardNumber := "1111222233334444"
	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", userID, cardNumber).Error
	s.Require().NoError(err)

	var merchantID int32
	err = s.gormDB.WithContext(ctx).Raw("INSERT INTO merchants (name, api_key, user_id, status) VALUES ('Merchant Stats', 'test_api_key', ?, 'active') RETURNING merchant_id", userID).Scan(&merchantID).Error
	s.Require().NoError(err)

	// Successful transaction
	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO transactions (card_number, amount, payment_method, merchant_id, transaction_time, status) VALUES (?, ?, 'credit_card', ?, ?, 'success')",
		cardNumber, 100000, merchantID, time.Date(s.testYear, time.Month(s.testMonth), 10, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)

	// Failed transaction
	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO transactions (card_number, amount, payment_method, merchant_id, transaction_time, status) VALUES (?, ?, 'credit_card', ?, ?, 'failed')",
		cardNumber, 50000, merchantID, time.Date(s.testYear, time.Month(s.testMonth), 11, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)

	// Global Monthly
	reqMonth := &requests.MonthStatusTransaction{Year: s.testYear, Month: s.testMonth}
	resSuccess, err := s.repo.GetMonthTransactionStatusSuccess(ctx, reqMonth)
	s.NoError(err)
	s.NotEmpty(resSuccess)

	resFailed, err := s.repo.GetMonthTransactionStatusFailed(ctx, reqMonth)
	s.NoError(err)
	s.NotEmpty(resFailed)

	// Global Yearly
	resYearSuccess, err := s.repo.GetYearlyTransactionStatusSuccess(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(resYearSuccess)

	resYearFailed, err := s.repo.GetYearlyTransactionStatusFailed(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(resYearFailed)

	// Card Monthly
	reqCard := &requests.MonthStatusTransactionCardNumber{
		Year:       s.testYear,
		Month:      s.testMonth,
		CardNumber: cardNumber,
	}
	resCardSuccess, err := s.repoByCard.GetMonthTransactionStatusSuccessByCardNumber(ctx, reqCard)
	s.NoError(err)
	s.NotEmpty(resCardSuccess)

	resCardFailed, err := s.repoByCard.GetMonthTransactionStatusFailedByCardNumber(ctx, reqCard)
	s.NoError(err)
	s.NotEmpty(resCardFailed)

}

func TestTransactionStatsRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsRepositoryTestSuite))
}
