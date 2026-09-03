package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/card.go
type CardCommandRepository interface {
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
	ProcessBillingCycles(ctx context.Context, billingCycleDay int) ([]*models.BillingCycle, error)
}

type CardQueryRepository interface {
	FindAllCards(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListWithDeletedRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListWithDeletedRow, error)
	FindById(ctx context.Context, card_id int) (*models.CardAllFieldsRow, error)
	FindCardByUserId(ctx context.Context, user_id int) (*models.CardAllFieldsRow, error)
	FindCardByCardNumber(ctx context.Context, card_number string) (*models.CardAllFieldsRow, error)
	FindUserCardByCardNumber(ctx context.Context, card_number string) (*models.CardByEmailRow, error)
}

type UserRepository interface {
	FindById(ctx context.Context, user_id int) (*models.UserByIDRow, error)
}

// CardAuthTransactionRepository handles auth transaction persistence.
type CardAuthTransactionRepository interface {
	InsertPending(ctx context.Context, req *requests.AuthorizeCardRequest) (*models.CardAuthTransaction, error)
	Approve(ctx context.Context, txnID string) (*models.CardAuthTransaction, error)
	Decline(ctx context.Context, txnID string) (*models.CardAuthTransaction, error)
	Reverse(ctx context.Context, txnID string) (*models.CardAuthTransaction, error)
	FindByIdempotencyKey(ctx context.Context, key string) (*models.CardAuthTransaction, error)
	FindByTxnID(ctx context.Context, txnID string) (*models.CardAuthTransaction, error)
	FindByCardNumber(ctx context.Context, cardNumber string, page, pageSize int) ([]*models.GetAuthTxnByCardNumberRow, error)
	CountRecentByCardNumber(ctx context.Context, cardNumber string, since time.Time) (int, error)
	UpdateRiskScore(ctx context.Context, txnID string, score int) error
}

// CardPaymentResult contains a payment and whether it was returned from an
// idempotent replay. Replays must not publish duplicate side-effects.
type CardPaymentResult struct {
	Payment  *models.CardPayment
	Replayed bool
}

// CardPaymentRepository handles card payment persistence.
type CardPaymentRepository interface {
	PostPayment(ctx context.Context, req *requests.PostPaymentRequest) (*models.CardPayment, error)
	PostPaymentIdempotent(ctx context.Context, req *requests.PostPaymentRequest, expectedBillingStatus, targetBillingStatus string) (*CardPaymentResult, error)
	GetPaymentByReferenceID(ctx context.Context, referenceID string) (*models.CardPayment, error)
	GetPaymentHistory(ctx context.Context, cardNumber string, page, pageSize int) ([]*models.CardPaymentListRow, error)
	CountPayments(ctx context.Context, cardNumber string) (int, error)
}

// CardRewardRepository handles reward points persistence.
type CardRewardRepository interface {
	EarnRewards(ctx context.Context, req *requests.EarnRewardsRequest) (*models.CardReward, error)
	GetBalance(ctx context.Context, cardNumber string) (int64, error)
	GetHistory(ctx context.Context, cardNumber string) ([]*models.CardReward, error)
	RedeemRewards(ctx context.Context, cardNumber string, points int64) (int64, error)
}

// BillingCycleRepository handles billing cycle persistence.
type BillingCycleRepository interface {
	GetBillingCyclesByCardNumber(ctx context.Context, cardNumber string) ([]*models.BillingCycle, error)
	GetBillingCycleByID(ctx context.Context, billingID int) (*models.BillingCycle, error)
}
