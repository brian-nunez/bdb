package postgres

import "testing"

func TestConfigDriverName(t *testing.T) {
	cfg := Config{}

	got := cfg.DriverName()
	if got != DriverName {
		t.Fatalf("expected driver name %q, got %q", DriverName, got)
	}
}

func TestDriverNameConstant(t *testing.T) {
	if DriverName != "postgres" {
		t.Fatalf("expected DriverName to be %q, got %q", "postgres", DriverName)
	}
}

func TestConfigGetDSN(t *testing.T) {
	cfg := Config{
		User:     "testuser",
		Password: "testpassword",
		DBName:   "testdb",
	}

	expected := "host=localhost port=5432 user=testuser password=testpassword dbname=testdb sslmode=disable"
	got := cfg.getDSN()
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}

	cfgCustom := Config{
		Host:     "127.0.0.1",
		Port:     5433,
		User:     "testuser",
		Password: "testpassword",
		DBName:   "testdb",
		SSLMode:  "require",
	}
	expectedCustom := "host=127.0.0.1 port=5433 user=testuser password=testpassword dbname=testdb sslmode=require"
	gotCustom := cfgCustom.getDSN()
	if gotCustom != expectedCustom {
		t.Fatalf("expected %q, got %q", expectedCustom, gotCustom)
	}

	cfgDSN := Config{
		DSN: "postgresql://testuser:testpassword@127.0.0.1:5432/testdb?sslmode=disable",
	}
	gotDSN := cfgDSN.getDSN()
	if gotDSN != cfgDSN.DSN {
		t.Fatalf("expected %q, got %q", cfgDSN.DSN, gotDSN)
	}
}
