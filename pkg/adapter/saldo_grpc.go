package adapter

import (
	"context"
	"time"

	pbcard "github.com/MamangRust/monolith-payment-gateway-pb/card"
	pbsaldo "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type SaldoAdapter struct {
	QueryClient   pbsaldo.SaldoQueryServiceClient
	CommandClient pbsaldo.SaldoCommandServiceClient
}

func NewSaldoAdapter(queryClient pbsaldo.SaldoQueryServiceClient, commandClient pbsaldo.SaldoCommandServiceClient) *SaldoAdapter {
	return &SaldoAdapter{
		QueryClient:   queryClient,
		CommandClient: commandClient,
	}
}

func (a *SaldoAdapter) FindByCardNumber(ctx context.Context, card_number string) (*models.Saldo, error) {
	resp, err := a.QueryClient.FindByCardNumber(ctx, &pbcard.FindByCardNumberRequest{
		CardNumber: card_number,
	})
	if err != nil {
		return nil, err
	}

	return MapSaldoResponseToModel(resp.Data), nil
}

func (a *SaldoAdapter) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*models.UpdateSaldoBalanceRow, error) {
	resp, err := a.CommandClient.UpdateSaldoBalance(ctx, &pbsaldo.UpdateSaldoBalanceRequest{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	})
	if err != nil {
		return nil, err
	}

	return &models.UpdateSaldoBalanceRow{
		SaldoID:      resp.Data.SaldoId,
		CardNumber:   resp.Data.CardNumber,
		TotalBalance: resp.Data.TotalBalance,
	}, nil
}

func (a *SaldoAdapter) UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*models.UpdateSaldoWithdrawRow, error) {
	withdrawTime := ""
	if request.WithdrawTime != nil {
		withdrawTime = request.WithdrawTime.Format(time.RFC3339)
	}

	withdrawAmount := int32(0)
	if request.WithdrawAmount != nil {
		withdrawAmount = int32(*request.WithdrawAmount)
	}

	resp, err := a.CommandClient.UpdateSaldoWithdraw(ctx, &pbsaldo.UpdateSaldoWithdrawRequest{
		CardNumber:     request.CardNumber,
		TotalBalance:   int32(request.TotalBalance),
		WithdrawTime:   withdrawTime,
		WithdrawAmount: withdrawAmount,
	})
	if err != nil {
		return nil, err
	}

	return &models.UpdateSaldoWithdrawRow{
		SaldoID:        resp.Data.SaldoId,
		CardNumber:     resp.Data.CardNumber,
		TotalBalance:   resp.Data.TotalBalance,
		WithdrawAmount: &resp.Data.WithdrawAmount,
		WithdrawTime:   parseTimePtr(resp.Data.WithdrawTime),
		CreatedAt:      parseTime(resp.Data.CreatedAt),
		UpdatedAt:      parseTime(resp.Data.UpdatedAt),
	}, nil
}

func MapSaldoResponseToModel(s *pbsaldo.SaldoResponse) *models.Saldo {
	if s == nil {
		return nil
	}

	saldo := &models.Saldo{
		SaldoID:      s.SaldoId,
		CardNumber:   s.CardNumber,
		TotalBalance: s.TotalBalance,
	}

	if s.WithdrawAmount != 0 {
		wa := int32(s.WithdrawAmount)
		saldo.WithdrawAmount = &wa
	}

	saldo.WithdrawTime = parseTimePtr(s.WithdrawTime)
	saldo.CreatedAt = parseTime(s.CreatedAt)
	saldo.UpdatedAt = parseTime(s.UpdatedAt)

	return saldo
}

func parseTimePtr(ts string) *time.Time {
	if ts == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil
	}
	return &t
}
