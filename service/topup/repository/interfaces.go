package repository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoRepository interface {
	FindByCardNumber(ctx context.Context, card_number string) (*models.Saldo, error)
	UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*models.UpdateSaldoBalanceRow, error)
}

type TopupQueryRepository interface {
	FindAllTopups(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTopups) ([]*models.TopupListWithDeletedRow, error)
	FindAllTopupByCardNumber(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*models.TopupListRow, error)
	FindById(ctx context.Context, topup_id int) (*models.TopupAllFieldsRow, error)
}

type TopupCommandRepository interface {
	CreateTopup(ctx context.Context, request *requests.CreateTopupRequest) (*models.TopupAllFieldsRow, error)
	CreateTopupAtomic(ctx context.Context, request *requests.CreateTopupRequest) (*TopupAtomicResult, error)
	UpdateTopup(ctx context.Context, request *requests.UpdateTopupRequest) (*models.TopupAllFieldsRow, error)
	UpdateTopupAtomic(ctx context.Context, request *requests.UpdateTopupRequest) (*models.TopupAllFieldsRow, error)
	UpdateSaldoBalanceDelta(ctx context.Context, cardNumber string, delta int) error
	UpdateTopupAmount(ctx context.Context, request *requests.UpdateTopupAmount) (*models.TopupAllFieldsRow, error)
	UpdateTopupStatus(ctx context.Context, request *requests.UpdateTopupStatus) (*models.TopupAllFieldsRow, error)
	TrashedTopup(ctx context.Context, topup_id int) (*models.Topup, error)
	RestoreTopup(ctx context.Context, topup_id int) (*models.Topup, error)
	DeleteTopupPermanent(ctx context.Context, topup_id int) (bool, error)
	RestoreAllTopup(ctx context.Context) (bool, error)
	DeleteAllTopupPermanent(ctx context.Context) (bool, error)
}

type CardRepository interface {
	FindUserCardByCardNumber(ctx context.Context, card_number string) (*models.CardByEmailRow, error)
	FindCardByCardNumber(ctx context.Context, card_number string) (*models.CardAllFieldsRow, error)
	UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*models.CardUpdateRow, error)
}
