# bdb PostgreSQL Database Driver

`postgres` is a SQL-backed database connection wrapper for the `bdb` abstraction, utilizing the PostgreSQL driver `github.com/lib/pq`.

---

## Installation

```bash
go get github.com/brian-nunez/bdb/drivers/postgres
```

## Features

- **Production Ready**: Full support for connection pooling, transaction isolation levels, and SQL prepared statements via PostgreSQL.
- **TLS/SSL Encryption**: Configurable SSL mode (`disable`, `require`, `verify-ca`, `verify-full`) to secure in-transit data.
- **DSN Override**: Accepts raw connection string URI (Data Source Name) or structured parameters.
- **Connection Pools**: Implements all `sql.DB` connection pool configurations and statistics methods exposed by the `bdb.DB` interface.

## Configuration

The driver uses the [postgres.Config](./config.go) struct:

| Field | Type | Description |
|---|---|---|
| `Host` | `string` | Database server host address. Defaults to `"localhost"`. |
| `Port` | `int` | Database server port. Defaults to `5432`. |
| `User` | `string` | Database authentication user. |
| `Password`| `string` | Database authentication password. |
| `DBName` | `string` | Database name. |
| `SSLMode` | `string` | TLS configuration mode (e.g. `"disable"`, `"require"`). Defaults to `"disable"`. |
| `DSN` | `string` | Complete connection string (e.g. `"postgres://user:password@localhost:5432/mydb?sslmode=disable"`). If provided, all other fields are ignored. |

```go
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	DSN      string
}
```

## Usage Example

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
	// Initialize PostgreSQL connection pool wrapper
	db, err := bdb.New(postgres.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "secretpassword",
		DBName:   "app_db",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("Failed to open Postgres database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Query for user
	var name string
	err = db.QueryOne(ctx, "SELECT name FROM users WHERE id = $1", []any{1}, &name)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	fmt.Printf("User: %s\n", name)
}
```
