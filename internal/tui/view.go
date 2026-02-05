package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rsanzone/clawdbay/internal/tmux"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170"))

	selectedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212"))

	repoStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("141"))

	sessionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	windowStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	idleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))

	workingStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("82"))

	doneStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
)

// View implements tea.Model.
func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render("- ClawdBay ") + strings.Repeat("-", 50) + "\n\n")

	if len(m.Nodes) == 0 {
		b.WriteString("  No active sessions.\n")
		b.WriteString("  Start one with: cb start <branch-name>\n")
	} else {
		for i, node := range m.Nodes {
			cursor := "  "
			if i == m.Cursor {
				cursor = "> "
			}

			line := m.renderNode(node, cursor)

			if i == m.Cursor {
				b.WriteString(selectedStyle.Render(line) + "\n")
			} else {
				b.WriteString(line + "\n")
			}
		}
	}

	// Dynamic footer based on selected node type
	b.WriteString("\n")
	footer := m.renderFooter()
	b.WriteString(footerStyle.Render(footer) + "\n")

	return b.String()
}

func (m Model) renderNode(node TreeNode, cursor string) string {
	switch node.Type {
	case NodeRepo:
		repo := m.Groups[node.RepoIndex]
		icon := "▸"
		if repo.Expanded {
			icon = "▼"
		}
		return fmt.Sprintf("%s%s %s", cursor, icon, repoStyle.Render(repo.Name))

	case NodeSession:
		session := m.Groups[node.RepoIndex].Sessions[node.SessionIndex]
		icon := "▸"
		if session.Expanded {
			icon = "▼"
		}
		statusBadge := renderStatus(session.Status)
		// Right-align status badge (pad to 60 chars)
		name := fmt.Sprintf("  %s %s", icon, session.Name)
		padding := 60 - len(name) - len(statusBadge)
		if padding < 1 {
			padding = 1
		}
		return fmt.Sprintf("%s%s%s%s", cursor, name, strings.Repeat(" ", padding), statusBadge)

	case NodeWindow:
		session := m.Groups[node.RepoIndex].Sessions[node.SessionIndex]
		window := session.Windows[node.WindowIndex]
		statusBadge := ""
		if strings.HasPrefix(window.Name, "claude") {
			key := session.Name + ":" + window.Name
			if status, ok := m.WindowStatuses[key]; ok {
				statusBadge = renderStatus(status)
			}
		}
		// Right-align status badge (pad to 60 chars)
		name := fmt.Sprintf("      %s", window.Name)
		padding := 60 - len(name) - len(statusBadge)
		if padding < 1 {
			padding = 1
		}
		return fmt.Sprintf("%s%s%s%s", cursor, name, strings.Repeat(" ", padding), statusBadge)

	default:
		return cursor + "Unknown"
	}
}

func (m Model) renderFooter() string {
	if m.Cursor >= len(m.Nodes) {
		return "[q] quit"
	}

	node := m.Nodes[m.Cursor]
	switch node.Type {
	case NodeRepo:
		return "[enter] expand  [n] new  [q] quit"
	case NodeSession:
		return "[enter] attach  [c] add claude  [x] archive  [r] refresh  [q] quit"
	case NodeWindow:
		return "[enter] attach  [r] refresh  [q] quit"
	default:
		return "[q] quit"
	}
}

func renderStatus(status tmux.Status) string {
	switch status {
	case tmux.StatusWorking:
		return workingStyle.Render("● WORKING")
	case tmux.StatusIdle:
		return idleStyle.Render("○ IDLE")
	case tmux.StatusDone:
		return doneStyle.Render("◌ DONE")
	default:
		return doneStyle.Render("◌ DONE")
	}
}
