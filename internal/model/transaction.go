package model

import "time"

// Transaction represents a single offline payment record.
type Transaction struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    int       `json:"amount"` // stored in smallest currency unit (paisa)
	Nonce     string    `json:"nonce"`  // unique one-time value for idempotency
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// SyncRequest holds a batch of offline transactions submitted for reconciliation.
type SyncRequest struct {
	Transactions []Transaction `json:"transactions"`
}

// SyncResponse reports which transactions were accepted or rejected after sync.
type SyncResponse struct {
	Accepted []string `json:"accepted"`
	Rejected []string `json:"rejected"`
}