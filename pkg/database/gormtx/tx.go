package gormtx

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// RunInTx executes the given function within a database transaction. If fn
// returns nil the transaction is committed; on any error it is rolled back.
// A new *gorm.DB bound to the transaction is passed to fn, ensuring all
// operations inside fn participate in the same atomic unit.
func RunInTx(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := fn(tx); err != nil {
			return fmt.Errorf("transaction function failed: %w", err)
		}
		return nil
	})
}
