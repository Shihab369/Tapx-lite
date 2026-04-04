package model

import "time"

type Transaction struct {
	ID 	  string        `json:"id"`
	UserID int          `json:"user_id"`
	Amount int          `json:"amount"`
	Nonce string        `json:"nonce"`
	Status string       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type SyncRequest struct {
	Transactions []Transaction `json:"transactions"`
}

type SyncResponse struct {
	Accepted []string `json:"accepted"`
	Rejected []string  `json:"rejected"`
	}