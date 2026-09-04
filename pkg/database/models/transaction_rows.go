package models

import (
	"time"

	"github.com/google/uuid"
)

// TransactionListRow - paginated transaction listing
type TransactionListRow struct {
	TransactionID   int32     `gorm:"column:transaction_id" json:"transaction_id"`
	TransactionNo   uuid.UUID `gorm:"column:transaction_no" json:"transaction_no"`
	CardNumber      string    `gorm:"column:card_number" json:"card_number"`
	Amount          int32     `gorm:"column:amount" json:"amount"`
	PaymentMethod   string    `gorm:"column:payment_method" json:"payment_method"`
	MerchantID      int32     `gorm:"column:merchant_id" json:"merchant_id"`
	TransactionTime time.Time `gorm:"column:transaction_time" json:"transaction_time"`
	Status          string    `gorm:"column:status" json:"status"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount      int64     `gorm:"column:total_count" json:"total_count"`
}

// TransactionListWithDeletedRow - transaction listing with deleted_at
type TransactionListWithDeletedRow struct {
	TransactionID   int32      `gorm:"column:transaction_id" json:"transaction_id"`
	TransactionNo   uuid.UUID  `gorm:"column:transaction_no" json:"transaction_no"`
	CardNumber      string     `gorm:"column:card_number" json:"card_number"`
	Amount          int32      `gorm:"column:amount" json:"amount"`
	PaymentMethod   string     `gorm:"column:payment_method" json:"payment_method"`
	MerchantID      int32      `gorm:"column:merchant_id" json:"merchant_id"`
	TransactionTime time.Time  `gorm:"column:transaction_time" json:"transaction_time"`
	Status          string     `gorm:"column:status" json:"status"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount      int64      `gorm:"column:total_count" json:"total_count"`
}

// TransactionAllFieldsRow - transaction with all fields
type TransactionAllFieldsRow struct {
	TransactionID      int32     `gorm:"column:transaction_id" json:"transaction_id"`
	TransactionNo      uuid.UUID `gorm:"column:transaction_no" json:"transaction_no"`
	CardNumber         string    `gorm:"column:card_number" json:"card_number"`
	Amount             int32     `gorm:"column:amount" json:"amount"`
	PaymentMethod      string    `gorm:"column:payment_method" json:"payment_method"`
	MerchantID         int32     `gorm:"column:merchant_id" json:"merchant_id"`
	TransactionTime    time.Time `gorm:"column:transaction_time" json:"transaction_time"`
	Status             string    `gorm:"column:status" json:"status"`
	IdempotencyKey     string    `gorm:"column:idempotency_key" json:"idempotency_key"`
	CardNumber2        string    `gorm:"column:card_number_2" json:"card_number_2"`
	MerchantCardNumber string    `gorm:"column:merchant_card_number" json:"merchant_card_number"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TransactionTrashRestoreRow - transaction with deleted_at for trash/restore
type TransactionTrashRestoreRow struct {
	TransactionID   int32      `gorm:"column:transaction_id" json:"transaction_id"`
	TransactionNo   uuid.UUID  `gorm:"column:transaction_no" json:"transaction_no"`
	CardNumber      string     `gorm:"column:card_number" json:"card_number"`
	Amount          int32      `gorm:"column:amount" json:"amount"`
	PaymentMethod   string     `gorm:"column:payment_method" json:"payment_method"`
	MerchantID      int32      `gorm:"column:merchant_id" json:"merchant_id"`
	TransactionTime time.Time  `gorm:"column:transaction_time" json:"transaction_time"`
	Status          string     `gorm:"column:status" json:"status"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TransactionByIdempotencyKeyRow - transaction for idempotency check
type TransactionByIdempotencyKeyRow struct {
	TransactionID   int32     `gorm:"column:transaction_id" json:"transaction_id"`
	TransactionNo   uuid.UUID `gorm:"column:transaction_no" json:"transaction_no"`
	CardNumber      string    `gorm:"column:card_number" json:"card_number"`
	Amount          int32     `gorm:"column:amount" json:"amount"`
	PaymentMethod   string    `gorm:"column:payment_method" json:"payment_method"`
	MerchantID      int32     `gorm:"column:merchant_id" json:"merchant_id"`
	TransactionTime time.Time `gorm:"column:transaction_time" json:"transaction_time"`
	Status          string    `gorm:"column:status" json:"status"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}
