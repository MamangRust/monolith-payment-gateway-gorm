package apps

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-auth/handler"
	"github.com/MamangRust/monolith-payment-gateway-auth/repository"
	"github.com/MamangRust/monolith-payment-gateway-auth/service"

	pb "github.com/MamangRust/monolith-payment-gateway-pb"
	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/outbox"
	"github.com/MamangRust/monolith-payment-gateway-pkg/server"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to create token manager: %w", err)
	}

	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})
	srv.AddCleanupHook(myKafka.Close)
	relay, relayErr := outbox.NewRelay(srv.GormDB, myKafka, outbox.RelayConfig{})
	if relayErr != nil {
		srv.Cleanup()
		return nil, fmt.Errorf("initialize auth outbox relay: %w", relayErr)
	}
	{
		relayDone := make(chan struct{})
		go func() {
			defer close(relayDone)
			if relayErr := relay.Run(srv.Ctx); relayErr != nil && !errors.Is(relayErr, context.Canceled) {
				srv.Logger.Error("auth outbox relay stopped", zap.Error(relayErr))
			}
		}()
		srv.AddCleanupHook(func() error {
			select {
			case <-relayDone:
				return nil
			case <-time.After(5 * time.Second):
				return errors.New("timed out waiting for auth outbox relay")
			}
		})
	}

	hasher := hash.NewHashingPassword()
	repositories := repository.NewRepositories(srv.GormDB)
	services := service.NewService(&service.Deps{
		Cache:        srv.CacheStore,
		Repositories: repositories,
		Token:        tokenManager,
		Hash:         hasher,
		Logger:       srv.Logger,
	})

	handlers := handler.NewHandler(&handler.Deps{Service: services, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterAuthServiceServer(gs, handlers.Auth)
	}

	return srv, nil
}
