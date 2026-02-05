package prompt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ListTemplates returns a list of .md template names in the given directory.
func ListTemplates(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var templates []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			templates = append(templates, strings.TrimSuffix(e.Name(), ".md"))
		}
	}
	return templates, nil
}

// CopyTemplate copies a template from srcDir to dstDir.
func CopyTemplate(srcDir, dstDir, name string) error {
	srcPath := filepath.Join(srcDir, name+".md")
	dstPath := filepath.Join(dstDir, name+".md")

	// Ensure destination directory exists
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("template %q not found: %w", name, err)
	}
	defer func() {
		_ = src.Close()
	}()

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		_ = dst.Close()
	}()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}
	return nil
}
