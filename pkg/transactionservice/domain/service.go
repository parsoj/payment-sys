package domain

import (
	"github.com/parsoj/paymentsys/pkg/id"
)

type TransactionService struct {
	db    PostgresDB
	idgen *id.SpecialIdGenerator
}

func NewTransactionService(db PostgresDB) (*TransactionService, error) {

	return &TransactionService{
		db:    db,
		idgen: &id.SpecialIdGenerator{},
	}, nil
}
