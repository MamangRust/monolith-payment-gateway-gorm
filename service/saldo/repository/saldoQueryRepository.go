package repository

import (
	"context"
	"strings"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
	"gorm.io/gorm"
)

type saldoQueryRepository struct {
	db *gorm.DB
}

func NewSaldoQueryRepository(db *gorm.DB) SaldoQueryRepository {
	return &saldoQueryRepository{db: db}
}

func (r *saldoQueryRepository) FindAllSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.SaldoRow
	query := r.db.WithContext(ctx).Table("saldos").
		Select("saldo_id, card_number, total_balance, withdraw_amount, withdraw_time, created_at, updated_at, COUNT(*) OVER () AS total_count").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		query = query.Where("card_number ILIKE ?", "%"+req.Search+"%")
	}

	query = query.Order("saldo_id ASC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, saldo_errors.ErrFindAllSaldosFailed.WithInternal(err)
	}

	return results, nil
}

func (r *saldoQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.SaldoActiveRow
	query := r.db.WithContext(ctx).Table("saldos").
		Select("saldo_id, card_number, total_balance, withdraw_amount, withdraw_time, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		query = query.Where("card_number ILIKE ?", "%"+req.Search+"%")
	}

	query = query.Order("saldo_id ASC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, saldo_errors.ErrFindActiveSaldosFailed.WithInternal(err)
	}

	return results, nil
}

func (r *saldoQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*models.SaldoTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.SaldoTrashedRow
	query := r.db.WithContext(ctx).Table("saldos").
		Select("saldo_id, card_number, total_balance, withdraw_amount, withdraw_time, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count").
		Where("deleted_at IS NOT NULL")

	if req.Search != "" {
		query = query.Where("card_number ILIKE ?", "%"+strings.TrimSpace(req.Search)+"%")
	}

	query = query.Order("deleted_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, saldo_errors.ErrFindTrashedSaldosFailed.WithInternal(err)
	}

	return results, nil
}

func (r *saldoQueryRepository) FindById(ctx context.Context, saldo_id int) (*models.SaldoByIDRow, error) {
	if saldo_id <= 0 {
		return nil, sharedErrors.NewBadRequestError("saldo ID must be greater than zero")
	}

	var saldo models.SaldoByIDRow
	if err := r.db.WithContext(ctx).
		Table("saldos").
		Select("saldo_id, card_number, total_balance, withdraw_amount, withdraw_time, created_at, updated_at").
		Where("saldo_id = ? AND deleted_at IS NULL", saldo_id).
		First(&saldo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, saldo_errors.ErrSaldoNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &saldo, nil
}

func (r *saldoQueryRepository) FindByCardNumber(ctx context.Context, card_number string) (*models.Saldo, error) {
	var saldo models.Saldo
	if err := r.db.WithContext(ctx).
		Where("card_number = ? AND deleted_at IS NULL", card_number).
		Order("saldo_id DESC").
		First(&saldo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, saldo_errors.ErrSaldoNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &saldo, nil
}
