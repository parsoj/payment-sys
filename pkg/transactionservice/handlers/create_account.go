package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/parsoj/paymentsys/pkg/transactionservice/domain"
)

type CreateAccountRequest struct {
	UserID string `json:"user_id"`
}

type CreateAccountResponse struct {
	Account domain.Account `json:"account"`
	Error   string         `json:"error,omitempty"`
}

func CreateAccountHandler(svc *domain.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// error if userid is null
		if req.UserID == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		account, err := svc.CreateAccount(context.Background(), req.UserID)
		if err != nil {
			HandleInternalError(w, err)
			return
		}

		resp := CreateAccountResponse{
			Account: account,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
