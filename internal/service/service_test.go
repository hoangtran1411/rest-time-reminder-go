package service

import (
	"testing"

	"github.com/hoangtran1411/rest-time-reminder-go/internal/config"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{}
	svc := New(cfg)
	if svc == nil {
		t.Fatal("New returned nil")
	}
}
