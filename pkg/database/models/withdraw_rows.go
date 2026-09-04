package models

import (
	"time"

	"github.com/google/uuid"
)

// WithdrawListRow - paginated withdraw listing
type WithdrawListRow struct {
	WithdrawID     int32     `gorm:"column:withdraw_id" json:"withdraw_id"`
	WithdrawNo     uuid.UUID `gorm:"column:withdraw_no" json:"withdraw_no"`
	CardNumber     string    `gorm:"column:card_number" json:"card_number"`
	WithdrawAmount int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	Status         string    `gorm:"column:status" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount     int64     `gorm:"column:total_count" json:"total_count"`
}

// WithdrawListWithDeletedRow - withdraw listing with deleted_at
type WithdrawListWithDeletedRow struct {
	WithdrawID     int32      `gorm:"column:withdraw_id" json:"withdraw_id"`
	WithdrawNo     uuid.UUID  `gorm:"column:withdraw_no" json:"withdraw_no"`
	CardNumber     string     `gorm:"column:card_number" json:"card_number"`
	WithdrawAmount int32      `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   time.Time  `gorm:"column:withdraw_time" json:"withdraw_time"`
	Status         string     `gorm:"column:status" json:"status"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount     int64      `gorm:"column:total_count" json:"total_count"`
}

// WithdrawAllFieldsRow - withdraw with all fields
type WithdrawAllFieldsRow struct {
	WithdrawID     int32     `gorm:"column:withdraw_id" json:"withdraw_id"`
	WithdrawNo     uuid.UUID `gorm:"column:withdraw_no" json:"withdraw_no"`
	CardNumber     string    `gorm:"column:card_number" json:"card_number"`
	WithdrawAmount int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	Status         string    `gorm:"column:status" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// WithdrawTrashRestoreRow - withdraw with deleted_at for trash/restore
type WithdrawTrashRestoreRow struct {
	WithdrawID     int32      `gorm:"column:withdraw_id" json:"withdraw_id"`
	WithdrawNo     uuid.UUID  `gorm:"column:withdraw_no" json:"withdraw_no"`
	CardNumber     string     `gorm:"column:card_number" json:"card_number"`
	WithdrawAmount int32      `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   time.Time  `gorm:"column:withdraw_time" json:"withdraw_time"`
	Status         string     `gorm:"column:status" json:"status"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// WithdrawByIdempotencyKeyRow - withdraw for idempotency check
type WithdrawByIdempotencyKeyRow struct {
	WithdrawID     int32     `gorm:"column:withdraw_id" json:"withdraw_id"`
	WithdrawNo     uuid.UUID `gorm:"column:withdraw_no" json:"withdraw_no"`
	CardNumber     string    `gorm:"column:card_number" json:"card_number"`
	WithdrawAmount int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	Status         string    `gorm:"column:status" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// WithdrawByCardNumberRow - withdraw by card number
type WithdrawByCardNumberRow struct {
	WithdrawID     int32     `gorm:"column:withdraw_id" json:"withdraw_id"`
	WithdrawNo     uuid.UUID `gorm:"column:withdraw_no" json:"withdraw_no"`
	CardNumber     string    `gorm:"column:card_number" json:"card_number"`
	WithdrawAmount int32     `gorm:"column:withdraw_amount" json:"withdraw_amount"`
	WithdrawTime   time.Time `gorm:"column:withdraw_time" json:"withdraw_time"`
	Status         string    `gorm:"column:status" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount     int64     `gorm:"column:total_count" json:"total_count"`
}
