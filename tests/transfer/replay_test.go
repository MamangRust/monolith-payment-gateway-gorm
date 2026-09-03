package transfer_test

import (
	"context"
	"testing"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-transfer/repository"
	"github.com/MamangRust/monolith-payment-gateway-transfer/service"
	"github.com/stretchr/testify/require"
)

// replayTransferCommandRepo mirrors the repository contract for
// CreateTransferAtomic: the same idempotency key returns the ORIGINAL result
// (Replayed=true), never a second settlement.
type replayTransferCommandRepo struct {
	seq     int32
	created map[string]*repository.TransferAtomicResult
}

func (r *replayTransferCommandRepo) CreateTransferAtomic(_ context.Context, req *requests.CreateTransferRequest) (*repository.TransferAtomicResult, error) {
	if req.IdempotencyKey != "" {
		if existing, ok := r.created[req.IdempotencyKey]; ok {
			existing.Replayed = true
			return existing, nil
		}
	}
	res := &repository.TransferAtomicResult{
		Row: &models.TransferAllFieldsRow{
			TransferID:     r.seq + 1,
			TransferFrom:   req.TransferFrom,
			TransferTo:     req.TransferTo,
			TransferAmount: int32(req.TransferAmount),
			Status:         "success",
		},
	}
	r.seq++
	if req.IdempotencyKey != "" {
		r.created[req.IdempotencyKey] = res
	}
	return res, nil
}

func (replayTransferCommandRepo) CreateTransfer(context.Context, *requests.CreateTransferRequest) (*models.TransferAllFieldsRow, error) {
	return &models.TransferAllFieldsRow{}, nil
}
func (replayTransferCommandRepo) UpdateTransferStatus(context.Context, *requests.UpdateTransferStatus) (*models.TransferAllFieldsRow, error) {
	return &models.TransferAllFieldsRow{}, nil
}
func (replayTransferCommandRepo) UpdateTransferAtomic(context.Context, models.TransferAllFieldsRow) (*models.TransferAllFieldsRow, error) {
	return &models.TransferAllFieldsRow{}, nil
}
func (replayTransferCommandRepo) UpdateTransferSettlementDelta(context.Context, string, string, int, int) error {
	return nil
}
func (replayTransferCommandRepo) UpdateTransfer(context.Context, *requests.UpdateTransferRequest) (*models.TransferAllFieldsRow, error) {
	return &models.TransferAllFieldsRow{}, nil
}
func (replayTransferCommandRepo) UpdateTransferAmount(context.Context, *requests.UpdateTransferAmountRequest) (*models.TransferAllFieldsRow, error) {
	return &models.TransferAllFieldsRow{}, nil
}
func (replayTransferCommandRepo) TrashedTransfer(context.Context, int) (*models.Transfer, error) {
	return &models.Transfer{}, nil
}
func (replayTransferCommandRepo) RestoreTransfer(context.Context, int) (*models.Transfer, error) {
	return &models.Transfer{}, nil
}
func (replayTransferCommandRepo) DeleteTransferPermanent(context.Context, int) (bool, error) {
	return true, nil
}
func (replayTransferCommandRepo) RestoreAllTransfer(context.Context) (bool, error) {
	return true, nil
}
func (replayTransferCommandRepo) DeleteAllTransferPermanent(context.Context) (bool, error) {
	return true, nil
}

func newReplayTransferService(command repository.TransferCommandRepository) service.TransferCommandService {
	return newP0TransferService(&p0TransferSaldoRepo{}, command)
}

func TestCreateTransferReplayReturnsOriginalSettlement(t *testing.T) {
	command := &replayTransferCommandRepo{created: map[string]*repository.TransferAtomicResult{}}
	svc := newReplayTransferService(command)

	req := &requests.CreateTransferRequest{
		TransferFrom:   "sender",
		TransferTo:     "receiver",
		TransferAmount: 50000,
		IdempotencyKey: "transfer:create:user:7:op-1",
	}

	first, err := svc.CreateTransaction(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, first)
	require.Equal(t, int32(1), first.TransferID)
	require.Equal(t, "success", first.Status)

	second, err := svc.CreateTransaction(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, second)
	require.Equal(t, first.TransferID, second.TransferID, "replay must return the original settlement id")
	require.Equal(t, first.TransferAmount, second.TransferAmount)
	require.Len(t, command.created, 1, "the key must be persisted exactly once")
}

func TestCreateTransferWithoutKeySettlesTwice(t *testing.T) {
	command := &replayTransferCommandRepo{created: map[string]*repository.TransferAtomicResult{}}
	svc := newReplayTransferService(command)

	req := &requests.CreateTransferRequest{
		TransferFrom:   "sender",
		TransferTo:     "receiver",
		TransferAmount: 50000,
	}

	first, err := svc.CreateTransaction(context.Background(), req)
	require.NoError(t, err)
	second, err := svc.CreateTransaction(context.Background(), req)
	require.NoError(t, err)

	require.NotEqual(t, first.TransferID, second.TransferID, "without a key every request creates a new settlement")
}
