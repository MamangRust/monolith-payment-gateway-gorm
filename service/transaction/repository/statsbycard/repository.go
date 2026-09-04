package transactionbycardrepository

import "gorm.io/gorm"

type TransactionStatsByCardRepository interface {
	TransactionStatsByCardAmountRepository
	TransactionStatsByCardMethodRepository
	TransactionStatsByCardStatusRepository
}

type repository struct {
	TransactionStatsByCardAmountRepository
	TransactionStatsByCardMethodRepository
	TransactionStatsByCardStatusRepository
}

func NewTransactionStatsByCardRepository(db *gorm.DB) TransactionStatsByCardRepository {
	return &repository{
		TransactionStatsByCardAmountRepository: NewTransactionStatsByCardAmountRepository(db),
		TransactionStatsByCardMethodRepository: NewTransactionStatsByCardMethodRepository(db),
		TransactionStatsByCardStatusRepository: NewTransactionStatsByCardStatusRepository(db),
	}
}
