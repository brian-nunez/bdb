package sqlite

// DriverName is the registered name of this database driver.
const DriverName = "sqlite"

// Config contains configuration options for the SQLite driver.
type Config struct {
	// Path is the SQLite database path.
	// Examples:
	//   Path: "mydb.db"
	//   Path: ":memory:"
	// If empty, ":memory:" is used.
	Path string
}

// DriverName returns the name of this driver.
func (Config) DriverName() string {
	return DriverName
}
