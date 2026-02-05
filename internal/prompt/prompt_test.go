package prompt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListTemplates(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test templates
	os.WriteFile(filepath.Join(tmpDir, "research.md"), []byte("# Research"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "plan.md"), []byte("# Plan"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "not-md.txt"), []byte("ignored"), 0644)

	templates, err := ListTemplates(tmpDir)
	if err != nil {
		t.Fatalf("ListTemplates() error = %v", err)
	}

	if len(templates) != 2 {
		t.Errorf("got %d templates, want 2", len(templates))
	}
}
