package router

import (
	"net/http"

	"github.com/parsoj/paymentsys/pkg/transactionservice/domain"
	"github.com/parsoj/paymentsys/pkg/transactionservice/handlers"
)

func SetupRouter(svc *domain.TransactionService) *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the handlers
	mux.HandleFunc("/create-user", handlers.CreateUserHandler(svc))
	mux.HandleFunc("/get-user", handlers.GetUserHandler(svc))
	mux.HandleFunc("/create-account", handlers.CreateAccountHandler(svc))
	mux.HandleFunc("/get-account", handlers.GetAccountHandler(svc))
	mux.HandleFunc("/deposit-funds", handlers.DepositFundsHandler(svc))
	mux.HandleFunc("/withdraw-funds", handlers.WithdrawFundsHandler(svc))
	mux.HandleFunc("/get-transaction", handlers.GetTransactionHandler(svc))
	mux.HandleFunc("/transfer-funds", handlers.TransferFundsHandler(svc))

	return mux
}

