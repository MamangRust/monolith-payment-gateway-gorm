package saldo_test

import (
	"context"
	"net"
	"testing"
	"time"

	pbsaldo "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/handler"
	"github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-saldo/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SaldoStatsHandlerGapiTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	lis        *bufconn.Listener
	conn       *grpc.ClientConn
	client     pb.SaldoStatsBalanceServiceClient
	userID     int32
	cardNumber string
	testYear   int
}

func (s *SaldoStatsHandlerGapiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}
	repos := repository.NewRepositories(gormDB)

	zapLog := zap.NewNop()
	myLogger := &logger.Logger{Log: zapLog}

	redisOption, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(redisOption)
	cacheStore := cache.NewCacheStore(redisClient, myLogger, &dummyCacheMetrics{})

	svc := service.NewService(&service.Deps{
		Repositories: repos,
		Logger:       myLogger,
		Cache:        cacheStore,
	})

	h := handler.NewHandler(svc)

	s.lis = bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pb.RegisterSaldoStatsBalanceServiceServer(server, h)
	pb.RegisterSaldoStatsTotalBalanceServer(server, h)

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
	s.client = pb.NewSaldoStatsBalanceServiceClient(conn)

	s.testYear = time.Now().Year()

	ctx := context.Background()
	err = gormDB.WithContext(ctx).Raw("INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('SaldoGapi', 'Stats', 'saldo_gapi_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&s.userID).Error
	s.Require().NoError(err)

	s.cardNumber = "7777888899990000"
	err = gormDB.WithContext(ctx).Exec("INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES (?, ?, 'debit', '123', 'visa', '2030-01-01')", s.userID, s.cardNumber).Error
	s.Require().NoError(err)
	err = gormDB.WithContext(ctx).Exec("INSERT INTO saldos (card_number, total_balance) VALUES (?, ?)", s.cardNumber, 2000000).Error
	s.Require().NoError(err)
}

func (s *SaldoStatsHandlerGapiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.lis != nil {
		s.lis.Close()
	}
	s.ts.Teardown()
}

func (s *SaldoStatsHandlerGapiTestSuite) TestFindMonthlySaldoBalances() {
	ctx := context.Background()
	res, err := s.client.FindMonthlySaldoBalances(ctx, &pbsaldo.FindYearlySaldo{
		Year: int32(s.testYear),
	})
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func TestSaldoStatsHandlerGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SaldoStatsHandlerGapiTestSuite))
}
