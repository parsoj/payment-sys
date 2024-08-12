package domain

import (
	"context"
	"database/sql"
)

type PostgresDB interface {
	Begin(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Exec(string, ...interface{}) (sql.Result, error)
	QueryRow(string, ...interface{}) *sql.Row
	Query(string, ...interface{}) (*sql.Rows, error)
	Ping() error
	Prepare(string) (*sql.Stmt, error)
	Close() error
}
