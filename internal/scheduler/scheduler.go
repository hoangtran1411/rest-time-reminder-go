// Package scheduler implements time-based reminder scheduling logic.
package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/hoangtran1411/rest-time-reminder-go/internal/config"
)

// Player defines the interface for audio playback.
type Player interface {
	Play() error
	Stop()
}

// Notifier defines the interface for desktop notifications.
type Notifier interface {
	Notify() error
}

// Scheduler manages the reminder timing and triggers notifications.
type Scheduler struct {
	config   config.ReminderConfig
	player   Player
	notifier Notifier
	interval time.Duration
	lastPlay time.Time
	mu       sync.Mutex
}

// New creates a new Scheduler instance.
func New(cfg config.ReminderConfig, player Player, notifier Notifier) *Scheduler {
	return &Scheduler{
		config:   cfg,
		player:   player,
		notifier: notifier,
	}
}

// Run starts the scheduler loop and blocks until the context is cancelled.
func (s *Scheduler) Run(ctx context.Context) error {
	// Parse interval
	interval, err := time.ParseDuration(s.config.Interval)
	if err != nil {
		return fmt.Errorf("invalid interval format: %w", err)
	}
	s.interval = interval

	slog.Info("scheduler started",
		"interval", s.interval,
		"trigger_minutes", s.config.TriggerMinutes,
	)

	// Use 1-second ticker for precise timing
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("scheduler stopping")
			return nil
		case now := <-ticker.C:
			if s.shouldTrigger(now) {
				// Run trigger asynchronously to prevent blocking the loop
				go s.trigger(now)
			}
		}
	}
}

// shouldTrigger determines if a reminder should be triggered at the given time.
func (s *Scheduler) shouldTrigger(now time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if we already played in this minute
	// This is more robust than checking now.Second() == 0, as we might miss the exact second
	// under load, but we won't miss the minute.
	if !s.lastPlay.IsZero() && now.Truncate(time.Minute).Equal(s.lastPlay.Truncate(time.Minute)) {
		return false
	}

	// If specific trigger minutes are configured, check against them
	if len(s.config.TriggerMinutes) > 0 {
		currentMinute := now.Minute()
		for _, m := range s.config.TriggerMinutes {
			if currentMinute == m {
				return true
			}
		}
		return false
	}

	// Otherwise, use interval-based triggering
	minutes := now.Minute()
	intervalMinutes := int(s.interval.Minutes())
	if intervalMinutes <= 0 {
		intervalMinutes = 30 // Default to 30 minutes
	}

	// Trigger if the current minute aligns with the interval
	// Note: We don't check seconds here because the debouncing (lastPlay check) above handles it.
	return minutes%intervalMinutes == 0
}

// trigger executes the reminder notification.
func (s *Scheduler) trigger(now time.Time) {
	s.mu.Lock()
	s.lastPlay = now
	s.mu.Unlock()

	slog.Info("🔔 reminder triggered",
		"time", now.Format("15:04:05"),
	)

	// Play sound
	if err := s.player.Play(); err != nil {
		slog.Error("failed to play sound", "error", err)
	}

	// Show desktop notification
	if err := s.notifier.Notify(); err != nil {
		slog.Error("failed to show notification", "error", err)
	}
}
