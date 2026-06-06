package main

import (
	"context"
	"fmt"
	"log"

	"github.com/brian-nunez/bdb"
	"github.com/brian-nunez/bdb/drivers/sqlite"
)

func sqlite_example() {
	log.Println("--- SQLite Example ---")
	conn, err := bdb.New(sqlite.Config{
		Path: "sqlite.db",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()

	// Perform database operations
	_, err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// Clean table for consistent runs
	_, _ = conn.Exec(ctx, "DELETE FROM users")

	_, err = conn.Exec(ctx, "INSERT INTO users (name) VALUES (?)", "Alice")
	if err != nil {
		log.Fatal(err)
	}

	var name string
	err = conn.QueryOne(ctx, "SELECT name FROM users WHERE id = 1", nil, &name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Got name from SQLite: %s\n", name)

	// Show connection stats
	fmt.Printf("SQLite Open Connections: %d\n", conn.OpenConnections())
}
