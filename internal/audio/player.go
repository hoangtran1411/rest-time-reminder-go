// Package audio provides audio playback functionality for the reminder.
package audio

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/hoangtran1411/rest-time-reminder-go/internal/config"
)

//go:embed bell.wav
var defaultSound []byte

// Player handles audio playback for reminder notifications.
type Player struct {
	config   config.SoundConfig
	initOnce sync.Once
	initErr  error
	mu       sync.Mutex
}

// NewPlayer creates a new Player instance.
func NewPlayer(cfg config.SoundConfig) *Player {
	return &Player{
		config: cfg,
	}
}

// Play plays the configured sound file.
func (p *Player) Play() error {
	if !p.config.Enabled {
		slog.Debug("sound is disabled, skipping playback")
		return nil
	}

	// Ensure only one sound plays at a time
	p.mu.Lock()
	defer p.mu.Unlock()

	var streamer beep.StreamSeekCloser
	var format beep.Format
	var err error

	// 1. Try to load from custom file if configured
	// 2. Fallback to embedded sound if file is "bell.wav" or empty
	if p.config.File != "" && p.config.File != "bell.wav" {
		streamer, format, err = p.loadFromFile(p.config.File)
		if err != nil {
			slog.Warn("failed to load custom sound, falling back to default", "path", p.config.File, "error", err)
			streamer = nil // Reset to ensure fallback
		}
	}

	// Fallback to embedded sound
	if streamer == nil {
		streamer, format, err = p.loadEmbedded()
		if err != nil {
			return fmt.Errorf("failed to load default sound: %w", err)
		}
	}
	defer func() { _ = streamer.Close() }()

	// Initialize speaker if not already done (thread-safe)
	p.initOnce.Do(func() {
		// Use a buffer size of 1/10th of a second
		if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
			p.initErr = fmt.Errorf("failed to initialize speaker: %w", err)
		}
	})

	if p.initErr != nil {
		return p.initErr
	}

	// Apply Volume control
	// Beep volume is logarithmic: 0 is silent, 1 is 100%, but it uses values like -1.0, -2.0...
	// We map 0.0-1.0 to Beep volume (where 1.0 is original, < 1.0 is quieter)
	// Base 2: volume of -1.0 is 50%, -2.0 is 25%, etc.
	// We use the formula: Volume = log2(p.config.Volume)
	vol := 0.0
	if p.config.Volume > 0 && p.config.Volume < 1.0 {
		// Example: Volume 0.5 -> log2(0.5) = -1.0 (50% volume)
		// Volume 0.25 -> log2(0.25) = -2.0 (25% volume)
		// We can simplify or use math.Log2, but for simplicity let's use a basic mapping:
		// We'll use p.config.Volume directly for a linear-ish feel if we set Base properly,
		// or just use 0.0 for 100% and negative values for lower.
		// A common way in beep is volume - 1.0 if using Base 2 roughly.
		// Let's use a simpler mapping: Volume 1.0 -> 0, Volume 0.5 -> -1, etc.
		// If volume is 0.5, we want it to be -1.0
		// But let's just use the built-in Gain if available or Volume effect.
		vol = p.config.Volume - 1.0 // Simple mapping: 1.0 -> 0 (normal), 0.5 -> -0.5 (quieter)
	}

	volumeControl := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   vol,
		Silent:   p.config.Volume <= 0,
	}

	// Play the sound
	done := make(chan bool)
	speaker.Play(beep.Seq(volumeControl, beep.Callback(func() {
		done <- true
	})))

	// Wait for playback to complete (or context cancel if implemented)
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
		_ = f.Close()
		return nil, beep.Format{}, fmt.Errorf("failed to decode WAV: %w", err)
	}

	return streamer, format, nil
}

// loadEmbedded loads the embedded WAV sound.
func (p *Player) loadEmbedded() (beep.StreamSeekCloser, beep.Format, error) {
	// Create a NopCloser because bytes.Reader doesn't have Close()
	reader := io.NopCloser(bytes.NewReader(defaultSound))
	streamer, format, err := wav.Decode(reader)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("failed to decode embedded WAV: %w", err)
	}
	return streamer, format, nil
}
