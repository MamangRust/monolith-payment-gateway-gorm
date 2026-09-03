package transactionstatsrepository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransactionStatsStatusRepository interface {
	GetMonthTransactionStatusSuccess(ctx context.Context, req *requests.MonthStatusTransaction) ([]*models.TransactionMonthlyStatusSuccessRow, error)
	GetYearlyTransactionStatusSuccess(ctx context.Context, year int) ([]*models.TransactionYearlyStatusSuccessRow, error)
	GetMonthTransactionStatusFailed(ctx context.Context, req *requests.MonthStatusTransaction) ([]*models.TransactionMonthlyStatusFailedRow, error)
	GetYearlyTransactionStatusFailed(ctx context.Context, year int) ([]*models.TransactionYearlyStatusFailedRow, error)
}

type TransactionStatsMethodRepository interface {
	GetMonthlyPaymentMethods(ctx context.Context, year int) ([]*models.TransactionMonthlyPaymentMethodRow, error)
	GetYearlyPaymentMethods(ctx context.Context, year int) ([]*models.TransactionYearlyPaymentMethodRow, error)
}

type TransactionStatsAmountRepository interface {
	GetMonthlyAmounts(ctx context.Context, year int) ([]*models.TransactionMonthlyAmountRow, error)
	GetYearlyAmounts(ctx context.Context, year int) ([]*models.TransactionYearlyAmountRow, error)
}
