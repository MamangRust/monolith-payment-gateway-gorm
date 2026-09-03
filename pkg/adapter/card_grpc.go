package adapter

import (
	"context"
	"time"

	pbcard "github.com/MamangRust/monolith-payment-gateway-pb/card"
	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CardAdapter struct {
	QueryClient   pbcard.CardQueryServiceClient
	CommandClient pbcard.CardCommandServiceClient
}

func NewCardAdapter(queryClient pbcard.CardQueryServiceClient, commandClient pbcard.CardCommandServiceClient) *CardAdapter {
	return &CardAdapter{
		QueryClient:   queryClient,
		CommandClient: commandClient,
	}
}

func (a *CardAdapter) FindCardByUserId(ctx context.Context, user_id int) (*models.CardAllFieldsRow, error) {
	resp, err := a.QueryClient.FindByUserIdCard(ctx, &pbcard.FindByUserIdCardRequest{
		UserId: int32(user_id),
	})
	if err != nil {
		return nil, err
	}

	return &models.CardAllFieldsRow{
		CardID:       resp.Data.Id,
		UserID:       resp.Data.UserId,
		CardNumber:   resp.Data.CardNumber,
		CardType:     resp.Data.CardType,
		ExpireDate:   parseTime(resp.Data.ExpireDate),
		Cvv:          resp.Data.Cvv,
		CardProvider: resp.Data.CardProvider,
		Status:       resp.Data.Status,
	}, nil
}

func (a *CardAdapter) FindUserCardByCardNumber(ctx context.Context, card_number string) (*models.CardByEmailRow, error) {
	resp, err := a.QueryClient.FindUserCardByCardNumber(ctx, &pbcard.FindByCardNumberRequest{
		CardNumber: card_number,
	})
	if err != nil {
		return nil, err
	}

	return &models.CardByEmailRow{
		CardID:       resp.Id,
		UserID:       resp.UserId,
		CardNumber:   resp.CardNumber,
		CardType:     resp.CardType,
		ExpireDate:   parseTime(resp.ExpireDate),
		Cvv:          resp.Cvv,
		CardProvider: resp.CardProvider,
		Email:        resp.Email,
	}, nil
}

func (a *CardAdapter) FindCardByCardNumber(ctx context.Context, card_number string) (*models.CardAllFieldsRow, error) {
	resp, err := a.QueryClient.FindByCardNumber(ctx, &pbcard.FindByCardNumberRequest{
		CardNumber: card_number,
	})
	if err != nil {
		return nil, err
	}

	return &models.CardAllFieldsRow{
		CardID:       resp.Data.Id,
		UserID:       resp.Data.UserId,
		CardNumber:   resp.Data.CardNumber,
		CardType:     resp.Data.CardType,
		ExpireDate:   parseTime(resp.Data.ExpireDate),
		Cvv:          resp.Data.Cvv,
		CardProvider: resp.Data.CardProvider,
		Status:       resp.Data.Status,
	}, nil
}

func (a *CardAdapter) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*models.CardUpdateRow, error) {
	resp, err := a.CommandClient.UpdateCard(ctx, &pbcard.UpdateCardRequest{
		CardId:       int32(request.CardID),
		UserId:       int32(request.UserID),
		CardType:     request.CardType,
		ExpireDate:   timestamppb.New(request.ExpireDate),
		Cvv:          request.CVV,
		CardProvider: request.CardProvider,
	})

	if err != nil {
		return nil, err
	}

	return &models.CardUpdateRow{
		CardID:       resp.Data.Id,
		UserID:       resp.Data.UserId,
		CardNumber:   resp.Data.CardNumber,
		CardType:     resp.Data.CardType,
		ExpireDate:   parseTime(resp.Data.ExpireDate),
		Cvv:          resp.Data.Cvv,
		CardProvider: resp.Data.CardProvider,
	}, nil
}

func (a *CardAdapter) AddRewardPoints(ctx context.Context, cardID int, points int) (*models.AddRewardPointsRow, error) {
	resp, err := a.CommandClient.UpdateCreditLimit(ctx, &pbcard.UpdateCreditLimitRequest{
		CardId:      int32(cardID),
		CreditLimit: int32(points),
	})

	if err != nil {
		return nil, err
	}

	return &models.AddRewardPointsRow{
		CardID:             resp.Data.Id,
		UserID:             resp.Data.UserId,
		CardNumber:         resp.Data.CardNumber,
		CardType:           resp.Data.CardType,
		ExpireDate:         parseTime(resp.Data.ExpireDate),
		Cvv:                resp.Data.Cvv,
		CardProvider:       resp.Data.CardProvider,
		Status:             resp.Data.Status,
		CreditLimit:        resp.Data.CreditLimit,
		OutstandingBalance: resp.Data.OutstandingBalance,
		RewardPoints:       resp.Data.RewardPoints,
	}, nil
}

func (a *CardAdapter) UpdateCardOutstandingBalance(ctx context.Context, cardID int, outstandingBalance int) (*models.UpdateOutstandingBalanceRow, error) {
	return &models.UpdateOutstandingBalanceRow{}, nil
}

func (a *CardAdapter) UpdateCardStatus(ctx context.Context, cardID int, status string) (*models.UpdateCardStatusRow, error) {
	resp, err := a.CommandClient.ToggleCardStatus(ctx, &pbcard.ToggleCardStatusRequest{
		CardId: int32(cardID),
	})

	if err != nil {
		return nil, err
	}

	return &models.UpdateCardStatusRow{
		CardID:             resp.Data.Id,
		UserID:             resp.Data.UserId,
		CardNumber:         resp.Data.CardNumber,
		CardType:           resp.Data.CardType,
		ExpireDate:         parseTime(resp.Data.ExpireDate),
		Cvv:                resp.Data.Cvv,
		CardProvider:       resp.Data.CardProvider,
		Status:             resp.Data.Status,
		CreditLimit:        resp.Data.CreditLimit,
		OutstandingBalance: resp.Data.OutstandingBalance,
		RewardPoints:       resp.Data.RewardPoints,
	}, nil
}

func parseTime(ts string) time.Time {
	if ts == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return time.Time{}
	}
	return t
}
