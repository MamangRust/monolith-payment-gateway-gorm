package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
)

type topupQueryRepository struct {
	db *gorm.DB
}

func NewTopupQueryRepository(db *gorm.DB) TopupQueryRepository {
	return &topupQueryRepository{db: db}
}

func (r *topupQueryRepository) FindAllTopups(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.TopupListRow
	base := `SELECT topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at, COUNT(*) OVER () AS total_count FROM topups WHERE deleted_at IS NULL`
	args := []interface{}{}
	if req.Search != "" {
		base += " AND (card_number ILIKE ? OR topup_method ILIKE ? OR status ILIKE ?)"
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	base += " ORDER BY topup_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)
	if err := r.db.WithContext(ctx).Raw(base, args...).Scan(&results).Error; err != nil {
		return nil, topup_errors.ErrFindAllTopupsFailed.WithInternal(err)
	}
	return results, nil
}

func (r *topupQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.TopupListRow
	base := `SELECT topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at, COUNT(*) OVER () AS total_count FROM topups WHERE deleted_at IS NULL`
	args := []interface{}{}
	if req.Search != "" {
		base += " AND (card_number ILIKE ? OR topup_method ILIKE ?)"
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}
	base += " ORDER BY topup_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)
	if err := r.db.WithContext(ctx).Raw(base, args...).Scan(&results).Error; err != nil {
		return nil, topup_errors.ErrFindTopupsByActiveFailed.WithInternal(err)
	}
	return results, nil
}

func (r *topupQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListWithDeletedRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.TopupListWithDeletedRow
	base := `SELECT topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count FROM topups WHERE deleted_at IS NOT NULL`
	args := []interface{}{}
	if req.Search != "" {
		base += " AND (card_number ILIKE ? OR topup_method ILIKE ?)"
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}
	base += " ORDER BY topup_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)
	if err := r.db.WithContext(ctx).Raw(base, args...).Scan(&results).Error; err != nil {
		return nil, topup_errors.ErrFindTopupsByTrashedFailed.WithInternal(err)
	}
	return results, nil
}

func (r *topupQueryRepository) FindAllTopupByCardNumber(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*models.TopupListRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.TopupListRow
	base := `SELECT topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at, COUNT(*) OVER () AS total_count FROM topups WHERE deleted_at IS NULL AND card_number = ?`
	args := []interface{}{req.CardNumber}
	if req.Search != "" {
		base += " AND (topup_method ILIKE ? OR status ILIKE ?)"
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}
	base += " ORDER BY topup_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)
	if err := r.db.WithContext(ctx).Raw(base, args...).Scan(&results).Error; err != nil {
		return nil, topup_errors.ErrFindAllTopupsFailed.WithInternal(err)
	}
	return results, nil
}

func (r *topupQueryRepository) FindById(ctx context.Context, topup_id int) (*models.TopupAllFieldsRow, error) {
	if topup_id <= 0 {
		return nil, topup_errors.ErrFindTopupByIdFailed.WithInternal(errors.New("invalid id"))
	}
	var result models.TopupAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`SELECT topup_id, topup_no, card_number, topup_amount, topup_method, topup_time, status, created_at, updated_at FROM topups WHERE topup_id = ? AND deleted_at IS NULL`, topup_id).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, topup_errors.ErrFindTopupByIdFailed.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &result, nil
}
