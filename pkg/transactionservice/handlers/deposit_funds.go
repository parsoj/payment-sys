package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/parsoj/paymentsys/pkg/transactionservice/domain"
)

type DepositFundsRequest struct {
	AccountID      string  `json:"account_id"`
	Amount         float64 `json:"amount"`
	IdempotencyKey *string `json:"idempotency_key,omitempty"`
}

type DepositFundsResponse struct {
	Account domain.Account `json:"account"`
	Error   string         `json:"error,omitempty"`
}

func DepositFundsHandler(svc *domain.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DepositFundsRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			fmt.Println(r)
			fmt.Println(err)
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

		account, err := svc.DepositFunds(context.Background(), req.AccountID, req.Amount, req.IdempotencyKey)
		if err != nil {
			HandleInternalError(w, err)
			return
		}

		resp := DepositFundsResponse{
			Account: account,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
