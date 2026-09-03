package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
	"gorm.io/gorm"
)

type saldoCommandRepository struct {
	db *gorm.DB
}

func NewSaldoCommandRepository(db *gorm.DB) SaldoCommandRepository {
	return &saldoCommandRepository{db: db}
}

func (r *saldoCommandRepository) CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*models.CreateSaldoRow, error) {
	var saldo models.Saldo
	err := r.db.WithContext(ctx).Raw(`
		INSERT INTO saldos (card_number, total_balance, created_at, updated_at)
		VALUES (?, ?, current_timestamp, current_timestamp)
		ON CONFLICT (card_number) WHERE deleted_at IS NULL
		DO UPDATE SET total_balance = EXCLUDED.total_balance, updated_at = current_timestamp
		RETURNING saldo_id, card_number, total_balance, withdraw_amount, withdraw_time, created_at, updated_at`,
		request.CardNumber, int32(request.TotalBalance),
	).Scan(&saldo).Error
	if err != nil {
		return nil, sharedErrors.ErrConstraintOrFailed(err, "Saldo", "create saldo")
	}

	return &models.CreateSaldoRow{
		SaldoID:      saldo.SaldoID,
		CardNumber:   saldo.CardNumber,
		TotalBalance: saldo.TotalBalance,
		WithdrawAmount: saldo.WithdrawAmount,
		WithdrawTime:   saldo.WithdrawTime,
		CreatedAt:      saldo.CreatedAt,
		UpdatedAt:      saldo.UpdatedAt,
	}, nil
}

func (r *saldoCommandRepository) CreateSaldoIfNotExists(ctx context.Context, request *requests.CreateSaldoRequest) error {
	// Use raw SQL to match ON CONFLICT DO NOTHING behavior
	err := r.db.WithContext(ctx).Exec(
		`INSERT INTO saldos (card_number, total_balance, created_at, updated_at)
		VALUES (?, ?, current_timestamp, current_timestamp)
		ON CONFLICT (card_number) WHERE deleted_at IS NULL DO NOTHING`,
		request.CardNumber, int32(request.TotalBalance),
	).Error
	if err != nil {
		return sharedErrors.ErrConstraintOrFailed(err, "Saldo", "create saldo")
	}
	return nil
}

func (r *saldoCommandRepository) UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*models.UpdateSaldoRow, error) {
	var saldo models.Saldo
	if err := r.db.WithContext(ctx).
		Where("saldo_id = ? AND deleted_at IS NULL", *request.SaldoID).
		First(&saldo).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Saldo", "update saldo")
	}

	saldo.CardNumber = request.CardNumber
	saldo.TotalBalance = int32(request.TotalBalance)
	saldo.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&saldo).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Saldo", "update saldo")
	}

	return &models.UpdateSaldoRow{
		SaldoID:        saldo.SaldoID,
		CardNumber:     saldo.CardNumber,
		TotalBalance:   saldo.TotalBalance,
		WithdrawAmount: saldo.WithdrawAmount,
		WithdrawTime:   saldo.WithdrawTime,
		CreatedAt:      saldo.CreatedAt,
		UpdatedAt:      saldo.UpdatedAt,
	}, nil
}

func (r *saldoCommandRepository) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*models.UpdateSaldoBalanceRow, error) {
	var saldo models.Saldo
	if err := r.db.WithContext(ctx).
		Where("card_number = ? AND deleted_at IS NULL", request.CardNumber).
		First(&saldo).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Saldo", "update saldo balance")
	}

	saldo.TotalBalance = int32(request.TotalBalance)
	saldo.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&saldo).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Saldo", "update saldo balance")
	}

	return &models.UpdateSaldoBalanceRow{
		SaldoID:        saldo.SaldoID,
		CardNumber:     saldo.CardNumber,
		TotalBalance:   saldo.TotalBalance,
		WithdrawAmount: saldo.WithdrawAmount,
		WithdrawTime:   saldo.WithdrawTime,
		CreatedAt:      saldo.CreatedAt,
		UpdatedAt:      saldo.UpdatedAt,
	}, nil
}

func (r *saldoCommandRepository) UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*models.UpdateSaldoWithdrawRow, error) {
	var saldo models.Saldo
	if err := r.db.WithContext(ctx).
		Where("card_number = ? AND deleted_at IS NULL", request.CardNumber).
		First(&saldo).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Saldo", "update saldo withdraw")
	}

	var withdrawAmount int32
	if request.WithdrawAmount != nil {
		withdrawAmount = int32(*request.WithdrawAmount)
	}

	saldo.WithdrawAmount = &withdrawAmount
	saldo.WithdrawTime = request.WithdrawTime
	saldo.TotalBalance -= withdrawAmount
	saldo.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&saldo).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Saldo", "update saldo withdraw")
	}

	return &models.UpdateSaldoWithdrawRow{
		SaldoID:        saldo.SaldoID,
		CardNumber:     saldo.CardNumber,
		TotalBalance:   saldo.TotalBalance,
		WithdrawAmount: saldo.WithdrawAmount,
		WithdrawTime:   saldo.WithdrawTime,
		CreatedAt:      saldo.CreatedAt,
		UpdatedAt:      saldo.UpdatedAt,
	}, nil
}

func (r *saldoCommandRepository) TrashedSaldo(ctx context.Context, saldo_id int) (*models.Saldo, error) {
	var saldo models.Saldo
	if err := r.db.WithContext(ctx).
		Where("saldo_id = ? AND deleted_at IS NULL", saldo_id).
		First(&saldo).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Saldo", "trash saldo")
	}

	now := time.Now()
	saldo.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	saldo.UpdatedAt = now

	if err := r.db.WithContext(ctx).Save(&saldo).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Saldo", "trash saldo")
	}

	return &saldo, nil
}

func (r *saldoCommandRepository) RestoreSaldo(ctx context.Context, saldo_id int) (*models.Saldo, error) {
	var saldo models.Saldo
	if err := r.db.WithContext(ctx).Unscoped().
		Where("saldo_id = ? AND deleted_at IS NOT NULL", saldo_id).
		First(&saldo).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Saldo", "restore saldo")
	}

	saldo.DeletedAt = gorm.DeletedAt{Valid: false}
	saldo.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Unscoped().Save(&saldo).Error; err != nil {
		return nil, sharedErrors.ErrNoRowsOrFailed(err, "Saldo", "restore saldo")
	}

	return &saldo, nil
}

func (r *saldoCommandRepository) DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("saldo_id = ? AND deleted_at IS NOT NULL", saldo_id).
		Delete(&models.Saldo{})

	if result.Error != nil {
		return false, sharedErrors.ErrNoRowsOrFailed(result.Error, "Saldo", "delete saldo")
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *saldoCommandRepository) RestoreAllSaldo(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Model(&models.Saldo{}).
		Where("deleted_at IS NOT NULL").
		Update("deleted_at", nil)

	if result.Error != nil {
		return false, saldo_errors.ErrRestoreAllSaldosFailed.WithInternal(result.Error)
	}

	return true, nil
}

func (r *saldoCommandRepository) DeleteAllSaldoPermanent(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("deleted_at IS NOT NULL").
		Delete(&models.Saldo{})

	if result.Error != nil {
		return false, saldo_errors.ErrDeleteAllSaldosPermanentFailed.WithInternal(result.Error)
	}

	return true, nil
}
