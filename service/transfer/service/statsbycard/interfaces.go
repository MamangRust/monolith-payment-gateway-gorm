package transferstatsbycardservice

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransferStatsByCardAmountService interface {
	FindMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferMonthlyAmountBySenderCardRow, error)
	FindYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferYearlyAmountBySenderCardRow, error)
	FindMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferMonthlyAmountByReceiverCardRow, error)
	FindYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferYearlyAmountByReceiverCardRow, error)
}

type TransferStatsByCardStatusService interface {
	FindMonthTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*models.TransferMonthlyStatusSuccessByCardRow, error)
	FindYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*models.TransferYearlyStatusSuccessByCardRow, error)
	FindMonthTransferStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*models.TransferMonthlyStatusFailedByCardRow, error)
	FindYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*models.TransferYearlyStatusFailedByCardRow, error)
}
