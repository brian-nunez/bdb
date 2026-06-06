package mariadb

import "testing"

func TestConfigDriverName(t *testing.T) {
	cfg := Config{}

	got := cfg.DriverName()
	if got != DriverName {
		t.Fatalf("expected driver name %q, got %q", DriverName, got)
	}
}

func TestDriverNameConstant(t *testing.T) {
	if DriverName != "mariadb" {
		t.Fatalf("expected DriverName to be %q, got %q", "mariadb", DriverName)
	}
}

func TestConfigGetDSN(t *testing.T) {
	cfg := Config{
		User:     "testuser",
		Password: "testpassword",
		DBName:   "testdb",
	}

	expected := "testuser:testpassword@tcp(localhost:3306)/testdb"
	got := cfg.getDSN()
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}

	cfgCustom := Config{
		Host:     "127.0.0.1",
		Port:     3307,
		User:     "testuseronly",
		DBName:   "testdb",
		Params: map[string]string{
			"parseTime": "true",
			"charset":   "utf8mb4",
		},
	}
	expectedCustom := "testuseronly@tcp(127.0.0.1:3307)/testdb?charset=utf8mb4&parseTime=true"
	gotCustom := cfgCustom.getDSN()
	if gotCustom != expectedCustom {
		t.Fatalf("expected %q, got %q", expectedCustom, gotCustom)
	}

	cfgDSN := Config{
		DSN: "testuser:testpassword@tcp(127.0.0.1:3306)/testdb",
	}
	gotDSN := cfgDSN.getDSN()
	if gotDSN != cfgDSN.DSN {
		t.Fatalf("expected %q, got %q", cfgDSN.DSN, gotDSN)
	}
}
