package models

import (
	"time"

	"github.com/google/uuid"
)

// ---- Merchant rows ----

// MerchantListRow - paginated merchant listing
type MerchantListRow struct {
	MerchantID int32     `gorm:"column:merchant_id" json:"merchant_id"`
	MerchantNo uuid.UUID `gorm:"column:merchant_no" json:"merchant_no"`
	Name       string    `gorm:"column:name" json:"name"`
	ApiKey     string    `gorm:"column:api_key" json:"api_key"`
	UserID     int32     `gorm:"column:user_id" json:"user_id"`
	Status     string    `gorm:"column:status" json:"status"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount int64     `gorm:"column:total_count" json:"total_count"`
}

// MerchantListWithDeletedRow - paginated merchant listing with deleted_at
type MerchantListWithDeletedRow struct {
	MerchantID int32      `gorm:"column:merchant_id" json:"merchant_id"`
	MerchantNo uuid.UUID  `gorm:"column:merchant_no" json:"merchant_no"`
	Name       string     `gorm:"column:name" json:"name"`
	ApiKey     string     `gorm:"column:api_key" json:"api_key"`
	UserID     int32      `gorm:"column:user_id" json:"user_id"`
	Status     string     `gorm:"column:status" json:"status"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount int64      `gorm:"column:total_count" json:"total_count"`
}

// MerchantAllFieldsRow - merchant with all fields (no deleted_at)
type MerchantAllFieldsRow struct {
	MerchantID int32     `gorm:"column:merchant_id" json:"merchant_id"`
	MerchantNo uuid.UUID `gorm:"column:merchant_no" json:"merchant_no"`
	Name       string    `gorm:"column:name" json:"name"`
	ApiKey     string    `gorm:"column:api_key" json:"api_key"`
	UserID     int32     `gorm:"column:user_id" json:"user_id"`
	Status     string    `gorm:"column:status" json:"status"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// MerchantListByUserRow - merchant list by user ID
type MerchantListByUserRow struct {
	MerchantID int32     `gorm:"column:merchant_id" json:"merchant_id"`
	MerchantNo uuid.UUID `gorm:"column:merchant_no" json:"merchant_no"`
	Name       string    `gorm:"column:name" json:"name"`
	ApiKey     string    `gorm:"column:api_key" json:"api_key"`
	UserID     int32     `gorm:"column:user_id" json:"user_id"`
	Status     string    `gorm:"column:status" json:"status"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// ---- Merchant Document rows ----

// MerchantDocumentListRow - paginated merchant document listing
type MerchantDocumentListRow struct {
	DocumentID   int32     `gorm:"column:document_id" json:"document_id"`
	MerchantID   int32     `gorm:"column:merchant_id" json:"merchant_id"`
	DocumentType string    `gorm:"column:document_type" json:"document_type"`
	DocumentUrl  string    `gorm:"column:document_url" json:"document_url"`
	Status       string    `gorm:"column:status" json:"status"`
	Note         *string   `gorm:"column:note" json:"note"`
	UploadedAt   time.Time `gorm:"column:uploaded_at" json:"uploaded_at"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount   int64     `gorm:"column:total_count" json:"total_count"`
}

// MerchantDocumentListWithDeletedRow - merchant document listing with deleted_at
type MerchantDocumentListWithDeletedRow struct {
	DocumentID   int32      `gorm:"column:document_id" json:"document_id"`
	MerchantID   int32      `gorm:"column:merchant_id" json:"merchant_id"`
	DocumentType string     `gorm:"column:document_type" json:"document_type"`
	DocumentUrl  string     `gorm:"column:document_url" json:"document_url"`
	Status       string     `gorm:"column:status" json:"status"`
	Note         *string    `gorm:"column:note" json:"note"`
	UploadedAt   time.Time  `gorm:"column:uploaded_at" json:"uploaded_at"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount   int64      `gorm:"column:total_count" json:"total_count"`
}

// MerchantDocumentAllFieldsRow - merchant document with all fields
type MerchantDocumentAllFieldsRow struct {
	DocumentID   int32     `gorm:"column:document_id" json:"document_id"`
	MerchantID   int32     `gorm:"column:merchant_id" json:"merchant_id"`
	DocumentType string    `gorm:"column:document_type" json:"document_type"`
	DocumentUrl  string    `gorm:"column:document_url" json:"document_url"`
	Status       string    `gorm:"column:status" json:"status"`
	Note         *string   `gorm:"column:note" json:"note"`
	UploadedAt   time.Time `gorm:"column:uploaded_at" json:"uploaded_at"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// MerchantDocumentCreateRow - created merchant document
type MerchantDocumentCreateRow struct {
	DocumentID   int32     `gorm:"column:document_id" json:"document_id"`
	MerchantID   int32     `gorm:"column:merchant_id" json:"merchant_id"`
	DocumentType string    `gorm:"column:document_type" json:"document_type"`
	DocumentUrl  string    `gorm:"column:document_url" json:"document_url"`
	Status       string    `gorm:"column:status" json:"status"`
	Note         *string   `gorm:"column:note" json:"note"`
	UploadedAt   time.Time `gorm:"column:uploaded_at" json:"uploaded_at"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// MerchantDocumentUpdateRow - updated merchant document
type MerchantDocumentUpdateRow struct {
	DocumentID   int32     `gorm:"column:document_id" json:"document_id"`
	MerchantID   int32     `gorm:"column:merchant_id" json:"merchant_id"`
	DocumentType string    `gorm:"column:document_type" json:"document_type"`
	DocumentUrl  string    `gorm:"column:document_url" json:"document_url"`
	Status       string    `gorm:"column:status" json:"status"`
	Note         *string   `gorm:"column:note" json:"note"`
	UploadedAt   time.Time `gorm:"column:uploaded_at" json:"uploaded_at"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// ---- Merchant Transaction rows ----

// MerchantTransactionRow - transaction listing for merchant
type MerchantTransactionRow struct {
	TransactionID   int32      `gorm:"column:transaction_id" json:"transaction_id"`
	CardNumber      string     `gorm:"column:card_number" json:"card_number"`
	Amount          int32      `gorm:"column:amount" json:"amount"`
	PaymentMethod   string     `gorm:"column:payment_method" json:"payment_method"`
	MerchantID      int32      `gorm:"column:merchant_id" json:"merchant_id"`
	MerchantName    string     `gorm:"column:merchant_name" json:"merchant_name"`
	TransactionTime time.Time  `gorm:"column:transaction_time" json:"transaction_time"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount      int64      `gorm:"column:total_count" json:"total_count"`
}

// MerchantApiKeyRow - merchant found by API key
type MerchantApiKeyRow struct {
	MerchantID int32  `gorm:"column:merchant_id" json:"merchant_id"`
	Name       string `gorm:"column:name" json:"name"`
	ApiKey     string `gorm:"column:api_key" json:"api_key"`
	UserID     int32  `gorm:"column:user_id" json:"user_id"`
	Status     string `gorm:"column:status" json:"status"`
}
