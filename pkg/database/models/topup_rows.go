package models

import (
	"time"

	"github.com/google/uuid"
)

// TopupListRow - paginated topup listing
type TopupListRow struct {
	TopupID     int32     `gorm:"column:topup_id" json:"topup_id"`
	TopupNo     uuid.UUID `gorm:"column:topup_no" json:"topup_no"`
	CardNumber  string    `gorm:"column:card_number" json:"card_number"`
	TopupAmount int32     `gorm:"column:topup_amount" json:"topup_amount"`
	TopupMethod string    `gorm:"column:topup_method" json:"topup_method"`
	TopupTime   time.Time `gorm:"column:topup_time" json:"topup_time"`
	Status      string    `gorm:"column:status" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount  int64     `gorm:"column:total_count" json:"total_count"`
}

// TopupListWithDeletedRow - topup listing with deleted_at
type TopupListWithDeletedRow struct {
	TopupID     int32      `gorm:"column:topup_id" json:"topup_id"`
	TopupNo     uuid.UUID  `gorm:"column:topup_no" json:"topup_no"`
	CardNumber  string     `gorm:"column:card_number" json:"card_number"`
	TopupAmount int32      `gorm:"column:topup_amount" json:"topup_amount"`
	TopupMethod string     `gorm:"column:topup_method" json:"topup_method"`
	TopupTime   time.Time  `gorm:"column:topup_time" json:"topup_time"`
	Status      string     `gorm:"column:status" json:"status"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount  int64      `gorm:"column:total_count" json:"total_count"`
}

// TopupAllFieldsRow - topup with all fields
type TopupAllFieldsRow struct {
	TopupID     int32     `gorm:"column:topup_id" json:"topup_id"`
	TopupNo     uuid.UUID `gorm:"column:topup_no" json:"topup_no"`
	CardNumber  string    `gorm:"column:card_number" json:"card_number"`
	TopupAmount int32     `gorm:"column:topup_amount" json:"topup_amount"`
	TopupMethod string    `gorm:"column:topup_method" json:"topup_method"`
	TopupTime   time.Time `gorm:"column:topup_time" json:"topup_time"`
	Status      string    `gorm:"column:status" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TopupTrashRestoreRow - topup with deleted_at for trash/restore
type TopupTrashRestoreRow struct {
	TopupID     int32      `gorm:"column:topup_id" json:"topup_id"`
	TopupNo     uuid.UUID  `gorm:"column:topup_no" json:"topup_no"`
	CardNumber  string     `gorm:"column:card_number" json:"card_number"`
	TopupAmount int32      `gorm:"column:topup_amount" json:"topup_amount"`
	TopupMethod string     `gorm:"column:topup_method" json:"topup_method"`
	TopupTime   time.Time  `gorm:"column:topup_time" json:"topup_time"`
	Status      string     `gorm:"column:status" json:"status"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TopupByIdempotencyKeyRow - topup for idempotency check
type TopupByIdempotencyKeyRow struct {
	TopupID     int32     `gorm:"column:topup_id" json:"topup_id"`
	TopupNo     uuid.UUID `gorm:"column:topup_no" json:"topup_no"`
	CardNumber  string    `gorm:"column:card_number" json:"card_number"`
	TopupAmount int32     `gorm:"column:topup_amount" json:"topup_amount"`
	TopupMethod string    `gorm:"column:topup_method" json:"topup_method"`
	TopupTime   time.Time `gorm:"column:topup_time" json:"topup_time"`
	Status      string    `gorm:"column:status" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}
