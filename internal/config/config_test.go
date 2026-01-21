package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Reminder.Interval != "30m" {
		t.Errorf("expected default interval '30m', got %s", cfg.Reminder.Interval)
	}
	if !cfg.Sound.Enabled {
		t.Error("expected sound enabled by default")
	}
}

func TestLoad_CustomFile(t *testing.T) {
	// Create a temporary config file
	content := []byte(`
reminder:
  interval: 45m
sound:
  enabled: false
`)
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load the config
	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Reminder.Interval != "45m" {
		t.Errorf("expected interval '45m', got %s", cfg.Reminder.Interval)
	}
	if cfg.Sound.Enabled {
		t.Error("expected sound disabled")
	}
}
