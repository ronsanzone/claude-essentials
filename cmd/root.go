package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cb",
	Short: "ClawdBay - A harbor for your Claude sessions",
	Long: `ClawdBay manages multi-session Claude Code workflows.

Start workflows with git worktrees, manage multiple Claude sessions
per worktree, and track session status from an interactive dashboard.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default to dashboard when no subcommand
		fmt.Println("ClawdBay dashboard would open here")
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
