package withdraw_test

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	withdrawstatsrepository "github.com/MamangRust/monolith-payment-gateway-withdraw/repository/stats"
	withdrawstatsbycardrepository "github.com/MamangRust/monolith-payment-gateway-withdraw/repository/statsbycard"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type WithdrawStatsRepositoryTestSuite struct {
	suite.Suite
	ts               *tests.TestSuite
	repo             withdrawstatsrepository.WithdrawStatsStatusRepository
	repoAmount       withdrawstatsrepository.WithdrawStatsAmountRepository
	repoByCard       withdrawstatsbycardrepository.WithdrawStatsByCardStatusRepository
	repoAmountByCard withdrawstatsbycardrepository.WithdrawStatsByCardAmountRepository
	testYear         int
}

func (s *WithdrawStatsRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}
	s.repo = withdrawstatsrepository.NewWithdrawStatsStatusRepository(gormDB)
	s.repoAmount = withdrawstatsrepository.NewWithdrawStatsAmountRepository(gormDB)
	s.repoByCard = withdrawstatsbycardrepository.NewWithdrawStatsByCardStatusRepository(gormDB)
	s.repoAmountByCard = withdrawstatsbycardrepository.NewWithdrawStatsByCardAmountRepository(gormDB)
	s.testYear = time.Now().Year()

	ctx := context.Background()

	var userID int32
	err = gormDB.WithContext(ctx).Raw("INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('WithdrawRepo', 'Stats', 'withdraw_repo_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&userID).Error
	s.Require().NoError(err)

	cardNumber1 := "1111222233334444"
	err = gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", userID, cardNumber1).Error
	s.Require().NoError(err)

	// Jan Success (1000)
	err = gormDB.WithContext(ctx).Exec("INSERT INTO withdraws (card_number, withdraw_amount, withdraw_time, status) VALUES (?, ?, ?, 'success')", cardNumber1, 1000, time.Date(s.testYear, 1, 10, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)
	// Feb Failed (500)
	err = gormDB.WithContext(ctx).Exec("INSERT INTO withdraws (card_number, withdraw_amount, withdraw_time, status) VALUES (?, ?, ?, 'failed')", cardNumber1, 500, time.Date(s.testYear, 2, 15, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)
}

func (s *WithdrawStatsRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *WithdrawStatsRepositoryTestSuite) TestGlobalStats() {
	ctx := context.Background()

	reqMonth := &requests.MonthStatusWithdraw{Year: s.testYear, Month: 1}
	resS, err := s.repo.GetMonthWithdrawStatusSuccess(ctx, reqMonth)
	s.NoError(err)
	s.NotEmpty(resS)

	resF, err := s.repo.GetMonthWithdrawStatusFailed(ctx, reqMonth)
	s.NoError(err)
	s.NotEmpty(resF)

	resA, err := s.repoAmount.GetMonthlyWithdraws(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(resA)

}

func (s *WithdrawStatsRepositoryTestSuite) TestByCardStats() {
	ctx := context.Background()
	cardNumber1 := "1111222233334444"

	reqMonth := &requests.MonthStatusWithdrawCardNumber{Year: s.testYear, Month: 1, CardNumber: cardNumber1}
	reqYearMonth := &requests.YearMonthCardNumber{Year: s.testYear, CardNumber: cardNumber1}

	resS, err := s.repoByCard.GetMonthWithdrawStatusSuccessByCardNumber(ctx, reqMonth)
	s.NoError(err)
	s.NotEmpty(resS)

	resF, err := s.repoByCard.GetMonthWithdrawStatusFailedByCardNumber(ctx, reqMonth)
	s.NoError(err)
	s.NotEmpty(resF)

	resA, err := s.repoAmountByCard.GetMonthlyWithdrawsByCardNumber(ctx, reqYearMonth)
	s.NoError(err)
	s.NotEmpty(resA)

	resY, err := s.repoAmountByCard.GetYearlyWithdrawsByCardNumber(ctx, reqYearMonth)
	s.NoError(err)
	s.NotEmpty(resY)

}

func TestWithdrawStatsRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(WithdrawStatsRepositoryTestSuite))
}
