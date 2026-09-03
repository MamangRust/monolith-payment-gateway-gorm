package mencache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type MerchantQueryCache interface {
	GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListRow, *int, bool)
	SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants, data []*models.MerchantListRow, total *int)

	GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListWithDeletedRow, *int, bool)
	SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants, data []*models.MerchantListWithDeletedRow, total *int)

	GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantListWithDeletedRow, *int, bool)
	SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants, data []*models.MerchantListWithDeletedRow, total *int)

	GetCachedMerchant(ctx context.Context, id int) (*models.MerchantAllFieldsRow, bool)
	SetCachedMerchant(ctx context.Context, data *models.MerchantAllFieldsRow)

	GetCachedMerchantsByUserId(ctx context.Context, userId int) ([]*models.MerchantListByUserRow, bool)
	SetCachedMerchantsByUserId(ctx context.Context, userId int, data []*models.MerchantListByUserRow)

	GetCachedMerchantByApiKey(ctx context.Context, apiKey string) (*models.MerchantAllFieldsRow, bool)
	SetCachedMerchantByApiKey(ctx context.Context, apiKey string, data *models.MerchantAllFieldsRow)
}

type MerchantDocumentQueryCache interface {
	GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListRow, *int, bool)
	SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*models.MerchantDocumentListRow, total *int)

	GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListWithDeletedRow, *int, bool)
	SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*models.MerchantDocumentListWithDeletedRow, total *int)

	GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentListWithDeletedRow, *int, bool)
	SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*models.MerchantDocumentListWithDeletedRow, total *int)

	GetCachedMerchantDocument(ctx context.Context, id int) (*models.MerchantDocumentAllFieldsRow, bool)
	SetCachedMerchantDocument(ctx context.Context, id int, data *models.MerchantDocumentAllFieldsRow)
}

type MerchantCommandCache interface {
	DeleteCachedMerchant(ctx context.Context, id int)
	InvalidateMerchantListCaches(ctx context.Context)
	DeleteCachedMerchantByUserID(ctx context.Context, userID int)
	DeleteCachedMerchantByAPIKey(ctx context.Context, apiKey string)
}

type MerchantDocumentCommandCache interface {
	DeleteCachedMerchantDocuments(ctx context.Context, id int)
}

type MerchantTransactionCache interface {
	GetCacheAllMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*models.MerchantTransactionRow, *int, bool)
	SetCacheAllMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions, data []*models.MerchantTransactionRow, total *int)

	GetCacheMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*models.MerchantTransactionRow, *int, bool)
	SetCacheMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactionsById, data []*models.MerchantTransactionRow, total *int)

	GetCacheMerchantTransactionApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*models.MerchantTransactionRow, *int, bool)
	SetCacheMerchantTransactionApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey, data []*models.MerchantTransactionRow, total *int)
}
