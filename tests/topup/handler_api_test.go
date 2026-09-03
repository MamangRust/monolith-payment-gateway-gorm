package topup_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	topup_handler "github.com/MamangRust/monolith-payment-gateway-apigateway/handler/topup"
	"github.com/MamangRust/monolith-payment-gateway-pkg/adapter"
	pbcard "github.com/MamangRust/monolith-payment-gateway-pb/card"
	pbsaldo "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/topup"
	pbuser "github.com/MamangRust/monolith-payment-gateway-pb/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	card_handler "github.com/MamangRust/monolith-payment-gateway-card/handler"
	card_repository "github.com/MamangRust/monolith-payment-gateway-card/repository"
	card_service "github.com/MamangRust/monolith-payment-gateway-card/service"
	saldo_handler "github.com/MamangRust/monolith-payment-gateway-saldo/handler"
	saldo_repository "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	saldo_service "github.com/MamangRust/monolith-payment-gateway-saldo/service"
	"github.com/MamangRust/monolith-payment-gateway-topup/handler"
	"github.com/MamangRust/monolith-payment-gateway-topup/repository"
	"github.com/MamangRust/monolith-payment-gateway-topup/service"
	user_handler "github.com/MamangRust/monolith-payment-gateway-user/handler"
	user_repository "github.com/MamangRust/monolith-payment-gateway-user/repository"
	user_service "github.com/MamangRust/monolith-payment-gateway-user/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TopupApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	echo        *echo.Echo
	userID      int32
	cardID      int32
	cardNumber  string
	topupID     int
	grpcServer  *grpc.Server
	lis         *bufconn.Listener
	conn        *grpc.ClientConn
}

func (s *TopupApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(opts)

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}
	
	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	hasher := hash.NewHashingPassword()
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)

	// Dependency services handlers
	cardSvc := card_service.NewService(&card_service.Deps{
		Cache:        cacheStore,
		Repositories: card_repository.NewRepositories(gormDB),
		Logger:       log,
		Kafka:        nil,
	})
	cardH := card_handler.NewHandler(cardSvc)

	saldoSvc := saldo_service.NewService(&saldo_service.Deps{
		Cache:        cacheStore,
		Repositories: saldo_repository.NewRepositories(gormDB),
		Logger:       log,
	})
	saldoH := saldo_handler.NewHandler(saldoSvc)

	userSvc := user_service.NewService(&user_service.Deps{
		Cache:        cacheStore,
		Repositories: user_repository.NewRepositories(gormDB),
		Hash:         hasher,
		Logger:       log,
	})
	userH := user_handler.NewHandler(userSvc)

	// Setup gRPC server for adapters and Topup GAPI
	s.lis = bufconn.Listen(1024 * 1024)
	s.grpcServer = grpc.NewServer()
	
	// Register dependencies
	pbcard.RegisterCardQueryServiceServer(s.grpcServer, cardH)
	pbcard.RegisterCardCommandServiceServer(s.grpcServer, cardH)
	pbsaldo.RegisterSaldoQueryServiceServer(s.grpcServer, saldoH)
	pbsaldo.RegisterSaldoCommandServiceServer(s.grpcServer, saldoH)

	// Dial to ourselves
	s.conn, err = grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return s.lis.Dial()
		}),
		grpc.WithInsecure())
	s.Require().NoError(err)

	cardAdapter := adapter.NewCardAdapter(pbcard.NewCardQueryServiceClient(s.conn), pbcard.NewCardCommandServiceClient(s.conn))
	saldoAdapter := adapter.NewSaldoAdapter(pbsaldo.NewSaldoQueryServiceClient(s.conn), pbsaldo.NewSaldoCommandServiceClient(s.conn))

	// Topup Service
	repos := repository.NewRepositories(gormDB, cardAdapter, saldoAdapter)
	topupSvc := service.NewService(&service.Deps{
		Kafka:        nil,
		Cache:        cacheStore,
		Repositories: repos,
		Logger:       log,
	})
	topupH := handler.NewHandler(topupSvc)

	// Register Topup GAPI on the same server
	pb.RegisterTopupQueryServiceServer(s.grpcServer, topupH)
	pb.RegisterTopupCommandServiceServer(s.grpcServer, topupH)

	go func() {
		if err := s.grpcServer.Serve(s.lis); err != nil {
		}
	}()

	s.echo = echo.New()
	obs, _ := observability.NewObservability("test", log)
	apiHandler := errors.NewApiHandler(obs, log)

	// Register Topup API Handler (from apigateway)
	topup_handler.RegisterTopupHandler(&topup_handler.DepsTopup{
		Client:     s.conn,
		E:          s.echo,
		Logger:     log,
		Cache:      cacheStore,
		ApiHandler: apiHandler,
	})

	// Create test environment
	ctx := context.Background()
	userRes, err := userH.Create(ctx, &pbuser.CreateUserRequest{
		Firstname:       "Topup",
		Lastname:        "Api",
		Email:           fmt.Sprintf("topup.api.%d@example.com", time.Now().UnixNano()),
		Password:        "Password123!",
		ConfirmPassword: "Password123!",
	})
	s.Require().NoError(err)
	s.userID = userRes.Data.Id

	cardRes, err := cardH.CreateCard(ctx, &pbcard.CreateCardRequest{
		UserId:       s.userID,
		CardType:     "debit",
		ExpireDate:   timestamppb.New(time.Now().AddDate(5, 0, 0)),
		Cvv:          "123",
		CardProvider: "visa",
	})
	s.Require().NoError(err)
	s.cardID = cardRes.Data.Id
	s.cardNumber = cardRes.Data.CardNumber

	_, err = saldoH.CreateSaldo(ctx, &pbsaldo.CreateSaldoRequest{
		CardNumber:   s.cardNumber,
		TotalBalance: 100000,
	})
	s.Require().NoError(err)
}

func (s *TopupApiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *TopupApiTestSuite) Test1_CreateTopup() {
	body := map[string]interface{}{
		"card_number":  s.cardNumber,
		"topup_amount": 50000,
		"topup_method": "bca",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/topup-command/create", strings.NewReader(string(jsonBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	s.topupID = int(data["id"].(float64))
	s.Equal(float64(50000), data["topup_amount"])
}

func (s *TopupApiTestSuite) Test2_UpdateTopup() {
	s.Require().NotZero(s.topupID)
	body := map[string]interface{}{
		"topup_id":     s.topupID,
		"card_number":  s.cardNumber,
		"topup_amount": 75000,
		"topup_method": "mandiri",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/topup-command/update/"+strconv.Itoa(s.topupID), strings.NewReader(string(jsonBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *TopupApiTestSuite) Test3_FindOperations() {
	// FindById
	req := httptest.NewRequest(http.MethodGet, "/api/topup-query/"+strconv.Itoa(s.topupID), nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// FindAll
	reqA := httptest.NewRequest(http.MethodGet, "/api/topup-query?page=1&page_size=10", nil)
	recA := httptest.NewRecorder()
	s.echo.ServeHTTP(recA, reqA)
	s.Equal(http.StatusOK, recA.Code)
}

func (s *TopupApiTestSuite) Test4_TrashAndRestore() {
	// Trash
	reqT := httptest.NewRequest(http.MethodPost, "/api/topup-command/trashed/"+strconv.Itoa(s.topupID), nil)
	recT := httptest.NewRecorder()
	s.echo.ServeHTTP(recT, reqT)
	s.Equal(http.StatusOK, recT.Code)

	// Restore
	reqR := httptest.NewRequest(http.MethodPost, "/api/topup-command/restore/"+strconv.Itoa(s.topupID), nil)
	recR := httptest.NewRecorder()
	s.echo.ServeHTTP(recR, reqR)
	s.Equal(http.StatusOK, recR.Code)
}

func TestTopupApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TopupApiTestSuite))
}
