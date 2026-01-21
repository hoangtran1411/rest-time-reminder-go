package updater

import (
	"testing"
)

func TestCheck_InvalidVersion(t *testing.T) {
	_, _, err := Check("invalid-version", "owner/repo")
	if err == nil {
		t.Error("expected error for invalid version, got nil")
	}
}

func TestUpdate_InvalidVersion(t *testing.T) {
	err := Update("invalid-version", "owner/repo")
	if err == nil {
		t.Error("expected error for invalid version, got nil")
	}
}
