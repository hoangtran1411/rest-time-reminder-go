// Package audio provides audio playback functionality for the reminder.
package audio

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/hoangtran1411/rest-time-reminder-go/internal/config"
)

// Player handles audio playback for reminder notifications.
type Player struct {
	config      config.SoundConfig
	initialized bool
}

// NewPlayer creates a new Player instance.
func NewPlayer(cfg config.SoundConfig) *Player {
	return &Player{
		config:      cfg,
		initialized: false,
	}
}

// Play plays the configured sound file.
func (p *Player) Play() error {
	if !p.config.Enabled {
		slog.Debug("sound is disabled, skipping playback")
		return nil
	}

	// If no sound file configured, just log it
	if p.config.File == "" {
		slog.Info("🔔 reminder triggered (no sound file configured)")
		return nil
	}

	streamer, format, err := p.loadFromFile(p.config.File)
	if err != nil {
		return fmt.Errorf("failed to load sound: %w", err)
	}
	defer streamer.Close()

	// Initialize speaker if not already done
	if !p.initialized {
		if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
			return fmt.Errorf("failed to initialize speaker: %w", err)
		}
		p.initialized = true
	}

	// Play the sound
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Wait for playback to complete
	<-done
	slog.Debug("sound playback completed")

	return nil
}

// Stop stops any currently playing sound.
func (p *Player) Stop() {
	speaker.Clear()
}

// loadFromFile loads a WAV file from the filesystem.
func (p *Player) loadFromFile(path string) (beep.StreamSeekCloser, beep.Format, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("invalid path: %w", err)
	}

	f, err := os.Open(absPath)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("failed to open file: %w", err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		f.Close()
		return nil, beep.Format{}, fmt.Errorf("failed to decode WAV: %w", err)
	}

	return streamer, format, nil
}
