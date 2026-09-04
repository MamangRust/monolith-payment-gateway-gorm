package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User maps to the "users" table.
type User struct {
	UserID           int32          `gorm:"column:user_id;primaryKey;autoIncrement" json:"user_id"`
	Firstname        string         `gorm:"column:firstname;size:100;not null" json:"firstname"`
	Lastname         string         `gorm:"column:lastname;size:100;not null" json:"lastname"`
	Email            string         `gorm:"column:email;size:100;uniqueIndex;not null" json:"email"`
	Password         string         `gorm:"column:password;size:100;not null" json:"password"`
	VerificationCode string         `gorm:"column:verification_code;size:100;not null" json:"verification_code"`
	IsVerified       *bool          `gorm:"column:is_verified;default:false" json:"is_verified"`
	CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (User) TableName() string { return "users" }

// Role maps to the "roles" table.
type Role struct {
	RoleID    int32          `gorm:"column:role_id;primaryKey;autoIncrement" json:"role_id"`
	RoleName  string         `gorm:"column:role_name;size:50;uniqueIndex;not null" json:"role_name"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Role) TableName() string { return "roles" }

// UserRole maps to the "user_roles" table.
type UserRole struct {
	UserRoleID int32          `gorm:"column:user_role_id;primaryKey;autoIncrement" json:"user_role_id"`
	UserID     int32          `gorm:"column:user_id;not null" json:"user_id"`
	RoleID     int32          `gorm:"column:role_id;not null" json:"role_id"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (UserRole) TableName() string { return "user_roles" }

// Card maps to the "cards" table.
type Card struct {
	CardID             int32          `gorm:"column:card_id;primaryKey;autoIncrement" json:"card_id"`
	UserID             int32          `gorm:"column:user_id;not null" json:"user_id"`
	CardNumber         string         `gorm:"column:card_number;size:16;uniqueIndex;not null" json:"card_number"`
	CardType           string         `gorm:"column:card_type;size:50;not null" json:"card_type"`
	ExpireDate         time.Time      `gorm:"column:expire_date;not null" json:"expire_date"`
	Cvv                string         `gorm:"column:cvv;size:3;not null" json:"cvv"`
	CardProvider       string         `gorm:"column:card_provider;size:50;not null" json:"card_provider"`
	Status             string         `gorm:"column:status;size:20;default:active" json:"status"`
	CreditLimit        int32          `gorm:"column:credit_limit;default:0" json:"credit_limit"`
	OutstandingBalance int32          `gorm:"column:outstanding_balance;default:0" json:"outstanding_balance"`
	RewardPoints       int32          `gorm:"column:reward_points;default:0" json:"reward_points"`
	EventVersion       int64          `gorm:"column:event_version;default:0" json:"event_version"`
	CreatedAt          time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Card) TableName() string { return "cards" }

// Merchant maps to the "merchants" table.
type Merchant struct {
	MerchantID   int32          `gorm:"column:merchant_id;primaryKey;autoIncrement" json:"merchant_id"`
	MerchantNo   uuid.UUID      `gorm:"column:merchant_no;type:uuid;default:gen_random_uuid()" json:"merchant_no"`
	Name         string         `gorm:"column:name;size:255;not null" json:"name"`
	ApiKey       string         `gorm:"column:api_key;size:255;uniqueIndex;not null" json:"api_key"`
	UserID       int32          `gorm:"column:user_id;not null" json:"user_id"`
	Status       string         `gorm:"column:status;size:20;default:pending" json:"status"`
	EventVersion int64          `gorm:"column:event_version;default:0" json:"event_version"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Merchant) TableName() string { return "merchants" }

// MerchantDocument maps to the "merchant_documents" table.
type MerchantDocument struct {
	DocumentID   int32          `gorm:"column:document_id;primaryKey;autoIncrement" json:"document_id"`
	MerchantID   int32          `gorm:"column:merchant_id;not null" json:"merchant_id"`
	DocumentType string         `gorm:"column:document_type;size:50;not null" json:"document_type"`
	DocumentUrl  string         `gorm:"column:document_url;type:text;not null" json:"document_url"`
	Status       string         `gorm:"column:status;size:20;default:pending" json:"status"`
	Note         *string        `gorm:"column:note;type:text" json:"note"`
	UploadedAt   time.Time      `gorm:"column:uploaded_at;autoCreateTime" json:"uploaded_at"`
	EventVersion int64          `gorm:"column:event_version;default:0" json:"event_version"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (MerchantDocument) TableName() string { return "merchant_documents" }

// Saldo maps to the "saldos" table.
type Saldo struct {
	SaldoID        int32          `gorm:"column:saldo_id;primaryKey;autoIncrement" json:"saldo_id"`
	CardNumber     string         `gorm:"column:card_number;size:16;not null" json:"card_number"`
	TotalBalance   int32          `gorm:"column:total_balance;not null" json:"total_balance"`
	WithdrawAmount *int32         `gorm:"column:withdraw_amount;default:0" json:"withdraw_amount"`
	WithdrawTime   *time.Time     `gorm:"column:withdraw_time" json:"withdraw_time"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Saldo) TableName() string { return "saldos" }

// Topup maps to the "topups" table.
type Topup struct {
	TopupID        int32          `gorm:"column:topup_id;primaryKey;autoIncrement" json:"topup_id"`
	TopupNo        uuid.UUID      `gorm:"column:topup_no;type:uuid;default:gen_random_uuid()" json:"topup_no"`
	CardNumber     string         `gorm:"column:card_number;size:16;not null" json:"card_number"`
	TopupAmount    int32          `gorm:"column:topup_amount;not null" json:"topup_amount"`
	TopupMethod    string         `gorm:"column:topup_method;size:50;not null" json:"topup_method"`
	TopupTime      time.Time      `gorm:"column:topup_time;not null" json:"topup_time"`
	Status         string         `gorm:"column:status;size:32;default:pending" json:"status"`
	IdempotencyKey string         `gorm:"column:idempotency_key;size:64;default:''" json:"idempotency_key"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Topup) TableName() string { return "topups" }

// Transaction maps to the "transactions" table.
type Transaction struct {
	TransactionID   int32          `gorm:"column:transaction_id;primaryKey;autoIncrement" json:"transaction_id"`
	TransactionNo   uuid.UUID      `gorm:"column:transaction_no;type:uuid;default:gen_random_uuid()" json:"transaction_no"`
	CardNumber      string         `gorm:"column:card_number;size:16;not null" json:"card_number"`
	Amount          int32          `gorm:"column:amount;not null" json:"amount"`
	PaymentMethod   string         `gorm:"column:payment_method;size:50;not null" json:"payment_method"`
	MerchantID      int32          `gorm:"column:merchant_id;not null" json:"merchant_id"`
	TransactionTime time.Time      `gorm:"column:transaction_time;not null" json:"transaction_time"`
	Status          string         `gorm:"column:status;size:32;default:pending" json:"status"`
	FraudScore      int32          `gorm:"column:fraud_score;default:0" json:"fraud_score"`
	IdempotencyKey  string         `gorm:"column:idempotency_key;size:64;default:''" json:"idempotency_key"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Transaction) TableName() string { return "transactions" }

// Transfer maps to the "transfers" table.
type Transfer struct {
	TransferID     int32          `gorm:"column:transfer_id;primaryKey;autoIncrement" json:"transfer_id"`
	TransferNo     uuid.UUID      `gorm:"column:transfer_no;type:uuid;default:gen_random_uuid()" json:"transfer_no"`
	TransferFrom   string         `gorm:"column:transfer_from;size:16;not null" json:"transfer_from"`
	TransferTo     string         `gorm:"column:transfer_to;size:16;not null" json:"transfer_to"`
	TransferAmount int32          `gorm:"column:transfer_amount;not null" json:"transfer_amount"`
	TransferTime   time.Time      `gorm:"column:transfer_time;not null" json:"transfer_time"`
	Status         string         `gorm:"column:status;size:32;default:pending" json:"status"`
	IdempotencyKey string         `gorm:"column:idempotency_key;size:64;default:''" json:"idempotency_key"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Transfer) TableName() string { return "transfers" }

// Withdraw maps to the "withdraws" table.
type Withdraw struct {
	WithdrawID     int32          `gorm:"column:withdraw_id;primaryKey;autoIncrement" json:"withdraw_id"`
	WithdrawNo     uuid.UUID      `gorm:"column:withdraw_no;type:uuid;default:gen_random_uuid()" json:"withdraw_no"`
	CardNumber     string         `gorm:"column:card_number;size:16;not null" json:"card_number"`
	WithdrawAmount int32          `gorm:"column:withdraw_amount;not null" json:"withdraw_amount"`
	WithdrawTime   time.Time      `gorm:"column:withdraw_time;not null" json:"withdraw_time"`
	Status         string         `gorm:"column:status;size:32;default:pending" json:"status"`
	IdempotencyKey string         `gorm:"column:idempotency_key;size:64;default:''" json:"idempotency_key"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Withdraw) TableName() string { return "withdraws" }

// RefreshToken maps to the "refresh_tokens" table.
type RefreshToken struct {
	RefreshTokenID int32          `gorm:"column:refresh_token_id;primaryKey;autoIncrement" json:"refresh_token_id"`
	UserID         int32          `gorm:"column:user_id;not null" json:"user_id"`
	Token          string         `gorm:"column:token;size:255;uniqueIndex;not null" json:"token"`
	Expiration     time.Time      `gorm:"column:expiration;not null" json:"expiration"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (RefreshToken) TableName() string { return "refresh_tokens" }

// ResetToken maps to the "reset_tokens" table.
type ResetToken struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"column:user_id;uniqueIndex;not null" json:"user_id"`
	Token      string    `gorm:"column:token;uniqueIndex;not null" json:"token"`
	ExpiryDate time.Time `gorm:"column:expiry_date;not null" json:"expiry_date"`
}

func (ResetToken) TableName() string { return "reset_tokens" }

// BillingCycle maps to the "billing_cycles" table.
type BillingCycle struct {
	BillingID  int32     `gorm:"column:billing_id;primaryKey;autoIncrement" json:"billing_id"`
	CardNumber string    `gorm:"column:card_number;size:16;not null" json:"card_number"`
	CycleStart time.Time `gorm:"column:cycle_start;not null" json:"cycle_start"`
	CycleEnd   time.Time `gorm:"column:cycle_end;not null" json:"cycle_end"`
	AmountDue  int32     `gorm:"column:amount_due;default:0" json:"amount_due"`
	DueDate    time.Time `gorm:"column:due_date;not null" json:"due_date"`
	Status     string    `gorm:"column:status;size:20;default:unpaid" json:"status"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (BillingCycle) TableName() string { return "billing_cycles" }

// CardAuthTransaction maps to the "card_auth_transactions" table.
type CardAuthTransaction struct {
	AuthID         int32     `gorm:"column:auth_id;primaryKey;autoIncrement" json:"auth_id"`
	TxnID          string    `gorm:"column:txn_id;size:36;uniqueIndex;not null" json:"txn_id"`
	CardNumber     string    `gorm:"column:card_number;size:16;not null" json:"card_number"`
	MerchantID     int32     `gorm:"column:merchant_id;default:0" json:"merchant_id"`
	Amount         int64     `gorm:"column:amount;default:0" json:"amount"`
	Currency       string    `gorm:"column:currency;size:3;default:IDR" json:"currency"`
	Mcc            string    `gorm:"column:mcc;size:4;default:''" json:"mcc"`
	PosEntryMode   string    `gorm:"column:pos_entry_mode;size:3;default:''" json:"pos_entry_mode"`
	Status         string    `gorm:"column:status;size:20;default:pending" json:"status"`
	IdempotencyKey string    `gorm:"column:idempotency_key;size:64;default:''" json:"idempotency_key"`
	RiskScore      int32     `gorm:"column:risk_score;default:0" json:"risk_score"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (CardAuthTransaction) TableName() string { return "card_auth_transactions" }

// CardPayment maps to the "card_payments" table.
type CardPayment struct {
	PaymentID      int32     `gorm:"column:payment_id;primaryKey;autoIncrement" json:"payment_id"`
	PaymentUuid    string    `gorm:"column:payment_uuid;size:36;uniqueIndex;not null" json:"payment_uuid"`
	CardNumber     string    `gorm:"column:card_number;size:16;not null" json:"card_number"`
	BillingID      *int32    `gorm:"column:billing_id" json:"billing_id"`
	Amount         int64     `gorm:"column:amount;default:0" json:"amount"`
	PaymentChannel string    `gorm:"column:payment_channel;size:20;default:bank_transfer" json:"payment_channel"`
	ReferenceID    string    `gorm:"column:reference_id;size:64;default:''" json:"reference_id"`
	Status         string    `gorm:"column:status;size:20;default:completed" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (CardPayment) TableName() string { return "card_payments" }

// CardReward maps to the "card_rewards" table.
type CardReward struct {
	RewardID     int32     `gorm:"column:reward_id;primaryKey;autoIncrement" json:"reward_id"`
	CardNumber   string    `gorm:"column:card_number;size:16;not null" json:"card_number"`
	TxnID        string    `gorm:"column:txn_id;size:36;default:''" json:"txn_id"`
	Amount       int64     `gorm:"column:amount;default:0" json:"amount"`
	Mcc          string    `gorm:"column:mcc;size:4;default:''" json:"mcc"`
	PointsEarned int32     `gorm:"column:points_earned;default:0" json:"points_earned"`
	ExpiresAt    time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	Redeemed     bool      `gorm:"column:redeemed;default:false" json:"redeemed"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (CardReward) TableName() string { return "card_rewards" }

// OutboxEvent maps to the "outbox_events" table.
type OutboxEvent struct {
	ID            int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	EventKey      string         `gorm:"column:event_key;size:128;uniqueIndex;not null" json:"event_key"`
	Topic         string         `gorm:"column:topic;size:255;not null" json:"topic"`
	MessageKey    string         `gorm:"column:message_key;size:255;default:''" json:"message_key"`
	Payload       []byte         `gorm:"column:payload;type:jsonb;not null" json:"payload"`
	Status        string         `gorm:"column:status;size:16;default:pending" json:"status"`
	Attempts      int32          `gorm:"column:attempts;default:0" json:"attempts"`
	NextAttemptAt time.Time      `gorm:"column:next_attempt_at;not null" json:"next_attempt_at"`
	LastError     string         `gorm:"column:last_error;type:text;default:''" json:"last_error"`
	ClaimVersion  int64          `gorm:"column:claim_version;default:0" json:"claim_version"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	PublishedAt   *time.Time     `gorm:"column:published_at" json:"published_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (OutboxEvent) TableName() string { return "outbox_events" }

// ConsumerInbox maps to the "consumer_inbox" table.
type ConsumerInbox struct {
	ConsumerName       string         `gorm:"column:consumer_name;primaryKey;size:128" json:"consumer_name"`
	EventKey           string         `gorm:"column:event_key;primaryKey;size:255" json:"event_key"`
	Topic              string         `gorm:"column:topic;size:255;default:''" json:"topic"`
	PartitionID        int32          `gorm:"column:partition_id;default:-1" json:"partition_id"`
	MessageOffset      int64          `gorm:"column:message_offset;default:-1" json:"message_offset"`
	Status             string         `gorm:"column:status;size:16;default:processing" json:"status"`
	Attempts           int32          `gorm:"column:attempts;default:1" json:"attempts"`
	LeaseUntil         time.Time      `gorm:"column:lease_until;not null" json:"lease_until"`
	LastError          string         `gorm:"column:last_error;type:text;default:''" json:"last_error"`
	ProcessedAt        *time.Time     `gorm:"column:processed_at" json:"processed_at"`
	ReservationVersion int64          `gorm:"column:reservation_version;default:0" json:"reservation_version"`
	CreatedAt          time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (ConsumerInbox) TableName() string { return "consumer_inbox" }

// CardEventLog maps to the "card_event_logs" table.
type CardEventLog struct {
	EventID     int64     `gorm:"column:event_id;primaryKey;autoIncrement" json:"event_id"`
	Topic       string    `gorm:"column:topic;size:80;not null" json:"topic"`
	EventType   string    `gorm:"column:event_type;size:80;not null" json:"event_type"`
	CardNumber  *string   `gorm:"column:card_number;size:16" json:"card_number"`
	ReferenceID *string   `gorm:"column:reference_id;size:64" json:"reference_id"`
	Payload     []byte    `gorm:"column:payload;type:jsonb;not null" json:"payload"`
	ReceivedAt  time.Time `gorm:"column:received_at;not null" json:"received_at"`
}

func (CardEventLog) TableName() string { return "card_event_logs" }
