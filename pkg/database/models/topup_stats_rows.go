package models

// Topup stats row types

// TopupMonthlyAmountRow - monthly topup amount stats
type TopupMonthlyAmountRow struct {
	Month       string `gorm:"column:month" json:"month"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TopupYearlyAmountRow - yearly topup amount stats
type TopupYearlyAmountRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalAmount int64  `gorm:"column:total_amount" json:"total_amount"`
}

// TopupMonthlyMethodRow - monthly topup by method
type TopupMonthlyMethodRow struct {
	Month       string `gorm:"column:month" json:"month"`
	TopupMethod string `gorm:"column:topup_method" json:"topup_method"`
	TotalTopups int32  `gorm:"column:total_topups" json:"total_topups"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TopupYearlyMethodRow - yearly topup by method
type TopupYearlyMethodRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TopupMethod string `gorm:"column:topup_method" json:"topup_method"`
	TotalTopups int64  `gorm:"column:total_topups" json:"total_topups"`
	TotalAmount int64  `gorm:"column:total_amount" json:"total_amount"`
}

// TopupMonthlyStatusRow - monthly status stats
type TopupMonthlyStatusRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalSuccess int64  `gorm:"column:total_success" json:"total_success"`
	TotalFailed  int64  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TopupYearlyStatusRow - yearly status stats
type TopupYearlyStatusRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalSuccess int32  `gorm:"column:total_success" json:"total_success"`
	TotalFailed  int32  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}
