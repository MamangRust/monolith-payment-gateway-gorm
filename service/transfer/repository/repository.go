package repository

import (
	"gorm.io/gorm"

	transferstatsrepository "github.com/MamangRust/monolith-payment-gateway-transfer/repository/stats"
	transferstatsbycardrepository "github.com/MamangRust/monolith-payment-gateway-transfer/repository/statsbycard"
)

type Repositories interface {
	SaldoRepository
	TransferQueryRepository
	TransferCommandRepository
	CardRepository
	transferstatsrepository.TransferStatsRepository
	transferstatsbycardrepository.TransferStatsByCardRepository
}

type repositories struct {
	SaldoRepository
	TransferQueryRepository
	TransferCommandRepository
	CardRepository
	transferstatsrepository.TransferStatsRepository
	transferstatsbycardrepository.TransferStatsByCardRepository
}

func NewRepositories(
	db *gorm.DB,
	saldo SaldoRepository,
	card CardRepository,
) Repositories {
	return &repositories{
		SaldoRepository:               saldo,
		TransferQueryRepository:       NewTransferQueryRepository(db),
		TransferCommandRepository:     NewTransferCommandRepository(db),
		TransferStatsRepository:       transferstatsrepository.NewTransferStatsRepository(db),
		TransferStatsByCardRepository: transferstatsbycardrepository.NewTransferStatsByCardRepository(db),
		CardRepository:                card,
	}
}
