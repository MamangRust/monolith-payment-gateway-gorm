package transactionstatsbycarcache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransactionStatsByCardAmountCache interface {
	GetMonthlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionMonthlyAmountByCardRow, bool)
	SetMonthlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*models.TransactionMonthlyAmountByCardRow)

	GetYearlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionYearlyAmountByCardRow, bool)
	SetYearlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*models.TransactionYearlyAmountByCardRow)
}

type TransactionStatsByCardMethodCache interface {
	GetMonthlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionMonthlyPaymentMethodByCardRow, bool)
	SetMonthlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*models.TransactionMonthlyPaymentMethodByCardRow)

	GetYearlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionYearlyPaymentMethodByCardRow, bool)
	SetYearlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*models.TransactionYearlyPaymentMethodByCardRow)
}

type TransactionStatsByCardStatusCache interface {
	GetMonthTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*models.TransactionMonthlyStatusSuccessByCardRow, bool)
	SetMonthTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber, data []*models.TransactionMonthlyStatusSuccessByCardRow)

	GetYearTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*models.TransactionYearlyStatusSuccessByCardRow, bool)
	SetYearTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber, data []*models.TransactionYearlyStatusSuccessByCardRow)

	GetMonthTransactionStatusFailedByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*models.TransactionMonthlyStatusFailedByCardRow, bool)
	SetMonthTransactionStatusFailedByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber, data []*models.TransactionMonthlyStatusFailedByCardRow)

	GetYearTransactionStatusFailedByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*models.TransactionYearlyStatusFailedByCardRow, bool)
	SetYearTransactionStatusFailedByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber, data []*models.TransactionYearlyStatusFailedByCardRow)
}
