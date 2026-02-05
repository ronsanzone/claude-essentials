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
