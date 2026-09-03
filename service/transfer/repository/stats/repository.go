package transferstatsrepository

import "gorm.io/gorm"

type TransferStatsRepository interface {
	TransferStatsAmountRepository
	TransferStatsStatusRepository
}

type repositories struct {
	TransferStatsAmountRepository
	TransferStatsStatusRepository
}

func NewTransferStatsRepository(db *gorm.DB) TransferStatsRepository {
	return &repositories{
		TransferStatsAmountRepository: NewTransferStatsAmountRepository(db),
		TransferStatsStatusRepository: NewTransferStatsStatusRepository(db),
	}
}
