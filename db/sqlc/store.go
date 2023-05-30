package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provide all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

// SQLStore provide all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries // by embedding this inside Store, all functions in Queries is available in Store
	db       *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{Queries: New(db), db: db}
}

// execTx executes a function within a db transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("txErr: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
