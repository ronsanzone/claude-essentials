package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() error = %v", err)
	}
	expected := filepath.Join(home, ".config", "cb")

	cfg, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if cfg.ConfigDir != expected {
		t.Errorf("ConfigDir = %q, want %q", cfg.ConfigDir, expected)
	}
}

func TestPromptsDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() error = %v", err)
	}
	expected := filepath.Join(home, ".config", "cb", "prompts")

	cfg, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if cfg.PromptsDir != expected {
		t.Errorf("PromptsDir = %q, want %q", cfg.PromptsDir, expected)
	}
}
