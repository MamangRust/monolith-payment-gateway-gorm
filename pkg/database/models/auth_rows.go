package models

import "time"

// ResetTokenRow - reset token response
type ResetTokenRow struct {
	UserID     int64     `gorm:"column:user_id" json:"user_id"`
	Token      string    `gorm:"column:token" json:"token"`
	ExpiryDate time.Time `gorm:"column:expiry_date" json:"expiry_date"`
}
