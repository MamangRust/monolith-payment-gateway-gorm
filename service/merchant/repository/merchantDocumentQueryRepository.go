package repository

import (
	"fmt"
	"context"
	"errors"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors/repository"
	"gorm.io/gorm"
)

type merchantDocumentQueryRepository struct {
	db *gorm.DB
}

func NewMerchantDocumentQueryRepository(db *gorm.DB) MerchantDocumentQueryRepository {
	return &merchantDocumentQueryRepository{db: db}
}

func (r *merchantDocumentQueryRepository) FindAllDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var totalCount int64
	countQuery := r.db.WithContext(ctx).Table("merchant_documents").Where("deleted_at IS NULL")
	if req.Search != "" {
		countQuery = countQuery.Where("(document_type ILIKE ? OR status ILIKE ?)", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, merchantdocument_errors.ErrFindAllMerchantDocumentsFailed.WithInternal(err)
	}

	var docs []*models.MerchantDocumentListRow
	query := r.db.WithContext(ctx).Table("merchant_documents").
		Select(fmt.Sprintf("merchant_documents.*, %d as total_count", totalCount)).
		Where("merchant_documents.deleted_at IS NULL")
	if req.Search != "" {
		query = query.Where("(merchant_documents.document_type ILIKE ? OR merchant_documents.status ILIKE ?)", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := query.Order("merchant_documents.document_id DESC").Limit(req.PageSize).Offset(offset).Scan(&docs).Error; err != nil {
		return nil, merchantdocument_errors.ErrFindAllMerchantDocumentsFailed.WithInternal(err)
	}
	return docs, nil
}

func (r *merchantDocumentQueryRepository) FindByIdDocument(ctx context.Context, id int) (*models.MerchantDocumentAllFieldsRow, error) {
	var doc models.MerchantDocumentAllFieldsRow
	err := r.db.WithContext(ctx).Table("merchant_documents").
		Where("document_id = ? AND deleted_at IS NULL", id).First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merchantdocument_errors.ErrMerchantDocumentNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentQueryRepository) FindByActiveDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListWithDeletedRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var totalCount int64
	countQuery := r.db.WithContext(ctx).Table("merchant_documents").Where("deleted_at IS NULL")
	if req.Search != "" {
		countQuery = countQuery.Where("(document_type ILIKE ? OR status ILIKE ?)", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, merchantdocument_errors.ErrFindActiveMerchantDocumentsFailed.WithInternal(err)
	}

	var docs []*models.MerchantDocumentListWithDeletedRow
	query := r.db.WithContext(ctx).Table("merchant_documents").
		Select(fmt.Sprintf("merchant_documents.*, %d as total_count", totalCount)).
		Where("merchant_documents.deleted_at IS NULL")
	if req.Search != "" {
		query = query.Where("(merchant_documents.document_type ILIKE ? OR merchant_documents.status ILIKE ?)", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := query.Order("merchant_documents.document_id DESC").Limit(req.PageSize).Offset(offset).Scan(&docs).Error; err != nil {
		return nil, merchantdocument_errors.ErrFindActiveMerchantDocumentsFailed.WithInternal(err)
	}
	return docs, nil
}

func (r *merchantDocumentQueryRepository) FindByTrashedDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListWithDeletedRow, error) {
	offset := (req.Page - 1) * req.PageSize
	var totalCount int64
	countQuery := r.db.WithContext(ctx).Table("merchant_documents").Where("deleted_at IS NOT NULL")
	if req.Search != "" {
		countQuery = countQuery.Where("(document_type ILIKE ? OR status ILIKE ?)", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, merchantdocument_errors.ErrFindTrashedMerchantDocumentsFailed.WithInternal(err)
	}

	var docs []*models.MerchantDocumentListWithDeletedRow
	query := r.db.WithContext(ctx).Table("merchant_documents").
		Select(fmt.Sprintf("merchant_documents.*, %d as total_count", totalCount)).
		Where("merchant_documents.deleted_at IS NOT NULL")
	if req.Search != "" {
		query = query.Where("(merchant_documents.document_type ILIKE ? OR merchant_documents.status ILIKE ?)", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := query.Order("merchant_documents.document_id DESC").Limit(req.PageSize).Offset(offset).Scan(&docs).Error; err != nil {
		return nil, merchantdocument_errors.ErrFindTrashedMerchantDocumentsFailed.WithInternal(err)
	}
	return docs, nil
}
