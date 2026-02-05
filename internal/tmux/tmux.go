package tmux

import (
	"strings"
)

// Session represents a tmux session.
type Session struct {
	Name string
}

// ParseSessionList parses tmux list-sessions output and returns only cb: prefixed sessions.
func ParseSessionList(output string) []Session {
	var sessions []Session
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		// Only include cb: prefixed sessions
		if !strings.HasPrefix(line, "cb:") {
			continue
		}

		// Parse: "cb:proj-123-auth: 3 windows (created ...)"
		// Session name is everything before the colon-space pattern " N windows"
		colonSpace := strings.Index(line, ": ")
		if colonSpace == -1 {
			continue
		}
		name := line[:colonSpace]

		sessions = append(sessions, Session{
			Name: name,
		})
	}

	return sessions
}
