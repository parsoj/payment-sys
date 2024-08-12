package domain

import (
	"context"
	"testing"
)

func TestCreateAndGetAccount(t *testing.T) {
	// Initialize sqlmock
	db, err := NewSQLiteDB()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Initialize the service
	svc := &TransactionService{db: db}

	account, err := svc.CreateAccount(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error during CreateAccount: %v", err)
	}

	// Call GetAccount
	retrievedAccount, err := svc.GetAccount(context.Background(), account.ID)
	if err != nil {

		t.Fatalf("unexpected error during GetAccount: %v", err)
	}

	// Compare the results
	if account.ID != retrievedAccount.ID || account.UserID != retrievedAccount.UserID || account.Balance != retrievedAccount.Balance {
		t.Errorf("expected account %+v, got %+v", account, retrievedAccount)
	}

}

func TestDepositAndWithdrawFunds(t *testing.T) {
	// Initialize sqlmock
	db, err := NewSQLiteDB()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := "2"

	//********************************************************************************
	// Create Account
	svc, err := NewTransactionService(db)
	if err != nil {
		t.Errorf("Failed to start Service: %v", err)
	}

	account, err := svc.CreateAccount(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error during CreateAccount: %v", err)
	}

	//********************************************************************************
	// Add Funds
	//
	id, err := svc.idgen.New()
	if err != nil {
		t.Fatalf("failed to generate id for transaction: %v", err)
	}
	idp_id := string(id)

	accountAfterDeposit, err := svc.DepositFunds(context.Background(), account.ID, 10.0, &idp_id)
	if err != nil {
		t.Fatalf("unexpected error during DepositFunds: %v", err)
	}

	if accountAfterDeposit.Balance != 10.0 {
		t.Errorf("expected balance 10.0, got %v", accountAfterDeposit.Balance)
	}

	//********************************************************************************
	// Remove Funds

	id, err = svc.idgen.New()
	if err != nil {
		t.Fatalf("failed to generate id for transaction: %v", err)
	}
	idp_id = string(id)

	accountAfterWithdrawal, err := svc.WithdrawFunds(context.Background(), account.ID, 5.0, &idp_id)
	if err != nil {
		t.Fatalf("unexpected error during WithdrawFunds: %v", err)
	}

	if accountAfterWithdrawal.Balance != 5.0 {
		t.Errorf("expected balance 5.0, got %v", accountAfterWithdrawal.Balance)
	}

}
