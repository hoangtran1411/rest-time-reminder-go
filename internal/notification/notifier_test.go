package notification

import (
	"testing"

	"github.com/hoangtran1411/rest-time-reminder-go/internal/config"
)

func TestNewNotifier(t *testing.T) {
	cfg := config.NotificationConfig{
		Title:   "Test Title",
		Message: "Test Message",
		Desktop: true,
	}
	n := NewNotifier(cfg)
	if n == nil {
		t.Fatal("NewNotifier returned nil")
	}
	if n.config != cfg {
		t.Errorf("expected config %+v, got %+v", cfg, n.config)
	}
}

func TestNotifier_Notify_Disabled(t *testing.T) {
	cfg := config.NotificationConfig{
		Desktop: false,
	}
	n := NewNotifier(cfg)
	err := n.Notify()
	if err != nil {
		t.Errorf("expected nil error when disabled, got %v", err)
	}
}

// Note: Testing Notify() when enabled requires a desktop environment or mocking beeep,
// which is not straightforward without refactoring.
// Skipping "Enabled" test to avoid CI failure in headless environments.
