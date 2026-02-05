package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rsanzone/clawdbay/internal/config"
	"github.com/rsanzone/clawdbay/templates"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ClawdBay configuration and templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		// Create directories
		if err := cfg.EnsureDirs(); err != nil {
			return fmt.Errorf("failed to create config directories: %w", err)
		}
		fmt.Printf("Created: %s\n", cfg.ConfigDir)

		// Copy templates
		err = fs.WalkDir(templates.FS, "prompts", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}

			content, err := templates.FS.ReadFile(path)
			if err != nil {
				return err
			}

			dstPath := filepath.Join(cfg.PromptsDir, filepath.Base(path))

			// Don't overwrite existing
			if _, err := os.Stat(dstPath); err == nil {
				fmt.Printf("Skipped (exists): %s\n", dstPath)
				return nil
			}

			if err := os.WriteFile(dstPath, content, 0644); err != nil {
				return err
			}
			fmt.Printf("Created: %s\n", dstPath)
			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to copy templates: %w", err)
		}

		fmt.Println("\nClawdBay initialized! Run 'cb start <ticket>' to begin.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
