package mysql

import (
	"context"
	"database/sql"
)

// interface definitions borrowed from github.com/volatiletech/sqlboiler

// executor can perform SQL queries.
type executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// contextExecutor can perform SQL queries with context
type contextExecutor interface {
	executor

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// beginner begins transactions.
type beginner interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	contextExecutor
}
