package repository

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	"gorm.io/gorm"
)

type merchantTransactionRepository struct {
	db *gorm.DB
}

func NewMerchantTransactionRepository(db *gorm.DB) MerchantTransactionRepository {
	return &merchantTransactionRepository{db: db}
}

func (r *merchantTransactionRepository) FindAllTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*models.MerchantTransactionRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var totalCount int64
	countQuery := r.db.WithContext(ctx).Table("transactions").Where("transactions.deleted_at IS NULL")
	if req.Search != "" {
		countQuery = countQuery.Where("(transactions.card_number ILIKE ? OR transactions.payment_method ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, merchant_errors.ErrFindAllTransactionsFailed.WithInternal(err)
	}

	var txns []*models.MerchantTransactionRow
	query := r.db.WithContext(ctx).Table("transactions").
		Select(fmt.Sprintf("transactions.*, merchants.name as merchant_name, %d as total_count", totalCount)).
		Joins("JOIN merchants ON merchants.merchant_id = transactions.merchant_id AND merchants.deleted_at IS NULL").
		Where("transactions.deleted_at IS NULL")
	if req.Search != "" {
		query = query.Where("(transactions.card_number ILIKE ? OR transactions.payment_method ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := query.Order("transactions.transaction_id DESC").Limit(req.PageSize).Offset(offset).Scan(&txns).Error; err != nil {
		return nil, merchant_errors.ErrFindAllTransactionsFailed.WithInternal(err)
	}
	return txns, nil
}

func (r *merchantTransactionRepository) FindAllTransactionsByMerchant(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*models.MerchantTransactionRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var totalCount int64
	countQuery := r.db.WithContext(ctx).Table("transactions").Where("transactions.deleted_at IS NULL AND transactions.merchant_id = ?", req.MerchantID)
	if req.Search != "" {
		countQuery = countQuery.Where("(transactions.card_number ILIKE ? OR transactions.payment_method ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, merchant_errors.ErrFindAllTransactionsByMerchantFailed.WithInternal(err)
	}

	var txns []*models.MerchantTransactionRow
	query := r.db.WithContext(ctx).Table("transactions").
		Select(fmt.Sprintf("transactions.*, merchants.name as merchant_name, %d as total_count", totalCount)).
		Joins("JOIN merchants ON merchants.merchant_id = transactions.merchant_id AND merchants.deleted_at IS NULL").
		Where("transactions.deleted_at IS NULL AND transactions.merchant_id = ?", req.MerchantID)
	if req.Search != "" {
		query = query.Where("(transactions.card_number ILIKE ? OR transactions.payment_method ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := query.Order("transactions.transaction_id DESC").Limit(req.PageSize).Offset(offset).Scan(&txns).Error; err != nil {
		return nil, merchant_errors.ErrFindAllTransactionsByMerchantFailed.WithInternal(err)
	}
	return txns, nil
}

func (r *merchantTransactionRepository) FindAllTransactionsByApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*models.MerchantTransactionRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var totalCount int64
	countQuery := r.db.WithContext(ctx).Table("transactions").
		Joins("JOIN merchants ON merchants.merchant_id = transactions.merchant_id AND merchants.deleted_at IS NULL").
		Where("transactions.deleted_at IS NULL AND merchants.api_key = ?", req.ApiKey)
	if req.Search != "" {
		countQuery = countQuery.Where("(transactions.card_number ILIKE ? OR transactions.payment_method ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, merchant_errors.ErrFindAllTransactionsByApiKeyFailed.WithInternal(err)
	}

	var txns []*models.MerchantTransactionRow
	query := r.db.WithContext(ctx).Table("transactions").
		Select(fmt.Sprintf("transactions.*, merchants.name as merchant_name, %d as total_count", totalCount)).
		Joins("JOIN merchants ON merchants.merchant_id = transactions.merchant_id AND merchants.deleted_at IS NULL").
		Where("transactions.deleted_at IS NULL AND merchants.api_key = ?", req.ApiKey)
	if req.Search != "" {
		query = query.Where("(transactions.card_number ILIKE ? OR transactions.payment_method ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := query.Order("transactions.transaction_id DESC").Limit(req.PageSize).Offset(offset).Scan(&txns).Error; err != nil {
		return nil, merchant_errors.ErrFindAllTransactionsByApiKeyFailed.WithInternal(err)
	}
	return txns, nil
}
