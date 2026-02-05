package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds ClawdBay configuration paths.
type Config struct {
	ConfigDir  string
	PromptsDir string
}

// New creates a Config with default paths.
func New() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	configDir := filepath.Join(home, ".config", "cb")

	return &Config{
		ConfigDir:  configDir,
		PromptsDir: filepath.Join(configDir, "prompts"),
	}, nil
}
