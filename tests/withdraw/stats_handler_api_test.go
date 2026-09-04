package withdraw_test

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"

	withdraw_handler "github.com/MamangRust/monolith-payment-gateway-apigateway/handler/withdraw"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/withdraw/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/handler"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/repository"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/service"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type WithdrawStatsHandlerApiTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	echo       *echo.Echo
	lis        *bufconn.Listener
	conn       *grpc.ClientConn
	userID     int32
	cardNumber string
	testYear   int
	testMonth  int
}

func (s *WithdrawStatsHandlerApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts
	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}

	realSaldo := &realSaldoRepo{repo: saldo_repo.NewRepositories(gormDB)}
	realCard := &realCardRepo{
		query: card_repo.NewCardQueryRepository(gormDB),
	}

	repos := repository.NewRepositories(gormDB, realCard, realSaldo)

	zapLog := zap.NewNop()
	myLogger := &logger.Logger{Log: zapLog}

	redisOption, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(redisOption)
	cacheStore := cache.NewCacheStore(redisClient, myLogger, &dummyCacheMetrics{})

	svc := service.NewService(&service.Deps{
		Kafka:        nil,
		Repositories: repos,
		Logger:       myLogger,
		Cache:        cacheStore,
	})

	h := handler.NewHandler(svc)

	s.lis = bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pb.RegisterWithdrawStatsStatusServiceServer(server, h)
	pb.RegisterWithdrawStatsAmountServiceServer(server, h)

	go func() {
		if err := server.Serve(s.lis); err != nil {
		}
	}()

	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return s.lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn

	s.echo = echo.New()
	obs, _ := observability.NewObservability("test", myLogger)
	apiHandler := errors.NewApiHandler(obs, myLogger)

	withdraw_handler.RegisterWithdrawHandler(&withdraw_handler.DepsWithdraw{
		Client:     s.conn,
		E:          s.echo,
		Logger:     myLogger,
		Cache:      cacheStore,
		ApiHandler: apiHandler,
	})

	s.testYear = time.Now().Year()
	s.testMonth = int(time.Now().Month())

	ctx := context.Background()
	err = gormDB.WithContext(ctx).Raw("INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('WithdrawApi', 'Stats', 'withdraw_api_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&s.userID).Error
	s.Require().NoError(err)

	s.cardNumber = "9999000011112222"
	err = gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", s.userID, s.cardNumber).Error
	s.Require().NoError(err)

	err = gormDB.WithContext(ctx).Exec("INSERT INTO withdraws (card_number, withdraw_amount, withdraw_time, status) VALUES (?, ?, ?, 'success')",
		s.cardNumber, 300000, time.Date(s.testYear, time.Month(s.testMonth), 10, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)
}

func (s *WithdrawStatsHandlerApiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.lis != nil {
		s.lis.Close()
	}
	s.ts.Teardown()
}

func (s *WithdrawStatsHandlerApiTestSuite) TestFindMonthlyWithdrawStatusSuccess() {
	year := strconv.Itoa(s.testYear)
	month := strconv.Itoa(s.testMonth)
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw-stats-status/monthly-success?year="+year+"&month="+month, nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var resp struct {
		Status string        `json:"status"`
		Data   []interface{} `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	s.Equal("success", resp.Status)
	s.NotEmpty(resp.Data)
}

func (s *WithdrawStatsHandlerApiTestSuite) TestFindYearlyWithdrawStatusSuccess() {
	year := strconv.Itoa(s.testYear)
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw-stats-status/yearly-success?year="+year, nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var resp struct {
		Status string        `json:"status"`
		Data   []interface{} `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	s.Equal("success", resp.Status)
	s.NotEmpty(resp.Data)
}

func (s *WithdrawStatsHandlerApiTestSuite) TestFindMonthlyWithdrawAmount() {
	year := strconv.Itoa(s.testYear)
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw-stats-amount/monthly-amount?year="+year, nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var resp struct {
		Status string        `json:"status"`
		Data   []interface{} `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	s.Equal("success", resp.Status)
	s.NotEmpty(resp.Data)
}

func (s *WithdrawStatsHandlerApiTestSuite) TestFindMonthlyWithdrawAmountByCardNumber() {
	year := strconv.Itoa(s.testYear)
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw-stats-amount/monthly-amount-card?year="+year+"&card_number="+s.cardNumber, nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var resp struct {
		Status string        `json:"status"`
		Data   []interface{} `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	s.Equal("success", resp.Status)
	s.NotEmpty(resp.Data)
}

func TestWithdrawStatsHandlerApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(WithdrawStatsHandlerApiTestSuite))
}
