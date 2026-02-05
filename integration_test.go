//go:build integration

package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCLI_Version(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "version")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	if !strings.Contains(string(output), "ClawdBay") {
		t.Errorf("output = %q, want to contain 'ClawdBay'", output)
	}
}

func TestCLI_Help(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("help command failed: %v", err)
	}

	expected := []string{"start", "claude", "prompt", "list", "archive", "dash"}
	for _, sub := range expected {
		if !strings.Contains(string(output), sub) {
			t.Errorf("help missing subcommand: %s", sub)
		}
	}
}
