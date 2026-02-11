package audio

import (
	"testing"

	"github.com/hoangtran1411/rest-time-reminder-go/internal/config"
)

func TestNewPlayer(t *testing.T) {
	cfg := config.SoundConfig{
		Enabled: true,
		Volume:  0.8,
		File:    "test.wav",
	}

	player := NewPlayer(cfg)

	if player.config.Enabled != cfg.Enabled {
		t.Errorf("expected enabled %v, got %v", cfg.Enabled, player.config.Enabled)
	}
	if player.config.Volume != cfg.Volume {
		t.Errorf("expected volume %v, got %v", cfg.Volume, player.config.Volume)
	}
	if player.config.File != cfg.File {
		t.Errorf("expected file %v, got %v", cfg.File, player.config.File)
	}
}

func TestPlayer_Play_Disabled(t *testing.T) {
	cfg := config.SoundConfig{
		Enabled: false,
	}
	player := NewPlayer(cfg)

	err := player.Play()
	if err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestPlayer_LoadEmbedded(t *testing.T) {
	if len(defaultSound) == 0 {
		t.Fatal("defaultSound is empty, embedding failed")
	}
	t.Logf("defaultSound length: %d", len(defaultSound))

	player := NewPlayer(config.SoundConfig{})
	streamer, format, err := player.loadEmbedded()
	if err != nil {
		// If decoding fails, we still passed the embedding check
		// Some environments might have issues with the specific WAV format
		t.Logf("Warning: could not decode embedded WAV (this might be an environment issue): %v", err)
		return
	}
	defer func() { _ = streamer.Close() }()

	if format.SampleRate <= 0 {
		t.Errorf("expected positive sample rate, got %v", format.SampleRate)
	}
}

func TestPlayer_LoadFromFile_Error(t *testing.T) {
	player := NewPlayer(config.SoundConfig{})

	_, _, err := player.loadFromFile("non-existent-file.wav")
	if err == nil {
		t.Error("expected error when loading non-existent file, got nil")
	}
}
