package main

import (
	"log"
	"net/http"

	"github.com/parsoj/paymentsys/pkg/transactionservice/domain"
	"github.com/parsoj/paymentsys/pkg/transactionservice/router"
	"github.com/parsoj/paymentsys/pkg/transactionservice/sqldb"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

func main() {

	connString := "postgres://transaction_svc:dev_pass@txn_db:5432/transactions_db?sslmode=disable"
	db, err := sqldb.NewSqlDB(connString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	svc, err := domain.NewTransactionService(db)
	if err != nil {
		log.Fatalf("Failed to start TransactionService: %v", err)
	}

	// Set up the router
	mux := router.SetupRouter(svc)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
