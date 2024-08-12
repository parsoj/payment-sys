package domain

import (
	"context"
	"fmt"
)

type Account struct {
	ID      string
	UserID  string
	Balance float64
}

// CreateAccount creates a new account for the given user ID with an initial balance of 0.
func (svc *TransactionService) CreateAccount(ctx context.Context, userID string) (Account, error) {
	query := "INSERT INTO accounts (id, user_id, balance) VALUES ($1, $2, $3)"

	accountId, err := svc.idgen.New()
	if err != nil {
		return Account{}, fmt.Errorf("Failed to generate id for account: %w", err)
	}

	// Execute the query without expecting a returned row
	_, err = svc.db.Exec(query, accountId, userID, 0.0) // Initial balance is set to 0.0
	if err != nil {
		return Account{}, fmt.Errorf("failed to insert account: %w", err)
	}

	// Return the Account struct with the ID
	return Account{
		ID:      string(accountId), // Convert the ID to a string if needed
		UserID:  userID,
		Balance: 0.0,
	}, nil
}

// GetAccount retrieves the account details for the given account ID.
func (svc *TransactionService) GetAccount(ctx context.Context, id string) (Account, error) {
	query := "SELECT id, user_id, balance FROM accounts WHERE id = $1"

	var account Account
	err := svc.db.QueryRow(query, id).Scan(&account.ID, &account.UserID, &account.Balance)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

// DepositFunds adds the specified amount to the account's balance.
// The operation is idempotent, identified by the idempotency key.
func (svc *TransactionService) DepositFunds(ctx context.Context, accountID string, amount float64, idempotencyKey *string) (Account, error) {
	tx, err := svc.db.Begin(ctx, nil)
	if err != nil {
		return Account{}, err
	}
	defer tx.Rollback()

	// Check if a transaction with the same idempotency key has already been created
	_, err = svc.createTransaction(accountID, "", amount, idempotencyKey, tx)
	if err != nil {
		return Account{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	var account Account
	query := "SELECT id, user_id, balance FROM accounts WHERE id = $1"
	err = tx.QueryRow(query, accountID).Scan(&account.ID, &account.UserID, &account.Balance)
	if err != nil {
		return Account{}, err
	}

	// Update the balance
	newBalance := account.Balance + amount
	query = "UPDATE accounts SET balance = $1 WHERE id = $2"
	_, err = tx.Exec(query, newBalance, accountID)
	if err != nil {
		return Account{}, fmt.Errorf("failed to update account balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return Account{}, err
	}

	account.Balance = newBalance
	return account, nil
}

// WithdrawFunds subtracts the specified amount from the account's balance.
// The operation is idempotent, identified by the idempotency key.
func (svc *TransactionService) WithdrawFunds(ctx context.Context, accountID string, amount float64, idempotencyKey *string) (Account, error) {
	tx, err := svc.db.Begin(ctx, nil)
	if err != nil {
		return Account{}, err
	}
	defer tx.Rollback()

	// Check if a transaction with the same idempotency key has already been created
	_, err = svc.createTransaction("", accountID, amount, idempotencyKey, tx)
	if err != nil {
		return Account{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	var account Account
	query := "SELECT id, user_id, balance FROM accounts WHERE id = $1"
	err = tx.QueryRow(query, accountID).Scan(&account.ID, &account.UserID, &account.Balance)
	if err != nil {
		return Account{}, err
	}

	// Check if the account has sufficient balance
	if account.Balance < amount {
		return Account{}, fmt.Errorf("insufficient funds. Balance: %.2f, Withdrawal: %.2f", account.Balance, amount)
	}

	// Update the balance
	newBalance := account.Balance - amount
	query = "UPDATE accounts SET balance = $1 WHERE id = $2"
	_, err = tx.Exec(query, newBalance, accountID)
	if err != nil {
		return Account{}, err
	}

	if err := tx.Commit(); err != nil {
		return Account{}, err
	}

	account.Balance = newBalance
	return account, nil
}
