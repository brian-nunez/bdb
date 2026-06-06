# bdb MariaDB / MySQL Database Driver

`mariadb` is a SQL-backed database connection wrapper for the `bdb` abstraction, utilizing the official MySQL driver `github.com/go-sql-driver/mysql`.

---

## Installation

```bash
go get github.com/brian-nunez/bdb/drivers/mariadb
```

## Features

- **High Compatibility**: Fully compatible with both MySQL and MariaDB servers.
- **Custom Connection Parameters**: Pass raw MySQL connection flags (e.g. `parseTime=true`, `loc=Local`) via the `Params` configuration map.
- **DSN Override**: Accepts raw connection string (Data Source Name) or structured parameters.
- **Connection Pools**: Implements all `sql.DB` connection pool configurations and statistics methods exposed by the `bdb.DB` interface.

## Configuration

The driver uses the [mariadb.Config](./config.go) struct:

| Field | Type | Description |
|---|---|---|
| `Host` | `string` | Database server host address. Defaults to `"localhost"`. |
| `Port` | `int` | Database server port. Defaults to `3306`. |
| `User` | `string` | Database authentication user. |
| `Password`| `string` | Database authentication password. |
| `DBName` | `string` | Database name. |
| `Params` | `map[string]string` | Key-value pairs of connection parameters (e.g., `parseTime=true`, `timeout=5s`). |
| `DSN` | `string` | Complete connection DSN (e.g., `"user:password@tcp(127.0.0.1:3306)/dbname?parseTime=true"`). If provided, all other fields are ignored. |

```go
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Params   map[string]string
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
	"github.com/brian-nunez/bdb/drivers/mariadb"
)

func main() {
	// Initialize MariaDB connection pool wrapper
	db, err := bdb.New(mariadb.Config{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "rootpassword",
		DBName:   "app_db",
		Params: map[string]string{
			"parseTime": "true",
			"charset":   "utf8mb4",
		},
	})
	if err != nil {
		log.Fatalf("Failed to open MariaDB database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Query for user
	var name string
	err = db.QueryOne(ctx, "SELECT name FROM users WHERE id = ?", []any{1}, &name)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	fmt.Printf("User: %s\n", name)
}
```
