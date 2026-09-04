package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
	"gorm.io/gorm"
)

var errInsufficientBalance = fmt.Errorf("insufficient balance")

type withdrawCommandRepository struct {
	db *gorm.DB
}

func NewWithdrawCommandRepository(db *gorm.DB, dailyLimit ...int64) WithdrawCommandRepository {
	return &withdrawCommandRepository{db: db}
}

type WithdrawAtomicResult struct {
	Row      *models.WithdrawAllFieldsRow
	Replayed bool
}

func (r *withdrawCommandRepository) CreateWithdraw(ctx context.Context, request *requests.CreateWithdrawRequest) (*models.WithdrawAllFieldsRow, error) {
	var result models.WithdrawAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`INSERT INTO withdraws (card_number, withdraw_amount, withdraw_time, idempotency_key, status, created_at, updated_at) VALUES (?, ?, ?, ?, 'pending', current_timestamp, current_timestamp) RETURNING withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at`,
		request.CardNumber, int32(request.WithdrawAmount), request.WithdrawTime, request.IdempotencyKey).Scan(&result).Error
	if err != nil {
		return nil, withdraw_errors.ErrCreateWithdrawFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *withdrawCommandRepository) CreateWithdrawAtomic(ctx context.Context, request *requests.CreateWithdrawRequest) (*WithdrawAtomicResult, error) {
	// Fast path: if idempotency key is set, check for existing record first
	if request.IdempotencyKey != "" {
		var existing models.WithdrawAllFieldsRow
		err := r.db.WithContext(ctx).Raw(`SELECT withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at FROM withdraws WHERE idempotency_key = ? AND deleted_at IS NULL`, request.IdempotencyKey).Scan(&existing).Error
		if err == nil && existing.WithdrawID > 0 {
			return &WithdrawAtomicResult{Row: &existing, Replayed: true}, nil
		}
	}
	var result models.WithdrawAllFieldsRow
	cn := request.CardNumber
	amt := int32(request.WithdrawAmount)
	wt := request.WithdrawTime
	ik := request.IdempotencyKey

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Atomic check-and-debit
		res := tx.Exec(`UPDATE saldos SET total_balance = total_balance - ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL AND total_balance >= ? AND ? > 0`,
			amt, cn, amt, amt)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errInsufficientBalance
		}
		return tx.Raw(`INSERT INTO withdraws (card_number, withdraw_amount, withdraw_time, idempotency_key, status, created_at, updated_at) VALUES (?, ?, ?, ?, 'success', current_timestamp, current_timestamp) RETURNING withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at`, cn, amt, wt, ik).Scan(&result).Error
	})
	if err != nil {
		return nil, withdraw_errors.ErrCreateWithdrawFailed.WithInternal(err)
	}
	return &WithdrawAtomicResult{Row: &result}, nil
}

func (r *withdrawCommandRepository) UpdateWithdraw(ctx context.Context, request *requests.UpdateWithdrawRequest) (*models.WithdrawAllFieldsRow, error) {
	var result models.WithdrawAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE withdraws SET card_number = ?, withdraw_amount = ?, withdraw_time = ?, updated_at = current_timestamp WHERE withdraw_id = ? AND deleted_at IS NULL RETURNING withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at`,
		request.CardNumber, int32(request.WithdrawAmount), request.WithdrawTime, *request.WithdrawID).Scan(&result).Error
	if err != nil {
		return nil, withdraw_errors.ErrUpdateWithdrawFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *withdrawCommandRepository) UpdateWithdrawAtomic(ctx context.Context, request *requests.UpdateWithdrawRequest) (*models.WithdrawAllFieldsRow, error) {
	var result models.WithdrawAllFieldsRow
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the old withdraw to compute delta atomically
		var oldAmount int32
		if err := tx.Raw(`SELECT withdraw_amount FROM withdraws WHERE withdraw_id = ? AND deleted_at IS NULL FOR UPDATE`, *request.WithdrawID).Scan(&oldAmount).Error; err != nil {
			return err
		}
		// Adjust saldo by delta (debit increases, credit decreases)
		delta := int32(request.WithdrawAmount) - oldAmount
		if err := tx.Exec(`UPDATE saldos SET total_balance = total_balance - ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, delta, request.CardNumber).Error; err != nil {
			return err
		}
		return tx.Raw(`UPDATE withdraws SET card_number = ?, withdraw_amount = ?, status = 'success', updated_at = current_timestamp WHERE withdraw_id = ? AND deleted_at IS NULL RETURNING withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at`,
			request.CardNumber, int32(request.WithdrawAmount), *request.WithdrawID).Scan(&result).Error
	})
	if err != nil {
		return nil, withdraw_errors.ErrUpdateWithdrawFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *withdrawCommandRepository) UpdateSaldoBalanceDelta(ctx context.Context, cardNumber string, delta int) error {
	return r.db.WithContext(ctx).Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, delta, cardNumber).Error
}

func (r *withdrawCommandRepository) UpdateWithdrawStatus(ctx context.Context, request *requests.UpdateWithdrawStatus) (*models.WithdrawAllFieldsRow, error) {
	var result models.WithdrawAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE withdraws SET status = ?, updated_at = current_timestamp WHERE withdraw_id = ? AND deleted_at IS NULL RETURNING withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at`,
		request.Status, request.WithdrawID).Scan(&result).Error
	if err != nil {
		return nil, withdraw_errors.ErrUpdateWithdrawStatusFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *withdrawCommandRepository) TrashedWithdraw(ctx context.Context, withdrawID int) (*models.Withdraw, error) {
	var result models.Withdraw
	err := r.db.WithContext(ctx).Raw(`UPDATE withdraws SET deleted_at = current_timestamp WHERE withdraw_id = ? AND deleted_at IS NULL RETURNING withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at, deleted_at`, withdrawID).Scan(&result).Error
	if err != nil {
		return nil, withdraw_errors.ErrTrashedWithdrawFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *withdrawCommandRepository) RestoreWithdraw(ctx context.Context, withdrawID int) (*models.Withdraw, error) {
	var result models.Withdraw
	err := r.db.WithContext(ctx).Raw(`UPDATE withdraws SET deleted_at = NULL WHERE withdraw_id = ? AND deleted_at IS NOT NULL RETURNING withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at, deleted_at`, withdrawID).Scan(&result).Error
	if err != nil {
		return nil, withdraw_errors.ErrRestoreWithdrawFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *withdrawCommandRepository) DeleteWithdrawPermanent(ctx context.Context, withdrawID int) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`DELETE FROM withdraws WHERE withdraw_id = ? AND deleted_at IS NOT NULL`, withdrawID).Error
	if err != nil {
		return false, withdraw_errors.ErrDeleteWithdrawPermanentFailed.WithInternal(err)
	}
	return true, nil
}

func (r *withdrawCommandRepository) RestoreAllWithdraw(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`UPDATE withdraws SET deleted_at = NULL WHERE deleted_at IS NOT NULL`).Error
	if err != nil {
		return false, withdraw_errors.ErrRestoreAllWithdrawsFailed.WithInternal(err)
	}
	return true, nil
}

func (r *withdrawCommandRepository) DeleteAllWithdrawPermanent(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`DELETE FROM withdraws WHERE deleted_at IS NOT NULL`).Error
	if err != nil {
		return false, withdraw_errors.ErrDeleteAllWithdrawsPermanentFailed.WithInternal(err)
	}
	return true, nil
}

var _ = time.Now
