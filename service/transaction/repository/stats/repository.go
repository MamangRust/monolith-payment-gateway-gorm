package transactionstatsrepository

import "gorm.io/gorm"

type TransactionStatsRepository interface {
	TransactionStatsAmountRepository
	TransactionStatsMethodRepository
	TransactionStatsStatusRepository
}

type repository struct {
	TransactionStatsAmountRepository
	TransactionStatsMethodRepository
	TransactionStatsStatusRepository
}

func NewTransactionStatsRepository(db *gorm.DB) TransactionStatsRepository {
	return &repository{
		TransactionStatsAmountRepository: NewTransactionStatsAmountRepository(db),
		TransactionStatsMethodRepository: NewTransactionStatsMethodRepository(db),
		TransactionStatsStatusRepository: NewTransactionStatsStatusRepository(db),
	}
}
