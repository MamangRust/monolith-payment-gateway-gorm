package mencache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TopupQueryCache interface {
	GetCachedTopupsCache(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListRow, *int, bool)
	SetCachedTopupsCache(ctx context.Context, req *requests.FindAllTopups, data []*models.TopupListRow, total *int)

	GetCacheTopupByCardCache(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*models.TopupListRow, *int, bool)
	SetCacheTopupByCardCache(ctx context.Context, req *requests.FindAllTopupsByCardNumber, data []*models.TopupListRow, total *int)

	GetCachedTopupActiveCache(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListRow, *int, bool)
	SetCachedTopupActiveCache(ctx context.Context, req *requests.FindAllTopups, data []*models.TopupListRow, total *int)

	GetCachedTopupTrashedCache(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListWithDeletedRow, *int, bool)
	SetCachedTopupTrashedCache(ctx context.Context, req *requests.FindAllTopups, data []*models.TopupListWithDeletedRow, total *int)

	GetCachedTopupCache(ctx context.Context, id int) (*models.TopupAllFieldsRow, bool)
	SetCachedTopupCache(ctx context.Context, data *models.TopupAllFieldsRow)
}

type TopupCommandCache interface {
	DeleteCachedTopupCache(ctx context.Context, id int)
}
