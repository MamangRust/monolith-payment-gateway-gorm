package transactionstatsservice

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransactionStatsAmountService interface {
	FindMonthlyAmounts(ctx context.Context, year int) ([]*models.TransactionMonthlyAmountRow, error)
	FindYearlyAmounts(ctx context.Context, year int) ([]*models.TransactionYearlyAmountRow, error)
}

type TransactionStatsMethodService interface {
	FindMonthlyPaymentMethods(ctx context.Context, year int) ([]*models.TransactionMonthlyPaymentMethodRow, error)
	FindYearlyPaymentMethods(ctx context.Context, year int) ([]*models.TransactionYearlyPaymentMethodRow, error)
}

type TransactionStatsStatusService interface {
	FindMonthTransactionStatusSuccess(ctx context.Context, req *requests.MonthStatusTransaction) ([]*models.TransactionMonthlyStatusSuccessRow, error)
	FindYearlyTransactionStatusSuccess(ctx context.Context, year int) ([]*models.TransactionYearlyStatusSuccessRow, error)
	FindMonthTransactionStatusFailed(ctx context.Context, req *requests.MonthStatusTransaction) ([]*models.TransactionMonthlyStatusFailedRow, error)
	FindYearlyTransactionStatusFailed(ctx context.Context, year int) ([]*models.TransactionYearlyStatusFailedRow, error)
}
