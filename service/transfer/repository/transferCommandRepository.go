package repository

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
	"gorm.io/gorm"
)

var errInsufficientBalance = fmt.Errorf("insufficient balance")

type transferCommandRepository struct {
	db *gorm.DB
}

func NewTransferCommandRepository(db *gorm.DB, dailyLimit ...int64) TransferCommandRepository {
	return &transferCommandRepository{db: db}
}

type TransferAtomicResult struct {
	Row      *models.TransferAllFieldsRow
	Replayed bool
}

func (r *transferCommandRepository) CreateTransferAtomic(ctx context.Context, request *requests.CreateTransferRequest) (*TransferAtomicResult, error) {
	// Fast path: if idempotency key is set, check for existing record first
	if request.IdempotencyKey != "" {
	var existing models.TransferAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`SELECT transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at FROM transfers WHERE idempotency_key = ? AND deleted_at IS NULL`, request.IdempotencyKey).Scan(&existing).Error
	if err == nil && existing.TransferID > 0 {
		return &TransferAtomicResult{Row: &existing, Replayed: true}, nil
	}
	}
	var result models.TransferAllFieldsRow
	amount := int32(request.TransferAmount)
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if request.TransferFrom == request.TransferTo {
			return errInsufficientBalance
		}
		// Atomic check-and-debit sender
		res := tx.Exec(`UPDATE saldos SET total_balance = total_balance - ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL AND total_balance >= ? AND ? > 0`,
			amount, request.TransferFrom, amount, amount)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errInsufficientBalance
		}
		// Credit receiver
		if err := tx.Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, amount, request.TransferTo).Error; err != nil {
			return err
		}
		// Insert transfer
		return tx.Raw(`INSERT INTO transfers (transfer_from, transfer_to, transfer_amount, transfer_time, idempotency_key, status, created_at, updated_at) VALUES (?, ?, ?, current_timestamp, ?, 'success', current_timestamp, current_timestamp) RETURNING transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at`,
			request.TransferFrom, request.TransferTo, amount, request.IdempotencyKey).Scan(&result).Error
	})
	if err != nil {
		return nil, transfer_errors.ErrCreateTransferFailed.WithInternal(err)
	}
	return &TransferAtomicResult{Row: &result}, nil
}

func (r *transferCommandRepository) UpdateTransferSettlementDelta(ctx context.Context, senderCard, receiverCard string, senderDelta, receiverDelta int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, senderDelta, senderCard).Error; err != nil {
			return err
		}
		return tx.Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, receiverDelta, receiverCard).Error
	})
}

func (r *transferCommandRepository) UpdateTransferAtomic(ctx context.Context, params models.TransferAllFieldsRow) (*models.TransferAllFieldsRow, error) {
	var result models.TransferAllFieldsRow
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the old transfer to compute delta atomically
		var oldAmount int32
		if err := tx.Raw(`SELECT transfer_amount FROM transfers WHERE transfer_id = ? AND deleted_at IS NULL FOR UPDATE`, params.TransferID).Scan(&oldAmount).Error; err != nil {
			return err
		}
		delta := params.TransferAmount - oldAmount
		// Debit sender by delta
		if err := tx.Exec(`UPDATE saldos SET total_balance = total_balance - ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, delta, params.TransferFrom).Error; err != nil {
			return err
		}
		// Credit receiver by delta
		if err := tx.Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, delta, params.TransferTo).Error; err != nil {
			return err
		}
		return tx.Raw(`UPDATE transfers SET transfer_from = ?, transfer_to = ?, transfer_amount = ?, updated_at = current_timestamp WHERE transfer_id = ? AND deleted_at IS NULL RETURNING transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at`,
			params.TransferFrom, params.TransferTo, params.TransferAmount, params.TransferID).Scan(&result).Error
	})
	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transferCommandRepository) UpdateTransfer(ctx context.Context, request *requests.UpdateTransferRequest) (*models.TransferAllFieldsRow, error) {
	var result models.TransferAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE transfers SET transfer_from = ?, transfer_to = ?, transfer_amount = ?, updated_at = current_timestamp WHERE transfer_id = ? AND deleted_at IS NULL RETURNING transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at`,
		request.TransferFrom, request.TransferTo, int32(request.TransferAmount), *request.TransferID).Scan(&result).Error
	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transferCommandRepository) UpdateTransferAmount(ctx context.Context, request *requests.UpdateTransferAmountRequest) (*models.TransferAllFieldsRow, error) {
	var result models.TransferAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE transfers SET transfer_amount = ?, updated_at = current_timestamp WHERE transfer_id = ? AND deleted_at IS NULL RETURNING transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at`,
		int32(request.TransferAmount), request.TransferID).Scan(&result).Error
	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferAmountFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transferCommandRepository) UpdateTransferStatus(ctx context.Context, request *requests.UpdateTransferStatus) (*models.TransferAllFieldsRow, error) {
	var result models.TransferAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE transfers SET status = ?, updated_at = current_timestamp WHERE transfer_id = ? AND deleted_at IS NULL RETURNING transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at`,
		request.Status, request.TransferID).Scan(&result).Error
	if err != nil {
		return nil, transfer_errors.ErrUpdateTransferStatusFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transferCommandRepository) TrashedTransfer(ctx context.Context, transferID int) (*models.Transfer, error) {
	var result models.Transfer
	err := r.db.WithContext(ctx).Raw(`UPDATE transfers SET deleted_at = current_timestamp WHERE transfer_id = ? AND deleted_at IS NULL RETURNING transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at, deleted_at`, transferID).Scan(&result).Error
	if err != nil {
		return nil, transfer_errors.ErrTrashedTransferFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transferCommandRepository) RestoreTransfer(ctx context.Context, transferID int) (*models.Transfer, error) {
	var result models.Transfer
	err := r.db.WithContext(ctx).Raw(`UPDATE transfers SET deleted_at = NULL WHERE transfer_id = ? AND deleted_at IS NOT NULL RETURNING transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at, deleted_at`, transferID).Scan(&result).Error
	if err != nil {
		return nil, transfer_errors.ErrRestoreTransferFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transferCommandRepository) DeleteTransferPermanent(ctx context.Context, transferID int) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`DELETE FROM transfers WHERE transfer_id = ? AND deleted_at IS NOT NULL`, transferID).Error
	if err != nil {
		return false, transfer_errors.ErrDeleteTransferPermanentFailed.WithInternal(err)
	}
	return true, nil
}

func (r *transferCommandRepository) RestoreAllTransfer(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`UPDATE transfers SET deleted_at = NULL WHERE deleted_at IS NOT NULL`).Error
	if err != nil {
		return false, transfer_errors.ErrRestoreAllTransfersFailed.WithInternal(err)
	}
	return true, nil
}

func (r *transferCommandRepository) DeleteAllTransferPermanent(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`DELETE FROM transfers WHERE deleted_at IS NOT NULL`).Error
	if err != nil {
		return false, transfer_errors.ErrDeleteAllTransfersPermanentFailed.WithInternal(err)
	}
	return true, nil
}
