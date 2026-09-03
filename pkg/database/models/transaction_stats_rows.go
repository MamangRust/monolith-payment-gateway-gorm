package models

// Transaction stats row types

// TransactionMonthlyStatusSuccessRow - monthly status success stats
type TransactionMonthlyStatusSuccessRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalSuccess int64  `gorm:"column:total_success" json:"total_success"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyStatusSuccessRow - yearly status success stats
type TransactionYearlyStatusSuccessRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalSuccess int32  `gorm:"column:total_success" json:"total_success"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyStatusFailedRow - monthly status failed stats
type TransactionMonthlyStatusFailedRow struct {
	Year        string `gorm:"column:year" json:"year"`
	Month       string `gorm:"column:month" json:"month"`
	TotalFailed int64  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyStatusFailedRow - yearly status failed stats
type TransactionYearlyStatusFailedRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalFailed int32  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyPaymentMethodRow - monthly payment method stats
type TransactionMonthlyPaymentMethodRow struct {
	Month         string `gorm:"column:month" json:"month"`
	PaymentMethod string `gorm:"column:payment_method" json:"payment_method"`
	TotalCount    int32  `gorm:"column:total_count" json:"total_count"`
	TotalAmount   int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyPaymentMethodRow - yearly payment method stats
type TransactionYearlyPaymentMethodRow struct {
	Year        string `gorm:"column:year" json:"year"`
	PaymentMethod string `gorm:"column:payment_method" json:"payment_method"`
	TotalCount    int64  `gorm:"column:total_count" json:"total_count"`
	TotalAmount   int64  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyAmountRow - monthly amount stats
type TransactionMonthlyAmountRow struct {
	Month       string `gorm:"column:month" json:"month"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyAmountRow - yearly amount stats
type TransactionYearlyAmountRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalAmount int64       `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyStatusSuccessByCardRow - monthly status success by card
type TransactionMonthlyStatusSuccessByCardRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalSuccess int64  `gorm:"column:total_success" json:"total_success"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyStatusSuccessByCardRow - yearly status success by card
type TransactionYearlyStatusSuccessByCardRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalSuccess int32  `gorm:"column:total_success" json:"total_success"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyStatusFailedByCardRow - monthly status failed by card
type TransactionMonthlyStatusFailedByCardRow struct {
	Year        string `gorm:"column:year" json:"year"`
	Month       string `gorm:"column:month" json:"month"`
	TotalFailed int64  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyStatusFailedByCardRow - yearly status failed by card
type TransactionYearlyStatusFailedByCardRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalFailed int32  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyPaymentMethodByCardRow - monthly payment method by card
type TransactionMonthlyPaymentMethodByCardRow struct {
	Month         string `gorm:"column:month" json:"month"`
	PaymentMethod string `gorm:"column:payment_method" json:"payment_method"`
	TotalCount    int32  `gorm:"column:total_count" json:"total_count"`
	TotalAmount   int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyPaymentMethodByCardRow - yearly payment method by card
type TransactionYearlyPaymentMethodByCardRow struct {
	Year        string `gorm:"column:year" json:"year"`
	PaymentMethod string `gorm:"column:payment_method" json:"payment_method"`
	TotalCount    int64  `gorm:"column:total_count" json:"total_count"`
	TotalAmount   int64  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyAmountByCardRow - monthly amount by card
type TransactionMonthlyAmountByCardRow struct {
	Month       string `gorm:"column:month" json:"month"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyAmountByCardRow - yearly amount by card
type TransactionYearlyAmountByCardRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalAmount int64       `gorm:"column:total_amount" json:"total_amount"`
}
