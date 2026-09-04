package models

// Withdraw stats row types

// WithdrawMonthlyAmountRow - monthly withdraw amount stats
type WithdrawMonthlyAmountRow struct {
	Month               string `gorm:"column:month" json:"month"`
	TotalWithdrawAmount int32  `gorm:"column:total_withdraw_amount" json:"total_withdraw_amount"`
}

// WithdrawYearlyAmountRow - yearly withdraw amount stats
type WithdrawYearlyAmountRow struct {
	Year                string `gorm:"column:year" json:"year"`
	TotalWithdrawAmount int64  `gorm:"column:total_withdraw_amount" json:"total_withdraw_amount"`
}

// WithdrawMonthlyStatusSuccessRow - monthly status success stats
type WithdrawMonthlyStatusSuccessRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalSuccess int64  `gorm:"column:total_success" json:"total_success"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// WithdrawYearlyStatusSuccessRow - yearly status success stats
type WithdrawYearlyStatusSuccessRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalSuccess int32  `gorm:"column:total_success" json:"total_success"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// WithdrawMonthlyStatusFailedRow - monthly status failed stats
type WithdrawMonthlyStatusFailedRow struct {
	Year        string `gorm:"column:year" json:"year"`
	Month       string `gorm:"column:month" json:"month"`
	TotalFailed int64  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// WithdrawYearlyStatusFailedRow - yearly status failed stats
type WithdrawYearlyStatusFailedRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalFailed int32  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// WithdrawMonthlyStatusSuccessByCardRow - monthly status success by card
type WithdrawMonthlyStatusSuccessByCardRow = WithdrawMonthlyStatusSuccessRow

// WithdrawYearlyStatusSuccessByCardRow - yearly status success by card
type WithdrawYearlyStatusSuccessByCardRow = WithdrawYearlyStatusSuccessRow

// WithdrawMonthlyStatusFailedByCardRow - monthly status failed by card
type WithdrawMonthlyStatusFailedByCardRow = WithdrawMonthlyStatusFailedRow

// WithdrawYearlyStatusFailedByCardRow - yearly status failed by card
type WithdrawYearlyStatusFailedByCardRow = WithdrawYearlyStatusFailedRow

// WithdrawMonthlyAmountByCardRow - monthly amount by card
type WithdrawMonthlyAmountByCardRow = WithdrawMonthlyAmountRow

// WithdrawYearlyAmountByCardRow - yearly amount by card
type WithdrawYearlyAmountByCardRow = WithdrawYearlyAmountRow
