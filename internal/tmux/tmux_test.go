package tmux

import (
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
