package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rsanzone/clawdbay/internal/config"
	"github.com/rsanzone/clawdbay/internal/prompt"
	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Manage prompt templates",
}

var promptListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available prompt templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		templates, err := prompt.ListTemplates(cfg.PromptsDir)
		if err != nil {
			return err
		}

		if len(templates) == 0 {
			fmt.Println("No templates found. Create templates in:", cfg.PromptsDir)
			return nil
		}

		fmt.Println("Available templates:")
		for _, t := range templates {
			fmt.Printf("  - %s\n", t)
		}
		return nil
	},
}

var promptAddCmd = &cobra.Command{
	Use:   "add <template-name>",
	Short: "Copy template to .prompts/ and open in editor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := args[0]
		cfg, err := config.New()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		dstDir := filepath.Join(cwd, ".prompts")

		// Copy template
		if err := prompt.CopyTemplate(cfg.PromptsDir, dstDir, templateName); err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, templateName+".md")
		fmt.Printf("Created: %s\n", dstPath)

		// Open in editor
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nvim"
		}

		editorCmd := exec.Command(editor, dstPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr
		return editorCmd.Run()
	},
}

var promptRunCmd = &cobra.Command{
	Use:   "run <prompt-file>",
	Short: "Execute prompt file with Claude",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		promptFile := args[0]
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		promptPath := filepath.Join(cwd, ".prompts", promptFile)

		if _, err := os.Stat(promptPath); os.IsNotExist(err) {
			return fmt.Errorf("prompt not found: %s", promptPath)
		}

		// Execute: claude < prompt.md
		file, err := os.Open(promptPath)
		if err != nil {
			return fmt.Errorf("failed to open prompt: %w", err)
		}
		defer file.Close()

		claudeCmd := exec.Command("claude")
		claudeCmd.Stdin = file
		claudeCmd.Stdout = os.Stdout
		claudeCmd.Stderr = os.Stderr
		return claudeCmd.Run()
	},
}

func init() {
	promptCmd.AddCommand(promptListCmd)
	promptCmd.AddCommand(promptAddCmd)
	promptCmd.AddCommand(promptRunCmd)
	rootCmd.AddCommand(promptCmd)
}
