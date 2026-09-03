package transaction_test

import (
	"context"
	"net"
	"testing"
	"time"

	pbtransaction "github.com/MamangRust/monolith-payment-gateway-pb/transaction"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transaction/stats"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/MamangRust/monolith-payment-gateway-transaction/handler"
	"github.com/MamangRust/monolith-payment-gateway-transaction/repository"
	"github.com/MamangRust/monolith-payment-gateway-transaction/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	merchant_repo "github.com/MamangRust/monolith-payment-gateway-merchant/repository"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
)

type TransactionStatsHandlerGapiTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	dbPool     *pgxpool.Pool
	lis        *bufconn.Listener
	conn       *grpc.ClientConn
	client     pb.TransactionStatsStatusServiceClient
	userID     int32
	cardNumber string
	merchantID int32
	testYear   int
	testMonth  int
}

func (s *TransactionStatsHandlerGapiTestSuite) SetupSuite() {
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

	realSaldo := &realSaldoRepo{repo: saldo_repo.NewRepositories(gormDB)}
	realCard := &realCardRepo{
		query:   card_repo.NewCardQueryRepository(gormDB),
		command: card_repo.NewCardCommandRepository(gormDB),
	}
	realMerchant := &realMerchantRepo{repo: merchant_repo.NewMerchantQueryRepository(gormDB)}

	repos := repository.NewRepositories(gormDB, realSaldo, realCard, realMerchant)

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
	pb.RegisterTransactionStatsStatusServiceServer(server, h)

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
	s.client = pb.NewTransactionStatsStatusServiceClient(conn)

	s.testYear = time.Now().Year()
	s.testMonth = int(time.Now().Month())

	ctx := context.Background()
	err = s.dbPool.QueryRow(ctx, "INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ('TransactionGapi', 'Stats', 'transaction_gapi_stats@example.com', 'pass', '123', true) RETURNING user_id").Scan(&s.userID)
	s.Require().NoError(err)

	s.cardNumber = "2222333344445555"
	_, err = s.dbPool.Exec(ctx, "INSERT INTO cards (user_id, card_number, card_type, cvv, card_provider, expire_date) VALUES ($1, $2, 'debit', '123', 'visa', '2030-01-01')", s.userID, s.cardNumber)
	s.Require().NoError(err)

	err = s.dbPool.QueryRow(ctx, "INSERT INTO merchants (name, api_key, user_id, status) VALUES ('Merchant Gapi', 'test_key_gapi', $1, 'active') RETURNING merchant_id", s.userID).Scan(&s.merchantID)
	s.Require().NoError(err)

	_, err = s.dbPool.Exec(ctx, "INSERT INTO transactions (card_number, merchant_id, amount, payment_method, transaction_time, status) VALUES ($1, $2, $3, 'credit_card', $4, 'success')",
		s.cardNumber, s.merchantID, 200000, time.Date(s.testYear, time.Month(s.testMonth), 10, 10, 0, 0, 0, time.UTC))
	s.Require().NoError(err)
}

func (s *TransactionStatsHandlerGapiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.lis != nil {
		s.lis.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	s.ts.Teardown()
}

func (s *TransactionStatsHandlerGapiTestSuite) TestFindMonthlyTransactionStatusSuccess() {
	ctx := context.Background()
	res, err := s.client.FindMonthlyTransactionStatusSuccess(ctx, &pbtransaction.FindMonthlyTransactionStatus{
		Year:  int32(s.testYear),
		Month: int32(s.testMonth),
	})
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func TestTransactionStatsHandlerGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsHandlerGapiTestSuite))
}
