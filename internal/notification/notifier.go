// Package notification provides desktop notification functionality.
package notification

import (
	"fmt"
	"log/slog"

	"github.com/gen2brain/beeep"
	"github.com/yourusername/rest-time-reminder-go/internal/config"
)

// Notifier handles desktop notifications for the reminder.
type Notifier struct {
	config config.NotificationConfig
}

// NewNotifier creates a new Notifier instance.
func NewNotifier(cfg config.NotificationConfig) *Notifier {
	return &Notifier{
		config: cfg,
	}
}

// Notify displays a desktop notification if enabled.
func (n *Notifier) Notify() error {
	if !n.config.Desktop {
		slog.Debug("desktop notifications disabled, skipping")
		return nil
	}

	slog.Debug("showing desktop notification",
		"title", n.config.Title,
		"message", n.config.Message,
	)

	// Show notification using beeep
	// Empty string for icon will use system default
	if err := beeep.Notify(n.config.Title, n.config.Message, ""); err != nil {
		return fmt.Errorf("failed to show notification: %w", err)
	}

	return nil
}

// Alert displays an alert notification (more prominent than Notify).
func (n *Notifier) Alert(title, message string) error {
	if err := beeep.Alert(title, message, ""); err != nil {
		return fmt.Errorf("failed to show alert: %w", err)
	}
	return nil
}
