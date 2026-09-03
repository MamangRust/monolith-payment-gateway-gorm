package transfer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-transfer/repository"
	"github.com/MamangRust/monolith-payment-gateway-transfer/service"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type p0TransferCardRepo struct{}

func (p0TransferCardRepo) FindUserCardByCardNumber(context.Context, string) (*models.CardByEmailRow, error) {
	return &models.CardByEmailRow{CardNumber: "sender", Email: "sender@example.com"}, nil
}

func (p0TransferCardRepo) FindCardByCardNumber(context.Context, string) (*models.CardAllFieldsRow, error) {
	return &models.CardAllFieldsRow{CardNumber: "receiver"}, nil
}

type p0TransferSaldoRepo struct {
	balances map[string]int32
}

func (r *p0TransferSaldoRepo) FindByCardNumber(_ context.Context, cardNumber string) (*models.Saldo, error) {
	balance, ok := r.balances[cardNumber]
	if !ok {
		return nil, errors.New("saldo not found")
	}
	return &models.Saldo{CardNumber: cardNumber, TotalBalance: balance}, nil
}

func (r *p0TransferSaldoRepo) UpdateSaldoBalance(_ context.Context, request *requests.UpdateSaldoBalance) (*models.UpdateSaldoBalanceRow, error) {
	if _, ok := r.balances[request.CardNumber]; !ok {
		return nil, errors.New("saldo not found")
	}
	r.balances[request.CardNumber] = int32(request.TotalBalance)
	return &models.UpdateSaldoBalanceRow{CardNumber: request.CardNumber, TotalBalance: int32(request.TotalBalance)}, nil
}

type p0TransferCommandRepo struct {
	createErr        error
	statusErr        error
	createTransferID int32
	statusCalls      int
	balances         map[string]int32
}

func (r *p0TransferCommandRepo) CreateTransfer(context.Context, *requests.CreateTransferRequest) (*models.TransferAllFieldsRow, error) {
	if r.createErr != nil {
		return nil, r.createErr
	}
	return &models.TransferAllFieldsRow{
		TransferID:     r.createTransferID,
		Status:         "success",
		TransferAmount: 50000,
	}, nil
}

func (r *p0TransferCommandRepo) CreateTransferAtomic(_ context.Context, req *requests.CreateTransferRequest) (*repository.TransferAtomicResult, error) {
	if r.createErr != nil {
		return nil, r.createErr
	}
	if r.balances != nil {
		r.balances[req.TransferFrom] -= int32(req.TransferAmount)
		r.balances[req.TransferTo] += int32(req.TransferAmount)
	}
	return &repository.TransferAtomicResult{
		Row: &models.TransferAllFieldsRow{
			TransferID:     r.createTransferID,
			TransferFrom:   req.TransferFrom,
			TransferTo:     req.TransferTo,
			TransferAmount: int32(req.TransferAmount),
			Status:         "success",
		},
	}, nil
}

func (r *p0TransferCommandRepo) UpdateTransferStatus(context.Context, *requests.UpdateTransferStatus) (*models.TransferAllFieldsRow, error) {
	r.statusCalls++
	if r.statusErr != nil {
		return nil, r.statusErr
	}
	return &models.TransferAllFieldsRow{}, nil
}

func (r *p0TransferCommandRepo) UpdateTransferAtomic(context.Context, models.TransferAllFieldsRow) (*models.TransferAllFieldsRow, error) {
	return &models.TransferAllFieldsRow{}, nil
}

func (r *p0TransferCommandRepo) UpdateTransferSettlementDelta(_ context.Context, senderCard, receiverCard string, senderDelta, receiverDelta int) error {
	if r.balances != nil {
		r.balances[senderCard] += int32(senderDelta)
		r.balances[receiverCard] += int32(receiverDelta)
	}
	return nil
}

func (p0TransferCommandRepo) UpdateTransfer(context.Context, *requests.UpdateTransferRequest) (*models.TransferAllFieldsRow, error) {
	return &models.TransferAllFieldsRow{}, nil
}
func (p0TransferCommandRepo) UpdateTransferAmount(context.Context, *requests.UpdateTransferAmountRequest) (*models.TransferAllFieldsRow, error) {
	return &models.TransferAllFieldsRow{}, nil
}
func (p0TransferCommandRepo) TrashedTransfer(context.Context, int) (*models.Transfer, error) {
	return &models.Transfer{}, nil
}
func (p0TransferCommandRepo) RestoreTransfer(context.Context, int) (*models.Transfer, error) {
	return &models.Transfer{}, nil
}
func (p0TransferCommandRepo) DeleteTransferPermanent(context.Context, int) (bool, error) {
	return true, nil
}
func (p0TransferCommandRepo) RestoreAllTransfer(context.Context) (bool, error) {
	return true, nil
}
func (p0TransferCommandRepo) DeleteAllTransferPermanent(context.Context) (bool, error) {
	return true, nil
}

type p0TransferQueryRepo struct{}

func (p0TransferQueryRepo) FindAll(context.Context, *requests.FindAllTransfers) ([]*models.TransferListRow, error) {
	return []*models.TransferListRow{}, nil
}
func (p0TransferQueryRepo) FindByActive(context.Context, *requests.FindAllTransfers) ([]*models.TransferListRow, error) {
	return []*models.TransferListRow{}, nil
}
func (p0TransferQueryRepo) FindByTrashed(context.Context, *requests.FindAllTransfers) ([]*models.TransferListWithDeletedRow, error) {
	return []*models.TransferListWithDeletedRow{}, nil
}
func (p0TransferQueryRepo) FindById(context.Context, int) (*models.TransferAllFieldsRow, error) {
	return &models.TransferAllFieldsRow{}, nil
}
func (p0TransferQueryRepo) FindTransferByTransferFrom(context.Context, string) ([]*models.TransferBySourceCardRow, error) {
	return []*models.TransferBySourceCardRow{}, nil
}
func (p0TransferQueryRepo) FindTransferByTransferTo(context.Context, string) ([]*models.TransferByDestinationCardRow, error) {
	return []*models.TransferByDestinationCardRow{}, nil
}

type p0TransferObservability struct{}

func (p0TransferObservability) StartTracingAndLogging(ctx context.Context, _ string, _ ...attribute.KeyValue) (context.Context, trace.Span, func(string), string, func(string, ...zap.Field)) {
	ctx, span := trace.NewNoopTracerProvider().Tracer("transfer-p0").Start(ctx, "transfer-p0")
	return ctx, span, func(string) { span.End() }, "success", func(string, ...zap.Field) {}
}
func (p0TransferObservability) RecordMetrics(context.Context, string, string, time.Time) {}

func newP0TransferService(saldo *p0TransferSaldoRepo, command repository.TransferCommandRepository) service.TransferCommandService {
	return service.NewTransferCommandService(&service.TransferCommandDeps{
		CardRepository:            p0TransferCardRepo{},
		SaldoRepository:           saldo,
		TransferQueryRepository:   p0TransferQueryRepo{},
		TransferCommandRepository: command,
		Logger:                    &logger.Logger{Log: zap.NewNop()},
		Observability:             p0TransferObservability{},
	})
}

func TestCreateTransferAtomicFailureLeavesBalancesUntouched(t *testing.T) {
	balances := map[string]int32{"sender": 100000, "receiver": 50000}
	saldo := &p0TransferSaldoRepo{balances: balances}
	svc := newP0TransferService(saldo, &p0TransferCommandRepo{createErr: errors.New("create transfer failed")})

	result, err := svc.CreateTransaction(context.Background(), &requests.CreateTransferRequest{
		TransferFrom:   "sender",
		TransferTo:     "receiver",
		TransferAmount: 50000,
	})

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, int32(100000), balances["sender"])
	require.Equal(t, int32(50000), balances["receiver"])
}

func TestCreateTransferAtomicSettlesOnceWithoutStatusUpdate(t *testing.T) {
	balances := map[string]int32{"sender": 100000, "receiver": 50000}
	saldo := &p0TransferSaldoRepo{balances: balances}
	command := &p0TransferCommandRepo{
		createTransferID: 7,
		statusErr:        errors.New("status update must not be called"),
		balances:         balances,
	}
	svc := newP0TransferService(saldo, command)

	result, err := svc.CreateTransaction(context.Background(), &requests.CreateTransferRequest{
		TransferFrom:   "sender",
		TransferTo:     "receiver",
		TransferAmount: 50000,
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, int32(7), result.TransferID)
	require.Equal(t, "success", result.Status)
	require.Zero(t, command.statusCalls)
	require.Equal(t, int32(50000), balances["sender"])
	require.Equal(t, int32(100000), balances["receiver"])
}
