package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brian-nunez/bdb"
)

func TestNewInvalidConfig(t *testing.T) {
	conn, err := New("bad config")
	if !errors.Is(err, bdb.ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}
	if conn != nil {
		t.Fatalf("expected nil conn, got %v", conn)
	}
}

func TestNewOpenFailure(t *testing.T) {
	expectedErr := errors.New("open failed")
	originalOpenDB := openDB
	defer func() { openDB = originalOpenDB }()

	openDB = func(driverName string, dataSourceName string) (*sql.DB, error) {
		return nil, expectedErr
	}

	conn, err := New(Config{DSN: "some-dsn"})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected open error, got %v", err)
	}
	if conn != nil {
		t.Fatalf("expected nil conn, got %v", conn)
	}
}

func TestNewPingFailure(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("sqlmock new failed: %v", err)
	}
	defer db.Close()

	originalOpenDB := openDB
	defer func() { openDB = originalOpenDB }()

	openDB = func(driverName string, dataSourceName string) (*sql.DB, error) {
		return db, nil
	}

	mock.ExpectPing().WillReturnError(errors.New("ping failed"))

	conn, err := New(Config{DSN: "some-dsn"})
	if err == nil {
		t.Fatal("expected ping error")
	}
	if conn != nil {
		t.Fatalf("expected nil conn, got %v", conn)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestNewSuccessAndInterfaceMethods(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("sqlmock new failed: %v", err)
	}
	defer db.Close()

	originalOpenDB := openDB
	defer func() { openDB = originalOpenDB }()

	openDB = func(driverName string, dataSourceName string) (*sql.DB, error) {
		if dataSourceName != "host=localhost port=5432 user=testuser password=testpassword dbname=testdb sslmode=disable" {
			t.Fatalf("unexpected DSN: %s", dataSourceName)
		}
		return db, nil
	}

	mock.ExpectPing().WillReturnError(nil)

	conn, err := New(Config{
		User:     "testuser",
		Password: "testpassword",
		DBName:   "testdb",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if conn == nil {
		t.Fatal("expected conn, got nil")
	}

	ctx := context.Background()

	// Test Ping
	mock.ExpectPing().WillReturnError(nil)
	if err := conn.Ping(ctx); err != nil {
		t.Fatalf("Ping failed: %v", err)
	}

	// Test Exec
	mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
	res, err := conn.Exec(ctx, "INSERT")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
	if res == nil {
		t.Fatal("expected result, got nil")
	}

	// Test Prepare
	mock.ExpectPrepare("SELECT")
	stmt, err := conn.Prepare(ctx, "SELECT")
	if err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}
	if stmt == nil {
		t.Fatal("expected statement, got nil")
	}
	stmt.Close()

	// Test Query
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}))
	rows, err := conn.Query(ctx, "SELECT")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if rows == nil {
		t.Fatal("expected rows, got nil")
	}
	rows.Close()

	// Test BeginTx
	mock.ExpectBegin()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("BeginTx failed: %v", err)
	}
	if tx == nil {
		t.Fatal("expected tx, got nil")
	}
	tx.Rollback()

	// Test QueryOne
	mock.ExpectQuery("SELECT ONE").WillReturnRows(sqlmock.NewRows([]string{"val"}).AddRow("hello"))
	var val string
	err = conn.QueryOne(ctx, "SELECT ONE", nil, &val)
	if err != nil {
		t.Fatalf("QueryOne failed: %v", err)
	}
	if val != "hello" {
		t.Fatalf("expected hello, got %s", val)
	}

	// Test Sets
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)
	conn.SetConnMaxIdleTime(time.Minute)

	// Test Stats
	stats := db.Stats()
	if conn.OpenConnections() != stats.OpenConnections {
		t.Fatalf("OpenConnections mismatch")
	}
	if conn.MaxOpenConnections() != stats.MaxOpenConnections {
		t.Fatalf("MaxOpenConnections mismatch")
	}
	if conn.IdleConnections() != stats.Idle {
		t.Fatalf("IdleConnections mismatch")
	}
	if conn.InUseConnections() != stats.InUse {
		t.Fatalf("InUseConnections mismatch")
	}

	// Test GetConnection
	if conn.GetConnection() != db {
		t.Fatalf("GetConnection mismatch")
	}

	// Test Close
	mock.ExpectClose()
	if err := conn.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
