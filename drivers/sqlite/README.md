# bdb SQLite Database Driver

`sqlite` is a SQL-backed database connection wrapper for the `bdb` abstraction, utilizing the pure-Go SQLite driver `modernc.org/sqlite` (no CGO required).

---

## Installation

```bash
go get github.com/brian-nunez/bdb/drivers/sqlite
```

## Features

- **CGO-Free**: Runs entirely in pure Go, making cross-compilation simple and straightforward.
- **In-Memory & File-Backed Support**: Supports transient in-memory databases (e.g. `":memory:"`) as well as local database files (e.g. `"sqlite.db"`).
- **In-Memory Connection Pinning**: Automatically restricts database connections to `MaxOpenConns = 1` for in-memory databases (`path == ":memory:"`) to ensure that all database operations refer to the same shared in-memory database workspace.
- **Connection Pools**: Implements all `sql.DB` connection pool configurations and statistics methods exposed by the `bdb.DB` interface.

## Configuration

The driver uses the [sqlite.Config](./config.go) struct:

| Field | Type | Description |
|---|---|---|
| `Path` | `string` | SQLite database path. E.g., `"mydb.db"` or `":memory:"`. If empty, defaults to `":memory:"`. |

```go
type Config struct {
	Path string
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
	"github.com/brian-nunez/bdb/drivers/sqlite"
)

func main() {
	// Initialize SQLite connection pool wrapper
	db, err := bdb.New(sqlite.Config{
		Path: "app.db",
	})
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Perform database commands
	_, err = db.Exec(ctx, "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Fatalf("Table creation failed: %v", err)
	}

	_, err = db.Exec(ctx, "INSERT INTO users (name) VALUES (?)", "Bob")
	if err != nil {
		log.Fatalf("Insert failed: %v", err)
	}

	var name string
	err = db.QueryOne(ctx, "SELECT name FROM users WHERE id = ?", []any{1}, &name)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	fmt.Printf("User: %s\n", name)
}
```
