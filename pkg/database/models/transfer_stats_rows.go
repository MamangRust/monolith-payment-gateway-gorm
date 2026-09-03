package models

// Transfer stats row types

// TransferMonthlyAmountRow - monthly transfer amount stats
type TransferMonthlyAmountRow struct {
	Month               string `gorm:"column:month" json:"month"`
	TotalTransferAmount int32  `gorm:"column:total_transfer_amount" json:"total_transfer_amount"`
}

// TransferYearlyAmountRow - yearly transfer amount stats
type TransferYearlyAmountRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalTransferAmount int64       `gorm:"column:total_transfer_amount" json:"total_transfer_amount"`
}

// TransferMonthlyStatusSuccessRow - monthly status success stats
type TransferMonthlyStatusSuccessRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalSuccess int64  `gorm:"column:total_success" json:"total_success"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransferYearlyStatusSuccessRow - yearly status success stats
type TransferYearlyStatusSuccessRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalSuccess int32  `gorm:"column:total_success" json:"total_success"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransferMonthlyStatusFailedRow - monthly status failed stats
type TransferMonthlyStatusFailedRow struct {
	Year        string `gorm:"column:year" json:"year"`
	Month       string `gorm:"column:month" json:"month"`
	TotalFailed int64  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransferYearlyStatusFailedRow - yearly status failed stats
type TransferYearlyStatusFailedRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalFailed int32  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransferMonthlyStatusSuccessByCardRow - monthly status success by card
type TransferMonthlyStatusSuccessByCardRow = TransferMonthlyStatusSuccessRow

// TransferYearlyStatusSuccessByCardRow - yearly status success by card
type TransferYearlyStatusSuccessByCardRow = TransferYearlyStatusSuccessRow

// TransferMonthlyStatusFailedByCardRow - monthly status failed by card
type TransferMonthlyStatusFailedByCardRow = TransferMonthlyStatusFailedRow

// TransferYearlyStatusFailedByCardRow - yearly status failed by card
type TransferYearlyStatusFailedByCardRow = TransferYearlyStatusFailedRow

// TransferMonthlyAmountBySenderCardRow - monthly amount by sender card
type TransferMonthlyAmountBySenderCardRow = TransferMonthlyAmountRow

// TransferYearlyAmountBySenderCardRow - yearly amount by sender card
type TransferYearlyAmountBySenderCardRow = TransferYearlyAmountRow

// TransferMonthlyAmountByReceiverCardRow - monthly amount by receiver card
type TransferMonthlyAmountByReceiverCardRow = TransferMonthlyAmountRow

// TransferYearlyAmountByReceiverCardRow - yearly amount by receiver card
type TransferYearlyAmountByReceiverCardRow = TransferYearlyAmountRow
