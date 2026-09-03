package topupstatsrepository

import "gorm.io/gorm"

type TopupStatsRepository interface {
	TopupStatsAmountRepository
	TopupStatsStatusRepository
	TopupStatsMethodRepository
}

type repository struct {
	TopupStatsAmountRepository
	TopupStatsStatusRepository
	TopupStatsMethodRepository
}

func NewTopupStatsRepository(db *gorm.DB) TopupStatsRepository {
	return &repository{
		TopupStatsAmountRepository: NewTopupStatsAmountRepository(db),
		TopupStatsStatusRepository: NewTopupStatsStatusRepository(db),
		TopupStatsMethodRepository: NewTopupStatsMethodRepository(db),
	}
}
