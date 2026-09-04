package withdraw_test

import (
	"context"
	"net"
	"testing"
	"time"

	pbwithdraw "github.com/MamangRust/monolith-payment-gateway-pb/withdraw"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/withdraw/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/handler"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/repository"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
)

type WithdrawStatsHandlerGapiTestSuite struct {
	suite.Suite
	ts           *tests.TestSuite
	lis          *bufconn.Listener
	conn         *grpc.ClientConn
	client       pb.WithdrawStatsStatusServiceClient
	clientAmount pb.WithdrawStatsAmountServiceClient
	userID       int32
	cardNumber   string
	testYear     int
	testMonth    int
}

func (s *WithdrawStatsHandlerGapiTestSuite) SetupSuite() {
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
	s.client = pb.NewWithdrawStatsStatusServiceClient(conn)
	s.clientAmount = pb.NewWithdrawStatsAmountServiceClient(conn)

	s.testYear = time.Now().Year()
	s.testMonth = int(time.Now().Month())

	ctx := context.Background()
	err = gormDB.WithContext(ctx).Raw("INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('WithdrawGapi', 'Stats', 'withdraw_gapi_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&s.userID).Error
	s.Require().NoError(err)

	s.cardNumber = "8888999900001111"
	err = gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", s.userID, s.cardNumber).Error
	s.Require().NoError(err)

	err = gormDB.WithContext(ctx).Exec("INSERT INTO withdraws (card_number, withdraw_amount, withdraw_time, status) VALUES (?, ?, ?, 'success')",
		s.cardNumber, 200000, time.Date(s.testYear, time.Month(s.testMonth), 10, 10, 0, 0, 0, time.UTC)).Error
	s.Require().NoError(err)
}

func (s *WithdrawStatsHandlerGapiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.lis != nil {
		s.lis.Close()
	}
	s.ts.Teardown()
}

func (s *WithdrawStatsHandlerGapiTestSuite) TestFindMonthlyWithdrawStatusSuccess() {
	ctx := context.Background()
	res, err := s.client.FindMonthlyWithdrawStatusSuccess(ctx, &pbwithdraw.FindMonthlyWithdrawStatus{
		Year:  int32(s.testYear),
		Month: int32(s.testMonth),
	})
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data) // generate_series returns 12 months
}

func (s *WithdrawStatsHandlerGapiTestSuite) TestFindMonthlyWithdraws() {
	ctx := context.Background()
	res, err := s.clientAmount.FindMonthlyWithdraws(ctx, &pbwithdraw.FindYearWithdrawStatus{
		Year: int32(s.testYear),
	})
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *WithdrawStatsHandlerGapiTestSuite) TestFindMonthlyWithdrawsByCardNumber() {
	ctx := context.Background()
	res, err := s.clientAmount.FindMonthlyWithdrawsByCardNumber(ctx, &pbwithdraw.FindYearWithdrawCardNumber{
		Year:       int32(s.testYear),
		CardNumber: s.cardNumber,
	})
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func TestWithdrawStatsHandlerGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(WithdrawStatsHandlerGapiTestSuite))
}
