package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ctxKey struct{}

// DBTX is the subset of pgx pool/tx used by repositories (sqlc-style).
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Querier returns the transaction bound to ctx, or the pool.
func (p *Postgres) Querier(ctx context.Context) DBTX {
	if p == nil {
		return nil
	}
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}
	return p.Pool
}

// WithTx runs fn inside a single PostgreSQL transaction.
// Nested WithTx reuses the existing transaction (no savepoints).
func (p *Postgres) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := TxFromContext(ctx); ok {
		return fn(ctx)
	}

	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(ContextWithTx(ctx, tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// ContextWithTx stores tx on ctx for Querier.
func ContextWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, ctxKey{}, tx)
}

// TxFromContext returns the transaction bound by WithTx, if any.
func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(ctxKey{}).(pgx.Tx)
	return tx, ok
}
