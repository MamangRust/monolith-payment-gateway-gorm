package repositorydashboard

import "gorm.io/gorm"

type CardDashboardRepository interface {
	CardDashboardBalanceRepository
	CardDashboardTopupRepository
	CardDashboardTransactionRepository
	CardDashboardTransferRepository
	CardDashboardWithdrawRepository
}

type repository struct {
	CardDashboardBalanceRepository
	CardDashboardTopupRepository
	CardDashboardTransactionRepository
	CardDashboardTransferRepository
	CardDashboardWithdrawRepository
}

func NewCardDashboardRepository(db *gorm.DB) CardDashboardRepository {
	return &repository{
		CardDashboardBalanceRepository:     NewCardDashboardBalanceRepository(db),
		CardDashboardTopupRepository:       NewCardDashboardTopupRepository(db),
		CardDashboardTransactionRepository: NewCardDashboardTransactionRepository(db),
		CardDashboardTransferRepository:    NewCardDashboardTransferRepository(db),
		CardDashboardWithdrawRepository:    NewCardDashboardWithdrawRepository(db),
	}
}
