package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/parsoj/paymentsys/pkg/transactionservice/domain"
)

type GetUserResponse struct {
	User  domain.User `json:"user"`
	Error string      `json:"error,omitempty"`
}

func GetUserHandler(svc *domain.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the user ID from the URL query parameters
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
			return
		}

		// Call the GetUser service function
		user, err := svc.GetUser(context.Background(), userID)
		if err != nil {
			HandleInternalError(w, err)
			return
		}

		// Prepare the response
		resp := GetUserResponse{
			User: user,
		}

		// Encode the response as JSON and send it
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
