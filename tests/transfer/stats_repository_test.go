package transfer_test

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	repository "github.com/MamangRust/monolith-payment-gateway-transfer/repository/stats"
	statsbycard_repository "github.com/MamangRust/monolith-payment-gateway-transfer/repository/statsbycard"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TransferStatsRepositoryTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	repo       repository.TransferStatsStatusRepository
	repoByCard statsbycard_repository.TransferStatsByCardStatusRepository
	gormDB     *gorm.DB
	testYear   int
	testMonth  int
}

func (s *TransferStatsRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}
	s.gormDB = gormDB
	s.repo = repository.NewTransferStatsStatusRepository(gormDB)
	s.repoByCard = statsbycard_repository.NewTransferStatsByCardStatusRepository(gormDB)
	s.testYear = time.Now().Year()
	s.testMonth = int(time.Now().Month())
}

func (s *TransferStatsRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *TransferStatsRepositoryTestSuite) TestTransferStatusStats() {
	ctx := context.Background()

	// Seed data
	var userID int32
	err := s.gormDB.WithContext(ctx).Raw("INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('Transfer', 'Stats', 'transfer_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&userID).Error
	s.Require().NoError(err)

	cardFrom := "1111222233334444"
	cardTo := "5555666677778888"
	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", userID, cardFrom).Error
	s.Require().NoError(err)
	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", userID, cardTo).Error
	s.Require().NoError(err)

	// Successful transfer
	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO transfers (transfer_from, transfer_to, transfer_amount, transfer_time, status) VALUES (?, ?, ?, ?, 'success')",
		cardFrom, cardTo, 100000, time.Date(s.testYear, time.Month(s.testMonth), 10, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)

	// Failed transfer
	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO transfers (transfer_from, transfer_to, transfer_amount, transfer_time, status) VALUES (?, ?, ?, ?, 'failed')",
		cardFrom, cardTo, 50000, time.Date(s.testYear, time.Month(s.testMonth), 11, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)

	// Global Monthly
	reqMonth := &requests.MonthStatusTransfer{Year: s.testYear, Month: s.testMonth}
	resSuccess, err := s.repo.GetMonthTransferStatusSuccess(ctx, reqMonth)
	s.NoError(err)
	s.NotEmpty(resSuccess)

	resFailed, err := s.repo.GetMonthTransferStatusFailed(ctx, reqMonth)
	s.NoError(err)
	s.NotEmpty(resFailed)

	// Global Yearly
	resYearSuccess, err := s.repo.GetYearlyTransferStatusSuccess(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(resYearSuccess)

	resYearFailed, err := s.repo.GetYearlyTransferStatusFailed(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(resYearFailed)

	// Card Monthly
	reqCard := &requests.MonthStatusTransferCardNumber{
		Year:       s.testYear,
		Month:      s.testMonth,
		CardNumber: cardFrom,
	}
	resCardSuccess, err := s.repoByCard.GetMonthTransferStatusSuccessByCardNumber(ctx, reqCard)
	s.NoError(err)
	s.NotEmpty(resCardSuccess)

	resCardFailed, err := s.repoByCard.GetMonthTransferStatusFailedByCardNumber(ctx, reqCard)
	s.NoError(err)
	s.NotEmpty(resCardFailed)
}

func TestTransferStatsRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransferStatsRepositoryTestSuite))
}
