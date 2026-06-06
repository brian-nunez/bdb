package bdb

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
)

type testConfig struct {
	name string
}

func (c testConfig) DriverName() string {
	return c.name
}

type testDB struct{}

func (t *testDB) Ping(ctx context.Context) error { return nil }
func (t *testDB) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return nil, nil
}
func (t *testDB) Close() error                                                 { return nil }
func (t *testDB) Prepare(ctx context.Context, query string) (*sql.Stmt, error) { return nil, nil }
func (t *testDB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return nil, nil
}
func (t *testDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) { return nil, nil }
func (t *testDB) QueryOne(ctx context.Context, query string, args []any, dest ...any) error {
	return nil
}
func (t *testDB) SetMaxIdleConns(n int)              {}
func (t *testDB) SetMaxOpenConns(n int)              {}
func (t *testDB) SetConnMaxLifetime(d time.Duration) {}
func (t *testDB) SetConnMaxIdleTime(d time.Duration) {}
func (t *testDB) OpenConnections() int               { return 0 }
func (t *testDB) MaxOpenConnections() int            { return 0 }
func (t *testDB) IdleConnections() int               { return 0 }
func (t *testDB) InUseConnections() int              { return 0 }
func (t *testDB) GetConnection() *sql.DB             { return nil }

func TestRegisterAndNew(t *testing.T) {
	originalDrivers := drivers
	drivers = make(map[string]Driver)
	defer func() {
		drivers = originalDrivers
	}()

	var mockDriver Driver = func(config any) (DB, error) {
		return &testDB{}, nil
	}

	Register("mock", mockDriver)

	db, err := New(testConfig{name: "mock"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if db == nil {
		t.Fatal("expected db instance, got nil")
	}

	_, err = New(testConfig{name: "unknown"})
	if !errors.Is(err, ErrUnknownDriver) {
		t.Fatalf("expected ErrUnknownDriver, got %v", err)
	}
}
