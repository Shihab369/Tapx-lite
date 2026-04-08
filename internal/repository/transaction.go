package repository

import (
	"context"
	"database/sql"

	"github.com/Shihab369/Tapx-lite/internal/model"
)

// TransactionRepository handles all database operations for transactions.
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository returns a new TransactionRepository with the given db connection.
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// NonceExists checks whether a given nonce has already been processed.
// Used to prevent duplicate transaction submissions.
func (r *TransactionRepository) NonceExists(ctx context.Context, nonce string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM transactions WHERE nonce = $1)`
	err := r.db.QueryRowContext(ctx, query, nonce).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// BeginTx starts a new database transaction for atomic multi-step operations.
func (r *TransactionRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// InsertTransactionTx inserts a transaction record within an existing database transaction.
// Relies on the UNIQUE constraint on nonce to reject duplicates at the database level.
func (r *TransactionRepository) InsertTransactionTx(ctx context.Context, tx *sql.Tx, t model.Transaction) error {
	query := `
		INSERT INTO transactions (id, user_id, amount, nonce, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.ExecContext(ctx, query,
		t.ID,
		t.UserID,
		t.Amount,
		t.Nonce,
		t.Status,
		t.CreatedAt,
	)
	return err
}

// UpdateBalanceTx deducts amount from the user's balance within a database transaction.
// The WHERE balance >= $1 guard prevents negative balances at the database level.
// Returns false (not an error) when the account has insufficient funds.
func (r *TransactionRepository) UpdateBalanceTx(ctx context.Context, tx *sql.Tx, userID int, amount int) (bool, error) {
	query := `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2 AND balance >= $1
	`
	result, err := tx.ExecContext(ctx, query, amount, userID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

// GetBalance returns the current balance for the given account ID.
func (r *TransactionRepository) GetBalance(ctx context.Context, userID int) (int, error) {
	var balance int
	query := `SELECT balance FROM accounts WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
