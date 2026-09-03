package topupstatsbycardcache

import (
	"context"

	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TopupStatsStatusByCardCache interface {
	GetMonthTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*models.TopupMonthlyStatusRow, bool)
	SetMonthTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber, data []*models.TopupMonthlyStatusRow)

	GetYearlyTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*models.TopupYearlyStatusRow, bool)
	SetYearlyTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber, data []*models.TopupYearlyStatusRow)

	GetMonthTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*models.TopupMonthlyStatusRow, bool)
	SetMonthTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber, data []*models.TopupMonthlyStatusRow)

	GetYearlyTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*models.TopupYearlyStatusRow, bool)
	SetYearlyTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber, data []*models.TopupYearlyStatusRow)
}

type TopupStatsMethodByCardCache interface {
	GetMonthlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupMonthlyMethodByCardRow, bool)
	SetMonthlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*models.TopupMonthlyMethodByCardRow)

	GetYearlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupYearlyMethodByCardRow, bool)
	SetYearlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*models.TopupYearlyMethodByCardRow)
}

type TopupStatsAmountByCardCache interface {
	GetMonthlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupMonthlyAmountRow, bool)
	SetMonthlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*models.TopupMonthlyAmountRow)

	GetYearlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*models.TopupYearlyAmountRow, bool)
	SetYearlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*models.TopupYearlyAmountRow)
}
