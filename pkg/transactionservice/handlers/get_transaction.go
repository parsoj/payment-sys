package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/parsoj/paymentsys/pkg/transactionservice/domain"
)

type GetTransactionResponse struct {
	Transaction domain.Transaction `json:"transaction"`
	Error       string             `json:"error,omitempty"`
}

func GetTransactionHandler(svc *domain.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the transaction ID from the URL query parameters
		transactionID := r.URL.Query().Get("transaction_id")
		if transactionID == "" {
			http.Error(w, "Missing transaction_id parameter", http.StatusBadRequest)
			return
		}

		// Call the GetTransaction service function
		transaction, err := svc.GetTransaction(context.Background(), transactionID)
		if err != nil {
			HandleInternalError(w, err)
			return
		}

		// Prepare the response
		resp := GetTransactionResponse{
			Transaction: transaction,
		}

		// Encode the response as JSON and send it
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
