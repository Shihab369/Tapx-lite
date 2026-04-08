package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Shihab369/Tapx-lite/internal/model"
	"github.com/Shihab369/Tapx-lite/internal/repository"
)

// TransactionService handles business logic for transaction processing.
type TransactionService struct {
	repo *repository.TransactionRepository
}

// NewTransactionService returns a new TransactionService with the given repository.
func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

// ProcessSync validates and reconciles a batch of offline transactions.
// Each transaction runs in its own database transaction for atomicity.
// A failure on one transaction does not affect others in the batch.
func (s *TransactionService) ProcessSync(ctx context.Context, req model.SyncRequest) (model.SyncResponse, error) {
	var accepted []string
	var rejected []string

	for _, tx := range req.Transactions {
		ok := s.processSingle(ctx, tx)
		if ok {
			accepted = append(accepted, tx.ID)
		} else {
			rejected = append(rejected, tx.ID)
		}
	}

	return model.SyncResponse{
		Accepted: accepted,
		Rejected: rejected,
	}, nil
}

// processSingle handles one transaction atomically.
// Extracted from the loop to allow safe use of defer for rollback.
// Returns true if the transaction was committed successfully.
func (s *TransactionService) processSingle(ctx context.Context, tx model.Transaction) bool {
	dbTx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Printf("failed to begin db transaction for %s: %v", tx.ID, err)
		return false
	}

	// Rollback is safe to call even after a successful commit.
	defer dbTx.Rollback()

	tx.Status = "accepted"
	tx.CreatedAt = time.Now().UTC()

	// Step 1: Insert transaction record.
	// The UNIQUE constraint on nonce rejects duplicates at the database level.
	if err := s.repo.InsertTransactionTx(ctx, dbTx, tx); err != nil {
		log.Printf("insert failed for %s (possible duplicate nonce): %v", tx.ID, err)
		return false
	}

	// Step 2: Deduct balance atomically.
	// Returns false without error when balance is insufficient.
	ok, err := s.repo.UpdateBalanceTx(ctx, dbTx, tx.UserID, tx.Amount)
	if err != nil {
		log.Printf("balance update error for tx %s: %v", tx.ID, err)
		return false
	}
	if !ok {
		log.Printf("insufficient balance for user %d in tx %s", tx.UserID, tx.ID)
		return false
	}

	// Step 3: Commit.
	if err := dbTx.Commit(); err != nil {
		log.Printf("commit failed for tx %s: %v", tx.ID, err)
		return false
	}

	return true
}

// GetBalance retrieves the current balance for a given user account.
func (s *TransactionService) GetBalance(ctx context.Context, userID int) (int, error) {
	balance, err := s.repo.GetBalance(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("account not found for user %d", userID)
	}
	return balance, nil
}
