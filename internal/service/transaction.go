package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Shihab369/Tapx-lite/internal/model"
	"github.com/Shihab369/Tapx-lite/internal/repository"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) ProcessSync(ctx context.Context, req model.SyncRequest) (model.SyncResponse, error) {
	var accepted []string
	var rejected []string

	for _, tx := range req.Transactions {

		// Step 1: Nonce check koro
		exists, err := s.repo.NonceExists(ctx, tx.Nonce)
		if err != nil {
			rejected = append(rejected, tx.ID)
			continue
		}
		if exists {
			rejected = append(rejected, tx.ID)
			continue
		}

		// Step 2: Balance check koro
		balance, err := s.repo.GetBalance(ctx, tx.UserID)
		if err != nil {
			rejected = append(rejected, tx.ID)
			continue
		}
		if balance < tx.Amount {
			rejected = append(rejected, tx.ID)
			continue
		}

		// Step 3: Transaction save koro
		tx.Status = "accepted"
		tx.CreatedAt = time.Now()
		err = s.repo.InsertTransaction(ctx, tx)
		if err != nil {
			rejected = append(rejected, tx.ID)
			continue
		}

		// Step 4: Balance update koro
		err = s.repo.UpdateBalance(ctx, tx.UserID, tx.Amount)
		if err != nil {
			rejected = append(rejected, tx.ID)
			continue
		}

		accepted = append(accepted, tx.ID)
	}

	return model.SyncResponse{
		Accepted: accepted,
		Rejected: rejected,
	}, nil
}

func (s *TransactionService) GetBalance(ctx context.Context, userID int) (int, error) {
	balance, err := s.repo.GetBalance(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("user %d not found", userID)
	}
	return balance, nil
}