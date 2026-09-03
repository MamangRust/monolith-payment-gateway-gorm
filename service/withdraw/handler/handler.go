package handler

import (
	"time"

	"gorm.io/gorm"

	withdrawstatshandler "github.com/MamangRust/monolith-payment-gateway-withdraw/handler/stats"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/service"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Handler interface {
	WithdrawQueryHandlerGrpc
	WithdrawCommandHandlerGrpc
	withdrawstatshandler.HandleStats
}

type handler struct {
	WithdrawQueryHandlerGrpc
	WithdrawCommandHandlerGrpc
	withdrawstatshandler.HandleStats
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
		WithdrawQueryHandlerGrpc:   NewWithdrawQueryHandleGrpc(service),
		WithdrawCommandHandlerGrpc: NewWithdrawCommandHandleGrpc(service),
		HandleStats:                withdrawstatshandler.NewWithdrawStatsHandleGrpc(service),
	}
}
