package models

import "time"

// UserRow - all users listing with total_count
type UserRow struct {
	UserID     int32     `gorm:"column:user_id" json:"user_id"`
	Firstname  string    `gorm:"column:firstname" json:"firstname"`
	Lastname   string    `gorm:"column:lastname" json:"lastname"`
	Email      string    `gorm:"column:email" json:"email"`
	Password   string    `gorm:"column:password" json:"password"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount int64     `gorm:"column:total_count" json:"total_count"`
}

// UserActiveRow - active users with deleted_at and total_count
type UserActiveRow struct {
	UserID     int32      `gorm:"column:user_id" json:"user_id"`
	Firstname  string     `gorm:"column:firstname" json:"firstname"`
	Lastname   string     `gorm:"column:lastname" json:"lastname"`
	Email      string     `gorm:"column:email" json:"email"`
	Password   string     `gorm:"column:password" json:"password"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount int64      `gorm:"column:total_count" json:"total_count"`
}

// UserTrashedRow - trashed users with deleted_at and total_count
type UserTrashedRow struct {
	UserID     int32      `gorm:"column:user_id" json:"user_id"`
	Firstname  string     `gorm:"column:firstname" json:"firstname"`
	Lastname   string     `gorm:"column:lastname" json:"lastname"`
	Email      string     `gorm:"column:email" json:"email"`
	Password   string     `gorm:"column:password" json:"password"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount int64      `gorm:"column:total_count" json:"total_count"`
}

// UserByIDRow - single user by ID
type UserByIDRow struct {
	UserID    int32     `gorm:"column:user_id" json:"user_id"`
	Firstname string    `gorm:"column:firstname" json:"firstname"`
	Lastname  string    `gorm:"column:lastname" json:"lastname"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// UserByEmailRow - user by email
type UserByEmailRow struct {
	UserID int32  `gorm:"column:user_id" json:"user_id"`
	Email  string `gorm:"column:email" json:"email"`
}

// UserByEmailWithPasswordRow - user by email with password
type UserByEmailWithPasswordRow struct {
	UserID   int32  `gorm:"column:user_id" json:"user_id"`
	Email    string `gorm:"column:email" json:"email"`
	Password string `gorm:"column:password" json:"password"`
}

// UserByEmailAndVerifiedRow - verified user by email
type UserByEmailAndVerifiedRow struct {
	UserID    int32     `gorm:"column:user_id" json:"user_id"`
	Firstname string    `gorm:"column:firstname" json:"firstname"`
	Lastname  string    `gorm:"column:lastname" json:"lastname"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// UserByVerificationCodeRow - user by verification code
type UserByVerificationCodeRow struct {
	UserID    int32     `gorm:"column:user_id" json:"user_id"`
	Firstname string    `gorm:"column:firstname" json:"firstname"`
	Lastname  string    `gorm:"column:lastname" json:"lastname"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// CreateUserRow - created user response
type CreateUserRow struct {
	UserID    int32     `gorm:"column:user_id" json:"user_id"`
	Firstname string    `gorm:"column:firstname" json:"firstname"`
	Lastname  string    `gorm:"column:lastname" json:"lastname"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// UpdateUserRow - updated user response
type UpdateUserRow struct {
	UserID    int32     `gorm:"column:user_id" json:"user_id"`
	Firstname string    `gorm:"column:firstname" json:"firstname"`
	Lastname  string    `gorm:"column:lastname" json:"lastname"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TrashUserRow - trashed user response
type TrashUserRow struct {
	UserID    int32      `gorm:"column:user_id" json:"user_id"`
	Firstname string     `gorm:"column:firstname" json:"firstname"`
	Lastname  string     `gorm:"column:lastname" json:"lastname"`
	Email     string     `gorm:"column:email" json:"email"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// RestoreUserRow - restored user response
type RestoreUserRow struct {
	UserID    int32      `gorm:"column:user_id" json:"user_id"`
	Firstname string     `gorm:"column:firstname" json:"firstname"`
	Lastname  string     `gorm:"column:lastname" json:"lastname"`
	Email     string     `gorm:"column:email" json:"email"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// UserIsVerifiedRow - user is verified response
type UserIsVerifiedRow struct {
	UserID    int32     `gorm:"column:user_id" json:"user_id"`
	Firstname string    `gorm:"column:firstname" json:"firstname"`
	Lastname  string    `gorm:"column:lastname" json:"lastname"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// UserPasswordRow - user password update response
type UserPasswordRow struct {
	UserID    int32     `gorm:"column:user_id" json:"user_id"`
	Firstname string    `gorm:"column:firstname" json:"firstname"`
	Lastname  string    `gorm:"column:lastname" json:"lastname"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
