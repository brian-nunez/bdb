package sqlite

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/brian-nunez/bdb"

	_ "modernc.org/sqlite"
)

func init() {
	bdb.Register(DriverName, New)
}

var openDB = sql.Open

type sqliteConn struct {
	db *sql.DB
}

// New creates and registers a new sqlite connection wrapper.
func New(config any) (bdb.DB, error) {
	cfg, ok := config.(Config)
	if !ok {
		return nil, bdb.ErrInvalidConfig
	}

	path := strings.TrimSpace(cfg.Path)
	if path == "" {
		path = ":memory:"
	}

	db, err := openDB("sqlite", path)
	if err != nil {
		return nil, err
	}

	if path == ":memory:" {
		db.SetMaxOpenConns(1)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &sqliteConn{db: db}, nil
}

func (c *sqliteConn) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c *sqliteConn) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

func (c *sqliteConn) Close() error {
	return c.db.Close()
}

func (c *sqliteConn) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return c.db.PrepareContext(ctx, query)
}

func (c *sqliteConn) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}

func (c *sqliteConn) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, opts)
}

func (c *sqliteConn) QueryOne(ctx context.Context, query string, args []any, dest ...any) error {
	return c.db.QueryRowContext(ctx, query, args...).Scan(dest...)
}

func (c *sqliteConn) SetMaxIdleConns(n int) {
	c.db.SetMaxIdleConns(n)
}

func (c *sqliteConn) SetMaxOpenConns(n int) {
	c.db.SetMaxOpenConns(n)
}

func (c *sqliteConn) SetConnMaxLifetime(d time.Duration) {
	c.db.SetConnMaxLifetime(d)
}

func (c *sqliteConn) SetConnMaxIdleTime(d time.Duration) {
	c.db.SetConnMaxIdleTime(d)
}

func (c *sqliteConn) OpenConnections() int {
	return c.db.Stats().OpenConnections
}

func (c *sqliteConn) MaxOpenConnections() int {
	return c.db.Stats().MaxOpenConnections
}

func (c *sqliteConn) IdleConnections() int {
	return c.db.Stats().Idle
}

func (c *sqliteConn) InUseConnections() int {
	return c.db.Stats().InUse
}

func (c *sqliteConn) GetConnection() *sql.DB {
	return c.db
}
