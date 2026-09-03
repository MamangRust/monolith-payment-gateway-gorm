package transactionstatscache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransactionStatsAmountCache interface {
	GetMonthlyAmountsCache(ctx context.Context, year int) ([]*models.TransactionMonthlyAmountRow, bool)
	SetMonthlyAmountsCache(ctx context.Context, year int, data []*models.TransactionMonthlyAmountRow)

	GetYearlyAmountsCache(ctx context.Context, year int) ([]*models.TransactionYearlyAmountRow, bool)
	SetYearlyAmountsCache(ctx context.Context, year int, data []*models.TransactionYearlyAmountRow)
}

type TransactionStatsMethodCache interface {
	GetMonthlyPaymentMethodsCache(ctx context.Context, year int) ([]*models.TransactionMonthlyPaymentMethodRow, bool)
	SetMonthlyPaymentMethodsCache(ctx context.Context, year int, data []*models.TransactionMonthlyPaymentMethodRow)

	GetYearlyPaymentMethodsCache(ctx context.Context, year int) ([]*models.TransactionYearlyPaymentMethodRow, bool)
	SetYearlyPaymentMethodsCache(ctx context.Context, year int, data []*models.TransactionYearlyPaymentMethodRow)
}

type TransactionStatsStatusCache interface {
	GetMonthTransactionStatusSuccessCache(ctx context.Context, req *requests.MonthStatusTransaction) ([]*models.TransactionMonthlyStatusSuccessRow, bool)
	SetMonthTransactionStatusSuccessCache(ctx context.Context, req *requests.MonthStatusTransaction, data []*models.TransactionMonthlyStatusSuccessRow)

	GetYearTransactionStatusSuccessCache(ctx context.Context, year int) ([]*models.TransactionYearlyStatusSuccessRow, bool)
	SetYearTransactionStatusSuccessCache(ctx context.Context, year int, data []*models.TransactionYearlyStatusSuccessRow)

	GetMonthTransactionStatusFailedCache(ctx context.Context, req *requests.MonthStatusTransaction) ([]*models.TransactionMonthlyStatusFailedRow, bool)
	SetMonthTransactionStatusFailedCache(ctx context.Context, req *requests.MonthStatusTransaction, data []*models.TransactionMonthlyStatusFailedRow)

	GetYearTransactionStatusFailedCache(ctx context.Context, year int) ([]*models.TransactionYearlyStatusFailedRow, bool)
	SetYearTransactionStatusFailedCache(ctx context.Context, year int, data []*models.TransactionYearlyStatusFailedRow)
}
