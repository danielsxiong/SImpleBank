package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provide all functions to execute db queries and transactions
type Store struct {
	*Queries // by embedding this inside Store, all functions in Queries is available in Store
	db       *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{Queries: New(db), db: db}
}

// execTx executes a function within a db transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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

// TransferTxParams contains input params for transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from 1 account to the other
// Creates transfer record, account entries, and update accounts balance within a single db transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx,
			CreateTransferParams{
				FromAccountID: arg.FromAccountID,
				ToAccountID:   arg.ToAccountID,
				Amount:        arg.Amount,
			})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx,
			CreateEntryParams{
				AccountID: arg.FromAccountID,
				Amount:    -arg.Amount,
			})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx,
			CreateEntryParams{
				AccountID: arg.ToAccountID,
				Amount:    arg.Amount,
			})
		if err != nil {
			return err
		}

		// update account balance
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
