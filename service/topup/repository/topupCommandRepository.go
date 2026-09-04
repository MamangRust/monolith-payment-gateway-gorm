package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
	"gorm.io/gorm"
)

type topupCommandRepository struct {
	db *gorm.DB
}

func NewTopupCommandRepository(db *gorm.DB) TopupCommandRepository {
	return &topupCommandRepository{db: db}
}

type TopupAtomicResult struct {
	Row      *models.TopupAllFieldsRow
	Replayed bool
}

func (r *topupCommandRepository) CreateTopup(ctx context.Context, request *requests.CreateTopupRequest) (*models.TopupAllFieldsRow, error) {
	var result models.TopupAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`INSERT INTO topups (card_number, topup_amount, topup_method, topup_time, idempotency_key, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, 'pending', current_timestamp, current_timestamp) RETURNING topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at`,
		request.CardNumber, int32(request.TopupAmount), request.TopupMethod, time.Now(), request.IdempotencyKey).Scan(&result).Error
	if err != nil {
		return nil, topup_errors.ErrCreateTopupFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *topupCommandRepository) CreateTopupAtomic(ctx context.Context, request *requests.CreateTopupRequest) (*TopupAtomicResult, error) {
	// Fast path: if idempotency key is set, check for existing record first
	if request.IdempotencyKey != "" {
		var existing models.TopupAllFieldsRow
		err := r.db.WithContext(ctx).Raw(`SELECT topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at FROM topups WHERE idempotency_key = ? AND deleted_at IS NULL`, request.IdempotencyKey).Scan(&existing).Error
		if err == nil && existing.TopupID > 0 {
			return &TopupAtomicResult{Row: &existing, Replayed: true}, nil
		}
	}

	var result models.TopupAllFieldsRow
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, int32(request.TopupAmount), request.CardNumber).Error; err != nil {
			return err
		}
		return tx.Raw(`INSERT INTO topups (card_number, topup_amount, topup_method, topup_time, idempotency_key, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, 'success', current_timestamp, current_timestamp) RETURNING topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at`,
			request.CardNumber, int32(request.TopupAmount), request.TopupMethod, time.Now(), request.IdempotencyKey).Scan(&result).Error
	})
	if err != nil {
		return nil, topup_errors.ErrCreateTopupFailed.WithInternal(err)
	}
	return &TopupAtomicResult{Row: &result}, nil
}

func (r *topupCommandRepository) UpdateTopup(ctx context.Context, request *requests.UpdateTopupRequest) (*models.TopupAllFieldsRow, error) {
	var result models.TopupAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE topups SET card_number = ?, topup_amount = ?, topup_method = ?, updated_at = current_timestamp WHERE topup_id = ? AND deleted_at IS NULL RETURNING topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at`,
		request.CardNumber, int32(request.TopupAmount), request.TopupMethod, *request.TopupID).Scan(&result).Error
	if err != nil {
		return nil, topup_errors.ErrUpdateTopupFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *topupCommandRepository) UpdateTopupAtomic(ctx context.Context, request *requests.UpdateTopupRequest) (*models.TopupAllFieldsRow, error) {
	var result models.TopupAllFieldsRow
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the old topup to compute delta atomically
		var oldAmount int32
		if err := tx.Raw(`SELECT topup_amount FROM topups WHERE topup_id = ? AND deleted_at IS NULL FOR UPDATE`, *request.TopupID).Scan(&oldAmount).Error; err != nil {
			return err
		}
		delta := int32(request.TopupAmount) - oldAmount
		// Adjust saldo by delta
		if err := tx.Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, delta, request.CardNumber).Error; err != nil {
			return err
		}
		return tx.Raw(`UPDATE topups SET card_number = ?, topup_amount = ?, topup_method = ?, updated_at = current_timestamp WHERE topup_id = ? AND deleted_at IS NULL RETURNING topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at`,
			request.CardNumber, int32(request.TopupAmount), request.TopupMethod, *request.TopupID).Scan(&result).Error
	})
	if err != nil {
		return nil, topup_errors.ErrUpdateTopupFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *topupCommandRepository) UpdateSaldoBalanceDelta(ctx context.Context, cardNumber string, delta int) error {
	return r.db.WithContext(ctx).Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = current_timestamp WHERE card_number = ? AND deleted_at IS NULL`, delta, cardNumber).Error
}

func (r *topupCommandRepository) UpdateTopupAmount(ctx context.Context, request *requests.UpdateTopupAmount) (*models.TopupAllFieldsRow, error) {
	var result models.TopupAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE topups SET topup_amount = ?, updated_at = current_timestamp WHERE topup_id = ? AND deleted_at IS NULL RETURNING topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at`,
		int32(request.TopupAmount), request.TopupID).Scan(&result).Error
	if err != nil {
		return nil, topup_errors.ErrUpdateTopupAmountFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *topupCommandRepository) UpdateTopupStatus(ctx context.Context, request *requests.UpdateTopupStatus) (*models.TopupAllFieldsRow, error) {
	var result models.TopupAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE topups SET status = ?, updated_at = current_timestamp WHERE topup_id = ? AND deleted_at IS NULL RETURNING topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at`,
		request.Status, request.TopupID).Scan(&result).Error
	if err != nil {
		return nil, topup_errors.ErrUpdateTopupStatusFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *topupCommandRepository) TrashedTopup(ctx context.Context, topup_id int) (*models.Topup, error) {
	var result models.Topup
	err := r.db.WithContext(ctx).Raw(`UPDATE topups SET deleted_at = current_timestamp WHERE topup_id = ? AND deleted_at IS NULL RETURNING topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at, deleted_at`, topup_id).Scan(&result).Error
	if err != nil {
		return nil, topup_errors.ErrTrashedTopupFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *topupCommandRepository) RestoreTopup(ctx context.Context, topup_id int) (*models.Topup, error) {
	var result models.Topup
	err := r.db.WithContext(ctx).Raw(`UPDATE topups SET deleted_at = NULL WHERE topup_id = ? AND deleted_at IS NOT NULL RETURNING topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at, deleted_at`, topup_id).Scan(&result).Error
	if err != nil {
		return nil, topup_errors.ErrRestoreTopupFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *topupCommandRepository) DeleteTopupPermanent(ctx context.Context, topup_id int) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`DELETE FROM topups WHERE topup_id = ? AND deleted_at IS NOT NULL`, topup_id).Error
	if err != nil {
		return false, topup_errors.ErrDeleteTopupPermanentFailed.WithInternal(err)
	}
	return true, nil
}

func (r *topupCommandRepository) RestoreAllTopup(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`UPDATE topups SET deleted_at = NULL WHERE deleted_at IS NOT NULL`).Error
	if err != nil {
		return false, topup_errors.ErrRestoreAllTopupFailed.WithInternal(err)
	}
	return true, nil
}

func (r *topupCommandRepository) DeleteAllTopupPermanent(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`DELETE FROM topups WHERE deleted_at IS NOT NULL`).Error
	if err != nil {
		return false, topup_errors.ErrDeleteAllTopupPermanentFailed.WithInternal(err)
	}
	return true, nil
}
