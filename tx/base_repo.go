package tx

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// contextKey is a private type to avoid context key collisions.
type contextKey string

const (
	txKey          contextKey = "postgres_tx"
	timescaleTxKey contextKey = "timescale_tx"
)

// BaseRepo provides transaction management for PostgreSQL and TimescaleDB.
//
// It enables context-based transaction propagation, allowing multiple
// repositories to share the same transaction across layers.
type BaseRepo struct {
	postgresDB  *sql.DB
	timescaleDB *pgxpool.Pool
}

// Compile-time assertion to ensure BaseRepo implements TxRepository.
var _ TxRepository = (*BaseRepo)(nil)

// NewBaseRepo creates a new BaseRepo instance.
//
// postgresDB   → *sql.DB for PostgreSQL
// timescaleDB  → *pgxpool.Pool for TimescaleDB
func NewBaseRepo(postgresDB *sql.DB, timescaleDB *pgxpool.Pool) *BaseRepo {
	return &BaseRepo{
		postgresDB:  postgresDB,
		timescaleDB: timescaleDB,
	}
}

// -----------------------------
// TimescaleDB Transaction
// -----------------------------

// WithTimescaleDBTx executes the given function within a TimescaleDB transaction.
//
// If a transaction already exists in the context, it will be reused.
// The transaction is automatically committed on success or rolled
// back on error or panic.
func (r *BaseRepo) WithTimescaleDBTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {

	// Reuse existing transaction if present
	if _, ok := r.GetTimescaleTx(ctx); ok {
		return fn(ctx)
	}

	tx, err := r.timescaleDB.Begin(ctx)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, timescaleTxKey, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

// -----------------------------
// PostgreSQL Transaction
// -----------------------------

// WithPostgresDBTx executes the given function within a PostgreSQL transaction.
//
// If a transaction already exists in the context, it will be reused.
// The transaction is automatically committed on success or rolled
// back on error or panic.
func (r *BaseRepo) WithPostgresDBTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {

	// Reuse existing transaction if present
	if _, ok := r.GetTxFromContext(ctx); ok {
		return fn(ctx)
	}

	tx, err := r.postgresDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// -----------------------------
// Transaction Extractors
// -----------------------------

// GetTxFromContext retrieves a PostgreSQL transaction from the context.
func (r *BaseRepo) GetTxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	return tx, ok
}

// GetTimescaleTx retrieves a TimescaleDB transaction from the context.
func (r *BaseRepo) GetTimescaleTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(timescaleTxKey).(pgx.Tx)
	return tx, ok
}

// -----------------------------
// Query Executors
// -----------------------------

// PostgresQueryExecutor returns a PostgreSQL query executor.
//
// If a transaction exists in the context, it is returned.
// Otherwise, the base *sql.DB instance is used.
func (r *BaseRepo) PostgresQueryExecutor(ctx context.Context) interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
} {
	if tx, ok := r.GetTxFromContext(ctx); ok {
		return tx
	}
	return r.postgresDB
}

// TimescaleQueryExecutor returns a TimescaleDB query executor.
//
// If a transaction exists in the context, it is returned.
// Otherwise, the base *pgxpool.Pool instance is used.
func (r *BaseRepo) TimescaleQueryExecutor(ctx context.Context) interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
} {
	if tx, ok := r.GetTimescaleTx(ctx); ok {
		return tx
	}
	return r.timescaleDB
}
