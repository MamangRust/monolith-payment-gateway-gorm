package topupstatsbycardrepository

import "gorm.io/gorm"

type TopupStatsByCardRepository interface {
	TopupStatsByCardAmountRepository
	TopupStatsByCardMethodRepository
	TopupStatsByCardStatusRepository
}

type repository struct {
	TopupStatsByCardAmountRepository
	TopupStatsByCardMethodRepository
	TopupStatsByCardStatusRepository
}

func NewTopupStatsByCardRepository(db *gorm.DB) TopupStatsByCardRepository {
	return &repository{
		TopupStatsByCardAmountRepository: NewTopupStatsByCardAmountRepository(db),
		TopupStatsByCardMethodRepository: NewTopupStatsByCardMethodRepository(db),
		TopupStatsByCardStatusRepository: NewTopupStatsByCardStatusRepository(db),
	}
}
