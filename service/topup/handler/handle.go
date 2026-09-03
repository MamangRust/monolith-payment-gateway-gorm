package handler

import (
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/convert"

	topupstatshandler "github.com/MamangRust/monolith-payment-gateway-topup/handler/stats"
	"github.com/MamangRust/monolith-payment-gateway-topup/service"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

func formatGormDeletedAt(dt gorm.DeletedAt) *wrapperspb.StringValue {
	if dt.Time.IsZero() {
		return wrapperspb.String("")
	}
	return wrapperspb.String(dt.Time.Format(time.RFC3339))
}

func formatPtrTime(t *time.Time) *wrapperspb.StringValue {
	if t == nil {
		return wrapperspb.String("")
	}
	return convert.TimeToWrappers(t)
}

type Handler interface {
	TopupQueryHandleGrpc
	TopupCommandHandleGrpc
	topupstatshandler.HandleStats
}

type handler struct {
	TopupQueryHandleGrpc
	TopupCommandHandleGrpc
	topupstatshandler.HandleStats
}

func NewHandler(service service.Service) Handler {
	return &handler{
		TopupQueryHandleGrpc:   NewTopupQueryHandleGrpc(service),
		TopupCommandHandleGrpc: NewTopupCommandHandleGrpc(service),
		HandleStats:            topupstatshandler.NewTopupStatsHandleGrpc(service),
	}
}
