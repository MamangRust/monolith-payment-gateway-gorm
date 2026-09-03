package models

// Merchant stats row types

// MerchantMonthlyPaymentMethodRow - monthly transaction stats by payment method
type MerchantMonthlyPaymentMethodRow struct {
	Month         string `gorm:"column:month" json:"month"`
	PaymentMethod string `gorm:"column:payment_method" json:"payment_method"`
	TotalAmount   int32  `gorm:"column:total_amount" json:"total_amount"`
}

// MerchantYearlyPaymentMethodRow - yearly transaction stats by payment method
type MerchantYearlyPaymentMethodRow struct {
	Year        string `gorm:"column:year" json:"year"`
	PaymentMethod string      `gorm:"column:payment_method" json:"payment_method"`
	TotalAmount   int64       `gorm:"column:total_amount" json:"total_amount"`
}

// MerchantMonthlyAmountRow - monthly transaction amount
type MerchantMonthlyAmountRow struct {
	Month       string `gorm:"column:month" json:"month"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// MerchantYearlyAmountRow - yearly transaction amount
type MerchantYearlyAmountRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalAmount int64       `gorm:"column:total_amount" json:"total_amount"`
}

// MerchantMonthlyTotalAmountRow - monthly total amount (with year)
type MerchantMonthlyTotalAmountRow struct {
	Year        string `gorm:"column:year" json:"year"`
	Month       string `gorm:"column:month" json:"month"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// MerchantYearlyTotalAmountRow - yearly total amount
type MerchantYearlyTotalAmountRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}
