package tui

import (
	"testing"

	"github.com/rsanzone/clawdbay/internal/tmux"
)

func TestModel_GroupSessionsByWorktree(t *testing.T) {
	sessions := []tmux.Session{
		{Name: "cb:proj-123-auth"},
		{Name: "cb:proj-456-bug"},
	}

	windows := map[string][]tmux.Window{
		"cb:proj-123-auth": {
			{Name: "shell"},
			{Name: "claude:default"},
			{Name: "claude:research"},
		},
		"cb:proj-456-bug": {
			{Name: "shell"},
			{Name: "claude:default"},
		},
	}

	// Pass nil for tmuxClient in tests (status detection skipped)
	groups := GroupByWorktree(sessions, windows, nil)

	if len(groups) != 2 {
		t.Fatalf("got %d groups, want 2", len(groups))
	}

	if groups[0].SessionCount != 2 {
		t.Errorf("first group sessions = %d, want 2", groups[0].SessionCount)
	}
}
