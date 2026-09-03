package merchantstatsrepository

import "gorm.io/gorm"

type MerchantStatsRepository interface {
	MerchantStatsAmountRepository
	MerchantStatsMethodRepository
	MerchantStatsTotalAmountRepository
}

type repository struct {
	MerchantStatsAmountRepository
	MerchantStatsMethodRepository
	MerchantStatsTotalAmountRepository
}

func NewMerchantStatsRepository(db *gorm.DB) MerchantStatsRepository {
	return &repository{
		MerchantStatsAmountRepository:      NewMerchantStatsAmountRepository(db),
		MerchantStatsMethodRepository:      NewMerchantStatsMethodRepository(db),
		MerchantStatsTotalAmountRepository: NewMerchantStatsTotalAmountRepository(db),
	}
}
