package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/parsoj/paymentsys/pkg/transactionservice/domain"
)

type WithdrawFundsRequest struct {
	AccountID      string  `json:"account_id"`
	Amount         float64 `json:"amount"`
	IdempotencyKey *string `json:"idempotency_key,omitempty"`
}

type WithdrawFundsResponse struct {
	Account domain.Account `json:"account"`
	Error   string         `json:"error,omitempty"`
}

func WithdrawFundsHandler(svc *domain.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req WithdrawFundsRequest

		// Decode the JSON request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// error if accountid is null
		if req.AccountID == "" {
			http.Error(w, "account_id is required", http.StatusBadRequest)
			return
		}

		//error if amount is null
		if req.Amount == 0 {
			http.Error(w, "amount is required and must be non-zero", http.StatusBadRequest)
			return
		}

		// Call the WithdrawFunds service function
		account, err := svc.WithdrawFunds(context.Background(), req.AccountID, req.Amount, req.IdempotencyKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare the response
		resp := WithdrawFundsResponse{
			Account: account,
		}

		// Encode the response as JSON and send it
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
