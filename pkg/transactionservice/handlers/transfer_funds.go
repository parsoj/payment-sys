package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/parsoj/paymentsys/pkg/transactionservice/domain"
)

type TransferFundsRequest struct {
	ToAccount      string  `json:"to_account"`
	FromAccount    string  `json:"from_account"`
	Amount         float64 `json:"amount"`
	IdempotencyKey *string `json:"idempotency_key,omitempty"`
}

type TransferFundsResponse struct {
	TransactionID string `json:"transaction_id"`
	Error         string `json:"error,omitempty"`
}

func TransferFundsHandler(svc *domain.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TransferFundsRequest

		// Decode the JSON request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if req.ToAccount == "" {
			http.Error(w, "to_account is required", http.StatusBadRequest)
		}

		if req.FromAccount == "" {
			http.Error(w, "from_account is required", http.StatusBadRequest)
		}

		if req.Amount == 0.0 {
			http.Error(w, "amount field is required and must be non-zero", http.StatusBadRequest)
		}

		// Call the TransferFunds service function
		transactionID, err := svc.TransferFunds(context.Background(), req.ToAccount, req.FromAccount, req.Amount, req.IdempotencyKey)
		if err != nil {
			HandleInternalError(w, err)
			return
		}

		// Prepare the response
		resp := TransferFundsResponse{
			TransactionID: transactionID,
		}

		// Encode the response as JSON and send it
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
