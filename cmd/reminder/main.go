// Package main is the entry point for the RestTimeReminder application.
// It provides a CLI interface for running the reminder as a console app,
// system service, or with system tray support.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/yourusername/rest-time-reminder-go/internal/audio"
	"github.com/yourusername/rest-time-reminder-go/internal/config"
	"github.com/yourusername/rest-time-reminder-go/internal/notification"
	"github.com/yourusername/rest-time-reminder-go/internal/scheduler"
	"github.com/yourusername/rest-time-reminder-go/internal/service"
	"github.com/yourusername/rest-time-reminder-go/internal/updater"
)

var (
	// Version information (set via ldflags during build)
	// Default to 0.0.0 for dev builds to allow semantic version parsing
	version   = "0.0.0" 
	commit    = "none"
	buildDate = "unknown"

	// Repository for updates
	repoSlug = "yourusername/rest-time-reminder-go"

	// CLI flags
	cfgFile  string
	interval string
	sound    string
	verbose  bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "rest-time-reminder",
		Short: "A reminder application for taking regular breaks",
		Long: `RestTimeReminder is a lightweight, cross-platform application
that reminds you to take breaks at regular intervals by playing
a pleasant bell sound.

Run without arguments to start in console mode, or use service
commands to install and manage as a system service.`,
		Run: runApp,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default: config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&interval, "interval", "i", "", "reminder interval (e.g., 30m, 1h)")
	rootCmd.PersistentFlags().StringVarP(&sound, "sound", "s", "", "path to custom sound file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")

	// Version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("RestTimeReminder %s\n", version)
			fmt.Printf("  Commit: %s\n", commit)
			fmt.Printf("  Built: %s\n", buildDate)
			
			// Check for updates
			if latest, found, err := updater.Check(version, repoSlug); err == nil && found {
				fmt.Printf("\nUpdate available: %s\n", latest.Version)
				fmt.Println("Run 'rest-time-reminder update' to upgrade.")
			}
		},
	})

	// Update command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "update",
		Short: "Update to the latest version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := updater.Update(version, repoSlug); err != nil {
				slog.Error("update failed", "error", err)
				os.Exit(1)
			}
		},
	})

	// Service commands
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "install",
			Short: "Install as system service",
			Run:   func(cmd *cobra.Command, args []string) { runServiceCommand("install") },
		},
		&cobra.Command{
			Use:   "uninstall",
			Short: "Remove system service",
			Run:   func(cmd *cobra.Command, args []string) { runServiceCommand("uninstall") },
		},
		&cobra.Command{
			Use:   "start",
			Short: "Start the service",
			Run:   func(cmd *cobra.Command, args []string) { runServiceCommand("start") },
		},
		&cobra.Command{
			Use:   "stop",
			Short: "Stop the service",
			Run:   func(cmd *cobra.Command, args []string) { runServiceCommand("stop") },
		},
		&cobra.Command{
			Use:   "status",
			Short: "Show service status",
			Run:   func(cmd *cobra.Command, args []string) { runServiceCommand("status") },
		},
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// runApp is the main application logic for console mode
func runApp(cmd *cobra.Command, args []string) {
	// Setup logging
	logLevel := slog.LevelInfo
	if verbose {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Override config with CLI flags
	if interval != "" {
		cfg.Reminder.Interval = interval
	}
	if sound != "" {
		cfg.Sound.File = sound
	}

	slog.Info("starting RestTimeReminder",
		"version", version,
		"interval", cfg.Reminder.Interval,
	)

	// Update check in background
	go func() {
		if _, found, err := updater.Check(version, repoSlug); err == nil && found {
			slog.Info("new version available", "command", "rest-time-reminder update")
		}
	}()

	// Initialize components
	player := audio.NewPlayer(cfg.Sound)
	notifier := notification.NewNotifier(cfg.Notification)
	sched := scheduler.New(cfg.Reminder, player, notifier)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		slog.Info("received shutdown signal", "signal", sig)
		cancel()
	}()

	// Run the scheduler
	if err := sched.Run(ctx); err != nil {
		slog.Error("scheduler error", "error", err)
		os.Exit(1)
	}

	slog.Info("RestTimeReminder stopped gracefully")
}

// runServiceCommand handles service management commands
func runServiceCommand(command string) {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	svc := service.New(cfg)
	if err := svc.Execute(command); err != nil {
		slog.Error("service command failed", "command", command, "error", err)
		os.Exit(1)
	}
}
