package service

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type MerchantQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListRow, *int, error)
	FindById(ctx context.Context, merchant_id int) (*models.MerchantAllFieldsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListWithDeletedRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListWithDeletedRow, *int, error)
	FindByApiKey(ctx context.Context, api_key string) (*models.MerchantAllFieldsRow, error)
	FindByMerchantUserId(ctx context.Context, user_id int) ([]*models.MerchantListByUserRow, error)
}

type MerchantDocumentQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListRow, *int, error)

	FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListWithDeletedRow, *int, error)

	FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListWithDeletedRow, *int, error)

	FindById(ctx context.Context, document_id int) (*models.MerchantDocumentAllFieldsRow, error)
}

type MerchantTransactionService interface {
	FindAllTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*models.MerchantTransactionRow, *int, error)
	FindAllTransactionsByApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*models.MerchantTransactionRow, *int, error)
	FindAllTransactionsByMerchant(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*models.MerchantTransactionRow, *int, error)
}

type MerchantCommandService interface {
	CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*models.MerchantAllFieldsRow, error)
	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*models.MerchantAllFieldsRow, error)
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*models.MerchantAllFieldsRow, error)
	TrashedMerchant(ctx context.Context, merchant_id int) (*models.Merchant, error)
	RestoreMerchant(ctx context.Context, merchant_id int) (*models.Merchant, error)
	DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error)

	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type MerchantDocumentCommandService interface {
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocumentCreateRow, error)
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocumentUpdateRow, error)
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocumentUpdateRow, error)
	TrashedMerchantDocument(ctx context.Context, document_id int) (*models.MerchantDocument, error)
	RestoreMerchantDocument(ctx context.Context, document_id int) (*models.MerchantDocument, error)
	DeleteMerchantDocumentPermanent(ctx context.Context, document_id int) (bool, error)
	RestoreAllMerchantDocument(ctx context.Context) (bool, error)
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error)
}
