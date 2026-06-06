package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brian-nunez/bdb"
	"github.com/brian-nunez/bdb/drivers/mariadb"
)

func mariadb_example() {
	log.Println("--- MariaDB Example ---")

	var conn bdb.DB
	var err error

	// Retry connection up to 15 times to allow the container database to initialize
	for i := 0; i < 15; i++ {
		conn, err = bdb.New(mariadb.Config{
			Host:     "localhost",
			Port:     3306,
			User:     "testuser",
			Password: "testpassword",
			DBName:   "testdb",
		})
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Println("Could not connect to MariaDB (make sure it's running via Docker):", err)
		return
	}
	defer conn.Close()

	ctx := context.Background()

	// Perform database operations
	_, err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS users (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255))")
	if err != nil {
		log.Fatal(err)
	}

	// Clean table for consistent runs
	_, _ = conn.Exec(ctx, "TRUNCATE TABLE users")

	_, err = conn.Exec(ctx, "INSERT INTO users (name) VALUES (?)", "Charlie")
	if err != nil {
		log.Fatal(err)
	}

	var name string
	err = conn.QueryOne(ctx, "SELECT name FROM users WHERE id = 1", nil, &name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Got name from MariaDB: %s\n", name)

	// Show connection stats
	fmt.Printf("MariaDB Open Connections: %d\n", conn.OpenConnections())
}
