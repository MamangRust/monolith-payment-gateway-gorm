package merchant_test

import (
	"context"
	"testing"
	"time"

	stats_cache "github.com/MamangRust/monolith-payment-gateway-merchant/redis/stats"
	apikey_cache "github.com/MamangRust/monolith-payment-gateway-merchant/redis/statsbyapikey"
	merchant_cache "github.com/MamangRust/monolith-payment-gateway-merchant/redis/statsbymerchant"
	stats_repo "github.com/MamangRust/monolith-payment-gateway-merchant/repository/stats"
	apikey_repo "github.com/MamangRust/monolith-payment-gateway-merchant/repository/statsbyapikey"
	merchant_repo "github.com/MamangRust/monolith-payment-gateway-merchant/repository/statsbymerchant"
	stats_service "github.com/MamangRust/monolith-payment-gateway-merchant/service/stats"
	apikey_service "github.com/MamangRust/monolith-payment-gateway-merchant/service/statsbyapikey"
	merchant_service "github.com/MamangRust/monolith-payment-gateway-merchant/service/statsbymerchant"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MerchantStatsServiceTestSuite struct {
	suite.Suite
	ts              *tests.TestSuite
	service         stats_service.MerchantStatsAmountService
	apikeyService   apikey_service.MerchantStatsByApiKeyAmountService
	merchantService merchant_service.MerchantStatsByMerchantAmountService
	merchantID1     int32
	apiKey1         string
	testYear        int
}

func (s *MerchantStatsServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(opts)

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}

	// Logger & Observability
	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	obs, _ := observability.NewObservability("test", log)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)

	// Caches
	s_cache := stats_cache.NewMerchantStatsCache(cacheStore)
	a_cache := apikey_cache.NewMerchantStatsByApiKeyCache(cacheStore)
	m_cache := merchant_cache.NewMerchantStatsByMerchantCache(cacheStore)

	// Repositories
	s_repo := stats_repo.NewMerchantStatsRepository(gormDB)
	a_repo := apikey_repo.NewMerchantStatsByApiKeyRepository(gormDB)
	m_repo := merchant_repo.NewMerchantStatsByMerchantRepository(gormDB)

	// Services
	s.service = stats_service.NewMerchantStatsAmountService(&stats_service.MerchantStatsAmountDeps{
		Cache:         s_cache,
		Logger:        log,
		Repository:    s_repo,
		Observability: obs,
	})
	s.apikeyService = apikey_service.NewMerchantStatsAmountByApiKeyService(&apikey_service.MerchantStatsAmountByApiKeyDeps{
		Cache:         a_cache,
		Logger:        log,
		Repository:    a_repo,
		Observability: obs,
	})
	s.merchantService = merchant_service.NewMerchantStatsByMerchantAmountService(&merchant_service.MerchantStatsByMerchantAmountServiceDeps{
		Cache:         m_cache,
		Logger:        log,
		Repository:    m_repo,
		Observability: obs,
	})

	s.testYear = time.Now().Year()

	// Seed Data
	ctx := context.Background()
	var userID int32
	err = gormDB.WithContext(ctx).Raw("INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('MerchantSvc', 'Stats', 'merchant_svc_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&userID).Error
	s.Require().NoError(err)

	s.apiKey1 = "svc-merchant-key-1"
	err = gormDB.WithContext(ctx).Raw("INSERT INTO merchants (name, api_key, user_id, status) VALUES ('Svc Merchant 1', ?, ?, 'active') RETURNING merchant_id", s.apiKey1, userID).Scan(&s.merchantID1).Error
	s.Require().NoError(err)

	cardNumber1 := "9999888877776666"
	err = gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", userID, cardNumber1).Error
	s.Require().NoError(err)

	// Transaction: Jan (1000)
	err = gormDB.WithContext(ctx).Exec("INSERT INTO transactions (card_number, amount, payment_method, merchant_id, transaction_time, status) VALUES (?, ?, ?, ?, ?, 'success')", cardNumber1, 1000, "visa", s.merchantID1, time.Date(s.testYear, 1, 10, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)
}

func (s *MerchantStatsServiceTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *MerchantStatsServiceTestSuite) TestGlobalService() {
	ctx := context.Background()
	res, err := s.service.FindMonthlyAmountMerchant(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)

}

func (s *MerchantStatsServiceTestSuite) TestMerchantService() {
	ctx := context.Background()
	res, err := s.merchantService.FindMonthlyAmountByMerchants(ctx, &requests.MonthYearAmountMerchant{Year: s.testYear, MerchantID: int(s.merchantID1)})
	s.NoError(err)
	s.NotEmpty(res)

}

func (s *MerchantStatsServiceTestSuite) TestApiKeyService() {
	ctx := context.Background()
	res, err := s.apikeyService.FindMonthlyAmountByApikey(ctx, &requests.MonthYearAmountApiKey{Year: s.testYear, Apikey: s.apiKey1})
	s.NoError(err)
	s.NotEmpty(res)

}

func TestMerchantStatsServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantStatsServiceTestSuite))
}
