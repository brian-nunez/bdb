package postgres

import "fmt"

// DriverName is the registered name of this database driver.
const DriverName = "postgres"

// Config contains configuration options for the Postgres driver.
type Config struct {
	// Host is the host address of the database.
	// If empty, defaults to "localhost".
	Host string

	// Port is the port of the database.
	// If 0, defaults to 5432.
	Port int

	// User is the database user.
	User string

	// Password is the password for the database user.
	Password string

	// DBName is the database name.
	DBName string

	// SSLMode controls SSL communication.
	// If empty, defaults to "disable".
	SSLMode string

	// DSN is the connection string (Data Source Name).
	// If provided, Host, Port, User, Password, DBName, and SSLMode are ignored.
	DSN string
}

// DriverName returns the name of this driver.
func (Config) DriverName() string {
	return DriverName
}

func (c Config) getDSN() string {
	if c.DSN != "" {
		return c.DSN
	}
	host := c.Host
	if host == "" {
		host = "localhost"
	}
	port := c.Port
	if port == 0 {
		port = 5432
	}
	sslmode := c.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, c.User, c.Password, c.DBName, sslmode)
}
