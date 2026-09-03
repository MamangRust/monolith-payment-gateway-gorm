package transactionbycardrepository

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransactionStatsByCardStatusRepository interface {
	GetMonthTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*models.TransactionMonthlyStatusSuccessByCardRow, error)
	GetYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*models.TransactionYearlyStatusSuccessByCardRow, error)
	GetMonthTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*models.TransactionMonthlyStatusFailedByCardRow, error)
	GetYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*models.TransactionYearlyStatusFailedByCardRow, error)
}

type TransactionStatsByCardMethodRepository interface {
	GetMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionMonthlyPaymentMethodByCardRow, error)
	GetYearlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionYearlyPaymentMethodByCardRow, error)
}

type TransactionStatsByCardAmountRepository interface {
	GetMonthlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionMonthlyAmountByCardRow, error)
	GetYearlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*models.TransactionYearlyAmountByCardRow, error)
}
