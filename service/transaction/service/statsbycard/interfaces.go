package transactionstatsbycardservice

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransactionStatsByCardAmountService interface {
	FindMonthlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionMonthlyAmountByCardRow, error)
	FindYearlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionYearlyAmountByCardRow, error)
}

type TransactionStatsByCardMethodService interface {
	FindMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionMonthlyPaymentMethodByCardRow, error)
	FindYearlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionYearlyPaymentMethodByCardRow, error)
}

type TransactionStatsByCardStatusService interface {
	FindMonthTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*models.TransactionMonthlyStatusSuccessByCardRow, error)
	FindYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*models.TransactionYearlyStatusSuccessByCardRow, error)
	FindMonthTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*models.TransactionMonthlyStatusFailedByCardRow, error)
	FindYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*models.TransactionYearlyStatusFailedByCardRow, error)
}
