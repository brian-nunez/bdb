# bdb

`bdb` is a small driver-based database connection abstraction for Go.

It provides a single common interface for database operations while allowing different database backends to be plugged in through drivers.

Current drivers:

- `sqlite` — SQLite driver using pure-Go `modernc.org/sqlite`
- `postgres` — PostgreSQL driver using `github.com/lib/pq`
- `mariadb` — MariaDB/MySQL driver using `github.com/go-sql-driver/mysql`

## Features

- Driver registration pattern
- SQLite support (in-memory & file-based)
- PostgreSQL support
- MariaDB/MySQL support
- Connection pool statistics (open, max, idle, in-use)
- Underlying connection exposure (`GetConnection`)
- 100% test branch coverage

## Install

```bash
go get github.com/brian-nunez/bdb
```

Install the drivers you want to use:

```bash
go get github.com/brian-nunez/bdb/drivers/sqlite
go get github.com/brian-nunez/bdb/drivers/postgres
go get github.com/brian-nunez/bdb/drivers/mariadb
```

## Usage

### SQLite

Import the root package and the SQLite driver.

The SQLite driver registers itself when imported.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/brian-nunez/bdb"
	"github.com/brian-nunez/bdb/drivers/sqlite"
)

func main() {
	conn, err := bdb.New(sqlite.Config{
		Path: "sqlite.db",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()

	_, err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(ctx, "INSERT INTO users (name) VALUES (?)", "Alice")
	if err != nil {
		log.Fatal(err)
	}

	var name string
	err = conn.QueryOne(ctx, "SELECT name FROM users WHERE id = 1", nil, &name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User:", name)
}
```

### PostgreSQL

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/brian-nunez/bdb"
	"github.com/brian-nunez/bdb/drivers/postgres"
)

func main() {
	conn, err := bdb.New(postgres.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpassword",
		DBName:   "testdb",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()

	var name string
	err = conn.QueryOne(ctx, "SELECT name FROM users WHERE id = $1", []any{1}, &name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User:", name)
}
```

### MariaDB

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/brian-nunez/bdb"
	"github.com/brian-nunez/bdb/drivers/mariadb"
)

func main() {
	conn, err := bdb.New(mariadb.Config{
		Host:     "localhost",
		Port:     3306,
		User:     "testuser",
		Password: "testpassword",
		DBName:   "testdb",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()

	var name string
	err = conn.QueryOne(ctx, "SELECT name FROM users WHERE id = ?", []any{1}, &name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User:", name)
}
```

## Interface

All drivers implement the same `bdb.DB` interface.

```go
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
```

## Driver Registration

Drivers register themselves with the root `bdb` package.

Example:

```go
func init() {
	bdb.Register(DriverName, New)
}
```

This allows the app to create connections through the root package using the registration pattern.

## Running the Examples

Start the database servers and run the examples:

```bash
cd examples
make run
```

Stop the database servers and remove volumes:

```bash
make clean
```

## License

MIT
