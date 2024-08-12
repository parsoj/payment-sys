package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/parsoj/paymentsys/pkg/transactionservice/domain"
)

type CreateUserRequest struct {
	Username string `json:"username"`
}

type CreateUserResponse struct {
	User  domain.User `json:"user"`
	Error string      `json:"error,omitempty"`
}

func CreateUserHandler(svc *domain.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user, err := svc.CreateUser(context.Background(), req.Username)
		if err != nil {
			HandleInternalError(w, err)
			return
		}

		//if username is null - error
		if user.Username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		resp := CreateUserResponse{
			User: user,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
