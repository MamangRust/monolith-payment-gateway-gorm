package transaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-transaction/repository"
	"github.com/MamangRust/monolith-payment-gateway-transaction/service"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type replayTransactionMerchantRepo struct{}

func (replayTransactionMerchantRepo) FindByApiKey(context.Context, string) (*models.MerchantAllFieldsRow, error) {
	return &models.MerchantAllFieldsRow{MerchantID: 1, UserID: 7, ApiKey: "key"}, nil
}

type replayTransactionCardRepo struct{}

func (replayTransactionCardRepo) FindUserCardByCardNumber(context.Context, string) (*models.CardByEmailRow, error) {
	return &models.CardByEmailRow{CardNumber: "4111111111111111", Email: "user@example.com"}, nil
}

func (replayTransactionCardRepo) FindCardByUserId(context.Context, int) (*models.CardAllFieldsRow, error) {
	return &models.CardAllFieldsRow{CardNumber: "5111111111111111"}, nil
}

func (replayTransactionCardRepo) FindCardByCardNumber(context.Context, string) (*models.CardAllFieldsRow, error) {
	return &models.CardAllFieldsRow{}, nil
}

func (replayTransactionCardRepo) UpdateCard(context.Context, *requests.UpdateCardRequest) (*models.CardUpdateRow, error) {
	return &models.CardUpdateRow{}, nil
}

func (replayTransactionCardRepo) UpdateCardOutstandingBalance(context.Context, int, int) (*models.UpdateOutstandingBalanceRow, error) {
	return &models.UpdateOutstandingBalanceRow{}, nil
}

func (replayTransactionCardRepo) AddRewardPoints(context.Context, int, int) (*models.AddRewardPointsRow, error) {
	return &models.AddRewardPointsRow{}, nil
}

type replayTransactionCommandRepo struct {
	seq     int32
	created map[string]*models.TransactionAllFieldsRow
}

func (replayTransactionCommandRepo) AuthorizeTransactionAtomic(context.Context, models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error) {
	return &models.TransactionAllFieldsRow{}, nil
}

func (replayTransactionCommandRepo) AuthorizeTransactionAtomicIdempotent(context.Context, models.TransactionAllFieldsRow) (*repository.AuthorizeTransactionResult, error) {
	return &repository.AuthorizeTransactionResult{Row: &models.TransactionAllFieldsRow{}}, nil
}

func (r *replayTransactionCommandRepo) CreateTransactionAtomicIdempotent(_ context.Context, params models.TransactionAllFieldsRow) (*repository.TransactionAtomicResult, error) {
	if params.IdempotencyKey != "" {
		if existing, ok := r.created[params.IdempotencyKey]; ok {
			return &repository.TransactionAtomicResult{Row: existing, Replayed: true}, nil
		}
	}
	res := &models.TransactionAllFieldsRow{
		TransactionID:   r.seq + 1,
		CardNumber:      params.CardNumber,
		Amount:          params.Amount,
		PaymentMethod:   params.PaymentMethod,
		MerchantID:      params.MerchantID,
		TransactionTime: params.TransactionTime,
		Status:          "success",
	}
	r.seq++
	if params.IdempotencyKey != "" {
		r.created[params.IdempotencyKey] = res
	}
	return &repository.TransactionAtomicResult{Row: res}, nil
}

func (replayTransactionCommandRepo) CreateTransactionAtomic(context.Context, models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error) {
	return &models.TransactionAllFieldsRow{}, nil
}

func (replayTransactionCommandRepo) UpdateSaldoBalanceDelta(context.Context, string, int) error {
	return nil
}

func (replayTransactionCommandRepo) CreateTransaction(context.Context, *requests.CreateTransactionRequest) (*models.TransactionAllFieldsRow, error) {
	return &models.TransactionAllFieldsRow{}, nil
}

func (replayTransactionCommandRepo) UpdateTransaction(context.Context, *requests.UpdateTransactionRequest) (*models.TransactionAllFieldsRow, error) {
	return &models.TransactionAllFieldsRow{}, nil
}

func (replayTransactionCommandRepo) UpdateTransactionStatus(context.Context, *requests.UpdateTransactionStatus) (*models.TransactionAllFieldsRow, error) {
	return &models.TransactionAllFieldsRow{}, nil
}

func (replayTransactionCommandRepo) UpdateTransactionAtomic(context.Context, models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error) {
	return &models.TransactionAllFieldsRow{}, nil
}

func (replayTransactionCommandRepo) CaptureTransactionAtomic(context.Context, models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error) {
	return &models.TransactionAllFieldsRow{}, nil
}

func (replayTransactionCommandRepo) VoidTransactionAtomic(context.Context, int32) (*models.TransactionAllFieldsRow, error) {
	return &models.TransactionAllFieldsRow{}, nil
}

func (replayTransactionCommandRepo) RefundTransactionAtomic(context.Context, models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error) {
	return &models.TransactionAllFieldsRow{}, nil
}

func (replayTransactionCommandRepo) TrashedTransaction(context.Context, int) (*models.TransactionTrashRestoreRow, error) {
	return &models.TransactionTrashRestoreRow{}, nil
}

func (replayTransactionCommandRepo) RestoreTransaction(context.Context, int) (*models.TransactionTrashRestoreRow, error) {
	return &models.TransactionTrashRestoreRow{}, nil
}

func (replayTransactionCommandRepo) DeleteTransactionPermanent(context.Context, int) (bool, error) {
	return true, nil
}

func (replayTransactionCommandRepo) RestoreAllTransaction(context.Context) (bool, error) {
	return true, nil
}

func (replayTransactionCommandRepo) DeleteAllTransactionPermanent(context.Context) (bool, error) {
	return true, nil
}

type replayTransactionQueryRepo struct{}

func (replayTransactionQueryRepo) FindAllTransactions(context.Context, *requests.FindAllTransactions) ([]*models.TransactionListRow, error) {
	return []*models.TransactionListRow{}, nil
}

func (replayTransactionQueryRepo) FindByActive(context.Context, *requests.FindAllTransactions) ([]*models.TransactionListRow, error) {
	return []*models.TransactionListRow{}, nil
}

func (replayTransactionQueryRepo) FindByTrashed(context.Context, *requests.FindAllTransactions) ([]*models.TransactionListWithDeletedRow, error) {
	return []*models.TransactionListWithDeletedRow{}, nil
}

func (replayTransactionQueryRepo) FindAllTransactionByCardNumber(context.Context, *requests.FindAllTransactionCardNumber) ([]*models.TransactionListRow, error) {
	return []*models.TransactionListRow{}, nil
}

func (replayTransactionQueryRepo) FindById(context.Context, int) (*models.TransactionAllFieldsRow, error) {
	return &models.TransactionAllFieldsRow{}, nil
}

func (replayTransactionQueryRepo) FindTransactionByMerchantId(context.Context, int) ([]*models.TransactionListRow, error) {
	return []*models.TransactionListRow{}, nil
}

type replayTransactionSaldoRepo struct{}

func (replayTransactionSaldoRepo) FindByCardNumber(context.Context, string) (*models.Saldo, error) {
	return &models.Saldo{}, nil
}

func (replayTransactionSaldoRepo) UpdateSaldoBalance(context.Context, *requests.UpdateSaldoBalance) (*models.UpdateSaldoBalanceRow, error) {
	return &models.UpdateSaldoBalanceRow{}, nil
}

type replayTransactionObservability struct{}

func (replayTransactionObservability) StartTracingAndLogging(ctx context.Context, _ string, _ ...attribute.KeyValue) (context.Context, trace.Span, func(string), string, func(string, ...zap.Field)) {
	ctx, span := trace.NewNoopTracerProvider().Tracer("transaction-replay").Start(ctx, "transaction-replay")
	return ctx, span, func(string) { span.End() }, "success", func(string, ...zap.Field) {}
}
func (replayTransactionObservability) RecordMetrics(context.Context, string, string, time.Time) {}

func newReplayTransactionService(command *replayTransactionCommandRepo) service.TransactionCommandService {
	return service.NewTransactionCommandService(&service.TransactionCommandServiceDeps{
		MerchantRepository:           replayTransactionMerchantRepo{},
		CardRepository:               replayTransactionCardRepo{},
		SaldoRepository:              replayTransactionSaldoRepo{},
		TransactionQueryRepository:   replayTransactionQueryRepo{},
		TransactionCommandRepository: command,
		Logger:                       &logger.Logger{Log: zap.NewNop()},
		Observability:                replayTransactionObservability{},
	})
}

func TestCreateTransactionReplayReturnsOriginal(t *testing.T) {
	command := &replayTransactionCommandRepo{created: map[string]*models.TransactionAllFieldsRow{}}
	svc := newReplayTransactionService(command)

	req := &requests.CreateTransactionRequest{
		CardNumber:      "4111111111111111",
		Amount:          100000,
		PaymentMethod:   "gopay",
		TransactionTime: time.Now(),
		IdempotencyKey:  "transaction:create:merchant:1:op-1",
	}

	first, err := svc.Create(context.Background(), "api-key", req)
	require.NoError(t, err)
	require.NotNil(t, first)
	require.Equal(t, int32(1), first.TransactionID)
	require.Equal(t, "success", first.Status)

	second, err := svc.Create(context.Background(), "api-key", req)
	require.NoError(t, err)
	require.NotNil(t, second)
	require.Equal(t, first.TransactionID, second.TransactionID, "replay must return the original transaction id")
	require.Equal(t, first.Amount, second.Amount)
	require.Len(t, command.created, 1, "the key must be persisted exactly once")
}

func TestCreateTransactionWithoutKeyCreatesTwice(t *testing.T) {
	command := &replayTransactionCommandRepo{created: map[string]*models.TransactionAllFieldsRow{}}
	svc := newReplayTransactionService(command)

	req := &requests.CreateTransactionRequest{
		CardNumber:      "4111111111111111",
		Amount:          100000,
		PaymentMethod:   "gopay",
		TransactionTime: time.Now(),
	}

	first, err := svc.Create(context.Background(), "api-key", req)
	require.NoError(t, err)
	second, err := svc.Create(context.Background(), "api-key", req)
	require.NoError(t, err)

	require.NotEqual(t, first.TransactionID, second.TransactionID, "without a key every request creates a new transaction")
}
