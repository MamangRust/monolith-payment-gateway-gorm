package transferstatsservice

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransferStatsStatusService interface {
	FindMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*models.TransferMonthlyStatusSuccessRow, error)
	FindYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*models.TransferYearlyStatusSuccessRow, error)
	FindMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*models.TransferMonthlyStatusFailedRow, error)
	FindYearlyTransferStatusFailed(ctx context.Context, year int) ([]*models.TransferYearlyStatusFailedRow, error)
}

type TransferStatsAmountService interface {
	FindMonthlyTransferAmounts(ctx context.Context, year int) ([]*models.TransferMonthlyAmountRow, error)
	FindYearlyTransferAmounts(ctx context.Context, year int) ([]*models.TransferYearlyAmountRow, error)
}
