package bdb

type Driver func(config any) (DB, error)

var drivers = make(map[string]Driver)

func Register(name string, driver Driver) {
	drivers[name] = driver
}

type NamedConfig interface {
	DriverName() string
}

func New(config NamedConfig) (DB, error) {
	driverName := config.DriverName()

	driver, ok := drivers[driverName]
	if !ok {
		return nil, ErrUnknownDriver
	}

	return driver(config)
}
