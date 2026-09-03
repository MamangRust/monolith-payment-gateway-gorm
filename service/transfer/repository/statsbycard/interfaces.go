package transferstatsbycardrepository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransferStatsByCardAmountSenderRepository interface {
	GetMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferMonthlyAmountBySenderCardRow, error)
	GetYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferYearlyAmountBySenderCardRow, error)
}

type TransferStatsByCardAmountReceiverRepository interface {
	GetMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferMonthlyAmountByReceiverCardRow, error)
	GetYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*models.TransferYearlyAmountByReceiverCardRow, error)
}

type TransferStatsByCardStatusRepository interface {
	GetMonthTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*models.TransferMonthlyStatusSuccessByCardRow, error)
	GetYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*models.TransferYearlyStatusSuccessByCardRow, error)
	GetMonthTransferStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*models.TransferMonthlyStatusFailedByCardRow, error)
	GetYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*models.TransferYearlyStatusFailedByCardRow, error)
}
