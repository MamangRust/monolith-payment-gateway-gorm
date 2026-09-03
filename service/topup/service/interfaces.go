package service

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// TopupQueryService defines the read-only operations for querying topup data.
type TopupQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListRow, *int, error)
	FindAllByCardNumber(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*models.TopupListRow, *int, error)
	FindById(ctx context.Context, topupID int) (*models.TopupAllFieldsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListWithDeletedRow, *int, error)
}

type TopupCommandService interface {
	CreateTopup(ctx context.Context, request *requests.CreateTopupRequest) (*models.TopupAllFieldsRow, error)
	UpdateTopup(ctx context.Context, request *requests.UpdateTopupRequest) (*models.TopupAllFieldsRow, error)
	TrashedTopup(ctx context.Context, topup_id int) (*models.Topup, error)
	RestoreTopup(ctx context.Context, topup_id int) (*models.Topup, error)
	DeleteTopupPermanent(ctx context.Context, topup_id int) (bool, error)

	RestoreAllTopup(ctx context.Context) (bool, error)
	DeleteAllTopupPermanent(ctx context.Context) (bool, error)
}
