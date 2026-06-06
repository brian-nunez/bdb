# bdb Database Drivers

This directory contains the database driver implementations for `bdb` (a lightweight, driver-based database connection abstraction for Go).

Drivers are separate Go packages that implement the [bdb.DB](../bdb.go) interface and register themselves dynamically with the core [bdb](../bdb.go) package.

---

## Existing Drivers

Each driver is hosted in its own package/sub-module. Click the links below to view their configurations and implementations:

1. **[sqlite](./sqlite)**
   - **Type**: SQL-backed database wrapper using pure-Go `modernc.org/sqlite` (CGO-free).
   - **Use Case**: Embedded database storage, testing, or serverless applications where a separate server is not needed.
   - **Configuration**: [sqlite.Config](./sqlite/config.go).
   - **Features**: Automatically creates and pins connection pools to `MaxOpenConns = 1` for in-memory (`:memory:`) databases to guarantee connection state alignment.

2. **[postgres](./postgres)**
   - **Type**: PostgreSQL client wrapper using `github.com/lib/pq`.
   - **Use Case**: Robust production deployments, relational data architectures, and transactional enterprise backends.
   - **Configuration**: [postgres.Config](./postgres/config.go).
   - **Features**: Supports connection string (DSN) overrides, host/port/SSL configurations, and PostgreSQL TLS integration.

3. **[mariadb](./mariadb)**
   - **Type**: MySQL/MariaDB client wrapper using `github.com/go-sql-driver/mysql`.
   - **Use Case**: High-performance read-heavy web apps, MySQL-based clusters, and legacy database systems.
   - **Configuration**: [mariadb.Config](./mariadb/config.go).
   - **Features**: Supports fully configurable query parameters (via `Params` map), DSN overrides, and structured host/credentials formatting.

---

## Driver Lifecycle & Registration

`bdb` drivers register themselves with the core registry in their package `init()` function:

```go
package mydriver

import "github.com/brian-nunez/bdb"

const DriverName = "mydriver"

func init() {
	bdb.Register(DriverName, New)
}
```

To use a driver, import it with a blank import to trigger `init()` and construct the connection with [bdb.New](../driver.go):

```go
import (
	"github.com/brian-nunez/bdb"
	_ "github.com/brian-nunez/bdb/drivers/postgres" // registers the postgres driver
)

func main() {
	db, err := bdb.New(postgres.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "mydb",
	})
	// ...
}
```

---

## How to Implement a Custom Driver

To add a new SQL database backend to `bdb`, follow these steps:

### 1. Implement `bdb.NamedConfig`
Define a configuration struct that contains connection parameters. It must implement the [bdb.NamedConfig](../driver.go) interface:

```go
type Config struct {
	DSN string
}

func (Config) DriverName() string {
	return "custom-sql-driver"
}
```

### 2. Implement the `bdb.DB` Interface
Create a connection wrapper type that implements the [bdb.DB](../bdb.go) interface:

```go
type customConn struct {
	db *sql.DB
}

func (c *customConn) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c *customConn) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

func (c *customConn) Close() error {
	return c.db.Close()
}

func (c *customConn) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return c.db.PrepareContext(ctx, query)
}

func (c *customConn) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}

func (c *customConn) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, opts)
}

func (c *customConn) QueryOne(ctx context.Context, query string, args []any, dest ...any) error {
	return c.db.QueryRowContext(ctx, query, args...).Scan(dest...)
}

// Implement connection pool wrappers
func (c *customConn) SetMaxIdleConns(n int)          { c.db.SetMaxIdleConns(n) }
func (c *customConn) SetMaxOpenConns(n int)          { c.db.SetMaxOpenConns(n) }
func (c *customConn) SetConnMaxLifetime(d time.Duration) { c.db.SetConnMaxLifetime(d) }
func (c *customConn) SetConnMaxIdleTime(d time.Duration) { c.db.SetConnMaxIdleTime(d) }

// Implement connection pool statistics
func (c *customConn) OpenConnections() int    { return c.db.Stats().OpenConnections }
func (c *customConn) MaxOpenConnections() int { return c.db.Stats().MaxOpenConnections }
func (c *customConn) IdleConnections() int    { return c.db.Stats().Idle }
func (c *customConn) InUseConnections() int   { return c.db.Stats().InUse }

// Expose the underlying *sql.DB connection
func (c *customConn) GetConnection() *sql.DB {
	return c.db
}
```

### 3. Provide the Constructor and Register the Driver
Define a factory function that takes an `any` interface, asserts it to your `Config`, opens the SQL connection, verifies it with a ping, and returns your connection wrapper:

```go
func New(config any) (bdb.DB, error) {
	cfg, ok := config.(Config)
	if !ok {
		return nil, bdb.ErrInvalidConfig
	}

	db, err := sql.Open("custom-driver", cfg.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &customConn{db: db}, nil
}

func init() {
	bdb.Register("custom-sql-driver", New)
}
```

---

## Shared Error Handling

Drivers should map their internal errors to the standard errors defined in [errors.go](../errors.go):

- Return [bdb.ErrInvalidConfig](../errors.go) in the factory constructor if config validation fails.
- Return [bdb.ErrDBClosed](../errors.go) if actions are attempted on a connection pool that has already been closed.
