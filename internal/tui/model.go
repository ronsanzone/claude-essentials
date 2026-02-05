package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsanzone/clawdbay/internal/tmux"
)

// WorktreeGroup represents a worktree with its Claude sessions.
type WorktreeGroup struct {
	Name         string
	SessionCount int
	Sessions     []ClaudeSession
	Expanded     bool
}

// ClaudeSession represents a single Claude session window.
type ClaudeSession struct {
	Name       string
	Status     string
	LastActive time.Time
}

// Model is the Bubbletea model for the dashboard.
type Model struct {
	Groups   []WorktreeGroup
	Cursor   int
	Quitting bool
}

// GroupByWorktree groups sessions and their Claude windows.
func GroupByWorktree(sessions []tmux.Session, windows map[string][]tmux.Window, tmuxClient *tmux.Client) []WorktreeGroup {
	var groups []WorktreeGroup

	for _, session := range sessions {
		wins := windows[session.Name]

		var claudeSessions []ClaudeSession
		for _, w := range wins {
			if strings.HasPrefix(w.Name, "claude:") {
				status := string(tmux.StatusIdle)
				if tmuxClient != nil {
					status = string(tmuxClient.GetPaneStatus(session.Name, w.Name))
				}
				claudeSessions = append(claudeSessions, ClaudeSession{
					Name:   w.Name,
					Status: status,
				})
			}
		}

		groups = append(groups, WorktreeGroup{
			Name:         session.Name,
			SessionCount: len(claudeSessions),
			Sessions:     claudeSessions,
			Expanded:     true,
		})
	}

	return groups
}

// InitialModel creates the initial dashboard model.
func InitialModel() Model {
	return Model{
		Groups: []WorktreeGroup{},
		Cursor: 0,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Groups)-1 {
				m.Cursor++
			}
		case "enter":
			// Would attach to session
			return m, tea.Quit
		}
	}
	return m, nil
}
