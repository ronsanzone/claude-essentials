package tmux

import (
	"fmt"
	"os/exec"
	"strings"
)

// Session represents a tmux session.
type Session struct {
	Name string
}

// Client provides tmux operations.
type Client struct {
	execCommand func(name string, args ...string) ([]byte, error)
}

// NewClient creates a Client that executes real tmux commands.
func NewClient() *Client {
	return &Client{
		execCommand: func(name string, args ...string) ([]byte, error) {
			return exec.Command(name, args...).Output()
		},
	}
}

// ListSessions returns all ClawdBay tmux sessions.
func (c *Client) ListSessions() ([]Session, error) {
	output, err := c.execCommand("tmux", "list-sessions")
	if err != nil {
		// tmux not running or no sessions is expected, return empty list
		errMsg := err.Error()
		if strings.Contains(errMsg, "no server running") ||
			strings.Contains(errMsg, "no sessions") {
			return []Session{}, nil
		}
		return nil, fmt.Errorf("failed to list tmux sessions: %w", err)
	}
	return ParseSessionList(string(output)), nil
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
