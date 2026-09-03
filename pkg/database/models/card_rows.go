package models

import "time"

// ---- Card CRUD rows ----

// CardCreateRow - response from creating a card
type CardCreateRow struct {
	CardID       int32     `gorm:"column:card_id" json:"card_id"`
	UserID       int32     `gorm:"column:user_id" json:"user_id"`
	CardNumber   string    `gorm:"column:card_number" json:"card_number"`
	CardType     string    `gorm:"column:card_type" json:"card_type"`
	ExpireDate   time.Time `gorm:"column:expire_date" json:"expire_date"`
	Cvv          string    `gorm:"column:cvv" json:"cvv"`
	CardProvider string    `gorm:"column:card_provider" json:"card_provider"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// CardUpdateRow - response from updating a card
type CardUpdateRow struct {
	CardID       int32     `gorm:"column:card_id" json:"card_id"`
	UserID       int32     `gorm:"column:user_id" json:"user_id"`
	CardNumber   string    `gorm:"column:card_number" json:"card_number"`
	CardType     string    `gorm:"column:card_type" json:"card_type"`
	ExpireDate   time.Time `gorm:"column:expire_date" json:"expire_date"`
	Cvv          string    `gorm:"column:cvv" json:"cvv"`
	CardProvider string    `gorm:"column:card_provider" json:"card_provider"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// CardTrashRow - response from trashing a card (with deleted_at)
type CardTrashRow struct {
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
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt          *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// CardRestoreRow - response from restoring a card
type CardRestoreRow struct {
	CardID       int32      `gorm:"column:card_id" json:"card_id"`
	UserID       int32      `gorm:"column:user_id" json:"user_id"`
	CardNumber   string     `gorm:"column:card_number" json:"card_number"`
	CardType     string     `gorm:"column:card_type" json:"card_type"`
	ExpireDate   time.Time  `gorm:"column:expire_date" json:"expire_date"`
	Cvv          string     `gorm:"column:cvv" json:"cvv"`
	CardProvider string     `gorm:"column:card_provider" json:"card_provider"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// CardAllFieldsRow - card with all fields (for status/credit/redeem operations)
type CardAllFieldsRow struct {
	CardID             int32     `gorm:"column:card_id" json:"card_id"`
	UserID             int32     `gorm:"column:user_id" json:"user_id"`
	CardNumber         string    `gorm:"column:card_number" json:"card_number"`
	CardType           string    `gorm:"column:card_type" json:"card_type"`
	ExpireDate         time.Time `gorm:"column:expire_date" json:"expire_date"`
	Cvv                string    `gorm:"column:cvv" json:"cvv"`
	CardProvider       string    `gorm:"column:card_provider" json:"card_provider"`
	Status             string    `gorm:"column:status" json:"status"`
	CreditLimit        int32     `gorm:"column:credit_limit" json:"credit_limit"`
	OutstandingBalance int32     `gorm:"column:outstanding_balance" json:"outstanding_balance"`
	RewardPoints       int32     `gorm:"column:reward_points" json:"reward_points"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// ---- Paginated card rows ----

// CardListRow - paginated card listing without deleted_at
type CardListRow struct {
	CardID       int32     `gorm:"column:card_id" json:"card_id"`
	UserID       int32     `gorm:"column:user_id" json:"user_id"`
	CardNumber   string    `gorm:"column:card_number" json:"card_number"`
	CardType     string    `gorm:"column:card_type" json:"card_type"`
	ExpireDate   time.Time `gorm:"column:expire_date" json:"expire_date"`
	Cvv          string    `gorm:"column:cvv" json:"cvv"`
	CardProvider string    `gorm:"column:card_provider" json:"card_provider"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount   int64     `gorm:"column:total_count" json:"total_count"`
}

// CardListWithDeletedRow - paginated card listing with deleted_at
type CardListWithDeletedRow struct {
	CardID       int32      `gorm:"column:card_id" json:"card_id"`
	UserID       int32      `gorm:"column:user_id" json:"user_id"`
	CardNumber   string     `gorm:"column:card_number" json:"card_number"`
	CardType     string     `gorm:"column:card_type" json:"card_type"`
	ExpireDate   time.Time  `gorm:"column:expire_date" json:"expire_date"`
	Cvv          string     `gorm:"column:cvv" json:"cvv"`
	CardProvider string     `gorm:"column:card_provider" json:"card_provider"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount   int64      `gorm:"column:total_count" json:"total_count"`
}

// CardByEmailRow - card with user email
type CardByEmailRow struct {
	CardID       int32     `gorm:"column:card_id" json:"card_id"`
	Email        string    `gorm:"column:email" json:"email"`
	UserID       int32     `gorm:"column:user_id" json:"user_id"`
	CardNumber   string    `gorm:"column:card_number" json:"card_number"`
	CardType     string    `gorm:"column:card_type" json:"card_type"`
	ExpireDate   time.Time `gorm:"column:expire_date" json:"expire_date"`
	Cvv          string    `gorm:"column:cvv" json:"cvv"`
	CardProvider string    `gorm:"column:card_provider" json:"card_provider"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// ---- Auth transaction rows ----

// GetAuthTxnByCardNumberRow - auth transactions by card number
type GetAuthTxnByCardNumberRow struct {
	AuthID         int32     `gorm:"column:auth_id" json:"auth_id"`
	TxnID          string    `gorm:"column:txn_id" json:"txn_id"`
	CardNumber     string    `gorm:"column:card_number" json:"card_number"`
	MerchantID     int32     `gorm:"column:merchant_id" json:"merchant_id"`
	Amount         int64     `gorm:"column:amount" json:"amount"`
	Currency       string    `gorm:"column:currency" json:"currency"`
	Mcc            string    `gorm:"column:mcc" json:"mcc"`
	PosEntryMode   string    `gorm:"column:pos_entry_mode" json:"pos_entry_mode"`
	Status         string    `gorm:"column:status" json:"status"`
	IdempotencyKey string    `gorm:"column:idempotency_key" json:"idempotency_key"`
	RiskScore      int32     `gorm:"column:risk_score" json:"risk_score"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// ---- Payment rows ----

// CardPaymentListRow - card payments with total count
type CardPaymentListRow struct {
	PaymentID      int32     `gorm:"column:payment_id" json:"payment_id"`
	PaymentUuid    string    `gorm:"column:payment_uuid" json:"payment_uuid"`
	CardNumber     string    `gorm:"column:card_number" json:"card_number"`
	BillingID      *int32    `gorm:"column:billing_id" json:"billing_id"`
	Amount         int64     `gorm:"column:amount" json:"amount"`
	PaymentChannel string    `gorm:"column:payment_channel" json:"payment_channel"`
	ReferenceID    string    `gorm:"column:reference_id" json:"reference_id"`
	Status         string    `gorm:"column:status" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount     int64     `gorm:"column:total_count" json:"total_count"`
}

// ---- Reward rows ----

// RedeemableRewardRow - reward IDs for redemption
type RedeemableRewardRow struct {
	RewardID     int32 `gorm:"column:reward_id" json:"reward_id"`
	PointsEarned int32 `gorm:"column:points_earned" json:"points_earned"`
}

// ---- Stats rows (monthly/yearly) ----

// MonthlyAmountRow - monthly stats with single amount
type MonthlyAmountRow struct {
	Month  string `gorm:"column:month" json:"month"`
	Amount int32  `gorm:"column:amount" json:"amount"`
}

// YearlyAmountRow - yearly stats with single amount
type YearlyAmountRow struct {
	Year        string `gorm:"column:year" json:"year"`
	Amount int64       `gorm:"column:amount" json:"amount"`
}

// MonthlyBalanceRow - monthly balance stats
type MonthlyBalanceRow struct {
	Month        string `gorm:"column:month" json:"month"`
	TotalBalance int32  `gorm:"column:total_balance" json:"total_balance"`
}

// YearlyBalanceRow - yearly balance stats
type YearlyBalanceRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalBalance int64       `gorm:"column:total_balance" json:"total_balance"`
}

// MonthlyTopupRow - monthly topup stats
type MonthlyTopupRow struct {
	Month            string `gorm:"column:month" json:"month"`
	TotalTopupAmount int32  `gorm:"column:total_topup_amount" json:"total_topup_amount"`
}

// YearlyTopupRow - yearly topup stats
type YearlyTopupRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalTopupAmount int64       `gorm:"column:total_topup_amount" json:"total_topup_amount"`
}

// MonthlyTransactionRow - monthly transaction stats
type MonthlyTransactionRow struct {
	Month                  string `gorm:"column:month" json:"month"`
	TotalTransactionAmount int32  `gorm:"column:total_transaction_amount" json:"total_transaction_amount"`
}

// YearlyTransactionRow - yearly transaction stats
type YearlyTransactionRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalTransactionAmount int64       `gorm:"column:total_transaction_amount" json:"total_transaction_amount"`
}

// MonthlyTransferSentRow - monthly transfer sent stats
type MonthlyTransferSentRow struct {
	Month           string `gorm:"column:month" json:"month"`
	TotalSentAmount int32  `gorm:"column:total_sent_amount" json:"total_sent_amount"`
}

// YearlyTransferSentRow - yearly transfer sent stats
type YearlyTransferSentRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalSentAmount int64       `gorm:"column:total_sent_amount" json:"total_sent_amount"`
}

// MonthlyTransferReceivedRow - monthly transfer received stats
type MonthlyTransferReceivedRow struct {
	Month               string `gorm:"column:month" json:"month"`
	TotalReceivedAmount int32  `gorm:"column:total_received_amount" json:"total_received_amount"`
}

// YearlyTransferReceivedRow - yearly transfer received stats
type YearlyTransferReceivedRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalReceivedAmount int64       `gorm:"column:total_received_amount" json:"total_received_amount"`
}

// MonthlyWithdrawRow - monthly withdraw stats
type MonthlyWithdrawRow struct {
	Month               string `gorm:"column:month" json:"month"`
	TotalWithdrawAmount int32  `gorm:"column:total_withdraw_amount" json:"total_withdraw_amount"`
}

// YearlyWithdrawRow - yearly withdraw stats
type YearlyWithdrawRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalWithdrawAmount int64       `gorm:"column:total_withdraw_amount" json:"total_withdraw_amount"`
}

// ---- Shared rows (used by multiple services) ----

// GetUserEmailByCardNumberRow - card with user email
type GetUserEmailByCardNumberRow struct {
	CardID       int32     `gorm:"column:card_id" json:"card_id"`
	Email        string    `gorm:"column:email" json:"email"`
	UserID       int32     `gorm:"column:user_id" json:"user_id"`
	CardNumber   string    `gorm:"column:card_number" json:"card_number"`
	CardType     string    `gorm:"column:card_type" json:"card_type"`
	ExpireDate   time.Time `gorm:"column:expire_date" json:"expire_date"`
	Cvv          string    `gorm:"column:cvv" json:"cvv"`
	CardProvider string    `gorm:"column:card_provider" json:"card_provider"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// UpdateCardRow - response from updating a card
type UpdateCardRow struct {
	CardID       int32     `gorm:"column:card_id" json:"card_id"`
	UserID       int32     `gorm:"column:user_id" json:"user_id"`
	CardNumber   string    `gorm:"column:card_number" json:"card_number"`
	CardType     string    `gorm:"column:card_type" json:"card_type"`
	ExpireDate   time.Time `gorm:"column:expire_date" json:"expire_date"`
	Cvv          string    `gorm:"column:cvv" json:"cvv"`
	CardProvider string    `gorm:"column:card_provider" json:"card_provider"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// AddRewardPointsRow - response from adding reward points
type AddRewardPointsRow struct {
	CardID             int32     `gorm:"column:card_id" json:"card_id"`
	UserID             int32     `gorm:"column:user_id" json:"user_id"`
	CardNumber         string    `gorm:"column:card_number" json:"card_number"`
	CardType           string    `gorm:"column:card_type" json:"card_type"`
	ExpireDate         time.Time `gorm:"column:expire_date" json:"expire_date"`
	Cvv                string    `gorm:"column:cvv" json:"cvv"`
	CardProvider       string    `gorm:"column:card_provider" json:"card_provider"`
	Status             string    `gorm:"column:status" json:"status"`
	CreditLimit        int32     `gorm:"column:credit_limit" json:"credit_limit"`
	OutstandingBalance int32     `gorm:"column:outstanding_balance" json:"outstanding_balance"`
	RewardPoints       int32     `gorm:"column:reward_points" json:"reward_points"`
}

// UpdateOutstandingBalanceRow - response from updating outstanding balance
type UpdateOutstandingBalanceRow struct {
	CardID             int32     `gorm:"column:card_id" json:"card_id"`
	UserID             int32     `gorm:"column:user_id" json:"user_id"`
	CardNumber         string    `gorm:"column:card_number" json:"card_number"`
	CardType           string    `gorm:"column:card_type" json:"card_type"`
	ExpireDate         time.Time `gorm:"column:expire_date" json:"expire_date"`
	Cvv                string    `gorm:"column:cvv" json:"cvv"`
	CardProvider       string    `gorm:"column:card_provider" json:"card_provider"`
	Status             string    `gorm:"column:status" json:"status"`
	CreditLimit        int32     `gorm:"column:credit_limit" json:"credit_limit"`
	OutstandingBalance int32     `gorm:"column:outstanding_balance" json:"outstanding_balance"`
	RewardPoints       int32     `gorm:"column:reward_points" json:"reward_points"`
}

// UpdateCardStatusRow - response from toggling card status
type UpdateCardStatusRow struct {
	CardID             int32     `gorm:"column:card_id" json:"card_id"`
	UserID             int32     `gorm:"column:user_id" json:"user_id"`
	CardNumber         string    `gorm:"column:card_number" json:"card_number"`
	CardType           string    `gorm:"column:card_type" json:"card_type"`
	ExpireDate         time.Time `gorm:"column:expire_date" json:"expire_date"`
	Cvv                string    `gorm:"column:cvv" json:"cvv"`
	CardProvider       string    `gorm:"column:card_provider" json:"card_provider"`
	Status             string    `gorm:"column:status" json:"status"`
	CreditLimit        int32     `gorm:"column:credit_limit" json:"credit_limit"`
	OutstandingBalance int32     `gorm:"column:outstanding_balance" json:"outstanding_balance"`
	RewardPoints       int32     `gorm:"column:reward_points" json:"reward_points"`
}
