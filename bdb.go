package bdb

import (
	"context"
	"database/sql"
	"time"
)

// DB represents a database connection wrapper.
type DB interface {
	Ping(ctx context.Context) error
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Close() error
	Prepare(ctx context.Context, query string) (*sql.Stmt, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	QueryOne(ctx context.Context, query string, args []any, dest ...any) error

	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)

	OpenConnections() int
	MaxOpenConnections() int
	IdleConnections() int
	InUseConnections() int

	GetConnection() *sql.DB
}
