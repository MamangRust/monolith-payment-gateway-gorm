package transferstatscache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransferStatsAmountCache interface {
	GetCachedMonthTransferAmounts(ctx context.Context, year int) ([]*models.TransferMonthlyAmountRow, bool)
	SetCachedMonthTransferAmounts(ctx context.Context, year int, data []*models.TransferMonthlyAmountRow)

	GetCachedYearlyTransferAmounts(ctx context.Context, year int) ([]*models.TransferYearlyAmountRow, bool)
	SetCachedYearlyTransferAmounts(ctx context.Context, year int, data []*models.TransferYearlyAmountRow)
}

type TransferStatsStatusCache interface {
	GetCachedMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*models.TransferMonthlyStatusSuccessRow, bool)
	SetCachedMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer, data []*models.TransferMonthlyStatusSuccessRow)

	GetCachedYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*models.TransferYearlyStatusSuccessRow, bool)
	SetCachedYearlyTransferStatusSuccess(ctx context.Context, year int, data []*models.TransferYearlyStatusSuccessRow)

	GetCachedMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*models.TransferMonthlyStatusFailedRow, bool)
	SetCachedMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer, data []*models.TransferMonthlyStatusFailedRow)

	GetCachedYearlyTransferStatusFailed(ctx context.Context, year int) ([]*models.TransferYearlyStatusFailedRow, bool)
	SetCachedYearlyTransferStatusFailed(ctx context.Context, year int, data []*models.TransferYearlyStatusFailedRow)
}
