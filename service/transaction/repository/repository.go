package repository

import (
	"gorm.io/gorm"

	transactionstatsrepository "github.com/MamangRust/monolith-payment-gateway-transaction/repository/stats"
	transactionbycardrepository "github.com/MamangRust/monolith-payment-gateway-transaction/repository/statsbycard"
)

type Repositories interface {
	SaldoRepository
	MerchantRepository
	CardRepository
	TransactionQueryRepository
	TransactionCommandRepository
	transactionstatsrepository.TransactionStatsRepository
	transactionbycardrepository.TransactionStatsByCardRepository
}

type repositories struct {
	SaldoRepository
	MerchantRepository
	CardRepository
	TransactionQueryRepository
	TransactionCommandRepository
	transactionstatsrepository.TransactionStatsRepository
	transactionbycardrepository.TransactionStatsByCardRepository
}

func NewRepositories(
	db *gorm.DB,
	saldo SaldoRepository,
	card CardRepository,
	merchant MerchantRepository,
) Repositories {
	return &repositories{
		SaldoRepository:                  saldo,
		MerchantRepository:               merchant,
		CardRepository:                   card,
		TransactionQueryRepository:       NewTransactionQueryRepository(db),
		TransactionCommandRepository:     NewTransactionCommandRepository(db),
		TransactionStatsRepository:       transactionstatsrepository.NewTransactionStatsRepository(db),
		TransactionStatsByCardRepository: transactionbycardrepository.NewTransactionStatsByCardRepository(db),
	}
}
