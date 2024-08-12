package domain

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type TransactionState string

const (
	Completed TransactionState = "Completed"
	Pending   TransactionState = "Pending"
	Cancelled TransactionState = "Cancelled"
)

type Transaction struct {
	ID             string
	ToAccount      string
	FromAccount    string
	Amount         float64
	State          TransactionState
	IdempotencyKey string
	CreatedAt      time.Time
}

func (svc *TransactionService) GetTransaction(ctx context.Context, id string) (Transaction, error) {
	query := `
        SELECT id, to_account, from_account, amount, state
        FROM transactions
        WHERE id = $1
    `

	var transaction Transaction
	err := svc.db.QueryRow(query, id).Scan(
		&transaction.ID,
		&transaction.ToAccount,
		&transaction.FromAccount,
		&transaction.Amount,
		&transaction.State,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Transaction{}, fmt.Errorf("transaction with id %s not found", id)
		}
		return Transaction{}, err
	}

	return transaction, nil
}

func (svc *TransactionService) TransferFunds(ctx context.Context, toAccount, fromAccount string, amount float64, idempotencyKey *string) (string, error) {

	if fromAccount == "" {
		return "", fmt.Errorf("Invalid account id for source account: '%s'", fromAccount)
	}

	if toAccount == "" {
		return "", fmt.Errorf("Invalid account id for destination account: '%s'", toAccount)
	}

	txn_id, err := svc.createTransaction(toAccount, fromAccount, amount, idempotencyKey, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create Transaction row: %w", err)
	}

	tx, err := svc.db.Begin(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to open DB transaction: %w", err)
	}
	defer tx.Rollback()

	//********************************************************************************
	//fetch & check account balances
	var balances map[string]float64
	var queryErr error

	// Define the query to fetch balances for both accounts in a single query
	query := `
    SELECT id, balance 
    FROM accounts 
    WHERE id IN ($1, $2);
  `

	rows, queryErr := tx.Query(query, fromAccount, toAccount)
	if queryErr != nil {
		return txn_id, fmt.Errorf("Failed to execute balance query: %w", queryErr)
	}
	defer rows.Close()

	balances = make(map[string]float64)

	for rows.Next() {
		var accountID string
		var balance float64
		if err := rows.Scan(&accountID, &balance); err != nil {
			return txn_id, fmt.Errorf("Failed to scan row: %w", err)
		}
		balances[accountID] = balance
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return txn_id, fmt.Errorf("Error iterating over rows: %w", err)
	}

	toBalance, toFound := balances[toAccount]
	fromBalance, fromFound := balances[fromAccount]

	if !toFound {
		return txn_id, fmt.Errorf("Failed to fetch balance for account '%s'", toAccount)
	}

	if !fromFound {
		return txn_id, fmt.Errorf("Failed to fetch balance for account '%s'", fromAccount)
	}

	if amount > fromBalance {
		return txn_id, fmt.Errorf("Insufficient funds in source account. Balance: %f -- Transfer Amount: %f", fromBalance, amount)
	}
	//********************************************************************************
	// perform the transfer

	fromNewBalance := fromBalance - amount
	toNewBalance := toBalance + amount

	updateQuery := `
    UPDATE accounts
    SET balance = CASE 
        WHEN id = $1 THEN $2::numeric
        WHEN id = $3 THEN $4::numeric
    END
    WHERE id IN ($1, $3);
`

	_, updateErr := tx.Exec(updateQuery, fromAccount, fromNewBalance, toAccount, toNewBalance)
	if updateErr != nil {
		return txn_id, fmt.Errorf("Failed to update balances for accounts: %w", updateErr)
	}
	//********************************************************************************
	// mark transaction as completed

	err = svc.completeTransaction(tx, txn_id)
	if err != nil {
		return txn_id, fmt.Errorf("Failed to mark transaction as completed: %w", err)
	}

	tx.Commit()

	return txn_id, nil

}

//********************************************************************************
// helper functions

func (svc *TransactionService) fetchAccountBalance(tx *sql.Tx, accountID string) (float64, error) {
	var balance float64

	query := "SELECT balance FROM accounts WHERE id = $1"
	err := tx.QueryRow(query, accountID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no account found with id %s", accountID)
		}
		return 0, err
	}

	return balance, nil
}

func (svc *TransactionService) setBalance(tx *sql.Tx, accountID string, newBalance float64) error {
	query := "UPDATE accounts SET balance = $1 WHERE id = $2"
	result, err := tx.Exec(query, newBalance, accountID)
	if err != nil {
		return fmt.Errorf("failed to set balance for account %s: %v", accountID, err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve result for account %s: %v", accountID, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no account found with id %s", accountID)
	}

	return nil
}

func (svc *TransactionService) createTransaction(toAccount, fromAccount string, amount float64, idempotencyKey *string, tx *sql.Tx) (string, error) {
	transactionID, err := svc.idgen.New()
	if err != nil {
		return "", fmt.Errorf("failed to generate id for transaction: %w", err)
	}

	createdAt := time.Now()

	var idp_key string
	if idempotencyKey == nil {
		idp, err := svc.idgen.New()
		if err != nil {
			return "", fmt.Errorf("Failed to generate idempotencyKey: %w", err)
		}
		idp_key = string(idp)
	} else {
		idp_key = *idempotencyKey
	}

	query := `
        INSERT INTO transactions (id, to_account, from_account, amount, state, idempotency_key, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	if tx == nil {
		_, err = svc.db.Exec(query, transactionID, toAccount, fromAccount, amount, Pending, idp_key, createdAt)
	} else {

		_, err = tx.Exec(query, transactionID, toAccount, fromAccount, amount, Pending, idp_key, createdAt)
	}
	if err != nil {
		return "", fmt.Errorf("failed to insert transaction: %v", err)
	}

	return string(transactionID), nil
}

func (svc *TransactionService) cancelTransaction(tx *sql.Tx, transactionID string) error {
	query := `
        UPDATE transactions SET state = $1 WHERE id = $2
    `
	_, err := tx.Exec(query, Cancelled, transactionID)
	if err != nil {
		return fmt.Errorf("failed to cancel transaction with id %s: %v", transactionID, err)
	}
	return nil
}

func (svc *TransactionService) completeTransaction(tx *sql.Tx, transactionID string) error {
	query := `
        UPDATE transactions SET state = $1 WHERE id = $2
    `
	_, err := tx.Exec(query, Completed, transactionID)
	if err != nil {
		return fmt.Errorf("failed to complete transaction with id %s: %v", transactionID, err)
	}
	return nil
}
