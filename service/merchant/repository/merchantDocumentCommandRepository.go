package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors/repository"
	"gorm.io/gorm"
)

type merchantDocumentCommandRepository struct {
	db *gorm.DB
}

func NewMerchantDocumentCommandRepository(db *gorm.DB) MerchantDocumentCommandRepository {
	return &merchantDocumentCommandRepository{db: db}
}

func (r *merchantDocumentCommandRepository) CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocumentCreateRow, error) {
	doc := models.MerchantDocument{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
	}
	if err := r.db.WithContext(ctx).Create(&doc).Error; err != nil {
		return nil, merchantdocument_errors.ErrCreateMerchantDocumentFailed.WithInternal(err)
	}
	return &models.MerchantDocumentCreateRow{
		DocumentID:   doc.DocumentID,
		MerchantID:   doc.MerchantID,
		DocumentType: doc.DocumentType,
		DocumentUrl:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         doc.Note,
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}, nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocumentUpdateRow, error) {
	result := r.db.WithContext(ctx).Model(&models.MerchantDocument{}).
		Where("document_id = ? AND deleted_at IS NULL", *request.DocumentID).
		Updates(map[string]interface{}{
			"document_type": request.DocumentType,
			"document_url":  request.DocumentUrl,
			"status":        request.Status,
		})
	if result.Error != nil {
		return nil, merchantdocument_errors.ErrUpdateMerchantDocumentFailed.WithInternal(result.Error)
	}
	var doc models.MerchantDocumentUpdateRow
	if err := r.db.WithContext(ctx).Table("merchant_documents").
		Where("document_id = ?", *request.DocumentID).Scan(&doc).Error; err != nil {
		return nil, merchantdocument_errors.ErrUpdateMerchantDocumentFailed.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocumentUpdateRow, error) {
	result := r.db.WithContext(ctx).Model(&models.MerchantDocument{}).
		Where("document_id = ? AND deleted_at IS NULL", *request.DocumentID).
		Update("status", request.Status)
	if result.Error != nil {
		return nil, merchantdocument_errors.ErrUpdateMerchantDocumentStatusFailed.WithInternal(result.Error)
	}
	var doc models.MerchantDocumentUpdateRow
	if err := r.db.WithContext(ctx).Table("merchant_documents").
		Where("document_id = ?", *request.DocumentID).Scan(&doc).Error; err != nil {
		return nil, merchantdocument_errors.ErrUpdateMerchantDocumentStatusFailed.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) TrashedMerchantDocument(ctx context.Context, documentID int) (*models.MerchantDocument, error) {
	result := r.db.WithContext(ctx).Where("document_id = ? AND deleted_at IS NULL", documentID).Delete(&models.MerchantDocument{})
	if result.Error != nil {
		return nil, merchantdocument_errors.ErrTrashedMerchantDocumentFailed.WithInternal(result.Error)
	}
	var doc models.MerchantDocument
	if err := r.db.WithContext(ctx).Unscoped().Where("document_id = ?", documentID).First(&doc).Error; err != nil {
		return nil, merchantdocument_errors.ErrTrashedMerchantDocumentFailed.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) RestoreMerchantDocument(ctx context.Context, documentID int) (*models.MerchantDocument, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.MerchantDocument{}).
		Where("document_id = ?", documentID).Update("deleted_at", nil)
	if result.Error != nil {
		return nil, merchantdocument_errors.ErrRestoreMerchantDocumentFailed.WithInternal(result.Error)
	}
	var doc models.MerchantDocument
	if err := r.db.WithContext(ctx).Where("document_id = ?", documentID).First(&doc).Error; err != nil {
		return nil, merchantdocument_errors.ErrRestoreMerchantDocumentFailed.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("document_id = ?", documentID).Delete(&models.MerchantDocument{})
	if result.Error != nil {
		return false, merchantdocument_errors.ErrDeleteMerchantDocumentPermanentFailed.WithInternal(result.Error)
	}
	return true, nil
}

func (r *merchantDocumentCommandRepository) RestoreAllMerchantDocument(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.MerchantDocument{}).
		Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, merchantdocument_errors.ErrRestoreAllMerchantDocumentsFailed.WithInternal(result.Error)
	}
	return true, nil
}

func (r *merchantDocumentCommandRepository) DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.MerchantDocument{})
	if result.Error != nil {
		return false, merchantdocument_errors.ErrDeleteAllMerchantDocumentsPermanentFailed.WithInternal(result.Error)
	}
	return true, nil
}
