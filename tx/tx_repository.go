package tx


import "context"

// TxRepository defines a contract for database transaction handling.
//
// This interface is intended to be used by the usecase layer to
// execute business logic within a transactional context without
// directly depending on database implementations.
type TxRepository interface {

	// WithPostgresDBTx executes the given function within a PostgreSQL transaction.
	//
	// If a transaction already exists in the context, it will be reused.
	// Otherwise, a new transaction is started and automatically committed
	// or rolled back based on the function result.
	WithPostgresDBTx(ctx context.Context, fn func(ctx context.Context) error) error

	// WithTimescaleDBTx executes the given function within a TimescaleDB transaction.
	//
	// If a transaction already exists in the context, it will be reused.
	// Otherwise, a new transaction is started and automatically committed
	// or rolled back based on the function result.
	WithTimescaleDBTx(ctx context.Context, fn func(ctx context.Context) error) error
}
