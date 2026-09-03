package merchantstatsmerchantrepository

import "gorm.io/gorm"

type MerchantStatsByMerchantRepository interface {
	MerchantStatsAmountByMerchantRepository
	MerchantStatsMethodByMerchantRepository
	MerchantStatsTotalAmountByMerchantRepository
}

type repository struct {
	MerchantStatsAmountByMerchantRepository
	MerchantStatsMethodByMerchantRepository
	MerchantStatsTotalAmountByMerchantRepository
}

func NewMerchantStatsByMerchantRepository(db *gorm.DB) MerchantStatsByMerchantRepository {
	return &repository{
		MerchantStatsAmountByMerchantRepository:      NewMerchantStatsAmountByMerchantRepository(db),
		MerchantStatsMethodByMerchantRepository:      NewMerchantStatsMethodByMerchantRepository(db),
		MerchantStatsTotalAmountByMerchantRepository: NewMerchantStatsTotalAmountByMerchantRepository(db),
	}
}
