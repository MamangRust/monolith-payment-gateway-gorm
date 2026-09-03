package apps

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/user"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/server"
	"github.com/MamangRust/monolith-payment-gateway-user/handler"
	"github.com/MamangRust/monolith-payment-gateway-user/repository"
	"github.com/MamangRust/monolith-payment-gateway-user/service"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.GormDB)
	svc := service.NewService(&service.Deps{
		Cache:        srv.CacheStore,
		Logger:       srv.Logger,
		Repositories: repos,
		Hash:         hash.NewHashingPassword(),
	})
	h := handler.NewHandler(svc)

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterUserQueryServiceServer(gs, h)
		pb.RegisterUserCommandServiceServer(gs, h)
	}

	return srv, nil
}
