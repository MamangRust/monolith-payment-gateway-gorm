package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Begin starts a transaction when the underlying sqlc DBTX is a pgx pool or
// connection. Repositories use WithTx to keep multi-step financial mutations
// atomic while generated query files remain source-controlled output.
func (q *Queries) Begin(ctx context.Context) (pgx.Tx, error) {
	switch db := q.db.(type) {
	case interface {
		Begin(context.Context) (pgx.Tx, error)
	}:
		return db.Begin(ctx)
	default:
		return nil, pgx.ErrTxClosed
	}
}
