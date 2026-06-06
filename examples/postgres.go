package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brian-nunez/bdb"
	"github.com/brian-nunez/bdb/drivers/postgres"
)

func postgres_example() {
	log.Println("--- Postgres Example ---")

	var conn bdb.DB
	var err error

	// Retry connection up to 15 times to allow the container database to initialize
	for i := 0; i < 15; i++ {
		conn, err = bdb.New(postgres.Config{
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpassword",
			DBName:   "testdb",
			SSLMode:  "disable",
		})
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Println("Could not connect to Postgres (make sure it's running via Docker):", err)
		return
	}
	defer conn.Close()

	ctx := context.Background()

	// Perform database operations
	_, err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// Clean table for consistent runs
	_, _ = conn.Exec(ctx, "TRUNCATE TABLE users RESTART IDENTITY")

	_, err = conn.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "Bob")
	if err != nil {
		log.Fatal(err)
	}

	var name string
	err = conn.QueryOne(ctx, "SELECT name FROM users WHERE id = 1", nil, &name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Got name from Postgres: %s\n", name)

	// Show connection stats
	fmt.Printf("Postgres Open Connections: %d\n", conn.OpenConnections())
}
