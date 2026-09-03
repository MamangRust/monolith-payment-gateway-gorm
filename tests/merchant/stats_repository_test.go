package merchant_test

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	stats_repo "github.com/MamangRust/monolith-payment-gateway-merchant/repository/stats"
	apikey_repo "github.com/MamangRust/monolith-payment-gateway-merchant/repository/statsbyapikey"
	merchant_repo "github.com/MamangRust/monolith-payment-gateway-merchant/repository/statsbymerchant"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type MerchantStatsRepositoryTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	repo        stats_repo.MerchantStatsRepository
	apikeyRepo  apikey_repo.MerchantStatsByApiKeyRepository
	merchantRepo merchant_repo.MerchantStatsByMerchantRepository
	merchantID1 int32
	merchantID2 int32
	apiKey1     string
	apiKey2     string
	testYear    int
}

func (s *MerchantStatsRepositoryTestSuite) SetupSuite() {
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
	s.repo = stats_repo.NewMerchantStatsRepository(gormDB)
	s.apikeyRepo = apikey_repo.NewMerchantStatsByApiKeyRepository(gormDB)
	s.merchantRepo = merchant_repo.NewMerchantStatsByMerchantRepository(gormDB)

	s.testYear = time.Now().Year()

	// Seed Data
	ctx := context.Background()
	var userID int32
	err = s.dbPool.QueryRow(ctx, "INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('Merchant', 'Stats', 'merchant_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&userID)
	s.Require().NoError(err)

	s.apiKey1 = "merchant-key-1"
	s.apiKey2 = "merchant-key-2"

	err = s.dbPool.QueryRow(ctx, "INSERT INTO merchants (name, api_key, user_id, status) VALUES ('Merchant 1', $1, $2, 'active') RETURNING merchant_id", s.apiKey1, userID).Scan(&s.merchantID1)
	s.Require().NoError(err)

	err = s.dbPool.QueryRow(ctx, "INSERT INTO merchants (name, api_key, user_id, status) VALUES ('Merchant 2', $1, $2, 'active') RETURNING merchant_id", s.apiKey2, userID).Scan(&s.merchantID2)
	s.Require().NoError(err)

	// Seed cards for transactions
	cardNumber1 := "1111222233334444"
	cardNumber2 := "5555666677778888"
	_, err = s.dbPool.Exec(ctx, "INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES ($1, $2, 'debit', '123', 'visa', '2030-01-01')", userID, cardNumber1)
	s.NoError(err)
	_, err = s.dbPool.Exec(ctx, "INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES ($1, $2, 'debit', '123', 'mastercard', '2030-01-01')", userID, cardNumber2)
	s.NoError(err)

	// Seed Transactions
	// Merchant 1: Jan (500, visa), Feb (300, mastercard)
	s.dbPool.Exec(ctx, "INSERT INTO transactions (card_number, amount, payment_method, merchant_id, transaction_time, status) VALUES ($1, $2, $3, $4, $5, 'success')", cardNumber1, 500, "visa", s.merchantID1, time.Date(s.testYear, 1, 10, 10, 0, 0, 0, time.UTC))
	s.dbPool.Exec(ctx, "INSERT INTO transactions (card_number, amount, payment_method, merchant_id, transaction_time, status) VALUES ($1, $2, $3, $4, $5, 'success')", cardNumber2, 300, "mastercard", s.merchantID1, time.Date(s.testYear, 2, 10, 10, 0, 0, 0, time.UTC))

	// Merchant 2: Jan (1000, visa)
	s.dbPool.Exec(ctx, "INSERT INTO transactions (card_number, amount, payment_method, merchant_id, transaction_time, status) VALUES ($1, $2, $3, $4, $5, 'success')", cardNumber1, 1000, "visa", s.merchantID2, time.Date(s.testYear, 1, 15, 10, 0, 0, 0, time.UTC))
}

func (s *MerchantStatsRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	s.ts.Teardown()
}

func (s *MerchantStatsRepositoryTestSuite) TestGlobalStats() {
	ctx := context.Background()
	
	// Monthly Amount
	res, err := s.repo.GetMonthlyAmountMerchant(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)
	// Jan: 500 (M1) + 1000 (M2) = 1500
	s.Equal(int32(1500), res[0].TotalAmount)
	// Feb: 300
	s.Equal(int32(300), res[1].TotalAmount)

	// Monthly Method
	methRes, err := s.repo.GetMonthlyPaymentMethodsMerchant(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(methRes)
}


func (s *MerchantStatsRepositoryTestSuite) TestMerchantStats() {
	ctx := context.Background()

	// Merchant 1 Monthly
	res1, err := s.merchantRepo.GetMonthlyAmountByMerchants(ctx, &requests.MonthYearAmountMerchant{Year: s.testYear, MerchantID: int(s.merchantID1)})
	s.NoError(err)
	s.NotEmpty(res1)
	s.Equal(int32(500), res1[0].TotalAmount)
	s.Equal(int32(300), res1[1].TotalAmount)

	// Merchant 2 Monthly
	res2, err := s.merchantRepo.GetMonthlyAmountByMerchants(ctx, &requests.MonthYearAmountMerchant{Year: s.testYear, MerchantID: int(s.merchantID2)})
	s.NoError(err)
	s.Equal(int32(1000), res2[0].TotalAmount)
	s.Equal(int32(0), res2[1].TotalAmount)
}

func (s *MerchantStatsRepositoryTestSuite) TestApiKeyStats() {
	ctx := context.Background()

	// API Key 1 Monthly
	res1, err := s.apikeyRepo.GetMonthlyAmountByApikey(ctx, &requests.MonthYearAmountApiKey{Year: s.testYear, Apikey: s.apiKey1})
	s.NoError(err)
	s.NotEmpty(res1)
	s.Equal(int32(500), res1[0].TotalAmount)

	// API Key 2 Monthly
	res2, err := s.apikeyRepo.GetMonthlyAmountByApikey(ctx, &requests.MonthYearAmountApiKey{Year: s.testYear, Apikey: s.apiKey2})
	s.NoError(err)
	s.Equal(int32(1000), res2[0].TotalAmount)
}

func TestMerchantStatsRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantStatsRepositoryTestSuite))
}
