package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170"))

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("212"))

	idleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))

	workingStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("82"))

	doneStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))
)

// View implements tea.Model.
func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render("- ClawdBay ") + strings.Repeat("-", 50) + "\n\n")

	if len(m.Groups) == 0 {
		b.WriteString("  No active workflows.\n")
		b.WriteString("  Start one with: cb start <ticket-id>\n")
	} else {
		for i, group := range m.Groups {
			cursor := "  "
			if i == m.Cursor {
				cursor = "> "
			}

			// Group header
			expandIcon := "v"
			if !group.Expanded {
				expandIcon = ">"
			}

			line := fmt.Sprintf("%s%s %s    %d sessions",
				cursor, expandIcon, group.Name, group.SessionCount)

			if i == m.Cursor {
				b.WriteString(selectedStyle.Render(line) + "\n")
			} else {
				b.WriteString(line + "\n")
			}

			// Sessions
			if group.Expanded {
				for _, session := range group.Sessions {
					statusIcon := "*"
					var statusStyle lipgloss.Style

					switch session.Status {
					case "IDLE":
						statusStyle = idleStyle
					case "WORKING":
						statusStyle = workingStyle
					case "DONE":
						statusStyle = doneStyle
					default:
						statusStyle = idleStyle
					}

					sessionLine := fmt.Sprintf("      %s %s",
						statusStyle.Render(statusIcon+" "+session.Status),
						session.Name)
					b.WriteString(sessionLine + "\n")
				}
			}
			b.WriteString("\n")
		}
	}

	// Footer
	b.WriteString("\n")
	b.WriteString("  [Enter] Attach  [n] New  [c] Add Claude  [p] Add Prompt  [x] Archive  [q] Quit\n")

	return b.String()
}
