package transferstatsbycardrepository

import "gorm.io/gorm"

type TransferStatsByCardRepository interface {
	TransferStatsByCardAmountSenderRepository
	TransferStatsByCardAmountReceiverRepository
	TransferStatsByCardStatusRepository
}

type repository struct {
	TransferStatsByCardAmountSenderRepository
	TransferStatsByCardAmountReceiverRepository
	TransferStatsByCardStatusRepository
}

func NewTransferStatsByCardRepository(db *gorm.DB) TransferStatsByCardRepository {
	return &repository{
		TransferStatsByCardAmountSenderRepository:   NewTransferStatsByCardAmountSenderRepository(db),
		TransferStatsByCardAmountReceiverRepository: NewTransferStatsByCardAmountReceiverRepository(db),
		TransferStatsByCardStatusRepository:         NewTransferStatsByCardStatusRepository(db),
	}
}
