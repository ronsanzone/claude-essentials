package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsanzone/clawdbay/internal/tmux"
	"github.com/rsanzone/clawdbay/internal/tui"
	"github.com/spf13/cobra"
)

var dashCmd = &cobra.Command{
	Use:   "dash",
	Short: "Open interactive dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		tmuxClient := tmux.NewClient()
		sessions, err := tmuxClient.ListSessions()
		if err != nil {
			return err
		}

		// Get windows for each session
		windows := make(map[string][]tmux.Window)
		for _, s := range sessions {
			wins, err := tmuxClient.ListWindows(s.Name)
			if err == nil {
				windows[s.Name] = wins
			}
		}

		// Build model
		model := tui.InitialModel()
		model.Groups = tui.GroupByWorktree(sessions, windows, tmuxClient)

		// Run TUI
		p := tea.NewProgram(model)
		finalModel, err := p.Run()
		if err != nil {
			return err
		}

		// Handle selection
		if m, ok := finalModel.(tui.Model); ok && !m.Quitting {
			if m.Cursor < len(m.Groups) {
				sessionName := m.Groups[m.Cursor].Name
				fmt.Printf("Attaching to %s...\n", sessionName)
				return tmuxClient.SwitchClient(sessionName)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dashCmd)
}
