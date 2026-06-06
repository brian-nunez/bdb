package sqlite

import "testing"

func TestConfigDriverName(t *testing.T) {
	cfg := Config{}

	got := cfg.DriverName()
	if got != DriverName {
		t.Fatalf("expected driver name %q, got %q", DriverName, got)
	}
}

func TestDriverNameConstant(t *testing.T) {
	if DriverName != "sqlite" {
		t.Fatalf("expected DriverName to be %q, got %q", "sqlite", DriverName)
	}
}

func TestConfigFields(t *testing.T) {
	cfg := Config{
		Path: "/tmp/test.db",
	}

	if cfg.Path != "/tmp/test.db" {
		t.Fatalf("expected Path %q, got %q", "/tmp/test.db", cfg.Path)
	}
}
