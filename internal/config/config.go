// Package config handles configuration loading and management
// for the RestTimeReminder application.
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the complete application configuration.
type Config struct {
	Reminder     ReminderConfig     `mapstructure:"reminder"`
	Sound        SoundConfig        `mapstructure:"sound"`
	Notification NotificationConfig `mapstructure:"notification"`
	Logging      LoggingConfig      `mapstructure:"logging"`
	Service      ServiceConfig      `mapstructure:"service"`
}

// ReminderConfig holds settings for the reminder scheduler.
type ReminderConfig struct {
	// Interval between reminders (e.g., "30m", "1h")
	Interval string `mapstructure:"interval"`
	// TriggerMinutes are specific minutes to trigger (e.g., [0, 30])
	TriggerMinutes []int `mapstructure:"trigger_minutes"`
}

// SoundConfig holds settings for audio playback.
type SoundConfig struct {
	// Enabled indicates whether sound notifications are enabled
	Enabled bool `mapstructure:"enabled"`
	// File is the path to a custom sound file (empty for embedded)
	File string `mapstructure:"file"`
	// Volume is the playback volume (0.0 - 1.0)
	Volume float64 `mapstructure:"volume"`
}

// NotificationConfig holds settings for desktop notifications.
type NotificationConfig struct {
	// Desktop indicates whether desktop notifications are enabled
	Desktop bool `mapstructure:"desktop"`
	// Title is the notification title
	Title string `mapstructure:"title"`
	// Message is the notification message body
	Message string `mapstructure:"message"`
}

// LoggingConfig holds settings for application logging.
type LoggingConfig struct {
	// Level is the log level (debug, info, warn, error)
	Level string `mapstructure:"level"`
	// File is the log file path (empty for stdout only)
	File string `mapstructure:"file"`
}

// ServiceConfig holds settings for the system service.
type ServiceConfig struct {
	// DisplayName is the service display name (Windows)
	DisplayName string `mapstructure:"display_name"`
	// Description is the service description
	Description string `mapstructure:"description"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Reminder: ReminderConfig{
			Interval:       "30m",
			TriggerMinutes: nil,
		},
		Sound: SoundConfig{
			Enabled: true,
			File:    "",
			Volume:  1.0,
		},
		Notification: NotificationConfig{
			Desktop: false,
			Title:   "Break Time!",
			Message: "Time to take a short break and rest your eyes.",
		},
		Logging: LoggingConfig{
			Level: "info",
			File:  "",
		},
		Service: ServiceConfig{
			DisplayName: "Rest Time Reminder",
			Description: "A background service that reminds you to take regular breaks",
		},
	}
}

// Load reads configuration from the specified file or default locations.
// It returns the loaded configuration merged with defaults.
func Load(configFile string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Configuration file settings
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.rest-time-reminder")
		v.AddConfigPath("/etc/rest-time-reminder")
	}

	// Environment variable support
	v.SetEnvPrefix("RTR")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		// If user explicitly specified a config file, failure to read it is an error
		if configFile != "" {
			return nil, fmt.Errorf("error reading config file %q: %w", configFile, err)
		}

		// Otherwise, it's okay if the default config file doesn't exist; use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal configuration
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error parsing configuration: %w", err)
	}

	return cfg, nil
}

// setDefaults sets default values in the viper instance.
func setDefaults(v *viper.Viper) {
	defaults := DefaultConfig()

	v.SetDefault("reminder.interval", defaults.Reminder.Interval)
	v.SetDefault("sound.enabled", defaults.Sound.Enabled)
	v.SetDefault("sound.volume", defaults.Sound.Volume)
	v.SetDefault("notification.desktop", defaults.Notification.Desktop)
	v.SetDefault("notification.title", defaults.Notification.Title)
	v.SetDefault("notification.message", defaults.Notification.Message)
	v.SetDefault("logging.level", defaults.Logging.Level)
	v.SetDefault("service.display_name", defaults.Service.DisplayName)
	v.SetDefault("service.description", defaults.Service.Description)
}
