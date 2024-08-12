package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/parsoj/paymentsys/pkg/transactionservice/domain"
)

type GetAccountResponse struct {
	Account domain.Account `json:"account"`
	Error   string         `json:"error,omitempty"`
}

func GetAccountHandler(svc *domain.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the account ID from the URL query parameters
		accountID := r.URL.Query().Get("account_id")
		if accountID == "" {
			http.Error(w, "Missing account_id parameter", http.StatusBadRequest)
			return
		}

		// Call the GetAccount service function
		account, err := svc.GetAccount(context.Background(), accountID)
		if err != nil {
			HandleInternalError(w, err)
			return
		}

		// Prepare the response
		resp := GetAccountResponse{
			Account: account,
		}

		// Encode the response as JSON and send it
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
