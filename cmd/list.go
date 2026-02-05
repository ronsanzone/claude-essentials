package cmd

import (
	"fmt"

	"github.com/rsanzone/clawdbay/internal/tmux"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active ClawdBay workflows",
	RunE: func(cmd *cobra.Command, args []string) error {
		tmuxClient := tmux.NewClient()
		sessions, err := tmuxClient.ListSessions()
		if err != nil {
			return err
		}

		if len(sessions) == 0 {
			fmt.Println("No active workflows. Start one with: cb start <ticket-id>")
			return nil
		}

		fmt.Println("Active workflows:")
		for _, s := range sessions {
			fmt.Printf("  %s\n", s.Name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
