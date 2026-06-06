package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/brian-nunez/bdb"

	_ "github.com/lib/pq"
)

func init() {
	bdb.Register(DriverName, New)
}

var openDB = sql.Open

type postgresConn struct {
	db *sql.DB
}

// New creates and registers a new postgres connection wrapper.
func New(config any) (bdb.DB, error) {
	cfg, ok := config.(Config)
	if !ok {
		return nil, bdb.ErrInvalidConfig
	}

	dsn := cfg.getDSN()

	db, err := openDB("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &postgresConn{db: db}, nil
}

func (c *postgresConn) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c *postgresConn) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

func (c *postgresConn) Close() error {
	return c.db.Close()
}

func (c *postgresConn) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return c.db.PrepareContext(ctx, query)
}

func (c *postgresConn) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}

func (c *postgresConn) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, opts)
}

func (c *postgresConn) QueryOne(ctx context.Context, query string, args []any, dest ...any) error {
	return c.db.QueryRowContext(ctx, query, args...).Scan(dest...)
}

func (c *postgresConn) SetMaxIdleConns(n int) {
	c.db.SetMaxIdleConns(n)
}

func (c *postgresConn) SetMaxOpenConns(n int) {
	c.db.SetMaxOpenConns(n)
}

func (c *postgresConn) SetConnMaxLifetime(d time.Duration) {
	c.db.SetConnMaxLifetime(d)
}

func (c *postgresConn) SetConnMaxIdleTime(d time.Duration) {
	c.db.SetConnMaxIdleTime(d)
}

func (c *postgresConn) OpenConnections() int {
	return c.db.Stats().OpenConnections
}

func (c *postgresConn) MaxOpenConnections() int {
	return c.db.Stats().MaxOpenConnections
}

func (c *postgresConn) IdleConnections() int {
	return c.db.Stats().Idle
}

func (c *postgresConn) InUseConnections() int {
	return c.db.Stats().InUse
}

func (c *postgresConn) GetConnection() *sql.DB {
	return c.db
}
