package withdrawstatsrepository

import "gorm.io/gorm"

type WithdrawStatsRepository interface {
	WithdrawStatsAmountRepository
	WithdrawStatsStatusRepository
}

type repositories struct {
	WithdrawStatsAmountRepository
	WithdrawStatsStatusRepository
}

func NewWithdrawStatsRepository(db *gorm.DB) WithdrawStatsRepository {
	return &repositories{
		WithdrawStatsAmountRepository: NewWithdrawStatsAmountRepository(db),
		WithdrawStatsStatusRepository: NewWithdrawStatsStatusRepository(db),
	}
}
