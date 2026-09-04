package models

// TopupMonthlyMethodByCardRow - monthly topup method stats by card number
type TopupMonthlyMethodByCardRow struct {
	Month       string `gorm:"column:month" json:"month"`
	TopupMethod string `gorm:"column:topup_method" json:"topup_method"`
	TotalTopups int32  `gorm:"column:total_topups" json:"total_topups"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TopupYearlyMethodByCardRow - yearly topup method stats by card number
type TopupYearlyMethodByCardRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TopupMethod string `gorm:"column:topup_method" json:"topup_method"`
	TotalTopups int64  `gorm:"column:total_topups" json:"total_topups"`
	TotalAmount int64  `gorm:"column:total_amount" json:"total_amount"`
}
