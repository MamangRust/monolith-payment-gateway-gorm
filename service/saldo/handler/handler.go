package handler

import (
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/convert"

	saldostatshandler "github.com/MamangRust/monolith-payment-gateway-saldo/handler/stats"
	"github.com/MamangRust/monolith-payment-gateway-saldo/service"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

func formatPtrTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func formatGormDeletedAt(d gorm.DeletedAt) *wrapperspb.StringValue {
	if !d.Valid {
		return nil
	}
	return convert.TimeToWrappers(&d.Time)
}

func formatPtrDeletedAt(d *time.Time) *wrapperspb.StringValue {
	if d == nil {
		return nil
	}
	return convert.TimeToWrappers(d)
}

type handler struct {
	SaldoQueryHandleGrpc
	SaldoCommandHandleGrpc
	saldostatshandler.HandleStats
}

type Handler interface {
	SaldoQueryHandleGrpc
	SaldoCommandHandleGrpc
	saldostatshandler.HandleStats
}

func NewHandler(service service.Service) Handler {
	return &handler{
		SaldoQueryHandleGrpc:   NewSaldoQueryHandleGrpc(service),
		SaldoCommandHandleGrpc: NewSaldoCommandHandleGrpc(service),
		HandleStats:            saldostatshandler.NewSaldoStatsHandle(service),
	}
}
