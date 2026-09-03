package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
)

type withdrawQueryRepository struct {
	db *gorm.DB
}

func NewWithdrawQueryRepository(db *gorm.DB) WithdrawQueryRepository {
	return &withdrawQueryRepository{db: db}
}

func (r *withdrawQueryRepository) FindAll(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.WithdrawListRow
	base := `SELECT withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at, COUNT(*) OVER () AS total_count FROM withdraws WHERE deleted_at IS NULL`
	args := []interface{}{}
	if req.Search != "" {
		base += " AND (card_number ILIKE ? OR status ILIKE ?)"
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}
	base += " ORDER BY withdraw_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)
	if err := r.db.WithContext(ctx).Raw(base, args...).Scan(&results).Error; err != nil {
		return nil, withdraw_errors.ErrFindAllWithdrawsFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.WithdrawListRow
	base := `SELECT withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at, COUNT(*) OVER () AS total_count FROM withdraws WHERE deleted_at IS NULL`
	args := []interface{}{}
	if req.Search != "" {
		base += " AND (card_number ILIKE ? OR status ILIKE ?)"
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}
	base += " ORDER BY withdraw_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)
	if err := r.db.WithContext(ctx).Raw(base, args...).Scan(&results).Error; err != nil {
		return nil, withdraw_errors.ErrFindActiveWithdrawsFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllWithdraws) ([]*models.WithdrawListWithDeletedRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.WithdrawListWithDeletedRow
	base := `SELECT withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count FROM withdraws WHERE deleted_at IS NOT NULL`
	args := []interface{}{}
	if req.Search != "" {
		base += " AND (card_number ILIKE ? OR status ILIKE ?)"
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}
	base += " ORDER BY withdraw_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)
	if err := r.db.WithContext(ctx).Raw(base, args...).Scan(&results).Error; err != nil {
		return nil, withdraw_errors.ErrFindTrashedWithdrawsFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawQueryRepository) FindAllByCardNumber(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*models.WithdrawByCardNumberRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.WithdrawByCardNumberRow
	base := `SELECT withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at, COUNT(*) OVER () AS total_count FROM withdraws WHERE deleted_at IS NULL AND card_number = ?`
	args := []interface{}{req.CardNumber}
	if req.Search != "" {
		base += " AND (status ILIKE ?)"
		args = append(args, "%"+req.Search+"%")
	}
	base += " ORDER BY withdraw_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)
	if err := r.db.WithContext(ctx).Raw(base, args...).Scan(&results).Error; err != nil {
		return nil, withdraw_errors.ErrFindAllWithdrawsFailed.WithInternal(err)
	}
	return results, nil
}

func (r *withdrawQueryRepository) FindById(ctx context.Context, id int) (*models.WithdrawAllFieldsRow, error) {
	if id <= 0 {
		return nil, withdraw_errors.ErrFindWithdrawByIdFailed.WithInternal(errors.New("invalid id"))
	}
	var result models.WithdrawAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`SELECT withdraw_id, withdraw_no, card_number, withdraw_amount, withdraw_time, status, created_at, updated_at FROM withdraws WHERE withdraw_id = ? AND deleted_at IS NULL`, id).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, withdraw_errors.ErrFindWithdrawByIdFailed.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &result, nil
}
