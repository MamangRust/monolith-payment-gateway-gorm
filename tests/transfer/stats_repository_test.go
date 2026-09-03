package transfer_test

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	repository "github.com/MamangRust/monolith-payment-gateway-transfer/repository/stats"
	statsbycard_repository "github.com/MamangRust/monolith-payment-gateway-transfer/repository/statsbycard"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TransferStatsRepositoryTestSuite struct {
	suite.Suite
	ts             *tests.TestSuite
	dbPool         *pgxpool.Pool
	repo           repository.TransferStatsStatusRepository
	repoByCard     statsbycard_repository.TransferStatsByCardStatusRepository
	testYear       int
	testMonth      int
}

func (s *TransferStatsRepositoryTestSuite) SetupSuite() {
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
	s.repo = repository.NewTransferStatsStatusRepository(gormDB)
	s.repoByCard = statsbycard_repository.NewTransferStatsByCardStatusRepository(gormDB)
	s.testYear = time.Now().Year()
	s.testMonth = int(time.Now().Month())
}

func (s *TransferStatsRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	s.ts.Teardown()
}

func (s *TransferStatsRepositoryTestSuite) TestTransferStatusStats() {
	ctx := context.Background()

	// Seed data
	var userID int32
	err := s.dbPool.QueryRow(ctx, "INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('Transfer', 'Stats', 'transfer_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&userID)
	s.Require().NoError(err)

	cardFrom := "1111222233334444"
	cardTo := "5555666677778888"
	_, err = s.dbPool.Exec(ctx, "INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES ($1, $2, 'debit', '123', 'visa', '2030-01-01')", userID, cardFrom)
	s.Require().NoError(err)
	_, err = s.dbPool.Exec(ctx, "INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES ($1, $2, 'debit', '123', 'visa', '2030-01-01')", userID, cardTo)
	s.Require().NoError(err)

	// Successful transfer
	_, err = s.dbPool.Exec(ctx, "INSERT INTO transfers (transfer_from, transfer_to, transfer_amount, transfer_time, status) VALUES ($1, $2, $3, $4, 'success')", 
		cardFrom, cardTo, 100000, time.Date(s.testYear, time.Month(s.testMonth), 10, 10, 0, 0, 0, time.UTC))
	s.Require().NoError(err)

	// Failed transfer
	_, err = s.dbPool.Exec(ctx, "INSERT INTO transfers (transfer_from, transfer_to, transfer_amount, transfer_time, status) VALUES ($1, $2, $3, $4, 'failed')", 
		cardFrom, cardTo, 50000, time.Date(s.testYear, time.Month(s.testMonth), 11, 10, 0, 0, 0, time.UTC))
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
