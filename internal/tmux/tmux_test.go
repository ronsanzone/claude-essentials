package tmux

import (
	"errors"
	"testing"
)

func TestParseSessionList(t *testing.T) {
	output := `cb:proj-123-auth: 3 windows (created Tue Feb  4 10:30:00 2025)
cb:proj-456-bug: 1 windows (created Tue Feb  4 11:00:00 2025)
other-session: 2 windows (created Tue Feb  4 09:00:00 2025)`

	sessions := ParseSessionList(output)

	// Should only include cb: prefixed sessions
	if len(sessions) != 2 {
		t.Errorf("got %d sessions, want 2", len(sessions))
	}

	if sessions[0].Name != "cb:proj-123-auth" {
		t.Errorf("first session = %q, want %q", sessions[0].Name, "cb:proj-123-auth")
	}
}

func TestClient_ListSessions_Success(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) ([]byte, error) {
			return []byte(`cb:test-session: 1 windows (created Tue Feb  4 10:30:00 2025)
other-session: 2 windows (created Tue Feb  4 09:00:00 2025)`), nil
		},
	}

	sessions, err := client.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions() error = %v", err)
	}
	if len(sessions) != 1 {
		t.Errorf("got %d sessions, want 1", len(sessions))
	}
	if sessions[0].Name != "cb:test-session" {
		t.Errorf("session name = %q, want %q", sessions[0].Name, "cb:test-session")
	}
}

func TestClient_ListSessions_NoTmux(t *testing.T) {
	// Test graceful handling when tmux not running
	client := &Client{
		execCommand: func(name string, args ...string) ([]byte, error) {
			return nil, &mockError{msg: "no server running"}
		},
	}

	sessions, err := client.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions() error = %v, want nil", err)
	}
	if len(sessions) != 0 {
		t.Errorf("got %d sessions, want 0", len(sessions))
	}
}

type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}

func TestParseWindowList(t *testing.T) {
	// Format from: tmux list-windows -F "#{window_index}:#{window_name}:#{window_active}"
	output := `0:shell:1
1:claude:default:0
2:claude:research:0`

	windows := ParseWindowList(output)

	if len(windows) != 3 {
		t.Fatalf("got %d windows, want 3", len(windows))
	}

	if windows[0].Name != "shell" {
		t.Errorf("window 0 name = %q, want %q", windows[0].Name, "shell")
	}
	if !windows[0].Active {
		t.Error("window 0 should be active")
	}
	if windows[1].Name != "claude:default" {
		t.Errorf("window 1 name = %q, want %q", windows[1].Name, "claude:default")
	}
}

func TestWindow_IsClaudeSession(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"shell", false},
		{"claude:default", true},
		{"claude:research", true},
		{"vim", false},
	}

	for _, tt := range tests {
		w := Window{Name: tt.name}
		if got := w.IsClaudeSession(); got != tt.want {
			t.Errorf("Window{%q}.IsClaudeSession() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestClient_GetPaneStatus(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		err      error
		expected Status
	}{
		{"claude running", "claude", nil, StatusIdle},
		{"shell running", "zsh", nil, StatusDone},
		{"bash running", "bash", nil, StatusDone},
		{"error", "", errors.New("error"), StatusDone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				execCommand: func(name string, args ...string) ([]byte, error) {
					return []byte(tt.output), tt.err
				},
			}
			status := client.GetPaneStatus("session", "window")
			if status != tt.expected {
				t.Errorf("GetPaneStatus() = %v, want %v", status, tt.expected)
			}
		})
	}
}
