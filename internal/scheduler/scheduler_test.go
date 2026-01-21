package scheduler

import (
	"testing"
	"time"

	"github.com/yourusername/rest-time-reminder-go/internal/config"
)

// MockPlayer implements Player interface for testing
type MockPlayer struct {
	PlayCount int
}

func (m *MockPlayer) Play() error {
	m.PlayCount++
	return nil
}
func (m *MockPlayer) Stop() {}

// MockNotifier implements Notifier interface for testing
type MockNotifier struct {
	NotifyCount int
}

func (m *MockNotifier) Notify() error {
	m.NotifyCount++
	return nil
}

func (m *MockNotifier) Alert(title, message string) error {
	return nil
}

func TestScheduler_shouldTrigger(t *testing.T) {
	tests := []struct {
		name           string
		cfg            config.ReminderConfig
		now            time.Time
		lastPlay       time.Time
		expectedResult bool
	}{
		{
			name: "Interval 30m - At :00 - Trigger",
			cfg:  config.ReminderConfig{Interval: "30m"},
			now:  time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedResult: true,
		},
		{
			name: "Interval 30m - At :30 - Trigger",
			cfg:  config.ReminderConfig{Interval: "30m"},
			now:  time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			expectedResult: true,
		},
		{
			name: "Interval 30m - At :15 - No Trigger",
			cfg:  config.ReminderConfig{Interval: "30m"},
			now:  time.Date(2023, 1, 1, 10, 15, 0, 0, time.UTC),
			expectedResult: false,
		},
		{
			name: "Interval 30m - Not Trigger Second != 0",
			cfg:  config.ReminderConfig{Interval: "30m"},
			now:  time.Date(2023, 1, 1, 10, 0, 5, 0, time.UTC),
			expectedResult: false,
		},
		{
			name: "Specific Trigger Minutes - At :00",
			cfg:  config.ReminderConfig{TriggerMinutes: []int{0, 15}},
			now:  time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedResult: true,
		},
		{
			name: "Specific Trigger Minutes - At :45 (Not in list)",
			cfg:  config.ReminderConfig{TriggerMinutes: []int{0, 15}},
			now:  time.Date(2023, 1, 1, 10, 45, 0, 0, time.UTC),
			expectedResult: false,
		},
		{
			name: "Double Trigger check (Played recently)",
			cfg:  config.ReminderConfig{Interval: "30m"},
			now:  time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			lastPlay: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC).Add(-10 * time.Second),
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.cfg, &MockPlayer{}, &MockNotifier{})
			
			// Manual setup for test since they are private fields in same package
			if tt.cfg.Interval != "" {
				d, _ := time.ParseDuration(tt.cfg.Interval)
				s.interval = d
			}
			s.lastPlay = tt.lastPlay
			
			if got := s.shouldTrigger(tt.now); got != tt.expectedResult {
				t.Errorf("shouldTrigger() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}
