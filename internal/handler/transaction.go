package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Shihab369/Tapx-lite/internal/model"
	"github.com/Shihab369/Tapx-lite/internal/service"
)

// TransactionHandler handles HTTP requests for transaction operations.
type TransactionHandler struct {
	service *service.TransactionService
}

// NewTransactionHandler returns a new TransactionHandler with the given service.
func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// HealthCheck godoc
// GET /health
// Returns server liveness status. Used by load balancers and container orchestrators.
func (h *TransactionHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "tapx-lite",
	})
}

// Sync godoc
// POST /sync
// Accepts a batch of offline transactions and reconciles them against the ledger.
// Returns lists of accepted and rejected transaction IDs.
func (h *TransactionHandler) Sync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.SyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Transactions) == 0 {
		http.Error(w, "no transactions provided", http.StatusBadRequest)
		return
	}

	resp, err := h.service.ProcessSync(r.Context(), req)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetBalance godoc
// GET /balance?user_id={id}
// Returns the current account balance for the specified user.
func (h *TransactionHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	balance, err := h.service.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{
		"user_id": userID,
		"balance": balance,
	})
}