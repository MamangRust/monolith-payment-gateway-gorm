package repository

import (
	merchantstatsrepository "github.com/MamangRust/monolith-payment-gateway-merchant/repository/stats"
	merchantstatsapikeyrepository "github.com/MamangRust/monolith-payment-gateway-merchant/repository/statsbyapikey"
	merchantstatsmerchantrepository "github.com/MamangRust/monolith-payment-gateway-merchant/repository/statsbymerchant"
	"gorm.io/gorm"
)

type Repositories interface {
	MerchantQueryRepository
	MerchantCommandRepository
	MerchantDocumentQueryRepository
	MerchantDocumentCommandRepository
	MerchantTransactionRepository
	merchantstatsrepository.MerchantStatsRepository
	merchantstatsapikeyrepository.MerchantStatsByApiKeyRepository
	merchantstatsmerchantrepository.MerchantStatsByMerchantRepository
	UserRepository
}

type repositories struct {
	MerchantQueryRepository
	MerchantCommandRepository
	MerchantDocumentQueryRepository
	MerchantDocumentCommandRepository
	MerchantTransactionRepository
	merchantstatsrepository.MerchantStatsRepository
	merchantstatsapikeyrepository.MerchantStatsByApiKeyRepository
	merchantstatsmerchantrepository.MerchantStatsByMerchantRepository
	UserRepository
}

func NewRepositories(db *gorm.DB) Repositories {
	return &repositories{
		NewMerchantQueryRepository(db),
		NewMerchantCommandRepository(db),
		NewMerchantDocumentQueryRepository(db),
		NewMerchantDocumentCommandRepository(db),
		NewMerchantTransactionRepository(db),
		merchantstatsrepository.NewMerchantStatsRepository(db),
		merchantstatsapikeyrepository.NewMerchantStatsByApiKeyRepository(db),
		merchantstatsmerchantrepository.NewMerchantStatsByMerchantRepository(db),
		NewUserRepository(db),
	}
}
