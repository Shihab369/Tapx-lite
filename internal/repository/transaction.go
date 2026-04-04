package repository

import (
	"context"
	"database/sql"
	
	"github.com/Shihab369/Tapx-lite/internal/model"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository { 
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) NonceExists (cxt context, nonce string) (bool, error) {
	var exists bool
	query := 'SELECT EXISTS(SELECT 1 FROM transactions WHERE nonce = $1)'
	err := r.db.QueryRowContext(cxt, query, nonce).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *TransactionRepository) InsertTransaction(cxt context.Context, userID int) (int, error){
	var balance int
	query := 'SELECT balance FROM users WHERE id = $1'
	err := r.db.QueryRowContext(cxt, query, userID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *TransactionRepository) UpdateBalance(ctx context.Context, userID int, amount int) error {
	query := 'UPDATE accounts SET balance = balance - $1 WHERE id = $2'
	_, err := r.db.ExecContext(ctx, query, amount, userID)
	return err
}