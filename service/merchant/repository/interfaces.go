package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type UserRepository interface {
	FindById(ctx context.Context, user_id int) (*models.UserByIDRow, error)
}

type MerchantQueryRepository interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListWithDeletedRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListWithDeletedRow, error)
	FindByApiKey(ctx context.Context, api_key string) (*models.MerchantAllFieldsRow, error)
	FindByMerchantId(ctx context.Context, merchant_id int) (*models.MerchantAllFieldsRow, error)
	FindByName(ctx context.Context, name string) (*models.MerchantAllFieldsRow, error)
	FindByMerchantUserId(ctx context.Context, user_id int) ([]*models.MerchantListByUserRow, error)
}

type MerchantDocumentQueryRepository interface {
	FindAllDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListRow, error)
	FindByIdDocument(ctx context.Context, id int) (*models.MerchantDocumentAllFieldsRow, error)
	FindByActiveDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListWithDeletedRow, error)
	FindByTrashedDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListWithDeletedRow, error)
}

type MerchantTransactionRepository interface {
	FindAllTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*models.MerchantTransactionRow, error)
	FindAllTransactionsByMerchant(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*models.MerchantTransactionRow, error)
	FindAllTransactionsByApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*models.MerchantTransactionRow, error)
}

type MerchantCommandRepository interface {
	CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*models.MerchantAllFieldsRow, error)
	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*models.MerchantAllFieldsRow, error)
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*models.MerchantAllFieldsRow, error)
	TrashedMerchant(ctx context.Context, merchantId int) (*models.Merchant, error)
	RestoreMerchant(ctx context.Context, merchantId int) (*models.Merchant, error)
	DeleteMerchantPermanent(ctx context.Context, merchantId int) (bool, error)
	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type MerchantDocumentCommandRepository interface {
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocumentCreateRow, error)
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocumentUpdateRow, error)
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocumentUpdateRow, error)
	TrashedMerchantDocument(ctx context.Context, merchant_document_id int) (*models.MerchantDocument, error)
	RestoreMerchantDocument(ctx context.Context, merchant_document_id int) (*models.MerchantDocument, error)
	DeleteMerchantDocumentPermanent(ctx context.Context, merchant_document_id int) (bool, error)
	RestoreAllMerchantDocument(ctx context.Context) (bool, error)
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error)
}
