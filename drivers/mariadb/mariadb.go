package mariadb

import (
	"context"
	"database/sql"
	"time"

	"github.com/brian-nunez/bdb"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	bdb.Register(DriverName, New)
}

var openDB = sql.Open

type mariadbConn struct {
	db *sql.DB
}

// New creates and registers a new mariadb connection wrapper.
func New(config any) (bdb.DB, error) {
	cfg, ok := config.(Config)
	if !ok {
		return nil, bdb.ErrInvalidConfig
	}

	dsn := cfg.getDSN()

	db, err := openDB("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &mariadbConn{db: db}, nil
}

func (c *mariadbConn) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c *mariadbConn) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

func (c *mariadbConn) Close() error {
	return c.db.Close()
}

func (c *mariadbConn) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return c.db.PrepareContext(ctx, query)
}

func (c *mariadbConn) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}

func (c *mariadbConn) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, opts)
}

func (c *mariadbConn) QueryOne(ctx context.Context, query string, args []any, dest ...any) error {
	return c.db.QueryRowContext(ctx, query, args...).Scan(dest...)
}

func (c *mariadbConn) SetMaxIdleConns(n int) {
	c.db.SetMaxIdleConns(n)
}

func (c *mariadbConn) SetMaxOpenConns(n int) {
	c.db.SetMaxOpenConns(n)
}

func (c *mariadbConn) SetConnMaxLifetime(d time.Duration) {
	c.db.SetConnMaxLifetime(d)
}

func (c *mariadbConn) SetConnMaxIdleTime(d time.Duration) {
	c.db.SetConnMaxIdleTime(d)
}

func (c *mariadbConn) OpenConnections() int {
	return c.db.Stats().OpenConnections
}

func (c *mariadbConn) MaxOpenConnections() int {
	return c.db.Stats().MaxOpenConnections
}

func (c *mariadbConn) IdleConnections() int {
	return c.db.Stats().Idle
}

func (c *mariadbConn) InUseConnections() int {
	return c.db.Stats().InUse
}

func (c *mariadbConn) GetConnection() *sql.DB {
	return c.db
}
