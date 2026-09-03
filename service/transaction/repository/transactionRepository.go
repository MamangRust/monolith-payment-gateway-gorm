package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
	"gorm.io/gorm"
)

var errInsufficientBalance = fmt.Errorf("insufficient balance")

type transactionCommandRepository struct {
	db *gorm.DB
}

func NewTransactionCommandRepository(db *gorm.DB) TransactionCommandRepository {
	return &transactionCommandRepository{db: db}
}

type AuthorizeTransactionResult struct {
	Row      *models.TransactionAllFieldsRow
	Replayed bool
}

type TransactionAtomicResult struct {
	Row      *models.TransactionAllFieldsRow
	Replayed bool
}

func (r *transactionCommandRepository) CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*models.TransactionAllFieldsRow, error) {
	var result models.TransactionAllFieldsRow
	mid := 0
	if request.MerchantID != nil {
		mid = *request.MerchantID
	}
	err := r.db.WithContext(ctx).Raw(`INSERT INTO transactions (card_number, amount, payment_method, merchant_id, transaction_time, idempotency_key, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, 'pending', current_timestamp, current_timestamp) RETURNING transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, idempotency_key, created_at, updated_at`,
		request.CardNumber, int32(request.Amount), request.PaymentMethod, mid, request.TransactionTime, request.IdempotencyKey).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrCreateTransactionFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transactionCommandRepository) UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*models.TransactionAllFieldsRow, error) {
	var result models.TransactionAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE transactions SET card_number = ?, amount = ?, payment_method = ?, merchant_id = ?, transaction_time = ?, updated_at = current_timestamp WHERE transaction_id = ? AND deleted_at IS NULL RETURNING transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, idempotency_key, created_at, updated_at`,
		request.CardNumber, int32(request.Amount), request.PaymentMethod, request.MerchantID, request.TransactionTime, *request.TransactionID).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transactionCommandRepository) UpdateTransactionStatus(ctx context.Context, request *requests.UpdateTransactionStatus) (*models.TransactionAllFieldsRow, error) {
	var result models.TransactionAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE transactions SET status = ?, updated_at = current_timestamp WHERE transaction_id = ? AND deleted_at IS NULL RETURNING transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, idempotency_key, created_at, updated_at`,
		request.Status, request.TransactionID).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionStatusFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transactionCommandRepository) TrashedTransaction(ctx context.Context, transaction_id int) (*models.TransactionTrashRestoreRow, error) {
	var result models.TransactionTrashRestoreRow
	err := r.db.WithContext(ctx).Raw(`UPDATE transactions SET deleted_at = current_timestamp WHERE transaction_id = ? AND deleted_at IS NULL RETURNING transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, created_at, updated_at, deleted_at`, transaction_id).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrTrashedTransactionFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transactionCommandRepository) RestoreTransaction(ctx context.Context, transaction_id int) (*models.TransactionTrashRestoreRow, error) {
	var result models.TransactionTrashRestoreRow
	err := r.db.WithContext(ctx).Raw(`UPDATE transactions SET deleted_at = NULL WHERE transaction_id = ? AND deleted_at IS NOT NULL RETURNING transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, created_at, updated_at, deleted_at`, transaction_id).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrRestoreTransactionFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transactionCommandRepository) DeleteTransactionPermanent(ctx context.Context, transaction_id int) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`DELETE FROM transactions WHERE transaction_id = ? AND deleted_at IS NOT NULL`, transaction_id).Error
	if err != nil {
		return false, transaction_errors.ErrDeleteTransactionPermanentFailed.WithInternal(err)
	}
	return true, nil
}

func (r *transactionCommandRepository) RestoreAllTransaction(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`UPDATE transactions SET deleted_at = NULL WHERE deleted_at IS NOT NULL`).Error
	if err != nil {
		return false, transaction_errors.ErrRestoreAllTransactionsFailed.WithInternal(err)
	}
	return true, nil
}

func (r *transactionCommandRepository) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`DELETE FROM transactions WHERE deleted_at IS NOT NULL`).Error
	if err != nil {
		return false, transaction_errors.ErrDeleteAllTransactionsPermanentFailed.WithInternal(err)
	}
	return true, nil
}

func (r *transactionCommandRepository) UpdateSaldoBalanceDelta(ctx context.Context, cardNumber string, delta int) error {
	return r.db.WithContext(ctx).Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, delta, cardNumber).Error
}

func (r *transactionCommandRepository) CreateTransactionAtomic(ctx context.Context, params models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error) {
	var result models.TransactionAllFieldsRow
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Atomic check-and-debit: single UPDATE with WHERE balance >= amount
		res := tx.Exec(`UPDATE saldos SET total_balance = total_balance - ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL AND total_balance >= ? AND ? > 0`,
			params.Amount, params.CardNumber, params.Amount, params.Amount)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errInsufficientBalance
		}
		if err := tx.Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, params.Amount, params.CardNumber2).Error; err != nil {
			return err
		}
		return tx.Raw(`INSERT INTO transactions (card_number, amount, payment_method, merchant_id, transaction_time, idempotency_key, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, 'success', current_timestamp, current_timestamp) RETURNING transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, idempotency_key, created_at, updated_at`,
			params.CardNumber, params.Amount, params.PaymentMethod, params.MerchantID, params.TransactionTime, params.IdempotencyKey).Scan(&result).Error
	})
	if err != nil {
		return nil, transaction_errors.ErrCreateTransactionFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transactionCommandRepository) AuthorizeTransactionAtomic(ctx context.Context, params models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error) {
	var result models.TransactionAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`INSERT INTO transactions (card_number, amount, payment_method, merchant_id, transaction_time, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, 'authorized', current_timestamp, current_timestamp) RETURNING transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, idempotency_key, created_at, updated_at`,
		params.CardNumber, params.Amount, params.PaymentMethod, params.MerchantID, params.TransactionTime).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrCreateTransactionFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transactionCommandRepository) AuthorizeTransactionAtomicIdempotent(ctx context.Context, params models.TransactionAllFieldsRow) (*AuthorizeTransactionResult, error) {
	// Check for existing by idempotency key (skip if empty)
	if params.IdempotencyKey != "" {
		var existing models.TransactionAllFieldsRow
		err := r.db.WithContext(ctx).Raw(`SELECT transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, idempotency_key, created_at, updated_at FROM transactions WHERE idempotency_key = ? AND deleted_at IS NULL`, params.IdempotencyKey).Scan(&existing).Error
		if err == nil && existing.TransactionID > 0 {
			return &AuthorizeTransactionResult{Row: &existing, Replayed: true}, nil
		}
	}
	result, err := r.AuthorizeTransactionAtomic(ctx, params)
	if err != nil {
		return nil, err
	}
	return &AuthorizeTransactionResult{Row: result}, nil
}

func (r *transactionCommandRepository) CreateTransactionAtomicIdempotent(ctx context.Context, params models.TransactionAllFieldsRow) (*TransactionAtomicResult, error) {
	// Check for existing by idempotency key (skip if empty)
	if params.IdempotencyKey != "" {
		var existing models.TransactionAllFieldsRow
		err := r.db.WithContext(ctx).Raw(`SELECT transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, idempotency_key, created_at, updated_at FROM transactions WHERE idempotency_key = ? AND deleted_at IS NULL`, params.IdempotencyKey).Scan(&existing).Error
		if err == nil && existing.TransactionID > 0 {
			return &TransactionAtomicResult{Row: &existing, Replayed: true}, nil
		}
	}
	result, err := r.CreateTransactionAtomic(ctx, params)
	if err != nil {
		return nil, err
	}
	return &TransactionAtomicResult{Row: result}, nil
}

func (r *transactionCommandRepository) UpdateTransactionAtomic(ctx context.Context, params models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error) {
	var result models.TransactionAllFieldsRow
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the old transaction to compute delta atomically
		var oldAmount int32
		if err := tx.Raw(`SELECT amount FROM transactions WHERE transaction_id = ? AND deleted_at IS NULL FOR UPDATE`, params.TransactionID).Scan(&oldAmount).Error; err != nil {
			return err
		}
		delta := params.Amount - oldAmount
		// Debit user (card_number) by delta
		if err := tx.Exec(`UPDATE saldos SET total_balance = total_balance - ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, delta, params.CardNumber).Error; err != nil {
			return err
		}
		// Credit merchant by delta
		if err := tx.Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, delta, params.MerchantCardNumber).Error; err != nil {
			return err
		}
		return tx.Raw(`UPDATE transactions SET status = 'success', amount = ?, payment_method = ?, merchant_id = ?, updated_at = current_timestamp WHERE transaction_id = ? AND deleted_at IS NULL RETURNING transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, idempotency_key, created_at, updated_at`,
			params.Amount, params.PaymentMethod, params.MerchantID, params.TransactionID).Scan(&result).Error
	})
	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transactionCommandRepository) CaptureTransactionAtomic(ctx context.Context, params models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error) {
	var result models.TransactionAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE transactions SET status = 'captured', updated_at = current_timestamp WHERE transaction_id = ? AND deleted_at IS NULL RETURNING transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, idempotency_key, created_at, updated_at`,
		params.TransactionID).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionStatusFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transactionCommandRepository) VoidTransactionAtomic(ctx context.Context, transactionID int32) (*models.TransactionAllFieldsRow, error) {
	var result models.TransactionAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE transactions SET status = 'voided', updated_at = current_timestamp WHERE transaction_id = ? AND deleted_at IS NULL RETURNING transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, idempotency_key, created_at, updated_at`,
		transactionID).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionStatusFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *transactionCommandRepository) RefundTransactionAtomic(ctx context.Context, params models.TransactionAllFieldsRow) (*models.TransactionAllFieldsRow, error) {
	var result models.TransactionAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE transactions SET status = 'refunded', updated_at = current_timestamp WHERE transaction_id = ? AND deleted_at IS NULL RETURNING transaction_id, transaction_no, card_number, amount, payment_method, merchant_id, transaction_time, status, idempotency_key, created_at, updated_at`,
		params.TransactionID).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrUpdateTransactionStatusFailed.WithInternal(err)
	}
	return &result, nil
}

var _ = time.Now
