package domain

import (
	"context"
	"fmt"
)

type User struct {
	ID       string
	Username string
}

func (svc *TransactionService) CreateUser(ctx context.Context, username string) (User, error) {
	query := "INSERT INTO users (username, id) VALUES ($1, $2) RETURNING id"

	userID, err := svc.idgen.New()
	if err != nil {
		return User{}, fmt.Errorf("failed to generate id for user: %w", err)
	}

	err = svc.db.QueryRow(query, username, userID).Scan(&userID)
	if err != nil {
		return User{}, err
	}

	createdUser := User{
		ID:       string(userID),
		Username: username,
	}

	return createdUser, nil
}

func (svc *TransactionService) GetUser(ctx context.Context, userID string) (User, error) {
	query := "SELECT id, username FROM users WHERE id = $1"

	var user User
	err := svc.db.QueryRow(query, userID).Scan(&user.ID, &user.Username)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
