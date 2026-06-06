package bdb

import "errors"

var (
	// ErrUnknownDriver is returned when a driver is requested that has not been registered.
	ErrUnknownDriver = errors.New("UNKNOWN_DRIVER")

	// ErrInvalidConfig is returned when a driver configuration is invalid.
	ErrInvalidConfig = errors.New("INVALID_CONFIG")

	// ErrDBClosed is returned when an operation is performed on a closed database connection.
	ErrDBClosed = errors.New("DB_CLOSED")
)
