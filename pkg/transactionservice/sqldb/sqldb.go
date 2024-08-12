package sqldb

import (
	"context"
	"database/sql"
)

type SqlDB struct {
	db *sql.DB
}

// NewSqlDB is the constructor that initializes a connection to PostgreSQL using database/sql.
func NewSqlDB(connString string) (*SqlDB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// Optionally, ping the database to verify the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &SqlDB{db: db}, nil
}

func (db *SqlDB) Begin(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.db.BeginTx(ctx, opts)
}

func (db *SqlDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

func (db *SqlDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.db.QueryRow(query, args...)
}

func (db *SqlDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.Query(query, args...)
}

func (db *SqlDB) Ping() error {
	return db.db.Ping()
}

func (db *SqlDB) Prepare(query string) (*sql.Stmt, error) {
	return db.db.Prepare(query)
}

func (db *SqlDB) Close() error {
	return db.db.Close()
}
