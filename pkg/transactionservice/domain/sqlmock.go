package domain

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const schemaFilePath = "../../../db/schema_sqlite.sql"

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB() (*SQLiteDB, error) {
	// Open an in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Read the schema SQL from the external file
	schemaSQL, err := ioutil.ReadFile(filepath.Clean(schemaFilePath))
	if err != nil {
		db.Close() // Close the database if there's an error during setup
		return nil, fmt.Errorf("failed to read SQL schema file: %w", err)
	}

	// Execute the SQL setup commands from the file
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		db.Close() // Close the database if there's an error during setup
		return nil, fmt.Errorf("failed to set up SQLite database schema: %w", err)
	}

	return &SQLiteDB{db: db}, nil
}

// Begin starts a transaction with the given options
func (s *SQLiteDB) Begin(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, opts)
}

// Exec executes a query without returning any rows
func (s *SQLiteDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

// QueryRow executes a query that is expected to return at most one row
func (s *SQLiteDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.db.QueryRow(query, args...)
}

// Query executes a query that returns rows, typically a SELECT
func (s *SQLiteDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

// Ping verifies a connection to the database is still alive
func (s *SQLiteDB) Ping() error {
	return s.db.Ping()
}

// Prepare creates a prepared statement for later queries or executions
func (s *SQLiteDB) Prepare(query string) (*sql.Stmt, error) {
	return s.db.Prepare(query)
}

// Close closes the database and prevents new queries from starting
func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

// DumpTable prints the contents of the specified table.
func (sdb *SQLiteDB) DumpTable(tableName string) {
	rows, err := sdb.db.Query(fmt.Sprintf("SELECT * FROM %s;", tableName))
	if err != nil {
		log.Printf("Error querying table %s: %v", tableName, err)
		return
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("Error getting columns for table %s: %v", tableName, err)
		return
	}

	fmt.Printf("Table: %s\n", tableName)

	// Create a slice of interface{}'s to represent each column, and a second slice to contain pointers to each item in the columns slice
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Iterate through the rows
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Printf("Error scanning row in table %s: %v", tableName, err)
			return
		}

		for i, col := range columns {
			val := values[i]

			var v interface{}
			switch val := val.(type) {
			case []byte:
				v = string(val)
			default:
				v = val
			}

			fmt.Printf("%s: %v\n", col, v)
		}
		fmt.Println("-------------------------------")
	}
}

// DumpAllTables dumps the contents of all specified tables in the database.
func (sdb *SQLiteDB) DumpAllTables() {
	tables := []string{"users", "accounts", "transactions"}

	for _, table := range tables {
		sdb.DumpTable(table)
	}
}
