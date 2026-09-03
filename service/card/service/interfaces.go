package service

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type CardAuthorizationService interface {
	Authorize(ctx context.Context, request *requests.AuthorizeCardRequest) (*models.CardAuthTransaction, error)
	Reverse(ctx context.Context, request *requests.ReverseTransactionRequest) (*models.CardAuthTransaction, error)
	GetAuthTransaction(ctx context.Context, txnID string) (*models.CardAuthTransaction, error)
	GetAuthTransactionsByCardNumber(ctx context.Context, cardNumber string, page, pageSize int) ([]*models.GetAuthTxnByCardNumberRow, error)
}

type CardPaymentService interface {
	PostPayment(ctx context.Context, request *requests.PostPaymentRequest) (*models.CardPayment, error)
	GetPaymentHistory(ctx context.Context, cardNumber string, page, pageSize int) ([]*models.CardPaymentListRow, error)
	CountPayments(ctx context.Context, cardNumber string) (int, error)
}

type CardRewardService interface {
	EarnRewards(ctx context.Context, request *requests.EarnRewardsRequest) (*models.CardReward, error)
	GetBalance(ctx context.Context, cardNumber string) (int64, error)
	GetHistory(ctx context.Context, cardNumber string) ([]*models.CardReward, error)
	RedeemRewards(ctx context.Context, cardNumber string, points int64) (int64, error)
}

type BillingEngineService interface {
	TriggerBillingCycle(ctx context.Context, billingCycleDay int) (int, error)
	GetStatement(ctx context.Context, cardNumber string) (*models.BillingCycle, error)
	GetStatementsByCard(ctx context.Context, cardNumber string, page, pageSize int) ([]*models.BillingCycle, error)
	GetBillingCyclesByCardNumber(ctx context.Context, cardNumber string) ([]*models.BillingCycle, error)
}

type CardQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListWithDeletedRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListWithDeletedRow, *int, error)
	FindById(ctx context.Context, card_id int) (*models.CardAllFieldsRow, error)
	FindByUserID(ctx context.Context, userID int) (*models.CardAllFieldsRow, error)
	FindByCardNumber(ctx context.Context, card_number string) (*models.CardAllFieldsRow, error)
	FindUserCardByCardNumber(ctx context.Context, card_number string) (*models.CardByEmailRow, error)
}

type CardDashboardService interface {
	DashboardCard(ctx context.Context) (*response.DashboardCard, error)
	DashboardCardCardNumber(ctx context.Context, cardNumber string) (*response.DashboardCardCardNumber, error)
}

type CardCommandService interface {
	CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*models.CardCreateRow, error)
	UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*models.CardUpdateRow, error)
	TrashedCard(ctx context.Context, cardId int) (*models.CardTrashRow, error)
	RestoreCard(ctx context.Context, cardId int) (*models.CardRestoreRow, error)
	DeleteCardPermanent(ctx context.Context, cardId int) (bool, error)

	RestoreAllCard(ctx context.Context) (bool, error)
	DeleteAllCardPermanent(ctx context.Context) (bool, error)
	ToggleCardStatus(ctx context.Context, request *requests.ToggleCardStatusRequest) (*models.CardAllFieldsRow, error)
	UpdateCreditLimit(ctx context.Context, request *requests.UpdateCreditLimitRequest) (*models.CardAllFieldsRow, error)
	RedeemPoints(ctx context.Context, request *requests.RedeemPointsRequest) (*models.CardAllFieldsRow, error)
	ProcessBillingCycles(ctx context.Context) error
}
