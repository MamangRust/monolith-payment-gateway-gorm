package models

import "time"

// RoleRows - paginated role listing without deleted_at
type RoleRow struct {
	RoleID     int32     `gorm:"column:role_id" json:"role_id"`
	RoleName   string    `gorm:"column:role_name" json:"role_name"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount int64     `gorm:"column:total_count" json:"total_count"`
}

// RoleActiveRow - paginated active role listing with deleted_at
type RoleActiveRow struct {
	RoleID     int32      `gorm:"column:role_id" json:"role_id"`
	RoleName   string     `gorm:"column:role_name" json:"role_name"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount int64      `gorm:"column:total_count" json:"total_count"`
}

// RoleTrashedRow - paginated trashed role listing with deleted_at
type RoleTrashedRow struct {
	RoleID     int32      `gorm:"column:role_id" json:"role_id"`
	RoleName   string     `gorm:"column:role_name" json:"role_name"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount int64      `gorm:"column:total_count" json:"total_count"`
}
