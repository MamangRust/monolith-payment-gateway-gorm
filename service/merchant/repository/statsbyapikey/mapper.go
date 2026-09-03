package merchantstatsapikeyrepository

import "gorm.io/gorm"

type MerchantStatsByApiKeyRepository interface {
	MerchantStatsAmountByApiKeyRepository
	MerchantStatsMethodByApiKeyRepository
	MerchantStatsTotalAmountByApiKeyRepository
}

type repository struct {
	MerchantStatsAmountByApiKeyRepository
	MerchantStatsMethodByApiKeyRepository
	MerchantStatsTotalAmountByApiKeyRepository
}

func NewMerchantStatsByApiKeyRepository(db *gorm.DB) MerchantStatsByApiKeyRepository {
	return &repository{
		MerchantStatsAmountByApiKeyRepository:      NewMerchantStatsAmountByApiKeyRepository(db),
		MerchantStatsMethodByApiKeyRepository:      NewMerchantStatsMethodByApiKeyRepository(db),
		MerchantStatsTotalAmountByApiKeyRepository: NewMerchantStatsTotalAmountByApiKeyRepository(db),
	}
}
