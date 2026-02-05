package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rsanzone/clawdbay/internal/tmux"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <branch-name>",
	Short: "Start a new workflow with a git worktree and tmux session",
	Long: `Creates a git worktree and tmux session for the given branch name.

Example:
  cb start proj-123-auth-feature
  cb start feature/add-login`,
	Args: cobra.ExactArgs(1),
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) error {
	branchName := sanitizeBranchName(args[0])

	// Verify we're in a git repository
	if _, err := exec.Command("git", "rev-parse", "--git-dir").Output(); err != nil {
		return fmt.Errorf("not in a git repository")
	}

	// Get current directory info
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	projectName := filepath.Base(cwd)
	worktreeDir := filepath.Join(filepath.Dir(cwd), projectName+"-"+branchName)

	// Check if worktree directory already exists
	if _, err := os.Stat(worktreeDir); err == nil {
		return fmt.Errorf("worktree directory already exists: %s", worktreeDir)
	}

	// Check if branch already exists
	checkBranch := exec.Command("git", "rev-parse", "--verify", branchName)
	if checkBranch.Run() == nil {
		// Branch exists, create worktree without -b flag
		fmt.Printf("Branch %s exists, creating worktree...\n", branchName)
		gitCmd := exec.Command("git", "worktree", "add", worktreeDir, branchName)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		if err := gitCmd.Run(); err != nil {
			return fmt.Errorf("failed to create worktree: %w", err)
		}
	} else {
		// Create new branch and worktree
		fmt.Printf("Creating worktree: %s\n", worktreeDir)
		gitCmd := exec.Command("git", "worktree", "add", worktreeDir, "-b", branchName)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		if err := gitCmd.Run(); err != nil {
			return fmt.Errorf("failed to create worktree: %w", err)
		}
	}

	// Create tmux session
	sessionName := "cb:" + branchName
	tmuxClient := tmux.NewClient()

	fmt.Printf("Creating tmux session: %s\n", sessionName)
	if err := tmuxClient.CreateSession(sessionName, worktreeDir); err != nil {
		return fmt.Errorf("failed to create tmux session: %w", err)
	}

	// Switch to the session
	if os.Getenv("TMUX") != "" {
		return tmuxClient.SwitchClient(sessionName)
	}
	return tmuxClient.AttachSession(sessionName)
}

// sanitizeBranchName converts a string to a valid git branch name.
func sanitizeBranchName(name string) string {
	// Replace spaces and special chars with dashes
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")

	// Remove characters not allowed in branch names
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '/' {
			result.WriteRune(r)
		}
	}

	// Clean up multiple dashes
	cleaned := result.String()
	for strings.Contains(cleaned, "--") {
		cleaned = strings.ReplaceAll(cleaned, "--", "-")
	}

	return strings.Trim(cleaned, "-")
}
