package models

import (
	"time"

	"github.com/google/uuid"
)

// TransferListRow - paginated transfer listing
type TransferListRow struct {
	TransferID     int32     `gorm:"column:transfer_id" json:"transfer_id"`
	TransferNo     uuid.UUID `gorm:"column:transfer_no" json:"transfer_no"`
	TransferFrom   string    `gorm:"column:transfer_from" json:"transfer_from"`
	TransferTo     string    `gorm:"column:transfer_to" json:"transfer_to"`
	TransferAmount int32     `gorm:"column:transfer_amount" json:"transfer_amount"`
	TransferTime   time.Time `gorm:"column:transfer_time" json:"transfer_time"`
	Status         string    `gorm:"column:status" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount     int64     `gorm:"column:total_count" json:"total_count"`
}

// TransferListWithDeletedRow - transfer listing with deleted_at
type TransferListWithDeletedRow struct {
	TransferID     int32      `gorm:"column:transfer_id" json:"transfer_id"`
	TransferNo     uuid.UUID  `gorm:"column:transfer_no" json:"transfer_no"`
	TransferFrom   string     `gorm:"column:transfer_from" json:"transfer_from"`
	TransferTo     string     `gorm:"column:transfer_to" json:"transfer_to"`
	TransferAmount int32      `gorm:"column:transfer_amount" json:"transfer_amount"`
	TransferTime   time.Time  `gorm:"column:transfer_time" json:"transfer_time"`
	Status         string     `gorm:"column:status" json:"status"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount     int64      `gorm:"column:total_count" json:"total_count"`
}

// TransferAllFieldsRow - transfer with all fields
type TransferAllFieldsRow struct {
	TransferID         int32     `gorm:"column:transfer_id" json:"transfer_id"`
	TransferNo         uuid.UUID `gorm:"column:transfer_no" json:"transfer_no"`
	TransferFrom       string    `gorm:"column:transfer_from" json:"transfer_from"`
	TransferTo         string    `gorm:"column:transfer_to" json:"transfer_to"`
	TransferAmount     int32     `gorm:"column:transfer_amount" json:"transfer_amount"`
	TransferTime       time.Time `gorm:"column:transfer_time" json:"transfer_time"`
	Status             string    `gorm:"column:status" json:"status"`
	TransferID2        int32     `gorm:"column:transfer_id_2" json:"transfer_id_2"`
	MerchantCardNumber string    `gorm:"column:merchant_card_number" json:"merchant_card_number"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TransferTrashRestoreRow - transfer with deleted_at for trash/restore
type TransferTrashRestoreRow struct {
	TransferID     int32      `gorm:"column:transfer_id" json:"transfer_id"`
	TransferNo     uuid.UUID  `gorm:"column:transfer_no" json:"transfer_no"`
	TransferFrom   string     `gorm:"column:transfer_from" json:"transfer_from"`
	TransferTo     string     `gorm:"column:transfer_to" json:"transfer_to"`
	TransferAmount int32      `gorm:"column:transfer_amount" json:"transfer_amount"`
	TransferTime   time.Time  `gorm:"column:transfer_time" json:"transfer_time"`
	Status         string     `gorm:"column:status" json:"status"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TransferByIdempotencyKeyRow - transfer for idempotency check
type TransferByIdempotencyKeyRow struct {
	TransferID     int32     `gorm:"column:transfer_id" json:"transfer_id"`
	TransferNo     uuid.UUID `gorm:"column:transfer_no" json:"transfer_no"`
	TransferFrom   string    `gorm:"column:transfer_from" json:"transfer_from"`
	TransferTo     string    `gorm:"column:transfer_to" json:"transfer_to"`
	TransferAmount int32     `gorm:"column:transfer_amount" json:"transfer_amount"`
	TransferTime   time.Time `gorm:"column:transfer_time" json:"transfer_time"`
	Status         string    `gorm:"column:status" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TransferBySourceCardRow - transfers from a source card
type TransferBySourceCardRow struct {
	TransferID     int32     `gorm:"column:transfer_id" json:"transfer_id"`
	TransferNo     uuid.UUID `gorm:"column:transfer_no" json:"transfer_no"`
	TransferFrom   string    `gorm:"column:transfer_from" json:"transfer_from"`
	TransferTo     string    `gorm:"column:transfer_to" json:"transfer_to"`
	TransferAmount int32     `gorm:"column:transfer_amount" json:"transfer_amount"`
	TransferTime   time.Time `gorm:"column:transfer_time" json:"transfer_time"`
	Status         string    `gorm:"column:status" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TransferByDestinationCardRow - transfers to a destination card
type TransferByDestinationCardRow = TransferBySourceCardRow
