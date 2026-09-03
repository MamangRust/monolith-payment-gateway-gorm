package saldo_test

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

	saldo_handler "github.com/MamangRust/monolith-payment-gateway-apigateway/handler/saldo"
	pbcard "github.com/MamangRust/monolith-payment-gateway-pb/card"
	pbsaldo "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	pbuser "github.com/MamangRust/monolith-payment-gateway-pb/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	card_handler "github.com/MamangRust/monolith-payment-gateway-card/handler"
	card_repository "github.com/MamangRust/monolith-payment-gateway-card/repository"
	card_service "github.com/MamangRust/monolith-payment-gateway-card/service"
	"github.com/MamangRust/monolith-payment-gateway-saldo/handler"
	"github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-saldo/service"
	user_handler "github.com/MamangRust/monolith-payment-gateway-user/handler"
	user_repository "github.com/MamangRust/monolith-payment-gateway-user/repository"
	user_service "github.com/MamangRust/monolith-payment-gateway-user/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	redis_client "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SaldoApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	echo        *echo.Echo
	grpcServer  *grpc.Server
	lis         *bufconn.Listener
	userID      int32
	cardID      int32
	cardNumber  string
	saldoID     int32
}

func (s *SaldoApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis_client.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis_client.NewClient(opts)

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}
	repos := repository.NewRepositories(gormDB)
	cardRepos := card_repository.NewRepositories(gormDB)
	userRepos := user_repository.NewRepositories(gormDB)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	hasher := hash.NewHashingPassword()
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)

	saldoSvc := service.NewService(&service.Deps{
		Cache:        cacheStore,
		Repositories: repos,
		Logger:       log,
	})

	cardSvc := card_service.NewService(&card_service.Deps{
		Cache:        cacheStore,
		Repositories: cardRepos,
		Logger:       log,
		Kafka:        nil,
	})

	userSvc := user_service.NewService(&user_service.Deps{
		Cache:        cacheStore,
		Repositories: userRepos,
		Hash:         hasher,
		Logger:       log,
	})

	saldoH := handler.NewHandler(saldoSvc)
	cardH := card_handler.NewHandler(cardSvc)
	userH := user_handler.NewHandler(userSvc)

	// Setup gRPC server
	s.lis = bufconn.Listen(1024 * 1024)
	s.grpcServer = grpc.NewServer()
	pbsaldo.RegisterSaldoQueryServiceServer(s.grpcServer, saldoH)
	pbsaldo.RegisterSaldoCommandServiceServer(s.grpcServer, saldoH)

	go func() {
		if err := s.grpcServer.Serve(s.lis); err != nil {
		}
	}()

	s.echo = echo.New()

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return s.lis.Dial()
		}),
		grpc.WithInsecure())
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test", log)
	s.Require().NoError(err)
	apiHandler := errors.NewApiHandler(obs, log)

	saldo_handler.RegisterSaldoHandler(&saldo_handler.DepsSaldo{
		Client:     conn,
		E:          s.echo,
		Logger:     log,
		Cache:      cacheStore,
		ApiHandler: apiHandler,
	})

	// Create user
	ctx := context.Background()
	userRes, err := userH.Create(ctx, &pbuser.CreateUserRequest{
		Firstname:       "Saldo",
		Lastname:        "Api",
		Email:           fmt.Sprintf("saldo.api.%d@example.com", time.Now().UnixNano()),
		Password:        "Password123!",
		ConfirmPassword: "Password123!",
	})
	s.Require().NoError(err)
	s.userID = userRes.Data.Id

	// Create card
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
}

func (s *SaldoApiTestSuite) TearDownSuite() {
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

func (s *SaldoApiTestSuite) Test1_SaldoLifecycle() {
	// Create
	reqJSON := fmt.Sprintf(`{"card_number": "%s", "total_balance": 100000}`, s.cardNumber)
	req := httptest.NewRequest(http.MethodPost, "/api/saldo-command/create", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	s.NoError(err)
	data := res["data"].(map[string]interface{})
	s.saldoID = int32(data["id"].(float64))
	s.Equal(float64(100000), data["total_balance"])

	// FindById
	req = httptest.NewRequest(http.MethodGet, "/api/saldo-query/"+strconv.Itoa(int(s.saldoID)), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// Update
	updateJSON := fmt.Sprintf(`{"card_number": "%s", "total_balance": 200000}`, s.cardNumber)
	req = httptest.NewRequest(http.MethodPost, "/api/saldo-command/update/"+strconv.Itoa(int(s.saldoID)), strings.NewReader(updateJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *SaldoApiTestSuite) Test2_QueryOperations() {
	// FindAll
	req := httptest.NewRequest(http.MethodGet, "/api/saldo-query?page=1&page_size=10&search="+s.cardNumber, nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// FindActive
	req = httptest.NewRequest(http.MethodGet, "/api/saldo-query/active?page=1&page_size=10&search="+s.cardNumber, nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	
	// FindByCardNumber
	req = httptest.NewRequest(http.MethodGet, "/api/saldo-query/card_number/"+s.cardNumber, nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *SaldoApiTestSuite) Test3_TrashAndRestore() {
	// Trash
	req := httptest.NewRequest(http.MethodPost, "/api/saldo-command/trashed/"+strconv.Itoa(int(s.saldoID)), nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	
	// FindTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/saldo-query/trashed?page=1&page_size=10&search="+s.cardNumber, nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// Restore
	req = httptest.NewRequest(http.MethodPost, "/api/saldo-command/restore/"+strconv.Itoa(int(s.saldoID)), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestSaldoApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SaldoApiTestSuite))
}
