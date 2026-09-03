package adapter

import (
	"context"

	pbmerchant "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
)

type MerchantAdapter struct {
	QueryClient pbmerchant.MerchantQueryServiceClient
}

func NewMerchantAdapter(queryClient pbmerchant.MerchantQueryServiceClient) *MerchantAdapter {
	return &MerchantAdapter{
		QueryClient: queryClient,
	}
}

func (a *MerchantAdapter) FindByApiKey(ctx context.Context, api_key string) (*models.MerchantAllFieldsRow, error) {
	resp, err := a.QueryClient.FindByApiKey(ctx, &pbmerchant.FindByApiKeyRequest{
		ApiKey: api_key,
	})
	if err != nil {
		return nil, err
	}

	return &models.MerchantAllFieldsRow{
		MerchantID: resp.Data.Id,
		Name:       resp.Data.Name,
		ApiKey:     resp.Data.ApiKey,
		UserID:     resp.Data.UserId,
		Status:     resp.Data.Status,
	}, nil
}
