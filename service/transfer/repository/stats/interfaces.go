package transferstatsrepository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransferStatsAmountRepository interface {
	GetMonthlyTransferAmounts(ctx context.Context, year int) ([]*models.TransferMonthlyAmountRow, error)
	GetYearlyTransferAmounts(ctx context.Context, year int) ([]*models.TransferYearlyAmountRow, error)
}

type TransferStatsStatusRepository interface {
	GetMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*models.TransferMonthlyStatusSuccessRow, error)
	GetYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*models.TransferYearlyStatusSuccessRow, error)
	GetMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*models.TransferMonthlyStatusFailedRow, error)
	GetYearlyTransferStatusFailed(ctx context.Context, year int) ([]*models.TransferYearlyStatusFailedRow, error)
}
