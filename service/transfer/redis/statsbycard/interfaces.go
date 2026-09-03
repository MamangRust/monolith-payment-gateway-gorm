package transferstatsbycardcache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransferStatsByCardAmountCache interface {
	GetMonthlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferMonthlyAmountBySenderCardRow, bool)
	SetMonthlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*models.TransferMonthlyAmountBySenderCardRow)

	GetMonthlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferMonthlyAmountByReceiverCardRow, bool)
	SetMonthlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*models.TransferMonthlyAmountByReceiverCardRow)

	GetYearlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferYearlyAmountBySenderCardRow, bool)
	SetYearlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*models.TransferYearlyAmountBySenderCardRow)

	GetYearlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferYearlyAmountByReceiverCardRow, bool)
	SetYearlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*models.TransferYearlyAmountByReceiverCardRow)
}

type TransferStatsByCardStatusCache interface {
	GetMonthTransferStatusSuccessByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*models.TransferMonthlyStatusSuccessByCardRow, bool)
	SetMonthTransferStatusSuccessByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber, data []*models.TransferMonthlyStatusSuccessByCardRow)

	GetYearlyTransferStatusSuccessByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*models.TransferYearlyStatusSuccessByCardRow, bool)
	SetYearlyTransferStatusSuccessByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber, data []*models.TransferYearlyStatusSuccessByCardRow)

	GetMonthTransferStatusFailedByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*models.TransferMonthlyStatusFailedByCardRow, bool)
	SetMonthTransferStatusFailedByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber, data []*models.TransferMonthlyStatusFailedByCardRow)

	GetYearlyTransferStatusFailedByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*models.TransferYearlyStatusFailedByCardRow, bool)
	SetYearlyTransferStatusFailedByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber, data []*models.TransferYearlyStatusFailedByCardRow)
}
