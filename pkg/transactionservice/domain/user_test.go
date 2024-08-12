package domain

import (
	"context"
	"testing"
)

func TestCreateUserAndGetUser(t *testing.T) {
	// Create a new mock database connection using your specified setup
	db, err := NewSQLiteDB()
	if err != nil {
		t.Fatalf("failed to create sqliteDB: %v", err)
	}
	defer db.Close()

	// Initialize your service using the mock database
	service, err := NewTransactionService(db)
	if err != nil {
		t.Fatalf("failed to initialize TransactionService: %v", err)
	}

	// Set up expectations
	username := "testuser"

	// Call your methods as usual
	createdUser, err := service.CreateUser(context.Background(), username)
	if err != nil {
		t.Errorf("unexpected error during CreateUser: %v", err)
	}

	readUser, err := service.GetUser(context.Background(), createdUser.ID)
	if err != nil {
		t.Errorf("unexpected error during GetUser: %v", err)
	}

	// Compare the results
	if createdUser.Username != readUser.Username || createdUser.ID != readUser.ID {
		t.Errorf("expected user %+v, got %+v", createdUser, readUser)
	}

}
