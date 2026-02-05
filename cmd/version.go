package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the current version of ClawdBay.
var Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print ClawdBay version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ClawdBay v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
