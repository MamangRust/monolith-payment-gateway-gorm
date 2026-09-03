package handler

import (
	"time"

	"gorm.io/gorm"

	transferstatshandler "github.com/MamangRust/monolith-payment-gateway-transfer/handler/stats"
	"github.com/MamangRust/monolith-payment-gateway-transfer/service"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Handler interface {
	TransferQueryHandleGrpc
	TransferCommandHandleGrpc
	transferstatshandler.HandleStats
}

type handler struct {
	TransferQueryHandleGrpc
	TransferCommandHandleGrpc
	transferstatshandler.HandleStats
}

func formatTimeWrapper(t *time.Time) *wrapperspb.StringValue {
	if t == nil {
		return nil
	}
	return &wrapperspb.StringValue{Value: t.Format(time.RFC3339)}
}

func formatGormDeletedAt(d gorm.DeletedAt) *wrapperspb.StringValue {
	if d.Valid {
		return &wrapperspb.StringValue{Value: d.Time.Format(time.RFC3339)}
	}
	return nil
}

func NewHandler(service service.Service) Handler {
	return &handler{
		TransferQueryHandleGrpc:   NewTransferQueryHandler(service),
		TransferCommandHandleGrpc: NewTransferCommandHandler(service),
		HandleStats:               transferstatshandler.NewTransferStatsHandleGrpc(service),
	}
}
