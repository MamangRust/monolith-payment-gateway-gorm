package gormerr

import (
	"errors"

	"gorm.io/gorm"
)

// IsRecordNotFound returns true when the error indicates that the requested
// record does not exist in the database.
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// IsDuplicateKey returns true when the error indicates a unique constraint
// violation (postgres error code 23505).
func IsDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, gorm.ErrDuplicatedKey)
}
