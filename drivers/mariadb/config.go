package mariadb

import (
	"fmt"
	"sort"
	"strings"
)

// DriverName is the registered name of this database driver.
const DriverName = "mariadb"

// Config contains configuration options for the MariaDB driver.
type Config struct {
	// Host is the host address of the database.
	// If empty, defaults to "localhost".
	Host string

	// Port is the port of the database.
	// If 0, defaults to 3306.
	Port int

	// User is the database user.
	User string

	// Password is the password for the database user.
	Password string

	// DBName is the database name.
	DBName string

	// Params is a map of key-value connection parameters.
	Params map[string]string

	// DSN is the connection string (Data Source Name).
	// If provided, Host, Port, User, Password, DBName, and Params are ignored.
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
		port = 3306
	}

	var creds string
	if c.User != "" {
		if c.Password != "" {
			creds = c.User + ":" + c.Password + "@"
		} else {
			creds = c.User + "@"
		}
	}

	dsn := fmt.Sprintf("%stcp(%s:%d)/%s", creds, host, port, c.DBName)

	if len(c.Params) > 0 {
		var paramStrings []string
		for k, v := range c.Params {
			paramStrings = append(paramStrings, fmt.Sprintf("%s=%s", k, v))
		}
		sort.Strings(paramStrings)
		dsn = dsn + "?" + strings.Join(paramStrings, "&")
	}

	return dsn
}
