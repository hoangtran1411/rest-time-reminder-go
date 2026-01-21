// Package service provides Windows Service and Linux daemon support.
package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kardianos/service"
	"github.com/yourusername/rest-time-reminder-go/internal/audio"
	"github.com/yourusername/rest-time-reminder-go/internal/config"
	"github.com/yourusername/rest-time-reminder-go/internal/notification"
	"github.com/yourusername/rest-time-reminder-go/internal/scheduler"
)

// Service wraps the application for running as a system service.
type Service struct {
	config *config.Config
	svc    service.Service
	cancel context.CancelFunc
}

// program implements the service.Interface
type program struct {
	cfg      *config.Config
	cancel   context.CancelFunc
	doneChan chan struct{}
}

// New creates a new Service instance.
func New(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

// Execute runs the specified service command.
func (s *Service) Execute(command string) error {
	// Create service configuration
	svcConfig := &service.Config{
		Name:        "RestTimeReminder",
		DisplayName: s.config.Service.DisplayName,
		Description: s.config.Service.Description,
	}

	// Create program instance
	prg := &program{
		cfg:      s.config,
		doneChan: make(chan struct{}),
	}

	// Create the service
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	s.svc = svc

	// Execute the command
	switch command {
	case "install":
		return s.install()
	case "uninstall":
		return s.uninstall()
	case "start":
		return s.start()
	case "stop":
		return s.stop()
	case "status":
		return s.status()
	case "run":
		return svc.Run()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (s *Service) install() error {
	slog.Info("installing service...")
	if err := s.svc.Install(); err != nil {
		return fmt.Errorf("failed to install service: %w", err)
	}
	slog.Info("service installed successfully")
	return nil
}

func (s *Service) uninstall() error {
	slog.Info("uninstalling service...")
	if err := s.svc.Uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall service: %w", err)
	}
	slog.Info("service uninstalled successfully")
	return nil
}

func (s *Service) start() error {
	slog.Info("starting service...")
	if err := s.svc.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}
	slog.Info("service started successfully")
	return nil
}

func (s *Service) stop() error {
	slog.Info("stopping service...")
	if err := s.svc.Stop(); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}
	slog.Info("service stopped successfully")
	return nil
}

func (s *Service) status() error {
	status, err := s.svc.Status()
	if err != nil {
		return fmt.Errorf("failed to get service status: %w", err)
	}

	statusStr := "unknown"
	switch status {
	case service.StatusRunning:
		statusStr = "running"
	case service.StatusStopped:
		statusStr = "stopped"
	case service.StatusUnknown:
		statusStr = "unknown"
	}

	fmt.Printf("Service Status: %s\n", statusStr)
	return nil
}

// Start is called when the service starts
func (p *program) Start(s service.Service) error {
	slog.Info("service starting...")

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	// Initialize components
	player := audio.NewPlayer(p.cfg.Sound)
	notifier := notification.NewNotifier(p.cfg.Notification)
	sched := scheduler.New(p.cfg.Reminder, player, notifier)

	// Start scheduler in background
	go func() {
		defer close(p.doneChan)
		if err := sched.Run(ctx); err != nil {
			slog.Error("scheduler error", "error", err)
		}
	}()

	return nil
}

// Stop is called when the service stops
func (p *program) Stop(s service.Service) error {
	slog.Info("service stopping...")

	if p.cancel != nil {
		p.cancel()
	}

	// Wait for scheduler to finish
	<-p.doneChan

	slog.Info("service stopped")
	return nil
}
