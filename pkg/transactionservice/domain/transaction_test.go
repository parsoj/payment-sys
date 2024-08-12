package domain

import (
	"context"
	"fmt"
	"testing"

	"github.com/segmentio/ksuid" // For generating unique IDs
)

// setupTest initializes two users and two accounts.
func setupTest(t *testing.T, svc *TransactionService) (User, Account, User, Account) {
	// Create two unique usernames
	user1Name := ksuid.New().String()[:8]
	user2Name := ksuid.New().String()[:8]

	user1, err := svc.CreateUser(context.Background(), user1Name)
	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}

	account1, err := svc.CreateAccount(context.Background(), user1.ID)
	if err != nil {
		t.Fatalf("Failed to create account for user1: %v", err)
	}

	user2, err := svc.CreateUser(context.Background(), user2Name)
	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	account2, err := svc.CreateAccount(context.Background(), user2.ID)
	if err != nil {
		t.Fatalf("Failed to create account for user2: %v", err)
	}

	return user1, account1, user2, account2
}

func initializeAccountBalances(t *testing.T, svc *TransactionService, account1 Account, amount1 float64, account2 Account, amount2 float64) {

	id, err := svc.idgen.New()
	if err != nil {
		t.Fatalf("failed to generate id for transaction: %v", err)
	}
	idp_id := string(id)

	_, err = svc.DepositFunds(context.Background(), account1.ID, amount1, &idp_id)
	if err != nil {
		t.Fatalf("Failed to set balance for account1: %v", err)
	}

	id, err = svc.idgen.New()
	if err != nil {
		t.Fatalf("failed to generate id for transaction: %v", err)
	}
	idp_id = string(id)

	_, err = svc.DepositFunds(context.Background(), account2.ID, amount2, &idp_id)
	if err != nil {
		t.Fatalf("Failed to set balance for account2: %v", err)
	}

}

func TestTransferFunds(t *testing.T) {
	db, err := NewSQLiteDB()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	svc, err := NewTransactionService(db)
	if err != nil {
		t.Fatalf("failed to initialize TransactionService: %v", err)
	}

	_, account1, _, account2 := setupTest(t, svc)

	initializeAccountBalances(t, svc, account1, 10.0, account2, 10.0)

	_, err = svc.TransferFunds(context.Background(), account2.ID, account1.ID, 5.0, nil)
	if err != nil {
		fmt.Printf("%v", err)
		t.Fatalf("unexpected error during TransferFunds: %v", err)
	}

	newAccount1, err := svc.GetAccount(context.Background(), account1.ID)
	if err != nil {
		t.Fatalf("unexpected error fetching balance for account1: %v", err)
	}

	newAccount2, err := svc.GetAccount(context.Background(), account2.ID)
	if err != nil {
		t.Fatalf("unexpected error fetching balance for account2: %v", err)
	}

	if newAccount1.Balance != 5.0 {
		t.Errorf("expected balance 5.0 for account1, got %v", newAccount1.Balance)
	}

	if newAccount2.Balance != 15.0 {
		t.Errorf("expected balance 15.0 for account2, got %v", newAccount2.Balance)
	}

}

func TestTransferFunds_InsufficientBalance(t *testing.T) {
	db, err := NewSQLiteDB()
	if err != nil {
		t.Fatalf("failed to create SQLiteDB: %v", err)
	}
	defer db.Close()

	svc, err := NewTransactionService(db)
	if err != nil {
		t.Fatalf("failed to initialize TransactionService: %v", err)
	}

	_, account1, _, account2 := setupTest(t, svc)

	// Initialize both account balances to 5
	initializeAccountBalances(t, svc, account1, 5.0, account2, 5.0)

	// Attempt to transfer 10 from account1 to account2, which should fail due to insufficient funds
	_, err = svc.TransferFunds(context.Background(), account2.ID, account1.ID, 10.0, nil)
	if err == nil {
		t.Fatalf("expected error due to insufficient funds, but got no error")
	}

	if err.Error() != fmt.Sprintf("Insufficient funds in source account. Balance: %f -- Transfer Amount: %f", 5.0, 10.0) {
		t.Errorf("expected insufficient funds error, got: %v", err)
	}
}
