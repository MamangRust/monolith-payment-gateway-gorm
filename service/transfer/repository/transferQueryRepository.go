package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferQueryRepository struct {
	db *gorm.DB
}

func NewTransferQueryRepository(db *gorm.DB) TransferQueryRepository {
	return &transferQueryRepository{db: db}
}

func (r *transferQueryRepository) FindAll(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.TransferListRow
	base := `SELECT transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at, COUNT(*) OVER () AS total_count FROM transfers WHERE deleted_at IS NULL`
	args := []interface{}{}
	if req.Search != "" {
		base += " AND (transfer_from ILIKE ? OR transfer_to ILIKE ?)"
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}
	base += " ORDER BY transfer_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)
	if err := r.db.WithContext(ctx).Raw(base, args...).Scan(&results).Error; err != nil {
		return nil, transfer_errors.ErrFindAllTransfersFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transferQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.TransferListRow
	base := `SELECT transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at, COUNT(*) OVER () AS total_count FROM transfers WHERE deleted_at IS NULL`
	args := []interface{}{}
	if req.Search != "" {
		base += " AND (transfer_from ILIKE ? OR transfer_to ILIKE ?)"
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}
	base += " ORDER BY transfer_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)
	if err := r.db.WithContext(ctx).Raw(base, args...).Scan(&results).Error; err != nil {
		return nil, transfer_errors.ErrFindActiveTransfersFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transferQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTransfers) ([]*models.TransferListWithDeletedRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.TransferListWithDeletedRow
	base := `SELECT transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count FROM transfers WHERE deleted_at IS NOT NULL`
	args := []interface{}{}
	if req.Search != "" {
		base += " AND (transfer_from ILIKE ? OR transfer_to ILIKE ?)"
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}
	base += " ORDER BY transfer_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)
	if err := r.db.WithContext(ctx).Raw(base, args...).Scan(&results).Error; err != nil {
		return nil, transfer_errors.ErrFindTrashedTransfersFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transferQueryRepository) FindById(ctx context.Context, id int) (*models.TransferAllFieldsRow, error) {
	if id <= 0 {
		return nil, transfer_errors.ErrFindTransferByIdFailed.WithInternal(errors.New("invalid id"))
	}
	var result models.TransferAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`SELECT transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at FROM transfers WHERE transfer_id = ? AND deleted_at IS NULL`, id).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, transfer_errors.ErrFindTransferByIdFailed.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &result, nil
}

func (r *transferQueryRepository) FindTransferByTransferFrom(ctx context.Context, transfer_from string) ([]*models.TransferBySourceCardRow, error) {
	var results []*models.TransferBySourceCardRow
	err := r.db.WithContext(ctx).Raw(`SELECT transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at FROM transfers WHERE transfer_from = ? AND deleted_at IS NULL ORDER BY transfer_time DESC`, transfer_from).Scan(&results).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, transfer_errors.ErrFindTransferByTransferFromFailed.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return results, nil
}

func (r *transferQueryRepository) FindTransferByTransferTo(ctx context.Context, transfer_to string) ([]*models.TransferByDestinationCardRow, error) {
	var results []*models.TransferByDestinationCardRow
	err := r.db.WithContext(ctx).Raw(`SELECT transfer_id, transfer_no, transfer_from, transfer_to, transfer_amount, transfer_time, status, created_at, updated_at FROM transfers WHERE transfer_to = ? AND deleted_at IS NULL ORDER BY transfer_time DESC`, transfer_to).Scan(&results).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, transfer_errors.ErrFindTransferByTransferToFailed.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return results, nil
}
