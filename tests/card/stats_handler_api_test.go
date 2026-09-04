package card_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apihandler "github.com/MamangRust/monolith-payment-gateway-apigateway/handler/card"
	"github.com/MamangRust/monolith-payment-gateway-card/handler"
	"github.com/MamangRust/monolith-payment-gateway-card/repository"
	"github.com/MamangRust/monolith-payment-gateway-card/service"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
	pbstats "github.com/MamangRust/monolith-payment-gateway-pb/card/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CardStatsApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	gormDB      *gorm.DB
	echo        *echo.Echo
	cardNumber1 string
	testYear    int
	grpcServer  *grpc.Server
	lis         *bufconn.Listener
	conn        *grpc.ClientConn
}

func (s *CardStatsApiTestSuite) SetupSuite() {
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
	s.gormDB = gormDB
	repos := repository.NewRepositories(gormDB)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)

	cardSvc := service.NewService(&service.Deps{
		Repositories: repos,
		Logger:       log,
		Cache:        cacheStore,
		Kafka:        nil,
	})
	cardH := handler.NewHandler(cardSvc)

	// Setup gRPC server
	s.lis = bufconn.Listen(1024 * 1024)
	s.grpcServer = grpc.NewServer()

	pb.RegisterCardQueryServiceServer(s.grpcServer, cardH)
	pbstats.RegisterCardStatsBalanceServiceServer(s.grpcServer, cardH)
	pbstats.RegisterCardStatsTopupServiceServer(s.grpcServer, cardH)
	pbstats.RegisterCardStatsTransferServiceServer(s.grpcServer, cardH)
	pbstats.RegisterCardStatsWithdrawServiceServer(s.grpcServer, cardH)
	pbstats.RegisterCardStatsTransactionServiceServer(s.grpcServer, cardH)

	go func() {
		if err := s.grpcServer.Serve(s.lis); err != nil {
		}
	}()

	s.conn, err = grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return s.lis.Dial()
		}),
		grpc.WithInsecure())
	s.Require().NoError(err)

	s.echo = echo.New()
	obs, _ := observability.NewObservability("test", log)
	apiHandler := errors.NewApiHandler(obs, log)

	apihandler.RegisterCardHandler(&apihandler.DepsCard{
		Client:     s.conn,
		E:          s.echo,
		Logger:     log,
		Cache:      cacheStore,
		ApiHandler: apiHandler,
	})

	s.testYear = time.Now().Year()

	// Seed data
	ctx := context.Background()
	s.cardNumber1 = "4444555566667777"

	var userID int32
	err = s.gormDB.WithContext(ctx).Raw("INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('Api', 'Stats', 'api_stats_card@example.com', 'pass', '123', true) RETURNING user_id").Scan(&userID).Error
	s.Require().NoError(err)

	err = s.gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", userID, s.cardNumber1).Error
	s.Require().NoError(err)

	s.seedHistoricalData()
}

func (s *CardStatsApiTestSuite) seedHistoricalData() {
	err := s.gormDB.WithContext(context.Background()).Exec("INSERT INTO saldos (card_number, total_balance, created_at, updated_at) VALUES (?, ?, ?, ?)", s.cardNumber1, 1000, time.Date(s.testYear, 1, 15, 10, 0, 0, 0, time.UTC), time.Date(s.testYear, 1, 15, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)
	err = s.gormDB.WithContext(context.Background()).Exec("INSERT INTO topups (card_number, topup_amount, topup_time, topup_method, status) VALUES (?, ?, ?, 'api', 'success')", s.cardNumber1, 500, time.Date(s.testYear, 1, 10, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)
}

func (s *CardStatsApiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *CardStatsApiTestSuite) TestBalanceApi() {
	// Global Monthly
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/card-stats-balance/monthly-balance?year=%d", s.testYear), nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	firstMonth := data[0].(map[string]interface{})
	s.Equal(float64(1000), firstMonth["total_balance"])

	// By Card Monthly
	reqByCard := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/card-stats-balance/monthly-balance-by-card?year=%d&card_number=%s", s.testYear, s.cardNumber1), nil)
	recByCard := httptest.NewRecorder()
	s.echo.ServeHTTP(recByCard, reqByCard)

	s.Equal(http.StatusOK, recByCard.Code)
}

func (s *CardStatsApiTestSuite) TestTopupApi() {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/card-stats-topup/monthly-topup-amount?year=%d", s.testYear), nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	firstMonth := data[0].(map[string]interface{})
	s.Equal(float64(500), firstMonth["total_amount"])
}

func TestCardStatsApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CardStatsApiTestSuite))
}
