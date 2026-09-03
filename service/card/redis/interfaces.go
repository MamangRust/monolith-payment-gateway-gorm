package mencache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type CardQueryCache interface {
	GetByIdCache(ctx context.Context, cardID int) (*models.CardAllFieldsRow, bool)
	GetByUserIDCache(ctx context.Context, userID int) (*models.CardAllFieldsRow, bool)
	GetByCardNumberCache(ctx context.Context, cardNumber string) (*models.CardAllFieldsRow, bool)
	GetUserCardByCardNumberCache(ctx context.Context, cardNumber string) (*models.CardByEmailRow, bool)
	GetFindAllCache(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListRow, *int, bool)
	GetByActiveCache(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListWithDeletedRow, *int, bool)
	GetByTrashedCache(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListWithDeletedRow, *int, bool)
	SetByIdCache(ctx context.Context, cardID int, data *models.CardAllFieldsRow)
	SetByUserIDCache(ctx context.Context, userID int, data *models.CardAllFieldsRow)
	SetByCardNumberCache(ctx context.Context, cardNumber string, data *models.CardAllFieldsRow)
	SetFindAllCache(ctx context.Context, req *requests.FindAllCards, data []*models.CardListRow, totalRecords *int)
	SetByActiveCache(ctx context.Context, req *requests.FindAllCards, data []*models.CardListWithDeletedRow, totalRecords *int)
	SetUserCardByCardNumberCache(ctx context.Context, cardNumber string, data *models.CardByEmailRow)
	SetByTrashedCache(ctx context.Context, req *requests.FindAllCards, data []*models.CardListWithDeletedRow, totalRecords *int)
	DeleteByIdCache(ctx context.Context, cardID int)
	DeleteByUserIDCache(ctx context.Context, userID int)
	DeleteByCardNumberCache(ctx context.Context, cardNumber string)
}

type CardCommandCache interface {
	DeleteCardCommandCache(ctx context.Context, id int)
	DeleteCardCache(ctx context.Context, id, userID int, cardNumber string)
}
