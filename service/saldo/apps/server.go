package apps

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	pbstats "github.com/MamangRust/monolith-payment-gateway-pb/saldo/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/server"
	"github.com/MamangRust/monolith-payment-gateway-saldo/handler"
	saldokafka "github.com/MamangRust/monolith-payment-gateway-saldo/kafka"
	"github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-saldo/service"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.GormDB)

	mykafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})
	srv.AddCleanupHook(mykafka.Close)

	svc := service.NewService(&service.Deps{
		Cache:        srv.CacheStore,
		Logger:       srv.Logger,
		Repositories: repos,
	})

	kafkaHandler := saldokafka.NewSaldoKafkaHandler(svc, srv.Logger, context.Background())
	_, err = mykafka.StartConsumersWithContext(srv.Ctx, []string{
		"saldo-service-topic-create-saldo",
		"saldo-service-topic-invalidate-cache",
	}, "saldo-service-group", kafkaHandler)
	if err != nil {
		srv.Logger.Error("Failed to start kafka consumers", zap.Error(err))
	}

	h := handler.NewHandler(svc)

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterSaldoQueryServiceServer(gs, h)
		pb.RegisterSaldoCommandServiceServer(gs, h)
		pbstats.RegisterSaldoStatsBalanceServiceServer(gs, h)
		pbstats.RegisterSaldoStatsTotalBalanceServer(gs, h)
	}

	return srv, nil
}
