package models

import (
	"time"
)

// SaldoRow - paginated saldo listing
type SaldoRow struct {
	SaldoID        int32      `gorm:"column:saldo_id" json:"saldo_id"`
	CardNumber     string     `gorm:"column:card_number" json:"card_number"`
	TotalBalance   int32      `gorm:"column:total_balance" json:"total_balance"`
	WithdrawAmount *int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   *time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
	TotalCount     int64      `gorm:"column:total_count" json:"total_count"`
}

// SaldoActiveRow - active saldo listing with deleted_at
type SaldoActiveRow struct {
	SaldoID        int32      `gorm:"column:saldo_id" json:"saldo_id"`
	CardNumber     string     `gorm:"column:card_number" json:"card_number"`
	TotalBalance   int32      `gorm:"column:total_balance" json:"total_balance"`
	WithdrawAmount *int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   *time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount     int64      `gorm:"column:total_count" json:"total_count"`
}

// SaldoTrashedRow - trashed saldo listing
type SaldoTrashedRow struct {
	SaldoID        int32      `gorm:"column:saldo_id" json:"saldo_id"`
	CardNumber     string     `gorm:"column:card_number" json:"card_number"`
	TotalBalance   int32      `gorm:"column:total_balance" json:"total_balance"`
	WithdrawAmount *int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   *time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount     int64      `gorm:"column:total_count" json:"total_count"`
}

// SaldoByIDRow - single saldo by ID
type SaldoByIDRow struct {
	SaldoID        int32      `gorm:"column:saldo_id" json:"saldo_id"`
	CardNumber     string     `gorm:"column:card_number" json:"card_number"`
	TotalBalance   int32      `gorm:"column:total_balance" json:"total_balance"`
	WithdrawAmount *int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   *time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// CreateSaldoRow - created saldo response
type CreateSaldoRow struct {
	SaldoID        int32      `gorm:"column:saldo_id" json:"saldo_id"`
	CardNumber     string     `gorm:"column:card_number" json:"card_number"`
	TotalBalance   int32      `gorm:"column:total_balance" json:"total_balance"`
	WithdrawAmount *int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   *time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// UpdateSaldoRow - updated saldo response
type UpdateSaldoRow struct {
	SaldoID        int32      `gorm:"column:saldo_id" json:"saldo_id"`
	CardNumber     string     `gorm:"column:card_number" json:"card_number"`
	TotalBalance   int32      `gorm:"column:total_balance" json:"total_balance"`
	WithdrawAmount *int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   *time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// UpdateSaldoBalanceRow - saldo balance update response
type UpdateSaldoBalanceRow struct {
	SaldoID        int32      `gorm:"column:saldo_id" json:"saldo_id"`
	CardNumber     string     `gorm:"column:card_number" json:"card_number"`
	TotalBalance   int32      `gorm:"column:total_balance" json:"total_balance"`
	WithdrawAmount *int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   *time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// UpdateSaldoWithdrawRow - saldo withdraw update response
type UpdateSaldoWithdrawRow struct {
	SaldoID        int32      `gorm:"column:saldo_id" json:"saldo_id"`
	CardNumber     string     `gorm:"column:card_number" json:"card_number"`
	TotalBalance   int32      `gorm:"column:total_balance" json:"total_balance"`
	WithdrawAmount *int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   *time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// MonthlySaldoBalanceRow - monthly saldo balance stats
type MonthlySaldoBalanceRow struct {
	Month        string `gorm:"column:month" json:"month"`
	TotalBalance int32  `gorm:"column:total_balance" json:"total_balance"`
}

// YearlySaldoBalanceRow - yearly saldo balance stats
type YearlySaldoBalanceRow struct {
	YYear        string `gorm:"column:y_year" json:"y_year"`
	TotalBalance int64  `gorm:"column:total_balance" json:"total_balance"`
}

// MonthlyTotalSaldoBalanceRow - monthly total saldo balance comparison
type MonthlyTotalSaldoBalanceRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalBalance int32  `gorm:"column:total_balance" json:"total_balance"`
}

// YearlyTotalSaldoBalancesRow - yearly total saldo balance comparison
type YearlyTotalSaldoBalancesRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalBalance int32  `gorm:"column:total_balance" json:"total_balance"`
}

// GetCardByCardNumberRow - card by card number for saldo
type SaldoCardByCardNumberRow struct {
	CardID             int32      `gorm:"column:card_id" json:"card_id"`
	UserID             int32      `gorm:"column:user_id" json:"user_id"`
	CardNumber         string     `gorm:"column:card_number" json:"card_number"`
	CardType           string     `gorm:"column:card_type" json:"card_type"`
	ExpireDate         time.Time  `gorm:"column:expire_date" json:"expire_date"`
	Cvv                string     `gorm:"column:cvv" json:"cvv"`
	CardProvider       string     `gorm:"column:card_provider" json:"card_provider"`
	Status             string     `gorm:"column:status" json:"status"`
	CreditLimit        int32      `gorm:"column:credit_limit" json:"credit_limit"`
	OutstandingBalance int32      `gorm:"column:outstanding_balance" json:"outstanding_balance"`
	RewardPoints       int32      `gorm:"column:reward_points" json:"reward_points"`
	EventVersion       int64      `gorm:"column:event_version" json:"event_version"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt          *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// UpdateSaldoBalanceDeltaRow - saldo balance delta response
type UpdateSaldoBalanceDeltaRow struct {
	SaldoID        int32      `gorm:"column:saldo_id" json:"saldo_id"`
	CardNumber     string     `gorm:"column:card_number" json:"card_number"`
	TotalBalance   int32      `gorm:"column:total_balance" json:"total_balance"`
	WithdrawAmount *int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   *time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
}
