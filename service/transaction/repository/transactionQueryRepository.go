package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
)

type transactionQueryRepository struct {
	db *gorm.DB
}

func NewTransactionQueryRepository(db *gorm.DB) TransactionQueryRepository {
	return &transactionQueryRepository{db: db}
}

func (r *transactionQueryRepository) FindAllTransactions(ctx context.Context, req *requests.FindAllTransactions) ([]*models.TransactionListRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.TransactionListRow

	baseQuery := `
		SELECT t.transaction_id, t.transaction_no, t.card_number, t.amount, t.payment_method,
			t.merchant_id, t.transaction_time, t.status, t.created_at, t.updated_at,
			COUNT(*) OVER () AS total_count
		FROM transactions t
		WHERE t.deleted_at IS NULL`
	args := []interface{}{}

	if req.Search != "" {
		baseQuery += " AND (t.payment_method ILIKE ?)"
		args = append(args, "%"+req.Search+"%")
	}

	baseQuery += " ORDER BY t.transaction_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)

	err := r.db.WithContext(ctx).Raw(baseQuery, args...).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrFindAllTransactionsFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionQueryRepository) FindAllTransactionByCardNumber(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*models.TransactionListRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.TransactionListRow

	baseQuery := `
		SELECT t.transaction_id, t.transaction_no, t.card_number, t.amount, t.payment_method,
			t.merchant_id, t.transaction_time, t.status, t.created_at, t.updated_at,
			COUNT(*) OVER () AS total_count
		FROM transactions t
		WHERE t.deleted_at IS NULL AND t.card_number = ?`
	args := []interface{}{req.CardNumber}

	if req.Search != "" {
		baseQuery += " AND (t.payment_method ILIKE ?)"
		args = append(args, "%"+req.Search+"%")
	}

	baseQuery += " ORDER BY t.transaction_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)

	err := r.db.WithContext(ctx).Raw(baseQuery, args...).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrFindTransactionsByCardNumberFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTransactions) ([]*models.TransactionListRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.TransactionListRow

	baseQuery := `
		SELECT t.transaction_id, t.transaction_no, t.card_number, t.amount, t.payment_method,
			t.merchant_id, t.transaction_time, t.status, t.created_at, t.updated_at,
			COUNT(*) OVER () AS total_count
		FROM transactions t
		WHERE t.deleted_at IS NULL`
	args := []interface{}{}

	if req.Search != "" {
		baseQuery += " AND (t.payment_method ILIKE ?)"
		args = append(args, "%"+req.Search+"%")
	}

	baseQuery += " ORDER BY t.transaction_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)

	err := r.db.WithContext(ctx).Raw(baseQuery, args...).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrFindActiveTransactionsFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTransactions) ([]*models.TransactionListWithDeletedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.TransactionListWithDeletedRow

	baseQuery := `
		SELECT t.transaction_id, t.transaction_no, t.card_number, t.amount, t.payment_method,
			t.merchant_id, t.transaction_time, t.status, t.created_at, t.updated_at, t.deleted_at,
			COUNT(*) OVER () AS total_count
		FROM transactions t
		WHERE t.deleted_at IS NOT NULL`
	args := []interface{}{}

	if req.Search != "" {
		baseQuery += " AND (t.payment_method ILIKE ?)"
		args = append(args, "%"+req.Search+"%")
	}

	baseQuery += " ORDER BY t.transaction_time DESC LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, offset)

	err := r.db.WithContext(ctx).Raw(baseQuery, args...).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrFindTrashedTransactionsFailed.WithInternal(err)
	}

	return results, nil
}

func (r *transactionQueryRepository) FindById(ctx context.Context, transaction_id int) (*models.TransactionAllFieldsRow, error) {
	if transaction_id <= 0 {
		return nil, sharedErrors.NewBadRequestError("transaction ID must be greater than zero")
	}

	var result models.TransactionAllFieldsRow

	err := r.db.WithContext(ctx).Raw(`
		SELECT t.transaction_id, t.transaction_no, t.card_number, t.amount, t.payment_method,
			t.merchant_id, t.transaction_time, t.status, t.created_at, t.updated_at
		FROM transactions t
		WHERE t.transaction_id = ? AND t.deleted_at IS NULL
	`, transaction_id).First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, transaction_errors.ErrTransactionNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &result, nil
}

func (r *transactionQueryRepository) FindTransactionByMerchantId(ctx context.Context, merchant_id int) ([]*models.TransactionListRow, error) {
	var results []*models.TransactionListRow

	err := r.db.WithContext(ctx).Raw(`
		SELECT t.transaction_id, t.transaction_no, t.card_number, t.amount, t.payment_method,
			t.merchant_id, t.transaction_time, t.status, t.created_at, t.updated_at,
			COUNT(*) OVER () AS total_count
		FROM transactions t
		WHERE t.merchant_id = ? AND t.deleted_at IS NULL
		ORDER BY t.transaction_time DESC
	`, merchant_id).Scan(&results).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, transaction_errors.ErrFindTransactionByMerchantIdFailed.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return results, nil
}
